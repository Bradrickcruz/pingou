# Graph Report - pingou  (2026-11-04)

## Corpus Check
- 216 files · ~140,183 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 10, .woff2 4, .example 1)

## Summary
- 2006 nodes · 6457 edges · 136 communities (114 shown, 22 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 430 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `583bbf37`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- MaintenanceWindow
- net/http.ResponseWriter
- Testing Strategy & Framework
- MonitorDetail.jsx
- Scheduler
- 🧩 Task: Aproveitar melhor a largura no desktop
- newTestServer
- Etapa 3: Página de manutenções
- time.Time
- Pingou Health Checker - Project Vision & Goals
- UnitOfWork
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
- time.Duration
- Plano de Desenvolvimento: {{Titulo}}
- MaintenanceConflictError
- Hero Image
- React Logo
- timezone.js
- middleware.go
- newCertEnv
- .NotifyRecovery
- Incident
- newMaintenanceEnv
- context.Context
- MonitorService
- newFakeSMTP
- Server
- vite
- Pingou Health Checker - Project Roadmap
- SendTestWebhook
- Pingou - Divergencias entre PRD e implementacao
- .listIncidents
- 14. Manutenção: pausa com motivo e janela agendada
- newTestRegistry
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- validateMonitor
- Task Breakdown
- Pingou Project State & Memory
- 10. Pausar monitor fecha o incidente
- Histórico de implementações (planos aplicados)
- TimezoneCard.jsx
- Status: **In Progress**
- Incidents.jsx
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
- net/http.Request
- scanMonitor
- 3. Cores
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- database/sql.DB
- MaintenanceService
- commands/info.go
- memtest.sh
- env
- newCertServer
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- newStatsService
- runServe
- testing.T
- root.go
- incidents_test.go
- NewWebhookNotifier
- hotQueries
- newMaintenanceServiceDB
- SettingsRepo
- NewNtfyNotifier
- Status: **Planned**
- 5. Remodelagem mobile
- fmtTime
- SendTestEmail
- Pingou README
- 2. Página de Incidentes
- 9. Eixo Y dinâmico + gráfico de pings
- spa.go
- 7. Largura no desktop + busca parcial
- .handleExportDB
- Resource Allocation
- migrations_test.go
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- Status: **Planned**
- devDependencies
- NewUnitOfWork
- format.js
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- Dashboard.jsx
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- dependencies
- SMTPConfig
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
- 8. Notificação via ntfy

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 112 edges
2. `Incident` - 61 edges
3. `MaintenanceWindow` - 37 edges
4. `newTestServer()` - 34 edges
5. `react` - 33 edges
6. `do()` - 32 edges
7. `MonitorDetail()` - 31 edges
8. `Task Breakdown` - 31 edges
9. `runServe()` - 31 edges
10. `lucide-react` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Task T26: Handlers `GET` e `PATCH` + 409 `MAINTENANCE_LOCKED`` --references--> `writeMaintenanceError()`  [INFERRED]
  plan-maintenance-windows.md → internal/handler/maintenance.go
- `Descartado` --references--> `NewServer()`  [INFERRED]
  docs/implementacoes.md → internal/handler/server.go
- `Task T2: Domínio + repositórios da pausa e do CloseReason` --references--> `incidentsListSQL()`  [INFERRED]
  plan-maintenance-windows.md → internal/repository/queries.go
- `Overview` --references--> `MaintenanceCard()`  [INFERRED]
  plan-maintenance-windows.md → web/src/components/monitor-detail/MaintenanceCard.jsx
- `Task T29: API client, rota, menu e `lib/maintenance.js`` --references--> `validateWindow()`  [INFERRED]
  plan-maintenance-windows.md → web/src/lib/maintenance.js

## Import Cycles
- 2-file cycle: `internal/repository/dbbench_test.go -> internal/service/maintenance_service_test.go -> internal/repository/dbbench_test.go`
- 2-file cycle: `internal/service/cert_state_machine_test.go -> internal/service/maintenance_service_test.go -> internal/service/cert_state_machine_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/maintenance_service_test.go -> internal/repository/dbbench_test.go`

## Communities (136 total, 22 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.12
Nodes (22): MonitorEnabled, createMaintenanceRequest, createMonitorRequest, downWindowResponse, infoResponse, monitorConfigRequest, paginatedResponse, statsBucketResponse (+14 more)

### Community 1 - "Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)"
Cohesion: 0.13
Nodes (14): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, File Structure, Notes, Overview (+6 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.10
Nodes (23): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm() (+15 more)

### Community 3 - "package.json"
Cohesion: 0.12
Nodes (19): autoprefixer, axios, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+11 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "MaintenanceWindow"
Cohesion: 0.14
Nodes (9): maintenanceListItem, MaintenanceWindow, MaintenanceRepository, toMaintenanceListItem(), scanMaintenance(), Task T12: Gate de manutenção no scheduler e no ping, Task T14: Handlers + wiring + bloco `maintenance`, Task T9: Domínio + `MaintenanceRepo` (+1 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.30
Nodes (7): Server, Server, writeServiceError(), writeError(), writeJSON(), Server, Server

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.17
Nodes (21): Entregue, ConfigCard(), DetailTopBar(), DownAlert(), RangeTabs(), RecentIncidents(), ChartLegend(), MonitorModals() (+13 more)

### Community 9 - "Scheduler"
Cohesion: 0.09
Nodes (11): Registry, certNotAfter(), requestErrorMessage(), NewRegistry(), tcpErrorMessage(), Checker, CheckResult, MonitorType (+3 more)

### Community 10 - "🧩 Task: Aproveitar melhor a largura no desktop"
Cohesion: 0.05
Nodes (37): 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🎯 Goal, 🎯 Goal, 🎯 Goal, 🎯 Goal (+29 more)

### Community 11 - "newTestServer"
Cohesion: 0.12
Nodes (36): settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode(), do(), newTestServer(), newTestServerSMTP() (+28 more)

### Community 12 - "Etapa 3: Página de manutenções"
Cohesion: 0.12
Nodes (14): MaintenanceWindowView, affected(), Etapa 3: Página de manutenções, Task T22: Criar branch, Task T23: Contrato `GET` e `PATCH /maintenance-windows`, Task T24: Repositório: `ListLive` e `Update`, Task T25: `MaintenanceService`: `List` e `Update`, Task T26: Handlers `GET` e `PATCH` + 409 `MAINTENANCE_LOCKED` (+6 more)

### Community 13 - "time.Time"
Cohesion: 0.12
Nodes (16): CheckRepository, DownWindow, MonitorStats, StatsBucket, MonitorSummary, CheckRepo, startWindowTx(), Scheduler (+8 more)

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 16 - "Monitor"
Cohesion: 0.07
Nodes (15): noopNotifier, CertAlert, Monitor, MonitorState, Notifier, effectiveThreshold(), noopNotifier, slowNotifier (+7 more)

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

### Community 29 - "time.Duration"
Cohesion: 0.19
Nodes (10): rateLimiter, TestRateLimiterWindowResets(), Server, clientIP(), newRateLimiter(), rejectRateLimited(), measure(), ms() (+2 more)

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 34 - "timezone.js"
Cohesion: 0.10
Nodes (42): RFC-3339, Desvios do plano, Entregue, Task T16: `zonedToUtc` + status MAINTENANCE, Task T17: Modal de nova janela, Task T28: Fluxo e estados da página, Task T31: `MaintenanceModal` em modo editar, Task T32: "Editar" no `MaintenanceCard` do detalhe (+34 more)

### Community 35 - "middleware.go"
Cohesion: 0.19
Nodes (7): responseWriter, constantTimeEqual(), Server, loggingMiddleware(), newRequestID(), recoverMiddleware(), redactPath()

### Community 36 - "newCertEnv"
Cohesion: 0.37
Nodes (13): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+5 more)

### Community 37 - ".NotifyRecovery"
Cohesion: 0.47
Nodes (4): toWebhookMonitor(), WebhookNotifier, webhookPayload, webhookPayloadMonitor

### Community 38 - "Incident"
Cohesion: 0.12
Nodes (13): IncidentRepository, CloseReason, Incident, IncidentFilter, IncidentRepo, scanIncident(), scanIncidentWithMonitor(), parseNullTime() (+5 more)

### Community 39 - "newMaintenanceEnv"
Cohesion: 0.19
Nodes (19): newMaintenanceEnv(), TestMaintenanceAppliedOnStart(), TestMaintenanceBlocksChecksAndClosesIncident(), TestMaintenanceCancel(), TestMaintenanceCancelAtStartSecond(), TestMaintenanceDiscardsCheckInFlight(), TestMaintenanceEndKeepsManualPause(), TestMaintenanceEndResumesMonitor() (+11 more)

### Community 40 - "context.Context"
Cohesion: 0.05
Nodes (17): Check, MonitorRepository, IncidentRepoTx, boolToInt(), encodeConfig(), MonitorRepo, nullableString(), nullableTime() (+9 more)

### Community 41 - "MonitorService"
Cohesion: 0.18
Nodes (5): MonitorService, newPingToken(), CreateMonitorInput, Recorder, Reloader

### Community 42 - "newFakeSMTP"
Cohesion: 0.24
Nodes (10): TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), newFakeSMTP(), TestEmailNotifierSendsDown(), TestEmailNotifierSkipsWithoutRecipients(), TestSendMailNone(), TestSendMailTimeout() (+2 more)

### Community 43 - "Server"
Cohesion: 0.25
Nodes (6): maintenanceResponse, Server, toMaintenanceResponse(), writeMaintenanceError(), NewServer(), IncidentService

### Community 45 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Critical Dependencies, Dependencies & Blockers, Development Phases Overview, Pingou Health Checker - Project Roadmap, Potential Blockers

### Community 46 - "SendTestWebhook"
Cohesion: 0.16
Nodes (12): updateSettingsRequest, toUpdateSettingsInput(), newWebhookTestClient(), SendTestWebhook(), TestSendTestWebhook(), webhookTestErrorMessage(), UpdateSettingsInput, validateWebhookURL() (+4 more)

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - ".listIncidents"
Cohesion: 0.16
Nodes (9): checkResponse, incidentMonitorResponse, incidentResponse, Server, toCheckResponse(), Server, incidentStatus(), toIncidentResponse() (+1 more)

### Community 49 - "14. Manutenção: pausa com motivo e janela agendada"
Cohesion: 0.40
Nodes (5): 14. Manutenção: pausa com motivo e janela agendada, Decisões, Descartado, Limitações conhecidas, Validação

### Community 50 - "newTestRegistry"
Cohesion: 0.24
Nodes (10): TCPChecker, newTestRegistry(), tcpMonitor(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost(), TestTCPCheckerRefused(), TestTCPCheckerTimeout() (+2 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

### Community 53 - "validateMonitor"
Cohesion: 0.13
Nodes (15): isHTTPS(), TestValidateMonitorByType(), TestValidateURL(), validateCertAlert(), validateGrace(), validateHeartbeatInterval(), validateHost(), validateInterval() (+7 more)

### Community 54 - "Task Breakdown"
Cohesion: 0.14
Nodes (14): Task Breakdown, Task T0: Criar branch, Task T10: Contrato da API de manutenção, Task T15: Testes da etapa 2 (backend), Task T19: Revisão de bordas e segurança, Task T1: Migração 0008 (pausa + close_note), Task T20: Docs, Task T21: Verificação final e PR (+6 more)

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 56 - "10. Pausar monitor fecha o incidente"
Cohesion: 0.40
Nodes (5): 10. Pausar monitor fecha o incidente, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 57 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.20
Nodes (10): 4. Detalhe do monitor, 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Descartado, Desvios na implementação, Entregue, Histórico de implementações (planos aplicados) (+2 more)

### Community 58 - "TimezoneCard.jsx"
Cohesion: 0.24
Nodes (10): Task T30: Página `Maintenance.jsx`, Tech Stack, pad(), TimezoneCard(), Combobox(), ConnectionOverlay(), useConnection(), useNow() (+2 more)

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "Incidents.jsx"
Cohesion: 0.09
Nodes (19): Success Criteria, client, incidentsApi, monitorsApi, BottomNav(), nav, Shell(), useOpenIncidents() (+11 more)

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

### Community 71 - "net/http.Request"
Cohesion: 0.32
Nodes (5): monitorListItem, monitorResponse, formatTime(), Server, toMonitorResponse()

### Community 72 - "scanMonitor"
Cohesion: 0.22
Nodes (6): MonitorFilter, MonitorConfig, scanMonitor(), nameScanner, scanner, storedConfig

### Community 73 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 74 - "bodyFor"
Cohesion: 0.13
Nodes (23): TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), oneLine(), splitRecipients(), subjectCert() (+15 more)

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
Cohesion: 0.27
Nodes (14): runMigrateDown(), ExportHandler, Down(), Status(), sizes(), lastStep(), main(), monitorKind() (+6 more)

### Community 81 - "MaintenanceService"
Cohesion: 0.24
Nodes (7): dedupe(), MaintenanceService, validateWindowTimes(), normalizeReason(), Task T11: `MaintenanceService`, CreateMaintenanceInput, MaintenanceEnder

### Community 82 - "commands/info.go"
Cohesion: 0.39
Nodes (8): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim()

### Community 84 - "env"
Cohesion: 0.25
Nodes (3): env, dbNotifier, fakeChecker

### Community 85 - "newCertServer"
Cohesion: 0.22
Nodes (8): HTTPChecker, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerNoCertWithoutTLS(), TestHTTPCheckerReadsCertExpiry(), TestHTTPCheckerErrorMessage(), NewHTTPChecker(), newNtfyClient()

### Community 86 - "react"
Cohesion: 0.13
Nodes (26): react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+18 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 88 - "newStatsService"
Cohesion: 0.31
Nodes (9): TestPauseClosesOpenIncidentSilently(), TestPauseReasonSavedAndClearedOnResume(), TestPauseWithoutIncidentAndResumeKeepState(), updateEnabled(), UpdateMonitorInput, at(), newStatsService(), TestStatsBucketsAndTotals() (+1 more)

### Community 89 - "runServe"
Cohesion: 0.22
Nodes (21): runAdd(), copyFile(), runExportDB(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), newTestServerEnv() (+13 more)

### Community 90 - "testing.T"
Cohesion: 0.11
Nodes (25): TestNtfyHost(), TestMaintenanceRoutesNeedAPIKey(), env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown() (+17 more)

### Community 91 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 92 - "incidents_test.go"
Cohesion: 0.35
Nodes (9): incidentsPage, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse(), TestListIncidentsResolvedOrderAndFilters(), TestListIncidentsTotalIgnoresLimit() (+1 more)

### Community 93 - "NewWebhookNotifier"
Cohesion: 0.50
Nodes (5): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType()

### Community 94 - "hotQueries"
Cohesion: 0.15
Nodes (19): Entregue, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 95 - "newMaintenanceServiceDB"
Cohesion: 0.46
Nodes (8): NewMaintenanceRepo(), addWindowRaw(), newMaintenanceServiceDB(), TestMaintenanceEditedIntervalIsUsed(), TestMaintenanceList(), TestMaintenanceUpdateRaceWithStart(), TestMaintenanceUpdateScheduled(), TestMaintenanceUpdateStartedLocked()

### Community 96 - "SettingsRepo"
Cohesion: 0.31
Nodes (3): SettingsRepo, NewRetentionWorker(), RetentionWorker

### Community 97 - "NewNtfyNotifier"
Cohesion: 0.15
Nodes (19): TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), ntfyErrorMessage(), SendTestNtfy(), decodedTitle(), newFakeNtfy(), TestMultiNotifierNtfyFailureDoesNotBlockOthers(), TestNtfyNotifierDown() (+11 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 100 - "fmtTime"
Cohesion: 0.43
Nodes (3): fmtTime(), Risk Areas, Task T13: Tick de início/fim da janela

### Community 101 - "SendTestEmail"
Cohesion: 0.33
Nodes (6): SendTestEmail(), smtpErrorMessage(), TestSendTestEmail(), parseRecipients(), TestParseRecipients(), EmailTestResult

### Community 103 - "2. Página de Incidentes"
Cohesion: 0.50
Nodes (4): 2. Página de Incidentes, Decisões, Descartado, Entregue

### Community 104 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 106 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 108 - "Resource Allocation"
Cohesion: 0.67
Nodes (3): Development Focus, Infrastructure Costs, Resource Allocation

### Community 235 - "migrations_test.go"
Cohesion: 0.21
Nodes (17): runMigrateStatus(), runMigrateUp(), ListMigrations(), Open(), OpenDB(), RunMigrations(), sqliteDSN(), assertFKOn() (+9 more)

### Community 236 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 240 - "NewUnitOfWork"
Cohesion: 0.42
Nodes (14): CheckRepoTx, NewCheckRepoTx(), envInt(), TestDBBench(), NewIncidentRepoTx(), NewMonitorRepoTx(), TestPauseIsSilentAndResumeReopensIncident(), TestResumeIgnoresFailuresBeforePause() (+6 more)

### Community 241 - "format.js"
Cohesion: 0.08
Nodes (44): Task T6: Modal de pausa + selos com motivo, OpenIncidentCard(), closeLabel, CloseNote(), groupByDay(), ResolvedTable(), Actions(), stateBadge (+36 more)

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
Cohesion: 0.10
Nodes (39): Task T18: Selos, card de detalhes, cancelar, lucide-react, react-router-dom, MobileHeader(), HistoryBar(), MonitorCard(), stateLabel, PingBlock() (+31 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.25
Nodes (8): scripts, build, dev, format, format:check, lint, preview, test

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 250 - "SMTPConfig"
Cohesion: 0.17
Nodes (13): Config, NtfyConfig, SMTPConfig, Load(), LoadConfig(), loadConfigWithAPIKey(), LoadNtfy(), LoadSMTP() (+5 more)

### Community 252 - "serve_test.go"
Cohesion: 0.24
Nodes (8): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep

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
Cohesion: 0.40
Nodes (5): 12. Índices e desempenho do banco, Decisões, Descartado, Limitações conhecidas, Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

### Community 265 - "13. Aviso de expiração do certificado SSL"
Cohesion: 0.25
Nodes (7): 13. Aviso de expiração do certificado SSL, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação

### Community 267 - "siteCert"
Cohesion: 0.47
Nodes (3): loadOrCreateCA(), main(), siteCert()

### Community 272 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

## Knowledge Gaps
- **24 isolated node(s):** `RFC-3339`, `MonitorService`, `github.com/Bradrickcruz/pingou`, `TMPDIR`, `autoprefixer` (+19 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 593 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Task Breakdown` connect `Task Breakdown` to `Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)`, `timezone.js`, `fmtTime`, `MaintenanceWindow`, `Etapa 3: Página de manutenções`, `MaintenanceService`, `format.js`, `Dashboard.jsx`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **What connects `RFC-3339`, `MonitorService`, `github.com/Bradrickcruz/pingou` to the rest of the system?**
  _24 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12239467849223946 - nodes in this community are weakly interconnected._
- **Why does `MaintenanceService` connect `MaintenanceService` to `serve.go`, `timezone.js`, `fmtTime`, `MaintenanceWindow`, `newMaintenanceEnv`, `context.Context`, `13. Aviso de expiração do certificado SSL`, `Server`, `Etapa 3: Página de manutenções`, `time.Time`, `Task Breakdown`, `runServe`, `newMaintenanceServiceDB`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Should `Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._
- **Why does `MaintenanceModal()` connect `timezone.js` to `MonitorDetail.jsx`, `TimezoneCard.jsx`, `Dashboard.jsx`, `Task Breakdown`?**
  _High betweenness centrality (0.105) - this node is a cross-community bridge._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0967741935483871 - nodes in this community are weakly interconnected._