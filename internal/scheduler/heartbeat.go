package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

const missedPingMessage = "no ping received within interval + grace"

// runHeartbeats varre os heartbeats ativos a cada heartbeatTick e marca DOWN
// os que passaram do prazo. Um ticker único, sem goroutine por monitor.
func (s *Scheduler) runHeartbeats(ctx context.Context) {
	ticker := time.NewTicker(s.heartbeatTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkHeartbeats(ctx, time.Now())
		}
	}
}

func (s *Scheduler) checkHeartbeats(ctx context.Context, now time.Time) {
	enabled := true
	typ := domain.TypeHeartbeat
	monitors, _, err := s.monitors.FindAll(ctx, domain.MonitorFilter{
		Enabled: &enabled,
		Type:    &typ,
		Limit:   100,
	})
	if err != nil {
		slog.Error("heartbeat evaluator error", "err", err)
		return
	}

	for _, m := range monitors {
		if !s.heartbeatExpired(m, now) {
			continue
		}
		// confere de novo dentro do lock: um ping pode ter chegado nesse meio tempo
		err := s.withMonitor(ctx, m.ID, func(cur *domain.Monitor) *domain.CheckResult {
			if !s.heartbeatExpired(cur, now) {
				return nil
			}
			msg := missedPingMessage
			return &domain.CheckResult{Success: false, ErrorMessage: &msg}
		})
		if err != nil {
			slog.Error("heartbeat evaluator error", "monitor_id", m.ID, "err", err)
		}
	}
}

// heartbeatExpired: venceu interval + grace desde o último ping. Gera no máximo
// uma falha por queda (já DOWN não gera de novo).
func (s *Scheduler) heartbeatExpired(m *domain.Monitor, now time.Time) bool {
	if m.CurrentState == domain.StateDown {
		return false
	}
	// sem nenhum ping, conta da criação ou da última edição/retomada (updated_at;
	// sem ping não há checks que o mexam): "esqueci o curl" também alerta
	base := m.UpdatedAt
	if m.LastPingAt != nil {
		base = *m.LastPingAt
	}
	window := time.Duration(m.IntervalSeconds+m.Config.GraceSeconds) * s.intervalUnit
	deadline := base.Add(window)
	// venceu com o Pingou fora do ar: o ping pode ter se perdido. Espera a
	// próxima execução da tarefa em vez de dar alarme falso.
	if deadline.Before(s.startedAt) {
		deadline = s.startedAt.Add(window)
	}
	return now.After(deadline)
}
