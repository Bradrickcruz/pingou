package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

func TestWebhookPayloadByType(t *testing.T) {
	bodies := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	n := NewWebhookNotifier(func() string { return srv.URL }, "")

	token := "segredo-do-token"
	incident := &domain.Incident{ID: "i1", StartedAt: time.Now()}
	str := func(s string) *string { return &s }

	cases := []struct {
		m           *domain.Monitor
		target, url *string
	}{
		{&domain.Monitor{ID: "1", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test/health"}, str("https://x.test/health"), str("https://x.test/health")},
		{&domain.Monitor{ID: "2", Name: "DB", Type: domain.TypeTCP, Config: domain.MonitorConfig{Host: "db", Port: 5432}}, str("db:5432"), nil},
		{&domain.Monitor{ID: "3", Name: "Backup", Type: domain.TypeHeartbeat, PingToken: &token}, nil, nil},
	}
	for _, c := range cases {
		for _, send := range []func(){
			func() { n.NotifyDown(t.Context(), c.m, incident) },
			func() { n.NotifyRecovery(t.Context(), c.m, incident) },
		} {
			send()
			body := <-bodies
			if strings.Contains(body, token) || strings.Contains(body, "ping_token") {
				t.Fatalf("%s: webhook vazou o token: %s", c.m.Type, body)
			}
			var p struct {
				Monitor struct {
					Type   string  `json:"type"`
					Target *string `json:"target"`
					URL    *string `json:"url"`
				} `json:"monitor"`
			}
			if err := json.Unmarshal([]byte(body), &p); err != nil {
				t.Fatal(err)
			}
			if p.Monitor.Type != string(c.m.Type) || !eqPtr(p.Monitor.Target, c.target) || !eqPtr(p.Monitor.URL, c.url) {
				t.Fatalf("%s: payload %s", c.m.Type, body)
			}
		}
	}
}

func TestWebhookDashboardURL(t *testing.T) {
	bodies := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	m := &domain.Monitor{ID: "abc", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}
	incident := &domain.Incident{ID: "i1", StartedAt: time.Now()}
	str := func(s string) *string { return &s }

	for publicURL, want := range map[string]*string{
		"":                           nil,
		"https://pingou.exemplo.com": str("https://pingou.exemplo.com/monitors/abc"),
	} {
		NewWebhookNotifier(func() string { return srv.URL }, publicURL).NotifyDown(t.Context(), m, incident)
		var p struct {
			DashboardURL *string `json:"dashboard_url"`
		}
		if err := json.Unmarshal([]byte(<-bodies), &p); err != nil {
			t.Fatal(err)
		}
		if !eqPtr(p.DashboardURL, want) {
			t.Fatalf("publicURL %q: dashboard_url = %v", publicURL, p.DashboardURL)
		}
	}
}

func eqPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func TestSendTestWebhook(t *testing.T) {
	prev := webhookTestTimeout
	webhookTestTimeout = 100 * time.Millisecond
	defer func() { webhookTestTimeout = prev }()

	bodies := make(chan string, 10)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
		switch r.URL.Path {
		case "/404":
			w.WriteHeader(http.StatusNotFound)
		case "/302":
			http.Redirect(w, r, "/redirected", http.StatusFound)
		case "/redirected":
			t.Error("o teste seguiu o redirect")
		case "/slow":
			time.Sleep(300 * time.Millisecond)
		case "/big":
			w.Write([]byte(strings.Repeat("x", 1<<20)))
		}
	}))
	defer srv.Close()

	cases := []struct {
		path   string
		ok     bool
		status int // 0 = status_code null
		errMsg string
	}{
		{"/", true, 200, ""},
		{"/big", true, 200, ""},
		{"/404", false, 404, "unexpected status code: 404"},
		{"/302", false, 302, "unexpected status code: 302"},
		{"/slow", false, 0, "timeout após 100ms"},
	}
	for _, c := range cases {
		r := SendTestWebhook(t.Context(), srv.URL+c.path)
		if r.OK != c.ok {
			t.Fatalf("%s: ok=%v, want %v", c.path, r.OK, c.ok)
		}
		if (c.status == 0) != (r.StatusCode == nil) || (r.StatusCode != nil && *r.StatusCode != c.status) {
			t.Fatalf("%s: status_code %v, want %d", c.path, r.StatusCode, c.status)
		}
		if (c.errMsg == "") != (r.Error == nil) || (r.Error != nil && *r.Error != c.errMsg) {
			t.Fatalf("%s: error %v, want %q", c.path, r.Error, c.errMsg)
		}

		var p struct {
			Event   string `json:"event"`
			Monitor struct {
				Name string `json:"name"`
			} `json:"monitor"`
		}
		if err := json.Unmarshal([]byte(<-bodies), &p); err != nil {
			t.Fatal(err)
		}
		if p.Event != "test" || p.Monitor.Name != "Teste do Pingou" {
			t.Fatalf("%s: payload %+v", c.path, p)
		}
	}

	// recusa de conexão: erro sem o prefixo com a URL
	r := SendTestWebhook(t.Context(), "http://127.0.0.1:1/hook")
	if r.OK || r.StatusCode != nil || r.Error == nil || strings.Contains(*r.Error, "127.0.0.1:1/hook") {
		t.Fatalf("conexão recusada: %+v %v", r, *r.Error)
	}
}
