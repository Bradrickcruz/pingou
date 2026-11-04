# Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)

## Overview

**O quê.** Uma placa de "em manutenção" num monitor. Enquanto ela está lá, o Pingou não confere nada, não abre incidente e não avisa ninguém. Quando sai, o monitor volta ao normal.

**Duas etapas, um PR:**

1. **Pausa com motivo.** Motivo opcional (até 200 caracteres), instante da pausa e da retomada guardados. A contagem de falhas passa a olhar só checks feitos depois da retomada. Isso corrige um bug real: hoje falhas de antes da pausa somam com as de depois (`countConsecutiveFails` em `internal/service/state_machine.go` lê os últimos N checks sem filtro de data).
2. **Janela agendada.** Início, fim e motivo opcional. Janela única (sem recorrência), por monitor, com criação em lote para vários monitores.

**Decisões fechadas (aprovadas na análise de viabilidade):**

| #   | Decisão                                                                                                                       |
| --- | ----------------------------------------------------------------------------------------------------------------------------- |
| D1  | Teto de **7 dias** por janela. Início no futuro, fim depois do início.                                                        |
| D2  | **Uma janela viva por monitor** (agendada ou em curso). "Próxima janela marcada" é essa mesma janela.                         |
| D3  | **Pausa manual vence a janela.** Monitor pausado continua pausado quando a janela termina.                                    |
| D4  | Motivo **opcional** na pausa e na janela, máx. 200 caracteres.                                                                |
| D5  | Criação em lote **tudo ou nada**: se um monitor já tem janela viva, nada é criado e o erro diz qual.                          |
| D6  | **"Em manutenção" é derivado do relógio** (`starts_at <= agora < ends_at`, não cancelada), não de flag. Sobrevive a reinício. |
| D7  | Pingou roda com **uma instância** (`dbMu` + SQLite `MaxOpenConns=1`). Sem lock distribuído.                                   |
| D8  | Uptime de 24 h não muda: sem checks na janela, não há o que contar.                                                           |
| D9  | `closed_manually` (INTEGER sem CHECK) passa a guardar 0/1/2 = nenhum/pausa/manutenção. Sem migração de esquema para isso.     |
| D10 | No modal de agendar, o rótulo de Início e Fim mostra o **nome do fuso** das Configurações (ex.: `America/Sao_Paulo`), não o offset (`-03`). |
| D11 | No modal de agendar, monitores **pausados** ou com **janela viva** (agendada ou em curso) aparecem na lista **desabilitados** (selo "pausado" / "já agendada"), pois não podem entrar numa janela. |

**Fora do escopo:** janela recorrente, aviso de "começou/terminou", página pública de status, limpeza de incidentes órfãos, edição de janela (cancela e cria outra), revisão dos baldes do `/stats` em UTC.

**Fluxo de git:** branch `feature/maintenance-windows` a partir de `develop`, commits por etapa, PR para `develop`.

## Project Type

**FULLSTACK** — backend Go (SQLite + goose, scheduler próprio, state machine) e SPA React/Vite embutida no binário.

## Definition of Done (DoD)

Este plano está concluído quando:

- [x] Todos os Success Criteria estão marcados
- [x] Build, lint e type-check passam
- [x] Nenhum item pendente no Status de Execução
- [x] Evidências de verificação registradas
- [x] PR aberto contra `develop` com descrição e checklist de teste manual

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

**Etapa 1**

- [x] Pausar abre um modal com motivo opcional; motivo e `paused_at` são gravados
- [x] Lista e detalhes mostram "Pausado" com o motivo
- [x] Retomar grava `resumed_at`; falhas anteriores à retomada não contam para o threshold (teste cobre: threshold 3, 2 falhas antes da pausa, 1 depois → sem incidente)
- [x] Incidente fechado pela pausa mostra "Fechado manualmente" com o motivo

**Etapa 2**

- [x] `POST /api/maintenance-windows` cria uma janela por monitor, tudo ou nada, validando D1, D2, D4
- [x] Durante a janela: nenhum check gravado, nenhum incidente, nenhum aviso (e-mail, webhook, ntfy, SSL); ping de heartbeat responde 200 sem registrar
- [x] Incidente aberto no início da janela fecha em silêncio com "Fechado por manutenção" e o motivo
- [x] Ao terminar: estado `UNKNOWN`, primeiro check roda na hora, contagem de falhas recomeça, prazo do heartbeat recomeça de "agora"
- [x] Pausa manual vence a janela (D3)
- [x] Cancelar antes do início apaga a janela; durante, encerra na hora e o monitor volta
- [x] Reiniciar o Pingou no meio da janela mantém o monitor em manutenção e aplica início/fim pendentes
- [x] Lista mostra "Em manutenção até HH:MM" (fuso das Configurações); detalhes mostram motivo, início, fim e botão de cancelar
- [x] Datas digitadas são interpretadas no fuso das Configurações, não no do navegador (teste cobre fuso diferente de UTC)
- [x] Rótulos de Início/Fim do modal mostram o nome do fuso das Configurações (D10)
- [x] Monitores pausados ou com janela viva aparecem desabilitados na lista do modal (D11)

## Tech Stack

| Technology                               | Choice Rationale                                                                             |
| ---------------------------------------- | -------------------------------------------------------------------------------------------- |
| Go + `net/http` (padrão 1.22 de rotas)   | Já usado; rotas novas seguem `server.go`                                                     |
| SQLite + goose (migrações embutidas)     | Já usado; migrações 0008 e 0009                                                              |
| Scheduler próprio (`internal/scheduler`) | O tick de manutenção segue o molde do `runHeartbeats` (um ticker, sem goroutine por monitor) |
| React 19 + Vite + Tailwind               | Já usado; reaproveita `Modal`, `Combobox`, `Badge`                                           |
| `<input type="datetime-local">`          | Nativo, sem lib de date picker. Conversão para UTC feita com helper em `lib/timezone.js`     |
| `node --test`                            | Testes de `web/src/lib/*.test.js` já existentes                                              |

Nenhuma dependência nova.

## File Structure

```
internal/
  database/migrations/
    0008_pause_reason.sql            # NOVO: monitors.pause_reason/paused_at/resumed_at, incidents.close_note
    0009_maintenance_windows.sql     # NOVO: tabela maintenance_windows + índices
  domain/
    monitor.go                       # + PauseReason, PausedAt, ResumedAt
    incident.go                      # ClosedManually bool → CloseReason (enum) + CloseNote
    maintenance.go                   # NOVO: MaintenanceWindow + Active(now)
    repository.go                    # + MaintenanceRepository; CloseManually ganha reason/note
  repository/
    monitor_repo.go                  # colunas novas no Update/scan
    incident_repo.go                 # CloseManually(reason, note); scan de close_note
    check_repo_tx.go                 # Latest(…, since *time.Time)
    maintenance_repo.go              # NOVO
    queries.go                       # SQL quente da janela ativa (entra no queryplan_test)
  service/
    monitor_service.go               # motivo, paused_at/resumed_at; Ping respeita janela
    maintenance_service.go           # NOVO: CreateBatch, Cancel, validação
    state_machine.go                 # countConsecutiveFails com since = resumed_at
    validate.go                      # validações de motivo e janela
    unit_of_work.go                  # assinatura de Latest
  scheduler/
    scheduler.go                     # gate de manutenção em withMonitor; relógio injetável
    maintenance.go                   # NOVO: tick de início/fim de janela
  handler/
    monitors.go                      # pause_reason no PATCH; bloco maintenance na resposta
    incidents.go                     # close_reason + close_note na resposta
    maintenance.go                   # NOVO: POST/DELETE /api/maintenance-windows
    server.go                        # rotas novas
cmd/pingou/commands/serve.go         # wiring do MaintenanceService/repo

web/src/
  api/maintenance.js                 # NOVO
  lib/timezone.js                    # + zonedToUtc (wall time no fuso → instante UTC)
  lib/monitors.js                    # status MAINTENANCE, severidade, contadores, rótulo
  components/monitors/MonitorModals.jsx  # + PauseModal, MaintenanceModal
  components/monitors/{RowMenu,MonitorRow,MonitorCard,StatusBanner}.jsx
  components/monitor-detail/MaintenanceCard.jsx  # NOVO
  components/incidents/ResolvedTable.jsx         # selo por close_reason + motivo
  pages/{Dashboard,MonitorDetail}.jsx

docs/implementacoes.md, README.md (API)
```

## Task Breakdown

| Task ID        | Agent                           | Name                                                      | Priority | Complexity | Dependencies  |
| -------------- | ------------------------------- | --------------------------------------------------------- | -------- | ---------- | ------------- |
| T0             | —                               | Criar branch `feature/maintenance-windows`                | P0       | XS         | -             |
| **Etapa 1**    |                                 |                                                           |          |            |               |
| T1             | `design-db`                     | Migração 0008 (pausa + close_note)                        | P0       | S          | T0            |
| T2             | `go-backend-specialist`         | Domínio + repositórios da pausa e do CloseReason          | P0       | M          | T1            |
| T3             | `go-backend-specialist`         | Contagem de falhas desde `resumed_at`                     | P0       | S          | T2            |
| T4             | `go-backend-specialist`         | Service + handler: motivo da pausa                        | P1       | S          | T2            |
| T5             | `code-test`                     | Testes da etapa 1 (backend)                               | P1       | M          | T3, T4        |
| T6             | `frontend-specialist`           | Modal de pausa + selos com motivo                         | P2       | M          | T4            |
| T7             | —                               | Commit da etapa 1                                         | P1       | XS         | T5, T6        |
| **Etapa 2**    |                                 |                                                           |          |            |               |
| T8             | `design-db`                     | Migração 0009 `maintenance_windows`                       | P0       | S          | T7            |
| T9             | `go-backend-specialist`         | Domínio + `MaintenanceRepo`                               | P0       | M          | T8            |
| T10            | `design-api`                    | Contrato da API de manutenção                             | P1       | S          | T9            |
| T11            | `go-backend-specialist`         | `MaintenanceService` (validação, lote, cancelar)          | P1       | M          | T9, T10       |
| T12            | `go-backend-specialist`         | Gate de manutenção no scheduler e no ping                 | P0       | M          | T9            |
| T13            | `go-backend-specialist`         | Tick de início/fim da janela                              | P0       | L          | T11, T12      |
| T14            | `go-backend-specialist`         | Handlers + wiring + bloco `maintenance` nos monitores     | P1       | M          | T11           |
| T15            | `code-test`                     | Testes da etapa 2 (backend)                               | P1       | L          | T13, T14      |
| T16            | `frontend-specialist`           | `zonedToUtc` + status MAINTENANCE em `lib/`               | P2       | S          | T14           |
| T17            | `frontend-specialist`           | Modal de nova janela (multi-seleção)                      | P2       | L          | T16           |
| T18            | `frontend-specialist`           | Selos, card de detalhes, cancelar, selo do incidente      | P2       | M          | T16           |
| **Fechamento** |                                 |                                                           |          |            |               |
| T19            | `qa-review` + `security-review` | Revisão de bordas e de entrada                            | P3       | M          | T15, T17, T18 |
| T20            | `document`                      | Docs e README                                             | P3       | S          | T19           |
| T21            | —                               | Verificação final, `graphify update .`, PR para `develop` | P3       | S          | T20           |

---

### Task T0: Criar branch

- **Agent**: —
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - `develop` limpo e atualizado (`git pull origin develop`)
- **INPUT**: branch `develop`
- **OUTPUT**:
  1. `git switch -c feature/maintenance-windows`
- **VERIFY**:
  - `git branch --show-current` → `feature/maintenance-windows`

### Task T1: Migração 0008 (pausa + close_note)

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T0
- **PRE-CONDITIONS**:
  - Última migração é `0007_monitor_cert_expiry.sql`
- **INPUT**: `monitors` sem campos de pausa; `incidents.closed_manually` INTEGER 0/1
- **OUTPUT**:
  1. `0008_pause_reason.sql` com Up/Down:
     - `monitors.pause_reason TEXT`, `monitors.paused_at TEXT`, `monitors.resumed_at TEXT` (RFC 3339 UTC, como o resto)
     - `incidents.close_note TEXT` (cópia do motivo no momento do fechamento; o motivo da pausa pode mudar depois e o incidente guarda o que valia)
  2. Comentário no topo explicando que `closed_manually` passa a ser 0/1/2 (D9)
- **VERIFY**:
  - `go test ./internal/database/...` (o `migrations_test.go` sobe e desce todas)
- **Rollback**: Down da migração remove as colunas.

### Task T2: Domínio + repositórios da pausa e do CloseReason

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 aplicada
- **INPUT**: `domain.Monitor`, `domain.Incident.ClosedManually bool`, `IncidentRepo.CloseManually(ctx, monitorID, endedAt)`
- **OUTPUT**:
  1. `domain.Monitor`: `PauseReason *string`, `PausedAt *time.Time`, `ResumedAt *time.Time`
  2. `domain.Incident`: `ClosedManually bool` → `CloseReason` (`CloseNone=0`, `ClosePause=1`, `CloseMaintenance=2`) + `CloseNote *string`. Por quê: dois motivos de fechamento com selos diferentes.
  3. `IncidentRepository.CloseManually(ctx, monitorID, endedAt, reason CloseReason, note *string)`
  4. `MonitorRepo.Update` e o scan gravam e leem as três colunas novas; `IncidentRepo` lê `close_note` e o enum (inclusive em `incidentsListSQL`)
  5. Ajustar todos os chamadores (grep `ClosedManually`, `CloseManually`, inclusive `scripts/seed`)
- **VERIFY**:
  - `go build ./... && go vet ./...`
  - `go test ./internal/repository/... ./internal/service/...` sem regressão

### Task T3: Contagem de falhas desde `resumed_at`

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - `Monitor.ResumedAt` disponível
- **INPUT**: `CheckRepoTx.Latest(ctx, tx, monitorID, n)` sem filtro de data
- **OUTPUT**:
  1. `Latest(…, n, since *time.Time)`: com `since != nil`, adiciona `AND checked_at > ?`
  2. `countConsecutiveFails` passa `m.ResumedAt`
  3. Manter o SQL em `queries.go` para o `queryplan_test.go` continuar cobrindo o plano (o índice `(monitor_id, checked_at)` da 0006 atende)
- **VERIFY**:
  - `go test ./internal/repository/ -run QueryPlan`
  - Teste novo em `state_machine_test.go` (T5)
- **Why**: corrige o bug atual da retomada, que vale para a pausa comum e será reusado no fim da janela.

### Task T4: Service + handler do motivo da pausa

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - T2 concluída
- **INPUT**: `MonitorService.Update` já detecta `wasEnabled` → `Enabled`
- **OUTPUT**:
  1. `UpdateMonitorInput.PauseReason *string`; validação em `validate.go`: `strings.TrimSpace`, máx. 200 runes, vazio vira nil
  2. Ao pausar (`true→false`): grava `PausedAt=now`, `PauseReason`; chama `CloseManually(…, ClosePause, PauseReason)`
  3. Ao retomar (`false→true`): grava `ResumedAt=now`, limpa `PauseReason` e `PausedAt`
  4. Handler: `PATCH /api/monitors/{id}` aceita `pause_reason`; resposta do monitor inclui `pause_reason`, `paused_at`, `resumed_at`
  5. Resposta de incidente: `close_reason` (`"pause" | "maintenance" | null`) e `close_note`. Manter `closed_manually` booleano na resposta por compatibilidade (`close_reason != null`)
- **VERIFY**:
  - `go test ./internal/handler/... ./internal/service/...`
  - Motivo com 201 caracteres → 400 `VALIDATION_ERROR`

### Task T5: Testes da etapa 1 (backend)

- **Agent**: `code-test`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T3, T4
- **PRE-CONDITIONS**:
  - T3 e T4 compilando
- **INPUT**: `monitor_pause_test.go`, `state_machine_test.go`, `incidents_test.go` existentes
- **OUTPUT**:
  1. Retomada: threshold 3, 2 falhas antes da pausa, 1 depois → sem incidente; 3 depois → incidente
  2. Pausa grava motivo, `paused_at`; retomar limpa motivo e grava `resumed_at`
  3. Incidente fechado pela pausa tem `close_reason=pause` e `close_note` igual ao motivo
  4. Motivo acima de 200 caracteres é recusado
- **VERIFY**:
  - `go test ./...` verde

### Task T6: Modal de pausa + selos com motivo

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T4
- **PRE-CONDITIONS**:
  - API aceita `pause_reason`
- **INPUT**: pausa é clique direto em `Dashboard.jsx:118` e `MonitorDetail.jsx:385`
- **OUTPUT**:
  1. `PauseModal` em `MonitorModals.jsx`: textarea opcional com `maxLength=200` e contador; Enter confirma; foco inicial no campo; `aria-describedby` no contador
  2. Pausar (RowMenu, MonitorCard, MonitorDetail) abre o modal; retomar continua direto
  3. Motivo visível junto do selo "Pausado" na linha, no card (truncado com `title`) e no detalhe (inteiro)
  4. `ResolvedTable`: selo por `close_reason` com `close_note` abaixo (o selo "Fechado por manutenção" já fica pronto para a etapa 2)
- **VERIFY**:
  - `cd web && npm run lint && npm test && npm run build`
  - Manual: pausar com e sem motivo, ver na lista e no detalhe, retomar

### Task T7: Commit da etapa 1

- **Agent**: —
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T5, T6
- **PRE-CONDITIONS**:
  - `go test ./...` e build do web verdes
- **INPUT**: mudanças da etapa 1
- **OUTPUT**:
  1. `gofumpt -w .`, `graphify update .`
  2. Commit `feat(monitors): motivo da pausa e contagem de falhas desde a retomada`, terminando com `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`
- **VERIFY**:
  - `git log -1` mostra o commit; `git status` limpo

### Task T8: Migração 0009 `maintenance_windows`

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T7
- **PRE-CONDITIONS**:
  - 0008 aplicada
- **INPUT**: nenhuma tabela de janela
- **OUTPUT**:
  1. Tabela:
     ```sql
     CREATE TABLE maintenance_windows (
       id           TEXT PRIMARY KEY,           -- UUIDv7
       monitor_id   TEXT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
       starts_at    TEXT NOT NULL,              -- RFC 3339 UTC
       ends_at      TEXT NOT NULL,
       reason       TEXT,
       started_at   TEXT,                       -- efeitos de início aplicados (tick)
       finalized_at TEXT,                       -- efeitos de fim aplicados (tick ou cancelamento)
       created_at   TEXT NOT NULL,
       CHECK (ends_at > starts_at)
     );
     ```
  2. Índice parcial único `ON maintenance_windows(monitor_id) WHERE finalized_at IS NULL`: garante D2 no banco, além da validação no service
  3. Índice `(monitor_id, starts_at)` para a consulta da janela ativa
  4. Down remove a tabela
- **VERIFY**:
  - `go test ./internal/database/...`
  - Confirmar que `PRAGMA foreign_keys` está ligado em `database.go` (senão o `ON DELETE CASCADE` não roda; nesse caso apagar as janelas no `MonitorService.Delete`)
- **Why**: cancelar antes do início apaga a linha; cancelar durante seta `ends_at=now` e `finalized_at`. Por isso não há coluna `canceled_at`.

### Task T9: Domínio + `MaintenanceRepo`

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T8
- **PRE-CONDITIONS**:
  - Tabela criada
- **INPUT**: padrão de `incident_repo.go` e `queries.go`
- **OUTPUT**:
  1. `domain.MaintenanceWindow` + `func (w *MaintenanceWindow) Active(now time.Time) bool`
  2. `domain.MaintenanceRepository`:
     - `CreateBatch(ctx, []*MaintenanceWindow) error` numa transação
     - `FindLiveByMonitor(ctx, monitorID)` (com `finalized_at IS NULL`)
     - `FindLiveByMonitors(ctx, ids)` em lote, como `FindOpenByMonitors`
     - `ActiveForMonitor(ctx, monitorID, now) (bool, error)`: consulta quente, vai para `queries.go`
     - `FindByID`, `Delete`, `EndNow(ctx, id, now)`
     - `DueStart(ctx, now)` (`starts_at <= now AND started_at IS NULL AND finalized_at IS NULL`)
     - `DueEnd(ctx, now)` (`ends_at <= now AND finalized_at IS NULL`)
     - `MarkStarted`, `MarkFinalized`
- **VERIFY**:
  - `go test ./internal/repository/...`, com `ActiveForMonitor` incluído no `queryplan_test.go` (usa índice, sem SCAN)

### Task T10: Contrato da API de manutenção

- **Agent**: `design-api`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T9
- **PRE-CONDITIONS**:
  - Formato de erro atual (`writeError(…, code)`) conhecido
- **INPUT**: rotas existentes em `server.go`
- **OUTPUT**:
  1. `POST /api/maintenance-windows`
     - corpo: `{monitor_ids: [..], starts_at: RFC3339, ends_at: RFC3339, reason?: string}`
     - 201 → `{windows: [...]}`
     - 400 `VALIDATION_ERROR` com mensagem específica (início no passado, fim antes do início, mais de 7 dias, motivo longo, lista vazia)
     - 404 se algum monitor não existe
     - 409 `MAINTENANCE_CONFLICT` com `monitor_ids` em conflito (D5)
  2. `DELETE /api/maintenance-windows/{id}` → 204 (antes do início: apaga; durante: encerra)
  3. `GET /api/monitors` e `GET /api/monitors/{id}` ganham `maintenance: {id, starts_at, ends_at, reason, active} | null`. `active` é calculado no servidor, para não depender do relógio do navegador
- **VERIFY**:
  - Contrato registrado na seção API do README (T20) e revisado antes de T14

### Task T11: `MaintenanceService`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T9, T10
- **PRE-CONDITIONS**:
  - Repositório pronto
- **INPUT**: padrão de `MonitorService` (`ErrValidation`, `ErrNotFound`, `reload`)
- **OUTPUT**:
  1. `CreateBatch(ctx, in)`:
     - deduplica os IDs
     - valida início `> now`, fim `> início`, duração `<= 7*24h`, motivo `<= 200` runes e trim
     - todos os monitores existem
     - nenhum tem janela viva → senão `ErrMaintenanceConflict{MonitorIDs}`
     - insere tudo numa transação
  2. `Cancel(ctx, id)`:
     - antes do início: `Delete`
     - durante: `EndNow` + aplica o fim na hora (mesma função do tick, T13)
     - já finalizada: `ErrNotFound`
  3. `Active(ctx, monitorID, now)`: helper usado pelo gate
  4. Relógio injetável (`clock func() time.Time`), no padrão da `StateMachine.WithClock`
- **VERIFY**:
  - Testes de validação em T15

### Task T12: Gate de manutenção no scheduler e no ping

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T9
- **PRE-CONDITIONS**:
  - `ActiveForMonitor` disponível
- **INPUT**: `withMonitor` descarta com `!m.Enabled`; `MonitorService.Ping` idem
- **OUTPUT**:
  1. `Scheduler` ganha `maintenance domain.MaintenanceRepository` e `now func() time.Time` (padrão `time.Now`, trocável nos testes)
  2. `withMonitor`: dentro do `dbMu`, depois do `FindByID`, se `ActiveForMonitor` → `return nil`. É o ponto único: cobre check HTTP/TCP, heartbeat vencido e ping, e por consequência incidente e todos os avisos, inclusive SSL (`processCert` roda dentro de `Process`)
  3. `runLoop`: no tick, se em manutenção, pula o `runCheck` sem sair do loop (evita requisição de rede inútil)
  4. `checkHeartbeats`: pula monitor em manutenção antes de `heartbeatExpired`
  5. `MonitorService.Ping`: em manutenção → `nil` (200, sem registrar)
- **VERIFY**:
  - Teste: check em andamento quando a janela começa é descartado
  - Teste: ping durante a janela → 200 e nenhum check gravado
- **Why**: o bloqueio vem do relógio e é exato; só os efeitos laterais (T13) dependem do tick.

### Task T13: Tick de início/fim da janela

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: L
- **Dependencies**: T11, T12
- **PRE-CONDITIONS**:
  - Gate funcionando
- **INPUT**: `runHeartbeats` como molde (um ticker, `wg`, cancelamento no `Stop`)
- **OUTPUT**:
  1. `internal/scheduler/maintenance.go`: `runMaintenance(ctx)` com ticker (`maintenanceTick`, padrão 10 s, trocável em teste), iniciado em `Start`
  2. **Início** (`DueStart`), por janela, sob `dbMu`, numa transação:
     - `CloseManually(…, CloseMaintenance, reason)`
     - `current_state = UNKNOWN`
     - `MarkStarted`
  3. **Fim** (`DueEnd`), por janela:
     - `MarkFinalized`
     - `resumed_at = now`
     - heartbeat: `LastPingAt = now` se não nil, senão `UpdatedAt = now` (base do prazo em `heartbeatExpired`)
     - `current_state = UNKNOWN`
     - depois do commit, `Reload(monitorID)`, que roda o primeiro check na hora e respeita `Enabled=false` (D3 de graça)
  4. Janela que começou e terminou com o Pingou desligado: o fim só roda depois do início (ordem início → fim no mesmo tick)
  5. Rodar uma vez no `Start` antes do primeiro tick (reinício no meio da janela aplica o pendente)
  6. Expor a função de fim para o `Cancel` (T11)
  7. Log `slog.Info("maintenance started|ended", "monitor_id", …, "window_id", …)`
- **VERIFY**:
  - Testes de T15 com relógio injetado e tick curto
  - `shutdown_test.go` continua verde (o loop novo entra no `wg`)

### Task T14: Handlers + wiring + bloco `maintenance`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T11
- **PRE-CONDITIONS**:
  - Contrato T10 fechado
- **INPUT**: `server.go`, `monitors.go`, `serve.go`
- **OUTPUT**:
  1. `handler/maintenance.go` com POST/DELETE; corpo limitado com `http.MaxBytesReader`; mapeamento de erros: `ErrValidation`→400, `ErrNotFound`→404, conflito→409
  2. Rotas em `server.go`, dentro do `authMiddleware`
  3. `GET /monitors` e `GET /monitors/{id}`: uma consulta em lote `FindLiveByMonitors` para a página (sem N+1)
  4. `serve.go`: cria repo e service, injeta no `Scheduler` e no `MonitorService`
- **VERIFY**:
  - `go test ./internal/handler/...`
  - `cmd/pingou/commands/serve_test.go` verde

### Task T15: Testes da etapa 2 (backend)

- **Agent**: `code-test`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: T13, T14
- **PRE-CONDITIONS**:
  - Relógio e ticks injetáveis
- **INPUT**: `scheduler_test.go`, `heartbeat_test.go`, `monitor_pause_test.go` como molde
- **OUTPUT**:
  1. Validação: início no passado, fim antes do início, 7 dias + 1 s, motivo longo, monitor inexistente, conflito em lote (nada criado)
  2. Durante a janela: zero checks gravados, zero notificações (notifier fake conta as chamadas, inclusive `NotifyCertExpiring`)
  3. Incidente aberto fecha com `close_reason=maintenance` e motivo
  4. Fim: estado `UNKNOWN`, primeiro check imediato, falhas antigas não contam
  5. Heartbeat: prazo recomeça no fim da janela (sem DOWN imediato)
  6. D3: monitor pausado durante a janela segue pausado depois
  7. Cancelar antes do início apaga; cancelar durante encerra e volta
  8. Reinício: `Start` com janela em curso mantém manutenção; janela vencida com o Pingou fora é finalizada
  9. Excluir monitor com janela viva não quebra nada
- **VERIFY**:
  - `go test -race ./...` verde

### Task T16: `zonedToUtc` + status MAINTENANCE

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T14
- **PRE-CONDITIONS**:
  - API devolve `maintenance`
- **INPUT**: `lib/timezone.js` só converte UTC → fuso (`zonedParts`, `wallMs`); `lib/monitors.js` com `PAUSED` vencendo o estado
- **OUTPUT**:
  1. `zonedToUtc("2026-11-10T02:00", tz) → Date`: chute inicial e correção pelo offset de `zonedParts`, cobrindo DST. Testes em `timezone.test.js` com `America/Sao_Paulo`, `Asia/Kolkata` e `UTC`
  2. `monitorStatus`: `PAUSED` > `MAINTENANCE` (se `maintenance?.active`) > `current_state`
  3. Severidade e contadores com `MAINTENANCE`; filtro "Em manutenção" no Dashboard; rótulo "Em manutenção até HH:MM" no fuso das Configurações (com data quando não é hoje)
- **VERIFY**:
  - `cd web && npm test`

### Task T17: Modal de nova janela

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T16
- **PRE-CONDITIONS**:
  - `api/maintenance.js` com `create` e `cancel`
- **INPUT**: `Modal`, `Combobox`, `MonitorForm` (padrão de erros)
- **OUTPUT**:
  1. `MaintenanceModal`:
     - multi-seleção de monitores (pré-selecionado quando aberto do detalhe); pausados (`!enabled`) e com janela viva (`maintenance`) ficam desabilitados com selo "pausado" / "já agendada" (D11)
     - início e fim com `datetime-local`, com o nome do fuso das Configurações no rótulo (`getTimeZone()`, ex.: `America/Sao_Paulo`) (D10)
     - motivo opcional com contador
  2. Validação no cliente espelhando o servidor (início futuro, fim depois, 7 dias); erro 409 lista os monitores em conflito pelo nome
  3. Entradas: botão "Agendar manutenção" no Dashboard e no detalhe (menu)
  4. Acessibilidade: labels, erros com `aria-invalid` e `aria-describedby`, foco no primeiro erro
- **VERIFY**:
  - `npm run lint && npm run build`
  - Manual: criar janela para 3 monitores com fuso das Configurações ≠ fuso do navegador e conferir o UTC gravado

### Task T18: Selos, card de detalhes, cancelar

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T16
- **PRE-CONDITIONS**:
  - Status `MAINTENANCE` disponível
- **INPUT**: `MonitorRow`, `MonitorCard`, `StatusBanner`, `MonitorDetail`, `ResolvedTable`
- **OUTPUT**:
  1. Lista e card: selo "Em manutenção até HH:MM" (tom neutro, distinto de "Pausado"); `HistoryBar` esmaecida
  2. `StatusBanner`: conta "em manutenção" e não trata o monitor como fora do ar
  3. `MaintenanceCard` no detalhe: motivo, início, fim, estado (agendada / em curso) e botão "Cancelar" com confirmação
  4. `ResolvedTable`: "Fechado por manutenção" com motivo (selo preparado em T6)
- **VERIFY**:
  - `npm run lint && npm test && npm run build`
  - Manual: janela curta (2 min) passando de agendada → em curso → encerrada sem recarregar a página

### Task T19: Revisão de bordas e segurança

- **Agent**: `qa-review` + `security-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T15, T17, T18
- **PRE-CONDITIONS**:
  - Feature completa na branch
- **INPUT**: branch inteira
- **OUTPUT**:
  1. QA: pausar durante a janela, retomar durante a janela, cancelar no segundo exato do início, excluir monitor com janela, editar URL durante a janela, fuso trocado nas Configurações com janela agendada
  2. Segurança: rotas novas atrás do `authMiddleware`; corpo limitado; motivo renderizado como texto (sem `dangerouslySetInnerHTML`); IDs validados; sem SQL concatenado com entrada do usuário
- **VERIFY**:
  - Achados corrigidos ou registrados em Notes

### Task T20: Docs

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T19
- **PRE-CONDITIONS**:
  - Comportamento final estável
- **INPUT**: `docs/implementacoes.md`, `README.md`, `.specs/project/ROADMAP.md`
- **OUTPUT**:
  1. `docs/implementacoes.md`: seção da feature, decisões D1–D9, nota de que os baldes do `/stats` seguem pendentes
  2. README: endpoints novos e campos novos de monitor e incidente
  3. ROADMAP: marcar "Maintenance window support"
- **VERIFY**:
  - Revisão visual dos arquivos

### Task T21: Verificação final e PR

- **Agent**: —
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T20
- **PRE-CONDITIONS**:
  - Tudo verde localmente
- **INPUT**: branch com etapas 1 e 2
- **OUTPUT**:
  1. Rodar a Phase X inteira
  2. `graphify update .`
  3. Commit da etapa 2: `feat(maintenance): janela de manutenção agendada por monitor`, com a linha `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`
  4. `git push -u origin feature/maintenance-windows`
  5. `gh pr create --base develop`, com:
     - resumo das duas etapas
     - decisões D1–D9
     - checklist de teste manual
     - nota sobre a correção do bug da retomada
     - última linha: `🤖 Generated with [Claude Code](https://claude.com/claude-code)`
- **VERIFY**:
  - `gh pr view --web` abre o PR contra `develop`; CI verde

## Risk Areas

| Risk                                                 | Impact                                                 | Mitigation                                                                                             |
| ---------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| Fuso: `datetime-local` usa o fuso do navegador       | Janela deslocada em horas                              | `zonedToUtc` com o fuso das Configurações + testes com fusos de meia hora e DST; fuso visível no modal |
| Consulta extra por check no `withMonitor`            | Latência no caminho quente                             | Consulta indexada coberta pelo `queryplan_test`; 100 monitores no máximo                               |
| Corrida entre o tick e um `Update`/`Cancel`          | Estado gravado por cima                                | Tudo sob `dbMu`; o tick recarrega o monitor dentro do lock, como `withMonitor`                         |
| Efeitos de início/fim atrasam até um tick            | Incidente fecha ou primeiro check roda até 10 s depois | Aceitável; o bloqueio em si é exato (relógio). Documentado                                             |
| `foreign_keys` desligado no SQLite                   | Janelas órfãs ao excluir monitor                       | Verificar em T8; apagar no service se preciso                                                          |
| Troca de `ClosedManually` por enum quebra chamadores | Build quebrado                                         | Grep completo em T2; manter `closed_manually` na resposta JSON                                         |
| Teste de scheduler instável por tempo real           | CI intermitente                                        | Relógio e ticks injetáveis; sem `time.Sleep` longo                                                     |

## Rollback Strategy

1. **Antes do merge:** fechar o PR e apagar a branch; `develop` não muda.
2. **Depois do merge, só etapa 2:** `git revert` do commit da etapa 2. A 0009 tem Down (`pingou migrate down` uma vez) e a tabela nova não afeta as antigas.
3. **Depois do merge, tudo:** reverter os dois commits e descer as migrações 0009 e 0008. Incidentes com `closed_manually=2` continuam lidos como "fechado" (valor ≠ 0).
4. Banco do usuário: recomendar `GET /api/export` antes de atualizar (já existe).

## Phase X: Verification

```bash
# Backend
gofumpt -l .            # deve sair vazio
go vet ./...
go test -race ./...

# Frontend
cd web && npm ci && npm run lint && npm run format:check && npm test && npm run build && cd ..

# Binário completo (web embutido)
make build

# Grafo
graphify update .
```

Teste manual mínimo (registrar no PR):

- Pausar com motivo → selo e motivo na lista e no detalhe → retomar → falha única não abre incidente (threshold 3)
- Agendar janela de 2 min em 2 monitores (um DOWN com incidente aberto) → ao começar, o incidente fecha com "Fechado por manutenção" e nenhum aviso sai → ao terminar, o primeiro check roda e o estado se resolve
- Cancelar uma janela agendada e uma em curso
- Heartbeat em manutenção: `curl /ping/<token>` → 200, sem check novo

## Cronograma Sugerido

| Fase                         | Tarefas | Complexity | Esforço | Target                                       | Status      |
| ---------------------------- | ------- | ---------- | ------- | -------------------------------------------- | ----------- |
| Fase 0                       | T0      | XS         | 5 min   | Branch criada                                | ✅ Concluído |
| Fase 1 (etapa 1)             | T1–T7   | M          | ~3 h    | Pausa com motivo + bug da retomada corrigido | ✅ Concluído |
| Fase 2 (backend da etapa 2)  | T8–T15  | XL         | ~6 h    | Janela funcionando via API, com testes       | ✅ Concluído |
| Fase 3 (frontend da etapa 2) | T16–T18 | L          | ~3–4 h  | UI completa                                  | ✅ Concluído |
| Fase 4 (fechamento)          | T19–T21 | M          | ~1,5 h  | PR aberto contra `develop`                   | ✅ Concluído |

## Status de Execução

### Concluído

- Análise de viabilidade e decisões D1–D9 aprovadas
- T0–T21: etapas 1 e 2, testes, QA manual, docs e PR (desvios registrados em `docs/implementacoes.md`, seção 14)

### Em andamento

-

### Pendente

- QA em dispositivo real / leitor de tela (registrado em Pendências abertas)

## Notes

- O gate único em `Scheduler.withMonitor` é o coração da feature: qualquer resultado de check novo no futuro passa por ali e herda o bloqueio.
- `resumed_at` serve à pausa e ao fim da janela; não há campo separado para "fim da manutenção".
- `close_note` copia o motivo no fechamento para o histórico não mudar se o motivo da pausa for alterado depois.
- Baldes do `/stats` em fusos de meia hora seguem fora do escopo (registrado em `docs/implementacoes.md`).

### Alternativas consideradas e descartadas

- `Flag is_in_maintenance em monitors` → descartada por exigir que o tick esteja sempre certo; derivar do relógio (D6) é exato e sobrevive a reinício.
- `Um timer por janela (time.AfterFunc)` → descartada: some no reinício e precisa de reconstrução; o tick único é o padrão do `runHeartbeats`.
- `Coluna canceled_at` → descartada: cancelar antes apaga, cancelar durante encerra (`ends_at=now`); o histórico não pede mais que isso.
- `Biblioteca de date picker` → descartada: `datetime-local` nativo + helper de fuso resolve.
- `Criação em lote parcial` → descartada por D5 (tudo ou nada é mais previsível).
- `Coluna nova close_reason TEXT` → descartada: `closed_manually` INTEGER sem CHECK já comporta 0/1/2 (D9).
