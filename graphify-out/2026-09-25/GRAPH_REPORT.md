# Graph Report - pingou-health-checker  (2026-09-25)

## Corpus Check
- 108 files · ~80,566 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 9, .css 2, .example 1)

## Summary
- 478 nodes · 1378 edges · 35 communities (18 shown, 17 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `759c9e61`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- main.jsx
- database.go
- package.json
- database/sql.DB
- MonitorService
- net/http.ResponseWriter
- middleware.go
- Monitor
- Incident
- context.Context
- UnitOfWork
- Check
- MonitorRepo
- CheckRepository
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
9. `writeJSON()` - 13 edges
10. `sqliteUnitOfWork` - 12 edges

## Surprising Connections (you probably didn't know these)
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go
- `runServe()` --calls--> `NewHTTPChecker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/checker/http_checker.go
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `Load()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/config/config.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Social Media Icons** — web_public_icons_bluesky_icon, web_public_icons_discord_icon, web_public_icons_github_icon, web_public_icons_x_icon [EXTRACTED 1.00]
- **Brand Assets** — web_public_favicon, web_src_assets_vite, web_src_assets_hero [INFERRED 0.80]

## Communities (35 total, 17 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (22): acquireLock(), checkExistingLock(), MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors (+14 more)

### Community 1 - "main.jsx"
Cohesion: 0.08
Nodes (34): axios, react, react-router-dom, base, client, incidentsApi, monitorsApi, settingsApi (+26 more)

### Community 2 - "database.go"
Cohesion: 0.06
Nodes (38): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+30 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (40): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, react-dom (+32 more)

### Community 4 - "database/sql.DB"
Cohesion: 0.10
Nodes (28): runAdd(), runRM(), releaseLock(), runServe(), context.CancelFunc, database/sql.DB, github.com/spf13/cobra.Command, sync.Mutex (+20 more)

### Community 5 - "MonitorService"
Cohesion: 0.33
Nodes (3): MonitorService, newUUIDv7(), Reloader

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.17
Nodes (18): net/http.Request, net/http.ResponseWriter, checkResponse, incidentResponse, monitorResponse, Server, toCheckResponse(), Server (+10 more)

### Community 7 - "middleware.go"
Cohesion: 0.08
Nodes (23): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+15 more)

### Community 8 - "Monitor"
Cohesion: 0.22
Nodes (9): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, Notifier, NewStateMachine(), monitorRepoWithTx (+1 more)

### Community 9 - "Incident"
Cohesion: 0.19
Nodes (6): Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo

### Community 10 - "context.Context"
Cohesion: 0.22
Nodes (4): context.Context, MonitorRepository, checkRepoWithTx, incidentRepoWithTx

### Community 11 - "UnitOfWork"
Cohesion: 0.13
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 12 - "Check"
Cohesion: 0.22
Nodes (3): Check, CheckRepo, CheckRepositoryTx

### Community 13 - "MonitorRepo"
Cohesion: 0.29
Nodes (3): scanMonitor(), MonitorRepo, scanner

### Community 14 - "CheckRepository"
Cohesion: 0.31
Nodes (3): CheckRepository, NewRetentionWorker(), RetentionWorker

### Community 15 - "boolToInt"
Cohesion: 0.36
Nodes (3): database/sql.Tx, boolToInt(), nullableTime()

### Community 16 - "WebhookNotifier"
Cohesion: 0.22
Nodes (7): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 34 - "validate.go"
Cohesion: 0.39
Nodes (11): go_pkg_regexp, validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold(), validateTimeout(), validateUpdateInput() (+3 more)

## Knowledge Gaps
- **59 isolated node(s):** `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse`, `updateMonitorRequest`, `base` (+54 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 98 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `database/sql.DB`, `MonitorService`, `net/http.ResponseWriter`, `context.Context`, `UnitOfWork`, `MonitorRepo`, `boolToInt`, `WebhookNotifier`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Why does `runServe()` connect `database/sql.DB` to `serve.go`, `database.go`, `middleware.go`, `Incident`, `CheckRepository`, `WebhookNotifier`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `database.go`, `database/sql.DB`, `MonitorService`, `Incident`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `MonitorEnabled`, `createMonitorRequest`, `paginatedResponse` to the rest of the system?**
  _59 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08416130917592052 - nodes in this community are weakly interconnected._
- **Should `database.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0602322206095791 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.05204872646733112 - nodes in this community are weakly interconnected._