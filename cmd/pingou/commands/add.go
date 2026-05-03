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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adiciona um novo monitor",
	Long: `Cria um novo monitor de health check.

Exemplo:
  pingou add --name "Meu Servico" --url "https://api.exemplo.com/health"

Para proteger o comando, use a flag --key:
  pingou add --name "Meu Servico" --url "https://api.exemplo.com/health" --key SUA_API_KEY`,
	RunE:              runAdd,
	PersistentPreRunE: requireKey,
	DisableAutoGenTag: true,
}

// flags do comando add
var (
	flagName             string
	flagURL              string
	flagIntervalSeconds  int
	flagTimeoutSeconds   int
	flagFailureThreshold int
)

// valores padrao
const (
	defaultInterval  = 60
	defaultTimeout   = 5
	defaultThreshold = 3
)

func init() {
	addCmd.Flags().StringVarP(&flagName, "name", "n", "", "nome do monitor (obrigatorio)")
	addCmd.Flags().StringVarP(&flagURL, "url", "u", "", "URL para verificar health (obrigatorio)")
	addCmd.Flags().IntVarP(&flagIntervalSeconds, "interval", "i", defaultInterval,
		"intervalo em segundos entre checagens")
	addCmd.Flags().IntVarP(&flagTimeoutSeconds, "timeout", "t", defaultTimeout,
		"timeout em segundos para cada checagem")
	addCmd.Flags().IntVarP(&flagFailureThreshold, "threshold", "K", defaultThreshold,
		"numero de falhas antes de marcar como down")

	// Requerer explicitamente estas flags
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("url")
}

func runAdd(cmd *cobra.Command, args []string) error {
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

	// Validar entrada usando servico existente
	input := service.CreateMonitorInput{
		Name:             flagName,
		URL:              flagURL,
		IntervalSeconds:  flagIntervalSeconds,
		TimeoutSeconds:   flagTimeoutSeconds,
		FailureThreshold: flagFailureThreshold,
		Enabled:          true,
	}

	// Criar monitor
	m, err := monitorSvc.Create(context.Background(), input)
	if err != nil {
		return fmt.Errorf("falha ao criar monitor: %w", err)
	}

	fmt.Printf("Monitor criado com sucesso!\n")
	fmt.Printf("  ID:     %s\n", m.ID)
	fmt.Printf("  Nome:  %s\n", m.Name)
	fmt.Printf("  URL:    %s\n", m.URL)
	fmt.Printf("  Estado: %s\n", m.CurrentState)

	return nil
}
