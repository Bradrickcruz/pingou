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
	FindAll(ctx context.Context, filter IncidentFilter) ([]*Incident, int, error)
	Close(ctx context.Context, id string, endedAt string, durationSeconds int) error
	DeleteResolvedOlderThan(ctx context.Context, before string) (int64, error)
	FindInWindow(ctx context.Context, monitorID string, since time.Time) ([]*Incident, error)
}

type MonitorFilter struct {
	Enabled *bool
	State   *MonitorState
	Type    *MonitorType
	Limit   int
	Offset  int
}
