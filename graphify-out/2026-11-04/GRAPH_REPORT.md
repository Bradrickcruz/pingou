# Graph Report - pingou-health-checker  (2026-11-04)

## Corpus Check
- 221 files · ~238,429 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2008 nodes · 6180 edges · 132 communities (113 shown, 19 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 306 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1bd0527a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- Task Breakdown
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- What You Must Do When Invoked
- net/http.ResponseWriter
- Testing Strategy & Framework
- MonitorDetail.jsx
- Scheduler
- 🧩 Task: Aproveitar melhor a largura no desktop
- newTestServer
- ResponseChart.jsx
- time.Time
- Pingou Health Checker - Project Vision & Goals
- validateMonitor
- Monitor
- dbbench: desempenho do banco com volume real
- Vite Logo
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Code Conventions & Standards
- Architecture Overview
- Technical Concerns & Risk Assessment
- Integrations & External Dependencies
- Technology Stack Analysis
- Project Structure Analysis
- MonitorRepo
- Plano de Desenvolvimento: {{Titulo}}
- MaintenanceConflictError
- Hero Image
- React Logo
- Dashboard.jsx
- graphify reference: extra exports and benchmark
- CheckResult
- .NotifyRecovery
- Incident
- testing.T
- Check
- MonitorService
- format.js
- newFakeSMTP
- vite
- Pingou Health Checker - Project Roadmap
- graphify reference: query, path, explain
- Pingou - Divergencias entre PRD e implementacao
- Dependencies & Blockers
- 14. Manutenção: pausa com motivo e janela agendada
- newTestRegistry
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- MonitorType
- NewIncidentRepo
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- 6. E-mail (SMTP) e envio fora do lock
- graphify reference: add a URL and watch a folder
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
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- 3. Cores
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- database/sql.DB
- graphify reference: GitHub clone and cross-repo merge
- graphify reference: transcribe video and audio
- memtest.sh
- extraction-spec.md
- newCertServer
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- SendTestWebhook
- runServe
- newEnv
- SendTestNtfy
- incidents_test.go
- NewWebhookNotifier
- hotQueries
- serve_test.go
- context.Context
- NewNtfyNotifier
- Status: **Planned**
- Histórico de implementações (planos aplicados)
- commands/info.go
- SendTestEmail
- memtest: prova da promessa de memória
- 2. Página de Incidentes
- 9. Eixo Y dinâmico + gráfico de pings
- spa.go
- migrations_test.go
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- Status: **Planned**
- devDependencies
- TestDBBench
- timezone.js
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- lib/monitors.js
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- dependencies
- config.go
- migrate.go
- Success Metrics
- Release Schedule
- Risk Management
- 11. Fuso horário + desligamento ordenado
- 12. Índices e desempenho do banco
- 13. Aviso de expiração do certificado SSL
- siteCert
- dbbench.sh
- dbmem.sh
- 8. Notificação via ntfy
- NewMultiNotifier

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 112 edges
2. `Incident` - 61 edges
3. `newTestServer()` - 33 edges
4. `MaintenanceWindow` - 32 edges
5. `do()` - 31 edges
6. `react` - 31 edges
7. `runServe()` - 31 edges
8. `Scheduler` - 29 edges
9. `newEnv()` - 29 edges
10. `lucide-react` - 27 edges

## Surprising Connections (you probably didn't know these)
- `Risk Areas` --references--> `zonedToUtc()`  [INFERRED]
  plan-maintenance-windows.md → web/src/lib/timezone.js
- `Desvios do plano` --references--> `MaintenanceModal()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitors/MaintenanceModal.jsx
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js

## Import Cycles
- None detected.

## Communities (132 total, 19 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.15
Nodes (63): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_tls, go_pkg_crypto_x509 (+55 more)

### Community 1 - "Task Breakdown"
Cohesion: 0.06
Nodes (35): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, File Structure, Notes, Overview (+27 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.10
Nodes (23): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm() (+15 more)

### Community 3 - "package.json"
Cohesion: 0.13
Nodes (18): autoprefixer, axios, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+10 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (43): go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, net/http.Request, net/http.ResponseWriter, net/url.Values, checkResponse, incidentMonitorResponse (+35 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.13
Nodes (18): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, ChecksTable(), detailOf(), DetailTopBar(), RangeTabs() (+10 more)

### Community 9 - "Scheduler"
Cohesion: 0.24
Nodes (3): context.CancelFunc, Scheduler, job

### Community 10 - "🧩 Task: Aproveitar melhor a largura no desktop"
Cohesion: 0.05
Nodes (37): 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🎯 Goal, 🎯 Goal, 🎯 Goal, 🎯 Goal (+29 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (40): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+32 more)

### Community 12 - "ResponseChart.jsx"
Cohesion: 0.27
Nodes (12): ref_node_assert, axisTicks(), ChartLegend(), gapBuckets(), isDown(), linePath(), ResponseChart(), summary() (+4 more)

### Community 13 - "time.Time"
Cohesion: 0.07
Nodes (22): DownWindow, time.Time, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceRepository, fmtTime(), MaintenanceRepo (+14 more)

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "validateMonitor"
Cohesion: 0.17
Nodes (12): TestValidateMonitorByType(), TestValidateURL(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval(), validateMonitor(), validateName() (+4 more)

### Community 16 - "Monitor"
Cohesion: 0.07
Nodes (18): sync/atomic.Int64, sync.WaitGroup, noopNotifier, CertAlert, Monitor, MonitorState, Notifier, effectiveThreshold() (+10 more)

### Community 17 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.15
Nodes (10): MVP Plan: Pingou, Pingou README, certserver: testar o aviso de certificado SSL, Como interpretar, Como usar, Como usar (Linux), dbbench: desempenho do banco com volume real, O que é medido (+2 more)

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

### Community 29 - "MonitorRepo"
Cohesion: 0.18
Nodes (6): MonitorConfig, MonitorFilter, MonitorRepo, scanMonitor(), scanner, storedConfig

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 34 - "Dashboard.jsx"
Cohesion: 0.11
Nodes (24): lucide-react, react-toastify, client, incidentsApi, maintenanceApi, monitorsApi, web_src_assets_brand_pingou_simbolo, MobileHeader() (+16 more)

### Community 35 - "graphify reference: extra exports and benchmark"
Cohesion: 0.22
Nodes (8): graphify reference: extra exports and benchmark, Step 6b - Wiki (only if --wiki flag), Step 7 - Neo4j export (only if --neo4j or --neo4j-push flag), Step 7a - FalkorDB export (only if --falkordb or --falkordb-push flag), Step 7b - SVG export (only if --svg flag), Step 7c - GraphML export (only if --graphml flag), Step 7d - MCP server (only if --mcp flag), Step 8 - Token reduction benchmark (only if total_words > 5000)

### Community 36 - "CheckResult"
Cohesion: 0.27
Nodes (16): CheckResult, TestWaitRespectsDeadline(), fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation() (+8 more)

### Community 37 - ".NotifyRecovery"
Cohesion: 0.47
Nodes (4): toWebhookMonitor(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 38 - "Incident"
Cohesion: 0.11
Nodes (14): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+6 more)

### Community 39 - "testing.T"
Cohesion: 0.17
Nodes (27): testing.T, TestNtfyHost(), TestMaintenanceRoutesNeedAPIKey(), env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel() (+19 more)

### Community 40 - "Check"
Cohesion: 0.05
Nodes (14): database/sql.Tx, Check, IncidentRepoTx, boolToInt(), encodeConfig(), nullableString(), nullableTime(), UnitOfWork (+6 more)

### Community 41 - "MonitorService"
Cohesion: 0.15
Nodes (12): net/http.ServeMux, net/http.Server, Server, NewServer(), CheckRepo, env, IncidentService, MonitorService (+4 more)

### Community 42 - "format.js"
Cohesion: 0.12
Nodes (33): Task T6: Modal de pausa + selos com motivo, OpenIncidentCard(), closeLabel, groupByDay(), ResolvedTable(), ConfigCard(), DownAlert(), RecentIncidents() (+25 more)

### Community 43 - "newFakeSMTP"
Cohesion: 0.23
Nodes (13): net.Conn, net.Listener, TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), newFakeSMTP(), TestEmailNotifierSendsDown(), TestEmailNotifierSkipsWithoutRecipients() (+5 more)

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 46 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 49 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.33
Nodes (6): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Desvios do plano, Limitações conhecidas, Validação

### Community 50 - "newTestRegistry"
Cohesion: 0.19
Nodes (12): TCPChecker, net.Dialer, newTestRegistry(), tcpMonitor(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused() (+4 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

### Community 53 - "MonitorType"
Cohesion: 0.14
Nodes (10): Registry, NewRegistry(), Checker, MonitorType, newPingToken(), isHTTPS(), normalizeReason(), validateCertAlert() (+2 more)

### Community 54 - "NewIncidentRepo"
Cohesion: 0.26
Nodes (13): runAdd(), runMigrateDown(), runMigrateStatus(), runRM(), github.com/spf13/cobra.Command, Down(), ListMigrations(), NewCheckRepo() (+5 more)

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 56 - "10. Pausar monitor fecha o incidente"
Cohesion: 0.40
Nodes (5): 10. Pausar monitor fecha o incidente, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 57 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 58 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.11
Nodes (16): react-router-dom, version, web_src_assets_brand_pingou_marca_horizontal, web_src_assets_brand_pingou_marca_horizontal_negativa, BottomNav(), nav, Shell(), useOpenIncidents() (+8 more)

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

### Community 71 - "graphify reference: commit hook and native CLAUDE.md integration"
Cohesion: 0.50
Nodes (3): For git commit hook, For native CLAUDE.md integration, graphify reference: commit hook and native CLAUDE.md integration

### Community 72 - "graphify reference: incremental update and cluster-only"
Cohesion: 0.50
Nodes (3): For --cluster-only, For --update (incremental re-extraction), graphify reference: incremental update and cluster-only

### Community 73 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 74 - "bodyFor"
Cohesion: 0.13
Nodes (24): time.Location, TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), oneLine(), splitRecipients() (+16 more)

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

### Community 80 - "database/sql.DB"
Cohesion: 0.37
Nodes (12): database/sql.DB, ExportHandler, Status(), lastStep(), main(), monitorKind(), mustExec(), printSummary() (+4 more)

### Community 85 - "newCertServer"
Cohesion: 0.18
Nodes (11): HTTPChecker, net/http/httptest.Server, net/http.Response, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage() (+3 more)

### Community 86 - "react"
Cohesion: 0.17
Nodes (20): react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+12 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 88 - "SendTestWebhook"
Cohesion: 0.18
Nodes (11): updateSettingsRequest, toUpdateSettingsInput(), SendTestWebhook(), TestSendTestWebhook(), webhookTestErrorMessage(), UpdateSettingsInput, validateWebhookURL(), normalizeTimezone() (+3 more)

### Community 89 - "runServe"
Cohesion: 0.19
Nodes (30): forceExitOnSecondSignal(), releaseLock(), runServe(), newTestServerEnv(), CheckRepoTx, NewCheckRepoTx(), NewIncidentRepoTx(), NewMaintenanceRepo() (+22 more)

### Community 90 - "newEnv"
Cohesion: 0.31
Nodes (11): env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording() (+3 more)

### Community 91 - "SendTestNtfy"
Cohesion: 0.20
Nodes (11): net/http.Client, newWebhookTestClient(), newNtfyClient(), ntfyErrorMessage(), publishNtfy(), SendTestNtfy(), TestSendTestNtfyTimeout(), TestValidNtfyTopic() (+3 more)

### Community 92 - "incidents_test.go"
Cohesion: 0.35
Nodes (11): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+3 more)

### Community 93 - "NewWebhookNotifier"
Cohesion: 0.50
Nodes (5): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType()

### Community 94 - "hotQueries"
Cohesion: 0.15
Nodes (19): testing.TB, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 95 - "serve_test.go"
Cohesion: 0.24
Nodes (9): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep (+1 more)

### Community 96 - "context.Context"
Cohesion: 0.07
Nodes (12): context.Context, MonitorSummary, CheckRepository, MonitorRepository, SettingsRepo, Scheduler, Scheduler, NewRetentionWorker() (+4 more)

### Community 97 - "NewNtfyNotifier"
Cohesion: 0.22
Nodes (15): net/http.Header, sync.Mutex, TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), decodedTitle(), newFakeNtfy(), TestMultiNotifierNtfyFailureDoesNotBlockOthers(), TestNtfyNotifierDown() (+7 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.20
Nodes (10): 5. Remodelagem mobile, 7. Largura no desktop + busca parcial, Decisões, Decisões, Descartado, Descartado, Entregue, Entregue (+2 more)

### Community 100 - "commands/info.go"
Cohesion: 0.39
Nodes (8): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim()

### Community 101 - "SendTestEmail"
Cohesion: 0.40
Nodes (5): SendTestEmail(), smtpErrorMessage(), parseRecipients(), TestParseRecipients(), EmailTestResult

### Community 102 - "memtest: prova da promessa de memória"
Cohesion: 0.40
Nodes (5): Como interpretar, Como usar, Limitações, memtest: prova da promessa de memória, O que ele faz

### Community 103 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 104 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 105 - "spa.go"
Cohesion: 0.50
Nodes (3): go_pkg_embed, go_pkg_io_fs, spaHandler()

### Community 235 - "migrations_test.go"
Cohesion: 0.25
Nodes (16): runMigrateUp(), go_pkg_github_com_pressly_goose_v3, Open(), OpenDB(), RunMigrations(), sqliteDSN(), assertFKOn(), assertNoFKViolations() (+8 more)

### Community 236 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 240 - "TestDBBench"
Cohesion: 0.16
Nodes (14): time.Duration, rateLimiter, TestRateLimiterWindowResets(), newRateLimiter(), envInt(), measure(), ms(), printTable() (+6 more)

### Community 241 - "timezone.js"
Cohesion: 0.09
Nodes (35): RFC-3339, Entregue, Task T16: `zonedToUtc` + status MAINTENANCE, Task T17: Modal de nova janela, Tech Stack, ref_node_test, CloseNote(), MaintenanceModal() (+27 more)

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
Cohesion: 0.11
Nodes (35): Task T18: Selos, card de detalhes, cancelar, MaintenanceCard(), HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon (+27 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.25
Nodes (8): scripts, build, dev, format, format:check, lint, preview, test

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 250 - "config.go"
Cohesion: 0.19
Nodes (14): log/slog.Level, Config, NtfyConfig, SMTPConfig, Load(), LoadConfig(), loadConfigWithAPIKey(), LoadNtfy() (+6 more)

### Community 252 - "migrate.go"
Cohesion: 0.12
Nodes (9): copyFile(), runExportDB(), Execute(), requireKey(), main(), go_pkg_github_com_joho_godotenv, go_pkg_github_com_spf13_cobra, go_pkg_log (+1 more)

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

### Community 267 - "siteCert"
Cohesion: 0.47
Nodes (6): crypto/ecdsa.PrivateKey, crypto/tls.Certificate, crypto/x509.Certificate, loadOrCreateCA(), main(), siteCert()

### Community 272 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 273 - "NewMultiNotifier"
Cohesion: 0.31
Nodes (7): TestMultiNotifierCertExpiringFanOut(), NewMultiNotifier(), TestMultiNotifierDoesNotBlock(), TestMultiNotifierIgnoresCanceledContext(), TestMultiNotifierIsolatesFailures(), TestMultiNotifierWait(), ctxNotifier

## Knowledge Gaps
- **547 isolated node(s):** `Overview`, `Project Type`, `Definition of Done (DoD)`, `Success Criteria`, `File Structure` (+542 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 649 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `context.Context`, `serve.go`, `CheckResult`, `.NotifyRecovery`, `net/http.ResponseWriter`, `testing.T`, `Check`, `Scheduler`, `bodyFor`, `Incident`, `time.Time`, `validateMonitor`, `NewMultiNotifier`, `newTestRegistry`, `MonitorType`, `newCertServer`, `newEnv`, `MonitorRepo`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `3. Página de Configurações`, `11. Fuso horário + desligamento ordenado`, `12. Índices e desempenho do banco`, `2. Página de Incidentes`, `MonitorDetail.jsx`, `13. Aviso de expiração do certificado SSL`, `9. Eixo Y dinâmico + gráfico de pings`, `8. Notificação via ntfy`, `dbbench: desempenho do banco com volume real`, `14. Manutenção: pausa com motivo e janela agendada`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `10. Pausar monitor fecha o incidente`, `6. E-mail (SMTP) e envio fora do lock`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `react` connect `react` to `Dashboard.jsx`, `package.json`, `MonitorForm.jsx`, `MonitorDetail.jsx`, `format.js`, `ResponseChart.jsx`, `timezone.js`, `lib/monitors.js`, `main.jsx`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `Overview`, `Project Type`, `Definition of Done (DoD)` to the rest of the system?**
  _547 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1457464175584981 - nodes in this community are weakly interconnected._
- **Should `Task Breakdown` be split into smaller, more focused modules?**
  _Cohesion score 0.05555555555555555 - nodes in this community are weakly interconnected._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0967741935483871 - nodes in this community are weakly interconnected._