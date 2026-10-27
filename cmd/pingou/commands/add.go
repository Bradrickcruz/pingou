package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adiciona um novo monitor",
	Long: `Cria um novo monitor de health check.

Exemplos:
  pingou add --name "Meu Servico" --url "https://api.exemplo.com/health"
  pingou add --name "Postgres" --type tcp --host db.interno --port 5432
  pingou add --name "Backup diario" --type heartbeat -i 86400 --grace 3600

Para proteger o comando, use a flag --key:
  pingou add --name "Meu Servico" --url "https://api.exemplo.com/health" --key SUA_API_KEY`,
	RunE:              runAdd,
	PersistentPreRunE: requireKey,
	DisableAutoGenTag: true,
}

// flags do comando add
var (
	flagName             string
	flagType             string
	flagURL              string
	flagHost             string
	flagPort             int
	flagGraceSeconds     int
	flagIntervalSeconds  int
	flagTimeoutSeconds   int
	flagFailureThreshold int
)

// valores padrao
const (
	defaultInterval  = 60
	defaultTimeout   = 5
	defaultThreshold = 3
	defaultGrace     = 60
)

func init() {
	addCmd.Flags().StringVarP(&flagName, "name", "n", "", "nome do monitor (obrigatorio)")
	addCmd.Flags().StringVar(&flagType, "type", "http", "tipo do monitor: http, tcp ou heartbeat")
	addCmd.Flags().StringVarP(&flagURL, "url", "u", "", "URL para verificar health (obrigatorio para http)")
	addCmd.Flags().StringVar(&flagHost, "host", "", "host (obrigatorio para tcp)")
	addCmd.Flags().IntVar(&flagPort, "port", 0, "porta (obrigatorio para tcp)")
	addCmd.Flags().IntVar(&flagGraceSeconds, "grace", defaultGrace, "tolerancia em segundos (heartbeat)")
	addCmd.Flags().IntVarP(&flagIntervalSeconds, "interval", "i", defaultInterval,
		"intervalo em segundos entre checagens")
	addCmd.Flags().IntVarP(&flagTimeoutSeconds, "timeout", "t", defaultTimeout,
		"timeout em segundos para cada checagem (nao se aplica a heartbeat)")
	addCmd.Flags().IntVarP(&flagFailureThreshold, "threshold", "K", defaultThreshold,
		"numero de falhas antes de marcar como down (nao se aplica a heartbeat)")

	// url/host/port dependem do tipo: o service valida
	addCmd.MarkFlagRequired("name")
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
		Type:             flagType,
		URL:              flagURL,
		Host:             flagHost,
		Port:             flagPort,
		GraceSeconds:     &flagGraceSeconds,
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
	fmt.Printf("  Tipo:   %s\n", m.Type)
	switch m.Type {
	case domain.TypeHTTP:
		fmt.Printf("  URL:    %s\n", m.URL)
	case domain.TypeTCP:
		fmt.Printf("  Alvo:   %s:%d\n", m.Config.Host, m.Config.Port)
	case domain.TypeHeartbeat:
		// so o caminho: o endereco publico do Pingou depende de onde ele roda
		fmt.Printf("  Ping:   /ping/%s\n", *m.PingToken)
	}
	fmt.Printf("  Estado: %s\n", m.CurrentState)

	return nil
}
