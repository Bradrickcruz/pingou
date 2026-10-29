package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// now fixo no meio de um balde de 10 min: until = 12:10, since = ontem 12:10
var statsNow = time.Date(2026, 10, 29, 12, 5, 0, 0, time.UTC)

func newStatsService(t *testing.T) (*service.MonitorService, *sql.DB) {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(`INSERT INTO monitors (id, name, url, interval_seconds, created_at, updated_at)
		VALUES ('m','x','http://x',60,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	return service.NewMonitorService(repository.NewMonitorRepo(db), repository.NewCheckRepo(db), repository.NewIncidentRepo(db)), db
}

func at(d time.Duration) string { return statsNow.Add(d).Format(time.RFC3339) }

func TestStatsBucketsAndTotals(t *testing.T) {
	svc, db := newStatsService(t)
	// 5 ok (100..500 ms) e 1 falha no balde 11:50–12:00
	for i, ms := range []int{100, 200, 300, 400, 500} {
		db.Exec(`INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES ('m',1,?,?)`,
			ms, at(-14*time.Minute+time.Duration(i)*time.Minute))
	}
	db.Exec(`INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES ('m',0,9999,?)`, at(-8*time.Minute))
	// fora da janela: não conta
	db.Exec(`INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES ('m',0,1,?)`, at(-25*time.Hour))

	st, err := svc.Stats(context.Background(), "m", "24h", statsNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Series) != 144 || st.BucketSeconds != 600 {
		t.Fatalf("series=%d bucket=%d, want 144/600", len(st.Series), st.BucketSeconds)
	}
	if want := time.Date(2026, 10, 29, 12, 10, 0, 0, time.UTC); !st.Until.Equal(want) || !st.Since.Equal(want.Add(-24*time.Hour)) {
		t.Fatalf("since=%s until=%s", st.Since, st.Until)
	}
	b := st.Series[142] // 11:50
	if !b.Start.Equal(time.Date(2026, 10, 29, 11, 50, 0, 0, time.UTC)) || b.Checks != 6 || b.Failures != 1 || b.AvgLatencyMs == nil || *b.AvgLatencyMs != 300 {
		t.Fatalf("balde 11:50 = %+v", b)
	}
	if st.Series[0].AvgLatencyMs != nil || st.Series[0].Checks != 0 || st.Series[143].AvgLatencyMs != nil {
		t.Fatalf("baldes vazios deveriam vir sem média: %+v %+v", st.Series[0], st.Series[143])
	}
	if st.Checks != 6 || st.Uptime == nil || *st.Uptime < 83.33 || *st.Uptime > 83.34 || *st.AvgLatencyMs != 300 {
		t.Fatalf("totais: checks=%d uptime=%v avg=%v", st.Checks, st.Uptime, st.AvgLatencyMs)
	}
}

func TestStatsIncidentsClipped(t *testing.T) {
	svc, db := newStatsService(t)
	db.Exec(`INSERT INTO incidents (id, monitor_id, started_at, ended_at) VALUES
		('antes','m',?,?), ('aberto','m',?,NULL), ('velho','m',?,?)`,
		at(-30*time.Hour), at(-23*time.Hour),
		at(-2*time.Hour),
		at(-72*time.Hour), at(-48*time.Hour))

	st, err := svc.Stats(context.Background(), "m", "24h", statsNow)
	if err != nil {
		t.Fatal(err)
	}
	if st.Incidents != 2 || len(st.Down) != 2 {
		t.Fatalf("incidents=%d down=%v, want 2", st.Incidents, st.Down)
	}
	// "antes": recortado em since (ontem 12:10) até 13:05 = 55 min; "aberto": 2 h até now
	if !st.Down[0].Start.Equal(st.Since) || st.Down[0].End.Sub(st.Down[0].Start) != 55*time.Minute {
		t.Fatalf("antes = %+v", st.Down[0])
	}
	if !st.Down[1].End.Equal(statsNow) {
		t.Fatalf("aberto termina em now: %+v", st.Down[1])
	}
	if want := int64((55*time.Minute + 2*time.Hour) / time.Second); st.DowntimeSeconds != want {
		t.Fatalf("downtime=%d, want %d", st.DowntimeSeconds, want)
	}
}

func TestStatsEmptyAndErrors(t *testing.T) {
	svc, _ := newStatsService(t)
	st, err := svc.Stats(context.Background(), "m", "30d", statsNow)
	if err != nil {
		t.Fatal(err)
	}
	if st.Uptime != nil || st.AvgLatencyMs != nil || st.Checks != 0 || len(st.Series) != 120 {
		t.Fatalf("vazio: %+v", st)
	}
	for rng, n := range map[string]int{"7d": 168, "14d": 112} {
		st, _ := svc.Stats(context.Background(), "m", rng, statsNow)
		if len(st.Series) != n {
			t.Fatalf("%s: %d pontos, want %d", rng, len(st.Series), n)
		}
	}
	if _, err := svc.Stats(context.Background(), "m", "1y", statsNow); !errors.Is(err, service.ErrValidation) {
		t.Fatalf("range inválido: %v", err)
	}
	if _, err := svc.Stats(context.Background(), "nope", "24h", statsNow); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("inexistente: %v", err)
	}
}
