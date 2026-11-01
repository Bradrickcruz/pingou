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

// dbNotifier imita o WebhookNotifier real: lê a URL do webhook no banco.
type dbNotifier struct {
	settings *repository.SettingsRepo
	downs    atomic.Int64
}

func (n *dbNotifier) NotifyDown(ctx context.Context, _ *domain.Monitor, _ *domain.Incident) {
	n.settings.Get(ctx, "webhook_url")
	n.downs.Add(1)
}

func (n *dbNotifier) NotifyRecovery(ctx context.Context, _ *domain.Monitor, _ *domain.Incident) {
	n.settings.Get(ctx, "webhook_url")
}

type env struct {
	db       *sql.DB
	monitors *repository.MonitorRepo
	checks   *repository.CheckRepo
	sched    *Scheduler
	svc      *service.MonitorService
	checker  *fakeChecker
	notifier *dbNotifier
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
		notifier: &dbNotifier{settings: repository.NewSettingsRepo(db)},
	}
	e.sched = NewScheduler(e.monitors, e.checker, db,
		repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db),
		e.notifier)
	// 1 "segundo" de intervalo vira 1ms para o teste rodar rápido
	e.sched.intervalUnit = time.Millisecond
	t.Cleanup(e.sched.Stop)

	e.svc = service.NewMonitorService(e.monitors, e.checks, repository.NewIncidentRepo(db))
	e.svc.SetReloader(e.sched)
	return e
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

// Regressão: o notifier era chamado com a transação aberta. Como o SQLite usa
// 1 conexão só, o notifier (que lê o webhook no banco) travava para sempre e
// bloqueava todos os outros monitores.
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
	if n := e.notifier.downs.Load(); n != 1 {
		t.Fatalf("notificações de down = %d, want 1", n)
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

// slowNotifier segura o envio até release fechar, como um webhook lento.
type slowNotifier struct {
	started chan struct{}
	release chan struct{}
}

func (n *slowNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident) {
	n.started <- struct{}{}
	<-n.release
}
func (n *slowNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) {}

// Regressão: o notifier rodava dentro do lock dbMu. Um destino lento no monitor
// A pausava a gravação de todos os outros monitores.
func TestSlowNotifierDoesNotBlockOtherMonitors(t *testing.T) {
	e := newEnv(t, true)
	slow := &slowNotifier{started: make(chan struct{}, 1), release: make(chan struct{})}
	multi := service.NewMultiNotifier(slow, e.notifier) // e.notifier lê o banco, como o webhook real
	e.sched.notifier = multi
	defer func() {
		close(slow.release)
		if !multi.Wait(time.Second) {
			t.Error("envios pendentes no fim do teste")
		}
	}()

	// service sem reloader: sem jobs rodando, só os Record do teste gravam checks
	ctx := context.Background()
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
		t.Fatal("Record do monitor B esperou o notifier do monitor A")
	}
	if e.countChecks(t, b.ID) != 1 {
		t.Fatal("check de B não foi gravado")
	}
}
