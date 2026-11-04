package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
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

// applyMaintenance aplica os inícios e depois os fins vencidos. A ordem importa
// para a janela que começou e terminou com o Pingou fora (Finalize também
// aplica o início que faltar).
func (s *Scheduler) applyMaintenance(ctx context.Context) {
	now := s.now().UTC()
	// sem cancelamento: a gravação que começou termina inteira
	writeCtx := context.WithoutCancel(ctx)

	starts, err := s.maintenance.DueStart(ctx, now)
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("maintenance tick error", "err", err)
		}
		return
	}
	for _, w := range starts {
		if ctx.Err() != nil {
			return
		}
		s.dbMu.Lock()
		ok, err := s.maintenance.Start(writeCtx, w, now)
		s.dbMu.Unlock()
		if err != nil {
			slog.Error("maintenance start error", "monitor_id", w.MonitorID, "window_id", w.ID, "err", err)
			continue
		}
		if ok {
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
	for _, w := range ends {
		if ctx.Err() != nil {
			return
		}
		if err := s.endWindow(writeCtx, w, now, nil); err != nil {
			slog.Error("maintenance end error", "monitor_id", w.MonitorID, "window_id", w.ID, "err", err)
		}
	}
}

// EndMaintenance encerra uma janela em curso na hora (cancelamento): ends_at
// passa a ser endsAt e o monitor volta como no fim normal.
func (s *Scheduler) EndMaintenance(ctx context.Context, w *domain.MaintenanceWindow, endsAt time.Time) error {
	return s.endWindow(ctx, w, s.now().UTC(), &endsAt)
}

// endWindow finaliza a janela sob dbMu e, depois do commit, recarrega o job:
// o primeiro check roda na hora. Reload respeita enabled=false, então um monitor
// pausado durante a janela continua pausado (D3).
func (s *Scheduler) endWindow(ctx context.Context, w *domain.MaintenanceWindow, now time.Time, endsAt *time.Time) error {
	s.dbMu.Lock()
	ok, err := s.maintenance.Finalize(ctx, w, now, endsAt)
	s.dbMu.Unlock()
	if err != nil || !ok {
		return err
	}
	slog.Info("maintenance ended", "monitor_id", w.MonitorID, "window_id", w.ID)
	return s.Reload(ctx, w.MonitorID)
}
