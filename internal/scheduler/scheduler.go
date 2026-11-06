package scheduler

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type job struct {
	monitor  *domain.Monitor
	cancelFn context.CancelFunc
}

type Scheduler struct {
	monitors           domain.MonitorRepository
	checker            domain.Checker
	db                 *sql.DB
	checkRepoTx        *repository.CheckRepoTx
	monitorRepoTx      *repository.MonitorRepoTx
	incidentRepoTx     *repository.IncidentRepoTx
	outboxRepoTx       *repository.OutboxRepoTx
	channels           []string                                // canais ativos: uma linha do outbox por canal
	wake               func()                                  // acorda o worker depois do commit
	maintenance        domain.MaintenanceRepository            // nil = sem janelas de manutenção
	maintenanceNotices func(ctx context.Context) (bool, error) // interruptor do aviso de manutenção; nil = desligado
	now                func() time.Time                        // relógio do gate e do tick de manutenção (testes trocam)
	jobs               map[string]*job
	dbMu               sync.Mutex // serializa acesso ao DB (SQLite com MaxOpenConns=1)
	mu                 sync.Mutex
	intervalUnit       time.Duration      // unidade de interval_seconds; testes reduzem para rodar rápido
	heartbeatTick      time.Duration      // frequência do avaliador de prazos do heartbeat
	maintenanceTick    time.Duration      // frequência do aplicador de início/fim das janelas
	rootCtx            context.Context    // contexto de vida dos jobs (definido em Start)
	startedAt          time.Time          // prazos vencidos antes disso venceram com o Pingou fora do ar
	stopHeartbeats     context.CancelFunc // para o avaliador de heartbeat e o tick de manutenção
	stopping           bool               // Stop foi chamado: nenhum job novo começa (protegido por mu)
	wg                 sync.WaitGroup     // loops de check e avaliador de heartbeat em andamento
}

func NewScheduler(
	monitors domain.MonitorRepository,
	checker domain.Checker,
	db *sql.DB,
	checkRepoTx *repository.CheckRepoTx,
	monitorRepoTx *repository.MonitorRepoTx,
	incidentRepoTx *repository.IncidentRepoTx,
	outboxRepoTx *repository.OutboxRepoTx,
	channels []string,
	wake func(),
) *Scheduler {
	return &Scheduler{
		monitors:        monitors,
		checker:         checker,
		db:              db,
		checkRepoTx:     checkRepoTx,
		monitorRepoTx:   monitorRepoTx,
		incidentRepoTx:  incidentRepoTx,
		outboxRepoTx:    outboxRepoTx,
		channels:        channels,
		wake:            wake,
		jobs:            make(map[string]*job),
		intervalUnit:    time.Second,
		heartbeatTick:   10 * time.Second,
		maintenanceTick: 10 * time.Second,
		now:             time.Now,
		rootCtx:         context.Background(),
	}
}

// WithMaintenance liga as janelas de manutenção: o gate em withMonitor e o tick
// de início/fim (iniciado em Start).
func (s *Scheduler) WithMaintenance(repo domain.MaintenanceRepository) *Scheduler {
	s.maintenance = repo
	return s
}

// WithMaintenanceNotices liga o aviso de início e fim das janelas agendadas.
// enabled lê o interruptor das settings; só é chamado quando há janela vencida.
func (s *Scheduler) WithMaintenanceNotices(enabled func(ctx context.Context) (bool, error)) *Scheduler {
	s.maintenanceNotices = enabled
	return s
}

// Start carrega todos os monitors ativos e inicia os jobs
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	s.rootCtx = ctx
	s.startedAt = s.now()
	s.stopping = false
	hbCtx, cancel := context.WithCancel(ctx)
	s.stopHeartbeats = cancel
	s.wg.Add(1)
	s.mu.Unlock()
	go func() {
		defer s.wg.Done()
		s.runHeartbeats(hbCtx)
	}()

	if s.maintenance != nil {
		// reinício no meio de uma janela, ou janela que venceu com o Pingou
		// fora: aplica o pendente antes de subir os jobs
		s.applyMaintenance(ctx)
		s.mu.Lock()
		s.wg.Add(1)
		s.mu.Unlock()
		go func() {
			defer s.wg.Done()
			s.runMaintenance(hbCtx)
		}()
	}

	enabled := true
	monitors, _, err := s.monitors.FindAll(ctx, domain.MonitorFilter{
		Enabled: &enabled,
		Limit:   100,
	})
	if err != nil {
		return err
	}

	for _, m := range monitors {
		s.startJob(m)
	}

	slog.Info("scheduler started", "monitors", len(monitors))
	return nil
}

// Reload sincroniza os jobs com o estado atual do banco
// chamado após create/update/delete de monitors.
// ctx é usado só para a leitura do banco: o job NÃO herda esse contexto,
// pois normalmente é o de uma requisição HTTP, cancelado ao fim da resposta.
func (s *Scheduler) Reload(ctx context.Context, monitorID string) error {
	m, err := s.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	s.stopJob(monitorID)

	if m != nil && m.Enabled {
		s.startJob(m)
	}

	return nil
}

func (s *Scheduler) startJob(m *domain.Monitor) {
	// heartbeat não tem loop ativo: os checks vêm do ping e do avaliador de prazos
	if m.Type == domain.TypeHeartbeat {
		return
	}

	s.mu.Lock()
	if s.stopping {
		// desligando: um Reload vindo de uma requisição em andamento não
		// pode começar check novo
		s.mu.Unlock()
		return
	}
	if old, ok := s.jobs[m.ID]; ok {
		old.cancelFn()
	}
	jobCtx, cancel := context.WithCancel(s.rootCtx)
	s.jobs[m.ID] = &job{monitor: m, cancelFn: cancel}
	s.wg.Add(1)
	s.mu.Unlock()

	go func() {
		defer s.wg.Done()
		s.runLoop(jobCtx, m)
	}()
	slog.Info("job started", "monitor_id", m.ID, "name", m.Name, "interval_seconds", m.IntervalSeconds)
}

func (s *Scheduler) stopJob(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if j, ok := s.jobs[id]; ok {
		j.cancelFn()
		delete(s.jobs, id)
		slog.Info("job stopped", "monitor_id", id)
	}
}

// runLoop: ctx é o do job (filho do contexto de Start, nunca o de uma
// requisição HTTP). É cancelado no desligamento e quando o monitor é editado,
// pausado ou excluído.
func (s *Scheduler) runLoop(ctx context.Context, m *domain.Monitor) {
	s.runCheck(ctx, m)

	ticker := time.NewTicker(time.Duration(m.IntervalSeconds) * s.intervalUnit)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updated, err := s.monitors.FindByID(ctx, m.ID)
			if err != nil || updated == nil || !updated.Enabled {
				return
			}
			s.runCheck(ctx, updated)
		}
	}
}

func (s *Scheduler) runCheck(jobCtx context.Context, m *domain.Monitor) {
	// em manutenção: nem faz a requisição (o gate em withMonitor descartaria)
	if s.inMaintenance(jobCtx, m.ID) {
		slog.Debug("check skipped: maintenance", "monitor_id", m.ID)
		return
	}
	slog.Debug("running check", "monitor_id", m.ID, "name", m.Name, "type", m.Type)

	checkCtx, cancel := context.WithTimeout(jobCtx, time.Duration(m.TimeoutSeconds)*time.Second)
	result := s.checker.Check(checkCtx, m)
	cancel()

	// Job parado durante o check (desligamento, edição, pausa ou exclusão): o
	// resultado pode ser um "context canceled" que não é queda de verdade.
	// Descarta, para não abrir incidente falso.
	if jobCtx.Err() != nil {
		slog.Debug("check discarded: job stopped", "monitor_id", m.ID)
		return
	}

	slog.Debug("check result",
		"monitor_id", m.ID,
		"success", result.Success,
		"latency_ms", result.LatencyMs,
		"status_code", result.StatusCode,
	)

	// a gravação não é cortada no meio: termina com commit ou rollback, e o
	// Wait do desligamento espera por ela
	if err := s.Record(context.WithoutCancel(jobCtx), m.ID, result, nil); err != nil {
		slog.Error("state machine error", "monitor_id", m.ID, "err", err)
	}
}

// Record grava um resultado de check pela state machine. pingedAt != nil indica
// um ping de heartbeat e atualiza last_ping_at.
func (s *Scheduler) Record(ctx context.Context, monitorID string, result domain.CheckResult, pingedAt *time.Time) error {
	return s.withMonitor(ctx, monitorID, func(m *domain.Monitor) *domain.CheckResult {
		if pingedAt != nil {
			m.LastPingAt = pingedAt
		}
		return &result
	})
}

// withMonitor serializa a parte de banco (SQLite com MaxOpenConns=1) e recarrega
// o monitor dentro do lock, para não gravar por cima de uma versão mais nova
// (ex.: last_ping_at de um ping que acabou de chegar). decide devolve nil para
// não registrar nada. O check em si fica fora do lock, para um site lento não
// atrasar os outros monitores.
func (s *Scheduler) withMonitor(ctx context.Context, monitorID string, decide func(m *domain.Monitor) *domain.CheckResult) error {
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	m, err := s.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}
	if m == nil || !m.Enabled {
		return nil // excluído ou pausado enquanto o check rodava
	}
	// Gate de manutenção: ponto único por onde passa todo resultado (check
	// HTTP/TCP, heartbeat vencido e ping). Sem gravar o check, não há incidente
	// nem aviso, inclusive o de certificado (processCert roda dentro de Process).
	if s.inMaintenance(ctx, m.ID) {
		return nil
	}

	result := decide(m)
	if result == nil {
		return nil
	}

	// Criar nova UnitOfWork para este check
	uow := service.NewUnitOfWork(s.db, s.checkRepoTx, s.monitorRepoTx, s.incidentRepoTx, s.outboxRepoTx)
	sm := service.NewStateMachine(uow, s.channels, s.wake)
	return sm.Process(ctx, m, *result)
}

func (s *Scheduler) hasJob(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.jobs[id]
	return ok
}

// Stop cancela todos os jobs e o avaliador de heartbeat, e impede jobs novos
// até o próximo Start. Não espera: use Wait para isso.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopping = true
	if s.stopHeartbeats != nil {
		s.stopHeartbeats()
		s.stopHeartbeats = nil
	}
	for id, j := range s.jobs {
		j.cancelFn()
		delete(s.jobs, id)
	}
	slog.Info("scheduler stopped")
}

// Wait espera os loops em andamento terminarem (check atual descartado e
// gravação em andamento concluída). Devolve ctx.Err() se o prazo acabar antes.
func (s *Scheduler) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
