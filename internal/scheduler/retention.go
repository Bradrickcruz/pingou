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
}

func NewRetentionWorker(checks domain.CheckRepository, settings *repository.SettingsRepo) *RetentionWorker {
	return &RetentionWorker{checks: checks, settings: settings}
}

func (w *RetentionWorker) Start(ctx context.Context) {
	go w.run(ctx)
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

func (w *RetentionWorker) purge(ctx context.Context) {
	retentionStr, err := w.settings.Get(ctx, "retention_days")
	if err != nil {
		slog.Error("retention worker: could not read settings", "err", err)
		return
	}

	days := 30
	fmt.Sscanf(retentionStr, "%d", &days)

	before := time.Now().UTC().
		AddDate(0, 0, -days).
		Format(time.RFC3339)

	deleted, err := w.checks.DeleteOlderThan(ctx, before)
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
