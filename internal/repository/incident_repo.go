package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type IncidentRepo struct{ db *sql.DB }

func NewIncidentRepo(db *sql.DB) *IncidentRepo { return &IncidentRepo{db: db} }

const incidentColumns = `id, monitor_id, started_at, ended_at, last_error, notification_sent, duration_seconds, failure_count, closed_manually, close_note`

func (r *IncidentRepo) Create(ctx context.Context, i *domain.Incident) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO incidents (id, monitor_id, started_at, last_error, notification_sent, duration_seconds, failure_count)
		VALUES (?,?,?,?,?,?,?)`,
		i.ID, i.MonitorID, i.StartedAt.UTC().Format(time.RFC3339),
		i.LastError, boolToInt(i.NotificationSent), i.DurationSeconds, i.FailureCount,
	)
	return err
}

func (r *IncidentRepo) FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error) {
	row := r.db.QueryRowContext(ctx, sqlIncidentOpenByMonitor, monitorID)
	return scanIncident(row)
}

// FindOpenByMonitors devolve o incidente aberto de cada monitor em uma query,
// indexado pelo monitor_id. Monitores sem incidente aberto ficam fora do mapa.
func (r *IncidentRepo) FindOpenByMonitors(ctx context.Context, monitorIDs []string) (map[string]*domain.Incident, error) {
	out := make(map[string]*domain.Incident, len(monitorIDs))
	if len(monitorIDs) == 0 {
		return out, nil
	}
	args := make([]any, len(monitorIDs))
	for i, id := range monitorIDs {
		args[i] = id
	}
	rows, err := r.db.QueryContext(ctx, sqlIncidentsOpenByMonitors(len(monitorIDs)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		i, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		out[i.MonitorID] = i
	}
	return out, rows.Err()
}

func (r *IncidentRepo) Close(ctx context.Context, id string, endedAt string, durationSeconds int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE incidents SET ended_at = ?, duration_seconds = ? WHERE id = ?`,
		endedAt,
		durationSeconds,
		id,
	)
	return err
}

// CloseManually fecha o incidente aberto do monitor; no-op se não há nenhum.
// reason vai para closed_manually (1 = pausa, 2 = manutenção) e note para close_note.
func (r *IncidentRepo) CloseManually(ctx context.Context, monitorID string, endedAt time.Time, reason domain.CloseReason, note *string) error {
	end := endedAt.UTC().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, sqlIncidentCloseManually, end, int(reason), note, end, monitorID)
	return err
}

// FindAll lista incidentes com os dados do monitor. O total vem de um COUNT(*)
// com o mesmo WHERE, independente de limit/offset.
func (r *IncidentRepo) FindAll(ctx context.Context, f domain.IncidentFilter) ([]*domain.Incident, int, error) {
	countSQL, listSQL, args := incidentsListSQL(f)

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	limit = min(limit, 100)
	offset := max(f.Offset, 0)

	rows, err := r.db.QueryContext(ctx, listSQL, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	incidents := []*domain.Incident{}
	for rows.Next() {
		i, err := scanIncidentWithMonitor(rows)
		if err != nil {
			return nil, 0, err
		}
		incidents = append(incidents, i)
	}
	return incidents, total, rows.Err()
}

type incidentFields struct {
	notifSent       int
	closedManually  int
	closeNote       sql.NullString
	startedAt       string
	endedAt         sql.NullString
	lastError       sql.NullString
	durationSeconds sql.NullInt64
}

func (f *incidentFields) dest(i *domain.Incident) []any {
	return []any{&i.ID, &i.MonitorID, &f.startedAt, &f.endedAt, &f.lastError, &f.notifSent, &f.durationSeconds, &i.FailureCount, &f.closedManually, &f.closeNote}
}

func (f *incidentFields) apply(i *domain.Incident) {
	i.StartedAt, _ = time.Parse(time.RFC3339, f.startedAt)
	i.NotificationSent = f.notifSent == 1
	i.CloseReason = domain.CloseReason(f.closedManually)
	if f.closeNote.Valid {
		i.CloseNote = &f.closeNote.String
	}
	i.EndedAt = parseNullTime(f.endedAt)
	if f.lastError.Valid {
		i.LastError = &f.lastError.String
	}
	if f.durationSeconds.Valid {
		i.DurationSeconds = &f.durationSeconds.Int64
	}
}

func scanIncident(s scanner) (*domain.Incident, error) {
	var i domain.Incident
	var f incidentFields
	err := s.Scan(f.dest(&i)...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	f.apply(&i)
	return &i, nil
}

func scanIncidentWithMonitor(s scanner) (*domain.Incident, error) {
	var i domain.Incident
	var f incidentFields
	var name, typ, url, config sql.NullString
	if err := s.Scan(append(f.dest(&i), &name, &typ, &url, &config)...); err != nil {
		return nil, err
	}
	f.apply(&i)

	if name.Valid { // LEFT JOIN: sem monitor, os campos ficam vazios
		m := domain.Monitor{Type: domain.MonitorType(typ.String), URL: url.String}
		var cfg storedConfig
		_ = json.Unmarshal([]byte(config.String), &cfg)
		m.Config = cfg.toDomain()
		i.MonitorName, i.MonitorType, i.MonitorTarget = name.String, m.Type, m.Target()
	}
	return &i, nil
}

// FindInWindow lista os incidentes do monitor ainda abertos ou encerrados
// depois de `since`, do mais antigo ao mais recente.
// ponytail: sem limite; um monitor com milhares de quedas em 30 d é patológico, pôr LIMIT 1000 se aparecer
func (r *IncidentRepo) FindInWindow(ctx context.Context, monitorID string, since time.Time) ([]*domain.Incident, error) {
	rows, err := r.db.QueryContext(ctx, sqlIncidentsInWindow,
		monitorID, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.Incident
	for rows.Next() {
		i, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}
