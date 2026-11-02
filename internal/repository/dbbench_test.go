package repository_test

// Medição do banco com volume real. Só roda com PINGOU_BENCH_DB apontando
// para um banco gerado por scripts/seed (make test continua rápido):
//
//	PINGOU_BENCH_DB=.bench/limite.db go test ./internal/repository -run TestDBBench -v -timeout 2h
//
// Trabalha numa CÓPIA do banco (o original nunca é alterado) e chama as
// funções reais de repositório e serviço, com a mesma conexão única do app.
// Variáveis opcionais: PINGOU_BENCH_REPS (repetições, padrão 20),
// PINGOU_BENCH_WRITES (gravações no custo de escrita, padrão 5000) e
// PINGOU_BENCH_CONCURRENCY (segundos do cenário de concorrência, padrão 30).

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type noopNotifier struct{}

func (noopNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident)     {}
func (noopNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) {}

type benchEnv struct {
	db        *sql.DB
	migration string // versões aplicadas ao abrir e quanto levou
	monitors  *repository.MonitorRepo
	checks    *repository.CheckRepo
	incidents *repository.IncidentRepo
	svc       *service.MonitorService
	now       time.Time // fim dos dados (o "agora" das consultas)
	first     time.Time // check mais antigo
	byType    map[domain.MonitorType]string
}

func envInt(name string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(name)); err == nil && v > 0 {
		return v
	}
	return def
}

func copyFile(t testing.TB, src, dst string) {
	t.Helper()
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}

func openBench(t *testing.T) *benchEnv {
	src := repository.BenchDBPath(t)
	dst := filepath.Join(t.TempDir(), "bench.db")
	t0 := time.Now()
	copyFile(t, src, dst)
	t.Logf("cópia de %s em %s", src, time.Since(t0).Round(time.Millisecond))

	// aplica as migrations pendentes, como o app ao subir: num banco gerado
	// antes de uma migration nova, este é o tempo da atualização
	db, err := database.OpenDB(dst)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	var v0, v1 int64
	db.QueryRow(`SELECT MAX(version_id) FROM goose_db_version`).Scan(&v0)
	t0 = time.Now()
	if err := database.RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	took := time.Since(t0)
	db.QueryRow(`SELECT MAX(version_id) FROM goose_db_version`).Scan(&v1)

	e := &benchEnv{
		db:        db,
		monitors:  repository.NewMonitorRepo(db),
		checks:    repository.NewCheckRepo(db),
		incidents: repository.NewIncidentRepo(db),
		byType:    map[domain.MonitorType]string{},
		migration: fmt.Sprintf("schema %d → %d em %s", v0, v1, took.Round(time.Millisecond)),
	}
	e.svc = service.NewMonitorService(e.monitors, e.checks, e.incidents)

	var last, first string
	if err := db.QueryRow(`SELECT MAX(checked_at), MIN(checked_at) FROM checks`).Scan(&last, &first); err != nil {
		t.Fatal(err)
	}
	e.now, _ = time.Parse(time.RFC3339, last)
	e.now = e.now.Add(time.Minute)
	e.first, _ = time.Parse(time.RFC3339, first)

	rows, err := db.Query(`SELECT type, MIN(id) FROM monitors GROUP BY type`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var typ, id string
		rows.Scan(&typ, &id)
		e.byType[domain.MonitorType(typ)] = id
	}
	rows.Close()
	return e
}

// timing guarda a primeira execução (cache frio para essa consulta) e as
// repetições (cache quente).
type timing struct {
	name  string
	first time.Duration
	reps  []time.Duration
}

func (tm timing) pct(p float64) time.Duration {
	if len(tm.reps) == 0 {
		return 0
	}
	s := append([]time.Duration(nil), tm.reps...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	i := int(p*float64(len(s))+0.5) - 1
	return s[max(0, min(i, len(s)-1))]
}

func ms(d time.Duration) string { return fmt.Sprintf("%.1f", float64(d.Microseconds())/1000) }

func measure(t *testing.T, name string, reps int, fn func() error) timing {
	t.Helper()
	tm := timing{name: name}
	for i := 0; i <= reps; i++ {
		t0 := time.Now()
		if err := fn(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		d := time.Since(t0)
		if i == 0 {
			tm.first = d
		} else {
			tm.reps = append(tm.reps, d)
		}
	}
	return tm
}

func printTable(title string, ts []timing) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n### %s\n\n| Consulta | 1ª (ms) | p50 (ms) | p95 (ms) | máx (ms) |\n|---|---:|---:|---:|---:|\n", title)
	for _, tm := range ts {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", tm.name, ms(tm.first), ms(tm.pct(.5)), ms(tm.pct(.95)), ms(tm.pct(1)))
	}
	return b.String()
}

// sizes mede o arquivo e cada índice. Sem dbstat no go-sqlite3: apaga o
// índice numa transação, vê quantas páginas foram para a lista livre e
// desfaz (rápido, nada é gravado).
func sizes(t *testing.T, db *sql.DB) string {
	var pageSize, pages int64
	db.QueryRow(`PRAGMA page_size`).Scan(&pageSize)
	db.QueryRow(`PRAGMA page_count`).Scan(&pages)
	mb := func(p int64) string { return fmt.Sprintf("%.1f", float64(p*pageSize)/1e6) }

	var b strings.Builder
	fmt.Fprintf(&b, "\n### Tamanho\n\narquivo: %s MB\n\n| Índice | Tabela | MB |\n|---|---|---:|\n", mb(pages))
	rows, err := db.Query(`SELECT name, tbl_name FROM sqlite_master WHERE type = 'index' AND sql IS NOT NULL ORDER BY tbl_name, name`)
	if err != nil {
		t.Fatal(err)
	}
	var idx [][2]string
	for rows.Next() {
		var n, tbl string
		rows.Scan(&n, &tbl)
		idx = append(idx, [2]string{n, tbl})
	}
	rows.Close()
	var sum int64
	for _, ix := range idx {
		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		var before, after int64
		tx.QueryRow(`PRAGMA freelist_count`).Scan(&before)
		if _, err := tx.Exec(`DROP INDEX "` + ix[0] + `"`); err != nil {
			t.Fatal(err)
		}
		tx.QueryRow(`PRAGMA freelist_count`).Scan(&after)
		tx.Rollback()
		sum += after - before
		fmt.Fprintf(&b, "| %s | %s | %s |\n", ix[0], ix[1], mb(after-before))
	}
	fmt.Fprintf(&b, "| (todos os índices) | | %s |\n", mb(sum))
	return b.String()
}

func TestDBBench(t *testing.T) {
	e := openBench(t)
	ctx := context.Background()
	reps := envInt("PINGOU_BENCH_REPS", 20)
	var report strings.Builder
	defer func() { fmt.Print(report.String()) }() // imprime o que deu, mesmo se falhar no meio

	var counts struct{ checks, incidents, monitors int }
	e.db.QueryRow(`SELECT (SELECT COUNT(*) FROM checks), (SELECT COUNT(*) FROM incidents), (SELECT COUNT(*) FROM monitors)`).
		Scan(&counts.checks, &counts.incidents, &counts.monitors)
	var sqliteVersion string
	e.db.QueryRow(`SELECT sqlite_version()`).Scan(&sqliteVersion)
	fmt.Fprintf(&report, "\nbanco: %d monitores, %d checks, %d incidentes · SQLite %s · %d repetições\n",
		counts.monitors, counts.checks, counts.incidents, sqliteVersion, reps)
	fmt.Fprintf(&report, "migrations ao abrir: %s\n", e.migration)

	report.WriteString(sizes(t, e.db))

	var all []string
	rows, err := e.db.Query(`SELECT id FROM monitors ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		all = append(all, id)
	}
	rows.Close()
	report.WriteString("\n### Planos de consulta (na cópia, depois das migrations)\n")
	report.WriteString(repository.QueryPlans(t, e.db, all, e.now))

	// 1. /stats nos 4 períodos, em um monitor de cada tipo
	var stats []timing
	for _, rng := range []string{"24h", "7d", "14d", "30d"} {
		for _, typ := range []domain.MonitorType{domain.TypeHTTP, domain.TypeTCP, domain.TypeHeartbeat} {
			id := e.byType[typ]
			if id == "" {
				continue
			}
			stats = append(stats, measure(t, fmt.Sprintf("/stats %s %s", rng, typ), reps, func() error {
				_, err := e.svc.Stats(ctx, id, rng, e.now)
				return err
			}))
		}
	}
	report.WriteString(printTable("/stats", stats))

	// 2. demais consultas quentes
	httpID := e.byType[domain.TypeHTTP]
	enabled, hb := true, domain.TypeHeartbeat
	var token string
	e.db.QueryRow(`SELECT ping_token FROM monitors WHERE ping_token IS NOT NULL LIMIT 1`).Scan(&token)
	others := []timing{
		measure(t, "GET /api/monitors (lista + resumo)", reps, func() error {
			mons, _, err := e.monitors.FindAll(ctx, domain.MonitorFilter{Limit: 100})
			if err != nil {
				return err
			}
			_, err = e.svc.Summaries(ctx, mons)
			return err
		}),
	}
	for _, off := range []int{0, 1000, 10000} {
		others = append(others, measure(t, fmt.Sprintf("checks limit 50 offset %d", off), reps, func() error {
			_, _, err := e.checks.FindByMonitor(ctx, httpID, 50, off)
			return err
		}))
	}
	for _, f := range []struct {
		name string
		f    domain.IncidentFilter
	}{
		{"incidentes: todos", domain.IncidentFilter{}},
		{"incidentes: abertos", domain.IncidentFilter{Status: domain.IncidentStatusOpen}},
		{"incidentes: resolvidos", domain.IncidentFilter{Status: domain.IncidentStatusResolved}},
		{"incidentes: do monitor", domain.IncidentFilter{MonitorID: httpID}},
	} {
		others = append(others, measure(t, f.name, reps, func() error {
			_, _, err := e.incidents.FindAll(ctx, f.f)
			return err
		}))
	}
	others = append(others,
		measure(t, "heartbeat: monitores ativos", reps, func() error {
			_, _, err := e.monitors.FindAll(ctx, domain.MonitorFilter{Enabled: &enabled, Type: &hb, Limit: 100})
			return err
		}),
		measure(t, "/ping/{token}: busca do token", reps, func() error {
			_, err := e.monitors.FindByPingToken(ctx, token)
			return err
		}),
	)
	report.WriteString(printTable("Outras consultas", others))

	// 3. custo de gravação pelo caminho real (state machine, um check por transação)
	writes := envInt("PINGOU_BENCH_WRITES", 5000)
	ids := []string{}
	rows, _ = e.db.Query(`SELECT id FROM monitors WHERE type <> 'heartbeat' ORDER BY id`)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rows.Close()
	writeResult := func(i int, ok bool) error {
		m, err := e.monitors.FindByID(ctx, ids[i%len(ids)])
		if err != nil {
			return err
		}
		res := domain.CheckResult{Success: ok, LatencyMs: 100}
		if !ok {
			msg := "bench failure"
			res.ErrorMessage = &msg
		}
		uow := service.NewUnitOfWork(e.db, repository.NewCheckRepoTx(e.db), repository.NewMonitorRepoTx(e.db), repository.NewIncidentRepoTx(e.db))
		return service.NewStateMachine(uow, noopNotifier{}).Process(ctx, m, res)
	}
	write := func(i int) error { return writeResult(i, true) }
	w := timing{name: fmt.Sprintf("gravação de check (%d)", writes)}
	t0 := time.Now()
	for i := 0; i < writes; i++ {
		s := time.Now()
		if err := write(i); err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			w.first = time.Since(s)
		} else {
			w.reps = append(w.reps, time.Since(s))
		}
	}
	total := time.Since(t0)
	// falha: a state machine lê os últimos checks para contar falhas seguidas;
	// abre incidente e marca DOWN nos monitores (é uma cópia do banco)
	failures := min(writes, 2000)
	f := timing{name: fmt.Sprintf("gravação de falha (%d)", failures)}
	for i := 0; i < failures; i++ {
		s := time.Now()
		if err := writeResult(i, false); err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			f.first = time.Since(s)
		} else {
			f.reps = append(f.reps, time.Since(s))
		}
	}
	for i := range ids { // volta todos a UP antes dos próximos cenários
		if err := write(i); err != nil {
			t.Fatal(err)
		}
	}
	report.WriteString(printTable("Gravação", []timing{w, f}))
	fmt.Fprintf(&report, "\nmédia por check com sucesso: %s ms (%d gravações em %s)\n", ms(total/time.Duration(writes)), writes, total.Round(time.Millisecond))

	// 4. concorrência: /stats 30d sem parar + gravação a 10 por segundo
	secs := envInt("PINGOU_BENCH_CONCURRENCY", 30)
	cctx, cancel := context.WithTimeout(ctx, time.Duration(secs)*time.Second)
	var wg sync.WaitGroup
	readers := timing{name: "/stats 30d (leitura contínua)"}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cctx.Err() == nil {
			s := time.Now()
			if _, err := e.svc.Stats(ctx, httpID, "30d", e.now); err != nil {
				t.Error(err)
				return
			}
			readers.reps = append(readers.reps, time.Since(s))
		}
	}()
	writer := timing{name: "gravação com leitura ao lado"}
	tick := time.NewTicker(100 * time.Millisecond)
	for i := 0; cctx.Err() == nil; i++ {
		<-tick.C
		s := time.Now()
		if err := write(i); err != nil {
			t.Fatal(err)
		}
		writer.reps = append(writer.reps, time.Since(s))
	}
	tick.Stop()
	cancel()
	wg.Wait()
	if len(readers.reps) > 0 {
		readers.first = readers.reps[0]
	}
	if len(writer.reps) > 0 {
		writer.first = writer.reps[0]
	}
	report.WriteString(printTable(fmt.Sprintf("Concorrência (%d s)", secs), []timing{readers, writer}))

	// 5. retenção: apaga o dia mais antigo (o mesmo que a limpeza horária faz
	// depois de um dia parado) com gravações a 10 por segundo ao lado. A espera
	// máxima de uma gravação mostra por quanto tempo a limpeza segura a conexão.
	before := e.first.Add(24 * time.Hour).UTC().Format(time.RFC3339)
	type result struct {
		n   int64
		d   time.Duration
		err error
	}
	done := make(chan result, 1)
	s := time.Now()
	go func() {
		n, err := e.checks.DeleteOlderThan(ctx, before)
		done <- result{n, time.Since(s), err}
	}()
	during := timing{name: "gravação durante a limpeza"}
	var res result
	tick = time.NewTicker(100 * time.Millisecond)
loop:
	for i := 0; ; i++ {
		select {
		case res = <-done:
			break loop
		case <-tick.C:
		}
		s := time.Now()
		if err := write(i); err != nil {
			t.Fatal(err)
		}
		during.reps = append(during.reps, time.Since(s))
	}
	tick.Stop()
	if res.err != nil {
		t.Fatal(res.err)
	}
	if len(during.reps) > 0 {
		during.first = during.reps[0]
	}
	fmt.Fprintf(&report, "\n### Retenção\n\nDELETE de 1 dia: %d linhas em %s ms\n", res.n, ms(res.d))
	report.WriteString(printTable("Gravação durante a retenção", []timing{during}))
}

// TestResultsSameWithOldIndex: as respostas das consultas quentes são as mesmas
// com o índice novo (0006) e depois de desfazê-lo (goose down), no mesmo banco.
func TestResultsSameWithOldIndex(t *testing.T) {
	db := repository.SmallBenchDB(t)
	ctx := context.Background()
	now := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	ids := []string{"m1", "m2", "m3"}

	snapshot := func() string {
		monitors := repository.NewMonitorRepo(db)
		checks := repository.NewCheckRepo(db)
		incidents := repository.NewIncidentRepo(db)
		svc := service.NewMonitorService(monitors, checks, incidents)
		var out []any
		add := func(v any, err error) {
			t.Helper()
			if err != nil {
				t.Fatal(err)
			}
			out = append(out, v)
		}
		for _, id := range ids {
			for _, rng := range []string{"24h", "7d", "14d", "30d"} {
				add(svc.Stats(ctx, id, rng, now))
			}
			for _, off := range []int{0, 7, 500} {
				c, total, err := checks.FindByMonitor(ctx, id, 50, off)
				add([]any{c, total}, err)
			}
			add(checks.Buckets(ctx, id, now.Add(-9*24*time.Hour), now, 600))
		}
		add(checks.Summaries(ctx, ids, 30, now.Add(-24*time.Hour)))
		ms, _, err := monitors.FindAll(ctx, domain.MonitorFilter{Limit: 100})
		add(ms, err)
		add(svc.Summaries(ctx, ms))
		for _, f := range []domain.IncidentFilter{
			{}, {Status: domain.IncidentStatusOpen}, {Status: domain.IncidentStatusResolved},
			{MonitorID: "m1"}, {Limit: 5, Offset: 10},
		} {
			l, total, err := incidents.FindAll(ctx, f)
			add([]any{l, total}, err)
		}
		b, err := json.Marshal(out)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	after := snapshot()
	if err := database.Down(db); err != nil { // desfaz a 0006
		t.Fatal(err)
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE name = 'idx_checks_monitor_time'`).Scan(&n)
	if n != 0 {
		t.Fatal("down não removeu idx_checks_monitor_time")
	}
	if before := snapshot(); before != after {
		t.Fatalf("respostas diferentes com e sem o índice novo\nantes: %.500s\ndepois: %.500s", before, after)
	}
}
