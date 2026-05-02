package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	APIKey      string
	DatabaseURL string
	Port        string
	LogLevel    slog.Level
	CORSAllowedOrigins []string
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

	// CORS allowed origins (comma separated). Empty = no CORS headers.
	cors := os.Getenv("PINGOU_CORS_ALLOWED_ORIGINS")
	var corsOrigins []string
	if cors != "" {
		for _, o := range strings.Split(cors, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				corsOrigins = append(corsOrigins, o)
			}
		}
	}

	return &Config{
		APIKey:      apiKey,
		DatabaseURL: dbURL,
		Port:        port,
		LogLevel:    logLevel,
		CORSAllowedOrigins: corsOrigins,
	}, nil
}
