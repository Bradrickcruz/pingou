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
	monitors       domain.MonitorRepository
	checker        domain.Checker
	db             *sql.DB
	checkRepoTx    *repository.CheckRepoTx
	monitorRepoTx  *repository.MonitorRepoTx
	incidentRepoTx *repository.IncidentRepoTx
	notifier       service.Notifier
	jobs           map[string]*job
	dbMu           sync.Mutex // serializa acesso ao DB (SQLite com MaxOpenConns=1)
	mu             sync.Mutex
	intervalUnit   time.Duration   // unidade de interval_seconds; testes reduzem para rodar rápido
	heartbeatTick  time.Duration   // frequência do avaliador de prazos do heartbeat
	rootCtx        context.Context // contexto de vida dos jobs (definido em Start)
	stopHeartbeats context.CancelFunc
}

func NewScheduler(
	monitors domain.MonitorRepository,
	checker domain.Checker,
	db *sql.DB,
	checkRepoTx *repository.CheckRepoTx,
	monitorRepoTx *repository.MonitorRepoTx,
	incidentRepoTx *repository.IncidentRepoTx,
	notifier service.Notifier,
) *Scheduler {
	return &Scheduler{
		monitors:       monitors,
		checker:        checker,
		db:             db,
		checkRepoTx:    checkRepoTx,
		monitorRepoTx:  monitorRepoTx,
		incidentRepoTx: incidentRepoTx,
		notifier:       notifier,
		jobs:           make(map[string]*job),
		intervalUnit:   time.Second,
		heartbeatTick:  10 * time.Second,
		rootCtx:        context.Background(),
	}
}

// Start carrega todos os monitors ativos e inicia os jobs
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	s.rootCtx = ctx
	hbCtx, cancel := context.WithCancel(ctx)
	s.stopHeartbeats = cancel
	s.mu.Unlock()
	go s.runHeartbeats(hbCtx)

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
	if old, ok := s.jobs[m.ID]; ok {
		old.cancelFn()
	}
	jobCtx, cancel := context.WithCancel(s.rootCtx)
	s.jobs[m.ID] = &job{monitor: m, cancelFn: cancel}
	s.mu.Unlock()

	go s.runLoop(jobCtx, m)
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

func (s *Scheduler) runLoop(ctx context.Context, m *domain.Monitor) {
	// primeiro check usa contexto próprio, independente da request HTTP
	checkCtx, cancel := context.WithTimeout(context.Background(), time.Duration(m.TimeoutSeconds)*time.Second)
	s.runCheck(checkCtx, m)
	cancel()

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
			checkCtx, cancel := context.WithTimeout(context.Background(), time.Duration(updated.TimeoutSeconds)*time.Second)
			s.runCheck(checkCtx, updated)
			cancel()
		}
	}
}

func (s *Scheduler) runCheck(ctx context.Context, m *domain.Monitor) {
	slog.Debug("running check", "monitor_id", m.ID, "name", m.Name, "type", m.Type)

	result := s.checker.Check(ctx, m)

	slog.Debug("check result",
		"monitor_id", m.ID,
		"success", result.Success,
		"latency_ms", result.LatencyMs,
		"status_code", result.StatusCode,
	)

	if err := s.Record(ctx, m.ID, result, nil); err != nil {
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

	result := decide(m)
	if result == nil {
		return nil
	}

	// Criar nova UnitOfWork para este check
	uow := service.NewUnitOfWork(s.db, s.checkRepoTx, s.monitorRepoTx, s.incidentRepoTx)
	sm := service.NewStateMachine(uow, s.notifier)
	return sm.Process(ctx, m, *result)
}

func (s *Scheduler) hasJob(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.jobs[id]
	return ok
}

// Stop encerra todos os jobs
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

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
