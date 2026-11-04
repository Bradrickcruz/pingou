-- Motivo da pausa e contagem de falhas desde a retomada.
--
-- monitors.pause_reason: motivo opcional da pausa (até 200 caracteres), limpo ao retomar.
-- monitors.paused_at / resumed_at: instante da última pausa e da última retomada
-- (RFC 3339, UTC). resumed_at também marca o fim de uma janela de manutenção: a
-- state machine só conta falhas de checks feitos depois dele.
-- incidents.close_note: cópia do motivo no momento do fechamento manual (o motivo
-- da pausa pode mudar depois; o incidente guarda o que valia).
--
-- incidents.closed_manually (INTEGER, sem CHECK) passa a guardar o motivo do
-- fechamento manual: 0 = nenhum, 1 = pausa, 2 = manutenção. Leitores antigos que
-- tratam "diferente de 0" como fechado manualmente continuam corretos.
-- +goose Up
ALTER TABLE monitors
ADD COLUMN pause_reason TEXT;

ALTER TABLE monitors
ADD COLUMN paused_at TEXT;

ALTER TABLE monitors
ADD COLUMN resumed_at TEXT;

ALTER TABLE incidents
ADD COLUMN close_note TEXT;

-- +goose Down
ALTER TABLE incidents
DROP COLUMN close_note;

ALTER TABLE monitors
DROP COLUMN resumed_at;

ALTER TABLE monitors
DROP COLUMN paused_at;

ALTER TABLE monitors
DROP COLUMN pause_reason;
