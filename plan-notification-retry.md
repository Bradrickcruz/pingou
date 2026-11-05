# Plano de Desenvolvimento: Retry de notificações com outbox (ideia 5)

## Overview

Hoje o aviso (queda, recuperação, SSL) sai **depois do commit** e **uma vez só**: `StateMachine.Process` devolve um `func()` que chama `MultiNotifier`, que dispara uma goroutine por canal. Se o canal falha (webhook/SMTP/ntfy fora do ar) ou o processo cai entre o commit e o envio, o aviso se perde. O webhook ainda registra "sent" mesmo com status 4xx/5xx (`notifier.go:179` só loga o status).

Esta mudança grava o aviso num **outbox** (tabela no SQLite) **na mesma transação** do `Process`, e um **worker** envia e repete até dar certo ou expirar.

**Decisões já tomadas (IDEAS.md, 2026-11-05):**

- Uma linha por canal (webhook, e-mail, ntfy); só o canal que falhou repete.
- Esperas: 30 s, 2 min, 10 min, 30 min. Queda e recuperação expiram em **60 min** (contados do evento); SSL em **24 h**.
- Campo `id` novo no JSON do webhook (mesmo `id` em todas as tentativas, para o receptor ignorar duplicata).
- "Expirado" aparece só no log, não no painel.
- O botão de **teste** (webhook, e-mail, ntfy) continua envio direto, sem fila.

**Restrições:**

- **Não guardar segredo na fila.** O destino (URL do webhook, destinatários, tópico) é lido das `settings` **na hora de enviar**.
- **Nunca segurar a conexão do SQLite durante o envio de rede** (`MaxOpenConns=1`; o worker disputa com a gravação dos checks).
- Texto de e-mail/ntfy montado **na hora de enviar** a partir do snapshot, não gravado pronto (compatível com a ideia 4, i18n).
- Entrega "pelo menos uma vez": duplicata é possível (crash entre enviar e marcar).
- Sem fila externa, sem tela de reenvio, sem canais novos.

## Project Type

**BACKEND** — Go + SQLite (migration nova, repositório, worker, refatoração dos 3 notificadores, wiring, shutdown). Nenhuma mudança no frontend.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Mudança de estado e linhas do outbox são gravadas na **mesma transação**: rollback não deixa aviso; commit sempre deixa.
- [ ] Canal que falha (erro de rede ou status não-2xx) é repetido nas esperas 30 s / 2 min / 10 min / 30 min; os demais canais do mesmo evento não repetem.
- [ ] Queda/recuperação passam para `expired` após 60 min do evento; SSL após 24 h; log `warn` com `event_id`, `channel`, `attempts`.
- [ ] Reiniciar o Pingou (inclusive `kill -9`) retoma o que estava pendente, sem perder e sem reenviar o que já saiu.
- [ ] JSON do webhook ganha `id` (UUID), igual em todas as tentativas do mesmo evento; demais campos intactos.
- [ ] Nenhuma chamada de rede com a conexão do banco presa (worker lê, solta, envia, grava).
- [ ] O outbox não guarda URL do webhook, destinatários, tópico, token do ntfy nem `ping_token`.
- [ ] Linhas `sent`/`expired` com mais de 7 dias são apagadas pela rotina de retenção.
- [ ] Botões de teste seguem diretos, sem passar pela fila.
- [ ] Desligamento limpo: worker para, envio em andamento termina dentro do `shutdownTimeout`, pendentes ficam no banco.
- [ ] README, `docs/implementacoes.md` e `IDEAS.md` atualizados.

## Tech Stack

| Technology                                              | Choice Rationale                                                                                                             |
| ------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| SQLite + goose (migration `0010`)                       | A ideia pede "fila no próprio SQLite"; migration nova é esperada (ver Context da ideia).                                     |
| Worker único (goroutine) com ticker + canal de "acorda" | Uma instância por banco (há lock de arquivo), então **não precisa de estado `sending`/claim**: sem linha "presa" após crash. |
| JSON em coluna `payload`                                | Snapshot do evento sem uma coluna por campo; eventos futuros (ideia 6) só trazem payload novo.                               |
| `newUUIDv7()` existente (`uuid.go`)                     | Reuso para o `id` do evento.                                                                                                 |
| `go test` + relógio injetável                           | Padrão já usado (`StateMachine.WithClock`). Sem biblioteca de fila/retry.                                                    |

## File Structure

```
internal/database/migrations/0010_notification_outbox.sql   # novo
internal/database/migrations_test.go                        # cobre 0010
internal/domain/outbox.go                                   # novo: OutboxEvent, snapshot, constantes de canal/status
internal/repository/outbox_repo.go                          # novo: leitura/atualização do worker + purge
internal/repository/outbox_repo_tx.go                       # novo: Enqueue dentro da transação
internal/service/unit_of_work.go / _impl.go                 # + OutboxRepo()
internal/service/state_machine.go                           # grava no outbox em vez de chamar notifier
internal/service/notifier.go                                # Channel interface; webhook devolve erro; campo id
internal/service/email_notifier.go, ntfy_notifier.go        # Send devolve erro; montam texto do snapshot
internal/service/multi_notifier.go                          # REMOVIDO (fanOut/wg viram o worker)
internal/service/outbox_worker.go                           # novo: loop, backoff, expiração
internal/scheduler/retention.go                             # + purge do outbox
internal/scheduler/scheduler.go                             # troca notifier por (canais, wake)
cmd/pingou/commands/serve.go                                # monta canais + worker; etapa de shutdown
README.md, docs/implementacoes.md, IDEAS.md                 # docs
```

## Task Breakdown

| Task ID | Agent                   | Name                                                     | Priority | Complexity | Dependencies |
| ------- | ----------------------- | -------------------------------------------------------- | -------- | ---------- | ------------ |
| T1      | `design-db`             | Migration `notification_outbox`                          | P0       | S          | -            |
| T2      | `go-backend-specialist` | Domínio e repositórios do outbox                         | P0       | M          | T1           |
| T3      | `go-backend-specialist` | Canais devolvem erro + campo `id`                        | P1       | L          | -            |
| T4      | `go-backend-specialist` | StateMachine grava no outbox (UoW)                       | P1       | M          | T2           |
| T5      | `go-backend-specialist` | Worker: envio, backoff e expiração                       | P1       | L          | T2, T3       |
| T6      | `go-backend-specialist` | Wiring, shutdown, retenção e remoção do MultiNotifier    | P1       | M          | T4, T5       |
| T7      | `code-test`             | Testes (unidade e integração)                            | P2       | L          | T4, T5, T6   |
| T8      | `security-review`       | Segredos no outbox e nos logs                            | P3       | S          | T6           |
| T9      | `performance`           | Carga com conexão única                                  | P3       | M          | T6           |
| T10     | `qa-review`             | Verificação ponta a ponta (receptor instável, `kill -9`) | P3       | M          | T6, T7       |
| T11     | `document`              | README, histórico e IDEAS                                | P3       | S          | T6           |

### Task T1: Migration `notification_outbox`

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Última migration é `0009_maintenance_windows.sql`.
- **INPUT**: banco sem tabela de fila.
- **OUTPUT**:
  1. `0010_notification_outbox.sql` (`-- +goose Up/Down`, datas RFC 3339 UTC como no resto do banco):
     ```sql
     CREATE TABLE notification_outbox (
       id              INTEGER PRIMARY KEY AUTOINCREMENT,
       event_id        TEXT NOT NULL,      -- UUID, igual em todos os canais e tentativas
       event_type      TEXT NOT NULL,      -- 'monitor_down' | 'monitor_up' | 'cert_expiring' (ver Notes)
       channel         TEXT NOT NULL,      -- 'webhook' | 'email' | 'ntfy'
       monitor_id      TEXT NOT NULL,      -- sem FK: apagar o monitor não cancela o aviso
       payload         TEXT NOT NULL,      -- snapshot JSON (sem segredos)
       status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sent','expired')),
       attempts        INTEGER NOT NULL DEFAULT 0,
       created_at      TEXT NOT NULL,
       expires_at      TEXT NOT NULL,
       next_attempt_at TEXT NOT NULL,
       last_error      TEXT,
       done_at         TEXT
     );
     CREATE INDEX idx_outbox_due ON notification_outbox(next_attempt_at) WHERE status = 'pending';
     CREATE INDEX idx_outbox_done ON notification_outbox(done_at) WHERE status <> 'pending';
     ```
  2. Teste em `migrations_test.go` (sobe e desce; índices parciais existem).
- **VERIFY**:
  - `go test ./internal/database/...` passa.
  - `EXPLAIN QUERY PLAN` da busca "vencidos" usa `idx_outbox_due` (padrão de `queryplan_test.go`).
- **ROLLBACK**: migration `Down` (`DROP INDEX` ×2, `DROP TABLE`).

### Task T2: Domínio e repositórios do outbox

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Migration `0010` aplicável.
- **INPUT**: padrão `CheckRepoTx`/`CheckRepo` (versão com `*sql.Tx` e versão com `*sql.DB`).
- **OUTPUT**:
  1. `domain/outbox.go`: `OutboxEvent` (linha), `EventSnapshot` (JSON: `monitor{id,name,type,target…}`, `incident{started_at,ended_at,duration_seconds,last_error,failure_count}`, `cert{expires_at,days_left,threshold_days}`, `occurred_at`). Campos do monitor = o mínimo para `Target()` e para os textos de e-mail/ntfy (conferir `domain/monitor.go`). **Sem** `ping_token`.
  2. `OutboxRepoTx.EnqueueWithTx(ctx, tx, events)`.
  3. `OutboxRepo`: `Due(ctx, now, limit)`, `MarkSent`, `MarkRetry(id, attempts, next, lastError)`, `MarkExpired`, `PurgeDoneBefore(ctx, before)`; `last_error` truncado em 200 caracteres.
- **VERIFY**:
  - Testes de repositório com SQLite em memória: enqueue → due → mark; `Due` respeita `status` e `next_attempt_at`; `Due` não devolve linha futura.
- **ROLLBACK**: apagar arquivos novos.

### Task T3: Canais devolvem erro + campo `id`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Nenhuma (pode andar em paralelo com T1/T2).
- **INPUT**: `WebhookNotifier`, `EmailNotifier`, `NtfyNotifier` com `NotifyX` sem retorno; `send` só loga.
- **OUTPUT**:
  1. Interface `Channel { Name() string; Send(ctx, ev domain.OutboxEvent) error }` (substitui `Notifier`).
  2. Webhook: monta o payload do snapshot; erro de rede **ou** status fora de 2xx vira erro (hoje 500 conta como "sent"). Destino lido de `getWebhookURL()` na hora; **vazio = `nil` sem enviar** (a linha vira `sent`/pulada, sem retry eterno).
  3. Campo `"id": "<event_id>"` no `webhookPayload`; no `SendTestWebhook` um UUID novo. `timestamp` do `cert_expiring` passa a ser o `occurred_at` do evento (antes `time.Now()` no envio); idem `NotifyRecovery`/`NotifyDown` (já usam o incidente).
  4. E-mail e ntfy: reconstroem `Monitor`/`Incident` a partir do snapshot e **reusam** `subjectDown`, `bodyFor`, `certBody`, `ntfyTitle` sem mudar texto; erro de `sendMail`/`publishNtfy` é devolvido (já passa por `smtpErrorMessage`/`ntfyErrorMessage`, que tiram URL/segredo).
  5. Botões de teste não mudam.
- **VERIFY**:
  - `go test ./internal/service/ -run 'Webhook|Email|Ntfy'` passa, com casos novos: 500 → erro; 204 → ok; URL vazia → nil sem chamada.
  - Textos de e-mail/ntfy idênticos aos testes atuais (golden dos testes existentes continua verde).
- **ROLLBACK**: `git checkout internal/service/*notifier*.go`.

### Task T4: StateMachine grava no outbox

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - `OutboxRepoTx` pronto.
- **INPUT**: `Process` devolve `notify func()` e `certNotify func()` chamados depois do commit; `NewStateMachine(uow, notifier)`.
- **OUTPUT**:
  1. `UnitOfWork.OutboxRepo()` + `NewUnitOfWork(... outboxRepo)` (único ponto de criação: `scheduler.go:276`, mais os testes).
  2. `handleSuccess`, `handleFailure` e `processCert` passam a devolver o **evento** (snapshot + tipo + prazo: 60 min / 24 h) em vez de `func()`. `Process` grava **uma linha por canal** (mesmo `event_id`) antes do `Commit`.
  3. `NewStateMachine(uow, channels []string, wake func())`: `wake` é chamado **depois** do commit, sem bloquear (canal com buffer 1 + `select default`).
  4. Os comentários sobre "notificar após commit por causa do lock" saem; o motivo agora é "o worker lê fora da transação".
- **VERIFY**:
  - Teste: erro depois de enfileirar (ex.: falha no `Commit`/`UpdateCertWithTx` simulada) → nenhuma linha no outbox.
  - Teste: queda gera N linhas (N canais) com o mesmo `event_id`; check sem mudança de estado gera 0 linhas.
  - `state_machine_test.go` e `cert_state_machine_test.go` adaptados e verdes.
- **ROLLBACK**: `git checkout internal/service/state_machine.go internal/service/unit_of_work*.go`.

### Task T5: Worker — envio, backoff e expiração

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: T2, T3
- **PRE-CONDITIONS**:
  - `OutboxRepo` e `Channel` prontos.
- **INPUT**: nada equivalente hoje.
- **OUTPUT**:
  1. `OutboxWorker{repo, channels map[string]Channel, now func() time.Time, wake chan struct{}}` com `Start(ctx)`, `Wait(ctx)` (mesmo estilo do `RetentionWorker`) e `Wake()`.
  2. Laço: roda no início (retoma pendentes), a cada `wake` e a cada 30 s (cobre o `next_attempt_at` futuro).
  3. Por rodada: `Due(now, 20)` (conexão solta ao voltar) → para cada linha: se `now > expires_at` → `MarkExpired` + `slog.Warn`; senão `Send` com timeout próprio → `MarkSent` ou `MarkRetry`. **Nenhum `Send` com transação aberta.**
  4. Espera = `[30s, 2m, 10m, 30m][min(attempts-1, 3)]`; se `now + espera > expires_at` → `expired` na hora (não adianta agendar). SSL (24 h) segue repetindo a cada 30 min.
  5. Canal desconhecido na linha (ex.: SMTP desligado depois) → `expired` com log, sem pânico. `recover()` em volta do `Send`.
  6. Logs por tentativa: `event_id`, `event_type`, `channel`, `attempt`, resultado e `err` já sanitizado. Nunca URL, tópico, e-mail.
  7. Constantes nomeadas (`outboxBackoff`, `outboxBatch`, `outboxPoll`) com comentário `ponytail:` onde houver teto conhecido.
- **VERIFY**:
  - Teste de tabela do backoff e da expiração com relógio falso (sem `time.Sleep`).
  - Teste: webhook falha, e-mail ok → só a linha do webhook segue `pending`.
- **ROLLBACK**: apagar `outbox_worker.go`.

### Task T6: Wiring, shutdown, retenção e remoção do MultiNotifier

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T4, T5
- **PRE-CONDITIONS**:
  - T4 e T5 verdes isoladamente.
- **INPUT**: `serve.go:91-120` monta webhook + `MultiNotifier`; etapa de shutdown `"notifications"` chama `notifier.Wait`; `scheduler.NewScheduler(..., notifier)`.
- **OUTPUT**:
  1. `serve.go`: monta `Channel`s (webhook sempre; e-mail se SMTP; ntfy se configurado), o `OutboxWorker` e passa `channels` + `worker.Wake` ao scheduler/StateMachine.
  2. Shutdown: etapa `"notifications"` vira `worker.Wait(ctx)` (para de pegar linhas novas; termina o envio em curso; o resto fica `pending` no banco). Mantém a ordem: scheduler → http → checks → **notificações** → database → lock.
  3. `RetentionWorker.purge` também chama `outbox.PurgeDoneBefore(now − 7 dias)` (valor fixo; sem setting nova).
  4. `multi_notifier.go` e seu teste removidos.
  5. Linha de log na subida: quantos pendentes foram retomados.
- **VERIFY**:
  - `go build ./... && go vet ./...` limpos; `grep -rn MultiNotifier` vazio.
  - `shutdown_test.go` adaptado e verde.
- **ROLLBACK**: `git revert` do commit; a tabela `notification_outbox` pode ficar (sem uso) ou descer a migration.

### Task T7: Testes

- **Agent**: `code-test`
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T4, T5, T6
- **PRE-CONDITIONS**:
  - Implementação completa.
- **INPUT**: testes unitários parciais das tasks anteriores.
- **OUTPUT**:
  1. Integração com `httptest.Server` que falha N vezes e depois responde 200: confere o número de tentativas, que o `id` do JSON é o mesmo em todas, e que a linha termina `sent`.
  2. "Reinício": cria o worker, enfileira, descarta o worker, cria outro sobre o mesmo banco → envia uma vez.
  3. Expiração: queda com falha contínua vira `expired` aos 60 min (relógio falso); SSL só aos 24 h.
  4. Status 4xx/5xx = falha; URL vazia = sem envio e sem retry.
  5. Purge apaga só `sent`/`expired` antigos.
  6. Teste de concorrência (`-race`): `Process` gravando checks enquanto o worker processa.
- **VERIFY**:
  - `go test -race ./...` verde; cada teste novo falha se o comportamento correspondente for quebrado (checar um por um revertendo a linha).
- **ROLLBACK**: apagar testes novos.

### Task T8: Segredos no outbox e nos logs

- **Agent**: `security-review`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T6
- **PRE-CONDITIONS**: T6 concluída.
- **INPUT**: tabela e logs novos.
- **OUTPUT**:
  1. Conferência: `payload` e `last_error` sem URL do webhook, destinatários, tópico, token do ntfy, `ping_token`, senha SMTP.
  2. Registrar o risco residual: `target` do monitor pode ter token na query (já é assim no payload do webhook e no banco). Vale para o export (ver Notes).
- **VERIFY**:
  - `sqlite3 pingou.db 'select payload,last_error from notification_outbox'` após um cenário com falha: nenhum segredo.
  - `grep` nos logs do cenário: sem URL completa.
- **ROLLBACK**: n/a.

### Task T9: Carga com conexão única

- **Agent**: `performance`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T6
- **PRE-CONDITIONS**: T6 concluída; perfil do `dbbench_test.go` disponível.
- **INPUT**: perfil de 100 monitores × 10 s.
- **OUTPUT**:
  1. Medir latência da gravação de checks com o outbox vazio e com 1 000 linhas `pending` (receptor fora do ar): `Due` + `MarkRetry` em lotes de 20 não deve atrasar `Process`.
  2. RAM no cenário grande não piora (a meta de 50 MB já está em "Candidatas").
- **VERIFY**:
  - Números antes/depois colados em Evidências. Se atrasar: reduzir lote ou aumentar `outboxPoll`.
- **ROLLBACK**: n/a.

### Task T10: Verificação ponta a ponta

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T6, T7
- **PRE-CONDITIONS**: build local; receptor webhook local controlável (script de 10 linhas).
- **INPUT**: Pingou real com um monitor que cai de propósito.
- **OUTPUT**:
  1. Receptor responde 500 por ~3 min e depois 200: o aviso chega depois, uma vez, com o mesmo `id`.
  2. `kill -9` logo após a queda, antes do envio: ao subir, o aviso sai.
  3. Receptor fora por mais de 60 min (relógio ajustado ou prazo reduzido em teste): vira `expired`, log claro, nada chega depois.
  4. Webhook com falha e ntfy ok: ntfy chega uma vez, não repete.
  5. Trocar a URL do webhook durante as tentativas: a próxima vai para a URL nova.
  6. Botões de teste continuam diretos.
- **VERIFY**:
  - Logs e JSONs recebidos colados em Evidências.
- **ROLLBACK**: n/a.

### Task T11: README, histórico e IDEAS

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T6
- **PRE-CONDITIONS**: comportamento final confirmado (T10 ideal).
- **INPUT**: README diz que o webhook sai depois do commit e não repete.
- **OUTPUT**:
  1. README (Webhook, e-mail, ntfy): tentativas e esperas, prazos de 60 min / 24 h, "pelo menos uma vez", campo `id` e como ignorar duplicata. Atualizar a "limitação conhecida" de aviso perdido.
  2. `docs/implementacoes.md`: item novo + linha na tabela do topo.
  3. `IDEAS.md`: tirar a ideia 5 do resumo e do corpo; nas ideias 4 e 6 trocar as referências para "já entregue"; na ideia 1, deixar explícito que `notification_outbox` fica fora do export.
- **VERIFY**:
  - `grep -n 'ideia 5' IDEAS.md` só em contexto histórico; links do README válidos.
- **ROLLBACK**: `git checkout README.md docs/implementacoes.md IDEAS.md`.

## Risk Areas

| Risk                                                        | Impact                                                 | Mitigação                                                                            |
| ----------------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| Worker segura a conexão única durante o envio               | Gravação de checks para; falso "fora do ar" em cascata | Ler lote → soltar → enviar → gravar; teste com `-race` e T9                          |
| Aviso duplicado (crash entre enviar e marcar)               | Receptor vê o mesmo evento 2×                          | `id` estável no webhook; documentado como "pelo menos uma vez"                       |
| Webhook devolve 4xx permanente (URL errada)                 | 5 tentativas inúteis por evento                        | Aceito: expira em 60 min, log mostra o status; sem retry infinito                    |
| Linha de canal desligado depois (SMTP removido do ambiente) | Fica `pending` até expirar                             | Worker marca `expired` com log quando o canal não existe                             |
| Evento enfileirado sem destino configurado                  | Linhas inúteis                                         | Destino vazio na hora do envio = encerra sem tentar                                  |
| Refatoração dos 3 canais muda texto/payload sem querer      | Quebra de integração                                   | Testes de texto existentes servem de golden; JSON só ganha `id`                      |
| Snapshot guarda `target` com token na query                 | Segredo extra no banco/export                          | Já está em `monitors`; fica fora do export quando a ideia 1 entrar; registrado em T8 |
| Ideia 7 (nomes de evento) entra depois e muda `event_type`  | Linhas antigas com nome velho                          | Linhas vivem ≤ 24 h; mapear nome no envio ou entregar a 7 antes                      |
| Tabela cresce                                               | Banco maior, `VACUUM INTO` mais lento                  | Purge de 7 dias na rotina horária; índice parcial                                    |

## Rollback Strategy

1. `git revert` do(s) commit(s). O `MultiNotifier` volta a ser o caminho de envio.
2. Migration: rodar o `Down` da `0010` (`DROP TABLE notification_outbox`). Sem ela no código, a tabela fica órfã e inofensiva.
3. Avisos pendentes no momento do rollback são perdidos (mesmo comportamento de hoje).
4. Quem consome o webhook ignora o campo `id` extra; nada a desfazer do lado deles.

## Phase X: Verification

```bash
# Build, vet e testes (com detector de corrida)
go build ./... && go vet ./...
go test -race ./...

# Nada do caminho antigo
grep -rn "MultiNotifier\|fanOut" internal cmd

# Front não é afetado, mas confirma
cd web && npm run build

# Esquema
sqlite3 pingou.db '.schema notification_outbox'
```

## Cronograma Sugerido

| Fase   | Tarefas                     | Complexity | Esforço | Target                                 | Status      |
| ------ | --------------------------- | ---------- | ------- | -------------------------------------- | ----------- |
| Fase 0 | T1, T2, T3 (T3 em paralelo) | L          | ~2 h    | Tabela, repositórios e canais com erro | ⬜ Pendente |
| Fase 1 | T4, T5                      | XL         | ~2 h    | Enfileira na transação e worker        | ⬜ Pendente |
| Fase 2 | T6, T7                      | XL         | ~2 h    | Wiring, shutdown e testes              | ⬜ Pendente |
| Fase 3 | T8, T9, T10, T11            | L          | ~1,5 h  | Segurança, carga, QA e docs            | ⬜ Pendente |

## Status de Execução

### Concluído

- Leitura de `state_machine.go`, `multi_notifier.go`, `notifier.go`, `unit_of_work.go`, `retention.go`, `serve.go` e das migrations `0008`/`0009`; mapa dos pontos a mudar. Só há 3 chamadores de notificação (todos em `state_machine.go`).

### Em andamento

-

### Pendente

- T1 a T11

### Evidências

-

## Notes

- **Ordem com a ideia 7 (`plan-webhook-event-names.md`).** Recomendado entregar a 7 antes (ou junto): quem consome o webhook ajusta o fluxo uma vez só (nomes novos + `id`). Este plano grava em `event_type` o nome que o webhook usa; se a 7 ainda não entrou, usa `down`/`up`/`cert_expiring` e a troca é só nas constantes.
- **Sem estado `sending`.** Existe um único processo por banco (lock de arquivo), então um único worker lê e envia em sequência. Ganho: crash nunca deixa linha "presa". Custo: envio sequencial (webhook 10 s de timeout pode atrasar o e-mail da fila). Aceitável para o volume do Pingou; se virar problema, um worker por canal (`ponytail:` no código).
- **Por que snapshot e não só `monitor_id`.** O monitor pode ser editado ou apagado antes da tentativa; o aviso deve dizer o que valia no evento. O destino, ao contrário, **é** lido na hora (decisão da ideia).
- **Ideia 6 (manutenção)** vira só mais um `event_type` com payload próprio; este desenho já comporta, e resolve o "não repetir aviso" dela.
- **Ideia 1 (backup):** o export atual (`VACUUM INTO`) leva a tabela; quando a ideia 1 entrar, `notification_outbox` fica de fora e o `pingou import` não a espera.
- **Ideia 4 (i18n):** e-mail e ntfy já nascem montados no envio, no idioma da hora.
- `timestamp` do `cert_expiring` muda de "hora do envio" para "hora do evento": diferença só aparece em retry, e é a semântica correta.

### Alternativas consideradas e descartadas

- Reenviar com goroutine + `time.AfterFunc` em memória → descartada: perde tudo se o processo cair (o problema que a ideia quer resolver).
- Gravar o texto pronto do e-mail/ntfy na fila → descartada: conflita com i18n e guardaria endereços.
- Linha única por evento, com estado de cada canal em JSON → descartada: atualização parcial e consulta mais difíceis do que uma linha por canal.
- Estado `sending` com reclamação de linha → descartada: desnecessário com worker único; só cria linha presa após crash.
- Fila externa (Redis, RabbitMQ) → fora do escopo da ideia.
- Retry infinito por canal → descartada: aviso de queda velho confunde; prazos de 60 min / 24 h já decididos.
