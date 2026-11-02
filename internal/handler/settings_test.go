package handler

import (
	"bufio"
	"bytes"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
)

func TestInfo(t *testing.T) {
	s, h, _ := newTestServer(t)

	if err := os.WriteFile(s.cfg.DatabaseURL+"-wal", make([]byte, 1000), 0o644); err != nil {
		t.Fatal(err)
	}
	want, err := databaseSize(s.cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}

	rec := do(t, h, "GET", "/api/info", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("info: %d %s", rec.Code, rec.Body)
	}
	got := decode[infoResponse](t, rec)
	if got.Version != "v1.2.3" || got.Commit != "abc1234" || got.BuildDate != "2026-10-29" || got.DatabaseSizeBytes != want {
		t.Fatalf("info: %+v, want size %d", got, want)
	}

	s.cfg.DatabaseURL = filepath.Join(t.TempDir(), "nao-existe.db")
	if rec := do(t, h, "GET", "/api/info", ""); rec.Code != http.StatusInternalServerError {
		t.Fatalf("sem banco: %d, want 500", rec.Code)
	}
}

func TestDatabaseSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "p.db")
	if err := os.WriteFile(path, make([]byte, 100), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path+"-wal", make([]byte, 20), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		dsn  string
		want int64
	}{
		{path, 120}, // sem -shm não é erro
		{"file:" + path, 120},
		{path + "?_journal_mode=WAL", 120},
		{"file:" + path + "?cache=shared", 120},
	}
	for _, c := range cases {
		got, err := databaseSize(c.dsn)
		if err != nil || got != c.want {
			t.Fatalf("%s: %d %v, want %d", c.dsn, got, err, c.want)
		}
	}
}

func TestExportFilename(t *testing.T) {
	_, h, _ := newTestServer(t)
	rec := do(t, h, "GET", "/api/export", "")
	want := `attachment; filename="pingou-` + time.Now().UTC().Format("2006-01-02") + `.db"`
	if rec.Code != http.StatusOK || rec.Header().Get("Content-Disposition") != want {
		t.Fatalf("export: %d %q, want %q", rec.Code, rec.Header().Get("Content-Disposition"), want)
	}
}

func TestWebhookTestRoute(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/404" {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer dest.Close()

	cases := []struct {
		name      string
		saved     string // webhook_url salva antes do teste
		body      string
		code      int
		ok        bool
		hasStatus int
	}{
		{"url do corpo 2xx", "", `{"url":"` + dest.URL + `"}`, 200, true, 200},
		{"url do corpo 404", "", `{"url":"` + dest.URL + `/404"}`, 200, false, 404},
		{"sem corpo usa a salva", dest.URL, "", 200, true, 200},
		{"corpo com url vazia usa a salva", dest.URL, `{"url":""}`, 200, true, 200},
		{"sem url", "", "", 400, false, 0},
		{"esquema ftp", "", `{"url":"ftp://x.test"}`, 400, false, 0},
		{"json inválido", "", `{`, 400, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, h, _ := newTestServer(t)
			if c.saved != "" {
				if rec := do(t, h, "PATCH", "/api/settings", `{"webhook_url":"`+c.saved+`"}`); rec.Code != 200 {
					t.Fatalf("salvar: %d %s", rec.Code, rec.Body)
				}
			}
			rec := do(t, h, "POST", "/api/settings/webhook/test", c.body)
			if rec.Code != c.code {
				t.Fatalf("%d %s, want %d", rec.Code, rec.Body, c.code)
			}
			if c.code != 200 {
				return
			}
			res := decode[struct {
				OK         bool    `json:"ok"`
				StatusCode *int    `json:"status_code"`
				Error      *string `json:"error"`
			}](t, rec)
			if res.OK != c.ok || res.StatusCode == nil || *res.StatusCode != c.hasStatus {
				t.Fatalf("resultado: %s", rec.Body)
			}
		})
	}
}

func TestWebhookTestRateLimit(t *testing.T) {
	_, h, _ := newTestServer(t)
	for i := range webhookTestRateLimit {
		if rec := do(t, h, "POST", "/api/settings/webhook/test", ""); rec.Code != http.StatusBadRequest {
			t.Fatalf("req %d: %d, want 400", i, rec.Code)
		}
	}
	rec := do(t, h, "POST", "/api/settings/webhook/test", "")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("excesso: %d Retry-After=%q, want 429", rec.Code, rec.Header().Get("Retry-After"))
	}
}

func TestSettingsRejectsNonHTTPWebhook(t *testing.T) {
	_, h, _ := newTestServer(t)
	for _, body := range []string{`{"webhook_url":"ftp://x.test"}`, `{"webhook_url":"mailto:a@b.c"}`, `{"webhook_url":"http://"}`} {
		if rec := do(t, h, "PATCH", "/api/settings", body); rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: %d, want 400", body, rec.Code)
		}
	}
	// vazio desliga o webhook
	if rec := do(t, h, "PATCH", "/api/settings", `{"webhook_url":""}`); rec.Code != http.StatusOK {
		t.Fatalf("vazio: %d %s", rec.Code, rec.Body)
	}
}

func TestWebhookTestURLNotLogged(t *testing.T) {
	_, h, _ := newTestServer(t)
	dest := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer dest.Close()

	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	do(t, h, "POST", "/api/settings/webhook/test", `{"url":"`+dest.URL+`/hook?token=segredo-do-hook"}`)
	// destino que não responde: o erro também não pode levar a URL
	rec := do(t, h, "POST", "/api/settings/webhook/test", `{"url":"http://127.0.0.1:1/hook?token=segredo-do-hook"}`)

	if strings.Contains(buf.String(), "segredo-do-hook") || strings.Contains(rec.Body.String(), "segredo-do-hook") {
		t.Fatalf("url completa vazou: log=%s resp=%s", buf.String(), rec.Body)
	}
	if !strings.Contains(buf.String(), "webhook test") {
		t.Fatalf("log sem o teste: %s", buf.String())
	}
}

// fakeSMTPServer responde 250 a tudo (354 no DATA) e conta as mensagens.
func fakeSMTPServer(t *testing.T) (config.SMTPConfig, *atomic.Int64) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	var sent atomic.Int64
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				r := bufio.NewReader(conn)
				conn.Write([]byte("220 fake\r\n"))
				for {
					line, err := r.ReadString('\n')
					if err != nil {
						return
					}
					switch strings.ToUpper(strings.TrimSpace(line)) {
					case "DATA":
						conn.Write([]byte("354 go\r\n"))
						for l, err := r.ReadString('\n'); l != ".\r\n"; l, err = r.ReadString('\n') {
							if err != nil {
								return
							}
						}
						sent.Add(1)
						conn.Write([]byte("250 queued\r\n"))
					case "QUIT":
						conn.Write([]byte("221 bye\r\n"))
						return
					default:
						conn.Write([]byte("250 ok\r\n"))
					}
				}
			}()
		}
	}()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	p, _ := strconv.Atoi(port)
	return config.SMTPConfig{Host: "127.0.0.1", Port: p, From: "Pingou <pingou@localhost>", TLS: config.SMTPTLSNone}, &sent
}

type settingsBody struct {
	WebhookURL      string `json:"webhook_url"`
	RetentionDays   int    `json:"retention_days"`
	EmailRecipients string `json:"email_recipients"`
	EmailAvailable  bool   `json:"email_available"`
	NtfyTopic       string `json:"ntfy_topic"`
	NtfyAvailable   bool   `json:"ntfy_available"`
	Timezone        string `json:"timezone"`
	TimezoneSet     bool   `json:"timezone_set"`
}

func TestSettingsEmailWithoutSMTP(t *testing.T) {
	_, h, _ := newTestServer(t)

	got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", ""))
	if got.EmailAvailable || got.EmailRecipients != "" {
		t.Fatalf("GET: %+v", got)
	}
	if rec := do(t, h, "PATCH", "/api/settings", `{"email_recipients":"ops@exemplo.com"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("PATCH sem SMTP: %d, want 400", rec.Code)
	}
	if rec := do(t, h, "PATCH", "/api/settings", `{"email_recipients":""}`); rec.Code != http.StatusOK {
		t.Fatalf("PATCH vazio: %d %s", rec.Code, rec.Body)
	}
	// salvar outra coisa não depende do e-mail
	if rec := do(t, h, "PATCH", "/api/settings", `{"retention_days":14}`); rec.Code != http.StatusOK {
		t.Fatalf("PATCH retenção: %d %s", rec.Code, rec.Body)
	}
	if rec := do(t, h, "POST", "/api/settings/email/test", `{"recipients":"ops@exemplo.com"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("teste sem SMTP: %d, want 400", rec.Code)
	}
}

func TestSettingsEmailWithSMTP(t *testing.T) {
	smtp, sent := fakeSMTPServer(t)
	smtp.Password = "senha-do-smtp"
	_, h, _ := newTestServerSMTP(t, smtp)

	if got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", "")); !got.EmailAvailable {
		t.Fatalf("GET: %+v", got)
	}

	rec := do(t, h, "PATCH", "/api/settings", `{"email_recipients":" ops@exemplo.com, Bryan <bryan@exemplo.com> "}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH: %d %s", rec.Code, rec.Body)
	}
	if got := decode[settingsBody](t, rec); got.EmailRecipients != "ops@exemplo.com,bryan@exemplo.com" {
		t.Fatalf("normalizado: %q", got.EmailRecipients)
	}
	for _, body := range []string{`{"email_recipients":"não é email"}`, `{"email_recipients":"a@x.com,b@x.com,c@x.com,d@x.com,e@x.com,f@x.com,g@x.com,h@x.com,i@x.com,j@x.com,k@x.com"}`} {
		if rec := do(t, h, "PATCH", "/api/settings", body); rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: %d, want 400", body, rec.Code)
		}
	}
	if strings.Contains(do(t, h, "GET", "/api/settings", "").Body.String(), "senha-do-smtp") {
		t.Fatal("GET vazou a senha")
	}

	// sem corpo: usa os destinatários salvos
	rec = do(t, h, "POST", "/api/settings/email/test", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("teste: %d %s", rec.Code, rec.Body)
	}
	res := decode[struct {
		OK    bool    `json:"ok"`
		Error *string `json:"error"`
	}](t, rec)
	if !res.OK || res.Error != nil || sent.Load() != 1 {
		t.Fatalf("teste: %s (enviados %d)", rec.Body, sent.Load())
	}

	if rec := do(t, h, "POST", "/api/settings/email/test", `{"recipients":"x"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("destinatário inválido: %d, want 400", rec.Code)
	}
}

func TestEmailTestWithoutRecipients(t *testing.T) {
	smtp, _ := fakeSMTPServer(t)
	_, h, _ := newTestServerSMTP(t, smtp)
	if rec := do(t, h, "POST", "/api/settings/email/test", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("sem destinatários: %d, want 400", rec.Code)
	}
}

func TestEmailTestRateLimit(t *testing.T) {
	_, h, _ := newTestServer(t)
	for i := range emailTestRateLimit {
		if rec := do(t, h, "POST", "/api/settings/email/test", ""); rec.Code != http.StatusBadRequest {
			t.Fatalf("req %d: %d, want 400", i, rec.Code)
		}
	}
	rec := do(t, h, "POST", "/api/settings/email/test", "")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("excesso: %d, want 429", rec.Code)
	}
}

func TestSettingsNtfyWithoutURL(t *testing.T) {
	_, h, _ := newTestServer(t)

	got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", ""))
	if got.NtfyAvailable || got.NtfyTopic != "" {
		t.Fatalf("GET: %+v", got)
	}
	if rec := do(t, h, "PATCH", "/api/settings", `{"ntfy_topic":"alertas"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("PATCH sem ntfy: %d, want 400", rec.Code)
	}
	if rec := do(t, h, "PATCH", "/api/settings", `{"ntfy_topic":""}`); rec.Code != http.StatusOK {
		t.Fatalf("PATCH vazio: %d %s", rec.Code, rec.Body)
	}
	if rec := do(t, h, "POST", "/api/settings/ntfy/test", `{"topic":"alertas"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("teste sem ntfy: %d, want 400", rec.Code)
	}
}

func TestSettingsNtfyWithURL(t *testing.T) {
	var published atomic.Int64
	var lastPath, lastAuth atomic.Value
	ntfy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		published.Add(1)
		lastPath.Store(r.URL.Path)
		lastAuth.Store(r.Header.Get("Authorization"))
	}))
	defer ntfy.Close()
	_, h, _ := newTestServerEnv(t, config.SMTPConfig{}, config.NtfyConfig{URL: ntfy.URL, Token: "tk_segredo"})

	if got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", "")); !got.NtfyAvailable {
		t.Fatalf("GET: %+v", got)
	}

	// sem tópico salvo
	if rec := do(t, h, "POST", "/api/settings/ntfy/test", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("teste sem tópico: %d, want 400", rec.Code)
	}

	rec := do(t, h, "PATCH", "/api/settings", `{"ntfy_topic":" alertas_prod "}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH: %d %s", rec.Code, rec.Body)
	}
	if got := decode[settingsBody](t, rec); got.NtfyTopic != "alertas_prod" {
		t.Fatalf("tópico: %q", got.NtfyTopic)
	}
	for _, topic := range []string{"com espaço", "a/b", strings.Repeat("x", 65)} {
		if rec := do(t, h, "PATCH", "/api/settings", `{"ntfy_topic":"`+topic+`"}`); rec.Code != http.StatusBadRequest {
			t.Fatalf("%q: %d, want 400", topic, rec.Code)
		}
	}
	body := do(t, h, "GET", "/api/settings", "").Body.String()
	if strings.Contains(body, "tk_segredo") || strings.Contains(body, ntfy.URL) {
		t.Fatalf("GET vazou token ou URL: %s", body)
	}

	// sem corpo: usa o tópico salvo
	rec = do(t, h, "POST", "/api/settings/ntfy/test", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("teste: %d %s", rec.Code, rec.Body)
	}
	res := decode[struct {
		OK    bool    `json:"ok"`
		Error *string `json:"error"`
	}](t, rec)
	if !res.OK || res.Error != nil || published.Load() != 1 || lastPath.Load() != "/alertas_prod" || lastAuth.Load() != "Bearer tk_segredo" {
		t.Fatalf("teste: %s (publicações %d, path %v)", rec.Body, published.Load(), lastPath.Load())
	}

	// tópico do corpo tem prioridade
	if rec := do(t, h, "POST", "/api/settings/ntfy/test", `{"topic":"outro"}`); rec.Code != http.StatusOK || lastPath.Load() != "/outro" {
		t.Fatalf("tópico do corpo: %d %v", rec.Code, lastPath.Load())
	}
	if rec := do(t, h, "POST", "/api/settings/ntfy/test", `{"topic":"a b"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("tópico inválido: %d, want 400", rec.Code)
	}
}

func TestNtfyTestRateLimit(t *testing.T) {
	_, h, _ := newTestServer(t)
	for i := range ntfyTestRateLimit {
		if rec := do(t, h, "POST", "/api/settings/ntfy/test", ""); rec.Code != http.StatusBadRequest {
			t.Fatalf("req %d: %d, want 400", i, rec.Code)
		}
	}
	rec := do(t, h, "POST", "/api/settings/ntfy/test", "")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("excesso: %d, want 429", rec.Code)
	}
}

func TestSettingsTimezone(t *testing.T) {
	_, h, _ := newTestServer(t)

	if got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", "")); got.Timezone != "UTC" || got.TimezoneSet {
		t.Fatalf("padrão = %q set=%v, want UTC não configurado", got.Timezone, got.TimezoneSet)
	}
	rec := do(t, h, "PATCH", "/api/settings", `{"timezone":"America/Sao_Paulo"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("válido: %d %s", rec.Code, rec.Body)
	}
	if got := decode[settingsBody](t, rec); got.Timezone != "America/Sao_Paulo" || !got.TimezoneSet {
		t.Fatalf("PATCH = %q set=%v", got.Timezone, got.TimezoneSet)
	}
	if rec := do(t, h, "PATCH", "/api/settings", `{"timezone":"Marte/Olimpo"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("inválido: %d, want 400", rec.Code)
	}
	// PATCH sem o campo não mexe no fuso
	do(t, h, "PATCH", "/api/settings", `{"retention_days":30}`)
	if got := decode[settingsBody](t, do(t, h, "GET", "/api/settings", "")); got.Timezone != "America/Sao_Paulo" {
		t.Fatalf("depois de outro PATCH = %q", got.Timezone)
	}
	// vazio volta para UTC
	if got := decode[settingsBody](t, do(t, h, "PATCH", "/api/settings", `{"timezone":""}`)); got.Timezone != "UTC" {
		t.Fatalf("vazio = %q, want UTC", got.Timezone)
	}
}
