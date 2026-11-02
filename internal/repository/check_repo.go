package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type CheckRepo struct{ db *sql.DB }

func NewCheckRepo(db *sql.DB) *CheckRepo { return &CheckRepo{db: db} }

func (r *CheckRepo) Create(ctx context.Context, c *domain.Check) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO checks (monitor_id, success, status_code, latency_ms, error_message, checked_at)
		VALUES (?,?,?,?,?,?)`,
		c.MonitorID, boolToInt(c.Success), c.StatusCode, c.LatencyMs,
		c.ErrorMessage, c.CheckedAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (r *CheckRepo) FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, sqlChecksPage,
		monitorID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var checks []*domain.Check
	for rows.Next() {
		var c domain.Check
		var success int
		var checkedAt string
		if err := rows.Scan(&c.ID, &c.MonitorID, &success, &c.StatusCode, &c.LatencyMs, &c.ErrorMessage, &checkedAt); err != nil {
			return nil, 0, err
		}
		c.Success = success == 1
		c.CheckedAt, _ = time.Parse(time.RFC3339, checkedAt)
		checks = append(checks, &c)
	}

	var total int
	r.db.QueryRowContext(ctx, sqlChecksCount, monitorID).Scan(&total)

	return checks, total, nil
}

func (r *CheckRepo) DeleteOlderThan(ctx context.Context, before string) (int64, error) {
	res, err := r.db.ExecContext(ctx, sqlChecksDeleteOlder, before)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Summaries monta o resumo do painel para vários monitores em duas queries
// (histórico e uptime), sem uma ida ao banco por monitor.
func (r *CheckRepo) Summaries(ctx context.Context, monitorIDs []string, historySize int, since time.Time) (map[string]*domain.MonitorSummary, error) {
	out := make(map[string]*domain.MonitorSummary, len(monitorIDs))
	if len(monitorIDs) == 0 {
		return out, nil
	}
	for _, id := range monitorIDs {
		out[id] = &domain.MonitorSummary{History: []bool{}}
	}

	args := make([]any, 0, len(monitorIDs)+1)
	for _, id := range monitorIDs {
		args = append(args, id)
	}

	// o SQLite aceita até 500 partes num UNION ALL; a lista do painel tem no
	// máximo 100 monitores, mas lotes deixam a função segura para qualquer N
	for start := 0; start < len(monitorIDs); start += historyBatch {
		batch := monitorIDs[start:min(start+historyBatch, len(monitorIDs))]
		histArgs := make([]any, 0, 2*len(batch))
		for _, id := range batch {
			histArgs = append(histArgs, id, historySize)
		}
		if err := r.scanHistory(ctx, out, sqlSummariesHistory(len(batch)), histArgs); err != nil {
			return nil, err
		}
	}

	rows, err := r.db.QueryContext(ctx, sqlSummariesUptime(len(monitorIDs)),
		append(args, since.UTC().Format(time.RFC3339))...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var uptime float64
		if err := rows.Scan(&id, &uptime); err != nil {
			return nil, err
		}
		out[id].Uptime24h = &uptime
	}
	return out, rows.Err()
}

const historyBatch = 100

func (r *CheckRepo) scanHistory(ctx context.Context, out map[string]*domain.MonitorSummary, q string, args []any) error {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var success int
		var latency int64
		if err := rows.Scan(&id, &success, &latency); err != nil {
			return err
		}
		s := out[id]
		s.History = append(s.History, success == 1)
		s.LastLatencyMs = &latency
	}
	return rows.Err()
}

// Buckets agrupa os checks de [since, until) em baldes de bucketSeconds.
// Devolve só os baldes com dados, em ordem; preencher os vazios é do service.
func (r *CheckRepo) Buckets(ctx context.Context, monitorID string, since, until time.Time, bucketSeconds int) ([]domain.StatsBucket, error) {
	rows, err := r.db.QueryContext(ctx, sqlChecksBuckets,
		bucketSeconds, monitorID,
		since.UTC().Format(time.RFC3339), until.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.StatsBucket
	for rows.Next() {
		var b int64
		var s domain.StatsBucket
		if err := rows.Scan(&b, &s.Checks, &s.Failures, &s.SumLatencyMs, &s.Successes); err != nil {
			return nil, err
		}
		s.Start = time.Unix(b*int64(bucketSeconds), 0).UTC()
		if s.Successes > 0 {
			avg := float64(s.SumLatencyMs) / float64(s.Successes)
			s.AvgLatencyMs = &avg
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
