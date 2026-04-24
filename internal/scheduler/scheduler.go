package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type job struct {
	monitor  *domain.Monitor
	cancelFn context.CancelFunc
}

type Scheduler struct {
	monitors     domain.MonitorRepository
	checker      domain.Checker
	stateMachine *service.StateMachine
	jobs         map[string]*job
	mu           sync.Mutex
}

func NewScheduler(
	monitors domain.MonitorRepository,
	checker domain.Checker,
	stateMachine *service.StateMachine,
) *Scheduler {
	return &Scheduler{
		monitors:     monitors,
		checker:      checker,
		stateMachine: stateMachine,
		jobs:         make(map[string]*job),
	}
}

// Start carrega todos os monitors ativos e inicia os jobs
func (s *Scheduler) Start(ctx context.Context) error {
	enabled := true
	monitors, _, err := s.monitors.FindAll(ctx, domain.MonitorFilter{
		Enabled: &enabled,
		Limit:   100,
	})
	if err != nil {
		return err
	}

	for _, m := range monitors {
		s.startJob(ctx, m)
	}

	slog.Info("scheduler started", "jobs", len(monitors))
	return nil
}

// Reload sincroniza os jobs com o estado atual do banco
// chamado após create/update/delete de monitors
func (s *Scheduler) Reload(ctx context.Context, monitorID string) error {
	m, err := s.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	s.stopJob(monitorID)

	if m != nil && m.Enabled {
		s.startJob(ctx, m)
	}

	return nil
}

func (s *Scheduler) startJob(ctx context.Context, m *domain.Monitor) {
	jobCtx, cancel := context.WithCancel(ctx)

	s.mu.Lock()
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
	// executa imediatamente na primeira vez
	s.runCheck(ctx, m)

	ticker := time.NewTicker(time.Duration(m.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// busca o monitor atualizado do banco antes de checar
			// (interval ou timeout podem ter mudado via PATCH)
			updated, err := s.monitors.FindByID(ctx, m.ID)
			if err != nil || updated == nil || !updated.Enabled {
				return
			}
			s.runCheck(ctx, updated)
		}
	}
}

func (s *Scheduler) runCheck(ctx context.Context, m *domain.Monitor) {
	result := s.checker.Check(ctx, m)

	if err := s.stateMachine.Process(ctx, m, result); err != nil {
		slog.Error("state machine error", "monitor_id", m.ID, "err", err)
	}
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
