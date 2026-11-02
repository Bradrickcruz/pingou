package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// blockingChecker segura o check até o contexto ser cancelado e então devolve
// a falha que um checker real devolveria ("context canceled").
type blockingChecker struct {
	started chan struct{}
}

func (c *blockingChecker) Check(ctx context.Context, _ *domain.Monitor) domain.CheckResult {
	select {
	case c.started <- struct{}{}:
	default:
	}
	<-ctx.Done()
	msg := ctx.Err().Error()
	return domain.CheckResult{ErrorMessage: &msg}
}

// Desligar no meio de um check não pode virar falha nem abrir incidente.
func TestStopDiscardsCheckCancelledByShutdown(t *testing.T) {
	e := newEnv(t, true)
	bc := &blockingChecker{started: make(chan struct{}, 1)}
	e.sched.checker = bc

	in := validInput
	in.TimeoutSeconds = 30 // o prazo do check não acaba antes do Stop
	m, err := e.svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	select {
	case <-bc.started:
	case <-time.After(3 * time.Second):
		t.Fatal("check não começou")
	}

	e.sched.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := e.sched.Wait(ctx); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}

	if n := e.countChecks(t, m.ID); n != 0 {
		t.Fatalf("checks gravados = %d, want 0 (check cancelado deve ser descartado)", n)
	}
	var incidents int
	if err := e.db.QueryRow(`SELECT COUNT(*) FROM incidents`).Scan(&incidents); err != nil {
		t.Fatal(err)
	}
	if incidents != 0 {
		t.Fatalf("incidentes = %d, want 0", incidents)
	}
	if n := e.notifier.downs.Load(); n != 0 {
		t.Fatalf("notificações de down = %d, want 0", n)
	}
}

// Depois do Stop, um Reload (ex.: requisição que terminava durante o
// desligamento) não começa job novo.
func TestReloadAfterStopStartsNothing(t *testing.T) {
	e := newEnv(t, true)
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	e.sched.Stop()

	m, err := e.svc.Create(context.Background(), validInput) // chama Reload
	if err != nil {
		t.Fatal(err)
	}
	if e.sched.hasJob(m.ID) {
		t.Fatal("job criado depois do Stop")
	}
	time.Sleep(50 * time.Millisecond)
	if n := e.countChecks(t, m.ID); n != 0 {
		t.Fatalf("checks = %d, want 0", n)
	}
}

// Wait devolve erro quando o prazo acaba com um check ainda preso.
func TestWaitRespectsDeadline(t *testing.T) {
	e := newEnv(t, true)
	release := make(chan struct{})
	e.sched.checker = checkerFunc(func(context.Context, *domain.Monitor) domain.CheckResult {
		<-release // ignora o cancelamento de propósito
		return domain.CheckResult{Success: true}
	})
	t.Cleanup(func() { close(release) })

	if _, err := e.svc.Create(context.Background(), validInput); err != nil {
		t.Fatal(err)
	}
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	e.sched.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := e.sched.Wait(ctx); err == nil {
		t.Fatal("Wait = nil, want erro de prazo")
	}
}

type checkerFunc func(context.Context, *domain.Monitor) domain.CheckResult

func (f checkerFunc) Check(ctx context.Context, m *domain.Monitor) domain.CheckResult {
	return f(ctx, m)
}
