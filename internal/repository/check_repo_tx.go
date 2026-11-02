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

// Latest retorna os n checks mais recentes do monitor usando a transação
// (SQLite tem 1 conexão: ler via r.db com tx aberta trava). Sem o COUNT(*) do
// total: a state machine chama a cada check e só quer as últimas linhas.
func (r *CheckRepoTx) Latest(ctx context.Context, tx *sql.Tx, monitorID string, n int) ([]*domain.Check, error) {
	rows, err := tx.QueryContext(ctx, sqlChecksPage, monitorID, n, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []*domain.Check
	for rows.Next() {
		var c domain.Check
		var success int
		var checkedAt string
		if err := rows.Scan(&c.ID, &c.MonitorID, &success, &c.StatusCode, &c.LatencyMs, &c.ErrorMessage, &checkedAt); err != nil {
			return nil, err
		}
		c.Success = success == 1
		c.CheckedAt, _ = time.Parse(time.RFC3339, checkedAt)
		checks = append(checks, &c)
	}
	return checks, rows.Err()
}
