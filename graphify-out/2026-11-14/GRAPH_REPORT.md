# Graph Report - pingou-health-checker  (2026-11-06)

## Corpus Check
- 243 files · ~271,602 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2295 nodes · 7698 edges · 146 communities (125 shown, 21 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 398 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `57c61968`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- theme.js
- MonitorForm.jsx
- package.json
- Plano de Desenvolvimento: Aviso de manutenção começou/terminou
- context.Context
- net/http.ResponseWriter
- Testing Strategy & Framework
- 16. Nomes claros nos eventos do webhook
- database/sql.Tx
- 🧩 Task: Backup sem configurações + importar
- testing.T
- What You Must Do When Invoked
- fakeSMTP
- Pingou Health Checker - Project Vision & Goals
- cert.js
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
- database/sql.DB
- MaintenanceConflictError
- Hero Image
- React Logo
- seed/main.go
- timezone.js
- 8. Notificação via ntfy
- Monitor
- Scheduler
- SMTPConfig
- Incident
- newMaintenanceEnv
- format.js
- graphify reference: query, path, explain
- RetentionWorker
- Pingou Health Checker - Project Roadmap
- 17. Retry de notificações (outbox)
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- Plano de Desenvolvimento: Status HTTP esperado
- newCertEnv
- UnitOfWork
- DESIGN.md — Design system do Pingou
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- owEnqueue
- NtfyConfig
- Status: **In Progress**
- MaintenanceModal.jsx
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
- NewNtfyNotifier
- 🧩 Task: Status HTTP esperado
- .Send
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- MonitorDetail.jsx
- graphify reference: transcribe video and audio
- notifyDown
- memtest.sh
- runServe
- Dependencies & Blockers
- Settings.jsx
- Scheduler
- extraction-spec.md
- 🧩 Task: Páginas públicas de status
- newHeartbeat
- 🧩 Task: Suporte a inglês (e outros idiomas)
- 2. Página de Incidentes
- WebhookNotifier
- migrations_test.go
- maintenanceEvent
- validateMonitor
- 🧩 Task: Menos RAM no cenário grande
- Status: **Planned**
- 🧩 Task: Métricas para Prometheus e Grafana
- check-tokens.sh
- Task Breakdown
- dbbench: desempenho do banco com volume real
- spa.go
- 18. Aviso de manutenção começou/terminou
- MonitorRepo
- newCertServer
- Histórico de implementações (planos aplicados)
- 3.6 Tema escuro (dark mode)
- 🧩 Task: Ver o que aconteceu com os avisos
- hotQueries
- env
- Pingou README
- 3. Cores
- siteCert
- 🧩 Task: CI no GitHub Actions
- IDEAS.md
- memtest: prova da promessa de memória
- Status de Execução
- 🧩 Task: Janela de manutenção recorrente
- 3. Página de Configurações
- time.Location
- commands/info.go
- 6. E-mail (SMTP) e envio fora do lock
- 7. Largura no desktop + busca parcial
- 5. Remodelagem mobile
- time.Time
- vite
- Status: **Planned**
- devDependencies
- react
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- Dashboard.jsx
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
1. `Monitor` - 79 edges
2. `MaintenanceWindow` - 46 edges
3. `react` - 36 edges
4. `Incident` - 36 edges
5. `newTestServer()` - 35 edges
6. `OutboxEvent` - 34 edges
7. `do()` - 33 edges
8. `runServe()` - 33 edges
9. `Scheduler` - 30 edges
10. `lucide-react` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `⚠️ Notes` --references--> `ThemedImage()`  [INFERRED]
  IDEAS.md → web/src/components/ui/ThemedImage.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Overview` --references--> `ConfigCard()`  [INFERRED]
  plan-status-http.md → web/src/components/monitor-detail/ConfigCard.jsx
- `Task T6: Linha no detalhe` --references--> `ConfigCard()`  [INFERRED]
  plan-status-http.md → web/src/components/monitor-detail/ConfigCard.jsx

## Import Cycles
- None detected.

## Communities (146 total, 21 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.17
Nodes (59): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_tls, go_pkg_crypto_x509 (+51 more)

### Community 1 - "theme.js"
Cohesion: 0.15
Nodes (20): 15. Modo escuro (Claro, Escuro, Sistema), Decisões, Descartado, Entregue, Limitações conhecidas, Validação, ref_node_fs, ThemeContext (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.16
Nodes (13): Task T5: Campo no formulário, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm(), PING_EVERY (+5 more)

### Community 3 - "package.json"
Cohesion: 0.14
Nodes (16): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+8 more)

### Community 4 - "Plano de Desenvolvimento: Aviso de manutenção começou/terminou"
Cohesion: 0.08
Nodes (25): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, File Structure, Notes, Overview (+17 more)

### Community 5 - "context.Context"
Cohesion: 0.11
Nodes (12): context.Context, database/sql.Result, MaintenanceWindow, MaintenanceWindowView, MaintenanceNoticeBuilder, MaintenanceRepository, affected(), finalizeWindowTx() (+4 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.05
Nodes (58): go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, net/http.Request, net/http.ResponseWriter, net/url.Values, time.Duration, checkResponse (+50 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "16. Nomes claros nos eventos do webhook"
Cohesion: 0.50
Nodes (4): 16. Nomes claros nos eventos do webhook, Entregue, Risco, Validação

### Community 9 - "database/sql.Tx"
Cohesion: 0.08
Nodes (15): database/sql.Tx, Check, CheckRepoTx, IncidentRepoTx, nullableTime(), MonitorRepoTx, OutboxRepoTx, NewScheduler() (+7 more)

### Community 10 - "🧩 Task: Backup sem configurações + importar"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Backup sem configurações + importar (+1 more)

### Community 11 - "testing.T"
Cohesion: 0.09
Nodes (62): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, testing.T, incidentsPage, settingsBody, TestNtfyHost(), TestMaintenanceSnapshotRoundTrip() (+54 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "fakeSMTP"
Cohesion: 0.20
Nodes (6): net.Conn, net.Listener, sync.Mutex, countingChannel, fakeSMTP, owChannel

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "cert.js"
Cohesion: 0.16
Nodes (14): ref_node_assert, ref_node_test, CertCard(), levels, CERT_ALERT_DAYS, CERT_WARN_DAYS, certDate(), certStatus() (+6 more)

### Community 16 - "MonitorService"
Cohesion: 0.08
Nodes (19): net/http.ServeMux, net/http.Server, Server, NewServer(), IncidentService, NewIncidentService(), dedupe(), MaintenanceService (+11 more)

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

### Community 30 - "database/sql.DB"
Cohesion: 0.16
Nodes (34): database/sql.DB, ExportHandler, Status(), sizes(), count(), fmtT(), noticeBuilder(), pendingFor() (+26 more)

### Community 34 - "seed/main.go"
Cohesion: 0.32
Nodes (10): go_pkg_log, lastStep(), main(), monitorKind(), mustExec(), printSummary(), seedChecks(), seedIncidents() (+2 more)

### Community 35 - "timezone.js"
Cohesion: 0.15
Nodes (21): Entregue, pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider(), DEFAULT_TIME_ZONE (+13 more)

### Community 36 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 37 - "Monitor"
Cohesion: 0.07
Nodes (34): Registry, SnapshotCert, SnapshotIncident, SnapshotMonitor, requestErrorMessage(), NewRegistry(), tcpErrorMessage(), CertAlert (+26 more)

### Community 38 - "Scheduler"
Cohesion: 0.19
Nodes (4): context.CancelFunc, sync.WaitGroup, Scheduler, job

### Community 39 - "SMTPConfig"
Cohesion: 0.14
Nodes (20): SMTPConfig, SettingsRepo, NewSettingsRepo(), TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage() (+12 more)

### Community 40 - "Incident"
Cohesion: 0.09
Nodes (15): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+7 more)

### Community 41 - "newMaintenanceEnv"
Cohesion: 0.17
Nodes (28): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceCancelInProgressSendsEndNotice(), TestMaintenanceCancelWithFlagErrorStillCancels() (+20 more)

### Community 42 - "format.js"
Cohesion: 0.10
Nodes (42): OpenIncidentCard(), closeLabel, groupByDay(), ResolvedTable(), ChecksTable(), detailOf(), ConfigCard(), DownAlert() (+34 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "RetentionWorker"
Cohesion: 0.21
Nodes (4): CheckRepository, NewRetentionWorker(), OutboxPurger, RetentionWorker

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

### Community 49 - "Plano de Desenvolvimento: Status HTTP esperado"
Cohesion: 0.15
Nodes (13): Alternativas consideradas e descartadas, Cronograma Sugerido, Definition of Done (DoD), File Structure, Notes, Overview, Phase X: Verification, Plano de Desenvolvimento: Status HTTP esperado (+5 more)

### Community 50 - "newCertEnv"
Cohesion: 0.12
Nodes (27): updateSettingsRequest, toUpdateSettingsInput(), fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation() (+19 more)

### Community 51 - "UnitOfWork"
Cohesion: 0.10
Nodes (6): OutboxRepositoryTx, UnitOfWork, erroringOutbox, failingCommit, failingEnqueueUoW, MonitorRepositoryTx

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
Cohesion: 0.39
Nodes (14): owAlways(), owDB(), owEnqueue(), owRows(), owWorker(), TestWorkerBackoffAndExpiryForDown(), TestWorkerCertNoticeLivesFor24Hours(), TestWorkerChannelPanicBecomesRetry() (+6 more)

### Community 58 - "NtfyConfig"
Cohesion: 0.16
Nodes (12): HTTPChecker, net/http.Client, NtfyConfig, newWebhookTestClient(), newNtfyClient(), ntfyErrorMessage(), publishNtfy(), SendTestNtfy() (+4 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "MaintenanceModal.jsx"
Cohesion: 0.11
Nodes (33): RFC-3339, 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação (+25 more)

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

### Community 72 - "NewNtfyNotifier"
Cohesion: 0.29
Nodes (13): net/http.Header, TestNtfyNotifierCertExpiring(), TestNtfyMaintenance(), NewNtfyNotifier(), decodedTitle(), newFakeNtfy(), TestNtfyNotifierDown(), TestNtfyNotifierHeartbeatNeverSendsPingToken() (+5 more)

### Community 73 - "🧩 Task: Status HTTP esperado"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Status HTTP esperado (+1 more)

### Community 74 - ".Send"
Cohesion: 0.13
Nodes (29): SnapshotMaintenance, TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), formatDuration(), joinUnits() (+21 more)

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
Cohesion: 0.14
Nodes (16): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, DetailTopBar(), RangeTabs(), MonitorModals(), typeLabel (+8 more)

### Community 82 - "notifyDown"
Cohesion: 0.31
Nodes (11): TestWebhookCertExpiring(), TestMaintenanceWithoutDataIsError(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookDownAndRecoveryEventNames(), TestWebhookPayloadByType(), TestWebhookSendEmptyURLIsNoop() (+3 more)

### Community 84 - "runServe"
Cohesion: 0.09
Nodes (59): runAdd(), runMigrateUp(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), github.com/spf13/cobra.Command, Open() (+51 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "Settings.jsx"
Cohesion: 0.13
Nodes (24): Task T3: Interruptor nas settings, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), MaintenanceNoticesCard(), NtfyCard() (+16 more)

### Community 87 - "Scheduler"
Cohesion: 0.22
Nodes (6): fmtKeyTime(), MaintenanceGroupKey(), TestMaintenanceGroupKey(), Scheduler, groupWindows(), windowIDs()

### Community 89 - "🧩 Task: Páginas públicas de status"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Páginas públicas de status (+1 more)

### Community 90 - "newHeartbeat"
Cohesion: 0.25
Nodes (9): env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording(), TestHeartbeatRegenerateToken() (+1 more)

### Community 91 - "🧩 Task: Suporte a inglês (e outros idiomas)"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Suporte a inglês (e outros idiomas) (+1 more)

### Community 92 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 93 - "WebhookNotifier"
Cohesion: 0.31
Nodes (5): toWebhookMonitor(), webhookMaintenancePayload, WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 94 - "migrations_test.go"
Cohesion: 0.26
Nodes (16): runMigrateDown(), go_pkg_github_com_pressly_goose_v3, Down(), OpenDB(), sqliteDSN(), assertFKOn(), assertNoFKViolations(), countRows() (+8 more)

### Community 95 - "maintenanceEvent"
Cohesion: 0.60
Nodes (5): maintenanceEvent(), TestMaintenanceTextCutsLongList(), TestWebhookMaintenanceEndedNilReasonAndNoPublicURL(), TestWebhookMaintenancePayload(), webhookBody()

### Community 96 - "validateMonitor"
Cohesion: 0.13
Nodes (15): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+7 more)

### Community 97 - "🧩 Task: Menos RAM no cenário grande"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Menos RAM no cenário grande (+1 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "🧩 Task: Métricas para Prometheus e Grafana"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Métricas para Prometheus e Grafana (+1 more)

### Community 101 - "Task Breakdown"
Cohesion: 0.25
Nodes (8): Task Breakdown, Task T1: Parser de status esperados, Task T2: Guardar na coluna `config`, Task T3: Validação e API, Task T4: Regra no checker e redirect, Task T6: Linha no detalhe, Task T7: Documentação, Task T8: QA com o binário real

### Community 102 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.25
Nodes (8): certserver: testar o aviso de certificado SSL, Como interpretar, Como usar, Como usar (Linux), dbbench: desempenho do banco com volume real, O que é medido, Os dados, scripts/

### Community 103 - "spa.go"
Cohesion: 0.40
Nodes (3): go_pkg_embed, go_pkg_io_fs, spaHandler()

### Community 104 - "18. Aviso de manutenção começou/terminou"
Cohesion: 0.33
Nodes (6): 18. Aviso de manutenção começou/terminou, Decisões, Descartado, Entregue, Limitações aceitas, Validação

### Community 105 - "MonitorRepo"
Cohesion: 0.10
Nodes (12): MonitorConfig, MonitorSummary, CheckRepo, boolToInt(), encodeConfig(), MonitorRepo, nullableString(), scanMonitor() (+4 more)

### Community 106 - "newCertServer"
Cohesion: 0.13
Nodes (18): TCPChecker, net.Dialer, net/http/httptest.Server, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), newTestRegistry() (+10 more)

### Community 107 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.33
Nodes (6): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue, Histórico de implementações (planos aplicados), Pendências abertas

### Community 108 - "3.6 Tema escuro (dark mode)"
Cohesion: 0.29
Nodes (7): 3.6.1 Tokens, 3.6.2 Papéis: claro → escuro, 3.6.3 Status no tema escuro, 3.6.4 Contraste no tema escuro, 3.6.5 Regras do tema escuro, 3.6.6 Papéis adicionais (aprovados em 2026-11-05), 3.6 Tema escuro (dark mode)

### Community 109 - "🧩 Task: Ver o que aconteceu com os avisos"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Ver o que aconteceu com os avisos (+1 more)

### Community 110 - "hotQueries"
Cohesion: 0.15
Nodes (19): testing.TB, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 111 - "env"
Cohesion: 0.22
Nodes (8): sync/atomic.Int64, env, TestDisableAndDeleteStopJob(), TestFailingMonitorDoesNotDeadlockScheduler(), TestMonitorCreatedViaRequestKeepsRunning(), TestStartLoadsEnabledMonitorsFromDB(), Scheduler, fakeChecker

### Community 113 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 114 - "siteCert"
Cohesion: 0.47
Nodes (6): crypto/ecdsa.PrivateKey, crypto/tls.Certificate, crypto/x509.Certificate, loadOrCreateCA(), main(), siteCert()

### Community 115 - "🧩 Task: CI no GitHub Actions"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: CI no GitHub Actions (+1 more)

### Community 116 - "IDEAS.md"
Cohesion: 0.40
Nodes (4): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🗑️ Descartadas, 💡 Ideias de evolução

### Community 117 - "memtest: prova da promessa de memória"
Cohesion: 0.40
Nodes (5): Como interpretar, Como usar, Limitações, memtest: prova da promessa de memória, O que ele faz

### Community 118 - "Status de Execução"
Cohesion: 0.50
Nodes (4): Concluído, Em andamento, Pendente, Status de Execução

### Community 119 - "🧩 Task: Janela de manutenção recorrente"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Janela de manutenção recorrente (+1 more)

### Community 120 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 121 - "time.Location"
Cohesion: 0.33
Nodes (3): time.Location, EmailNotifier, NtfyNotifier

### Community 122 - "commands/info.go"
Cohesion: 0.06
Nodes (37): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+29 more)

### Community 123 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 124 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 125 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 126 - "time.Time"
Cohesion: 0.05
Nodes (28): DownWindow, OutboxStatus, net/http.Response, time.Time, maintenanceListItem, certNotAfter(), MonitorStats, StatsBucket (+20 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 241 - "react"
Cohesion: 0.06
Nodes (35): axios, react, react-dom, react-router-dom, version, client, incidentsApi, monitorsApi (+27 more)

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
Cohesion: 0.40
Nodes (5): 11. Fuso horário + desligamento ordenado, Decisões, Descartado, Limitações conhecidas, Validação

### Community 262 - "12. Índices e desempenho do banco"
Cohesion: 0.33
Nodes (6): 12. Índices e desempenho do banco, Decisões, Descartado, Entregue, Limitações conhecidas, Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

### Community 265 - "13. Aviso de expiração do certificado SSL"
Cohesion: 0.29
Nodes (7): 13. Aviso de expiração do certificado SSL, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação

## Knowledge Gaps
- **601 isolated node(s):** `name`, `private`, `type`, `dev`, `build` (+596 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 716 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `ThemedImage()` connect `react` to `🧩 Task: Páginas públicas de status`, `theme.js`, `Dashboard.jsx`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `🧩 Task: Páginas públicas de status` connect `🧩 Task: Páginas públicas de status` to `IDEAS.md`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Why does `⚠️ Notes` connect `🧩 Task: Páginas públicas de status` to `react`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **What connects `name`, `private`, `type` to the rest of the system?**
  _601 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.1437908496732026 - nodes in this community are weakly interconnected._
- **Should `Plano de Desenvolvimento: Aviso de manutenção começou/terminou` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.11497584541062802 - nodes in this community are weakly interconnected._