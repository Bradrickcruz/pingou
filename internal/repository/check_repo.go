package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type CheckRepo struct{ db *sql.DB }

func NewCheckRepo(db *sql.DB) *CheckRepo { return &CheckRepo{db: db} }

func (r *CheckRepo) Create(ctx context.Context, c *domain.Check) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO checks (monitor_id, success, status_code, latency_ms, error_message, checked_at)
		VALUES (?,?,?,?,?,?)`,
		c.MonitorID, boolToInt(c.Success), c.StatusCode, c.LatencyMs,
		c.ErrorMessage, c.CheckedAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (r *CheckRepo) FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, monitor_id, success, status_code, latency_ms, error_message, checked_at
		FROM checks WHERE monitor_id = ?
		ORDER BY checked_at DESC LIMIT ? OFFSET ?`,
		monitorID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var checks []*domain.Check
	for rows.Next() {
		var c domain.Check
		var success int
		var checkedAt string
		if err := rows.Scan(&c.ID, &c.MonitorID, &success, &c.StatusCode, &c.LatencyMs, &c.ErrorMessage, &checkedAt); err != nil {
			return nil, 0, err
		}
		c.Success = success == 1
		c.CheckedAt, _ = time.Parse(time.RFC3339, checkedAt)
		checks = append(checks, &c)
	}

	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM checks WHERE monitor_id = ?`, monitorID).Scan(&total)

	return checks, total, nil
}

func (r *CheckRepo) DeleteOlderThan(ctx context.Context, before string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM checks WHERE checked_at < ?`, before,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
