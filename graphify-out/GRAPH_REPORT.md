# Graph Report - pingou-health-checker  (2026-10-05)

## Corpus Check
- 110 files · ~90,372 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 9, .woff2 4, .example 1)

## Summary
- 495 nodes · 1386 edges · 28 communities (16 shown, 12 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2a5299d6`
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
- Monitor
- context.Context
- SettingsRepo
- Check
- boolToInt
- WebhookNotifier
- MonitorRepo
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
7. `writeError()` - 14 edges
8. `Check` - 14 edges
9. `DESIGN.md — Design system do Pingou` - 13 edges
10. `writeJSON()` - 13 edges

## Surprising Connections (you probably didn't know these)
- `runServe()` --calls--> `NewSettingsRepo()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/repository/settings_repo.go
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go
- `runServe()` --calls--> `NewSettingsService()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/service/settings_service.go
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewHTTPChecker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/checker/http_checker.go

## Import Cycles
- None detected.

## Communities (28 total, 12 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.20
Nodes (20): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_google_uuid (+12 more)

### Community 1 - "main.jsx"
Cohesion: 0.07
Nodes (36): lucide-react, react, react-router-dom, base, client, incidentsApi, monitorsApi, settingsApi (+28 more)

### Community 2 - "runServe"
Cohesion: 0.06
Nodes (50): runAdd(), copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion() (+42 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (42): autoprefixer, axios, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss (+34 more)

### Community 4 - "database/sql.DB"
Cohesion: 0.15
Nodes (15): context.CancelFunc, database/sql.DB, sync.Mutex, Checker, CheckRepoTx, NewCheckRepoTx(), IncidentRepoTx, NewIncidentRepoTx() (+7 more)

### Community 5 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.10
Nodes (19): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+11 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.19
Nodes (16): net/http.Request, net/http.ResponseWriter, incidentResponse, monitorResponse, Server, Server, Server, Server (+8 more)

### Community 7 - "middleware.go"
Cohesion: 0.10
Nodes (18): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+10 more)

### Community 8 - "Monitor"
Cohesion: 0.16
Nodes (9): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, MonitorRepository, Notifier, NewStateMachine() (+1 more)

### Community 9 - "context.Context"
Cohesion: 0.18
Nodes (8): context.Context, Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo, IncidentRepositoryTx

### Community 10 - "SettingsRepo"
Cohesion: 0.13
Nodes (11): updateSettingsRequest, CheckRepository, toUpdateSettingsInput(), SettingsRepo, NewSettingsRepo(), NewRetentionWorker(), SettingsService, UpdateSettingsInput (+3 more)

### Community 11 - "Check"
Cohesion: 0.12
Nodes (7): checkResponse, Check, toCheckResponse(), CheckRepo, CheckRepositoryTx, checkRepoWithTx, incidentRepoWithTx

### Community 12 - "boolToInt"
Cohesion: 0.24
Nodes (4): database/sql.Tx, boolToInt(), nullableTime(), monitorRepoWithTx

### Community 13 - "WebhookNotifier"
Cohesion: 0.24
Nodes (7): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 14 - "MonitorRepo"
Cohesion: 0.29
Nodes (3): scanMonitor(), MonitorRepo, scanner

### Community 34 - "validate.go"
Cohesion: 0.19
Nodes (14): go_pkg_regexp, MonitorService, newUUIDv7(), validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold() (+6 more)

## Knowledge Gaps
- **69 isolated node(s):** `nav`, `dotColor`, `defaults`, `stateStyles`, `MonitorEnabled` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 110 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `validate.go`, `database/sql.DB`, `net/http.ResponseWriter`, `boolToInt`, `WebhookNotifier`, `MonitorRepo`, `UnitOfWork`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `runServe()` connect `runServe` to `serve.go`, `database/sql.DB`, `middleware.go`, `context.Context`, `SettingsRepo`, `WebhookNotifier`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `runServe`, `validate.go`, `database/sql.DB`, `context.Context`, `SettingsRepo`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `nav`, `dotColor`, `defaults` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0715846994535519 - nodes in this community are weakly interconnected._
- **Should `runServe` be split into smaller, more focused modules?**
  _Cohesion score 0.055288461538461536 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.049494949494949494 - nodes in this community are weakly interconnected._