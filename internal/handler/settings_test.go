package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
