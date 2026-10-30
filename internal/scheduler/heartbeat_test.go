package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

func newHeartbeat(t *testing.T, e *env) *domain.Monitor {
	t.Helper()
	e.svc.SetRecorder(e.sched)
	grace := 60
	m, err := e.svc.Create(context.Background(), service.CreateMonitorInput{
		Name: "backup", Type: "heartbeat", IntervalSeconds: 3600, GraceSeconds: &grace, Enabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.PingToken == nil || len(*m.PingToken) < 40 {
		t.Fatalf("token de ping ausente ou curto: %v", m.PingToken)
	}
	return m
}

func (e *env) load(t *testing.T, id string) *domain.Monitor {
	t.Helper()
	m, err := e.monitors.FindByID(context.Background(), id)
	if err != nil || m == nil {
		t.Fatalf("monitor %s: %v", id, err)
	}
	return m
}

// prazo vencido com folga, na unidade de intervalo do teste
func (e *env) afterDeadline(m *domain.Monitor) time.Time {
	return m.LastPingAt.Add(time.Duration(m.IntervalSeconds+m.Config.GraceSeconds+1) * e.sched.intervalUnit)
}

func TestHeartbeatLifecycle(t *testing.T) {
	e := newEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e)

	if e.sched.hasJob(m.ID) {
		t.Fatal("heartbeat não deveria ter job de check ativo")
	}

	// antes do primeiro ping, dentro do prazo desde a criação: UNKNOWN
	e.sched.checkHeartbeats(ctx, m.UpdatedAt)
	if got := e.load(t, m.ID); got.CurrentState != domain.StateUnknown || e.countChecks(t, m.ID) != 0 {
		t.Fatalf("sem ping: state=%s checks=%d, want UNKNOWN e 0", got.CurrentState, e.countChecks(t, m.ID))
	}

	// ping → UP, sem notificar (UNKNOWN→UP)
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}
	got := e.load(t, m.ID)
	if got.CurrentState != domain.StateUp || got.LastPingAt == nil || e.countChecks(t, m.ID) != 1 {
		t.Fatalf("após ping: state=%s last_ping_at=%v checks=%d", got.CurrentState, got.LastPingAt, e.countChecks(t, m.ID))
	}

	// dentro do prazo: nada muda
	e.sched.checkHeartbeats(ctx, *got.LastPingAt)
	if e.countChecks(t, m.ID) != 1 {
		t.Fatal("avaliador gerou check dentro do prazo")
	}

	// prazo vencido → DOWN uma vez, com incidente e notificação
	late := e.afterDeadline(got)
	e.sched.checkHeartbeats(ctx, late)
	got = e.load(t, m.ID)
	if got.CurrentState != domain.StateDown || e.countChecks(t, m.ID) != 2 || e.notifier.downs.Load() != 1 {
		t.Fatalf("após prazo: state=%s checks=%d downs=%d", got.CurrentState, e.countChecks(t, m.ID), e.notifier.downs.Load())
	}

	// segundo tick já em DOWN: sem check novo
	e.sched.checkHeartbeats(ctx, late.Add(time.Hour))
	if n := e.countChecks(t, m.ID); n != 2 {
		t.Fatalf("tick em DOWN gerou check: checks=%d, want 2", n)
	}

	// novo ping → UP e incidente fechado
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}
	if got = e.load(t, m.ID); got.CurrentState != domain.StateUp {
		t.Fatalf("após novo ping: state=%s, want UP", got.CurrentState)
	}
	open, _ := repository.NewIncidentRepo(e.db).FindOpenByMonitor(ctx, m.ID)
	if open != nil {
		t.Fatal("incidente continua aberto após o ping")
	}
}

// Criado e nunca pingado ("esqueci o curl"): vence o prazo contado da criação.
func TestHeartbeatNeverPingedGoesDown(t *testing.T) {
	e := newEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e)

	late := m.UpdatedAt.Add(time.Duration(m.IntervalSeconds+m.Config.GraceSeconds+1) * e.sched.intervalUnit)
	e.sched.checkHeartbeats(ctx, late)
	if got := e.load(t, m.ID); got.CurrentState != domain.StateDown || e.notifier.downs.Load() != 1 {
		t.Fatalf("sem ping após o prazo: state=%s downs=%d, want DOWN e 1", got.CurrentState, e.notifier.downs.Load())
	}
}

// Prazo que venceu com o Pingou fora do ar não vira DOWN no restart: ganha
// mais um intervalo + grace a partir da subida.
func TestHeartbeatDeadlineMissedWhileOffline(t *testing.T) {
	e := newEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e)
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}
	got := e.load(t, m.ID)

	window := time.Duration(m.IntervalSeconds+m.Config.GraceSeconds) * e.sched.intervalUnit
	e.sched.startedAt = e.afterDeadline(got).Add(time.Hour) // subiu bem depois do prazo

	e.sched.checkHeartbeats(ctx, e.sched.startedAt.Add(window/2))
	if got := e.load(t, m.ID); got.CurrentState != domain.StateUp {
		t.Fatalf("prazo vencido offline gerou alarme: state=%s, want UP", got.CurrentState)
	}

	e.sched.checkHeartbeats(ctx, e.sched.startedAt.Add(window+time.Millisecond))
	if got := e.load(t, m.ID); got.CurrentState != domain.StateDown {
		t.Fatalf("sem ping após a subida + prazo: state=%s, want DOWN", got.CurrentState)
	}
}

func TestHeartbeatPausedAcceptsPingWithoutRecording(t *testing.T) {
	e := newEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e)
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}

	off := false
	if _, err := e.svc.Update(ctx, m.ID, service.UpdateMonitorInput{Enabled: &off}); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatalf("ping pausado: %v, want aceito", err)
	}
	if n := e.countChecks(t, m.ID); n != 1 {
		t.Fatalf("ping pausado registrou check: checks=%d, want 1", n)
	}

	// pausado não vence prazo
	e.sched.checkHeartbeats(ctx, time.Now().Add(365*24*time.Hour))
	if got := e.load(t, m.ID); got.CurrentState != domain.StateUp {
		t.Fatalf("pausado mudou de estado: %s", got.CurrentState)
	}

	// retomar recomeça o prazo a partir de agora
	before := time.Now().Add(-time.Second)
	on := true
	if _, err := e.svc.Update(ctx, m.ID, service.UpdateMonitorInput{Enabled: &on}); err != nil {
		t.Fatal(err)
	}
	if got := e.load(t, m.ID); got.LastPingAt == nil || got.LastPingAt.Before(before) {
		t.Fatalf("retomar não recomeçou o prazo: last_ping_at=%v", got.LastPingAt)
	}
}

func TestHeartbeatRegenerateToken(t *testing.T) {
	e := newEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e)
	old := *m.PingToken

	updated, err := e.svc.RegenerateToken(ctx, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if *updated.PingToken == old {
		t.Fatal("token não mudou")
	}
	if err := e.svc.Ping(ctx, old); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("token antigo: %v, want ErrNotFound", err)
	}
	if err := e.svc.Ping(ctx, *updated.PingToken); err != nil {
		t.Fatalf("token novo: %v", err)
	}

	// só heartbeat tem token
	httpMon, err := e.svc.Create(ctx, validInput)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.RegenerateToken(ctx, httpMon.ID); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("regenerar em http: %v, want ErrValidation", err)
	}
}

func TestHeartbeatTypeCannotChange(t *testing.T) {
	e := newEnv(t, true)
	m := newHeartbeat(t, e)
	typ := "http"
	if _, err := e.svc.Update(context.Background(), m.ID, service.UpdateMonitorInput{Type: &typ}); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("trocar tipo: %v, want ErrValidation", err)
	}
}

// O avaliador roda sozinho depois do Start.
func TestHeartbeatEvaluatorRunsOnTicker(t *testing.T) {
	e := newEnv(t, true)
	e.sched.heartbeatTick = 5 * time.Millisecond
	m := newHeartbeat(t, e)
	if err := e.svc.Ping(context.Background(), *m.PingToken); err != nil {
		t.Fatal(err)
	}
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	// intervalo 3600 + grace 60 em ms ≈ 3,7 s; acima do waitFor. Encurta o prazo.
	interval := 60
	grace := 0
	if _, err := e.svc.Update(context.Background(), m.ID, service.UpdateMonitorInput{IntervalSeconds: &interval, GraceSeconds: &grace}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "heartbeat vencer e ficar DOWN", func() bool {
		return e.load(t, m.ID).CurrentState == domain.StateDown
	})
}
