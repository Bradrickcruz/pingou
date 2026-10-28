package database

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
)

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func setupGoose(t *testing.T) {
	t.Helper()
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}
}

// Regressão de risco: recriar monitors com as FKs ligadas apagaria checks e
// incidents via ON DELETE CASCADE.
func TestMonitorTypesMigrationKeepsHistory(t *testing.T) {
	db, err := OpenDB(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	setupGoose(t)

	if err := goose.UpTo(db, "migrations", 2); err != nil {
		t.Fatal(err)
	}
	now := "2026-01-01T00:00:00Z"
	mustExec(t, db, `INSERT INTO monitors (id, name, url, interval_seconds, timeout_seconds, failure_threshold, enabled, current_state, created_at, updated_at)
		VALUES ('m1','API','https://x.test',60,5,3,1,'DOWN',?,?)`, now, now)
	for range 3 {
		mustExec(t, db, `INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES ('m1',0,10,?)`, now)
	}
	mustExec(t, db, `INSERT INTO incidents (id, monitor_id, started_at) VALUES ('i1','m1',?)`, now)

	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	if c, i := countRows(t, db, "checks"), countRows(t, db, "incidents"); c != 3 || i != 1 {
		t.Fatalf("após up: checks=%d incidents=%d, want 3 e 1", c, i)
	}
	var typ, state string
	if err := db.QueryRow(`SELECT type, current_state FROM monitors WHERE id='m1'`).Scan(&typ, &state); err != nil {
		t.Fatal(err)
	}
	if typ != "http" || state != "DOWN" {
		t.Fatalf("type=%s state=%s, want http e DOWN", typ, state)
	}
	assertNoFKViolations(t, db)
	assertFKOn(t, db)

	// heartbeat de 1 semana passa a ser aceito
	mustExec(t, db, `INSERT INTO monitors (id, name, type, ping_token, interval_seconds, created_at, updated_at)
		VALUES ('m2','Backup','heartbeat','tok',604800,?,?)`, now, now)

	// down bloqueia com monitor não-http e não altera nada
	if err := Down(db); err == nil {
		t.Fatal("down com heartbeat deveria falhar")
	}
	mustExec(t, db, `DELETE FROM monitors WHERE id='m2'`)

	if err := Down(db); err != nil {
		t.Fatal(err)
	}
	if c := countRows(t, db, "checks"); c != 3 {
		t.Fatalf("após down: checks=%d, want 3", c)
	}
	assertNoFKViolations(t, db)

	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	if c := countRows(t, db, "checks"); c != 3 {
		t.Fatalf("após up de novo: checks=%d, want 3", c)
	}
}

func TestMonitorTypesConstraints(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := "2026-01-01T00:00:00Z"

	bad := []string{
		// http sem url
		`INSERT INTO monitors (id, name, type, interval_seconds, created_at, updated_at) VALUES ('a','x','http',60,'` + now + `','` + now + `')`,
		// heartbeat sem token
		`INSERT INTO monitors (id, name, type, interval_seconds, created_at, updated_at) VALUES ('b','x','heartbeat',60,'` + now + `','` + now + `')`,
		// token fora de heartbeat
		`INSERT INTO monitors (id, name, type, ping_token, interval_seconds, created_at, updated_at) VALUES ('c','x','tcp','t',60,'` + now + `','` + now + `')`,
		// tipo desconhecido
		`INSERT INTO monitors (id, name, type, interval_seconds, created_at, updated_at) VALUES ('d','x','dns',60,'` + now + `','` + now + `')`,
	}
	for _, q := range bad {
		if _, err := db.Exec(q); err == nil {
			t.Errorf("insert aceito, want erro de CHECK: %s", q)
		}
	}
}

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatal(err)
	}
}

func assertNoFKViolations(t *testing.T, db *sql.DB) {
	t.Helper()
	rows, err := db.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("PRAGMA foreign_key_check encontrou violações")
	}
}

func assertFKOn(t *testing.T, db *sql.DB) {
	t.Helper()
	var on int
	if err := db.QueryRow(`PRAGMA foreign_keys`).Scan(&on); err != nil {
		t.Fatal(err)
	}
	if on != 1 {
		t.Fatal("foreign_keys ficou desligado após a migration")
	}
}
