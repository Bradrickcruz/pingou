-- Janelas de manutenção agendadas, uma por monitor (sem recorrência).
--
-- "Em manutenção" é derivado do relógio (starts_at <= agora < ends_at, não
-- finalizada), não de uma flag: sobrevive a reinício.
-- started_at: o tick já aplicou os efeitos do início (fechar incidente, estado UNKNOWN).
-- finalized_at: o tick (ou o cancelamento) já aplicou os efeitos do fim.
-- Cancelar antes do início apaga a linha; cancelar durante seta ends_at = agora
-- e finalized_at. Por isso não há coluna canceled_at.
-- Datas em RFC 3339, UTC, como o resto do banco.
-- +goose Up
CREATE TABLE maintenance_windows (
    id           TEXT PRIMARY KEY,
    monitor_id   TEXT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    starts_at    TEXT NOT NULL,
    ends_at      TEXT NOT NULL,
    reason       TEXT,
    started_at   TEXT,
    finalized_at TEXT,
    created_at   TEXT NOT NULL,
    CHECK (ends_at > starts_at)
);

-- Uma janela viva (agendada ou em curso) por monitor.
CREATE UNIQUE INDEX idx_maintenance_live_monitor
    ON maintenance_windows(monitor_id) WHERE finalized_at IS NULL;

-- Janela ativa do monitor (a cada check).
CREATE INDEX idx_maintenance_monitor_starts
    ON maintenance_windows(monitor_id, starts_at);

-- +goose Down
DROP INDEX idx_maintenance_monitor_starts;

DROP INDEX idx_maintenance_live_monitor;

DROP TABLE maintenance_windows;
