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
	n := NewWebhookNotifier(func() string { return srv.URL })

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

func eqPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
