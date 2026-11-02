// Gerador de dados sintéticos para medir o banco com volume real.
//
//	go run ./scripts/seed --profile tipico --out .bench/tipico.db
//	go run ./scripts/seed --profile limite --out .bench/limite.db
//
// O banco nasce pelas migrations reais e os checks são inseridos com
// INSERT ... SELECT sobre uma CTE recursiva, em ordem cronológica e com os
// monitores intercalados (o mesmo layout em disco que o app produz).
//
// Determinístico: os valores vêm de contas inteiras sobre o número do passo e
// do monitor, sem random(). Duas execuções com o mesmo --end geram o mesmo
// banco. O --end padrão é o início do dia atual (UTC), para os dados
// terminarem perto de "agora" e as consultas do painel acharem dados.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
)

type profile struct {
	monitors int // total; os últimos viram TCP e Heartbeat
	tcp      int
	hb       int
	interval int // segundos (HTTP e TCP)
	hbEvery  int // segundos entre pings dos heartbeats
}

var profiles = map[string]profile{
	"tipico": {monitors: 10, tcp: 1, hb: 1, interval: 60, hbEvery: 300},
	"limite": {monitors: 100, tcp: 5, hb: 5, interval: 10, hbEvery: 300},
}

const (
	days       = 30
	week       = 7 * 24 * 3600
	histWeeks  = 52                              // incidentes antigos (mais velhos que a retenção) até 1 ano atrás
	monitorFmt = "00000000-0000-7000-8000-%012d" // mesmo tamanho de um UUID
)

func main() {
	name := flag.String("profile", "tipico", "tipico ou limite")
	out := flag.String("out", "", "arquivo .db de saída (é sobrescrito)")
	endFlag := flag.String("end", "", "fim dos dados, RFC3339 (padrão: hoje 00:00 UTC)")
	flag.Parse()

	p, ok := profiles[*name]
	if !ok || *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	end := time.Now().UTC().Truncate(24 * time.Hour)
	if *endFlag != "" {
		t, err := time.Parse(time.RFC3339, *endFlag)
		if err != nil {
			log.Fatalf("--end: %v", err)
		}
		end = t.UTC()
	}
	start := end.Add(-days * 24 * time.Hour)

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		log.Fatal(err)
	}
	for _, suf := range []string{"", "-wal", "-shm"} {
		_ = os.Remove(*out + suf)
	}

	began := time.Now()
	db, err := database.Open(*out) // migrations reais
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// acelera a carga; o padrão do app (WAL) é restaurado no fim
	mustExec(db, `PRAGMA journal_mode = OFF`)
	mustExec(db, `PRAGMA synchronous = OFF`)
	mustExec(db, `PRAGMA cache_size = -262144`) // 256 MB só durante a carga
	mustExec(db, `PRAGMA temp_store = MEMORY`)

	seedMonitors(db, p, start, end)
	seedChecks(db, p, start, end)
	seedIncidents(db, p, start, end)
	mustExec(db, `UPDATE settings SET value = '30' WHERE key = 'retention_days'`)

	mustExec(db, `PRAGMA journal_mode = WAL`)
	mustExec(db, `PRAGMA synchronous = FULL`)
	mustExec(db, `PRAGMA wal_checkpoint(TRUNCATE)`)

	printSummary(db, *out, *name, start, end, time.Since(began))
}

// seedMonitors: monitores 1..N. Os últimos `hb` são Heartbeat e os `tcp`
// antes deles são TCP; o resto é HTTP.
func seedMonitors(db *sql.DB, p profile, start, end time.Time) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO monitors
		(id, name, type, url, config, ping_token, interval_seconds, timeout_seconds,
		 failure_threshold, enabled, current_state, last_checked_at, last_ping_at, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,10,3,1,'UP',?,?,?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	created := start.Add(-histWeeks * week * time.Second).Format(time.RFC3339)
	for i := 1; i <= p.monitors; i++ {
		id := fmt.Sprintf(monitorFmt, i)
		typ, url, cfg, interval := monitorKind(p, i)
		var token, lastPing any
		last := lastStep(start, end, interval)
		if typ == "heartbeat" {
			token = fmt.Sprintf("seed-token-%04d", i)
			lastPing = last
		}
		if _, err := stmt.Exec(id, fmt.Sprintf("seed %s %03d", typ, i), typ, url, cfg, token,
			interval, last, lastPing, created, last); err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func monitorKind(p profile, i int) (typ string, url any, cfg string, interval int) {
	switch {
	case i > p.monitors-p.hb:
		return "heartbeat", nil, `{"grace_seconds":60}`, p.hbEvery
	case i > p.monitors-p.hb-p.tcp:
		return "tcp", nil, fmt.Sprintf(`{"host":"db%03d.example.com","port":5432}`, i), p.interval
	default:
		return "http", fmt.Sprintf("https://site%03d.example.com/health", i), `{}`, p.interval
	}
}

func lastStep(start, end time.Time, interval int) string {
	n := int(end.Sub(start)/time.Second) / interval
	return start.Add(time.Duration((n-1)*interval) * time.Second).Format(time.RFC3339)
}

// outageSQL: começo (offset dentro da semana) e duração (5 a 30 min) da queda
// semanal do monitor m na semana w. As mesmas contas geram checks e incidentes.
const (
	outageOff = `((%[1]s * 7919 + (%[2]s + 1000) * 104729) %% (604800 - 1800))`
	outageDur = `(300 + ((%[1]s * 31 + (%[2]s + 1000) * 17) %% 1501))`
)

// seedChecks gera um check por passo e por monitor. O CROSS JOIN força o
// passo no laço de fora: as linhas saem em ordem cronológica, com os
// monitores intercalados, como o scheduler grava.
func seedChecks(db *sql.DB, p profile, start, end time.Time) {
	steps := int(end.Sub(start)/time.Second) / p.interval
	off := fmt.Sprintf(outageOff, "m.n", "(s.t / 604800)")
	dur := fmt.Sprintf(outageDur, "m.n", "(s.t / 604800)")
	// down: dentro da queda semanal; isolada: ~0,3 % (resíduo do passo, nunca
	// duas seguidas, então não abre incidente com threshold 3)
	q := fmt.Sprintf(`
	WITH RECURSIVE
	  s(t) AS (SELECT 0 UNION ALL SELECT t + %[1]d FROM s WHERE t + %[1]d < %[2]d),
	  m(n, id, typ, every) AS (
	    SELECT CAST(substr(id, 25) AS INTEGER), id, type, interval_seconds FROM monitors
	  )
	INSERT INTO checks (monitor_id, success, status_code, latency_ms, error_message, checked_at)
	SELECT id,
	  CASE WHEN down OR iso THEN 0 ELSE 1 END,
	  CASE WHEN typ <> 'http' THEN NULL WHEN down THEN 503 WHEN iso THEN NULL ELSE 200 END,
	  CASE WHEN typ = 'heartbeat' THEN 0 WHEN iso THEN 10000 ELSE lat END,
	  CASE WHEN typ = 'heartbeat' AND down THEN 'no ping received within interval + grace'
	       WHEN down AND typ = 'tcp' THEN 'dial tcp: connection refused'
	       WHEN down THEN 'unexpected status code: 503'
	       WHEN iso THEN 'context deadline exceeded (Client.Timeout exceeded while awaiting headers)'
	  END,
	  strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', %[3]d + t, 'unixepoch')
	FROM (
	  SELECT m.id, m.typ, s.t,
	    ((s.t %% 604800) >= %[4]s AND (s.t %% 604800) < %[4]s + %[5]s) AS down,
	    (m.typ <> 'heartbeat' AND (((s.t / %[1]d) * 761 + m.n * 337) %% 1000) < 3) AS iso,
	    20 + (m.n * 37) %% 80 + (((s.t / %[1]d) * 48271 + m.n * 69621) %% 2147483647) %% 300 AS lat
	  FROM s CROSS JOIN m
	  WHERE s.t %% m.every = 0
	)`, p.interval, steps*p.interval, start.Unix(), off, dur)
	t0 := time.Now()
	mustExec(db, q)
	log.Printf("checks inseridos em %s", time.Since(t0).Round(time.Millisecond))
}

// seedIncidents: uma queda por monitor por semana, de 1 ano atrás até o fim
// dos dados. Os de antes do início não têm checks (a retenção já apagou),
// como num banco real: incidentes nunca expiram.
func seedIncidents(db *sql.DB, p profile, start, end time.Time) {
	off := fmt.Sprintf(outageOff, "m.n", "w.w")
	dur := fmt.Sprintf(outageDur, "m.n", "w.w")
	q := fmt.Sprintf(`
	WITH RECURSIVE
	  w(w) AS (SELECT -%[1]d UNION ALL SELECT w + 1 FROM w WHERE w < 5),
	  m(n, id, typ, every) AS (
	    SELECT CAST(substr(id, 25) AS INTEGER), id, type, interval_seconds FROM monitors
	  ),
	  o AS (
	    SELECT m.n, m.id, m.typ, m.every, w.w,
	      %[2]d + w.w * 604800 + %[3]s AS s0,
	      %[2]d + w.w * 604800 + %[3]s + %[4]s AS s1
	    FROM w CROSS JOIN m
	  )
	INSERT INTO incidents (id, monitor_id, started_at, ended_at, last_error,
	  notification_sent, duration_seconds, failure_count, closed_manually)
	SELECT printf('00000000-0000-7000-9000-%%06d%%06d', n, w + 1000), id,
	  strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', s0, 'unixepoch'),
	  strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', s1, 'unixepoch'),
	  CASE typ WHEN 'heartbeat' THEN 'no ping received within interval + grace'
	    WHEN 'tcp' THEN 'dial tcp: connection refused' ELSE 'unexpected status code: 503' END, 1, s1 - s0, MAX(1, (s1 - s0) / every), 0
	FROM o WHERE s1 <= %[5]d
	ORDER BY s0`, histWeeks, start.Unix(), off, dur, end.Unix())
	mustExec(db, q)
}

func printSummary(db *sql.DB, out, name string, start, end time.Time, took time.Duration) {
	var checks, incidents, monitors int
	db.QueryRow(`SELECT COUNT(*) FROM checks`).Scan(&checks)
	db.QueryRow(`SELECT COUNT(*) FROM incidents`).Scan(&incidents)
	db.QueryRow(`SELECT COUNT(*) FROM monitors`).Scan(&monitors)
	var size int64
	if fi, err := os.Stat(out); err == nil {
		size = fi.Size()
	}
	fmt.Printf("perfil:     %s\n", name)
	fmt.Printf("arquivo:    %s (%.1f MB)\n", out, float64(size)/1e6)
	fmt.Printf("período:    %s → %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339))
	fmt.Printf("monitores:  %d\n", monitors)
	fmt.Printf("checks:     %d\n", checks)
	fmt.Printf("incidentes: %d\n", incidents)
	fmt.Printf("tempo:      %s\n", took.Round(time.Millisecond))
}

func mustExec(db *sql.DB, q string) {
	if _, err := db.Exec(q); err != nil {
		log.Fatalf("%v\n%s", err, q)
	}
}
