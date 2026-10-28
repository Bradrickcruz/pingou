# Plano de Desenvolvimento: Reconstrução da página de Incidentes

## Overview

Reconstruir [Incidents.jsx](web/src/pages/Incidents.jsx) para chegar em [Incidentes.png](concept-design/Incidentes.png): incidentes abertos em cards ("Fora do ar há 1 dia e 4 h", falhas seguidas, erro em caixa própria, "Abrir URL"), resolvidos em tabela agrupada por dia com duração, abas Todos/Abertos/Resolvidos, filtro por monitor e "Carregar mais".

**O que já existe e ajuda** (verificado no código):

- `duration_seconds` já vem calculado em `toIncidentResponse` ([incidents.go](internal/handler/incidents.go)), para abertos (desde `started_at`) e resolvidos (coluna gravada no `Close`). Nada a fazer aqui.
- `formatDuration` em [format.js](web/src/lib/format.js) já produz "1 dia e 4 h". `targetOf` em [monitors.js](web/src/lib/monitors.js) e a lógica de `target` do webhook ([notifier.go:48-60](internal/service/notifier.go#L48-L60)) já resolvem URL / host:porta.
- `GET /api/monitors/{id}/incidents` já filtra por monitor.

**O que está quebrado ou faltando** (causa raiz, não sintoma):

1. **Contador da sidebar errado.** `queryIncidents` devolve `len(incidents)` como `total` ([incident_repo.go](internal/repository/incident_repo.go)). A sidebar pede `limit: 1`, então o total é sempre no máximo 1. O mesmo bug impede o "Mostrando 5 de 23". Correção: `COUNT(*)` no repositório, um lugar só.
2. **Sem nome do monitor.** A resposta só traz `monitor_id`. Precisa de JOIN com `monitors`.
3. **Sem "falhas seguidas".** Não existe contagem por incidente, e contar `checks` não serve: o retention apaga checks com mais de 30 dias. Precisa de coluna nova → **migration**.
4. **Sem filtro "Resolvidos".** A API só tem `open=true` (abertos) ou nada (todos).
5. **Erro repete a URL.** O `*url.Error` do Go formata `Get "https://…": <motivo>`. Corrigir na origem (checker) e limpar o que já está gravado (migration).
6. **"Guardados por 30 dias" é falso hoje.** O retention apaga só `checks`; incidentes nunca são apagados. A cópia da página promete isso, então o retention passa a apagar incidentes resolvidos antigos.

**Limites (fora do escopo):** busca textual, filtro por data, página/modal de detalhe do incidente, timeline de eventos, tempo real (WebSocket tem tarefa própria), tema escuro. Ideias de [tarefa-Incidents-melhorias.md](tasks/tarefa-Incidents-melhorias.md) que não estão no PNG ficam para depois.

**Restrições:**

- Sem dependência nova (backend stdlib; frontend React + Tailwind + lucide já instalados).
- SQLite com `MaxOpenConns(1)`: leitura dentro da state machine passa pela transação (`uow`).
- Branch `feat/incidents-page` a partir de `develop`. Backend (Fases 0-1) entra antes do frontend (Fase 2), mas a Fase 2 pode começar com dados mockados assim que o contrato da T5 estiver fechado.

## Project Type

**FULLSTACK** — Go + SQLite (`internal/`) e React/Vite/Tailwind (`web/`), front embutido via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Migration 0004 sobe e desce num banco com dados reais, sem perder incidentes nem checks
- [ ] `GET /api/incidents` devolve `total` real (`COUNT(*)`), independente de `limit`
- [ ] Badge da sidebar mostra **2** com 2 incidentes abertos (hoje mostra 1)
- [ ] Cada incidente traz `monitor.name`, `monitor.type`, `monitor.target` e `failure_count` numa consulta só (sem N+1)
- [ ] `?status=open|resolved` e `?monitor_id=` funcionam e combinam entre si; `?open=true` continua funcionando
- [ ] `failure_count` sobe a cada check com falha enquanto o incidente está aberto e para de subir ao fechar
- [ ] Erros novos e já gravados não começam mais com `Get "<url>": `
- [ ] Timeout aparece como `context deadline exceeded (timeout 10s)`
- [ ] Retention apaga incidentes **resolvidos** mais antigos que `retention_days`; abertos nunca
- [ ] Página bate com [Incidentes.png](concept-design/Incidentes.png) em desktop e é usável em 375 px
- [ ] Estados vazio, carregando e erro existem para abertos e resolvidos
- [ ] "Ver histórico" no menu ⋮ do painel abre `/incidents?monitor=<id>` já filtrado
- [ ] README atualizado (endpoint `/api/incidents` e retenção de incidentes)

## Tech Stack

| Technology                                     | Choice Rationale                                                                               |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| goose SQL migration (`0004`)                   | Mesmo padrão das 0001-0003; `ALTER TABLE ADD/DROP COLUMN` basta, sem recriar tabela            |
| `LEFT JOIN monitors` no `IncidentRepo`         | Nome e alvo numa consulta só; `LEFT` por segurança, embora o `ON DELETE CASCADE` garanta o par |
| `errors.As(err, *url.Error)`                   | Tira o prefixo `Get "url":` na origem, stdlib                                                  |
| `useSearchParams` (react-router, já instalado) | Aba e monitor na URL: "Ver histórico" vira um link simples e o filtro sobrevive a F5           |
| `Intl.DateTimeFormat("pt-BR")`                 | "Domingo · 26 out", "ontem, 11:43" sem lib de datas                                            |

## Contrato da API (alvo)

`GET /api/incidents?status=open|resolved&monitor_id=<id>&limit=20&offset=0`

- `status` vazio = todos. `open=true` continua aceito como sinônimo de `status=open`. `status` inválido → 400 `VALIDATION_ERROR`.
- Ordem: abertos por `started_at DESC`; resolvidos por `ended_at DESC` (o agrupamento por dia usa o fim).
- `GET /api/monitors/{id}/incidents` passa a chamar o mesmo código com `monitor_id` fixo.

```json
{
  "data": [
    {
      "id": "0192…",
      "monitor_id": "0191…",
      "monitor": {
        "id": "0191…",
        "name": "c10",
        "type": "http",
        "target": "https://api.c10.site"
      },
      "started_at": "2026-10-27T14:43:00Z",
      "ended_at": null,
      "last_error": "tls: failed to verify certificate: …",
      "failure_count": 10093,
      "duration_seconds": 101520,
      "notification_sent": true,
      "open": true
    }
  ],
  "total": 2,
  "limit": 20,
  "offset": 0
}
```

`monitor_id` fica no topo para não quebrar quem já consome a API.

## File Structure

```text
internal/
├── database/
│   ├── migrations/0004_incidents_page.sql   # novo: failure_count + backfill, índice, limpeza de prefixo
│   └── migrations_test.go                   # + caso 0004
├── checker/
│   ├── http_checker.go                      # desembrulha *url.Error, anota timeout
│   └── checker_test.go                      # + casos de mensagem
├── domain/
│   ├── incident.go                          # + FailureCount, MonitorName, MonitorType, MonitorTarget; IncidentFilter
│   ├── monitor.go                           # + (m *Monitor) Target() *string (extraído do notifier)
│   └── repository.go                        # FindAll(ctx, IncidentFilter); DeleteResolvedOlderThan
├── repository/
│   ├── incident_repo.go                     # JOIN + COUNT(*); FindByMonitor removido (vira filtro)
│   └── incident_repo_tx.go                  # + IncrementFailureCount; Create grava failure_count
├── service/
│   ├── incident_service.go                  # List(ctx, filter)
│   ├── notifier.go                          # usa m.Target()
│   ├── state_machine.go                     # incrementa failure_count enquanto DOWN
│   └── state_machine_test.go                # + casos
├── scheduler/retention.go                   # + apaga incidentes resolvidos antigos
└── handler/
    ├── incidents.go                         # status, monitor_id, objeto monitor, failure_count
    └── handler_test.go                      # + casos

web/src/
├── api/incidents.js                         # (sem mudança de forma; params novos)
├── lib/format.js                            # + dayLabel, startedLabel
├── components/
│   ├── layout/Shell.jsx                     # badge usa status=open
│   ├── monitors/RowMenu.jsx                 # "Ver histórico" → /incidents?monitor=<id>
│   └── incidents/
│       ├── OpenIncidentCard.jsx             # novo
│       └── ResolvedTable.jsx                # novo (tabela no desktop, lista no mobile)
└── pages/Incidents.jsx                      # reescrita
```

## Task Breakdown

| Task ID | Agent                 | Name                                                      | Priority | Complexity | Dependencies |
| ------- | --------------------- | --------------------------------------------------------- | -------- | ---------- | ------------ |
| T1      | design-db             | Migration 0004                                            | P0       | M          | -            |
| T2      | code-test             | Teste da migration 0004                                   | P0       | S          | T1           |
| T3      | go-backend-specialist | Mensagem de erro limpa no HTTP checker                    | P0       | S          | -            |
| T4      | go-backend-specialist | Repositório: filtro, JOIN e COUNT(\*)                     | P1       | M          | T1           |
| T5      | go-backend-specialist | Handler: contrato novo de `/incidents`                    | P1       | S          | T4, T6       |
| T6      | refactor              | Extrair `Monitor.Target()` do notifier                    | P1       | XS         | -            |
| T7      | go-backend-specialist | State machine incrementa `failure_count`                  | P1       | S          | T1           |
| T8      | go-backend-specialist | Retention apaga incidentes resolvidos                     | P1       | S          | T4           |
| T9      | code-test             | Testes de handler, repo e state machine                   | P1       | M          | T5, T7, T8   |
| T10     | frontend-specialist   | Corrigir badge da sidebar                                 | P2       | XS         | T4           |
| T11     | frontend-specialist   | Helpers de data (`dayLabel`, `startedLabel`)              | P2       | S          | -            |
| T12     | frontend-specialist   | `OpenIncidentCard`                                        | P2       | S          | T5, T11      |
| T13     | frontend-specialist   | `ResolvedTable` agrupada por dia                          | P2       | M          | T5, T11      |
| T14     | frontend-specialist   | Página `Incidents.jsx` (abas, filtro, paginação, estados) | P2       | L          | T12, T13     |
| T15     | frontend-specialist   | "Ver histórico" no menu ⋮ do painel                       | P2       | XS         | T14          |
| T16     | qa-review             | QA da página contra o PNG + mobile + teclado              | P3       | M          | T14, T15     |
| T17     | document              | README: endpoint e retenção                               | P3       | XS         | T5, T8       |

### Task T1: Migration 0004

- **Agent**: design-db
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Backup do `pingou.db` local (`cp pingou.db pingou.db.bak`)
- **INPUT**: `incidents` sem contagem de falhas, sem índice para listagem global; `last_error` e `checks.error_message` com prefixo `Get "<url>": `.
- **OUTPUT**:
  1. `internal/database/migrations/0004_incidents_page.sql` com `-- +goose Up`:
     - `ALTER TABLE incidents ADD COLUMN failure_count INTEGER NOT NULL DEFAULT 0;`
     - Backfill: `UPDATE incidents SET failure_count = (SELECT COUNT(*) FROM checks c WHERE c.monitor_id = incidents.monitor_id AND c.success = 0 AND c.checked_at >= incidents.started_at AND (incidents.ended_at IS NULL OR c.checked_at < incidents.ended_at));` — é aproximado (checks antigos já podem ter sido apagados); aceito.
     - `CREATE INDEX idx_incidents_started_at ON incidents (started_at DESC);` e `CREATE INDEX idx_incidents_ended_at ON incidents (ended_at DESC);` (ordenação das duas listas).
     - Limpeza do prefixo, em `incidents.last_error` e `checks.error_message`: `UPDATE … SET x = substr(x, instr(x, '": ') + 3) WHERE x LIKE 'Get "%": %';`
  2. `-- +goose Down`: `DROP INDEX` dos dois índices e `ALTER TABLE incidents DROP COLUMN failure_count;` (mesmo padrão da 0002). Comentário: a limpeza de prefixo **não é revertida** (era só ruído).
- **WHY**: contar `checks` na hora não serve porque o retention apaga checks; uma coluna incrementada na própria transação do check é exata e barata.
- **VERIFY**:
  - `go run ./cmd/pingou migrate up` no banco copiado; `sqlite3 pingou.db "select failure_count, substr(last_error,1,40) from incidents"` mostra números > 0 e erro sem `Get "`
  - `migrate down` + `migrate up` de novo sem erro
- **ROLLBACK**: `go run ./cmd/pingou migrate down` ou restaurar `pingou.db.bak`.

### Task T2: Teste da migration 0004

- **Agent**: code-test
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **INPUT**: [migrations_test.go](internal/database/migrations_test.go) cobre até a 0003.
- **OUTPUT**:
  1. Caso: sobe até 0003, insere monitor + 3 checks com falha dentro de um incidente aberto e 1 fora, com `last_error = 'Get "http://x": dial tcp …'`; sobe 0004; espera `failure_count = 3` e `last_error = 'dial tcp …'`; desce; coluna some e linhas continuam.
- **VERIFY**: `go test ./internal/database/...` verde.
- **ROLLBACK**: remover o caso.

### Task T3: Mensagem de erro limpa no HTTP checker

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **INPUT**: [http_checker.go:52-54](internal/checker/http_checker.go#L52-L54) usa `err.Error()` direto.
- **OUTPUT**:
  1. Se `errors.As(err, &urlErr)`, usar `urlErr.Err.Error()`.
  2. Se `errors.Is(err, context.DeadlineExceeded)`, mensagem `context deadline exceeded (timeout %ds)` com `m.TimeoutSeconds`.
- **WHY**: corrigir na origem vale para webhook, painel e incidentes de uma vez; cortar no frontend deixaria o webhook com ruído.
- **VERIFY**: caso em `checker_test.go` com `httptest` fechado (connection refused) e outro com handler que dorme além do timeout; nenhum começa com `Get "`.
- **ROLLBACK**: reverter o arquivo; mensagens voltam ao formato antigo, nada mais depende disso.

### Task T4: Repositório: filtro, JOIN e COUNT(\*)

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **INPUT**: `FindAll(ctx, onlyOpen, limit, offset)` e `FindByMonitor(...)` duplicados, sem JOIN, `total = len(página)`.
- **OUTPUT**:
  1. `domain.IncidentFilter{ MonitorID string; Status string /* "", "open", "resolved" */; Limit, Offset int }`.
  2. `domain.Incident` ganha `FailureCount int64`, `MonitorName string`, `MonitorType MonitorType`, `MonitorURL *string`, `MonitorConfig MonitorConfig` (o suficiente para `Target()`, ver T6).
  3. `IncidentRepository.FindAll(ctx, IncidentFilter)` monta um `WHERE` só, usado tanto pela listagem quanto pelo `COUNT(*)` (mesmos args). `LEFT JOIN monitors m ON m.id = i.monitor_id`. Ordem: `ended_at DESC` quando `status=resolved`, senão `started_at DESC`.
  4. `FindByMonitor` (não-tx) removido da interface; chamadores passam a usar `FindAll` com `MonitorID`. `FindOpenByMonitor` e o repo tx continuam iguais, exceto que `scanIncident` passa a ler `failure_count`.
  5. `LIMIT`/`OFFSET` como parâmetros `?` (hoje são `Sprintf`), `limit` limitado a 100.
- **VERIFY**: `go build ./... && go vet ./...`; `total` de uma consulta com `limit=1` e 2 abertos = 2 (coberto em T9).
- **ROLLBACK**: reverter commit; a coluna nova da 0004 é ignorada pelo código antigo.

### Task T5: Handler: contrato novo de `/incidents`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4, T6
- **INPUT**: [incidents.go](internal/handler/incidents.go) lê só `open`, `limit`, `offset`.
- **OUTPUT**:
  1. Lê `status`, `monitor_id` (e `open=true` como sinônimo de `status=open`); `status` fora de `"" | open | resolved` → 400.
  2. `incidentResponse` ganha `monitor {id,name,type,target}` e `failure_count`, conforme "Contrato da API".
  3. `handleListMonitorIncidents` monta o filtro com o `id` do path e reusa o mesmo código.
  4. Remove a linha morta `var _ *service.IncidentService`.
- **VERIFY**: `curl -H "Authorization: Bearer $TOKEN" 'localhost:8080/api/incidents?status=open&limit=1' | jq '.total, .data[0].monitor.name'` → `2` e `"c10"` (ou similar).
- **ROLLBACK**: reverter; o frontend antigo só usa campos que continuam existindo.

### Task T6: Extrair `Monitor.Target()` do notifier

- **Agent**: refactor
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: None
- **INPUT**: lógica de alvo inline em [notifier.go:55-62](internal/service/notifier.go#L55-L62).
- **OUTPUT**: `func (m *Monitor) Target() *string` em `domain/monitor.go` (URL para http, `host:porta` para tcp, `nil` para heartbeat); notifier passa a chamá-lo.
- **VERIFY**: `go test ./internal/service/...` (notifier_test já cobre o `target` do webhook).
- **ROLLBACK**: reverter commit.

### Task T7: State machine incrementa `failure_count`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T1
- **INPUT**: [state_machine.go](internal/service/state_machine.go) cria o incidente na transição →DOWN e não toca nele depois.
- **OUTPUT**:
  1. Ao abrir: `FailureCount = consecutiveFails` (as falhas do threshold contam, "falhas seguidas" começa em 3 e não em 1).
  2. Já `DOWN` e falhou de novo: `IncidentRepoTx.IncrementFailureCount(ctx, tx, monitorID)` → `UPDATE incidents SET failure_count = failure_count + 1 WHERE monitor_id = ? AND ended_at IS NULL`. Uma query, dentro da transação existente, sem SELECT antes.
  3. `last_error` continua sendo o erro **da abertura** (é o "Motivo da queda"; o webhook usa o mesmo valor).
- **VERIFY**: teste em `state_machine_test.go`: threshold 3, 5 falhas → `failure_count = 5`; 1 sucesso → incidente fechado; mais 1 falha não altera o incidente fechado.
- **ROLLBACK**: reverter; a coluna fica com o último valor e para de subir, sem efeito colateral.

### Task T8: Retention apaga incidentes resolvidos

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **INPUT**: [retention.go](internal/scheduler/retention.go) só chama `checks.DeleteOlderThan`.
- **OUTPUT**: `IncidentRepository.DeleteResolvedOlderThan(ctx, before)` → `DELETE FROM incidents WHERE ended_at IS NOT NULL AND ended_at < ?`; chamado no mesmo `purge`, com log igual ao dos checks. Abertos nunca são apagados, mesmo se começaram há 90 dias.
- **WHY**: a página diz "Guardados por 30 dias (mudar em Configurações)"; hoje isso não é verdade.
- **VERIFY**: teste do repo com um resolvido antigo, um resolvido recente e um aberto antigo → só o primeiro some.
- **ROLLBACK**: remover a chamada no `purge`.

### Task T9: Testes de handler, repo e state machine

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T5, T7, T8
- **INPUT**: [handler_test.go](internal/handler/handler_test.go) com banco real em memória.
- **OUTPUT** (casos, não suítes):
  1. 2 abertos + 3 resolvidos, `?status=open&limit=1` → `total=2`, 1 item, `monitor.name` preenchido.
  2. `?status=resolved` → só resolvidos, ordenados por `ended_at DESC`.
  3. `?monitor_id=X&status=resolved` combina os dois filtros.
  4. `?status=xyz` → 400.
  5. `?open=true` ainda funciona.
- **VERIFY**: `go test ./...` verde.
- **ROLLBACK**: n/a.

### Task T10: Corrigir badge da sidebar

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T4
- **INPUT**: [Shell.jsx:28-32](web/src/components/layout/Shell.jsx#L28-L32) já lê `res.total`; o bug era só no backend.
- **OUTPUT**: trocar `{ open: true, limit: 1 }` por `{ status: "open", limit: 1 }`. Nenhuma outra mudança.
- **VERIFY**: com 2 incidentes abertos, badge mostra **2**.
- **ROLLBACK**: reverter a linha (o `open=true` continua aceito).

### Task T11: Helpers de data

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **OUTPUT** em [format.js](web/src/lib/format.js):
  1. `dayLabel(iso, now)` → `"Hoje · 28 out"`, `"Ontem · 27 out"`, `"Domingo · 26 out"` (até 6 dias), depois `"20 out"`; datas de outro ano levam o ano.
  2. `startedLabel(iso, now)` → `"hoje, 08:03"`, `"ontem, 11:43"`, `"26 out, 19:05"`.
  3. `timeRange(start, end)` → `"23:58 → 00:04"`.
  - Agrupar por dia usa o **dia local** do `ended_at` (incidente que vira a meia-noite fica no dia em que terminou).
- **VERIFY**: um `assert` rápido no console do dev server ou um teste simples se já houver runner (hoje `web/` não tem; não adicionar um só para isso).
- **ROLLBACK**: remover as funções.

### Task T12: `OpenIncidentCard`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T5, T11
- **OUTPUT**: `components/incidents/OpenIncidentCard.jsx`:
  1. Ícone X em círculo `pg-down-bg`, nome, "Fora do ar há {formatDuration(now - started_at)}" atualizando com `useNow`, alvo em mono.
  2. Caixa de erro `font-mono`, `whitespace-pre-wrap break-words` (sem truncar).
  3. Lateral: "Começou {startedLabel}", "Falhas seguidas {failure_count.toLocaleString('pt-BR')}".
  4. "Abrir URL" (`<a target="_blank" rel="noopener noreferrer">` com ícone `ExternalLink`) **só para `type=http`**.
  5. Mobile: lateral desce para baixo do erro.
- **VERIFY**: visual lado a lado com o PNG; leitor de tela anuncia "c10, fora do ar há 1 dia e 4 h".
- **ROLLBACK**: apagar o arquivo.

### Task T13: `ResolvedTable` agrupada por dia

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T5, T11
- **OUTPUT**: `components/incidents/ResolvedTable.jsx`:
  1. `<table>` semântica com colunas Monitor / Período / Duração / Motivo da queda; linha de grupo por dia (`<th scope="rowgroup" colSpan=4>`).
  2. Alvo e motivo truncados com `truncate` + `title` com o texto inteiro.
  3. Rodapé "Mostrando {n} de {total} incidentes resolvidos" + botão "Carregar mais" (oculto quando `n >= total`).
  4. Abaixo de `md`: vira lista de cards (monitor, período, duração, motivo em 2 linhas), sem scroll horizontal.
- **VERIFY**: 5 itens em 3 dias renderizam 3 cabeçalhos; "Carregar mais" acrescenta sem duplicar.
- **ROLLBACK**: apagar o arquivo.

### Task T14: Página `Incidents.jsx`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T12, T13
- **INPUT**: página atual com checkbox "Somente abertos".
- **OUTPUT**:
  1. Cabeçalho + subtítulo com `retention_days` vindo de `useSettings` e link para `/settings`.
  2. Abas Todos / Abertos ({n}) / Resolvidos como botões `aria-pressed` (mesmo visual de chip escuro do PNG) e `<select>` de monitor (lista de `useMonitors`). Estado em `useSearchParams`: `?tab=open|resolved&monitor=<id>`.
  3. Duas buscas independentes: abertos (`status=open`, `limit=100`, recarrega a cada 15 s como o painel) e resolvidos (`status=resolved`, `limit=20`, `offset` crescente). Aba Abertos esconde a seção de resolvidos e vice-versa.
  4. Estados: carregando (Spinner por seção), vazio ("Nenhum incidente aberto. Tudo no ar." / "Nenhum incidente resolvido nos últimos N dias" / monitor sem incidentes), erro com "Tentar de novo".
  5. Trocar aba/monitor zera o `offset` e descarta respostas atrasadas (flag `alive` como em `Shell.jsx`).
- **VERIFY**: `cd web && npm run lint && npm run build`; navegar nas 3 abas + filtro "test" (estado vazio) + F5 mantém filtros.
- **ROLLBACK**: `git checkout develop -- web/src/pages/Incidents.jsx`.

### Task T15: "Ver histórico" no menu ⋮

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T14
- **INPUT**: [RowMenu.jsx](web/src/components/monitors/RowMenu.jsx) esconde o item até existir a tela.
- **OUTPUT**: item visível, navega para `/incidents?monitor=<id>`; atualizar o comentário do arquivo.
- **VERIFY**: clicar no painel abre a página já filtrada.
- **ROLLBACK**: esconder o item de novo.

### Task T16: QA

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T14, T15
- **OUTPUT**: checklist executado com o binário real (`make build && ./bin/pingou serve`):
  - Monitor TCP e heartbeat abertos (sem "Abrir URL", alvo correto / heartbeat sem alvo)
  - Monitor apagado com incidentes (some da lista, sem quebrar)
  - Incidente que cruza a meia-noite
  - Nome e erro muito longos em 375 px
  - Navegação só por teclado nas abas, select e "Carregar mais"; foco visível
  - Badge da sidebar acompanha abrir/fechar incidente em até 15 s
- **VERIFY**: lista de achados registrada em Notes; nenhum bloqueante aberto.
- **ROLLBACK**: n/a.

### Task T17: README

- **Agent**: document
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T5, T8
- **OUTPUT**: seção "Endpoints" com `status`, `monitor_id`, campos novos; seção de retenção dizendo que incidentes resolvidos também expiram.
- **VERIFY**: leitura; exemplos de `curl` rodam.
- **ROLLBACK**: reverter o arquivo.

## Risk Areas

| Risk                                                    | Impact                                             | Mitigation                                                                                                       |
| ------------------------------------------------------- | -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| Backfill de `failure_count` lento com muitos checks     | Startup demora na primeira subida                  | Índice `idx_checks_monitor_checked` já cobre o subselect; incidentes são poucos. Medir no `pingou.db` real em T1 |
| Backfill aproximado (checks já apagados pelo retention) | Incidentes antigos mostram menos falhas que o real | Aceito; só afeta incidentes anteriores à 0004                                                                    |
| Limpeza de prefixo corta texto errado                   | Erro antigo fica estranho                          | `LIKE 'Get "%": %'` só pega o formato do `*url.Error`; teste em T2                                               |
| Remover `FindByMonitor` quebra chamador escondido       | Build falha                                        | O compilador acha todos; `graphify query` antes de remover                                                       |
| Retention apagar incidente que o usuário queria         | Perda de histórico                                 | Mesmo `retention_days` que já vale para checks (7-90); abertos nunca são apagados; README avisa                  |
| Polling de 15 s em duas listas + sidebar                | Carga no SQLite                                    | Só a lista de abertos faz polling; resolvidos não mudam sozinhos de forma relevante                              |

## Rollback Strategy

1. Frontend: reverter os commits da Fase 2; a API nova é compatível com a tela antiga (campos só foram acrescentados, `open=true` segue aceito).
2. Backend: reverter os commits da Fase 1; o código antigo ignora `failure_count` e os índices.
3. Banco: `go run ./cmd/pingou migrate down` remove a coluna e os índices. A limpeza de prefixo não volta, o que é inofensivo. Em último caso, restaurar `pingou.db.bak`.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test ./...

# Migration num banco real copiado
cp pingou.db /tmp/pingou-test.db && PINGOU_DATABASE_URL=/tmp/pingou-test.db go run ./cmd/pingou migrate up && PINGOU_DATABASE_URL=/tmp/pingou-test.db go run ./cmd/pingou migrate down

# Frontend
cd web && npm run lint && npm run build

# Binário completo
make build && ./bin/pingou serve

# Grafo
graphify update .
```

## Cronograma Sugerido

| Fase                            | Tarefas                      | Complexity | Esforço | Target                                        | Status      |
| ------------------------------- | ---------------------------- | ---------- | ------- | --------------------------------------------- | ----------- |
| Fase 0 — Banco e origem do erro | T1, T2, T3                   | M          | ~1 h    | Migration sobe/desce; erros sem prefixo       | ⬜ Pendente |
| Fase 1 — API                    | T4, T5, T6, T7, T8, T9       | L          | ~2 h    | Contrato da API pronto e testado; badge certo | ⬜ Pendente |
| Fase 2 — Tela                   | T10, T11, T12, T13, T14, T15 | XL         | ~3 h    | Página igual ao PNG                           | ⬜ Pendente |
| Fase 3 — Polimento              | T16, T17                     | M          | ~45 min | QA e README                                   | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T17

## Notes

- **Não é bug do Pingou, mas vale olhar:** o incidente do `c10` é certificado de outro domínio (`aps.develop.plantvision.app`) servido em `api.c10.site`. Provável proxy (Traefik/Nginx) roteando errado ou certificado não renovado.
- `last_error` mantém o erro da **abertura** do incidente (motivo da queda). Se no futuro o card de aberto precisar do erro atual, ele já existe em `GET /monitors/{id}` (último check), sem mudar o schema.
- `duration_seconds` já era calculado pelo backend; o PRD de "duração nos resolvidos" está atendido sem trabalho.
- Contagem "Abertos (2)" na aba e "Abertos agora 2" saem do `total` da busca de abertos, a mesma que alimenta a sidebar.

### Alternativas consideradas e descartadas

- Contar falhas com `COUNT(*)` em `checks` na hora da leitura → descartada: o retention apaga checks, e a contagem de um incidente longo ficaria errada (e cara).
- Cortar o prefixo `Get "url":` só no frontend → descartada: o webhook continuaria mandando o ruído; corrigir no checker vale para tudo.
- Endpoint novo `/api/incidents/summary` para a sidebar → descartado: o bug era só o `total`; `limit=1` + `COUNT(*)` resolve.
- Paginação por cursor → descartada: `offset` basta para um histórico de 7-90 dias.
- Recriar a tabela `incidents` (como a 0003 fez com `monitors`) → descartada: só há coluna e índice novos, `ALTER TABLE` resolve sem `NO TRANSACTION`.
- Biblioteca de datas (date-fns/dayjs) → descartada: `Intl.DateTimeFormat` cobre os 3 formatos.
