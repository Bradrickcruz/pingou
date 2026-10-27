# Plano de Desenvolvimento: Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel

## Overview

Hoje o Pingou só monitora HTTP. O objetivo é aceitar **três tipos** de monitor (`http`, `tcp`, `heartbeat`) e chegar às telas de [Painel.png](concept-design/Painel.png), [Novo-monitor.png](concept-design/Novo-monitor.png) e [Editar-monitor.png](concept-design/Editar-monitor.png).

Fonte das regras de comportamento: [README.md](README.md) (seções "Tipos de monitor", "Modelo de dados", "Webhook", "Endpoints"). Este plano só adiciona a ordem de execução e as decisões técnicas.

**O que já existe e ajuda** (verificado no código):

- `domain.Checker` (`Check(ctx, *Monitor) CheckResult`) já existe em [checker.go](internal/domain/checker.go), e o `Scheduler` já recebe um `domain.Checker`. A "refatoração" é só criar um despachante por tipo, sem mexer no scheduler.
- `StateMachine.Process` já recebe um `CheckResult` sem saber o tipo. TCP e Heartbeat reaproveitam incidentes, transações e webhook.

**O que muda de verdade:**

- `monitors.url` é `NOT NULL CHECK (length BETWEEN 1 AND 2048)` e `interval_seconds` tem `CHECK` até 86400. TCP/Heartbeat não têm URL, e o design pede intervalo de **1 semana** (604800). Exige recriar a tabela.
- Heartbeat inverte a direção: precisa de rota pública, de um avaliador de prazos e de `last_ping_at`.
- O painel precisa de histórico (30 checks), uptime 24 h e latência dentro de `/api/monitors`.
- O webhook ganha `type` e `target`, sem nunca expor o token.

**Limites (fora do escopo):** WebSocket/tempo real (já há tarefa própria em `tasks/`), tema escuro, tela "Ver histórico" nova (o item do menu aponta para a tela de detalhe que já existe ou fica para depois, ver Notes), DNS/SSL/ICMP.

**Restrições:**

- Sem dependência nova no backend (stdlib: `net.DialTimeout`, `sync`, `time`). Rate limit em memória.
- SQLite com `MaxOpenConns(1)`: tudo que grava passa pelo `dbMu` do scheduler.
- Três PRs/entregas, em ordem, na branch `feat/check-types`. O redesign do frontend (Fase 4) pode correr em paralelo depois da Fase 1.

## Project Type

**FULLSTACK** — Go + SQLite (`internal/`, `cmd/`) e React/Vite/Tailwind (`web/`), com front embutido via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Monitores existentes continuam funcionando sem alteração (migration preserva checks e incidentes)
- [ ] `POST /api/monitors` cria `http`, `tcp` e `heartbeat`; sem `type` cria `http`
- [ ] TCP vira `UP`/`DOWN` conforme a porta, com erro legível em `last_error` e latência = tempo de conexão
- [ ] `/ping/{token}` válido → 200 e check `UP`; token inválido → 404; pausado → 200 sem registrar; excesso por IP → 429
- [ ] Heartbeat sem ping além de `interval + grace` vira `DOWN` **uma vez** (um incidente, sem checks `DOWN` repetidos) e volta a `UP` no próximo ping
- [ ] Heartbeat antes do primeiro ping fica `UNKNOWN`
- [ ] Webhook traz `type` e `target`, mantém `url` só para HTTP e **nunca** contém o token
- [ ] Regenerar token invalida o antigo imediatamente
- [ ] `GET /api/monitors` devolve `history` (últimos 30), `uptime_24h` e `last_latency_ms` sem N+1
- [ ] CLI `add` aceita `--type/--host/--port/--grace`
- [ ] Painel, Novo monitor e Editar monitor batem com os 3 PNGs (itens da Fase 4)
- [ ] README: remover os marcadores "(em desenvolvimento)" ao fim de cada fase

## Tech Stack

| Technology                                              | Choice Rationale                                                      |
| ------------------------------------------------------- | --------------------------------------------------------------------- |
| `net.Dialer.DialContext`                                | TCP check na stdlib, respeita `ctx` e timeout                         |
| `crypto/rand` + base64url (32 bytes)                    | Token de ping imprevisível, sem dependência                           |
| Rate limit em memória (map + mutex, janela fixa por IP) | Um binário só, sem Redis; limite conhecido registrado com `ponytail:` |
| goose `NO TRANSACTION`                                  | Necessário para `PRAGMA foreign_keys=OFF` na recriação da tabela      |
| Coluna `config` TEXT (JSON)                             | Campos por tipo sem uma coluna nova por tipo (como no README)         |
| React + Tailwind (já existente)                         | Reusar `Modal`, `Button`, `Badge`; classes `pg-*` do DESIGN.md        |

## File Structure

```text
internal/
├── database/migrations/
│   └── 0003_monitor_types.sql          # recria monitors (type, config, ping_token, last_ping_at)
├── domain/
│   └── monitor.go                      # + Type, Config, PingToken, LastPingAt; MonitorSummary (history/uptime)
├── checker/
│   ├── http_checker.go                 # (existente)
│   ├── tcp_checker.go                  # novo
│   └── registry.go                     # novo: implementa domain.Checker, despacha por m.Type
├── scheduler/
│   ├── scheduler.go                    # pula heartbeat; expõe Record() (check + dbMu)
│   └── heartbeat.go                    # novo: avaliador de prazos
├── service/
│   ├── validate.go                     # validação por tipo
│   ├── monitor_service.go              # Create/Update por tipo; RegenerateToken; summaries
│   ├── state_machine.go                # heartbeat: threshold efetivo = 1
│   └── notifier.go                     # + type, target
├── handler/
│   ├── ping.go                         # novo: GET/POST/HEAD /ping/{token}
│   ├── ratelimit.go                    # novo
│   └── monitors.go                     # + type/config, /token, resumo na listagem
└── repository/monitor_repo*.go         # colunas novas + FindByPingToken + query de resumo
cmd/pingou/commands/add.go              # --type --host --port --grace
web/src/
├── pages/Dashboard.jsx                 # faixa de situação, tabela, filtros, menu ⋮
├── components/monitors/
│   ├── MonitorForm.jsx                 # tipo, chips, avançado, bloco de ping URL
│   ├── MonitorRow.jsx                  # substitui MonitorCard
│   ├── StatusBanner.jsx, HistoryBar.jsx, RowMenu.jsx, ChipGroup.jsx
└── components/layout/Shell.jsx         # contador de incidentes abertos
```

## Task Breakdown

| Task ID | Agent                 | Name                                                           | Priority | Complexity | Dependencies |
| ------- | --------------------- | -------------------------------------------------------------- | -------- | ---------- | ------------ |
| T1      | design-db             | Migration 0003 (recriar `monitors`)                            | P0       | M          | -            |
| T2      | go-backend-specialist | Domain + repository com `type/config/ping_token/last_ping_at`  | P0       | M          | T1           |
| T3      | go-backend-specialist | Validação por tipo e service Create/Update                     | P0       | M          | T2           |
| T4      | go-backend-specialist | `checker.Registry` (despacho por tipo)                         | P1       | S          | T2           |
| T5      | code-test             | Teste de regressão HTTP (Fase 1)                               | P1       | S          | T3, T4       |
| T6      | go-backend-specialist | `TCPChecker`                                                   | P1       | S          | T4           |
| T7      | go-backend-specialist | API/CLI aceitam TCP                                            | P1       | S          | T3, T6       |
| T8      | code-test             | Testes TCP                                                     | P1       | S          | T6           |
| T9      | go-backend-specialist | Token: geração, regeneração, `FindByPingToken`                 | P1       | S          | T3           |
| T10     | go-backend-specialist | `Scheduler.Record` + state machine (threshold 1 p/ heartbeat)  | P1       | M          | T4           |
| T11     | go-backend-specialist | Rota `/ping/{token}` + rate limit                              | P1       | M          | T9, T10      |
| T12     | go-backend-specialist | Avaliador de prazos do heartbeat                               | P1       | M          | T10          |
| T13     | go-backend-specialist | Webhook com `type`/`target`                                    | P1       | S          | T3           |
| T14     | go-backend-specialist | CLI heartbeat + `POST /api/monitors/:id/token`                 | P1       | S          | T9           |
| T15     | code-test             | Testes heartbeat (ping, prazo, pausa, 404, 429)                | P1       | M          | T11, T12     |
| T16     | security-review       | Revisão da rota pública e do token                             | P1       | S          | T11, T13     |
| T17     | go-backend-specialist | `/api/monitors` com history/uptime/latência                    | P2       | M          | T2           |
| T18     | frontend-specialist   | Componentes base: ChipGroup, HistoryBar, RowMenu, StatusBanner | P2       | M          | T17          |
| T19     | frontend-specialist   | Painel (faixa, tabela, filtros, menu ⋮, ordem por gravidade)   | P2       | L          | T18          |
| T20     | frontend-specialist   | Novo monitor (tipo, chips, avançado, frase)                    | P2       | L          | T18, T7, T14 |
| T21     | frontend-specialist   | Editar monitor (ping URL, cron, regenerar, excluir)            | P2       | M          | T20          |
| T22     | frontend-specialist   | Sidebar: contador de incidentes abertos                        | P2       | XS         | T19          |
| T23     | qa-review             | QA contra os 3 PNGs                                            | P3       | M          | T19–T22      |
| T24     | document              | README/DESIGN: remover "(em desenvolvimento)"                  | P3       | XS         | T23          |

### Task T1: Migration 0003 (recriar `monitors`)

- **Agent**: design-db
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Backup de `pingou.db` feito (`cp pingou.db /tmp/...` ou `pingou export-db`)
- **INPUT**: `monitors.url NOT NULL CHECK(1..2048)`, `interval_seconds CHECK(10..86400)`. `checks` e `incidents` têm `ON DELETE CASCADE` para `monitors`.
- **OUTPUT**:
  1. `0003_monitor_types.sql` com `-- +goose NO TRANSACTION`, `PRAGMA foreign_keys=OFF`, criar `monitors_new`, copiar dados, `DROP monitors`, `RENAME`, recriar índices, `PRAGMA foreign_keys=ON`.
  2. Colunas novas: `type TEXT NOT NULL DEFAULT 'http' CHECK (type IN ('http','tcp','heartbeat'))`, `config TEXT NOT NULL DEFAULT '{}'`, `ping_token TEXT UNIQUE` (nulo fora de heartbeat), `last_ping_at TEXT`; `url` passa a nulo.
  3. `CHECK` de coerência: `(type='http' AND url IS NOT NULL) OR type<>'http'`; intervalo até **604800** (heartbeat de 1 semana).
  4. `Down` que recria o formato antigo (falha explícita se houver monitores não-http, ou converte/descarta; decidir e documentar no arquivo).
- **VERIFY**:
  - Copiar o `pingou.db` real, rodar `pingou migrate up`: contagem de `checks` e `incidents` **igual** antes/depois
  - `PRAGMA foreign_key_check;` vazio; `migrate down` e `up` de novo funcionam
- **Rollback**: restaurar o backup do `.db`.

### Task T2: Domain + repository

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: T1 aplicada
- **INPUT**: `domain.Monitor` sem tipo; `monitor_repo.go`/`monitor_repo_tx.go` com SELECT/INSERT/UPDATE de colunas fixas
- **OUTPUT**:
  1. `Monitor` ganha `Type`, `Config` (struct tipada: `Host`, `Port`, `GraceSeconds`; serializada em JSON), `PingToken *string`, `LastPingAt *time.Time`, `URL` como string vazia fora de HTTP.
  2. Repos leem/gravam as colunas novas; `FindByPingToken(ctx, token)`.
  3. `UpdateWithTx` grava `last_ping_at` (ver risco de sobrescrita).
- **VERIFY**: `go build ./... && go test ./...` passam; monitores antigos carregam com `Type=http`.

### Task T3: Validação por tipo e service

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T2
- **INPUT**: `validateCreateInput` exige URL; `Create/Update` só conhecem campos HTTP
- **OUTPUT**:
  1. `http`: igual a hoje. `tcp`: `host` não vazio (sem esquema), `port` 1–65535. `heartbeat`: `interval` 60–604800, `grace_seconds` 0–86400 (padrão 60), `timeout`/`threshold` ignorados (gravar 10 e 1).
  2. `type` ausente → `http`. `type` não pode mudar no `Update` (o design diz "o tipo não pode ser alterado").
  3. Heartbeat: `Create` gera o token (T9).
- **VERIFY**: testes em `validate_test.go` cobrindo cada tipo (válidos e inválidos, porta 0 e 65536, host com `http://`).

### Task T4: `checker.Registry`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **INPUT**: `serve.go` injeta `NewHTTPChecker` direto no scheduler
- **OUTPUT**:
  1. `Registry` implementa `domain.Checker` e escolhe por `m.Type`; tipo desconhecido → `CheckResult` com erro.
  2. `serve.go` injeta o `Registry`. Scheduler **não muda de assinatura**.
  3. `Scheduler.startJob` ignora `heartbeat` (não há loop ativo para ele).
- **VERIFY**: `make build`; monitor HTTP real continua gerando checks

### Task T5: Regressão HTTP (fecha o PR 1)

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3, T4
- **OUTPUT**: teste com `httptest.Server` (200 → UP, 500 → falha, threshold) passando pelo `Registry`; `scheduler_test.go` e `state_machine_test.go` existentes verdes
- **VERIFY**: `make test`. **Critério de merge do PR 1: nenhum comportamento HTTP mudou.**

### Task T6: `TCPChecker`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **OUTPUT**:
  1. `net.Dialer{}.DialContext(ctx, "tcp", net.JoinHostPort(host, port))`, fecha a conexão, mede latência até conectar.
  2. Erros traduzidos para `last_error`: `connection refused`, `timeout`, `no such host`.
  3. Registrado no `Registry`.
- **VERIFY**: ver T8

### Task T7: API/CLI aceitam TCP

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3, T6
- **OUTPUT**: `handler/monitors.go` lê/escreve `type` e `config`; `add.go` ganha `--type`, `--host`, `--port`
- **VERIFY**: `curl` do README (TCP) cria o monitor e o `Postgres` fica `UP`/`DOWN` conforme a porta

### Task T8: Testes TCP (fecha o PR 2)

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T6
- **OUTPUT**: `net.Listen("tcp","127.0.0.1:0")` → UP; porta fechada → DOWN com "refused"; timeout com IP não roteável (ou listener que não aceita) → DOWN
- **VERIFY**: `make test`

### Task T9: Token de ping

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3
- **OUTPUT**: `newPingToken()` (32 bytes `crypto/rand`, base64url); `MonitorService.RegenerateToken(id)`; só `heartbeat` aceita
- **VERIFY**: tokens únicos; regenerar em monitor HTTP retorna erro de validação

### Task T10: `Scheduler.Record` + state machine

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T4
- **INPUT**: `runCheck` mistura "rodar check" e "gravar com `dbMu`"; `handleFailure` usa `m.FailureThreshold`
- **OUTPUT**:
  1. `Record(ctx, monitorID, result, pingedAt *time.Time)` extrai a parte de banco de `runCheck`: pega `dbMu`, **recarrega o monitor do banco** e chama `sm.Process`. `runCheck` passa a usá-la.
  2. `state_machine`: para `heartbeat` o limite efetivo é 1 (`DOWN` direto).
  3. Quando é ping, setar `m.LastPingAt` antes de `Process`.
- **VERIFY**: testes existentes verdes; novo teste: heartbeat com 1 check `DOWN` já abre incidente

### Task T11: Rota `/ping/{token}` + rate limit

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T9, T10
- **OUTPUT**:
  1. `ping.go`: `GET|POST|HEAD /ping/{token}` **fora** do middleware de API Key. Token inexistente → 404 sem corpo útil. Monitor pausado → 200 sem registrar. Corpo da requisição ignorado (`http.MaxBytesReader` pequeno). Válido → `Record(UP)`.
  2. `ratelimit.go`: janela fixa por IP (ex.: 60 req/min), 429 com `Retry-After`; limpeza periódica do map. Usar `RemoteAddr`; só confiar em `X-Forwarded-For` se vier de configuração explícita (não adicionar agora).
  3. Também limita 404 repetidos (anti-varredura de token) pelo mesmo limiter.
- **VERIFY**: `curl -i /ping/<token>` 200; `/ping/xxx` 404; rajada >limite → 429; monitor pausado 200 e `checks` não cresce

### Task T12: Avaliador de prazos

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T10
- **OUTPUT**:
  1. `heartbeat.go`: ticker (10 s) que lista monitores `heartbeat` ativos.
  2. Vence quando `now > last_ping_at + interval + grace` **e** `current_state != DOWN`. Gera **um** `CheckResult` falho ("no ping received within interval + grace") via `Record`.
  3. Sem `last_ping_at` → continua `UNKNOWN` (não gera DOWN). Decisão registrada em Notes.
  4. Iniciado/parado junto com o scheduler (`Start`/`Stop`).
- **VERIFY**: teste com relógio/`intervalUnit` reduzido: sem ping → DOWN uma vez; com ping → UP e incidente fechado; segundo tick em DOWN não cria novo check

### Task T13: Webhook com `type`/`target`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3
- **OUTPUT**: `webhookPayloadMonitor` ganha `Type`, `Target *string`, `URL *string` (`omitempty`/`null` fora de HTTP). `target`: URL (http), `host:porta` (tcp), `null` (heartbeat). Token fora do payload.
- **VERIFY**: teste do notifier com os 3 tipos, conferindo que `ping_token` não aparece no JSON

### Task T14: CLI heartbeat e endpoint de token

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T9
- **OUTPUT**: `add --type heartbeat -i 86400 --grace 3600` imprime o caminho `/ping/<token>`; `POST /api/monitors/:id/token` devolve o monitor com o novo token
- **VERIFY**: token antigo passa a retornar 404 imediatamente

### Task T15: Testes heartbeat (fecha o PR 3)

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T11, T12
- **OUTPUT**: handler (200/404/429/pausado), avaliador (T12), regeneração (T14), `UNKNOWN → UP → DOWN → UP`
- **VERIFY**: `make test` verde

### Task T16: Revisão de segurança

- **Agent**: security-review
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T11, T13
- **OUTPUT**: checklist: token nunca em logs (o `loggingMiddleware` loga o path de `/ping/{token}`! redigir), nunca em webhook, limite por IP, 404 uniforme, comparação do token via índice único (sem timing relevante)
- **VERIFY**: grep nos logs após um ping; token não aparece

### Task T17: `/api/monitors` com resumo

- **Agent**: go-backend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T2
- **OUTPUT**:
  1. Listagem devolve por monitor: `history` (últimos 30 `success`, mais recente por último), `uptime_24h` (%), `last_latency_ms`. Em 2–3 queries para todos os monitores (`ROW_NUMBER() OVER (PARTITION BY monitor_id ...)`), não 1 por monitor.
  2. `ping_token` só no `GET /api/monitors/:id` e na listagem do tipo heartbeat (necessário ao painel); API já é protegida por API Key.
  3. Pausado: `history` continua mostrando o passado.
- **VERIFY**: `curl /api/monitors` traz os campos; log de queries mostra número constante com 7 monitores

### Task T18: Componentes base

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T17
- **OUTPUT**: `ChipGroup` (valor pronto + "Outro…" com input numérico e unidade), `HistoryBar` (30 barras verde/vermelho/cinza), `RowMenu` (Editar, Pausar/Retomar, Ver histórico, Excluir em vermelho), `StatusBanner`
- **VERIFY**: `npm run build`; cores só dos tokens `pg-*` do DESIGN.md

### Task T19: Painel

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T18
- **OUTPUT**:
  1. Faixa de situação no topo (substitui os 3 cards): frase + contagens (no ar, fora do ar, desconhecido, pausado). Verde/neutra quando tudo certo.
  2. Tabela: Status, Monitor (nome + alvo com ícone por tipo; heartbeat mostra "Espera um ping a cada 1 dia"), Últimos 30 checks, Uptime 24 h, Resposta, Último check ("há 4 s").
  3. Ordem por gravidade (DOWN, UNKNOWN, UP, pausado). Filtros Todos/Fora do ar/No ar/Pausados com contagem. Busca só aparece com ≥10 monitores.
  4. Linha pausada esmaecida. "Atualizado há 12 s" + botão de refresh.
- **VERIFY**: comparar com [Painel.png](concept-design/Painel.png) lado a lado; `Excluir` não aparece fora do menu ⋮

### Task T20: Novo monitor

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T18, T7, T14
- **OUTPUT**:
  1. Seletor de tipo (3 cartões com frase) primeiro; depois Nome; campos por tipo (HTTP: URL; TCP: host + porta; Heartbeat: "Esperar um ping a cada" 1 h/6 h/1 dia/1 semana/Outro e "Tolerância" 1 min/15 min/1 h/Outro).
  2. HTTP/TCP: "Verificar a cada" 30 s/1 min/5 min/15 min/1 h/Outro; bloco **Avançado** (timeout, falhas até cair) com a frase "Fica fora do ar depois de N falhas seguidas, ou seja, em cerca de X min".
  3. Frase do heartbeat: "Se nenhum ping chegar em 1 dia e 1 h, o monitor fica fora do ar."
  4. Sem campo "Ativado". Botão "Criar monitor".
- **VERIFY**: criar um monitor de cada tipo pela UI; validação de erro da API aparece no formulário

### Task T21: Editar monitor

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T20
- **OUTPUT**: mesmo form com tipo travado ("Heartbeat · o tipo não pode ser alterado"); heartbeat mostra bloco de URL de ping (`window.location.origin + /ping/<token>`), botão **Copiar**, exemplo de cron com a URL, "Último ping recebido há 9 h", **Gerar nova URL** (com confirmação, pois quebra o cron atual); rodapé com Excluir (vermelho, com confirmação), Cancelar, Salvar
- **VERIFY**: comparar com [Editar-monitor.png](concept-design/Editar-monitor.png); regenerar URL e confirmar que a antiga dá 404

### Task T22: Contador de incidentes na sidebar

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T19
- **OUTPUT**: badge vermelho no item "Incidentes" com a contagem de abertos (`GET /api/incidents?open=true`, usar o `total`); some quando 0
- **VERIFY**: derrubar um monitor → contador sobe

### Task T23: QA contra o design

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T19–T22
- **OUTPUT**: lista de divergências vs. os 3 PNGs; casos de borda (0 monitores, 1 monitor, heartbeat nunca pingado, monitor com nome longo, "Outro…" com valor inválido, URL de ping copiada em HTTP sem HTTPS)
- **VERIFY**: divergências corrigidas ou registradas em Notes

### Task T24: Documentação

- **Agent**: document
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T23
- **OUTPUT**: remover "(em desenvolvimento)" do README, marcar o Roadmap, rodar `graphify update .`
- **VERIFY**: `grep "em desenvolvimento" README.md` vazio

## Risk Areas

| Risk                                                                                                                              | Impact                                                 | Mitigation                                                                                               |
| --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------ | -------------------------------------------------------------------------------------------------------- |
| `DROP TABLE monitors` com FK ligado apaga `checks` e `incidents` via `CASCADE`                                                    | **Perda total de histórico**                           | `NO TRANSACTION` + `PRAGMA foreign_keys=OFF` antes do DROP; backup; conferir contagens antes/depois (T1) |
| `PRAGMA foreign_keys` não tem efeito dentro de transação (goose abre uma por padrão)                                              | Mesmo risco acima, silencioso                          | `-- +goose NO TRANSACTION` obrigatório; teste com `.db` real                                             |
| `UpdateWithTx` grava o monitor inteiro: ping chega enquanto outra rotina segura um `*Monitor` antigo e sobrescreve `last_ping_at` | Heartbeat cai sem motivo                               | `Record` recarrega o monitor dentro do `dbMu` (T10)                                                      |
| Avaliador gera check `DOWN` a cada tick                                                                                           | Milhares de checks, `countConsecutiveFails` distorcido | Só gera quando `current_state != DOWN` (T12)                                                             |
| Rota pública sem API Key                                                                                                          | Abuso/varredura de token, escrita no SQLite            | Rate limit por IP, 404 uniforme, token de 256 bits (T11, T16)                                            |
| Token no log de acesso (`loggingMiddleware`)                                                                                      | Vazamento do segredo                                   | Redigir `/ping/*` no log (T16)                                                                           |
| Heartbeat `UNKNOWN` para sempre se nunca pingar                                                                                   | Usuário não percebe que o cron nunca rodou             | Decisão: manter `UNKNOWN` (README). Banner mostra "desconhecido". Revisitar se virar queixa              |
| Intervalo de 1 semana (604800) vs `CHECK` atual (86400)                                                                           | Design não pode ser criado                             | Ampliar o `CHECK` na migration; `interval` de HTTP/TCP continua limitado a 86400 na validação            |
| `maxMonitors=100` e `FindAll(Limit:100)` no `Start`                                                                               | Teto conhecido                                         | Fora do escopo, mantido                                                                                  |
| Proxy reverso muda o IP visto pelo rate limit                                                                                     | Todos compartilham 1 IP e tomam 429                    | Limite generoso; suporte a `X-Forwarded-For` fica como melhoria futura, documentada                      |

## Rollback Strategy

1. Cada fase é um PR separado: reverter o PR da fase.
2. Banco: restaurar o backup feito antes da T1 (o `Down` da migration é melhor esforço, não garantia).
3. Se só o frontend falhar, o backend continua compatível (campos novos são aditivos; `type` ausente = `http`).
4. Heartbeat com problema: desabilitar registrando `/ping` e o avaliador num único ponto (`serve.go`), sem tocar HTTP/TCP.

## Phase X: Verification

```bash
# Go
gofumpt -l . && go vet ./... && go build ./...
make test

# Frontend
cd web && npm run lint && npm run build && cd ..

# Build completo + migration em cópia do banco real
make build
cp pingou.db /tmp/pingou-copia.db
sqlite3 /tmp/pingou-copia.db "select count(*) from checks; select count(*) from incidents;"
PINGOU_DATABASE_URL=/tmp/pingou-copia.db ./bin/pingou migrate up --key dev-api-key
sqlite3 /tmp/pingou-copia.db "select count(*) from checks; select count(*) from incidents; pragma foreign_key_check;"

# Smoke dos três tipos
make dev   # e criar um monitor de cada tipo; curl no /ping/<token>
```

## Cronograma Sugerido

| Fase                               | Tarefas | Complexity | Esforço | Target                                            | Status      |
| ---------------------------------- | ------- | ---------- | ------- | ------------------------------------------------- | ----------- |
| Fase 1 — Base + refatoração (PR 1) | T1–T5   | L          | ~2 h    | HTTP idêntico, schema novo                        | ⬜ Pendente |
| Fase 2 — TCP (PR 2)                | T6–T8   | M          | ~1 h    | TCP via API/CLI                                   | ⬜ Pendente |
| Fase 3 — Heartbeat (PR 3)          | T9–T16  | XL         | ~3 h    | Ping, prazo, webhook, segurança                   | ⬜ Pendente |
| Fase 4 — API de resumo + Frontend  | T17–T22 | XL         | ~4 h    | Telas dos 3 PNGs (pode iniciar T17 após a Fase 1) | ⬜ Pendente |
| Fase 5 — QA e docs                 | T23–T24 | M          | ~1 h    | Divergências fechadas                             | ⬜ Pendente |

## Status de Execução

### Concluído

- Levantamento do código e escopo aprovado no README

### Em andamento

-

### Pendente

- T1–T24

## Notes

- O README já descreve o comportamento final e marca o que falta com "(em desenvolvimento)". Ele é a especificação; se algo aqui divergir, o README vence.
- **Decisão assumida:** heartbeat nunca pingado fica `UNKNOWN` (não vira `DOWN`). O prazo só começa a contar a partir do primeiro ping.
- **Decisão assumida:** a URL de ping é montada no frontend com `window.location.origin`, sem nova variável de ambiente. O CLI imprime só o caminho `/ping/<token>`.
- **Decisão assumida:** busca no painel só aparece com 10+ monitores (o design mostra a caixa; a conversa de design aprovou esse comportamento).
- **Aberto:** o item "Ver histórico" do menu ⋮ não tem tela no design. Aponta para a tela de detalhe de [tarefa-detalhe-monitor-checks-incidentes.md](tasks/tarefa-detalhe-monitor-checks-incidentes.md) se já existir; senão, ocultar o item até lá.
- **Limitação conhecida:** uptime 24 h de heartbeat é calculado sobre checks gravados (1 por ping + 1 por queda), não sobre tempo. Aceitável para o painel; documentar se confundir.
- Ao terminar cada fase, rodar `graphify update .` (regra do CLAUDE.md do projeto).

### Alternativas consideradas e descartadas

- `Guardar URL fictícia ("heartbeat://id") para evitar recriar a tabela` → descartada: polui dados e a validação, e o `CHECK` de intervalo continuaria barrando 1 semana.
- `Uma tabela por tipo (tcp_monitors, heartbeat_monitors)` → descartada: duplica repositórios, scheduler, listagem e filtros; `config` JSON resolve com uma tabela.
- `Avaliador do heartbeat como um job por monitor no scheduler` → descartada: um único ticker varrendo é mais simples e não cria goroutines por monitor.
- `Heartbeat respeitando failure_threshold` → descartada: tolerância (`grace_seconds`) já cumpre esse papel (decisão no README).
- `Biblioteca de rate limit (x/time/rate, tollbooth)` → descartada: dependência nova para ~30 linhas; limite em memória basta para um binário só.
- `Websocket para atualizar o painel` → fora deste plano (tarefa própria em `tasks/`).
