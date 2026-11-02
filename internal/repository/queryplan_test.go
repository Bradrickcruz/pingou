package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

// hotQuery é uma consulta quente com argumentos de exemplo, usada pelo teste
// de plano (banco pequeno) e pela impressão de planos no banco de benchmark.
type hotQuery struct {
	name string
	sql  string
	args []any
}

func hotQueries(monitorID string, all []string, now time.Time) []hotQuery {
	ts := func(t time.Time) string { return t.UTC().Format(time.RFC3339) }
	ids := make([]any, len(all))
	for i, id := range all {
		ids[i] = id
	}
	countAll, listAll, argsAll := incidentsListSQL(domain.IncidentFilter{})
	countOpen, listOpen, argsOpen := incidentsListSQL(domain.IncidentFilter{Status: domain.IncidentStatusOpen})
	countRes, listRes, argsRes := incidentsListSQL(domain.IncidentFilter{Status: domain.IncidentStatusResolved})
	_, listMon, argsMon := incidentsListSQL(domain.IncidentFilter{MonitorID: monitorID})
	return []hotQuery{
		{"stats: buckets 30d", sqlChecksBuckets, []any{21600, monitorID, ts(now.Add(-30 * 24 * time.Hour)), ts(now)}},
		{"stats: down window", sqlIncidentsInWindow, []any{monitorID, ts(now.Add(-30 * 24 * time.Hour))}},
		{"lista: history", sqlSummariesHistory(len(ids)), append(append([]any{}, ids...), 30)},
		{"lista: uptime 24h", sqlSummariesUptime(len(ids)), append(append([]any{}, ids...), ts(now.Add(-24*time.Hour)))},
		{"lista: down_since", sqlIncidentsOpenByMonitors(len(ids)), ids},
		{"checks: página", sqlChecksPage, []any{monitorID, 50, 0}},
		{"checks: total", sqlChecksCount, []any{monitorID}},
		{"retenção: delete", sqlChecksDeleteOlder, []any{ts(now.Add(-29 * 24 * time.Hour))}},
		{"state machine: aberto", sqlIncidentOpenByMonitor, []any{monitorID}},
		{"state machine: +falha", sqlIncidentIncrementFailures, []any{monitorID}},
		{"incidentes: total", countAll, argsAll},
		{"incidentes: lista", listAll, append(argsAll, 20, 0)},
		{"incidentes: abertos total", countOpen, argsOpen},
		{"incidentes: abertos", listOpen, append(argsOpen, 20, 0)},
		{"incidentes: resolvidos total", countRes, argsRes},
		{"incidentes: resolvidos", listRes, append(argsRes, 20, 0)},
		{"incidentes: do monitor", listMon, append(argsMon, 20, 0)},
	}
}

// explain devolve as linhas de detalhe do EXPLAIN QUERY PLAN.
func explain(ctx context.Context, db *sql.DB, q string, args ...any) ([]string, error) {
	rows, err := db.QueryContext(ctx, "EXPLAIN QUERY PLAN "+q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id, parent, notused int
		var detail string
		if err := rows.Scan(&id, &parent, &notused, &detail); err != nil {
			return nil, err
		}
		out = append(out, detail)
	}
	return out, rows.Err()
}

// BenchDBPath devolve o banco de benchmark de PINGOU_BENCH_DB ou pula o teste.
// Caminho relativo vale a partir da raiz do repositório (go test roda na
// pasta do pacote). Exportado para os benchmarks do pacote repository_test.
func BenchDBPath(t testing.TB) string {
	path := os.Getenv("PINGOU_BENCH_DB")
	if path == "" {
		t.Skip("PINGOU_BENCH_DB não definido")
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join("..", "..", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("PINGOU_BENCH_DB: %v", err)
	}
	return path
}

// TestPrintQueryPlans imprime o plano de cada consulta quente no banco de
// benchmark (PINGOU_BENCH_DB). Só lê: abre em modo somente leitura.
func TestPrintQueryPlans(t *testing.T) {
	path := BenchDBPath(t)
	db, err := sql.Open("sqlite3", "file:"+path+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()

	var all []string
	rows, err := db.QueryContext(ctx, `SELECT id FROM monitors ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		all = append(all, id)
	}
	rows.Close()
	var last string
	db.QueryRowContext(ctx, `SELECT MAX(checked_at) FROM checks`).Scan(&last)
	now, _ := time.Parse(time.RFC3339, last)

	var b strings.Builder
	b.WriteString("\n| Consulta | Plano |\n|---|---|\n")
	for _, q := range hotQueries(all[0], all, now) {
		plan, err := explain(ctx, db, q.sql, q.args...)
		if err != nil {
			t.Fatalf("%s: %v", q.name, err)
		}
		fmt.Fprintf(&b, "| %s | %s |\n", q.name, strings.Join(plan, " · "))
	}
	fmt.Print(b.String())
}
