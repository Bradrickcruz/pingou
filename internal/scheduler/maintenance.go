package scheduler

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// inMaintenance diz se o monitor está numa janela de manutenção agora. Erro de
// leitura não bloqueia o check (o mesmo erro aparece na gravação).
func (s *Scheduler) inMaintenance(ctx context.Context, monitorID string) bool {
	if s.maintenance == nil {
		return false
	}
	active, err := s.maintenance.ActiveForMonitor(ctx, monitorID, s.now())
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("maintenance gate error", "monitor_id", monitorID, "err", err)
		}
		return false
	}
	return active
}

// runMaintenance aplica o início e o fim das janelas a cada maintenanceTick.
// Um ticker único, no molde do runHeartbeats. O bloqueio em si não depende do
// tick (vem do relógio, no gate); só os efeitos laterais podem atrasar um tick.
func (s *Scheduler) runMaintenance(ctx context.Context) {
	ticker := time.NewTicker(s.maintenanceTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.applyMaintenance(ctx)
		}
	}
}

// applyMaintenance aplica os inícios e depois os fins vencidos, em grupos de
// janelas com mesmo início, fim e motivo (um aviso por grupo). A ordem importa
// para a janela que começou e terminou com o Pingou fora (o fim também aplica o
// início que faltar).
func (s *Scheduler) applyMaintenance(ctx context.Context) {
	now := s.now().UTC()
	// sem cancelamento: a gravação que começou termina inteira
	writeCtx := context.WithoutCancel(ctx)
	notices := s.noticesReader(writeCtx)

	starts, err := s.maintenance.DueStart(ctx, now)
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("maintenance tick error", "err", err)
		}
		return
	}
	// janela que já passou do fim vai direto para o fim, que aplica o início
	// que faltou: sem aviso, porque o início nunca saiu
	starts = slices.DeleteFunc(starts, func(w *domain.MaintenanceWindow) bool { return !w.EndsAt.After(now) })
	for _, group := range groupWindows(starts) {
		if ctx.Err() != nil {
			return
		}
		enabled, err := notices()
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("maintenance notices setting error", "err", err)
			}
			return // tenta no próximo tick; o bloqueio dos checks vem do relógio
		}
		applied, err := s.startGroup(writeCtx, group, now, enabled)
		if err != nil {
			slog.Error("maintenance start error", "window_ids", windowIDs(group), "err", err)
			continue
		}
		for _, w := range applied {
			slog.Info("maintenance started", "monitor_id", w.MonitorID, "window_id", w.ID)
		}
	}

	ends, err := s.maintenance.DueEnd(ctx, now)
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("maintenance tick error", "err", err)
		}
		return
	}
	for _, group := range groupWindows(ends) {
		if ctx.Err() != nil {
			return
		}
		enabled, err := notices()
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("maintenance notices setting error", "err", err)
			}
			return
		}
		if err := s.endGroup(writeCtx, group, now, nil, enabled); err != nil {
			slog.Error("maintenance end error", "window_ids", windowIDs(group), "err", err)
		}
	}
}

// noticesReader lê o interruptor no máximo uma vez por tick, e só se houver
// janela vencida: sem janela, o tick não faz consulta extra.
func (s *Scheduler) noticesReader(ctx context.Context) func() (bool, error) {
	var done, enabled bool
	return func() (bool, error) {
		if done {
			return enabled, nil
		}
		if s.maintenanceNotices == nil {
			return false, nil
		}
		v, err := s.maintenanceNotices(ctx)
		if err != nil {
			return false, err
		}
		done, enabled = true, v
		return v, nil
	}
}

// groupWindows separa as janelas por MaintenanceGroupKey, na ordem em que aparecem.
func groupWindows(ws []*domain.MaintenanceWindow) [][]*domain.MaintenanceWindow {
	var groups [][]*domain.MaintenanceWindow
	index := map[string]int{}
	for _, w := range ws {
		k := domain.MaintenanceGroupKey(w)
		i, ok := index[k]
		if !ok {
			i = len(groups)
			index[k] = i
			groups = append(groups, nil)
		}
		groups[i] = append(groups[i], w)
	}
	return groups
}

func windowIDs(ws []*domain.MaintenanceWindow) []string {
	ids := make([]string, len(ws))
	for i, w := range ws {
		ids[i] = w.ID
	}
	return ids
}

// noticeBuilder devolve o montador do aviso, ou nil com o interruptor desligado.
func (s *Scheduler) noticeBuilder(eventType string, enabled bool, now time.Time) domain.MaintenanceNoticeBuilder {
	if !enabled {
		return nil
	}
	return func(applied []*domain.MaintenanceWindow, monitors []*domain.Monitor) ([]domain.OutboxEvent, error) {
		return service.NewMaintenanceEvents(eventType, applied, monitors, s.channels, now)
	}
}

// startGroup aplica o início do grupo sob dbMu; o aviso entra na mesma transação.
// Depois do commit, acorda o worker do outbox.
func (s *Scheduler) startGroup(ctx context.Context, group []*domain.MaintenanceWindow, now time.Time, notices bool) ([]*domain.MaintenanceWindow, error) {
	s.dbMu.Lock()
	applied, err := s.maintenance.StartGroup(ctx, group, now, s.noticeBuilder(domain.EventMaintenanceStarted, notices, now))
	s.dbMu.Unlock()
	if err == nil && notices && len(applied) > 0 {
		s.wakeWorker()
	}
	return applied, err
}

// endGroup aplica o fim do grupo sob dbMu e, depois do commit, acorda o worker
// e recarrega o job de cada monitor: o primeiro check roda na hora. Reload
// respeita enabled=false, então um monitor pausado durante a janela continua
// pausado (D3).
func (s *Scheduler) endGroup(ctx context.Context, group []*domain.MaintenanceWindow, now time.Time, endsAt *time.Time, notices bool) error {
	s.dbMu.Lock()
	applied, err := s.maintenance.FinalizeGroup(ctx, group, now, endsAt, s.noticeBuilder(domain.EventMaintenanceEnded, notices, now))
	s.dbMu.Unlock()
	if err != nil {
		return err
	}
	if notices && len(applied) > 0 {
		s.wakeWorker()
	}
	var firstErr error
	for _, w := range applied {
		slog.Info("maintenance ended", "monitor_id", w.MonitorID, "window_id", w.ID)
		if err := s.Reload(ctx, w.MonitorID); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Scheduler) wakeWorker() {
	if s.wake != nil {
		s.wake()
	}
}

// EndMaintenance encerra uma janela em curso na hora (cancelamento): ends_at
// passa a ser endsAt e o monitor volta como no fim normal. Vira um grupo de uma
// janela, com aviso de fim se o interruptor estiver ligado. Erro ao ler o
// interruptor não impede o cancelamento: segue sem aviso e registra no log.
func (s *Scheduler) EndMaintenance(ctx context.Context, w *domain.MaintenanceWindow, endsAt time.Time) error {
	notices := false
	if s.maintenanceNotices != nil {
		v, err := s.maintenanceNotices(ctx)
		if err != nil {
			slog.Error("maintenance notices setting error, ending without notice", "window_id", w.ID, "err", err)
		}
		notices = err == nil && v
	}
	return s.endGroup(ctx, []*domain.MaintenanceWindow{w}, s.now().UTC(), &endsAt, notices)
}
