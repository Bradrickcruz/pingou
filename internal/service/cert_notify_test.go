package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

var (
	certMonitor = &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://api.example.com"}
	certExpires = time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC)
)

func TestWebhookCertExpiring(t *testing.T) {
	bodies := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	n := NewWebhookNotifier(func() string { return srv.URL }, "https://pingou.test")

	notifyCertExpiring(t.Context(), n, certMonitor, domain.CertAlert{ExpiresAt: certExpires, DaysLeft: 7, ThresholdDays: 7})
	var p map[string]any
	if err := json.Unmarshal([]byte(<-bodies), &p); err != nil {
		t.Fatal(err)
	}
	mon, _ := p["monitor"].(map[string]any)
	if p["event"] != "cert_expiring" || p["cert_expires_at"] != "2026-11-08T09:00:00Z" ||
		p["days_left"] != float64(7) || p["threshold_days"] != float64(7) ||
		p["dashboard_url"] != "https://pingou.test/monitors/m1" ||
		mon["id"] != "m1" || mon["name"] != "API" || mon["url"] != "https://api.example.com" {
		t.Fatalf("payload: %v", p)
	}
	if _, err := time.Parse(time.RFC3339, p["timestamp"].(string)); err != nil {
		t.Fatalf("timestamp: %v", p["timestamp"])
	}

	// eventos down e up não ganham os campos novos
	notifyDown(t.Context(), n, certMonitor, &domain.Incident{StartedAt: time.Now()})
	down := <-bodies
	for _, k := range []string{"cert_expires_at", "days_left", "threshold_days"} {
		if strings.Contains(down, k) {
			t.Fatalf("down com %s: %s", k, down)
		}
	}
}

func TestCertMessageText(t *testing.T) {
	cases := []struct {
		days      int
		threshold int
		subject   string
		line      string
	}{
		{30, 30, "[Pingou] ⚠️ Certificado SSL de API expira em 30 dias", "Faltam: 30 dias"},
		{7, 7, "[Pingou] ⚠️ Certificado SSL de API expira em 7 dias", "Faltam: 7 dias"},
		{1, 7, "[Pingou] ⚠️ Certificado SSL de API expira em 1 dia", "Faltam: 1 dia"},
		{0, 7, "[Pingou] ⚠️ Certificado SSL de API expira hoje", "Faltam: 0 dias"},
		{-1, 7, "[Pingou] ⚠️ Certificado SSL de API expirou", "Situação: o certificado já expirou"},
	}
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	for _, c := range cases {
		a := domain.CertAlert{ExpiresAt: certExpires, DaysLeft: c.days, ThresholdDays: c.threshold}
		if got := subjectCert(certMonitor, a); got != c.subject {
			t.Errorf("assunto %d dias: %q, want %q", c.days, got, c.subject)
		}
		body := certBody(certMonitor, a, "https://pingou.test", loc)
		for _, s := range []string{
			"Monitor: API", "URL: https://api.example.com", "Expira em: 2026-11-08 06:00:00 -03", c.line,
			"Detalhes: https://pingou.test/monitors/m1",
		} {
			if !strings.Contains(body, s) {
				t.Errorf("%d dias: corpo sem %q:\n%s", c.days, s, body)
			}
		}
	}
}

func TestEmailNotifierSendsCertExpiring(t *testing.T) {
	f := newFakeSMTP(t)
	n := NewEmailNotifier(f.cfg(config.SMTPTLSNone), func() string { return "ops@exemplo.com" }, "")
	notifyCertExpiring(t.Context(), n, certMonitor, domain.CertAlert{ExpiresAt: certExpires, DaysLeft: 7, ThresholdDays: 7})
	<-f.done
	f.mu.Lock()
	defer f.mu.Unlock()
	if !strings.Contains(f.data, "Subject: =?utf-8?q?") || !strings.Contains(f.data, "Certificado") || len(f.rcpts) != 1 {
		t.Fatalf("rcpts=%v DATA:\n%s", f.rcpts, f.data)
	}
}

func TestNtfyNotifierCertExpiring(t *testing.T) {
	f := newFakeNtfy(t, http.StatusOK)
	n := NewNtfyNotifier(config.NtfyConfig{URL: f.srv.URL}, func() string { return "alertas" }, "https://pingou.test")

	for _, c := range []struct {
		threshold, days int
		priority, title string
	}{
		{30, 25, "default", "Certificado SSL de API expira em 25 dias"},
		{7, 7, "high", "Certificado SSL de API expira em 7 dias"},
	} {
		notifyCertExpiring(t.Context(), n, certMonitor, domain.CertAlert{ExpiresAt: certExpires, DaysLeft: c.days, ThresholdDays: c.threshold})
		<-f.got
		r := f.last(t)
		if got := decodedTitle(t, r.header); got != c.title {
			t.Fatalf("Title = %q, want %q", got, c.title)
		}
		if r.header.Get("Priority") != c.priority || r.header.Get("Tags") != "lock,warning" ||
			r.header.Get("Click") != "https://pingou.test/monitors/m1" {
			t.Fatalf("%d dias: headers %v", c.threshold, r.header)
		}
		if !strings.Contains(r.body, "Expira em: 2026-11-08 09:00:00 UTC") {
			t.Fatalf("corpo: %s", r.body)
		}
	}
}
