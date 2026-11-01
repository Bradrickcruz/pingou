package config

import (
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"strconv"
	"strings"
)

// Port é fixa: o Pingou sempre escuta nela (mapeie outra porta no host, se precisar).
const Port = "8956"

type Config struct {
	APIKey             string
	DatabaseURL        string
	Port               string
	LogLevel           slog.Level
	CORSAllowedOrigins []string
	MaxRedirects       int
	GlobalTimeout      int
	PublicURL          string // endereço público do dashboard, para links no webhook
	SMTP               SMTPConfig

	// Dados do build (-ldflags), preenchidos pelo comando serve. Expostos em /api/info.
	Version   string
	Commit    string
	BuildDate string
}

// Load carrega configuracao completa (requer API key)
func Load() (*Config, error) {
	apiKey := os.Getenv("PINGOU_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PINGOU_API_KEY is required")
	}

	cfg, err := loadConfigWithAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	// SMTP só no serve: um erro aqui não deve travar comandos como o migrate.
	if cfg.SMTP, err = LoadSMTP(); err != nil {
		return nil, err
	}
	return cfg, nil
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

	publicURL := strings.TrimRight(os.Getenv("PINGOU_PUBLIC_URL"), "/")

	return &Config{
		APIKey:             apiKey,
		DatabaseURL:        dbURL,
		Port:               Port,
		LogLevel:           logLevel,
		CORSAllowedOrigins: corsOrigins,
		MaxRedirects:       maxRedirects,
		GlobalTimeout:      globalTimeout,
		PublicURL:          publicURL,
	}, nil
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// Modos de TLS do SMTP.
const (
	SMTPTLSStartTLS = "starttls" // conexão aberta e STARTTLS obrigatório (587)
	SMTPTLSImplicit = "tls"      // TLS desde o início (465)
	SMTPTLSNone     = "none"     // sem TLS: só relay local, sem autenticação
)

// SMTPConfig vem toda do ambiente (PINGOU_SMTP_*). Nada disso vai para o banco
// nem para o export. Password nunca aparece em log, API ou `pingou config`.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      string
}

// Enabled diz se o e-mail está disponível (PINGOU_SMTP_HOST setado).
func (c SMTPConfig) Enabled() bool { return c.Host != "" }

// LoadSMTP lê PINGOU_SMTP_*. Sem HOST, devolve config vazia e ignora o resto.
// Com HOST e alguma regra violada, devolve erro: o serve não sobe, como quando
// falta PINGOU_API_KEY, para o usuário não achar que configurou e não receber nada.
func LoadSMTP() (SMTPConfig, error) {
	host := strings.TrimSpace(os.Getenv("PINGOU_SMTP_HOST"))
	if host == "" {
		return SMTPConfig{}, nil
	}

	c := SMTPConfig{
		Host:     host,
		Port:     587,
		Username: os.Getenv("PINGOU_SMTP_USERNAME"),
		Password: os.Getenv("PINGOU_SMTP_PASSWORD"),
		From:     strings.TrimSpace(os.Getenv("PINGOU_SMTP_FROM")),
		TLS:      strings.ToLower(strings.TrimSpace(os.Getenv("PINGOU_SMTP_TLS"))),
	}

	if v := strings.TrimSpace(os.Getenv("PINGOU_SMTP_PORT")); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil || port < 1 || port > 65535 {
			return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_PORT must be a number between 1 and 65535")
		}
		c.Port = port
	}

	if c.From == "" {
		return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_FROM is required when PINGOU_SMTP_HOST is set")
	}
	if _, err := mail.ParseAddress(c.From); err != nil {
		return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_FROM must be a valid address (ex.: Pingou <alerta@exemplo.com>)")
	}

	switch c.TLS {
	case "":
		c.TLS = SMTPTLSStartTLS
		if c.Port == 465 {
			c.TLS = SMTPTLSImplicit
		}
	case SMTPTLSStartTLS, SMTPTLSImplicit, SMTPTLSNone:
	default:
		return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_TLS must be starttls, tls or none")
	}

	if c.Username != "" && c.Password == "" {
		return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_PASSWORD is required when PINGOU_SMTP_USERNAME is set")
	}
	if c.TLS == SMTPTLSNone && c.Username != "" {
		return SMTPConfig{}, fmt.Errorf("PINGOU_SMTP_TLS=none cannot be used with PINGOU_SMTP_USERNAME (the password would travel in plain text)")
	}

	return c, nil
}
