# Graph Report - pingou-health-checker  (2026-11-04)

## Corpus Check
- 223 files · ~238,979 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 1994 nodes · 6271 edges · 132 communities (110 shown, 22 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3d3e86f5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- Histórico de implementações (planos aplicados)
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- time.Time
- net/http.ResponseWriter
- Testing Strategy & Framework
- MonitorDetail.jsx
- Scheduler
- 🧩 Task: Backup sem configurações + importar
- newTestServer
- What You Must Do When Invoked
- UnitOfWork
- Pingou Health Checker - Project Vision & Goals
- newEnv
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
- graphify reference: extra exports and benchmark
- Plano de Desenvolvimento: {{Titulo}}
- MaintenanceConflictError
- Hero Image
- React Logo
- Dashboard.jsx
- timezone.js
- CheckResult
- newTestServerEnv
- Incident
- newMaintenanceEnv
- context.Context
- MonitorService
- SendTestWebhook
- graphify reference: query, path, explain
- vite
- Pingou Health Checker - Project Roadmap
- seed/main.go
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- 14. Manutenção: pausa com motivo e janela agendada
- NewTCPChecker
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- 2. Página de Incidentes
- newFakeSMTP
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
- database/sql.Tx
- 3. Cores
- NtfyNotifier
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- EmailNotifier
- graphify reference: transcribe video and audio
- SendTestNtfy
- memtest.sh
- TestDBBench
- Dependencies & Blockers
- Settings.jsx
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- extraction-spec.md
- database/sql.DB
- testing.T
- root.go
- incidents_test.go
- SendTestEmail
- hotQueries
- runServe
- NewMultiNotifier
- NewNtfyNotifier
- Status: **Planned**
- 5. Remodelagem mobile
- newCertServer
- commands/info.go
- Pingou README
- NewWebhookNotifier
- StateMachine
- 9. Eixo Y dinâmico + gráfico de pings
- 7. Largura no desktop + busca parcial
- bodyFor
- .checkHeartbeats
- .NotifyRecovery
- migrations_test.go
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
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
- serve_test.go
- Success Metrics
- Release Schedule
- Risk Management
- 11. Fuso horário + desligamento ordenado
- 12. Índices e desempenho do banco
- 13. Aviso de expiração do certificado SSL
- siteCert
- dbbench.sh
- dbmem.sh

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 112 edges
2. `Incident` - 61 edges
3. `MaintenanceWindow` - 36 edges
4. `newTestServer()` - 34 edges
5. `react` - 33 edges
6. `do()` - 32 edges
7. `runServe()` - 31 edges
8. `newEnv()` - 29 edges
9. `lucide-react` - 29 edges
10. `Scheduler` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Desvios do plano` --references--> `CancelWindowDialog()`  [INFERRED]
  docs/implementacoes.md → web/src/components/maintenance/CancelWindowDialog.jsx

## Import Cycles
- None detected.

## Communities (132 total, 22 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.14
Nodes (59): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_tls, go_pkg_crypto_x509, go_pkg_crypto_x509_pkix (+51 more)

### Community 1 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.20
Nodes (10): 6. E-mail (SMTP) e envio fora do lock, 8. Notificação via ntfy, Decisões, Decisões, Descartado, Descartado, Entregue, Histórico de implementações (planos aplicados) (+2 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.10
Nodes (23): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm() (+15 more)

### Community 3 - "package.json"
Cohesion: 0.13
Nodes (18): autoprefixer, axios, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+10 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "time.Time"
Cohesion: 0.07
Nodes (25): DownWindow, database/sql.Result, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceWindowView (+17 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (49): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, net/http.Request, net/http.ResponseWriter, net/url.Values, checkResponse (+41 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.16
Nodes (14): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, DetailTopBar(), RangeTabs(), useSettings(), clampRange() (+6 more)

### Community 9 - "Scheduler"
Cohesion: 0.13
Nodes (6): context.CancelFunc, sync.WaitGroup, Scheduler, Notifier, job, MultiNotifier

### Community 10 - "🧩 Task: Backup sem configurações + importar"
Cohesion: 0.09
Nodes (22): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🔍 Context, 🔍 Context, 🗑️ Descartadas, 🎯 Goal, 🎯 Goal, 💡 Idea (+14 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (41): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+33 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "UnitOfWork"
Cohesion: 0.11
Nodes (4): UnitOfWork, failingCommit, MonitorRepositoryTx, sqliteUnitOfWork

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newEnv"
Cohesion: 0.18
Nodes (20): runAdd(), runRM(), NewCheckRepo(), openBench(), TestResultsSameWithOldIndex(), NewIncidentRepo(), env, newEnv() (+12 more)

### Community 16 - "Monitor"
Cohesion: 0.06
Nodes (16): sync/atomic.Int64, noopNotifier, CertAlert, Monitor, MonitorRepository, TestMultiNotifierIsolatesFailures(), noopNotifier, dbNotifier (+8 more)

### Community 17 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.15
Nodes (13): certserver: testar o aviso de certificado SSL, Como interpretar, Como interpretar, Como usar, Como usar, Como usar (Linux), dbbench: desempenho do banco com volume real, Limitações (+5 more)

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

### Community 34 - "Dashboard.jsx"
Cohesion: 0.11
Nodes (28): react, react-toastify, maintenanceApi, monitorsApi, web_src_assets_brand_pingou_simbolo, MobileHeader(), CancelWindowDialog(), MaintenanceCard() (+20 more)

### Community 35 - "timezone.js"
Cohesion: 0.09
Nodes (37): RFC-3339, Desvios do plano, Entregue, ref_node_assert, ref_node_test, CloseNote(), stateBadge, WindowsTable() (+29 more)

### Community 36 - "CheckResult"
Cohesion: 0.29
Nodes (15): CheckResult, fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal() (+7 more)

### Community 37 - "newTestServerEnv"
Cohesion: 0.12
Nodes (15): go_pkg_embed, go_pkg_io_fs, net/http.ServeMux, net/http.Server, NtfyConfig, SMTPConfig, newTestServerEnv(), Server (+7 more)

### Community 38 - "Incident"
Cohesion: 0.11
Nodes (15): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+7 more)

### Community 39 - "newMaintenanceEnv"
Cohesion: 0.23
Nodes (15): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+7 more)

### Community 40 - "context.Context"
Cohesion: 0.06
Nodes (14): context.Context, Check, MonitorSummary, CheckRepository, MaintenanceRepository, CheckRepo, Scheduler, benchEnv (+6 more)

### Community 41 - "MonitorService"
Cohesion: 0.07
Nodes (28): Registry, newTestRegistry(), TestRegistryHTTP(), TestRegistryUnknownType(), NewRegistry(), Checker, MonitorType, MonitorService (+20 more)

### Community 42 - "SendTestWebhook"
Cohesion: 0.18
Nodes (11): updateSettingsRequest, toUpdateSettingsInput(), SendTestWebhook(), TestSendTestWebhook(), webhookTestErrorMessage(), UpdateSettingsInput, validateWebhookURL(), normalizeTimezone() (+3 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 46 - "seed/main.go"
Cohesion: 0.56
Nodes (8): lastStep(), main(), monitorKind(), mustExec(), seedChecks(), seedIncidents(), seedMonitors(), profile

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 49 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.40
Nodes (5): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Limitações conhecidas, Validação

### Community 50 - "NewTCPChecker"
Cohesion: 0.24
Nodes (9): TCPChecker, net.Dialer, tcpMonitor(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused(), TestTCPCheckerTimeout(), TestTCPCheckerUp(), NewTCPChecker() (+1 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

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

### Community 57 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 58 - "newFakeSMTP"
Cohesion: 0.23
Nodes (13): net.Conn, net.Listener, TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), newFakeSMTP(), TestEmailNotifierSendsDown(), TestEmailNotifierSkipsWithoutRecipients() (+5 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.10
Nodes (19): react-router-dom, version, client, incidentsApi, web_src_assets_brand_pingou_marca_horizontal, web_src_assets_brand_pingou_marca_horizontal_negativa, BottomNav(), nav (+11 more)

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

### Community 72 - "database/sql.Tx"
Cohesion: 0.12
Nodes (10): database/sql.Tx, MonitorConfig, boolToInt(), encodeConfig(), MonitorRepo, nullableString(), nullableTime(), scanMonitor() (+2 more)

### Community 73 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 74 - "NtfyNotifier"
Cohesion: 0.23
Nodes (10): time.Location, TestCertMessageText(), certBody(), certPhrase(), subjectCert(), dashboardLink(), ntfyTitle(), LocationLoader() (+2 more)

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

### Community 82 - "SendTestNtfy"
Cohesion: 0.18
Nodes (12): net/http.Client, newWebhookTestClient(), newNtfyClient(), ntfyErrorMessage(), publishNtfy(), SendTestNtfy(), TestSendTestNtfyFailures(), TestSendTestNtfyTimeout() (+4 more)

### Community 84 - "TestDBBench"
Cohesion: 0.20
Nodes (12): time.Duration, rateLimiter, TestRateLimiterWindowResets(), newRateLimiter(), envInt(), measure(), ms(), printTable() (+4 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "Settings.jsx"
Cohesion: 0.16
Nodes (20): Entregue, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+12 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "database/sql.DB"
Cohesion: 0.33
Nodes (19): database/sql.DB, ExportHandler, Status(), CheckRepoTx, NewCheckRepoTx(), IncidentRepoTx, NewIncidentRepoTx(), NewMonitorRepo() (+11 more)

### Community 90 - "testing.T"
Cohesion: 0.14
Nodes (24): testing.T, TestNtfyHost(), TestMaintenanceRoutesNeedAPIKey(), env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker() (+16 more)

### Community 91 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 92 - "incidents_test.go"
Cohesion: 0.35
Nodes (11): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+3 more)

### Community 93 - "SendTestEmail"
Cohesion: 0.29
Nodes (7): buildEmail(), SendTestEmail(), smtpErrorMessage(), TestBuildEmail(), parseRecipients(), TestParseRecipients(), EmailTestResult

### Community 94 - "hotQueries"
Cohesion: 0.15
Nodes (19): testing.TB, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 95 - "runServe"
Cohesion: 0.27
Nodes (13): forceExitOnSecondSignal(), releaseLock(), runServe(), NewMaintenanceRepo(), NewMaintenanceService(), addWindowRaw(), newMaintenanceService(), newMaintenanceServiceDB() (+5 more)

### Community 96 - "NewMultiNotifier"
Cohesion: 0.47
Nodes (5): TestMultiNotifierCertExpiringFanOut(), NewMultiNotifier(), TestMultiNotifierDoesNotBlock(), TestMultiNotifierIgnoresCanceledContext(), TestMultiNotifierWait()

### Community 97 - "NewNtfyNotifier"
Cohesion: 0.24
Nodes (14): net/http.Header, sync.Mutex, TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), decodedTitle(), newFakeNtfy(), TestMultiNotifierNtfyFailureDoesNotBlockOthers(), TestNtfyNotifierDown() (+6 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 100 - "newCertServer"
Cohesion: 0.18
Nodes (11): HTTPChecker, net/http/httptest.Server, net/http.Response, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage() (+3 more)

### Community 101 - "commands/info.go"
Cohesion: 0.11
Nodes (25): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+17 more)

### Community 103 - "NewWebhookNotifier"
Cohesion: 0.50
Nodes (5): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType()

### Community 104 - "StateMachine"
Cohesion: 0.17
Nodes (9): MonitorState, MonitorFilter, EvaluateCertAlert(), sameInstant(), TestEvaluateCertAlert(), TestEvaluateCertAlertRepeated(), effectiveThreshold(), CertAlertDecision (+1 more)

### Community 105 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 106 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 107 - "bodyFor"
Cohesion: 0.23
Nodes (13): bodyFor(), formatDuration(), joinUnits(), oneLine(), shortDuration(), subjectDown(), subjectUp(), TestBodyFor() (+5 more)

### Community 109 - ".NotifyRecovery"
Cohesion: 0.47
Nodes (4): toWebhookMonitor(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

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

### Community 241 - "format.js"
Cohesion: 0.09
Nodes (45): lucide-react, OpenIncidentCard(), closeLabel, groupByDay(), ResolvedTable(), ChecksTable(), detailOf(), ConfigCard() (+37 more)

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
Nodes (33): HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu(), joinNames() (+25 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.25
Nodes (8): scripts, build, dev, format, format:check, lint, preview, test

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 252 - "serve_test.go"
Cohesion: 0.24
Nodes (9): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep (+1 more)

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

## Knowledge Gaps
- **505 isolated node(s):** `💡 Ideias de evolução`, `💡 Idea`, `🎯 Goal`, `📦 Initial Scope`, `🚧 Out of Scope` (+500 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 611 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `time.Time`, `net/http.ResponseWriter`, `Scheduler`, `UnitOfWork`, `CheckResult`, `newMaintenanceEnv`, `context.Context`, `MonitorService`, `NewTCPChecker`, `database/sql.Tx`, `NtfyNotifier`, `EmailNotifier`, `testing.T`, `NewMultiNotifier`, `newCertServer`, `StateMachine`, `bodyFor`, `.checkHeartbeats`, `.NotifyRecovery`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `5. Remodelagem mobile`, `3. Página de Configurações`, `11. Fuso horário + desligamento ordenado`, `Pingou README`, `12. Índices e desempenho do banco`, `MonitorDetail.jsx`, `13. Aviso de expiração do certificado SSL`, `7. Largura no desktop + busca parcial`, `9. Eixo Y dinâmico + gráfico de pings`, `14. Manutenção: pausa com motivo e janela agendada`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `10. Pausar monitor fecha o incidente`, `2. Página de Incidentes`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Why does `react` connect `Dashboard.jsx` to `MonitorForm.jsx`, `package.json`, `timezone.js`, `MonitorDetail.jsx`, `format.js`, `lib/monitors.js`, `Settings.jsx`, `main.jsx`?**
  _High betweenness centrality (0.013) - this node is a cross-community bridge._
- **What connects `💡 Ideias de evolução`, `💡 Idea`, `🎯 Goal` to the rest of the system?**
  _505 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13757855048177628 - nodes in this community are weakly interconnected._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0967741935483871 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.12631578947368421 - nodes in this community are weakly interconnected._