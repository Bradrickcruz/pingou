package service

import (
	"testing"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

func TestValidateURL(t *testing.T) {
	valid := []string{
		"http://example.com",
		"https://example.com/health?x=1",
		"http://127.0.0.1:8080/",
		"http://localhost:65535/status", // porta com 5 dígitos
	}
	for _, u := range valid {
		if err := validateURL(u); err != nil {
			t.Errorf("validateURL(%q) = %v, want nil", u, err)
		}
	}

	invalid := []string{"", "example.com", "ftp://example.com", "http://example.com:123456"}
	for _, u := range invalid {
		if err := validateURL(u); err == nil {
			t.Errorf("validateURL(%q) = nil, want erro", u)
		}
	}
}

func TestValidateMonitorByType(t *testing.T) {
	base := func(typ domain.MonitorType) *domain.Monitor {
		return &domain.Monitor{
			Name: "x", Type: typ, URL: "https://x.test", IntervalSeconds: 60,
			TimeoutSeconds: 5, FailureThreshold: 3,
			Config: domain.MonitorConfig{Host: "db.interno", Port: 5432, GraceSeconds: 60},
		}
	}

	cases := []struct {
		name  string
		edit  func(m *domain.Monitor)
		typ   domain.MonitorType
		valid bool
	}{
		{"http válido", func(*domain.Monitor) {}, domain.TypeHTTP, true},
		{"http sem url", func(m *domain.Monitor) { m.URL = "" }, domain.TypeHTTP, false},
		{"tcp válido", func(m *domain.Monitor) { m.URL = "" }, domain.TypeTCP, true},
		{"tcp porta 0", func(m *domain.Monitor) { m.Config.Port = 0 }, domain.TypeTCP, false},
		{"tcp porta 65536", func(m *domain.Monitor) { m.Config.Port = 65536 }, domain.TypeTCP, false},
		{"tcp porta 65535", func(m *domain.Monitor) { m.Config.Port = 65535 }, domain.TypeTCP, true},
		{"tcp host com esquema", func(m *domain.Monitor) { m.Config.Host = "http://db" }, domain.TypeTCP, false},
		{"tcp host com caminho", func(m *domain.Monitor) { m.Config.Host = "db/x" }, domain.TypeTCP, false},
		{"tcp host vazio", func(m *domain.Monitor) { m.Config.Host = "" }, domain.TypeTCP, false},
		{"tcp intervalo 1 semana", func(m *domain.Monitor) { m.IntervalSeconds = 604800 }, domain.TypeTCP, false},
		{"heartbeat 1 semana", func(m *domain.Monitor) { m.IntervalSeconds = 604800 }, domain.TypeHeartbeat, true},
		{"heartbeat 30 s", func(m *domain.Monitor) { m.IntervalSeconds = 30 }, domain.TypeHeartbeat, false},
		{"heartbeat acima de 1 semana", func(m *domain.Monitor) { m.IntervalSeconds = 604801 }, domain.TypeHeartbeat, false},
		{"heartbeat grace 0", func(m *domain.Monitor) { m.Config.GraceSeconds = 0 }, domain.TypeHeartbeat, true},
		{"heartbeat grace negativo", func(m *domain.Monitor) { m.Config.GraceSeconds = -1 }, domain.TypeHeartbeat, false},
		{"heartbeat grace > 1 dia", func(m *domain.Monitor) { m.Config.GraceSeconds = 86401 }, domain.TypeHeartbeat, false},
		{"heartbeat ignora timeout", func(m *domain.Monitor) { m.TimeoutSeconds = 0 }, domain.TypeHeartbeat, true},
		{"tipo desconhecido", func(*domain.Monitor) {}, "dns", false},
	}
	for _, c := range cases {
		m := base(c.typ)
		c.edit(m)
		err := validateMonitor(m)
		if c.valid && err != nil {
			t.Errorf("%s: erro %v, want válido", c.name, err)
		}
		if !c.valid && err == nil {
			t.Errorf("%s: válido, want erro", c.name)
		}
	}
}
