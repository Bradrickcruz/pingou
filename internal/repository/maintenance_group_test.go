package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

var groupNow = time.Date(2026, 11, 8, 2, 0, 5, 0, time.UTC)

// seedGroup cria n monitores, cada um com uma janela (mesmo início, fim e motivo)
// e um incidente aberto.
func seedGroup(t *testing.T, db *sql.DB, n int, starts, ends time.Time) []*domain.MaintenanceWindow {
	t.Helper()
	reason := "troca de disco"
	var ws []*domain.MaintenanceWindow
	for i := 0; i < n; i++ {
		id := string(rune('a' + i))
		if err := repository.NewMonitorRepo(db).Create(context.Background(), &domain.Monitor{
			ID: "m" + id, Name: "mon-" + id, Type: domain.TypeHTTP, URL: "https://" + id + ".test",
			IntervalSeconds: 60, TimeoutSeconds: 10, FailureThreshold: 1, Enabled: true,
			CurrentState: domain.StateDown, CreatedAt: groupNow, UpdatedAt: groupNow,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`INSERT INTO incidents (id, monitor_id, started_at) VALUES (?,?,?)`,
			"i"+id, "m"+id, fmtT(groupNow.Add(-time.Hour))); err != nil {
			t.Fatal(err)
		}
		ws = append(ws, &domain.MaintenanceWindow{
			ID: "w" + id, MonitorID: "m" + id, StartsAt: starts, EndsAt: ends, Reason: &reason, CreatedAt: groupNow.Add(-time.Hour),
		})
	}
	if err := repository.NewMaintenanceRepo(db).CreateBatch(context.Background(), ws); err != nil {
		t.Fatal(err)
	}
	return ws
}

func fmtT(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func noticeBuilder(eventType string) domain.MaintenanceNoticeBuilder {
	return func(applied []*domain.MaintenanceWindow, monitors []*domain.Monitor) ([]domain.OutboxEvent, error) {
		snap := domain.NewMaintenanceSnapshot(applied, monitors, groupNow)
		return []domain.OutboxEvent{{
			EventID: "ev-" + eventType, EventType: eventType, Channel: domain.ChannelWebhook, MonitorID: applied[0].MonitorID,
			Snapshot: snap, CreatedAt: groupNow, ExpiresAt: groupNow.Add(time.Hour), NextAttemptAt: groupNow,
		}}, nil
	}
}

func count(t *testing.T, db *sql.DB, q string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestStartGroupAppliesAndEnqueuesOnce(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewMaintenanceRepo(db)
	ws := seedGroup(t, db, 3, groupNow.Add(-5*time.Second), groupNow.Add(time.Hour))

	applied, err := repo.StartGroup(ctx, ws, groupNow, noticeBuilder(domain.EventMaintenanceStarted))
	if err != nil || len(applied) != 3 {
		t.Fatalf("applied=%d err=%v, want 3", len(applied), err)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM notification_outbox WHERE event_type = ?`, domain.EventMaintenanceStarted); n != 1 {
		t.Fatalf("linhas no outbox=%d, want 1", n)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM incidents WHERE ended_at IS NOT NULL`); n != 3 {
		t.Fatalf("incidentes fechados=%d, want 3", n)
	}
	// o payload leva os 3 monitores, lidos pelo tx
	ev := pendingFor(t, db)
	if len(ev.Snapshot.Maintenance.Monitors) != 3 {
		t.Fatalf("monitores no aviso=%d, want 3", len(ev.Snapshot.Maintenance.Monitors))
	}

	// segunda chamada: nada aplicado, nada enfileirado
	applied, err = repo.StartGroup(ctx, ws, groupNow.Add(10*time.Second), noticeBuilder(domain.EventMaintenanceStarted))
	if err != nil || len(applied) != 0 {
		t.Fatalf("2ª chamada: applied=%d err=%v, want 0", len(applied), err)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM notification_outbox`); n != 1 {
		t.Fatalf("linhas no outbox=%d, want 1", n)
	}
}

func pendingFor(t *testing.T, db *sql.DB) domain.OutboxEvent {
	t.Helper()
	all, err := repository.NewOutboxRepo(db).Due(context.Background(), groupNow.Add(24*time.Hour), 10)
	if err != nil || len(all) == 0 {
		t.Fatalf("outbox vazio: %v", err)
	}
	return all[0]
}

func TestStartGroupBuildErrorRollsBackEverything(t *testing.T) {
	db := newOutboxDB(t)
	ws := seedGroup(t, db, 2, groupNow.Add(-5*time.Second), groupNow.Add(time.Hour))
	boom := errors.New("boom")

	_, err := repository.NewMaintenanceRepo(db).StartGroup(context.Background(), ws, groupNow,
		func([]*domain.MaintenanceWindow, []*domain.Monitor) ([]domain.OutboxEvent, error) { return nil, boom })
	if !errors.Is(err, boom) {
		t.Fatalf("err=%v, want boom", err)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM maintenance_windows WHERE started_at IS NOT NULL`); n != 0 {
		t.Fatalf("started_at gravado em %d janelas, want 0", n)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM incidents WHERE ended_at IS NOT NULL`); n != 0 {
		t.Fatalf("incidentes fechados=%d, want 0", n)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM notification_outbox`); n != 0 {
		t.Fatalf("outbox=%d, want 0", n)
	}
}

// Janela adiada entre a leitura do tick e a transação não inicia no horário velho.
func TestStartGroupSkipsWindowPostponedAfterDue(t *testing.T) {
	db := newOutboxDB(t)
	repo := repository.NewMaintenanceRepo(db)
	ws := seedGroup(t, db, 1, groupNow.Add(-5*time.Second), groupNow.Add(time.Hour))
	// o tick leu ws (vencida); depois o usuário adiou
	if _, err := db.Exec(`UPDATE maintenance_windows SET starts_at = ?, ends_at = ?`,
		fmtT(groupNow.Add(30*time.Minute)), fmtT(groupNow.Add(2*time.Hour))); err != nil {
		t.Fatal(err)
	}
	applied, err := repo.StartGroup(context.Background(), ws, groupNow, noticeBuilder(domain.EventMaintenanceStarted))
	if err != nil || len(applied) != 0 {
		t.Fatalf("applied=%d err=%v, want 0", len(applied), err)
	}
	if n := count(t, db, `SELECT COUNT(*) FROM notification_outbox`); n != 0 {
		t.Fatalf("outbox=%d, want 0", n)
	}
}

func TestFinalizeGroupNoticeOnlyForStartedWindows(t *testing.T) {
	db := newOutboxDB(t)
	ctx := context.Background()
	repo := repository.NewMaintenanceRepo(db)
	ws := seedGroup(t, db, 2, groupNow.Add(-2*time.Hour), groupNow.Add(-time.Minute))
	// só a primeira já tinha começado
	if _, err := db.Exec(`UPDATE maintenance_windows SET started_at = ? WHERE id = 'wa'`, fmtT(groupNow.Add(-2*time.Hour))); err != nil {
		t.Fatal(err)
	}

	applied, err := repo.FinalizeGroup(ctx, ws, groupNow, nil, noticeBuilder(domain.EventMaintenanceEnded))
	if err != nil || len(applied) != 2 {
		t.Fatalf("applied=%d err=%v, want 2", len(applied), err)
	}
	ev := pendingFor(t, db)
	if ms := ev.Snapshot.Maintenance.Monitors; len(ms) != 1 || ms[0].ID != "ma" {
		t.Fatalf("aviso de fim deve listar só a janela que já tinha começado: %+v", ms)
	}
	// a outra teve o início aplicado (incidente fechado) mesmo sem aviso; o
	// incidente da primeira ficou como estava (o início dela já constava)
	if n := count(t, db, `SELECT COUNT(*) FROM incidents WHERE ended_at IS NOT NULL AND monitor_id = 'mb'`); n != 1 {
		t.Fatalf("incidente da janela sem início: fechados=%d, want 1", n)
	}

	applied, err = repo.FinalizeGroup(ctx, ws, groupNow, nil, noticeBuilder(domain.EventMaintenanceEnded))
	if err != nil || len(applied) != 0 || count(t, db, `SELECT COUNT(*) FROM notification_outbox`) != 1 {
		t.Fatalf("2ª chamada: applied=%d err=%v", len(applied), err)
	}
}

// Janela que começou e terminou com o Pingou fora: nenhum aviso.
func TestFinalizeGroupNeverStartedHasNoNotice(t *testing.T) {
	db := newOutboxDB(t)
	repo := repository.NewMaintenanceRepo(db)
	ws := seedGroup(t, db, 1, groupNow.Add(-2*time.Hour), groupNow.Add(-time.Minute))
	built := false
	applied, err := repo.FinalizeGroup(context.Background(), ws, groupNow, nil,
		func([]*domain.MaintenanceWindow, []*domain.Monitor) ([]domain.OutboxEvent, error) {
			built = true
			return nil, nil
		})
	if err != nil || len(applied) != 1 || built {
		t.Fatalf("applied=%d err=%v built=%v", len(applied), err, built)
	}
}

func TestFinalizeGroupBuildErrorRollsBack(t *testing.T) {
	db := newOutboxDB(t)
	ws := seedGroup(t, db, 1, groupNow.Add(-2*time.Hour), groupNow.Add(-time.Minute))
	if _, err := db.Exec(`UPDATE maintenance_windows SET started_at = ?`, fmtT(groupNow.Add(-2*time.Hour))); err != nil {
		t.Fatal(err)
	}
	_, err := repository.NewMaintenanceRepo(db).FinalizeGroup(context.Background(), ws, groupNow, nil,
		func([]*domain.MaintenanceWindow, []*domain.Monitor) ([]domain.OutboxEvent, error) {
			return nil, errors.New("boom")
		})
	if err == nil {
		t.Fatal("want error")
	}
	if n := count(t, db, `SELECT COUNT(*) FROM maintenance_windows WHERE finalized_at IS NOT NULL`); n != 0 {
		t.Fatalf("finalized_at gravado em %d janelas, want 0", n)
	}
}

// Janela adiada depois da leitura do DueEnd não é finalizada.
func TestFinalizeGroupSkipsWindowExtendedAfterDue(t *testing.T) {
	db := newOutboxDB(t)
	ws := seedGroup(t, db, 1, groupNow.Add(-2*time.Hour), groupNow.Add(-time.Minute))
	if _, err := db.Exec(`UPDATE maintenance_windows SET ends_at = ?`, fmtT(groupNow.Add(time.Hour))); err != nil {
		t.Fatal(err)
	}
	applied, err := repository.NewMaintenanceRepo(db).FinalizeGroup(context.Background(), ws, groupNow, nil, nil)
	if err != nil || len(applied) != 0 {
		t.Fatalf("applied=%d err=%v, want 0", len(applied), err)
	}
}

// Cancelamento: ends_at passa a valer o novo, e o aviso sai se já tinha começado.
func TestFinalizeGroupCancelShortensWindow(t *testing.T) {
	db := newOutboxDB(t)
	repo := repository.NewMaintenanceRepo(db)
	ws := seedGroup(t, db, 1, groupNow.Add(-time.Hour), groupNow.Add(time.Hour))
	if _, err := db.Exec(`UPDATE maintenance_windows SET started_at = ?`, fmtT(groupNow.Add(-time.Hour))); err != nil {
		t.Fatal(err)
	}
	endsAt := groupNow
	applied, err := repo.FinalizeGroup(context.Background(), ws, groupNow, &endsAt, noticeBuilder(domain.EventMaintenanceEnded))
	if err != nil || len(applied) != 1 || !applied[0].EndsAt.Equal(endsAt) {
		t.Fatalf("applied=%+v err=%v", applied, err)
	}
	if ev := pendingFor(t, db); !ev.Snapshot.Maintenance.EndsAt.Equal(endsAt) {
		t.Fatalf("o aviso deve trazer o fim novo: %v", ev.Snapshot.Maintenance.EndsAt)
	}
}
