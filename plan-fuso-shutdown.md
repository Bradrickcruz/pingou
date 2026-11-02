# Plano de Desenvolvimento: Fuso horário e graceful shutdown

## Overview

Duas melhorias independentes no Pingou, planejadas juntas para dividir o levantamento do código, os testes, o QA e a documentação.

1. **Fuso horário nas configurações.** Hoje as datas aparecem em UTC (as mensagens de e-mail e ntfy mostram "UTC"). O usuário escolhe um fuso em Configurações. O painel e as mensagens de alerta passam a usar esse fuso. Padrão: `UTC`, igual a hoje. Também prepara a futura "janela de manutenção", que depende de fuso.
2. **Graceful shutdown.** Ao receber SIGTERM/SIGINT, o Pingou termina o que está em andamento (checks, transações da state machine, pings recebidos, notificações) antes de sair. Hoje o desligamento não é ordenado. A única espera conhecida é a das notificações (10 s).

**Decisões já tomadas (pelo dono do projeto):**

- Fuso vale no **painel + e-mail + ntfy**. API e banco continuam sempre em UTC. O JSON do webhook continua com `timestamp` em UTC (RFC3339).
- Fuso padrão: **UTC**.

**Decisões assumidas (revisar se discordar):**

- Orçamento do shutdown: **15 s no total**, com `stop_grace_period` de 20 s no Compose (o padrão do Docker é 10 s e mataria o processo antes).
- Segundo sinal durante o shutdown encerra na hora.

**Fora de escopo:** fuso por usuário, `PINGOU_TIMEZONE` por variável de ambiente, janela de manutenção, mudança dos baldes do `/stats` (continuam alinhados em UTC). O export sem segredos saiu deste plano e será repensado à parte.

**Restrições:** sem dependência nova (só biblioteca padrão do Go), preservar o princípio de escopo do README e o baixo uso de RAM.

## Project Type

**FULLSTACK** — Go (service, handler, scheduler, SQLite + goose) e SPA React em `web/`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e testes passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

**Fuso horário**

- [ ] `GET /api/settings` devolve `timezone` (padrão `"UTC"`); `PATCH` aceita nome IANA válido e responde `400` para inválido
- [ ] Instalações existentes continuam em UTC após a migration, sem nenhuma ação
- [ ] Painel mostra datas e horas no fuso escolhido (dashboard, detalhe do monitor, incidentes, últimos checks)
- [ ] Mensagens de e-mail e ntfy mostram o horário no fuso escolhido, com o sufixo do fuso (ex.: `-03`)
- [ ] Funciona na imagem Docker (sem depender de `tzdata` do sistema)
- [ ] Webhook e API continuam em UTC

**Graceful shutdown**

- [ ] SIGTERM/SIGINT para de aceitar novas conexões e aguarda as requisições em andamento
- [ ] Nenhum check novo começa depois do sinal; checks em andamento são descartados **sem** gerar falha falsa
- [ ] Transações da state machine em andamento terminam (commit ou rollback), nunca ficam pela metade
- [ ] Notificações em andamento são aguardadas dentro do orçamento
- [ ] O banco é fechado (WAL fechado) e o lock `~/.pingou/pingou.lock` é liberado
- [ ] Sai com código `0` quando termina limpo; `1` se estourar o prazo (com log explicando)
- [ ] `docker stop` conclui sem `SIGKILL` em condições normais

## Tech Stack

| Technology | Choice Rationale |
|---|---|
| Go `time` + `time/tzdata` | Biblioteca padrão. `tzdata` embute os fusos no binário (~450 KB), então funciona em imagem mínima sem `/usr/share/zoneinfo` |
| Go `os/signal` + `http.Server.Shutdown` | Biblioteca padrão. `signal.NotifyContext` e `sync.WaitGroup` bastam |
| goose | Migration no padrão das anteriores |
| `Intl.DateTimeFormat` (navegador) | Formata em qualquer fuso IANA sem biblioteca nova |
| React + Tailwind | Seletor de fuso, seguindo o `DESIGN.md` |

## File Structure

Arquivos prováveis. A task T1 confirma os caminhos reais antes de qualquer mudança.

```text
cmd/pingou/
  main.go / serve                        (shutdown ordenado; import _ "time/tzdata")
internal/
  database/migrations/00NN_settings_timezone.sql   (novo; número = próximo livre)
  domain/settings.go                     (+ Timezone)
  service/settings_service.go            (validação do fuso)
  service/notifier.go                    (horário local nas mensagens; espera de envios)
  handler/settings.go                    (+ timezone no JSON)
  scheduler/scheduler.go                 (Stop com espera; descarta resultado cancelado)
  scheduler/retention.go                 (parada ordenada)
web/src/
  lib/datetime.js                        (novo ou existente: formatação com fuso)
  pages/Settings*.jsx                    (card "Fuso horário")
docker-compose.yml                       (stop_grace_period)
README.md                                (documentação)
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
|---|---|---|---|---|---|
| T1 | `go-backend-specialist` | Levantamento do código atual | P0 | S | - |
| T2 | `design-db` | Migration `timezone` em settings | P0 | XS | T1 |
| T3 | `go-backend-specialist` | Settings: domínio, validação e API do fuso | P1 | S | T2 |
| T4 | `go-backend-specialist` | `time/tzdata` + horário local em e-mail e ntfy | P1 | M | T3 |
| T5 | `go-backend-specialist` | Parada ordenada do scheduler e dos jobs | P1 | M | T1 |
| T6 | `go-backend-specialist` | Shutdown ordenado no `serve` | P1 | M | T5 |
| T7 | `devops` | Compose e Dockerfile: `stop_grace_period` | P1 | XS | T6 |
| T8 | `frontend-specialist` | UI: seletor de fuso e datas via `Intl` | P2 | L | T3 |
| T9 | `code-test` | Testes backend | P3 | M | T4, T6 |
| T10 | `code-test` | Testes frontend | P3 | S | T8 |
| T11 | `qa-review` | Verificação manual ponta a ponta | P3 | M | T7, T9, T10 |
| T12 | `document` | README e documentação | P3 | S | T4, T7 |

**Paralelismo:** depois da T1, duas trilhas podem andar ao mesmo tempo: fuso (T2 → T3 → T4 e T8) e shutdown (T5 → T6 → T7). Se for um agente só, a ordem sugerida é shutdown e depois fuso.

### Task T1: Levantamento do código atual

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Repositório atualizado na `main`
  - Ler o plano `plan-pause-closes-incident.md` (ele reserva a migration `0005`)
- **INPUT**: código atual, sem confirmação de como o desligamento e as settings funcionam hoje
- **OUTPUT**:
  1. Nota curta em "Status de Execução" respondendo:
     - Como o `serve` trata SIGTERM/SIGINT hoje? Existe `srv.Shutdown`? Onde está a espera de 10 s das notificações?
     - Como o scheduler, o avaliador de heartbeat e a retenção são iniciados e parados?
     - Como as settings são guardadas (tabela, colunas, singleton)? Qual o último número de migration?
     - Onde fica o formatador de datas do frontend? Existe hook ou contexto de settings?
  2. Lista final de arquivos da seção "File Structure" corrigida
- **VERIFY**:
  - Nota registrada; nenhuma mudança de código nesta task

### Task T2: Migration `timezone`

- **Agent**: `design-db`
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Número da migration definido na T1 (se o plano de pausa ainda não foi mesclado, combinar a ordem para não haver dois `0005`)
- **INPUT**: tabela de settings sem fuso
- **OUTPUT**:
  1. `00NN_settings_timezone.sql`: adiciona `timezone TEXT NOT NULL DEFAULT 'UTC'`; Down com `DROP COLUMN` (mesmo padrão da `0004`)
  2. Linhas existentes ficam com `'UTC'`
- **VERIFY**:
  - `go test ./internal/database/...` passa, up e down

### Task T3: Settings: domínio, validação e API

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**: coluna existe
- **INPUT**: settings sem fuso
- **OUTPUT**:
  1. Campo `Timezone` no domínio e no repositório de settings
  2. Validação com `time.LoadLocation`; vazio vira `"UTC"`; inválido → erro de validação (`400`)
  3. `timezone` em `GET` e `PATCH /api/settings`
  4. Tamanho máximo do texto (ex.: 64) para evitar lixo
- **VERIFY**:
  - `curl -X PATCH /api/settings -d '{"timezone":"America/Sao_Paulo"}'` → `200`; `{"timezone":"Marte/Olimpo"}` → `400`

### Task T4: `time/tzdata` e horário local nas mensagens

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T3
- **PRE-CONDITIONS**: API de settings devolve o fuso
- **INPUT**: e-mail e ntfy formatam o horário em UTC com o texto fixo "UTC"
- **OUTPUT**:
  1. `import _ "time/tzdata"` no `main` (sem isso, `LoadLocation` falha em imagem mínima)
  2. Notifier lê o fuso das settings **no momento do envio** (mudar o fuso vale na próxima notificação, sem reiniciar)
  3. Layout `2006-01-02 15:04:05 MST`: em UTC continua mostrando `UTC`; em `America/Sao_Paulo` mostra `-03`
  4. Se o fuso salvo não carregar, usar UTC e registrar aviso no log (não derrubar a notificação)
  5. Campos `Caiu em` e `Voltou em` usam o fuso; `Fora do ar por` não muda
  6. Webhook **não muda** (continua RFC3339 em UTC)
- **VERIFY**:
  - Teste do notifier com fuso `America/Sao_Paulo`: `14:03:12 UTC` aparece como `11:03:12 -03`
  - Imagem Docker: `docker run` com fuso configurado e botão **Enviar teste** funcionando

### Task T5: Parada ordenada do scheduler e dos jobs

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: forma de iniciar/parar os jobs conhecida (T1)
- **INPUT**: jobs de check, avaliador de heartbeat e retenção sem espera de término
- **OUTPUT**:
  1. `Stop(ctx)` no scheduler: para de agendar, aguarda goroutines em andamento (`sync.WaitGroup`) até o prazo do `ctx`
  2. **Resultado de check cancelado pelo shutdown é descartado**, não vai para a state machine (senão `context canceled` conta como falha e pode abrir incidente falso). Reaproveitar a ideia do descarte já existente para monitor pausado
  3. Avaliador de heartbeat e retenção param no mesmo fluxo; o avaliador não gera `DOWN` durante o desligamento
  4. Transação da state machine em andamento termina antes de a task considerar o job parado
- **VERIFY**:
  - Teste: iniciar check lento, chamar `Stop` → não grava falha; transação em andamento termina

### Task T6: Shutdown ordenado no `serve`

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T5
- **PRE-CONDITIONS**: `Scheduler.Stop` disponível
- **INPUT**: `serve` sem sequência de desligamento definida
- **OUTPUT**:
  1. `signal.NotifyContext` para SIGTERM e SIGINT; orçamento total de **15 s** com `context.WithTimeout`
  2. Sequência:
     1. parar o scheduler (nenhum check novo)
     2. `srv.Shutdown(ctx)`: para de aceitar conexões e termina as requisições em andamento (inclui `/ping/{token}`)
     3. aguardar checks em andamento (T5)
     4. aguardar notificações em andamento (a espera de 10 s que já existe passa a usar o mesmo orçamento)
     5. fechar o banco (`db.Close()`, fecha o WAL)
     6. liberar o lock `~/.pingou/pingou.lock`
  3. Segundo sinal durante o shutdown → sai na hora com código `1`
  4. Logs: `shutdown started`, uma linha por etapa e `shutdown complete` (ou o motivo do estouro)
  5. Código de saída `0` se limpo, `1` se estourou o prazo
- **VERIFY**:
  - Subir com `make dev`, criar monitores, `kill -TERM <pid>`: logs na ordem esperada, lock removido, `echo $?` → `0`

### Task T7: Compose e Dockerfile

- **Agent**: `devops`
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T6
- **PRE-CONDITIONS**: orçamento de 15 s definido
- **INPUT**: Compose sem `stop_grace_period` (padrão do Docker: 10 s)
- **OUTPUT**:
  1. `stop_grace_period: 20s` no `docker-compose.yml`
  2. Confirmar que o processo do container é o próprio binário (PID 1 recebe o SIGTERM); se o `ENTRYPOINT` estiver em formato shell, trocar para formato exec
- **VERIFY**:
  - `docker compose up -d && time docker compose stop` termina em menos de 20 s, sem `Killed`

### Task T8: UI: seletor de fuso e datas

- **Agent**: `frontend-specialist`
- **Priority**: P2
- **Complexity**: L
- **Dependencies**: T3
- **PRE-CONDITIONS**: ler `DESIGN.md` (card, select, texto de ajuda)
- **INPUT**: datas formatadas em UTC ou no fuso do navegador, sem controle do usuário
- **OUTPUT**:
  1. Card **Fuso horário** em Configurações: lista de fusos com busca simples, `UTC` no topo, texto de ajuda ("As datas do painel e dos alertas usam este fuso")
  2. Lista via `Intl.supportedValuesOf('timeZone')`, com `UTC` adicionado à mão e uma lista fixa curta de reserva se o navegador não tiver a função
  3. Helper único de formatação (`Intl.DateTimeFormat` com `timeZone`) usado em todas as telas que mostram data/hora: dashboard, detalhe do monitor, gráficos, incidentes, últimos checks
  4. Fuso vem das settings (hook/contexto existente, se houver); troca vale sem recarregar a página
  5. Eixos dos gráficos usam o fuso (os baldes continuam alinhados em UTC; ver Risk Areas)
- **VERIFY**:
  - `cd web && npm run lint && npm run build`
  - No navegador: trocar para `Asia/Tokyo` e ver todas as telas mudarem de horário

### Task T9: Testes backend

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T4, T6
- **OUTPUT**:
  1. Settings: fuso válido, inválido, vazio; migration up/down
  2. Notifier: horário convertido; fuso inválido cai em UTC sem quebrar o envio
  3. Scheduler: `Stop` aguarda check em andamento; resultado cancelado não vira falha; sem incidente falso
  4. Shutdown: teste de integração que sobe o servidor, envia sinal e confere ordem e código de saída
- **VERIFY**:
  - `make test` verde

### Task T10: Testes frontend

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T8
- **OUTPUT**:
  1. Helper de formatação: UTC → `America/Sao_Paulo` e `Asia/Tokyo`, incluindo virada de dia
  2. Seletor de fuso salva e reflete a mudança
- **VERIFY**:
  - `cd web && npm test`

### Task T11: QA manual

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T7, T9, T10
- **OUTPUT**: checklist executado, com evidências em "Status de Execução":
  - Fuso: trocar para `America/Sao_Paulo`, provocar uma queda e conferir o horário no painel, no e-mail (Mailpit) e no ntfy
  - Fuso: instalação antiga continua em UTC após subir a migration
  - Shutdown: `docker compose stop` durante checks ativos → nenhum incidente novo aparece ao subir de novo
  - Shutdown: `kill -TERM` e depois `kill -9` em outra execução → o lock não impede subir de novo
  - Shutdown: heartbeat recebendo ping durante o desligamento → ping respondido ou falha de conexão limpa, sem corromper o banco
- **VERIFY**: evidências anotadas

### Task T12: Documentação

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T4, T7
- **OUTPUT**:
  1. README: seção sobre fuso horário (onde configurar, padrão UTC, efeito nas mensagens, webhook continua UTC)
  2. README: seção "Desligamento" (sinais, prazo de 15 s, `stop_grace_period`, `--stop-timeout` no `docker run`)
  3. Registrar no `DESIGN.md` o componente novo do seletor, se houver
- **VERIFY**:
  - Revisão de links e comandos copiando e colando

## Risk Areas

| Risk | Impact | Mitigation |
|---|---|---|
| Shutdown cancela um check no meio e o `context canceled` conta como falha | Incidente falso ao reiniciar | T5 descarta resultado cancelado pelo desligamento; critério de sucesso e teste dedicados |
| `docker stop` mata o processo em 10 s antes do fim do shutdown | Shutdown pela metade | Orçamento de 15 s + `stop_grace_period: 20s` (T7); documentar `--stop-timeout` |
| Imagem mínima sem `tzdata` | `LoadLocation` falha e o fuso não funciona | `import _ "time/tzdata"` (T4) e teste na imagem Docker |
| Baldes do `/stats` ficam alinhados em UTC | Em fusos com meia hora (ex.: `Asia/Kolkata`) os rótulos dos baldes de 1 h não caem em hora cheia | Aceitar nesta versão e registrar nas Notes; revisar junto com a janela de manutenção |
| Dois planos reservam o mesmo número de migration | Falha ao subir | T1 confere o último número; combinar a ordem com `plan-pause-closes-incident.md` |
| Mudar o fuso com notificações em andamento | Mensagem sai com o fuso antigo | Aceitável: vale a partir da próxima notificação |
| Datas formatadas em vários lugares do frontend | Alguma tela esquecida continua em UTC | Helper único (T8) e busca por `toLocale`/`new Date(` antes de fechar a task |

## Rollback Strategy

1. **Fuso:** reverter o commit; `goose down` remove a coluna `timezone` (SQLite ≥ 3.35 suporta `DROP COLUMN`, como a `0004` já usa). As mensagens voltam a UTC; nada se perde.
2. **Shutdown:** reverter o commit. O comportamento volta ao desligamento anterior. O `stop_grace_period` maior no Compose não causa dano.
3. Cada trilha tem commits próprios, então qualquer uma pode ser revertida sem a outra.

## Phase X: Verification

```bash
# Backend
go build ./... && go vet ./... && make test

# Frontend
cd web && npm run lint && npm run build && npm test

# Build completo e imagem
make build && make docker-build

# Shutdown (manual)
docker compose up -d && time docker compose stop
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 0 | T1 | S | ~15 min | Levantamento do código | ⬜ Pendente |
| Fase 1 | T5, T6, T7 | M | ~1 h 5 min | Graceful shutdown completo | ⬜ Pendente |
| Fase 2 | T2, T3, T4 | M | ~50 min | Fuso no backend e nas mensagens | ⬜ Pendente |
| Fase 3 | T8 | L | ~1 h | UI do fuso | ⬜ Pendente |
| Fase 4 | T9, T10, T11, T12 | M | ~1 h 45 min | Testes, QA e documentação | ⬜ Pendente |

Total estimado: cerca de 5 horas de trabalho. A ordem das fases 1 e 2 pode mudar, porque as trilhas são independentes depois da T1.

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T12

## Notes

- O fuso é uma configuração global do Pingou (uma por instalação), igual ao webhook. Fuso por pessoa não faz sentido com uma única API Key.
- Este plano prepara a futura **janela de manutenção**, que vai usar o fuso salvo aqui.
- Se na T1 descobrirmos que parte do shutdown já existe, a T5 e a T6 diminuem de tamanho. Ajustar o cronograma.
- O export sem segredos foi retirado deste plano. Quando for retomado, vai precisar de um plano próprio.

### Alternativas consideradas e descartadas

- `Fuso por variável de ambiente (PINGOU_TIMEZONE)` → descartada: o pedido é configurar nas Configurações, e a variável exigiria reiniciar para mudar.
- `Fuso do navegador como padrão` → descartada: muda o comportamento de quem já usa e faz e-mail e painel mostrarem horários diferentes.
- `os.Exit logo após o timeout, sem log` → descartada: dificulta descobrir por que o desligamento estourou.
- `Esperar mais de 15 s no shutdown` → descartada: o Docker mataria o processo antes, e um prazo longo atrasa deploys.