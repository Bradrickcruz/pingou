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
	// ListLive lista as janelas vivas com o nome do monitor, por início.
	ListLive(ctx context.Context) ([]*MaintenanceWindowView, error)
	// UpdateTimes troca início, fim e motivo de uma janela que ainda não
	// começou (starts_at > now, guardado no próprio UPDATE). false = nenhuma
	// linha mudou (começou, foi finalizada ou não existe).
	UpdateTimes(ctx context.Context, id string, startsAt, endsAt time.Time, reason *string, now time.Time) (bool, error)
	// UpdateReason troca só o motivo de uma janela viva. false = nenhuma linha mudou.
	UpdateReason(ctx context.Context, id string, reason *string) (bool, error)
	// Delete apaga uma janela que ainda não começou.
	Delete(ctx context.Context, id string) error
	// DueStart: janelas cujo início chegou e cujos efeitos de início não rodaram.
	DueStart(ctx context.Context, now time.Time) ([]*MaintenanceWindow, error)
	// DueEnd: janelas cujo fim chegou e que ainda não foram finalizadas.
	DueEnd(ctx context.Context, now time.Time) ([]*MaintenanceWindow, error)
	// StartGroup aplica o início de várias janelas numa transação só: para cada
	// uma, fecha o incidente aberto do monitor (motivo manutenção), põe o estado
	// em UNKNOWN e marca started_at. Só aplica se starts_at <= now ainda valer
	// (uma janela adiada depois do DueStart fica de fora). Devolve as janelas em
	// que o efeito foi aplicado; as já aplicadas não voltam.
	//
	// build (nil = sem aviso) roda depois de aplicar, com as janelas aplicadas e
	// seus monitores (lidos na própria transação), e as linhas que devolver vão
	// para o outbox na mesma transação: se build falha, nada é aplicado.
	StartGroup(ctx context.Context, windows []*MaintenanceWindow, now time.Time, build MaintenanceNoticeBuilder) ([]*MaintenanceWindow, error)
	// FinalizeGroup aplica o fim de várias janelas numa transação só (e o início
	// antes, se não rodou): marca finalized_at, grava resumed_at = now, estado
	// UNKNOWN e recomeça o prazo do heartbeat. endsAt != nil encurta as janelas
	// (cancelamento). Devolve as janelas finalizadas agora. build recebe só as que
	// já tinham started_at antes da transação: janela que começa e termina com o
	// Pingou fora não avisa (o início nunca saiu).
	FinalizeGroup(ctx context.Context, windows []*MaintenanceWindow, now time.Time, endsAt *time.Time, build MaintenanceNoticeBuilder) ([]*MaintenanceWindow, error)
}

type MonitorFilter struct {
	Enabled *bool
	State   *MonitorState
	Type    *MonitorType
	Limit   int
	Offset  int
}

// MaintenanceNoticeBuilder monta as linhas do outbox de um grupo de janelas
// aplicadas. Roda dentro da transação: não pode ler o banco (MaxOpenConns=1).
type MaintenanceNoticeBuilder func(applied []*MaintenanceWindow, monitors []*Monitor) ([]OutboxEvent, error)
