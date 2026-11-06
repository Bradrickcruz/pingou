package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

var noticeAt = time.Date(2026, 11, 8, 2, 0, 4, 0, time.UTC)

func maintenanceEvent(eventType string, n int, reason *string) domain.OutboxEvent {
	token := "segredo-do-token"
	var ws []*domain.MaintenanceWindow
	var ms []*domain.Monitor
	for i := range n {
		id := fmt.Sprintf("m%d", i)
		ws = append(ws, &domain.MaintenanceWindow{
			MonitorID: id, StartsAt: noticeAt.Add(-4 * time.Second), EndsAt: noticeAt.Add(2 * time.Hour), Reason: reason,
		})
		ms = append(ms, &domain.Monitor{ID: id, Name: fmt.Sprintf("API %d", i), Type: domain.TypeHTTP,
			URL: "https://x.test/" + id, PingToken: &token})
	}
	return testEvent(eventType, domain.NewMaintenanceSnapshot(ws, ms, noticeAt))
}

func webhookBody(t *testing.T, ev domain.OutboxEvent, publicURL string) (string, map[string]any) {
	t.Helper()
	bodies := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies <- string(b)
	}))
	defer srv.Close()
	if err := NewWebhookNotifier(func() string { return srv.URL }, publicURL).Send(t.Context(), ev); err != nil {
		t.Fatal(err)
	}
	body := <-bodies
	var p map[string]any
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatal(err)
	}
	return body, p
}

func TestWebhookMaintenancePayload(t *testing.T) {
	reason := "troca de disco"
	ev := maintenanceEvent(domain.EventMaintenanceStarted, 1, &reason)
	ev.EventID = "0193-abc"
	body, _ := webhookBody(t, ev, "https://pingou.exemplo.com")

	want := `{"id":"0193-abc","event":"maintenance_started","monitors":[{"id":"m0","name":"API 0","type":"http",` +
		`"target":"https://x.test/m0","url":"https://x.test/m0"}],"starts_at":"2026-11-08T02:00:00Z",` +
		`"ends_at":"2026-11-08T04:00:04Z","reason":"troca de disco","timestamp":"2026-11-08T02:00:04Z",` +
		`"dashboard_url":"https://pingou.exemplo.com/maintenance"}`
	if body != want {
		t.Fatalf("payload:\n got %s\nwant %s", body, want)
	}
	if strings.Contains(body, "ping_token") || strings.Contains(body, "segredo") {
		t.Fatalf("webhook vazou o token: %s", body)
	}
}

func TestWebhookMaintenanceEndedNilReasonAndNoPublicURL(t *testing.T) {
	_, p := webhookBody(t, maintenanceEvent(domain.EventMaintenanceEnded, 3, nil), "")
	if p["event"] != "maintenance_ended" {
		t.Fatalf("event=%v", p["event"])
	}
	if _, ok := p["monitor"]; ok {
		t.Fatal("evento de manutenção não leva o campo monitor")
	}
	if v, ok := p["reason"]; !ok || v != nil {
		t.Fatalf("reason deve sair null: %v (presente=%v)", v, ok)
	}
	if v, ok := p["dashboard_url"]; !ok || v != nil {
		t.Fatalf("dashboard_url deve sair null sem PINGOU_PUBLIC_URL: %v", v)
	}
	if ms, _ := p["monitors"].([]any); len(ms) != 3 {
		t.Fatalf("monitors=%v, want 3", p["monitors"])
	}
}

// Os eventos existentes seguem com o mesmo formato.
func TestWebhookMonitorDownUnchangedByMaintenance(t *testing.T) {
	m := &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}
	i := &domain.Incident{StartedAt: noticeAt}
	body, p := webhookBody(t, testEvent(domain.EventDown, domain.NewDownSnapshot(m, i)), "")
	for _, k := range []string{"monitor", "last_error", "downtime_duration_seconds", "dashboard_url"} {
		if _, ok := p[k]; !ok {
			t.Errorf("campo %q sumiu: %s", k, body)
		}
	}
	for _, k := range []string{"monitors", "starts_at", "ends_at", "reason"} {
		if _, ok := p[k]; ok {
			t.Errorf("campo %q não deveria existir em monitor_down: %s", k, body)
		}
	}
}

func TestMaintenanceWithoutDataIsError(t *testing.T) {
	ev := testEvent(domain.EventMaintenanceStarted, domain.EventSnapshot{OccurredAt: noticeAt})
	n := NewWebhookNotifier(func() string { return "http://127.0.0.1:1" }, "")
	if err := n.Send(t.Context(), ev); err == nil {
		t.Fatal("snapshot sem manutenção deveria dar erro")
	}
}

func TestMaintenanceTextSingleAndMany(t *testing.T) {
	reason := "troca de disco\r\nBcc: x@y"
	one := maintenanceEvent(domain.EventMaintenanceStarted, 1, &reason).Snapshot.Maintenance
	many := maintenanceEvent(domain.EventMaintenanceEnded, 3, nil).Snapshot.Maintenance

	if got := subjectMaintenance(domain.EventMaintenanceStarted, one); got != "[Pingou] 🔧 Manutenção começou: API 0" {
		t.Fatalf("assunto (1): %q", got)
	}
	if got := subjectMaintenance(domain.EventMaintenanceEnded, many); got != "[Pingou] 🔧 Manutenção terminou em 3 monitores" {
		t.Fatalf("assunto (3): %q", got)
	}

	sp := time.FixedZone("-03", -3*3600)
	body := maintenanceBody(domain.EventMaintenanceStarted, one, "https://p.test", sp)
	for _, want := range []string{"Monitor: API 0", "Início: 2026-11-07 23:00:00 -03", "Motivo: troca de disco Bcc: x@y", "Detalhes: https://p.test/maintenance"} {
		if !strings.Contains(body, want) {
			t.Errorf("corpo sem %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "terminou") {
		t.Errorf("aviso de início não deve dizer que terminou:\n%s", body)
	}

	// reason nil: sem a linha; fim: só diz que terminou
	end := maintenanceBody(domain.EventMaintenanceEnded, many, "", time.UTC)
	if strings.Contains(end, "Motivo") || strings.Contains(end, "Detalhes") {
		t.Errorf("sem motivo e sem link não deve ter as linhas:\n%s", end)
	}
	if !strings.Contains(end, "A manutenção terminou.") || strings.Contains(strings.ToLower(end), "fora do ar") {
		t.Errorf("aviso de fim:\n%s", end)
	}
}

func TestMaintenanceTextCutsLongList(t *testing.T) {
	mt := maintenanceEvent(domain.EventMaintenanceStarted, 100, nil).Snapshot.Maintenance
	body := maintenanceBody(domain.EventMaintenanceStarted, mt, "", time.UTC)
	if n := strings.Count(body, "- API "); n != maxNoticeMonitors {
		t.Fatalf("nomes listados=%d, want %d", n, maxNoticeMonitors)
	}
	if !strings.Contains(body, "- e mais 80") {
		t.Fatalf("falta o 'e mais 80':\n%s", body)
	}
	// o webhook leva a lista inteira
	_, p := webhookBody(t, maintenanceEvent(domain.EventMaintenanceStarted, 100, nil), "")
	if ms, _ := p["monitors"].([]any); len(ms) != 100 {
		t.Fatalf("webhook monitors=%d, want 100", len(ms))
	}
}

// Motivo com CR/LF não vira cabeçalho no e-mail.
func TestMaintenanceEmailNoHeaderInjection(t *testing.T) {
	reason := "x\r\nBcc: invasor@x.com"
	mt := maintenanceEvent(domain.EventMaintenanceStarted, 1, &reason).Snapshot.Maintenance
	raw := string(buildEmail("Pingou <a@x.com>", []string{"b@x.com"},
		subjectMaintenance(domain.EventMaintenanceStarted, mt), maintenanceBody(domain.EventMaintenanceStarted, mt, "", time.UTC), noticeAt))
	head, _, _ := strings.Cut(raw, "\r\n\r\n")
	if strings.Contains(head, "Bcc") {
		t.Fatalf("injeção de cabeçalho:\n%s", head)
	}
}

func TestNtfyMaintenance(t *testing.T) {
	f := newFakeNtfy(t, 200)
	n := NewNtfyNotifier(config.NtfyConfig{URL: f.srv.URL}, func() string { return "topico" }, "https://p.test")
	if err := n.Send(t.Context(), maintenanceEvent(domain.EventMaintenanceStarted, 2, nil)); err != nil {
		t.Fatal(err)
	}
	req := f.last(t)
	if got := decodedTitle(t, req.header); got != "Manutenção começou em 2 monitores" {
		t.Fatalf("título: %q", got)
	}
	if req.header.Get("Priority") != "default" || req.header.Get("Tags") != "wrench" || req.header.Get("Click") != "https://p.test/maintenance" {
		t.Fatalf("headers: %v", req.header)
	}
}

func TestNewMaintenanceEvents(t *testing.T) {
	var ws []*domain.MaintenanceWindow
	var ms []*domain.Monitor
	for i := range 3 {
		id := fmt.Sprintf("m%d", i)
		ws = append(ws, &domain.MaintenanceWindow{MonitorID: id, StartsAt: noticeAt, EndsAt: noticeAt.Add(time.Hour)})
		ms = append(ms, &domain.Monitor{ID: id, Name: id, Type: domain.TypeHTTP, URL: "https://x"})
	}
	rows, err := NewMaintenanceEvents(domain.EventMaintenanceEnded, ws, ms, []string{"webhook", "email"}, noticeAt)
	if err != nil || len(rows) != 2 {
		t.Fatalf("rows=%d err=%v, want 2 (um por canal)", len(rows), err)
	}
	if rows[0].EventID == "" || rows[0].EventID != rows[1].EventID || rows[0].MonitorID != "m0" {
		t.Fatalf("linhas: %+v", rows)
	}
	if !rows[0].ExpiresAt.Equal(noticeAt.Add(incidentNoticeTTL)) {
		t.Fatalf("prazo: %v", rows[0].ExpiresAt)
	}
	if rows, _ := NewMaintenanceEvents(domain.EventMaintenanceEnded, ws, ms, nil, noticeAt); len(rows) != 0 {
		t.Fatal("sem canais não há linha")
	}
}
