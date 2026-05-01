# 🏗️ MVP Plan: Pingou

> **"Rodou, Pingou"** — Health checker self-hosted, leve e open-source em Go.

---

## Overview

Pingou é um monitor de uptime self-hosted em Go que executa checks HTTP periódicos em URLs configuradas, persiste resultados em SQLite, e dispara webhooks em transições de estado (UP↔DOWN). Distribuído como binário único (com UI React embutida via `embed.FS`) e via Docker.

Este projeto é também um **veículo de aprendizado de Go** para um dev senior Node.js. Cada fase introduz conceitos novos de forma incremental.

---

## Success Criteria

- [ ] `docker compose up` sobe o sistema em < 30s
- [ ] Suporta 100 URLs com checagem a cada 10s consumindo < 50MB RAM
- [ ] Webhook dispara corretamente em transições UP↔DOWN respeitando threshold
- [ ] UI React embutida no binário Go acessível em `http://localhost:8080`
- [ ] Roda 7 dias ininterruptos sem leak de memória ou goroutines
- [ ] Outro dev clona o repo e sobe em < 10 minutos via README
- [ ] Export de dump do SQLite funciona via API e CLI

---

## Tech Stack

| Layer       | Technology                                                        | Rationale                                            |
| ----------- | ----------------------------------------------------------------- | ---------------------------------------------------- |
| Linguagem   | Go 1.25+                                                          | Concorrência nativa, binário único, aprendizado      |
| HTTP Router | `net/http` (stdlib)                                               | Zero deps, routing nativo com path params (Go 1.22+) |
| DB          | SQLite via `github.com/mattn/go-sqlite3`                          | Driver SQLite maduro; requer CGO no build            |
| Migrations  | `goose`                                                           | Padrão de mercado, simples, embedável                |
| Scheduler   | `time.Ticker` + goroutines                                        | Idiomático Go, sem libs                              |
| Logs        | `log/slog` (stdlib)                                               | Estruturado, zero deps                               |
| Config      | env vars + `godotenv`                                             | KISS                                                 |
| Validação   | Manual no service; futura migração para `go-playground/validator` | KISS agora; validação declarativa depois             |
| HTTP Client | `net/http` (stdlib)                                               | Idiomático                                           |
| UI          | React 18 + Vite + TypeScript + Tailwind                           | Familiaridade do dev                                 |
| UI bundling | `embed.FS` (stdlib)                                               | 1 binário único                                      |
| Container   | Docker multi-stage + docker-compose                               | Solicitado                                           |

---

## File Structure

```
pingou/
├── cmd/
│   └── pingou/
│       └── main.go                    # Entry point: parse env, wire deps, start server
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Carrega e valida env vars
│   │
│   ├── database/
│   │   ├── database.go                # Abre SQLite, configura WAL, pool e embute migrations
│   │   └── migrations/
│   │       └── 0001_init.sql          # monitors, checks, incidents, settings
│   │
│   ├── monitors/                      # Domínio: URLs monitoradas
│   │   ├── model.go                   # Struct Monitor
│   │   ├── repository.go              # CRUD no SQLite
│   │   ├── service.go                 # Regras de negócio
│   │   └── handler.go                 # HTTP handlers /api/monitors
│   │
│   ├── checks/                        # Domínio: execução de checagens
│   │   ├── model.go                   # Struct Check, CheckResult
│   │   ├── repository.go              # Insert + queries de histórico
│   │   ├── checker.go                 # Executa HTTP GET, retorna resultado
│   │   ├── scheduler.go               # Orquestra goroutines por monitor
│   │   ├── retention.go               # Job de limpeza
│   │   └── handler.go                 # GET /api/monitors/:id/checks
│   │
│   ├── incidents/                     # Domínio: state machine UP/DOWN/UNKNOWN
│   │   ├── model.go                   # Struct Incident, State
│   │   ├── repository.go              # CRUD incidents
│   │   ├── state_machine.go           # Lógica de transições
│   │   └── handler.go                 # GET /api/incidents
│   │
│   ├── notifications/                 # Domínio: webhook
│   │   ├── model.go                   # Payload do webhook
│   │   ├── webhook.go                 # Client HTTP com retry
│   │   └── dispatcher.go              # Recebe events, envia em goroutine
│   │
│   ├── settings/                      # Domínio: configs runtime (webhook URL, retention)
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── handler/
│   │   ├── server.go                  # http.Server, rotas /api/* e SPA
│   │   ├── middleware.go              # Logging com request ID, auth, recover/CORS futuros
│   │   └── response.go                # Helpers JSON/erro
│   │
│   └── export/
│       └── dump.go                    # Gera dump do SQLite
│
├── web/                               # Frontend React
│   ├── src/
│   │   ├── pages/
│   │   ├── components/
│   │   ├── api/
│   │   └── App.tsx
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
│
├── ui/
│   └── embed.go                       # //go:embed dist  → embute build do React
│
├── scripts/
│   └── build.sh                       # Build completo: web → embed → go build
│
├── .env.example
├── .gitignore
├── .dockerignore
├── Dockerfile                         # Multi-stage: node + go + alpine
├── docker-compose.yml
├── Makefile                           # make dev, make build, make test, make docker
├── go.mod
├── go.sum
├── README.md
├── LICENSE                            # MIT ou Apache 2.0
└── pingou-mvp.md                      # Este arquivo
```

---

## Dependency Graph

```
Fase 1 (Setup)
    ↓
Fase 2 (Database + Migrations)
    ↓
Fase 3 (Domínio Monitors: model + repo + service)
    ↓
    ├─→ Fase 4 (HTTP Server + API Monitors)
    │       ↓
    │   Fase 5 (Auth middleware)
    │
    └─→ Fase 6 (Domínio Checks: checker + scheduler)
            ↓
        Fase 7 (Domínio Incidents: state machine)
            ↓
        Fase 8 (Notifications: webhook dispatcher)
            ↓
        Fase 9 (Settings + Retention + Export)
            ↓
        Fase 10 (Frontend React + embed.FS)
            ↓
        Fase 11 (Docker + docker-compose)
            ↓
        Fase X (Verificação final)
```

---

## Fases e Subetapas

### 🔧 Fase 1 — Setup do Projeto Go (P0)

**Objetivo de aprendizado:** Estrutura de projeto Go, `go mod`, organização `cmd/internal`, ferramentas básicas.

| #   | Subetapa                                                                                                                              | Output              | Verify                             |
| --- | ------------------------------------------------------------------------------------------------------------------------------------- | ------------------- | ---------------------------------- |
| 1.1 | Criar repo, `go mod init github.com/seu-user/pingou`                                                                                  | `go.mod`            | `go version` e `go mod tidy` rodam |
| 1.2 | Criar estrutura de diretórios vazia (com `.gitkeep`)                                                                                  | Árvore acima        | `tree internal/` mostra estrutura  |
| 1.3 | Criar `cmd/pingou/main.go` com "Hello, Pingou"                                                                                        | Binário compila     | `go run ./cmd/pingou` imprime msg  |
| 1.4 | Criar `internal/config/config.go` lendo env vars básicas (`PINGOU_PORT`, `PINGOU_API_KEY`, `PINGOU_DATABASE_URL`, `PINGOU_LOG_LEVEL`) | Struct `Config`     | Test unitário com env mockado      |
| 1.5 | Adicionar `godotenv` e `.env.example`                                                                                                 | Carrega .env em dev | `make dev` carrega vars            |
| 1.6 | Configurar `log/slog` JSON pra stdout                                                                                                 | Logger global       | Logs aparecem estruturados         |
| 1.7 | Criar `Makefile` com `make dev`, `make build`, `make test`                                                                            | Makefile funcional  | Cada target executa                |
| 1.8 | Adicionar `.gitignore`, `.dockerignore`, `.editorconfig`                                                                              | Arquivos            | Git ignora `bin/`, `*.db`, `.env`  |

**🎓 Conceitos novos:** módulos Go, `cmd/internal`, structs, env vars, slog, build flags.

---

### 🗄️ Fase 2 — Database & Migrations (P0)

**Objetivo de aprendizado:** `database/sql`, SQLite com WAL, migrations embutidas, pool de conexões.

| #   | Subetapa                                                                                                | Output               | Verify                          |
| --- | ------------------------------------------------------------------------------------------------------- | -------------------- | ------------------------------- |
| 2.1 | Adicionar dep `github.com/mattn/go-sqlite3` (driver SQLite via CGO)                                     | `go.sum` atualizado  | `go build` compila com CGO      |
| 2.2 | Adicionar dep `github.com/pressly/goose/v3`                                                             | Lib disponível       | Import funciona                 |
| 2.3 | Criar `internal/database/database.go`: abre DB, configura WAL, busy_timeout, max conns                  | `*sql.DB` retornado  | `db.Ping()` passa               |
| 2.4 | Criar migration `0001_init.sql` com tabelas: `monitors`, `checks`, `incidents`, `settings`              | Arquivo SQL          | Schema sintaticamente válido    |
| 2.5 | Embutir migrations SQL em `internal/database/database.go` com `//go:embed migrations/*.sql`             | FS embutido          | Build inclui SQL                |
| 2.6 | Integrar goose pra rodar migrations no startup                                                          | Migrations aplicadas | DB criado com tabelas           |
| 2.7 | Criar índices: `monitors(enabled)`, `checks(monitor_id, checked_at)`, `incidents(monitor_id, ended_at)` | Migration adicional  | `EXPLAIN QUERY PLAN` usa índice |
| 2.8 | Adicionar comando CLI `pingou migrate up/down/status` (subcomando)                                      | CLI funcional        | `pingou migrate status` lista   |

**🎓 Conceitos novos:** `database/sql`, drivers Go, `embed.FS`, WAL mode, migrations declarativas.

**⚠️ Pitfalls Node→Go:**

- Em Go, `*sql.DB` é **pool de conexões**, não conexão única. Não feche depois de cada query.
- SQLite + concurrent writes = use WAL (`PRAGMA journal_mode=WAL`) e limite writers a 1 (`SetMaxOpenConns(1)` pra writes, ou serialize via mutex).

---

### 📦 Fase 3 — Domínio Monitors (P0/P1)

**Objetivo de aprendizado:** Layered architecture, repository pattern, validação, enum em Go.

| #   | Subetapa                                                                                                                                     | Output                       | Verify                 |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- | ---------------------- |
| 3.1 | `domain/monitor.go`: struct `Monitor`, type `MonitorState`, const                                                                            | Tipos definidos              | Compila                |
| 3.2 | Validação manual no service: `interval >= 10 && <= 86400`, `failure_threshold >= 1 && <= 10`; URL por tamanho hoje, formato em tarefa futura | Validação funcional          | Input inválido falha   |
| 3.3 | `repository/monitor_repo.go`: implementação SQLite (Create/Get/List/Update/Delete)                                                           | CRUD completo                | Compila                |
| 3.4 | `service/monitor_service.go`: regras de negócio (validar antes de salvar, enforce limite de 100 monitors)                                    | Service layer                | Rejeita o 101º monitor |
| 3.5 | DTOs: `CreateMonitorDTO`, `UpdateMonitorDTO` separados do model                                                                              | Separação domínio/transporte | Compila                |
| 3.6 | `handler/monitors.go`: handlers HTTP para criação, listagem, atualização e remoção                                                           | Endpoints de monitors        | `go test ./...` passa  |

**🎓 Conceitos novos:** interfaces Go, dependency injection manual, enums via const+type, table-driven tests.

**⚠️ Pitfall Node→Go:** Em Go, **interfaces são satisfeitas implicitamente** (duck typing estrutural). Você não declara "implements Repository" — basta ter os métodos certos.

---

### 🌐 Fase 4 — HTTP Server + API Monitors (P1)

**Objetivo de aprendizado:** `net/http` ServeMux Go 1.22+ com path params, handlers, JSON encoding, graceful shutdown.

| #   | Subetapa                                                                                    | Output                  | Verify                      |
| --- | ------------------------------------------------------------------------------------------- | ----------------------- | --------------------------- |
| 4.1 | `handler/server.go`: `http.Server` com timeouts (read/write/idle) configurados              | Server estruturado      | Compila                     |
| 4.2 | Graceful shutdown via `signal.NotifyContext(ctx, SIGTERM, SIGINT)` e `http.Server.Shutdown` | Shutdown limpo          | Ctrl+C não corrompe DB      |
| 4.3 | `handler/middleware.go`: logging com request ID e latency; recover/CORS em tarefa futura    | Middlewares             | Logs mostram request        |
| 4.4 | `handler/monitors.go`: handlers para POST/GET/PATCH/DELETE `/api/monitors`                  | Endpoints REST          | curl funciona               |
| 4.5 | Helper `respondJSON(w, status, data)` e `respondError(w, status, msg)`                      | Response consistente    | Erros sempre em JSON        |
| 4.6 | Integration tests com `httptest.NewServer`                                                  | Tests E2E nos endpoints | `go test ./...` passa       |
| 4.7 | Endpoint `/healthz` (sem auth, retorna `{status:"ok"}`)                                     | Health endpoint         | `curl /healthz` retorna 200 |

**🎓 Conceitos novos:** `http.Handler`, `http.HandlerFunc`, ServeMux 1.22 (`POST /api/monitors/{id}`), middleware composition, `httptest`.

---

### 🔐 Fase 5 — Autenticação API Key (P2)

**Objetivo de aprendizado:** Middleware pattern, context, header validation.

| #   | Subetapa                                                              | Output           | Verify                  |
| --- | --------------------------------------------------------------------- | ---------------- | ----------------------- |
| 5.1 | `handler/middleware.go`: lê `X-API-Key` e compara com `config.APIKey` | Middleware       | Test rejeita key errada |
| 5.2 | Aplicar middleware em todas rotas `/api/*` exceto `/healthz`          | Rotas protegidas | curl sem header → 401   |
| 5.3 | Documentar header esperado no README                                  | Doc              | README atualizado       |

**🎓 Conceitos novos:** middleware decoration, composição com `http.Handler`.

---

### 🔄 Fase 6 — Domínio Checks: Engine de Checagem (P1)

**Objetivo de aprendizado:** Goroutines, channels, `context.Context`, `time.Ticker`, HTTP client com timeout.

| #   | Subetapa                                                                                                      | Output             | Verify                               |
| --- | ------------------------------------------------------------------------------------------------------------- | ------------------ | ------------------------------------ |
| 6.1 | `domain/check.go`: struct `Check`, `CheckResult`                                                              | Tipos              | Compila                              |
| 6.2 | `repository/check_repo.go`: `Create(check)`, `FindByMonitor(monitorID, limit, offset)`                        | Repo CRUD          | Compila                              |
| 6.3 | `checker/http_checker.go`: função `Check(ctx, monitor) CheckResult` — faz HTTP GET com timeout, mede latência | Checker puro       | Test futuro com `httptest.Server`    |
| 6.4 | `scheduler/scheduler.go`: struct `Scheduler` com jobs cancelados por `context.CancelFunc`                     | Estrutura          | Compila                              |
| 6.5 | `Scheduler.Start(ctx)`: lança 1 goroutine por monitor ACTIVE com `time.Ticker(monitor.Interval)`              | Scheduler ativo    | Logs mostram checks executando       |
| 6.6 | Integração: cada tick → `checker.Check` → `stateMachine.Process` → persiste check/transição                   | Pipeline funcional | DB acumula checks                    |
| 6.7 | `Scheduler.Reload(monitorID)`: para goroutine antiga e cria nova (chamado quando monitor é editado)           | Hot reload         | Editar interval reflete no scheduler |
| 6.8 | `Scheduler.Stop(monitorID)`: para goroutine quando monitor vira INACTIVE ou é deletado                        | Stop limpo         | Sem goroutine leak                   |
| 6.9 | Endpoint `GET /api/monitors/:id/checks?limit=N`                                                               | Histórico via API  | curl retorna lista                   |

**🎓 Conceitos novos:** goroutines, `context.Context` (cancelamento), `time.Ticker`, channels, `sync.RWMutex`, `http.Client` com `Timeout`.

**⚠️ Pitfalls críticos Node→Go:**

- **Sempre** passe `context.Context` em operações I/O. É o equivalente Go a `AbortController`.
- Goroutine sem cancelamento = leak. Toda goroutine precisa de "como morrer".
- `http.Client` default tem timeout = 0 (nunca expira). **Sempre** configure `Timeout`.
- Nunca compartilhe `time.Ticker` entre goroutines sem mutex.

---

### 📊 Fase 7 — State Machine de Incidentes (P1)

**Objetivo de aprendizado:** State machine integrada ao service, consistência de persistência e transições de incidentes.

| #   | Subetapa                                                                                                                                                                                                                              | Output                | Verify                         |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- | ------------------------------ |
| 7.1 | `domain/incident.go`: struct `Incident`, estados no `MonitorState` (UNKNOWN, UP, DOWN)                                                                                                                                                | Tipos                 | Compila                        |
| 7.2 | `repository/incident_repo.go`: `Create`, `Close`, `FindOpenByMonitor`, listagens                                                                                                                                                      | CRUD                  | Compila                        |
| 7.3 | `service/state_machine.go`: `Process(ctx, monitor, result)` persiste check, calcula transição e abre/fecha incidentes                                                                                                                 | Lógica integrada      | `go test ./...` passa          |
| 7.4 | Transições especificadas: <br>• UNKNOWN→UP: silencioso<br>• UNKNOWN→DOWN (após N falhas): notifica `down`<br>• UP→DOWN (após N falhas): notifica `down` + abre incident<br>• DOWN→UP (após 1 sucesso): notifica `up` + fecha incident | Comportamento correto | Logs mostram transições        |
| 7.5 | Integração no scheduler: após cada check → chamada direta para `stateMachine.Process`                                                                                                                                                 | Pipeline direto       | DB acumula checks e incidentes |
| 7.6 | Endpoint `GET /api/incidents` (listar) e `GET /api/monitors/:id/incidents`                                                                                                                                                            | Visibilidade          | curl retorna histórico         |

**🎓 Conceitos novos:** state machine em camada de service, transições de estado, consistência entre checks, monitors e incidents.

---

### 📢 Fase 8 — Notifications: Webhook Dispatcher (P1)

**Objetivo de aprendizado:** Worker pattern, retry com backoff, fan-out via channel.

| #   | Subetapa                                                                                                                                        | Output            | Verify                                   |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- | ---------------------------------------- |
| 8.1 | `notifications/model.go`: struct `WebhookPayload` com contrato real: `event`, `monitor`, `timestamp`, `last_error`, `downtime_duration_seconds` | Payload tipado    | Compila                                  |
| 8.2 | `notifications/webhook.go`: `Send(ctx, url, payload)` — POST JSON com timeout 10s                                                               | Client HTTP       | Test contra `httptest.Server`            |
| 8.3 | Retry: 3 tentativas, espera fixa 5s entre elas, `context.Done()` aborta                                                                         | Retry resiliente  | Test simula falha e sucesso              |
| 8.4 | `notifications/dispatcher.go`: goroutine consumindo channel de events, chama webhook em goroutine separada (não bloqueia state machine)         | Dispatcher async  | State machine não trava se webhook lento |
| 8.5 | Integração: state machine emite event → dispatcher recebe → busca webhook URL em settings → envia                                               | Pipeline completo | Webhook real recebe POST                 |
| 8.6 | Logs estruturados de envio (sucesso/falha + tentativa)                                                                                          | Observabilidade   | Logs mostram eventos                     |

Contrato real do payload:

```json
{
  "event": "down",
  "monitor": {
    "id": "018f2f7a-...",
    "name": "API",
    "url": "https://api.example.com/health"
  },
  "timestamp": "2026-05-01T12:00:00Z",
  "last_error": "unexpected status code: 500",
  "downtime_duration_seconds": null
}
```

```json
{
  "event": "up",
  "monitor": {
    "id": "018f2f7a-...",
    "name": "API",
    "url": "https://api.example.com/health"
  },
  "timestamp": "2026-05-01T12:05:00Z",
  "last_error": null,
  "downtime_duration_seconds": 300
}
```

**🎓 Conceitos novos:** producer/consumer com channel, fan-out, `time.After` vs `time.Sleep` (cancelável).

---

### ⚙️ Fase 9 — Settings, Retention & Export (P2)

**Objetivo de aprendizado:** Tabela key-value, jobs periódicos, streaming de arquivo via HTTP.

| #   | Subetapa                                                                                                                       | Output              | Verify                |
| --- | ------------------------------------------------------------------------------------------------------------------------------ | ------------------- | --------------------- |
| 9.1 | `settings/`: CRUD em tabela `settings(key, value)`. Keys: `webhook_url`, `retention_days`                                      | Settings funcionais | API GET/PUT funciona  |
| 9.2 | Validação: `retention_days >= 7 && <= 90`                                                                                      | Limite enforced     | API rejeita 5 ou 100  |
| 9.3 | `checks/retention.go`: goroutine que roda diariamente (`time.Ticker(24h)`) e deleta `checks` mais antigos que `retention_days` | Cleanup automático  | Logs mostram execução |
| 9.4 | `export/dump.go`: gera dump SQLite usando `VACUUM INTO 'tmp.db'` ou copia o arquivo com lock                                   | Export funcional    | Arquivo `.db` gerado  |
| 9.5 | Endpoint `GET /api/export` (auth required): retorna `application/octet-stream` com o dump                                      | Download via HTTP   | curl baixa o `.db`    |
| 9.6 | Comando CLI `pingou export --output backup.db`                                                                                 | CLI funcional       | Comando gera arquivo  |

**🎓 Conceitos novos:** jobs cron-like com `Ticker`, streaming HTTP response, `io.Copy`.

---

### 🎨 Fase 10 — Frontend React + Embed (P2)

**Objetivo de aprendizado:** `embed.FS`, servir SPA via Go, build pipeline integrado.

| #     | Subetapa                                                                                                                                                  | Output               | Verify                      |
| ----- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- | --------------------------- |
| 10.1  | `web/`: `npm create vite@latest -- --template react-ts`                                                                                                   | Projeto Vite         | `npm run dev` roda          |
| 10.2  | Setup Tailwind v4 + estrutura de pastas (`pages/`, `components/`, `api/`)                                                                                 | Stack pronta         | Hot reload funciona         |
| 10.3  | Página: lista de monitors com status badge (UP verde, DOWN vermelho, UNKNOWN cinza)                                                                       | UI funcional         | Renderiza mock              |
| 10.4  | Página: detalhe do monitor com últimos N checks + lista de incidentes                                                                                     | UI funcional         | Navega entre rotas          |
| 10.5  | Form: criar/editar monitor (validações client-side espelhando backend)                                                                                    | CRUD via UI          | Submit funciona             |
| 10.6  | Página: settings (webhook URL, retention days)                                                                                                            | UI configurável      | Salva                       |
| 10.7  | Página: export → botão que baixa o dump                                                                                                                   | Download             | Arquivo baixa               |
| 10.8  | Cliente API: passa `X-API-Key` (lida de localStorage; tela de login simples pede key na 1ª visita)                                                        | Auth funcional       | Requests autenticados       |
| 10.9  | Build de produção: `vite build` gera `web/dist/`                                                                                                          | Build estático       | `dist/index.html` existe    |
| 10.10 | `ui/embed.go`: `//go:embed all:dist` (assumindo cópia de `web/dist` pra `ui/dist`) + handler que serve assets com fallback pra `index.html` (SPA routing) | Embed funcional      | `go build` inclui assets    |
| 10.11 | Integrar no router: `/` e `/assets/*` servidos pelo handler do embed; `/api/*` pelos handlers Go                                                          | Tudo no mesmo server | UI carrega na porta 8080    |
| 10.12 | Script `scripts/build.sh`: builda React → copia pra `ui/dist` → builda Go                                                                                 | Build integrado      | `make build` gera 1 binário |

**🎓 Conceitos novos:** `embed.FS`, `http.FileServerFS`, SPA fallback routing em Go.

---

### 🐳 Fase 11 — Docker & Docker Compose (P2)

**Objetivo de aprendizado:** Multi-stage build (node + go + alpine), CGO em container, volumes, healthcheck.

| #    | Subetapa                                                                                                               | Output                         | Verify                         |
| ---- | ---------------------------------------------------------------------------------------------------------------------- | ------------------------------ | ------------------------------ |
| 11.1 | `Dockerfile` stage 1: `node:20-alpine` builda React                                                                    | Stage builda                   | Layer cacheado                 |
| 11.2 | `Dockerfile` stage 2: `golang:1.25-alpine` builda Go com `CGO_ENABLED=1`, `gcc`/`musl-dev`, copiando assets do stage 1 | Binário compila com SQLite CGO | Binário compila                |
| 11.3 | `Dockerfile` stage 3: `alpine` com binário + ca-certs                                                                  | Imagem final enxuta            | `docker images` mostra tamanho |
| 11.4 | `HEALTHCHECK` no Dockerfile chamando `/healthz`                                                                        | Healthcheck                    | `docker ps` mostra "healthy"   |
| 11.5 | `docker-compose.yml`: serviço `pingou`, volume nomeado `pingou-data:/data`, env vars, port 8080                        | Compose funcional              | `docker compose up` sobe       |
| 11.6 | `.dockerignore`: exclui `node_modules`, `bin`, `*.db`, `.env`, `.git`                                                  | Build rápido                   | Context pequeno                |
| 11.7 | Test: `docker compose up`, criar monitor, verificar persistência após `down`+`up`                                      | E2E funcional                  | Dados persistem                |

**🎓 Conceitos novos:** multi-stage builds, CGO em containers Alpine, volumes Docker.

---

### ✅ Fase X — Verificação Final & Documentação (P3)

| #   | Subetapa                                                                                       | Output                 | Verify                   |
| --- | ---------------------------------------------------------------------------------------------- | ---------------------- | ------------------------ |
| X.1 | README completo: o que é, instalação Docker, configuração env, exemplos curl, screenshot da UI | README didático        | Outro dev consegue subir |
| X.2 | `.env.example` com TODAS as vars documentadas e valores default                                | Docs                   | Cópia funciona           |
| X.3 | Documentar payload do webhook + exemplo de receiver (Discord/Slack) no README                  | Integração documentada | Webhook compreensível    |
| X.4 | LICENSE (MIT recomendado pra OSS)                                                              | Licença                | Arquivo presente         |
| X.5 | GitHub Action básica: `go test ./...` + `golangci-lint` em PR                                  | CI funcional           | PR mostra checks         |
| X.6 | Checklist de verificação manual (ver abaixo) executado                                         | Tudo verde             | Lista preenchida         |

---

## Phase X: Verification Checklist

### Build

- [ ] `make build` gera binário < 20MB
- [ ] `docker compose build` constrói imagem < 30MB
- [ ] `go test ./...` passa sem flakes
- [ ] `golangci-lint run` sem warnings

### Funcionalidade Core

- [ ] CRUD completo de monitors via API e UI
- [ ] `Enabled` controla execução do scheduler
- [ ] Editar interval reflete no scheduler sem reiniciar
- [ ] 100 URLs simultâneas rodando estável por 1h
- [ ] Histórico de checks visível (paginado se necessário)

### State Machine & Webhooks

- [ ] `UNKNOWN → UP` **não** dispara webhook
- [ ] `UNKNOWN → DOWN` (após N falhas) dispara webhook
- [ ] `UP → DOWN` (após N falhas) dispara webhook + abre incident
- [ ] `DOWN → UP` (após 1 sucesso) dispara webhook + fecha incident
- [ ] Payload do webhook contém todos campos esperados
- [ ] Retry de webhook funciona (testar com endpoint fake que falha)

### Auth & Segurança

- [ ] Rotas `/api/*` exigem `X-API-Key`
- [ ] `/healthz` é público
- [ ] Sem secrets hardcoded no código
- [ ] Rotas protegidas rejeitam API key inválida

### Persistência & Retenção

- [ ] SQLite usa WAL mode
- [ ] Retention job remove checks antigos
- [ ] Limites enforced: interval 10s–86400s, retention 7d–90d, máx 100 monitors
- [ ] Export gera dump válido (abre em outro SQLite)
- [ ] Volume Docker persiste após `down` + `up`

### Operação

- [ ] Graceful shutdown não corrompe DB
- [ ] Logs estruturados JSON em stdout
- [ ] Sem goroutine leak após 24h (`runtime.NumGoroutine()` estável)
- [ ] Consumo de RAM < 50MB com 100 monitors

### Documentação

- [ ] README permite setup em < 10min
- [ ] `.env.example` cobre todas vars
- [ ] Exemplo de webhook payload documentado
- [ ] Comandos CLI documentados (`migrate`, `export`)

---

## Decisões Arquiteturais Registradas

> **D1 — Embed React via `embed.FS`:** 1 binário único. Trade-off: build em 2 passos.
> **D2 — Limites rígidos como feature:** 10s–24h, 7d–90d, 100 monitors. Trade-off: perde power users; ganha foco.
> **D3 — `UNKNOWN→UP` silencioso:** Evita spam no startup. Apenas transições entre estados conhecidos são incidentes.
> **D4 — `net/http` puro (sem chi):** Aprende routing nativo. Migrar pra chi depois é trivial se necessário.
> **D5 — Layered por domínio:** Boundaries claros desde o dia 1. Custo de setup similar; ganho de manutenibilidade alto.
> **D6 — Webhook global apenas:** Simplicidade. Per-URL fica pra v1.1 se houver demanda.
> **D7 — `github.com/mattn/go-sqlite3` (CGO):** Driver SQLite maduro e amplamente usado. Trade-off: build exige CGO/toolchain C; Docker usa Alpine com `gcc`/`musl-dev`.
> **D8 — Scheduler via goroutine-per-monitor:** Simples, idiomático. Limite de 100 monitors mantém isso viável (~100 goroutines = trivial).

---

## Pós-MVP (não implementar agora)

- v1.1: Validação de body/latência configurável
- v1.2: TCP port check
- v1.3: Status page pública (sem auth, read-only)
- v1.4: Agregações (gráficos hora/dia)
- v1.5: Webhook per-URL + templates Discord/Slack
- v1.6: Métricas Prometheus em `/metrics`
- v2.0: Multi-tenancy

---

**Última atualização:** 2026-04-23
**Status:** Plano aprovado, aguardando início da Fase 1
