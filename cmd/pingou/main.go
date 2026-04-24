package main

import (
	"log/slog"
	"os"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/handler"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
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
	incidentRepo := repository.NewIncidentRepo(db)

	// services
	monitorSvc := service.NewMonitorService(monitorRepo, incidentRepo)

	// server
	srv := handler.NewServer(cfg, monitorSvc)
	if err := srv.Start(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
