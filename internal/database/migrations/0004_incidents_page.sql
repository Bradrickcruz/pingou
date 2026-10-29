-- Página de incidentes: contagem de falhas por incidente, índices para as
-- listas de abertos/resolvidos e limpeza do prefixo `Get "<url>": ` que o
-- *url.Error do Go colocava nas mensagens de erro.
-- +goose Up
ALTER TABLE incidents
ADD COLUMN failure_count INTEGER NOT NULL DEFAULT 0;

-- Backfill aproximado: checks mais antigos que a retenção já foram apagados.
UPDATE incidents
SET
  failure_count = (
    SELECT
      COUNT(*)
    FROM
      checks c
    WHERE
      c.monitor_id = incidents.monitor_id
      AND c.success = 0
      AND c.checked_at >= incidents.started_at
      AND (
        incidents.ended_at IS NULL
        OR c.checked_at < incidents.ended_at
      )
  );

CREATE INDEX idx_incidents_started_at ON incidents (started_at DESC);

CREATE INDEX idx_incidents_ended_at ON incidents (ended_at DESC);

UPDATE incidents
SET
  last_error = substr (last_error, instr (last_error, '": ') + 3)
WHERE
  last_error LIKE 'Get "%": %';

UPDATE checks
SET
  error_message = substr (error_message, instr (error_message, '": ') + 3)
WHERE
  error_message LIKE 'Get "%": %';

-- +goose Down
-- A limpeza do prefixo não é revertida: era só ruído.
DROP INDEX idx_incidents_ended_at;

DROP INDEX idx_incidents_started_at;

ALTER TABLE incidents
DROP COLUMN failure_count;
