package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bradrickcruz/pingou/internal/checker"
	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/handler"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/scheduler"
	"github.com/Bradrickcruz/pingou/internal/service"
	"github.com/joho/godotenv"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Warn("failed to load .env", "err", err)
	}

	slog.Info("Pingou starting", "version", version, "commit", commit, "buildDate", buildDate)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("database error", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// repositories
	monitorRepo := repository.NewMonitorRepo(db)
	checkRepo := repository.NewCheckRepo(db)
	incidentRepo := repository.NewIncidentRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)

	// notifier — lê webhook_url do banco em runtime
	notifier := service.NewWebhookNotifier(func() string {
		url, _ := settingsRepo.Get(context.Background(), "webhook_url")
		return url
	})

	// state machine
	stateMachine := service.NewStateMachine(monitorRepo, checkRepo, incidentRepo, notifier)

	// checker
	httpChecker := checker.NewHTTPChecker()

	// services
	monitorSvc := service.NewMonitorService(monitorRepo, incidentRepo)

	// scheduler
	sched := scheduler.NewScheduler(monitorRepo, httpChecker, stateMachine)
	monitorSvc.SetReloader(sched)

	// inicia scheduler
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := sched.Start(ctx); err != nil {
		slog.Error("scheduler error", "err", err)
		os.Exit(1)
	}

	// server
	srv := handler.NewServer(cfg, monitorSvc)

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	// aguarda sinal de encerramento
	<-ctx.Done()
	slog.Info("shutting down...")
	sched.Stop()
}
