package scheduler

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type fakeChecker struct {
	calls   atomic.Int64
	success bool
}

func (c *fakeChecker) Check(context.Context, *domain.Monitor) domain.CheckResult {
	c.calls.Add(1)
	if c.success {
		code := 200
		return domain.CheckResult{Success: true, StatusCode: &code}
	}
	msg := "boom"
	return domain.CheckResult{ErrorMessage: &msg}
}

type env struct {
	db       *sql.DB
	monitors *repository.MonitorRepo
	checks   *repository.CheckRepo
	sched    *Scheduler
	svc      *service.MonitorService
	checker  *fakeChecker
}

func newEnv(t *testing.T, success bool) *env {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	e := &env{
		db:       db,
		monitors: repository.NewMonitorRepo(db),
		checks:   repository.NewCheckRepo(db),
		checker:  &fakeChecker{success: success},
	}
	e.sched = NewScheduler(e.monitors, e.checker, db,
		repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db),
		repository.NewOutboxRepoTx(), []string{domain.ChannelWebhook}, nil)
	// 1 "segundo" de intervalo vira 1ms para o teste rodar rápido
	e.sched.intervalUnit = time.Millisecond
	t.Cleanup(e.sched.Stop)

	e.svc = service.NewMonitorService(e.monitors, e.checks, repository.NewIncidentRepo(db))
	e.svc.SetReloader(e.sched)
	return e
}

// countNotices conta os avisos gravados no outbox para o tipo de evento.
func (e *env) countNotices(t *testing.T, eventType string) int {
	t.Helper()
	var n int
	if err := e.db.QueryRow(`SELECT COUNT(*) FROM notification_outbox WHERE event_type = ?`, eventType).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func (e *env) countChecks(t *testing.T, id string) int {
	t.Helper()
	var n int
	if err := e.db.QueryRow(`SELECT COUNT(*) FROM checks WHERE monitor_id = ?`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timeout esperando: %s", what)
}

var validInput = service.CreateMonitorInput{
	Name: "site", URL: "http://example.com", IntervalSeconds: 10,
	TimeoutSeconds: 1, FailureThreshold: 1, Enabled: true,
}

// Regressão: monitor cadastrado pela web usava o contexto da requisição HTTP.
// Quando a requisição terminava, o contexto era cancelado e o job morria
// depois do primeiro check.
func TestMonitorCreatedViaRequestKeepsRunning(t *testing.T) {
	e := newEnv(t, true)
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	reqCtx, endRequest := context.WithCancel(context.Background())
	m, err := e.svc.Create(reqCtx, validInput)
	if err != nil {
		t.Fatal(err)
	}
	endRequest() // requisição HTTP terminou

	waitFor(t, "vários checks após o fim da requisição", func() bool {
		return e.countChecks(t, m.ID) >= 3
	})

	got, _ := e.monitors.FindByID(context.Background(), m.ID)
	if got.CurrentState != domain.StateUp || got.LastCheckedAt == nil {
		t.Fatalf("state = %s, last_checked_at = %v; want UP com last_checked_at", got.CurrentState, got.LastCheckedAt)
	}
}

// Regressão: ler o banco com a transação aberta trava (o SQLite usa 1 conexão
// só) e bloqueava todos os outros monitores. A queda vira um aviso no outbox.
func TestFailingMonitorDoesNotDeadlockScheduler(t *testing.T) {
	e := newEnv(t, false)
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	m, err := e.svc.Create(context.Background(), validInput)
	if err != nil {
		t.Fatal(err)
	}

	waitFor(t, "monitor ficar DOWN e continuar checando", func() bool {
		return e.countChecks(t, m.ID) >= 3
	})
	if n := e.countNotices(t, domain.EventDown); n != 1 {
		t.Fatalf("avisos de down no outbox = %d, want 1", n)
	}
	got, _ := e.monitors.FindByID(context.Background(), m.ID)
	if got.CurrentState != domain.StateDown {
		t.Fatalf("state = %s, want DOWN", got.CurrentState)
	}
}

// Monitores já salvos no banco devem iniciar junto com o servidor.
func TestStartLoadsEnabledMonitorsFromDB(t *testing.T) {
	e := newEnv(t, true)
	on, err := e.svc.Create(context.Background(), validInput) // ainda sem scheduler rodando
	if err != nil {
		t.Fatal(err)
	}
	offInput := validInput
	offInput.Enabled = false
	off, err := e.svc.Create(context.Background(), offInput)
	if err != nil {
		t.Fatal(err)
	}
	e.sched.Stop() // descarta jobs criados pelo Reload acima

	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "monitor ativo checado", func() bool { return e.countChecks(t, on.ID) >= 2 })
	if n := e.countChecks(t, off.ID); n != 0 {
		t.Fatalf("monitor desativado teve %d checks, want 0", n)
	}
}

// Desativar e excluir pela web deve parar o job.
func TestDisableAndDeleteStopJob(t *testing.T) {
	e := newEnv(t, true)
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	m, err := e.svc.Create(context.Background(), validInput)
	if err != nil {
		t.Fatal(err)
	}
	waitFor(t, "primeiro check", func() bool { return e.countChecks(t, m.ID) >= 1 })

	disabled := false
	if _, err := e.svc.Update(context.Background(), m.ID, service.UpdateMonitorInput{Enabled: &disabled}); err != nil {
		t.Fatal(err)
	}
	if e.sched.hasJob(m.ID) {
		t.Fatal("job continua ativo após desativar")
	}

	enabled := true
	if _, err := e.svc.Update(context.Background(), m.ID, service.UpdateMonitorInput{Enabled: &enabled}); err != nil {
		t.Fatal(err)
	}
	if !e.sched.hasJob(m.ID) {
		t.Fatal("job não voltou após reativar")
	}

	if err := e.svc.Delete(context.Background(), m.ID); err != nil {
		t.Fatal(err)
	}
	if e.sched.hasJob(m.ID) {
		t.Fatal("job continua ativo após excluir")
	}
}

// slowChannel segura o envio até release fechar, como um webhook lento.
type slowChannel struct {
	started chan struct{}
	release chan struct{}
}

func (c *slowChannel) Name() string { return domain.ChannelWebhook }

func (c *slowChannel) Send(context.Context, domain.OutboxEvent) error {
	c.started <- struct{}{}
	<-c.release
	return nil
}

// Regressão: o envio rodava dentro do lock dbMu. Um destino lento no monitor A
// pausava a gravação de todos os outros monitores. Agora o envio é do worker do
// outbox, fora do lock e sem segurar a conexão do banco.
func TestSlowChannelDoesNotBlockOtherMonitors(t *testing.T) {
	e := newEnv(t, true)
	slow := &slowChannel{started: make(chan struct{}, 1), release: make(chan struct{})}
	worker := service.NewOutboxWorker(repository.NewOutboxRepo(e.db), slow)
	e.sched.wake = worker.Wake
	ctx, cancel := context.WithCancel(context.Background())
	worker.Start(ctx)
	defer func() {
		close(slow.release)
		cancel()
		waitCtx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		if err := worker.Wait(waitCtx); err != nil {
			t.Error("envios pendentes no fim do teste")
		}
	}()

	// service sem reloader: sem jobs rodando, só os Record do teste gravam checks
	ctx = context.Background()
	svc := service.NewMonitorService(e.monitors, e.checks, repository.NewIncidentRepo(e.db))
	a, err := svc.Create(ctx, validInput)
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.Create(ctx, validInput)
	if err != nil {
		t.Fatal(err)
	}

	msg := "boom"
	if err := e.sched.Record(ctx, a.ID, domain.CheckResult{ErrorMessage: &msg}, nil); err != nil {
		t.Fatal(err)
	}
	<-slow.started // o envio de A está preso

	done := make(chan error, 1)
	go func() {
		code := 200
		done <- e.sched.Record(ctx, b.ID, domain.CheckResult{Success: true, StatusCode: &code}, nil)
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Record do monitor B esperou o envio do monitor A")
	}
	if e.countChecks(t, b.ID) != 1 {
		t.Fatal("check de B não foi gravado")
	}
}
