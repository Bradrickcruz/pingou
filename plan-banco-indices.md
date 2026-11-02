# Plano de Desenvolvimento: Otimizar o banco com índices

## Overview

O teste de memória do README rodou **100 monitores a cada 10 s durante 1 hora**. Com essa mesma carga, 30 dias geram cerca de **26 milhões de linhas** na tabela `checks` (100 × 8.640 por dia × 30). Ninguém mediu como as consultas se comportam nesse volume: estatísticas (`/stats`), lista do painel, histórico de checks, incidentes e limpeza de retenção.

Há um agravante: o SQLite abre com `SetMaxOpenConns(1)`. Com uma conexão só, **uma consulta lenta atrasa também a gravação dos checks**. Então o custo de uma consulta ruim não é só o painel lento.

**Objetivo:** medir primeiro, criar só os índices (e pequenos ajustes de consulta) que os números justificarem, e medir de novo. Sem mudar nenhum comportamento visível.

**Fluxo:** medir → decidir → implementar → medir de novo. O plano pode terminar na T5, se a medição mostrar que o banco já está bom.

**Perfis de teste (assumidos):**

| Perfil | Monitores | Intervalo | Checks em 30 dias |
|---|---|---|---|
| típico | 10 | 60 s | ≈ 430 mil |
| limite | 100 | 10 s | ≈ 26 milhões |

O perfil limite é o mesmo cenário do teste de memória do README.

**Decisões já tomadas (pelo dono do projeto):**

- Não há decisão de regra de negócio neste plano.
- A primeira medição é o `/stats` com 30 dias de dados.

**Decisões assumidas (revisar se discordar):**

- Dois perfis de dados (típico e limite), gerados de forma determinística.
- Medir o código real (repositórios e serviços), não cópias do SQL.
- Alvos numéricos da seção Success Criteria. Eles são suposições e a T5 os confirma depois de ver a linha de base.

**Fora de escopo:** trocar o pool de conexões ou criar pool de leitura, tabelas de agregados (rollup), mudar a retenção, paginação por cursor, cache de estatísticas, mudar endpoints ou formato das respostas.

**Restrições:** sem dependência nova, respostas da API idênticas, pico de RAM continua abaixo de 50 MB (meta do README), migrations pelo `goose`.

## Project Type

**BACKEND** — Go e SQLite (migrations, repositórios, retenção), com ferramentas de medição em `scripts/`. Nenhuma mudança no frontend.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

**Medição**

- [ ] Existe um gerador de dados que cria os perfis típico e limite de forma repetível
- [ ] Linha de base registrada (antes dos índices) para os dois perfis, começando pelo `/stats` nos 4 períodos, com 30 dias de dados
- [ ] Planos de consulta (`EXPLAIN QUERY PLAN`) de todas as consultas quentes registrados antes e depois

**Resultado (alvos assumidos, no perfil limite; a T5 confirma)**

- [ ] Nenhuma consulta quente faz varredura completa de `checks` (`SCAN checks`)
- [ ] `/stats?range=30d`: p95 ≤ 500 ms
- [ ] `GET /api/monitors` (lista com `history`, `uptime_24h` e `last_latency_ms`): p95 ≤ 200 ms
- [ ] `GET /api/monitors/:id/checks` (limit 50, offset 0): p95 ≤ 50 ms
- [ ] Gravação de check durante uma leitura pesada de `/stats`: p95 ≤ 1 s
- [ ] Limpeza de retenção não segura a conexão por mais de 2 s seguidos

**Custos**

- [ ] Custo de gravação por check sobe no máximo 20 % em relação à linha de base
- [ ] Tamanho do arquivo antes e depois, e tamanho de cada índice, registrados
- [ ] Tempo da migration em um banco do perfil limite medido e documentado
- [ ] Pico de RAM continua abaixo de 50 MB, inclusive com o banco do perfil limite aberto e `/stats?range=30d` sendo consultado

**Qualidade**

- [ ] Banco existente atualiza sem perda de dados; `migrate down` remove os índices
- [ ] Respostas da API idênticas antes e depois (teste de equivalência)
- [ ] Teste automático falha se uma consulta quente voltar a varrer `checks`
- [ ] `make test` e `make build` passam, e `make test` continua rápido (benchmarks só rodam com variável de ambiente)

## Tech Stack

| Technology | Choice Rationale |
|---|---|
| SQLite (`EXPLAIN QUERY PLAN`, índices parciais, `ANALYZE`/`PRAGMA optimize`) | Já é o banco do projeto. O plano de consulta mostra se o índice é usado |
| Go `testing` (benchmarks) | Biblioteca padrão. Mede o código real, incluindo o efeito da conexão única |
| Go + `go-sqlite3` no gerador de dados | Dependência que já existe. Evita exigir o `sqlite3` instalado na máquina |
| `goose` | Migration no padrão das anteriores |
| `scripts/memtest.sh` | Já existe. Reaproveitado para conferir a meta de RAM |

## File Structure

Arquivos prováveis. A task T1 confirma os caminhos reais.

```text
scripts/
  seed/main.go                          (novo: gerador de dados sintéticos)
  dbbench.sh                            (novo: gera dados, roda benchmarks e planos)
  README.md                             (+ seção de benchmark do banco)
internal/
  database/migrations/00NN_checks_indexes.sql   (novo; número = próximo livre)
  database/database.go                  (+ log ao aplicar migrations pendentes)
  repository/*_bench_test.go            (novo: benchmarks; pulam sem PINGOU_BENCH_DB)
  repository/*_plan_test.go             (novo: teste de plano de consulta)
  repository/*.go                       (ajustes de consulta só se a T5 indicar)
  scheduler/retention.go                (limpeza em lotes, só se a T5 indicar)
.bench/                                 (dados gerados; fora do git e do Docker)
.gitignore / .dockerignore              (+ .bench/)
Dockerfile                              (HEALTHCHECK --start-period, só se a T8 indicar)
README.md                               (documentação)
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
|---|---|---|---|---|---|
| T1 | `go-backend-specialist` | Levantamento: consultas quentes e índices atuais | P0 | S | - |
| T2 | `performance` | Gerador de dados sintéticos (perfis típico e limite) | P0 | M | T1 |
| T3 | `go-backend-specialist` | Ferramentas de medição (benchmarks e planos) | P0 | M | T1 |
| T4 | `performance` | Linha de base antes dos índices | P0 | M | T2, T3 |
| T5 | `performance` | Decisão: quais índices e ajustes valem a pena | P1 | S | T4 |
| T6 | `design-db` | Migration dos índices | P1 | S | T5 |
| T7 | `go-backend-specialist` | Ajustes de código (log de migration, consultas, retenção) | P1 | M | T5 |
| T8 | `performance` | Medição depois e comparação | P1 | M | T6, T7 |
| T9 | `code-test` | Testes de regressão | P3 | M | T6, T7 |
| T10 | `document` | Documentação e operação | P3 | S | T8 |

**Sobre a ordem:** a primeira medição pedida (o `/stats` com 30 dias) é a T4. As tasks T1 a T3 existem porque não dá para medir sem saber quais são as consultas, sem dados e sem ferramenta. T2 e T3 podem andar ao mesmo tempo.

### Task T1: Levantamento: consultas quentes e índices atuais

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Repositório atualizado na `main`
  - Ler `plan-pause-closes-incident.md` (reserva a migration `0005`) e `plan-fuso-shutdown.md` (também cria uma migration)
- **INPUT**: schema e consultas conhecidos só pelo README
- **OUTPUT**:
  1. Tabela de **consultas quentes**, com: nome, arquivo e linha, tabelas, filtros, ordenação, quem chama e com que frequência. Cobrir pelo menos:
     - `/stats`: totais (uptime, latência média, contagem) e série por balde
     - `/stats`: janela `down` (incidentes que se sobrepõem ao período)
     - `GET /api/monitors`: resumo por monitor (`history` dos últimos 30 checks, `uptime_24h`, `last_latency_ms`, `down_since`)
     - `GET /api/monitors/:id/checks`: paginado (`limit`, `offset`)
     - `GET /api/incidents` e `/api/monitors/:id/incidents`: lista e total, com `status` e `monitor_id`
     - retenção: `DELETE` por data
     - avaliador de heartbeat (a cada 10 s), carga do scheduler e busca por `ping_token` em `/ping/{token}`
  2. Índices que já existem em cada tabela (migrations `0001` em diante)
  3. Tipo das colunas de data (texto ISO ou inteiro), de sucesso e de latência, e como o balde do `/stats` é calculado (expressão sobre a data?)
  4. Último número de migration e quais números os outros planos já reservaram
  5. De quanto em quanto tempo o painel chama `GET /api/monitors` (isso define o peso da consulta do resumo)
- **VERIFY**:
  - Nota registrada em "Status de Execução"; nenhuma mudança de código

### Task T2: Gerador de dados sintéticos

- **Agent**: `performance`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: nomes e tipos das colunas conhecidos (T1)
- **INPUT**: não existe como testar com 30 dias de dados
- **OUTPUT**:
  1. `scripts/seed/main.go`, chamado com `go run ./scripts/seed --profile tipico|limite --out .bench/<perfil>.db`
  2. O banco nasce pelas **migrations reais** (mesmo caminho do app) e recebe os dados com `INSERT ... SELECT` usando CTE recursiva (rápido, sem uma linha por chamada)
  3. **Determinístico:** os valores vêm de uma função de hash do monitor e do instante, não de `random()`. Duas execuções geram o mesmo banco
  4. **Ordem cronológica, com os monitores intercalados.** Isso importa: se os dados saírem agrupados por monitor, as consultas ficam artificialmente rápidas, porque o layout em disco não seria o real
  5. Conteúdo realista:
     - latência entre ~20 e ~400 ms, com variação
     - cerca de 0,3 % de falhas isoladas
     - uma queda de 5 a 30 minutos por monitor por semana (checks com falha seguidos, incidente aberto e fechado, estado do monitor coerente)
     - cerca de 5 mil incidentes no total, incluindo os mais antigos que a retenção (incidentes nunca expiram)
     - alguns monitores Heartbeat, com poucos checks e sem latência
     - `retention_days = 30` nas settings
  6. Durante a geração, desligar o que acelera a carga (`synchronous`, `journal_mode`) e **restaurar** o padrão do app (WAL) no fim
  7. Resumo impresso: linhas por tabela, tamanho do arquivo e tempo gasto
  8. `.bench/` no `.gitignore` e no `.dockerignore`
- **VERIFY**:
  - Gerar o perfil típico, conferir `COUNT(*)` (≈ 430 mil checks), repetir e comparar contagens
  - Subir o app apontando para o banco gerado e ver o painel carregar

### Task T3: Ferramentas de medição

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**: lista de consultas quentes pronta (T1)
- **INPUT**: nenhuma forma repetível de medir tempo e plano das consultas
- **OUTPUT**:
  1. **Benchmarks em Go** que abrem uma **cópia** do banco indicado em `PINGOU_BENCH_DB` (nunca o original) e chamam as funções reais de repositório e serviço:
     - `/stats` nos 4 períodos, em 3 monitores de tipos diferentes
     - resumo da lista de monitores
     - checks paginado (offset 0, 1.000 e 10.000)
     - incidentes: abertos, resolvidos e total
     - retenção, apagando 1 dia de dados na cópia
  2. Sem `PINGOU_BENCH_DB`, os testes são pulados (`t.Skip`), então `make test` não fica lento
  3. **Planos de consulta:** um teste ou comando que imprime o `EXPLAIN QUERY PLAN` de cada consulta quente usando as **mesmas strings SQL do código** (exportar constantes se for preciso, sem mudar comportamento)
  4. **Cenário de concorrência:** uma rotina repete `/stats` 30d enquanto outra grava checks pelo caminho real da state machine a 10 por segundo. Registrar p50, p95 e máximo da latência de gravação
  5. **Custo de gravação:** tempo médio para gravar 50 mil checks pelo caminho real
  6. Primeira execução separada das repetições (cache frio e cache quente)
  7. Saída em tabela de texto, fácil de colar no plano
  8. `scripts/dbbench.sh`: gera os dados se faltarem, roda tudo e salva em `.bench/results-<data>.md`
- **VERIFY**:
  - Rodar no perfil típico com saída completa
  - `make test` sem a variável continua verde e rápido

### Task T4: Linha de base antes dos índices

- **Agent**: `performance`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T2, T3
- **PRE-CONDITIONS**:
  - Bancos dos dois perfis gerados
  - Máquina sem outra carga pesada. A mesma máquina será usada na T8
- **INPUT**: nenhum número de desempenho com volume real
- **OUTPUT**:
  1. **Primeiro item: `/stats` nos 4 períodos com 30 dias de dados**, nos dois perfis (p50, p95, primeira execução)
  2. Demais consultas quentes, nos dois perfis
  3. Plano de cada consulta, marcando as que fazem `SCAN`
  4. Latência de gravação com leitura pesada ao lado (cenário de concorrência)
  5. Tempo da limpeza de retenção e por quanto tempo ela segura a conexão
  6. Tamanho do arquivo (`.db` e `-wal`) e por tabela e índice
  7. RAM de pico (RSS) com o banco do perfil limite aberto durante `/stats?range=30d`
  8. Tudo registrado em "Status de Execução" e em `.bench/results-baseline.md`, com máquina, CPU, disco e versão do SQLite
- **VERIFY**:
  - 3 execuções, com variação menor que 15 %. Se passar disso, a máquina está ruidosa e a medição não vale

### Task T5: Decisão: quais índices e ajustes valem a pena

- **Agent**: `performance`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **PRE-CONDITIONS**: linha de base completa
- **INPUT**: números da T4
- **OUTPUT**:
  1. Para cada candidato: o índice, a consulta que o justifica (com plano e tempo atuais), o custo (tamanho e tempo de criação medidos **criando o índice numa cópia do perfil limite**) e a decisão sim ou não. Candidatos iniciais, que são hipóteses a confirmar:
     - `checks(monitor_id, checked_at)`: base do `/stats`, da lista e da paginação (se já não existir)
     - `checks(checked_at)`: a retenção apagar sem varrer a tabela
     - variante que também cobre sucesso e latência: **só** se o `/stats` 30d continuar acima do alvo, pois o índice fica bem maior
     - `incidents(monitor_id, started_at)` e índice parcial para incidentes abertos (`WHERE ended_at IS NULL`)
     - confirmar que `monitors(ping_token)` é único (a tabela é pequena, valor baixo)
  2. Alvos numéricos confirmados ou ajustados
  3. Decisão sobre a T7: precisa reescrever consultas ou fazer a retenção em lotes?
  4. **Se nenhuma consulta falha nos alvos e não há `SCAN` em `checks`, encerrar o plano aqui** e só documentar o resultado
- **VERIFY**:
  - Cada índice aprovado tem um número antes e depois (no protótipo) que o justifica

### Task T6: Migration dos índices

- **Agent**: `design-db`
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T5
- **PRE-CONDITIONS**:
  - Número da migration conferido (T1)
  - Ordem de merge combinada com `plan-pause-closes-incident.md` (`0005`) e `plan-fuso-shutdown.md`
- **INPUT**: lista de índices aprovados
- **OUTPUT**:
  1. `00NN_checks_indexes.sql` com `CREATE INDEX IF NOT EXISTS` para cada índice aprovado
  2. Comentário em cada índice dizendo qual consulta o justifica
  3. Down com `DROP INDEX IF EXISTS`
  4. `ANALYZE` (ou `PRAGMA optimize`) no fim, para o planejador usar os índices. Decidir pelo custo medido na T5
- **VERIFY**:
  - `go test ./internal/database/...` passa, subindo e descendo
  - Aplicar numa cópia do perfil limite e anotar tempo e tamanho

### Task T7: Ajustes de código

- **Agent**: `go-backend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T5
- **PRE-CONDITIONS**: decisões da T5
- **INPUT**: consultas e retenção como estão hoje
- **OUTPUT**:
  1. **Sempre:** log claro ao aplicar migrations pendentes (quais versões e quanto tempo levou), para quem atualiza um banco grande não achar que o app travou
  2. **Só se a T5 indicar:** reescrever consultas para o índice ser usável, sem mudar o resultado. Exemplos de causa (hipóteses): função aplicada sobre a coluna de data no filtro, `ORDER BY` diferente da ordem do índice, `COUNT` do total calculado quando não é pedido
  3. **Só se a T5 indicar:** limpeza de retenção em lotes (um bloco de linhas por transação), para não segurar a única conexão por segundos. O resultado final é o mesmo conjunto apagado. Checar o `ctx` entre lotes, para parar rápido num desligamento
- **VERIFY**:
  - Testes existentes passam
  - Teste de equivalência: mesmas respostas antes e depois, num banco pequeno determinístico
  - Benchmark da T3 mostra ganho ou neutralidade

### Task T8: Medição depois e comparação

- **Agent**: `performance`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T6, T7
- **PRE-CONDITIONS**: mesma máquina e mesmos perfis da T4
- **INPUT**: código e migration prontos
- **OUTPUT**:
  1. Tabela comparativa antes, depois e variação, nos dois perfis, com o `/stats` 30d em primeiro lugar
  2. Planos de consulta depois (sem `SCAN checks`)
  3. Tempo da migration num banco do perfil limite que ainda não tem os índices
  4. Crescimento do arquivo por índice
  5. Custo de gravação por check e latência de gravação com leitura pesada ao lado
  6. RAM de pico com o banco limite aberto. Rodar também `HTTPS_N=50 scripts/memtest.sh 3600` para conferir a meta de 50 MB
  7. Cada alvo marcado como cumprido ou com o desvio justificado. Se algum falhar, voltar à T5
- **VERIFY**:
  - Tabela preenchida e registrada em "Status de Execução"

### Task T9: Testes de regressão

- **Agent**: `code-test`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T6, T7
- **OUTPUT**:
  1. Migration sobe e desce sem perder dados
  2. **Teste de plano de consulta** num banco pequeno: roda `EXPLAIN QUERY PLAN` das consultas quentes e falha se o índice esperado não aparecer. Verificar pelo **nome do índice**, não só pela ausência de `SCAN`, porque com poucos dados o planejador pode escolher diferente entre versões do SQLite. Rodar `ANALYZE` antes, para o teste ser estável
  3. Equivalência de resultados: `/stats`, lista e incidentes iguais antes e depois
  4. Retenção em lotes (se existir): apaga exatamente o mesmo conjunto que a versão antiga e respeita o cancelamento
- **VERIFY**:
  - `make test` verde

### Task T10: Documentação e operação

- **Agent**: `document`
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T8
- **OUTPUT**:
  1. README: seção "Desempenho e tamanho do banco", com o que foi medido (perfil, máquina, antes e depois) e o aviso de que o resultado depende do hardware, como já é feito na seção de memória
  2. README: aviso para quem atualiza: a primeira subida pode demorar no banco grande, quanto espaço livre é preciso (aproximadamente o tamanho dos novos índices) e como acompanhar pelo log
  3. `scripts/README.md`: como gerar os dados e rodar os benchmarks (`scripts/dbbench.sh`)
  4. Se a migration no perfil limite passar de 30 s: ajustar o `HEALTHCHECK --start-period` no Dockerfile (agent `devops`) e documentar. Sem isso, o Docker pode marcar o container como `unhealthy` no meio da atualização
  5. Atualizar a tabela de variáveis ou comandos do README se alguma mudança exigir
- **VERIFY**:
  - Copiar e colar os comandos do README e conferir que funcionam

## Risk Areas

| Risk | Impact | Mitigation |
|---|---|---|
| Dados de teste agrupados por monitor deixam as consultas artificialmente rápidas | Conclusão errada: parece não precisar de índice | T2 gera em ordem cronológica com monitores intercalados, igual ao app real |
| `CREATE INDEX` em banco com milhões de linhas demora e usa espaço temporário | Subida lenta após atualizar; healthcheck do Docker pode derrubar o container | Medir o tempo na T8, log antes de aplicar (T7), `--start-period` se passar de 30 s (T10), documentar espaço livre |
| Índice grande aumenta muito o arquivo (no limite, na ordem de centenas de MB a GB, a medir) | Disco de VPS pequena enche | Medir cada índice na T5; preferir índices estreitos; índice "cobrindo" só se o ganho justificar |
| Mais índices deixam a gravação mais lenta | Atraso nos checks | Critério de no máximo 20 % a mais, medido na T8 |
| Pool com uma conexão: leitura lenta bloqueia gravação | Checks atrasados ou perdidos durante o painel aberto | Medir o cenário de concorrência (T3, T4, T8); trocar o pool está fora do escopo e só é revisto se os índices não bastarem |
| Planejador não usa o índice sem estatísticas | Índice criado e ignorado | `ANALYZE` ou `PRAGMA optimize` na migration (T6); teste de plano com `ANALYZE` (T9) |
| Três planos pedem migration ao mesmo tempo (pausa, fuso e este) | Dois arquivos com o mesmo número; falha ao subir | T1 confere o último número; combinar a ordem de merge antes da T6 |
| Reescrever consulta muda o resultado | Painel mostra número diferente | T7 sem mudança de comportamento; teste de equivalência (T9) |
| Resultado depende do hardware | Números absolutos viram promessa | Comparar antes e depois só na mesma máquina; documentar como medição, não garantia |
| Índices aumentam o uso de cache e RAM | Estoura a meta de 50 MB | Conferir RSS e rodar `scripts/memtest.sh` na T8 |
| Paginação com `offset` grande continua lenta mesmo com índice | Página 1.000 do histórico demora | Fora do escopo (cursor); registrar o tempo na T4 e decidir depois |

## Rollback Strategy

1. **Índices:** `goose down` remove os índices (`DROP INDEX IF EXISTS`), de forma rápida e sem perder dados. Os índices só aceleram consultas, não guardam informação própria. O arquivo `.db` só encolhe de volta depois de um `VACUUM`.
2. **Código:** reverter o commit da T7. O log de migration, o ajuste de consultas e a retenção em lotes podem ser commits separados e revertidos um a um.
3. **Ferramentas de medição:** não afetam o app em produção. Podem ficar ou ser removidas sem efeito.
4. Se uma atualização em banco grande estiver demorando demais, dá para parar o app e restaurar o backup do arquivo `.db` feito antes da atualização.

## Phase X: Verification

```bash
# Build, lint e testes
go build ./... && go vet ./... && make test

# Gerar os dados (o perfil limite leva vários minutos)
go run ./scripts/seed --profile tipico --out .bench/tipico.db
go run ./scripts/seed --profile limite --out .bench/limite.db

# Benchmarks e planos de consulta
PINGOU_BENCH_DB=.bench/limite.db go test ./internal/repository -run '^$' -bench . -benchtime 20x
scripts/dbbench.sh

# Tempo da migration num banco grande (cópia do banco gerado ANTES dos índices)
cp .bench/limite-antes.db /tmp/teste.db
time PINGOU_DATABASE_URL=/tmp/teste.db PINGOU_API_KEY=dev-api-key ./bin/pingou migrate up --key dev-api-key

# Meta de memória do README
make build && HTTPS_N=50 scripts/memtest.sh 3600
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 0 | T1 | S | ~15 min | Levantamento das consultas | ⬜ Pendente |
| Fase 1 | T2, T3 | M | ~1 h | Dados e ferramentas de medição | ⬜ Pendente |
| Fase 2 | T4 | M | ~30 min (+ tempo de máquina) | Linha de base, começando pelo `/stats` 30d | ⬜ Pendente |
| Fase 3 | T5 | S | ~15 min | Decisão (pode encerrar o plano) | ⬜ Pendente |
| Fase 4 | T6, T7 | M | ~45 min | Índices e ajustes | ⬜ Pendente |
| Fase 5 | T8, T9, T10 | M | ~1 h 15 min (+ 1 h do memtest) | Comparação, testes e documentação | ⬜ Pendente |

Total estimado: cerca de 4 horas de trabalho, fora o tempo de máquina (geração do perfil limite e memtest de 1 hora). Se a T5 encerrar o plano, cai para cerca de 2 horas.

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T10

## Notes

- A conta do volume: 100 monitores × 8.640 checks por dia × 30 dias ≈ 26 milhões de linhas. O teste de memória do README cobre só 1 hora de dados, então o comportamento em 30 dias ainda não foi medido.
- O pedido era começar medindo o `/stats` com 30 dias. A medição é a T4. Antes dela, T1 a T3 preparam o que é preciso para medir (lista de consultas, dados e ferramenta).
- **Possível resultado "não precisa":** se a T4 mostrar que os índices que já existem bastam, o plano termina na T5 com um registro dos números. Isso também é um resultado válido e evita custo de disco e de gravação à toa.
- Os índices listados na T5 são hipóteses. Quem decide são os números da linha de base.
- Este é o terceiro plano que pode criar migration (junto com `plan-pause-closes-incident.md` e `plan-fuso-shutdown.md`). Combinar a ordem de merge antes da T6.
- A retenção em lotes (T7) combina com o shutdown ordenado do plano de fuso, mas **não depende** dele.

### Alternativas consideradas e descartadas

- `Criar os índices "por precaução", sem medir` → descartada: cada índice custa disco e tempo de gravação, e sem número não dá para saber se vale.
- `Medir só com banco pequeno` → descartada: o problema só aparece com volume, e o README mediu apenas 1 hora.
- `Gerar os dados com o sqlite3 CLI` → descartada: exige uma ferramenta externa instalada. O Go e o `go-sqlite3` já estão no projeto.
- `Copiar o SQL para um arquivo de benchmark` → descartada: a cópia sai de sincronia com o código. Os benchmarks chamam as funções reais.
- `Tabela de agregados por balde (rollup)` → descartada por ora: é a saída se os índices não bastarem, mas adiciona estado e complexidade e foge da simplicidade do projeto. Revisar só se a T8 falhar nos alvos.
- `Pool separado para leitura` → descartada: muda a arquitetura de conexão e está fora do pedido. Revisar só se a medição mostrar bloqueio de gravação mesmo com os índices.
- `Índice cobrindo todas as colunas desde o início` → descartada: fica muito grande. Só entra se o `/stats` 30d continuar acima do alvo com o índice estreito.