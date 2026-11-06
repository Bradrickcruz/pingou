package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// OutboxRepoTx grava avisos dentro da transação do check que os gerou: se a
// transação desfaz, o aviso não existe; se confirma, o aviso existe.
type OutboxRepoTx struct{}

func NewOutboxRepoTx() *OutboxRepoTx { return &OutboxRepoTx{} }

// EnqueueWithTx insere uma linha pending por evento.
func (r *OutboxRepoTx) EnqueueWithTx(ctx context.Context, tx *sql.Tx, events []domain.OutboxEvent) error {
	for _, ev := range events {
		payload, err := json.Marshal(ev.Snapshot)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO notification_outbox
				(event_id, event_type, channel, monitor_id, payload, status, attempts, created_at, expires_at, next_attempt_at)
			VALUES (?,?,?,?,?,'pending',0,?,?,?)`,
			ev.EventID, ev.EventType, ev.Channel, ev.MonitorID, string(payload),
			ev.CreatedAt.UTC().Format(time.RFC3339),
			ev.ExpiresAt.UTC().Format(time.RFC3339),
			ev.NextAttemptAt.UTC().Format(time.RFC3339),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
