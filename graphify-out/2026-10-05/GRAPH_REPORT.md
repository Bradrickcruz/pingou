# Graph Report - pingou-health-checker  (2026-10-05)

## Corpus Check
- 111 files · ~91,996 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 9, .woff2 4, .example 1)

## Summary
- 496 nodes · 1387 edges · 28 communities (15 shown, 13 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c8b5863a`
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
- `runServe()` --calls--> `NewWebhookNotifier()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/service/notifier.go
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go

## Import Cycles
- None detected.

## Communities (28 total, 13 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (22): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_google_uuid (+14 more)

### Community 1 - "main.jsx"
Cohesion: 0.07
Nodes (38): axios, lucide-react, react, react-dom, react-router-dom, base, client, incidentsApi (+30 more)

### Community 2 - "database/sql.DB"
Cohesion: 0.06
Nodes (52): runAdd(), copyFile(), runExportDB(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), runRM(), acquireLock() (+44 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (40): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, tailwindcss (+32 more)

### Community 4 - "Scheduler"
Cohesion: 0.20
Nodes (6): context.CancelFunc, sync.Mutex, Notifier, NewStateMachine(), job, Scheduler

### Community 5 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.10
Nodes (19): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+11 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.19
Nodes (16): net/http.Request, net/http.ResponseWriter, incidentResponse, monitorResponse, Server, Server, Server, Server (+8 more)

### Community 7 - "middleware.go"
Cohesion: 0.08
Nodes (23): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+15 more)

### Community 8 - "Monitor"
Cohesion: 0.17
Nodes (10): time.Time, Checker, CheckResult, Monitor, MonitorState, MonitorFilter, scanMonitor(), MonitorRepo (+2 more)

### Community 9 - "context.Context"
Cohesion: 0.18
Nodes (8): context.Context, Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo, IncidentRepositoryTx

### Community 10 - "info.go"
Cohesion: 0.16
Nodes (11): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim() (+3 more)

### Community 11 - "Check"
Cohesion: 0.12
Nodes (8): checkResponse, Check, CheckRepository, toCheckResponse(), NewRetentionWorker(), CheckRepo, RetentionWorker, CheckRepositoryTx

### Community 12 - "boolToInt"
Cohesion: 0.15
Nodes (6): database/sql.Tx, boolToInt(), nullableTime(), checkRepoWithTx, incidentRepoWithTx, monitorRepoWithTx

### Community 13 - "WebhookNotifier"
Cohesion: 0.28
Nodes (5): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier

### Community 34 - "validate.go"
Cohesion: 0.19
Nodes (14): go_pkg_regexp, MonitorService, newUUIDv7(), validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold() (+6 more)

## Knowledge Gaps
- **69 isolated node(s):** `0. Como usar este arquivo`, `1. Conceito`, `Arquivos oficiais`, `3.1 Tokens`, `3.2 Tailwind` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 111 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `Scheduler`, `net/http.ResponseWriter`, `boolToInt`, `WebhookNotifier`, `MonitorRepository`, `UnitOfWork`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `runServe()` connect `database/sql.DB` to `serve.go`, `middleware.go`, `context.Context`, `Check`, `WebhookNotifier`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `context.Context`, `database/sql.DB`, `validate.go`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `0. Como usar este arquivo`, `1. Conceito`, `Arquivos oficiais` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06810035842293907 - nodes in this community are weakly interconnected._
- **Should `database/sql.DB` be split into smaller, more focused modules?**
  _Cohesion score 0.06060606060606061 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.05204872646733112 - nodes in this community are weakly interconnected._