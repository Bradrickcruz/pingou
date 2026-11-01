# Plano de Desenvolvimento: Largura no desktop + busca parcial no dashboard

## Overview

Duas melhorias só de frontend (`web/`, React + Tailwind), sem mudança no backend:

1. **Largura no desktop**: hoje `Dashboard`, `Incidents` e `MonitorDetail` usam `max-w-[1180px]` e `Settings` usa `max-w-[760px]`, o que deixa vazio à direita. Aumentar o limite (não remover) e ajustar proporções onde ficar estranho.
2. **Busca parcial**: o dashboard **já tem** um campo de busca (`Dashboard.jsx`), mas ele só filtra por nome, diferencia acentos (`"servico"` não acha `"Serviço"`), só aparece com 10+ monitores (`SEARCH_FROM`) e não persiste na URL. O trabalho é evoluir esse campo, não criar um novo.

Restrições: seguir `DESIGN.md`; sem mudanças no layout mobile; sem busca no backend; sem novas dependências.

## Project Type

**FRONTEND** — SPA React/Vite em `web/`. Testes unitários via `node --test src/lib/*.test.js`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e testes passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Em 1920px, Painel, Incidentes e Detalhe do monitor ocupam bem a largura, sem faixa vazia grande à direita
- [ ] Em 1366px e 2560px+ o layout continua legível (sem linhas esticadas demais, sem quebra)
- [ ] Novo limite de largura registrado no `DESIGN.md`
- [ ] Busca acha por parte do nome, ignorando maiúsculas e acentos
- [ ] Busca também acha por alvo: URL (http) e `host:porta` (tcp)
- [ ] O termo fica em `?q=` e sobrevive a recarregar e a voltar da página de detalhe
- [ ] Sem resultado: mensagem clara + ação "Limpar busca"
- [ ] Teste unitário da função de busca passa

## Tech Stack

| Technology | Choice Rationale |
|---|---|
| React 18 + react-router-dom | já em uso; `useSearchParams` cobre o `?q=` sem dependência nova |
| Tailwind | já em uso; largura é só classe utilitária |
| `String.prototype.normalize("NFD")` | remove acentos sem lib |
| `node --test` | runner de testes já configurado em `web/package.json` |

## File Structure

```
web/src/
├── lib/
│   ├── monitors.js          # + matchesQuery()
│   └── monitors.test.js     # + casos de busca
├── pages/
│   ├── Dashboard.jsx        # busca com ?q=, estado vazio; largura
│   ├── Incidents.jsx        # largura
│   ├── MonitorDetail.jsx    # largura (3 ocorrências)
│   └── Settings.jsx         # largura (decisão em T2)
DESIGN.md                    # registrar limite de largura
```

## Task Breakdown

| Task ID | Agent | Name | Priority | Complexity | Dependencies |
|---|---|---|---|---|---|
| T1 | code-test | `matchesQuery` + testes | P1 | S | - |
| T2 | frontend-specialist | Busca no Dashboard (`?q=`, alvo, vazio) | P2 | M | T1 |
| T3 | frontend-specialist | Aumentar limite de largura nas páginas | P2 | S | - |
| T4 | frontend-specialist | Ajustar proporções por página | P2 | M | T3 |
| T5 | qa-review | Validação visual (1366 / 1920 / 2560) e da busca | P3 | S | T2, T4 |
| T6 | document | Registrar largura no `DESIGN.md` | P3 | XS | T4 |

### Task T1: `matchesQuery` + testes

- **Agent**: code-test
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - `web/src/lib/monitors.js` tem `targetOf(m)` e `monitors.test.js` existe
- **INPUT**: filtro atual no Dashboard faz só `m.name.toLocaleLowerCase("pt-BR").includes(q)`
- **OUTPUT**:
  1. `matchesQuery(m, q)` em `lib/monitors.js`: normaliza (minúsculas + NFD sem diacríticos) e testa `contém` no nome e no alvo. Para http usa `m.url`; para tcp `host:porta`. **Heartbeat: só nome** (o `targetOf` dele é uma frase, não um alvo).
  2. Query vazia → `true`.
  3. Testes: acento (`servico` ↔ `Serviço`), caixa, parcial, URL, `host:porta`, heartbeat não casa pela frase, query vazia.
- **VERIFY**:
  - `cd web && npm test` passa

### Task T2: Busca no Dashboard

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 concluída
- **INPUT**: `useState("")` para `query`; filtro inline no `visible`; mensagem "Nenhum monitor neste filtro." (mobile e tabela)
- **OUTPUT**:
  1. `query` vem de `useSearchParams` (`?q=`); `replace: true` ao digitar para não poluir o histórico; param removido quando vazio.
  2. `visible` usa `matchesQuery`.
  3. Placeholder/label: "Buscar por nome ou endereço".
  4. Estado vazio com busca ativa: "Nenhum monitor encontrado para “…”" + botão **Limpar busca** (mobile e tabela). Sem busca, mantém a mensagem atual do filtro.
  5. Decisão: manter `SEARCH_FROM = 10`? Ver Notes.
- **VERIFY**:
  - Digitar, recarregar a página e voltar do detalhe mantém o termo
  - `servico` acha `Serviço`; trecho de URL acha monitor http
  - Estado vazio mostra o botão e ele limpa a busca
  - `npm run lint` sem erros

### Task T3: Aumentar limite de largura

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - `main` em `Shell.jsx` não limita largura (`flex-1 min-w-0 p-4 md:p-8`); o limite vem de cada página
- **INPUT**: `max-w-[1180px]` em `Dashboard.jsx:120`, `Incidents.jsx:201`, `MonitorDetail.jsx:401/413/464`; `max-w-[760px]` em `Settings.jsx:41/97`
- **OUTPUT**:
  1. Trocar `1180px` por um valor entre 1600 e 1920px (começar em ~1680px; fechar na validação da T5).
  2. `Settings`: manter estreito por ser formulário (decisão: só subir se ficar estranho ao lado das outras). Registrar a decisão.
- **VERIFY**:
  - `grep -rn "max-w-\[1180px\]" web/src` não retorna nada
  - Em 1920px não há mais faixa vazia grande à direita

### Task T4: Ajustar proporções

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: M
- **Dependencies**: T3
- **PRE-CONDITIONS**:
  - T3 concluída
- **INPUT**: tabela do dashboard, cards/gráfico de `MonitorDetail`, lista de `Incidents` nas novas larguras
- **OUTPUT**:
  1. Dashboard: a coluna "Monitor" já usa `max-w-0 w-full` e absorve a folga; conferir que "Últimos 30 checks" e "Resposta" não ficam perdidos. Ajustar espaçamento/largura de colunas só se necessário.
  2. `MonitorDetail`: conferir que o gráfico (`ResponseChart`) e os blocos de resumo se distribuem bem (grid), sem cards esticados demais.
  3. `Incidents`: linhas legíveis, sem texto solto no canto esquerdo e vazio no direito.
- **VERIFY**:
  - Capturas em 1366, 1920 e 2560px sem quebra nem espaço desbalanceado
  - Mobile e tablet inalterados

### Task T5: Validação visual e de busca

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T2, T4
- **PRE-CONDITIONS**:
  - App rodando (`make` / `npm run dev` + backend)
- **INPUT**: telas implementadas
- **OUTPUT**:
  1. Checklist: 3 larguras × 3 páginas; busca com acento, URL, `host:porta`, vazio, recarregar, voltar do detalhe, combinada com os filtros de estado.
- **VERIFY**:
  - Lista de problemas encontrados (ou "nenhum") registrada em Status de Execução

### Task T6: Registrar no `DESIGN.md`

- **Agent**: document
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T4
- **PRE-CONDITIONS**:
  - Valor final do limite definido
- **INPUT**: `DESIGN.md` não menciona largura máxima
- **OUTPUT**:
  1. Uma linha com o limite de largura das páginas e o motivo (legibilidade em ultrawide).
- **VERIFY**:
  - `grep -n "largura" DESIGN.md` mostra a regra

## Risk Areas

| Risk | Impact | Mitigação |
|---|---|---|
| Linhas muito largas dificultam leitura da tabela | Médio | Manter limite (não remover); validar em 2560px |
| `?q=` some ao navegar pelo menu lateral | Baixo | Esperado (links do menu vão para `/`); só precisa sobreviver a recarregar e a voltar do detalhe |
| Busca escondida com menos de 10 monitores contraria a task | Baixo | Decidir em T2 (ver Notes) |
| Heartbeat casar pela frase "Espera um ping..." | Baixo | `matchesQuery` usa só o nome para heartbeat (T1) |

## Rollback Strategy

1. Mudanças só em `web/src` + `DESIGN.md`; `git revert` do commit de cada task.
2. Largura: voltar a classe para `max-w-[1180px]`.
3. Busca: voltar a `useState` e ao filtro por nome.

## Phase X: Verification

```bash
cd web
npm test
npm run lint
npm run build
```

## Cronograma Sugerido

| Fase | Tarefas | Complexity | Esforço | Target | Status |
|---|---|---|---|---|---|
| Fase 1 | T1, T3 | S | ~30 min | `matchesQuery` testada; limite novo aplicado | ⬜ Pendente |
| Fase 2 | T2, T4 | M | ~1 h | Busca e proporções prontas | ⬜ Pendente |
| Fase 3 | T5, T6 | S | ~30 min | Validado e documentado | ⬜ Pendente |

## Status de Execução

### Concluído

-

### Em andamento

-

### Pendente

- T1 a T6

## Notes

- A busca por nome já existe no dashboard; a task vira evolução dela. Isso reduz o esforço.
- **Pergunta aberta:** a task diz "campo de busca no topo", mas o código esconde o campo com menos de 10 monitores (`SEARCH_FROM`). O plano mantém o limite. Se quiser o campo sempre visível, é só remover a condição em T2.
- Busca é filtro combinado (E) com os filtros de estado já existentes (Todos / Fora do ar / No ar / Pausados).
- Limite de largura final (1600–1920px) fecha na T5, olhando as capturas.

### Alternativas consideradas e descartadas

- Remover o limite de largura → descartado: linhas ilegíveis em ultrawide.
- Busca no backend → descartado: `GET /api/monitors` já traz tudo (máx. 100).
- Debounce na busca → descartado: filtro em memória sobre ≤100 itens.
