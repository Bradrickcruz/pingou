package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type MonitorRepo struct{ db *sql.DB }

func NewMonitorRepo(db *sql.DB) *MonitorRepo {
	return &MonitorRepo{db: db}
}

func (r *MonitorRepo) Create(ctx context.Context, m *domain.Monitor) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO monitors
			(id, name, url, interval_seconds, timeout_seconds, failure_threshold,
			 enabled, current_state, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?)`,
		m.ID, m.Name, m.URL, m.IntervalSeconds, m.TimeoutSeconds, m.FailureThreshold,
		boolToInt(m.Enabled), string(m.CurrentState),
		m.CreatedAt.UTC().Format(time.RFC3339),
		m.UpdatedAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (r *MonitorRepo) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	row := r.db.QueryRowContext(ctx, `SELECT * FROM monitors WHERE id = ?`, id)
	return scanMonitor(row)
}

func (r *MonitorRepo) FindAll(ctx context.Context, f domain.MonitorFilter) ([]*domain.Monitor, int, error) {
	query := `SELECT * FROM monitors WHERE 1=1`
	args := []any{}

	if f.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, boolToInt(*f.Enabled))
	}
	if f.State != nil {
		query += ` AND current_state = ?`
		args = append(args, string(*f.State))
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, limit, f.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var monitors []*domain.Monitor
	for rows.Next() {
		m, err := scanMonitor(rows)
		if err != nil {
			return nil, 0, err
		}
		monitors = append(monitors, m)
	}

	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM monitors`).Scan(&total)

	return monitors, total, nil
}

func (r *MonitorRepo) Update(ctx context.Context, m *domain.Monitor) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE monitors SET
			name=?, url=?, interval_seconds=?, timeout_seconds=?,
			failure_threshold=?, enabled=?, current_state=?,
			last_checked_at=?, updated_at=?
		WHERE id=?`,
		m.Name, m.URL, m.IntervalSeconds, m.TimeoutSeconds,
		m.FailureThreshold, boolToInt(m.Enabled), string(m.CurrentState),
		nullableTime(m.LastCheckedAt), m.UpdatedAt.UTC().Format(time.RFC3339),
		m.ID,
	)
	return err
}

func (r *MonitorRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM monitors WHERE id = ?`, id)
	return err
}

func (r *MonitorRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM monitors`).Scan(&n)
	return n, err
}

// ── helpers ────────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

func scanMonitor(s scanner) (*domain.Monitor, error) {
	var m domain.Monitor
	var enabled int
	var state string
	var lastChecked, createdAt, updatedAt string

	err := s.Scan(
		&m.ID, &m.Name, &m.URL, &m.IntervalSeconds, &m.TimeoutSeconds,
		&m.FailureThreshold, &enabled, &state, &lastChecked, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Enabled = enabled == 1
	m.CurrentState = domain.MonitorState(state)
	m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	m.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if lastChecked != "" {
		t, _ := time.Parse(time.RFC3339, lastChecked)
		m.LastCheckedAt = &t
	}

	return &m, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}
