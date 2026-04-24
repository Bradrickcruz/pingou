-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE
  monitors (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK (length (name) BETWEEN 1 AND 100),
    url TEXT NOT NULL CHECK (length (url) BETWEEN 1 AND 2048),
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds BETWEEN 10 AND 86400),
    timeout_seconds INTEGER NOT NULL DEFAULT 10 CHECK (timeout_seconds BETWEEN 1 AND 60),
    failure_threshold INTEGER NOT NULL DEFAULT 3 CHECK (failure_threshold BETWEEN 1 AND 10),
    enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
    current_state TEXT NOT NULL DEFAULT 'UNKNOWN' CHECK (current_state IN ('UNKNOWN', 'UP', 'DOWN')),
    last_checked_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );

CREATE INDEX idx_monitors_enabled ON monitors (enabled);

CREATE INDEX idx_monitors_current_state ON monitors (current_state);

CREATE TABLE
  checks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    monitor_id TEXT NOT NULL REFERENCES monitors (id) ON DELETE CASCADE,
    success INTEGER NOT NULL CHECK (success IN (0, 1)),
    status_code INTEGER,
    latency_ms INTEGER NOT NULL,
    error_message TEXT,
    checked_at TEXT NOT NULL
  );

CREATE INDEX idx_checks_monitor_checked ON checks (monitor_id, checked_at DESC);

CREATE INDEX idx_checks_checked_at ON checks (checked_at);

CREATE TABLE
  incidents (
    id TEXT PRIMARY KEY,
    monitor_id TEXT NOT NULL REFERENCES monitors (id) ON DELETE CASCADE,
    started_at TEXT NOT NULL,
    ended_at TEXT,
    last_error TEXT,
    notification_sent INTEGER NOT NULL DEFAULT 0 CHECK (notification_sent IN (0, 1))
  );

CREATE INDEX idx_incidents_monitor_ended ON incidents (monitor_id, ended_at);

CREATE TABLE
  settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );

INSERT INTO
  settings (key, value, updated_at)
VALUES
  (
    'webhook_url',
    '',
    strftime ('%Y-%m-%dT%H:%M:%fZ', 'now')
  ),
  (
    'retention_days',
    '30',
    strftime ('%Y-%m-%dT%H:%M:%fZ', 'now')
  );

-- +goose Down
DROP TABLE settings;

DROP TABLE incidents;

DROP TABLE checks;

DROP TABLE monitors;