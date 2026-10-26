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
	rootCtx        context.Context // contexto de vida dos jobs (definido em Start)
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
		rootCtx:        context.Background(),
	}
}

// Start carrega todos os monitors ativos e inicia os jobs
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	s.rootCtx = ctx
	s.mu.Unlock()

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

	slog.Info("scheduler started", "jobs", len(monitors))
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
	slog.Debug("running check", "monitor_id", m.ID, "name", m.Name, "url", m.URL)

	result := s.checker.Check(ctx, m)

	slog.Debug("check result",
		"monitor_id", m.ID,
		"success", result.Success,
		"latency_ms", result.LatencyMs,
		"status_code", result.StatusCode,
	)

	// Serializa só a parte de banco (SQLite com MaxOpenConns=1); o check HTTP
	// fica fora do lock para um site lento não atrasar os outros monitores
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	// Criar nova UnitOfWork para este check
	uow := service.NewUnitOfWork(s.db, s.checkRepoTx, s.monitorRepoTx, s.incidentRepoTx)
	sm := service.NewStateMachine(uow, s.notifier)

	if err := sm.Process(ctx, m, result); err != nil {
		slog.Error("state machine error", "monitor_id", m.ID, "err", err)
	}
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

	for id, j := range s.jobs {
		j.cancelFn()
		delete(s.jobs, id)
	}
	slog.Info("scheduler stopped")
}
