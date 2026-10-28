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

// UpdateWithTx grava o resultado da state machine usando a transação fornecida.
// Só toca nos campos de estado: assim um check não desfaz uma edição feita
// pela API enquanto ele rodava (nome, intervalo, token etc.).
func (r *MonitorRepoTx) UpdateWithTx(ctx context.Context, tx *sql.Tx, m *domain.Monitor) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE monitors SET
			current_state=?, last_checked_at=?, last_ping_at=?, updated_at=?
		WHERE id=?`,
		string(m.CurrentState), nullableTime(m.LastCheckedAt), nullableTime(m.LastPingAt),
		m.UpdatedAt.UTC().Format(time.RFC3339),
		m.ID,
	)
	return err
}

// FindByID retorna um monitor pelo ID (lê do DB principal, não da transação).
func (r *MonitorRepoTx) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+monitorColumns+` FROM monitors WHERE id = ?`, id)
	return scanMonitor(row)
}
