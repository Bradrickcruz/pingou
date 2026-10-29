package domain

import "time"

type Incident struct {
	ID               string
	MonitorID        string
	StartedAt        time.Time
	EndedAt          *time.Time
	LastError        *string
	NotificationSent bool
	DurationSeconds  *int64
	FailureCount     int64 // checks com falha enquanto o incidente esteve aberto

	// Dados do monitor (JOIN), preenchidos só pelas listagens.
	MonitorName   string
	MonitorType   MonitorType
	MonitorTarget *string
}

const (
	IncidentStatusOpen     = "open"
	IncidentStatusResolved = "resolved"
)

// IncidentFilter: campos vazios não filtram. Status é "", "open" ou "resolved".
type IncidentFilter struct {
	MonitorID string
	Status    string
	Limit     int
	Offset    int
}
