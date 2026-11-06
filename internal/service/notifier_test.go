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
			func() { notifyDown(t.Context(), n, c.m, incident) },
			func() { notifyRecovery(t.Context(), n, c.m, incident) },
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

	dashboardURL := func() *string {
		t.Helper()
		var p struct {
			DashboardURL *string `json:"dashboard_url"`
		}
		if err := json.Unmarshal([]byte(<-bodies), &p); err != nil {
			t.Fatal(err)
		}
		return p.DashboardURL
	}

	for publicURL, want := range map[string]*string{
		"":                           nil,
		"https://pingou.exemplo.com": str("https://pingou.exemplo.com/monitors/abc"),
	} {
		n := NewWebhookNotifier(func() string { return srv.URL }, publicURL)
		notifyDown(t.Context(), n, m, incident)
		if got := dashboardURL(); !eqPtr(got, want) {
			t.Fatalf("down, publicURL %q: dashboard_url = %v", publicURL, got)
		}
		notifyRecovery(t.Context(), n, m, incident)
		if got := dashboardURL(); !eqPtr(got, want) {
			t.Fatalf("up, publicURL %q: dashboard_url = %v", publicURL, got)
		}
	}

	// evento "test": o monitor é fictício, então nunca leva link
	SendTestWebhook(t.Context(), srv.URL)
	if got := dashboardURL(); got != nil {
		t.Fatalf("test: dashboard_url = %v, want null", *got)
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

// Status fora de 2xx é falha (antes contava como "sent"); 2xx é sucesso.
func TestWebhookSendStatusCodes(t *testing.T) {
	for _, c := range []struct {
		status  int
		wantErr bool
	}{{200, false}, {204, false}, {301, true}, {400, true}, {404, true}, {500, true}, {503, true}} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "corpo do receptor", c.status)
			if c.status < 300 {
				w.WriteHeader(c.status)
			}
		}))
		n := NewWebhookNotifier(func() string { return srv.URL }, "")
		err := notifyDown(t.Context(), n, &domain.Monitor{ID: "1", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}, &domain.Incident{StartedAt: time.Now()})
		srv.Close()
		if (err != nil) != c.wantErr {
			t.Errorf("status %d: err = %v, wantErr %v", c.status, err, c.wantErr)
		}
		if err != nil && strings.Contains(err.Error(), srv.URL) {
			t.Errorf("status %d: erro com a URL: %v", c.status, err)
		}
	}
}

// Erro de rede devolve erro sem a URL (que pode ter token na query).
func TestWebhookSendNetworkErrorHidesURL(t *testing.T) {
	n := NewWebhookNotifier(func() string { return "http://127.0.0.1:1/hook?token=segredo" }, "")
	err := notifyDown(t.Context(), n, &domain.Monitor{ID: "1", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}, &domain.Incident{StartedAt: time.Now()})
	if err == nil {
		t.Fatal("erro esperado")
	}
	if strings.Contains(err.Error(), "segredo") || strings.Contains(err.Error(), "127.0.0.1:1/hook") {
		t.Fatalf("erro vazou a URL: %v", err)
	}
}

// Sem URL configurada o canal está desligado: nil, sem chamada (sem retry eterno).
func TestWebhookSendEmptyURLIsNoop(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer srv.Close()
	n := NewWebhookNotifier(func() string { return "" }, "")
	if err := notifyDown(t.Context(), n, &domain.Monitor{ID: "1", Name: "API"}, &domain.Incident{StartedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if calls != 0 {
		t.Fatalf("chamadas = %d, want 0", calls)
	}
}

// O JSON ganha só o campo id (igual ao event_id); o timestamp é a hora do
// evento, não a do envio.
func TestWebhookPayloadHasEventIDAndEventTime(t *testing.T) {
	bodies := make(chan string, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	n := NewWebhookNotifier(func() string { return srv.URL }, "")
	m := &domain.Monitor{ID: "1", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}
	occurred := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

	ev := testEvent(domain.EventCertExpiring, domain.NewCertSnapshot(m, domain.CertAlert{ExpiresAt: occurred.Add(72 * time.Hour), DaysLeft: 3, ThresholdDays: 7}, occurred))
	for range 2 { // duas tentativas: mesmo corpo
		if err := n.Send(t.Context(), ev); err != nil {
			t.Fatal(err)
		}
	}
	first, second := <-bodies, <-bodies
	if first != second {
		t.Fatalf("corpo mudou entre tentativas:\n%s\n%s", first, second)
	}
	var p map[string]any
	if err := json.Unmarshal([]byte(first), &p); err != nil {
		t.Fatal(err)
	}
	if p["id"] != ev.EventID || p["timestamp"] != "2026-03-04T05:06:07Z" || p["event"] != "cert_expiring" {
		t.Fatalf("payload: %s", first)
	}
	// os demais campos seguem como antes
	for _, k := range []string{"monitor", "last_error", "downtime_duration_seconds", "dashboard_url", "cert_expires_at", "days_left", "threshold_days"} {
		if _, ok := p[k]; !ok {
			t.Errorf("campo %q sumiu: %s", k, first)
		}
	}
}

// O botão de teste do webhook leva um id novo a cada envio.
func TestSendTestWebhookHasID(t *testing.T) {
	bodies := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	if res := SendTestWebhook(t.Context(), srv.URL); !res.OK {
		t.Fatalf("resultado: %+v", res)
	}
	var p struct{ ID, Event string }
	if err := json.Unmarshal([]byte(<-bodies), &p); err != nil || p.ID == "" || p.Event != "test" {
		t.Fatalf("payload do teste: %+v err=%v", p, err)
	}
}

// O campo "event" do webhook usa monitor_down/monitor_up, não os tipos internos do outbox.
func TestWebhookDownAndRecoveryEventNames(t *testing.T) {
	bodies := make(chan string, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	n := NewWebhookNotifier(func() string { return srv.URL }, "")
	m := &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP}
	i := &domain.Incident{ID: "i1", StartedAt: time.Now().Add(-90 * time.Second)}

	for _, c := range []struct {
		send func() error
		want string
	}{
		{func() error { return notifyDown(t.Context(), n, m, i) }, "monitor_down"},
		{func() error { return notifyRecovery(t.Context(), n, m, i) }, "monitor_up"},
	} {
		if err := c.send(); err != nil {
			t.Fatal(err)
		}
		var p struct{ Event string }
		if err := json.Unmarshal([]byte(<-bodies), &p); err != nil || p.Event != c.want {
			t.Fatalf("event = %q (err %v), quer %q", p.Event, err, c.want)
		}
	}
}
