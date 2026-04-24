package domain

import "context"

type MonitorRepository interface {
	Create(ctx context.Context, m *Monitor) error
	FindAll(ctx context.Context, filter MonitorFilter) ([]*Monitor, int, error)
	FindByID(ctx context.Context, id string) (*Monitor, error)
	Update(ctx context.Context, m *Monitor) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type CheckRepository interface {
	Create(ctx context.Context, c *Check) error
	FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*Check, int, error)
	DeleteOlderThan(ctx context.Context, before string) (int64, error) // <- adiciona
}

type IncidentRepository interface {
	Create(ctx context.Context, i *Incident) error
	FindOpenByMonitor(ctx context.Context, monitorID string) (*Incident, error)
	FindByMonitor(ctx context.Context, monitorID string, onlyOpen bool, limit, offset int) ([]*Incident, int, error)
	FindAll(ctx context.Context, onlyOpen bool, limit, offset int) ([]*Incident, int, error)
	Close(ctx context.Context, id string, endedAt string, durationSeconds int) error
}

type MonitorFilter struct {
	Enabled *bool
	State   *MonitorState
	Limit   int
	Offset  int
}
