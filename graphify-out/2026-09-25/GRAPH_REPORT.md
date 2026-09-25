# Graph Report - pingou-health-checker  (2026-09-25)

## Corpus Check
- 108 files · ~80,566 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 9, .css 2, .example 1)

## Summary
- 498 nodes · 1397 edges · 34 communities (17 shown, 17 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `68b02b18`
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
- Incident
- context.Context
- UnitOfWork
- info.go
- MonitorRepo
- boolToInt
- WebhookNotifier
- Divergências entre PRD e Implementação
- Favicon
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Pingou README
- Architecture Overview
- Technical Concerns & Risk Assessment
- Bluesky Icon
- Discord Icon
- Documentation Icon
- GitHub Icon
- Social Icon
- X (Twitter) Icon
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

## Hyperedges (group relationships)
- **Social Media Icons** — web_public_icons_bluesky_icon, web_public_icons_discord_icon, web_public_icons_github_icon, web_public_icons_x_icon [EXTRACTED 1.00]
- **Brand Assets** — web_public_favicon, web_src_assets_vite, web_src_assets_hero [INFERRED 0.80]

## Communities (34 total, 17 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.18
Nodes (24): acquireLock(), checkExistingLock(), MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors (+16 more)

### Community 1 - "main.jsx"
Cohesion: 0.09
Nodes (33): axios, react, base, client, incidentsApi, monitorsApi, settingsApi, nav (+25 more)

### Community 2 - "database/sql.DB"
Cohesion: 0.07
Nodes (47): runAdd(), copyFile(), runExportDB(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), runRM(), releaseLock() (+39 more)

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
Cohesion: 0.07
Nodes (26): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+18 more)

### Community 8 - "Monitor"
Cohesion: 0.24
Nodes (8): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, Notifier, NewStateMachine(), StateMachine

### Community 9 - "Incident"
Cohesion: 0.18
Nodes (6): Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo

### Community 10 - "context.Context"
Cohesion: 0.11
Nodes (10): context.Context, Check, CheckRepository, MonitorRepository, NewRetentionWorker(), CheckRepo, RetentionWorker, CheckRepositoryTx (+2 more)

### Community 11 - "UnitOfWork"
Cohesion: 0.12
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 12 - "info.go"
Cohesion: 0.16
Nodes (11): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim() (+3 more)

### Community 13 - "MonitorRepo"
Cohesion: 0.29
Nodes (3): scanMonitor(), MonitorRepo, scanner

### Community 15 - "boolToInt"
Cohesion: 0.27
Nodes (4): database/sql.Tx, boolToInt(), nullableTime(), monitorRepoWithTx

### Community 16 - "WebhookNotifier"
Cohesion: 0.32
Nodes (5): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier

### Community 34 - "validate.go"
Cohesion: 0.19
Nodes (14): go_pkg_regexp, MonitorService, newUUIDv7(), validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold() (+6 more)

## Knowledge Gaps
- **75 isolated node(s):** `0. Como usar este arquivo`, `1. Conceito`, `Arquivos oficiais`, `3.1 Tokens`, `3.2 Tailwind` (+70 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 115 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `Scheduler`, `net/http.ResponseWriter`, `context.Context`, `UnitOfWork`, `MonitorRepo`, `boolToInt`, `WebhookNotifier`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `runServe()` connect `database/sql.DB` to `serve.go`, `middleware.go`, `Incident`, `context.Context`, `WebhookNotifier`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `Incident`, `database/sql.DB`, `validate.go`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `0. Como usar este arquivo`, `1. Conceito`, `Arquivos oficiais` to the rest of the system?**
  _75 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08590441621294616 - nodes in this community are weakly interconnected._
- **Should `database/sql.DB` be split into smaller, more focused modules?**
  _Cohesion score 0.06836158192090395 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.0507399577167019 - nodes in this community are weakly interconnected._