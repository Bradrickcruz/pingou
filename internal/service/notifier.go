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
	Event                   string                `json:"event"`
	Monitor                 webhookPayloadMonitor `json:"monitor"`
	Timestamp               string                `json:"timestamp"`
	LastError               *string               `json:"last_error"`
	DowntimeDurationSeconds *int64                `json:"downtime_duration_seconds"`
}

type webhookPayloadMonitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (n *WebhookNotifier) NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	payload := webhookPayload{
		Event: "down",
		Monitor: webhookPayloadMonitor{
			ID:   m.ID,
			Name: m.Name,
			URL:  m.URL,
		},
		Timestamp:               i.StartedAt.UTC().Format(time.RFC3339),
		LastError:               i.LastError,
		DowntimeDurationSeconds: nil,
	}

	n.send(url, payload)
}

func (n *WebhookNotifier) NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	url := n.getWebhookURL()
	if url == "" {
		return
	}

	timestamp := time.Now().UTC()
	if i.EndedAt != nil {
		timestamp = i.EndedAt.UTC()
	}

	var downtimeDurationSeconds *int64
	if i.DurationSeconds != nil {
		downtimeDurationSeconds = i.DurationSeconds
	} else {
		d := int64(timestamp.Sub(i.StartedAt).Seconds())
		if d < 0 {
			d = 0
		}
		downtimeDurationSeconds = &d
	}

	payload := webhookPayload{
		Event: "up",
		Monitor: webhookPayloadMonitor{
			ID:   m.ID,
			Name: m.Name,
			URL:  m.URL,
		},
		Timestamp:               timestamp.Format(time.RFC3339),
		LastError:               nil,
		DowntimeDurationSeconds: downtimeDurationSeconds,
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
