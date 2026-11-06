package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type MaintenanceRepo struct {
	db     *sql.DB
	outbox OutboxRepoTx
}

func NewMaintenanceRepo(db *sql.DB) *MaintenanceRepo { return &MaintenanceRepo{db: db} }

const maintenanceColumns = `id, monitor_id, starts_at, ends_at, reason, started_at, finalized_at, created_at`

func fmtTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func (r *MaintenanceRepo) CreateBatch(ctx context.Context, windows []*domain.MaintenanceWindow) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, w := range windows {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO maintenance_windows (id, monitor_id, starts_at, ends_at, reason, created_at)
			VALUES (?,?,?,?,?,?)`,
			w.ID, w.MonitorID, fmtTime(w.StartsAt), fmtTime(w.EndsAt), w.Reason, fmtTime(w.CreatedAt),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *MaintenanceRepo) FindByID(ctx context.Context, id string) (*domain.MaintenanceWindow, error) {
	return scanMaintenance(r.db.QueryRowContext(ctx, `SELECT `+maintenanceColumns+` FROM maintenance_windows WHERE id = ?`, id))
}

func (r *MaintenanceRepo) FindLiveByMonitor(ctx context.Context, monitorID string) (*domain.MaintenanceWindow, error) {
	return scanMaintenance(r.db.QueryRowContext(ctx, sqlMaintenanceLiveByMonitor, monitorID))
}

func (r *MaintenanceRepo) FindLiveByMonitors(ctx context.Context, monitorIDs []string) (map[string]*domain.MaintenanceWindow, error) {
	out := make(map[string]*domain.MaintenanceWindow, len(monitorIDs))
	if len(monitorIDs) == 0 {
		return out, nil
	}
	args := make([]any, len(monitorIDs))
	for i, id := range monitorIDs {
		args[i] = id
	}
	windows, err := r.list(ctx, sqlMaintenanceLiveByMonitors(len(monitorIDs)), args...)
	if err != nil {
		return nil, err
	}
	for _, w := range windows {
		out[w.MonitorID] = w
	}
	return out, nil
}

func (r *MaintenanceRepo) ActiveForMonitor(ctx context.Context, monitorID string, now time.Time) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, sqlMaintenanceActive, monitorID, fmtTime(now)).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *MaintenanceRepo) ListLive(ctx context.Context) ([]*domain.MaintenanceWindowView, error) {
	rows, err := r.db.QueryContext(ctx, sqlMaintenanceListLive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.MaintenanceWindowView{}
	for rows.Next() {
		var v domain.MaintenanceWindowView
		w, err := scanMaintenance(nameScanner{rows, &v.MonitorName})
		if err != nil {
			return nil, err
		}
		v.MaintenanceWindow = *w
		out = append(out, &v)
	}
	return out, rows.Err()
}

// nameScanner lê as colunas da janela e, por último, o nome do monitor.
type nameScanner struct {
	s    scanner
	name *string
}

func (n nameScanner) Scan(dest ...any) error { return n.s.Scan(append(dest, n.name)...) }

func (r *MaintenanceRepo) UpdateTimes(ctx context.Context, id string, startsAt, endsAt time.Time, reason *string, now time.Time) (bool, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE maintenance_windows SET starts_at = ?, ends_at = ?, reason = ?
		WHERE id = ? AND finalized_at IS NULL AND starts_at > ?`,
		fmtTime(startsAt), fmtTime(endsAt), reason, id, fmtTime(now))
	return affected(res, err)
}

func (r *MaintenanceRepo) UpdateReason(ctx context.Context, id string, reason *string) (bool, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE maintenance_windows SET reason = ?
		WHERE id = ? AND finalized_at IS NULL`, reason, id)
	return affected(res, err)
}

func affected(res sql.Result, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func (r *MaintenanceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM maintenance_windows WHERE id = ?`, id)
	return err
}

func (r *MaintenanceRepo) DueStart(ctx context.Context, now time.Time) ([]*domain.MaintenanceWindow, error) {
	return r.list(ctx, `SELECT `+maintenanceColumns+` FROM maintenance_windows
		WHERE starts_at <= ? AND started_at IS NULL AND finalized_at IS NULL ORDER BY starts_at`, fmtTime(now))
}

func (r *MaintenanceRepo) DueEnd(ctx context.Context, now time.Time) ([]*domain.MaintenanceWindow, error) {
	return r.list(ctx, `SELECT `+maintenanceColumns+` FROM maintenance_windows
		WHERE ends_at <= ? AND finalized_at IS NULL ORDER BY ends_at`, fmtTime(now))
}

// StartGroup: ver domain.MaintenanceRepository.
func (r *MaintenanceRepo) StartGroup(ctx context.Context, windows []*domain.MaintenanceWindow, now time.Time, build domain.MaintenanceNoticeBuilder) ([]*domain.MaintenanceWindow, error) {
	if len(windows) == 0 {
		return nil, nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var applied []*domain.MaintenanceWindow
	for _, w := range windows {
		ok, err := startWindowTx(ctx, tx, w, now)
		if err != nil {
			return nil, err
		}
		if ok {
			applied = append(applied, w)
		}
	}
	if len(applied) == 0 {
		return nil, nil
	}
	if err := r.enqueueNotices(ctx, tx, applied, build); err != nil {
		return nil, err
	}
	return applied, tx.Commit()
}

// FinalizeGroup: ver domain.MaintenanceRepository.
func (r *MaintenanceRepo) FinalizeGroup(ctx context.Context, windows []*domain.MaintenanceWindow, now time.Time, endsAt *time.Time, build domain.MaintenanceNoticeBuilder) ([]*domain.MaintenanceWindow, error) {
	if len(windows) == 0 {
		return nil, nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var applied, notice []*domain.MaintenanceWindow
	for _, w := range windows {
		hadStarted, ok, err := finalizeWindowTx(ctx, tx, w, now, endsAt)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		done := *w
		if endsAt != nil {
			done.EndsAt = endsAt.UTC()
		}
		applied = append(applied, &done)
		if hadStarted {
			notice = append(notice, &done)
		}
	}
	if len(applied) == 0 {
		return nil, nil
	}
	if err := r.enqueueNotices(ctx, tx, notice, build); err != nil {
		return nil, err
	}
	return applied, tx.Commit()
}

// enqueueNotices lê os monitores das janelas pelo tx (ler pelo db travaria
// esperando a própria conexão) e grava no outbox as linhas que build devolver.
func (r *MaintenanceRepo) enqueueNotices(ctx context.Context, tx *sql.Tx, windows []*domain.MaintenanceWindow, build domain.MaintenanceNoticeBuilder) error {
	if build == nil || len(windows) == 0 {
		return nil
	}
	monitors := make([]*domain.Monitor, 0, len(windows))
	for _, w := range windows {
		m, err := scanMonitor(tx.QueryRowContext(ctx, `SELECT `+monitorColumns+` FROM monitors WHERE id = ?`, w.MonitorID))
		if err != nil {
			return err
		}
		if m != nil {
			monitors = append(monitors, m)
		}
	}
	events, err := build(windows, monitors)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	return r.outbox.EnqueueWithTx(ctx, tx, events)
}

// startWindowTx marca started_at e só então fecha o incidente e zera o estado:
// se outra chamada já aplicou o início, nada muda. A guarda starts_at <= now
// cobre a janela adiada entre a leitura do tick e esta transação: vale o
// horário novo.
func startWindowTx(ctx context.Context, tx *sql.Tx, w *domain.MaintenanceWindow, now time.Time) (bool, error) {
	ts := fmtTime(now)
	res, err := tx.ExecContext(ctx, `UPDATE maintenance_windows SET started_at = ?
		WHERE id = ? AND started_at IS NULL AND finalized_at IS NULL AND starts_at <= ?`, ts, w.ID, ts)
	if err != nil {
		return false, err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, sqlIncidentCloseManually,
		ts, int(domain.CloseMaintenance), w.Reason, ts, w.MonitorID); err != nil {
		return false, err
	}
	// UNKNOWN: ao voltar, handleFailure abre incidente novo (só abre se prevState != DOWN)
	if _, err := tx.ExecContext(ctx, `UPDATE monitors SET current_state = 'UNKNOWN' WHERE id = ?`, w.MonitorID); err != nil {
		return false, err
	}
	return true, nil
}

// finalizeWindowTx aplica o fim de uma janela. hadStarted diz se o início já
// tinha sido aplicado antes desta transação; ok, se o fim foi aplicado agora.
// Tudo é lido pelo tx, e com uma conexão só nada muda entre a leitura e a escrita.
func finalizeWindowTx(ctx context.Context, tx *sql.Tx, w *domain.MaintenanceWindow, now time.Time, endsAt *time.Time) (hadStarted, ok bool, err error) {
	ts := fmtTime(now)
	var started bool
	var curEnds string
	err = tx.QueryRowContext(ctx, `SELECT started_at IS NOT NULL, ends_at FROM maintenance_windows
		WHERE id = ? AND finalized_at IS NULL`, w.ID).Scan(&started, &curEnds)
	if err == sql.ErrNoRows {
		return false, false, nil // apagada ou já finalizada
	}
	if err != nil {
		return false, false, err
	}
	// janela adiada depois da leitura do tick: o fim ainda não chegou
	if endsAt == nil && curEnds > ts {
		return false, false, nil
	}

	// janela que começou e terminou sem o tick do início (Pingou desligado):
	// o início roda antes, para o incidente aberto fechar como manutenção
	if _, err := startWindowTx(ctx, tx, w, now); err != nil {
		return false, false, err
	}

	var ends *string
	if endsAt != nil {
		s := fmtTime(*endsAt)
		ends = &s
	}
	res, err := tx.ExecContext(ctx, `UPDATE maintenance_windows SET finalized_at = ?, ends_at = COALESCE(?, ends_at)
		WHERE id = ? AND finalized_at IS NULL`, ts, ends, w.ID)
	if err != nil {
		return false, false, err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return false, false, err
	}
	// resumed_at: só falhas a partir de agora contam. updated_at e last_ping_at
	// são a base do prazo do heartbeat: recomeça agora, sem DOWN imediato.
	if _, err := tx.ExecContext(ctx, `UPDATE monitors SET current_state = 'UNKNOWN', resumed_at = ?, updated_at = ?,
			last_ping_at = CASE WHEN last_ping_at IS NOT NULL THEN ? END
		WHERE id = ?`, ts, ts, ts, w.MonitorID); err != nil {
		return false, false, err
	}
	return started, true, nil
}

func (r *MaintenanceRepo) list(ctx context.Context, query string, args ...any) ([]*domain.MaintenanceWindow, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.MaintenanceWindow
	for rows.Next() {
		w, err := scanMaintenance(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func scanMaintenance(s scanner) (*domain.MaintenanceWindow, error) {
	var w domain.MaintenanceWindow
	var startsAt, endsAt, createdAt string
	var reason, startedAt, finalizedAt sql.NullString
	err := s.Scan(&w.ID, &w.MonitorID, &startsAt, &endsAt, &reason, &startedAt, &finalizedAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	w.StartsAt, _ = time.Parse(time.RFC3339, startsAt)
	w.EndsAt, _ = time.Parse(time.RFC3339, endsAt)
	w.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if reason.Valid {
		w.Reason = &reason.String
	}
	w.StartedAt = parseNullTime(startedAt)
	w.FinalizedAt = parseNullTime(finalizedAt)
	return &w, nil
}
