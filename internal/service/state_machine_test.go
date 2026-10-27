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

type noopNotifier struct{}

func (noopNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident)     {}
func (noopNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) {}

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

	uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
	msg := "boom"
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := service.NewStateMachine(uow, noopNotifier{}).Process(cctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
		t.Fatal(err)
	}
	if m.CurrentState != domain.StateDown {
		t.Fatalf("state = %s, want DOWN", m.CurrentState)
	}
}

// dbNotifier imita o WebhookNotifier real: lê a URL do webhook no banco.
type dbNotifier struct {
	settings         *repository.SettingsRepo
	downs, recovered int
	lastRecovery     *domain.Incident
}

func (n *dbNotifier) NotifyDown(ctx context.Context, _ *domain.Monitor, _ *domain.Incident) {
	if _, err := n.settings.Get(ctx, "webhook_url"); err == nil {
		n.downs++
	}
}

func (n *dbNotifier) NotifyRecovery(ctx context.Context, _ *domain.Monitor, i *domain.Incident) {
	if _, err := n.settings.Get(ctx, "webhook_url"); err == nil {
		n.recovered++
		n.lastRecovery = i
	}
}

// Regressão: o notifier rodava com a transação aberta e, ao ler o banco,
// travava esperando a única conexão do SQLite.
func TestNotifierRunsAfterCommit(t *testing.T) {
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

	n := &dbNotifier{settings: repository.NewSettingsRepo(db)}
	process := func(res domain.CheckResult) {
		t.Helper()
		uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
		cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := service.NewStateMachine(uow, n).Process(cctx, m, res); err != nil {
			t.Fatal(err)
		}
	}
	msg := "boom"
	fail := domain.CheckResult{ErrorMessage: &msg}
	code := 200
	ok := domain.CheckResult{Success: true, StatusCode: &code}

	// 1ª falha: abaixo do threshold, continua UNKNOWN
	process(fail)
	if m.CurrentState != domain.StateUnknown || n.downs != 0 {
		t.Fatalf("após 1 falha: state=%s downs=%d", m.CurrentState, n.downs)
	}

	// 2ª falha: atinge threshold → DOWN + 1 notificação + incidente aberto
	process(fail)
	if m.CurrentState != domain.StateDown || n.downs != 1 {
		t.Fatalf("após 2 falhas: state=%s downs=%d", m.CurrentState, n.downs)
	}
	open, err := repository.NewIncidentRepo(db).FindOpenByMonitor(ctx, m.ID)
	if err != nil || open == nil {
		t.Fatalf("incidente aberto esperado, err=%v", err)
	}

	// 3ª falha: já DOWN, não notifica de novo
	process(fail)
	if n.downs != 1 {
		t.Fatalf("notificou down repetido: downs=%d", n.downs)
	}

	// sucesso: DOWN → UP, fecha incidente e notifica recuperação
	process(ok)
	if m.CurrentState != domain.StateUp || n.recovered != 1 {
		t.Fatalf("após sucesso: state=%s recovered=%d", m.CurrentState, n.recovered)
	}
	if n.lastRecovery.EndedAt == nil || n.lastRecovery.DurationSeconds == nil {
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

	uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
	msg := "no ping"
	if err := service.NewStateMachine(uow, noopNotifier{}).Process(ctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
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
