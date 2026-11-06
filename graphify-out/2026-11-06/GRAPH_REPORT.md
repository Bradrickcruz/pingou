# Graph Report - pingou-health-checker  (2026-11-05)

## Corpus Check
- 233 files · ~256,760 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 17 file(s) not represented in the graph (top: (none) 11, .woff2 4, .example 1)

## Summary
- 2179 nodes · 6555 edges · 136 communities (115 shown, 21 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 347 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `270802a7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- theme.js
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- time.Time
- net/http.ResponseWriter
- Testing Strategy & Framework
- Histórico de implementações (planos aplicados)
- Scheduler
- 🧩 Task: Backup sem configurações + importar
- newTestServer
- What You Must Do When Invoked
- UnitOfWork
- Pingou Health Checker - Project Vision & Goals
- newHeartbeat
- .Check
- Plano de Desenvolvimento: Retry de notificações com outbox (ideia 5)
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
- Plano de Desenvolvimento: Modo escuro (dark mode) no painel
- NewWebhookNotifier
- Monitor
- newMaintenanceEnv
- Incident
- MonitorService
- Incidents.jsx
- graphify reference: query, path, explain
- context.Context
- Pingou Health Checker - Project Roadmap
- database/sql.DB
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
- Plano de Desenvolvimento: Nomes claros nos eventos do webhook (ideia 7)
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
- newTestRegistry
- Dependencies & Blockers
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- extraction-spec.md
- .Check
- testing.T
- root.go
- incidents_test.go
- TestWaitRespectsDeadline
- openBench
- Status: **Planned**
- 5. Remodelagem mobile
- check-tokens.sh
- commands/info.go
- SMTPConfig
- 🧩 Task: Status HTTP esperado
- 🧩 Task: Páginas públicas de status
- 🧩 Task: Suporte a inglês (e outros idiomas)
- 🧩 Task: Retry de notificações
- .NotifyRecovery
- 🧩 Task: Aviso de manutenção começou/terminou
- 14. Manutenção: pausa com motivo e janela agendada
- 🧩 Task: Nomes claros nos eventos do webhook
- 7. Largura no desktop + busca parcial
- 8. Notificação via ntfy
- 6. E-mail (SMTP) e envio fora do lock
- IDEAS.md
- 9. Eixo Y dinâmico + gráfico de pings
- vite
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
4. `react` - 35 edges
5. `newTestServer()` - 34 edges
6. `do()` - 32 edges
7. `runServe()` - 31 edges
8. `newEnv()` - 29 edges
9. `lucide-react` - 29 edges
10. `Scheduler` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Task T3: Núcleo do tema (lib, hook, script inline)` --references--> `getPreference()`  [INFERRED]
  tasks/plano-dark-mode.md → web/src/lib/theme.js
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `⚠️ Notes` --references--> `ThemedImage()`  [INFERRED]
  IDEAS.md → web/src/components/ui/ThemedImage.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx

## Import Cycles
- None detected.

## Communities (136 total, 21 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.12
Nodes (68): MonitorEnabled, go_pkg_bytes, go_pkg_context, go_pkg_crypto_ecdsa, go_pkg_crypto_elliptic, go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_crypto_tls (+60 more)

### Community 1 - "theme.js"
Cohesion: 0.15
Nodes (20): 15. Modo escuro (Claro, Escuro, Sistema), Decisões, Descartado, Entregue, Limitações conhecidas, Validação, ref_node_fs, ThemeContext (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.13
Nodes (17): blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm(), PING_EVERY, PingBlock() (+9 more)

### Community 3 - "package.json"
Cohesion: 0.13
Nodes (18): autoprefixer, axios, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+10 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "time.Time"
Cohesion: 0.05
Nodes (32): DownWindow, database/sql.Result, time.Time, maintenanceListItem, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceWindowView (+24 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (48): net/http.Request, net/http.ResponseWriter, net/url.Values, time.Duration, checkResponse, incidentMonitorResponse, incidentResponse, maintenanceResponse (+40 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.40
Nodes (5): 4. Detalhe do monitor, Descartado, Desvios na implementação, Histórico de implementações (planos aplicados), Pendências abertas

### Community 9 - "Scheduler"
Cohesion: 0.20
Nodes (4): context.CancelFunc, sync.WaitGroup, Scheduler, job

### Community 10 - "🧩 Task: Backup sem configurações + importar"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Backup sem configurações + importar (+1 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (41): go_pkg_bufio, net/http.Handler, net/http/httptest.ResponseRecorder, settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode() (+33 more)

### Community 12 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 13 - "UnitOfWork"
Cohesion: 0.15
Nodes (3): UnitOfWork, failingCommit, MonitorRepositoryTx

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "newHeartbeat"
Cohesion: 0.26
Nodes (10): env, env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording() (+2 more)

### Community 16 - ".Check"
Cohesion: 0.40
Nodes (4): HTTPChecker, net/http.Response, certNotAfter(), requestErrorMessage()

### Community 17 - "Plano de Desenvolvimento: Retry de notificações com outbox (ideia 5)"
Cohesion: 0.06
Nodes (32): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Evidências, File Structure, Notes (+24 more)

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
Nodes (32): RFC-3339, Desvios do plano, Entregue, lucide-react, react-toastify, Task T4: Migrar base (`ui/*`, layout, estilos globais), maintenanceApi, CancelWindowDialog() (+24 more)

### Community 35 - "timezone.js"
Cohesion: 0.12
Nodes (24): ref_node_assert, ref_node_test, pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider() (+16 more)

### Community 36 - "Plano de Desenvolvimento: Modo escuro (dark mode) no painel"
Cohesion: 0.11
Nodes (17): Alternativas consideradas e descartadas, Build Verification, Cronograma Sugerido, Decisão de arquitetura (o ponto mais importante), Definition of Done (DoD), Fase 1 — tokens e migração (T2, T4–T7), Fase 2 — ligar o tema (T3, T8), File Structure (+9 more)

### Community 37 - "NewWebhookNotifier"
Cohesion: 0.50
Nodes (5): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType()

### Community 38 - "Monitor"
Cohesion: 0.06
Nodes (21): noopNotifier, CertAlert, Monitor, MonitorConfig, MonitorState, scanMonitor(), TestMultiNotifierIsolatesFailures(), Notifier (+13 more)

### Community 39 - "newMaintenanceEnv"
Cohesion: 0.21
Nodes (18): env, newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause() (+10 more)

### Community 40 - "Incident"
Cohesion: 0.11
Nodes (15): database/sql.NullInt64, database/sql.NullString, CloseReason, Incident, IncidentFilter, IncidentRepository, IncidentRepo, scanIncident() (+7 more)

### Community 41 - "MonitorService"
Cohesion: 0.07
Nodes (27): Registry, NewRegistry(), Checker, MonitorType, MonitorFilter, MonitorService, UpdateMonitorInput, newPingToken() (+19 more)

### Community 42 - "Incidents.jsx"
Cohesion: 0.09
Nodes (21): Concluído, Em andamento, Pendente, Status de Execução, Task T1: Fechar lacunas de design do escuro, client, incidentsApi, OPTIONS (+13 more)

### Community 43 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 44 - "context.Context"
Cohesion: 0.07
Nodes (13): context.Context, Check, MonitorSummary, CheckRepository, MonitorRepository, CheckRepo, Scheduler, benchEnv (+5 more)

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 46 - "database/sql.DB"
Cohesion: 0.20
Nodes (17): database/sql.DB, sync/atomic.Int64, ExportHandler, Status(), sizes(), env, Scheduler, fakeChecker (+9 more)

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 49 - "runServe"
Cohesion: 0.16
Nodes (40): runAdd(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), Open(), newTestServerEnv(), NewCheckRepo() (+32 more)

### Community 50 - "CheckResult"
Cohesion: 0.35
Nodes (14): CheckResult, fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal() (+6 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

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

### Community 57 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 58 - "newFakeSMTP"
Cohesion: 0.08
Nodes (31): net.Conn, net.Listener, updateSettingsRequest, toUpdateSettingsInput(), TestEmailNotifierSendsCertExpiring(), buildEmail(), NewEmailNotifier(), sendMail() (+23 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.08
Nodes (30): react-router-dom, Success Criteria, Task Breakdown, Task T0: Corrigir `DESIGN.md` duplicado, Task T10: Verificação visual, contraste e piscar, Task T11: Documentação final, Task T2: Camada de tokens de papel, Task T3: Núcleo do tema (lib, hook, script inline) (+22 more)

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
Cohesion: 0.09
Nodes (10): database/sql.Tx, IncidentRepoTx, boolToInt(), encodeConfig(), MonitorRepo, nullableString(), nullableTime(), MonitorRepoTx (+2 more)

### Community 73 - "Plano de Desenvolvimento: Nomes claros nos eventos do webhook (ideia 7)"
Cohesion: 0.07
Nodes (26): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Evidências, File Structure, Notes (+18 more)

### Community 74 - "bodyFor"
Cohesion: 0.11
Nodes (28): net/http.Client, time.Location, TestCertMessageText(), bodyFor(), certBody(), certPhrase(), formatDuration(), joinUnits() (+20 more)

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
Cohesion: 0.10
Nodes (27): Entregue, OpenIncidentCard(), CertCard(), levels, DetailTopBar(), DownAlert(), RangeTabs(), ChartLegend() (+19 more)

### Community 82 - "NewNtfyNotifier"
Cohesion: 0.14
Nodes (21): net/http.Header, net/http/httptest.Server, sync.Mutex, TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), ntfyErrorMessage(), SendTestNtfy(), decodedTitle() (+13 more)

### Community 84 - "newTestRegistry"
Cohesion: 0.31
Nodes (9): newTestRegistry(), tcpMonitor(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused(), TestTCPCheckerTimeout(), TestTCPCheckerUp() (+1 more)

### Community 85 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 86 - "react"
Cohesion: 0.21
Nodes (17): react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+9 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - ".Check"
Cohesion: 0.50
Nodes (3): TCPChecker, net.Dialer, tcpErrorMessage()

### Community 90 - "testing.T"
Cohesion: 0.11
Nodes (31): testing.T, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage(), NewHTTPChecker(), TestNtfyHost() (+23 more)

### Community 91 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 92 - "incidents_test.go"
Cohesion: 0.35
Nodes (11): go_pkg_slices, incidentsPage, Server, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse() (+3 more)

### Community 94 - "openBench"
Cohesion: 0.14
Nodes (20): testing.TB, copyFile(), openBench(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime() (+12 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 101 - "commands/info.go"
Cohesion: 0.18
Nodes (15): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+7 more)

### Community 103 - "SMTPConfig"
Cohesion: 0.09
Nodes (25): go_pkg_embed, go_pkg_io_fs, log/slog.Level, net/http.ServeMux, net/http.Server, Config, NtfyConfig, SMTPConfig (+17 more)

### Community 105 - "🧩 Task: Status HTTP esperado"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Status HTTP esperado (+1 more)

### Community 106 - "🧩 Task: Páginas públicas de status"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Páginas públicas de status (+1 more)

### Community 107 - "🧩 Task: Suporte a inglês (e outros idiomas)"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Suporte a inglês (e outros idiomas) (+1 more)

### Community 108 - "🧩 Task: Retry de notificações"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Retry de notificações (+1 more)

### Community 109 - ".NotifyRecovery"
Cohesion: 0.47
Nodes (4): toWebhookMonitor(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 110 - "🧩 Task: Aviso de manutenção começou/terminou"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Aviso de manutenção começou/terminou (+1 more)

### Community 111 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.40
Nodes (5): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Limitações conhecidas, Validação

### Community 112 - "🧩 Task: Nomes claros nos eventos do webhook"
Cohesion: 0.22
Nodes (9): 🔍 Context, 🎯 Goal, 💡 Idea, 📦 Initial Scope, ⚠️ Notes, 🚧 Out of Scope, ⏱️ Priority, 🧩 Task: Nomes claros nos eventos do webhook (+1 more)

### Community 113 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 114 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 115 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 116 - "IDEAS.md"
Cohesion: 0.40
Nodes (4): ⏸️ Adiadas, 🌱 Candidatas (ainda não decididas), 🗑️ Descartadas, 💡 Ideias de evolução

### Community 118 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 235 - "migrations_test.go"
Cohesion: 0.25
Nodes (15): runMigrateUp(), go_pkg_github_com_pressly_goose_v3, OpenDB(), RunMigrations(), sqliteDSN(), assertFKOn(), assertNoFKViolations(), countRows() (+7 more)

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
Nodes (38): closeLabel, CloseNote(), groupByDay(), ResolvedTable(), ChecksTable(), detailOf(), ConfigCard(), RecentIncidents() (+30 more)

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
Cohesion: 0.12
Nodes (32): monitorsApi, HistoryBar(), MonitorCard(), stateLabel, MonitorRow(), typeIcon, joinNames(), StatusBanner() (+24 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.22
Nodes (9): scripts, build, check:tokens, dev, format, format:check, lint, preview (+1 more)

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
- **628 isolated node(s):** `Overview`, `Project Type`, `Definition of Done (DoD)`, `Success Criteria`, `Tech Stack` (+623 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 739 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `theme.js`, `5. Remodelagem mobile`, `3. Página de Configurações`, `11. Fuso horário + desligamento ordenado`, `12. Índices e desempenho do banco`, `13. Aviso de expiração do certificado SSL`, `14. Manutenção: pausa com motivo e janela agendada`, `7. Largura no desktop + busca parcial`, `8. Notificação via ntfy`, `6. E-mail (SMTP) e envio fora do lock`, `DESIGN.md — Design system do Pingou`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `9. Eixo Y dinâmico + gráfico de pings`, `10. Pausar monitor fecha o incidente`, `2. Página de Incidentes`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `ThemedImage()` connect `main.jsx` to `theme.js`, `🧩 Task: Páginas públicas de status`, `Incidents.jsx`, `Dashboard.jsx`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `serve.go`, `time.Time`, `net/http.ResponseWriter`, `Scheduler`, `UnitOfWork`, `newHeartbeat`, `.Check`, `newMaintenanceEnv`, `Incident`, `MonitorService`, `context.Context`, `database/sql.DB`, `CheckResult`, `database/sql.Tx`, `bodyFor`, `newTestRegistry`, `.Check`, `testing.T`, `TestWaitRespectsDeadline`, `.NotifyRecovery`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `Overview`, `Project Type`, `Definition of Done (DoD)` to the rest of the system?**
  _628 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1229120473022912 - nodes in this community are weakly interconnected._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.13043478260869565 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.12631578947368421 - nodes in this community are weakly interconnected._