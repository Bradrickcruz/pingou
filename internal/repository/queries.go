package repository

import (
	"strings"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// SQL das consultas quentes (chamadas a cada check ou a cada refresh do painel).
// Ficam aqui, e não inline, para o teste de plano de consulta
// (queryplan_test.go) rodar EXPLAIN QUERY PLAN sobre as mesmas strings do código.

const checkColumns = `id, monitor_id, success, status_code, latency_ms, error_message, checked_at`

const (
	// sqlChecksPage: histórico paginado do detalhe do monitor e, com limit
	// pequeno, as falhas seguidas da state machine (a cada check).
	sqlChecksPage = `
		SELECT ` + checkColumns + `
		FROM checks WHERE monitor_id = ?
		ORDER BY checked_at DESC, id DESC LIMIT ? OFFSET ?`

	// sqlChecksLatestSince: falhas seguidas da state machine contando só checks
	// a partir da última retomada (pausa ou fim de manutenção). ">=" porque a
	// precisão é de segundos e o primeiro check roda no mesmo segundo da retomada.
	sqlChecksLatestSince = `
		SELECT ` + checkColumns + `
		FROM checks WHERE monitor_id = ? AND checked_at >= ?
		ORDER BY checked_at DESC, id DESC LIMIT ?`

	sqlChecksCount = `SELECT COUNT(*) FROM checks WHERE monitor_id = ?`

	// sqlChecksDeleteOlderBatch: retenção (a cada hora), um lote por vez.
	sqlChecksDeleteOlderBatch = `
		DELETE FROM checks WHERE id IN (
			SELECT id FROM checks WHERE checked_at < ? ORDER BY checked_at LIMIT ?)`

	// sqlChecksBuckets: série e totais do /stats.
	sqlChecksBuckets = `
		SELECT CAST(strftime('%s', checked_at) AS INTEGER) / ? AS b,
			COUNT(*), SUM(1 - success),
			COALESCE(SUM(CASE WHEN success = 1 THEN latency_ms END), 0), SUM(success)
		FROM checks
		WHERE monitor_id = ? AND checked_at >= ? AND checked_at < ?
		GROUP BY b ORDER BY b`

	// sqlIncidentsInWindow: janelas "down" do /stats.
	sqlIncidentsInWindow = `
		SELECT ` + incidentColumns + `
		FROM incidents WHERE monitor_id = ? AND (ended_at IS NULL OR ended_at > ?)
		ORDER BY started_at`

	// sqlIncidentOpenByMonitor: state machine (DOWN→UP) a cada recuperação.
	sqlIncidentOpenByMonitor = `
		SELECT ` + incidentColumns + `
		FROM incidents WHERE monitor_id = ? AND ended_at IS NULL`

	// sqlIncidentCloseManually: pausa e início de manutenção. Args: (ended_at,
	// motivo 1|2, close_note, ended_at, monitor_id).
	sqlIncidentCloseManually = `
		UPDATE incidents SET ended_at = ?, closed_manually = ?, close_note = ?,
			duration_seconds = MAX(0, CAST(strftime('%s', ?) - strftime('%s', started_at) AS INTEGER))
		WHERE monitor_id = ? AND ended_at IS NULL`

	// sqlIncidentIncrementFailures: a cada falha com o monitor DOWN.
	sqlIncidentIncrementFailures = `UPDATE incidents SET failure_count = failure_count + 1 WHERE monitor_id = ? AND ended_at IS NULL`
)

// placeholders devolve "?,?,…" com n posições.
func placeholders(n int) string {
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

// sqlSummariesHistory: últimos N checks de cada monitor da lista do painel,
// agrupados por monitor e do mais antigo ao mais recente. Uma subconsulta com
// LIMIT por monitor: cada uma lê só N entradas do índice, de trás para frente.
// (Um ROW_NUMBER() sobre todos os checks dos monitores numerava a tabela
// inteira para ficar com 30 linhas de cada.)
// Argumentos: (monitor_id, N) para cada monitor.
func sqlSummariesHistory(monitors int) string {
	part := `SELECT * FROM (
			SELECT monitor_id, success, latency_ms, checked_at, id FROM checks
			WHERE monitor_id = ? ORDER BY checked_at DESC, id DESC LIMIT ?)`
	parts := make([]string, monitors)
	for i := range parts {
		parts[i] = part
	}
	return `
		SELECT monitor_id, success, latency_ms FROM (
		` + strings.Join(parts, "\n\t\tUNION ALL ") + `
		) ORDER BY monitor_id, checked_at, id`
}

// sqlSummariesUptime: uptime das últimas 24 h de cada monitor da lista.
func sqlSummariesUptime(monitors int) string {
	return `
		SELECT monitor_id, 100.0 * SUM(success) / COUNT(*)
		FROM checks WHERE monitor_id IN (` + placeholders(monitors) + `) AND checked_at >= ?
		GROUP BY monitor_id`
}

// sqlIncidentsOpenByMonitors: incidente aberto dos monitores DOWN da lista.
func sqlIncidentsOpenByMonitors(monitors int) string {
	return `
		SELECT ` + incidentColumns + `
		FROM incidents WHERE ended_at IS NULL AND monitor_id IN (` + placeholders(monitors) + `)`
}

// incidentsListSQL monta o COUNT(*) e a página de GET /api/incidents com o
// mesmo WHERE. A lista recebe limit e offset no fim de args.
func incidentsListSQL(f domain.IncidentFilter) (count, list string, args []any) {
	where := ` WHERE 1=1`
	args = []any{}
	if f.MonitorID != "" {
		where += ` AND i.monitor_id = ?`
		args = append(args, f.MonitorID)
	}
	order := ` ORDER BY i.started_at DESC`
	switch f.Status {
	case domain.IncidentStatusOpen:
		where += ` AND i.ended_at IS NULL`
	case domain.IncidentStatusResolved:
		where += ` AND i.ended_at IS NOT NULL`
		order = ` ORDER BY i.ended_at DESC` // a tela agrupa resolvidos pelo dia do fim
	}
	count = `SELECT COUNT(*) FROM incidents i` + where
	list = `
		SELECT i.id, i.monitor_id, i.started_at, i.ended_at, i.last_error, i.notification_sent,
			i.duration_seconds, i.failure_count, i.closed_manually, i.close_note, m.name, m.type, m.url, m.config
		FROM incidents i LEFT JOIN monitors m ON m.id = i.monitor_id` + where + order + ` LIMIT ? OFFSET ?`
	return count, list, args
}
