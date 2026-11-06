package domain

import "time"

// Canais de aviso (coluna channel do outbox).
const (
	ChannelWebhook = "webhook"
	ChannelEmail   = "email"
	ChannelNtfy    = "ntfy"
)

// Tipos de evento (coluna event_type do outbox; também o campo "event" do webhook).
const (
	EventDown         = "down"
	EventUp           = "up"
	EventCertExpiring = "cert_expiring"
)

// OutboxStatus é o estado de uma linha do outbox.
type OutboxStatus string

const (
	OutboxPending OutboxStatus = "pending"
	OutboxSent    OutboxStatus = "sent"
	OutboxExpired OutboxStatus = "expired"
)

// OutboxEvent é um aviso para um canal. EventID é igual em todos os canais e
// tentativas do mesmo evento (o receptor do webhook usa para ignorar duplicata).
type OutboxEvent struct {
	ID            int64
	EventID       string
	EventType     string
	Channel       string
	MonitorID     string
	Snapshot      EventSnapshot
	Status        OutboxStatus
	Attempts      int
	CreatedAt     time.Time
	ExpiresAt     time.Time
	NextAttemptAt time.Time
	LastError     *string
	DoneAt        *time.Time
}

// EventSnapshot é o retrato do evento no momento em que aconteceu (coluna
// payload, JSON). O monitor pode ser editado ou apagado antes do envio, e o
// aviso deve dizer o que valia no evento. Nunca leva segredo: nada de
// ping_token, URL do webhook, destinatários ou tópico.
type EventSnapshot struct {
	Monitor    SnapshotMonitor   `json:"monitor"`
	Incident   *SnapshotIncident `json:"incident,omitempty"`
	Cert       *SnapshotCert     `json:"cert,omitempty"`
	OccurredAt time.Time         `json:"occurred_at"`
}

// SnapshotMonitor guarda o mínimo para Monitor.Target() e para os textos.
type SnapshotMonitor struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	Type MonitorType `json:"type"`
	URL  string      `json:"url,omitempty"`  // http
	Host string      `json:"host,omitempty"` // tcp
	Port int         `json:"port,omitempty"` // tcp
}

type SnapshotIncident struct {
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	DurationSeconds *int64     `json:"duration_seconds,omitempty"`
	LastError       *string    `json:"last_error,omitempty"`
	FailureCount    int64      `json:"failure_count"`
}

type SnapshotCert struct {
	ExpiresAt     time.Time `json:"expires_at"`
	DaysLeft      int       `json:"days_left"`
	ThresholdDays int       `json:"threshold_days"`
}

func snapshotMonitor(m *Monitor) SnapshotMonitor {
	s := SnapshotMonitor{ID: m.ID, Name: m.Name, Type: m.Type}
	switch m.Type {
	case TypeHTTP:
		s.URL = m.URL
	case TypeTCP:
		s.Host, s.Port = m.Config.Host, m.Config.Port
	}
	return s
}

func snapshotIncident(i *Incident) *SnapshotIncident {
	return &SnapshotIncident{
		StartedAt:       i.StartedAt.UTC(),
		EndedAt:         utcPtr(i.EndedAt),
		DurationSeconds: i.DurationSeconds,
		LastError:       i.LastError,
		FailureCount:    i.FailureCount,
	}
}

func utcPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	u := t.UTC()
	return &u
}

// NewDownSnapshot é o retrato de uma queda: ocorreu no início do incidente.
func NewDownSnapshot(m *Monitor, i *Incident) EventSnapshot {
	return EventSnapshot{Monitor: snapshotMonitor(m), Incident: snapshotIncident(i), OccurredAt: i.StartedAt.UTC()}
}

// NewRecoverySnapshot é o retrato de uma recuperação: ocorreu no fim do
// incidente (ou em at, se o incidente não tiver fim).
func NewRecoverySnapshot(m *Monitor, i *Incident, at time.Time) EventSnapshot {
	if i.EndedAt != nil {
		at = *i.EndedAt
	}
	return EventSnapshot{Monitor: snapshotMonitor(m), Incident: snapshotIncident(i), OccurredAt: at.UTC()}
}

// NewCertSnapshot é o retrato de um aviso de certificado, decidido em at.
func NewCertSnapshot(m *Monitor, a CertAlert, at time.Time) EventSnapshot {
	return EventSnapshot{
		Monitor:    snapshotMonitor(m),
		Cert:       &SnapshotCert{ExpiresAt: a.ExpiresAt.UTC(), DaysLeft: a.DaysLeft, ThresholdDays: a.ThresholdDays},
		OccurredAt: at.UTC(),
	}
}

// AsMonitor reconstrói o monitor (sem ping_token) para os textos e o Target().
func (s EventSnapshot) AsMonitor() *Monitor {
	return &Monitor{
		ID: s.Monitor.ID, Name: s.Monitor.Name, Type: s.Monitor.Type, URL: s.Monitor.URL,
		Config: MonitorConfig{Host: s.Monitor.Host, Port: s.Monitor.Port},
	}
}

// AsIncident reconstrói o incidente; nunca nil (evento de certificado devolve vazio).
func (s EventSnapshot) AsIncident() *Incident {
	if s.Incident == nil {
		return &Incident{MonitorID: s.Monitor.ID, StartedAt: s.OccurredAt}
	}
	return &Incident{
		MonitorID:       s.Monitor.ID,
		StartedAt:       s.Incident.StartedAt,
		EndedAt:         s.Incident.EndedAt,
		DurationSeconds: s.Incident.DurationSeconds,
		LastError:       s.Incident.LastError,
		FailureCount:    s.Incident.FailureCount,
	}
}

// AsCertAlert reconstrói o aviso de certificado (vazio fora do evento cert_expiring).
func (s EventSnapshot) AsCertAlert() CertAlert {
	if s.Cert == nil {
		return CertAlert{}
	}
	return CertAlert{ExpiresAt: s.Cert.ExpiresAt, DaysLeft: s.Cert.DaysLeft, ThresholdDays: s.Cert.ThresholdDays}
}
