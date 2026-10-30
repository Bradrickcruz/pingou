# Graph Report - pingou-health-checker  (2026-10-05)

## Corpus Check
- 111 files · ~91,996 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 9, .woff2 4, .example 1)

## Summary
- 508 nodes · 1434 edges · 29 communities (16 shown, 13 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `13608499`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- main.jsx
- database/sql.DB
- package.json
- Scheduler
- DESIGN.md — Design system do Pingou
- net/http.ResponseWriter
- middleware.go
- Monitor
- context.Context
- info.go
- Check
- boolToInt
- WebhookNotifier
- MonitorRepository
- UnitOfWork
- Server
- Divergências entre PRD e Implementação
- Vite Logo
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Pingou README
- Architecture Overview
- Technical Concerns & Risk Assessment
- Hero Image
- React Logo
- validate.go

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 40 edges
2. `Incident` - 27 edges
3. `runServe()` - 21 edges
4. `Scheduler` - 18 edges
5. `MonitorService` - 16 edges
6. `react` - 14 edges
7. `Check` - 14 edges
8. `writeError()` - 14 edges
9. `DESIGN.md — Design system do Pingou` - 13 edges
10. `writeJSON()` - 13 edges

## Surprising Connections (you probably didn't know these)
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go
- `runServe()` --calls--> `NewHTTPChecker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/checker/http_checker.go
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go
- `runServe()` --calls--> `NewIncidentService()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/service/incident_service.go

## Import Cycles
- None detected.

## Communities (29 total, 13 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.18
Nodes (25): acquireLock(), checkExistingLock(), MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors (+17 more)

### Community 1 - "main.jsx"
Cohesion: 0.07
Nodes (37): axios, lucide-react, react, react-router-dom, base, client, incidentsApi, monitorsApi (+29 more)

### Community 2 - "database/sql.DB"
Cohesion: 0.06
Nodes (56): CheckRepositoryTx, runAdd(), copyFile(), runExportDB(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), runRM() (+48 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (41): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, react-dom (+33 more)

### Community 4 - "Scheduler"
Cohesion: 0.26
Nodes (5): context.CancelFunc, sync.Mutex, Checker, job, Scheduler

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
Cohesion: 0.18
Nodes (11): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, scanMonitor(), UnitOfWork, NewStateMachine() (+3 more)

### Community 9 - "context.Context"
Cohesion: 0.14
Nodes (9): context.Context, Incident, IncidentRepository, scanIncident(), Notifier, IncidentRepo, IncidentRepositoryTx, incidentRepoWithTx (+1 more)

### Community 10 - "info.go"
Cohesion: 0.16
Nodes (11): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim() (+3 more)

### Community 11 - "Check"
Cohesion: 0.14
Nodes (5): database/sql.Tx, Check, CheckRepositoryTx, checkRepoWithTx, monitorRepoWithTx

### Community 12 - "boolToInt"
Cohesion: 0.22
Nodes (3): boolToInt(), nullableTime(), MonitorRepo

### Community 13 - "WebhookNotifier"
Cohesion: 0.38
Nodes (4): HTTPChecker, net/http.Client, NewHTTPChecker(), WebhookNotifier

### Community 16 - "Server"
Cohesion: 0.11
Nodes (16): net/http.ServeMux, net/http.Server, updateSettingsRequest, CheckRepository, Server, NewServer(), toUpdateSettingsInput(), SettingsRepo (+8 more)

### Community 34 - "validate.go"
Cohesion: 0.19
Nodes (14): go_pkg_regexp, MonitorService, newUUIDv7(), validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold() (+6 more)

## Knowledge Gaps
- **69 isolated node(s):** `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse`, `updateMonitorRequest`, `base` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 116 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `Scheduler`, `net/http.ResponseWriter`, `context.Context`, `Check`, `boolToInt`, `WebhookNotifier`, `MonitorRepository`, `UnitOfWork`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `sqliteUnitOfWork` connect `database/sql.DB` to `serve.go`, `Check`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Why does `Incident` connect `context.Context` to `serve.go`, `net/http.ResponseWriter`, `Monitor`, `boolToInt`, `WebhookNotifier`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06980433632998413 - nodes in this community are weakly interconnected._
- **Should `database/sql.DB` be split into smaller, more focused modules?**
  _Cohesion score 0.05593561368209256 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.0507399577167019 - nodes in this community are weakly interconnected._