-- +goose NO TRANSACTION
-- Recria a tabela monitors para aceitar os tipos http, tcp e heartbeat.
--
-- Por que recriar: o SQLite não altera CHECK nem NOT NULL de coluna existente,
-- e url deixa de ser obrigatória (TCP e Heartbeat não têm URL), além de o
-- intervalo máximo subir para 1 semana (heartbeat semanal).
--
-- Por que NO TRANSACTION: PRAGMA foreign_keys não tem efeito dentro de uma
-- transação. Com as FKs ligadas, o DROP TABLE monitors apagaria checks e
-- incidents via ON DELETE CASCADE. Desligadas, os dados ficam intactos.
-- +goose Up
PRAGMA foreign_keys = OFF;

CREATE TABLE
  monitors_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK (length (name) BETWEEN 1 AND 100),
    type TEXT NOT NULL DEFAULT 'http' CHECK (type IN ('http', 'tcp', 'heartbeat')),
    url TEXT CHECK (url IS NULL OR length (url) BETWEEN 1 AND 2048),
    config TEXT NOT NULL DEFAULT '{}',
    ping_token TEXT UNIQUE,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds BETWEEN 10 AND 604800),
    timeout_seconds INTEGER NOT NULL DEFAULT 10 CHECK (timeout_seconds BETWEEN 1 AND 60),
    failure_threshold INTEGER NOT NULL DEFAULT 3 CHECK (failure_threshold BETWEEN 1 AND 10),
    enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
    current_state TEXT NOT NULL DEFAULT 'UNKNOWN' CHECK (current_state IN ('UNKNOWN', 'UP', 'DOWN')),
    last_checked_at TEXT,
    last_ping_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CHECK (type <> 'http' OR url IS NOT NULL),
    CHECK ((type = 'heartbeat') = (ping_token IS NOT NULL))
  );

INSERT INTO
  monitors_new (
    id,
    name,
    type,
    url,
    config,
    interval_seconds,
    timeout_seconds,
    failure_threshold,
    enabled,
    current_state,
    last_checked_at,
    created_at,
    updated_at
  )
SELECT
  id,
  name,
  'http',
  url,
  '{}',
  interval_seconds,
  timeout_seconds,
  failure_threshold,
  enabled,
  current_state,
  last_checked_at,
  created_at,
  updated_at
FROM
  monitors;

DROP TABLE monitors;

ALTER TABLE monitors_new
RENAME TO monitors;

CREATE INDEX idx_monitors_enabled ON monitors (enabled);

CREATE INDEX idx_monitors_current_state ON monitors (current_state);

CREATE INDEX idx_monitors_type ON monitors (type);

PRAGMA foreign_keys = ON;

-- +goose Down
-- Decisão: o Down FALHA se existir monitor tcp ou heartbeat, em vez de apagar
-- dados em silêncio. Remova (ou exporte) esses monitores antes de reverter.
-- json('...') inválido gera erro e interrompe a migration antes de qualquer mudança.
SELECT
  CASE
    WHEN EXISTS (
      SELECT
        1
      FROM
        monitors
      WHERE
        type <> 'http'
    ) THEN json ('down bloqueado: existem monitores tcp/heartbeat')
  END;

PRAGMA foreign_keys = OFF;

CREATE TABLE
  monitors_old (
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

INSERT INTO
  monitors_old (
    id,
    name,
    url,
    interval_seconds,
    timeout_seconds,
    failure_threshold,
    enabled,
    current_state,
    last_checked_at,
    created_at,
    updated_at
  )
SELECT
  id,
  name,
  url,
  interval_seconds,
  timeout_seconds,
  failure_threshold,
  enabled,
  current_state,
  last_checked_at,
  created_at,
  updated_at
FROM
  monitors;

DROP TABLE monitors;

ALTER TABLE monitors_old
RENAME TO monitors;

CREATE INDEX idx_monitors_enabled ON monitors (enabled);

CREATE INDEX idx_monitors_current_state ON monitors (current_state);

PRAGMA foreign_keys = ON;
