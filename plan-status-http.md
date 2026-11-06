# Plano de Desenvolvimento: Status HTTP esperado

## Overview

Hoje o monitor HTTP só considera sucesso um status `2xx`, depois de seguir os redirects ([http_checker.go:60](internal/checker/http_checker.go#L60)). Isso gera falso alarme em endpoints cujo "saudável" é outro status, como uma rota protegida que devolve `401`. Este plano deixa cada monitor HTTP declarar quais status contam como sucesso, com lista e faixas: `200-299,401`.

**O que já existe e será reaproveitado:**

- A coluna `config` de `monitors` guarda JSON por tipo, sem `CHECK` no conteúdo ([0003_monitor_types.sql:20](internal/database/migrations/0003_monitor_types.sql#L20)). O campo cabe nela, então **não precisa de migration**.
- `HTTPChecker.Check` é o único lugar que decide o sucesso por status. A state machine, as estatísticas e o painel só leem `Success`.
- `MonitorService.Update` já chama `s.reload`, que recria o job com a configuração nova ([monitor_service.go:271](internal/service/monitor_service.go#L271)).
- O formulário já tem a seção "Avançado" ([MonitorForm.jsx:539](web/src/components/monitors/MonitorForm.jsx#L539)), e o detalhe já tem o `ConfigCard`.

**Limites (decididos em `IDEAS.md`):**

- Sem o campo preenchido, nada muda.
- Só monitores HTTP.
- Fora do escopo: palavra-chave no corpo da resposta, cabeçalhos e autenticação customizados.

**Restrições:**

- Sem migration e sem dependência nova.
- Monitores existentes mantêm comportamento, resposta da API (`config: {}`) e mensagem de erro idênticos.
- Cores só por tokens `pg-*` (`check:tokens` barra cor crua).

## Project Type

**FULLSTACK** — checker, validação e API em Go; um campo no formulário e uma linha no detalhe em React.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Sem o campo, `200` é sucesso, `500` é falha com a mensagem `unexpected status code: 500`, e redirects continuam sendo seguidos.
- [ ] Com `200-299,401`, uma resposta `401` é sucesso e um `500` é falha com `unexpected status code: 500 (expected 200-299,401)`.
- [ ] Com `401`, uma resposta `200` é falha (a lista substitui o padrão, não soma a ele).
- [ ] Com `301`, o redirect não é seguido e o `301` é sucesso.
- [ ] A API recusa com `400` (`VALIDATION_ERROR`) valores fora de 100 a 599, faixa invertida ou igual (`299-200`, `200-200`), itens vazios (`200,,401`), sinais (`+200`) e texto com mais de 100 caracteres.
- [ ] `PATCH` sem o campo mantém o valor; `PATCH` com `""` limpa e volta ao padrão.
- [ ] O campo é ignorado em monitores TCP e heartbeat, como os outros campos de outro tipo.
- [ ] O formulário cria, edita e limpa o campo; o detalhe mostra o esperado.
- [ ] README documenta o campo e a regra de redirect.

## Tech Stack

| Technology               | Choice Rationale                                                           |
| ------------------------ | -------------------------------------------------------------------------- |
| Go `net/http`            | `http.ErrUseLastResponse` já resolve "não seguir redirect"                 |
| Go `strings` + `strconv` | Lista e faixas cabem em um parser pequeno, sem biblioteca                  |
| Coluna `config` (JSON)   | Campo por tipo sem migration; `omitempty` mantém `{}` nos monitores atuais |
| React + tokens `pg-*`    | Mesmo padrão dos campos da seção "Avançado"                                |

## File Structure

Só dois arquivos novos, marcados abaixo.

```text
internal/domain/
  expected_status.go         # (novo) ExpectedStatus, ParseExpectedStatus, Allows, IncludesRedirect
  expected_status_test.go    # (novo) tabela de casos do parser
  monitor.go                 # MonitorConfig.ExpectedStatus
internal/repository/
  monitor_repo.go            # storedConfig, toDomain, encodeConfig (caso http)
internal/service/
  validate.go                # validateExpectedStatus, no caso http de validateMonitor
  monitor_service.go         # ExpectedStatus nos inputs de Create e Update
internal/handler/
  monitors.go                # DTO de config e toMonitorResponse
internal/checker/
  http_checker.go            # regra de sucesso, cliente sem redirect, mensagem
web/src/components/
  monitors/MonitorForm.jsx   # campo em "Avançado"
  monitor-detail/ConfigCard.jsx  # linha "Status esperados"
README.md, docs/implementacoes.md, IDEAS.md
```

## Task Breakdown

| Task ID | Agent                   | Name                        | Priority | Complexity | Dependencies |
| ------- | ----------------------- | --------------------------- | -------- | ---------- | ------------ |
| T1      | `go-backend-specialist` | Parser de status esperados  | P0       | S          | -            |
| T2      | `go-backend-specialist` | Guardar na coluna `config`  | P0       | XS         | T1           |
| T3      | `go-backend-specialist` | Validação e API             | P1       | S          | T2           |
| T4      | `go-backend-specialist` | Regra no checker e redirect | P1       | M          | T1           |
| T5      | `frontend-specialist`   | Campo no formulário         | P2       | S          | T3           |
| T6      | `frontend-specialist`   | Linha no detalhe            | P2       | XS         | T3           |
| T7      | skill `document`        | Documentação                | P3       | S          | T3, T4       |
| T8      | skill `qa-review`       | QA com o binário real       | P3       | S          | T4, T5, T6   |

### Task T1: Parser de status esperados

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Nenhuma.
- **INPUT**: A regra de sucesso está fixa no checker. O checker e o service precisam interpretar a mesma lista, e os dois já importam `domain`.
- **OUTPUT**:
  1. `type ExpectedStatus []StatusRange` em `internal/domain/expected_status.go`.
  2. `ParseExpectedStatus(spec string) (ExpectedStatus, error)`: `""` devolve `nil`. Separa por vírgula e usa `strings.Cut` em `-`. Cada número tem exatamente 3 dígitos (o `strconv.Atoi` sozinho aceitaria `+200`), vai de 100 a 599, e a faixa exige início menor que o fim.
  3. `Allows(code int) bool`: lista `nil` mantém a regra atual (`200` a `299`).
  4. `IncludesRedirect() bool`: alguma faixa cruza `300` a `399`.
- **VERIFY**:
  - `go test ./internal/domain` com tabela: `200`, `200-299,401`, `""`, e as recusas `99`, `600`, `299-200`, `200-200`, `200,,401`, `200,`, `+200`, `-200`, `200-299-300`, `abc`.
  - `Allows` com lista `nil` e com `401` (o `200` deve dar `false`).
  - **Rollback:** apagar os dois arquivos novos.

### Task T2: Guardar na coluna `config`

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 entregue.
- **INPUT**: `encodeConfig` grava `{}` para HTTP ([monitor_repo.go:210](internal/repository/monitor_repo.go#L210)).
- **OUTPUT**:
  1. `MonitorConfig.ExpectedStatus string`. Guarda o texto, não a lista já interpretada: é o que a API devolve e o que a mensagem de erro mostra.
  2. `storedConfig.ExpectedStatus` com `json:"expected_status,omitempty"`, e `toDomain` copiando o valor.
  3. Caso `domain.TypeHTTP` em `encodeConfig`.
- **VERIFY**:
  - Teste de repositório: criar, ler e atualizar preserva o campo; monitor HTTP sem o campo continua gravando `{}`.
  - **Rollback:** reverter o commit. Sem migration, o banco não muda de forma.

### Task T3: Validação e API

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - T2 entregue.
- **INPUT**: `validateMonitor` valida por tipo. O `PATCH` usa ponteiros, em que `nil` mantém o valor.
- **OUTPUT**:
  1. `monitorConfigRequest.ExpectedStatus *string` com `json:"expected_status"`.
  2. `CreateMonitorInput.ExpectedStatus string` e `UpdateMonitorInput.ExpectedStatus *string`, aplicados só no caso HTTP.
  3. No service, tirar os espaços antes de guardar: `strings.Join(strings.Fields(s), "")`. Assim `200-299, 401` vira `200-299,401`.
  4. `validateExpectedStatus` em `validate.go`: no máximo 100 caracteres, mais `ParseExpectedStatus`, com erro embrulhando `ErrValidation`. Chamada no caso HTTP de `validateMonitor`.
  5. `toMonitorResponse`: `res.Config["expected_status"]` só quando preenchido, para a resposta dos monitores atuais não mudar.
- **VERIFY**:
  - `validate_test.go`: casos válidos e inválidos.
  - `handler_test.go`: `POST` com valor inválido responde `400` `VALIDATION_ERROR`; `PATCH` sem o campo mantém; `PATCH` com `""` limpa; monitor TCP ignora o campo; monitor HTTP sem o campo responde `config: {}`.
  - **Rollback:** reverter o commit.

### Task T4: Regra no checker e redirect

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 entregue.
- **INPUT**: O cliente segue redirects antes de o status ser comparado, então um `301` esperado nunca chegaria à comparação. Por isso a lista sozinha não cobre o caso citado na ideia.
- **OUTPUT**:
  1. Segundo cliente em `HTTPChecker`, `noRedirect`, com o mesmo `Timeout` e `CheckRedirect` devolvendo `http.ErrUseLastResponse`. Os dois usam o transporte padrão, então compartilham as conexões.
  2. Em `Check`: interpretar `m.Config.ExpectedStatus`, usar `noRedirect` quando `IncludesRedirect()` for verdadeiro, e decidir o sucesso com `Allows`.
  3. Mensagem de erro: com lista, `unexpected status code: 500 (expected 200-299,401)`; sem lista, o texto atual, sem mudança.
  4. Erro de parser no checker (só por edição manual do banco) cai na regra padrão, com comentário `ponytail:`.
- **VERIFY**:
  - `checker_test.go` com `httptest`: `401` com e sem lista; `200` com lista `401`; `301` para um destino com `200` usando as listas `301` (sucesso, destino não chamado) e vazia (sucesso, destino chamado); `500` com a mensagem nova e com a antiga.
  - `TestRegistryHTTP` e `TestHTTPCheckerErrorMessage` passam sem alteração.
  - **Rollback:** reverter o commit. Monitores com o campo preenchido voltam à regra `2xx`.

### Task T5: Campo no formulário

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T3
- **PRE-CONDITIONS**:
  - API de T3 disponível.
- **INPUT**: `MonitorForm` monta `payload()` por tipo e traduz erros em `friendlyError`.
- **OUTPUT**:
  1. `expected_status` em `blank` e em `fromMonitor` (`m.config?.expected_status ?? ""`).
  2. Campo "Status esperados" dentro de "Avançado", só para HTTP: fonte mono, `maxLength={100}`, `placeholder="200-299"`, `id` por `useId` e `aria-describedby` no texto de ajuda.
  3. Ajuda: "Deixe vazio para aceitar qualquer 2xx. Ex.: 200-299,401. Se a lista incluir um 3xx, o Pingou não segue o redirect."
  4. `payload()` de HTTP envia `config: { expected_status: form.expected_status.trim() }`, sempre, para que esvaziar o campo limpe o valor.
  5. `friendlyError`: `expected_status` vira "Status esperados inválidos. Use números de 100 a 599 e faixas, como 200-299,401."
- **VERIFY**:
  - `npm run lint`, `npm run check:tokens` e `npm run build` passam.
  - No navegador: criar, editar e limpar o campo; valor inválido mostra a mensagem em português; o campo não aparece em TCP nem em heartbeat.
  - **Rollback:** reverter o commit; a API aceita o monitor sem o campo.

### Task T6: Linha no detalhe

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T3
- **PRE-CONDITIONS**:
  - API de T3 disponível.
- **INPUT**: `ConfigCard` monta as linhas por tipo em uma lista ([ConfigCard.jsx:6](web/src/components/monitor-detail/ConfigCard.jsx#L6)).
- **OUTPUT**:
  1. Linha "Status esperados", só em HTTP, antes de "Criado em": `m.config?.expected_status || "2xx (padrão)"`.
- **VERIFY**:
  - No navegador, nos temas claro e escuro: monitor com lista, monitor sem lista, e TCP e heartbeat sem a linha.
  - **Rollback:** reverter o commit.

### Task T7: Documentação

- **Agent**: skill `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T3, T4
- **PRE-CONDITIONS**:
  - Comportamento final de T3 e T4 confirmado.
- **INPUT**: A tabela de `config` do README diz que HTTP "usa o campo `url`" ([README.md:232](README.md#L232)).
- **OUTPUT**:
  1. README: `expected_status` na tabela de `config`, formato aceito, exemplo de `curl`, regra do redirect e o exemplo da mensagem de erro.
  2. `docs/implementacoes.md`: item 18, com decisões, descartes e validação.
  3. `IDEAS.md`: remover a ideia e a linha do resumo, e atualizar a data.
- **VERIFY**:
  - `grep -n expected_status README.md docs/implementacoes.md` encontra os trechos.
  - O exemplo de `curl` do README funciona no binário real (T8).
  - **Rollback:** reverter o commit.

### Task T8: QA com o binário real

- **Agent**: skill `qa-review`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T4, T5, T6
- **PRE-CONDITIONS**:
  - `make build` passa.
- **INPUT**: Testes unitários verdes, com o ciclo completo do incidente ainda não exercitado.
- **OUTPUT**:
  1. Evidências registradas em `docs/implementacoes.md`.
- **VERIFY**:
  - Servidor local devolvendo `401`: o monitor cai sem a lista. Ao salvar `401`, o próximo check dá sucesso, o incidente fecha e o aviso de recuperação sai.
  - Redirect `301` com lista `301`: sobe sem chamar o destino.
  - Monitor `https://` com lista `301`: `cert_expires_at` continua sendo preenchido.
  - Painel em largura de celular, sem rolagem horizontal, e teclado no campo novo.
  - **Rollback:** nenhum; a task não altera código.

## Risk Areas

| Risk                                               | Impact                                                                                    | Mitigation                                                                                    |
| -------------------------------------------------- | ----------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Lista com 3xx deixa de olhar o destino do redirect | `200-399` marca "no ar" com o destino quebrado                                            | Regra explícita na ajuda do campo e no README                                                 |
| Certificado lido de outro host com lista 3xx       | Sem seguir o redirect, a data passa a ser a do host da URL, não a do destino              | É o certificado do endereço cadastrado, então o comportamento fica aceitável; conferido em T8 |
| Lista sem nenhum 2xx (`401`)                       | Uma rota que volta a responder `200` passa a contar como queda                            | É o pedido literal da lista; a mensagem de erro mostra o esperado                             |
| `100` a `199` na lista nunca casam                 | O `net/http` consome 1xx (menos `101`); um monitor só com `100` ficaria sempre fora do ar | Mantida a faixa da ideia; documentar no README                                                |
| Salvar a lista em monitor com incidente aberto     | O próximo check fecha o incidente e gera aviso de recuperação                             | Comportamento correto; verificado em T8                                                       |

## Rollback Strategy

1. Antes de voltar de versão, limpar o campo dos monitores que o usam: `PATCH` com `{"config":{"expected_status":""}}`. Sem isso, um monitor com `401` esperado cai assim que a versão antiga rodar.
2. Reverter os commits ou subir o binário anterior. Não há migration para desfazer: a versão antiga ignora a chave extra no JSON de `config`.

## Phase X: Verification

```bash
# Backend
gofumpt -l .
go vet ./...
go test -race ./...

# Frontend
cd web && npm run lint && npm test && npm run check:tokens && npm run format:check && npm run build

# Binário completo e grafo
cd .. && make build
graphify update .
```

## Cronograma Sugerido

| Fase   | Tarefas | Complexity | Esforço | Target                                   | Status      |
| ------ | ------- | ---------- | ------- | ---------------------------------------- | ----------- |
| Fase 0 | T1, T2  | S          | 20 min  | Parser testado e campo persistido        | ⬜ Pendente |
| Fase 1 | T3, T4  | M          | 45 min  | API e checker funcionando                | ⬜ Pendente |
| Fase 2 | T5, T6  | S          | 20 min  | Campo no painel                          | ⬜ Pendente |
| Fase 3 | T7, T8  | S          | 30 min  | Documentado e verificado no binário real | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T8

## Notes

- **A ideia diz "esperar um caso real antes de fazer".** Antes de implementar, vale anotar aqui qual é o caso concreto: ele pode mostrar que basta uma lista simples, ou que o problema é a palavra-chave no corpo, que está fora deste plano.
- **Regra de redirect (decisão deste plano, não estava em `IDEAS.md`):** lista com algum 3xx não segue redirect; qualquer outra lista segue, como hoje.
- **Mensagem padrão preservada:** a ideia pede que a mensagem fale do esperado, mas só muda quando há lista. Sem lista, o texto segue igual, e nada muda para quem já usa o Pingou.
- **Página pública (ideia 3):** `expected_status` não deve entrar na lista de permissão dos campos públicos.

### Alternativas consideradas e descartadas

- `Campo follow_redirects separado` → descartada por criar uma segunda opção para o único caso em que ela importa, que a lista já expressa.
- `Coluna nova com migration` → descartada porque `config` foi feita para campos por tipo.
- `Lista interpretada guardada no JSON` → descartada porque o texto já serve à API e à mensagem, e interpretar a cada check custa muito pouco perto da requisição HTTP.
- `Validar também no frontend` → descartada porque a API já valida e o `friendlyError` traduz a mensagem; duas validações poderiam divergir.
- `Mudar a mensagem também sem lista` → descartada por alterar o texto de todos os monitores e os exemplos do README, sem ganho.
