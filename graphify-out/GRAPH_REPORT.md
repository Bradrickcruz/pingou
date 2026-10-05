# Graph Report - pingou-health-checker  (2026-10-05)

## Corpus Check
- 111 files · ~91,996 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 15 file(s) not represented in the graph (top: (none) 9, .woff2 4, .example 1)

## Summary
- 1041 nodes · 2007 edges · 81 communities (66 shown, 15 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 92 edges (avg confidence: 0.91)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `96e7a4a1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- serve.go
- Dashboard.jsx
- database/sql.DB
- package.json
- Scheduler
- Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)
- net/http.ResponseWriter
- Testing Strategy & Framework
- Monitor
- Incident
- info.go
- Check
- boolToInt
- WebhookNotifier
- Pingou Health Checker - Project Vision & Goals
- UnitOfWork
- context.Context
- Divergências entre PRD e Implementação
- Vite Logo
- Docker Compose Configuration
- github.com/Bradrickcruz/pingou
- Code Conventions & Standards
- Architecture Overview
- Technical Concerns & Risk Assessment
- Integrations & External Dependencies
- Technology Stack Analysis
- Project Structure Analysis
- What You Must Do When Invoked
- Plano de Desenvolvimento: {{Titulo}}
- Plano de Desenvolvimento: Fase 4 — Migração para o tema claro do DESIGN.md
- Hero Image
- React Logo
- MonitorService
- Task Breakdown
- Task Breakdown
- Settings.jsx
- main.jsx
- 🧩 Tarefa: Melhorias na página de Incidentes
- Por task
- 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes
- 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores
- graphify reference: extra exports and benchmark
- validateCreateInput
- MonitorForm.jsx
- graphify reference: query, path, explain
- Pingou - Divergencias entre PRD e implementacao
- Pingou Health Checker - Project Roadmap
- Status: **In Progress**
- Status: **Planned**
- Status: **Planned**
- Status: **Planned**
- Status: **Future**
- Status: **Continuous**
- Pingou Project State & Memory
- graphify reference: add a URL and watch a folder
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- Success Metrics
- Release Schedule
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
- graphify reference: transcribe video and audio
- Dependencies & Blockers
- Risk Management
- Project Status
- Risk Management
- Team & Resources
- CLAUDE.md
- .claude/CLAUDE.md
- extraction-spec.md

## God Nodes (most connected - your core abstractions)
1. `Monitor` - 40 edges
2. `Incident` - 27 edges
3. `Task Breakdown` - 26 edges
4. `Task Breakdown` - 23 edges
5. `runServe()` - 21 edges
6. `Scheduler` - 18 edges
7. `Plano de Desenvolvimento: Fase 4 — Migração para o tema claro do DESIGN.md` - 17 edges
8. `MonitorService` - 16 edges
9. `Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)` - 15 edges
10. `Check` - 14 edges

## Surprising Connections (you probably didn't know these)
- `Task T5: Ícones no `Shell` ✅` --references--> `Shell()`  [INFERRED]
  plan-design-system-web.md → web/src/components/layout/Shell.jsx
- `Notes` --references--> `Shell()`  [INFERRED]
  plan-tema-claro.md → web/src/components/layout/Shell.jsx
- `Por task` --references--> `Shell()`  [INFERRED]
  plan-tema-claro.md → web/src/components/layout/Shell.jsx
- `Por task` --references--> `MonitorCard()`  [INFERRED]
  plan-tema-claro.md → web/src/components/monitors/MonitorCard.jsx
- `Por task` --references--> `MonitorForm()`  [INFERRED]
  plan-tema-claro.md → web/src/components/monitors/MonitorForm.jsx

## Import Cycles
- None detected.

## Communities (81 total, 15 thin omitted)

### Community 0 - "serve.go"
Cohesion: 0.12
Nodes (35): copyFile(), runExportDB(), Execute(), requireKey(), acquireLock(), checkExistingLock(), main(), MonitorEnabled (+27 more)

### Community 1 - "Dashboard.jsx"
Cohesion: 0.17
Nodes (13): Task T6: Ícones em `Dashboard`, `Modal`, `Login` ✅, Task T10: `Dashboard.jsx` + `MonitorCard.jsx` ✅, Task T11: `MonitorForm.jsx` + `Modal.jsx` em classes ✅, lucide-react, monitorsApi, web_src_assets_brand_pingou_simbolo, dotColor, MonitorCard() (+5 more)

### Community 2 - "database/sql.DB"
Cohesion: 0.06
Nodes (57): runAdd(), runMigrateDown(), runMigrateStatus(), runMigrateUp(), runRM(), releaseLock(), runServe(), go_pkg_embed (+49 more)

### Community 3 - "package.json"
Cohesion: 0.05
Nodes (40): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, tailwindcss (+32 more)

### Community 4 - "Scheduler"
Cohesion: 0.26
Nodes (5): context.CancelFunc, sync.Mutex, Checker, job, Scheduler

### Community 5 - "Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)"
Cohesion: 0.05
Nodes (43): 0. Como usar este arquivo, 10. Pendências (não definidas na v1.0), 11. Checklist antes de entregar, 1. Conceito, 2. Marca, 3.1 Tokens, 3.2 Tailwind, 3.3 Proporção (+35 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (31): go_pkg_crypto_rand, go_pkg_crypto_subtle, go_pkg_encoding_hex, go_pkg_runtime_debug, go_pkg_strings, net/http.Handler, net/http.Request, net/http.ResponseWriter (+23 more)

### Community 7 - "Testing Strategy & Framework"
Cohesion: 0.05
Nodes (42): 1. Unit Testing, 2. Integration Testing, 3. End-to-End Testing, Backend Integration Tests, Backend Mocking, Backend Test Environment, Backend Test Structure, Backend Testing Guidelines (+34 more)

### Community 8 - "Monitor"
Cohesion: 0.19
Nodes (9): time.Time, CheckResult, Monitor, MonitorState, MonitorFilter, Notifier, NewStateMachine(), monitorRepoWithTx (+1 more)

### Community 9 - "Incident"
Cohesion: 0.16
Nodes (7): Incident, IncidentRepository, scanIncident(), IncidentService, NewIncidentService(), IncidentRepo, noopNotifier

### Community 10 - "info.go"
Cohesion: 0.39
Nodes (8): getEnvInt(), getEnvList(), getEnvOr(), runConfig(), runVersion(), split(), splitAndTrim(), trim()

### Community 11 - "Check"
Cohesion: 0.13
Nodes (5): database/sql.Tx, Check, CheckRepositoryTx, checkRepoWithTx, incidentRepoWithTx

### Community 12 - "boolToInt"
Cohesion: 0.18
Nodes (5): boolToInt(), nullableTime(), scanMonitor(), MonitorRepo, scanner

### Community 13 - "WebhookNotifier"
Cohesion: 0.32
Nodes (5): HTTPChecker, net/http.Client, NewHTTPChecker(), NewWebhookNotifier(), WebhookNotifier

### Community 14 - "Pingou Health Checker - Project Vision & Goals"
Cohesion: 0.05
Nodes (42): 1. Convention Over Configuration, 1. Simplicity First, 2. Batteries Included, 2. Reliability & Performance, 3. Extensibility by Design, 3. Self-Hosted Independence, 4. Developer Experience, 5. Extensibility (+34 more)

### Community 15 - "UnitOfWork"
Cohesion: 0.13
Nodes (3): IncidentRepositoryTx, MonitorRepositoryTx, UnitOfWork

### Community 16 - "context.Context"
Cohesion: 0.21
Nodes (5): context.Context, CheckRepository, MonitorRepository, NewRetentionWorker(), RetentionWorker

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

### Community 29 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 30 - "Plano de Desenvolvimento: {{Titulo}}"
Cohesion: 0.09
Nodes (22): Alternativas consideradas e descartadas, Build Verification, Concluído, Cronograma Sugerido, Definition of Done (DoD), Em andamento, Fase 1, File Structure (+14 more)

### Community 31 - "Plano de Desenvolvimento: Fase 4 — Migração para o tema claro do DESIGN.md"
Cohesion: 0.10
Nodes (21): Alternativas consideradas e descartadas, Concluído, Cronograma Sugerido, Decisões pendentes, Definition of Done (DoD), Em andamento, File Structure, Mapeamento de tokens (fonte da T2 e T3) (+13 more)

### Community 34 - "MonitorService"
Cohesion: 0.26
Nodes (5): MonitorService, newUUIDv7(), CreateMonitorInput, Reloader, UpdateMonitorInput

### Community 35 - "Task Breakdown"
Cohesion: 0.13
Nodes (18): Fase 1: Tokens do DESIGN.md no Tailwind, Fase 3: Fontes locais, Task Breakdown, Task T10: `@font-face` + ativar fontes nos tokens ✅, Task T11: Remover Google Fonts do `globalStyles.js` ✅, Task T12: Atualizar `DESIGN.md` §5 ✅, Task T1: Variáveis `--pg-*` em `:root` ✅, Task T2: Tokens `pg` no `tailwind.config.js` ✅ (+10 more)

### Community 36 - "Task Breakdown"
Cohesion: 0.12
Nodes (17): Fase 4A: Troca centralizada para o tema claro, Fase 4B: Migração gradual dos componentes, Fase 4C: Limpeza e verificação, Task Breakdown, Task T0: Commitar Fases 1–3 ✅, Task T12: `Incidents.jsx` ✅, Task T15: Remover `tokens.js` e legados sem uso ✅ (exceto `darkMode: 'class'`, pendente), Task T16: Resolver os 2 erros de lint pré-existentes ✅ (+9 more)

### Community 37 - "Settings.jsx"
Cohesion: 0.20
Nodes (8): Fase 2: Ícones com lucide-react, axios, base, client, incidentsApi, settingsApi, useSettings(), Settings()

### Community 38 - "main.jsx"
Cohesion: 0.20
Nodes (9): react, react-dom, react-router-dom, web_src_assets_brand_pingou_marca_horizontal, nav, ConnectionContext, ConnectionProvider(), web_src_index (+1 more)

### Community 39 - "🧩 Tarefa: Melhorias na página de Incidentes"
Cohesion: 0.15
Nodes (12): 🔍 Contexto, Detalhamento, 📦 Escopo inicial, 🚧 Fora de escopo, Funcionalidades de visualização, 💡 Ideia, Informações adicionais, 🎯 Objetivo (+4 more)

### Community 40 - "Por task"
Cohesion: 0.29
Nodes (9): Correções pós-T18 (fora do plano original), Por task, Resultado da execução, Task T14: `Login.jsx`, `Spinner`, `ConnectionOverlay`, `TimeRangeSlider` ✅, Task T7: Verificação visual da troca de tema (4A) ✅, Verificação, ConnectionOverlay(), Spinner() (+1 more)

### Community 41 - "🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Tela de detalhe do monitor com checks e incidentes (+1 more)

### Community 42 - "🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores"
Cohesion: 0.20
Nodes (9): 🔍 Contexto, 📦 Escopo inicial, 🚧 Fora de escopo, 💡 Ideia, 🎯 Objetivo, ⚠️ Observações, ⏱️ Prioridade, 🧩 Tarefa: Sistema WebSocket para atualização em tempo real dos monitores (+1 more)

### Community 43 - "graphify reference: extra exports and benchmark"
Cohesion: 0.22
Nodes (8): graphify reference: extra exports and benchmark, Step 6b - Wiki (only if --wiki flag), Step 7 - Neo4j export (only if --neo4j or --neo4j-push flag), Step 7a - FalkorDB export (only if --falkordb or --falkordb-push flag), Step 7b - SVG export (only if --svg flag), Step 7c - GraphML export (only if --graphml flag), Step 7d - MCP server (only if --mcp flag), Step 8 - Token reduction benchmark (only if total_words > 5000)

### Community 44 - "validateCreateInput"
Cohesion: 0.54
Nodes (8): validateCreateInput(), validateInterval(), validateMonitor(), validateName(), validateThreshold(), validateTimeout(), validateUpdateInput(), validateURL()

### Community 45 - "MonitorForm.jsx"
Cohesion: 0.38
Nodes (4): defaults, variants, formatTime(), TimeRangeSlider()

### Community 46 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 47 - "Pingou - Divergencias entre PRD e implementacao"
Cohesion: 0.33
Nodes (5): Arquivos-chave analisados, Divergencias de criterio de sucesso, Pingou - Divergencias entre PRD e implementacao, Recomendacoes, Resumo executivo

### Community 48 - "Pingou Health Checker - Project Roadmap"
Cohesion: 0.33
Nodes (5): Development Focus, Development Phases Overview, Infrastructure Costs, Pingou Health Checker - Project Roadmap, Resource Allocation

### Community 49 - "Status: **In Progress**"
Cohesion: 0.40
Nodes (5): Milestone 1.1: Testing Infrastructure ✅, Milestone 1.2: Code Quality & Documentation ✅, Milestone 1.3: Security & Reliability 🔄, Phase 1: Foundation & Stability (Current - Q1 2024), Status: **In Progress**

### Community 50 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 2.1: Advanced HTTP Monitoring, Milestone 2.2: TCP & Network Monitoring, Milestone 2.3: Notification System Enhancement, Phase 2: Enhanced Monitoring Capabilities (Q2 2024), Status: **Planned**

### Community 51 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 3.1: Dashboard Enhancements, Milestone 3.2: Basic Analytics & Reporting, Milestone 3.3: Data Management, Phase 3: User Experience & Analytics (Q3 2024), Status: **Planned**

### Community 52 - "Status: **Planned**"
Cohesion: 0.40
Nodes (5): Milestone 4.1: Multi-tenancy & Access Control, Milestone 4.2: Advanced Authentication, Milestone 4.3: Scalability Improvements, Phase 4: Enterprise Features (Q4 2024), Status: **Planned**

### Community 53 - "Status: **Future**"
Cohesion: 0.40
Nodes (5): Milestone 5.1: Advanced Monitoring, Milestone 5.2: Integration Ecosystem, Milestone 5.3: Public Features, Phase 5: Advanced Features & Integrations (Q1-Q2 2025), Status: **Future**

### Community 54 - "Status: **Continuous**"
Cohesion: 0.40
Nodes (5): Milestone 6.1: Community Building, Milestone 6.2: Developer Experience, Milestone 6.3: Ecosystem Integration, Phase 6: Ecosystem & Community (Ongoing), Status: **Continuous**

### Community 55 - "Pingou Project State & Memory"
Cohesion: 0.40
Nodes (4): Current Metrics, Metrics & KPIs, Pingou Project State & Memory, Success Metrics

### Community 56 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 57 - "graphify reference: commit hook and native CLAUDE.md integration"
Cohesion: 0.50
Nodes (3): For git commit hook, For native CLAUDE.md integration, graphify reference: commit hook and native CLAUDE.md integration

### Community 58 - "graphify reference: incremental update and cluster-only"
Cohesion: 0.50
Nodes (3): For --cluster-only, For --update (incremental re-extraction), graphify reference: incremental update and cluster-only

### Community 59 - "Success Metrics"
Cohesion: 0.50
Nodes (4): Community Metrics, Success Metrics, Technical Metrics, User Metrics

### Community 60 - "Release Schedule"
Cohesion: 0.50
Nodes (4): Current Version Tracking, Release Frequency, Release Schedule, Versioning Strategy

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

### Community 73 - "Dependencies & Blockers"
Cohesion: 0.67
Nodes (3): Critical Dependencies, Dependencies & Blockers, Potential Blockers

### Community 74 - "Risk Management"
Cohesion: 0.67
Nodes (3): Project Risks, Risk Management, Technical Risks

### Community 75 - "Project Status"
Cohesion: 0.67
Nodes (3): Current Phase: Foundation & Stability (Q1 2024), Project Status, Recent Achievements

### Community 76 - "Risk Management"
Cohesion: 0.67
Nodes (3): Current Risks, Mitigation Status, Risk Management

### Community 77 - "Team & Resources"
Cohesion: 0.67
Nodes (3): Current Team, Resource Needs, Team & Resources

## Knowledge Gaps
- **434 isolated node(s):** `github.com/Bradrickcruz/pingou`, `MonitorEnabled`, `Server`, `Server`, `Server` (+429 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 503 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Task Breakdown` connect `Task Breakdown` to `Dashboard.jsx`, `Settings.jsx`, `Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)` connect `Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)` to `Task Breakdown`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Why does `Task Breakdown` connect `Task Breakdown` to `Dashboard.jsx`, `Task Breakdown`, `Por task`, `MonitorForm.jsx`, `Plano de Desenvolvimento: Fase 4 — Migração para o tema claro do DESIGN.md`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `Task Breakdown` (e.g. with `ConnectionOverlay()` and `Spinner()`) actually correct?**
  _`Task Breakdown` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/Bradrickcruz/pingou`, `MonitorEnabled`, `Server` to the rest of the system?**
  _434 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `serve.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11655266757865937 - nodes in this community are weakly interconnected._
- **Should `database/sql.DB` be split into smaller, more focused modules?**
  _Cohesion score 0.056338028169014086 - nodes in this community are weakly interconnected._