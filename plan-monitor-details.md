# Plano de Desenvolvimento: Página de detalhes do monitor

## Overview

Criar a rota `/monitors/:id` conforme [Detalhes.png](concept-design/Detalhes.png): cabeçalho com nome, badge de estado, alvo, intervalo e "último check há N s"; ações Pausar / Editar / ⋮; bloco **Desempenho** com seletor 24h / 7 dias / 14 dias / 30 dias (só habilitado até a retenção configurada, com tooltip explicando os desabilitados), 4 cards (Uptime, Resposta média, Incidentes, Tempo fora do ar) e gráfico de tempo de resposta com faixas "fora do ar" e tooltip por ponto; **Últimos checks** (tabela de 8 + "ver mais"); **Incidentes** recentes do monitor (3 + "Ver todos"); **Configuração** (tipo, intervalo, timeout, threshold, criado em). Concretiza [tarefa-detalhe-monitor-checks-incidentes.md](tasks/tarefa-detalhe-monitor-checks-incidentes.md).

**O que já existe e serve sem mudança** (verificado no código):

| Bloco da tela             | Fonte                                                                                                              |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Cabeçalho e Configuração  | `GET /api/monitors/{id}` ([monitors.go](internal/handler/monitors.go))                                             |
| Últimos checks + ver mais | `GET /api/monitors/{id}/checks?limit=8&offset=` ([checks.go](internal/handler/checks.go)), `total` já é `COUNT(*)` |
| Incidentes recentes       | `GET /api/monitors/{id}/incidents?limit=3` (já traz `duration_seconds`, `last_error`, `open`)                      |
| Editar / Excluir / Pausar | `MonitorForm` + `Modal` + `PATCH`/`DELETE` já usados no [Dashboard.jsx](web/src/pages/Dashboard.jsx)               |
| Menu ⋮                    | [RowMenu.jsx](web/src/components/monitors/RowMenu.jsx) (Editar, Pausar, Ver histórico, Excluir)                    |
| Badge, alvo, datas        | `Badge`, `statusOf`, `targetOf`, `formatDuration`, `timeAgo`, `startedLabel`                                       |

**O que falta** (causa raiz):

1. **Não existe agregação por janela.** O único resumo é `CheckRepo.Summaries` (últimos 30 checks + uptime 24 h fixo, para a lista). O bloco Desempenho precisa de uptime, média de resposta e série temporal para 24 h / 7 d / 14 d / 30 d. Baixar os checks brutos para o browser não serve: 30 d a 1 min = 43 200 linhas. → **endpoint novo `GET /api/monitors/{id}/stats`**, uma query agrupada por balde de tempo.
2. **"Incidentes" e "Tempo fora do ar" da janela** exigem incidentes que **se sobrepõem** à janela (um incidente que começou ontem e ainda está aberto conta em "24 h") e a duração **recortada** à janela. A listagem de incidentes filtra por status, não por período. → método de repositório pequeno, usado só pelo `stats`.
3. **Sem rota no frontend**, sem link do painel para ela, e o item "Painel" da sidebar usa `end` em `/`, então não fica ativo em `/monitors/:id` como no PNG.
4. **Modais de editar/excluir vivem dentro do Dashboard.** A página nova precisa dos mesmos; copiar ~60 linhas duplicaria a regra de `run`/`refetch`. → extrair.

**Sem migration.** Os índices existentes cobrem as consultas: `idx_checks_monitor_checked (monitor_id, checked_at DESC)` para a série e `idx_incidents_monitor_ended (monitor_id, ended_at)` para os incidentes da janela.

**Limites (fora do escopo):** zoom ou intervalo customizado, exportar CSV, filtrar checks por falha, editar inline, tempo real (WebSocket tem tarefa própria), tema escuro, percentis (p95/p99).

**Restrições:**

- Sem dependência nova. Gráfico e tooltips em SVG/HTML à mão (uma linha + retângulos + uma caixa posicionada); lib de chart seria a maior dependência do bundle para um desenho só.
- SQLite com `MaxOpenConns(1)`: o `stats` faz 2 queries curtas, sem transação.
- Períodos maiores que `retention_days` ficam **desabilitados** na tela (não há dados para eles). A API continua aceitando os 4 valores: a regra é de apresentação, e a retenção pode mudar entre um request e outro sem que isso vire erro.
- Branch `feat/monitor-details` a partir de `develop`. Fase 1 (API) antes da Fase 2 (tela); a tela pode começar com mock assim que o contrato abaixo estiver fechado.

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

- [ ] `GET /api/monitors/{id}/stats?range=24h|7d|14d|30d` responde em < 50 ms com 30 d de checks a 1 min (43 200 linhas) no SQLite local
- [ ] `range` inválido → 400 `VALIDATION_ERROR`; monitor inexistente → 404
- [ ] Uptime e média batem com um cálculo manual sobre os mesmos checks (teste)
- [ ] Média de resposta usa **só checks com sucesso**; balde sem sucesso vem `null` e vira buraco na linha
- [ ] Incidente que começou antes da janela ou ainda está aberto entra na contagem e tem a duração recortada à janela (teste)
- [ ] `series` tem sempre o mesmo tamanho por janela (144 / 168 / 112 / 120 pontos), baldes vazios incluídos
- [ ] Botões de período "24h", "7 dias", "14 dias", "30 dias"; os maiores que `retention_days` ficam desabilitados e mostram, no hover **e** no foco por teclado, o tooltip com o motivo e o link para Configurações
- [ ] `?range=` acima da retenção (link antigo, retenção reduzida) cai no maior período permitido, sem request inválido
- [ ] Tooltip do gráfico no hover, toque e setas do teclado: horário do intervalo, média, checks/falhas e "fora do ar" quando o ponto cai numa queda
- [ ] Rota `/monitors/:id` bate com [Detalhes.png](concept-design/Detalhes.png) em 1440 px e é usável em 375 px sem scroll horizontal
- [ ] Nome do monitor no painel e nos incidentes leva à página; "Painel" fica ativo na sidebar
- [ ] Pausar/Retomar, Editar e Excluir funcionam na página (Excluir volta para `/`)
- [ ] Heartbeat: sem card de resposta média e sem gráfico de latência (não há latência); TCP e HTTP com tudo
- [ ] Estados: carregando, monitor não encontrado, monitor sem checks ainda, pausado, erro com "Tentar de novo"
- [ ] Seletor de janela fica na URL (`?range=7d`) e sobrevive a F5
- [ ] README documenta o `/stats`

## Tech Stack

| Technology                                        | Choice Rationale                                                                                            |
| ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `GROUP BY strftime('%s', checked_at) / :balde`    | Série + totais numa query só, coberta pelo índice `(monitor_id, checked_at)`; `checked_at` já é RFC3339 UTC |
| Totais derivados dos baldes em Go                 | Uptime e média saem da soma dos baldes: sem segunda varredura em `checks`                                   |
| SVG à mão (`<path>` + `<rect>`)                   | Uma linha, eixos fixos e faixas vermelhas; ~100 linhas, zero dependência                                    |
| Pointer Events + índice por `x`                   | Tooltip do gráfico: ponto mais próximo = `round(x / largura * (n-1))`, O(1), funciona com mouse e toque     |
| `aria-disabled` + `aria-describedby` nos períodos | `disabled` nativo não recebe foco nem hover, e o tooltip com o motivo precisa aparecer para teclado e mouse |
| `useSearchParams`                                 | Janela na URL, mesmo padrão da página de incidentes                                                         |
| `MonitorModals` extraído do Dashboard             | Uma fonte para editar/excluir/regenerar token nas duas páginas                                              |

## Contrato da API (alvo)

`GET /api/monitors/{id}/stats?range=24h` (padrão `24h`; aceita `24h`, `7d`, `14d`, `30d`)

| `range` | Rótulo  | Balde  | Pontos |
| ------- | ------- | ------ | ------ |
| `24h`   | 24h     | 10 min | 144    |
| `7d`    | 7 dias  | 1 h    | 168    |
| `14d`   | 14 dias | 3 h    | 112    |
| `30d`   | 30 dias | 6 h    | 120    |

Como a retenção mínima é 7 dias, "24h" e "7 dias" estão sempre habilitados.

```json
{
  "range": "24h",
  "since": "2026-10-28T14:10:00Z",
  "until": "2026-10-29T14:10:00Z",
  "bucket_seconds": 600,
  "uptime": 99.93,
  "avg_latency_ms": 214,
  "checks": 1440,
  "incidents": 1,
  "downtime_seconds": 60,
  "series": [
    {
      "t": "2026-10-28T14:10:00Z",
      "avg_latency_ms": 203,
      "checks": 10,
      "failures": 0
    },
    {
      "t": "2026-10-28T14:20:00Z",
      "avg_latency_ms": null,
      "checks": 0,
      "failures": 0
    }
  ],
  "down": [{ "start": "2026-10-29T10:12:00Z", "end": "2026-10-29T10:13:00Z" }]
}
```

- `since` alinhado ao início do balde (`until - range`, arredondado para baixo) para os pontos não "andarem" a cada refresh.
- `uptime` e `avg_latency_ms` são `null` sem checks na janela. `avg_latency_ms` arredondado para inteiro.
- `down`: incidentes que se sobrepõem à janela, recortados a `[since, until]`; aberto usa `until` como fim. `downtime_seconds` é a soma deles.
- `incidents` = `len(down)`.

## File Structure

```text
internal/
├── domain/
│   ├── check.go                     # + MonitorStats, StatsBucket, DownWindow
│   └── repository.go                # CheckRepository.Buckets; IncidentRepository.FindInWindow
├── repository/
│   ├── check_repo.go                # + Buckets (GROUP BY balde)
│   └── incident_repo.go             # + FindInWindow (sobreposição à janela)
├── service/
│   ├── monitor_service.go           # + Stats(ctx, id, range): janelas, preenche baldes, recorta quedas
│   └── stats_test.go                # novo: recorte, baldes vazios, médias
└── handler/
    ├── server.go                    # + GET /monitors/{id}/stats
    ├── stats.go                     # novo: handler + DTO
    └── handler_test.go              # + 400 / 404 / formato

web/src/
├── main.jsx                         # + <Route path="/monitors/:id">
├── api/monitors.js                  # + stats(id, range)
├── lib/format.js                    # + fullDate, whenLabel, bucketLabel ("27 out, 14:00 – 15:00")
├── components/
│   ├── layout/Shell.jsx             # "Painel" ativo em /monitors/:id
│   ├── monitors/
│   │   ├── MonitorModals.jsx        # novo: extraído do Dashboard (editar, excluir, token)
│   │   └── MonitorRow.jsx           # nome vira <Link>
│   ├── incidents/                   # nome do monitor vira <Link> (cards e tabela)
│   └── monitor-detail/
│       ├── ResponseChart.jsx        # novo: SVG + tooltip por ponto
│       ├── RangeTabs.jsx            # novo: 24h / 7 dias / 14 dias / 30 dias, limitado pela retenção
│       ├── ChecksTable.jsx          # novo
│       ├── RecentIncidents.jsx      # novo
│       └── ConfigCard.jsx           # novo
└── pages/
    ├── Dashboard.jsx                # usa MonitorModals
    └── MonitorDetail.jsx            # novo
```

## Task Breakdown

| Task ID | Agent                 | Name                                                     | Priority | Complexity | Dependencies              |
| ------- | --------------------- | -------------------------------------------------------- | -------- | ---------- | ------------------------- |
| T1      | go-backend-specialist | `CheckRepo.Buckets`                                      | P1       | S          | -                         |
| T2      | go-backend-specialist | `IncidentRepo.FindInWindow`                              | P1       | S          | -                         |
| T3      | go-backend-specialist | `MonitorService.Stats`                                   | P1       | M          | T1, T2                    |
| T4      | go-backend-specialist | Handler e rota `/stats`                                  | P1       | S          | T3                        |
| T5      | code-test             | Testes de stats (service + handler)                      | P1       | M          | T4                        |
| T6      | refactor              | Extrair `MonitorModals` do Dashboard                     | P2       | S          | -                         |
| T7      | frontend-specialist   | Rota, API client, sidebar e links                        | P2       | S          | -                         |
| T8      | frontend-specialist   | Helpers de data (`fullDate`, `whenLabel`, `bucketLabel`) | P2       | S          | -                         |
| T9      | frontend-specialist   | `ResponseChart` (SVG)                                    | P2       | M          | T4                        |
| T10     | frontend-specialist   | `ChecksTable`, `RecentIncidents`, `ConfigCard`           | P2       | M          | T8                        |
| T11     | frontend-specialist   | Página `MonitorDetail.jsx`                               | P2       | L          | T6, T7, T9, T10, T14, T15 |
| T12     | qa-review             | QA contra o PNG, 3 tipos, mobile, teclado                | P3       | M          | T11                       |
| T13     | document              | README: `/stats` e rota de detalhe                       | P3       | XS         | T4                        |
| T14     | frontend-specialist   | Tooltip do gráfico                                       | P2       | M          | T8, T9                    |
| T15     | frontend-specialist   | `RangeTabs` limitado pela retenção + tooltip             | P2       | S          | -                         |

### Task T1: `CheckRepo.Buckets`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Branch `feat/monitor-details` criada a partir de `develop`
- **INPUT**: [check_repo.go](internal/repository/check_repo.go) só tem `FindByMonitor` e `Summaries`.
- **OUTPUT**:
  1. `domain.StatsBucket{ Start time.Time; Checks, Failures int; AvgLatencyMs *float64; SumLatencyMs int64; Successes int }` (soma e contagem ficam para o service derivar a média total sem recalcular).
  2. `Buckets(ctx, monitorID string, since, until time.Time, bucketSeconds int) ([]StatsBucket, error)`:
     ```sql
     SELECT CAST(strftime('%s', checked_at) AS INTEGER) / ? AS b,
            COUNT(*), SUM(1 - success),
            SUM(CASE WHEN success = 1 THEN latency_ms END), SUM(success)
     FROM checks
     WHERE monitor_id = ? AND checked_at >= ? AND checked_at < ?
     GROUP BY b ORDER BY b
     ```
     Só baldes com dados; preencher vazios é do service.
  3. Interface `CheckRepository` ganha o método.
- **WHY**: uma varredura pelo índice resolve série e totais; nenhum check bruto sai do banco.
- **VERIFY**: `go build ./... && go vet ./...`; `EXPLAIN QUERY PLAN` mostra `USING INDEX idx_checks_monitor_checked`.
- **ROLLBACK**: reverter o commit; nada depende ainda.

### Task T2: `IncidentRepo.FindInWindow`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: None
- **INPUT**: [incident_repo.go](internal/repository/incident_repo.go) lista por status/monitor com paginação; não filtra por período.
- **OUTPUT**:
  1. `FindInWindow(ctx, monitorID string, since time.Time) ([]*domain.Incident, error)`: `WHERE monitor_id = ? AND (ended_at IS NULL OR ended_at > ?) ORDER BY started_at`, sem JOIN (o monitor já é conhecido), reusando `scanIncident`.
  2. Sem `LIMIT`. `// ponytail: sem limite; um monitor com milhares de quedas em 30 d é patológico, pôr LIMIT 1000 se aparecer`.
- **VERIFY**: `go build ./...`; coberto em T5.
- **ROLLBACK**: reverter o commit.

### Task T3: `MonitorService.Stats`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1, T2
- **INPUT**: `MonitorService` tem `checks` e `incidents` injetados.
- **OUTPUT**:
  1. `domain.MonitorStats{ Range string; Since, Until time.Time; BucketSeconds int; Uptime, AvgLatencyMs *float64; Checks, Incidents int; DowntimeSeconds int64; Series []StatsBucket; Down []DownWindow }`.
  2. `var statsRanges = map[string]struct{ d time.Duration; bucket int }{"24h": {24h, 600}, "7d": {168h, 3600}, "14d": {336h, 10800}, "30d": {720h, 21600}}`; range fora do mapa → `ErrValidation`.
  3. `Stats(ctx, id, rng string, now time.Time)`: `GetByID` (404); `until = now.Truncate(bucket)` + 1 balde (o balde corrente entra parcial); `since = until - d`; chama `Buckets`, encaixa por índice num slice de tamanho fixo (`d / bucket`), soma totais; chama `FindInWindow`, recorta cada incidente a `[since, until]` (aberto → `until`), soma `DowntimeSeconds`.
  4. `now` como parâmetro para o teste ser determinístico (o handler passa `time.Now()`).
- **WHY**: regra de janela e recorte fica num lugar só e testável sem HTTP.
- **VERIFY**: `go test ./internal/service/...` (casos em T5).
- **ROLLBACK**: reverter o commit.

### Task T4: Handler e rota `/stats`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3
- **INPUT**: [server.go](internal/handler/server.go) registra as rotas de `/monitors/{id}/…`.
- **OUTPUT**:
  1. `api.HandleFunc("GET /monitors/{id}/stats", s.handleMonitorStats)`.
  2. `internal/handler/stats.go`: lê `range` (vazio → `24h`), chama o service, converte para o JSON do "Contrato da API" (datas `time.RFC3339` UTC, `avg_latency_ms` arredondado, `null` onde não há dado). Erros via `writeServiceError` (400 / 404 já mapeados).
- **VERIFY**: `curl -H "Authorization: Bearer $TOKEN" 'localhost:8080/api/monitors/<id>/stats?range=24h' | jq '.uptime, .avg_latency_ms, (.series|length)'` → números e `144`.
- **ROLLBACK**: remover a rota e o arquivo; o resto da API não muda.

### Task T5: Testes de stats

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T4
- **INPUT**: [handler_test.go](internal/handler/handler_test.go) já sobe banco real em memória.
- **OUTPUT** (casos, não suítes):
  1. Service, `now` fixo, 24 h: 6 checks no mesmo balde (5 ok a 100/200/300/400/500 ms, 1 falha) → balde com `avg=300`, `failures=1`; uptime `83,33`; `series` com 144 itens e os vazios com `avg_latency_ms=nil`.
  2. Incidente iniciado 30 h atrás e fechado 23 h atrás → conta 1, `downtime = 1 h` (recortado).
  3. Incidente aberto há 2 h → conta, `downtime = 2 h`. Incidente fechado há 2 dias → não conta.
  4. Handler: `?range=1y` → 400; id inexistente → 404; sem `range` → `"range":"24h"`; `?range=14d` → 112 pontos, `bucket_seconds=10800`.
  5. Monitor sem checks → `uptime` e `avg_latency_ms` `null`, `checks=0`.
- **VERIFY**: `go test ./...` verde.
- **ROLLBACK**: n/a.

### Task T6: Extrair `MonitorModals` do Dashboard

- **Agent**: refactor
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **INPUT**: [Dashboard.jsx](web/src/pages/Dashboard.jsx) contém `modal`, `selected`, `saving`, `run`, `handleEdit/Delete/Regenerate`, `TypeNote` e os 3 `<Modal>`.
- **OUTPUT**:
  1. `components/monitors/MonitorModals.jsx`: `<MonitorModals modal selected onClose onSaved onDeleted setSelected />` com os modais de editar, excluir e o `TypeNote`; "criar" fica no Dashboard (só ele cria).
  2. Dashboard passa a usar o componente; comportamento idêntico.
- **WHY**: a página de detalhe abre os mesmos modais; duas cópias de `run`/`refetch` divergiriam no primeiro ajuste.
- **VERIFY**: `npm run lint && npm run build`; no painel, editar, excluir e regenerar token continuam funcionando.
- **ROLLBACK**: `git checkout develop -- web/src/pages/Dashboard.jsx` e apagar o arquivo novo.

### Task T7: Rota, API client, sidebar e links

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **OUTPUT**:
  1. [main.jsx](web/src/main.jsx): `<Route path="/monitors/:id" element={<MonitorDetail />} />`.
  2. [monitors.js](web/src/api/monitors.js): `stats: (id, range) => client.get(\`/monitors/${id}/stats\`, { params: { range } }).then((r) => r.data)`.
  3. [Shell.jsx](web/src/components/layout/Shell.jsx): item "Painel" ativo também quando `pathname.startsWith("/monitors/")` (o `NavLink end` atual não cobre).
  4. [MonitorRow.jsx](web/src/components/monitors/MonitorRow.jsx): nome vira `<Link to={/monitors/:id}>` (mantém o ⋮ intacto); `OpenIncidentCard` e `ResolvedTable` idem para o nome do monitor.
- **VERIFY**: clicar no nome no painel abre a rota (página vazia por ora); sidebar destaca "Painel".
- **ROLLBACK**: reverter o commit.

### Task T8: Helpers de data

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **OUTPUT** em [format.js](web/src/lib/format.js), reusando `MONTHS`, `WEEKDAYS`, `daysAgo`, `clock`:
  1. `fullDate(iso)` → `"20 out 2026"` (Criado em).
  2. `whenLabel(iso, now)` → `"Hoje, 10:12"`, `"Ontem, 08:00"`, `"Domingo, 19:05"` (até 6 dias), `"24 out, 03:20"`.
  3. `bucketLabel(startIso, bucketSeconds, now)` → cabeçalho do tooltip do gráfico: `"Hoje, 14:10 – 14:20"`, `"Ontem, 14:00 – 15:00"`, `"27 out, 06:00 – 12:00"` (dia de `whenLabel` + `clock` do início e do fim). Balde que vira a meia-noite leva só a data do início.
- **VERIFY**: conferir no dev server com 4 datas (hoje, ontem, 3 dias, 10 dias) e um balde de 6 h cruzando a meia-noite.
- **ROLLBACK**: remover as funções.

### Task T9: `ResponseChart` (SVG)

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T4
- **INPUT**: `series`, `down`, `since`, `until` do `/stats`.
- **OUTPUT**: `components/monitor-detail/ResponseChart.jsx`:
  1. `viewBox` fixo e `preserveAspectRatio="none"` no plot, rótulos do eixo Y em HTML ao lado (texto não deforma). Altura ~220 px.
  2. Escala Y: 0 até o máximo da série arredondado para cima em múltiplos de 200 ms (mínimo 200); 4 linhas de grade `pg-borda` com rótulos `0 / 200 ms / …`.
  3. Linha `pg-ciano` (`stroke-width 2`, `vector-effect="non-scaling-stroke"`); balde `null` quebra o `path` (novo `M`), sem ligar por cima do buraco.
  4. Cada item de `down` vira `<rect>` `pg-down-bg` + linha tracejada `pg-down` no início; mínimo de 2 px de largura para quedas curtas aparecerem em 30 d.
  5. Legenda "média por intervalo" / "fora do ar" como no PNG.
  6. `role="img"` + `aria-label` "Tempo de resposta médio 214 ms, pico 402 ms, 1 período fora do ar".
  7. Série toda `null` → "Sem dados de resposta neste período."
  8. Expõe para a T14 a função `xOf(i)`/`yOf(ms)` e o container `relative` onde o tooltip é posicionado.
- **WHY**: lib de gráfico (recharts ~100 kB) para uma linha não se paga.
- **VERIFY**: com o JSON de exemplo e com um mock de 30 d, visual lado a lado com o PNG; buraco e faixa vermelha no lugar certo.
- **ROLLBACK**: apagar o arquivo.

### Task T10: `ChecksTable`, `RecentIncidents`, `ConfigCard`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T8
- **OUTPUT** em `components/monitor-detail/`:
  1. **ChecksTable**: `<table>` Horário (`clock` com segundos) / Status (✓ Sucesso `pg-up`, ✕ Falha `pg-down`) / Resposta (`tabular-nums`, `—` em heartbeat) / Detalhe (mono: `HTTP 200` | `error_message` truncado com `title` | `conexão aberta` (tcp) | `ping recebido` (heartbeat)). Rodapé "Mostrando os N mais recentes · ver mais" (`ver mais` acrescenta +20 via `offset`, some quando `N >= total`). Abaixo de `sm` esconde a coluna Detalhe e mostra o erro sob o status.
  2. **RecentIncidents**: 3 itens: ícone ✓ (resolvido) ou ✕ (aberto), `whenLabel(started_at)`, `last_error` mono truncado, duração à direita (`formatDuration`; aberto: "em andamento"). Link "Ver todos" → `/incidents?monitor=<id>`. Vazio: "Nenhum incidente registrado."
  3. **ConfigCard**: `<dl>` Tipo / Verifica a cada (`formatDuration(interval_seconds)`) / Tempo limite (http, tcp) ou Tolerância (`grace_seconds`, heartbeat) / Cai depois de (`N falhas seguidas`; omitido em heartbeat) / Criado em (`fullDate`). Link "Editar" chama `onEdit`.
- **VERIFY**: renderizar com dados dos 3 tipos; teclado alcança "ver mais", "Ver todos" e "Editar".
- **ROLLBACK**: apagar os arquivos.

### Task T11: Página `MonitorDetail.jsx`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T6, T7, T9, T10
- **OUTPUT**:
  1. Voltar "‹ Monitores" (`Link` para `/`).
  2. Cabeçalho: `h1` nome + `Badge state={statusOf(m)}`; linha com ícone + `targetOf(m)` (link `target="_blank" rel="noopener noreferrer"` só em http) e "a cada {formatDuration} · último check {timeAgo}" (`useNow(1000)`). Heartbeat: "último ping".
  3. Ações: botões Pausar/Retomar e Editar (escondidos abaixo de `sm`) + `RowMenu` reaproveitado sem mudança (no mobile ele é o único acesso; no desktop a repetição é aceitável). Excluir → `navigate("/")` após sucesso.
  4. Desempenho: `RangeTabs` (T15) com `retention_days` de `useSettings`, estado em `?range=`; enquanto `useSettings` carrega, só "24h" e "7 dias" ficam clicáveis (sempre válidos). 4 cards: Uptime (`formatUptime`), Resposta média ("só checks com sucesso"; **oculto em heartbeat**), Incidentes, Tempo fora do ar (`formatDuration`, "0 s" vira "nenhum"). Rodapé dos cards: "nas últimas 24 h" / "nos últimos 7 dias" / "nos últimos 14 dias" / "nos últimos 30 dias". Gráfico oculto em heartbeat.
  5. Grid inferior: Últimos checks (2/3) | Incidentes + Configuração (1/3); empilha abaixo de `lg`.
  6. Dados: `get`, `checks(limit 8)`, `incidents(limit 3)` a cada 15 s (mesmo ritmo do painel); `stats` ao trocar janela e a cada 60 s. Flag `alive` para descartar respostas atrasadas ao trocar de monitor/janela.
  7. Estados: Spinner na carga inicial; 404 → "Monitor não encontrado" + link para o painel; sem checks → cards com "—" e "Aguardando o primeiro check"; erro → mensagem + "Tentar de novo".
- **VERIFY**: `cd web && npm run lint && npm run build`; navegar nos 3 tipos; trocar janela; F5 em `?range=7d`; pausar e retomar; excluir volta ao painel.
- **ROLLBACK**: remover a rota em `main.jsx` e apagar o arquivo.

### Task T12: QA

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T11
- **OUTPUT**: checklist com o binário real (`make build && ./bin/pingou serve`), banco semeado com 30 d de checks:
  - Lado a lado com [Detalhes.png](concept-design/Detalhes.png) em 1440 px
  - 375 px: sem scroll horizontal, ações no ⋮, grid empilhado
  - HTTP, TCP e heartbeat; monitor pausado; monitor recém-criado
  - Queda em andamento: faixa vermelha vai até a borda direita, card "Incidentes" ≥ 1
  - `retention_days = 7`: "14 dias" e "30 dias" desabilitados, tooltip no hover e no Tab, clique não faz nada; `retention_days = 14`: só "30 dias" desabilitado; `?range=30d` com retenção 7 cai em "7 dias"
  - Tooltip do gráfico: hover, toque (375 px), setas do teclado; não sai da caixa nas bordas; ponto num buraco e ponto dentro de uma queda
  - Teclado: voltar, abas, ações, "ver mais", "Ver todos", "Editar"; foco visível
  - Tempo de `/stats?range=30d` no log do servidor
- **VERIFY**: achados registrados em Notes; nenhum bloqueante aberto.
- **ROLLBACK**: n/a.

### Task T13: README

- **Agent**: document
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T4
- **OUTPUT**: endpoint `/api/monitors/{id}/stats` (parâmetro, baldes, campos) na seção de endpoints; menção à página de detalhe.
- **VERIFY**: exemplo de `curl` roda.
- **ROLLBACK**: reverter o arquivo.

### Task T14: Tooltip do gráfico

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T8, T9
- **INPUT**: `ResponseChart` desenhado, sem interação.
- **OUTPUT** em `ResponseChart.jsx` (sem arquivo novo):
  1. `onPointerMove` / `onPointerDown` no plot: `i = round(offsetX / largura * (n - 1))`, guarda `active = i`; `onPointerLeave` limpa. Toque: `pointerdown` fixa o ponto, tocar fora limpa.
  2. Marca o ponto ativo: linha-guia vertical `pg-borda` + bolinha `pg-ciano` (sem bolinha em balde `null`).
  3. Caixa HTML (`absolute`, `pointer-events-none`, `rounded-pg-control border bg-pg-branco shadow`) acima do ponto, virando para a esquerda na metade direita para não sair do card:
     - `bucketLabel(t, bucket_seconds)`
     - "Média **214 ms**" (ou "Sem resposta" em balde `null`)
     - "{checks} checks · {failures} falha(s)" (`failures` em `pg-down` quando > 0)
     - "Fora do ar" em `pg-down` quando o intervalo do balde cruza algum item de `down`
  4. Teclado: plot com `tabIndex=0`, `←`/`→` movem `active`, `Home`/`End` vão às pontas, `Esc` limpa; foco visível. A caixa tem `role="status"` `aria-live="polite"` para leitor de tela anunciar o ponto.
- **WHY**: o `aria-label` geral resume, mas não deixa ler um ponto; teclado + `aria-live` resolvem sem tabela escondida.
- **VERIFY**: hover nas pontas (caixa não corta), num buraco e numa faixa vermelha; Tab até o gráfico + setas anunciam o ponto; toque em 375 px.
- **ROLLBACK**: remover os handlers e a caixa; o gráfico estático continua.

### Task T15: `RangeTabs` limitado pela retenção

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **OUTPUT**: `components/monitor-detail/RangeTabs.jsx`:
  1. `const RANGES = [{ key: "24h", label: "24h", days: 1 }, { key: "7d", label: "7 dias", days: 7 }, { key: "14d", label: "14 dias", days: 14 }, { key: "30d", label: "30 dias", days: 30 }]`.
  2. Props `value`, `onChange`, `retentionDays`. Habilitado se `days <= retentionDays`. Exporta `clampRange(key, retentionDays)` → maior período permitido ≤ o pedido; a página usa ao ler `?range=` (link antigo ou retenção reduzida) e reescreve a URL com `replace`.
  3. Grupo `role="group" aria-label="Período"`, botões `aria-pressed`, visual de segmento escuro do PNG.
  4. Desabilitados: **`aria-disabled="true"` em vez de `disabled`** (para continuar focável e receber hover), `onClick` ignora, estilo `text-pg-neblina cursor-not-allowed`. Tooltip `role="tooltip"` ligado por `aria-describedby`, visível em `group-hover` e `group-focus-visible` (CSS, sem estado):
     > Os checks são guardados por {retentionDays} dias. Aumente a retenção em Configurações para ver {label}.
  5. "Configurações" no texto não é link (tooltip some no hover-out); o caminho fica no texto. `// ponytail: tooltip em CSS; trocar por popover clicável se pedirem link`.
  6. Tooltip abaixo do botão, alinhado à direita no último para não sair da tela em 375 px.
- **VERIFY**: retenção 7, 14 e 30; Tab passa pelos desabilitados e mostra o tooltip; Enter/clique neles não muda a URL; leitor de tela lê "30 dias, indisponível, Os checks são guardados por 7 dias…".
- **ROLLBACK**: apagar o arquivo e voltar a botões fixos na página.

## Risk Areas

| Risk                                                  | Impact                                      | Mitigation                                                                                                       |
| ----------------------------------------------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `strftime('%s', checked_at)` por linha em 43 k linhas | `/stats?range=30d` lento                    | Filtro por índice antes do cálculo; medir em T12. Se > 50 ms, guardar `checked_at` epoch numa coluna (migration) |
| Fuso: baldes em UTC, usuário em UTC-3                 | Balde de 6 h começa às 21:00 local          | Aceito: gráfico sem rótulos no eixo X, como no PNG                                                               |
| Janela maior que a retenção                           | Uptime "30 dias" calculado sobre 7          | Botão desabilitado + `clampRange` na URL (T15); a API não bloqueia, só a tela                                    |
| Retenção acabou de subir (7 → 30)                     | "30 dias" habilitado com só 7 dias de dados | Aceito: o gráfico mostra buraco à esquerda, que é a verdade; os dados enchem com o tempo                         |
| Tooltip do gráfico em telas de toque                  | Hover não existe                            | `pointerdown` fixa o ponto (T14); verificado em 375 px na T12                                                    |
| `aria-disabled` sem `disabled`                        | Clique ainda dispara se o handler esquecer  | `onClick` checa o estado; caso de QA em T12                                                                      |
| Polling de 4 requests a cada 15 s + 1 a cada 60 s     | Carga no SQLite com `MaxOpenConns(1)`       | Só quando a página está aberta; consultas indexadas e pequenas                                                   |
| `RowMenu` repete Pausar/Editar no desktop             | Leve redundância visual                     | Aceito; trocar por prop `compact` só se incomodar no QA                                                          |
| Extração de `MonitorModals` quebra o painel           | Editar/excluir param no Dashboard           | T6 isolada, verificada antes da página nova                                                                      |

## Rollback Strategy

1. Frontend: reverter os commits da Fase 2; a rota some e o painel volta a ter os modais inline (ou mantém `MonitorModals`, que é compatível).
2. Backend: reverter os commits da Fase 1; só a rota `/stats` e dois métodos de repositório saem, nada mais depende deles.
3. Banco: nada a fazer; não há migration.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm run lint && npm run build

# Binário completo
make build && ./bin/pingou serve
curl -s -H "Authorization: Bearer $PINGOU_API_KEY" "localhost:8080/api/monitors/<id>/stats?range=30d" | jq '.checks, (.series|length)'

# Grafo
graphify update .
```

## Cronograma Sugerido

| Fase               | Tarefas                            | Complexity | Esforço | Target                                 | Status      |
| ------------------ | ---------------------------------- | ---------- | ------- | -------------------------------------- | ----------- |
| Fase 1 — API       | T1, T2, T3, T4, T5                 | L          | ~1 h 30 | `/stats` pronto e testado              | ⬜ Pendente |
| Fase 2 — Tela      | T6, T7, T8, T9, T14, T15, T10, T11 | XL         | ~4 h    | Página igual ao PNG, 3 tipos, tooltips | ⬜ Pendente |
| Fase 3 — Polimento | T12, T13                           | M          | ~45 min | QA e README                            | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T15

## Notes

- Mudanças pedidas após a primeira versão do plano: tooltip no gráfico (T14), período "14 dias" e rótulos "24h / 7 dias / 14 dias / 30 dias", períodos acima da retenção desabilitados com tooltip do motivo (T15).

- Os números do PNG (99,9 %, 214 ms, 1 incidente, 1 min) são consistentes com o contrato: 1 queda de 1 min em 1 440 checks.
- Os incidentes listados no PNG (domingo, 24 out) estão fora da janela de 24 h, mas aparecem porque a lista "Incidentes" é **independente da janela** (são os 3 mais recentes). O card "Incidentes: 1" segue a janela. Manter essa diferença.
- `last_error` dos incidentes já vem sem o prefixo `Get "<url>":` (corrigido no plano de incidentes).
- "ver mais" dos checks pagina na própria página; não há página de histórico de checks.
- `RowMenu` → "Ver histórico" continua indo para `/incidents?monitor=<id>`; o acesso à página de detalhe é pelo nome do monitor.

### Alternativas consideradas e descartadas

- Agregar no frontend a partir de `/checks?limit=…` → descartada: 30 d a 1 min = 43 200 linhas por request.
- Tabela de rollup (`checks_hourly`) preenchida pelo scheduler → descartada: migration + job + backfill para um ganho que a query indexada já entrega; reconsiderar se T12 medir > 50 ms.
- Recharts / Chart.js → descartada: dependência grande para uma linha e retângulos.
- Estender `GET /monitors/{id}` com as estatísticas → descartada: a janela muda independentemente do monitor e o polling de 15 s recalcularia 30 d à toa.
- Reusar `IncidentFilter` com campo `Since` → descartada: a listagem tem JOIN, paginação e `COUNT(*)` que o stats não precisa; um método de 10 linhas é mais claro.
- Percentis (p95) no card → descartado: o PNG mostra média; SQLite não tem percentil nativo.
- `disabled` nativo nos períodos indisponíveis → descartado: não recebe foco nem eventos de mouse, então o tooltip com o motivo nunca apareceria.
- Esconder os períodos acima da retenção → descartado: o usuário não descobriria que pode ter 30 dias mudando a retenção.
- API rejeitar `range` acima da retenção → descartado: a retenção pode mudar entre dois requests; a regra fica na tela.
- `title` nativo como tooltip → descartado: não aparece no foco por teclado nem em toque, e demora ~1 s.
