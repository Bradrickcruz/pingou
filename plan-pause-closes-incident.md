# Plano de Desenvolvimento: Pausar monitor fecha o incidente aberto (silenciosamente)

## Overview

**Bug**: ao pausar um monitor que está DOWN, o incidente continua aberto indefinidamente (e o monitor segue com `current_state = DOWN`), embora o scheduler não faça mais checks. O incidente só fecharia por uma transição DOWN→UP, que não acontece com o monitor pausado.

**Comportamento desejado**:
- Ao pausar, o incidente aberto do monitor é fechado **sem notificação** (nem e-mail, nem webhook/ntfy).
- O incidente ganha um estado distinto de "resolvido": **"fechado manualmente"**.
- O estado do monitor volta a `UNKNOWN` ao pausar.

**Fluxo atual (verificado no código)**:
- Pausa = `PATCH /api/monitors/{id}` com `enabled:false` → `MonitorService.Update` ([monitor_service.go:160](internal/service/monitor_service.go#L160)); é o único caminho de pausa (CLI `add`/`rm` não pausam).
- Incidente só fecha em `StateMachine.handleSuccess` ([state_machine.go:85](internal/service/state_machine.go#L85)), que é quem chama o notifier.
- "Aberto" = `ended_at IS NULL` em todas as queries; "resolvido" = `ended_at IS NOT NULL`.
- `handleFailure` só abre incidente se `prevState != DOWN`. **Consequência**: se pausar e retomar com o monitor ainda `DOWN` e a falha persistir, nenhum incidente novo seria aberto → por isso o estado precisa ir a `UNKNOWN` na pausa.

**Decisões (defaults, sem bloquear)**:
- Marcação do incidente: coluna nova `closed_manually INTEGER NOT NULL DEFAULT 0` (booleano, mesmo padrão de `notification_sent`). `ended_at`/`duration_seconds` são preenchidos normalmente (duração = início até a pausa), então stats de uptime continuam corretas sem mudança.
- API: `incidentResponse` ganha `closed_manually: bool`. O filtro `status=resolved` **continua incluindo** os fechados manualmente (aparecem na aba "Resolvidos" com selo "Fechado manualmente"). Sem novo valor de filtro.
- Idempotente: o fechamento roda sempre que `Update` termina com `Enabled == false`, não só na transição; assim um erro no meio se cura na próxima edição/pausa.

Restrições: sem nova dependência; sem mexer em `StateMachine` nem nos notifiers (o silêncio vem de não passar por eles).

## Project Type

**FULLSTACK** — Go (service/repository/SQLite+goose) e SPA React em `web/`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Pausar monitor DOWN: incidente aberto passa a ter `ended_at`, `duration_seconds` e `closed_manually = true`
- [ ] Nenhuma notificação é disparada na pausa (teste com notifier fake)
- [ ] Monitor pausado volta a `current_state = UNKNOWN`
- [ ] Pausar monitor sem incidente aberto não gera erro nem alteração
- [ ] Retomar monitor que seguia fora do ar abre um incidente novo e notifica DOWN normalmente
- [ ] Incidente recuperado normalmente continua com `closed_manually = false`
- [ ] `GET /api/incidents` expõe `closed_manually`; incidente fechado manualmente aparece em "Resolvidos" com selo "Fechado manualmente"
- [ ] Migration aplica e reverte; incidentes existentes ficam `closed_manually = 0`

## Tech Stack

| Technology | Choice Rationale |
|---|---|
| Go + `database/sql` (SQLite) | Stack existente; um `UPDATE` resolve o fechamento |
| goose | Migration `0005` no padrão das anteriores |
| React + Tailwind | Selo na tabela de resolvidos, seguindo `DESIGN.md` |

## File Structure

```
internal/
  database/migrations/0005_incident_closed_manually.sql   (novo)
  domain/incident.go            (+ campo ClosedManually)
  domain/repository.go          (+ IncidentRepository.CloseManually)
  repository/incident_repo.go   (+ CloseManually, coluna em incidentColumns/scan)
  repository/incident_repo_tx.go(+ coluna no scan, se usar incidentColumns)
  service/monitor_service.go    (Update: fecha incidente + UNKNOWN ao pausar)
  handler/incidents.go          (+ closed_manually no JSON)
web/src/components/incidents/ResolvedTable.jsx   (selo)
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
|---|---|---|---|---|---|
| T1 | `design-db` | Migration `closed_manually` | P0 | XS | - |
| T2 | `go-backend-specialist` | Domínio + repositório (`CloseManually`) | P1 | S | T1 |
| T3 | `go-backend-specialist` | Fechar incidente e zerar estado ao pausar | P1 | S | T2 |
| T4 | `go-backend-specialist` | Expor `closed_manually` na API | P1 | XS | T2 |
| T5 | `frontend-specialist` | Selo "Fechado manualmente" | P2 | S | T4 |
| T6 | `code-test` | Testes (service/repo/handler) | P3 | M | T3, T4 |
| T7 | `qa-review` | Verificação manual ponta a ponta | P3 | S | T5, T6 |

### Task T1: Migration

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Última migration é `0004_incidents_page.sql`
- **INPUT**: tabela `incidents` sem marcação de motivo de fechamento
- **OUTPUT**:
  1. `0005_incident_closed_manually.sql`: `ALTER TABLE incidents ADD COLUMN closed_manually INTEGER NOT NULL DEFAULT 0;` e Down com `DROP COLUMN` (mesmo padrão da 0004)
- **VERIFY**:
  - `go test ./internal/database/...` (migrations_test) passa, up e down

### Task T2: Domínio + repositório

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**: coluna existe
- **INPUT**: `Incident` sem o campo; `incidentColumns` e `incidentFields.dest` em [incident_repo.go](internal/repository/incident_repo.go)
- **OUTPUT**:
  1. `Incident.ClosedManually bool`
  2. `closed_manually` adicionado a `incidentColumns`, ao SELECT de `FindAll` e a `dest/apply` (escaneia int → bool)
  3. `IncidentRepository.CloseManually(ctx, monitorID string, endedAt time.Time) error`, um único UPDATE sem leitura prévia:
     ```sql
     UPDATE incidents SET ended_at = ?, closed_manually = 1,
       duration_seconds = MAX(0, CAST(strftime('%s', ?) - strftime('%s', started_at) AS INTEGER))
     WHERE monitor_id = ? AND ended_at IS NULL
     ```
     (no-op se não há incidente aberto)
  4. Conferir `incident_repo_tx.go`: se algum SELECT usar `incidentColumns`, o scan compartilhado já cobre; senão nada a fazer
  5. Atualizar fakes/mocks que implementam `IncidentRepository` (grep nos `_test.go`)
- **VERIFY**:
  - `go build ./...` e `go vet ./...`

### Task T3: Pausa fecha incidente

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**: `CloseManually` disponível
- **INPUT**: `Update` em [monitor_service.go:160-224](internal/service/monitor_service.go#L160)
- **OUTPUT**:
  1. Antes de `s.monitors.Update`: se `!m.Enabled` → `m.CurrentState = domain.StateUnknown` (grava no mesmo UPDATE do monitor; checar que `monitorRepo.Update` persiste `current_state`, caso contrário incluir)
  2. Depois do `s.monitors.Update` bem-sucedido: se `!m.Enabled` → `s.incidents.CloseManually(ctx, m.ID, now)`; erro é retornado
  3. Nenhuma chamada a notifier (o `MonitorService` nem tem acesso a um)
  4. Comentário curto explicando o porquê do `UNKNOWN` (senão retomar DOWN não reabre incidente)
- **VERIFY**:
  - Teste de T6; `go test ./internal/service/...`

### Task T4: API

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T2
- **OUTPUT**:
  1. `ClosedManually bool \`json:"closed_manually"\`` em `incidentResponse` e preenchido em `toIncidentResponse` ([incidents.go](internal/handler/incidents.go))
- **VERIFY**:
  - `curl /api/incidents?status=resolved` mostra o campo

### Task T5: Frontend

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T4
- **PRE-CONDITIONS**: ler `DESIGN.md` para o estilo de selo/badge existente
- **INPUT**: [ResolvedTable.jsx](web/src/components/incidents/ResolvedTable.jsx) (versão desktop e mobile)
- **OUTPUT**:
  1. Selo neutro "Fechado manualmente" nas duas variações (tabela e cartão) quando `i.closed_manually`, reaproveitando o componente de badge de `components/ui` se houver
  2. Tom neutro (não verde de "recuperado")
- **VERIFY**:
  - `cd web && npm run build` e `npm run lint`; conferir no navegador

### Task T6: Testes

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T3, T4
- **OUTPUT**:
  1. Service: monitor DOWN com incidente aberto → `Update(enabled:false)` → incidente fechado, `closed_manually`, duração > 0, monitor `UNKNOWN`
  2. Service: pausar sem incidente aberto → sem erro
  3. Service: nenhuma notificação (notifier fake com contador, igual aos de `state_machine_test.go`)
  4. Scheduler/State machine: após pausa+retomada, nova falha acima do threshold abre incidente novo e notifica
  5. Handler: `closed_manually` no JSON de incidentes (estende `incidents_test.go`); incidente recuperado normal → `false`
- **VERIFY**:
  - `go test ./...` verde

### Task T7: QA manual

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T5, T6
- **OUTPUT**: checklist executado: derrubar um monitor (porta fechada) → incidente aberto → pausar → some de "Abertos", aparece em "Resolvidos" com selo, sem e-mail/ntfy; retomar → volta a abrir incidente se continua fora; heartbeat também
- **VERIFY**: evidências anotadas em "Status de Execução"

## Risk Areas

| Risk | Impact | Mitigation |
|---|---|---|
| Check em andamento termina depois da pausa e reabre incidente | Incidente aberto de monitor pausado | `scheduler.withMonitor` já descarta resultado se `!Enabled` ([scheduler.go:203](internal/scheduler/scheduler.go#L203)); T6 cobre |
| `monitors.Update` ok e `CloseManually` falha | Pausado com incidente aberto | Fechamento idempotente roda a cada `Update` com `Enabled=false`; erro é devolvido ao usuário, que repete |
| Uptime/stats contam o incidente | Distorção | Incidente termina na pausa; `stats.go` usa `ended_at`, sem mudança. Verificar no QA |
| Fakes de `IncidentRepository` quebram no build | Testes não compilam | T2 item 5 |
| Pausas anteriores já deixaram incidentes órfãos abertos no banco | Dado antigo continua errado | Fora de escopo; opcional: no Up da migration fechar incidentes abertos de monitores com `enabled = 0` (ver Notes) |

## Rollback Strategy

1. Reverter o commit (código); `goose down` remove a coluna `closed_manually` (SQLite ≥ 3.35 suporta `DROP COLUMN`, como a 0004 já usa)
2. Incidentes fechados manualmente permanecem como resolvidos normais (têm `ended_at`), sem perda de dados

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm run lint && npm run build && npm test
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 0 | T1 | XS | ~5 min | Migration | ⬜ Pendente |
| Fase 1 | T2, T3, T4 | S | ~45 min | Backend completo | ⬜ Pendente |
| Fase 2 | T5 | S | ~15 min | Selo na UI | ⬜ Pendente |
| Fase 3 | T6, T7 | M | ~45 min | Testes + QA | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T7

## Notes

- T3 é a correção em si; T1, T2, T4 e T5 existem só para o status distinto pedido. Sem o status distinto, T3 sozinho (usando `Close` existente) bastaria.
- Opcional na migration: `UPDATE incidents SET ended_at=..., closed_manually=1 WHERE ended_at IS NULL AND monitor_id IN (SELECT id FROM monitors WHERE enabled=0)` para limpar órfãos já existentes. Decidir na execução; a duração precisaria usar `updated_at` do monitor como aproximação.

### Alternativas consideradas e descartadas

- `Reutilizar StateMachine para fechar` → descartada: ela é acoplada a notifier e a resultado de check; o fechamento manual não é um check.
- `Coluna enum resolution ('recovered'|'manual')` → descartada: um booleano cobre o único caso novo (YAGNI); migra-se para enum se surgir um terceiro motivo.
- `Novo valor de filtro status=closed_manually` → descartada: a UI não precisa filtrar separadamente; selo basta.
- `Deixar ended_at nulo e marcar só por flag` → descartada: quebraria a definição "aberto = ended_at IS NULL" usada em todas as queries.
