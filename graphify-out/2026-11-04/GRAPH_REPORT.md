# Graph Report - pingou  (2026-11-04)

## Corpus Check
- 200 files · ~120,631 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 10, .woff2 4, .example 1)

## Summary
- 1827 nodes · 5511 edges · 118 communities (99 shown, 19 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 367 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `550c4bd9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- Task Breakdown
- MonitorForm.jsx
- package.json
- 3. Página de Configurações
- incidents_test.go
- net/http.ResponseWriter
- Testing Strategy & Framework
- MonitorDetail.jsx
- MonitorType
- 🧩 Task: Aproveitar melhor a largura no desktop
- settings_test.go
- fakeSMTP
- EvaluateCertAlert
- Pingou Health Checker - Project Vision & Goals
- NewWebhookNotifier
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
- 7. Largura no desktop + busca parcial
- Plano de Desenvolvimento: {{Titulo}}
- Server
- Hero Image
- React Logo
- SettingsRepo
- newFakeSMTP
- newCertEnv
- 10. Pausar monitor fecha o incidente
- Incident
- newEnv
- sqliteUnitOfWork
- cert.js
- format.js
- at
- vite
- Resource Allocation
- Pingou - Divergencias entre PRD e implementacao
- Pingou Health Checker - Project Roadmap
- Notifier
- newTestRegistry
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- context.Context
- 6. E-mail (SMTP) e envio fora do lock
- Pingou Project State & Memory
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
- 3. Cores
- bodyFor
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- memtest.sh
- react
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- openBench
- 11. Fuso horário + desligamento ordenado
- time.Time
- NewNtfyNotifier
- Status: **Planned**
- Histórico de implementações (planos aplicados)
- migrations_test.go
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- Status: **Planned**
- devDependencies
- NtfyConfig
- database/sql.DB
- Dashboard.jsx
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- MonitorDetail.jsx
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- 9. Eixo Y dinâmico + gráfico de pings
- dependencies
- SendTestNtfy
- timezone.js
- commands/info.go
- Success Metrics
- Release Schedule
- Risk Management
- seed/main.go
- root.go
- serve_test.go
- 12. Índices e desempenho do banco
- SendTestWebhook
- 13. Aviso de expiração do certificado SSL
- Check
- siteCert
- Pingou README
- dbbench.sh
- dbmem.sh
- 8. Notificação via ntfy
- NewMultiNotifier

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 110 edges
2. `Incident` - 61 edges
3. `newTestServer()` - 31 edges
4. `do()` - 30 edges
5. `runServe()` - 29 edges
6. `react` - 29 edges
7. `newEnv()` - 28 edges
8. `MonitorDetail()` - 27 edges
9. `Scheduler` - 27 edges
10. `lucide-react` - 25 edges

## Surprising Connections (you probably didn't know these)
- `Task T6: Modal de pausa + selos com motivo` --references--> `ResolvedTable()`  [INFERRED]
  plan-maintenance-windows.md → web/src/components/incidents/ResolvedTable.jsx
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Entregue` --references--> `TimeZoneProvider()`  [INFERRED]
  docs/implementacoes.md → web/src/hooks/useTimeZone.jsx
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx

## Import Cycles
- 2-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/repository/dbbench_test.go`
- 2-file cycle: `internal/repository/dbbench_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 2-file cycle: `internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`

## Communities (118 total, 19 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.10
Nodes (36): MonitorEnabled, createMonitorRequest, downWindowResponse, infoResponse, monitorConfigRequest, paginatedResponse, statsBucketResponse, statsResponse (+28 more)

### Community 1 - "Task Breakdown"
Cohesion: 0.05
Nodes (36): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, File Structure, Notes, Overview (+28 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.15
Nodes (12): blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, MonitorForm(), PING_EVERY, TYPES (+4 more)

### Community 3 - "package.json"
Cohesion: 0.14
Nodes (16): autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals (+8 more)

### Community 4 - "3. Página de Configurações"
Cohesion: 0.50
Nodes (4): 3. Página de Configurações, Decisões, Descartado, Entregue

### Community 5 - "incidents_test.go"
Cohesion: 0.38
Nodes (8): incidentsPage, ids(), listIncidents(), seedIncidents(), TestClosedManuallyInResolved(), TestCloseReasonAndNoteInResponse(), TestListIncidentsResolvedOrderAndFilters(), TestListIncidentsTotalIgnoresLimit()

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (39): checkResponse, incidentMonitorResponse, incidentResponse, monitorListItem, monitorResponse, rateLimiter, responseWriter, Server (+31 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.14
Nodes (24): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, ConfigCard(), DetailTopBar(), DownAlert(), RangeTabs() (+16 more)

### Community 9 - "MonitorType"
Cohesion: 0.15
Nodes (13): Registry, NewRegistry(), Checker, MonitorConfig, MonitorType, MonitorFilter, scanIncident(), scanIncidentWithMonitor() (+5 more)

### Community 10 - "🧩 Task: Aproveitar melhor a largura no desktop"
Cohesion: 0.05
Nodes (37): 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🎯 Goal, 🎯 Goal, 🎯 Goal, 🎯 Goal (+29 more)

### Community 11 - "settings_test.go"
Cohesion: 0.13
Nodes (40): settingsBody, newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerReadsCertExpiry(), TestNtfyHost(), TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat() (+32 more)

### Community 13 - "EvaluateCertAlert"
Cohesion: 0.40
Nodes (5): EvaluateCertAlert(), sameInstant(), TestEvaluateCertAlert(), TestEvaluateCertAlertRepeated(), CertAlertDecision

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "NewWebhookNotifier"
Cohesion: 0.50
Nodes (5): TestWebhookCertExpiring(), NewWebhookNotifier(), eqPtr(), TestWebhookDashboardURL(), TestWebhookPayloadByType()

### Community 16 - "Monitor"
Cohesion: 0.07
Nodes (16): noopNotifier, CertAlert, Monitor, toWebhookMonitor(), noopNotifier, dbNotifier, fakeChecker, slowNotifier (+8 more)

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

### Community 29 - "7. Largura no desktop + busca parcial"
Cohesion: 0.50
Nodes (4): 7. Largura no desktop + busca parcial, Decisões, Descartado, Entregue

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 31 - "Server"
Cohesion: 0.22
Nodes (4): Server, NewServer(), spaHandler(), IncidentService

### Community 34 - "SettingsRepo"
Cohesion: 0.29
Nodes (4): CheckRepository, SettingsRepo, NewRetentionWorker(), RetentionWorker

### Community 35 - "newFakeSMTP"
Cohesion: 0.27
Nodes (13): SMTPConfig, TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage(), newFakeSMTP(), TestEmailNotifierSendsDown() (+5 more)

### Community 36 - "newCertEnv"
Cohesion: 0.37
Nodes (13): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+5 more)

### Community 37 - "10. Pausar monitor fecha o incidente"
Cohesion: 0.40
Nodes (5): 10. Pausar monitor fecha o incidente, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 38 - "Incident"
Cohesion: 0.11
Nodes (8): CloseReason, Incident, IncidentFilter, IncidentRepository, incidentsListSQL(), IncidentRepo, dbNotifier, incidentRepoWithTx

### Community 39 - "newEnv"
Cohesion: 0.14
Nodes (21): env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline(), TestHeartbeatEvaluatorRunsOnTicker(), TestHeartbeatLifecycle(), TestHeartbeatNeverPingedGoesDown(), TestHeartbeatPausedAcceptsPingWithoutRecording(), TestHeartbeatRegenerateToken() (+13 more)

### Community 40 - "sqliteUnitOfWork"
Cohesion: 0.15
Nodes (8): MonitorState, effectiveThreshold(), CheckRepositoryTx, IncidentRepositoryTx, MonitorRepositoryTx, sqliteUnitOfWork, StateMachine, UnitOfWork

### Community 41 - "cert.js"
Cohesion: 0.16
Nodes (12): CertCard(), levels, CERT_ALERT_DAYS, CERT_WARN_DAYS, certDate(), certStatus(), hasCert(), isHttps() (+4 more)

### Community 42 - "format.js"
Cohesion: 0.11
Nodes (39): Task T16: `zonedToUtc` + status MAINTENANCE, OpenIncidentCard(), closeLabel, CloseNote(), groupByDay(), ManualBadge(), MonitorCell(), ResolvedTable() (+31 more)

### Community 43 - "at"
Cohesion: 0.67
Nodes (3): at(), TestStatsBucketsAndTotals(), TestStatsIncidentsClipped()

### Community 45 - "Resource Allocation"
Cohesion: 0.67
Nodes (3): Development Focus, Infrastructure Costs, Resource Allocation

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Critical Dependencies, Dependencies & Blockers, Development Phases Overview, Pingou Health Checker - Project Roadmap, Potential Blockers

### Community 50 - "newTestRegistry"
Cohesion: 0.24
Nodes (13): TCPChecker, TestHTTPCheckerNoCertWithoutTLS(), newTestRegistry(), tcpMonitor(), TestHTTPCheckerErrorMessage(), TestRegistryHTTP(), TestRegistryUnknownType(), TestTCPCheckerNoSuchHost() (+5 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

### Community 53 - "context.Context"
Cohesion: 0.07
Nodes (14): requestErrorMessage(), tcpErrorMessage(), CheckResult, MonitorRepository, Scheduler, newPingToken(), normalizeReason(), MonitorRepo (+6 more)

### Community 54 - "6. E-mail (SMTP) e envio fora do lock"
Cohesion: 0.40
Nodes (5): 6. E-mail (SMTP) e envio fora do lock, Decisões, Descartado, Entregue, Limitações conhecidas

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.09
Nodes (22): axios, react-dom, react-router-dom, client, incidentsApi, BottomNav(), nav, Shell() (+14 more)

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

### Community 73 - "3. Cores"
Cohesion: 0.33
Nodes (6): 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção, 3.4 Regras de cor, 3.5 Contraste sobre branco, 3. Cores

### Community 74 - "bodyFor"
Cohesion: 0.11
Nodes (25): HTTPChecker, TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), oneLine(), splitRecipients() (+17 more)

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

### Community 86 - "react"
Cohesion: 0.16
Nodes (21): react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+13 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "openBench"
Cohesion: 0.19
Nodes (13): copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath(), explain(), hotQueries() (+5 more)

### Community 92 - "11. Fuso horário + desligamento ordenado"
Cohesion: 0.33
Nodes (6): 11. Fuso horário + desligamento ordenado, Decisões, Descartado, Entregue, Limitações conhecidas, Validação

### Community 96 - "time.Time"
Cohesion: 0.13
Nodes (11): DownWindow, certNotAfter(), MonitorStats, StatsBucket, MonitorSummary, CheckRepo, Scheduler, MonitorService (+3 more)

### Community 97 - "NewNtfyNotifier"
Cohesion: 0.15
Nodes (18): TestNtfyNotifierCertExpiring(), newNtfyClient(), NewNtfyNotifier(), ntfyErrorMessage(), SendTestNtfy(), decodedTitle(), newFakeNtfy(), TestMultiNotifierNtfyFailureDoesNotBlockOthers() (+10 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.20
Nodes (10): 2. Página de Incidentes, 5. Remodelagem mobile, Decisões, Decisões, Descartado, Descartado, Entregue, Entregue (+2 more)

### Community 235 - "migrations_test.go"
Cohesion: 0.16
Nodes (19): copyFile(), runExportDB(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), Down(), ListMigrations(), OpenDB() (+11 more)

### Community 236 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 239 - "NtfyConfig"
Cohesion: 0.24
Nodes (8): Config, NtfyConfig, Load(), LoadConfig(), loadConfigWithAPIKey(), LoadNtfy(), TestLoadFailsOnInvalidSMTP(), TestLoadNtfy()

### Community 240 - "database/sql.DB"
Cohesion: 0.13
Nodes (49): runAdd(), runRM(), forceExitOnSecondSignal(), releaseLock(), runServe(), ExportHandler, Open(), RunMigrations() (+41 more)

### Community 241 - "Dashboard.jsx"
Cohesion: 0.12
Nodes (19): lucide-react, react-toastify, version, monitorsApi, MobileHeader(), MonitorModals(), PauseModal(), REASON_MAX (+11 more)

### Community 242 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 4.1: Multi-tenancy & Access Control, Milestone 4.2: Advanced Authentication, Milestone 4.3: Scalability Improvements, Phase 4: Enterprise Features (Q4 2024), Status: **Planned**

### Community 243 - "Status: **Future**"
Cohesion: 0.40
Nodes (5): Milestone 5.1: Advanced Monitoring, Milestone 5.2: Integration Ecosystem, Milestone 5.3: Public Features, Phase 5: Advanced Features & Integrations (Q1-Q2 2025), Status: **Future**

### Community 244 - "Status: **Continuous**"
Cohesion: 0.40
Nodes (5): Milestone 6.1: Community Building, Milestone 6.2: Developer Experience, Milestone 6.3: Ecosystem Integration, Phase 6: Ecosystem & Community (Ongoing), Status: **Continuous**

### Community 245 - "MonitorDetail.jsx"
Cohesion: 0.13
Nodes (30): Task T18: Selos, card de detalhes, cancelar, HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu() (+22 more)

### Community 246 - "1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel"
Cohesion: 0.33
Nodes (6): 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel, Decisões, Descartado, Entregue, Limitações conhecidas, http()

### Community 247 - "scripts"
Cohesion: 0.25
Nodes (8): scripts, build, dev, format, format:check, lint, preview, test

### Community 248 - "9. Eixo Y dinâmico + gráfico de pings"
Cohesion: 0.50
Nodes (4): 9. Eixo Y dinâmico + gráfico de pings, Decisões, Descartado, Entregue

### Community 249 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, axios, lucide-react, react, react-dom, react-router-dom, react-toastify

### Community 250 - "SendTestNtfy"
Cohesion: 0.17
Nodes (12): updateSettingsRequest, toUpdateSettingsInput(), TestValidNtfyTopic(), validNtfyTopic(), SettingsService, UpdateSettingsInput, normalizeTimezone(), TestNormalizeTimezone() (+4 more)

### Community 251 - "timezone.js"
Cohesion: 0.15
Nodes (22): Task T17: Modal de nova janela, Tech Stack, pad(), TimezoneCard(), Combobox(), TimeZoneContext, loadTimeZone(), TimeZoneProvider() (+14 more)

### Community 252 - "commands/info.go"
Cohesion: 0.29
Nodes (10): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim() (+2 more)

### Community 253 - "Success Metrics"
Cohesion: 0.50
Nodes (4): Community Metrics, Success Metrics, Technical Metrics, User Metrics

### Community 254 - "Release Schedule"
Cohesion: 0.50
Nodes (4): Current Version Tracking, Release Frequency, Release Schedule, Versioning Strategy

### Community 257 - "Risk Management"
Cohesion: 0.67
Nodes (3): Project Risks, Risk Management, Technical Risks

### Community 258 - "seed/main.go"
Cohesion: 0.42
Nodes (9): lastStep(), main(), monitorKind(), mustExec(), printSummary(), seedChecks(), seedIncidents(), seedMonitors() (+1 more)

### Community 260 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 261 - "serve_test.go"
Cohesion: 0.24
Nodes (8): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep

### Community 262 - "12. Índices e desempenho do banco"
Cohesion: 0.33
Nodes (6): 12. Índices e desempenho do banco, Decisões, Descartado, Entregue, Limitações conhecidas, Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

### Community 263 - "SendTestWebhook"
Cohesion: 0.33
Nodes (6): newWebhookTestClient(), SendTestWebhook(), TestSendTestWebhook(), webhookTestErrorMessage(), validateWebhookURL(), WebhookTestResult

### Community 265 - "13. Aviso de expiração do certificado SSL"
Cohesion: 0.29
Nodes (7): 13. Aviso de expiração do certificado SSL, Decisões, Descartado, Desvios do plano, Entregue, Limitações conhecidas, Validação

### Community 266 - "Check"
Cohesion: 0.14
Nodes (7): Check, boolToInt(), encodeConfig(), nullableString(), nullableTime(), checkRepoWithTx, monitorRepoWithTx

### Community 267 - "siteCert"
Cohesion: 0.47
Nodes (3): loadOrCreateCA(), main(), siteCert()

### Community 272 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 273 - "NewMultiNotifier"
Cohesion: 0.38
Nodes (6): TestMultiNotifierCertExpiringFanOut(), NewMultiNotifier(), TestMultiNotifierDoesNotBlock(), TestMultiNotifierIgnoresCanceledContext(), TestMultiNotifierIsolatesFailures(), TestMultiNotifierWait()

## Knowledge Gaps
- **23 isolated node(s):** `github.com/Bradrickcruz/pingou`, `TMPDIR`, `autoprefixer`, `eslint-config-prettier`, `postcss` (+18 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 579 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Monitor` connect `Monitor` to `serve.go`, `time.Time`, `newCertEnv`, `net/http.ResponseWriter`, `newEnv`, `sqliteUnitOfWork`, `MonitorType`, `Check`, `bodyFor`, `Incident`, `NewMultiNotifier`, `newTestRegistry`, `Notifier`, `context.Context`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `github.com/Bradrickcruz/pingou`, `TMPDIR`, `autoprefixer` to the rest of the system?**
  _23 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09923011120615911 - nodes in this community are weakly interconnected._
- **Why does `Histórico de implementações (planos aplicados)` connect `Histórico de implementações (planos aplicados)` to `3. Página de Configurações`, `10. Pausar monitor fecha o incidente`, `12. Índices e desempenho do banco`, `MonitorDetail.jsx`, `13. Aviso de expiração do certificado SSL`, `Pingou README`, `8. Notificação via ntfy`, `1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel`, `6. E-mail (SMTP) e envio fora do lock`, `9. Eixo Y dinâmico + gráfico de pings`, `11. Fuso horário + desligamento ordenado`, `7. Largura no desktop + busca parcial`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `Task Breakdown` be split into smaller, more focused modules?**
  _Cohesion score 0.05405405405405406 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.1437908496732026 - nodes in this community are weakly interconnected._
- **Should `net/http.ResponseWriter` be split into smaller, more focused modules?**
  _Cohesion score 0.060833902939166094 - nodes in this community are weakly interconnected._