package domain

import (
	"context"
	"time"
)

type MonitorRepository interface {
	Create(ctx context.Context, m *Monitor) error
	FindAll(ctx context.Context, filter MonitorFilter) ([]*Monitor, int, error)
	FindByID(ctx context.Context, id string) (*Monitor, error)
	FindByPingToken(ctx context.Context, token string) (*Monitor, error)
	Update(ctx context.Context, m *Monitor) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type CheckRepository interface {
	Create(ctx context.Context, c *Check) error
	FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*Check, int, error)
	DeleteOlderThan(ctx context.Context, before string) (int64, error) // <- adiciona
	Summaries(ctx context.Context, monitorIDs []string, historySize int, since time.Time) (map[string]*MonitorSummary, error)
	Buckets(ctx context.Context, monitorID string, since, until time.Time, bucketSeconds int) ([]StatsBucket, error)
}

type IncidentRepository interface {
	Create(ctx context.Context, i *Incident) error
	FindOpenByMonitor(ctx context.Context, monitorID string) (*Incident, error)
	FindOpenByMonitors(ctx context.Context, monitorIDs []string) (map[string]*Incident, error)
	FindAll(ctx context.Context, filter IncidentFilter) ([]*Incident, int, error)
	Close(ctx context.Context, id string, endedAt string, durationSeconds int) error
	// CloseManually fecha o incidente aberto do monitor, se houver (sem notificar),
	// gravando o motivo do fechamento e uma cópia do texto do motivo.
	CloseManually(ctx context.Context, monitorID string, endedAt time.Time, reason CloseReason, note *string) error
	FindInWindow(ctx context.Context, monitorID string, since time.Time) ([]*Incident, error)
}

// MaintenanceRepository guarda as janelas de manutenção. "Viva" = ainda não
// finalizada (agendada ou em curso); há no máximo uma por monitor.
type MaintenanceRepository interface {
	// CreateBatch insere todas as janelas numa transação (tudo ou nada).
	CreateBatch(ctx context.Context, windows []*MaintenanceWindow) error
	FindByID(ctx context.Context, id string) (*MaintenanceWindow, error)
	FindLiveByMonitor(ctx context.Context, monitorID string) (*MaintenanceWindow, error)
	// FindLiveByMonitors devolve a janela viva de cada monitor, indexada pelo id.
	FindLiveByMonitors(ctx context.Context, monitorIDs []string) (map[string]*MaintenanceWindow, error)
	// ActiveForMonitor: consulta quente do gate (a cada check e ping).
	ActiveForMonitor(ctx context.Context, monitorID string, now time.Time) (bool, error)
	// Delete apaga uma janela que ainda não começou.
	Delete(ctx context.Context, id string) error
	// DueStart: janelas cujo início chegou e cujos efeitos de início não rodaram.
	DueStart(ctx context.Context, now time.Time) ([]*MaintenanceWindow, error)
	// DueEnd: janelas cujo fim chegou e que ainda não foram finalizadas.
	DueEnd(ctx context.Context, now time.Time) ([]*MaintenanceWindow, error)
	// Start aplica o início numa transação: fecha o incidente aberto do monitor
	// (motivo manutenção), põe o estado em UNKNOWN e marca started_at.
	// Devolve false se já tinha sido aplicado.
	Start(ctx context.Context, w *MaintenanceWindow, now time.Time) (bool, error)
	// Finalize aplica o fim numa transação (e o início antes, se não rodou):
	// marca finalized_at, grava resumed_at = now, estado UNKNOWN e recomeça o
	// prazo do heartbeat. endsAt != nil encurta a janela (cancelamento).
	// Devolve false se já estava finalizada.
	Finalize(ctx context.Context, w *MaintenanceWindow, now time.Time, endsAt *time.Time) (bool, error)
}

type MonitorFilter struct {
	Enabled *bool
	State   *MonitorState
	Type    *MonitorType
	Limit   int
	Offset  int
}
