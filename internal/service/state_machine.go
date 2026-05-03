package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type StateMachine struct {
	uow      UnitOfWork
	notifier Notifier
}

func NewStateMachine(uow UnitOfWork, notifier Notifier) *StateMachine {
	return &StateMachine{
		uow:      uow,
		notifier: notifier,
	}
}

// Process recebe o resultado de um check e aplica as transições de estado.
// Todas as operações de persistência são executadas dentro de uma transação.
func (sm *StateMachine) Process(ctx context.Context, m *domain.Monitor, result domain.CheckResult) error {
	now := time.Now().UTC()

	// Inicia transação
	if err := sm.uow.Begin(ctx); err != nil {
		return err
	}

	// Garante rollback em caso de erro
	success := false
	defer func() {
		if !success {
			_ = sm.uow.Rollback()
		}
	}()

	// 1. persiste o check dentro da transação
	check := &domain.Check{
		MonitorID:    m.ID,
		Success:      result.Success,
		StatusCode:   result.StatusCode,
		LatencyMs:    result.LatencyMs,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    now,
	}
	if err := sm.uow.CheckRepo().CreateWithTx(ctx, nil, check); err != nil {
		return err
	}

	// 2. aplica transição de estado
	prevState := m.CurrentState

	var notificationErr error
	if result.Success {
		notificationErr = sm.handleSuccess(ctx, m, prevState, now)
	} else {
		notificationErr = sm.handleFailure(ctx, m, prevState, result, now)
	}

	// Commit da transação antes de notificar
	if err := sm.uow.Commit(); err != nil {
		return err
	}

	success = true

	// Notificações APÓS commit (fora da transação)
	if notificationErr != nil {
		return notificationErr
	}

	return nil
}

func (sm *StateMachine) handleSuccess(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, now time.Time) error {
	m.CurrentState = domain.StateUp
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	if err := sm.uow.MonitorRepo().UpdateWithTx(ctx, nil, m); err != nil {
		return err
	}

	// DOWN→UP: fecha incidente (dentro da transação)
	if prevState == domain.StateDown {
		incident, err := sm.uow.IncidentRepo().FindOpenByMonitor(ctx, m.ID)
		if err != nil {
			return err
		}
		if incident != nil {
			endedAt := now.Format(time.RFC3339)
			durationSeconds := int(now.Sub(incident.StartedAt).Seconds())
			if durationSeconds < 0 {
				durationSeconds = 0
			}
			if err := sm.uow.IncidentRepo().Close(ctx, incident.ID, endedAt, durationSeconds); err != nil {
				return err
			}
			d := int64(durationSeconds)
			incident.DurationSeconds = &d
			slog.Info("monitor recovered", "monitor_id", m.ID, "name", m.Name)
			// Notificação fora da transação (após commit)
			sm.notifier.NotifyRecovery(ctx, m, incident)
		}
	}
	// UNKNOWN→UP: não notifica (D na PRD)

	return nil
}

func (sm *StateMachine) handleFailure(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, result domain.CheckResult, now time.Time) error {
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	// conta falhas consecutivas recentes (lê do DB principal, não da transação)
	consecutiveFails, err := sm.countConsecutiveFails(ctx, m)
	if err != nil {
		return err
	}

	// threshold ainda não atingido — mantém estado atual
	if consecutiveFails < m.FailureThreshold {
		return sm.uow.MonitorRepo().UpdateWithTx(ctx, nil, m)
	}

	// threshold atingido — vai pra DOWN
	m.CurrentState = domain.StateDown
	if err := sm.uow.MonitorRepo().UpdateWithTx(ctx, nil, m); err != nil {
		return err
	}

	// só abre incidente na transição →DOWN (não repete se já estava DOWN)
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
		if err := sm.uow.IncidentRepo().CreateWithTx(ctx, nil, incident); err != nil {
			return err
		}
		slog.Warn("monitor down", "monitor_id", m.ID, "name", m.Name, "error", result.ErrorMessage)
		// Notificação fora da transação (após commit)
		sm.notifier.NotifyDown(ctx, m, incident)
	}

	return nil
}

// countConsecutiveFails conta checks recentes até encontrar um sucesso
// Lê do DB principal (fora da transação) para evitar locking
func (sm *StateMachine) countConsecutiveFails(ctx context.Context, m *domain.Monitor) (int, error) {
	// busca os últimos N checks (N = failure_threshold + 1 para ter margem)
	checks, _, err := sm.uow.CheckRepo().FindByMonitor(ctx, m.ID, m.FailureThreshold+1, 0)
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