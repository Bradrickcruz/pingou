package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type StateMachine struct {
	monitors  domain.MonitorRepository
	checks    domain.CheckRepository
	incidents domain.IncidentRepository
	notifier  Notifier
}

func NewStateMachine(
	monitors domain.MonitorRepository,
	checks domain.CheckRepository,
	incidents domain.IncidentRepository,
	notifier Notifier,
) *StateMachine {
	return &StateMachine{
		monitors:  monitors,
		checks:    checks,
		incidents: incidents,
		notifier:  notifier,
	}
}

// Process recebe o resultado de um check e aplica as transições de estado
func (sm *StateMachine) Process(ctx context.Context, m *domain.Monitor, result domain.CheckResult) error {
	now := time.Now().UTC()

	// 1. persiste o check
	check := &domain.Check{
		MonitorID:    m.ID,
		Success:      result.Success,
		StatusCode:   result.StatusCode,
		LatencyMs:    result.LatencyMs,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    now,
	}
	if err := sm.checks.Create(ctx, check); err != nil {
		return err
	}

	// 2. aplica transição de estado
	prevState := m.CurrentState

	if result.Success {
		return sm.handleSuccess(ctx, m, prevState, now)
	}
	return sm.handleFailure(ctx, m, prevState, result, now)
}

func (sm *StateMachine) handleSuccess(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, now time.Time) error {
	m.CurrentState = domain.StateUp
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	if err := sm.monitors.Update(ctx, m); err != nil {
		return err
	}

	// DOWN→UP: fecha incidente e notifica
	if prevState == domain.StateDown {
		incident, err := sm.incidents.FindOpenByMonitor(ctx, m.ID)
		if err != nil {
			return err
		}
		if incident != nil {
			endedAt := now.Format(time.RFC3339)
			durationSeconds := int(now.Sub(incident.StartedAt).Seconds())
			if durationSeconds < 0 {
				durationSeconds = 0
			}
			if err := sm.incidents.Close(ctx, incident.ID, endedAt, durationSeconds); err != nil {
				return err
			}
			d := int64(durationSeconds)
			incident.DurationSeconds = &d
			slog.Info("monitor recovered", "monitor_id", m.ID, "name", m.Name)
			sm.notifier.NotifyRecovery(ctx, m, incident)
		}
	}
	// UNKNOWN→UP: não notifica (D na PRD)

	return nil
}

func (sm *StateMachine) handleFailure(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, result domain.CheckResult, now time.Time) error {
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	// conta falhas consecutivas recentes
	consecutiveFails, err := sm.countConsecutiveFails(ctx, m)
	if err != nil {
		return err
	}

	// threshold ainda não atingido — mantém estado atual
	if consecutiveFails < m.FailureThreshold {
		return sm.monitors.Update(ctx, m)
	}

	// threshold atingido — vai pra DOWN
	m.CurrentState = domain.StateDown
	if err := sm.monitors.Update(ctx, m); err != nil {
		return err
	}

	// só abre incidente e notifica na transição →DOWN (não repete se já estava DOWN)
	if prevState != domain.StateDown {
		id, err := newUUIDv7()
		if err != nil {
			return err
		}
		incident := &domain.Incident{
			ID:        id,
			MonitorID: m.ID,
			StartedAt: now,
			LastError: result.ErrorMessage,
		}
		if err := sm.incidents.Create(ctx, incident); err != nil {
			return err
		}
		slog.Warn("monitor down", "monitor_id", m.ID, "name", m.Name, "error", result.ErrorMessage)
		sm.notifier.NotifyDown(ctx, m, incident)
	}

	return nil
}

// countConsecutiveFails conta checks recentes até encontrar um sucesso
func (sm *StateMachine) countConsecutiveFails(ctx context.Context, m *domain.Monitor) (int, error) {
	// busca os últimos N checks (N = failure_threshold + 1 para ter margem)
	checks, _, err := sm.checks.FindByMonitor(ctx, m.ID, m.FailureThreshold+1, 0)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, c := range checks {
		if !c.Success {
			count++
		} else {
			break // encontrou sucesso, para de contar
		}
	}
	return count, nil
}
