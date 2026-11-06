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

// ErrMaintenanceLocked: a janela já começou e o horário não pode mudar (D12).
var ErrMaintenanceLocked = errors.New("maintenance window already started: only the reason can change")

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
	starts, ends, err := validateWindowTimes(in.StartsAt, in.EndsAt, now)
	if err != nil {
		return nil, err
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

// validateWindowTimes aplica D1 (início no futuro, fim depois do início, até 7
// dias) e devolve as datas na precisão do banco (segundos, UTC). Regra única
// para criar e editar.
func validateWindowTimes(startsAt, endsAt, now time.Time) (time.Time, time.Time, error) {
	starts, ends := startsAt.UTC().Truncate(time.Second), endsAt.UTC().Truncate(time.Second)
	if !starts.After(now) {
		return starts, ends, fmt.Errorf("%w: starts_at must be in the future", ErrValidation)
	}
	if !ends.After(starts) {
		return starts, ends, fmt.Errorf("%w: ends_at must be after starts_at", ErrValidation)
	}
	if ends.Sub(starts) > maxMaintenanceDuration {
		return starts, ends, fmt.Errorf("%w: a maintenance window can last at most 7 days", ErrValidation)
	}
	return starts, ends, nil
}

// List devolve as janelas vivas (agendadas e em curso) com o nome do monitor,
// por início (D14).
func (s *MaintenanceService) List(ctx context.Context) ([]*domain.MaintenanceWindowView, error) {
	return s.windows.ListLive(ctx)
}

// UpdateMaintenanceInput: campo nil não muda. ReasonSet com Reason nil ou ""
// apaga o motivo.
type UpdateMaintenanceInput struct {
	StartsAt  *time.Time
	EndsAt    *time.Time
	Reason    *string
	ReasonSet bool
}

// Update edita uma janela viva (D12, D13). Agendada: início, fim e motivo.
// Iniciada (pelo relógio): só o motivo; horário → ErrMaintenanceLocked. Trocar
// o motivo não reaplica nenhum efeito de início ou fim (isso é do tick).
func (s *MaintenanceService) Update(ctx context.Context, id string, in UpdateMaintenanceInput) (*domain.MaintenanceWindow, error) {
	timesSet := in.StartsAt != nil || in.EndsAt != nil
	if !timesSet && !in.ReasonSet {
		return nil, fmt.Errorf("%w: nothing to update", ErrValidation)
	}
	reason, err := normalizeReason("reason", in.Reason)
	if err != nil {
		return nil, err
	}

	w, err := s.windows.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w == nil || w.FinalizedAt != nil {
		return nil, ErrNotFound
	}
	if !in.ReasonSet {
		reason = w.Reason // mantém
	}
	now := s.clock().UTC().Truncate(time.Second)

	if !timesSet {
		ok, err := s.windows.UpdateReason(ctx, id, reason)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrNotFound
		}
		w.Reason = reason
		return w, nil
	}

	if w.Started(now) {
		return nil, ErrMaintenanceLocked
	}
	starts, ends := w.StartsAt, w.EndsAt
	if in.StartsAt != nil {
		starts = *in.StartsAt
	}
	if in.EndsAt != nil {
		ends = *in.EndsAt
	}
	starts, ends, err = validateWindowTimes(starts, ends, now)
	if err != nil {
		return nil, err
	}
	ok, err := s.windows.UpdateTimes(ctx, id, starts, ends, reason, now)
	if err != nil {
		return nil, err
	}
	if !ok {
		// a janela começou (ou sumiu) entre a leitura e o UPDATE
		cur, err := s.windows.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if cur == nil || cur.FinalizedAt != nil {
			return nil, ErrNotFound
		}
		return nil, ErrMaintenanceLocked
	}
	w.StartsAt, w.EndsAt, w.Reason = starts, ends, reason
	return w, nil
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
	_, err = s.windows.FinalizeGroup(ctx, []*domain.MaintenanceWindow{w}, now, &endsAt, nil)
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
