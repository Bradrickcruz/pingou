package service

import (
	"bytes"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

// ntfyRequest é o que o servidor fake recebeu.
type ntfyRequest struct {
	method, path, body string
	header             http.Header
}

// fakeNtfy grava cada publicação e responde com status.
type fakeNtfy struct {
	srv    *httptest.Server
	mu     sync.Mutex
	reqs   []ntfyRequest
	got    chan struct{}
	status int
}

func newFakeNtfy(t *testing.T, status int) *fakeNtfy {
	t.Helper()
	f := &fakeNtfy{got: make(chan struct{}, 10), status: status}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.reqs = append(f.reqs, ntfyRequest{r.Method, r.URL.EscapedPath(), string(b), r.Header.Clone()})
		f.mu.Unlock()
		w.WriteHeader(f.status)
		f.got <- struct{}{}
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeNtfy) last(t *testing.T) ntfyRequest {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.reqs) == 0 {
		t.Fatal("nenhuma publicação")
	}
	return f.reqs[len(f.reqs)-1]
}

func decodedTitle(t *testing.T, h http.Header) string {
	t.Helper()
	s, err := new(mime.WordDecoder).DecodeHeader(h.Get("Title"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestNtfyNotifierDown(t *testing.T) {
	f := newFakeNtfy(t, http.StatusOK)
	n := NewNtfyNotifier(config.NtfyConfig{URL: f.srv.URL, Token: "tk_segredo"}, func() string { return "alertas" }, "https://pingou.exemplo.com")
	m := &domain.Monitor{ID: "m1", Name: "API\r\nX-Evil: 1", Type: domain.TypeHTTP, URL: "https://api.test"}
	errMsg := "connection refused"
	notifyDown(t.Context(), n, m, &domain.Incident{StartedAt: time.Now(), LastError: &errMsg})

	r := f.last(t)
	if r.method != http.MethodPost || r.path != "/alertas" {
		t.Fatalf("%s %s", r.method, r.path)
	}
	if got := decodedTitle(t, r.header); got != "API X-Evil: 1 está fora do ar" {
		t.Fatalf("Title = %q", got)
	}
	if r.header.Get("X-Evil") != "" {
		t.Fatal("injeção de header pelo nome do monitor")
	}
	want := map[string]string{
		"Priority":      "high",
		"Tags":          "rotating_light",
		"Click":         "https://pingou.exemplo.com/monitors/m1",
		"Authorization": "Bearer tk_segredo",
	}
	for k, v := range want {
		if r.header.Get(k) != v {
			t.Fatalf("%s = %q, want %q", k, r.header.Get(k), v)
		}
	}
	for _, s := range []string{"Monitor: API X-Evil: 1", "Alvo: https://api.test", "Erro: connection refused", "Detalhes: https://pingou.exemplo.com/monitors/m1"} {
		if !strings.Contains(r.body, s) {
			t.Fatalf("corpo sem %q:\n%s", s, r.body)
		}
	}
}

func TestNtfyNotifierRecoveryWithoutTokenOrLink(t *testing.T) {
	f := newFakeNtfy(t, http.StatusOK)
	n := NewNtfyNotifier(config.NtfyConfig{URL: f.srv.URL}, func() string { return "alertas" }, "")
	started := time.Now().Add(-5 * time.Minute)
	ended := time.Now()
	notifyRecovery(t.Context(), n, &domain.Monitor{Name: "API", Type: domain.TypeHTTP, URL: "https://api.test"},
		&domain.Incident{StartedAt: started, EndedAt: &ended})

	r := f.last(t)
	if got := decodedTitle(t, r.header); got != "API voltou (fora por 5 min)" {
		t.Fatalf("Title = %q", got)
	}
	if r.header.Get("Priority") != "default" || r.header.Get("Tags") != "white_check_mark" {
		t.Fatalf("headers: %v", r.header)
	}
	if _, ok := r.header["Authorization"]; ok {
		t.Fatal("Authorization sem token")
	}
	if _, ok := r.header["Click"]; ok {
		t.Fatal("Click sem PINGOU_PUBLIC_URL")
	}
	if !strings.Contains(r.body, "Fora do ar por: 5 min") {
		t.Fatalf("corpo:\n%s", r.body)
	}
}

func TestNtfyNotifierHeartbeatNeverSendsPingToken(t *testing.T) {
	f := newFakeNtfy(t, http.StatusOK)
	n := NewNtfyNotifier(config.NtfyConfig{URL: f.srv.URL}, func() string { return "alertas" }, "https://p.test")
	token := "ping-token-secreto"
	m := &domain.Monitor{ID: "h1", Name: "Backup", Type: domain.TypeHeartbeat, PingToken: &token}
	notifyDown(t.Context(), n, m, &domain.Incident{StartedAt: time.Now()})

	r := f.last(t)
	if strings.Contains(r.body, "Alvo") {
		t.Fatalf("heartbeat com alvo:\n%s", r.body)
	}
	dump := r.body + r.path
	for k, v := range r.header {
		dump += k + strings.Join(v, ",")
	}
	if strings.Contains(dump, token) {
		t.Fatal("ping_token vazou na mensagem")
	}
}

func TestNtfyNotifierSkipsWithoutTopic(t *testing.T) {
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer srv.Close()
	n := NewNtfyNotifier(config.NtfyConfig{URL: srv.URL}, func() string { return "" }, "")
	notifyDown(t.Context(), n, &domain.Monitor{Name: "x"}, &domain.Incident{StartedAt: time.Now()})
	if calls.Load() != 0 {
		t.Fatal("publicou sem tópico")
	}
}

func TestSendTestNtfy(t *testing.T) {
	f := newFakeNtfy(t, http.StatusOK)
	res := SendTestNtfy(t.Context(), config.NtfyConfig{URL: f.srv.URL, Token: "tk"}, "meu_topico-1")
	if !res.OK || res.Error != nil || res.StatusCode == nil || *res.StatusCode != 200 {
		t.Fatalf("resultado: %+v", res)
	}
	if r := f.last(t); r.path != "/meu_topico-1" || r.header.Get("Authorization") != "Bearer tk" {
		t.Fatalf("req: %+v", r)
	}
}

func TestSendTestNtfyFailures(t *testing.T) {
	forbidden := newFakeNtfy(t, http.StatusForbidden)
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://outro.test/", http.StatusMovedPermanently)
	}))
	defer redirect.Close()

	cases := []struct {
		name, url, wantErr string
		status             int
	}{
		{"403", forbidden.srv.URL, "403 (token inválido", 403},
		{"redirect", redirect.URL, "unexpected status code: 301", 301},
		{"fora do ar", "http://127.0.0.1:1", "connection refused", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			prev := slog.Default()
			slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
			defer slog.SetDefault(prev)

			res := SendTestNtfy(t.Context(), config.NtfyConfig{URL: c.url, Token: "tk_segredo"}, "topico-secreto")
			if res.OK || res.Error == nil || !strings.Contains(*res.Error, c.wantErr) {
				t.Fatalf("resultado: %+v (%v)", res, res.Error)
			}
			if (res.StatusCode == nil) != (c.status == 0) || (res.StatusCode != nil && *res.StatusCode != c.status) {
				t.Fatalf("status: %v", res.StatusCode)
			}
			for _, secret := range []string{"tk_segredo", "topico-secreto", c.url} {
				if strings.Contains(*res.Error, secret) || strings.Contains(buf.String(), secret) {
					t.Fatalf("%q vazou: erro=%s log=%s", secret, *res.Error, buf.String())
				}
			}
		})
	}
}

func TestSendTestNtfyTimeout(t *testing.T) {
	prev := ntfyTimeout
	ntfyTimeout = 50 * time.Millisecond
	defer func() { ntfyTimeout = prev }()
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	defer srv.Close()
	defer close(release)

	res := SendTestNtfy(t.Context(), config.NtfyConfig{URL: srv.URL}, "t")
	if res.Error == nil || !strings.Contains(*res.Error, "timeout após") {
		t.Fatalf("resultado: %+v", res)
	}
}

func TestValidNtfyTopic(t *testing.T) {
	for _, s := range []string{"a", "alertas_prod-1", strings.Repeat("x", 64)} {
		if !validNtfyTopic(s) {
			t.Fatalf("%q deveria valer", s)
		}
	}
	for _, s := range []string{"", "com espaço", "a/b", "acentuação", strings.Repeat("x", 65), "a\r\nb"} {
		if validNtfyTopic(s) {
			t.Fatalf("%q não deveria valer", s)
		}
	}
}
