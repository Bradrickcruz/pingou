package repository_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

var outboxT0 = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func newOutboxDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func outboxEvent(eventID, channel string) domain.OutboxEvent {
	m := &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://api.test"}
	inc := &domain.Incident{StartedAt: outboxT0}
	return domain.OutboxEvent{
		EventID: eventID, EventType: domain.EventDown, Channel: channel, MonitorID: "m1",
		Snapshot:  domain.NewDownSnapshot(m, inc),
		CreatedAt: outboxT0, ExpiresAt: outboxT0.Add(time.Hour), NextAttemptAt: outboxT0,
	}
}

func enqueue(t *testing.T, db *sql.DB, events ...domain.OutboxEvent) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.NewOutboxRepoTx().EnqueueWithTx(context.Background(), tx, events); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestOutboxEnqueueAndDue(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewOutboxRepo(db)

	later := outboxEvent("e2", domain.ChannelEmail)
	later.NextAttemptAt = outboxT0.Add(time.Minute)
	enqueue(t, db, outboxEvent("e1", domain.ChannelWebhook), later)

	got, err := repo.Due(ctx, outboxT0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].EventID != "e1" || got[0].Channel != domain.ChannelWebhook {
		t.Fatalf("vencidos em t0 = %+v, want só e1 (e2 é futuro)", got)
	}
	ev := got[0]
	if ev.Status != domain.OutboxPending || ev.Attempts != 0 || ev.Snapshot.Monitor.URL != "https://api.test" ||
		!ev.Snapshot.OccurredAt.Equal(outboxT0) || !ev.ExpiresAt.Equal(outboxT0.Add(time.Hour)) {
		t.Fatalf("linha lida: %+v", ev)
	}

	if got, _ = repo.Due(ctx, outboxT0.Add(time.Minute), 10); len(got) != 2 {
		t.Fatalf("vencidos em t0+1min = %d, want 2", len(got))
	}
	if got, _ = repo.Due(ctx, outboxT0.Add(time.Minute), 1); len(got) != 1 {
		t.Fatalf("limit 1 devolveu %d", len(got))
	}
}

func TestOutboxMarks(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewOutboxRepo(db)
	enqueue(t, db, outboxEvent("e1", domain.ChannelWebhook), outboxEvent("e2", domain.ChannelEmail), outboxEvent("e3", domain.ChannelNtfy))
	all, _ := repo.Due(ctx, outboxT0, 10)

	// retry: sai da rodada atual e volta quando next_attempt_at chegar; last_error truncado
	long := strings.Repeat("é", 500)
	if err := repo.MarkRetry(ctx, all[0].ID, 1, outboxT0.Add(30*time.Second), long); err != nil {
		t.Fatal(err)
	}
	if err := repo.MarkSent(ctx, all[1].ID, 1, outboxT0); err != nil {
		t.Fatal(err)
	}
	if err := repo.MarkExpired(ctx, all[2].ID, 3, outboxT0, "desistiu"); err != nil {
		t.Fatal(err)
	}

	if got, _ := repo.Due(ctx, outboxT0, 10); len(got) != 0 {
		t.Fatalf("sent/expired/retry futuro voltaram: %+v", got)
	}
	got, _ := repo.Due(ctx, outboxT0.Add(30*time.Second), 10)
	if len(got) != 1 || got[0].EventID != "e1" || got[0].Attempts != 1 {
		t.Fatalf("após o prazo do retry: %+v", got)
	}
	if got[0].LastError == nil || len([]rune(*got[0].LastError)) != 200 {
		t.Fatalf("last_error não foi truncado em 200 caracteres: %v", got[0].LastError)
	}
	if n, _ := repo.CountPending(ctx); n != 1 {
		t.Fatalf("pendentes = %d, want 1", n)
	}
}

func TestOutboxPurgeDoneBefore(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewOutboxRepo(db)
	enqueue(t, db, outboxEvent("old-sent", "a"), outboxEvent("old-expired", "b"), outboxEvent("new-sent", "c"), outboxEvent("old-pending", "d"))
	all, _ := repo.Due(ctx, outboxT0, 10)
	old, fresh := outboxT0.Add(-10*24*time.Hour), outboxT0
	repo.MarkSent(ctx, all[0].ID, 1, old)
	repo.MarkExpired(ctx, all[1].ID, 5, old, "x")
	repo.MarkSent(ctx, all[2].ID, 1, fresh)

	n, err := repo.PurgeDoneBefore(ctx, outboxT0.Add(-7*24*time.Hour))
	if err != nil || n != 2 {
		t.Fatalf("apagou %d (err=%v), want 2", n, err)
	}
	var left int
	db.QueryRow(`SELECT COUNT(*) FROM notification_outbox`).Scan(&left)
	if left != 2 { // new-sent e old-pending
		t.Fatalf("sobraram %d linhas, want 2 (a recente e a pendente)", left)
	}
}

// Linha com payload ilegível não pode ocupar o lote para sempre.
func TestOutboxDueExpiresUnreadableRow(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewOutboxRepo(db)
	enqueue(t, db, outboxEvent("good", domain.ChannelWebhook))
	ts := outboxT0.Format(time.RFC3339)
	if _, err := db.Exec(`INSERT INTO notification_outbox (event_id,event_type,channel,monitor_id,payload,created_at,expires_at,next_attempt_at)
		VALUES ('bad','down','webhook','m1','não é json',?,?,?)`, ts, ts, ts); err != nil {
		t.Fatal(err)
	}
	got, err := repo.Due(ctx, outboxT0, 10)
	if err != nil || len(got) != 1 || got[0].EventID != "good" {
		t.Fatalf("Due = %+v, err=%v, want só a linha boa", got, err)
	}
	var status string
	db.QueryRow(`SELECT status FROM notification_outbox WHERE event_id='bad'`).Scan(&status)
	if status != "expired" {
		t.Fatalf("linha ilegível: status=%s, want expired", status)
	}
}

// O snapshot nunca leva segredo: monitor heartbeat com ping_token não vaza.
func TestOutboxSnapshotHasNoSecrets(t *testing.T) {
	db := newOutboxDB(t)
	token := "segredo-do-ping"
	m := &domain.Monitor{ID: "m1", Name: "Backup", Type: domain.TypeHeartbeat, PingToken: &token}
	ev := outboxEvent("e1", domain.ChannelWebhook)
	ev.Snapshot = domain.NewDownSnapshot(m, &domain.Incident{StartedAt: outboxT0})
	enqueue(t, db, ev)
	var payload string
	db.QueryRow(`SELECT payload FROM notification_outbox`).Scan(&payload)
	if strings.Contains(payload, token) || strings.Contains(payload, "ping_token") {
		t.Fatalf("payload vazou o token: %s", payload)
	}
}
