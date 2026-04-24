package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	APIKey      string
	DatabaseURL string
	Port        string
	LogLevel    slog.Level
}

func Load() (*Config, error) {
	apiKey := os.Getenv("PINGOU_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PINGOU_API_KEY is required")
	}

	dbURL := os.Getenv("PINGOU_DATABASE_URL")
	if dbURL == "" {
		dbURL = "pingou.db"
	}

	port := os.Getenv("PINGOU_PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := slog.LevelInfo
	if os.Getenv("PINGOU_LOG_LEVEL") == "DEBUG" {
		logLevel = slog.LevelDebug
	}

	return &Config{
		APIKey:      apiKey,
		DatabaseURL: dbURL,
		Port:        port,
		LogLevel:    logLevel,
	}, nil
}
