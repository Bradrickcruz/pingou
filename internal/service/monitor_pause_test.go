package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/service"
)

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
