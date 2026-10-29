package repository_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

func TestDeleteResolvedOlderThanKeepsOpen(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO monitors (id, name, url, interval_seconds, created_at, updated_at)
		VALUES ('m','x','http://x',60,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`)
	exec(`INSERT INTO incidents (id, monitor_id, started_at, ended_at) VALUES
		('old-resolved','m','2026-01-01T00:00:00Z','2026-01-01T01:00:00Z'),
		('new-resolved','m','2026-03-01T00:00:00Z','2026-03-01T01:00:00Z'),
		('old-open','m','2026-01-02T00:00:00Z',NULL)`)

	deleted, err := repository.NewIncidentRepo(db).DeleteResolvedOlderThan(context.Background(), "2026-02-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 1 {
		t.Fatalf("deleted=%d, want 1", deleted)
	}
	rows, err := db.Query(`SELECT id FROM incidents ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var left []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		left = append(left, id)
	}
	if len(left) != 2 || left[0] != "new-resolved" || left[1] != "old-open" {
		t.Fatalf("restaram %v, want [new-resolved old-open]", left)
	}
}
