-- Aviso de expiração do certificado SSL (monitores http com URL https).
--
-- cert_alert_enabled: o checkbox "avisar quando o certificado estiver perto de expirar".
-- cert_expires_at: NotAfter do último certificado visto (RFC 3339, UTC).
-- cert_alerted_30_for / cert_alerted_7_for: a data de expiração para a qual o
-- aviso de 30 / 7 dias já saiu. Guardar a data, e não um sim/não, faz o aviso
-- "rearmar" sozinho quando o certificado é renovado (a data nova não bate).
-- +goose Up
ALTER TABLE monitors
ADD COLUMN cert_alert_enabled INTEGER NOT NULL DEFAULT 0 CHECK (cert_alert_enabled IN (0, 1));

ALTER TABLE monitors
ADD COLUMN cert_expires_at TEXT;

ALTER TABLE monitors
ADD COLUMN cert_alerted_30_for TEXT;

ALTER TABLE monitors
ADD COLUMN cert_alerted_7_for TEXT;

-- +goose Down
ALTER TABLE monitors
DROP COLUMN cert_alerted_7_for;

ALTER TABLE monitors
DROP COLUMN cert_alerted_30_for;

ALTER TABLE monitors
DROP COLUMN cert_expires_at;

ALTER TABLE monitors
DROP COLUMN cert_alert_enabled;
