package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

func smOutboxEnv(t *testing.T, threshold int) (*sql.DB, *domain.Monitor) {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	now := time.Now().UTC()
	m := &domain.Monitor{
		ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "http://x.invalid", IntervalSeconds: 60, TimeoutSeconds: 5,
		FailureThreshold: threshold, Enabled: true, CurrentState: domain.StateUnknown, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.NewMonitorRepo(db).Create(context.Background(), m); err != nil {
		t.Fatal(err)
	}
	return db, m
}

func failResult() domain.CheckResult {
	msg := "connection refused"
	return domain.CheckResult{ErrorMessage: &msg}
}

// Queda com 3 canais: 3 linhas, mesmo event_id, snapshot sem segredo, prazo de
// 60 min contado do evento.
func TestProcessWritesOneRowPerChannelWithSameEventID(t *testing.T) {
	db, m := smOutboxEnv(t, 1)
	channels := []string{domain.ChannelWebhook, domain.ChannelEmail, domain.ChannelNtfy}
	woke := 0
	sm := service.NewStateMachine(newUoW(db), channels, func() { woke++ })
	if err := sm.Process(context.Background(), m, failResult()); err != nil {
		t.Fatal(err)
	}

	rows := pendingNotices(t, db, domain.EventDown)
	if len(rows) != 3 {
		t.Fatalf("linhas = %d, want 3", len(rows))
	}
	seen := map[string]bool{}
	for _, r := range rows {
		seen[r.Channel] = true
		if r.EventID != rows[0].EventID {
			t.Fatalf("event_id diferente entre canais: %s e %s", r.EventID, rows[0].EventID)
		}
		// o banco guarda o prazo em segundos inteiros
		if d := r.ExpiresAt.Sub(r.Snapshot.OccurredAt); d < 59*time.Minute+59*time.Second || d > 60*time.Minute {
			t.Fatalf("prazo de queda = %v, want 60 min", d)
		}
		if r.Snapshot.Incident == nil || r.Snapshot.Incident.LastError == nil || *r.Snapshot.Incident.LastError != "connection refused" {
			t.Fatalf("snapshot sem o erro: %+v", r.Snapshot.Incident)
		}
	}
	if !seen[domain.ChannelWebhook] || !seen[domain.ChannelEmail] || !seen[domain.ChannelNtfy] {
		t.Fatalf("canais gravados: %v", seen)
	}
	if woke != 1 {
		t.Fatalf("wake chamado %d vezes, want 1 (depois do commit)", woke)
	}
}

// Check sem mudança de estado não grava nada nem acorda o worker.
func TestProcessWithoutStateChangeWritesNothing(t *testing.T) {
	db, m := smOutboxEnv(t, 3)
	woke := 0
	sm := service.NewStateMachine(newUoW(db), testChannels, func() { woke++ })
	ctx := context.Background()
	code := 200
	for _, res := range []domain.CheckResult{{Success: true, StatusCode: &code}, failResult(), failResult(), {Success: true, StatusCode: &code}} {
		if err := sm.Process(ctx, m, res); err != nil {
			t.Fatal(err)
		}
	}
	if n := len(pendingNotices(t, db, "")); n != 0 {
		t.Fatalf("%d linhas no outbox, want 0", n)
	}
	if woke != 0 {
		t.Fatalf("wake chamado %d vezes, want 0", woke)
	}
}

// Rollback não deixa aviso: se o commit falha, a mudança de estado e o aviso somem juntos.
func TestProcessCommitFailureLeavesNoNotice(t *testing.T) {
	db, m := smOutboxEnv(t, 1)
	woke := 0
	sm := service.NewStateMachine(failingCommit{newUoW(db)}, testChannels, func() { woke++ })
	if err := sm.Process(context.Background(), m, failResult()); err == nil {
		t.Fatal("erro de commit esperado")
	}
	if n := len(pendingNotices(t, db, "")); n != 0 {
		t.Fatalf("%d linhas no outbox após rollback, want 0", n)
	}
	if woke != 0 {
		t.Fatal("acordou o worker sem commit")
	}
	var incidents int
	db.QueryRow(`SELECT COUNT(*) FROM incidents`).Scan(&incidents)
	if incidents != 0 {
		t.Fatalf("incidente gravado sem commit: %d", incidents)
	}
}

// Erro depois de enfileirar também desfaz o aviso.
type failingEnqueueUoW struct{ service.UnitOfWork }

func (u failingEnqueueUoW) OutboxRepo() service.OutboxRepositoryTx {
	return erroringOutbox{u.UnitOfWork.OutboxRepo()}
}

type erroringOutbox struct{ service.OutboxRepositoryTx }

func (o erroringOutbox) EnqueueWithTx(ctx context.Context, ev []domain.OutboxEvent) error {
	if err := o.OutboxRepositoryTx.EnqueueWithTx(ctx, ev); err != nil {
		return err
	}
	return errors.New("falha depois de gravar")
}

func TestProcessErrorAfterEnqueueRollsBack(t *testing.T) {
	db, m := smOutboxEnv(t, 1)
	sm := service.NewStateMachine(failingEnqueueUoW{newUoW(db)}, testChannels, nil)
	if err := sm.Process(context.Background(), m, failResult()); err == nil {
		t.Fatal("erro esperado")
	}
	if n := len(pendingNotices(t, db, "")); n != 0 {
		t.Fatalf("%d linhas no outbox, want 0", n)
	}
}

// Recuperação: snapshot com duração, prazo de 60 min contado do fim do incidente.
func TestProcessRecoveryNoticeSnapshot(t *testing.T) {
	db, m := smOutboxEnv(t, 1)
	t0 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	now := t0
	sm := service.NewStateMachine(newUoW(db), testChannels, nil).WithClock(func() time.Time { return now })
	ctx := context.Background()
	if err := sm.Process(ctx, m, failResult()); err != nil {
		t.Fatal(err)
	}
	now = t0.Add(90 * time.Second)
	code := 200
	if err := sm.Process(ctx, m, domain.CheckResult{Success: true, StatusCode: &code}); err != nil {
		t.Fatal(err)
	}
	up := pendingNotices(t, db, domain.EventUp)
	if len(up) != 1 {
		t.Fatalf("avisos de recuperação = %d, want 1", len(up))
	}
	inc := up[0].Snapshot.AsIncident()
	if inc.DurationSeconds == nil || *inc.DurationSeconds != 90 || inc.EndedAt == nil || !inc.EndedAt.Equal(now) {
		t.Fatalf("snapshot de recuperação: %+v", inc)
	}
	if !up[0].Snapshot.OccurredAt.Equal(now) || !up[0].ExpiresAt.Equal(now.Add(60*time.Minute)) {
		t.Fatalf("occurred_at=%v expires_at=%v", up[0].Snapshot.OccurredAt, up[0].ExpiresAt)
	}
}

// Concorrência (-race): vários monitores gravando checks e avisos enquanto o
// worker envia. Nada trava (conexão única) e todo aviso sai.
func TestProcessAndWorkerRunConcurrently(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	now := time.Now().UTC()
	const monitors = 5
	for i := range monitors {
		m := &domain.Monitor{
			ID: string(rune('a' + i)), Name: "m", URL: "http://x.invalid", IntervalSeconds: 60, TimeoutSeconds: 5,
			FailureThreshold: 1, Enabled: true, CurrentState: domain.StateUnknown, CreatedAt: now, UpdatedAt: now,
		}
		if err := repository.NewMonitorRepo(db).Create(ctx, m); err != nil {
			t.Fatal(err)
		}
	}

	sink := &countingChannel{}
	worker := service.NewOutboxWorker(repository.NewOutboxRepo(db), sink)
	wctx, cancel := context.WithCancel(ctx)
	worker.Start(wctx)

	var wg sync.WaitGroup
	for i := range monitors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			code := 200
			for range 10 { // cai, volta, cai, volta...
				m, err := repository.NewMonitorRepo(db).FindByID(ctx, string(rune('a'+i)))
				if err != nil {
					t.Error(err)
					return
				}
				sm := service.NewStateMachine(newUoW(db), testChannels, worker.Wake)
				if err := sm.Process(ctx, m, failResult()); err != nil {
					t.Error(err)
					return
				}
				m, _ = repository.NewMonitorRepo(db).FindByID(ctx, string(rune('a'+i)))
				if err := sm.Process(ctx, m, domain.CheckResult{Success: true, StatusCode: &code}); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(20 * time.Second):
		t.Fatal("Process travou com o worker rodando")
	}

	want := monitors * 10 * 2 // uma queda e uma recuperação por rodada
	deadline := time.Now().Add(10 * time.Second)
	for sink.count() < want && time.Now().Before(deadline) {
		worker.Wake()
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	waitCtx, stop := context.WithTimeout(ctx, 5*time.Second)
	defer stop()
	if err := worker.Wait(waitCtx); err != nil {
		t.Fatal(err)
	}
	if sink.count() != want {
		t.Fatalf("avisos enviados = %d, want %d", sink.count(), want)
	}
	var pending int
	db.QueryRow(`SELECT COUNT(*) FROM notification_outbox WHERE status <> 'sent'`).Scan(&pending)
	if pending != 0 {
		t.Fatalf("%d linhas não enviadas", pending)
	}
}

type countingChannel struct {
	mu sync.Mutex
	n  int
}

func (c *countingChannel) Name() string { return domain.ChannelWebhook }
func (c *countingChannel) Send(context.Context, domain.OutboxEvent) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
	return nil
}
func (c *countingChannel) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
