# Graph Report - pingou-health-checker  (2026-09-24)

## Corpus Check
- Corpus is ~47,717 words - fits in a single context window. You may not need a graph.

## Summary
- 479 nodes · 1379 edges · 34 communities (17 shown, 17 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 85,308 input · 2,643 output

## Community Hubs (Navigation)
- CLI Lock & Init
- Frontend API Client
- DB Export/Import CLI
- Frontend Build Tooling
- CLI Commands (add/rm/serve)
- HTTP Server & Settings
- HTTP Handlers (REST API)
- Server Middleware & Embed
- Monitor Domain Model
- Monitor Repository Ops
- Incident Repository
- Transactional Repositories
- Check Repository
- Monitor Repo Core
- Check Retention Worker
- DB Tx Helpers
- HTTP Checker & Notifier
- PRD & MVP Docs
- Favicon & Vite Logo
- Docker Compose Doc
- GitHub Repo Link
- Pingou README
- Architecture Overview Doc
- Risk Assessment Doc
- Bluesky Icon
- Discord Icon
- Documentation Icon
- GitHub Icon
- Social Icon
- X (Twitter) Icon
- Hero Image
- React Logo

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 38 edges
2. `Incident` - 25 edges
3. `runServe()` - 21 edges
4. `Scheduler` - 18 edges
5. `MonitorService` - 16 edges
6. `Check` - 14 edges
7. `writeError()` - 14 edges
8. `tokens` - 14 edges
9. `writeJSON()` - 13 edges
10. `Server` - 12 edges

## Surprising Connections (you probably didn't know these)
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewHTTPChecker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/checker/http_checker.go
- `runServe()` --calls--> `Load()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/config/config.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Brand Assets** — web_public_favicon, web_src_assets_vite, web_src_assets_hero [INFERRED 0.80]
- **Social Media Icons** — web_public_icons_bluesky_icon, web_public_icons_discord_icon, web_public_icons_github_icon, web_public_icons_x_icon [EXTRACTED 1.00]

## Communities (34 total, 17 thin omitted)

### Community 0 - "CLI Lock & Init"
Cohesion: 0.18
Nodes (25): acquireLock(), checkExistingLock(), MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors (+17 more)

### Community 1 - "Frontend API Client"
Cohesion: 0.09
Nodes (33): axios, react, base, client, incidentsApi, monitorsApi, settingsApi, nav (+25 more)

### Community 2 - "DB Export/Import CLI"
Cohesion: 0.06
Nodes (38): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+30 more)

### Community 3 - "Frontend Build Tooling"
Cohesion: 0.05
Nodes (41): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, react-dom (+33 more)

### Community 4 - "CLI Commands (add/rm/serve)"
Cohesion: 0.10
Nodes (29): runAdd(), runRM(), releaseLock(), runServe(), database/sql.DB, github.com/spf13/cobra.Command, sync.Mutex, Open() (+21 more)

### Community 5 - "HTTP Server & Settings"
Cohesion: 0.10
Nodes (22): net/http.ServeMux, net/http.Server, updateSettingsRequest, Server, NewServer(), toUpdateSettingsInput(), MonitorService, SettingsService (+14 more)

### Community 6 - "HTTP Handlers (REST API)"
Cohesion: 0.17
Nodes (18): net/http.Request, net/http.ResponseWriter, checkResponse, incidentResponse, monitorResponse, Server, toCheckResponse(), Server (+10 more)

### Community 7 - "Server Middleware & Embed"
Cohesion: 0.14
Nodes (14): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, responseWriter (+6 more)

### Community 8 - "Monitor Domain Model"
Cohesion: 0.21
Nodes (9): context.CancelFunc, time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, job, monitorRepoWithTx (+1 more)

### Community 9 - "Monitor Repository Ops"
Cohesion: 0.19
Nodes (5): context.Context, IncidentRepository, MonitorRepository, IncidentService, NewIncidentService()

### Community 10 - "Incident Repository"
Cohesion: 0.23
Nodes (4): Incident, scanIncident(), IncidentRepo, incidentRepoWithTx

### Community 11 - "Transactional Repositories"
Cohesion: 0.13
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 12 - "Check Repository"
Cohesion: 0.18
Nodes (4): Check, CheckRepo, CheckRepositoryTx, checkRepoWithTx

### Community 13 - "Monitor Repo Core"
Cohesion: 0.29
Nodes (3): scanMonitor(), MonitorRepo, scanner

### Community 14 - "Check Retention Worker"
Cohesion: 0.31
Nodes (3): CheckRepository, NewRetentionWorker(), RetentionWorker

### Community 15 - "DB Tx Helpers"
Cohesion: 0.36
Nodes (3): database/sql.Tx, boolToInt(), nullableTime()

### Community 16 - "HTTP Checker & Notifier"
Cohesion: 0.38
Nodes (4): HTTPChecker, net/http.Client, NewHTTPChecker(), WebhookNotifier

## Knowledge Gaps
- **59 isolated node(s):** `github.com/Bradrickcruz/pingou`, `MonitorEnabled`, `Server`, `Server`, `Server` (+54 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 99 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor Domain Model` to `CLI Lock & Init`, `CLI Commands (add/rm/serve)`, `HTTP Server & Settings`, `HTTP Handlers (REST API)`, `Monitor Repository Ops`, `Incident Repository`, `Transactional Repositories`, `Monitor Repo Core`, `DB Tx Helpers`, `HTTP Checker & Notifier`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Why does `runServe()` connect `CLI Commands (add/rm/serve)` to `CLI Lock & Init`, `DB Export/Import CLI`, `HTTP Server & Settings`, `Monitor Repository Ops`, `Check Retention Worker`, `HTTP Checker & Notifier`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Why does `Server` connect `HTTP Server & Settings` to `CLI Lock & Init`, `DB Export/Import CLI`, `CLI Commands (add/rm/serve)`, `Server Middleware & Embed`, `Monitor Repository Ops`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `github.com/Bradrickcruz/pingou`, `MonitorEnabled`, `Server` to the rest of the system?**
  _59 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Frontend API Client` be split into smaller, more focused modules?**
  _Cohesion score 0.08590441621294616 - nodes in this community are weakly interconnected._
- **Should `DB Export/Import CLI` be split into smaller, more focused modules?**
  _Cohesion score 0.0602322206095791 - nodes in this community are weakly interconnected._
- **Should `Frontend Build Tooling` be split into smaller, more focused modules?**
  _Cohesion score 0.0507399577167019 - nodes in this community are weakly interconnected._