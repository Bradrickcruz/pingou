# Plano de Desenvolvimento: Reconstrução da página de Configurações

## Overview

Reconstruir [Settings.jsx](web/src/pages/Settings.jsx) para chegar em [Configuracoes.png](concept-design/Configuracoes.png). A tela tem três cards:

- **Notificações**: URL do webhook, botão "Enviar teste" e link "Ver formato do JSON".
- **Retenção de dados**: chips de 7 a 90 dias e o tamanho atual do banco.
- **Backup**: aviso de segredo, "Baixar backup (.db)" com nome e tamanho do arquivo.

Abaixo deles, "Salvar alterações" (desabilitado sem mudanças) e um rodapé com versão, commit, licença e links.

### Conferência do briefing da IA de design contra o código (2026-10-29)

O briefing colado cobre o redesign inteiro (login, painel, incidentes, configurações). Boa parte já foi entregue pelos planos anteriores. Este plano só inclui o que a página de Configurações precisa.

| Item do briefing | Situação real | Neste plano? |
|---|---|---|
| 1.1 `GET /api/session` | Não existe. Serve só ao login | **Não**. Fica para o plano do login |
| 1.2 `GET /api/info` | Não existe. A tela precisa de versão, commit e tamanho do banco | **Sim** (T2) |
| 1.3 Nome do arquivo no export | Hoje é `pingou-backup-YYYYMMDD-HHMMSS.db` ([export.go](internal/handler/export.go)), e o front ignora e inventa `pingou_bkp_...` | **Sim** (T3) |
| 1.4 Erro HTTP sem prefixo | Feito: `requestErrorMessage` em [http_checker.go](internal/checker/http_checker.go) + migration 0004 | Não |
| 1.5 Badge da sidebar | Feito: `COUNT(*)` no repositório ([plan-incidents-page.md](plan-incidents-page.md)) | Não |
| Entrega 2 (histórico, uptime) | **Feito com outros nomes**: `history` (bool[]), `uptime_24h`, `last_latency_ms` ([README.md](README.md)). O briefing pede `recent_checks`/`last_response_time_ms` | Não. Criar campos duplicados violaria a regra "não quebre/siga convenções" |
| Entrega 3 (incidentes com monitor, duração, filtro, total) | Feito: `monitor`, `duration_seconds`, `?monitor_id=`, `?status=`, `total` no corpo | Não |
| Entrega 4 (teste do webhook) | Não existe | **Sim** (T4–T6) |
| Pergunta: a retenção apaga o quê? | Apaga **checks** e **incidentes resolvidos** mais antigos que `retention_days`. Incidentes abertos nunca. Roda **na subida e a cada 1 h**, não "uma vez por dia" ([retention.go](internal/scheduler/retention.go)) | **Sim**: a cópia da tela é corrigida (T9). O comportamento não muda |

**Divergências entre PNG e código, já decididas:**

1. A cópia do PNG diz "Os mais antigos são apagados uma vez por dia". É falso. Ajusto a **cópia** (regra 8 do briefing: não mexer na retenção): "Checks e incidentes resolvidos mais antigos que isso são apagados automaticamente. Incidentes abertos nunca."
2. O PNG mostra a mesma medida "12,4 MB" para o banco e para o backup. O backup (`VACUUM INTO`) sai **menor** que banco + WAL. Uso `database_size_bytes` nos dois lugares, com o rótulo "≈" no backup. Medir o arquivo exportado exigiria gerar o dump só para mostrar o tamanho.
3. A sidebar do PNG de Configurações não mostra "pingou v1.0" no rodapé, mas a do Painel mostra. Mantenho a sidebar como está (é compartilhada). Fora do escopo.

**Limites (fora do escopo):** `/api/session` e a tela de login; webhook por monitor; vários webhooks; troca de API Key pela tela; restaurar backup (upload); tema escuro; mudar a frequência da retenção.

**Restrições:**

- Só adicionar, não quebrar: `GET/PATCH /api/settings` mantêm o formato.
- Camadas handler → service → repository; `fmt.Errorf("...: %w")`; `slog` sem URL completa de webhook.
- Sem dependência nova. O rate limit reaproveita `rateLimiter` de [ratelimit.go](internal/handler/ratelimit.go), sem `x/time/rate`.
- Branch `feat/settings-page` a partir de `develop`. Backend (Fases 0–1) antes do frontend (Fase 2). Um commit por fase.

## Project Type

**FULLSTACK**: Go + SQLite (`internal/`) e React/Vite/Tailwind (`web/`), front embutido via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] `GET /api/info` devolve `version`, `commit`, `build_date` (os mesmos de `pingou version`) e `database_size_bytes` = banco + `-wal` + `-shm`
- [ ] `GET /api/export` responde `Content-Disposition: attachment; filename="pingou-YYYY-MM-DD.db"` (UTC)
- [ ] `POST /api/settings/webhook/test` testa a URL do corpo ou, sem corpo, a salva; responde `400` sem URL ou com esquema que não seja http(s)
- [ ] O teste responde sempre `200` com `{ok, status_code, duration_ms, error}` quando roda, sem devolver o corpo do destino
- [ ] O cliente do teste tem timeout de 5 s, não segue redirects e lê no máximo 4 KB da resposta
- [ ] Mais de 5 testes por minuto do mesmo IP → `429` com `Retry-After`
- [ ] O payload de teste tem o mesmo formato dos eventos reais, com `"event": "test"` e `"name": "Teste do Pingou"`; nada é gravado no banco
- [ ] Os logs do teste mostram só o host, nunca a URL completa
- [ ] `PATCH /api/settings` recusa `webhook_url` com esquema diferente de http/https (hoje aceita `ftp://`)
- [ ] A página bate com [Configuracoes.png](concept-design/Configuracoes.png) em desktop e é usável em 375 px
- [ ] "Salvar alterações" fica desabilitado até haver mudança e volta a desabilitar depois de salvar
- [ ] "Enviar teste" usa o valor **digitado** (sem salvar) e mostra sucesso ou falha inline, com `aria-live`
- [ ] O nome e o tamanho do backup aparecem antes do download; o arquivo baixado tem o mesmo nome
- [ ] A cópia da retenção diz a verdade: checks e incidentes resolvidos, abertos nunca
- [ ] README: evento `"test"`, aviso para ignorar `event` desconhecido, rotas `/api/info` e `/api/settings/webhook/test`

## Tech Stack

| Technology | Choice Rationale |
|---|---|
| Go stdlib (`net/http`, `os.Stat`, `io.LimitReader`) | Tudo o que o backend precisa. Sem dependência nova |
| `rateLimiter` existente | Já é janela fixa por IP, com teste. Basta uma segunda instância com `limit=5` |
| `config.Config` para versão/commit/data | `handler` não pode importar `commands` (ciclo). Três campos na config evitam mudar a assinatura de `NewServer` (2 chamadores) |
| React + Tailwind + lucide (já instalados) | Mesmo padrão das páginas de Painel e Incidentes |
| Classes `chip` como em [ChipGroup.jsx](web/src/components/ui/ChipGroup.jsx) | Os chips de retenção do PNG têm o mesmo visual arredondado. Reaproveitar as classes, não o componente (sem "Outro…") |

## File Structure

```
internal/
  config/config.go                 # + Version, Commit, BuildDate
  handler/
    server.go                      # rotas /info e /settings/webhook/test; webhookTestLimiter
    info.go                        # NOVO: handleInfo + databaseSize
    export.go                      # nome pingou-YYYY-MM-DD.db
    settings.go                    # handleTestWebhook
    handler_test.go                # casos novos (table-driven)
  service/
    notifier.go                    # send() devolve (status, err); TestWebhook()
    notifier_test.go               # casos do teste do webhook
    settings_service.go            # validateWebhookURL (http/https) + TestWebhook
cmd/pingou/commands/serve.go       # preenche cfg.Version/Commit/BuildDate
web/src/
  api/settings.js                  # + info(), testWebhook(url), export()
  hooks/useSettings.js             # + info
  pages/Settings.jsx               # reescrita
  components/settings/
    WebhookCard.jsx                # NOVO
    RetentionCard.jsx              # NOVO
    BackupCard.jsx                 # NOVO
  lib/format.js                    # + formatBytes ("12,4 MB")
README.md                          # webhook "test", /api/info, /api/settings/webhook/test
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
|---|---|---|---|---|---|
| T1 | design-api | Fechar contrato de `/api/info` e do teste do webhook | P0 | XS | - |
| T2 | go-backend-specialist | `GET /api/info` | P1 | S | T1 |
| T3 | go-backend-specialist | Nome do arquivo no export | P1 | XS | - |
| T4 | go-backend-specialist | Validar esquema http/https da `webhook_url` | P1 | XS | - |
| T5 | go-backend-specialist | `TestWebhook` no notifier/service | P1 | M | T1, T4 |
| T6 | go-backend-specialist | Rota `POST /api/settings/webhook/test` + rate limit | P1 | S | T5 |
| T7 | code-test | Testes table-driven do backend | P1 | S | T2, T3, T6 |
| T8 | frontend-specialist | Camada de API, hook e `formatBytes` | P2 | S | T1 |
| T9 | frontend-specialist | Cards Notificações, Retenção e Backup | P2 | M | T8 |
| T10 | frontend-specialist | Página: dirty state, salvar e rodapé | P2 | S | T9 |
| T11 | security-review | Revisão da rota de teste (SSRF, logs, limites) | P3 | S | T6 |
| T12 | qa-review | QA manual no binário | P3 | S | T10 |
| T13 | document | README | P3 | XS | T6 |

### Task T1: Contrato

- **Agent**: design-api
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**: briefing lido; convenção do projeto = snake_case e datas RFC3339 (confirmado em `settings_service.go` e `notifier.go`).
- **INPUT**: briefing (Entregas 1.2 e 4).
- **OUTPUT**:
  1. `GET /api/info` → `200 {"version","commit","build_date","database_size_bytes"}`. `build_date` vai como string (vem de `-ldflags`, pode ser `"unknown"` em dev).
  2. `POST /api/settings/webhook/test`, corpo opcional `{"url"}` → `200 {"ok","status_code","duration_ms","error"}`. `400 BAD_REQUEST` sem URL ou com esquema inválido. `429 RATE_LIMITED`.
  3. Mensagens de `error`: `unexpected status code: N` para respostas fora de 2xx; `timeout após 5s`; nos outros casos, a causa do `*url.Error` sem o prefixo `Post "<url>": `.
- **VERIFY**: contrato colado em Notes e revisado antes da T2.
- **ROLLBACK**: n/a.

### Task T2: `GET /api/info`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**: `commands.version/commit/date` existem ([root.go](cmd/pingou/commands/root.go)).
- **INPUT**: nada expõe versão nem tamanho do banco por HTTP.
- **OUTPUT**:
  1. `config.Config` ganha `Version`, `Commit` e `BuildDate`, preenchidos em [serve.go](cmd/pingou/commands/serve.go) antes de `NewServer`.
  2. Novo `internal/handler/info.go` com `handleInfo` e `databaseSize(dsn)`:
     - tira `?query` e o prefixo `file:` do DSN;
     - aplica `os.Stat` no arquivo principal (se falhar → `500`);
     - soma `-wal` e `-shm` quando existirem (`os.IsNotExist` → ignora; outro erro → `500`).
  3. Rota `GET /info` dentro do mux protegido.
  4. Nenhum SQL: é só `os.Stat`, então não precisa de service/repository (a regra de camadas vale para acesso a banco).
- **VERIFY**: `curl -H "X-API-Key: $K" localhost:8080/api/info`. O valor bate com `stat -c%s pingou.db pingou.db-wal pingou.db-shm | paste -sd+ | bc`.
- **ROLLBACK**: remover a rota e o arquivo.

### Task T3: Nome do arquivo no export

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: None
- **INPUT**: `pingou-backup-20261029-153000.db`.
- **OUTPUT**: `"pingou-" + time.Now().UTC().Format("2006-01-02") + ".db"` em [export.go](internal/handler/export.go).
- **VERIFY**: `curl -sD - -o /dev/null -H "X-API-Key: $K" localhost:8080/api/export | grep -i disposition`.
- **ROLLBACK**: reverter a linha.

### Task T4: Validar esquema da `webhook_url`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: None
- **INPUT**: `Update` usa `url.ParseRequestURI`, que aceita `ftp://x` e `mailto:`.
- **OUTPUT**: `validateWebhookURL(s string) error` em [settings_service.go](internal/service/settings_service.go). Exige esquema `http`/`https` e host não vazio, e devolve `ErrValidation`. É usada pelo `Update` e pela T5 (uma regra só).
- **VERIFY**: `PATCH {"webhook_url":"ftp://x"}` → `400`; `""` continua aceito (desliga o webhook).
- **ROLLBACK**: reverter a função. Uma URL `ftp://` já salva nunca funcionou, então ninguém perde nada.

### Task T5: `TestWebhook`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1, T4
- **PRE-CONDITIONS**: `webhookPayload` e `toWebhookMonitor` em [notifier.go](internal/service/notifier.go).
- **INPUT**: `send()` não devolve nada e usa um cliente de 10 s que segue redirects.
- **OUTPUT**:
  1. `type WebhookTestResult struct { OK bool; StatusCode *int; DurationMS int64; Error *string }` com as tags JSON do contrato.
  2. Em `notifier.go`, `testClient` próprio no nível do pacote: `Timeout: 5s` e `CheckRedirect: func(...) error { return http.ErrUseLastResponse }`.
  3. `func SendTestWebhook(ctx, rawURL) WebhookTestResult`:
     - payload com `Event: "test"`, monitor `{ID: "00000000-0000-0000-0000-000000000000", Name: "Teste do Pingou", Type: "http", Target/URL: "https://exemplo.com"}`, `Timestamp: now`, `LastError: nil` e `DowntimeDurationSeconds: nil`;
     - `io.Copy(io.Discard, io.LimitReader(resp.Body, 4<<10))`;
     - `slog.Info("webhook test", "host", u.Host, "status", ...)`, só com o host.
  4. `SettingsService.TestWebhook(ctx, url *string) (*WebhookTestResult, error)` com esta ordem: URL do corpo → URL salva → `ErrValidation` se as duas estiverem vazias → `validateWebhookURL` → `SendTestWebhook`.
  5. `NotifyDown`/`NotifyRecovery` **não mudam**. Tocar no caminho real de notificação está fora do escopo.
- **VERIFY**: testes da T7 com `httptest` (2xx, 404, 302 sem seguir, timeout, corpo grande).
- **ROLLBACK**: remover as funções novas; nada existente foi alterado.

### Task T6: Rota + rate limit

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T5
- **INPUT**: só o `/ping` tem limiter.
- **OUTPUT**:
  1. `Server.webhookTestLimiter = newRateLimiter(5, time.Minute)`.
  2. `api.Handle("POST /settings/webhook/test", s.webhookTestLimiter.middleware(http.HandlerFunc(s.handleTestWebhook)))`.
  3. `handleTestWebhook`: um corpo vazio (`io.EOF`) é válido; JSON inválido → `400`; erro do service → `writeServiceError`; sucesso → `200`.
- **VERIFY**: `curl` nos três cenários do briefing (200 ok, 404, timeout), mais `400` sem URL, `400` com `ftp://` e a 6ª chamada em 1 min → `429`. Respostas coladas em Notes.
- **ROLLBACK**: remover a rota.

### Task T7: Testes do backend

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2, T3, T6
- **OUTPUT**: casos table-driven em [handler_test.go](internal/handler/handler_test.go) e [notifier_test.go](internal/service/notifier_test.go):
  - `/api/info`: tamanho = soma dos arquivos (banco em `t.TempDir()`, criando um `-wal` falso); sem `-shm` não dá erro; sem o arquivo principal → `500`.
  - export: header com a data UTC.
  - teste do webhook: 200 → `ok:true`; 404 → `ok:false, status_code:404`; 302 → `status_code:302`, sem seguir o redirect; servidor lento → `status_code:null`, `error:"timeout após 5s"` (o timeout entra como variável do pacote para o teste rodar rápido); o payload recebido tem `event:"test"`; `ftp://` → `400`; sem URL salva e sem corpo → `400`; 6ª chamada → `429`.
  - log: capturar o `slog` e conferir que um token na query não aparece (mesmo padrão de `TestPingTokenNotLogged`).
- **VERIFY**: `go test ./internal/...` verde.
- **ROLLBACK**: n/a.

### Task T8: API, hook e formatBytes

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T1
- **OUTPUT**:
  1. [settings.js](web/src/api/settings.js): `info()`, `testWebhook(url)` e `export()` (blob). O download sai de `Settings.jsx`.
  2. `useSettings` busca `settings` e `info` em paralelo. Se `info` falhar, a tela continua (tamanho e versão somem).
  3. `formatBytes` em [format.js](web/src/lib/format.js): `Intl.NumberFormat("pt-BR", {maximumFractionDigits: 1})` → "12,4 MB".
  4. `backupFilename()` = `pingou-${new Date().toISOString().slice(0,10)}.db` (UTC, igual ao backend).
- **VERIFY**: `npm run lint`; os valores aparecem no React DevTools.
- **ROLLBACK**: reverter os arquivos.

### Task T9: Cards

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T8
- **OUTPUT** (em `components/settings/`, textos do PNG):
  1. **WebhookCard**:
     - título, descrição e o link "Ver formato do JSON" → `https://github.com/Bradrickcruz/pingou#webhook` (`target="_blank" rel="noreferrer"`);
     - `input type="url"` em fonte mono e "Enviar teste" com o ícone `Send`. Fica desabilitado com o campo vazio ou durante o teste (spinner);
     - resultado inline com `role="status"`: ✓ "Entregue · 200 · 142 ms" em `text-pg-up`, ou ✗ "Falhou · 404 · unexpected status code: 404" em `text-pg-down`. O resultado some quando a URL muda;
     - erro `429` → "Muitos testes seguidos. Tente de novo em 1 minuto.";
     - ajuda: "Deixe vazio para não enviar notificações."
  2. **RetentionCard**:
     - `fieldset` com rádios `sr-only` em chips arredondados (classes do `ChipGroup`);
     - cópia corrigida (ver Overview);
     - "Tamanho atual do banco: 12,4 MB" (ou nada, se `info` falhar).
  3. **BackupCard**:
     - caixa `bg-pg-ciano-nevoa` com o ícone `Lock` e o aviso do PNG;
     - botão "Baixar backup (.db)" com `Download` e um estado "Gerando…" durante o download (o `VACUUM INTO` pode demorar);
     - ao lado, em mono, `{backupFilename()} · ≈ {formatBytes(size)}`;
     - erro de download inline (`role="alert"`).
- **VERIFY**: comparação lado a lado com o PNG em 1440 px; 375 px sem rolagem horizontal (o input e o botão de teste empilham).
- **ROLLBACK**: reverter os arquivos.

### Task T10: Página, dirty state e rodapé

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T9
- **OUTPUT**:
  1. [Settings.jsx](web/src/pages/Settings.jsx):
     - `max-w-[760px]`, título "Configurações" e o subtítulo "Valem para todos os monitores.";
     - um `<form>` envolve Notificações e Retenção. Backup fica fora (não é salvo).
  2. `dirty = form.webhook_url !== settings.webhook_url || form.retention_days !== settings.retention_days`.
  3. "Salvar alterações" alinhado à direita:
     - `disabled={!dirty || saving}`;
     - depois de salvar, ✓ "Salvo" por 2 s e o formulário volta a ficar limpo;
     - erro de validação da API inline (`role="alert"`).
  4. Rodapé com `border-t`:
     - `pingou {version} · commit {commit} · Apache 2.0`, sem o commit quando ele for `none`;
     - links "Código no GitHub" (`https://github.com/Bradrickcruz/pingou`) e "Relatar um problema" (`…/issues/new`).
  5. Estados: carregando (Spinner, como hoje) e falha ao carregar settings (mensagem + "Tentar de novo").
- **VERIFY**: `npm run lint && npm run build`; o botão fica cinza sem mudanças, ativa ao digitar e volta a cinza ao desfazer a mudança.
- **ROLLBACK**: `git revert` do commit da Fase 2. A API nova é compatível com a tela antiga.

### Task T11: Revisão de segurança

- **Agent**: security-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T6
- **OUTPUT**: checklist em Notes com estes pontos:
  - a rota exige API Key;
  - redirects estão desligados;
  - timeout de 5 s;
  - leitura limitada a 4 KB;
  - o corpo do destino não é devolvido;
  - só o host vai para o log;
  - rate limit ativo.

  **SSRF para a rede interna (localhost, 169.254.x)** fica aceito e documentado: quem tem a API Key já pode salvar qualquer URL de webhook, e o envio real já chama essa URL. Bloquear IPs privados quebraria o uso comum (n8n na mesma rede).
- **VERIFY**: nenhum achado bloqueante.
- **ROLLBACK**: n/a.

### Task T12: QA manual

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T10
- **OUTPUT**: no binário real (`make build && ./bin/pingou serve`), conferir:
  - teste contra o webhook.site (2xx), contra uma URL 404 e contra um IP que não responde (timeout de 5 s);
  - testar sem salvar, depois recarregar a página: a URL antiga continua lá;
  - apagar a URL e salvar → as notificações param;
  - o nome do backup baixado é igual ao exibido;
  - navegação só por teclado (chips com as setas, foco visível nos botões e links);
  - 375 px.
- **VERIFY**: achados registrados em Notes; nenhum bloqueante aberto.
- **ROLLBACK**: n/a.

### Task T13: README

- **Agent**: document
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T6
- **OUTPUT**:
  - na seção Webhook, o evento `"test"` com um exemplo e o aviso: "integrações que tratam `event` devem ignorar valores desconhecidos";
  - na lista de endpoints, `GET /api/info` e `POST /api/settings/webhook/test`;
  - na retenção, deixar claro que ela roda na subida e a cada hora.
- **VERIFY**: leitura; os `curl` do README funcionam.
- **ROLLBACK**: reverter o arquivo.

## Risk Areas

| Risk | Impact | Mitigation |
|---|---|---|
| Teste sem redirect falha onde o envio real "funciona" (o real segue redirect) | Usuário vê falha num 301 http→https | Uma falha honesta: o `http.Client` transforma POST em GET num 301/302, então o envio real também não entrega. A mensagem mostra `status_code: 301` |
| DSN com `file:` ou `?query` | `os.Stat` falha → `500` em `/api/info` | `databaseSize` limpa o DSN; teste cobre os dois formatos |
| Tamanho exibido ≠ tamanho do arquivo baixado | Pequena confusão | Rótulo "≈" no backup |
| Rate limit por IP atrás de proxy | Todos os usuários dividem 5/min | Aceito: painel de um admin só; o mesmo `ponytail` do `/ping` já documenta isso |
| `build_date` "unknown" em `go run` | Rodapé estranho em dev | Rodapé só mostra versão e commit; commit some se for `none` |
| Teste com timeout prende a request por 5 s | `WriteTimeout` do servidor é 30 s | Cabe com folga |

## Rollback Strategy

1. Frontend: reverter o commit da Fase 2. A tela antiga continua funcionando com a API nova (só houve acréscimos).
2. Backend: reverter o commit da Fase 1. Sem migration, então o banco não muda.
3. A validação http/https (T4) é a única mudança de comportamento numa rota existente. Se quebrar alguém, reverter só ela.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test ./...

# Frontend
cd web && npm run lint && npm run build && cd ..

# Binário completo
make build && ./bin/pingou serve

# Rotas novas
curl -s -H "X-API-Key: $K" localhost:8080/api/info
curl -s -X POST -H "X-API-Key: $K" -d '{"url":"https://webhook.site/<id>"}' localhost:8080/api/settings/webhook/test
curl -sD - -o /dev/null -H "X-API-Key: $K" localhost:8080/api/export | grep -i disposition

# Grafo
graphify update .
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 0: Contrato | T1 | XS | ~5 min | Contrato fechado | ⬜ Pendente |
| Fase 1: API | T2, T3, T4, T5, T6, T7 | M | ~1 h 15 | Rotas novas com testes e `curl` | ⬜ Pendente |
| Fase 2: Tela | T8, T9, T10 | M | ~1 h 15 | Página igual ao PNG | ⬜ Pendente |
| Fase 3: Polimento | T11, T12, T13 | S | ~40 min | Segurança, QA e README | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T13

## Notes

- **Resposta à pergunta em aberto do briefing:** a retenção apaga **checks e incidentes resolvidos** mais antigos que `retention_days` (7–90, padrão 30). Incidentes abertos nunca expiram. O worker roda na subida e a cada 1 h ([retention.go](internal/scheduler/retention.go)).
- **Conflito com a Entrega 2 do briefing:** o painel já recebe `history` / `uptime_24h` / `last_latency_ms`. Não crio `recent_checks` / `last_response_time_ms`. Se o novo front precisar de `status`/`checked_at` por check, isso vira um acréscimo em outro plano.
- `/api/session` (Entrega 1.1) pertence ao plano do login.

### Alternativas consideradas e descartadas

- `golang.org/x/time/rate` para o rate limit → descartada: o `rateLimiter` do projeto já resolve, e seria uma dependência nova.
- Bloquear IPs privados no teste do webhook → descartada: quebra o n8n na rede local, e o envio real já chama qualquer URL salva.
- Medir o tamanho real do backup com um `VACUUM INTO` prévio → descartada: gera o dump inteiro só para mostrar um número.
- `BuildInfo` como parâmetro novo de `NewServer` → descartada: três campos em `Config` não mexem nos chamadores.
- Reaproveitar `send()` do notifier para o teste → descartada: exigiria mudar timeout e redirect do envio real.
- Mudar a retenção para "uma vez por dia" e bater com o PNG → descartada: o briefing proíbe mexer na retenção. A cópia é que muda.
