# Plano de Desenvolvimento: Nomes claros nos eventos do webhook (ideia 7)

## Overview

Hoje o campo `event` do JSON do webhook vale `down` ou `up` (além de `cert_expiring` e `test`). São nomes vagos e deixam pouco espaço para eventos novos (ex.: `maintenance_started`, ideia 6).

Esta mudança renomeia `down` → `monitor_down` e `up` → `monitor_up`, no padrão `<objeto>_<acontecimento>`, minúsculas, separado por `_`.

**Limites (decididos em 2026-11-05):**

- **Quebra limpa.** Sem período de transição e sem enviar o nome antigo.
- Os outros campos do JSON não mudam (`monitor`, `timestamp`, `last_error`, `downtime_duration_seconds`, `dashboard_url`).
- Sem versionar o formato do JSON.
- Só o dono usa o webhook hoje (fluxo no n8n), então esse é o único fluxo a ajustar.

**Restrição:** a quebra é silenciosa. Fluxo "se `down`, alerta; senão, voltou" passa a tratar toda queda como "voltou". Por isso o README e a nota da versão precisam destacar isso, e o README deve trazer a regra "evento desconhecido: ignore".

**Ordem:** sai antes ou junto das ideias 5 (`id` no JSON) e 6 (eventos de manutenção), para quem integra ajustar o fluxo uma vez só.

## Project Type

**BACKEND** — mudança de contrato do JSON no notificador Go, com testes e documentação. O frontend não referencia os nomes de evento (verificado por grep em `web/src`).

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e testes passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] `NotifyDown` envia `"event": "monitor_down"` e `NotifyRecovery` envia `"event": "monitor_up"`.
- [ ] Nenhum `"down"` ou `"up"` restou como valor de `event` no código, nos testes, no README ou em exemplos.
- [ ] Testes novos cobrem os dois nomes (hoje nenhum teste Go confere `down`/`up`).
- [ ] README: exemplos atualizados, padrão de nomes documentado, aviso de quebra e regra "ignore o que não conhece".
- [ ] `docs/implementacoes.md` e `IDEAS.md` refletem a entrega.
- [ ] Decisão sobre `cert_expiring` e `test` registrada (ver Notes).

## Tech Stack

| Technology                            | Choice Rationale                                                                                    |
| ------------------------------------- | --------------------------------------------------------------------------------------------------- |
| Go (`internal/service/notifier.go`)   | O `Event` é montado só em `NotifyDown`, `NotifyRecovery`, `NotifyCertExpiring` e `SendTestWebhook`. |
| Constantes de string no mesmo arquivo | Os nomes hoje são literais espalhados; constantes evitam novo erro de digitação. Sem pacote novo.   |
| `go test`                             | Já usado em `notifier_test.go` e `cert_notify_test.go`.                                             |

## File Structure

Arquivos tocados (nenhum arquivo novo de código):

```
internal/service/notifier.go        # constantes + 2 literais trocados
internal/service/notifier_test.go   # testes de monitor_down e monitor_up
README.md                           # seção Webhook (L~835-965)
docs/implementacoes.md              # novo item de histórico
IDEAS.md                            # tirar a ideia 7 da lista de pendentes
```

E-mail e ntfy **não** são tocados: o `event` neles é só rótulo de log (`email_notifier.go:68`, `ntfy_notifier.go:114`) e o `fanOut` (`multi_notifier.go:49`) também. Ver T3.

## Task Breakdown

| Task ID | Agent                   | Name                                         | Priority | Complexity | Dependencies   |
| ------- | ----------------------- | -------------------------------------------- | -------- | ---------- | -------------- |
| T0      | `devops`                | Criar a branch `feature/webhook-event-names` | P0       | XS         | -              |
| T1      | `go-backend-specialist` | Renomear eventos no notificador              | P1       | XS         | T0             |
| T2      | `code-test`             | Testes dos eventos novos                     | P1       | S          | T1             |
| T3      | `go-backend-specialist` | Revisar rótulos de log (e-mail/ntfy/fanOut)  | P3       | XS         | T1             |
| T4      | `document`              | README, histórico e IDEAS                    | P3       | S          | T1             |
| T5      | `qa-review`             | Verificação ponta a ponta com receptor local | P3       | S          | T1, T2         |
| T6      | `devops`                | Abrir a Pull Request para `develop`          | P3       | XS         | T2, T3, T4, T5 |

### Task T0: Criar a branch de feature

- **Agent**: `devops`
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Estar em `develop`, com a árvore limpa e atualizada (`git status` limpo; `git pull --ff-only`).
- **INPUT**: Branch `develop` (as features anteriores usaram o prefixo `feature/` ou `feat/`).
- **OUTPUT**:
  1. Branch `feature/webhook-event-names` criada a partir de `develop` e já ativa.
  2. O arquivo do plano (`plan-webhook-event-names.md`) vai no primeiro commit da branch (ou fica fora dele, se preferir; o histórico diz que os planos são removidos depois).
- **VERIFY**:
  - `git branch --show-current` mostra `feature/webhook-event-names`.
- **ROLLBACK**: `git switch develop && git branch -D feature/webhook-event-names` (ainda sem push).

### Task T1: Renomear eventos no notificador

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T0
- **PRE-CONDITIONS**:
  - Na branch `feature/webhook-event-names` (T0).
- **INPUT**: `notifier.go` com `Event: "down"` (L110) e `Event: "up"` (L131).
- **OUTPUT**:
  1. Constantes `eventMonitorDown = "monitor_down"` e `eventMonitorUp = "monitor_up"` (e, por consistência, `eventCertExpiring`, `eventTest` com os valores atuais).
  2. `NotifyDown`, `NotifyRecovery`, `NotifyCertExpiring` e `SendTestWebhook` usam as constantes.
  3. Comentário curto no `webhookPayload` com o padrão `<objeto>_<acontecimento>`.
- **VERIFY**:
  - `go build ./...` passa.
  - `grep -n 'Event:' internal/service/notifier.go` não mostra literal solto.
- **ROLLBACK**: `git checkout internal/service/notifier.go`.

### Task T2: Testes dos eventos novos

- **Agent**: `code-test`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 concluída.
- **INPUT**: `notifier_test.go` tem `TestWebhookPayloadByType` (L15) e o teste do `test`; nenhum confere `down`/`up`.
- **OUTPUT**:
  1. Teste com `httptest.Server` que chama `NotifyDown` e `NotifyRecovery` e confere `event == "monitor_down"` / `"monitor_up"`, no estilo do `cert_notify_test.go:36`.
  2. Na recuperação, confere também `downtime_duration_seconds` preenchido; na queda, `null` (garante que só o nome mudou).
- **VERIFY**:
  - `go test ./internal/service/ -run 'Webhook|Notify' -v` passa.
  - Reverter T1 localmente faz o teste novo falhar (prova que ele confere o nome).
- **ROLLBACK**: apagar os testes novos.

### Task T3: Revisar rótulos de log

- **Agent**: `go-backend-specialist`
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 concluída.
- **INPUT**: `multi_notifier.go` chama `fanOut(ctx, event, ...)` com um nome de evento só para o log de pânico; e-mail e ntfy recebem `event` para o log.
- **OUTPUT**:
  1. Conferir quais strings esses chamadores passam. Se forem `"down"`/`"up"`, usar as mesmas constantes para o log bater com o webhook. Se já forem outra coisa, **não mexer**.
- **VERIFY**:
  - `grep -rn 'fanOut(' internal/service` mostra rótulos coerentes.
  - `go test ./internal/service/...` passa.
- **ROLLBACK**: `git checkout` dos arquivos tocados.

### Task T4: README, histórico e IDEAS

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Nomes finais confirmados (ver Notes sobre `cert_expiring`).
- **INPUT**: README com `"event": "down"` (L839), `"up"` (L858), parágrafo L916 "Os eventos `down` e `up` não mudam" e aviso L918.
- **OUTPUT**:
  1. Exemplos de queda e recuperação com `monitor_down` e `monitor_up`.
  2. Nova subseção "Nomes dos eventos": tabela (`monitor_down`, `monitor_up`, `cert_expiring`, `test`) e o padrão `<objeto>_<acontecimento>`.
  3. Bloco de destaque **"Mudança que quebra compatibilidade"**, com exemplo do risco ("se `down`, alerta; senão, voltou" vira tudo "voltou") e como migrar.
  4. Regra para quem integra: **evento desconhecido deve ser ignorado**, nunca tratado como `up`. Já existe a frase em L963; promover para o aviso principal.
  5. Atualizar L916 e L918 (que citam `down`/`up` como nomes atuais).
  6. Item novo em `docs/implementacoes.md` (data 2026-11-05, "quebra limpa") e linha na tabela do topo.
  7. `IDEAS.md`: remover a ideia 7 do resumo e do corpo; ajustar as referências nas ideias 4, 5 e 6 que dizem "ideia 7" para "já entregue".
- **VERIFY**:
  - `grep -nE '"event": "(down|up)"' README.md` não retorna nada.
  - `grep -n 'monitor_down\|monitor_up' README.md` mostra exemplos e tabela.
  - Links internos do README continuam válidos.
- **ROLLBACK**: `git checkout README.md docs/implementacoes.md IDEAS.md`.

### Task T5: Verificação ponta a ponta

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T1, T2
- **PRE-CONDITIONS**:
  - Build local funcionando.
- **INPUT**: Pingou rodando com um receptor local (ex.: `nc -l` ou um servidor Python de 5 linhas) como `webhook_url`.
- **OUTPUT**:
  1. Evidência de que uma queda real e a volta chegam com `monitor_down` e `monitor_up`.
  2. O botão "Testar webhook" continua mandando `test`.
- **VERIFY**:
  - JSON capturado colado em "Evidências" no Status de Execução.
- **ROLLBACK**: n/a (só leitura).

### Task T6: Abrir a Pull Request

- **Agent**: `devops`
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T2, T3, T4, T5
- **PRE-CONDITIONS**:
  - Phase X (verificação) passou na branch.
  - Mudanças commitadas (mensagem começando com `BREAKING:`, ver Notes).
  - **Confirmar com o dono antes do `push` e de abrir a PR**: são ações visíveis fora da máquina.
- **INPUT**: Branch `feature/webhook-event-names` com os commits das tasks anteriores.
- **OUTPUT**:
  1. `git push -u origin feature/webhook-event-names`.
  2. PR para `develop` via `gh pr create`, com título como "BREAKING: nomes claros nos eventos do webhook (`monitor_down`/`monitor_up`)".
  3. Descrição da PR: o que mudou, tabela antigo → novo, **aviso de quebra** (fluxo "se `down`... senão `up`" falha em silêncio), regra "ignore evento desconhecido", como foi testado (evidências da T5) e checklist do n8n.
- **VERIFY**:
  - `gh pr view --web` ou `gh pr view` mostra a PR aberta, base `develop`, e os checks rodando.
- **ROLLBACK**: `gh pr close <n>` e, se preciso, `git push origin --delete feature/webhook-event-names`.

## Risk Areas

| Risk                                               | Impact                                     | Mitigação                                                                                             |
| -------------------------------------------------- | ------------------------------------------ | ----------------------------------------------------------------------------------------------------- |
| Fluxo do dono (n8n) trata "não é `down`" como `up` | Toda queda vira "voltou", sem erro visível | Ajustar o fluxo **antes** de subir a versão nova; README com aviso e regra "ignore o que não conhece" |
| Esquecer algum literal `"down"`/`"up"`             | Parte dos avisos sai com nome antigo       | Constantes (T1) + grep de verificação + testes (T2)                                                   |
| Pessoa futura reintroduz literal solto             | Nomes inconsistentes                       | Constantes e comentário do padrão no `webhookPayload`                                                 |
| Mudar `cert_expiring`/`test` sem necessidade       | Segunda quebra desnecessária               | Manter como estão (ver Notes)                                                                         |

## Rollback Strategy

1. `git revert` do commit da mudança (é pequena e isolada em `notifier.go`).
2. Voltar o fluxo do n8n para os nomes antigos.
3. Como é quebra limpa e sem estado no banco, não há migration nem dado a desfazer.

## Phase X: Verification

```bash
# Build e testes
go build ./... && go vet ./...
go test ./internal/service/...

# Nada de nome antigo
grep -rnE 'Event:\s+"(down|up)"' internal
grep -nE '"event": "(down|up)"' README.md

# Front não é afetado, mas confirma
cd web && npm run build
```

## Cronograma Sugerido

| Fase   | Tarefas | Complexity | Esforço | Target          | Status      |
| ------ | ------- | ---------- | ------- | --------------- | ----------- |
| Fase 0 | T0      | XS         | ~5 min  | Branch criada   | ⬜ Pendente |
| Fase 1 | T1, T2  | S          | ~20 min | Código e testes | ⬜ Pendente |
| Fase 2 | T3, T4  | S          | ~20 min | Docs e rótulos  | ⬜ Pendente |
| Fase 3 | T5      | S          | ~15 min | Evidência       | ⬜ Pendente |
| Fase 4 | T6      | XS         | ~10 min | PR aberta       | ⬜ Pendente |

## Status de Execução

### Concluído

- Leitura de `notifier.go`, README (seção Webhook) e IDEAS.md; mapa dos pontos a mudar.

### Em andamento

-

### Pendente

- T0 a T6

### Evidências

-

## Notes

- **`cert_expiring` e `test` ficam como estão (proposta).** `cert_expiring` já segue `<objeto>_<acontecimento>` (certificado + expirando), e renomear de novo seria uma segunda quebra sem ganho. `test` não é um evento de objeto. **Confirmar antes da T4.** Se preferir `monitor_cert_expiring`, é só trocar a constante e o README.
- A ideia 7 pergunta "se o aviso de SSL também sai por webhook": sim, sai (`cert_expiring`, já entregue). A resposta está acima.
- Não há frontend a mexer: a tela de Webhook não mostra JSON de exemplo com `event` (grep em `web/src` sem ocorrências).
- O `id` do evento (ideia 5) e os eventos de manutenção (ideia 6) **não** entram aqui. Quando entrarem, seguem o padrão documentado em T4.
- Mensagem de commit/nota da versão deve começar com "BREAKING:".

### Alternativas consideradas e descartadas

- Enviar `event` antigo e novo juntos por um período → descartada: decisão de quebra limpa (2026-11-05), e só o dono usa o webhook.
- Campo `event_version` ou versionar o JSON → descartada: fora do escopo da ideia.
- Pacote/tipo `EventType` com validação → descartada: constantes de string bastam para 4 valores.
