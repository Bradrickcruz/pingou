package service_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

func newMaintenanceService(t *testing.T, now time.Time) (*service.MaintenanceService, func() int) {
	t.Helper()
	_, db := newStatsService(t) // cria o monitor 'm'
	if _, err := db.Exec(`INSERT INTO monitors (id, name, url, interval_seconds, created_at, updated_at)
		VALUES ('m2','y','http://y',60,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	svc := service.NewMaintenanceService(repository.NewMaintenanceRepo(db), repository.NewMonitorRepo(db)).
		WithClock(func() time.Time { return now })
	count := func() int {
		var n int
		db.QueryRow(`SELECT COUNT(*) FROM maintenance_windows`).Scan(&n)
		return n
	}
	return svc, count
}

func TestMaintenanceValidation(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, count := newMaintenanceService(t, now)
	long := strings.Repeat("a", 201)
	cases := []struct {
		name string
		in   service.CreateMaintenanceInput
		want error
	}{
		{"lista vazia", service.CreateMaintenanceInput{StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour)}, service.ErrValidation},
		{"início no passado", service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}, service.ErrValidation},
		{"início agora", service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now, EndsAt: now.Add(time.Hour)}, service.ErrValidation},
		{"fim antes do início", service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now.Add(2 * time.Hour), EndsAt: now.Add(time.Hour)}, service.ErrValidation},
		{"7 dias + 1 s", service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now.Add(time.Hour), EndsAt: now.Add(time.Hour + 7*24*time.Hour + time.Second)}, service.ErrValidation},
		{"motivo longo", service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Reason: &long}, service.ErrValidation},
		{"monitor inexistente", service.CreateMaintenanceInput{MonitorIDs: []string{"m", "nope"}, StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour)}, service.ErrNotFound},
	}
	for _, c := range cases {
		if _, err := svc.CreateBatch(context.Background(), c.in); !errors.Is(err, c.want) {
			t.Errorf("%s: err=%v, want %v", c.name, err, c.want)
		}
	}
	if n := count(); n != 0 {
		t.Fatalf("%d janelas criadas por pedidos inválidos", n)
	}

	// exatamente 7 dias é aceito; ids repetidos viram uma janela só
	reason := "  deploy  "
	ws, err := svc.CreateBatch(context.Background(), service.CreateMaintenanceInput{
		MonitorIDs: []string{"m", "m"}, StartsAt: now.Add(time.Hour), EndsAt: now.Add(time.Hour + 7*24*time.Hour), Reason: &reason,
	})
	if err != nil || len(ws) != 1 || *ws[0].Reason != "deploy" {
		t.Fatalf("7 dias: err=%v windows=%d", err, len(ws))
	}
}

// Tudo ou nada (D5): com um monitor em conflito, nenhum é criado e o erro diz qual.
func TestMaintenanceBatchConflict(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, count := newMaintenanceService(t, now)
	ctx := context.Background()
	in := service.CreateMaintenanceInput{MonitorIDs: []string{"m"}, StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour)}
	if _, err := svc.CreateBatch(ctx, in); err != nil {
		t.Fatal(err)
	}

	in.MonitorIDs = []string{"m2", "m"}
	_, err := svc.CreateBatch(ctx, in)
	var conflict *service.MaintenanceConflictError
	if !errors.As(err, &conflict) || !errors.Is(err, service.ErrMaintenanceConflict) {
		t.Fatalf("err=%v, want MaintenanceConflictError", err)
	}
	if len(conflict.MonitorIDs) != 1 || conflict.MonitorIDs[0] != "m" {
		t.Fatalf("conflito=%v, want [m]", conflict.MonitorIDs)
	}
	if n := count(); n != 1 {
		t.Fatalf("janelas=%d, want 1 (nada criado para m2)", n)
	}
}

// janela gravada direto no repositório, sem a regra de "início no futuro"
func addWindowRaw(t *testing.T, db *sql.DB, id, monitorID string, starts, ends time.Time) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO maintenance_windows (id, monitor_id, starts_at, ends_at, reason, created_at) VALUES (?,?,?,?,?,?)`,
		id, monitorID, starts.UTC().Format(time.RFC3339), ends.UTC().Format(time.RFC3339), "antigo", starts.UTC().Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}
}

func newMaintenanceServiceDB(t *testing.T, now *time.Time) (*service.MaintenanceService, *sql.DB) {
	t.Helper()
	_, db := newStatsService(t)
	if _, err := db.Exec(`INSERT INTO monitors (id, name, url, interval_seconds, created_at, updated_at)
		VALUES ('m2','Banco','http://y',60,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	svc := service.NewMaintenanceService(repository.NewMaintenanceRepo(db), repository.NewMonitorRepo(db)).
		WithClock(func() time.Time { return *now })
	return svc, db
}

func TestMaintenanceList(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, db := newMaintenanceServiceDB(t, &now)
	ctx := context.Background()
	addWindowRaw(t, db, "late", "m", now.Add(2*time.Hour), now.Add(3*time.Hour))
	addWindowRaw(t, db, "early", "m2", now.Add(time.Hour), now.Add(2*time.Hour))

	views, err := svc.List(ctx)
	if err != nil || len(views) != 2 || views[0].ID != "early" || views[0].MonitorName != "Banco" {
		t.Fatalf("lista: err=%v %+v", err, views)
	}
	db.Exec(`UPDATE maintenance_windows SET finalized_at = ? WHERE id = 'late'`, now.Format(time.RFC3339))
	if views, _ = svc.List(ctx); len(views) != 1 {
		t.Fatalf("finalizada continua na lista: %d", len(views))
	}
}

func TestMaintenanceUpdateScheduled(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, db := newMaintenanceServiceDB(t, &now)
	ctx := context.Background()
	addWindowRaw(t, db, "w", "m", now.Add(time.Hour), now.Add(2*time.Hour))

	starts, ends, reason := now.Add(3*time.Hour), now.Add(5*time.Hour), "  novo  "
	w, err := svc.Update(ctx, "w", service.UpdateMaintenanceInput{StartsAt: &starts, EndsAt: &ends, Reason: &reason, ReasonSet: true})
	if err != nil || !w.StartsAt.Equal(starts) || !w.EndsAt.Equal(ends) || *w.Reason != "novo" {
		t.Fatalf("editar tudo: err=%v %+v", err, w)
	}
	// só o fim: início e motivo mantidos
	ends2 := now.Add(4 * time.Hour)
	if w, err = svc.Update(ctx, "w", service.UpdateMaintenanceInput{EndsAt: &ends2}); err != nil || !w.StartsAt.Equal(starts) || *w.Reason != "novo" {
		t.Fatalf("só o fim: err=%v %+v", err, w)
	}
	// limpar o motivo
	if w, err = svc.Update(ctx, "w", service.UpdateMaintenanceInput{ReasonSet: true}); err != nil || w.Reason != nil {
		t.Fatalf("limpar motivo: err=%v %+v", err, w)
	}
	var stored string
	db.QueryRow(`SELECT ends_at FROM maintenance_windows WHERE id='w'`).Scan(&stored)
	if stored != ends2.Format(time.RFC3339) {
		t.Fatalf("ends_at gravado=%s", stored)
	}

	long := strings.Repeat("a", 201)
	past, before, tooLong := now.Add(-time.Minute), now.Add(2*time.Hour), now.Add(3*time.Hour+7*24*time.Hour+time.Second)
	for name, in := range map[string]service.UpdateMaintenanceInput{
		"corpo vazio":         {},
		"início no passado":   {StartsAt: &past},
		"fim antes do início": {EndsAt: &before},
		"7 dias + 1 s":        {EndsAt: &tooLong},
		"motivo longo":        {Reason: &long, ReasonSet: true},
	} {
		if _, err := svc.Update(ctx, "w", in); !errors.Is(err, service.ErrValidation) {
			t.Errorf("%s: err=%v, want ErrValidation", name, err)
		}
	}
	if _, err := svc.Update(ctx, "nope", service.UpdateMaintenanceInput{EndsAt: &ends}); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("inexistente: %v", err)
	}
}

// D12: iniciada (pelo relógio, mesmo antes do tick) trava o horário; o motivo
// muda sem mexer no monitor nem no incidente.
func TestMaintenanceUpdateStartedLocked(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, db := newMaintenanceServiceDB(t, &now)
	ctx := context.Background()
	addWindowRaw(t, db, "w", "m", now, now.Add(time.Hour)) // começa neste segundo exato
	db.Exec(`UPDATE monitors SET current_state='UNKNOWN', resumed_at='2026-11-01T00:00:00Z' WHERE id='m'`)
	db.Exec(`INSERT INTO incidents (id, monitor_id, started_at, ended_at, closed_manually, close_note) VALUES ('i','m','2026-11-10T11:00:00Z','2026-11-10T12:00:00Z',2,'antigo')`)

	ends := now.Add(2 * time.Hour)
	for _, in := range []service.UpdateMaintenanceInput{{EndsAt: &ends}, {StartsAt: &ends}} {
		if _, err := svc.Update(ctx, "w", in); !errors.Is(err, service.ErrMaintenanceLocked) {
			t.Fatalf("horário em iniciada: err=%v, want ErrMaintenanceLocked", err)
		}
	}
	var gotEnds string
	db.QueryRow(`SELECT ends_at FROM maintenance_windows WHERE id='w'`).Scan(&gotEnds)
	if gotEnds != now.Add(time.Hour).Format(time.RFC3339) {
		t.Fatalf("ends_at mudou: %s", gotEnds)
	}

	reason := "motivo novo"
	w, err := svc.Update(ctx, "w", service.UpdateMaintenanceInput{Reason: &reason, ReasonSet: true})
	if err != nil || *w.Reason != "motivo novo" {
		t.Fatalf("motivo em iniciada: err=%v", err)
	}
	var state, resumed, note string
	db.QueryRow(`SELECT current_state, resumed_at FROM monitors WHERE id='m'`).Scan(&state, &resumed)
	db.QueryRow(`SELECT close_note FROM incidents WHERE id='i'`).Scan(&note)
	if state != "UNKNOWN" || resumed != "2026-11-01T00:00:00Z" || note != "antigo" {
		t.Fatalf("efeitos reaplicados: state=%s resumed=%s note=%s", state, resumed, note)
	}

	// finalizada: 404
	db.Exec(`UPDATE maintenance_windows SET finalized_at = ? WHERE id='w'`, now.Format(time.RFC3339))
	if _, err := svc.Update(ctx, "w", service.UpdateMaintenanceInput{Reason: &reason, ReasonSet: true}); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("finalizada: %v", err)
	}
}

// A guarda starts_at > now no UPDATE fecha a corrida: se a janela começa entre a
// leitura e a gravação, o horário não muda e o erro é ErrMaintenanceLocked.
func TestMaintenanceUpdateRaceWithStart(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	_, db := newMaintenanceServiceDB(t, &now)
	repo := repository.NewMaintenanceRepo(db)
	addWindowRaw(t, db, "w", "m", now.Add(time.Minute), now.Add(time.Hour))
	ok, err := repo.UpdateTimes(context.Background(), "w", now.Add(2*time.Minute), now.Add(2*time.Hour), nil, now.Add(time.Minute))
	if err != nil || ok {
		t.Fatalf("UpdateTimes com a janela já iniciada: ok=%v err=%v, want false", ok, err)
	}
	if _, err := repo.UpdateTimes(context.Background(), "w", now.Add(2*time.Hour), now.Add(time.Hour), nil, now); err == nil {
		t.Fatal("UpdateTimes com fim antes do início passou pelo CHECK")
	}
}

// O gate e o tick leem o horário da linha: depois de editar, vale o intervalo novo.
func TestMaintenanceEditedIntervalIsUsed(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	svc, db := newMaintenanceServiceDB(t, &now)
	ctx := context.Background()
	repo := repository.NewMaintenanceRepo(db)
	addWindowRaw(t, db, "w", "m", now.Add(time.Hour), now.Add(2*time.Hour))
	starts, ends := now.Add(3*time.Hour), now.Add(4*time.Hour)
	if _, err := svc.Update(ctx, "w", service.UpdateMaintenanceInput{StartsAt: &starts, EndsAt: &ends}); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		at   time.Time
		want bool
	}{{now.Add(90 * time.Minute), false}, {starts, true}} {
		if got, _ := repo.ActiveForMonitor(ctx, "m", c.at); got != c.want {
			t.Errorf("ativa em %s = %v, want %v", c.at.Format(time.Kitchen), got, c.want)
		}
		due, _ := repo.DueStart(ctx, c.at)
		if (len(due) == 1) != c.want {
			t.Errorf("DueStart em %s = %d janelas", c.at.Format(time.Kitchen), len(due))
		}
	}
}
