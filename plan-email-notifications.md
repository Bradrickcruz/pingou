# Plano de Desenvolvimento: Notificação por e-mail (SMTP) e envio fora do lock

## Overview

Hoje o Pingou avisa só por webhook. Este plano adiciona o **e-mail (SMTP)** como segundo canal. Os dois podem ficar ligados ao mesmo tempo, e cada evento (`down`/`up`) vai para todos os canais ativos.

Antes disso, o plano corrige um problema que já existe: **o envio da notificação trava todos os monitores**. O `notify()` roda dentro de `StateMachine.Process` ([state_machine.go:74](internal/service/state_machine.go#L74)), que roda dentro de `withMonitor` com o lock global `s.dbMu` seguro ([scheduler.go:196](internal/scheduler/scheduler.go#L196)). Um webhook lento (timeout de 10 s) pausa a gravação de todos os monitores e heartbeats. Com SMTP, que pode levar 30 s, isso piora.

### Decisões tomadas (conversa de 2026-11-01)

| Tema                               | Decisão                                                                                                            |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Quantos canais                     | **Um de cada tipo**: 1 webhook + 1 configuração de e-mail (com vários destinatários)                               |
| Onde fica a config do SMTP         | **Toda no ambiente**: `PINGOU_SMTP_*`, inclusive a senha. Nada do servidor SMTP vai para o banco nem para o export |
| O que a tela configura             | Só os **destinatários**, separados por vírgula, guardados em `settings` com a chave `email_recipients`             |
| Quando o e-mail aparece habilitado | Só com o SMTP configurado no ambiente. Senão, o card aparece **bloqueado** e explica quais variáveis faltam        |
| Biblioteca                         | **stdlib** (`net/smtp`, `crypto/tls`, `net/mail`, `mime`, `mime/quotedprintable`). Nenhuma dependência nova        |
| Formato do e-mail                  | **Texto simples** (`text/plain; charset=utf-8`)                                                                    |
| Correção do lock                   | **Primeira fase** e commit próprio. Pode ir para `develop` antes do SMTP                                           |

### Reaproveitado da análise (não criar de novo)

- A interface `Notifier` ([notifier.go:18](internal/service/notifier.go#L18)) **já é** o contrato de canal, e `Monitor` + `Incident` **já são** o evento neutro. Não haverá `NotificationEvent`, `Channel` nem `Dispatcher` com repositório.
- O fan-out é um `MultiNotifier`, que também implementa `Notifier`. `StateMachine` e `Scheduler` não mudam.
- Ler a config na hora do envio, como já faz o `getWebhookURL`, também serve para os destinatários do e-mail.
- O teste do e-mail copia o padrão do teste do webhook: `POST /settings/webhook/test`, `rateLimiter`, `WebhookTestResult` e `WebhookCard`.
- O link do dashboard reaproveita `PINGOU_PUBLIC_URL` / `dashboardURL(m)`.
- **Não precisa de migration.** `SettingsRepo.Get` devolve `""` quando a chave não existe, e `Set` já faz upsert.

### Limites (fora do escopo)

Vários webhooks ou vários servidores SMTP; canal por monitor; e-mail em HTML; outbox e retry; Slack, Discord ou Teams nativos; templates personalizados.

### Restrições

- Só adicionar, não quebrar: `GET/PATCH /api/settings` ganham campos novos, e `webhook_url` continua igual.
- Camadas handler → service → repository; `fmt.Errorf("...: %w")`; `slog` sem senha, sem lista completa de destinatários e sem URL de webhook.
- O e-mail nunca leva o `ping_token`, seguindo a mesma regra do `webhookPayloadMonitor`.
- Branch `feat/email-notifications` a partir de `develop`. Um commit por fase. A Fase 0 pode virar um PR separado.

## Project Type

**FULLSTACK**: Go + SQLite (`internal/`, `cmd/`) e React/Vite/Tailwind (`web/`), com o front embutido via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Um notifier lento (ex.: 10 s) não atrasa o `Process` de outros monitores. Há teste automatizado
- [ ] A falha ou o panic de um canal não impede o envio pelo outro
- [ ] No SIGTERM, envios em andamento têm até 10 s para terminar antes do processo sair
- [ ] Sem `PINGOU_SMTP_HOST`: `email_available: false`, o card de e-mail fica bloqueado e um `PATCH` com destinatários devolve `400`
- [ ] Com `PINGOU_SMTP_HOST` setado e a config incompleta ou inválida: o `serve` não sobe e mostra uma mensagem clara (ex.: `PINGOU_SMTP_FROM is required when PINGOU_SMTP_HOST is set`)
- [ ] Destinatários são validados com `net/mail` (no máximo 10), normalizados e salvos em `email_recipients`
- [ ] Queda e volta de um monitor geram um e-mail em texto simples, com assunto legível (acentos e emoji via `mime.QEncoding`) e link para o dashboard quando há `PINGOU_PUBLIC_URL`
- [ ] `POST /api/settings/email/test` envia um e-mail de teste e devolve `{ok, duration_ms, error}`, com rate limit de 5/min
- [ ] Os modos `starttls` (587), `tls` (465) e `none` (relay local, sem autenticação) funcionam
- [ ] Senha, `ping_token` e lista de destinatários não aparecem em log, resposta da API, `pingou config` ou export
- [ ] README, `.env.example` e `pingou config` documentam ou mostram as variáveis `PINGOU_SMTP_*`

## Tech Stack

| Technology                                                              | Choice Rationale                                                                                                                                               |
| ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `net/smtp` + `crypto/tls`                                               | Cobre STARTTLS e TLS implícito em cerca de 30 linhas. O usuário pediu para evitar dependências. `PlainAuth` já recusa autenticar sem TLS (exceto em localhost) |
| `net/mail`                                                              | Valida destinatários (`ParseAddress`) e o `From`, e gera o `Message-ID`/`Date` corretos                                                                        |
| `mime.QEncoding` + `mime/quotedprintable`                               | Assunto e corpo em UTF-8 sem biblioteca externa                                                                                                                |
| `context.WithoutCancel` (Go 1.21+)                                      | O envio assíncrono não pode morrer quando o `ctx` do job é cancelado (edição ou pausa do monitor)                                                              |
| `sync.WaitGroup`                                                        | Espera os envios no shutdown sem fila nem worker                                                                                                               |
| `rateLimiter` existente ([ratelimit.go](internal/handler/ratelimit.go)) | Mesmo limite do teste do webhook, sem dependência nova                                                                                                         |

## File Structure

```
cmd/pingou/commands/
├── serve.go                 # MODIFICAR: monta MultiNotifier (webhook + e-mail se configurado); Wait no shutdown
└── info.go                  # MODIFICAR: `pingou config` mostra SMTP (sem senha, só "password_set")
internal/config/
├── config.go                # MODIFICAR: struct SMTPConfig + loadSMTP() com validação
└── config_test.go           # NOVO: casos de SMTP (vazio, completo, incompleto, porta/tls inválidos)
internal/service/
├── notifier.go              # MODIFICAR: MultiNotifier + helper downtimeSeconds(i) compartilhado
├── notifier_test.go         # MODIFICAR: testes do MultiNotifier (lento, panic, Wait)
├── email_notifier.go        # NOVO: EmailNotifier, buildEmail(), sendMail(), SendTestEmail()
├── email_notifier_test.go   # NOVO: formato da mensagem + envio contra um servidor SMTP falso em memória
├── settings_service.go      # MODIFICAR: email_recipients, email_available, TestEmail()
└── validate.go              # MODIFICAR: parseRecipients()
internal/handler/
├── settings.go              # MODIFICAR: campo novo no PATCH + handleTestEmail
├── server.go                # MODIFICAR: rota POST /settings/email/test + emailTestLimiter
├── settings_test.go         # MODIFICAR: GET/PATCH com e-mail, rota de teste, 429
└── handler_test.go          # MODIFICAR: newTestServer com a nova assinatura de NewSettingsService
web/src/
├── api/settings.js          # MODIFICAR: testEmail(recipients)
├── components/settings/
│   ├── EmailCard.jsx        # NOVO: input de destinatários, estado bloqueado, "Enviar teste"
│   └── WebhookCard.jsx      # MODIFICAR: título "Webhook" (o grupo vira "Notificações")
└── pages/Settings.jsx       # MODIFICAR: EmailCard, dirty e payload com email_recipients
.env.example                 # MODIFICAR: bloco PINGOU_SMTP_* comentado
README.md                    # MODIFICAR: variáveis + seção "E-mail" + endpoint de teste
```

## Contrato

### Variáveis de ambiente

| Variável               | Obrigatória           | Padrão                                     | Regra                                                                    |
| ---------------------- | --------------------- | ------------------------------------------ | ------------------------------------------------------------------------ |
| `PINGOU_SMTP_HOST`     | liga o e-mail         | (vazio)                                    | Vazio = e-mail indisponível. As outras variáveis são ignoradas           |
| `PINGOU_SMTP_PORT`     | não                   | `587`                                      | 1–65535                                                                  |
| `PINGOU_SMTP_USERNAME` | não                   | (vazio)                                    | Vazio = sem autenticação (relay)                                         |
| `PINGOU_SMTP_PASSWORD` | se houver `USERNAME`  | (vazio)                                    | Nunca aparece em log, API ou `pingou config`                             |
| `PINGOU_SMTP_FROM`     | sim, se houver `HOST` | -                                          | Validado com `net/mail.ParseAddress`. Ex.: `Pingou <alerta@exemplo.com>` |
| `PINGOU_SMTP_TLS`      | não                   | `tls` se a porta for 465, senão `starttls` | `starttls` \| `tls` \| `none`. `none` só sem `USERNAME`                  |

Com `HOST` setado e qualquer regra violada, `config.Load()` devolve erro e o `serve` não sobe. É o mesmo comportamento de quando falta `PINGOU_API_KEY`.

### API

```jsonc
// GET /api/settings: campos novos
{ "webhook_url": "", "retention_days": 30,
  "email_recipients": "ops@exemplo.com,bryan@exemplo.com",  // "" = e-mail desligado
  "email_available": true }                                 // false = SMTP não configurado no ambiente

// PATCH /api/settings
{ "email_recipients": "ops@exemplo.com, bryan@exemplo.com" }
// 200 → salvo normalizado ("ops@exemplo.com,bryan@exemplo.com")
// 400 VALIDATION → endereço inválido, mais de 10, ou não vazio com email_available=false
// "" sempre aceito (desliga o e-mail)

// POST /api/settings/email/test   (corpo opcional; sem corpo usa os destinatários salvos)
{ "recipients": "ops@exemplo.com" }
// 200 → { "ok": true, "duration_ms": 812, "error": null }
//       { "ok": false, "duration_ms": 30001, "error": "timeout após 30s" }
// 400 → SMTP indisponível, sem destinatários, endereço inválido
// 429 → mais de 5 testes por minuto
```

O e-mail é enviado quando `email_available` é verdadeiro **e** `email_recipients` não está vazio. Se a env for removida depois, os destinatários continuam salvos e o envio simplesmente não acontece. O card mostra o valor salvo, bloqueado.

### Mensagem

```
Subject: [Pingou] 🔴 API de pagamentos está fora do ar
Subject: [Pingou] 🟢 API de pagamentos voltou (fora por 5 min)
Subject: [Pingou] Teste de e-mail

Monitor: API de pagamentos
Tipo: http
Alvo: https://api.exemplo.com/health     (omitido no heartbeat)
Caiu em: 2026-11-01 14:03:12 UTC
Erro: connection refused                  (só no down)
Fora do ar por: 5 min 12 s                (só no up)
Detalhes: https://pingou.exemplo.com/monitors/<id>   (só com PINGOU_PUBLIC_URL)
```

## Task Breakdown

| Task ID | Agent                 | Name                                                             | Priority | Complexity | Dependencies |
| ------- | --------------------- | ---------------------------------------------------------------- | -------- | ---------- | ------------ |
| T1      | go-backend-specialist | `MultiNotifier` assíncrono                                       | P0       | S          | -            |
| T2      | code-test             | Testes do `MultiNotifier`                                        | P0       | S          | T1           |
| T3      | go-backend-specialist | Ligar `MultiNotifier` no `serve` + Wait no shutdown              | P0       | XS         | T1           |
| T4      | go-backend-specialist | `SMTPConfig` no `config` + `pingou config`                       | P0       | S          | -            |
| T5      | go-backend-specialist | `buildEmail` (mensagem) + helper de duração                      | P1       | S          | T4           |
| T6      | go-backend-specialist | `sendMail` (stdlib, 3 modos de TLS) + `EmailNotifier`            | P1       | M          | T5           |
| T7      | code-test             | Testes do e-mail (mensagem + SMTP falso)                         | P1       | M          | T6           |
| T8      | go-backend-specialist | `SettingsService`: destinatários, `email_available`, `TestEmail` | P1       | S          | T6           |
| T9      | go-backend-specialist | Handler, rota de teste e montagem no `serve`                     | P1       | S          | T3, T8       |
| T10     | code-test             | Testes de handler                                                | P1       | S          | T9           |
| T11     | frontend-specialist   | `settingsApi.testEmail` + `EmailCard`                            | P2       | M          | T9           |
| T12     | frontend-specialist   | `Settings.jsx`: dirty, salvar e títulos                          | P2       | S          | T11          |
| T13     | security-review       | Revisão de segurança                                             | P3       | S          | T10, T12     |
| T14     | qa-review             | QA manual (Mailpit + provedor real)                              | P3       | M          | T12          |
| T15     | document              | README + `.env.example`                                          | P3       | S          | T9           |

### Task T1: `MultiNotifier` assíncrono

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Branch `feat/email-notifications` criada a partir de `develop`
- **INPUT**: `notify()` roda de forma síncrona dentro do lock global `dbMu`. Um destino lento trava todos os monitores.
- **OUTPUT**:
  1. Em [notifier.go](internal/service/notifier.go): `type MultiNotifier struct { notifiers []Notifier; wg sync.WaitGroup }`, que implementa `Notifier`
  2. `NotifyDown`/`NotifyRecovery`: para cada notifier, uma `go func` com `ctx := context.WithoutCancel(ctx)`, `defer wg.Done()` e `defer recover()` com `slog.Error("notifier panic", ...)`
  3. `Wait(timeout time.Duration) bool`: espera o `wg` ou desiste depois do timeout
  4. Comentário explicando o porquê: o envio sai do lock do scheduler; `m` e `i` são seguros entre goroutines porque `withMonitor` recarrega o monitor a cada chamada e ninguém mais escreve neles
- **VERIFY**:
  - `go build ./... && go vet ./...`
- **ROLLBACK**: reverter o commit da Fase 0. O `serve` volta a passar o `WebhookNotifier` direto.

### Task T2: Testes do `MultiNotifier`

- **Agent**: code-test
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **INPUT**: `MultiNotifier` sem testes.
- **OUTPUT**:
  1. `TestMultiNotifierDoesNotBlock`: um notifier que dorme 2 s, `NotifyDown` volta em menos de 50 ms
  2. `TestMultiNotifierIsolatesFailures`: um notifier com panic e outro que conta chamadas; o segundo recebe 1 chamada
  3. `TestMultiNotifierWait`: `Wait(1s)` é `true` com envio rápido e `false` com envio preso
  4. Em [scheduler_test.go](internal/scheduler/scheduler_test.go): um teste com um `MultiNotifier` lento no `Process` do monitor A, e o `Process` do monitor B termina sem esperar por ele
- **VERIFY**:
  - `go test -race ./internal/service/... ./internal/scheduler/...`
- **ROLLBACK**: apagar os testes junto com o T1.

### Task T3: Ligar no `serve` + Wait no shutdown

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T1
- **INPUT**: [serve.go:96](cmd/pingou/commands/serve.go#L96) passa o `webhookNotifier` direto para o scheduler.
- **OUTPUT**:
  1. `notifier := service.NewMultiNotifier(webhookNotifier)`, passado para o `NewScheduler`
  2. Depois de `monitorScheduler.Stop()`: `if !notifier.Wait(10 * time.Second) { slog.Warn("notifications still pending at shutdown") }`
- **VERIFY**:
  - `make build && ./bin/pingou`: derrubar um monitor apontado para um webhook lento (`httpbin /delay/8`). Os outros monitores continuam gravando checks no meio tempo (`GET /api/monitors/:id/checks`)
  - **Commit da Fase 0** (opcional: PR separado para `develop`)
- **ROLLBACK**: voltar a passar o `webhookNotifier` direto.

### Task T4: `SMTPConfig` no `config` + `pingou config`

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **INPUT**: `Config` não conhece SMTP. O `pingou config` ([info.go:40](cmd/pingou/commands/info.go#L40)) imprime um mapa com uma lista explícita de campos.
- **OUTPUT**:
  1. `type SMTPConfig struct { Host string; Port int; Username, Password, From, TLS string }` e `Config.SMTP SMTPConfig`
  2. `(SMTPConfig) Enabled() bool { return c.Host != "" }`
  3. `loadSMTP() (SMTPConfig, error)` com as regras da tabela do contrato (porta padrão, TLS deduzido da porta, `FROM` validado com `net/mail`, `USERNAME` sem `PASSWORD` = erro, `none` com `USERNAME` = erro)
  4. `pingou config` passa a mostrar `SMTPHost`, `SMTPPort`, `SMTPFrom`, `SMTPTLS` e `SMTPPasswordSet` (bool). Nunca a senha
  5. `config_test.go` com tabela: vazio, completo, 465 → `tls`, sem `FROM`, porta `abc`, `TLS=foo`, usuário sem senha, `none` com usuário
- **VERIFY**:
  - `go test ./internal/config/...`
  - `PINGOU_SMTP_HOST=x ./bin/pingou` → sai com `config error: PINGOU_SMTP_FROM is required when PINGOU_SMTP_HOST is set`
- **ROLLBACK**: remover o campo; ninguém mais depende dele até o T6.

### Task T5: `buildEmail` + helper de duração

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **INPUT**: o cálculo da duração da queda está duplicado inline em `NotifyRecovery` ([notifier.go:100-112](internal/service/notifier.go#L100-L112)).
- **OUTPUT**:
  1. `downtimeSeconds(i *domain.Incident, end time.Time) int64` extraído e usado pelo webhook (sem mudar o comportamento)
  2. Em `email_notifier.go`: a função pura `buildEmail(from string, to []string, subject, body string, now time.Time) []byte`, com os cabeçalhos `From`, `To`, `Subject` (`mime.QEncoding.Encode("utf-8", ...)`), `Date` (RFC 5322), `Message-ID`, `MIME-Version`, `Content-Type: text/plain; charset=utf-8` e `Content-Transfer-Encoding: quoted-printable`; linhas com CRLF
  3. `subjectDown/subjectUp/bodyFor(m, i, dashboardURL)` no formato da seção "Mensagem". CR/LF removidos do nome do monitor e do erro antes de montar o assunto
- **VERIFY**:
  - `go test ./internal/service/...` (os testes do webhook continuam verdes depois da extração)
- **ROLLBACK**: reverter o arquivo. O webhook volta ao cálculo inline.

### Task T6: `sendMail` + `EmailNotifier`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T5
- **INPUT**: mensagem pronta, mas sem transporte.
- **OUTPUT**:
  1. `sendMail(ctx, cfg config.SMTPConfig, to []string, msg []byte) error` usando a stdlib:
     - `tls`: `tls.DialWithDialer(&net.Dialer{Timeout: 10s}, ...)` + `smtp.NewClient`
     - `starttls`/`none`: `net.Dialer.DialContext` + `smtp.NewClient`. Em `starttls`, `c.StartTLS(&tls.Config{ServerName: host})` é **obrigatório** (erro se o servidor não anunciar)
     - `conn.SetDeadline(now + 30s)` para cobrir a conversa inteira, já que `net/smtp` não aceita `context`
     - `smtp.PlainAuth` só com `Username`; depois `Mail`, `Rcpt` para cada destinatário, `Data`, `Quit`
  2. `EmailNotifier{cfg, getRecipients func() string, publicURL}` implementando `Notifier`. Com destinatários vazios, não faz nada. O log diz `email sent` com `event` e `recipients_count` (sem endereços); o erro é logado sem a senha
  3. `SendTestEmail(ctx, cfg, to []string) EmailTestResult{OK, DurationMS, Error}`, no mesmo molde do `SendTestWebhook`
- **VERIFY**:
  - `go build ./... && go vet ./...`
- **ROLLBACK**: apagar `email_notifier.go`. O T3 continua funcionando só com o webhook.

### Task T7: Testes do e-mail

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T6
- **INPUT**: e-mail sem testes.
- **OUTPUT**:
  1. `TestBuildEmail`: cabeçalhos presentes; assunto com acento e emoji codificado; corpo sem `ping_token`; um nome de monitor com `\r\nBcc: x@y` não gera cabeçalho novo
  2. `TestBodyFor`: down com erro, up com duração, heartbeat sem "Alvo", e "Detalhes" só com `publicURL`
  3. Um servidor SMTP falso em memória (`net.Listen` em 127.0.0.1, cerca de 40 linhas respondendo `220/250/354/221`) para testar o modo `none`: o envio grava `MAIL FROM`, um `RCPT TO` por destinatário e o `DATA`
  4. `TestSendMailTimeout`: servidor que aceita a conexão e não responde; `sendMail` volta com erro dentro do deadline (deadline reduzido por variável de teste, como o `webhookTestTimeout`)
  5. `TestStartTLSRequired`: o servidor falso não anuncia `STARTTLS` e o modo `starttls` devolve erro
- **VERIFY**:
  - `go test -race ./internal/service/...`
- **ROLLBACK**: apagar o arquivo de teste.

### Task T8: `SettingsService` com e-mail

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T6
- **INPUT**: `Settings` só tem `webhook_url` e `retention_days`.
- **OUTPUT**:
  1. `NewSettingsService(repo, smtp config.SMTPConfig)`, com `smtp.Enabled()` decidindo o `email_available`
  2. `Settings` ganha `EmailRecipients string \`json:"email_recipients"\``e`EmailAvailable bool \`json:"email_available"\``. `UpdateSettingsInput`ganha`EmailRecipients \*string`
  3. `parseRecipients(s string) ([]string, error)` em [validate.go](internal/service/validate.go): separa por vírgula, faz `TrimSpace`, ignora vazios, valida com `mail.ParseAddress`, usa só o `.Address`, remove duplicados e aceita no máximo 10 (`ErrValidation`)
  4. `Update`: com valor não vazio e `!Enabled()`, devolve `ErrValidation` ("email is not configured on the server (PINGOU*SMTP*\*)"). Senão salva o valor normalizado (`strings.Join(list, ",")`)
  5. `TestEmail(ctx, raw *string)`: usa os destinatários informados ou os salvos, aplica as mesmas validações e chama `SendTestEmail`
  6. Testes de `parseRecipients` em `validate_test.go`
- **VERIFY**:
  - `go test ./internal/service/...`
- **ROLLBACK**: reverter o arquivo. A chave `email_recipients` órfã no banco não atrapalha nada.

### Task T9: Handler, rota e montagem no `serve`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T3, T8
- **INPUT**: rotas e montagem só conhecem o webhook.
- **OUTPUT**:
  1. [settings.go](internal/handler/settings.go): `EmailRecipients *string` no `updateSettingsRequest`; `handleTestEmail`, com corpo opcional `{recipients}`, como o `handleTestWebhook`
  2. [server.go](internal/handler/server.go): `emailTestLimiter: newRateLimiter(webhookTestRateLimit, time.Minute)` e `api.Handle("POST /settings/email/test", ...)`
  3. [serve.go](cmd/pingou/commands/serve.go): `NewSettingsService(settingsRepository, loadedConfig.SMTP)`. Se `loadedConfig.SMTP.Enabled()`, adiciona `service.NewEmailNotifier(cfg.SMTP, func() string { v, _ := settingsRepository.Get(ctx, "email_recipients"); return v }, cfg.PublicURL)` ao `MultiNotifier`. Log de subida: `email notifications available` com host e porta
- **VERIFY**:
  - `go build ./...`
  - `curl -H "X-API-Key: ..." localhost:8956/api/settings` mostra `email_available`
- **ROLLBACK**: reverter os três arquivos. Os T1–T8 continuam compilando.

### Task T10: Testes de handler

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T9
- **INPUT**: [settings_test.go](internal/handler/settings_test.go) cobre só o webhook.
- **OUTPUT**:
  1. Ajustar o `newTestServer` ([handler_test.go:50](internal/handler/handler_test.go#L50)) para a nova assinatura, com um `SMTPConfig` vazio por padrão e uma variante com SMTP apontando para o servidor falso do T7
  2. GET com `email_available` `false` e `true`; PATCH válido normalizado; PATCH inválido → `400`; PATCH não vazio sem SMTP → `400`; PATCH `""` → `200`
  3. `POST /settings/email/test`: `200` contra o servidor falso; sem destinatários → `400`; sem SMTP → `400`; 6ª chamada → `429`
- **VERIFY**:
  - `go test ./...`
  - **Commit da Fase 1**
- **ROLLBACK**: reverter os testes junto com o T9.

### Task T11: `testEmail` + `EmailCard`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T9
- **INPUT**: [WebhookCard.jsx](web/src/components/settings/WebhookCard.jsx) é o modelo (input, "Enviar teste", `resultText`, mensagem de `429`).
- **OUTPUT**:
  1. [settings.js](web/src/api/settings.js): `testEmail: (recipients) => client.post("/settings/email/test", recipients ? { recipients } : {})`
  2. `EmailCard({ value, onChange, available })`:
     - `available=false`: input `disabled` com o valor salvo visível, um aviso "Para ativar, configure `PINGOU_SMTP_HOST` e `PINGOU_SMTP_FROM` no servidor." e o link para `README#e-mail`; botão de teste escondido
     - `available=true`: input `type="text"` (não `email`, porque são vários), `inputMode="email"`, `autoComplete="off"`, placeholder `ops@exemplo.com, voce@exemplo.com`, dica "Separe por vírgula. Vazio desliga o e-mail."; "Enviar teste" com os destinatários digitados; resultado "Teste enviado em 812 ms." ou "Teste falhou. <erro>"
     - `aria-describedby` ligando o input à dica e ao resultado; o resultado em `role="status"`
- **VERIFY**:
  - `cd web && npm run lint && npm run build`
- **ROLLBACK**: apagar o componente e a função.

### Task T12: `Settings.jsx`

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T11
- **INPUT**: o `dirty` e o `update` em [Settings.jsx:55-76](web/src/pages/Settings.jsx#L55-L76) conhecem só `webhook_url` e `retention_days`.
- **OUTPUT**:
  1. `EmailCard` logo abaixo do `WebhookCard`, com `available={settings.email_available}`
  2. `dirty` passa a comparar também `email_recipients`; o `update` envia `email_recipients` só quando `email_available` (evita o `400` ao salvar outra coisa com o e-mail bloqueado)
  3. Depois de salvar, o valor exibido vem da resposta (normalizado pelo backend)
  4. Títulos: `WebhookCard` → "Webhook"; `EmailCard` → "E-mail". Um subtítulo "Notificações: os canais ligados recebem os mesmos eventos." acima dos dois
- **VERIFY**:
  - `cd web && npm run lint && npm run build && npm test`
  - `make build` e conferir a tela em 1440 e 375 px, com e sem `PINGOU_SMTP_HOST`
  - **Commit da Fase 2**
- **ROLLBACK**: reverter a página. A API continua aceitando o formato antigo.

### Task T13: Revisão de segurança

- **Agent**: security-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T10, T12
- **INPUT**: canal novo com credencial e saída de rede.
- **OUTPUT** (checklist registrado no Status):
  1. A senha não aparece em `slog`, `pingou config`, `GET /settings`, mensagens de erro nem no export (`grep` no dump do `VACUUM INTO`)
  2. Sem injeção de cabeçalho: nome do monitor e erro com CR/LF (teste do T7) e destinatários validados por `net/mail`
  3. `starttls` falha se o servidor não oferecer TLS, sem rebaixar a conexão para texto aberto
  4. `tls.Config` sem `InsecureSkipVerify`
  5. Rota de teste com `authMiddleware` e rate limit; limite de 10 destinatários (impede usar a rota para spam)
  6. Log sem endereços (só `recipients_count`)
- **VERIFY**:
  - Checklist preenchido no Status de Execução
- **ROLLBACK**: n/a (só revisão; os achados viram tasks novas).

### Task T14: QA manual

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T12
- **INPUT**: feature completa no binário.
- **OUTPUT**:
  1. **Mailpit** local (`docker run -p 1025:1025 -p 8025:8025 axllent/mailpit`) com `PINGOU_SMTP_TLS=none`: teste, queda e volta de um monitor; conferir assunto, acentos, emoji e link
  2. **Provedor real** em 587 (`starttls`), por exemplo Gmail com senha de app; e em 465 (`tls`), se houver um disponível
  3. Webhook e e-mail ligados juntos: os dois recebem; derrubar o SMTP (host errado) e o webhook continua saindo
  4. Sem env: card bloqueado; destinatário salvo antes continua visível; salvar retenção não dá `400`
  5. SIGTERM logo depois de uma queda: o e-mail ainda chega (Wait)
  6. Só teclado: tab até o input, digitar, "Enviar teste", ler o resultado
- **VERIFY**:
  - Evidências (logs e prints do Mailpit) registradas no Status de Execução
- **ROLLBACK**: n/a.

### Task T15: Documentação

- **Agent**: document
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T9
- **INPUT**: README sem e-mail; `.env.example` sem SMTP.
- **OUTPUT**:
  1. Tabela "Variáveis de ambiente" ([README.md:354](README.md#L354)) com as 6 `PINGOU_SMTP_*`
  2. Seção `## E-mail` depois de `## Webhook`: como ligar, exemplos (Gmail 587 com senha de app, relay local com `none`, Mailpit para testar), formato da mensagem e `POST /api/settings/email/test`
  3. Exemplo do `pingou config` atualizado
  4. Bloco comentado em `.env.example`; uma linha comentada em `docker-compose.yml`
  5. Nota no README: o Pingou continua avisando "uma vez por transição", sem retry (envio em memória)
- **VERIFY**:
  - Seguir o README do zero com o Mailpit funciona
  - **Commit da Fase 3**
- **ROLLBACK**: reverter os docs.

## Risk Areas

| Risk                                                                                            | Impact                                                                                     | Mitigation                                                                                                                           |
| ----------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| Goroutine de envio lê `settings` enquanto o scheduler segura a única conexão (`MaxOpenConns=1`) | Envio atrasado, sem deadlock (a goroutine só espera a conexão; o lock é solto logo depois) | Coberto pelo teste de regressão em `scheduler_test.go`, que imita leitura no banco (`dbNotifier`)                                    |
| Processo morre entre o commit e o envio                                                         | Aviso perdido                                                                              | `Wait(10s)` no shutdown cobre o SIGTERM. Crash ou `kill -9` continuam perdendo o aviso: aceito, e o outbox fica registrado nas Notes |
| `net/smtp` sem `context`                                                                        | Envio pendurado                                                                            | `conn.SetDeadline(30s)` + `Dialer.Timeout`, testado no T7                                                                            |
| Provedores exigem OAuth2 (Microsoft 365 em desativação de basic auth)                           | Usuário não consegue autenticar                                                            | Documentar: usar senha de app, relay ou um serviço transacional (SES, Brevo, Mailgun). OAuth2 fica fora do escopo                    |
| `starttls` obrigatório quebra relays internos sem TLS                                           | Erro na subida ou no teste                                                                 | Modo `none` explícito, só sem autenticação                                                                                           |
| E-mail caindo no spam                                                                           | Aviso não visto                                                                            | `Date`, `Message-ID` e `From` válidos; no README, recomendar um `FROM` do mesmo domínio autenticado (SPF/DKIM é com o usuário)       |
| Mudar a assinatura de `NewSettingsService`                                                      | Quebra o `newTestServer`                                                                   | Ajustado no T10 no mesmo commit                                                                                                      |
| Monitor oscilando gera muitos e-mails                                                           | Caixa lotada                                                                               | Já é limitado pelo `failure_threshold`. Rate limit de notificação fica no roadmap (Milestone 2.3)                                    |

## Rollback Strategy

1. Cada fase é um commit. `git revert <commit>` desfaz uma fase sem afetar as anteriores (a Fase 0 se sustenta sozinha).
2. Não há migration: a chave `email_recipients` órfã em `settings` não atrapalha, e um `DELETE FROM settings WHERE key='email_recipients'` limpa se necessário.
3. Em produção, desligar o e-mail sem deploy: remover `PINGOU_SMTP_HOST` e reiniciar. O card fica bloqueado e nada é enviado.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test -race ./...
gofumpt -l .            # sem saída

# Frontend
cd web && npm run lint && npm run format:check && npm test && npm run build && cd ..

# Binário completo
make build

# Fumaça com Mailpit
docker run -d --name mailpit -p 1025:1025 -p 8025:8025 axllent/mailpit
PINGOU_SMTP_HOST=localhost PINGOU_SMTP_PORT=1025 PINGOU_SMTP_TLS=none \
PINGOU_SMTP_FROM="Pingou <pingou@localhost>" ./bin/pingou &
curl -s -X PATCH -H "X-API-Key: $PINGOU_API_KEY" -d '{"email_recipients":"ops@exemplo.com"}' localhost:8956/api/settings
curl -s -X POST  -H "X-API-Key: $PINGOU_API_KEY" localhost:8956/api/settings/email/test
# conferir em http://localhost:8025
```

## Cronograma Sugerido

| Fase                       | Tarefas                     | Complexity | Esforço | Target                                                                    | Status       |
| -------------------------- | --------------------------- | ---------- | ------- | ------------------------------------------------------------------------- | ------------ |
| Fase 0: Envio fora do lock | T1, T2, T3                  | S          | ~35 min | Webhook lento não trava outros monitores (pode ir para `develop` sozinha) | ✅ Feito     |
| Fase 1: E-mail no backend  | T4, T5, T6, T7, T8, T9, T10 | L          | ~2 h 30 | `curl` envia o teste ao Mailpit; testes verdes                            | ✅ Feito     |
| Fase 2: Tela               | T11, T12                    | M          | ~45 min | Card bloqueado e desbloqueado funcionando                                 | ✅ Feito     |
| Fase 3: Polimento          | T13, T14, T15               | M          | ~1 h 15 | Segurança, QA com provedor real e README                                  | 🟡 Falta T14 |

## Status de Execução

### Concluído

- T1–T3 (Fase 0): `MultiNotifier` em `internal/service/multi_notifier.go` (arquivo próprio, em vez de `notifier.go`), com `Add()` para o `serve` incluir o e-mail. `Wait(10s)` no shutdown
- T4: `config.LoadSMTP()` exportada e chamada só no `Load()` do `serve`, para um SMTP inválido não travar o `migrate`. `pingou config` mostra `SMTPError` quando a config é inválida
- T5–T7: `email_notifier.go` + testes (mensagem, SMTP falso nos modos `none`, timeout e STARTTLS obrigatório). O corpo da volta também mostra `Voltou em`
- T8–T10: settings, rota `POST /api/settings/email/test` (`emailTestRateLimit = 5`) e testes de handler
- T11–T12: `EmailCard`, `Settings.jsx`; botão e resultado de teste extraídos para `TestControls.jsx` (usados pelos dois cards) e `testErrorText` para `lib/testError.js`
- T13 (checklist de segurança):
  1. Senha só usada em `smtp.PlainAuth`; não aparece em `slog`, `pingou config` (só `SMTPPasswordSet`), `GET /settings` (teste `TestSettingsEmailWithSMTP`) nem no export (não vai para o banco) ✅
  2. Sem injeção de cabeçalho: `oneLine()` em todos os cabeçalhos; teste com `\r\nBcc:`; destinatários validados por `net/mail` ✅
  3. `starttls` falha sem `STARTTLS` (`TestStartTLSRequired`), sem rebaixar ✅
  4. `tls.Config` sem `InsecureSkipVerify` (com `MinVersion: TLS 1.2`) ✅
  5. Rota de teste atrás do `authMiddleware` + rate limit; máximo de 10 destinatários ✅
  6. Log só com `recipients_count` ✅
- T15: README (variáveis, seção `## E-mail`, endpoint, `pingou config`), `.env.example` e `docker-compose.yml`
- Verificação: `go vet ./...`, `go test -race ./...`, `gofumpt -l` (só o `rm.go`, que já estava assim), `npm run lint`, `format:check`, `npm test`, `npm run build`. Tela conferida em 1440 e 375 px com e sem `PINGOU_SMTP_HOST`

### Em andamento

-

### Pendente

- T14: QA manual com Mailpit e provedor real (587/465). O ambiente de desenvolvimento não tinha Docker nem acesso a um SMTP real

## Notes

- **Próximo passo natural (fora deste plano): outbox.** Uma tabela `notification_outbox` gravada na mesma transação do `Process` e um worker com retry. Isso resolve o aviso perdido em crash e o item do PRD "Retry de webhook" (ainda desmarcado). O `MultiNotifier` vira o consumidor do worker, sem mudar os canais.
- O `webhook_url` também é um segredo (URLs do Slack e do Discord carregam token) e já vai no export. Não muda aqui; fica registrado.
- Roadmap: este plano marca "Email notification support" (Milestone 2.3 em [ROADMAP.md](.specs/project/ROADMAP.md)).

### Alternativas consideradas e descartadas

- Tabela `notification_channels` com CRUD (proposta da IA externa) → descartada: só existe um canal de cada tipo. Seria uma tabela, um repo, 5 endpoints e uma migration de dados para resolver algo que 2 chaves em `settings` resolvem. Volta à mesa se aparecer canal por monitor.
- `NotificationEvent` + interface `Channel` + `Dispatcher` → descartados: duplicam o `Notifier`, o `Monitor` e o `Incident`, que já existem.
- `github.com/wneessen/go-mail` → descartada: dependência nova para texto simples, que a stdlib cobre.
- Host, porta e usuário do SMTP editáveis na tela, com só a senha no ambiente → descartada pelo usuário: tudo no ambiente.
- Config SMTP incompleta só desabilita o card, sem derrubar a subida → descartada: erro silencioso; o usuário acha que configurou e não recebe nada. Falhar na subida é o padrão do `PINGOU_API_KEY`.
- `email_available` em `GET /api/info` → descartada: a página já carrega `/settings`, e o campo descreve uma configuração, não o build.
- Destinatários como array JSON → descartada: o usuário edita uma string com vírgulas; o backend normaliza.
- Tirar o `notify()` de dentro do `withMonitor` (devolvendo a função para o scheduler chamar depois do unlock) → descartada: mexe na assinatura do `Process` e dos testes da `StateMachine`. A goroutine dentro do `MultiNotifier` resolve sem tocar neles e também isola falhas entre canais.
