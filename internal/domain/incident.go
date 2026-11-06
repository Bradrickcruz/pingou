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
	FailureCount     int64       // checks com falha enquanto o incidente esteve aberto
	CloseReason      CloseReason // por que foi fechado sem recuperação (CloseNone = recuperou ou aberto)
	CloseNote        *string     // motivo da pausa ou da manutenção no momento do fechamento

	// Dados do monitor (JOIN), preenchidos só pelas listagens.
	MonitorName   string
	MonitorType   MonitorType
	MonitorTarget *string
}

// CloseReason é guardado em incidents.closed_manually (0/1/2).
type CloseReason int

const (
	CloseNone        CloseReason = 0 // aberto ou fechado por recuperação
	ClosePause       CloseReason = 1 // fechado ao pausar o monitor
	CloseMaintenance CloseReason = 2 // fechado no início de uma janela de manutenção
)

// String devolve o nome usado na API ("" para CloseNone).
func (r CloseReason) String() string {
	switch r {
	case ClosePause:
		return "pause"
	case CloseMaintenance:
		return "maintenance"
	}
	return ""
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
