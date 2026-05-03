package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "dev"
	commit  = "none"
	date    = "unknown"
	flagKey string
)

// RootCmd representa o comando root do CLI
var RootCmd = &cobra.Command{
	Use:   "pingou",
	Short: "Pingou CLI - Gestao de saude de servicos",
	Long: `Pingou e uma ferramenta CLI para gerenciamento de health checks
e monitoracao de servicos.

Uso sem subcomando executa 'serve' por padrao.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Executa serve por padrao
		if err := runServe(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute executa todos os comandos
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "modo verboso")
	RootCmd.PersistentFlags().StringVar(&flagKey, "key", "", "API key para comandos protegidos (PINGOU_API_KEY)")
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(exportCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(configCmd)
}

// requireKey valida a flag --key contra PINGOU_API_KEY
func requireKey(cmd *cobra.Command, args []string) error {
	// Carregar .env se existir
	godotenv.Load()

	// Se nao passou --key, tentar usar o do .env
	requiredKey := os.Getenv("PINGOU_API_KEY")
	if requiredKey == "" {
		// Sem API key configurada - permitir (modo desenvolvimento)
		return nil
	}

	// Validar key fornecida
	if flagKey == "" {
		return fmt.Errorf("flag --key e obrigatoria. Use: pingou %s --key <API_KEY>", cmd.Name())
	}

	if flagKey != requiredKey {
		return fmt.Errorf("API key invalida")
	}

	return nil
}

// GetVersion retorna a versao atual
func GetVersion() string {
	return version
}

// GetCommit retorna o commit atual
func GetCommit() string {
	return commit
}

// GetBuildDate retorna a data de build
func GetBuildDate() string {
	return date
}