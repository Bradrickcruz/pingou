# Plano de Desenvolvimento: Notificação via ntfy nas configurações

## Overview

O Pingou avisa por webhook e por e-mail. Este plano adiciona o **ntfy** como terceiro canal, com um card próprio em Configurações, no mesmo molde do e-mail.

**Por que um canal e não só documentação.** O teste da etapa 1 (2026-11-01, ntfy 2.36.0 em Docker) mostrou que o webhook atual já funciona com o ntfy usando `?tpl=yes&auth=<base64url "Bearer tk_…">&t=…&m=…`. Funciona, mas o token vai dentro do `webhook_url`, que fica no banco e no export (o README já avisa que o export vaza o segredo do webhook), e o usuário precisa montar um template na mão. O canal próprio resolve os dois pontos: o token fica só no ambiente e a mensagem sai pronta.

> Se o objetivo for só custo zero, a alternativa é a seção "ntfy" no README com a URL de exemplo (ver *Alternativas*). Este plano segue o canal próprio, como pedido.

### Decisões

| Tema | Decisão |
| --- | --- |
| Config do servidor | **Ambiente**: `PINGOU_NTFY_URL` (liga o canal) e `PINGOU_NTFY_TOKEN` (opcional). Nada disso vai para o banco nem para o export |
| O que a tela configura | Só o **tópico**, em `settings` com a chave `ntfy_topic`. Vazio = desligado |
| Disponibilidade | Sem `PINGOU_NTFY_URL`, o card aparece **bloqueado** (igual ao e-mail sem SMTP) |
| Autenticação | Header `Authorization: Bearer <token>`. Sem token = servidor aberto |
| Envio | POST de texto em `<URL>/<tópico>` com `Title`, `Priority`, `Tags`, `Click`. Sem template, sem JSON |
| Entrega | Igual aos outros canais: assíncrona via `MultiNotifier`, uma vez por transição, sem retry |
| Teste | `POST /api/settings/ntfy/test`, rate limit de 5 por minuto por IP |
| Dependência | Nenhuma. `net/http` da stdlib |
| Migration | Não precisa. `SettingsRepo.Get` devolve `""` para chave ausente |

### Reaproveitado (não criar de novo)

- `Notifier` + `MultiNotifier.Add` ([multi_notifier.go](internal/service/multi_notifier.go)): o fan-out, o panic isolado e o `Wait` do shutdown já cobrem o novo canal.
- `subjectDown`/`subjectUp`, `bodyFor`, `oneLine`, `dashboardLink`, `downtimeSeconds` ([email_notifier.go](internal/service/email_notifier.go)): o texto da notificação sai das mesmas funções do e-mail.
- `rateLimiter` ([server.go](internal/handler/server.go)), `TestButton`/`TestResult` ([TestControls.jsx](web/src/components/settings/TestControls.jsx)), `testErrorText`, `Card`.
- `LoadSMTP` ([config.go](internal/config/config.go)) como modelo de `LoadNtfy`.

### Limites (fora do escopo)

Vários tópicos ou servidores; configuração por monitor; prioridade e ações avançadas do ntfy; outros canais (Telegram, Slack, Discord); retry; autenticação por usuário e senha (só token).

### Restrições

- Só adicionar: `GET/PATCH /api/settings` ganham `ntfy_topic` e `ntfy_available`; `webhook_url` e `email_*` não mudam.
- O `ping_token` do Heartbeat **nunca** vai na mensagem.
- `slog` nunca leva token, URL do servidor completa nem o tópico (o tópico funciona como senha em servidores abertos). Logar só `host`.
- Camadas handler → service → repository; `fmt.Errorf("...: %w")`.
- Branch `feat/ntfy-notifications` a partir de `develop`. Um commit por fase.

## Project Type

**FULLSTACK**: Go + SQLite (`internal/`, `cmd/`) e React/Vite/Tailwind (`web/`), com o front embutido via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [x] Todos os Success Criteria estão marcados
- [x] Build, lint e type-check passam
- [x] Nenhum item pendente no Status de Execução
- [x] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [x] Com `PINGOU_NTFY_URL` e um tópico salvo, uma queda e uma recuperação chegam ao ntfy, com título legível e link do dashboard (se houver `PINGOU_PUBLIC_URL`)
- [x] Com `PINGOU_NTFY_TOKEN`, o envio usa `Authorization: Bearer`; servidor com `deny-all` e token errado devolve erro no teste
- [x] Sem `PINGOU_NTFY_URL`, o card fica bloqueado e `PATCH` com `ntfy_topic` devolve 400
- [x] Token e URL do servidor não aparecem em log, em `/api/settings`, no export nem no `pingou config`
- [x] O `ping_token` não aparece em nenhuma mensagem
- [x] Botão **Enviar teste** funciona e respeita o rate limit (5 por minuto por IP)
- [x] Falha ou panic do ntfy não impede webhook nem e-mail (teste automatizado)
- [x] README documenta as variáveis, o card e um exemplo com ntfy.sh e com servidor protegido

## Tech Stack

| Technology | Choice Rationale |
| --- | --- |
| `net/http` (stdlib) | O ntfy publica com um POST simples. Não precisa de SDK |
| Headers `Title`/`Priority`/`Tags`/`Click` | Formato nativo do ntfy; evita template e JSON |
| `binwiederhier/ntfy` (Docker) | Só nos testes manuais; servidor real para validar auth e formato |

## File Structure

```
internal/
  config/config.go              # + NtfyConfig, LoadNtfy, campo Ntfy em Config
  config/config_test.go         # + casos de LoadNtfy
  service/ntfy_notifier.go      # NOVO: NtfyNotifier, SendTestNtfy
  service/ntfy_notifier_test.go # NOVO
  service/settings_service.go   # + ntfy_topic, ntfy_available, TestNtfy
  handler/settings.go           # + campo ntfy_topic, handleTestNtfy
  handler/server.go             # + ntfyTestLimiter e rota POST /settings/ntfy/test
  handler/settings_test.go      # + casos de ntfy
cmd/pingou/commands/serve.go    # + notifier.Add(NewNtfyNotifier(...))
web/src/
  api/settings.js               # + testNtfy
  components/settings/NtfyCard.jsx  # NOVO
  pages/Settings.jsx            # + estado ntfy_topic e <NtfyCard>
README.md                       # + seção "ntfy", variáveis na tabela
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
| --- | --- | --- | --- | --- | --- |
| T1 | go-backend-specialist | `NtfyConfig` e `LoadNtfy` | P0 | S | - |
| T2 | go-backend-specialist | `NtfyNotifier` e `SendTestNtfy` | P1 | M | T1 |
| T3 | go-backend-specialist | Settings: `ntfy_topic`, `ntfy_available`, `TestNtfy` | P1 | M | T1, T2 |
| T4 | go-backend-specialist | Rota de teste, rate limit e ligação no `serve` | P1 | S | T2, T3 |
| T5 | code-test | Testes de backend | P1 | M | T2, T3, T4 |
| T6 | frontend-specialist | `NtfyCard` e integração em Settings | P2 | M | T4 |
| T7 | qa-review | Teste manual contra ntfy real | P3 | S | T6 |
| T8 | document | README | P3 | S | T7 |
| T9 | security-review | Revisão de vazamento de token/tópico | P3 | S | T7 |

### Task T1: `NtfyConfig` e `LoadNtfy`

- **Agent**: go-backend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Branch `feat/ntfy-notifications` criada a partir de `develop`
- **INPUT**: `Config` só conhece `SMTP`. `LoadSMTP` é o modelo.
- **OUTPUT**:
  1. `NtfyConfig{URL, Token string}` com `Enabled()` (`URL != ""`)
  2. `LoadNtfy()`: `PINGOU_NTFY_URL` sem `/` final; exige `http://` ou `https://` com host; token opcional. Sem URL, ignora o token
  3. Chamado só no `serve` (como o SMTP), para um erro não travar `migrate`
- **VERIFY**:
  - `go test ./internal/config/...` cobre: vazio, URL válida, URL sem esquema (erro), URL com `/` final, token sem URL (ignorado)

### Task T2: `NtfyNotifier` e `SendTestNtfy`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 concluída
- **INPUT**: `EmailNotifier` já monta assunto e corpo a partir de `Monitor` + `Incident`.
- **OUTPUT**:
  1. `NtfyNotifier{cfg, getTopic func() string, publicURL}` implementando `Notifier`; sem tópico, retorna sem enviar
  2. `publish(ctx, topic, title, body, priority, tags, click)`: `POST <URL>/<tópico>` (tópico com `url.PathEscape`), corpo = texto, headers `Title`, `Priority` (queda `high`, recuperação `default`), `Tags` (`rotating_light` / `white_check_mark`), `Click` (só se houver dashboard), `Authorization: Bearer` (só se houver token)
  3. Título e corpo via `subjectDown`/`subjectUp`/`bodyFor` (sem o prefixo `[Pingou]` se ficar redundante; decidir na hora). Header `Title` passa por `oneLine`, e o valor é limitado a ASCII seguro ou codificado em UTF-8 conforme o Go aceitar (o emoji precisa ir no corpo ou em `Tags`, não no header)
  4. Cliente HTTP com timeout de 10 s, `CheckRedirect` sem seguir (redirect vira falha, como no teste do webhook)
  5. `SendTestNtfy(ctx, cfg, topic) NtfyTestResult` no formato de `WebhookTestResult`; corpo da resposta lido até 4 KiB e descartado
  6. Log só com `host` e `status`; nunca o token nem o tópico
- **VERIFY**:
  - `httptest.Server` confirma método, path, headers e corpo; status 4xx/5xx vira erro; header `Authorization` ausente quando não há token

### Task T3: Settings: `ntfy_topic`, `ntfy_available`, `TestNtfy`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1, T2
- **PRE-CONDITIONS**:
  - `SettingsService` recebe hoje só `config.SMTPConfig`
- **INPUT**: [settings_service.go](internal/service/settings_service.go) com `email_recipients`/`email_available`.
- **OUTPUT**:
  1. `NewSettingsService` recebe também `config.NtfyConfig`
  2. `Settings` ganha `NtfyTopic string json:"ntfy_topic"` e `NtfyAvailable bool json:"ntfy_available"`
  3. `Update`: com `ntfy_topic` informado e ntfy indisponível, `ErrValidation` (400). Valida o tópico: 1–64 chars, `[A-Za-z0-9_-]`, que é o que o ntfy aceita. Vazio desliga
  4. `TestNtfy(ctx, raw *string)`: usa o tópico digitado ou o salvo; sem tópico, `ErrValidation`
  5. Chamadores de `NewSettingsService` atualizados (serve e testes)
- **VERIFY**:
  - `go build ./...` e testes de `settings_test.go` passam; tópico inválido devolve 400

### Task T4: Rota de teste, rate limit e ligação no `serve`

- **Agent**: go-backend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2, T3
- **PRE-CONDITIONS**:
  - T2 e T3 concluídas
- **INPUT**: `emailTestLimiter` e a rota `POST /settings/email/test`.
- **OUTPUT**:
  1. `ntfyTestLimiter` (5 por minuto) e `POST /settings/ntfy/test` → `handleTestNtfy` (corpo opcional `{"topic": "..."}`)
  2. `updateSettingsRequest` ganha `NtfyTopic *string json:"ntfy_topic"`, repassado no `UpdateSettingsInput`
  3. Em [serve.go](cmd/pingou/commands/serve.go): se `loadedConfig.Ntfy.Enabled()`, `notifier.Add(NewNtfyNotifier(...))` lendo `ntfy_topic` do repositório; `slog.Info("ntfy notifications available", "host", …)`
  4. Conferir que `pingou config` e o export não mostram `PINGOU_NTFY_*`
- **VERIFY**:
  - `curl` na rota: 200 com `ok`, 429 na 6ª chamada, 400 sem tópico

### Task T5: Testes de backend

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2, T3, T4
- **PRE-CONDITIONS**:
  - Rotas e notifier prontos
- **INPUT**: `email_notifier_test.go`, `multi_notifier_test.go`, `settings_test.go`.
- **OUTPUT**:
  1. `ntfy_notifier_test.go`: down, up, heartbeat sem target, sem tópico (não envia), com e sem token, servidor fora do ar, redirect
  2. Teste de que o `ping_token` não aparece no corpo nem nos headers
  3. Teste no `MultiNotifier`: ntfy com panic/erro não impede os outros canais
  4. `settings_test.go`: `ntfy_available` falso/verdadeiro, PATCH com e sem ntfy, tópico inválido
- **VERIFY**:
  - `go test -race ./...` passa

### Task T6: `NtfyCard` e integração em Settings

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T4
- **PRE-CONDITIONS**:
  - API de settings devolve `ntfy_topic` e `ntfy_available`
- **INPUT**: [EmailCard.jsx](web/src/components/settings/EmailCard.jsx), [Settings.jsx](web/src/pages/Settings.jsx), [settings.js](web/src/api/settings.js).
- **OUTPUT**:
  1. `settingsApi.testNtfy(topic)` → `POST /settings/ntfy/test`
  2. `NtfyCard` copiado do molde do `EmailCard`: campo "Tópico", ajuda ("Vazio desliga o ntfy"), `TestButton`, `TestResult`, estado bloqueado com a lista das variáveis quando `available=false`, link para a seção do README
  3. `Settings.jsx`: estado `ntfyTopic`, detecção de mudança para a `UnsavedBar`, `ntfy_topic` enviado no PATCH só quando `ntfy_available`
  4. Estados cobertos: vazio, carregando, sucesso, erro, bloqueado. Label ligada ao input, `aria-describedby`, foco visível
- **VERIFY**:
  - `cd web && npm run build && npm run lint`; card renderiza bloqueado sem env e habilitado com env

### Task T7: Teste manual contra ntfy real

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T6
- **PRE-CONDITIONS**:
  - Front buildado e embutido
- **INPUT**: Docker com `binwiederhier/ntfy` (`NTFY_AUTH_DEFAULT_ACCESS=deny-all`, volume em `/var/lib/ntfy`, token via `ntfy token add`).
- **OUTPUT**:
  1. Roteiro executado: sem env (card bloqueado); servidor aberto; servidor com token certo e errado; tópico inválido; monitor HTTP que cai e volta; heartbeat que expira; Enviar teste 6 vezes seguidas
  2. Evidências anotadas no Status de Execução
- **VERIFY**:
  - Mensagens aparecem no tópico (`/<tópico>/json?poll=1`) com título, corpo e link corretos

### Task T8: README

- **Agent**: document
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T7
- **PRE-CONDITIONS**:
  - Comportamento validado na T7
- **INPUT**: Seção de e-mail e tabela de variáveis do [README.md](README.md).
- **OUTPUT**:
  1. `PINGOU_NTFY_URL` e `PINGOU_NTFY_TOKEN` na tabela de variáveis
  2. Seção "ntfy": como configurar, exemplo com ntfy.sh e com servidor próprio protegido, formato da notificação, aviso de que o tópico em servidor aberto funciona como senha
  3. Roadmap atualizado
- **VERIFY**:
  - Os exemplos do README funcionam copiados e colados

### Task T9: Revisão de vazamento

- **Agent**: security-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T7
- **PRE-CONDITIONS**:
  - Fluxo completo funcionando
- **INPUT**: Diff da branch.
- **OUTPUT**:
  1. Conferência de que token, URL do servidor e tópico não saem em log, erro de API, export, `pingou config` nem no `webhookTestErrorMessage`-equivalente do ntfy (erro do `url.Error` carrega a URL)
  2. Conferência de injeção de header (`\r\n` em nome de monitor e erro) no `Title`
- **VERIFY**:
  - `grep` por `Token` e `Topic` nos `slog.*` do diff sem ocorrência; teste com nome de monitor contendo `\r\n`

## Risk Areas

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Nome de monitor com `\r\n` ou emoji no header `Title` | Header inválido ou injeção de header | `oneLine` no título; emoji vai em `Tags`, não no header |
| `url.Error` do Go inclui a URL (e o tópico) | Tópico/URL vazam na resposta de teste e no log | Reaproveitar o tratamento de `webhookTestErrorMessage` |
| Tópico em servidor aberto é público para quem o adivinhar | Alertas expostos | Aviso no card e no README; recomendar servidor com token |
| Servidor ntfy lento | Atraso nos alertas | Já isolado pelo `MultiNotifier` e timeout de 10 s |
| `NewSettingsService` muda de assinatura | Testes e `serve` quebram | T3 atualiza todos os chamadores; `go build ./...` na verificação |

## Rollback Strategy

1. Sem `PINGOU_NTFY_URL`, o canal não é registrado e o card fica bloqueado: desligar é só remover a variável.
2. A chave `ntfy_topic` fica órfã em `settings`, sem efeito. Não há migration para desfazer.
3. Reverter o commit da fase com problema (`git revert <sha>`); as fases são independentes a partir da T2.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && go test -race ./...

# Frontend
cd web && npm run lint && npm run build

# Vazamento
grep -rn "slog\." internal/service/ntfy_notifier.go | grep -i "token\|topic"   # sem resultado
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
| --- | --- | --- | --- | --- | --- |
| Fase 0 | T1 | S | ~15 min | Config do ambiente | ✅ Concluído |
| Fase 1 | T2, T3, T4 | M | ~1 h 15 | Canal funcionando por API | ✅ Concluído |
| Fase 2 | T5 | M | ~30 min | Testes verdes | ✅ Concluído |
| Fase 3 | T6 | M | ~30 min | Card em Configurações | ✅ Concluído |
| Fase 4 | T7, T8, T9 | S | ~45 min | Validado, documentado e revisado | ✅ Concluído |

## Status de Execução

### Concluído

- Etapa 1 do `IDEAS.md`: webhook com template funciona com ntfy 2.36.0 (down, up, test, token, 403 sem token). Heartbeat mostra `<no value>` no target.
- T1 a T9, na branch `feat/ntfy-notifications`, um commit por fase.

### Evidências (2026-11-01, ntfy 2.36.0 em Docker, `deny-all` + token)

- `go build ./... && go vet ./... && go test -race ./...`: tudo passa. `npm run build` passa; `eslint` limpo nos arquivos novos/alterados (o único erro do `npm run lint` é de formatação no `EmailCard.jsx`, anterior a esta branch).
- Sem `PINGOU_NTFY_URL`: `GET /api/settings` → `ntfy_available:false`; `PATCH {"ntfy_topic":"x"}` → 400; teste → 400; card bloqueado com as variáveis.
- `PINGOU_NTFY_URL=ntfy.sh` → `serve` não sobe: `PINGOU_NTFY_URL must be a valid http(s) URL`.
- Token certo: teste `ok:true, status 200`. Token errado ou sem token: `ok:false`, `401 (token inválido ou sem permissão no tópico)`.
- Tópico `a b` → 400. 6ª chamada de teste no mesmo minuto → 429.
- Monitor HTTP derrubado e levantado: chegaram `Site local está fora do ar` (prioridade 4, `rotating_light`, `Click` com o link) e `Site local voltou (fora por 20 s)` (prioridade 3, `white_check_mark`). Acento no título chega certo (RFC 2047).
- Heartbeat expirado: `Backup diário está fora do ar`, sem `Alvo`; o `ping_token` não aparece em nenhuma mensagem do tópico.
- Log do `serve`, resposta da API, export (`/api/export`) e `pingou config` sem token, URL completa nem tópico (`pingou config` mostra só `NtfyHost` e `NtfyTokenSet`).
- Exemplo do README (servidor próprio com `deny-all`) testado: sem token 403, com token 200.

### Em andamento

-

### Pendente

-

### Desvios do plano

- 401/403 ganham uma dica no erro (`token inválido ou sem permissão no tópico`).
- `LoadNtfy` também recusa URL com query, fragmento ou usuário (o token vai só em `PINGOU_NTFY_TOKEN`).
- Título sem `[Pingou]` e sem emoji (o emoji vai em `Tags`); acentos codificados em RFC 2047, que o ntfy decodifica.
- "Entrega" virou a seção geral `## Entrega das notificações` no README, valendo para os três canais.

## Notes

- A mensagem do ntfy usa o mesmo conteúdo do e-mail para manter os canais coerentes. Se o corpo ficar longo para celular, encurtar na T2.
- Auth por query (`?auth=`) não é usada: o canal próprio manda `Authorization: Bearer`, que não aparece em log de proxy.

### Alternativas consideradas e descartadas

- `Só documentar a URL com template no README` → viável (etapa 1 passou), mas deixa o token no banco e no export, e exige template manual. Pode ser feito além deste plano, como atalho para quem já usa webhook.
- `Templates genéricos de webhook (Discord/Slack/ntfy)` → resolve mais de um canal, mas é um projeto maior e o roadmap já o trata à parte. Reavaliar quando vier o Discord/Slack.
- `Token e servidor no banco, editáveis na tela` → descartado: repete o problema do segredo no export que o README já avisa.
- `Mandar JSON para a API do ntfy (publish via JSON)` → descartado: headers + texto são mais simples e suficientes.
