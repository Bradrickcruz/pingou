package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// maxOutboxError é o tamanho máximo de last_error (em caracteres).
const maxOutboxError = 200

// OutboxRepo é o lado do worker: lê os avisos vencidos e grava o resultado de
// cada tentativa. Cada chamada é curta e solta a conexão ao voltar; o envio de
// rede acontece entre elas, nunca com a conexão presa (MaxOpenConns=1).
type OutboxRepo struct {
	db *sql.DB
}

func NewOutboxRepo(db *sql.DB) *OutboxRepo { return &OutboxRepo{db: db} }

// Due devolve até limit avisos pending com next_attempt_at <= now, os mais
// antigos primeiro. Linha com payload ilegível é marcada expired (senão
// ocuparia o lote para sempre) e não é devolvida.
func (r *OutboxRepo) Due(ctx context.Context, now time.Time, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, sqlOutboxDue, now.UTC().Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}

	var out []domain.OutboxEvent
	var corrupt []int64
	for rows.Next() {
		ev, err := scanOutbox(rows)
		if err != nil {
			var id int64
			if ev != nil {
				id = ev.ID
			}
			if id == 0 {
				_ = rows.Close()
				return nil, err
			}
			corrupt = append(corrupt, id)
			continue
		}
		out = append(out, *ev)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	// fecha antes de gravar: com uma conexão só, escrever com o cursor aberto trava
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for _, id := range corrupt {
		slog.Error("outbox row unreadable, expiring", "outbox_id", id)
		if err := r.MarkExpired(ctx, id, 0, now, "unreadable payload"); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MarkSent encerra a linha como enviada.
func (r *OutboxRepo) MarkSent(ctx context.Context, id int64, attempts int, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_outbox
		SET status = 'sent', attempts = ?, last_error = NULL, done_at = ?
		WHERE id = ?`, attempts, at.UTC().Format(time.RFC3339), id)
	return err
}

// MarkRetry agenda a próxima tentativa.
func (r *OutboxRepo) MarkRetry(ctx context.Context, id int64, attempts int, next time.Time, lastError string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_outbox
		SET attempts = ?, next_attempt_at = ?, last_error = ?
		WHERE id = ?`, attempts, next.UTC().Format(time.RFC3339), truncateError(lastError), id)
	return err
}

// MarkExpired desiste do aviso.
func (r *OutboxRepo) MarkExpired(ctx context.Context, id int64, attempts int, at time.Time, lastError string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_outbox
		SET status = 'expired', attempts = ?, last_error = ?, done_at = ?
		WHERE id = ?`, attempts, truncateError(lastError), at.UTC().Format(time.RFC3339), id)
	return err
}

// CountPending devolve quantos avisos aguardam envio.
func (r *OutboxRepo) CountPending(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notification_outbox WHERE status = 'pending'`).Scan(&n)
	return n, err
}

// PurgeDoneBefore apaga linhas sent/expired encerradas antes de before.
func (r *OutboxRepo) PurgeDoneBefore(ctx context.Context, before time.Time) (int64, error) {
	res, err := r.db.ExecContext(ctx, sqlOutboxPurge, before.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func truncateError(s string) string {
	if r := []rune(s); len(r) > maxOutboxError {
		return string(r[:maxOutboxError])
	}
	return s
}

// scanOutbox lê uma linha. Se só o payload estiver ruim, devolve o evento
// parcial (com ID) junto do erro, para o chamador poder descartá-lo.
func scanOutbox(row interface{ Scan(...any) error }) (*domain.OutboxEvent, error) {
	var (
		ev                           domain.OutboxEvent
		payload, status              string
		createdAt, expiresAt, nextAt string
		lastError, doneAt            sql.NullString
	)
	if err := row.Scan(&ev.ID, &ev.EventID, &ev.EventType, &ev.Channel, &ev.MonitorID, &payload, &status,
		&ev.Attempts, &createdAt, &expiresAt, &nextAt, &lastError, &doneAt); err != nil {
		return nil, err
	}
	ev.Status = domain.OutboxStatus(status)
	if lastError.Valid {
		ev.LastError = &lastError.String
	}

	var err error
	if ev.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return &ev, fmt.Errorf("created_at: %w", err)
	}
	if ev.ExpiresAt, err = time.Parse(time.RFC3339, expiresAt); err != nil {
		return &ev, fmt.Errorf("expires_at: %w", err)
	}
	if ev.NextAttemptAt, err = time.Parse(time.RFC3339, nextAt); err != nil {
		return &ev, fmt.Errorf("next_attempt_at: %w", err)
	}
	if doneAt.Valid {
		t, err := time.Parse(time.RFC3339, doneAt.String)
		if err != nil {
			return &ev, fmt.Errorf("done_at: %w", err)
		}
		ev.DoneAt = &t
	}
	if err := json.Unmarshal([]byte(payload), &ev.Snapshot); err != nil {
		return &ev, fmt.Errorf("payload: %w", err)
	}
	return &ev, nil
}
