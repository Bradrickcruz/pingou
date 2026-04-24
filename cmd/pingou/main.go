package main

import (
	"log/slog"
	"os"
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

	slog.Info("Pingou starting",
		"version", version,
		"commit", commit,
		"buildDate", buildDate,
	)
}
