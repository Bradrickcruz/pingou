package service_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// Regressão: com MaxOpenConns=1, ler fora da transação travava o caminho de falha
// e o monitor nunca saía de UNKNOWN.
func TestFailingMonitorGoesDown(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	m := &domain.Monitor{
		ID: "m1", Name: "x", URL: "http://x.invalid", IntervalSeconds: 60, TimeoutSeconds: 5,
		FailureThreshold: 1, Enabled: true, CurrentState: domain.StateUnknown, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.NewMonitorRepo(db).Create(ctx, m); err != nil {
		t.Fatal(err)
	}

	uow := newUoW(db)
	msg := "boom"
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := service.NewStateMachine(uow, testChannels, nil).Process(cctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
		t.Fatal(err)
	}
	if m.CurrentState != domain.StateDown {
		t.Fatalf("state = %s, want DOWN", m.CurrentState)
	}
}

// Cada mudança de estado grava um aviso no outbox (um por canal), na mesma
// transação do check; check sem mudança não grava nada.
func TestStateChangesGoToOutbox(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	m := &domain.Monitor{
		ID: "m1", Name: "x", URL: "http://x.invalid", IntervalSeconds: 60, TimeoutSeconds: 5,
		FailureThreshold: 2, Enabled: true, CurrentState: domain.StateUnknown, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.NewMonitorRepo(db).Create(ctx, m); err != nil {
		t.Fatal(err)
	}

	downs := func() int { return len(pendingNotices(t, db, domain.EventDown)) }
	process := func(res domain.CheckResult) {
		t.Helper()
		uow := newUoW(db)
		cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := service.NewStateMachine(uow, testChannels, nil).Process(cctx, m, res); err != nil {
			t.Fatal(err)
		}
	}
	msg := "boom"
	fail := domain.CheckResult{ErrorMessage: &msg}
	code := 200
	ok := domain.CheckResult{Success: true, StatusCode: &code}

	// 1ª falha: abaixo do threshold, continua UNKNOWN
	process(fail)
	if m.CurrentState != domain.StateUnknown || downs() != 0 {
		t.Fatalf("após 1 falha: state=%s downs=%d", m.CurrentState, downs())
	}

	// 2ª falha: atinge threshold → DOWN + 1 notificação + incidente aberto
	process(fail)
	if m.CurrentState != domain.StateDown || downs() != 1 {
		t.Fatalf("após 2 falhas: state=%s downs=%d", m.CurrentState, downs())
	}
	open, err := repository.NewIncidentRepo(db).FindOpenByMonitor(ctx, m.ID)
	if err != nil || open == nil {
		t.Fatalf("incidente aberto esperado, err=%v", err)
	}

	// 3ª falha: já DOWN, não notifica de novo
	process(fail)
	if downs() != 1 {
		t.Fatalf("notificou down repetido: downs=%d", downs())
	}

	// sucesso: DOWN → UP, fecha incidente e notifica recuperação
	process(ok)
	if m.CurrentState != domain.StateUp || len(pendingNotices(t, db, domain.EventUp)) != 1 {
		t.Fatalf("após sucesso: state=%s recovered=%d", m.CurrentState, len(pendingNotices(t, db, domain.EventUp)))
	}
	if rec := pendingNotices(t, db, domain.EventUp)[0].Snapshot.AsIncident(); rec.EndedAt == nil || rec.DurationSeconds == nil {
		t.Fatal("incidente de recuperação sem ended_at/duration")
	}
	open, _ = repository.NewIncidentRepo(db).FindOpenByMonitor(ctx, m.ID)
	if open != nil {
		t.Fatal("incidente continua aberto após recuperação")
	}

	stored, _ := repository.NewMonitorRepo(db).FindByID(ctx, m.ID)
	if stored.CurrentState != domain.StateUp || stored.LastCheckedAt == nil {
		t.Fatalf("banco: state=%s last_checked_at=%v", stored.CurrentState, stored.LastCheckedAt)
	}
}

// Heartbeat ignora failure_threshold: a primeira falha (prazo perdido) já abre incidente.
func TestHeartbeatGoesDownOnFirstFailure(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	token := "tok"
	m := &domain.Monitor{
		ID: "hb", Name: "backup", Type: domain.TypeHeartbeat, PingToken: &token, IntervalSeconds: 3600,
		TimeoutSeconds: 10, FailureThreshold: 5, Enabled: true, CurrentState: domain.StateUp, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.NewMonitorRepo(db).Create(ctx, m); err != nil {
		t.Fatal(err)
	}

	uow := newUoW(db)
	msg := "no ping"
	if err := service.NewStateMachine(uow, testChannels, nil).Process(ctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
		t.Fatal(err)
	}
	if m.CurrentState != domain.StateDown {
		t.Fatalf("state = %s, want DOWN", m.CurrentState)
	}
	open, err := repository.NewIncidentRepo(db).FindOpenByMonitor(ctx, m.ID)
	if err != nil || open == nil {
		t.Fatalf("incidente aberto esperado, err=%v", err)
	}
}

// failure_count começa nas falhas do threshold, sobe enquanto DOWN e para ao fechar.
func TestFailureCountFollowsOpenIncident(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	m := &domain.Monitor{
		ID: "m1", Name: "x", URL: "http://x.invalid", IntervalSeconds: 60, TimeoutSeconds: 5,
		FailureThreshold: 3, Enabled: true, CurrentState: domain.StateUnknown, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.NewMonitorRepo(db).Create(ctx, m); err != nil {
		t.Fatal(err)
	}
	process := func(res domain.CheckResult) {
		t.Helper()
		uow := newUoW(db)
		if err := service.NewStateMachine(uow, testChannels, nil).Process(ctx, m, res); err != nil {
			t.Fatal(err)
		}
	}
	failureCount := func() int64 {
		t.Helper()
		var n int64
		if err := db.QueryRow(`SELECT failure_count FROM incidents WHERE monitor_id = ?`, m.ID).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	msg := "boom"
	fail := domain.CheckResult{ErrorMessage: &msg}
	code := 200
	ok := domain.CheckResult{Success: true, StatusCode: &code}

	for range 5 {
		process(fail)
	}
	if n := failureCount(); n != 5 {
		t.Fatalf("após 5 falhas (threshold 3): failure_count=%d, want 5", n)
	}

	process(ok)
	process(fail) // nova falha após recuperar: não mexe no incidente fechado
	if n := failureCount(); n != 5 {
		t.Fatalf("incidente fechado mudou: failure_count=%d, want 5", n)
	}
}
