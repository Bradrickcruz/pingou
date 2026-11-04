# Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)

## Overview

**O quê.** Uma placa de "em manutenção" num monitor. Enquanto ela está lá, o Pingou não confere nada, não abre incidente e não avisa ninguém. Quando sai, o monitor volta ao normal.

**Duas etapas, um PR:**

1. **Pausa com motivo.** Motivo opcional (até 200 caracteres), instante da pausa e da retomada guardados. A contagem de falhas passa a olhar só checks feitos depois da retomada. Isso corrige um bug real: hoje falhas de antes da pausa somam com as de depois (`countConsecutiveFails` em `internal/service/state_machine.go` lê os últimos N checks sem filtro de data).
2. **Janela agendada.** Início, fim e motivo opcional. Janela única (sem recorrência), por monitor, com criação em lote para vários monitores.
3. **Página de manutenções (Etapa 3, nova).** Uma página `/maintenance` que lista as janelas vivas (agendadas e em curso) de todos os monitores, permite **editar** reabrindo o modal e **desligar** (cancelar) qualquer uma. Janela agendada tem tudo editável; janela já iniciada tem o **horário travado** (só o motivo muda) e pode ser desligada. A edição também fica disponível no card de manutenção do detalhe do monitor.

**Decisões fechadas (aprovadas na análise de viabilidade):**

| #   | Decisão                                                                                                                                                                                                                                                                                                                                                                                                                      |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| D1  | Teto de **7 dias** por janela. Início no futuro, fim depois do início.                                                                                                                                                                                                                                                                                                                                                       |
| D2  | **Uma janela viva por monitor** (agendada ou em curso). "Próxima janela marcada" é essa mesma janela.                                                                                                                                                                                                                                                                                                                        |
| D3  | **Pausa manual vence a janela.** Monitor pausado continua pausado quando a janela termina.                                                                                                                                                                                                                                                                                                                                   |
| D4  | Motivo **opcional** na pausa e na janela, máx. 200 caracteres.                                                                                                                                                                                                                                                                                                                                                               |
| D5  | Criação em lote **tudo ou nada**: se um monitor já tem janela viva, nada é criado e o erro diz qual.                                                                                                                                                                                                                                                                                                                         |
| D6  | **"Em manutenção" é derivado do relógio** (`starts_at <= agora < ends_at`, não cancelada), não de flag. Sobrevive a reinício.                                                                                                                                                                                                                                                                                                |
| D7  | Pingou roda com **uma instância** (`dbMu` + SQLite `MaxOpenConns=1`). Sem lock distribuído.                                                                                                                                                                                                                                                                                                                                  |
| D8  | Uptime de 24 h não muda: sem checks na janela, não há o que contar.                                                                                                                                                                                                                                                                                                                                                          |
| D9  | `closed_manually` (INTEGER sem CHECK) passa a guardar 0/1/2 = nenhum/pausa/manutenção. Sem migração de esquema para isso.                                                                                                                                                                                                                                                                                                    |
| D10 | No modal de agendar, o rótulo de Início e Fim mostra o **nome do fuso** das Configurações (ex.: `America/Sao_Paulo`), não o offset (`-03`).                                                                                                                                                                                                                                                                                  |
| D11 | No modal de agendar, monitores **pausados** ou com **janela viva** (agendada ou em curso) aparecem na lista **desabilitados** (selo "pausado" / "já agendada"), pois não podem entrar numa janela.                                                                                                                                                                                                                           |
| D12 | **Janela iniciada tem o horário travado; só o motivo muda.** "Iniciada" é decidido pelo relógio do servidor (`starts_at <= agora`, igual a D6), não pelo tick (`started_at`). Iniciada: `starts_at` e `ends_at` não editáveis, só `reason`; para parar antes, **desligar** (cancelar, que encerra na hora). Agendada: início, fim e motivo editáveis.                                                                        |
| D13 | **Editar não troca o monitor.** Cada janela pertence a um monitor; editar mexe só em `starts_at`, `ends_at` e `reason` dela. Para mudar o monitor: desligar e agendar outra. A edição revalida D1 (início no futuro, fim depois do início, teto de 7 dias) e D4. D2 não é afetada (a janela já é a única viva do monitor).                                                                                                   |
| D14 | **A página lista só janelas vivas** (agendadas e em curso), ordenadas por início. Janelas encerradas ou canceladas não aparecem (não há histórico: cancelar antes do início apaga a linha, ver T8). Máx. 100 linhas (uma por monitor, limite de 100 monitores), sem paginação.                                                                                                                                               |
| D15 | **O modal é reaproveitado** em dois modos: criar (como hoje) e editar (monitor fixo, campos preenchidos, fuso das Configurações). Em janela iniciada, os campos de início e fim ficam desabilitados com a explicação "A manutenção já começou: só o motivo pode mudar". Corrida: se a janela começar enquanto o modal está aberto com horário alterado, o servidor recusa com 409 `MAINTENANCE_LOCKED` e a página recarrega. |
| D16 | **Etapa 3 em branch e PR próprios, empilhados.** O PR das etapas 1 e 2 ainda não foi mergeado: `feature/maintenance-page` parte de `feature/maintenance-windows` e o PR da Etapa 3 abre com base nela. Quando o PR das etapas 1 e 2 for mergeado, rebasear a branch em `develop` e trocar a base do PR para `develop`.                                                                                                       |
| D17 | **Editar também no card de manutenção do detalhe do monitor** (`MaintenanceCard`): botão "Editar" abre o mesmo modal em modo editar (D15), com as mesmas regras.                                                                                                                                                                                                                                                             |

**Fora do escopo:** janela recorrente, aviso de "começou/terminou", página pública de status, limpeza de incidentes órfãos, histórico de janelas encerradas, estender ou encurtar o horário de janela em curso, trocar o monitor de uma janela, revisão dos baldes do `/stats` em UTC.

**Fluxo de git:** branch `feature/maintenance-windows` a partir de `develop`, commits por etapa, PR para `develop`. Etapa 3: branch `feature/maintenance-page` a partir de `feature/maintenance-windows`, commits por camada (backend, frontend), PR empilhado com base em `feature/maintenance-windows` (D16).

## Project Type

**FULLSTACK** — backend Go (SQLite + goose, scheduler próprio, state machine) e SPA React/Vite embutida no binário.

## Definition of Done (DoD)

Este plano está concluído quando (reaberto pela Etapa 3; as etapas 1 e 2 já cumpriram estes itens):

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas
- [ ] PR aberto contra `develop` com descrição e checklist de teste manual (um por etapa: etapas 1–2 já abertos; Etapa 3 pendente)

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

**Etapa 3**

- [ ] `GET /api/maintenance-windows` devolve as janelas vivas com nome do monitor, estado (`scheduled` | `active`), início, fim e motivo, ordenadas por início (D14)
- [ ] `PATCH /api/maintenance-windows/{id}` edita início, fim e motivo de janela **agendada**, revalidando D1 e D4 (D12, D13)
- [ ] PATCH em janela **iniciada** aceita só `reason`; com `starts_at` ou `ends_at` → 409 `MAINTENANCE_LOCKED`, nada é alterado; janela inexistente ou já finalizada → 404 (D12)
- [ ] Mudar o motivo de janela em curso não reaplica nenhum efeito de início/fim (incidente, estado, `resumed_at` ficam como estão)
- [ ] Página `/maintenance` acessível pelo menu (desktop e `BottomNav`), com estados de carregando, vazio e erro
- [ ] A lista mostra monitor, estado, início e fim no fuso das Configurações, e motivo
- [ ] Toda janela tem "Editar" (abre o modal preenchido, monitor fixo). Agendada: início, fim e motivo editáveis, e botão "Cancelar". Em curso: início e fim desabilitados no modal, motivo editável, e botão "Desligar" com confirmação (D12, D15)
- [ ] Salvar a edição atualiza a linha sem recarregar a página; janela que começa com o modal aberto e horário alterado mostra o erro do 409 e a lista se atualiza (D15)
- [ ] Linha passa de "Agendada" para "Em curso" sozinha, sem recarregar a página (o botão "Cancelar" vira "Desligar")
- [ ] Botão "Agendar manutenção" na página abre o modal em modo criar
- [ ] O `MaintenanceCard` do detalhe do monitor tem "Editar", que abre o mesmo modal em modo editar (D17)

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

**Etapa 3 (arquivos novos ou alterados)**

```
internal/
  domain/
    maintenance.go                   # + Started(now) (starts_at <= now); MaintenanceRepository ganha ListLive e Update
  repository/
    maintenance_repo.go              # + ListLive (JOIN com monitors.name), UpdateTimes (guarda starts_at > now), UpdateReason
    queries.go                       # + SQL da listagem (entra no queryplan_test)
  service/
    maintenance_service.go           # + List, Update (agendada: tudo; iniciada: só reason, horário → ErrMaintenanceLocked)
    validate.go                      # reaproveita validação da janela (extrair função comum a CreateBatch e Update)
  handler/
    maintenance.go                   # + GET /maintenance-windows, PATCH /maintenance-windows/{id}; 409 MAINTENANCE_LOCKED
    server.go                        # rotas novas

web/src/
  api/maintenance.js                 # + list, update
  lib/maintenance.js                 # + windowState(w, now), validateWindow aceita edição (início futuro só se mudou)
  pages/Maintenance.jsx              # NOVO: página /maintenance
  components/maintenance/            # NOVO
    WindowsTable.jsx                 # tabela (desktop) e cartões (celular)
    CancelWindowDialog.jsx           # confirmação de cancelar / desligar
  components/monitors/MaintenanceModal.jsx  # + modo editar (prop `window`; horário desabilitado se iniciada)
  components/monitor-detail/MaintenanceCard.jsx  # + botão "Editar" (abre o modal em modo editar)
  components/layout/{Shell,BottomNav}.jsx   # + item "Manutenções"
  main.jsx                           # + rota /maintenance
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
| **Etapa 3**    |                                 |                                                           |          |            |               |
| T22            | —                               | Criar branch `feature/maintenance-page`                   | P0       | XS         | T21           |
| T23            | `design-api`                    | Contrato `GET` e `PATCH /maintenance-windows`             | P0       | S          | T22           |
| T24            | `go-backend-specialist`         | Repositório: `ListLive` e `Update`                        | P0       | M          | T23           |
| T25            | `go-backend-specialist`         | `MaintenanceService`: `List` e `Update` (regra D12/D13)   | P1       | M          | T24           |
| T26            | `go-backend-specialist`         | Handlers `GET` e `PATCH` + erro 409 `MAINTENANCE_LOCKED`  | P1       | M          | T25           |
| T27            | `code-test`                     | Testes da Etapa 3 (backend)                               | P1       | M          | T26           |
| T28            | `design-ux`                     | Fluxo e estados da página `/maintenance`                  | P2       | S          | T23           |
| T29            | `frontend-specialist`           | API client, rota, menu e `lib/maintenance.js`             | P2       | S          | T26, T28      |
| T30            | `frontend-specialist`           | Página `Maintenance.jsx`: lista, cancelar, atualização    | P2       | L          | T29           |
| T31            | `frontend-specialist`           | `MaintenanceModal` em modo editar                         | P2       | M          | T29           |
| T32            | `frontend-specialist`           | "Editar" no `MaintenanceCard` do detalhe                  | P2       | S          | T31           |
| T33            | `qa-review` + `security-review` | Revisão de bordas e de entrada da Etapa 3                 | P3       | M          | T27, T30, T32 |
| T34            | `document`                      | Docs e README da Etapa 3                                  | P3       | S          | T33           |
| T35            | —                               | Verificação final, `graphify update .`, PR da Etapa 3     | P3       | S          | T34           |

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

---

## Etapa 3: Página de manutenções

### Task T22: Criar branch

- **Agent**: —
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T21
- **PRE-CONDITIONS**:
  - `feature/maintenance-windows` com as etapas 1 e 2 commitadas (o PR delas ainda não foi mergeado, D16)
- **INPUT**: branch `feature/maintenance-windows`
- **OUTPUT**:
  1. `git switch feature/maintenance-windows && git switch -c feature/maintenance-page`
- **VERIFY**:
  - `git branch --show-current` → `feature/maintenance-page`
- **Rollback**: apagar a branch; `feature/maintenance-windows` não muda.

### Task T23: Contrato `GET` e `PATCH /maintenance-windows`

- **Agent**: `design-api`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T22
- **PRE-CONDITIONS**:
  - Contrato da T10 (`POST`, `DELETE`, bloco `maintenance`) conhecido
- **INPUT**: só `POST` e `DELETE` existem; a lista de janelas só aparece embutida em cada monitor
- **OUTPUT**:
  1. `GET /api/maintenance-windows` → 200 `{windows: [{id, monitor_id, monitor_name, starts_at, ends_at, reason, state}]}`
     - `state`: `"scheduled"` | `"active"`, calculado no servidor (D12, mesmo relógio de D6)
     - só janelas vivas (`finalized_at IS NULL`), por `starts_at` crescente (D14)
  2. `PATCH /api/maintenance-windows/{id}` → 200 com a janela atualizada
     - corpo: `{starts_at?, ends_at?, reason?}`; campo ausente = não muda; `reason: ""` ou `null` limpa o motivo
     - janela **agendada**: os três campos valem
     - janela **iniciada**: só `reason` vale; corpo com `starts_at` ou `ends_at` → 409 (mesmo que o valor seja igual ao atual, para o contrato ficar simples)
     - 400 `VALIDATION_ERROR`: início no passado, fim antes do início, mais de 7 dias, motivo longo, corpo vazio
     - 404: não existe ou já finalizada
     - 409 `MAINTENANCE_LOCKED`: tentativa de mudar horário de janela iniciada (D12)
  3. `DELETE` não muda (já cobre "desligar" a janela em curso)
- **VERIFY**:
  - Contrato revisado antes da T24 e registrado no README (T34)
- **Why**: o front precisa de uma lista única, sem montar a página com N chamadas a `GET /monitors`, e de um erro distinguível para a corrida do D15.

### Task T24: Repositório: `ListLive` e `Update`

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T23
- **PRE-CONDITIONS**:
  - `MaintenanceRepository` da T9 funcionando
- **INPUT**: `FindLiveByMonitors(ids)` exige os IDs; não há listagem geral nem update de horários
- **OUTPUT**:
  1. `domain.MaintenanceWindow.Started(now) bool` (`!now.Before(StartsAt)`), ao lado de `Active`
  2. `MaintenanceRepository.ListLive(ctx) ([]MaintenanceWindowView, error)`: `JOIN monitors` para o nome, `finalized_at IS NULL`, `ORDER BY starts_at`; SQL em `queries.go`
  3. Dois métodos, porque as regras são diferentes:
     - `UpdateTimes(ctx, id, startsAt, endsAt, reason, now)`: `UPDATE ... WHERE id = ? AND finalized_at IS NULL AND starts_at > ?` (o `now` entra como parâmetro); nenhuma linha mudou → `ErrNotFound`. A guarda no SQL fecha a corrida entre o service e o início da janela
     - `UpdateReason(ctx, id, reason)`: `UPDATE ... WHERE id = ? AND finalized_at IS NULL`; vale para agendada e em curso; nenhuma linha → `ErrNotFound`
  4. `ListLive` incluído no `queryplan_test.go` (sem SCAN indesejado; com 100 linhas máx., SCAN da tabela pequena é aceitável se o teste documentar)
- **VERIFY**:
  - `go test ./internal/repository/...`
  - `CHECK (ends_at > starts_at)` continua valendo (teste: `UpdateTimes` com fim antes do início falha)

### Task T25: `MaintenanceService`: `List` e `Update`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T24
- **PRE-CONDITIONS**:
  - Repositório pronto
- **INPUT**: `CreateBatch` valida inline; `Cancel` já distingue antes/durante
- **OUTPUT**:
  1. `List(ctx)`: chama `ListLive` e preenche `state` com `Started(now)` do relógio injetável
  2. `Update(ctx, id, in UpdateMaintenanceInput)`:
     - `FindByID`; ausente ou finalizada → `ErrNotFound`
     - corpo sem nenhum campo → `ErrValidation`; motivo com trim e `<= 200` runes (D4)
     - `Started(now)` e o corpo traz `starts_at` ou `ends_at` → `ErrMaintenanceLocked` (novo erro; D12)
     - `Started(now)` e só `reason` → `UpdateReason`; **não** toca em incidente, estado do monitor nem `resumed_at` (os efeitos de início/fim são do tick, T13)
     - agendada: monta os valores finais (campo ausente mantém o atual) e valida início `> now`, fim `> início`, duração `<= 7*24h`; grava via `UpdateTimes`; se o repo devolver "nenhuma linha" (a janela começou no meio), reler: iniciada → `ErrMaintenanceLocked` se o corpo trazia horário (mas aplicar `UpdateReason` quando o corpo só tinha `reason`), sumiu → `ErrNotFound`
  3. Extrair a validação de janela para uma função usada por `CreateBatch` e `Update` (uma regra só para D1 e D4)
- **VERIFY**:
  - Testes de validação e de trava em T27

### Task T26: Handlers `GET` e `PATCH` + 409 `MAINTENANCE_LOCKED`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T25
- **PRE-CONDITIONS**:
  - Contrato da T23 fechado
- **INPUT**: `handler/maintenance.go` com `POST` e `DELETE`; `writeMaintenanceError`
- **OUTPUT**:
  1. `handleListMaintenance` e `handleUpdateMaintenance`; corpo do `PATCH` limitado com `http.MaxBytesReader`, como o `POST`
  2. `writeMaintenanceError`: `ErrMaintenanceLocked` → 409 com código `MAINTENANCE_LOCKED`; os demais mapeamentos não mudam
  3. Rotas em `server.go`, dentro do `authMiddleware` e do mesmo `if s.maintenanceService != nil`
  4. `PATCH` distingue "campo ausente" de `reason: null` (ponteiros ou `json.RawMessage`), já que ambos significam coisas diferentes
- **VERIFY**:
  - `go test ./internal/handler/...`
  - `cmd/pingou/commands/serve_test.go` verde

### Task T27: Testes da Etapa 3 (backend)

- **Agent**: `code-test`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T26
- **PRE-CONDITIONS**:
  - Relógio injetável (T11)
- **INPUT**: testes da T15 como molde
- **OUTPUT**:
  1. `List`: janelas vivas ordenadas; `state` correto antes e depois do `starts_at`; finalizadas fora da lista
  2. `Update` em agendada: troca início, fim e motivo; troca só o motivo; limpa o motivo
  3. `Update` recusado: início no passado, fim antes do início, 7 dias + 1 s, motivo longo, corpo vazio
  4. `Update` em janela iniciada com `starts_at` ou `ends_at` → `ErrMaintenanceLocked`, valores intactos (inclusive no segundo exato do início)
     4b. `Update` em janela iniciada só com `reason` → ok; motivo novo gravado; incidente, `current_state` e `resumed_at` do monitor intactos; a janela continua bloqueando o monitor
  5. `Update` em janela inexistente ou finalizada → `ErrNotFound`
  6. Janela editada continua com o monitor em manutenção só no novo intervalo (gate de T12 usa os valores novos)
  7. Handler: `GET` 200, `PATCH` 200/400/404/409 e sem autenticação → 401
- **VERIFY**:
  - `go test -race ./...` verde

### Task T28: Fluxo e estados da página

- **Agent**: `design-ux`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T23
- **PRE-CONDITIONS**:
  - Decisões D12–D15 aprovadas
- **INPUT**: padrão visual de `Incidents.jsx` (abas, tabela, cartões no celular) e `MaintenanceCard`
- **OUTPUT**:
  1. Colunas: Monitor (link para o detalhe), Estado (selo "Agendada" / "Em curso"), Início, Fim, Motivo, Ações
  2. Ações por linha: agendada → "Editar" e "Cancelar"; em curso → "Editar" (só o motivo, horário desabilitado no modal) e "Desligar" (rótulo diferente porque o efeito é diferente: encerra uma manutenção que já está valendo)
  3. Confirmação: texto distinto para cancelar (agendada: "A janela será removida") e desligar (em curso: "O monitor volta a ser checado agora")
  4. Estados: carregando (`Spinner`), vazio ("Nenhuma manutenção agendada" + botão "Agendar manutenção"), erro com "Tentar de novo"
  5. Celular: uma linha vira cartão; ações em menu ou botões de largura toda; foco e `aria-label` com o nome do monitor ("Editar manutenção de <monitor>")
  6. Botão "Agendar manutenção" no topo
- **VERIFY**:
  - Revisão do fluxo antes da T30; checklist de acessibilidade (foco, teclado, `aria-live` na troca de estado) anexado à T33

### Task T29: API client, rota, menu e `lib/maintenance.js`

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T26, T28
- **PRE-CONDITIONS**:
  - `GET` e `PATCH` respondendo
- **INPUT**: `api/maintenance.js` só tem `create` e `cancel`; `main.jsx` com rotas `/`, `/monitors/:id`, `/incidents`, `/settings`
- **OUTPUT**:
  1. `api/maintenance.js`: `list()` e `update(id, patch)`
  2. `lib/maintenance.js`: `windowState(w, now)` (usa `state` do servidor e vira "em curso" quando o relógio local passa de `starts_at`, só para a UI não esperar o próximo poll) e `validateWindow` ajustada para edição (início futuro exigido só quando o início mudou; janela iniciada não valida horário; mensagens iguais às do servidor)
  3. Teste em `lib/maintenance.test.js`: edição sem mudar o início não acusa "início no passado"; troca de estado na virada do `starts_at`; janela iniciada ignora a validação de horário
  4. Rota `/maintenance` em `main.jsx` e item "Manutenções" em `Shell.jsx` e `BottomNav.jsx`
- **VERIFY**:
  - `cd web && npm test && npm run lint`

### Task T30: Página `Maintenance.jsx`

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T29
- **PRE-CONDITIONS**:
  - Fluxo da T28 aprovado
- **INPUT**: `Incidents.jsx` e `MaintenanceCard.jsx` como molde; `useNow` para o relógio
- **OUTPUT**:
  1. `pages/Maintenance.jsx` + `components/maintenance/WindowsTable.jsx`: carrega com `maintenanceApi.list()`, mostra datas no fuso das Configurações (`formatDateTime` já existente) e o nome do fuso no cabeçalho (D10)
  2. `CancelWindowDialog.jsx`: confirmação com os dois textos da T28; chama `maintenanceApi.cancel` e recarrega; toast de sucesso
  3. Atualização sem recarregar: `useNow` vira a linha para "Em curso" (o "Cancelar" vira "Desligar") e recarrega a lista quando passa do `starts_at` ou `ends_at` de alguma janela (a que terminou sai da lista)
  4. "Agendar manutenção" abre o `MaintenanceModal` em modo criar; ao criar, recarrega
  5. "Editar" abre o modal da T31 e, ao salvar, recarrega
  6. Erros de rede da lista aparecem na página; erro de cancelar aparece em toast
- **VERIFY**:
  - `npm run lint && npm test && npm run build`
  - Manual: janela de 2 min cruza agendada → em curso → some, sem recarregar a página

### Task T31: `MaintenanceModal` em modo editar

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T29
- **PRE-CONDITIONS**:
  - `maintenanceApi.update` disponível
- **INPUT**: `MaintenanceModal` com lista de monitores, datas e motivo; só cria
- **OUTPUT**:
  1. Prop `window` (opcional). Com ela: título "Editar manutenção", botão "Salvar", **sem a lista de monitores** (mostra o nome do monitor como texto fixo, D13), início/fim preenchidos com `utcToZonedInput` no fuso das Configurações, motivo preenchido
     1b. Janela iniciada (`windowState` = em curso): início e fim `disabled` com o texto de apoio "A manutenção já começou: só o motivo pode mudar" ligado por `aria-describedby`; só o motivo é enviado
  2. Envia `maintenanceApi.update(id, patch)` só com os campos que mudaram; nada mudou → fecha sem chamar a API
  3. 409 `MAINTENANCE_LOCKED` (a janela começou com o modal aberto e o horário foi alterado): mensagem "Esta manutenção começou enquanto você editava, então o horário não pode mais mudar. O motivo ainda pode", mantém o modal aberto com o horário desabilitado, e chama `onStale()` para quem abriu recarregar (D15)
  4. 404: "Esta manutenção não existe mais" + `onStale()`
  5. Reaproveita validação, foco no primeiro erro, `aria-invalid`/`aria-describedby` e o rótulo com o nome do fuso já existentes
  6. Modo criar continua igual (inclusive monitores pausados e com janela viva desabilitados, D11)
- **VERIFY**:
  - `npm run lint && npm test && npm run build`
  - Manual: editar horário e motivo de uma agendada; editar só o motivo de uma em curso; deixar o `starts_at` passar com o modal aberto e salvar o horário → mensagem do 409 e horário desabilitado

### Task T32: "Editar" no `MaintenanceCard` do detalhe

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T31
- **PRE-CONDITIONS**:
  - Modal em modo editar pronto (T31)
- **INPUT**: `MaintenanceCard.jsx` mostra motivo, início, fim, estado e "Cancelar"; `MonitorDetail.jsx` já recarrega o monitor depois de cancelar
- **OUTPUT**:
  1. Botão "Editar" no card, para janela agendada e em curso (mesmas regras da página, D12), ao lado de "Cancelar"/"Desligar"
  2. Abre `MaintenanceModal` com `window` montada a partir de `monitor.maintenance` (`id`, `starts_at`, `ends_at`, `reason`, `active`) e o nome do monitor; o estado vem de `maintenance.active`, sem consulta nova
  3. Ao salvar, `MonitorDetail` recarrega o monitor (mesmo caminho do cancelar); o card mostra o motivo e o horário novos
  4. Rótulo acessível "Editar manutenção" no botão
- **VERIFY**:
  - `npm run lint && npm test && npm run build`
  - Manual: editar pelo detalhe uma janela agendada e uma em curso; o resultado aparece também na página `/maintenance`
- **Why**: quem está no detalhe de um monitor em manutenção quer ajustar ali, sem ir à página (D17).

### Task T33: Revisão de bordas e segurança da Etapa 3

- **Agent**: `qa-review` + `security-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T27, T30, T32
- **PRE-CONDITIONS**:
  - Etapa 3 completa na branch
- **INPUT**: diff da branch
- **OUTPUT**:
  1. QA: editar no segundo exato do início; editar o motivo de janela em curso e conferir que nada mais mudou (incidente, estado, check); editar duas vezes seguidas; editar para um intervalo que já está em curso; pausar o monitor com janela agendada e editar; excluir o monitor com janela e abrir a página; fuso trocado nas Configurações com a página aberta; cancelar a janela que acabou de terminar; duas abas editando a mesma janela; editar pelo card do detalhe e pela página
  2. Segurança: `GET` e `PATCH` atrás do `authMiddleware`; corpo do `PATCH` limitado; ID do caminho validado; motivo e nome do monitor renderizados como texto (sem `dangerouslySetInnerHTML`); SQL do `Update` e do `ListLive` parametrizado; o `PATCH` não permite mudar `monitor_id` (campo ignorado ou 400)
- **VERIFY**:
  - Achados corrigidos ou registrados em Notes

### Task T34: Docs da Etapa 3

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T33
- **PRE-CONDITIONS**:
  - Comportamento final estável
- **INPUT**: `docs/implementacoes.md`, `README.md`
- **OUTPUT**:
  1. `docs/implementacoes.md`: seção da Etapa 3, decisões D12–D17 e as limitações (sem histórico, horário de janela em curso não muda)
  2. README: `GET` e `PATCH /api/maintenance-windows`, erro `MAINTENANCE_LOCKED` e o campo `state`
- **VERIFY**:
  - Revisão visual dos arquivos

### Task T35: Verificação final e PR da Etapa 3

- **Agent**: —
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T34
- **PRE-CONDITIONS**:
  - Tudo verde localmente
- **INPUT**: branch `feature/maintenance-page`
- **OUTPUT**:
  1. Rodar a Phase X inteira
  2. `graphify update .`
  3. Commit(s) por camada, seguindo o padrão do projeto (`feat(maintenance): …`), com a linha de `Co-Authored-By` indicada pela sessão
  4. `git push -u origin feature/maintenance-page` e `gh pr create --base feature/maintenance-windows` (PR empilhado, D16), com resumo, decisões D12–D17, checklist de teste manual da Etapa 3 e aviso no topo: "Depende do PR das etapas 1 e 2"
  5. Depois que o PR das etapas 1 e 2 for mergeado: `git rebase --onto develop feature/maintenance-windows` e trocar a base do PR para `develop` (`gh pr edit --base develop`)
- **VERIFY**:
  - `gh pr view --web` abre o PR com a base certa; CI verde

## Risk Areas

| Risk                                                   | Impact                                                                          | Mitigation                                                                                                                                                    |
| ------------------------------------------------------ | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Fuso: `datetime-local` usa o fuso do navegador         | Janela deslocada em horas                                                       | `zonedToUtc` com o fuso das Configurações + testes com fusos de meia hora e DST; fuso visível no modal                                                        |
| Consulta extra por check no `withMonitor`              | Latência no caminho quente                                                      | Consulta indexada coberta pelo `queryplan_test`; 100 monitores no máximo                                                                                      |
| Corrida entre o tick e um `Update`/`Cancel`            | Estado gravado por cima                                                         | Tudo sob `dbMu`; o tick recarrega o monitor dentro do lock, como `withMonitor`                                                                                |
| Efeitos de início/fim atrasam até um tick              | Incidente fecha ou primeiro check roda até 10 s depois                          | Aceitável; o bloqueio em si é exato (relógio). Documentado                                                                                                    |
| `foreign_keys` desligado no SQLite                     | Janelas órfãs ao excluir monitor                                                | Verificar em T8; apagar no service se preciso                                                                                                                 |
| Troca de `ClosedManually` por enum quebra chamadores   | Build quebrado                                                                  | Grep completo em T2; manter `closed_manually` na resposta JSON                                                                                                |
| Teste de scheduler instável por tempo real             | CI intermitente                                                                 | Relógio e ticks injetáveis; sem `time.Sleep` longo                                                                                                            |
| Editar janela que começa no meio da operação (Etapa 3) | Horário de janela iniciada alterado, ou incidente já fechado com horário antigo | Trava por relógio do servidor (D12) + guarda `starts_at > now` no próprio `UPDATE` do horário; o 409 recarrega a página                                       |
| Editar o motivo de janela em curso (Etapa 3)           | `close_note` dos incidentes fechados no início continua com o motivo antigo     | Aceito: `close_note` é cópia de propósito (ver Notes); só a janela e o card mostram o motivo novo. Documentar na T34                                          |
| Editar mexe no tick de início/fim (Etapa 3)            | Efeitos de início/fim aplicados no horário antigo                               | `DueStart`/`DueEnd` leem `starts_at`/`ends_at` da linha a cada tick; editar antes do início não deixa efeito pendente (`started_at` ainda nulo). Teste na T27 |
| Duas abas editando a mesma janela (Etapa 3)            | Última gravação vence sem aviso                                                 | Aceito: janela pequena, uso de uma pessoa. O `PATCH` só envia campos alterados; registrar em Notes                                                            |
| Página e modal duplicam regras de validação (Etapa 3)  | Mensagens divergentes entre cliente e servidor                                  | `validateWindow` único em `lib/maintenance.js`; o servidor continua sendo a fonte da verdade (T25)                                                            |

## Rollback Strategy

1. **Antes do merge:** fechar o PR e apagar a branch; `develop` não muda.
2. **Depois do merge, só etapa 2:** `git revert` do commit da etapa 2. A 0009 tem Down (`pingou migrate down` uma vez) e a tabela nova não afeta as antigas.
3. **Depois do merge, tudo:** reverter os dois commits e descer as migrações 0009 e 0008. Incidentes com `closed_manually=2` continuam lidos como "fechado" (valor ≠ 0).
4. Banco do usuário: recomendar `GET /api/export` antes de atualizar (já existe).
5. **Etapa 3:** não há migração nova. Antes do merge, fechar o PR e apagar `feature/maintenance-page`; as etapas 1 e 2 não são afetadas. Depois do merge, `git revert` do(s) commit(s) da Etapa 3: as rotas `GET`/`PATCH`, a página, o modo editar e o botão "Editar" do card somem; janelas, `POST` e `DELETE` continuam como estavam. Se o PR das etapas 1 e 2 for revertido, a Etapa 3 vai junto (ela depende dele).

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

Teste manual da Etapa 3 (registrar no PR):

- Página `/maintenance` com 0 janelas → estado vazio; agendar duas (uma em 2 min, outra amanhã) → aparecem ordenadas
- Editar a de amanhã (horário e motivo) → linha atualiza sem recarregar; conferir o UTC gravado com o fuso das Configurações ≠ fuso do navegador
- Deixar a de 2 min começar → linha vira "Em curso" sozinha e o "Cancelar" vira "Desligar"
- "Editar" na janela em curso → início e fim desabilitados; mudar o motivo e salvar → motivo novo na lista e no card, monitor segue em manutenção, sem incidente novo
- "Desligar" → monitor volta a ser checado
- Abrir "Editar" na janela de 2 min antes de começar, esperar passar do início e salvar o horário → mensagem do 409, horário desabilitado, lista atualizada
- Detalhe do monitor em manutenção → "Editar" no card abre o mesmo modal; salvar atualiza o card
- `curl -X PATCH` com `starts_at` em janela iniciada → 409 `MAINTENANCE_LOCKED`; só com `reason` → 200

## Cronograma Sugerido

| Fase                         | Tarefas | Complexity | Esforço | Target                                                                      | Status       |
| ---------------------------- | ------- | ---------- | ------- | --------------------------------------------------------------------------- | ------------ |
| Fase 0                       | T0      | XS         | 5 min   | Branch criada                                                               | ✅ Concluído |
| Fase 1 (etapa 1)             | T1–T7   | M          | ~3 h    | Pausa com motivo + bug da retomada corrigido                                | ✅ Concluído |
| Fase 2 (backend da etapa 2)  | T8–T15  | XL         | ~6 h    | Janela funcionando via API, com testes                                      | ✅ Concluído |
| Fase 3 (frontend da etapa 2) | T16–T18 | L          | ~3–4 h  | UI completa                                                                 | ✅ Concluído |
| Fase 4 (fechamento)          | T19–T21 | M          | ~1,5 h  | PR aberto contra `develop`                                                  | ✅ Concluído |
| Fase 5 (backend da Etapa 3)  | T22–T27 | L          | ~3 h    | `GET`/`PATCH` com testes                                                    | ⬜ Pendente  |
| Fase 6 (frontend da Etapa 3) | T28–T32 | L          | ~3,5 h  | Página `/maintenance`, modal em modo editar e "Editar" no card              | ⬜ Pendente  |
| Fase 7 (fechamento Etapa 3)  | T33–T35 | M          | ~1,5 h  | PR da Etapa 3 aberto (base `feature/maintenance-windows`, depois `develop`) | ⬜ Pendente  |

## Status de Execução

### Concluído

- Análise de viabilidade e decisões D1–D9 aprovadas
- T0–T21: etapas 1 e 2, testes, QA manual, docs e PR (desvios registrados em `docs/implementacoes.md`, seção 14)

### Em andamento

-

### Pendente

- QA em dispositivo real / leitor de tela (registrado em Pendências abertas)
- Etapa 3 inteira (T22–T35): página de manutenções, edição de janela (horário só da agendada, motivo também da em curso), desligar janela em curso, "Editar" no card do detalhe. Decisões D12–D17 aprovadas; aguardando o aval final do plano para construir
- Merge do PR das etapas 1 e 2 (a Etapa 3 está empilhada nele, D16)

## Notes

- O gate único em `Scheduler.withMonitor` é o coração da feature: qualquer resultado de check novo no futuro passa por ali e herda o bloqueio.
- `resumed_at` serve à pausa e ao fim da janela; não há campo separado para "fim da manutenção".
- `close_note` copia o motivo no fechamento para o histórico não mudar se o motivo da pausa for alterado depois.
- Baldes do `/stats` em fusos de meia hora seguem fora do escopo (registrado em `docs/implementacoes.md`).
- Etapa 3: "iniciada" vem do relógio (`starts_at <= agora`), não de `started_at`, porque o tick pode atrasar até 10 s e a janela já está bloqueando o monitor nesse intervalo (D6, D12).
- Etapa 3: editar uma janela agendada não precisa tocar em incidentes nem em estado do monitor, já que os efeitos de início só rodam quando `starts_at` chega (T13).
- Etapa 3: editar o motivo de janela em curso grava só em `maintenance_windows.reason`. O `close_note` dos incidentes fechados no início da janela é cópia do motivo da época e não muda (mesma razão da nota anterior sobre `close_note`).
- Etapa 3: `PATCH` com `starts_at` ou `ends_at` em janela iniciada recusa mesmo quando o valor é igual ao atual. Contrato simples; o front só envia o que mudou.

### Alternativas consideradas e descartadas

- `Flag is_in_maintenance em monitors` → descartada por exigir que o tick esteja sempre certo; derivar do relógio (D6) é exato e sobrevive a reinício.
- `Um timer por janela (time.AfterFunc)` → descartada: some no reinício e precisa de reconstrução; o tick único é o padrão do `runHeartbeats`.
- `Coluna canceled_at` → descartada: cancelar antes apaga, cancelar durante encerra (`ends_at=now`); o histórico não pede mais que isso.
- `Biblioteca de date picker` → descartada: `datetime-local` nativo + helper de fuso resolve.
- `Criação em lote parcial` → descartada por D5 (tudo ou nada é mais previsível).
- `Coluna nova close_reason TEXT` → descartada: `closed_manually` INTEGER sem CHECK já comporta 0/1/2 (D9).
- `Montar a página com GET /monitors e filtrar maintenance != null` → descartada (Etapa 3): paginada por monitor, devolve payload grande e não dá o nome do monitor ordenado por início; um `GET /maintenance-windows` dedicado é uma consulta só.
- `Editar = cancelar e criar outra` → descartada (Etapa 3): perde o ID, pode falhar no meio (cancelou e não criou) e briga com a trava D2; `PATCH` é atômico.
- `Permitir editar o fim de janela em curso (estender/encurtar)` → descartada (Etapa 3): o horário de janela iniciada é travado por decisão; mexer no fim em curso exigiria reaplicar efeitos do tick (T13). Para parar antes, "Desligar".
- `PR da Etapa 3 direto contra develop, sem esperar o merge` → descartada (Etapa 3): o diff mostraria as etapas 1 e 2 junto. PR empilhado (D16) mantém a revisão só do que é novo.
- `Tela de histórico de janelas encerradas` → descartada (Etapa 3): cancelar antes do início apaga a linha e não há coluna de histórico; precisaria de migração.
