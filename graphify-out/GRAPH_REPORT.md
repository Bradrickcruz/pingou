# Graph Report - pingou-health-checker  (2026-11-04)

## Corpus Check
- 223 files · ~243,719 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2046 nodes · 6323 edges · 141 communities (117 shown, 24 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3e747ae3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- 6. E-mail (SMTP) e envio fora do lock
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- context.Context
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
- MaintenanceModal.jsx
- timezone.js
- SendTestWebhook
- runServe
- Incident
- newMaintenanceEnv
- database/sql.Tx
- MonitorService
- validateMonitor
- graphify reference: query, path, explain
- Check
- Pingou Health Checker - Project Roadmap
- seed/main.go
- Pingou - Divergencias entre PRD e implementacao
- graphify reference: add a URL and watch a folder
- database/sql.DB
- newCertEnv
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
- MonitorRepo
- 3. Cores
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- NtfyNotifier
- graphify reference: transcribe video and audio
- NewNtfyNotifier
- memtest.sh
- newTestRegistry
- Dependencies & Blockers
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- extraction-spec.md
- newHeartbeat
- testing.T
- root.go
- incidents_test.go
- 🧩 Task: Status HTTP esperado
- hotQueries
- newMaintenanceServiceDB
- 🧩 Task: Páginas públicas de status
- 🧩 Task: Modo escuro (dark mode)
- Status: **Planned**
- 5. Remodelagem mobile
- 🧩 Task: Suporte a inglês (e outros idiomas)
- commands/info.go
- Pingou README
- 3. Cores
- time.Time
- 9. Eixo Y dinâmico + gráfico de pings
- IDEAS.md
- 8. Notificação via ntfy
- 3.6 Tema escuro (dark mode)
- .NotifyRecovery
- 14. Manutenção: pausa com motivo e janela agendada
- Histórico de implementações (planos aplicados)
- memtest: prova da promessa de memória
- 7. Largura no desktop + busca parcial
- formatDuration
- 2. Marca
- 2. Marca
- 4. Status
- 4. Status
- migrations_test.go
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
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `CloseNote()`  [INFERRED]
  docs/implementacoes.md → web/src/components/incidents/ResolvedTable.jsx

## Import Cycles
- None detected.

## Communities (141 total, 24 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.13
Nodes (64): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_tls, go_pkg_crypto_x509, go_pkg_crypto_x509_pkix (+56 more)

### Community 1 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.50
Nodes (4): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Limitações conhecidas

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.15
Nodes (13): blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm(), PING_EVERY, TYPES (+5 more)

### Community 3 - "package.json"
Cohesion: 0.09
Nodes (22): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+14 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "context.Context"
Cohesion: 0.09
Nodes (12): context.Context, database/sql.Result, MaintenanceWindow, MaintenanceWindowView, MaintenanceRepository, affected(), fmtTime(), MaintenanceRepo (+4 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (52): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, net/http.Request, net/http.ResponseWriter, net/url.Values, time.Duration (+44 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.10
Nodes (27): Entregue, react-router-dom, OpenIncidentCard(), CertCard(), levels, DetailTopBar(), DownAlert(), RangeTabs() (+19 more)

### Community 9 - "Scheduler"
Cohesion: 0.13
Nodes (6): context.CancelFunc, sync.WaitGroup, Scheduler, Notifier, job, MultiNotifier

### Community 10 - "🧩 Task: Backup sem configurações + importar"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Backup sem configurações + importar (+1 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (41): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+33 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newEnv"
Cohesion: 0.22
Nodes (21): runAdd(), runRM(), Open(), RunMigrations(), NewCheckRepo(), openBench(), TestResultsSameWithOldIndex(), NewIncidentRepo() (+13 more)

### Community 16 - "Monitor"
Cohesion: 0.07
Nodes (22): HTTPChecker, Registry, net/http.Response, sync/atomic.Int64, noopNotifier, certNotAfter(), requestErrorMessage(), NewRegistry() (+14 more)

### Community 17 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.25
Nodes (8): certserver: testar o aviso de certificado SSL, Como interpretar, Como usar, Como usar (Linux), dbbench: desempenho do banco com volume real, O que é medido, Os dados, scripts/

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
Cohesion: 0.09
Nodes (33): RFC-3339, Desvios do plano, Entregue, axios, react-toastify, client, maintenanceApi, monitorsApi (+25 more)

### Community 35 - "timezone.js"
Cohesion: 0.12
Nodes (23): ref_node_assert, ref_node_test, pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider() (+15 more)

### Community 36 - "SendTestWebhook"
Cohesion: 0.15
Nodes (13): updateSettingsRequest, toUpdateSettingsInput(), SendTestWebhook(), TestSendTestWebhook(), webhookTestErrorMessage(), UpdateSettingsInput, validateWebhookURL(), normalizeTimezone() (+5 more)

### Community 37 - "runServe"
Cohesion: 0.10
Nodes (28): forceExitOnSecondSignal(), releaseLock(), runServe(), go_pkg_embed, go_pkg_io_fs, log/slog.Level, net/http.ServeMux, net/http.Server (+20 more)

### Community 38 - "Incident"
Cohesion: 0.06
Nodes (18): database/sql.NullInt64, CertAlert, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+10 more)

### Community 39 - "newMaintenanceEnv"
Cohesion: 0.25
Nodes (15): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+7 more)

### Community 40 - "database/sql.Tx"
Cohesion: 0.07
Nodes (10): database/sql.Tx, IncidentRepoTx, boolToInt(), encodeConfig(), nullableString(), nullableTime(), incidentRepoWithTx, MonitorRepositoryTx (+2 more)

### Community 41 - "MonitorService"
Cohesion: 0.18
Nodes (7): MonitorService, newPingToken(), newUUIDv7(), normalizeReason(), CreateMonitorInput, Recorder, Reloader

### Community 42 - "validateMonitor"
Cohesion: 0.13
Nodes (15): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+7 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "Check"
Cohesion: 0.11
Nodes (7): Check, CheckRepository, SettingsRepo, NewRetentionWorker(), RetentionWorker, CheckRepositoryTx, checkRepoWithTx

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 46 - "seed/main.go"
Cohesion: 0.49
Nodes (9): lastStep(), main(), monitorKind(), mustExec(), printSummary(), seedChecks(), seedIncidents(), seedMonitors() (+1 more)

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 49 - "database/sql.DB"
Cohesion: 0.26
Nodes (21): database/sql.DB, ExportHandler, Status(), CheckRepoTx, NewCheckRepoTx(), envInt(), sizes(), TestDBBench() (+13 more)

### Community 50 - "newCertEnv"
Cohesion: 0.37
Nodes (13): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+5 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.11
Nodes (19): 0. Como usar este arquivo, 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 10. Pendências (não definidas na v1.1), 11. Checklist antes de entregar, 11. Checklist antes de entregar, 1. Conceito, 1. Conceito (+11 more)

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
Cohesion: 0.16
Nodes (18): net.Conn, net.Listener, TestEmailNotifierSendsCertExpiring(), buildEmail(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage() (+10 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.09
Nodes (21): react-dom, incidentsApi, web_src_assets_brand_pingou_simbolo, BottomNav(), MobileHeader(), nav, Shell(), useOpenIncidents() (+13 more)

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

### Community 72 - "MonitorRepo"
Cohesion: 0.15
Nodes (9): database/sql.NullString, MonitorConfig, MonitorRepo, parseNullTime(), scanMonitor(), nameScanner, scanner, storedConfig (+1 more)

### Community 73 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 74 - "bodyFor"
Cohesion: 0.17
Nodes (18): time.Location, TestCertMessageText(), bodyFor(), certBody(), certPhrase(), oneLine(), splitRecipients(), subjectCert() (+10 more)

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

### Community 80 - "NtfyNotifier"
Cohesion: 0.27
Nodes (8): net/http.Client, dashboardLink(), newWebhookTestClient(), newNtfyClient(), ntfyTitle(), publishNtfy(), ntfyMessage, NtfyNotifier

### Community 82 - "NewNtfyNotifier"
Cohesion: 0.14
Nodes (21): net/http.Header, net/http/httptest.Server, sync.Mutex, TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), ntfyErrorMessage(), SendTestNtfy(), decodedTitle() (+13 more)

### Community 84 - "newTestRegistry"
Cohesion: 0.24
Nodes (11): TCPChecker, net.Dialer, newTestRegistry(), tcpMonitor(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused() (+3 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "react"
Cohesion: 0.15
Nodes (22): Entregue, react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard() (+14 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "newHeartbeat"
Cohesion: 0.26
Nodes (10): env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording() (+2 more)

### Community 90 - "testing.T"
Cohesion: 0.10
Nodes (30): testing.T, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage(), NewHTTPChecker(), TestNtfyHost() (+22 more)

### Community 91 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 92 - "incidents_test.go"
Cohesion: 0.35
Nodes (11): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+3 more)

### Community 93 - "🧩 Task: Status HTTP esperado"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Status HTTP esperado (+1 more)

### Community 94 - "hotQueries"
Cohesion: 0.15
Nodes (19): testing.TB, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 95 - "newMaintenanceServiceDB"
Cohesion: 0.46
Nodes (8): NewMaintenanceRepo(), addWindowRaw(), newMaintenanceServiceDB(), TestMaintenanceEditedIntervalIsUsed(), TestMaintenanceList(), TestMaintenanceUpdateRaceWithStart(), TestMaintenanceUpdateScheduled(), TestMaintenanceUpdateStartedLocked()

### Community 96 - "🧩 Task: Páginas públicas de status"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Páginas públicas de status (+1 more)

### Community 97 - "🧩 Task: Modo escuro (dark mode)"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Modo escuro (dark mode) (+1 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 100 - "🧩 Task: Suporte a inglês (e outros idiomas)"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Suporte a inglês (e outros idiomas) (+1 more)

### Community 101 - "commands/info.go"
Cohesion: 0.17
Nodes (16): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+8 more)

### Community 103 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 104 - "time.Time"
Cohesion: 0.07
Nodes (24): DownWindow, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MonitorSummary, toMaintenanceListItem(), CheckRepo (+16 more)

### Community 105 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 106 - "IDEAS.md"
Cohesion: 0.40
Nodes (4): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🗑️ Descartadas, 💡 Ideias de evolução

### Community 107 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 108 - "3.6 Tema escuro (dark mode)"
Cohesion: 0.33
Nodes (6): 3.6.1 Tokens, 3.6.2 Papéis: claro → escuro, 3.6.3 Status no tema escuro, 3.6.4 Contraste no tema escuro, 3.6.5 Regras do tema escuro, 3.6 Tema escuro (dark mode)

### Community 109 - ".NotifyRecovery"
Cohesion: 0.26
Nodes (9): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType(), toWebhookMonitor(), WebhookNotifier, webhookPayload (+1 more)

### Community 110 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.40
Nodes (5): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Limitações conhecidas, Validação

### Community 111 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.40
Nodes (5): 4. Detalhe do monitor, Descartado, Desvios na implementação, Histórico de implementações (planos aplicados), Pendências abertas

### Community 112 - "memtest: prova da promessa de memória"
Cohesion: 0.40
Nodes (5): Como interpretar, Como usar, Limitações, memtest: prova da promessa de memória, O que ele faz

### Community 113 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 114 - "formatDuration"
Cohesion: 0.50
Nodes (4): formatDuration(), joinUnits(), shortDuration(), TestFormatDuration()

### Community 235 - "migrations_test.go"
Cohesion: 0.30
Nodes (13): go_pkg_github_com_pressly_goose_v3, OpenDB(), sqliteDSN(), assertFKOn(), assertNoFKViolations(), countRows(), mustExec(), setupGoose() (+5 more)

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
Nodes (39): closeLabel, CloseNote(), groupByDay(), ResolvedTable(), ChecksTable(), detailOf(), ConfigCard(), RecentIncidents() (+31 more)

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
Nodes (36): lucide-react, HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu() (+28 more)

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
- **550 isolated node(s):** `0. Como usar este arquivo`, `0. Como usar este arquivo`, `1. Conceito`, `Arquivos oficiais`, `3.1 Tokens` (+545 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 656 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `6. E-mail (SMTP) e envio fora do lock`, `5. Remodelagem mobile`, `3. Página de Configurações`, `11. Fuso horário + desligamento ordenado`, `Pingou README`, `12. Índices e desempenho do banco`, `13. Aviso de expiração do certificado SSL`, `9. Eixo Y dinâmico + gráfico de pings`, `8. Notificação via ntfy`, `14. Manutenção: pausa com motivo e janela agendada`, `7. Largura no desktop + busca parcial`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `10. Pausar monitor fecha o incidente`, `2. Página de Incidentes`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `serve.go`, `net/http.ResponseWriter`, `Incident`, `MonitorRepo`, `time.Time`, `database/sql.Tx`, `bodyFor`, `validateMonitor`, `.NotifyRecovery`, `Scheduler`, `newMaintenanceEnv`, `NtfyNotifier`, `MonitorService`, `newCertEnv`, `newTestRegistry`, `newHeartbeat`, `testing.T`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `react` connect `react` to `MaintenanceModal.jsx`, `package.json`, `MonitorForm.jsx`, `timezone.js`, `MonitorDetail.jsx`, `format.js`, `Dashboard.jsx`, `main.jsx`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **What connects `0. Como usar este arquivo`, `0. Como usar este arquivo`, `1. Conceito` to the rest of the system?**
  _550 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12979559748427674 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.08994708994708994 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.08944793850454227 - nodes in this community are weakly interconnected._