# Graph Report - pingou  (2026-11-04)

## Corpus Check
- 213 files · ~130,727 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 16 file(s) not represented in the graph (top: (none) 10, .woff2 4, .example 1)

## Summary
- 1943 nodes · 6230 edges · 112 communities (93 shown, 19 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 530 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c7efea03`
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
- Scheduler
- 🧩 Task: Aproveitar melhor a largura no desktop
- settings_test.go
- dbbench: desempenho do banco com volume real
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
- 5. Remodelagem mobile
- Plano de Desenvolvimento: {{Titulo}}
- MaintenanceConflictError
- Hero Image
- React Logo
- SettingsRepo
- newCertEnv
- Incident
- testing.T
- sqliteUnitOfWork
- format.js
- newEnv
- vite
- Resource Allocation
- Pingou - Divergencias entre PRD e implementacao
- Pingou Health Checker - Project Roadmap
- Notifier
- newTestRegistry
- 🧩 Tarefa: Melhorias na página de Incidentes
- DESIGN.md — Design system do Pingou
- service.MonitorService
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
- context.Context
- NewNtfyNotifier
- Status: **Planned**
- Histórico de implementações (planos aplicados)
- migrations_test.go
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- Status: **Planned**
- devDependencies
- runServe
- NewUnitOfWork
- Dashboard.jsx
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- MonitorDetail.jsx
- 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel
- scripts
- 9. Eixo Y dinâmico + gráfico de pings
- dependencies
- newFakeSMTP
- Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)
- commands/info.go
- Success Metrics
- Release Schedule
- Risk Management
- seed/main.go
- root.go
- serve_test.go
- 12. Índices e desempenho do banco
- 13. Aviso de expiração do certificado SSL
- Check
- siteCert
- dbbench.sh
- dbmem.sh
- 8. Notificação via ntfy
- NewMultiNotifier

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 113 edges
2. `Incident` - 62 edges
3. `Scheduler` - 41 edges
4. `MaintenanceWindow` - 33 edges
5. `newTestServer()` - 33 edges
6. `runServe()` - 31 edges
7. `do()` - 31 edges
8. `react` - 31 edges
9. `MonitorDetail()` - 29 edges
10. `newEnv()` - 29 edges

## Surprising Connections (you probably didn't know these)
- `Descartado` --references--> `NewServer()`  [INFERRED]
  docs/implementacoes.md → internal/handler/server.go
- `Entregue` --references--> `matchesQuery()`  [INFERRED]
  docs/implementacoes.md → web/src/lib/monitors.js
- `Entregue` --references--> `MobileHeader()`  [INFERRED]
  docs/implementacoes.md → web/src/components/layout/MobileHeader.jsx
- `Decisões` --references--> `ResponseChart()`  [INFERRED]
  docs/implementacoes.md → web/src/components/monitor-detail/ResponseChart.jsx
- `Limitações conhecidas` --references--> `runShutdown()`  [INFERRED]
  docs/implementacoes.md → cmd/pingou/commands/serve.go

## Import Cycles
- 2-file cycle: `internal/repository/dbbench_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 2-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/repository/dbbench_test.go`
- 2-file cycle: `internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 3-file cycle: `internal/repository/dbbench_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 4-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/monitor_pause_test.go -> internal/service/state_machine_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`
- 5-file cycle: `internal/repository/dbbench_test.go -> internal/service/cert_state_machine_test.go -> internal/service/monitor_pause_test.go -> internal/service/stats_test.go -> internal/service/unit_of_work_impl.go -> internal/repository/dbbench_test.go`

## Communities (112 total, 19 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.11
Nodes (33): MonitorEnabled, createMaintenanceRequest, createMonitorRequest, downWindowResponse, infoResponse, monitorConfigRequest, paginatedResponse, statsBucketResponse (+25 more)

### Community 1 - "Task Breakdown"
Cohesion: 0.11
Nodes (20): Task Breakdown, Task T0: Criar branch, Task T10: Contrato da API de manutenção, Task T15: Testes da etapa 2 (backend), Task T19: Revisão de bordas e segurança, Task T1: Migração 0008 (pausa + close_note), Task T20: Docs, Task T21: Verificação final e PR (+12 more)

### Community 2 - "MonitorForm.jsx"
Cohesion: 0.09
Nodes (22): CertCard(), levels, blank, CHECK_EVERY, friendlyError(), fromMonitor(), GRACE, PING_EVERY (+14 more)

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
Cohesion: 0.05
Nodes (43): checkResponse, incidentMonitorResponse, incidentResponse, maintenanceResponse, monitorListItem, monitorResponse, rateLimiter, responseWriter (+35 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "MonitorDetail.jsx"
Cohesion: 0.13
Nodes (27): 4. Detalhe do monitor, Descartado, Desvios na implementação, Entregue, lucide-react, ChecksTable(), detailOf(), ConfigCard() (+19 more)

### Community 9 - "Scheduler"
Cohesion: 0.14
Nodes (11): Desvios do plano, Entregue, MonitorRepository, Alternativas consideradas e descartadas, Notes, Risk Areas, Task T12: Gate de manutenção no scheduler e no ping, Task T13: Tick de início/fim da janela (+3 more)

### Community 10 - "🧩 Task: Aproveitar melhor a largura no desktop"
Cohesion: 0.05
Nodes (37): 🔍 Context, 🔍 Context, 🔍 Context, 🔍 Context, 🎯 Goal, 🎯 Goal, 🎯 Goal, 🎯 Goal (+29 more)

### Community 11 - "settings_test.go"
Cohesion: 0.12
Nodes (36): settingsBody, TestMonitorCertAlertField(), TestMonitorCertAlertValidation(), createHeartbeat(), decode(), do(), newTestServer(), newTestServerSMTP() (+28 more)

### Community 12 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.40
Nodes (5): Como interpretar, Como usar, dbbench: desempenho do banco com volume real, O que é medido, Os dados

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
Cohesion: 0.05
Nodes (22): noopNotifier, certNotAfter(), requestErrorMessage(), tcpErrorMessage(), CertAlert, CheckResult, Monitor, toWebhookMonitor() (+14 more)

### Community 17 - "dbbench: desempenho do banco com volume real"
Cohesion: 0.15
Nodes (10): MVP Plan: Pingou, Pingou README, certserver: testar o aviso de certificado SSL, Como interpretar, Como usar, Como usar (Linux), Limitações, memtest: prova da promessa de memória (+2 more)

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

### Community 29 - "5. Remodelagem mobile"
Cohesion: 0.50
Nodes (4): 5. Remodelagem mobile, Decisões, Descartado, Entregue

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 34 - "SettingsRepo"
Cohesion: 0.27
Nodes (4): SettingsRepo, NewSettingsRepo(), NewRetentionWorker(), RetentionWorker

### Community 36 - "newCertEnv"
Cohesion: 0.37
Nodes (13): fail(), newCertEnv(), ok(), TestCertAlertDisabledStillStoresDate(), TestCertAlertGuardedByURL(), TestCertAlertLateActivation(), TestCertAlertRenewal(), TestCertAlertRollback() (+5 more)

### Community 38 - "Incident"
Cohesion: 0.07
Nodes (18): Entregue, IncidentRepository, CloseReason, Incident, IncidentFilter, MonitorConfig, MonitorType, scanIncident() (+10 more)

### Community 39 - "testing.T"
Cohesion: 0.08
Nodes (41): newCertServer(), TestHTTPCheckerExpiredCertFails(), TestHTTPCheckerReadsCertExpiry(), TestNtfyHost(), TestMaintenanceRoutesNeedAPIKey(), env, newHeartbeat(), TestHeartbeatDeadlineMissedWhileOffline() (+33 more)

### Community 40 - "sqliteUnitOfWork"
Cohesion: 0.11
Nodes (17): 10. Pausar monitor fecha o incidente, 6. E-mail (SMTP) e envio fora do lock, Decisões, Decisões, Descartado, Descartado, Entregue, Entregue (+9 more)

### Community 42 - "format.js"
Cohesion: 0.09
Nodes (45): Task T16: `zonedToUtc` + status MAINTENANCE, OpenIncidentCard(), axisTicks(), gapBuckets(), isDown(), linePath(), ResponseChart(), summary() (+37 more)

### Community 43 - "newEnv"
Cohesion: 0.23
Nodes (19): runAdd(), runRM(), NewCheckRepo(), openBench(), TestResultsSameWithOldIndex(), NewIncidentRepo(), NewMonitorRepo(), newEnv() (+11 more)

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
Cohesion: 0.20
Nodes (16): Registry, TCPChecker, TestHTTPCheckerNoCertWithoutTLS(), newTestRegistry(), tcpMonitor(), TestHTTPCheckerErrorMessage(), TestRegistryHTTP(), TestRegistryUnknownType() (+8 more)

### Community 51 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 52 - "DESIGN.md — Design system do Pingou"
Cohesion: 0.14
Nodes (14): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 4.1 Certificado SSL (página de detalhe), 4. Status, 5. Tipografia (+6 more)

### Community 53 - "service.MonitorService"
Cohesion: 0.08
Nodes (11): boolToInt(), encodeConfig(), nullableString(), nullableTime(), newPingToken(), normalizeReason(), Task T2: Domínio + repositórios da pausa e do CloseReason, MonitorRepo (+3 more)

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 59 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 60 - "main.jsx"
Cohesion: 0.08
Nodes (22): axios, react-dom, react-router-dom, version, client, incidentsApi, BottomNav(), nav (+14 more)

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
Nodes (26): TestCertMessageText(), bodyFor(), buildEmail(), certBody(), certPhrase(), formatDuration(), joinUnits(), oneLine() (+18 more)

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
Cohesion: 0.11
Nodes (31): react, settingsApi, BackupCard(), Card(), EmailCard(), resultText(), NtfyCard(), resultText() (+23 more)

### Community 87 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 89 - "openBench"
Cohesion: 0.15
Nodes (19): Entregue, copyFile(), placeholders(), sqlIncidentsOpenByMonitors(), sqlMaintenanceLiveByMonitors(), sqlSummariesHistory(), sqlSummariesUptime(), BenchDBPath() (+11 more)

### Community 96 - "context.Context"
Cohesion: 0.09
Nodes (19): DownWindow, MonitorStats, StatsBucket, MaintenanceWindow, MaintenanceRepository, fmtTime(), MaintenanceRepo, scanMaintenance() (+11 more)

### Community 97 - "NewNtfyNotifier"
Cohesion: 0.26
Nodes (12): TestNtfyNotifierCertExpiring(), NewNtfyNotifier(), decodedTitle(), newFakeNtfy(), TestNtfyNotifierDown(), TestNtfyNotifierHeartbeatNeverSendsPingToken(), TestNtfyNotifierRecoveryWithoutTokenOrLink(), TestNtfyNotifierSkipsWithoutTopic() (+4 more)

### Community 98 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 99 - "Histórico de implementações (planos aplicados)"
Cohesion: 0.14
Nodes (14): 14. Manutenção: pausa com motivo e janela agendada, 2. Página de Incidentes, 7. Largura no desktop + busca parcial, Decisões, Decisões, Decisões, Descartado, Descartado (+6 more)

### Community 235 - "migrations_test.go"
Cohesion: 0.25
Nodes (15): runMigrateUp(), Open(), OpenDB(), RunMigrations(), sqliteDSN(), assertFKOn(), assertNoFKViolations(), countRows() (+7 more)

### Community 236 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 237 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 238 - "devDependencies"
Cohesion: 0.12
Nodes (16): devDependencies, autoprefixer, eslint, eslint-config-prettier, @eslint/js, eslint-plugin-prettier, eslint-plugin-react-hooks, eslint-plugin-react-refresh (+8 more)

### Community 239 - "runServe"
Cohesion: 0.13
Nodes (23): forceExitOnSecondSignal(), releaseLock(), runServe(), Config, NtfyConfig, SMTPConfig, Load(), LoadConfig() (+15 more)

### Community 240 - "NewUnitOfWork"
Cohesion: 0.26
Nodes (18): CheckRepoTx, NewCheckRepoTx(), envInt(), TestDBBench(), IncidentRepoTx, NewIncidentRepoTx(), MonitorRepoTx, NewMonitorRepoTx() (+10 more)

### Community 241 - "Dashboard.jsx"
Cohesion: 0.13
Nodes (29): RFC-3339, Task T17: Modal de nova janela, react-toastify, maintenanceApi, monitorsApi, MobileHeader(), MaintenanceCard(), MaintenanceModal() (+21 more)

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
Cohesion: 0.11
Nodes (32): Task T18: Selos, card de detalhes, cancelar, HistoryBar(), MonitorCard(), stateLabel, PingBlock(), MonitorRow(), typeIcon, RowMenu() (+24 more)

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

### Community 250 - "newFakeSMTP"
Cohesion: 0.06
Nodes (38): HTTPChecker, updateSettingsRequest, toUpdateSettingsInput(), TestEmailNotifierSendsCertExpiring(), NewEmailNotifier(), sendMail(), SendTestEmail(), smtpErrorMessage() (+30 more)

### Community 251 - "Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)"
Cohesion: 0.13
Nodes (14): Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, File Structure, Overview, Pendente, Phase X: Verification (+6 more)

### Community 252 - "commands/info.go"
Cohesion: 0.18
Nodes (14): copyFile(), runExportDB(), getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split() (+6 more)

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
Cohesion: 0.24
Nodes (13): ExportHandler, Status(), sizes(), env, lastStep(), main(), monitorKind(), mustExec() (+5 more)

### Community 260 - "root.go"
Cohesion: 0.25
Nodes (3): Execute(), requireKey(), main()

### Community 261 - "serve_test.go"
Cohesion: 0.15
Nodes (14): acquireLock(), checkExistingLock(), runShutdown(), TestCheckExistingLockIgnoresOwnPID(), TestRunShutdownKeepsGoingAfterTimeout(), TestRunShutdownRunsStepsInOrder(), writeLockForTest(), shutdownStep (+6 more)

### Community 262 - "12. Índices e desempenho do banco"
Cohesion: 0.40
Nodes (5): 12. Índices e desempenho do banco, Decisões, Descartado, Limitações conhecidas, Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

### Community 265 - "13. Aviso de expiração do certificado SSL"
Cohesion: 0.33
Nodes (6): 13. Aviso de expiração do certificado SSL, Decisões, Descartado, Desvios do plano, Limitações conhecidas, Validação

### Community 266 - "Check"
Cohesion: 0.15
Nodes (7): CheckRepository, Check, MonitorSummary, CheckRepo, benchEnv, CheckRepositoryTx, checkRepoWithTx

### Community 267 - "siteCert"
Cohesion: 0.47
Nodes (3): loadOrCreateCA(), main(), siteCert()

### Community 272 - "8. Notificação via ntfy"
Cohesion: 0.50
Nodes (4): 8. Notificação via ntfy, Decisões, Descartado, Entregue

### Community 273 - "NewMultiNotifier"
Cohesion: 0.32
Nodes (7): TestMultiNotifierCertExpiringFanOut(), NewMultiNotifier(), TestMultiNotifierDoesNotBlock(), TestMultiNotifierIgnoresCanceledContext(), TestMultiNotifierIsolatesFailures(), TestMultiNotifierWait(), TestMultiNotifierNtfyFailureDoesNotBlockOthers()

## Knowledge Gaps
- **24 isolated node(s):** `RFC-3339`, `github.com/Bradrickcruz/pingou`, `TMPDIR`, `autoprefixer`, `eslint-config-prettier` (+19 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 570 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Task Breakdown` connect `Task Breakdown` to `context.Context`, `Scheduler`, `format.js`, `Dashboard.jsx`, `MonitorDetail.jsx`, `service.MonitorService`, `Plano de Desenvolvimento: Manutenção de monitores (pausa com motivo e janela agendada)`?**
  _High betweenness centrality (0.072) - this node is a cross-community bridge._
- **What connects `RFC-3339`, `github.com/Bradrickcruz/pingou`, `TMPDIR` to the rest of the system?**
  _24 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10723514211886305 - nodes in this community are weakly interconnected._
- **Why does `MaintenanceService` connect `context.Context` to `serve.go`, `Task Breakdown`, `net/http.ResponseWriter`, `testing.T`, `Scheduler`, `runServe`?**
  _High betweenness centrality (0.061) - this node is a cross-community bridge._
- **Should `Task Breakdown` be split into smaller, more focused modules?**
  _Cohesion score 0.11428571428571428 - nodes in this community are weakly interconnected._
- **Why does `Monitor` connect `Monitor` to `serve.go`, `context.Context`, `newCertEnv`, `net/http.ResponseWriter`, `Incident`, `sqliteUnitOfWork`, `Scheduler`, `testing.T`, `bodyFor`, `Check`, `NewMultiNotifier`, `newTestRegistry`, `Notifier`, `service.MonitorService`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `MonitorForm.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.09462365591397849 - nodes in this community are weakly interconnected._