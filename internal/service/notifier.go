package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// Channel é um canal de aviso (webhook, e-mail, ntfy). O OutboxWorker chama
// Send para cada linha do outbox e repete se devolver erro.
//
// Send lê o destino (URL, destinatários, tópico) das settings na hora e monta
// o texto a partir do snapshot do evento. Destino vazio = canal desligado:
// devolve nil sem enviar (a linha encerra, sem retry eterno). O erro devolvido
// nunca leva URL, tópico, endereço nem token.
type Channel interface {
	Name() string
	Send(ctx context.Context, ev domain.OutboxEvent) error
}

// WebhookNotifier envia payloads para a URL configurada nas settings
type WebhookNotifier struct {
	getWebhookURL func() string
	publicURL     string // PINGOU_PUBLIC_URL sem "/" final; vazio = sem dashboard_url
	client        *http.Client
}

func NewWebhookNotifier(getWebhookURL func() string, publicURL string) *WebhookNotifier {
	return &WebhookNotifier{
		getWebhookURL: getWebhookURL,
		publicURL:     publicURL,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

// dashboardURL é o link para a página de detalhe do monitor, ou nil sem PINGOU_PUBLIC_URL.
func (n *WebhookNotifier) dashboardURL(m *domain.Monitor) *string {
	return dashboardLink(n.publicURL, m)
}

// dashboardLink monta o link do monitor no dashboard. Vale para todos os canais.
func dashboardLink(publicURL string, m *domain.Monitor) *string {
	if publicURL == "" {
		return nil
	}
	u := publicURL + "/monitors/" + m.ID
	return &u
}

// recoveredAt é o momento da volta: o fim do incidente ou, sem ele, agora.
func recoveredAt(i *domain.Incident) time.Time {
	if i.EndedAt != nil {
		return i.EndedAt.UTC()
	}
	return time.Now().UTC()
}

// downtimeSeconds é quanto o monitor ficou fora: a duração gravada no incidente
// ou, sem ela, de StartedAt até end. Nunca negativo.
func downtimeSeconds(i *domain.Incident, end time.Time) int64 {
	if i.DurationSeconds != nil {
		return *i.DurationSeconds
	}
	return max(int64(end.Sub(i.StartedAt).Seconds()), 0)
}

// Nomes do campo "event" do webhook. Padrão: <objeto>_<acontecimento>, minúsculas,
// separado por "_". Também servem de rótulo de log nos outros canais.
const (
	eventMonitorDown  = "monitor_down"
	eventMonitorUp    = "monitor_up"
	eventCertExpiring = "cert_expiring"
	eventMaintStart   = "maintenance_started"
	eventMaintEnd     = "maintenance_ended"
	eventTest         = "test"
)

type webhookPayload struct {
	// ID é o event_id do outbox: igual em todas as tentativas do mesmo evento,
	// para o receptor ignorar duplicata (a entrega é "pelo menos uma vez").
	ID                      string                `json:"id"`
	Event                   string                `json:"event"`
	Monitor                 webhookPayloadMonitor `json:"monitor"`
	Timestamp               string                `json:"timestamp"`
	LastError               *string               `json:"last_error"`
	DowntimeDurationSeconds *int64                `json:"downtime_duration_seconds"`
	DashboardURL            *string               `json:"dashboard_url"`
	// só no evento cert_expiring (omitidos nos demais, que não mudam)
	CertExpiresAt *string `json:"cert_expires_at,omitempty"`
	DaysLeft      *int    `json:"days_left,omitempty"`
	ThresholdDays *int    `json:"threshold_days,omitempty"`
}

// webhookMaintenancePayload é o payload dos eventos maintenance_started e
// maintenance_ended. O aviso é de um grupo de monitores: usa "monitors" (lista,
// mesmo com um só) e não leva "monitor". reason sai null sem motivo.
type webhookMaintenancePayload struct {
	ID           string                  `json:"id"`
	Event        string                  `json:"event"`
	Monitors     []webhookPayloadMonitor `json:"monitors"`
	StartsAt     string                  `json:"starts_at"`
	EndsAt       string                  `json:"ends_at"`
	Reason       *string                 `json:"reason"`
	Timestamp    string                  `json:"timestamp"`
	DashboardURL *string                 `json:"dashboard_url"` // página /maintenance
}

// webhookPayloadMonitor nunca leva o ping_token: webhooks costumam ir para
// terceiros (Slack, Discord) e o token é o segredo da URL de ping.
type webhookPayloadMonitor struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Target *string `json:"target"` // URL (http), host:porta (tcp), null (heartbeat)
	URL    *string `json:"url"`    // só http, mantido por compatibilidade
}

func toWebhookMonitor(m *domain.Monitor) webhookPayloadMonitor {
	out := webhookPayloadMonitor{ID: m.ID, Name: m.Name, Type: string(m.Type), Target: m.Target()}
	if m.Type == domain.TypeHTTP {
		out.URL = out.Target
	}
	return out
}

func (n *WebhookNotifier) Name() string { return domain.ChannelWebhook }

// Send entrega o evento. Erro de rede ou status fora de 2xx é falha (o worker repete).
func (n *WebhookNotifier) Send(ctx context.Context, ev domain.OutboxEvent) error {
	url := n.getWebhookURL()
	if url == "" {
		return nil
	}
	if ev.EventType == domain.EventMaintenanceStarted || ev.EventType == domain.EventMaintenanceEnded {
		payload, err := n.maintenancePayloadFor(ev)
		if err != nil {
			return err
		}
		return n.post(ctx, url, payload)
	}
	payload, err := n.payloadFor(ev)
	if err != nil {
		return err
	}
	return n.post(ctx, url, payload)
}

// maintenancePayloadFor monta o JSON dos eventos de manutenção a partir do snapshot.
func (n *WebhookNotifier) maintenancePayloadFor(ev domain.OutboxEvent) (webhookMaintenancePayload, error) {
	mt, err := maintenanceSnapshot(ev)
	if err != nil {
		return webhookMaintenancePayload{}, err
	}
	event := eventMaintStart
	if ev.EventType == domain.EventMaintenanceEnded {
		event = eventMaintEnd
	}
	monitors := make([]webhookPayloadMonitor, 0, len(mt.Monitors))
	for _, sm := range mt.Monitors {
		monitors = append(monitors, toWebhookMonitor((&domain.EventSnapshot{Monitor: sm}).AsMonitor()))
	}
	return webhookMaintenancePayload{
		ID:           ev.EventID,
		Event:        event,
		Monitors:     monitors,
		StartsAt:     mt.StartsAt.UTC().Format(time.RFC3339),
		EndsAt:       mt.EndsAt.UTC().Format(time.RFC3339),
		Reason:       mt.Reason,
		Timestamp:    ev.Snapshot.OccurredAt.UTC().Format(time.RFC3339),
		DashboardURL: maintenanceLink(n.publicURL),
	}, nil
}

// payloadFor monta o JSON do evento a partir do snapshot. O timestamp é o
// momento do evento, não o do envio: numa nova tentativa a mensagem não muda.
func (n *WebhookNotifier) payloadFor(ev domain.OutboxEvent) (webhookPayload, error) {
	snap := ev.Snapshot
	m := snap.AsMonitor()
	payload := webhookPayload{
		ID:           ev.EventID,
		Monitor:      toWebhookMonitor(m),
		Timestamp:    snap.OccurredAt.UTC().Format(time.RFC3339),
		DashboardURL: n.dashboardURL(m),
	}
	switch ev.EventType {
	case domain.EventDown:
		payload.Event = eventMonitorDown
		payload.LastError = snap.AsIncident().LastError
	case domain.EventUp:
		payload.Event = eventMonitorUp
		i := snap.AsIncident()
		downtime := downtimeSeconds(i, recoveredAt(i))
		payload.DowntimeDurationSeconds = &downtime
	case domain.EventCertExpiring:
		payload.Event = eventCertExpiring
		a := snap.AsCertAlert()
		expiresAt := a.ExpiresAt.UTC().Format(time.RFC3339)
		daysLeft, threshold := a.DaysLeft, a.ThresholdDays
		payload.CertExpiresAt = &expiresAt
		payload.DaysLeft = &daysLeft
		payload.ThresholdDays = &threshold
	default:
		return payload, fmt.Errorf("unknown event type %q", ev.EventType)
	}
	return payload, nil
}

// webhookSendTimeout limita cada POST do outbox.
const webhookSendTimeout = 10 * time.Second

func (n *WebhookNotifier) post(ctx context.Context, rawURL string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.New("could not encode payload")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return errors.New("invalid url")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return errors.New(urlErrorMessage(err, webhookSendTimeout))
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, webhookTestBodyLimit))

	// antes, qualquer resposta contava como "sent"; agora 4xx/5xx é falha
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// urlErrorMessage tira o prefixo `Post "<url>": ` do erro, para a URL (que pode
// ter token na query) não ir para o log nem para o last_error.
func urlErrorMessage(err error, timeout time.Duration) string {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Timeout() {
			return fmt.Sprintf("timeout após %s", timeout)
		}
		return ue.Err.Error()
	}
	return err.Error()
}

// WebhookTestResult é a resposta de POST /api/settings/webhook/test. O corpo
// devolvido pelo destino nunca entra aqui.
type WebhookTestResult struct {
	OK         bool    `json:"ok"`
	StatusCode *int    `json:"status_code"`
	DurationMS int64   `json:"duration_ms"`
	Error      *string `json:"error"`
}

// webhookTestTimeout é variável só para os testes rodarem rápido.
var webhookTestTimeout = 5 * time.Second

// webhookTestBodyLimit: quanto do corpo da resposta é lido (e descartado).
const webhookTestBodyLimit = 4 << 10

// newWebhookTestClient não segue redirects: o envio real transforma POST em GET
// num 301/302 e não entrega, então o teste mostra o redirect como falha.
func newWebhookTestClient() *http.Client {
	return &http.Client{
		Timeout: webhookTestTimeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// SendTestWebhook envia um payload de teste (event "test") para rawURL, no mesmo
// formato dos eventos reais. Não grava nada no banco. O log leva só o host.
func SendTestWebhook(ctx context.Context, rawURL string) WebhookTestResult {
	target := "https://exemplo.com"
	testID, _ := newUUIDv7()
	payload := webhookPayload{
		ID:    testID,
		Event: eventTest,
		Monitor: webhookPayloadMonitor{
			ID:     "00000000-0000-0000-0000-000000000000",
			Name:   "Teste do Pingou",
			Type:   string(domain.TypeHTTP),
			Target: &target,
			URL:    &target,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	fail := func(start time.Time, status *int, msg string) WebhookTestResult {
		return WebhookTestResult{StatusCode: status, DurationMS: time.Since(start).Milliseconds(), Error: &msg}
	}

	host := ""
	if u, err := url.Parse(rawURL); err == nil {
		host = u.Host
	}

	start := time.Now()
	body, err := json.Marshal(payload)
	if err != nil {
		return fail(start, nil, "could not encode payload")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return fail(start, nil, "invalid url")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := newWebhookTestClient().Do(req)
	if err != nil {
		msg := webhookTestErrorMessage(err)
		slog.Info("webhook test", "host", host, "error", msg)
		return fail(start, nil, msg)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, webhookTestBodyLimit))

	status := resp.StatusCode
	slog.Info("webhook test", "host", host, "status", status)
	if status < 200 || status > 299 {
		return fail(start, &status, fmt.Sprintf("unexpected status code: %d", status))
	}
	return WebhookTestResult{OK: true, StatusCode: &status, DurationMS: time.Since(start).Milliseconds()}
}

// webhookTestErrorMessage: ver urlErrorMessage.
func webhookTestErrorMessage(err error) string {
	return urlErrorMessage(err, webhookTestTimeout)
}
