package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type countNotifier struct{ downs, recoveries int }

func (n *countNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident) { n.downs++ }
func (n *countNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) {
	n.recoveries++
}
func (n *countNotifier) NotifyCertExpiring(context.Context, *domain.Monitor, domain.CertAlert) {}

func updateEnabled(b *bool) service.UpdateMonitorInput {
	return service.UpdateMonitorInput{Enabled: b}
}

func TestPauseClosesOpenIncidentSilently(t *testing.T) {
	svc, db := newStatsService(t)
	started := time.Now().UTC().Add(-10 * time.Minute).Format(time.RFC3339)
	db.Exec(`UPDATE monitors SET current_state='DOWN' WHERE id='m'`)
	db.Exec(`INSERT INTO incidents (id, monitor_id, started_at) VALUES ('i','m',?)`, started)

	off := false
	for range 2 { // segunda pausa é no-op (idempotente)
		if _, err := svc.Update(context.Background(), "m", updateEnabled(&off)); err != nil {
			t.Fatal(err)
		}
	}

	var ended *string
	var manual, dur int
	if err := db.QueryRow(`SELECT ended_at, closed_manually, duration_seconds FROM incidents WHERE id='i'`).Scan(&ended, &manual, &dur); err != nil {
		t.Fatal(err)
	}
	if ended == nil || manual != 1 || dur < 590 {
		t.Fatalf("ended=%v manual=%d dur=%d, want fechado manualmente com ~600 s", ended, manual, dur)
	}
	var state string
	db.QueryRow(`SELECT current_state FROM monitors WHERE id='m'`).Scan(&state)
	if state != "UNKNOWN" {
		t.Fatalf("estado=%s, want UNKNOWN", state)
	}
}

func TestPauseWithoutIncidentAndResumeKeepState(t *testing.T) {
	svc, db := newStatsService(t)
	off, on := false, true
	if _, err := svc.Update(context.Background(), "m", updateEnabled(&off)); err != nil {
		t.Fatal(err)
	}
	db.Exec(`UPDATE monitors SET current_state='UP' WHERE id='m'`)
	if _, err := svc.Update(context.Background(), "m", updateEnabled(&on)); err != nil {
		t.Fatal(err)
	}
	var state string
	db.QueryRow(`SELECT current_state FROM monitors WHERE id='m'`).Scan(&state)
	if state != "UP" {
		t.Fatalf("retomar não deve mexer no estado, got %s", state)
	}
}

// Pausar não notifica; retomar um monitor que segue fora do ar abre incidente novo
// na primeira falha (as falhas anteriores à pausa ainda contam para o threshold).
func TestPauseIsSilentAndResumeReopensIncident(t *testing.T) {
	svc, db := newStatsService(t)
	ctx := context.Background()
	notifier := &countNotifier{}
	process := func() {
		t.Helper()
		m, err := repository.NewMonitorRepo(db).FindByID(ctx, "m")
		if err != nil {
			t.Fatal(err)
		}
		uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
		msg := "boom"
		if err := service.NewStateMachine(uow, notifier).Process(ctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
			t.Fatal(err)
		}
	}
	open := func() int {
		var n int
		db.QueryRow(`SELECT COUNT(*) FROM incidents WHERE monitor_id='m' AND ended_at IS NULL`).Scan(&n)
		return n
	}

	for range 5 { // acima de qualquer threshold
		process()
	}
	if notifier.downs != 1 || open() != 1 {
		t.Fatalf("antes da pausa: downs=%d abertos=%d, want 1 e 1", notifier.downs, open())
	}

	off, on := false, true
	if _, err := svc.Update(ctx, "m", updateEnabled(&off)); err != nil {
		t.Fatal(err)
	}
	if notifier.downs != 1 || notifier.recoveries != 0 || open() != 0 {
		t.Fatalf("pausa: downs=%d recoveries=%d abertos=%d, want 1, 0 e 0", notifier.downs, notifier.recoveries, open())
	}

	if _, err := svc.Update(ctx, "m", updateEnabled(&on)); err != nil {
		t.Fatal(err)
	}
	process()
	if notifier.downs != 2 || open() != 1 {
		t.Fatalf("retomada: downs=%d abertos=%d, want 2 e 1", notifier.downs, open())
	}
}
