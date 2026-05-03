package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove um monitor",
	Long: `Remove um monitor de health check pelo ID.

Exemplo:
  pingou rm --id 019def55-0c92-76a5-a7c4-b573ab447ac2

Para proteger o comando, use a flag --key:
  pingou rm --id 019def55-0c92-76a5-a7c4-b573ab447ac2 --key SUA_API_KEY`,
	RunE:              runRM,
	PersistentPreRunE:  requireKey,
	DisableAutoGenTag: true,
}

// flags do comando rm
var flagMonitorID string

func init() {
	rmCmd.Flags().StringVarP(&flagMonitorID, "id", "i", "", "ID do monitor a remover (obrigatorio)")
	rmCmd.MarkFlagRequired("id")
}

func runRM(cmd *cobra.Command, args []string) error {
	// Carregar .env
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "WARNING: failed to load .env: %v\n", err)
	}

	// Obter DSN do banco
	dsn := os.Getenv("PINGOU_DATABASE_URL")
	if dsn == "" {
		dsn = "pingou.db"
	}

	// Abrir banco
	db, err := database.Open(dsn)
	if err != nil {
		return fmt.Errorf("falha ao abrir banco: %w", err)
	}
	defer db.Close()

	// Criar repositorios
	monitorRepo := repository.NewMonitorRepo(db)
	checkRepo := repository.NewCheckRepo(db)
	incidentRepo := repository.NewIncidentRepo(db)

	// Criar servico
	monitorSvc := service.NewMonitorService(monitorRepo, checkRepo, incidentRepo)

	// Verificar se monitor existe antes de remover
	m, err := monitorSvc.GetByID(context.Background(), flagMonitorID)
	if err != nil {
		return fmt.Errorf("monitor nao encontrado: %w", err)
	}

	// Remover monitor
	if err := monitorSvc.Delete(context.Background(), flagMonitorID); err != nil {
		return fmt.Errorf("falha ao remover monitor: %w", err)
	}

	fmt.Printf("Monitor removido com sucesso!\n")
	fmt.Printf("  ID:    %s\n", m.ID)
	fmt.Printf("  Nome: %s\n", m.Name)

	return nil
}