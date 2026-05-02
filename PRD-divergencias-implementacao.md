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
- Webhook implementado com chamada **direta e síncrona** (decisão registrada em `state-machine-channel-event-necessidade-e-melhorias.md`).

Resolvido nesta rodada:

- SQLite padronizado como `github.com/mattn/go-sqlite3` com CGO no PRD.
- SQLite agora configura `_busy_timeout=5000` junto com `_foreign_keys=on` e `_journal_mode=WAL`, preservando query params existentes no DSN.
- Env vars publicas padronizadas como `PINGOU_*` no PRD e README.
- Export padronizado como `GET /api/export` no PRD, README, backend e frontend.
- Payload real do webhook documentado no PRD e README.
- `.editorconfig` implementado na raiz do projeto e documentado no PRD/README.

## Estado real por area

### Settings, retention e export

Implementado:

- Settings `webhook_url` e `retention_days`.
- Validacao `retention_days` entre 7 e 90.
- Retention worker roda na inicializacao e depois a cada 1h.
- Export HTTP baixa arquivo SQLite em `GET /api/export`.

Divergencias:

- PRD pede retention diario (`24h`); implementacao roda a cada 1h.
- PRD pede export com `VACUUM INTO` ou copia com lock; implementacao abre e copia o arquivo direto, sem lock aparente.

### Frontend e embed

Implementado:

- Frontend React/Vite existe em `web/`.
- Paginas: Login, Dashboard, Incidents, Settings.
- API client usa `X-API-Key`.
- API key usa `localStorage` com chave `pingou_api_key`.
- Build gera assets em `internal/handler/dist`.
- Backend embute `internal/handler/dist` via `embed.FS` e serve SPA fallback.

Divergencias:

- PRD pede detalhe do monitor com ultimos checks + incidentes; rotas reais sao `/`, `/incidents` e `/settings`, sem rota de detalhe por monitor.
- PRD pede pagina/export com botao de download; nao apareceu pagina dedicada de export nos arquivos principais.
- Botao de export existente em Settings aponta para `GET /api/export`.
- `web/src/App.jsx` ainda contem app padrao do Vite, mas nao e usado por `web/src/main.jsx`.
- PRD pede `ui/embed.go` e `ui/dist`; implementacao usa `internal/handler/spa.go` e `internal/handler/dist`.
- PRD pede `scripts/build.sh`; nao existe.

### Docker e distribuicao

Implementado:

- `Dockerfile` multi-stage: Node build, Go build, Alpine final.
- `docker-compose.yml` com volume `/data`, porta `8080`, env vars `PINGOU_*`.
- `HEALTHCHECK` no Dockerfile via `GET /healthz` (30s interval, 10s timeout, 3 retries).

Divergencias:

- PRD espera imagem final < 30MB; com Alpine + CGO isso precisa ser medido, nao garantido.
- Dockerfile contem copia suspeita: `COPY --from=web-builder /app/internal/handler/dist ./internal/handler/dist`. O stage `web-builder` trabalha em `/app/web` e Vite gera `../internal/handler/dist`, logo o caminho pode existir por efeito do `outDir`, mas e fragil e diferente do desenho do PRD.

### Documentacao e verificacao

Implementado:

- README existe e descreve app, endpoints, Docker, env vars e comandos.
- LICENSE existe.
- `.env.example` existe.
- `.editorconfig` existe e esta documentado como convencao basica de edicao.

Divergencias:

- PRD pede GitHub Action com `go test ./...` e `golangci-lint`; nao ha workflow visivel.
- PRD pede checklist manual executado; nao ha registro.
- Nao ha testes automatizados Go.

## Divergencias de criterio de sucesso

| Criterio do PRD                                     | Estado real observado                                               |
| --------------------------------------------------- | ------------------------------------------------------------------- |
| `docker compose up` sobe em < 30s                   | Nao validado nesta analise                                          |
| 100 URLs a cada 10s com < 50MB RAM                  | Nao validado; scheduler suporta ate 100 por regra                   |
| Webhook em transicoes UP/DOWN respeitando threshold | Implementado: transicoes existem, webhook disparado, chamada direta |
| UI React embutida em `http://localhost:8080`        | Implementado via `internal/handler/dist`                            |
| Roda 7 dias sem leak                                | Nao validado                                                        |
| Outro dev sobe em < 10 min via README               | Env vars principais foram padronizadas; setup completo nao validado |
| Export SQLite via API e CLI                         | Parcial: API existe em `/api/export`; CLI nao existe                |

## Recomendacoes

1. Decidir se o PRD continua sendo plano ideal ou se deve virar documentacao da realidade atual.
2. Atualizar PRD/README para refletir arquitetura real se a decisao for preservar implementacao atual.
3. Corrigir divergencias de seguranca/operacao antes de evoluir features:
   - graceful shutdown real do `http.Server`.
   - `HEALTHCHECK` no Dockerfile.
4. Criar testes minimos para state machine, monitor service, checker HTTP e handlers principais.
5. Escolha SQLite resolvida: manter `mattn/go-sqlite3` com CGO.
6. Operacao SQLite resolvida nesta rodada: `_busy_timeout=5000` implementado com montagem de DSN compativel com query params existentes.
7. Contratos publicos resolvidos nesta rodada: env vars `PINGOU_*`, `GET /api/export` e payload real do webhook.
8. Convencao basica de edicao resolvida nesta rodada: `.editorconfig` implementado.
9. **Arquitetura de notificacoes resolvida (2026-05-02)**: Manter webhook com chamada **direta e síncrona** (sem channels/dispatcher). Justificativa em `state-machine-channel-event-necessidade-e-melhorias.md`.
10. **PRD.md atualizado (2026-05-02)**: Fase 8 renomeada para "Webhook Simples", conceitos e tabela refletem implementacao atual.
11. **PRD-divergencias-implementacao.md atualizado (2026-05-02)**: Todas as mencoes a channel/event/dispatcher/retry removidas ou atualizadas.
12. **Middleware recover e CORS implementados (2026-05-02)**: Adicionados em `internal/handler/middleware.go` com documentação em `http-recover-cors-necessidade-e-melhorias.md`.
13. **HEALTHCHECK implementado (2026-05-02)**: Dockerfile agora inclui HEALTHCHECK que verifica `/healthz` a cada 30s.
14. **Frontend React stack resolvido (2026-05-02)**: PRD.md atualizado para refletir React 19 + JavaScript + CSS/tokens (sem TypeScript/Tailwind) como implementação real.

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
