package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/google/uuid"
)

type MonitorService struct {
	monitors  domain.MonitorRepository
	incidents domain.IncidentRepository
}

func NewMonitorService(monitors domain.MonitorRepository, incidents domain.IncidentRepository) *MonitorService {
	return &MonitorService{monitors: monitors, incidents: incidents}
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
	count, err := s.monitors.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count >= maxMonitors {
		return nil, fmt.Errorf("%w", ErrLimitReached)
	}

	now := time.Now().UTC()
	m := &domain.Monitor{
		ID:               uuid.New().String(),
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

	if err := validate(m); err != nil {
		return nil, err
	}

	if err := s.monitors.Create(ctx, m); err != nil {
		return nil, err
	}
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

	if err := validate(m); err != nil {
		return nil, err
	}

	if err := s.monitors.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MonitorService) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.monitors.Delete(ctx, id)
}
