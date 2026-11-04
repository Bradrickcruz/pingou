package service_test

import (
	"context"
	"errors"
	"strings"
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

// Pausar não notifica; retomar um monitor que segue fora do ar abre incidente
// novo, mas só quando as falhas DEPOIS da retomada atingem o threshold (3): as
// de antes da pausa não somam.
func TestPauseIsSilentAndResumeReopensIncident(t *testing.T) {
	svc, db := newStatsService(t)
	ctx := context.Background()
	notifier := &countNotifier{}
	// checks de antes da pausa ficam 10 min no passado; os de depois, 1 s à frente
	// (a retomada grava resumed_at com o relógio real, com precisão de segundos)
	clock := time.Now().UTC().Add(-10 * time.Minute)
	process := func() {
		t.Helper()
		m, err := repository.NewMonitorRepo(db).FindByID(ctx, "m")
		if err != nil {
			t.Fatal(err)
		}
		uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
		msg := "boom"
		now := clock
		sm := service.NewStateMachine(uow, notifier).WithClock(func() time.Time { return now })
		if err := sm.Process(ctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
			t.Fatal(err)
		}
		clock = clock.Add(time.Second)
	}
	open := func() int {
		var n int
		db.QueryRow(`SELECT COUNT(*) FROM incidents WHERE monitor_id='m' AND ended_at IS NULL`).Scan(&n)
		return n
	}

	for range 5 { // acima do threshold
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
	clock = time.Now().UTC().Add(time.Second)
	process()
	process()
	if notifier.downs != 1 || open() != 0 {
		t.Fatalf("2 falhas após retomar: downs=%d abertos=%d, want 1 e 0", notifier.downs, open())
	}
	process()
	if notifier.downs != 2 || open() != 1 {
		t.Fatalf("3 falhas após retomar: downs=%d abertos=%d, want 2 e 1", notifier.downs, open())
	}
}

// Bug da retomada: threshold 3, 2 falhas antes da pausa e 1 depois não abrem
// incidente (antes, as três somavam).
func TestResumeIgnoresFailuresBeforePause(t *testing.T) {
	svc, db := newStatsService(t)
	ctx := context.Background()
	before := time.Now().UTC().Add(-time.Hour)
	for i := range 2 {
		db.Exec(`INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES ('m',0,10,?)`,
			before.Add(time.Duration(i)*time.Minute).Format(time.RFC3339))
	}

	off, on := false, true
	if _, err := svc.Update(ctx, "m", updateEnabled(&off)); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Update(ctx, "m", updateEnabled(&on)); err != nil {
		t.Fatal(err)
	}

	m, err := repository.NewMonitorRepo(db).FindByID(ctx, "m")
	if err != nil {
		t.Fatal(err)
	}
	if m.ResumedAt == nil {
		t.Fatal("resumed_at não gravado")
	}
	notifier := &countNotifier{}
	uow := service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db))
	msg := "boom"
	if err := service.NewStateMachine(uow, notifier).Process(ctx, m, domain.CheckResult{ErrorMessage: &msg}); err != nil {
		t.Fatal(err)
	}
	if notifier.downs != 0 || m.CurrentState == domain.StateDown {
		t.Fatalf("downs=%d state=%s, want 0 e não DOWN", notifier.downs, m.CurrentState)
	}
}

func TestPauseReasonSavedAndClearedOnResume(t *testing.T) {
	svc, db := newStatsService(t)
	ctx := context.Background()
	db.Exec(`UPDATE monitors SET current_state='DOWN' WHERE id='m'`)
	db.Exec(`INSERT INTO incidents (id, monitor_id, started_at) VALUES ('i','m',?)`, time.Now().UTC().Add(-time.Minute).Format(time.RFC3339))

	off, on := false, true
	reason := "  troca do servidor  "
	m, err := svc.Update(ctx, "m", service.UpdateMonitorInput{Enabled: &off, PauseReason: &reason})
	if err != nil {
		t.Fatal(err)
	}
	if m.PauseReason == nil || *m.PauseReason != "troca do servidor" || m.PausedAt == nil {
		t.Fatalf("pausa: reason=%v paused_at=%v", m.PauseReason, m.PausedAt)
	}

	inc, _, err := repository.NewIncidentRepo(db).FindAll(ctx, domain.IncidentFilter{MonitorID: "m"})
	if err != nil || len(inc) != 1 {
		t.Fatalf("incidentes: %v %d", err, len(inc))
	}
	if inc[0].CloseReason != domain.ClosePause || inc[0].CloseNote == nil || *inc[0].CloseNote != "troca do servidor" {
		t.Fatalf("incidente: reason=%v note=%v", inc[0].CloseReason, inc[0].CloseNote)
	}

	// pausado: o motivo pode ser trocado; ausente mantém
	other := "migração"
	if m, err = svc.Update(ctx, "m", service.UpdateMonitorInput{PauseReason: &other}); err != nil || *m.PauseReason != "migração" {
		t.Fatalf("trocar motivo: %v %v", err, m.PauseReason)
	}
	name := "novo nome"
	if m, err = svc.Update(ctx, "m", service.UpdateMonitorInput{Name: &name}); err != nil || m.PauseReason == nil {
		t.Fatalf("editar sem motivo apagou o motivo: %v %v", err, m.PauseReason)
	}

	m, err = svc.Update(ctx, "m", updateEnabled(&on))
	if err != nil {
		t.Fatal(err)
	}
	stored, _ := repository.NewMonitorRepo(db).FindByID(ctx, "m")
	for _, got := range []*domain.Monitor{m, stored} {
		if got.PauseReason != nil || got.PausedAt != nil || got.ResumedAt == nil {
			t.Fatalf("retomar: reason=%v paused_at=%v resumed_at=%v", got.PauseReason, got.PausedAt, got.ResumedAt)
		}
	}
}

func TestPauseReasonTooLong(t *testing.T) {
	svc, _ := newStatsService(t)
	off := false
	ok := strings.Repeat("é", 200) // 200 caracteres, 400 bytes
	if _, err := svc.Update(context.Background(), "m", service.UpdateMonitorInput{Enabled: &off, PauseReason: &ok}); err != nil {
		t.Fatalf("200 caracteres: %v", err)
	}
	long := strings.Repeat("a", 201)
	_, err := svc.Update(context.Background(), "m", service.UpdateMonitorInput{PauseReason: &long})
	if !errors.Is(err, service.ErrValidation) {
		t.Fatalf("201 caracteres: err=%v, want ErrValidation", err)
	}
}
