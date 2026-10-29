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

const incidentColumns = `id, monitor_id, started_at, ended_at, last_error, notification_sent, duration_seconds, failure_count`

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
	row := r.db.QueryRowContext(ctx, `
		SELECT `+incidentColumns+`
		FROM incidents WHERE monitor_id = ? AND ended_at IS NULL`, monitorID)
	return scanIncident(row)
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

// FindAll lista incidentes com os dados do monitor. O total vem de um COUNT(*)
// com o mesmo WHERE, independente de limit/offset.
func (r *IncidentRepo) FindAll(ctx context.Context, f domain.IncidentFilter) ([]*domain.Incident, int, error) {
	where := ` WHERE 1=1`
	args := []any{}
	if f.MonitorID != "" {
		where += ` AND i.monitor_id = ?`
		args = append(args, f.MonitorID)
	}
	order := ` ORDER BY i.started_at DESC`
	switch f.Status {
	case domain.IncidentStatusOpen:
		where += ` AND i.ended_at IS NULL`
	case domain.IncidentStatusResolved:
		where += ` AND i.ended_at IS NOT NULL`
		order = ` ORDER BY i.ended_at DESC` // a tela agrupa resolvidos pelo dia do fim
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM incidents i`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	limit = min(limit, 100)
	offset := max(f.Offset, 0)

	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, i.monitor_id, i.started_at, i.ended_at, i.last_error, i.notification_sent,
			i.duration_seconds, i.failure_count, m.name, m.type, m.url, m.config
		FROM incidents i LEFT JOIN monitors m ON m.id = i.monitor_id`+where+order+` LIMIT ? OFFSET ?`,
		append(args, limit, offset)...)
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

// DeleteResolvedOlderThan apaga incidentes resolvidos antes de `before`.
// Incidentes abertos nunca são apagados.
func (r *IncidentRepo) DeleteResolvedOlderThan(ctx context.Context, before string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM incidents WHERE ended_at IS NOT NULL AND ended_at < ?`, before)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

type incidentFields struct {
	notifSent       int
	startedAt       string
	endedAt         sql.NullString
	lastError       sql.NullString
	durationSeconds sql.NullInt64
}

func (f *incidentFields) dest(i *domain.Incident) []any {
	return []any{&i.ID, &i.MonitorID, &f.startedAt, &f.endedAt, &f.lastError, &f.notifSent, &f.durationSeconds, &i.FailureCount}
}

func (f *incidentFields) apply(i *domain.Incident) {
	i.StartedAt, _ = time.Parse(time.RFC3339, f.startedAt)
	i.NotificationSent = f.notifSent == 1
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
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+incidentColumns+`
		FROM incidents WHERE monitor_id = ? AND (ended_at IS NULL OR ended_at > ?)
		ORDER BY started_at`,
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
