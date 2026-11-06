package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Bradrickcruz/pingou/internal/checker"
	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
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
	loadedConfig.Version, loadedConfig.Commit, loadedConfig.BuildDate = version, commit, date

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
	maintenanceRepository := repository.NewMaintenanceRepo(db)

	// repositories com suporte a transação
	checkRepoTx := repository.NewCheckRepoTx(db)
	monitorRepoTx := repository.NewMonitorRepoTx(db)
	incidentRepoTx := repository.NewIncidentRepoTx(db)
	outboxRepoTx := repository.NewOutboxRepoTx()
	outboxRepository := repository.NewOutboxRepo(db)

	// webhookNotifier
	webhookNotifier := service.NewWebhookNotifier(func() string {
		url, _ := settingsRepository.Get(context.Background(), "webhook_url")
		return url
	}, loadedConfig.PublicURL)

	// fuso das datas no e-mail e no ntfy (o webhook segue em UTC, RFC3339)
	messageLocation := service.LocationLoader(settingsRepository.Get)

	// canais: o webhook sempre; e-mail e ntfy se configurados no ambiente
	channels := []service.Channel{webhookNotifier}
	if loadedConfig.SMTP.Enabled() {
		channels = append(channels, service.NewEmailNotifier(loadedConfig.SMTP, func() string {
			recipients, _ := settingsRepository.Get(context.Background(), "email_recipients")
			return recipients
		}, loadedConfig.PublicURL).WithLocation(messageLocation))
		slog.Info("email notifications available", "host", loadedConfig.SMTP.Host, "port", loadedConfig.SMTP.Port, "tls", loadedConfig.SMTP.TLS)
	}
	if loadedConfig.Ntfy.Enabled() {
		channels = append(channels, service.NewNtfyNotifier(loadedConfig.Ntfy, func() string {
			topic, _ := settingsRepository.Get(context.Background(), "ntfy_topic")
			return topic
		}, loadedConfig.PublicURL).WithLocation(messageLocation))
		// só o host: a URL pode ter caminho e o token nunca vai para o log
		slog.Info("ntfy notifications available", "host", loadedConfig.Ntfy.Host(), "auth", loadedConfig.Ntfy.Token != "")
	}
	channelNames := make([]string, len(channels))
	for i, c := range channels {
		channelNames[i] = c.Name()
	}

	// worker do outbox: envia os avisos gravados na transação dos checks e
	// repete os que falharam
	outboxWorker := service.NewOutboxWorker(outboxRepository, channels...)

	// checkers: o registry escolhe pelo tipo (heartbeat não tem check ativo)
	monitorChecker := checker.NewRegistry(map[domain.MonitorType]domain.Checker{
		domain.TypeHTTP: checker.NewHTTPChecker(loadedConfig.MaxRedirects, loadedConfig.GlobalTimeout),
		domain.TypeTCP:  checker.NewTCPChecker(),
	})

	// services
	monitorService := service.NewMonitorService(monitorRepository, checkRepository, incidentRepository)
	incidentService := service.NewIncidentService(incidentRepository)
	settingsService := service.NewSettingsService(settingsRepository, loadedConfig.SMTP, loadedConfig.Ntfy)
	maintenanceService := service.NewMaintenanceService(maintenanceRepository, monitorRepository)

	// scheduler
	monitorScheduler := scheduler.NewScheduler(monitorRepository, monitorChecker, db, checkRepoTx, monitorRepoTx, incidentRepoTx, outboxRepoTx, channelNames, outboxWorker.Wake).
		WithMaintenance(maintenanceRepository)
	monitorService.SetReloader(monitorScheduler)
	maintenanceService.SetEnder(monitorScheduler)
	monitorService.SetRecorder(monitorScheduler)

	// Context com signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Inicia o worker do outbox antes do scheduler: retoma os avisos pendentes de antes
	outboxWorker.Start(ctx)

	// Inicia scheduler
	if err := monitorScheduler.Start(ctx); err != nil {
		return fmt.Errorf("scheduler error: %w", err)
	}

	// Retention worker
	retention := scheduler.NewRetentionWorker(checkRepository, settingsRepository).WithOutbox(outboxRepository)
	retention.Start(ctx)

	// Inicia server
	server := handler.NewServer(loadedConfig, db, monitorService, incidentService, settingsService, maintenanceService)

	serverErr := make(chan error, 1)
	go func() {
		// ErrServerClosed é o retorno normal depois do Shutdown
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Espera sinal de encerramento (ou falha ao subir o servidor)
	var startErr error
	select {
	case <-ctx.Done():
	case err := <-serverErr:
		startErr = fmt.Errorf("server error: %w", err)
	}
	stop()
	forceExitOnSecondSignal()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = runShutdown(shutdownCtx, []shutdownStep{
		{"scheduler", func(context.Context) error {
			monitorScheduler.Stop() // nenhum check novo começa
			return nil
		}},
		{"http server", server.Shutdown}, // inclui pings /ping/{token} em andamento
		{"checks", func(ctx context.Context) error {
			if err := monitorScheduler.Wait(ctx); err != nil {
				return err
			}
			return retention.Wait(ctx)
		}},
		{"notifications", func(ctx context.Context) error {
			// para de pegar avisos novos e termina o envio em curso; o resto
			// fica pending no banco para a próxima subida
			return outboxWorker.Wait(ctx)
		}},
		{"database", func(context.Context) error { return db.Close() }},
		{"lock", func(context.Context) error { releaseLock(); return nil }},
	})
	if startErr != nil {
		return startErr
	}
	return err
}

// shutdownTimeout é o orçamento total do desligamento. Fica abaixo do
// stop_grace_period do docker-compose.yml (20 s); o padrão do Docker (10 s)
// mataria o processo antes.
const shutdownTimeout = 15 * time.Second

type shutdownStep struct {
	name string
	run  func(ctx context.Context) error
}

// runShutdown executa as etapas em ordem, todas no mesmo orçamento (ctx). Uma
// etapa que falha ou estoura o prazo não pula as seguintes: banco e lock são
// sempre fechados. Devolve erro se alguma etapa falhou.
func runShutdown(ctx context.Context, steps []shutdownStep) error {
	slog.Info("shutdown started", "timeout", shutdownTimeout.String())
	began := time.Now()
	var failed []string
	for _, st := range steps {
		t := time.Now()
		if err := st.run(ctx); err != nil {
			slog.Error("shutdown step failed", "step", st.name, "err", err, "elapsed", time.Since(t).String())
			failed = append(failed, st.name)
			continue
		}
		slog.Info("shutdown step done", "step", st.name, "elapsed", time.Since(t).String())
	}
	if len(failed) > 0 {
		slog.Error("shutdown incomplete", "failed_steps", failed, "elapsed", time.Since(began).String())
		return fmt.Errorf("shutdown incomplete: %v", failed)
	}
	slog.Info("shutdown complete", "elapsed", time.Since(began).String())
	return nil
}

// forceExitOnSecondSignal: um segundo Ctrl+C ou SIGTERM durante o desligamento
// encerra na hora, com código 1.
func forceExitOnSecondSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-ch
		slog.Error("second signal during shutdown, exiting now", "signal", sig.String())
		releaseLock()
		os.Exit(1)
	}()
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

	// Verifica se o processo ainda esta rodando. O proprio PID conta como lock
	// velho: num container reiniciado o processo volta como PID 1, igual ao
	// lock deixado por uma execucao que morreu sem liberar (ex.: kill -9).
	if pid > 0 && pid != os.Getpid() {
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
