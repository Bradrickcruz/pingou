package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

type RetentionWorker struct {
	checks   domain.CheckRepository
	settings *repository.SettingsRepo
	done     chan struct{} // fechado quando o loop termina (nil antes do Start)
}

func NewRetentionWorker(checks domain.CheckRepository, settings *repository.SettingsRepo) *RetentionWorker {
	return &RetentionWorker{checks: checks, settings: settings}
}

func (w *RetentionWorker) Start(ctx context.Context) {
	w.done = make(chan struct{})
	go func() {
		defer close(w.done)
		w.run(ctx)
	}()
}

// Wait espera o loop terminar depois que o ctx de Start for cancelado.
func (w *RetentionWorker) Wait(ctx context.Context) error {
	if w.done == nil {
		return nil
	}
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *RetentionWorker) run(ctx context.Context) {
	// roda imediatamente na inicialização
	w.purge(ctx)

	// depois a cada 1 hora
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("retention worker stopped")
			return
		case <-ticker.C:
			w.purge(ctx)
		}
	}
}

// purge apaga só os checks; incidentes são histórico e nunca expiram.
func (w *RetentionWorker) purge(ctx context.Context) {
	retentionStr, err := w.settings.Get(ctx, "retention_days")
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		slog.Error("retention worker: could not read settings", "err", err)
		return
	}

	days := 30
	fmt.Sscanf(retentionStr, "%d", &days)

	before := time.Now().UTC().
		AddDate(0, 0, -days).
		Format(time.RFC3339)

	deleted, err := w.checks.DeleteOlderThan(ctx, before)
	if err != nil && ctx.Err() != nil {
		return // desligando: para entre lotes; o que sobrou sai na próxima vez
	}
	if err != nil {
		slog.Error("retention worker: purge failed", "err", err)
		return
	}

	if deleted > 0 {
		slog.Info("retention worker: purged old checks",
			"deleted", deleted,
			"retention_days", days,
			"before", before,
		)
	} else {
		slog.Debug("retention worker: nothing to purge", "retention_days", days)
	}
}
