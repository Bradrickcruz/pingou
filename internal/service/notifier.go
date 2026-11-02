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

type Notifier interface {
	NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident)
	NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident)
	// NotifyCertExpiring avisa que o certificado SSL de um monitor http está
	// perto de expirar (limite de 30 ou 7 dias).
	NotifyCertExpiring(ctx context.Context, m *domain.Monitor, a domain.CertAlert)
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

type webhookPayload struct {
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

func (n *WebhookNotifier) NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	payload := webhookPayload{
		Event:                   "down",
		Monitor:                 toWebhookMonitor(m),
		Timestamp:               i.StartedAt.UTC().Format(time.RFC3339),
		LastError:               i.LastError,
		DowntimeDurationSeconds: nil,
		DashboardURL:            n.dashboardURL(m),
	}

	n.send(url, payload)
}

func (n *WebhookNotifier) NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	timestamp := recoveredAt(i)
	downtimeDurationSeconds := downtimeSeconds(i, timestamp)

	payload := webhookPayload{
		Event:                   "up",
		Monitor:                 toWebhookMonitor(m),
		Timestamp:               timestamp.Format(time.RFC3339),
		LastError:               nil,
		DowntimeDurationSeconds: &downtimeDurationSeconds,
		DashboardURL:            n.dashboardURL(m),
	}

	n.send(url, payload)
}

// NotifyCertExpiring envia o evento "cert_expiring". timestamp é o momento do
// aviso; cert_expires_at é a validade do certificado (UTC, RFC 3339).
func (n *WebhookNotifier) NotifyCertExpiring(ctx context.Context, m *domain.Monitor, a domain.CertAlert) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	expiresAt := a.ExpiresAt.UTC().Format(time.RFC3339)
	daysLeft, threshold := a.DaysLeft, a.ThresholdDays
	payload := webhookPayload{
		Event:         "cert_expiring",
		Monitor:       toWebhookMonitor(m),
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		DashboardURL:  n.dashboardURL(m),
		CertExpiresAt: &expiresAt,
		DaysLeft:      &daysLeft,
		ThresholdDays: &threshold,
	}

	n.send(url, payload)
}

func (n *WebhookNotifier) send(url string, payload webhookPayload) {
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("webhook marshal error", "err", err)
		return
	}

	resp, err := n.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Error("webhook send error", "err", err)
		return
	}
	defer resp.Body.Close()

	slog.Info("webhook sent", "event", payload.Event, "status", resp.StatusCode)
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
	payload := webhookPayload{
		Event: "test",
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

// webhookTestErrorMessage tira o prefixo `Post "<url>": ` do erro, para a URL
// (que pode ter token na query) não voltar na resposta.
func webhookTestErrorMessage(err error) string {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Timeout() {
			return fmt.Sprintf("timeout após %s", webhookTestTimeout)
		}
		return ue.Err.Error()
	}
	return err.Error()
}
