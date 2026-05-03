package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// IncidentRepoTx implementa IncidentRepositoryTx com suporte a transações.
type IncidentRepoTx struct {
	db *sql.DB
}

// NewIncidentRepoTx cria um repositório de incidents com suporte a transações.
func NewIncidentRepoTx(db *sql.DB) *IncidentRepoTx {
	return &IncidentRepoTx{db: db}
}

// CreateWithTx cria um novo incidente usando a transação fornecida.
func (r *IncidentRepoTx) CreateWithTx(ctx context.Context, tx *sql.Tx, i *domain.Incident) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO incidents (id, monitor_id, started_at, last_error, notification_sent, duration_seconds)
		VALUES (?,?,?,?,?,?)`,
		i.ID, i.MonitorID, i.StartedAt.UTC().Format(time.RFC3339),
		i.LastError, boolToInt(i.NotificationSent), i.DurationSeconds,
	)
	return err
}

// FindOpenByMonitor retorna o incidente aberto de um monitor (lê do DB principal).
func (r *IncidentRepoTx) FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, monitor_id, started_at, ended_at, last_error, notification_sent, duration_seconds
		FROM incidents WHERE monitor_id = ? AND ended_at IS NULL`, monitorID)
	return scanIncident(row)
}

// Close fecha um incidente (lê do DB principal).
func (r *IncidentRepoTx) Close(ctx context.Context, id string, endedAt string, durationSeconds int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE incidents SET ended_at = ?, duration_seconds = ? WHERE id = ?`,
		endedAt,
		durationSeconds,
		id,
	)
	return err
}