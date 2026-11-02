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

	"github.com/Bradrickcruz/pingou/internal/database"
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
	hist := make([]any, 0, 2*len(all))
	for i, id := range all {
		ids[i] = id
		hist = append(hist, id, 30)
	}
	countAll, listAll, argsAll := incidentsListSQL(domain.IncidentFilter{})
	countOpen, listOpen, argsOpen := incidentsListSQL(domain.IncidentFilter{Status: domain.IncidentStatusOpen})
	countRes, listRes, argsRes := incidentsListSQL(domain.IncidentFilter{Status: domain.IncidentStatusResolved})
	_, listMon, argsMon := incidentsListSQL(domain.IncidentFilter{MonitorID: monitorID})
	return []hotQuery{
		{"stats: buckets 30d", sqlChecksBuckets, []any{21600, monitorID, ts(now.Add(-30 * 24 * time.Hour)), ts(now)}},
		{"stats: down window", sqlIncidentsInWindow, []any{monitorID, ts(now.Add(-30 * 24 * time.Hour))}},
		{"lista: history", sqlSummariesHistory(len(ids)), hist},
		{"lista: uptime 24h", sqlSummariesUptime(len(ids)), append(append([]any{}, ids...), ts(now.Add(-24*time.Hour)))},
		{"lista: down_since", sqlIncidentsOpenByMonitors(len(ids)), ids},
		{"checks: página", sqlChecksPage, []any{monitorID, 50, 0}},
		{"checks: total", sqlChecksCount, []any{monitorID}},
		{"retenção: delete", sqlChecksDeleteOlderBatch, []any{ts(now.Add(-29 * 24 * time.Hour)), 5000}},
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

	fmt.Print(QueryPlans(t, db, all, now))
}

// QueryPlans devolve a tabela de planos das consultas quentes no banco db
// (usado também pelo TestDBBench, na cópia já migrada). O histórico da lista
// usa só 2 monitores: o plano se repete igual para cada um.
func QueryPlans(t testing.TB, db *sql.DB, monitorIDs []string, now time.Time) string {
	ctx := context.Background()
	var b strings.Builder
	b.WriteString("\n| Consulta | Plano |\n|---|---|\n")
	for _, q := range hotQueries(monitorIDs[0], monitorIDs, now) {
		if q.name == "lista: history" {
			q.sql, q.args = sqlSummariesHistory(2), q.args[:4]
		}
		plan, err := explain(ctx, db, q.sql, q.args...)
		if err != nil {
			t.Fatalf("%s: %v", q.name, err)
		}
		fmt.Fprintf(&b, "| %s | %s |\n", q.name, strings.Join(plan, " · "))
	}
	return b.String()
}

// TestHotQueriesUseIndexes: regressão do plano de consulta num banco pequeno.
// Confere pelo NOME do índice (com poucos dados o planejador pode mudar de
// ideia entre versões do SQLite; o ANALYZE antes deixa o teste estável).
func TestHotQueriesUseIndexes(t *testing.T) {
	ctx := context.Background()
	db := smallBenchDB(t)

	ids := []string{"m1", "m2", "m3"}
	now := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	want := map[string]string{
		"stats: buckets 30d":    "COVERING INDEX idx_checks_monitor_time",
		"lista: history":        "COVERING INDEX idx_checks_monitor_time",
		"lista: uptime 24h":     "COVERING INDEX idx_checks_monitor_time",
		"checks: página":        "INDEX idx_checks_monitor_time",
		"checks: total":         "COVERING INDEX idx_checks_monitor_time",
		"retenção: delete":      "INDEX idx_checks_checked_at",
		"stats: down window":    "INDEX idx_incidents_monitor_ended",
		"lista: down_since":     "INDEX idx_incidents_monitor_ended",
		"state machine: aberto": "INDEX idx_incidents_monitor_ended",
		"state machine: +falha": "INDEX idx_incidents_monitor_ended",
	}
	for _, q := range hotQueries(ids[0], ids, now) {
		plan, err := explain(ctx, db, q.sql, q.args...)
		if err != nil {
			t.Fatalf("%s: %v", q.name, err)
		}
		joined := strings.Join(plan, " · ")
		for _, line := range plan {
			// varredura de tabela sem índice (SCAN de subconsulta é só leitura do resultado)
			if strings.HasPrefix(line, "SCAN ") && !strings.Contains(line, "USING") && !strings.HasPrefix(line, "SCAN (subquery") {
				t.Errorf("%s: varredura completa: %s", q.name, joined)
			}
			// ordenação à parte das linhas do monitor: o índice não está servindo para o ORDER BY
			if strings.Contains(line, "LAST TERM OF ORDER BY") {
				t.Errorf("%s: ordenação fora do índice: %s", q.name, joined)
			}
		}
		if w, ok := want[q.name]; ok && !strings.Contains(joined, w) {
			t.Errorf("%s: plano sem %q: %s", q.name, w, joined)
		}
	}
}

// smallBenchDB cria um banco pequeno e determinístico pelas migrations reais,
// com ANALYZE, para os testes de plano e de equivalência.
func smallBenchDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "small.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}
	for i, typ := range []string{"http", "tcp", "heartbeat"} {
		id := fmt.Sprintf("m%d", i+1)
		var url, token any
		if typ == "http" {
			url = "https://x.test"
		}
		if typ == "heartbeat" {
			token = "tok"
		}
		exec(`INSERT INTO monitors (id, name, type, url, ping_token, interval_seconds, created_at, updated_at)
			VALUES (?,?,?,?,?,60,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`, id, id, typ, url, token)
	}
	// 9 dias de checks a cada 5 min, monitores intercalados; vários checks no
	// mesmo segundo para o desempate por id contar
	exec(`WITH RECURSIVE s(t) AS (SELECT 0 UNION ALL SELECT t + 300 FROM s WHERE t < 9 * 86400)
		INSERT INTO checks (monitor_id, success, status_code, latency_ms, checked_at)
		SELECT m.id, (s.t / 300) % 17 <> 0, 200, 20 + (s.t / 300) % 90,
			strftime('%Y-%m-%dT%H:%M:%SZ', 1767225600 + s.t - (s.t / 300) % 2 * 300, 'unixepoch')
		FROM s CROSS JOIN monitors m`)
	exec(`WITH RECURSIVE d(n) AS (SELECT 0 UNION ALL SELECT n + 1 FROM d WHERE n < 40)
		INSERT INTO incidents (id, monitor_id, started_at, ended_at, duration_seconds, failure_count)
		SELECT 'i' || m.id || '-' || d.n, m.id,
			strftime('%Y-%m-%dT%H:%M:%SZ', 1767225600 + d.n * 20000, 'unixepoch'),
			CASE WHEN d.n = 40 AND m.id = 'm1' THEN NULL
				ELSE strftime('%Y-%m-%dT%H:%M:%SZ', 1767225600 + d.n * 20000 + 600, 'unixepoch') END,
			600, 3
		FROM d CROSS JOIN monitors m`)
	exec(`ANALYZE`)
	return db
}

// SmallBenchDB expõe o banco pequeno ao pacote repository_test.
var SmallBenchDB = smallBenchDB

// legacySummariesHistory é a consulta do histórico antes da 0006 (ROW_NUMBER
// sobre todos os checks), mantida só como referência de equivalência.
func legacySummariesHistory(monitors int) string {
	return `
		SELECT monitor_id, success, latency_ms FROM (
			SELECT monitor_id, success, latency_ms,
				ROW_NUMBER() OVER (PARTITION BY monitor_id ORDER BY checked_at DESC, id DESC) AS rn
			FROM checks WHERE monitor_id IN (` + placeholders(monitors) + `)
		) WHERE rn <= ?
		ORDER BY monitor_id, rn DESC`
}

func TestSummariesHistoryMatchesLegacy(t *testing.T) {
	db := smallBenchDB(t)
	ctx := context.Background()
	ids := []string{"m1", "m2", "m3", "sem-checks"}
	for _, n := range []int{1, 30, 5000} { // 5000: mais que o total de checks do monitor
		args := []any{}
		for _, id := range ids {
			args = append(args, id)
		}
		want := dumpRows(t, db, legacySummariesHistory(len(ids)), append(args, n)...)

		r := NewCheckRepo(db)
		got, err := r.Summaries(ctx, ids, n, time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		var b strings.Builder
		for _, id := range ids {
			s := got[id]
			for i, ok := range s.History {
				fmt.Fprintf(&b, "%s %v", id, ok)
				if i == len(s.History)-1 {
					fmt.Fprintf(&b, " last=%d", *s.LastLatencyMs)
				}
				b.WriteString("\n")
			}
		}
		// a referência devolve a latência de todas as linhas; compara status e a última
		var ref strings.Builder
		lines := strings.Split(strings.TrimSpace(want), "\n")
		for i, l := range lines {
			f := strings.Fields(l)
			fmt.Fprintf(&ref, "%s %v", f[0], f[1] == "1")
			if i == len(lines)-1 || strings.Fields(lines[i+1])[0] != f[0] {
				fmt.Fprintf(&ref, " last=%s", f[2])
			}
			ref.WriteString("\n")
		}
		if b.String() != ref.String() {
			t.Fatalf("n=%d: histórico diferente da consulta antiga\n got:\n%s\nwant:\n%s", n, b.String(), ref.String())
		}
	}
}

func dumpRows(t *testing.T, db *sql.DB, q string, args ...any) string {
	t.Helper()
	rows, err := db.Query(q, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	var b strings.Builder
	for rows.Next() {
		v := make([]any, len(cols))
		p := make([]any, len(cols))
		for i := range v {
			p[i] = &v[i]
		}
		if err := rows.Scan(p...); err != nil {
			t.Fatal(err)
		}
		for i, x := range v {
			if i > 0 {
				b.WriteString(" ")
			}
			fmt.Fprint(&b, x)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// TestDeleteOlderThanBatches: em lotes, apaga exatamente o mesmo conjunto
// que o DELETE único de antes, e para com o ctx cancelado.
func TestDeleteOlderThanBatches(t *testing.T) {
	db := smallBenchDB(t)
	ctx := context.Background()
	before := "2026-01-07T00:00:00Z" // mais de um lote (retentionBatch)

	want := dumpRows(t, db, `SELECT id FROM checks WHERE checked_at >= ? ORDER BY id`, before)
	var old int64
	db.QueryRow(`SELECT COUNT(*) FROM checks WHERE checked_at < ?`, before).Scan(&old)
	if old <= retentionBatch {
		t.Fatalf("banco de teste pequeno demais: %d linhas para apagar", old)
	}

	cctx, cancel := context.WithCancel(ctx)
	cancel()
	r := NewCheckRepo(db)
	if n, err := r.DeleteOlderThan(cctx, before); err == nil || n != 0 {
		t.Fatalf("ctx cancelado: n=%d err=%v, want 0 e erro", n, err)
	}

	n, err := r.DeleteOlderThan(ctx, before)
	if err != nil {
		t.Fatal(err)
	}
	if n != old {
		t.Fatalf("apagou %d, want %d", n, old)
	}
	if got := dumpRows(t, db, `SELECT id FROM checks ORDER BY id`); got != want {
		t.Fatal("sobrou um conjunto diferente do esperado")
	}
}
