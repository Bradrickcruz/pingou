package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// Reloader é implementado pelo Scheduler — evita import circular
type Reloader interface {
	Reload(ctx context.Context, monitorID string) error
}

type MonitorService struct {
	monitors  domain.MonitorRepository
	checks    domain.CheckRepository
	incidents domain.IncidentRepository
	reloader  Reloader
}

func NewMonitorService(
	monitors domain.MonitorRepository,
	checks domain.CheckRepository, // <- adiciona
	incidents domain.IncidentRepository,
) *MonitorService {
	return &MonitorService{monitors: monitors, checks: checks, incidents: incidents}
}

func (s *MonitorService) SetReloader(r Reloader) {
	s.reloader = r
}

const maxMonitors = 100

type CreateMonitorInput struct {
	Name             string
	URL              string
	IntervalSeconds  int
	TimeoutSeconds   int
	FailureThreshold int
	Enabled          bool
}

type UpdateMonitorInput struct {
	Name             *string
	URL              *string
	IntervalSeconds  *int
	TimeoutSeconds   *int
	FailureThreshold *int
	Enabled          *bool
}

func (s *MonitorService) Create(ctx context.Context, in CreateMonitorInput) (*domain.Monitor, error) {
	if err := validateCreateInput(in); err != nil {
		return nil, err
	}

	count, err := s.monitors.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count >= maxMonitors {
		return nil, fmt.Errorf("%w", ErrLimitReached)
	}

	now := time.Now().UTC()
	id, err := newUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("generate monitor id: %w", err)
	}
	m := &domain.Monitor{
		ID:               id,
		Name:             in.Name,
		URL:              in.URL,
		IntervalSeconds:  in.IntervalSeconds,
		TimeoutSeconds:   in.TimeoutSeconds,
		FailureThreshold: in.FailureThreshold,
		Enabled:          in.Enabled,
		CurrentState:     domain.StateUnknown,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.monitors.Create(ctx, m); err != nil {
		return nil, err
	}

	s.reload(ctx, m.ID)
	return m, nil
}

func (s *MonitorService) GetByID(ctx context.Context, id string) (*domain.Monitor, error) {
	m, err := s.monitors.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MonitorService) List(ctx context.Context, f domain.MonitorFilter) ([]*domain.Monitor, int, error) {
	return s.monitors.FindAll(ctx, f)
}

func (s *MonitorService) Update(ctx context.Context, id string, in UpdateMonitorInput) (*domain.Monitor, error) {
	if err := validateUpdateInput(in); err != nil {
		return nil, err
	}

	m, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		m.Name = *in.Name
	}
	if in.URL != nil {
		m.URL = *in.URL
	}
	if in.IntervalSeconds != nil {
		m.IntervalSeconds = *in.IntervalSeconds
	}
	if in.TimeoutSeconds != nil {
		m.TimeoutSeconds = *in.TimeoutSeconds
	}
	if in.FailureThreshold != nil {
		m.FailureThreshold = *in.FailureThreshold
	}
	if in.Enabled != nil {
		m.Enabled = *in.Enabled
	}

	m.UpdatedAt = time.Now().UTC()

	if err := s.monitors.Update(ctx, m); err != nil {
		return nil, err
	}

	s.reload(ctx, m.ID)
	return m, nil
}

func (s *MonitorService) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.monitors.Delete(ctx, id); err != nil {
		return err
	}
	s.reload(ctx, id)
	return nil
}

func (s *MonitorService) reload(ctx context.Context, id string) {
	if s.reloader != nil {
		if err := s.reloader.Reload(ctx, id); err != nil {
			slog.Error("scheduler reload error", "monitor_id", id, "err", err)
		}
	}
}

func (s *MonitorService) ListChecks(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error) {
	if _, err := s.GetByID(ctx, monitorID); err != nil {
		return nil, 0, err
	}
	return s.checks.FindByMonitor(ctx, monitorID, limit, offset)
}
