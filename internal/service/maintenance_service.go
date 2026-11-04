package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// maxMaintenanceDuration é o teto de uma janela (D1).
const maxMaintenanceDuration = 7 * 24 * time.Hour

// ErrMaintenanceConflict: algum monitor do lote já tem janela viva (D2, D5).
var ErrMaintenanceConflict = errors.New("maintenance window conflict")

// MaintenanceConflictError lista os monitores que já têm janela viva.
type MaintenanceConflictError struct{ MonitorIDs []string }

func (e *MaintenanceConflictError) Error() string {
	return fmt.Sprintf("monitors already have a scheduled or running maintenance window: %s", strings.Join(e.MonitorIDs, ", "))
}

func (e *MaintenanceConflictError) Is(target error) bool { return target == ErrMaintenanceConflict }

// MaintenanceEnder aplica o fim de uma janela na hora (Scheduler): serializa com
// os checks, finaliza e recarrega o job, que roda o primeiro check.
type MaintenanceEnder interface {
	EndMaintenance(ctx context.Context, w *domain.MaintenanceWindow, endsAt time.Time) error
}

type MaintenanceService struct {
	windows  domain.MaintenanceRepository
	monitors domain.MonitorRepository
	ender    MaintenanceEnder
	clock    func() time.Time
}

func NewMaintenanceService(windows domain.MaintenanceRepository, monitors domain.MonitorRepository) *MaintenanceService {
	return &MaintenanceService{windows: windows, monitors: monitors, clock: time.Now}
}

func (s *MaintenanceService) SetEnder(e MaintenanceEnder) { s.ender = e }

// WithClock troca o relógio (testes).
func (s *MaintenanceService) WithClock(clock func() time.Time) *MaintenanceService {
	s.clock = clock
	return s
}

type CreateMaintenanceInput struct {
	MonitorIDs []string
	StartsAt   time.Time
	EndsAt     time.Time
	Reason     *string
}

// CreateBatch cria uma janela por monitor, tudo ou nada (D5).
func (s *MaintenanceService) CreateBatch(ctx context.Context, in CreateMaintenanceInput) ([]*domain.MaintenanceWindow, error) {
	now := s.clock().UTC()
	ids := dedupe(in.MonitorIDs)
	if len(ids) == 0 {
		return nil, fmt.Errorf("%w: monitor_ids must not be empty", ErrValidation)
	}
	if len(ids) > maxMonitors {
		return nil, fmt.Errorf("%w: at most %d monitors per request", ErrValidation, maxMonitors)
	}
	// o banco guarda segundos: compara já com a precisão gravada
	starts, ends := in.StartsAt.UTC().Truncate(time.Second), in.EndsAt.UTC().Truncate(time.Second)
	if !starts.After(now) {
		return nil, fmt.Errorf("%w: starts_at must be in the future", ErrValidation)
	}
	if !ends.After(starts) {
		return nil, fmt.Errorf("%w: ends_at must be after starts_at", ErrValidation)
	}
	if ends.Sub(starts) > maxMaintenanceDuration {
		return nil, fmt.Errorf("%w: a maintenance window can last at most 7 days", ErrValidation)
	}
	reason, err := normalizeReason("reason", in.Reason)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		m, err := s.monitors.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if m == nil {
			return nil, fmt.Errorf("%w: monitor %s", ErrNotFound, id)
		}
	}
	if err := s.checkConflicts(ctx, ids); err != nil {
		return nil, err
	}

	windows := make([]*domain.MaintenanceWindow, len(ids))
	for i, id := range ids {
		wid, err := newUUIDv7()
		if err != nil {
			return nil, fmt.Errorf("generate maintenance id: %w", err)
		}
		windows[i] = &domain.MaintenanceWindow{
			ID: wid, MonitorID: id, StartsAt: starts, EndsAt: ends, Reason: reason, CreatedAt: now,
		}
	}
	if err := s.windows.CreateBatch(ctx, windows); err != nil {
		// corrida com outro pedido: o índice único recusou; diz quais monitores
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if cerr := s.checkConflicts(ctx, ids); cerr != nil {
				return nil, cerr
			}
		}
		return nil, err
	}
	return windows, nil
}

func (s *MaintenanceService) checkConflicts(ctx context.Context, ids []string) error {
	live, err := s.windows.FindLiveByMonitors(ctx, ids)
	if err != nil {
		return err
	}
	if len(live) == 0 {
		return nil
	}
	conflict := &MaintenanceConflictError{}
	for _, id := range ids { // na ordem do pedido
		if live[id] != nil {
			conflict.MonitorIDs = append(conflict.MonitorIDs, id)
		}
	}
	return conflict
}

// Cancel apaga a janela que ainda não começou ou encerra na hora a que está em
// curso (o monitor volta como no fim normal). Janela finalizada: ErrNotFound.
func (s *MaintenanceService) Cancel(ctx context.Context, id string) error {
	w, err := s.windows.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if w == nil || w.FinalizedAt != nil {
		return ErrNotFound
	}
	now := s.clock().UTC().Truncate(time.Second)
	if now.Before(w.StartsAt) {
		return s.windows.Delete(ctx, id)
	}
	// ends_at precisa ficar depois de starts_at (CHECK): cancelar no segundo
	// exato do início encerra 1 s depois dele
	endsAt := now
	if !endsAt.After(w.StartsAt) {
		endsAt = w.StartsAt.Add(time.Second)
	}
	if s.ender != nil {
		return s.ender.EndMaintenance(ctx, w, endsAt)
	}
	_, err = s.windows.Finalize(ctx, w, now, &endsAt)
	return err
}

// LiveByMonitors devolve a janela viva (agendada ou em curso) de cada monitor.
func (s *MaintenanceService) LiveByMonitors(ctx context.Context, ids []string) (map[string]*domain.MaintenanceWindow, error) {
	return s.windows.FindLiveByMonitors(ctx, ids)
}

// Now é o relógio do serviço (a API calcula "active" com ele, não com o do navegador).
func (s *MaintenanceService) Now() time.Time { return s.clock().UTC() }

func dedupe(ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out
}
