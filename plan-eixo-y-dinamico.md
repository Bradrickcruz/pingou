# Plano de Desenvolvimento: Eixo Y dinâmico no gráfico de tempo de resposta

## Overview

**Contexto (verificado no código)**:

- O gráfico é SVG feito à mão em [ResponseChart.jsx](web/src/components/monitor-detail/ResponseChart.jsx); **não há lib de gráfico**, então não existe "domínio automático" a reaproveitar.
- O eixo Y **já não é fixo**: `scaleOf` usa o pico da série e arredonda em passos de 100 ms, sempre com 3 intervalos (pico 402 → 0/200/400/600). Problemas: o teto fica folgado demais (402 → 600), passos de 100 ms em séries pequenas (pico 30 ms → 0/100/200/300, linha achatada embaixo), e ticks sempre múltiplos de 100.
- **Heartbeat hoje não tem gráfico nenhum**: [MonitorDetail.jsx:~272](web/src/pages/MonitorDetail.jsx) esconde o card com `!heartbeat`. A parte "Heartbeat" do pedido é, portanto, um gráfico **novo** (pings por balde), não um ajuste.
- Um ping de heartbeat é gravado como check `Success: true` com latência nula ([monitor_service.go:251](internal/service/monitor_service.go#L251)); prazo perdido vira check com falha. Logo, **pings recebidos por balde = `checks - failures`** (não `checks`, que inclui prazos perdidos). Isso foge do "campo `checks`" do pedido — ver Notes.

**Objetivo**: eixo Y proporcional aos dados de cada monitor/período, recalculado ao trocar `range`.

**Regra do máximo (HTTP/TCP)**: `peak` = maior `avg_latency_ms` não nulo; `p = max(peak, 100)`; `step` = menor valor em {1,2,5}×10ⁿ ≥ `p/5`; `max = ceil(p/step)*step`; ticks de `0` a `max` de `step` em `step`.
Confere com os exemplos: 70 → 100; 350 (step 100) → 400; 2300 (step 500) → 2500. Série toda `null` → 100 (na prática o gráfico já mostra "Sem dados", ver T3).

**Heartbeat**: barras de pings recebidos por balde; eixo Y de contagem (inteiros, mínimo 1, mesma função com piso 1). Balde com zero pings fica vazio e **sem faixa cinza** (zero é normal). Vermelho só vem dos períodos `down`.

**Fora de escopo**: zoom/seleção, outros gráficos, mudanças no `/stats`.

## Project Type

**FRONTEND** — SPA React em `web/` (sem mudança de backend).

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] HTTP/TCP: max do eixo segue a regra acima; 70→100, 350→400, 2300→2500, nunca < 100
- [ ] Baldes `avg_latency_ms: null` ignorados no pico; tudo `null` não quebra
- [ ] Trocar 24h/7d/14d/30d recalcula o eixo (sem estado extra: derivado de `stats`)
- [ ] Heartbeat mostra gráfico de pings por balde, eixo inteiro dinâmico
- [ ] Heartbeat: balde com 0 pings não é pintado de erro nem cinza; vermelho só em `down`
- [ ] Tooltip/legenda/aria-label do heartbeat falam de pings, não de "tempo de resposta/média/falhas"
- [ ] Teste unitário da função de escala passa (`pnpm test`)
- [ ] Nenhuma dependência nova, `/stats` intacto

## Tech Stack

| Technology    | Choice Rationale                                                         |
| ------------- | ------------------------------------------------------------------------ |
| SVG existente | Sem lib de gráfico; reaproveitar `ResponseChart` evita nova dependência  |
| `node --test` | Já usado em `web/src/lib/*.test.js`; escala fica em módulo puro testável |

## File Structure

```
web/src/
├── lib/
│   ├── chartScale.js        # novo: niceScale(peak, floor) → { max, ticks }
│   └── chartScale.test.js   # novo
├── components/monitor-detail/
│   └── ResponseChart.jsx    # usa niceScale; modo heartbeat (barras)
└── pages/MonitorDetail.jsx  # renderiza o card também para heartbeat
```

## Task Breakdown

| Task ID | Agent               | Name                                                             | Priority | Complexity | Dependencies |
| ------- | ------------------- | ---------------------------------------------------------------- | -------- | ---------- | ------------ |
| T1      | frontend-specialist | `niceScale` puro + teste                                         | P1       | S          | -            |
| T2      | frontend-specialist | Trocar `scaleOf` por `niceScale` (HTTP/TCP)                      | P1       | S          | T1           |
| T3      | frontend-specialist | Modo heartbeat no `ResponseChart` (barras de pings)              | P2       | M          | T1, T2       |
| T4      | frontend-specialist | Exibir o card no `MonitorDetail` para heartbeat + título/legenda | P2       | S          | T3           |
| T5      | qa-review           | Verificação manual por tipo e range                              | P3       | S          | T4           |

### Task T1: `niceScale` puro + teste

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - `pnpm test` roda hoje (`node --test src/lib/*.test.js`)
- **INPUT**: `scaleOf` embutido em ResponseChart.jsx, sem teste
- **OUTPUT**:
  1. `web/src/lib/chartScale.js` exportando `niceScale(peak, floor = 100)` → `{ max, ticks }` com a regra do Overview (step ∈ {1,2,5}×10ⁿ ≥ max(peak,floor)/5; passo mínimo 1; ticks inteiros)
  2. `chartScale.test.js` cobrindo: 0/NaN→floor; 70→100; 350→400; 2300→2500; 100→100; 101→120 (step 20); `floor=1` com peak 3 → max 3 e ticks `[0,1,2,3]`
- **VERIFY**:
  - `cd web && pnpm test` verde com os casos acima
- **Rollback**: remover os dois arquivos novos

### Task T2: Trocar `scaleOf` por `niceScale` (HTTP/TCP)

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**: T1 concluída
- **INPUT**: `scaleOf(series.map(b => b.avg_latency_ms ?? 0))` ([ResponseChart.jsx:13](web/src/components/monitor-detail/ResponseChart.jsx#L13))
- **OUTPUT**:
  1. `scaleOf` removido; chamada passa a `niceScale(Math.max(0, ...series.map(b => b.avg_latency_ms ?? 0)))` (`null`→0 é equivalente a ignorar; piso 100 cobre o resto)
  2. Rótulos/linhas de grade seguem usando `ticks` (já iteram sobre o array; `key={t}` continua único)
- **VERIFY**:
  - Monitor HTTP com pico ~35 ms: eixo 0–100 com grade de 20 em 20 e linha ocupando o gráfico
  - Pico ~2300 ms: eixo termina em 2500; trocar range recalcula
- **Rollback**: `git checkout` do ResponseChart.jsx

### Task T3: Modo heartbeat no `ResponseChart` (barras de pings)

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T1, T2
- **PRE-CONDITIONS**: confirmar no `/api/monitors/:id/stats` de um heartbeat real que `series[i].checks`/`failures` refletem pings/prazos perdidos (ver Risk Areas)
- **INPUT**: componente só desenha linha de latência, faixa cinza em baldes sem média e tooltip de latência
- **OUTPUT**:
  1. Prop `heartbeat` (ou derivada de tipo) no `ResponseChart`; valor por balde = `checks - failures`; escala `niceScale(max, 1)`
  2. Em vez do `<path>`, uma `<rect>` por balde com altura proporcional (balde 0 → sem barra); `gapBuckets` **não** é aplicado (sem faixa cinza)
  3. Faixas `down` inalteradas (única fonte de vermelho)
  4. Estado vazio: "Sem pings neste período." só se todos os baldes tiverem 0 **e** não houver `down`; caso contrário desenha o eixo
  5. Tooltip/`tipText`/`summary` em termos de pings ("3 pings"), sem "média"/"falhas"; ponto destacado vira destaque da barra
  6. Extrair o que variar por modo (valor do balde, texto) em pequenas funções locais; sem novo componente
- **VERIFY**:
  - Heartbeat diário em 30d: barras de altura 1, baldes vazios sem cor de erro
  - Período com `down`: faixa vermelha presente; fora dele, nada vermelho
  - Navegação por teclado e leitor de tela (aria-label) falam de pings
- **Rollback**: reverter ResponseChart.jsx (T2 permanece válido isolado)

### Task T4: Exibir o card no `MonitorDetail` para heartbeat

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T3
- **PRE-CONDITIONS**: T3 concluída
- **INPUT**: `{!heartbeat && (<div ...>)}` em [MonitorDetail.jsx](web/src/pages/MonitorDetail.jsx) esconde o card
- **OUTPUT**:
  1. Card renderizado para os 3 tipos; título "Tempo de resposta" (HTTP/TCP) vs "Pings recebidos" (heartbeat)
  2. `ChartLegend` com variante heartbeat: "pings por intervalo", "fora do ar"; **sem** "sem dados" (não se aplica)
  3. Passa `heartbeat` ao `ResponseChart`
- **VERIFY**:
  - Página de heartbeat mostra o card; HTTP/TCP inalterados
  - `pnpm build` sem erros
- **Rollback**: restaurar o condicional `!heartbeat`

### Task T5: Verificação manual

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T4
- **PRE-CONDITIONS**: app rodando com monitores http, tcp e heartbeat
- **INPUT**: implementação completa
- **OUTPUT**: checklist registrada em "Status de Execução"
- **VERIFY**:
  - Os 4 ranges nos 3 tipos; monitor novo sem checks; monitor com todos `null`; tema claro/escuro e mobile (rótulos do eixo Y com até 6 ticks cabem em 200 px)

## Risk Areas

| Risk                                                              | Impact                         | Mitigation                                                          |
| ----------------------------------------------------------------- | ------------------------------ | ------------------------------------------------------------------- |
| Heartbeat é gráfico novo, não só ajuste de eixo                   | Escopo maior que o descrito    | T3/T4 isoláveis; T1+T2 entregam o pedido HTTP/TCP sozinhos          |
| `checks` inclui prazos perdidos no heartbeat                      | Barra contaria falha como ping | Usar `checks - failures`; validar na pré-condição de T3             |
| Até 6 rótulos no eixo Y em 200 px                                 | Rótulos apertados no mobile    | `p/5` limita a 5 intervalos (~40 px cada); revisar em T5            |
| Exemplos do pedido (350→400, 2300→2500) não seguem só {1,2,5}×10ⁿ | Interpretação divergente       | Regra: passo {1,2,5}×10ⁿ com ≤5 intervalos; teste prova os exemplos |

## Rollback Strategy

1. Mudança só no frontend, sem migração/API: `git revert` do(s) commit(s).
2. Parcial: reverter T3/T4 mantém T1/T2 (eixo dinâmico HTTP/TCP).

## Phase X: Verification

```bash
cd web
pnpm test
pnpm run lint
pnpm run build
```

## Cronograma Sugerido

| Fase   | Tarefas | Complexity | Esforço | Target                       | Status      |
| ------ | ------- | ---------- | ------- | ---------------------------- | ----------- |
| Fase 1 | T1, T2  | S          | ~30 min | Eixo dinâmico HTTP/TCP       | ⬜ Pendente |
| Fase 2 | T3, T4  | M          | ~45 min | Gráfico de pings (heartbeat) | ⬜ Pendente |
| Fase 3 | T5      | S          | ~15 min | QA                           | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1–T5

## Notes

- Pedido diz usar `checks` para heartbeat; o plano usa `checks - failures` porque o scheduler grava prazo perdido como check com falha. Se o scheduler não gravar falhas para heartbeat, `failures` será 0 e vira equivalente.
- Outliers já suavizados (média por balde, timeouts fora de `avg_latency_ms`): máximo da série basta, sem p95.
- Todos `null` em HTTP/TCP: o componente já exibe "Sem dados de resposta neste período."; o piso de 100 ms vale para o caso de haver só `down` sem latência.

### Alternativas consideradas e descartadas

- Lib de gráfico (Recharts etc.) → descartada: dependência nova para algo que o SVG atual resolve.
- Percentil p95 como teto → descartado: pontos já são médias de balde.
- Manter 3 intervalos fixos no eixo → descartado: gera tetos folgados (402 → 600).
