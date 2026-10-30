package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type MonitorRepo struct{ db *sql.DB }

func NewMonitorRepo(db *sql.DB) *MonitorRepo {
	return &MonitorRepo{db: db}
}

func (r *MonitorRepo) Create(ctx context.Context, m *domain.Monitor) error {
	if m.Type == "" {
		m.Type = domain.TypeHTTP // mesmo padrão da coluna
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO monitors
			(id, name, type, url, config, ping_token, interval_seconds, timeout_seconds,
			 failure_threshold, enabled, current_state, last_ping_at, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		m.ID, m.Name, string(m.Type), nullableString(m.URL), encodeConfig(m), m.PingToken,
		m.IntervalSeconds, m.TimeoutSeconds, m.FailureThreshold,
		boolToInt(m.Enabled), string(m.CurrentState), nullableTime(m.LastPingAt),
		m.CreatedAt.UTC().Format(time.RFC3339),
		m.UpdatedAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (r *MonitorRepo) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+monitorColumns+` FROM monitors WHERE id = ?`, id)
	return scanMonitor(row)
}

// FindByPingToken busca o heartbeat dono do token (índice único em ping_token).
func (r *MonitorRepo) FindByPingToken(ctx context.Context, token string) (*domain.Monitor, error) {
	if token == "" {
		return nil, nil
	}
	row := r.db.QueryRowContext(ctx, `SELECT `+monitorColumns+` FROM monitors WHERE ping_token = ?`, token)
	return scanMonitor(row)
}

func (r *MonitorRepo) FindAll(ctx context.Context, f domain.MonitorFilter) ([]*domain.Monitor, int, error) {
	query := `SELECT ` + monitorColumns + ` FROM monitors WHERE 1=1`
	args := []any{}

	if f.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, boolToInt(*f.Enabled))
	}
	if f.State != nil {
		query += ` AND current_state = ?`
		args = append(args, string(*f.State))
	}
	if f.Type != nil {
		query += ` AND type = ?`
		args = append(args, string(*f.Type))
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, limit, f.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var monitors []*domain.Monitor
	for rows.Next() {
		m, err := scanMonitor(rows)
		if err != nil {
			return nil, 0, err
		}
		monitors = append(monitors, m)
	}

	var total int
	r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM monitors`).Scan(&total)

	return monitors, total, nil
}

// Update grava os campos editáveis pelo usuário. Estado e last_checked_at
// ficam com a state machine (UpdateWithTx), para uma edição não desfazer um check.
func (r *MonitorRepo) Update(ctx context.Context, m *domain.Monitor) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE monitors SET
			name=?, url=?, config=?, ping_token=?, interval_seconds=?, timeout_seconds=?,
			failure_threshold=?, enabled=?, last_ping_at=?, updated_at=?
		WHERE id=?`,
		m.Name, nullableString(m.URL), encodeConfig(m), m.PingToken, m.IntervalSeconds, m.TimeoutSeconds,
		m.FailureThreshold, boolToInt(m.Enabled), nullableTime(m.LastPingAt),
		m.UpdatedAt.UTC().Format(time.RFC3339),
		m.ID,
	)
	return err
}

func (r *MonitorRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM monitors WHERE id = ?`, id)
	return err
}

func (r *MonitorRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM monitors`).Scan(&n)
	return n, err
}

// ── helpers ────────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

const monitorColumns = `id, name, type, url, config, ping_token, interval_seconds, timeout_seconds,
	failure_threshold, enabled, current_state, last_checked_at, last_ping_at, created_at, updated_at`

func scanMonitor(s scanner) (*domain.Monitor, error) {
	var m domain.Monitor
	var enabled int
	var typ, state, config string
	var url, lastChecked, lastPing sql.NullString
	var createdAt, updatedAt string

	err := s.Scan(
		&m.ID, &m.Name, &typ, &url, &config, &m.PingToken, &m.IntervalSeconds, &m.TimeoutSeconds,
		&m.FailureThreshold, &enabled, &state, &lastChecked, &lastPing, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Type = domain.MonitorType(typ)
	m.URL = url.String
	var cfg storedConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("monitor %s: invalid config: %w", m.ID, err)
	}
	m.Config = cfg.toDomain()
	m.Enabled = enabled == 1
	m.CurrentState = domain.MonitorState(state)
	m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	m.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	m.LastCheckedAt = parseNullTime(lastChecked)
	m.LastPingAt = parseNullTime(lastPing)

	return &m, nil
}

// storedConfig é o formato da coluna config; só os campos do tipo são gravados.
type storedConfig struct {
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	GraceSeconds *int   `json:"grace_seconds,omitempty"`
}

func (c *storedConfig) toDomain() domain.MonitorConfig {
	out := domain.MonitorConfig{Host: c.Host, Port: c.Port}
	if c.GraceSeconds != nil {
		out.GraceSeconds = *c.GraceSeconds
	}
	return out
}

func encodeConfig(m *domain.Monitor) string {
	var c storedConfig
	switch m.Type {
	case domain.TypeTCP:
		c.Host, c.Port = m.Config.Host, m.Config.Port
	case domain.TypeHeartbeat:
		g := m.Config.GraceSeconds
		c.GraceSeconds = &g
	}
	b, _ := json.Marshal(c)
	return string(b)
}

func parseNullTime(s sql.NullString) *time.Time {
	if !s.Valid {
		return nil
	}
	t, _ := time.Parse(time.RFC3339, s.String)
	return &t
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}
