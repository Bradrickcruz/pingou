package scheduler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// mEnv: env com janelas de manutenção ligadas. intervalUnit em segundos de
// verdade: só o check imediato do Reload roda durante o teste.
type mEnv struct {
	*env
	windows *repository.MaintenanceRepo
	maint   *service.MaintenanceService
}

func newMaintenanceEnv(t *testing.T, success bool) *mEnv {
	t.Helper()
	e := newEnv(t, success)
	e.sched.intervalUnit = time.Second
	windows := repository.NewMaintenanceRepo(e.db)
	e.sched.WithMaintenance(windows)
	e.svc.SetRecorder(e.sched)
	maint := service.NewMaintenanceService(windows, e.monitors)
	maint.SetEnder(e.sched)
	return &mEnv{env: e, windows: windows, maint: maint}
}

// addWindow grava a janela direto no repositório (sem a regra de "início no futuro").
func (e *mEnv) addWindow(t *testing.T, monitorID string, starts, ends time.Time) *domain.MaintenanceWindow {
	t.Helper()
	reason := "troca de disco"
	w := &domain.MaintenanceWindow{
		ID: "w-" + monitorID, MonitorID: monitorID, Reason: &reason,
		StartsAt: starts.UTC().Truncate(time.Second), EndsAt: ends.UTC().Truncate(time.Second),
		CreatedAt: time.Now().UTC(),
	}
	if err := e.windows.CreateBatch(context.Background(), []*domain.MaintenanceWindow{w}); err != nil {
		t.Fatal(err)
	}
	return w
}

// expire leva o fim da janela para o passado (o tempo "passou").
func (e *mEnv) expire(t *testing.T, w *domain.MaintenanceWindow) {
	t.Helper()
	if _, err := e.db.Exec(`UPDATE maintenance_windows SET ends_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339), w.ID); err != nil {
		t.Fatal(err)
	}
}

func (e *mEnv) window(t *testing.T, id string) *domain.MaintenanceWindow {
	t.Helper()
	w, err := e.windows.FindByID(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

type incidentRow struct {
	ended  sql.NullString
	reason int
	note   sql.NullString
}

// incident devolve o primeiro incidente do monitor (o aberto antes da janela).
func (e *mEnv) incident(t *testing.T, monitorID string) incidentRow {
	t.Helper()
	var r incidentRow
	if err := e.db.QueryRow(`SELECT ended_at, closed_manually, close_note FROM incidents
		WHERE monitor_id = ? ORDER BY started_at, id LIMIT 1`, monitorID).Scan(&r.ended, &r.reason, &r.note); err != nil {
		t.Fatal(err)
	}
	return r
}

// monitor DOWN com incidente aberto (threshold 1, checker falhando).
func (e *mEnv) downMonitor(t *testing.T, threshold int) *domain.Monitor {
	t.Helper()
	in := validInput
	in.FailureThreshold = threshold
	m, err := e.svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestMaintenanceBlocksChecksAndClosesIncident(t *testing.T) {
	e := newMaintenanceEnv(t, false)
	ctx := context.Background()
	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	m := e.downMonitor(t, 1)
	waitFor(t, "monitor DOWN", func() bool { return e.load(t, m.ID).CurrentState == domain.StateDown })
	downs, checks := e.countNotices(t, domain.EventDown), e.countChecks(t, m.ID)

	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Second), now.Add(time.Hour))
	e.sched.applyMaintenance(ctx)

	if got := e.window(t, w.ID); got.StartedAt == nil {
		t.Fatal("started_at não gravado")
	}
	inc := e.incident(t, m.ID)
	if !inc.ended.Valid || inc.reason != int(domain.CloseMaintenance) || inc.note.String != "troca de disco" {
		t.Fatalf("incidente: %+v, want fechado por manutenção com o motivo", inc)
	}
	if got := e.load(t, m.ID); got.CurrentState != domain.StateUnknown {
		t.Fatalf("estado=%s, want UNKNOWN", got.CurrentState)
	}

	// edição durante a janela: o Reload roda o check imediato, que o gate barra
	name := "editado"
	if _, err := e.svc.Update(ctx, m.ID, service.UpdateMonitorInput{Name: &name}); err != nil {
		t.Fatal(err)
	}
	if err := e.sched.Record(ctx, m.ID, domain.CheckResult{}, nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	if n := e.countChecks(t, m.ID); n != checks {
		t.Fatalf("checks durante a janela: %d → %d, want nenhum novo", checks, n)
	}
	if n := e.countNotices(t, domain.EventDown); n != downs {
		t.Fatalf("avisos durante a janela: %d → %d", downs, n)
	}
}

// Ao terminar: UNKNOWN, primeiro check na hora e falhas antigas não contam.
func TestMaintenanceEndResumesMonitor(t *testing.T) {
	e := newMaintenanceEnv(t, false)
	ctx := context.Background()
	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	m := e.downMonitor(t, 2)
	// 2 falhas antes da janela: no limite do threshold
	waitFor(t, "check imediato", func() bool { return e.countChecks(t, m.ID) >= 1 })
	if err := e.sched.Record(ctx, m.ID, domain.CheckResult{}, nil); err != nil {
		t.Fatal(err)
	}
	if got := e.load(t, m.ID); got.CurrentState != domain.StateDown {
		t.Fatalf("antes: estado=%s, want DOWN", got.CurrentState)
	}
	// a janela começa no futuro: os checks de antes dela são de antes do fim
	// (no teste tudo cai no mesmo segundo, a precisão do banco)
	if _, err := e.db.Exec(`UPDATE checks SET checked_at = ? WHERE monitor_id = ?`,
		time.Now().UTC().Add(-time.Hour).Format(time.RFC3339), m.ID); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Minute), now.Add(time.Hour))
	e.sched.applyMaintenance(ctx)
	before := e.countChecks(t, m.ID)

	e.expire(t, w)
	e.sched.applyMaintenance(ctx)

	if got := e.window(t, w.ID); got.FinalizedAt == nil {
		t.Fatal("finalized_at não gravado")
	}
	waitFor(t, "primeiro check após a janela", func() bool { return e.countChecks(t, m.ID) == before+1 })
	got := e.load(t, m.ID)
	if got.ResumedAt == nil {
		t.Fatal("resumed_at não gravado")
	}
	// 1 falha depois da janela < threshold 2: as de antes não somam
	if got.CurrentState != domain.StateUnknown {
		t.Fatalf("estado=%s, want UNKNOWN (falhas de antes da janela não contam)", got.CurrentState)
	}
	var open int
	e.db.QueryRow(`SELECT COUNT(*) FROM incidents WHERE monitor_id = ? AND ended_at IS NULL`, m.ID).Scan(&open)
	if open != 0 {
		t.Fatalf("incidentes abertos=%d, want 0", open)
	}
}

// Pausa manual vence a janela (D3).
func TestMaintenanceEndKeepsManualPause(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	ctx := context.Background()
	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	m := e.downMonitor(t, 1)
	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Minute), now.Add(time.Hour))
	e.sched.applyMaintenance(ctx)

	off := false
	if _, err := e.svc.Update(ctx, m.ID, service.UpdateMonitorInput{Enabled: &off}); err != nil {
		t.Fatal(err)
	}
	e.expire(t, w)
	e.sched.applyMaintenance(ctx)

	if got := e.load(t, m.ID); got.Enabled {
		t.Fatal("monitor voltou ativo; a pausa manual deveria vencer")
	}
	if e.sched.hasJob(m.ID) {
		t.Fatal("job de check rodando com o monitor pausado")
	}
}

func TestMaintenanceHeartbeat(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	ctx := context.Background()
	m := newHeartbeat(t, e.env)
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}
	m = e.load(t, m.ID)
	pings := e.countChecks(t, m.ID)

	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Minute), now.Add(time.Hour))
	e.sched.applyMaintenance(ctx)

	// ping na janela: aceito (sem erro), nada registrado
	if err := e.svc.Ping(ctx, *m.PingToken); err != nil {
		t.Fatal(err)
	}
	// prazo vencido na janela: não vira DOWN
	e.sched.checkHeartbeats(ctx, e.afterDeadline(m))
	if n := e.countChecks(t, m.ID); n != pings {
		t.Fatalf("checks na janela: %d → %d, want nenhum", pings, n)
	}

	e.expire(t, w)
	e.sched.applyMaintenance(ctx)
	got := e.load(t, m.ID)
	if got.LastPingAt == nil || time.Since(*got.LastPingAt) > 5*time.Second {
		t.Fatalf("last_ping_at=%v, want ~agora (prazo recomeça no fim da janela)", got.LastPingAt)
	}
	// logo depois do fim, dentro do prazo novo: sem DOWN
	e.sched.checkHeartbeats(ctx, time.Now())
	if got := e.load(t, m.ID); got.CurrentState == domain.StateDown {
		t.Fatal("heartbeat caiu logo após o fim da janela")
	}
}

// Um check que estava rodando quando a janela começou é descartado.
func TestMaintenanceDiscardsCheckInFlight(t *testing.T) {
	e := newMaintenanceEnv(t, false)
	release := make(chan struct{})
	started := make(chan struct{}, 1)
	e.sched.checker = checkerFunc(func(ctx context.Context, _ *domain.Monitor) domain.CheckResult {
		select {
		case started <- struct{}{}:
		default:
		}
		<-release
		msg := "boom"
		return domain.CheckResult{ErrorMessage: &msg}
	})
	if err := e.sched.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	in := validInput
	in.TimeoutSeconds = 30
	m, err := e.svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("check não começou")
	}
	now := time.Now()
	e.addWindow(t, m.ID, now.Add(-time.Second), now.Add(time.Hour))
	close(release)
	time.Sleep(100 * time.Millisecond)
	if n := e.countChecks(t, m.ID); n != 0 {
		t.Fatalf("checks=%d, want 0 (check em andamento no início da janela)", n)
	}
}

func TestMaintenanceCancel(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	ctx := context.Background()
	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	a := e.downMonitor(t, 1)
	now := time.Now()

	// antes do início: apaga
	future := e.addWindow(t, a.ID, now.Add(time.Hour), now.Add(2*time.Hour))
	if err := e.maint.Cancel(ctx, future.ID); err != nil {
		t.Fatal(err)
	}
	if e.window(t, future.ID) != nil {
		t.Fatal("janela futura cancelada continua no banco")
	}

	// durante: encerra na hora e o monitor volta
	b := e.downMonitor(t, 1)
	running := e.addWindow(t, b.ID, now.Add(-time.Minute), now.Add(time.Hour))
	e.sched.applyMaintenance(ctx)
	before := e.countChecks(t, b.ID)
	if err := e.maint.Cancel(ctx, running.ID); err != nil {
		t.Fatal(err)
	}
	got := e.window(t, running.ID)
	if got.FinalizedAt == nil || got.EndsAt.After(time.Now()) {
		t.Fatalf("cancelada durante: finalized=%v ends=%v", got.FinalizedAt, got.EndsAt)
	}
	waitFor(t, "check após cancelar", func() bool { return e.countChecks(t, b.ID) > before })

	// já finalizada: 404
	if err := e.maint.Cancel(ctx, running.ID); err != service.ErrNotFound {
		t.Fatalf("cancelar de novo: %v, want ErrNotFound", err)
	}
}

// Cancelar no segundo exato do início: ends_at fica 1 s depois (CHECK do banco).
func TestMaintenanceCancelAtStartSecond(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	ctx := context.Background()
	m := e.downMonitor(t, 1)
	now := time.Now().UTC().Truncate(time.Second)
	w := e.addWindow(t, m.ID, now, now.Add(time.Hour))
	e.maint.WithClock(func() time.Time { return now })
	if err := e.maint.Cancel(ctx, w.ID); err != nil {
		t.Fatal(err)
	}
	if got := e.window(t, w.ID); got.FinalizedAt == nil || !got.EndsAt.Equal(now.Add(time.Second)) {
		t.Fatalf("finalized=%v ends=%v", got.FinalizedAt, got.EndsAt)
	}
}

// Reinício: janela em curso segue bloqueando; janela vencida com o Pingou fora
// é finalizada (e o início pendente aplicado antes).
func TestMaintenanceAppliedOnStart(t *testing.T) {
	e := newMaintenanceEnv(t, false)
	ctx := context.Background()
	now := time.Now()

	// direto no repositório: Create recarregaria o job antes do Start
	raw := func(id string) *domain.Monitor {
		m := &domain.Monitor{
			ID: id, Name: id, Type: domain.TypeHTTP, URL: "http://example.com", IntervalSeconds: 10,
			TimeoutSeconds: 1, FailureThreshold: 1, Enabled: true, CurrentState: domain.StateDown,
			CreatedAt: now, UpdatedAt: now,
		}
		if err := e.monitors.Create(ctx, m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	running, expired := raw("running"), raw("expired")
	for _, id := range []string{running.ID, expired.ID} {
		if _, err := e.db.Exec(`INSERT INTO incidents (id, monitor_id, started_at) VALUES (?,?,?)`,
			"i-"+id, id, now.Add(-3*time.Hour).UTC().Format(time.RFC3339)); err != nil {
			t.Fatal(err)
		}
	}
	wr := e.addWindow(t, running.ID, now.Add(-time.Minute), now.Add(time.Hour))
	we := e.addWindow(t, expired.ID, now.Add(-2*time.Hour), now.Add(-time.Hour))

	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if got := e.window(t, wr.ID); got.StartedAt == nil || got.FinalizedAt != nil {
		t.Fatalf("em curso: started=%v finalized=%v", got.StartedAt, got.FinalizedAt)
	}
	if got := e.window(t, we.ID); got.StartedAt == nil || got.FinalizedAt == nil {
		t.Fatalf("vencida: started=%v finalized=%v", got.StartedAt, got.FinalizedAt)
	}
	for _, id := range []string{running.ID, expired.ID} {
		if inc := e.incident(t, id); inc.reason != int(domain.CloseMaintenance) {
			t.Fatalf("%s: incidente com motivo %d, want manutenção", id, inc.reason)
		}
	}
	time.Sleep(50 * time.Millisecond)
	if n := e.countChecks(t, running.ID); n != 0 {
		t.Fatalf("monitor em manutenção conferido após o reinício: %d checks", n)
	}
	waitFor(t, "check do monitor cuja janela venceu", func() bool { return e.countChecks(t, expired.ID) > 0 })
}

func TestMaintenanceMonitorDeleted(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	ctx := context.Background()
	m := e.downMonitor(t, 1)
	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Second), now.Add(time.Hour))
	if err := e.svc.Delete(ctx, m.ID); err != nil {
		t.Fatal(err)
	}
	if e.window(t, w.ID) != nil {
		t.Fatal("janela do monitor excluído continua no banco (ON DELETE CASCADE)")
	}
	e.sched.applyMaintenance(ctx) // sem erro nem pânico
}

// O loop de manutenção entra no Wait do desligamento.
func TestMaintenanceTickRunsAndStops(t *testing.T) {
	e := newMaintenanceEnv(t, true)
	e.sched.maintenanceTick = 10 * time.Millisecond
	ctx := context.Background()
	if err := e.sched.Start(ctx); err != nil {
		t.Fatal(err)
	}
	m := e.downMonitor(t, 1)
	now := time.Now()
	w := e.addWindow(t, m.ID, now.Add(-time.Second), now.Add(time.Hour))
	waitFor(t, "tick aplica o início", func() bool { return e.window(t, w.ID).StartedAt != nil })

	e.sched.Stop()
	wctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := e.sched.Wait(wctx); err != nil {
		t.Fatalf("Wait = %v", err)
	}
}
