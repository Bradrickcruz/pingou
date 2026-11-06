# Graph Report - pingou-health-checker  (2026-11-06)

## Corpus Check
- 237 files · ~253,874 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2134 nodes · 7024 edges · 131 communities (112 shown, 19 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 346 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d346cee3`
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
- CheckRepo
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
- Plano de Desenvolvimento: {{Titulo}}
- MaintenanceConflictError
- Hero Image
- React Logo
- MaintenanceModal
- timezone.js
- notifyDown
- MonitorRepo
- Monitor
- newEnv
- Incident
- database/sql.DB
- ResponseChart.jsx
- graphify reference: query, path, explain
- context.Context
- Pingou Health Checker - Project Roadmap
- 17. Retry de notificações (outbox)
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- runServe
- newCertEnv
- SendTestNtfy
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
- SendTestWebhook
- time.Duration
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- MonitorDetail.jsx
- graphify reference: transcribe video and audio
- NewNtfyNotifier
- memtest.sh
- Open
- Dependencies & Blockers
- Settings.jsx
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- extraction-spec.md
- 3.6 Tema escuro (dark mode)
- testing.T
- database/sql.Tx
- Histórico de implementações (planos aplicados)
- memtest: prova da promessa de memória
- migrations_test.go
- validateMonitor
- hotQueries
- incidents_test.go
- Status: **Planned**
- RetentionWorker
- check-tokens.sh
- WebhookNotifier
- dbbench: desempenho do banco com volume real
- newTestServerEnv
- commands/info.go
- 5. Remodelagem mobile
- 8. Notificação via ntfy
- 9. Eixo Y dinâmico + gráfico de pings
- formatDuration
- seed/main.go
- NewOutboxRepo
- siteCert
- 6. E-mail (SMTP) e envio fora do lock
- Status: **Planned**
- devDependencies
- format.js
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
8. `lucide-react` - 29 edges
9. `newEnv()` - 29 edges
10. `Scheduler` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `EmailCard()`  [INFERRED]
  docs/implementacoes.md → web/src/components/settings/EmailCard.jsx

## Import Cycles
- None detected.

## Communities (131 total, 19 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.17
Nodes (56): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_tls, go_pkg_crypto_x509 (+48 more)

### Community 1 - "theme.js"
Cohesion: 0.15
Nodes (20): 15. Modo escuro (Claro, Escuro, Sistema), Decisões, Descartado, Entregue, Limitações conhecidas, Validação, ref_node_fs, ThemeContext (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.15
Nodes (13): blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm(), PING_EVERY, TYPES (+5 more)

### Community 3 - "package.json"
Cohesion: 0.12
Nodes (19): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+11 more)

### Community 4 - "UnitOfWork"
Cohesion: 0.13
Nodes (6): OutboxRepositoryTx, UnitOfWork, CheckRepositoryTx, erroringOutbox, failingCommit, failingEnqueueUoW

### Community 5 - "time.Time"
Cohesion: 0.06
Nodes (30): DownWindow, database/sql.NullString, database/sql.Result, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MaintenanceWindow (+22 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (48): go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, net/http.Request, net/http.ResponseWriter, net/url.Values, checkResponse, downWindowResponse (+40 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "16. Nomes claros nos eventos do webhook"
Cohesion: 0.50
Nodes (4): 16. Nomes claros nos eventos do webhook, Entregue, Risco, Validação

### Community 9 - "Scheduler"
Cohesion: 0.20
Nodes (4): context.CancelFunc, sync.WaitGroup, Scheduler, job

### Community 10 - "IDEAS.md"
Cohesion: 0.04
Nodes (48): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🗑️ Descartadas (+40 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (41): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+33 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "OutboxEvent"
Cohesion: 0.09
Nodes (9): OutboxStatus, sync.Mutex, OutboxEvent, scanOutbox(), slowChannel, countingChannel, OutboxStore, OutboxWorker (+1 more)

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newMaintenanceEnv"
Cohesion: 0.30
Nodes (14): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+6 more)

### Community 16 - "MonitorService"
Cohesion: 0.22
Nodes (5): MonitorService, newPingToken(), normalizeReason(), Recorder, Reloader

### Community 17 - "CheckRepo"
Cohesion: 0.16
Nodes (6): sync/atomic.Int64, MonitorSummary, CheckRepo, env, Scheduler, fakeChecker

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

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 34 - "MaintenanceModal"
Cohesion: 0.14
Nodes (20): RFC-3339, 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação (+12 more)

### Community 35 - "timezone.js"
Cohesion: 0.16
Nodes (20): pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider(), DEFAULT_TIME_ZONE, FALLBACK_ZONES (+12 more)

### Community 36 - "notifyDown"
Cohesion: 0.24
Nodes (15): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookDownAndRecoveryEventNames(), TestWebhookPayloadByType(), TestWebhookPayloadHasEventIDAndEventTime(), TestWebhookSendEmptyURLIsNoop() (+7 more)

### Community 37 - "MonitorRepo"
Cohesion: 0.18
Nodes (7): MonitorConfig, MonitorRepo, scanMonitor(), benchEnv, nameScanner, scanner, storedConfig

### Community 38 - "Monitor"
Cohesion: 0.07
Nodes (36): HTTPChecker, Registry, SnapshotCert, SnapshotIncident, SnapshotMonitor, net/http.Response, certNotAfter(), requestErrorMessage() (+28 more)

### Community 39 - "newEnv"
Cohesion: 0.15
Nodes (19): env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording() (+11 more)

### Community 40 - "Incident"
Cohesion: 0.14
Nodes (10): database/sql.NullInt64, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident(), scanIncidentWithMonitor() (+2 more)

### Community 41 - "database/sql.DB"
Cohesion: 0.15
Nodes (22): database/sql.DB, testing.TB, ExportHandler, Status(), CheckRepoTx, NewCheckRepoTx(), copyFile(), envInt() (+14 more)

### Community 42 - "ResponseChart.jsx"
Cohesion: 0.22
Nodes (14): ref_node_assert, axisTicks(), ChartLegend(), gapBuckets(), isDown(), linePath(), ResponseChart(), summary() (+6 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "context.Context"
Cohesion: 0.07
Nodes (10): context.Context, MonitorRepository, truncateError(), Scheduler, Scheduler, OutboxRepo, IncidentRepositoryTx, incidentRepoWithTx (+2 more)

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

### Community 49 - "runServe"
Cohesion: 0.14
Nodes (25): runAdd(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), github.com/spf13/cobra.Command, NewCheckRepo(), openBench() (+17 more)

### Community 50 - "newCertEnv"
Cohesion: 0.29
Nodes (14): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+6 more)

### Community 51 - "SendTestNtfy"
Cohesion: 0.18
Nodes (12): net/http.Client, newWebhookTestClient(), newNtfyClient(), ntfyErrorMessage(), publishNtfy(), SendTestNtfy(), TestSendTestNtfy(), TestSendTestNtfyTimeout() (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.1), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

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
Cohesion: 0.36
Nodes (15): NewOutboxRepoTx(), owAlways(), owDB(), owEnqueue(), owRows(), owWorker(), TestWorkerBackoffAndExpiryForDown(), TestWorkerCertNoticeLivesFor24Hours() (+7 more)

### Community 58 - "newFakeSMTP"
Cohesion: 0.14
Nodes (20): net.Conn, net.Listener, TestEmailNotifierSendsCertExpiring(), buildEmail(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage() (+12 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "Dashboard.jsx"
Cohesion: 0.06
Nodes (49): ⚠️ Notes, axios, react, react-router-dom, react-toastify, version, client, incidentsApi (+41 more)

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

### Community 72 - "SendTestWebhook"
Cohesion: 0.18
Nodes (11): updateSettingsRequest, toUpdateSettingsInput(), SendTestWebhook(), TestSendTestWebhook(), TestSendTestWebhookHasID(), UpdateSettingsInput, validateWebhookURL(), normalizeTimezone() (+3 more)

### Community 73 - "time.Duration"
Cohesion: 0.29
Nodes (8): time.Duration, rateLimiter, TestRateLimiterWindowResets(), newRateLimiter(), measure(), ms(), printTable(), timing

### Community 74 - "bodyFor"
Cohesion: 0.14
Nodes (20): time.Location, TestCertMessageText(), bodyFor(), certBody(), certPhrase(), oneLine(), splitRecipients(), subjectCert() (+12 more)

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
Cohesion: 0.09
Nodes (29): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, lucide-react, ref_node_test, CertCard(), levels (+21 more)

### Community 82 - "NewNtfyNotifier"
Cohesion: 0.28
Nodes (12): net/http.Header, TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), decodedTitle(), newFakeNtfy(), TestNtfyNotifierDown(), TestNtfyNotifierHeartbeatNeverSendsPingToken(), TestNtfyNotifierRecoveryWithoutTokenOrLink() (+4 more)

### Community 84 - "Open"
Cohesion: 0.32
Nodes (19): Open(), NewMonitorRepo(), TestPauseIsSilentAndResumeReopensIncident(), TestResumeIgnoresFailuresBeforePause(), newUoW(), pendingNotices(), failResult(), smOutboxEnv() (+11 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "Settings.jsx"
Cohesion: 0.13
Nodes (24): settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText(), OPTIONS (+16 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "3.6 Tema escuro (dark mode)"
Cohesion: 0.15
Nodes (13): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3.6.1 Tokens, 3.6.2 Papéis: claro → escuro, 3.6.3 Status no tema escuro (+5 more)

### Community 90 - "testing.T"
Cohesion: 0.11
Nodes (35): TCPChecker, net.Dialer, net/http/httptest.Server, testing.T, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry() (+27 more)

### Community 91 - "database/sql.Tx"
Cohesion: 0.11
Nodes (8): database/sql.Tx, Check, boolToInt(), encodeConfig(), nullableString(), nullableTime(), checkRepoWithTx, outboxRepoWithTx

### Community 92 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.14
Nodes (14): 2. Página de Incidentes, 3. Página de Configurações, 7. Largura no desktop + busca parcial, Decisões, Decisões, Decisões, Descartado, Descartado (+6 more)

### Community 93 - "memtest: prova da promessa de memória"
Cohesion: 0.15
Nodes (10): MVP Plan: Pingou, Pingou README, certserver: testar o aviso de certificado SSL, Como interpretar, Como usar, Como usar (Linux), Limitações, memtest: prova da promessa de memória (+2 more)

### Community 94 - "migrations_test.go"
Cohesion: 0.19
Nodes (20): runMigrateDown(), runMigrateStatus(), runMigrateUp(), go_pkg_github_com_pressly_goose_v3, Down(), ListMigrations(), OpenDB(), RunMigrations() (+12 more)

### Community 95 - "validateMonitor"
Cohesion: 0.14
Nodes (14): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+6 more)

### Community 96 - "hotQueries"
Cohesion: 0.21
Nodes (13): placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), dumpRows(), hotQueries(), legacySummariesHistory() (+5 more)

### Community 97 - "incidents_test.go"
Cohesion: 0.35
Nodes (11): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+3 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "RetentionWorker"
Cohesion: 0.21
Nodes (4): CheckRepository, NewRetentionWorker(), OutboxPurger, RetentionWorker

### Community 101 - "WebhookNotifier"
Cohesion: 0.24
Nodes (7): dashboardLink(), toWebhookMonitor(), urlErrorMessage(), webhookTestErrorMessage(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 102 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.40
Nodes (5): Como interpretar, Como usar, dbbench: desempenho do banco com volume real, O que é medido, Os dados

### Community 103 - "newTestServerEnv"
Cohesion: 0.09
Nodes (26): go_pkg_embed, go_pkg_io_fs, log/slog.Level, net/http.ServeMux, net/http.Server, Config, NtfyConfig, SMTPConfig (+18 more)

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

### Community 108 - "formatDuration"
Cohesion: 0.50
Nodes (4): formatDuration(), joinUnits(), shortDuration(), TestFormatDuration()

### Community 110 - "seed/main.go"
Cohesion: 0.49
Nodes (9): lastStep(), main(), monitorKind(), mustExec(), printSummary(), seedChecks(), seedIncidents(), seedMonitors() (+1 more)

### Community 111 - "NewOutboxRepo"
Cohesion: 0.53
Nodes (9): NewOutboxRepo(), enqueue(), newOutboxDB(), outboxEvent(), TestOutboxDueExpiresUnreadableRow(), TestOutboxEnqueueAndDue(), TestOutboxMarks(), TestOutboxPurgeDoneBefore() (+1 more)

### Community 114 - "siteCert"
Cohesion: 0.47
Nodes (6): crypto/ecdsa.PrivateKey, crypto/tls.Certificate, crypto/x509.Certificate, loadOrCreateCA(), main(), siteCert()

### Community 115 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 241 - "format.js"
Cohesion: 0.14
Nodes (30): OpenIncidentCard(), closeLabel, groupByDay(), ResolvedTable(), ConfigCard(), DownAlert(), RecentIncidents(), bucketLabel() (+22 more)

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
Nodes (31): HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu(), Badge() (+23 more)

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
- **532 isolated node(s):** `stateBadge`, `MonitorEnabled`, `createMaintenanceRequest`, `infoResponse`, `paginatedResponse` (+527 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 647 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `theme.js`, `MaintenanceModal`, `11. Fuso horário + desligamento ordenado`, `12. Índices e desempenho do banco`, `16. Nomes claros nos eventos do webhook`, `13. Aviso de expiração do certificado SSL`, `5. Remodelagem mobile`, `8. Notificação via ntfy`, `9. Eixo Y dinâmico + gráfico de pings`, `17. Retry de notificações (outbox)`, `MonitorDetail.jsx`, `6. E-mail (SMTP) e envio fora do lock`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `10. Pausar monitor fecha o incidente`, `memtest: prova da promessa de memória`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `ThemedImage()` connect `Dashboard.jsx` to `theme.js`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `serve.go`, `notifyDown`, `time.Time`, `net/http.ResponseWriter`, `MonitorRepo`, `newEnv`, `WebhookNotifier`, `bodyFor`, `Scheduler`, `context.Context`, `newMaintenanceEnv`, `MonitorService`, `CheckRepo`, `newCertEnv`, `Open`, `testing.T`, `database/sql.Tx`, `validateMonitor`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `stateBadge`, `MonitorEnabled`, `createMaintenanceRequest` to the rest of the system?**
  _532 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.11688311688311688 - nodes in this community are weakly interconnected._
- **Should `UnitOfWork` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._
- **Should `time.Time` be split into smaller, more focused modules?**
  _Cohesion score 0.0563165905631659 - nodes in this community are weakly interconnected._