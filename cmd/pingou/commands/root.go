package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Logica de pre-run se necessaria
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
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(migrateCmd)
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
