package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra a versao do binario",
	Long:  "Exibe informacoes de versao, commit e data de build do binario",
	RunE:  runVersion,
}

var configCmd = &cobra.Command{
	Use:               "config",
	Short:             "Mostra a configuracao atual",
	Long:              "Exibe as configuracoes do Pingou (sem secrets)",
	RunE:              runConfig,
	PersistentPreRunE: requireKey,
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("pingou %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("build date: %s\n", date)
	return nil
}

func runConfig(cmd *cobra.Command, args []string) error {
	_ = godotenv.Load()

	cfg := map[string]interface{}{
		"DatabaseURL":        getEnvOr("PINGOU_DATABASE_URL", "pingou.db"),
		"Port":               getEnvOr("PINGOU_PORT", "8080"),
		"LogLevel":           getEnvOr("PINGOU_LOG_LEVEL", "INFO"),
		"CORSAllowedOrigins": getEnvList("PINGOU_CORS_ALLOWED_ORIGINS"),
		"MaxRedirects":       getEnvInt("PINGOU_MAX_REDIRECTS", 5),
		"GlobalTimeout":      getEnvInt("PINGOU_GLOBAL_TIMEOUT", 60),
	}

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("json error: %w", err)
	}

	fmt.Println(string(out))
	return nil
}

func getEnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return def
}

func getEnvList(key string) []string {
	if v := os.Getenv(key); v != "" {
		var result []string
		for _, s := range splitAndTrim(v, ",") {
			if s != "" {
				result = append(result, s)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range split(s, sep) {
		result = append(result, trim(part))
	}
	return result
}

func split(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
