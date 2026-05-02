package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIKey             string
	DatabaseURL        string
	Port               string
	LogLevel           slog.Level
	CORSAllowedOrigins []string
	MaxRedirects       int
	GlobalTimeout      int
}

// Load carrega configuracao completa (requer API key)
func Load() (*Config, error) {
	apiKey := os.Getenv("PINGOU_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PINGOU_API_KEY is required")
	}

	return loadConfigWithAPIKey(apiKey)
}

// LoadConfig carrega configuracao minima para CLI (sem API key)
func LoadConfig() (*Config, error) {
	return loadConfigWithAPIKey("")
}

func loadConfigWithAPIKey(apiKey string) (*Config, error) {
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

	maxRedirects := 5
	if v := os.Getenv("PINGOU_MAX_REDIRECTS"); v != "" {
		if n := parseInt(v); n > 0 {
			maxRedirects = n
		}
	}

	globalTimeout := 60
	if v := os.Getenv("PINGOU_GLOBAL_TIMEOUT"); v != "" {
		if n := parseInt(v); n > 0 {
			globalTimeout = n
		}
	}

	return &Config{
		APIKey:             apiKey,
		DatabaseURL:        dbURL,
		Port:               port,
		LogLevel:           logLevel,
		CORSAllowedOrigins: corsOrigins,
		MaxRedirects:       maxRedirects,
		GlobalTimeout:      globalTimeout,
	}, nil
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
