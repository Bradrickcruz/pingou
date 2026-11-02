-- Índice de checks por monitor, na ordem das consultas e cobrindo o /stats.
--
-- Substitui idx_checks_monitor_checked (monitor_id, checked_at DESC):
--
-- 1. Ordem: as consultas pedem ORDER BY checked_at DESC, id DESC. O índice
--    antigo guarda checked_at DESC com o rowid crescente, que não serve para
--    essa ordem, e o SQLite ordenava à parte (USE TEMP B-TREE FOR LAST TERM OF
--    ORDER BY). Com (checked_at, id) crescentes, o SQLite lê o índice de trás
--    para frente e já sai na ordem certa: página de checks, histórico da
--    lista do painel e falhas seguidas da state machine (a cada check).
-- 2. Cobertura: success e latency_ms no índice deixam o /stats (série e
--    totais) e o uptime 24 h da lista lerem só o índice, sem buscar cada
--    linha na tabela.
--
-- idx_checks_checked_at (retenção) continua como está.
-- +goose Up
CREATE INDEX IF NOT EXISTS idx_checks_monitor_time ON checks (monitor_id, checked_at, id, success, latency_ms);

DROP INDEX IF EXISTS idx_checks_monitor_checked;

-- +goose Down
CREATE INDEX IF NOT EXISTS idx_checks_monitor_checked ON checks (monitor_id, checked_at DESC);

DROP INDEX IF EXISTS idx_checks_monitor_time;
