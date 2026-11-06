-- Fila de avisos (outbox): o aviso é gravado na mesma transação do check que o
-- gerou, e um worker envia e repete até dar certo ou expirar.
--
-- Uma linha por canal (webhook, email, ntfy): só o canal que falhou repete.
-- payload é um snapshot JSON do evento (sem segredos): o destino (URL do
-- webhook, destinatários, tópico) é lido das settings na hora de enviar.
-- Sem FK em monitor_id: apagar o monitor não cancela o aviso.
-- Sem estado "sending": existe um único worker por banco (lock de arquivo).
-- Datas em RFC 3339, UTC, como o resto do banco.
-- +goose Up
CREATE TABLE notification_outbox (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id        TEXT NOT NULL,  -- UUID, igual em todos os canais e tentativas
    event_type      TEXT NOT NULL,  -- 'down' | 'up' | 'cert_expiring'
    channel         TEXT NOT NULL,  -- 'webhook' | 'email' | 'ntfy'
    monitor_id      TEXT NOT NULL,
    payload         TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sent','expired')),
    attempts        INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    expires_at      TEXT NOT NULL,
    next_attempt_at TEXT NOT NULL,
    last_error      TEXT,
    done_at         TEXT
);

-- Avisos vencidos (a cada rodada do worker).
CREATE INDEX idx_outbox_due ON notification_outbox(next_attempt_at) WHERE status = 'pending';

-- Retenção: apaga sent/expired antigos.
CREATE INDEX idx_outbox_done ON notification_outbox(done_at) WHERE status <> 'pending';

-- +goose Down
DROP INDEX idx_outbox_done;
DROP INDEX idx_outbox_due;
DROP TABLE notification_outbox;
