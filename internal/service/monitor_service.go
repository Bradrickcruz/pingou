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

// Recorder grava um resultado de check pela state machine (Scheduler).
// pingedAt != nil indica um ping de heartbeat.
type Recorder interface {
	Record(ctx context.Context, monitorID string, result domain.CheckResult, pingedAt *time.Time) error
}

type MonitorService struct {
	monitors  domain.MonitorRepository
	checks    domain.CheckRepository
	incidents domain.IncidentRepository
	reloader  Reloader
	recorder  Recorder
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

func (s *MonitorService) SetRecorder(r Recorder) {
	s.recorder = r
}

const (
	maxMonitors = 100
	historySize = 30
	// heartbeat não usa timeout nem threshold; grava valores fixos válidos
	heartbeatTimeout   = 10
	heartbeatThreshold = 1
)

type CreateMonitorInput struct {
	Name             string
	Type             string // vazio = http
	URL              string
	Host             string
	Port             int
	GraceSeconds     *int // nil = padrão (60)
	IntervalSeconds  int
	TimeoutSeconds   int
	FailureThreshold int
	Enabled          bool
}

type UpdateMonitorInput struct {
	Name             *string
	Type             *string // só aceito se igual ao atual
	URL              *string
	Host             *string
	Port             *int
	GraceSeconds     *int
	IntervalSeconds  *int
	TimeoutSeconds   *int
	FailureThreshold *int
	Enabled          *bool
}

func (s *MonitorService) Create(ctx context.Context, in CreateMonitorInput) (*domain.Monitor, error) {
	now := time.Now().UTC()
	m := &domain.Monitor{
		Name:             in.Name,
		Type:             domain.MonitorType(in.Type),
		IntervalSeconds:  in.IntervalSeconds,
		TimeoutSeconds:   in.TimeoutSeconds,
		FailureThreshold: in.FailureThreshold,
		Enabled:          in.Enabled,
		CurrentState:     domain.StateUnknown,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if m.Type == "" {
		m.Type = domain.TypeHTTP
	}

	switch m.Type {
	case domain.TypeHTTP:
		m.URL = in.URL
	case domain.TypeTCP:
		m.Config.Host, m.Config.Port = in.Host, in.Port
	case domain.TypeHeartbeat:
		m.TimeoutSeconds, m.FailureThreshold = heartbeatTimeout, heartbeatThreshold
		m.Config.GraceSeconds = defaultGraceSeconds
		if in.GraceSeconds != nil {
			m.Config.GraceSeconds = *in.GraceSeconds
		}
	}

	if err := validateMonitor(m); err != nil {
		return nil, err
	}

	count, err := s.monitors.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count >= maxMonitors {
		return nil, fmt.Errorf("%w", ErrLimitReached)
	}

	id, err := newUUIDv7()
	if err != nil {
		return nil, fmt.Errorf("generate monitor id: %w", err)
	}
	m.ID = id

	if m.Type == domain.TypeHeartbeat {
		token, err := newPingToken()
		if err != nil {
			return nil, err
		}
		m.PingToken = &token
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
	m, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Type != nil && domain.MonitorType(*in.Type) != m.Type {
		return nil, fmt.Errorf("%w: type cannot be changed", ErrValidation)
	}

	wasEnabled := m.Enabled
	if in.Name != nil {
		m.Name = *in.Name
	}
	if in.IntervalSeconds != nil {
		m.IntervalSeconds = *in.IntervalSeconds
	}
	if in.Enabled != nil {
		m.Enabled = *in.Enabled
	}

	// campos de outros tipos são ignorados
	switch m.Type {
	case domain.TypeHTTP:
		if in.URL != nil {
			m.URL = *in.URL
		}
	case domain.TypeTCP:
		if in.Host != nil {
			m.Config.Host = *in.Host
		}
		if in.Port != nil {
			m.Config.Port = *in.Port
		}
	case domain.TypeHeartbeat:
		if in.GraceSeconds != nil {
			m.Config.GraceSeconds = *in.GraceSeconds
		}
	}
	if m.Type != domain.TypeHeartbeat {
		if in.TimeoutSeconds != nil {
			m.TimeoutSeconds = *in.TimeoutSeconds
		}
		if in.FailureThreshold != nil {
			m.FailureThreshold = *in.FailureThreshold
		}
	}

	if err := validateMonitor(m); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	// Heartbeat retomado: o ping recebido durante a pausa não foi registrado,
	// então o prazo recomeça agora (senão cairia logo após retomar).
	if m.Type == domain.TypeHeartbeat && !wasEnabled && m.Enabled && m.LastPingAt != nil {
		m.LastPingAt = &now
	}
	m.UpdatedAt = now

	if err := s.monitors.Update(ctx, m); err != nil {
		return nil, err
	}

	s.reload(ctx, m.ID)
	return m, nil
}

// RegenerateToken troca o token de ping de um heartbeat; o antigo deixa de valer na hora.
func (s *MonitorService) RegenerateToken(ctx context.Context, id string) (*domain.Monitor, error) {
	m, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m.Type != domain.TypeHeartbeat {
		return nil, fmt.Errorf("%w: only heartbeat monitors have a ping token", ErrValidation)
	}

	token, err := newPingToken()
	if err != nil {
		return nil, err
	}
	m.PingToken = &token
	m.UpdatedAt = time.Now().UTC()

	if err := s.monitors.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Ping registra o ping de um heartbeat. Token desconhecido → ErrNotFound.
// Monitor pausado aceita o ping, mas nada é registrado.
func (s *MonitorService) Ping(ctx context.Context, token string) error {
	m, err := s.monitors.FindByPingToken(ctx, token)
	if err != nil {
		return err
	}
	if m == nil || m.Type != domain.TypeHeartbeat {
		return ErrNotFound
	}
	if !m.Enabled || s.recorder == nil {
		return nil
	}
	now := time.Now().UTC()
	return s.recorder.Record(ctx, m.ID, domain.CheckResult{Success: true}, &now)
}

// Summaries devolve histórico, uptime 24 h e última latência dos monitores.
func (s *MonitorService) Summaries(ctx context.Context, monitors []*domain.Monitor) (map[string]*domain.MonitorSummary, error) {
	ids := make([]string, len(monitors))
	for i, m := range monitors {
		ids[i] = m.ID
	}
	return s.checks.Summaries(ctx, ids, historySize, time.Now().Add(-24*time.Hour))
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
