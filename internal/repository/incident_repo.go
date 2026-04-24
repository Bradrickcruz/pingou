package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type IncidentRepo struct{ db *sql.DB }

func NewIncidentRepo(db *sql.DB) *IncidentRepo { return &IncidentRepo{db: db} }

func (r *IncidentRepo) Create(ctx context.Context, i *domain.Incident) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO incidents (id, monitor_id, started_at, last_error, notification_sent)
		VALUES (?,?,?,?,?)`,
		i.ID, i.MonitorID, i.StartedAt.UTC().Format(time.RFC3339),
		i.LastError, boolToInt(i.NotificationSent),
	)
	return err
}

func (r *IncidentRepo) FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, monitor_id, started_at, ended_at, last_error, notification_sent
		FROM incidents WHERE monitor_id = ? AND ended_at IS NULL`, monitorID)
	return scanIncident(row)
}

func (r *IncidentRepo) Close(ctx context.Context, id string, endedAt string, _ int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE incidents SET ended_at = ? WHERE id = ?`, endedAt, id)
	return err
}

func (r *IncidentRepo) FindByMonitor(ctx context.Context, monitorID string, onlyOpen bool, limit, offset int) ([]*domain.Incident, int, error) {
	q := `SELECT id, monitor_id, started_at, ended_at, last_error, notification_sent FROM incidents WHERE monitor_id = ?`
	args := []any{monitorID}
	if onlyOpen {
		q += ` AND ended_at IS NULL`
	}
	if limit <= 0 {
		limit = 20
	}
	q += fmt.Sprintf(` ORDER BY started_at DESC LIMIT %d OFFSET %d`, limit, offset)

	return r.queryIncidents(ctx, q, args...)
}

func (r *IncidentRepo) FindAll(ctx context.Context, onlyOpen bool, limit, offset int) ([]*domain.Incident, int, error) {
	q := `SELECT id, monitor_id, started_at, ended_at, last_error, notification_sent FROM incidents WHERE 1=1`
	if onlyOpen {
		q += ` AND ended_at IS NULL`
	}
	if limit <= 0 {
		limit = 20
	}
	q += fmt.Sprintf(` ORDER BY started_at DESC LIMIT %d OFFSET %d`, limit, offset)

	return r.queryIncidents(ctx, q)
}

func (r *IncidentRepo) queryIncidents(ctx context.Context, q string, args ...any) ([]*domain.Incident, int, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var incidents []*domain.Incident
	for rows.Next() {
		i, err := scanIncident(rows)
		if err != nil {
			return nil, 0, err
		}
		incidents = append(incidents, i)
	}
	return incidents, len(incidents), nil
}

func scanIncident(s scanner) (*domain.Incident, error) {
	var i domain.Incident
	var notifSent int
	var startedAt string
	var endedAt sql.NullString
	var lastError sql.NullString

	err := s.Scan(&i.ID, &i.MonitorID, &startedAt, &endedAt, &lastError, &notifSent)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	i.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
	i.NotificationSent = notifSent == 1
	if endedAt.Valid {
		t, _ := time.Parse(time.RFC3339, endedAt.String)
		i.EndedAt = &t
	}
	if lastError.Valid {
		i.LastError = &lastError.String
	}

	return &i, nil
}
