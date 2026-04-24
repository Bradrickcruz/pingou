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
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Warn("failed to load .env", "err", err)
	}

	slog.Info("Pingou starting", "version", version, "commit", commit, "buildDate", buildDate)

	loadedConfig, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loadedConfig.LogLevel,
	}))
	slog.SetDefault(logger)

	db, err := database.Open(loadedConfig.DatabaseURL)
	if err != nil {
		slog.Error("database error", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// repositories
	monitorRepository := repository.NewMonitorRepo(db)
	checkRepository := repository.NewCheckRepo(db)
	incidentRepository := repository.NewIncidentRepo(db)
	settingsRepository := repository.NewSettingsRepo(db)

	// webhookNotifier — lê webhook_url do banco em runtime
	webhookNotifier := service.NewWebhookNotifier(func() string {
		url, _ := settingsRepository.Get(context.Background(), "webhook_url")
		return url
	})

	// state machine
	stateMachine := service.NewStateMachine(monitorRepository, checkRepository, incidentRepository, webhookNotifier)

	// checker
	httpChecker := checker.NewHTTPChecker()

	// services
	monitorService := service.NewMonitorService(monitorRepository, checkRepository, incidentRepository)
	incidentService := service.NewIncidentService(incidentRepository)
	settingsService := service.NewSettingsService(settingsRepository)

	// scheduler
	monitorScheduler := scheduler.NewScheduler(monitorRepository, httpChecker, stateMachine)
	monitorService.SetReloader(monitorScheduler)

	// inicia scheduler
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := monitorScheduler.Start(ctx); err != nil {
		slog.Error("scheduler error", "err", err)
		os.Exit(1)
	}

	// retention worker
	retention := scheduler.NewRetentionWorker(checkRepository, settingsRepository)
	retention.Start(ctx)

	// inicia server
	server := handler.NewServer(loadedConfig, monitorService, incidentService, settingsService)

	go func() {
		if err := server.Start(); err != nil {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	// aguarda sinal de encerramento
	<-ctx.Done()
	slog.Info("shutting down...")
	monitorScheduler.Stop()
}
