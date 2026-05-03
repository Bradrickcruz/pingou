package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// CheckRepoTx implementa CheckRepositoryTx com suporte a transações.
type CheckRepoTx struct {
	db *sql.DB
}

// NewCheckRepoTx cria um repositório de checks com suporte a transações.
func NewCheckRepoTx(db *sql.DB) *CheckRepoTx {
	return &CheckRepoTx{db: db}
}

// CreateWithTx insere um novo check usando a transação fornecida.
func (r *CheckRepoTx) CreateWithTx(ctx context.Context, tx *sql.Tx, c *domain.Check) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO checks (monitor_id, success, status_code, latency_ms, error_message, checked_at)
		VALUES (?,?,?,?,?,?)`,
		c.MonitorID, boolToInt(c.Success), c.StatusCode, c.LatencyMs,
		c.ErrorMessage, c.CheckedAt.UTC().Format(time.RFC3339),
	)
	return err
}

// FindByMonitor retorna checks de um monitor (lê do DB principal, não da transação).
func (r *CheckRepoTx) FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error) {
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