# Graph Report - pingou-health-checker  (2026-10-05)

## Corpus Check
- 110 files · ~90,025 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 9, .woff2 4, .example 1)

## Summary
- 495 nodes · 1384 edges · 22 communities (11 shown, 11 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e7f1377d`
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
- context.Context
- SettingsRepo
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
6. `Check` - 14 edges
7. `writeError()` - 14 edges
8. `DESIGN.md — Design system do Pingou` - 13 edges
9. `writeJSON()` - 13 edges
10. `react` - 12 edges

## Surprising Connections (you probably didn't know these)
- `runServe()` --calls--> `NewRetentionWorker()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/scheduler/retention.go
- `runServe()` --calls--> `NewWebhookNotifier()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/service/notifier.go
- `runMigrateStatus()` --calls--> `ListMigrations()`  [EXTRACTED]
  cmd/pingou/commands/migrate.go → internal/database/database.go
- `runServe()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/handler/server.go
- `runServe()` --calls--> `NewCheckRepoTx()`  [EXTRACTED]
  cmd/pingou/commands/serve.go → internal/repository/check_repo_tx.go

## Import Cycles
- None detected.

## Communities (22 total, 11 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (22): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_database_sql, go_pkg_encoding_json, go_pkg_errors, go_pkg_fmt, go_pkg_github_com_google_uuid (+14 more)

### Community 1 - "main.jsx"
Cohesion: 0.07
Nodes (36): lucide-react, react, react-router-dom, base, client, incidentsApi, monitorsApi, settingsApi (+28 more)

### Community 2 - "runServe"
Cohesion: 0.05
Nodes (53): HTTPChecker, runAdd(), copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig() (+45 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (42): autoprefixer, axios, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss (+34 more)

### Community 4 - "database/sql.DB"
Cohesion: 0.10
Nodes (19): context.CancelFunc, database/sql.DB, sync.Mutex, Checker, CheckRepoTx, NewCheckRepoTx(), IncidentRepoTx, NewIncidentRepoTx() (+11 more)

### Community 5 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.10
Nodes (19): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+11 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.17
Nodes (18): net/http.Request, net/http.ResponseWriter, checkResponse, incidentResponse, monitorResponse, Server, toCheckResponse(), Server (+10 more)

### Community 7 - "middleware.go"
Cohesion: 0.10
Nodes (18): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_embed, go_pkg_encoding_hex, go_pkg_io_fs, go_pkg_runtime_debug, net/http.Handler, net/http.ServeMux (+10 more)

### Community 9 - "context.Context"
Cohesion: 0.05
Nodes (29): context.Context, database/sql.Tx, time.Time, Check, CheckResult, Incident, Monitor, MonitorState (+21 more)

### Community 10 - "SettingsRepo"
Cohesion: 0.13
Nodes (11): updateSettingsRequest, CheckRepository, toUpdateSettingsInput(), SettingsRepo, NewSettingsRepo(), NewRetentionWorker(), SettingsService, UpdateSettingsInput (+3 more)

### Community 34 - "validate.go"
Cohesion: 0.19
Nodes (14): go_pkg_regexp, MonitorService, newUUIDv7(), validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold() (+6 more)

## Knowledge Gaps
- **69 isolated node(s):** `nav`, `dotColor`, `defaults`, `stateStyles`, `variants` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 110 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `context.Context` to `serve.go`, `validate.go`, `database/sql.DB`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `runServe()` connect `runServe` to `serve.go`, `database/sql.DB`, `middleware.go`, `context.Context`, `SettingsRepo`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `Server` connect `middleware.go` to `serve.go`, `runServe`, `validate.go`, `database/sql.DB`, `context.Context`, `SettingsRepo`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `nav`, `dotColor`, `defaults` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07049180327868852 - nodes in this community are weakly interconnected._
- **Should `runServe` be split into smaller, more focused modules?**
  _Cohesion score 0.05179982440737489 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.049494949494949494 - nodes in this community are weakly interconnected._