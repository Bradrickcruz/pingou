package handler

import (
	"net/http"
	"slices"
	"testing"
)

type incidentsPage struct {
	Data  []incidentResponse `json:"data"`
	Total int                `json:"total"`
}

// seedIncidents: monitor A (http) com 1 aberto e 2 resolvidos; B (tcp) com
// 1 aberto e 1 resolvido. Resolvidos terminam em ordem diferente da de início.
func seedIncidents(t *testing.T, s *Server) {
	t.Helper()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := s.db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	now := "2026-10-28T00:00:00Z"
	exec(`INSERT INTO monitors (id, name, type, url, interval_seconds, created_at, updated_at)
		VALUES ('a','API','http','https://api.test',60,?,?)`, now, now)
	exec(`INSERT INTO monitors (id, name, type, config, interval_seconds, created_at, updated_at)
		VALUES ('b','Postgres','tcp','{"host":"db.interno","port":5432}',60,?,?)`, now, now)

	for _, r := range [][]any{
		{"a-open", "a", "2026-10-28T10:00:00Z", nil, 7},
		{"a-r1", "a", "2026-10-27T08:00:00Z", "2026-10-27T08:10:00Z", 3},
		{"a-r2", "a", "2026-10-26T08:00:00Z", "2026-10-27T09:00:00Z", 4}, // começou antes, terminou depois
		{"b-open", "b", "2026-10-28T11:00:00Z", nil, 3},
		{"b-r1", "b", "2026-10-25T08:00:00Z", "2026-10-25T08:05:00Z", 3},
	} {
		exec(`INSERT INTO incidents (id, monitor_id, started_at, ended_at, last_error, failure_count)
			VALUES (?,?,?,?,'dial tcp: connection refused',?)`, r...)
	}
}

func listIncidents(t *testing.T, h http.Handler, path string) incidentsPage {
	t.Helper()
	rec := do(t, h, "GET", path, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("%s: %d %s", path, rec.Code, rec.Body)
	}
	return decode[incidentsPage](t, rec)
}

func ids(p incidentsPage) []string {
	out := make([]string, len(p.Data))
	for i, d := range p.Data {
		out[i] = d.ID
	}
	return out
}

func TestListIncidentsTotalIgnoresLimit(t *testing.T) {
	s, h, _ := newTestServer(t)
	seedIncidents(t, s)

	// regressão: o total era len(página), e a sidebar (limit=1) mostrava no máximo 1
	p := listIncidents(t, h, "/api/incidents?status=open&limit=1")
	if p.Total != 2 || len(p.Data) != 1 {
		t.Fatalf("total=%d itens=%d, want 2 e 1", p.Total, len(p.Data))
	}
	inc := p.Data[0]
	if inc.ID != "b-open" || inc.Monitor == nil || inc.Monitor.Name != "Postgres" ||
		inc.Monitor.Type != "tcp" || inc.Monitor.Target == nil || *inc.Monitor.Target != "db.interno:5432" {
		t.Fatalf("incidente: %+v monitor=%+v", inc, inc.Monitor)
	}
	if inc.FailureCount != 3 || !inc.Open {
		t.Fatalf("failure_count=%d open=%v", inc.FailureCount, inc.Open)
	}

	// ?open=true continua aceito
	if p := listIncidents(t, h, "/api/incidents?open=true"); p.Total != 2 {
		t.Fatalf("open=true: total=%d, want 2", p.Total)
	}
	if p := listIncidents(t, h, "/api/incidents"); p.Total != 5 {
		t.Fatalf("todos: total=%d, want 5", p.Total)
	}
}

func TestListIncidentsResolvedOrderAndFilters(t *testing.T) {
	s, h, _ := newTestServer(t)
	seedIncidents(t, s)

	p := listIncidents(t, h, "/api/incidents?status=resolved")
	if got, want := ids(p), []string{"a-r2", "a-r1", "b-r1"}; p.Total != 3 || !slices.Equal(got, want) {
		t.Fatalf("resolvidos por ended_at DESC: %v (total %d), want %v", got, p.Total, want)
	}
	if p.Data[0].Monitor.Target == nil || *p.Data[0].Monitor.Target != "https://api.test" {
		t.Fatalf("target http: %+v", p.Data[0].Monitor)
	}

	p = listIncidents(t, h, "/api/incidents?status=resolved&monitor_id=b")
	if got := ids(p); p.Total != 1 || !slices.Equal(got, []string{"b-r1"}) {
		t.Fatalf("monitor b resolvidos: %v (total %d)", got, p.Total)
	}

	// rota por monitor usa o mesmo código
	p = listIncidents(t, h, "/api/monitors/a/incidents?status=open")
	if got := ids(p); p.Total != 1 || !slices.Equal(got, []string{"a-open"}) {
		t.Fatalf("/monitors/a/incidents: %v (total %d)", got, p.Total)
	}

	if rec := do(t, h, "GET", "/api/incidents?status=xyz", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("status inválido: %d, want 400", rec.Code)
	}
}
