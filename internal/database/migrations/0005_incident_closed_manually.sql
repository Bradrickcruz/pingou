-- Incidente fechado ao pausar o monitor (sem notificação), distinto de "recuperado".
-- +goose Up
ALTER TABLE incidents
ADD COLUMN closed_manually INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE incidents
DROP COLUMN closed_manually;
