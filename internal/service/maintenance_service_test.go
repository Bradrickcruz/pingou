package service_test

import (
	"context"
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
