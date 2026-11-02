# Plano de Desenvolvimento: Aviso de expiração do certificado SSL

## Overview

Hoje o Pingou avisa quando um site cai, mas não avisa quando o certificado SSL está perto de vencer. Quando o certificado vence, o site "cai" para todos os visitantes. O aviso antecipado evita esse susto.

**O que será construído:**

1. Cada monitor do tipo **HTTP** ganha uma opção (checkbox): "Avisar quando o certificado SSL estiver perto de expirar".
2. Com a opção ligada, o Pingou envia **um aviso 30 dias antes** e **um aviso 7 dias antes** da expiração.
3. Os avisos vão por **todos os canais configurados** (webhook, e-mail e ntfy).
4. A **data de expiração aparece na página de detalhes** do monitor HTTP.

**Como funciona (visão simples):** o checker HTTP já faz a conexão TLS a cada check. O certificado do servidor já vem dentro da resposta (`resp.TLS`), então ler a data de validade (`NotAfter`) **não custa nenhuma conexão extra**. Depois é só comparar com a data de hoje.

**Decisões já tomadas (pelo dono do projeto):**

- Avisar 30 dias antes e 7 dias antes, uma vez cada
- Opcional por monitor, só para monitores HTTP, como checkbox
- Vai por todos os canais configurados
- Data de expiração visível na página de detalhes do monitor HTTP

**Decisões assumidas (revisar se discordar):**

| #   | Suposição                                                                                                                                           | Se discordar                                                            |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| 1   | A **data aparece em todo monitor HTTP com URL https**, mesmo com a opção desligada (ler a data é grátis). O checkbox controla só o **aviso**        | Mostrar a data só quando a opção estiver ligada (mudança pequena na T9) |
| 2   | O checkbox começa **desmarcado** para monitores novos e existentes                                                                                  | Marcar por padrão em monitores https novos (mudança de uma linha na T8) |
| 3   | O certificado lido é o da **resposta final**, depois dos redirects (o checker segue até `PINGOU_MAX_REDIRECTS`). É o certificado que o visitante vê | Capturar também o do primeiro host (exige `httptrace`, mais complexo)   |
| 4   | O recurso só vale para URL que começa com `https://`. Com `http://` o checkbox fica desabilitado                                                    | -                                                                       |
| 5   | "Uma vez" significa **uma vez por certificado e por limite**. Se o certificado for renovado (data nova), os avisos "rearmam" sozinhos               | -                                                                       |
| 6   | O limite é medido por duração: 30 dias = 30 × 24 h antes da data de expiração                                                                       | -                                                                       |

**Fora de escopo:**

- Validade de certificados intermediários (só o certificado do próprio site)
- Monitores TCP com TLS
- Limites configuráveis (30 e 7 são fixos)
- Flag no CLI (`pingou add --cert-alert`)
- Aviso de certificado inválido por outros motivos (isso já vira incidente DOWN)
- Ícone ou badge de certificado na lista de monitores

**Restrições:**

- Go com a biblioteca padrão, sem dependência nova
- Notificação enviada **depois** do commit, como já é feito no restante do projeto
- Nenhum check a mais, nenhuma conexão a mais, nenhuma escrita extra por check
- Pico de RAM continua abaixo de 50 MB (meta do README)
- Comportamento atual de queda e recuperação não muda

## Project Type

**FULLSTACK** — Go (checker, state machine, notificações, API, SQLite + goose) e SPA React em `web/` (checkbox no formulário e data na página de detalhes).

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

**Comportamento**

- [ ] O checker HTTP guarda a data de expiração (`NotAfter`) do certificado a cada resposta HTTPS, inclusive quando o status HTTP é de erro
- [ ] Monitor HTTP com a opção ligada recebe **exatamente 1 aviso** ao cruzar 30 dias e **exatamente 1 aviso** ao cruzar 7 dias, sem repetir nos checks seguintes
- [ ] Se a opção for ligada quando já faltam menos de 7 dias, sai **só o aviso de 7 dias** (não dois avisos de uma vez)
- [ ] Se o certificado for renovado (data nova), os avisos voltam a funcionar para o certificado novo
- [ ] Se a URL do monitor mudar, a data guardada e as marcas de aviso são zeradas
- [ ] O aviso chega por todos os canais configurados (webhook, e-mail e ntfy), com o texto correto em cada um
- [ ] Monitor com a opção desligada nunca envia aviso, mas continua mostrando a data
- [ ] Os avisos de queda e recuperação continuam iguais (nenhuma regressão)

**API**

- [ ] `POST` e `PATCH /api/monitors` aceitam `cert_alert_enabled`; `GET /api/monitors/:id` devolve `cert_alert_enabled` e `cert_expires_at`
- [ ] Ligar a opção em monitor que não é HTTP, ou cuja URL não é `https://`, retorna erro 400 claro

**Interface**

- [ ] Checkbox no formulário de monitor HTTP, desabilitado quando a URL não é https
- [ ] A página de detalhes mostra a data de expiração e os dias que faltam, com cor de atenção (30 dias) e de alerta (7 dias ou vencido)
- [ ] O visual segue o `DESIGN.md`

**Qualidade**

- [ ] Regra de aviso testada em função pura, com tabela de casos (limites, renovação, ativação tardia, certificado já vencido)
- [ ] Teste do checker com certificado de validade conhecida, gerado no próprio teste
- [ ] Teste de integração: sequência de checks ao longo do tempo gera o número certo de avisos
- [ ] Migration sobe e desce sem perder dados
- [ ] README atualizado (recurso, campos da API, novo evento do webhook)
- [ ] `make test` e `make build` passam

## Tech Stack

| Technology                                                             | Choice Rationale                                                                                                          |
| ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Go `net/http` + `crypto/tls` (`resp.TLS.PeerCertificates[0].NotAfter`) | O certificado já vem na resposta que o checker usa hoje. Sem dependência nova e sem conexão extra                         |
| Go `crypto/x509` (só nos testes)                                       | Gera um certificado com a validade que o teste precisar, sem esperar 30 dias                                              |
| SQLite + `goose`                                                       | Migration no mesmo padrão das anteriores (colunas novas na tabela `monitors`)                                             |
| Unit of Work existente (`StateMachine.Process`)                        | A atualização da data e a marca de "aviso enviado" entram na mesma transação do check. O envio fica para depois do commit |
| Notificador existente (webhook, e-mail, ntfy)                          | Reaproveita os canais configurados. Só ganha um evento novo                                                               |
| React + Tailwind                                                       | Checkbox e linha de detalhe, seguindo o `DESIGN.md`                                                                       |

## File Structure

Caminhos prováveis. A task T1 confirma os caminhos reais (o README mostra `migrations/` na raiz, mas planos anteriores usam `internal/database/migrations/`).

```text
internal/
  database/migrations/00NN_monitor_cert_expiry.sql   (novo; número = próximo livre)
  domain/monitor.go                  (+ campos de certificado)
  checker/                           (Result ganha CertNotAfter; checker HTTP preenche)
  repository/monitor_repo*.go        (persistir e ler os campos; update guardado pela URL)
  service/cert_alert.go              (novo: regra pura de aviso)
  service/cert_alert_test.go         (novo: testes em tabela)
  service/state_machine.go           (+ passo de certificado dentro da transação)
  service/notifier*.go               (novo evento nos 3 canais)
  handler/monitor*.go                (campo na API + validação)
web/src/
  components/                        (checkbox no formulário do monitor)
  pages/                             (data de expiração na página de detalhes)
scripts/certserver/main.go           (novo: servidor HTTPS local com validade escolhida, para QA)
README.md                            (recurso, API, evento do webhook)
DESIGN.md                            (se surgir componente novo)
```

## Task Breakdown

| Task ID | Agent                   | Name                                                 | Priority | Complexity | Dependencies   |
| ------- | ----------------------- | ---------------------------------------------------- | -------- | ---------- | -------------- |
| T1      | `go-backend-specialist` | Levantamento do código atual                         | P0       | S          | -              |
| T2      | `design-db`             | Migration dos campos de certificado                  | P0       | S          | T1             |
| T3      | `go-backend-specialist` | Checker HTTP captura a data do certificado           | P1       | M          | T1             |
| T4      | `design-api`            | Domínio, repositório e API (campo, validação, reset) | P1       | M          | T2             |
| T5      | `go-backend-specialist` | Regra pura de aviso (30 e 7 dias)                    | P1       | M          | T1             |
| T6      | `go-backend-specialist` | Notificações: novo evento nos 3 canais               | P1       | L          | T1             |
| T7      | `go-backend-specialist` | Integração na state machine                          | P1       | L          | T3, T4, T5, T6 |
| T8      | `frontend-specialist`   | Checkbox no formulário do monitor                    | P2       | S          | T4             |
| T9      | `frontend-specialist`   | Data de expiração na página de detalhes              | P2       | M          | T4             |
| T10     | `code-test`             | Testes de integração e regressão                     | P3       | M          | T7             |
| T11     | `qa-review`             | QA manual com certificado de validade escolhida      | P3       | M          | T7, T8, T9     |
| T12     | `document`              | Documentação                                         | P3       | S          | T7, T9         |

**Duas trilhas em paralelo:** depois da T4, o frontend (T8 e T9) anda sozinho. Enquanto isso, o backend faz T3, T5 e T6, que não dependem uma da outra, e depois a T7.

### Task T1: Levantamento do código atual

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Repositório atualizado na `main`
  - Ler `plan-fuso-shutdown.md`, `plan-banco-indices.md` e `plan-pause-closes-incident.md` para saber quais números de migration já estão reservados
- **INPUT**: o código atual, conhecido só pelo README e pelos outros planos
- **OUTPUT**: nota curta, sem mudar código, respondendo:
  1. **Checker HTTP:** onde está, qual é o formato do `Result`, se `resp` está disponível no ponto em que o resultado é montado, se existe opção de ignorar a verificação do certificado (TLS inseguro) e como os redirects são tratados
  2. **Tabela `monitors`:** colunas atuais, tipo usado para datas (texto ISO ou inteiro) e se existe coluna `type`
  3. **`StateMachine.Process`:** onde entra o `Result`, como o monitor é carregado, **se `UpdateWithTx` grava a linha inteira ou só campos específicos** (importante para o risco de sobrescrever a edição do usuário) e onde ficam as notificações pós-commit
  4. **Notificações:** interface do notificador, formato do payload do webhook (existe template configurável?), formato do e-mail, formato do ntfy (título, prioridade, tags), idioma dos textos
  5. **API e validação de monitor:** handler, DTOs de criação e atualização, onde a validação por tipo acontece
  6. **Frontend:** arquivos do formulário de monitor e da página de detalhes, o componente de data/hora existente e se o plano de fuso (T8 dele) já foi mesclado
  7. **Último número de migration** e os números reservados pelos outros planos (pausa: `0005`)
  8. Existe infraestrutura de teste no frontend (vitest ou similar) e uma forma de injetar relógio/"agora" nos testes do backend?
- **VERIFY**:
  - Nota registrada em "Status de Execução"; as respostas 3 e 4 mudam o desenho da T6 e da T7 e devem ser claras

### Task T2: Migration dos campos de certificado

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Número da migration conferido (T1)
  - Ordem de merge combinada com os planos de pausa, fuso e índices
- **INPUT**: tabela `monitors` sem informação de certificado
- **OUTPUT**: `00NN_monitor_cert_expiry.sql` com 4 colunas novas em `monitors`:
  1. `cert_alert_enabled` BOOLEAN NOT NULL DEFAULT 0 (o checkbox)
  2. `cert_expires_at` (data de expiração do último certificado visto; aceita NULL)
  3. `cert_alerted_30_for` (guarda **a data de expiração para a qual o aviso de 30 dias já foi enviado**; aceita NULL)
  4. `cert_alerted_7_for` (o mesmo, para o aviso de 7 dias; aceita NULL)
  - Down remove as colunas (conferir se a versão do SQLite do `go-sqlite3` aceita `DROP COLUMN`; senão, usar o padrão de recriar a tabela das migrations anteriores)
  - Motivo do desenho: guardar **a data do certificado alertado**, e não um "sim/não", faz o aviso rearmar sozinho quando o certificado é renovado. Não precisa de lógica de reset
- **VERIFY**:
  - `go test ./internal/database/...` passa, subindo e descendo
  - Banco existente atualiza sem perder dados; monitores antigos ficam com a opção desligada

### Task T3: Checker HTTP captura a data do certificado

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: formato do `Result` e ponto de montagem conhecidos (T1)
- **INPUT**: o checker devolve só status, latência e erro
- **OUTPUT**:
  1. `Result` ganha um campo opcional `CertNotAfter *time.Time`
  2. O checker HTTP preenche o campo quando `resp.TLS != nil` e há certificado: usa `PeerCertificates[0]` (o certificado do próprio site, o primeiro da cadeia). Fica em UTC
  3. Preenche **mesmo quando o status HTTP é de erro** (ex.: 500). O certificado foi lido de qualquer forma e isso mantém a data atualizada durante um incidente
  4. Falha de conexão ou de TLS (inclusive certificado já vencido, que o Go recusa) não traz resposta: o campo fica vazio. Nesse caso o monitor já cai como DOWN pelo fluxo atual
  5. Os outros tipos de checker (TCP e Heartbeat) não mudam
  6. Teste com `httptest.NewTLSServer` usando um certificado gerado no teste com validade conhecida (ex.: daqui a 20 dias), conferindo o valor lido. Teste para URL `http://` (campo vazio) e para resposta 500 (campo preenchido)
- **VERIFY**:
  - `go test ./internal/checker/...` passa
  - Nenhuma mudança de latência ou de comportamento nos checks existentes

### Task T4: Domínio, repositório e API

- **Agent**: `design-api`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**: migration aplicada; ponto de validação conhecido (T1)
- **INPUT**: API de monitores sem os campos novos
- **OUTPUT**:
  1. `domain.Monitor` ganha `CertAlertEnabled`, `CertExpiresAt`, `CertAlerted30For`, `CertAlerted7For`
  2. Repositório lê e grava os campos
  3. **API:**
     - `POST` e `PATCH /api/monitors` aceitam `cert_alert_enabled` (bool, opcional)
     - `GET /api/monitors/:id` devolve `cert_alert_enabled` e `cert_expires_at` (UTC, formato RFC 3339, `null` se ainda não visto). As marcas internas de aviso **não** vão para a API
     - Só a página de detalhes precisa dos campos; a lista `GET /api/monitors` não muda (evita respostas maiores)
  4. **Validação:** `cert_alert_enabled=true` só é aceito se o monitor é do tipo HTTP e a URL começa com `https://`. Caso contrário, `400` com mensagem clara (padrão de erro já usado, com `code`)
  5. **Troca de URL:** ao atualizar um monitor e a URL mudar, zerar `cert_expires_at`, `cert_alerted_30_for` e `cert_alerted_7_for`. O certificado antigo não vale mais para o endereço novo. Se a URL nova não for https, desligar a opção (ou rejeitar a combinação, conforme o padrão de validação da T1)
  6. Trocar o tipo do monitor (se o PATCH permitir) também zera os campos
- **VERIFY**:
  - Testes de handler: criar com a opção, ligar em URL http (400), ligar em heartbeat (400), trocar URL (campos zerados)
  - `curl` manual mostrando os campos novos no `GET /api/monitors/:id`

### Task T5: Regra pura de aviso (30 e 7 dias)

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: nada além do desenho abaixo
- **INPUT**: nenhuma regra de aviso de certificado
- **OUTPUT**: arquivo `service/cert_alert.go` com **uma função pura** (sem banco, sem relógio global, sem envio), por exemplo:
  `EvaluateCertAlert(now, expiresAt, alerted30For, alerted7For) -> (aviso, marcar30, marcar7)`
  Regras:
  1. `restante = expiresAt - now`
  2. Se `restante <= 7 × 24 h` (inclui certificado já vencido) e `alerted7For != expiresAt`: **aviso de 7 dias**; marca **os dois** (30 e 7) como enviados para esta data. Assim, quem liga a opção já na reta final recebe só um aviso, e o de 30 dias não chega atrasado depois
  3. Senão, se `restante <= 30 × 24 h` e `alerted30For != expiresAt`: **aviso de 30 dias**; marca só o de 30
  4. Senão: nada
  5. Se `expiresAt` for diferente do que foi alertado antes (certificado renovado), as marcas antigas não batem e o ciclo recomeça sozinho
  6. O aviso carrega o número de dias restantes (arredondado para baixo) e o limite (30 ou 7), para o texto das mensagens
  - **Testes em tabela:** exatamente 30 dias; 30 dias menos 1 segundo; entre 30 e 7 (já avisado e não avisado); exatamente 7 dias; ativação tardia dentro dos 7 dias; certificado já vencido; renovação (data nova maior); data mais antiga que a alertada (certificado trocado por um mais curto); repetição de 100 avaliações seguidas gerando 1 aviso só
- **VERIFY**:
  - `go test ./internal/service -run CertAlert -v` passa, incluindo todos os casos acima

### Task T6: Notificações: novo evento nos três canais

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: T1
- **PRE-CONDITIONS**: interface do notificador, payload do webhook e formatos de e-mail e ntfy conhecidos (T1)
- **INPUT**: o notificador só sabe avisar queda e recuperação
- **OUTPUT**:
  1. Novo método no notificador, por exemplo `NotifyCertExpiring(monitor, expiresAt, daysLeft, thresholdDays)`, chamado pelos mesmos canais configurados (webhook, e-mail, ntfy)
  2. **Webhook:** novo `event: "cert_expiring"`, no mesmo formato dos eventos atuais, com campos extras:
     ```json
     {
       "event": "cert_expiring",
       "monitor": {
         "id": "...",
         "name": "API",
         "url": "https://api.example.com"
       },
       "timestamp": "2026-11-01T12:00:00Z",
       "cert_expires_at": "2026-11-08T09:00:00Z",
       "days_left": 7,
       "threshold_days": 7
     }
     ```
     Se existir template de webhook configurável, definir o que ele recebe para o evento novo (e o que acontece com variáveis que só fazem sentido para queda). Os eventos `down` e `up` **não mudam**
  3. **E-mail:** assunto e corpo claros, seguindo o idioma e o estilo dos e-mails atuais. Ex.: assunto "Certificado SSL de API expira em 7 dias", com monitor, URL e a data de expiração
  4. **ntfy:** mensagem curta. Prioridade normal no aviso de 30 dias e mais alta no de 7 dias, se o canal já usar prioridade
  5. A data nas mensagens segue o fuso configurado quando o plano de fuso estiver mesclado; antes disso, fica em UTC com o sufixo "UTC"
  6. Envio assíncrono, uma vez por evento, sem retry (igual ao comportamento atual)
  7. Falha de um canal não impede os outros
  8. Testes com notificador falso e com servidor HTTP falso para o webhook, conferindo corpo e campos
- **VERIFY**:
  - `go test` dos pacotes de notificação passa
  - O botão "testar webhook" e os eventos `down` e `up` continuam iguais

### Task T7: Integração na state machine

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: T3, T4, T5, T6
- **PRE-CONDITIONS**: todas as peças anteriores prontas
- **INPUT**: `StateMachine.Process` trata só check, estado e incidente
- **OUTPUT**: dentro da mesma transação de `Process`, depois do check:
  1. Se o `Result` trouxe `CertNotAfter`:
     - se for diferente de `cert_expires_at`, atualizar `cert_expires_at`
  2. Se o monitor tem `cert_alert_enabled` ligado e `cert_expires_at` conhecido, chamar `EvaluateCertAlert`. Se pedir aviso:
     - gravar as marcas `cert_alerted_30_for` e/ou `cert_alerted_7_for` **na mesma transação**
     - guardar o aviso numa lista para enviar **depois do commit** (igual a queda e recuperação)
  3. **Escrita guardada pela URL:** o `UPDATE` dos campos de certificado usa `WHERE id = ? AND url = ?` (a URL com a qual o check foi feito). Se o usuário trocou a URL enquanto o check rodava, a escrita é ignorada e o reset da T4 vale. Isso evita gravar o certificado do endereço antigo no monitor novo
  4. Só escrever no banco quando algo mudou (data nova ou marca nova). Na maioria dos checks, **nenhuma escrita extra**
  5. Sem `CertNotAfter` (falha de conexão, URL http, TCP, Heartbeat): não altera nada. O último valor conhecido continua guardado
  6. Se a transação falhar e der rollback, **nenhum aviso é enviado** e as marcas não ficam gravadas; o próximo check tenta de novo
  7. Monitor pausado não faz check, então não gera aviso (comportamento natural, sem código extra). Registrar isso para o plano de pausa
  8. Log (`slog`) em nível info ao enviar um aviso: monitor, limite, data
- **VERIFY**:
  - Teste com relógio injetado: ver T10
  - Rodar com um monitor real e conferir no log e no banco que não há escrita a cada check

### Task T8: Checkbox no formulário do monitor

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T4
- **PRE-CONDITIONS**: ler `DESIGN.md` (campos de formulário, texto de ajuda)
- **INPUT**: formulário de monitor sem a opção
- **OUTPUT**:
  1. Checkbox visível **só para monitores do tipo HTTP**, com o texto "Avisar quando o certificado SSL estiver perto de expirar" e a ajuda "Envia um aviso 30 dias e outro 7 dias antes de expirar, por todos os canais configurados."
  2. **Desabilitado** quando a URL não começa com `https://`, com a ajuda "Disponível apenas para URLs https://"
  3. Se o usuário estava com a opção marcada e troca a URL para `http://`, a opção desmarca sozinha
  4. Na edição, carrega o valor atual de `cert_alert_enabled`; no envio, manda `cert_alert_enabled`
  5. Padrão desmarcado (suposição 2)
  6. Mostrar o erro `400` da API de forma amigável, se acontecer
- **VERIFY**:
  - Criar e editar um monitor pela tela e conferir o valor no `GET /api/monitors/:id`
  - `cd web && npm run lint && npm run build` passam

### Task T9: Data de expiração na página de detalhes

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T4
- **PRE-CONDITIONS**:
  - `DESIGN.md` lido
  - Saber se o helper de data com fuso do plano de fuso já existe (T1). Se existir, usar. Se não, formatar com o fuso do navegador e trocar quando o helper chegar
- **INPUT**: a página de detalhes não mostra nada sobre certificado
- **OUTPUT**:
  1. Uma linha "Certificado SSL" na página de detalhes do monitor HTTP com URL https, mostrando a data de expiração e os dias restantes. Ex.: "Expira em 23/12/2026 (daqui a 52 dias)"
  2. Cores com os tokens do `DESIGN.md` (nada de cor nova solta):
     - normal: mais de 30 dias
     - atenção: 30 dias ou menos
     - alerta (vermelho): 7 dias ou menos, ou já vencido ("Expirou em ...")
  3. Estados vazios: "Aguardando o primeiro check" quando `cert_expires_at` é `null`
  4. Se a opção de aviso está ligada, um texto discreto "Aviso ativo (30 e 7 dias antes)"
  5. Não aparece para monitor TCP, Heartbeat ou HTTP com URL `http://`
  6. Quando o site está fora do ar, o valor mostrado é o **último conhecido**. Um rótulo discreto "último valor conhecido" aparece enquanto o monitor estiver DOWN
  7. Cálculo de "dias restantes" no navegador a partir de `cert_expires_at`
  8. Layout funciona no celular, sem rolagem horizontal
- **VERIFY**:
  - Conferir visualmente os 4 estados (normal, atenção, alerta, vencido) e o estado vazio, usando valores no banco
  - `cd web && npm run lint && npm run build` passam

### Task T10: Testes de integração e regressão

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T7
- **PRE-CONDITIONS**: relógio injetável e notificador falso disponíveis (T1, T6)
- **INPUT**: peças testadas sozinhas
- **OUTPUT**: testes de `StateMachine.Process` com banco SQLite em memória, relógio falso e notificador falso, cobrindo:
  1. Simulação de 90 dias de checks (1 por hora) de um certificado: **1 aviso de 30 dias e 1 de 7 dias**, nenhum a mais
  2. Opção desligada: nenhum aviso, mas `cert_expires_at` é atualizado
  3. Opção ligada com 5 dias restantes: um único aviso (o de 7 dias)
  4. Renovação no meio do ciclo (data nova): ciclo recomeça e gera avisos de novo, na hora certa
  5. Troca de URL: campos zerados, sem aviso do certificado antigo
  6. Falha de transação (rollback): nenhum aviso enviado, marcas não gravadas
  7. Check sem certificado (falha de conexão) não apaga a data guardada
  8. Monitor em DOWN com certificado visível ainda envia o aviso de certificado
  9. Regressão: transições `UNKNOWN→UP`, `UP→DOWN`, `DOWN→UP` com os mesmos avisos de antes
- **VERIFY**:
  - `make test` verde

### Task T11: QA manual com certificado de validade escolhida

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T7, T8, T9
- **PRE-CONDITIONS**: backend e frontend integrados
- **INPUT**: só testes automáticos; ninguém viu a conexão TLS real funcionando
- **OUTPUT**:
  1. `scripts/certserver/main.go`: servidor HTTPS local que cria uma CA de teste e um certificado de site com **validade escolhida** (`--days 25`) e grava a CA em arquivo. Uso:
     `go run ./scripts/certserver --days 25 --port 9443`
  2. Rodar o Pingou confiando nessa CA. No Linux, o Go respeita `SSL_CERT_FILE=<ca.pem>`. No macOS isso não funciona, então rodar o Pingou em Linux ou Docker
  3. Roteiro manual:
     - criar monitor `https://localhost:9443` com a opção ligada: a data aparece na página de detalhes e o aviso de **30 dias sai na hora** (25 dias já é menos que 30)
     - reiniciar o servidor com `--days 5`: sai o aviso de **7 dias** uma vez; checks seguintes não repetem
     - reiniciar com `--days 90` (renovação): a data atualiza e os avisos rearmam
     - trocar a URL para outro endereço: data zerada
     - desligar a opção: sem aviso novo; a data continua visível
     - conferir webhook, e-mail e ntfy recebendo (com os três configurados)
  4. Registrar capturas de tela ou saídas em "Status de Execução"
- **VERIFY**:
  - Todos os passos do roteiro passam e os resultados estão registrados

### Task T12: Documentação

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T7, T9
- **OUTPUT**:
  1. README: nova funcionalidade "Aviso de expiração do certificado SSL" na lista de funcionalidades, como ligar e quando os avisos saem
  2. README: campos `cert_alert_enabled` e `cert_expires_at` na API de monitores
  3. README: novo evento do webhook `cert_expiring`, com exemplo, e o **aviso para quem consome o webhook**: o campo `event` agora pode ter um terceiro valor; integrações que tratam "tudo que não é `up` como queda" devem ser ajustadas
  4. README: nota sobre certificados com renovação automática (ver Notes) e sobre o que o recurso não cobre (só o certificado do site, só HTTP com https)
  5. `scripts/README.md` (se existir): como usar o `certserver` para testar
  6. `DESIGN.md`: registrar o componente novo, se houver
- **VERIFY**:
  - Copiar e colar os comandos e o exemplo do webhook e conferir que funcionam

## Risk Areas

| Risk                                                                                                                                                                        | Impact                                                                                                     | Mitigation                                                                                                                                                                                                             |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Renovação automática coincide com o aviso de 30 dias.** O Certbot/Let's Encrypt renova quando faltam 30 dias, então o aviso de 30 dias pode sair horas antes da renovação | Aviso "falso" a cada ~60 dias para quem usa renovação automática. Pode irritar e levar a ignorar os avisos | Aceitar no MVP e documentar na T12. O aviso de 7 dias continua sendo o importante. Se incomodar, avaliar uma carência de ~24 h antes do aviso de 30 dias (exige guardar a data em que o limite foi cruzado). Ver Notes |
| Aviso perdido: o envio não tem retry (comportamento atual do projeto)                                                                                                       | Se o e-mail ou webhook falhar naquele momento, o aviso daquele limite se perde                             | Os dois avisos funcionam como redundância (30 e 7). O erro de envio vai para o log. Retry geral fica fora de escopo                                                                                                    |
| Redirect para outro host: o certificado lido é o da resposta final                                                                                                          | A data pode não ser a do host configurado no monitor                                                       | Suposição 3, documentada. O que o visitante vê é o que importa. Capturar o primeiro host fica como melhoria futura                                                                                                     |
| Balanceador ou CDN com certificados diferentes por servidor                                                                                                                 | A data "oscila" entre dois certificados e os avisos podem repetir                                          | Raro. As marcas guardam a data alertada, então só repete quando a oscilação volta para um certificado já avisado depois de outro. Registrado na documentação                                                           |
| Check em andamento grava o certificado de uma URL antiga depois do usuário editar                                                                                           | Data errada no monitor novo                                                                                | `UPDATE` dos campos de certificado só vale se a URL ainda é a mesma (T7) e reset na troca de URL (T4)                                                                                                                  |
| `UpdateWithTx` grava a linha inteira e sobrescreve `cert_alert_enabled` editado durante o check                                                                             | O checkbox "volta" sozinho                                                                                 | T1 confirma como o update funciona. Se gravar a linha inteira, usar um `UPDATE` específico só para os campos de certificado                                                                                            |
| Novo `event` no webhook quebra integrações que assumem só `up` e `down`                                                                                                     | Aviso de certificado vira alerta de queda em quem consome o webhook                                        | Opção desligada por padrão (ninguém recebe sem ligar). Documentar na T12                                                                                                                                               |
| Template de webhook configurável (se existir) não prevê o evento novo                                                                                                       | Mensagem vazia ou quebrada                                                                                 | T1 descobre se existe; T6 define o que o template recebe                                                                                                                                                               |
| Conflito de número de migration com pausa (`0005`), fuso e índices                                                                                                          | Falha ao subir                                                                                             | T1 confere; combinar a ordem de merge antes da T2                                                                                                                                                                      |
| Site fora do ar não devolve certificado                                                                                                                                     | Aviso atrasado ou ausente durante a queda                                                                  | O incidente DOWN já avisa. A data guardada continua válida; a tela marca "último valor conhecido"                                                                                                                      |
| Certificado vencido: o Go recusa o handshake                                                                                                                                | Não dá para ler a data de um certificado já vencido                                                        | Não precisa: o monitor cai como DOWN com o erro do certificado. A regra ainda trata "vencido" por segurança (T5)                                                                                                       |
| Fuso do plano de fuso ainda não mesclado                                                                                                                                    | Datas das mensagens em UTC                                                                                 | T6 e T9 usam UTC ou fuso do navegador até o helper existir                                                                                                                                                             |

## Rollback Strategy

1. **Código:** reverter os commits do plano. O comportamento volta ao atual. As tasks são commits separados (checker, regra, notificações, integração, frontend), então dá para reverter uma parte só.
2. **Migration:** `pingou migrate down` remove as colunas. Nenhum dado existente depende delas. Vale fazer uma cópia do `.db` antes de atualizar.
3. **Corte rápido sem reverter:** como a opção começa desmarcada, nada é enviado até alguém ligar. Para parar os avisos, basta desmarcar o checkbox.
4. **Frontend:** reverter só T8 e T9 esconde a interface sem afetar o backend.
5. **Aviso indevido já enviado:** não tem como "desenviar". As marcas no banco evitam repetição.

## Phase X: Verification

```bash
# Build, vet e testes
go build ./... && go vet ./... && make test

# Testes específicos
go test ./internal/checker/... -run Cert -v
go test ./internal/service/... -run CertAlert -v

# Migration sobe e desce numa cópia do banco
cp pingou.db /tmp/teste.db
PINGOU_DATABASE_URL=/tmp/teste.db PINGOU_API_KEY=dev-api-key ./bin/pingou migrate up --key dev-api-key
PINGOU_DATABASE_URL=/tmp/teste.db PINGOU_API_KEY=dev-api-key ./bin/pingou migrate down --key dev-api-key

# Frontend
cd web && npm run lint && npm run build && cd ..

# Build completo
make build

# QA manual (Linux): certificado de 25 dias, confiando na CA de teste
go run ./scripts/certserver --days 25 --port 9443 &
SSL_CERT_FILE=.bench/ca.pem PINGOU_API_KEY=dev-api-key ./bin/pingou serve
```

> O caminho do arquivo da CA depende de onde o `certserver` gravar. Ajustar conforme a T11.

## Cronograma Sugerido

| Fase   | Tarefas       | Complexity | Esforço     | Target                                                            | Status      |
| ------ | ------------- | ---------- | ----------- | ----------------------------------------------------------------- | ----------- |
| Fase 0 | T1            | S          | ~15 min     | Levantamento do código                                            | ⬜ Pendente |
| Fase 1 | T2, T3, T5    | M          | ~1 h        | Migration, checker e regra pura (em paralelo)                     | ⬜ Pendente |
| Fase 2 | T4, T6        | L          | ~1 h 30 min | API e notificações (em paralelo)                                  | ⬜ Pendente |
| Fase 3 | T7            | L          | ~1 h        | Integração na state machine                                       | ⬜ Pendente |
| Fase 4 | T8, T9        | M          | ~45 min     | Frontend (pode começar após a T4, em paralelo com as fases 2 e 3) | ⬜ Pendente |
| Fase 5 | T10, T11, T12 | M          | ~1 h 15 min | Testes, QA manual e documentação                                  | ⬜ Pendente |

Total estimado: cerca de 6 horas de trabalho. Com o frontend em paralelo ao backend, o caminho mais longo fica em cerca de 4 h 30 min.

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T12

## Notes

- **Ponto de atenção, Let's Encrypt:** o Certbot normalmente renova quando faltam 30 dias, que é exatamente o limite do primeiro aviso. Para quem usa renovação automática, o aviso de 30 dias pode sair pouco antes da renovação. O de 7 dias é o que realmente indica "a renovação falhou". Decisão do dono do projeto foi manter 30 e 7; o plano segue assim e só registra o risco. Se o aviso de 30 dias incomodar na prática, a saída mais simples é uma carência de ~24 h (exige guardar quando o limite foi cruzado). Não faz parte deste plano.
- **Por que guardar a data do certificado e não "sim/não":** `cert_alerted_30_for = <data de expiração alertada>` rearma o aviso sozinho quando o certificado é renovado, pois a data nova não bate mais com a guardada. Evita código de reset e evita bug de "esqueci de zerar".
- **Por que a regra é uma função pura:** assim o ponto mais delicado do plano (quando avisar e quando não) é testado em tabela, sem banco, sem relógio e sem rede.
- **Ativação tardia:** se alguém ligar a opção quando já faltam 20 dias, recebe o aviso de 30 dias na hora (é a primeira informação útil). Se faltam 5 dias, recebe só o de 7. Se desligar e ligar de novo, as marcas permanecem e o mesmo aviso não repete.
- **Sem escrita extra por check:** a data só é gravada quando muda, e as marcas só quando um aviso sai. Importa porque o SQLite abre com uma conexão só (`SetMaxOpenConns(1)`) e há monitores a cada 10 s.
- **Coordenação com outros planos:** migration (pausa `0005`, fuso, índices); formatação de data usa o helper do plano de fuso quando existir (acoplamento leve, não bloqueia); monitor pausado não gera aviso porque não faz check.
- O README do projeto (versão anexada ao Project) está desatualizado em relação ao repositório (e-mail e ntfy existem hoje). A T1 e a T6 confirmam os formatos reais no código.

### Alternativas consideradas e descartadas

- `Conexão TLS separada só para ler o certificado` → descartada: o certificado já vem em `resp.TLS`, sem custo.
- `Gravar um "enviado sim/não" por limite` → descartada: precisaria de lógica de reset na renovação. Guardar a data do certificado alertado resolve sozinho.
- `Gravar data e hora da última verificação do certificado (cert_checked_at) a cada check` → descartada: gera uma escrita por check em cada monitor, com banco de conexão única. A data de expiração só muda quando o certificado muda.
- `Enviar os dois avisos quando a opção é ligada com menos de 7 dias restantes` → descartada: duas mensagens iguais em seguida, sem utilidade.
- `Limites configuráveis (ex.: 14 dias)` → descartada por ora: o pedido fixa 30 e 7, e configurar exige mais campos e telas.
- `Mostrar a data apenas quando o aviso está ligado` → descartada como padrão: ler a data é grátis e útil por si só. Fica como alternativa se o dono discordar (suposição 1).
- `Capturar o certificado do primeiro host antes dos redirects (httptrace)` → descartada: mais código, e o certificado da resposta final é o que o visitante vê.
- `Flag --cert-alert no pingou add` → descartada do escopo: o pedido fala de checkbox. Pode entrar depois.
- `Badge de certificado na lista de monitores` → descartada do escopo: o pedido é a página de detalhes.
