# Graph Report - pingou-health-checker  (2026-09-26)

## Corpus Check
- 109 files · ~83,371 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 9, .css 2, .example 1)

## Summary
- 494 nodes · 1407 edges · 29 communities (17 shown, 12 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `21c17506`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- main.jsx
- runServe
- package.json
- database/sql.DB
- DESIGN.md — Design system do Pingou
- net/http.ResponseWriter
- middleware.go
- .handleFailure
- Incident
- context.Context
- UnitOfWork
- Check
- Monitor
- MonitorService
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
- validate.go

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
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go

## Import Cycles
- None detected.

## Communities (29 total, 12 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (22): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_google_uuid (+14 more)

### Community 1 - "main.jsx"
Cohesion: 0.08
Nodes (36): axios, lucide-react, react, react-dom, base, client, incidentsApi, monitorsApi (+28 more)

### Community 2 - "runServe"
Cohesion: 0.06
Nodes (50): runAdd(), copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion() (+42 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (41): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, react-router-dom (+33 more)

### Community 4 - "database/sql.DB"
Cohesion: 0.14
Nodes (15): context.CancelFunc, database/sql.DB, sync.Mutex, Checker, CheckRepoTx, NewCheckRepoTx(), IncidentRepoTx, NewIncidentRepoTx() (+7 more)

### Community 5 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.10
Nodes (19): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+11 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.17
Nodes (18): net/http.Request, net/http.ResponseWriter, checkResponse, incidentResponse, monitorResponse, Server, toCheckResponse(), Server (+10 more)

### Community 7 - "middleware.go"
Cohesion: 0.07
Nodes (26): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+18 more)

### Community 8 - ".handleFailure"
Cohesion: 0.26
Nodes (6): time.Time, CheckResult, Notifier, NewStateMachine(), newUUIDv7(), StateMachine

### Community 9 - "Incident"
Cohesion: 0.19
Nodes (6): Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo

### Community 10 - "context.Context"
Cohesion: 0.17
Nodes (6): context.Context, CheckRepository, MonitorRepository, NewRetentionWorker(), RetentionWorker, incidentRepoWithTx

### Community 11 - "UnitOfWork"
Cohesion: 0.13
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 12 - "Check"
Cohesion: 0.18
Nodes (4): Check, CheckRepo, CheckRepositoryTx, checkRepoWithTx

### Community 13 - "Monitor"
Cohesion: 0.20
Nodes (7): Monitor, MonitorState, MonitorFilter, scanMonitor(), MonitorRepo, scanner, monitorRepoWithTx

### Community 15 - "boolToInt"
Cohesion: 0.36
Nodes (3): database/sql.Tx, boolToInt(), nullableTime()

### Community 16 - "WebhookNotifier"
Cohesion: 0.28
Nodes (5): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier

### Community 34 - "validate.go"
Cohesion: 0.39
Nodes (11): go_pkg_regexp, validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold(), validateTimeout(), validateUpdateInput() (+3 more)

## Knowledge Gaps
- **69 isolated node(s):** `name`, `private`, `version`, `type`, `dev` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 109 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `database/sql.DB`, `net/http.ResponseWriter`, `.handleFailure`, `context.Context`, `UnitOfWork`, `MonitorService`, `boolToInt`, `WebhookNotifier`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Why does `runServe()` connect `runServe` to `serve.go`, `database/sql.DB`, `middleware.go`, `Incident`, `context.Context`, `WebhookNotifier`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `runServe`, `database/sql.DB`, `Incident`, `MonitorService`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **What connects `name`, `private`, `version` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08306010928961749 - nodes in this community are weakly interconnected._
- **Should `runServe` be split into smaller, more focused modules?**
  _Cohesion score 0.055288461538461536 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.0507399577167019 - nodes in this community are weakly interconnected._