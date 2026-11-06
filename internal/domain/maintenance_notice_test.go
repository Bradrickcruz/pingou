package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMaintenanceGroupKey(t *testing.T) {
	s := time.Date(2026, 11, 8, 2, 0, 0, 0, time.UTC)
	e := s.Add(2 * time.Hour)
	str := func(v string) *string { return &v }
	w := func(reason *string) *MaintenanceWindow {
		return &MaintenanceWindow{StartsAt: s, EndsAt: e, Reason: reason}
	}

	if MaintenanceGroupKey(w(str("a"))) != MaintenanceGroupKey(w(str("a"))) {
		t.Fatal("mesmos valores devem dar a mesma chave")
	}
	if MaintenanceGroupKey(w(str("a"))) == MaintenanceGroupKey(w(str("b"))) {
		t.Fatal("motivos diferentes devem dar chaves diferentes")
	}
	if MaintenanceGroupKey(w(nil)) == MaintenanceGroupKey(w(str(""))) {
		t.Fatal("motivo nil deve ser diferente de motivo vazio")
	}
	other := w(str("a"))
	other.EndsAt = e.Add(time.Second)
	if MaintenanceGroupKey(w(str("a"))) == MaintenanceGroupKey(other) {
		t.Fatal("fim diferente deve dar chave diferente")
	}
	// o fuso não muda o instante
	local := w(str("a"))
	local.StartsAt = s.In(time.FixedZone("x", -3*3600))
	if MaintenanceGroupKey(w(str("a"))) != MaintenanceGroupKey(local) {
		t.Fatal("mesmo instante em outro fuso deve dar a mesma chave")
	}
}

func TestMaintenanceSnapshotRoundTrip(t *testing.T) {
	at := time.Date(2026, 11, 8, 2, 0, 4, 0, time.UTC)
	reason := "troca de disco"
	ws := []*MaintenanceWindow{
		{MonitorID: "m1", StartsAt: at.Add(-4 * time.Second), EndsAt: at.Add(2 * time.Hour), Reason: &reason},
		{MonitorID: "m2", StartsAt: at.Add(-4 * time.Second), EndsAt: at.Add(2 * time.Hour), Reason: &reason},
	}
	ms := []*Monitor{
		{ID: "m1", Name: "API", Type: TypeHTTP, URL: "https://api.test", PingToken: &reason},
		{ID: "m2", Name: "DB", Type: TypeTCP, Config: MonitorConfig{Host: "db", Port: 5432}},
	}
	snap := NewMaintenanceSnapshot(ws, ms, at)
	if !snap.OccurredAt.Equal(at) {
		t.Fatalf("OccurredAt=%v, want o momento do tick", snap.OccurredAt)
	}
	raw, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	var back EventSnapshot
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatal(err)
	}
	mt := back.Maintenance
	if mt == nil || len(mt.Monitors) != 2 || mt.Monitors[1].Port != 5432 || mt.Reason == nil || *mt.Reason != reason ||
		!mt.EndsAt.Equal(ws[0].EndsAt) {
		t.Fatalf("ida e volta: %+v", back)
	}
}

// Snapshots gravados antes do campo existir continuam lendo.
func TestOldSnapshotWithoutMaintenanceStillReads(t *testing.T) {
	old := `{"monitor":{"id":"m1","name":"API","type":"http","url":"https://x"},"occurred_at":"2026-01-01T12:00:00Z"}`
	var s EventSnapshot
	if err := json.Unmarshal([]byte(old), &s); err != nil {
		t.Fatal(err)
	}
	if s.Maintenance != nil || s.Monitor.Name != "API" {
		t.Fatalf("snapshot antigo: %+v", s)
	}
	raw, _ := json.Marshal(s)
	if string(raw) == "" || contains(string(raw), "maintenance") {
		t.Fatalf("evento comum não deve gravar o campo maintenance: %s", raw)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
