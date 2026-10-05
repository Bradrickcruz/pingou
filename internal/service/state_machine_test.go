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
