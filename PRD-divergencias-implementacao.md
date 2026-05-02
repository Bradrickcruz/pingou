# Pingou - Divergencias entre PRD e implementacao

Analise feita em 2026-05-01, comparando `PRD.md` com a arvore e codigo atuais do projeto.

## Resumo executivo

O projeto esta bem mais avancado que o status declarado no PRD (`Plano aprovado, aguardando inicio da Fase 1`). Ja existe backend Go funcional, SQLite com migrations, scheduler, state machine, webhooks, settings, retention, export, frontend React embutido e Docker.

As maiores divergencias nao sao falta total de implementacao, mas diferenca de arquitetura e contratos:

- Estrutura real usa pacotes horizontais (`domain`, `repository`, `service`, `handler`, `checker`, `scheduler`), enquanto PRD planeja pacotes por dominio (`monitors`, `checks`, `incidents`, `notifications`, `settings`, `server`, `export`, `auth`).
- Frontend real usa React 19 + Vite + JavaScript, sem TypeScript e sem Tailwind.
- UI embutida fica em `internal/handler/dist`, nao em `ui/dist` com `ui/embed.go`.
- Nao ha comandos CLI `migrate` e `export`.
- Nao ha testes Go no repositorio.
- Auth nao usa comparacao em tempo constante.

## Divergencias de criterio de sucesso

| Criterio do PRD                                     | Estado real observado                                               |
| --------------------------------------------------- | ------------------------------------------------------------------- |
| 100 URLs a cada 10s com < 50MB RAM                  | Nao validado; scheduler suporta ate 100 por regra                   |
| Webhook em transicoes UP/DOWN respeitando threshold | Implementado: transicoes existem, webhook disparado, chamada direta |
| Roda 7 dias sem leak                                | Nao validado                                                        |
| Export SQLite via API e CLI                         | Parcial: API existe em `/api/export`; CLI nao existe                |

## Recomendacoes

1. Decidir se o PRD continua sendo plano ideal ou se deve virar documentacao da realidade atual.
2. Atualizar PRD/README para refletir arquitetura real se a decisao for preservar implementacao atual.
3. Corrigir divergencias de seguranca/operacao antes de evoluir features:
   - graceful shutdown real do `http.Server`.
4. Criar testes minimos para state machine, monitor service, checker HTTP e handlers principais.

## Arquivos-chave analisados

- `PRD.md`
- `README.md`
- `go.mod`
- `cmd/pingou/main.go`
- `internal/config/config.go`
- `internal/database/database.go`
- `internal/database/database_test.go`
- `internal/database/migrations/0001_init.sql`
- `internal/database/migrations/0002_add_duration_seconds_to_incidents.sql`
- `internal/handler/server.go`
- `internal/handler/middleware.go`
- `internal/handler/monitors.go`
- `internal/handler/checks.go`
- `internal/handler/incidents.go`
- `internal/handler/settings.go`
- `internal/handler/export.go`
- `internal/handler/spa.go`
- `internal/checker/http_checker.go`
- `internal/scheduler/scheduler.go`
- `internal/scheduler/retention.go`
- `internal/service/monitor_service.go`
- `internal/service/state_machine.go`
- `internal/service/notifier.go`
- `internal/service/settings_service.go`
- `web/package.json`
- `web/vite.config.js`
- `Dockerfile`
- `docker-compose.yml`
- `.env.example`
- `.editorconfig`
