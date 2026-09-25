# Graph Report - pingou-health-checker  (2026-09-25)

## Corpus Check
- 108 files · ~79,763 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 9, .css 2, .example 1)

## Summary
- 491 nodes · 1396 edges · 29 communities (18 shown, 11 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `01338934`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- main.jsx
- database/sql.DB
- package.json
- config.go
- DESIGN.md — Design system do Pingou
- net/http.ResponseWriter
- middleware.go
- Monitor
- Incident
- context.Context
- UnitOfWork
- info.go
- validate.go
- CheckRepository
- boolToInt
- WebhookNotifier
- Divergências entre PRD e Implementação
- Vite Logo
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Pingou README
- Architecture Overview
- Technical Concerns & Risk Assessment
- Hero Image
- React Logo
- MonitorService

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 38 edges
2. `Incident` - 25 edges
3. `runServe()` - 21 edges
4. `Scheduler` - 18 edges
5. `MonitorService` - 16 edges
6. `tokens` - 14 edges
7. `Check` - 14 edges
8. `writeError()` - 14 edges
9. `DESIGN.md — Design system do Pingou` - 13 edges
10. `writeJSON()` - 13 edges

## Surprising Connections (you probably didn't know these)
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go
- `runServe()` --calls--> `NewHTTPChecker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/checker/http_checker.go
- `runServe()` --calls--> `NewWebhookNotifier()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/service/notifier.go
- `runMigrateDown()` --calls--> `LoadConfig()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/config/config.go
- `runMigrateStatus()` --calls--> `LoadConfig()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/config/config.go

## Import Cycles
- None detected.

## Communities (29 total, 11 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (22): acquireLock(), checkExistingLock(), MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors (+14 more)

### Community 1 - "main.jsx"
Cohesion: 0.09
Nodes (33): axios, react, base, client, incidentsApi, monitorsApi, settingsApi, nav (+25 more)

### Community 2 - "database/sql.DB"
Cohesion: 0.08
Nodes (37): runAdd(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), runRM(), releaseLock(), runServe(), go_pkg_github_com_mattn_go_sqlite3 (+29 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (41): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, react-dom (+33 more)

### Community 4 - "config.go"
Cohesion: 0.11
Nodes (18): go_pkg_strconv, log/slog.Level, net/http.ServeMux, net/http.Server, updateSettingsRequest, Config, Load(), LoadConfig() (+10 more)

### Community 5 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.10
Nodes (19): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+11 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.17
Nodes (18): net/http.Request, net/http.ResponseWriter, checkResponse, incidentResponse, monitorResponse, Server, toCheckResponse(), Server (+10 more)

### Community 7 - "middleware.go"
Cohesion: 0.14
Nodes (14): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, responseWriter (+6 more)

### Community 8 - "Monitor"
Cohesion: 0.15
Nodes (11): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, scanMonitor(), Notifier, NewStateMachine() (+3 more)

### Community 9 - "Incident"
Cohesion: 0.18
Nodes (6): Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo

### Community 10 - "context.Context"
Cohesion: 0.15
Nodes (7): context.Context, Check, MonitorRepository, CheckRepo, CheckRepositoryTx, checkRepoWithTx, incidentRepoWithTx

### Community 11 - "UnitOfWork"
Cohesion: 0.12
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 12 - "info.go"
Cohesion: 0.09
Nodes (20): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+12 more)

### Community 13 - "validate.go"
Cohesion: 0.53
Nodes (9): go_pkg_regexp, validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold(), validateTimeout(), validateUpdateInput() (+1 more)

### Community 14 - "CheckRepository"
Cohesion: 0.31
Nodes (3): CheckRepository, NewRetentionWorker(), RetentionWorker

### Community 15 - "boolToInt"
Cohesion: 0.27
Nodes (4): database/sql.Tx, boolToInt(), nullableTime(), monitorRepoWithTx

### Community 16 - "WebhookNotifier"
Cohesion: 0.24
Nodes (7): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 34 - "MonitorService"
Cohesion: 0.26
Nodes (5): MonitorService, newUUIDv7(), CreateMonitorInput, Reloader, UpdateMonitorInput

## Knowledge Gaps
- **68 isolated node(s):** `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse`, `updateMonitorRequest`, `base` (+63 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 108 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `database/sql.DB`, `MonitorService`, `net/http.ResponseWriter`, `context.Context`, `UnitOfWork`, `validate.go`, `boolToInt`, `WebhookNotifier`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Why does `runServe()` connect `database/sql.DB` to `serve.go`, `config.go`, `Incident`, `CheckRepository`, `WebhookNotifier`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `Server` connect `config.go` to `serve.go`, `database/sql.DB`, `MonitorService`, `middleware.go`, `Incident`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **What connects `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse` to the rest of the system?**
  _68 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08590441621294616 - nodes in this community are weakly interconnected._
- **Should `database/sql.DB` be split into smaller, more focused modules?**
  _Cohesion score 0.08408163265306122 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.0507399577167019 - nodes in this community are weakly interconnected._