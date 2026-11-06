# Graph Report - pingou-health-checker  (2026-11-06)

## Corpus Check
- 239 files · ~254,286 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2157 nodes · 7045 edges · 144 communities (121 shown, 23 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 346 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0b4578b9`
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
- Histórico de implementações (planos aplicados)
- Scheduler
- 🧩 Task: Backup sem configurações + importar
- testing.T
- What You Must Do When Invoked
- OutboxEvent
- Pingou Health Checker - Project Vision & Goals
- newMaintenanceEnv
- MonitorService
- Check
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
- MaintenanceModal.jsx
- timezone.js
- notifyDown
- MonitorType
- Monitor
- newEnv
- Incident
- root.go
- serve_test.go
- graphify reference: query, path, explain
- context.Context
- Pingou Health Checker - Project Roadmap
- 17. Retry de notificações (outbox)
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- runServe
- CheckResult
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- owEnqueue
- SMTPConfig
- Status: **In Progress**
- main.jsx
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
- NtfyConfig
- memtest.sh
- Open
- Dependencies & Blockers
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- extraction-spec.md
- 3.6 Tema escuro (dark mode)
- newMaintenanceServiceDB
- database/sql.Tx
- 7. Largura no desktop + busca parcial
- dbbench: desempenho do banco com volume real
- database/sql.DB
- validateMonitor
- .listIncidents
- Pingou README
- Status: **Planned**
- RetentionWorker
- check-tokens.sh
- WebhookNotifier
- .Check
- Server
- migrate.go
- 🧩 Task: Status HTTP esperado
- 🧩 Task: Páginas públicas de status
- 🧩 Task: Suporte a inglês (e outros idiomas)
- 🧩 Task: Aviso de manutenção começou/terminou
- NewTCPChecker
- seed/main.go
- NewOutboxRepo
- .handleListMonitors
- .handleMonitorStats
- siteCert
- 6. E-mail (SMTP) e envio fora do lock
- IDEAS.md
- 2. Página de Incidentes
- 3. Página de Configurações
- vite
- .checkHeartbeats
- .handleExportDB
- Config
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- Status: **Planned**
- devDependencies
- format.js
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- Dashboard.jsx
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- dependencies
- commands/info.go
- Success Metrics
- Release Schedule
- Risk Management
- 11. Fuso horário + desligamento ordenado
- 12. Índices e desempenho do banco
- 13. Aviso de expiração do certificado SSL
- middleware.go
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
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `⚠️ Notes` --references--> `ThemedImage()`  [INFERRED]
  IDEAS.md → web/src/components/ui/ThemedImage.jsx
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx

## Import Cycles
- None detected.

## Communities (144 total, 23 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.17
Nodes (56): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_tls, go_pkg_crypto_x509 (+48 more)

### Community 1 - "theme.js"
Cohesion: 0.15
Nodes (20): 15. Modo escuro (Claro, Escuro, Sistema), Decisões, Descartado, Entregue, Limitações conhecidas, Validação, ref_node_fs, ThemeContext (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.08
Nodes (27): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm() (+19 more)

### Community 3 - "package.json"
Cohesion: 0.13
Nodes (18): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+10 more)

### Community 4 - "UnitOfWork"
Cohesion: 0.10
Nodes (6): OutboxRepositoryTx, UnitOfWork, erroringOutbox, failingCommit, failingEnqueueUoW, MonitorRepositoryTx

### Community 5 - "time.Time"
Cohesion: 0.06
Nodes (26): DownWindow, database/sql.Result, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceWindowView (+18 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.25
Nodes (11): net/http.Request, net/http.ResponseWriter, Server, Server, Server, writeMaintenanceError(), Server, writeServiceError() (+3 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.14
Nodes (14): 16. Nomes claros nos eventos do webhook, 8. Notificação via ntfy, 9. Eixo Y dinâmico + gráfico de pings, Decisões, Decisões, Descartado, Descartado, Entregue (+6 more)

### Community 9 - "Scheduler"
Cohesion: 0.20
Nodes (4): context.CancelFunc, sync.WaitGroup, Scheduler, job

### Community 10 - "🧩 Task: Backup sem configurações + importar"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Backup sem configurações + importar (+1 more)

### Community 11 - "testing.T"
Cohesion: 0.09
Nodes (64): go_pkg_bufio, go_pkg_slices, net/http.Handler, net/http/httptest.ResponseRecorder, net/http/httptest.Server, testing.T, incidentsPage, settingsBody (+56 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "OutboxEvent"
Cohesion: 0.11
Nodes (8): OutboxStatus, OutboxEvent, scanOutbox(), slowChannel, Channel, OutboxStore, OutboxWorker, owChannel

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newMaintenanceEnv"
Cohesion: 0.30
Nodes (14): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+6 more)

### Community 16 - "MonitorService"
Cohesion: 0.20
Nodes (6): MonitorService, newPingToken(), normalizeReason(), CreateMonitorInput, Recorder, Reloader

### Community 17 - "Check"
Cohesion: 0.12
Nodes (6): Check, MonitorSummary, CheckRepository, CheckRepo, CheckRepositoryTx, checkRepoWithTx

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

### Community 34 - "MaintenanceModal.jsx"
Cohesion: 0.11
Nodes (31): RFC-3339, 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação (+23 more)

### Community 35 - "timezone.js"
Cohesion: 0.16
Nodes (20): pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider(), DEFAULT_TIME_ZONE, FALLBACK_ZONES (+12 more)

### Community 36 - "notifyDown"
Cohesion: 0.15
Nodes (21): CertAlert, NewCertSnapshot(), EvaluateCertAlert(), sameInstant(), TestEvaluateCertAlert(), TestEvaluateCertAlertRepeated(), TestWebhookCertExpiring(), NewWebhookNotifier() (+13 more)

### Community 37 - "MonitorType"
Cohesion: 0.12
Nodes (13): Registry, newTestRegistry(), TestRegistryHTTP(), TestRegistryUnknownType(), NewRegistry(), Checker, MonitorType, MonitorFilter (+5 more)

### Community 38 - "Monitor"
Cohesion: 0.13
Nodes (18): SnapshotCert, SnapshotIncident, SnapshotMonitor, Monitor, MonitorConfig, MonitorState, EventSnapshot, NewDownSnapshot() (+10 more)

### Community 39 - "newEnv"
Cohesion: 0.11
Nodes (23): sync/atomic.Int64, env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown() (+15 more)

### Community 40 - "Incident"
Cohesion: 0.12
Nodes (14): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+6 more)

### Community 41 - "root.go"
Cohesion: 0.22
Nodes (4): Execute(), requireKey(), main(), go_pkg_time_tzdata

### Community 42 - "serve_test.go"
Cohesion: 0.24
Nodes (9): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep (+1 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "context.Context"
Cohesion: 0.09
Nodes (8): context.Context, MonitorRepository, truncateError(), Scheduler, OutboxRepo, IncidentRepositoryTx, incidentRepoWithTx, monitorRepoWithTx

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
Cohesion: 0.09
Nodes (43): runAdd(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), github.com/spf13/cobra.Command, newTestServerEnv(), NewCheckRepo() (+35 more)

### Community 50 - "CheckResult"
Cohesion: 0.27
Nodes (15): CheckResult, fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal() (+7 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

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
Cohesion: 0.32
Nodes (16): NewOutboxWorker(), owAlways(), owDB(), owEnqueue(), owRows(), owWorker(), TestWorkerBackoffAndExpiryForDown(), TestWorkerCertNoticeLivesFor24Hours() (+8 more)

### Community 58 - "SMTPConfig"
Cohesion: 0.16
Nodes (19): net.Conn, net.Listener, SMTPConfig, TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage() (+11 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.06
Nodes (32): 5. Remodelagem mobile, Decisões, Descartado, Entregue, axios, react-router-dom, client, incidentsApi (+24 more)

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
Cohesion: 0.12
Nodes (16): updateSettingsRequest, toUpdateSettingsInput(), SendTestWebhook(), TestSendTestWebhook(), TestSendTestWebhookHasID(), TestValidNtfyTopic(), validNtfyTopic(), UpdateSettingsInput (+8 more)

### Community 73 - "time.Duration"
Cohesion: 0.15
Nodes (13): sync.Mutex, time.Duration, rateLimiter, TestRateLimiterWindowResets(), Server, clientIP(), newRateLimiter(), rejectRateLimited() (+5 more)

### Community 74 - "bodyFor"
Cohesion: 0.10
Nodes (27): time.Location, TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), formatDuration(), joinUnits() (+19 more)

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
Cohesion: 0.12
Nodes (19): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, ChecksTable(), detailOf(), ConfigCard(), DetailTopBar() (+11 more)

### Community 82 - "NtfyConfig"
Cohesion: 0.14
Nodes (21): net/http.Client, net/http.Header, NtfyConfig, TestNtfyNotifierCertExpiring(), newWebhookTestClient(), newNtfyClient(), NewNtfyNotifier(), publishNtfy() (+13 more)

### Community 84 - "Open"
Cohesion: 0.33
Nodes (19): Open(), NewMonitorRepo(), TestPauseIsSilentAndResumeReopensIncident(), TestResumeIgnoresFailuresBeforePause(), newUoW(), pendingNotices(), failResult(), smOutboxEnv() (+11 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "react"
Cohesion: 0.13
Nodes (25): Entregue, react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard() (+17 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "3.6 Tema escuro (dark mode)"
Cohesion: 0.15
Nodes (13): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3.6.1 Tokens, 3.6.2 Papéis: claro → escuro, 3.6.3 Status no tema escuro (+5 more)

### Community 90 - "newMaintenanceServiceDB"
Cohesion: 0.46
Nodes (8): NewMaintenanceRepo(), addWindowRaw(), newMaintenanceServiceDB(), TestMaintenanceEditedIntervalIsUsed(), TestMaintenanceList(), TestMaintenanceUpdateRaceWithStart(), TestMaintenanceUpdateScheduled(), TestMaintenanceUpdateStartedLocked()

### Community 91 - "database/sql.Tx"
Cohesion: 0.15
Nodes (6): database/sql.Tx, boolToInt(), encodeConfig(), nullableString(), nullableTime(), outboxRepoWithTx

### Community 92 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 93 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.15
Nodes (13): certserver: testar o aviso de certificado SSL, Como interpretar, Como interpretar, Como usar, Como usar, Como usar (Linux), dbbench: desempenho do banco com volume real, Limitações (+5 more)

### Community 94 - "database/sql.DB"
Cohesion: 0.09
Nodes (41): runMigrateDown(), runMigrateUp(), go_pkg_github_com_pressly_goose_v3, database/sql.DB, testing.TB, ExportHandler, Down(), OpenDB() (+33 more)

### Community 95 - "validateMonitor"
Cohesion: 0.13
Nodes (15): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+7 more)

### Community 96 - ".listIncidents"
Cohesion: 0.16
Nodes (10): net/url.Values, checkResponse, incidentMonitorResponse, incidentResponse, Server, toCheckResponse(), Server, incidentStatus() (+2 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "RetentionWorker"
Cohesion: 0.23
Nodes (4): SettingsRepo, NewRetentionWorker(), OutboxPurger, RetentionWorker

### Community 101 - "WebhookNotifier"
Cohesion: 0.24
Nodes (7): dashboardLink(), toWebhookMonitor(), urlErrorMessage(), webhookTestErrorMessage(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 102 - ".Check"
Cohesion: 0.40
Nodes (4): HTTPChecker, net/http.Response, certNotAfter(), requestErrorMessage()

### Community 103 - "Server"
Cohesion: 0.22
Nodes (8): go_pkg_embed, go_pkg_io_fs, net/http.ServeMux, net/http.Server, Server, NewServer(), spaHandler(), IncidentService

### Community 104 - "migrate.go"
Cohesion: 0.22
Nodes (5): copyFile(), runExportDB(), go_pkg_github_com_joho_godotenv, go_pkg_github_com_spf13_cobra, go_pkg_log

### Community 105 - "🧩 Task: Status HTTP esperado"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Status HTTP esperado (+1 more)

### Community 106 - "🧩 Task: Páginas públicas de status"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Páginas públicas de status (+1 more)

### Community 107 - "🧩 Task: Suporte a inglês (e outros idiomas)"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Suporte a inglês (e outros idiomas) (+1 more)

### Community 108 - "🧩 Task: Aviso de manutenção começou/terminou"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Aviso de manutenção começou/terminou (+1 more)

### Community 109 - "NewTCPChecker"
Cohesion: 0.24
Nodes (9): TCPChecker, net.Dialer, tcpMonitor(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused(), TestTCPCheckerTimeout(), TestTCPCheckerUp(), NewTCPChecker() (+1 more)

### Community 110 - "seed/main.go"
Cohesion: 0.49
Nodes (9): lastStep(), main(), monitorKind(), mustExec(), printSummary(), seedChecks(), seedIncidents(), seedMonitors() (+1 more)

### Community 111 - "NewOutboxRepo"
Cohesion: 0.53
Nodes (9): NewOutboxRepo(), enqueue(), newOutboxDB(), outboxEvent(), TestOutboxDueExpiresUnreadableRow(), TestOutboxEnqueueAndDue(), TestOutboxMarks(), TestOutboxPurgeDoneBefore() (+1 more)

### Community 112 - ".handleListMonitors"
Cohesion: 0.36
Nodes (6): maintenanceResponse, monitorListItem, monitorResponse, toMaintenanceResponse(), formatTime(), toMonitorResponse()

### Community 113 - ".handleMonitorStats"
Cohesion: 0.29
Nodes (6): downWindowResponse, statsBucketResponse, statsResponse, Server, roundMs(), toStatsResponse()

### Community 114 - "siteCert"
Cohesion: 0.47
Nodes (6): crypto/ecdsa.PrivateKey, crypto/tls.Certificate, crypto/x509.Certificate, loadOrCreateCA(), main(), siteCert()

### Community 115 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.50
Nodes (4): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Limitações conhecidas

### Community 116 - "IDEAS.md"
Cohesion: 0.40
Nodes (4): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🗑️ Descartadas, 💡 Ideias de evolução

### Community 117 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 118 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 235 - "Config"
Cohesion: 0.20
Nodes (12): runMigrateStatus(), log/slog.Level, Config, Load(), LoadConfig(), loadConfigWithAPIKey(), LoadNtfy(), LoadSMTP() (+4 more)

### Community 236 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 241 - "format.js"
Cohesion: 0.10
Nodes (42): ref_node_assert, ref_node_test, OpenIncidentCard(), closeLabel, groupByDay(), ResolvedTable(), DownAlert(), RecentIncidents() (+34 more)

### Community 242 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 4.1: Multi-tenancy & Access Control, Milestone 4.2: Advanced Authentication, Milestone 4.3: Scalability Improvements, Phase 4: Enterprise Features (Q4 2024), Status: **Planned**

### Community 243 - "Status: **Future**"
Cohesion: 0.40
Nodes (5): Milestone 5.1: Advanced Monitoring, Milestone 5.2: Integration Ecosystem, Milestone 5.3: Public Features, Phase 5: Advanced Features & Integrations (Q1-Q2 2025), Status: **Future**

### Community 244 - "Status: **Continuous**"
Cohesion: 0.40
Nodes (5): Milestone 6.1: Community Building, Milestone 6.2: Developer Experience, Milestone 6.3: Ecosystem Integration, Phase 6: Ecosystem & Community (Ongoing), Status: **Continuous**

### Community 245 - "Dashboard.jsx"
Cohesion: 0.11
Nodes (37): lucide-react, HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu() (+29 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.22
Nodes (9): scripts, build, check:tokens, dev, format, format:check, lint, preview (+1 more)

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 252 - "commands/info.go"
Cohesion: 0.39
Nodes (8): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim()

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

### Community 267 - "middleware.go"
Cohesion: 0.19
Nodes (10): go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, responseWriter, constantTimeEqual(), Server, loggingMiddleware(), newRequestID() (+2 more)

## Knowledge Gaps
- **550 isolated node(s):** `Decisões`, `Limitações conhecidas`, `Descartado`, `Entregue`, `Decisões` (+545 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 667 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `UnitOfWork`, `time.Time`, `Scheduler`, `newMaintenanceEnv`, `MonitorService`, `Check`, `notifyDown`, `MonitorType`, `newEnv`, `context.Context`, `CheckResult`, `bodyFor`, `Open`, `database/sql.Tx`, `validateMonitor`, `WebhookNotifier`, `.Check`, `NewTCPChecker`, `.handleListMonitors`, `.checkHeartbeats`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `Pingou README`, `MaintenanceModal.jsx`, `theme.js`, `11. Fuso horário + desligamento ordenado`, `12. Índices e desempenho do banco`, `13. Aviso de expiração do certificado SSL`, `17. Retry de notificações (outbox)`, `7. Largura no desktop + busca parcial`, `MonitorDetail.jsx`, `6. E-mail (SMTP) e envio fora do lock`, `2. Página de Incidentes`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `3. Página de Configurações`, `10. Pausar monitor fecha o incidente`, `main.jsx`?**
  _High betweenness centrality (0.015) - this node is a cross-community bridge._
- **Why does `ThemedImage()` connect `main.jsx` to `theme.js`, `🧩 Task: Páginas públicas de status`, `Dashboard.jsx`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **What connects `Decisões`, `Limitações conhecidas`, `Descartado` to the rest of the system?**
  _550 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07557354925775979 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.12631578947368421 - nodes in this community are weakly interconnected._
- **Should `UnitOfWork` be split into smaller, more focused modules?**
  _Cohesion score 0.09523809523809523 - nodes in this community are weakly interconnected._