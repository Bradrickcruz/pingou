package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// MonitorRepoTx implementa MonitorRepositoryTx com suporte a transações.
type MonitorRepoTx struct {
	db *sql.DB
}

// NewMonitorRepoTx cria um repositório de monitors com suporte a transações.
func NewMonitorRepoTx(db *sql.DB) *MonitorRepoTx {
	return &MonitorRepoTx{db: db}
}

// UpdateWithTx atualiza um monitor usando a transação fornecida.
func (r *MonitorRepoTx) UpdateWithTx(ctx context.Context, tx *sql.Tx, m *domain.Monitor) error {
	_, err := tx.ExecContext(ctx, `
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

// FindByID retorna um monitor pelo ID (lê do DB principal, não da transação).
func (r *MonitorRepoTx) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	row := r.db.QueryRowContext(ctx, `SELECT * FROM monitors WHERE id = ?`, id)
	return scanMonitor(row)
}