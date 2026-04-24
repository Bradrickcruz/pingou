package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type Notifier interface {
	NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident)
	NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident)
}

// WebhookNotifier envia payloads para a URL configurada nas settings
type WebhookNotifier struct {
	getWebhookURL func() string
	client        *http.Client
}

func NewWebhookNotifier(getWebhookURL func() string) *WebhookNotifier {
	return &WebhookNotifier{
		getWebhookURL: getWebhookURL,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	Event     string `json:"event"`
	MonitorID string `json:"monitor_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	State     string `json:"state"`
	Error     string `json:"error,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
	EndedAt   string `json:"ended_at,omitempty"`
	Timestamp string `json:"timestamp"`
}

func (n *WebhookNotifier) NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	payload := webhookPayload{
		Event:     "monitor.down",
		MonitorID: m.ID,
		Name:      m.Name,
		URL:       m.URL,
		State:     "DOWN",
		StartedAt: i.StartedAt.Format(time.RFC3339),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if i.LastError != nil {
		payload.Error = *i.LastError
	}

	n.send(url, payload)
}

func (n *WebhookNotifier) NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	payload := webhookPayload{
		Event:     "monitor.recovered",
		MonitorID: m.ID,
		Name:      m.Name,
		URL:       m.URL,
		State:     "UP",
		StartedAt: i.StartedAt.Format(time.RFC3339),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if i.EndedAt != nil {
		payload.EndedAt = i.EndedAt.Format(time.RFC3339)
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
