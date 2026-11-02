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

- [x] Todos os Success Criteria estão marcados (dois com desvio justificado: lista em 1 de 3 rodadas e RAM na 1ª subida após a migration)
- [x] Build, lint e type-check passam
- [x] Nenhum item pendente no Status de Execução
- [x] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

**Medição**

- [x] Existe um gerador de dados que cria os perfis típico e limite de forma repetível
- [x] Linha de base registrada (antes dos índices) para os dois perfis, começando pelo `/stats` nos 4 períodos, com 30 dias de dados
- [x] Planos de consulta (`EXPLAIN QUERY PLAN`) de todas as consultas quentes registrados antes e depois

**Resultado (alvos assumidos, no perfil limite; a T5 confirma)**

- [x] Nenhuma consulta quente faz varredura completa de `checks` (`SCAN checks`). Já não fazia antes; o problema era ordenar à parte e ler a tabela inteira pelo índice
- [x] `/stats?range=30d`: p95 ≤ 500 ms (208 a 340 ms; antes 815 a 855 ms)
- [x] `GET /api/monitors` (lista com `history`, `uptime_24h` e `last_latency_ms`): p95 ≤ 200 ms. **Desvio:** 114 e 160 ms em duas rodadas e 203 ms na terceira (p50 de 82 a 135 ms). Antes: 91 a 107 s
- [x] `GET /api/monitors/:id/checks` (limit 50, offset 0): p95 ≤ 50 ms (13 a 22 ms)
- [x] Gravação de check durante uma leitura pesada de `/stats`: p95 ≤ 1 s (247 a 258 ms; antes 740 a 759 ms)
- [x] Limpeza de retenção não segura a conexão por mais de 2 s seguidos (maior espera de uma gravação: 81 a 104 ms; antes 2,4 s)

**Custos**

- [x] Custo de gravação por check sobe no máximo 20 % em relação à linha de base
- [x] Tamanho do arquivo antes e depois, e tamanho de cada índice, registrados
- [x] Tempo da migration em um banco do perfil limite medido e documentado
- [x] Pico de RAM continua abaixo de 50 MB, inclusive com o banco do perfil limite aberto e `/stats?range=30d` sendo consultado. **Desvio:** 39 a 40 MB em regime, mas 47 a 51 MB na primeira subida logo depois da migration (ver evidências)

**Qualidade**

- [x] Banco existente atualiza sem perda de dados; `migrate down` remove os índices
- [x] Respostas da API idênticas antes e depois (teste de equivalência)
- [x] Teste automático falha se uma consulta quente voltar a varrer `checks`
- [x] `make test` e `make build` passam, e `make test` continua rápido (benchmarks só rodam com variável de ambiente)

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

# Benchmarks e planos de consulta (a cópia aplica as migrations pendentes e
# o relatório mostra o tempo; um banco gerado antes da 0006 mede a atualização)
PINGOU_BENCH_DB=.bench/limite.db go test ./internal/repository -run TestDBBench -v -timeout 2h
scripts/dbbench.sh limite

# RAM do app com o banco grande aberto e /stats 30d + lista sem parar
make build && scripts/dbmem.sh .bench/limite.db 120

# Meta de memória do README
make build && HTTPS_N=50 scripts/memtest.sh 3600
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 0 | T1 | S | ~15 min | Levantamento das consultas | ✅ Concluído |
| Fase 1 | T2, T3 | M | ~1 h | Dados e ferramentas de medição | ✅ Concluído |
| Fase 2 | T4 | M | ~30 min (+ tempo de máquina) | Linha de base, começando pelo `/stats` 30d | ✅ Concluído |
| Fase 3 | T5 | S | ~15 min | Decisão (pode encerrar o plano) | ✅ Concluído |
| Fase 4 | T6, T7 | M | ~45 min | Índices e ajustes | ✅ Concluído |
| Fase 5 | T8, T9, T10 | M | ~1 h 15 min (+ 1 h do memtest) | Comparação, testes e documentação | ✅ Concluído |

Total estimado: cerca de 4 horas de trabalho, fora o tempo de máquina (geração do perfil limite e memtest de 1 hora). Se a T5 encerrar o plano, cai para cerca de 2 horas.

## Status de Execução

### Concluído

- **T1 (levantamento)**:
  - Índices que já existiam: `checks(monitor_id, checked_at DESC)`, `checks(checked_at)`, `incidents(monitor_id, ended_at)`, `incidents(started_at DESC)`, `incidents(ended_at DESC)`, `monitors(enabled|current_state|type)` e `ping_token UNIQUE`. Os candidatos "básicos" da T5 já estavam lá.
  - Datas são texto RFC3339 em UTC (`2026-01-01T00:00:00Z`), que ordena igual ao tempo. `success` é 0/1; `latency_ms` inteiro. O balde do `/stats` é `strftime('%s', checked_at) / bucket` no `SELECT` (não no filtro, então não atrapalha o índice).
  - Consultas quentes: série do `/stats` (`Buckets`), janela down (`FindInWindow`), lista do painel (`Summaries`: histórico de 30 checks + uptime 24 h; `FindOpenByMonitors`), checks paginados, incidentes (`FindAll` com `COUNT`), retenção (`DELETE` a cada hora), state machine (últimos N checks a cada **falha**), heartbeat (a cada 10 s) e `/ping/{token}`. O SQL foi para `internal/repository/queries.go` (mesmo texto) para o teste de plano usar as mesmas strings.
  - O painel chama `GET /api/monitors` a cada 15 s (`Shell.jsx`, `useMonitors`).
  - Última migration era a `0005` (do plano de pausa, já mergeado). Esta é a `0006`. O plano de fuso não criou migration (settings são chave-valor).
- **T2** — `scripts/seed`. Perfil típico: 397.440 checks, 560 incidentes, 95 MB, 2,6 s. Perfil limite: 24.667.200 checks, 5.622 incidentes, 6,0 GB, 4 min 31 s. Duas execuções deram as mesmas contagens e somas. (Os heartbeats pingam a cada 5 min, por isso o total fica um pouco abaixo dos ≈ 430 mil / 26 milhões do plano.)
- **T3** — `internal/repository/dbbench_test.go` (`TestDBBench`), `queryplan_test.go` (`TestPrintQueryPlans`, `QueryPlans`), `scripts/dbbench.sh` e `scripts/dbmem.sh`. Sem `PINGOU_BENCH_DB` os testes pulam. Adaptação: tempos com p50/p95 num `Test` em vez de `Benchmark` do Go (o `b.N` não dá percentis). Tamanho de cada índice medido com `DROP INDEX` numa transação desfeita (o `go-sqlite3` não tem `dbstat`).
- **T4 (linha de base)** — evidências abaixo. O primeiro número pedido, `/stats` 30d no perfil limite: **p95 de 815 a 855 ms**. O maior problema apareceu na lista do painel: **91 a 107 s**, acima do `WriteTimeout` de 30 s do servidor. Com o binário antigo e o banco limite, a lista falhou em todas as chamadas: o dashboard não carregava.
- **T5 (decisão)**:
  - `ROW_NUMBER()` sobre todos os checks dos 100 monitores: a lista lia ~25 milhões de entradas para ficar com 30 de cada. **Sim:** reescrever com `LIMIT` por monitor (`UNION ALL`).
  - `ORDER BY checked_at DESC, id DESC` com o índice `checked_at DESC`: o rowid fica crescente dentro do índice, e o SQLite ordenava à parte (`USE TEMP B-TREE FOR LAST TERM OF ORDER BY`). **Sim:** índice com `checked_at` e `id` crescentes, lido de trás para frente.
  - `/stats` 30d buscava cada linha na tabela (~260 mil por monitor). **Sim:** `success` e `latency_ms` no mesmo índice (cobrindo). Ele **substitui** o índice antigo, que tem o mesmo prefixo, então não há dois índices para manter. O novo ficou com 1,9 GB contra 3,3 GB do antigo. Parte disso é porque um índice recém-criado é compacto; com as inserções intercaladas ele tende a crescer de novo. Mesmo assim o total de índices não aumenta.
  - Gravação de **falha**: a state machine fazia `COUNT(*)` de todos os checks do monitor e jogava fora (26 ms p50). **Sim:** `CheckRepositoryTx.Latest`, sem o total.
  - Retenção: `DELETE` de um dia segurou a conexão por 2,4 s (acima do alvo de 2 s). **Sim:** lotes de 5 mil.
  - `checks(checked_at)`, incidentes, `ping_token`: já atendem (≤ 1 ms nos incidentes). **Não** mexer.
  - `ANALYZE`/`PRAGMA optimize` na migration: **não**. O planejador já escolhe os índices certos sem estatísticas, e o `ANALYZE` varreria os 25 milhões de linhas a mais na atualização. O teste de plano roda `ANALYZE` no banco pequeno, para ficar estável.
  - Alvos mantidos como estavam.
- **T6** — `0006_checks_indexes.sql`: cria `idx_checks_monitor_time` e apaga `idx_checks_monitor_checked`; o down faz o contrário.
- **T7** — log de migrations pendentes (versões e tempo); histórico da lista com `LIMIT` por monitor; `Latest` na state machine; retenção em lotes com `ctx` conferido entre eles. **Bug achado e corrigido** durante a medição de RAM: com o banco em dia, o `goose.CollectMigrations` devolve `ErrNoMigrationFiles`, e o app não subia na segunda vez. Teste `TestOpenTwice`.
- **T8 (depois)** — evidências abaixo. Todos os alvos cumpridos, com dois desvios: a lista deu 203 ms de p95 em 1 de 3 rodadas; a RAM passou de 50 MB na primeira subida após a migration.
- **T9** — `TestHotQueriesUseIndexes` (pelo nome do índice; falha sem a `0006`, conferido), `TestSummariesHistoryMatchesLegacy` (histórico novo × consulta antiga), `TestResultsSameWithOldIndex` (respostas de `/stats`, lista, checks e incidentes iguais antes e depois do `goose down`; o down remove o índice), `TestDeleteOlderThanBatches` (mesmo conjunto apagado; para com `ctx` cancelado), `TestOpenTwice`.
- **T10** — README: "Desempenho e tamanho do banco" e "Atualizando um banco grande" (tempo, log, espaço livre, `VACUUM`, `migrate down`). `scripts/README.md`: seção do dbbench. Dockerfile: `HEALTHCHECK --start-period=5m` (a migration passou de 30 s).

### Evidências (mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

Perfil limite (100 monitores × 10 s, 30 dias), p95 de 10 repetições. Antes = 2 rodadas no commit das ferramentas (`f8d3597`); depois = 3 rodadas no código final. A variação entre rodadas iguais ficou abaixo de 15 % nas consultas pesadas (`/stats` 30d: 815/855 ms antes, 208/291/228 ms depois), exceto a lista depois (114 a 203 ms), que tem poucas dezenas de ms e sente qualquer ruído.

| Consulta | Antes | Depois |
|---|---:|---:|
| `/stats` 24h http | 30 a 32 ms | 5 a 10 ms |
| `/stats` 7d http | 159 a 224 ms | 38 a 46 ms |
| `/stats` 14d http | 341 a 347 ms | 82 a 97 ms |
| **`/stats` 30d http** | **815 a 855 ms** | **208 a 291 ms** |
| `/stats` 30d tcp | 819 a 822 ms | 295 a 340 ms |
| `/stats` 30d heartbeat | 34 a 36 ms | 5 a 10 ms |
| `GET /api/monitors` | 91.0 a 106.8 s | 114 a 203 ms |
| checks offset 0 / 1.000 / 10.000 | 36 a 41 / 28 a 34 / 72 a 77 ms | 13 a 22 / 13 a 21 / 13 a 19 ms |
| incidentes (todos os filtros) | ≤ 0,4 ms | ≤ 1,2 ms (sem mudança; ruído) |
| gravação com sucesso (média) | 0,28 a 0,30 ms | 0,29 a 0,35 ms (+9 % na média) |
| gravação com falha | 40 ms | 0,4 a 0,5 ms |
| gravação com `/stats` 30d ao lado | 740 a 759 ms | 247 a 258 ms |
| retenção: `DELETE` de 1 dia (822 mil linhas) | 2,5 a 2,7 s numa instrução só | 6,0 s em lotes |
| retenção: maior espera de uma gravação | 2.358 ms | 81 a 104 ms |

Tamanho (perfil limite): índices de `checks` de 4.091 MB para 2.719 MB (`idx_checks_monitor_checked` 3.251 MB → `idx_checks_monitor_time` 1.878 MB; `idx_checks_checked_at` 841 MB igual). O arquivo foi de 5.972 MB para 7.850 MB: o índice novo é escrito antes de o antigo ser apagado, e as páginas do antigo ficam livres dentro do arquivo (reaproveitadas pelos checks novos ou devolvidas com `VACUUM`). Perfil típico: índices de 65 MB para 43 MB, arquivo de 95 para 125 MB.

Migration `0006`: 1 min 21 s a 1 min 27 s no benchmark e 1 min 18 s no app real (`scripts/dbmem.sh`, com o log novo). Perfil típico: 0,7 s.

Plano depois (sem `SCAN checks` e sem ordenação à parte):
`stats: buckets` → `SEARCH checks USING COVERING INDEX idx_checks_monitor_time (monitor_id=? AND checked_at>? AND checked_at<?)`; `lista: history` → `SEARCH checks USING COVERING INDEX idx_checks_monitor_time (monitor_id=?)` por monitor; `checks: página` → `SEARCH checks USING INDEX idx_checks_monitor_time (monitor_id=?)`; `retenção` → `SEARCH checks USING COVERING INDEX idx_checks_checked_at (checked_at<?)`.

Perfil típico (10 × 60 s), p95: lista 839 → 4,3 ms; `/stats` 30d http 70 → 35 ms; gravação com falha 6,1 → 0,4 ms.

RAM (`scripts/dbmem.sh`, app real, banco limite, `/stats` 30d de 3 monitores + lista sem parar):

| Situação | Pico (VmHWM) |
|---|---:|
| Antes (binário antigo, 120 s) | 31 MB, mas a lista falhou em todas as chamadas (corte de 30 s) e só deu 2 ciclos |
| Depois, banco já migrado (120 s, 154 ciclos) | 39 a 40 MB (27 MB do processo + 11 MB do binário mapeado) |
| Depois, primeira subida com a migration (60 a 120 s) | 47 a 51 MB |

O pico durante a migration em si ficou em 29 a 30 MB. O excesso aparece nas consultas logo depois (RssAnon 31 MB e RssFile 15 MB, contra 27 e 11 em regime) e some depois de reiniciar. Não é uma condição de regime, e este cenário (consultar o painel sem parar, com 25 milhões de checks) é bem mais pesado que o do README. Fica registrado como desvio.

`make test` (3,5 s), `make build` e `go vet ./...` verdes.

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