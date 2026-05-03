package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Bradrickcruz/pingou/internal/checker"
	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/handler"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/scheduler"
	"github.com/Bradrickcruz/pingou/internal/service"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Lock file paths
var (
	lockDir  = filepath.Join(os.Getenv("HOME"), ".pingou")
	lockFile = filepath.Join(lockDir, "pingou.lock")
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor API e dashboard",
	Long:  "Inicia o servidor HTTP com API REST e dashboard SPA",
	RunE:  runServe,
	// Sem subcomandos pais para que 'pingou' execute 'serve' por padrao
}

func init() {
	serveCmd.Annotations = map[string]string{
		"supports": "DASHBOARD",
	}
}

// runServe executa o servidor
func runServe(cmd *cobra.Command, args []string) error {
	// Protecao anti-multinstancia
	if err := acquireLock(); err != nil {
		return fmt.Errorf("servidor ja esta em execucao na maquina atual: %w", err)
	}
	defer releaseLock()

	// Carrega .env se existir
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Warn("failed to load .env", "err", err)
	}

	slog.Info("Pingou starting", "version", version, "commit", commit, "buildDate", date)

	loadedConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loadedConfig.LogLevel,
	}))
	slog.SetDefault(logger)

	db, err := database.Open(loadedConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	defer db.Close()

	// repositories
	monitorRepository := repository.NewMonitorRepo(db)
	checkRepository := repository.NewCheckRepo(db)
	incidentRepository := repository.NewIncidentRepo(db)
	settingsRepository := repository.NewSettingsRepo(db)

	// repositories com suporte a transação
	checkRepoTx := repository.NewCheckRepoTx(db)
	monitorRepoTx := repository.NewMonitorRepoTx(db)
	incidentRepoTx := repository.NewIncidentRepoTx(db)

	// webhookNotifier
	webhookNotifier := service.NewWebhookNotifier(func() string {
		url, _ := settingsRepository.Get(context.Background(), "webhook_url")
		return url
	})

	// checker
	httpChecker := checker.NewHTTPChecker(loadedConfig.MaxRedirects, loadedConfig.GlobalTimeout)

	// services
	monitorService := service.NewMonitorService(monitorRepository, checkRepository, incidentRepository)
	incidentService := service.NewIncidentService(incidentRepository)
	settingsService := service.NewSettingsService(settingsRepository)

	// scheduler
	monitorScheduler := scheduler.NewScheduler(monitorRepository, httpChecker, db, checkRepoTx, monitorRepoTx, incidentRepoTx, webhookNotifier)
	monitorService.SetReloader(monitorScheduler)

	// Context com signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Inicia scheduler
	if err := monitorScheduler.Start(ctx); err != nil {
		return fmt.Errorf("scheduler error: %w", err)
	}

	// Retention worker
	retention := scheduler.NewRetentionWorker(checkRepository, settingsRepository)
	retention.Start(ctx)

	// Inicia server
	server := handler.NewServer(loadedConfig, db, monitorService, incidentService, settingsService)

	go func() {
		if err := server.Start(); err != nil {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	// Espera sinal de encerramento
	<-ctx.Done()
	slog.Info("shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("server shutdown error", "err", err)
	}

	monitorScheduler.Stop()

	return nil
}

// acquireLock adquire o lock file para evitar multiplas instancias
func acquireLock() error {
	// Cria diretorio se nao existir
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		return fmt.Errorf("nao foi possivel criar diretorio de lock: %w", err)
	}

	// Tenta abrir ou criar o lock file
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		// Se ja existe, verifica se o processo ainda esta rodando
		if os.IsExist(err) {
			return checkExistingLock()
		}
		return fmt.Errorf("erro ao criar lock file: %w", err)
	}

	// Escreve PID e timestamp no lock file
	pid := os.Getpid()
	_, err = fmt.Fprintf(file, "%d\n%d\n", pid, time.Now().Unix())
	if err != nil {
		file.Close()
		os.Remove(lockFile)
		return fmt.Errorf("erro ao escrever no lock file: %w", err)
	}
	file.Close()

	return nil
}

// checkExistingLock verifica se existe um lock de uma execucao anterior
func checkExistingLock() error {
	data, err := os.ReadFile(lockFile)
	if err != nil {
		return fmt.Errorf("lock file existe mas nao pode ser lido: %w", err)
	}

	// Se nao conseguir ler o PID, assume que ha uma execucao anterior
	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		// Lock file invalido, remove e permite execucao
		os.Remove(lockFile)
		return nil
	}

	// Verifica se o processo ainda esta rodando
	if pid > 0 {
		proc, err := os.FindProcess(pid)
		if err == nil && proc.Signal(syscall.Signal(0)) == nil {
			// Processo esta rodando
			return fmt.Errorf("processo %d ainda esta ativo", pid)
		}
	}

	// Processo nao esta mais rodando, remove lock antigo e permite execucao
	os.Remove(lockFile)
	return nil
}

// releaseLock remove o lock file
func releaseLock() {
	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		slog.Error("erro ao remover lock file", "err", err)
	}
}
