# Graph Report - pingou-health-checker  (2026-11-06)

## Corpus Check
- 236 files · ~253,698 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2124 nodes · 7014 edges · 124 communities (104 shown, 20 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 345 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `09db253f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- theme.js
- MonitorForm.jsx
- package.json
- UnitOfWork
- time.Time
- net/http.ResponseWriter
- Testing Strategy & Framework
- 16. Nomes claros nos eventos do webhook
- Scheduler
- IDEAS.md
- newTestServer
- What You Must Do When Invoked
- OutboxEvent
- Pingou Health Checker - Project Vision & Goals
- newMaintenanceEnv
- MonitorService
- Plano de Desenvolvimento: {{Titulo}}
- Vite Logo
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Code Conventions & Standards
- Architecture Overview
- Technical Concerns & Risk Assessment
- Integrations & External Dependencies
- Technology Stack Analysis
- Project Structure Analysis
- graphify reference: extra exports and benchmark
- newMaintenanceServiceDB
- MaintenanceConflictError
- Hero Image
- React Logo
- 14. Manutenção: pausa com motivo e janela agendada
- timezone.js
- 3. Página de Configurações
- MonitorRepo
- Monitor
- testing.T
- Incident
- 7. Largura no desktop + busca parcial
- format.js
- graphify reference: query, path, explain
- context.Context
- Pingou Health Checker - Project Roadmap
- 17. Retry de notificações (outbox)
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- NewIncidentRepo
- newCertEnv
- newCertServer
- DESIGN.md — Design system do Pingou
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- owEnqueue
- newFakeSMTP
- Status: **In Progress**
- Dashboard.jsx
- Key Decisions Made
- Lessons Learned
- Preferences & Guidelines
- Current Blockers
- Session Handoff Information
- Current Technical Debt
- Deferred Ideas & Future Considerations
- Memory Tags
- Next Steps & Immediate Actions
- React + Vite
- graphify reference: GitHub clone and cross-repo merge
- .Update
- TestDBBench
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- MonitorDetail.jsx
- graphify reference: transcribe video and audio
- notifyDown
- memtest.sh
- Open
- Dependencies & Blockers
- Settings.jsx
- .checkHeartbeats
- extraction-spec.md
- newTestRegistry
- Histórico de implementações (planos aplicados)
- migrations_test.go
- validateMonitor
- hotQueries
- incidents_test.go
- Status: **Planned**
- check-tokens.sh
- runServe
- commands/info.go
- 5. Remodelagem mobile
- 8. Notificação via ntfy
- 9. Eixo Y dinâmico + gráfico de pings
- database/sql.DB
- NewOutboxRepo
- certserver/main.go
- 6. E-mail (SMTP) e envio fora do lock
- Status: **Planned**
- devDependencies
- main.jsx
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- lib/monitors.js
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- dependencies
- Success Metrics
- Release Schedule
- Risk Management
- 11. Fuso horário + desligamento ordenado
- 12. Índices e desempenho do banco
- 13. Aviso de expiração do certificado SSL
- dbbench.sh
- dbmem.sh

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 76 edges
2. `Incident` - 36 edges
3. `MaintenanceWindow` - 36 edges
4. `react` - 35 edges
5. `newTestServer()` - 34 edges
6. `runServe()` - 33 edges
7. `do()` - 32 edges
8. `newEnv()` - 29 edges
9. `lucide-react` - 29 edges
10. `Scheduler` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `EmailCard()`  [INFERRED]
  docs/implementacoes.md → web/src/components/settings/EmailCard.jsx

## Import Cycles
- None detected.

## Communities (124 total, 20 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.19
Nodes (48): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_tls, go_pkg_database_sql, go_pkg_encoding_base64, go_pkg_encoding_json, go_pkg_errors (+40 more)

### Community 1 - "theme.js"
Cohesion: 0.15
Nodes (20): 15. Modo escuro (Claro, Escuro, Sistema), Decisões, Descartado, Entregue, Limitações conhecidas, Validação, ref_node_fs, ThemeContext (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.10
Nodes (23): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm() (+15 more)

### Community 3 - "package.json"
Cohesion: 0.12
Nodes (19): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+11 more)

### Community 4 - "UnitOfWork"
Cohesion: 0.10
Nodes (6): OutboxRepositoryTx, UnitOfWork, erroringOutbox, failingCommit, failingEnqueueUoW, MonitorRepositoryTx

### Community 5 - "time.Time"
Cohesion: 0.06
Nodes (27): DownWindow, database/sql.Result, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceWindowView (+19 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (41): net/http.Request, net/http.ResponseWriter, net/url.Values, checkResponse, downWindowResponse, incidentMonitorResponse, incidentResponse, maintenanceResponse (+33 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "16. Nomes claros nos eventos do webhook"
Cohesion: 0.50
Nodes (4): 16. Nomes claros nos eventos do webhook, Entregue, Risco, Validação

### Community 9 - "Scheduler"
Cohesion: 0.08
Nodes (14): context.CancelFunc, database/sql.Tx, sync.WaitGroup, CheckRepoTx, IncidentRepoTx, MonitorRepoTx, OutboxRepoTx, Scheduler (+6 more)

### Community 10 - "IDEAS.md"
Cohesion: 0.04
Nodes (48): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🗑️ Descartadas (+40 more)

### Community 11 - "newTestServer"
Cohesion: 0.11
Nodes (42): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+34 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "OutboxEvent"
Cohesion: 0.11
Nodes (7): OutboxStatus, OutboxEvent, scanOutbox(), slowChannel, Channel, OutboxStore, OutboxWorker

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newMaintenanceEnv"
Cohesion: 0.30
Nodes (14): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+6 more)

### Community 16 - "MonitorService"
Cohesion: 0.08
Nodes (12): Check, MonitorSummary, CheckRepository, CheckRepo, MonitorService, newPingToken(), normalizeReason(), benchEnv (+4 more)

### Community 17 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 22 - "Code Conventions & Standards"
Cohesion: 0.05
Nodes (40): API Conventions, API Documentation, API Integration, Authentication, Branch Strategy, Build and Deployment Conventions, CLI Conventions (Cobra & Viper), Code Comments (+32 more)

### Community 26 - "Integrations & External Dependencies"
Cohesion: 0.05
Nodes (38): Advanced Monitoring, API Integrations, API Security, Authentication & Authorization, Build Automation Integration, Build Tools, CLI Framework Integration, Configuration Management (+30 more)

### Community 27 - "Technology Stack Analysis"
Cohesion: 0.07
Nodes (27): Architecture Patterns, Backend Architecture, Backend Stack, Build & Distribution, Build Tools & Development, CLI Framework, Configuration & Environment, Containerization (+19 more)

### Community 28 - "Project Structure Analysis"
Cohesion: 0.07
Nodes (27): Application Entry Point (`cmd/`), Authentication Layer, Backend Structure (`internal/`), Background Processing, Build and Deployment, Build and Distribution Structure, Configuration Files, Configuration Management (+19 more)

### Community 29 - "graphify reference: extra exports and benchmark"
Cohesion: 0.22
Nodes (8): graphify reference: extra exports and benchmark, Step 6b - Wiki (only if --wiki flag), Step 7 - Neo4j export (only if --neo4j or --neo4j-push flag), Step 7a - FalkorDB export (only if --falkordb or --falkordb-push flag), Step 7b - SVG export (only if --svg flag), Step 7c - GraphML export (only if --graphml flag), Step 7d - MCP server (only if --mcp flag), Step 8 - Token reduction benchmark (only if total_words > 5000)

### Community 30 - "newMaintenanceServiceDB"
Cohesion: 0.29
Nodes (11): NewMaintenanceRepo(), addWindowRaw(), newMaintenanceService(), newMaintenanceServiceDB(), TestMaintenanceBatchConflict(), TestMaintenanceEditedIntervalIsUsed(), TestMaintenanceList(), TestMaintenanceUpdateRaceWithStart() (+3 more)

### Community 34 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.40
Nodes (5): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Limitações conhecidas, Validação

### Community 35 - "timezone.js"
Cohesion: 0.08
Nodes (44): RFC-3339, Entregue, ref_node_assert, ref_node_test, closeLabel, CloseNote(), groupByDay(), ResolvedTable() (+36 more)

### Community 36 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 37 - "MonitorRepo"
Cohesion: 0.14
Nodes (10): MonitorConfig, boolToInt(), encodeConfig(), MonitorRepo, nullableString(), nullableTime(), scanMonitor(), nameScanner (+2 more)

### Community 38 - "Monitor"
Cohesion: 0.07
Nodes (38): Registry, SnapshotCert, SnapshotIncident, SnapshotMonitor, net/http.Response, certNotAfter(), requestErrorMessage(), NewRegistry() (+30 more)

### Community 39 - "testing.T"
Cohesion: 0.13
Nodes (28): sync/atomic.Int64, testing.T, TestNtfyHost(), TestMaintenanceRoutesNeedAPIKey(), env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline() (+20 more)

### Community 40 - "Incident"
Cohesion: 0.17
Nodes (12): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepo, scanIncident(), scanIncidentWithMonitor() (+4 more)

### Community 41 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 42 - "format.js"
Cohesion: 0.13
Nodes (26): ChecksTable(), detailOf(), RecentIncidents(), axisTicks(), gapBuckets(), isDown(), linePath(), ResponseChart() (+18 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "context.Context"
Cohesion: 0.07
Nodes (11): context.Context, IncidentRepository, MonitorRepository, truncateError(), Scheduler, NewMaintenanceService(), OutboxRepo, OutboxPurger (+3 more)

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 46 - "17. Retry de notificações (outbox)"
Cohesion: 0.33
Nodes (6): 17. Retry de notificações (outbox), Decisões, Descartado, Entregue, Limitações conhecidas, Validação

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 49 - "NewIncidentRepo"
Cohesion: 0.35
Nodes (11): runAdd(), runMigrateUp(), runRM(), github.com/spf13/cobra.Command, RunMigrations(), NewCheckRepo(), openBench(), TestResultsSameWithOldIndex() (+3 more)

### Community 50 - "newCertEnv"
Cohesion: 0.19
Nodes (19): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+11 more)

### Community 51 - "newCertServer"
Cohesion: 0.20
Nodes (11): HTTPChecker, net/http.Client, net/http/httptest.Server, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage() (+3 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.04
Nodes (42): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.1), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+34 more)

### Community 53 - "graphify reference: commit hook and native CLAUDE.md integration"
Cohesion: 0.50
Nodes (3): For git commit hook, For native CLAUDE.md integration, graphify reference: commit hook and native CLAUDE.md integration

### Community 54 - "graphify reference: incremental update and cluster-only"
Cohesion: 0.50
Nodes (3): For --cluster-only, For --update (incremental re-extraction), graphify reference: incremental update and cluster-only

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 56 - "10. Pausar monitor fecha o incidente"
Cohesion: 0.40
Nodes (5): 10. Pausar monitor fecha o incidente, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 57 - "owEnqueue"
Cohesion: 0.30
Nodes (17): NewOutboxRepoTx(), NewOutboxWorker(), owAlways(), owDB(), owEnqueue(), owRows(), owWorker(), TestWorkerBackoffAndExpiryForDown() (+9 more)

### Community 58 - "newFakeSMTP"
Cohesion: 0.18
Nodes (16): net.Conn, net.Listener, TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage(), newFakeSMTP() (+8 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "Dashboard.jsx"
Cohesion: 0.10
Nodes (35): Desvios do plano, lucide-react, react, react-toastify, maintenanceApi, monitorsApi, web_src_assets_brand_pingou_simbolo, web_src_assets_brand_pingou_simbolo_invertido (+27 more)

### Community 61 - "Key Decisions Made"
Cohesion: 0.50
Nodes (4): Architecture Decisions, Development Strategy Decisions, Key Decisions Made, Scope Decisions

### Community 62 - "Lessons Learned"
Cohesion: 0.50
Nodes (4): Architecture Lessons, Lessons Learned, Process Lessons, Technical Lessons

### Community 63 - "Preferences & Guidelines"
Cohesion: 0.50
Nodes (4): Code Quality Preferences, Communication Preferences, Development Preferences, Preferences & Guidelines

### Community 64 - "Current Blockers"
Cohesion: 0.50
Nodes (4): Current Blockers, Decision Blockers, Resource Blockers, Technical Blockers

### Community 65 - "Session Handoff Information"
Cohesion: 0.50
Nodes (4): Current Session Context, Key Files to Reference, Resume Instructions, Session Handoff Information

### Community 66 - "Current Technical Debt"
Cohesion: 0.50
Nodes (4): Current Technical Debt, High Priority, Low Priority, Medium Priority

### Community 67 - "Deferred Ideas & Future Considerations"
Cohesion: 0.50
Nodes (4): Deferred Ideas & Future Considerations, Feature Deferrals, Integration Deferrals, Technical Deferrals

### Community 68 - "Memory Tags"
Cohesion: 0.50
Nodes (4): Development Planning, Memory Tags, Project Management, Technical Analysis

### Community 69 - "Next Steps & Immediate Actions"
Cohesion: 0.50
Nodes (4): Immediate (Next 1-2 weeks), Medium-term (Next 3-6 months), Next Steps & Immediate Actions, Short-term (Next 1-2 months)

### Community 70 - "React + Vite"
Cohesion: 0.50
Nodes (3): Expanding the ESLint configuration, React Compiler, React + Vite

### Community 72 - ".Update"
Cohesion: 0.22
Nodes (8): TestValidNtfyTopic(), validNtfyTopic(), normalizeTimezone(), TestNormalizeTimezone(), parseRecipients(), TestParseRecipients(), NtfyTestResult, Settings

### Community 73 - "TestDBBench"
Cohesion: 0.13
Nodes (16): sync.Mutex, time.Duration, rateLimiter, TestRateLimiterWindowResets(), newRateLimiter(), NewCheckRepoTx(), envInt(), measure() (+8 more)

### Community 74 - "bodyFor"
Cohesion: 0.07
Nodes (36): time.Location, TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), formatDuration(), joinUnits() (+28 more)

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

### Community 80 - "MonitorDetail.jsx"
Cohesion: 0.13
Nodes (20): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, ConfigCard(), DetailTopBar(), DownAlert(), RangeTabs() (+12 more)

### Community 82 - "notifyDown"
Cohesion: 0.14
Nodes (25): net/http.Header, TestNtfyNotifierCertExpiring(), TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookDownAndRecoveryEventNames(), TestWebhookPayloadByType() (+17 more)

### Community 84 - "Open"
Cohesion: 0.18
Nodes (30): Open(), NewMonitorRepo(), TestPauseClosesOpenIncidentSilently(), TestPauseIsSilentAndResumeReopensIncident(), TestPauseReasonSavedAndClearedOnResume(), TestPauseReasonTooLong(), TestPauseWithoutIncidentAndResumeKeepState(), TestResumeIgnoresFailuresBeforePause() (+22 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "Settings.jsx"
Cohesion: 0.13
Nodes (24): settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText(), OPTIONS (+16 more)

### Community 90 - "newTestRegistry"
Cohesion: 0.24
Nodes (11): TCPChecker, net.Dialer, newTestRegistry(), tcpMonitor(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused() (+3 more)

### Community 92 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.33
Nodes (6): 2. Página de Incidentes, Decisões, Descartado, Entregue, Histórico de implementações (planos aplicados), Pendências abertas

### Community 94 - "migrations_test.go"
Cohesion: 0.31
Nodes (14): go_pkg_github_com_pressly_goose_v3, OpenDB(), sqliteDSN(), assertFKOn(), assertNoFKViolations(), countRows(), mustExec(), setupGoose() (+6 more)

### Community 95 - "validateMonitor"
Cohesion: 0.13
Nodes (15): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+7 more)

### Community 96 - "hotQueries"
Cohesion: 0.18
Nodes (15): testing.TB, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+7 more)

### Community 97 - "incidents_test.go"
Cohesion: 0.38
Nodes (10): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+2 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 103 - "runServe"
Cohesion: 0.07
Nodes (35): runMigrateDown(), runMigrateStatus(), forceExitOnSecondSignal(), releaseLock(), runServe(), go_pkg_embed, go_pkg_io_fs, log/slog.Level (+27 more)

### Community 104 - "commands/info.go"
Cohesion: 0.07
Nodes (26): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+18 more)

### Community 105 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 106 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 107 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 110 - "database/sql.DB"
Cohesion: 0.26
Nodes (16): database/sql.DB, ExportHandler, Status(), dumpRows(), smallBenchDB(), TestDeleteOlderThanBatches(), TestSummariesHistoryMatchesLegacy(), lastStep() (+8 more)

### Community 111 - "NewOutboxRepo"
Cohesion: 0.53
Nodes (9): NewOutboxRepo(), enqueue(), newOutboxDB(), outboxEvent(), TestOutboxDueExpiresUnreadableRow(), TestOutboxEnqueueAndDue(), TestOutboxMarks(), TestOutboxPurgeDoneBefore() (+1 more)

### Community 114 - "certserver/main.go"
Cohesion: 0.09
Nodes (24): go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_crypto_x509, go_pkg_crypto_x509_pkix, go_pkg_encoding_hex, go_pkg_encoding_pem (+16 more)

### Community 115 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 241 - "main.jsx"
Cohesion: 0.09
Nodes (21): ⚠️ Notes, axios, react-router-dom, version, client, incidentsApi, web_src_assets_brand_pingou_marca_horizontal, web_src_assets_brand_pingou_marca_horizontal_negativa (+13 more)

### Community 242 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 4.1: Multi-tenancy & Access Control, Milestone 4.2: Advanced Authentication, Milestone 4.3: Scalability Improvements, Phase 4: Enterprise Features (Q4 2024), Status: **Planned**

### Community 243 - "Status: **Future**"
Cohesion: 0.40
Nodes (5): Milestone 5.1: Advanced Monitoring, Milestone 5.2: Integration Ecosystem, Milestone 5.3: Public Features, Phase 5: Advanced Features & Integrations (Q1-Q2 2025), Status: **Future**

### Community 244 - "Status: **Continuous**"
Cohesion: 0.40
Nodes (5): Milestone 6.1: Community Building, Milestone 6.2: Developer Experience, Milestone 6.3: Ecosystem Integration, Phase 6: Ecosystem & Community (Ongoing), Status: **Continuous**

### Community 245 - "lib/monitors.js"
Cohesion: 0.12
Nodes (32): HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu(), Badge() (+24 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.22
Nodes (9): scripts, build, check:tokens, dev, format, format:check, lint, preview (+1 more)

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 253 - "Success Metrics"
Cohesion: 0.50
Nodes (4): Community Metrics, Success Metrics, Technical Metrics, User Metrics

### Community 254 - "Release Schedule"
Cohesion: 0.50
Nodes (4): Current Version Tracking, Release Frequency, Release Schedule, Versioning Strategy

### Community 257 - "Risk Management"
Cohesion: 0.67
Nodes (3): Project Risks, Risk Management, Technical Risks

### Community 261 - "11. Fuso horário + desligamento ordenado"
Cohesion: 0.33
Nodes (6): 11. Fuso horário + desligamento ordenado, Decisões, Descartado, Entregue, Limitações conhecidas, Validação

### Community 262 - "12. Índices e desempenho do banco"
Cohesion: 0.33
Nodes (6): 12. Índices e desempenho do banco, Decisões, Descartado, Entregue, Limitações conhecidas, Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

### Community 265 - "13. Aviso de expiração do certificado SSL"
Cohesion: 0.29
Nodes (7): 13. Aviso de expiração do certificado SSL, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação

## Knowledge Gaps
- **525 isolated node(s):** `Overview`, `Project Type`, `Definition of Done (DoD)`, `Success Criteria`, `Tech Stack` (+520 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 639 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `theme.js`, `14. Manutenção: pausa com motivo e janela agendada`, `3. Página de Configurações`, `11. Fuso horário + desligamento ordenado`, `12. Índices e desempenho do banco`, `16. Nomes claros nos eventos do webhook`, `13. Aviso de expiração do certificado SSL`, `5. Remodelagem mobile`, `7. Largura no desktop + busca parcial`, `8. Notificação via ntfy`, `9. Eixo Y dinâmico + gráfico de pings`, `17. Retry de notificações (outbox)`, `MonitorDetail.jsx`, `6. E-mail (SMTP) e envio fora do lock`, `DESIGN.md — Design system do Pingou`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `10. Pausar monitor fecha o incidente`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `ThemedImage()` connect `main.jsx` to `theme.js`, `Dashboard.jsx`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `serve.go`, `UnitOfWork`, `time.Time`, `net/http.ResponseWriter`, `MonitorRepo`, `testing.T`, `Scheduler`, `bodyFor`, `context.Context`, `newMaintenanceEnv`, `MonitorService`, `notifyDown`, `newCertEnv`, `Open`, `.checkHeartbeats`, `newTestRegistry`, `validateMonitor`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `Overview`, `Project Type`, `Definition of Done (DoD)` to the rest of the system?**
  _525 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0967741935483871 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.11688311688311688 - nodes in this community are weakly interconnected._
- **Should `UnitOfWork` be split into smaller, more focused modules?**
  _Cohesion score 0.09523809523809523 - nodes in this community are weakly interconnected._