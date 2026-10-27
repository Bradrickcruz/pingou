# Plano de Desenvolvimento: Fase 4 — Migração para o tema claro do DESIGN.md

## Overview

O app inteiro hoje é escuro (paleta roxa de [tokens.js](web/src/theme/tokens.js) e variáveis `--text`/`--bg`/`--accent` de [index.css](web/src/index.css)). O DESIGN.md define o tema claro como padrão: fundo `--pg-nuvem`, superfícies `--pg-branco`, texto `--pg-ardosia`, botão primário `--pg-ciano-fundo`.

**Objetivo:** trocar o app para o tema claro do DESIGN.md, e tirar o uso de cores fixas do código para que os componentes usem as classes `pg-*`.

**Estratégia (duas camadas):**

1. **Troca de tema centralizada (Fase 4A).** Os valores antigos passam a apontar para os tokens `--pg-*`. Como 12 arquivos leem `t.colors.*` e 40 pontos usam `var(--...)`, uma única alteração em `tokens.js` e `index.css` deixa o app inteiro claro, sem editar componente. É o único momento em que a aparência muda de uma vez.
2. **Migração gradual dos componentes (Fase 4B).** Cada componente troca `style={{ color: t.colors.x }}` por classes `pg-*`, um por vez. Nessa etapa não há mudança visual esperada, é só refatoração. Isso cumpre o pedido de migrar "aos poucos" sem telas meio claras e meio escuras.

**Limites:**

- Não mexer em ícones, fontes nem no conteúdo textual (já concluídos nas Fases 2 e 3).
- Grafismos (DESIGN.md §8) ficam fora: o app não tem nenhum hoje, e o §8 diz que telas de trabalho ficam limpas.
- Tema escuro: não será mantido. A §10 do DESIGN.md deixa de ser pendência para este app. Se um dia for necessário, vira plano próprio.

**Restrições:**

- Cores fora de `DESIGN.md §3.1` só com aprovação (§3.4). Os tons de hover e de sombra de overlay que faltam são derivados da própria cor (ver Risk Areas).
- Mudanças de tema ficam em commits separados, para reverter a troca visual sem desfazer a refatoração.

## Project Type

**FRONTEND** — React 19 + Vite 6 + Tailwind 3 (`web/`). Build embutido no binário Go via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [x] Todos os Success Criteria estão marcados
- [x] `npm run build` e `npm run lint` passam em `web/` (os 2 erros de lint pré-existentes são resolvidos na T16)
- [x] Nenhum item pendente no Status de Execução (restam só as decisões de design abaixo)
- [x] Evidências de verificação registradas (screenshots das 4 rotas + Login, claro)

> Revisar o plano se: a aprovação de cores derivadas for negada, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [x] Nenhuma tela do app usa fundo escuro (`#0f1117`, `#1a1d27`, `#16171d`) em nenhum estado
- [x] `grep -rn "theme/tokens" web/src` retorna vazio (arquivo removido na T15)
- [x] `grep -rn "var(--text\|var(--bg\|var(--accent\|var(--border" web/src` retorna vazio (tokens antigos removidos)
- [x] Nenhum `style={{ color|background|borderColor: ... }}` com cor hardcoded ou vinda de `t.colors` nos componentes
- [x] Botão primário é `--pg-ciano-fundo` com texto branco, e só um por tela (DESIGN.md §7)
- [x] Badges de status (UP / DOWN / UNKNOWN) têm texto + ícone, fundo `*-bg` e pílula `rounded-full` (DESIGN.md §4)
- [x] Nenhum `bg-black`, `#000` ou `text-gray-*` no código (DESIGN.md §3.4)
- [x] Texto de corpo ≥ 4,5:1 sobre o fundo (DESIGN.md §11)
- [x] Checklist da §11 do DESIGN.md marcado para todas as telas (verificado por grep e leitura de código; ver "Resultado da execução")

## Tech Stack

| Technology                              | Choice Rationale                                                                        |
| --------------------------------------- | --------------------------------------------------------------------------------------- |
| Tokens `--pg-*` (CSS custom properties) | Já existem em `:root` (Fase 1). Um único lugar para trocar cores                        |
| Tailwind `pg-*` (classes)               | Já configuradas (Fase 1). Substituem `style={{}}` e `t.colors` nos componentes          |
| `tokens.js` como ponte temporária       | Permite trocar o tema sem editar componente. Removido na T15                            |
| Tailwind `hover:` / `focus-visible:`    | Substituem os `onMouseEnter`/`onMouseLeave` do `Button`. Menos JS, foco visível via CSS |

## File Structure

```text
web/src/
├── index.css                 # 4A: remove bloco prefers-color-scheme dark; color-scheme: light;
│                             #     legados (--text, --bg, --border, --accent…) apontam para --pg-*
│                             #     4C: remove legados e bloco .counter/.hero (sem uso)
├── theme/
│   ├── tokens.js             # 4A: valores passam a ser var(--pg-*); 4C: arquivo removido
│   └── globalStyles.js       # 4A: input/focus/scrollbar com tokens claros; 4C: sem tokens.js
├── components/
│   ├── ui/Button.jsx         # 4B: variantes primary/secondary/destructive em classes pg-*
│   ├── ui/Badge.jsx          # 4A: cores *-bg e pílula; 4B: sem stateColors hardcoded
│   ├── ui/Modal.jsx          # 4A: overlay derivado de ardósia; 4B: classes pg-*
│   ├── ui/Spinner.jsx        # 4B
│   ├── ui/ConnectionOverlay.jsx  # 4B
│   ├── ui/TimeRangeSlider.jsx    # 4B + lint (T16)
│   ├── layout/Shell.jsx      # 4B: sidebar branca, item ativo ciano-nevoa
│   ├── monitors/MonitorCard.jsx  # 4B
│   └── monitors/MonitorForm.jsx  # 4B
├── hooks/useConnection.jsx   # 4C: fast-refresh (lint)
└── pages/{Dashboard,Incidents,Settings,Login}.jsx  # 4B
```

## Mapeamento de tokens (fonte da T2 e T3)

| Token antigo (`tokens.js`) | Valor antigo        | Novo                                              | Motivo                                                                          |
| -------------------------- | ------------------- | ------------------------------------------------- | ------------------------------------------------------------------------------- |
| `bg`                       | `#0f1117`           | `var(--pg-nuvem)`                                 | Fundo geral (§7 Card, §3.3)                                                     |
| `surface`                  | `#1a1d27`           | `var(--pg-branco)`                                | Cards, sidebar, modal (§7 Card, Input)                                          |
| `surfaceAlt`               | `#22263a`           | `var(--pg-ciano-nevoa)`                           | Item ativo da nav, cabeçalho de tabela (§7). Inputs usam `--pg-branco` (ver T4) |
| `border`                   | `#2e3248`           | `var(--pg-borda)`                                 | Bordas (§6)                                                                     |
| `primary`                  | `#6c63ff` (violeta) | `var(--pg-ciano-fundo)`                           | Botão primário (§7). Violeta sai: não é cor do manual                           |
| `primaryHov`               | `#574fd6`           | `var(--pg-ciano-fundo)` com `hover:brightness-90` | **Derivado — requer aprovação** (ver Risk Areas)                                |
| `success`                  | `#22c55e`           | `var(--pg-up)`                                    | Status UP (§4). Contraste 4,7:1 sobre branco                                    |
| `danger`                   | `#ef4444`           | `var(--pg-down)`                                  | Status DOWN, ações destrutivas (§7)                                             |
| `warning`                  | `#f59e0b`           | **a definir**                                     | Não existe no DESIGN.md. Ver Risk Areas                                         |
| `unknown`                  | `#6b7280`           | `var(--pg-unknown)`                               | Status UNKNOWN (§4)                                                             |
| `textPrimary`              | `#f1f5f9`           | `var(--pg-ardosia)`                               | Texto principal (§3.1)                                                          |
| `textMuted`                | `#94a3b8`           | `var(--pg-neblina)`                               | Texto secundário, 5,5:1 sobre branco (§3.5)                                     |
| `radius.sm`                | `6px`               | `8px` (controles)                                 | §6                                                                              |
| `radius.md`                | `10px`              | `12px` (cards)                                    | §6                                                                              |
| `radius.lg`                | `16px`              | `12px`                                            | §6 não define 16px. Usar 12px para cards grandes                                |

Legados em `index.css` (`--text: #6b6375`, `--text-h`, `--bg`, `--border`, `--accent: #aa3bff` violeta, `--social-bg`…): passam a apontar para os tokens `--pg-*` correspondentes na T1.

## Task Breakdown

| Task ID | Agent                    | Name                                                           | Priority | Complexity | Dependencies |
| ------- | ------------------------ | -------------------------------------------------------------- | -------- | ---------- | ------------ |
| T0      | frontend-specialist      | Commitar Fases 1–3                                             | P0       | XS         | -            |
| T1      | frontend-specialist      | Tema claro em `index.css` + legados apontando para `--pg-*`    | P0       | S          | T0           |
| T2      | frontend-specialist      | `tokens.js` apontando para `var(--pg-*)`                       | P0       | S          | T1           |
| T3      | frontend-specialist      | Radius de `tokens.js` alinhado a §6                            | P0       | XS         | T2           |
| T4      | frontend-specialist      | `globalStyles.js`: input, foco, scrollbar                      | P0       | S          | T2           |
| T5      | frontend-specialist      | `Badge.jsx`: cores de status, pílula, sem hex                  | P0       | S          | T2           |
| T6      | frontend-specialist      | `Modal.jsx`: overlay derivado de ardósia                       | P0       | XS         | T2           |
| T7      | qa-review                | Verificação visual da troca de tema (4A)                       | P0       | M          | T1–T6        |
| T8      | frontend-specialist      | `Button.jsx` em classes `pg-*` + encaminhar `type`             | P1       | S          | T7           |
| T9      | frontend-specialist      | `Shell.jsx` em classes `pg-*`                                  | P1       | S          | T8           |
| T10     | frontend-specialist      | `Dashboard.jsx` + `MonitorCard.jsx`                            | P1       | M          | T8           |
| T11     | frontend-specialist      | `MonitorForm.jsx` + `Modal.jsx` em classes                     | P1       | M          | T8, T6       |
| T12     | frontend-specialist      | `Incidents.jsx`                                                | P1       | S          | T8           |
| T13     | frontend-specialist      | `Settings.jsx`                                                 | P1       | S          | T8           |
| T14     | frontend-specialist      | `Login.jsx`, `Spinner`, `ConnectionOverlay`, `TimeRangeSlider` | P1       | S          | T8           |
| T15     | frontend-specialist      | Remover `tokens.js` e legados sem uso                          | P2       | S          | T9–T14       |
| T16     | frontend-specialist      | Resolver os 2 erros de lint pré-existentes                     | P2       | XS         | -            |
| T17     | ux-accessibility-auditor | Auditoria de contraste e a11y no tema claro                    | P3       | S          | T15          |
| T18     | qa-review                | Verificação final (checklist DESIGN.md §11)                    | P3       | M          | T17          |

---

### Fase 4A: Troca centralizada para o tema claro

### Task T0: Commitar Fases 1–3 ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Working tree com as mudanças das Fases 1–3 ainda não commitadas (`git status` mostra `web/`, `DESIGN.md`)
- **INPUT**: 3 fases em um estado não commitado. Rollback difícil sem commit.
- **OUTPUT**:
  1. Um commit por fase (tokens, ícones, fontes) + commit do DESIGN.md. Só com aprovação do usuário, conforme regra do repositório.
- **WHY**: a Fase 4A muda a aparência inteira. Ter as fases anteriores commitadas permite reverter só a troca visual.
- **VERIFY**:
  - `git log --oneline` mostra os commits; `git status` limpo em `web/`
- **Rollback**: `git revert` do commit da fase.

### Task T1: Tema claro em `index.css` ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T0
- **PRE-CONDITIONS**:
  - Bloco `:root` com `--pg-*` da Fase 1 existe
- **INPUT**: `:root` com legados escuros; bloco `@media (prefers-color-scheme: dark)` sobrescreve os legados; `color-scheme: light dark`.
- **OUTPUT**:
  1. Remover o bloco `@media (prefers-color-scheme: dark)` inteiro. O tema escuro não é mais suportado (ver Overview).
  2. `color-scheme: light`.
  3. Legados em `:root` passam a apontar para os tokens claros: `--text: var(--pg-neblina)`, `--text-h: var(--pg-ardosia)`, `--bg: var(--pg-nuvem)`, `--border: var(--pg-borda)`, `--accent: var(--pg-ciano-fundo)`, `--accent-bg: var(--pg-ciano-nevoa)`, `--accent-border: var(--pg-ciano)`, `--code-bg: var(--pg-ciano-nevoa)`, `--social-bg: var(--pg-branco)`.
  4. `--shadow`: remover sombra, usar `none` (DESIGN §6: evite sombra).
  5. Não tocar em `h1`/`h2`/`code` nem nos `@apply`. Eles já leem `var(--text-h)`.
- **WHY**: 40 usos de `var(--...)` passam a ler cor clara sem editar componente. Esse é o motivo de a troca ser cheap.
- **VERIFY**:
  - `npm run dev`, abrir qualquer rota: fundo `#F3F7F8`, texto ardósia
  - Sistema em modo escuro (prefers-color-scheme: dark): app continua claro
- **Rollback**: restaurar `index.css` do commit T0.

### Task T2: `tokens.js` apontando para `var(--pg-*)` ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Mapeamento da tabela acima aprovado
- **INPUT**: `tokens.colors` com hex escuros.
- **OUTPUT**:
  1. Cada valor de `colors` vira a string `var(--pg-...)` da tabela de mapeamento. Exceção: `warning` e `primaryHov` (ver Risk Areas, ainda sem valor final).
  2. `font.mono` e `font.sans` não mudam.
- **WHY**: `t.colors.*` é lido como string em `style={{}}`. `var()` funciona nesses lugares e troca o tema sem editar nenhum componente com cor.
- **VERIFY**:
  - `grep -rn "t.colors" web/src` → nenhum uso faz comparação numérica ou `rgba()` com o valor (se fizer, tratar na T5 ou T6)
  - Rodar o app: todas as telas claras, nenhum texto invisível
- **Rollback**: restaurar `tokens.js`.

### Task T3: Radius alinhado a §6 ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T2
- **OUTPUT**: `radius.sm: 8px`, `radius.md: 12px`, `radius.lg: 12px`.
- **VERIFY**: cards e inputs com cantos de 8 e 12 px (inspeção no DevTools).
- **Rollback**: restaurar valores.

### Task T4: `globalStyles.js` — input, foco, scrollbar ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T2
- **OUTPUT**:
  1. `input, textarea, select`: `background: var(--pg-branco)` (DESIGN §7: input é branco, não névoa), `border-color: var(--pg-borda)`, `color: var(--pg-ardosia)`.
  2. `input:focus`: `border-color: var(--pg-ciano)` + anel `box-shadow: 0 0 0 2px var(--pg-ciano)` (DESIGN §7: anel ciano no foco). Trocar o `outline: none` por anel visível (regra a11y: foco não pode sumir).
  3. Scrollbar: track `var(--pg-nuvem)`, thumb `var(--pg-borda)`.
- **WHY**: o `outline: none` atual remove o foco visível no teclado. Sem anel de substituto, isso é bloqueante (ver `ux-accessibility-auditor`, regra 3).
- **VERIFY**:
  - Tab até um input: anel ciano visível
  - Digitar num input: texto legível em ardósia sobre branco
- **Rollback**: restaurar o bloco CSS.

### Task T5: `Badge.jsx` — cores de status e pílula ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T2
- **INPUT**: `stateColors` com hex escuros (`#14532d`, `#450a0a`, `#1f2937`) e `borderRadius: t.radius.sm`.
- **OUTPUT**:
  1. `UP`: `bg: var(--pg-up-bg)`, `color: var(--pg-up)`, ícone `Check`.
  2. `DOWN`: `bg: var(--pg-down-bg)`, `color: var(--pg-down)`, ícone `X`.
  3. `UNKNOWN`: `bg: var(--pg-unknown-bg)`, `color: var(--pg-unknown)`, ícone `CircleHelp`.
  4. Formato pílula: `rounded-full` (DESIGN §4 e §6, `9999px`). Remover o `borderRadius` inline.
  5. Fonte: Inter 500 (`font-medium`), 12–13 px. Hoje é `font-semibold` 11px, que usa peso 600 (proibido pela §5). Trocar para `font-medium` e `text-[12px]`.
- **WHY**: o Badge é a única informação de estado da tela. Se ele ficar com cor errada, o usuário não sabe se o monitor caiu.
- **VERIFY**:
  - Monitor UP, DOWN e UNKNOWN: pílula clara com ícone + texto nos três
  - Texto de badge com contraste ≥ 4,5:1 (`--pg-up` 4,7:1 e `--pg-down` 5,0:1 sobre branco, mas sobre `*-bg` conferir)
- **Rollback**: restaurar `Badge.jsx`.

### Task T6: `Modal.jsx` — overlay derivado de ardósia ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T2
- **INPUT**: `bg-black/60` (preto puro, proibido pela §3.4).
- **OUTPUT**: `bg-[rgba(30,43,51,0.5)]` (ardósia `#1E2B33` com 50% de opacidade). Derivado da própria cor, não de outra paleta (§3.4).
- **VERIFY**: overlay acinzentado azulado, não preto.
- **Rollback**: restaurar a classe.

### Task T7: Verificação visual da troca de tema (4A) ✅

- **Agent**: qa-review
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T1–T6
- **INPUT**: app com tema claro.
- **OUTPUT**: screenshots de Login, Dashboard (com e sem monitores), Incidents, Settings, modal aberto e cada status de badge, em 1280 px e 390 px.
- **VERIFY**:
  - Nenhuma área escura restante (buscar visualmente o fundo do `ConnectionOverlay` e do `Spinner`; eles ficam para a T14)
  - Texto legível em todas as telas
  - Responsivo em 390 px sem scroll horizontal
- **Rollback**: reverter o commit de 4A.

---

### Fase 4B: Migração gradual dos componentes

> Cada task desta fase é uma refatoração **sem mudança visual esperada**. Se a tela mudar, a task está errada.

### Task T8: `Button.jsx` em classes `pg-*` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T7
- **INPUT**: `variants` em JS, `onMouseEnter`/`onMouseLeave` para hover, `style={{}}` com cor. O `Button` também **descarta a prop `type`**, que `Login.jsx` e `Settings.jsx` passam como `type="submit"`.
- **OUTPUT**:
  1. Variantes por classe: `primary` = `bg-pg-ciano-fundo text-white hover:brightness-90`; `secondary` = `border border-pg-borda text-pg-ardosia bg-transparent`; `destructive` = `text-pg-down border border-pg-down` (fundo sólido só na confirmação, §7).
  2. Remover `onMouseEnter`/`onMouseLeave` e o `style={{}}`.
  3. Forwardar `type` (default `"button"`; o `Login` e o `Settings` passam `"submit"`).
  4. Raio: `rounded-pg-control` (8px). Fonte: `font-medium text-pg-ui`.
- **WHY**: `Button` é usado em todas as telas. Trocar aqui tem o maior retorno e reduz o JS de hover.
- **VERIFY**:
  - Login e Settings: botão submete o formulário (teste Enter e clique)
  - Hover escurece o primário; secundário e destrutivo mudam de fundo no hover
  - `grep -n "style=" web/src/components/ui/Button.jsx` vazio
- **Rollback**: restaurar `Button.jsx`.

### Task T9: `Shell.jsx` em classes `pg-*` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T8
- **INPUT**: sidebar com `t.colors.surface`, `borderRight`, `NavLink` com `style` por `isActive`.
- **OUTPUT**:
  1. Sidebar: `bg-pg-branco border-r border-pg-borda`.
  2. Item ativo: `bg-pg-ciano-nevoa text-pg-ciano-fundo font-medium`. Inativo: `text-pg-neblina hover:bg-pg-nuvem`.
  3. Logo: manter `pingou-simbolo.svg` (é o símbolo padrão, funciona em fundo claro, §2). Texto "Pingou" em `text-pg-ciano-fundo` (≥ 24 px para usar ciano puro não é exigido aqui, o texto da marca fica no ciano-fundo para passar AA).
  4. `NavLink` `style` removido: o estado ativo vem de `isActive` nas classes.
- **VERIFY**:
  - Item ativo visível em cada rota; navegação por Tab com anel de foco
- **Rollback**: restaurar `Shell.jsx`.

### Task T10: `Dashboard.jsx` + `MonitorCard.jsx` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T8
- **INPUT**: cards de estatística com `style={{ color, background, borderColor, borderRadius }}`, cabeçalho e lista com `t.colors`.
- **OUTPUT**:
  1. Cards de stat: `bg-pg-branco border-pg-borda rounded-pg-card`. Número em `font-medium text-pg-h2 tabular-nums` (§5: números com `tabular-nums`). Cor do número: `text-pg-up`, `text-pg-down`, `text-pg-unknown` conforme o estado.
  2. Botão de refresh: `bg-pg-branco border-pg-borda`. Manter `aria-label`, `aria-busy` e `motion-reduce:animate-none` da Fase 2.
  3. Título `h1` continua em Rubik 500 (regra já no `index.css`).
  4. `MonitorCard`: mesma troca de `style` por classes. Ícones e badges já usam os tokens da T5.
- **VERIFY**:
  - Os 3 contadores (up/down/unknown) com cor de estado
  - Refresh gira só fora de `prefers-reduced-motion`
- **Rollback**: restaurar os dois arquivos.

### Task T11: `MonitorForm.jsx` + `Modal.jsx` em classes ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T8, T6
- **INPUT**: `Modal` com `style` de fundo e borda; `MonitorForm` com labels e inputs estilizados por `t.colors`.
- **OUTPUT**:
  1. `Modal`: `bg-pg-branco border-pg-borda rounded-pg-card`, sem `style`.
  2. Labels: `text-pg-ui text-pg-neblina font-medium`.
  3. Inputs: classes da §7 (branco, borda `--pg-borda`, raio 8 px, foco com anel ciano). Como o `globalStyles` já cobre `input`, só remover o `style` dos inputs.
  4. Erro de formulário: `text-pg-down`, com texto que diz o que fazer (§9).
- **VERIFY**:
  - Abrir "Add Monitor", tabular pelos campos com anel ciano visível
  - Esc fecha o modal (comportamento atual preservado)
- **Rollback**: restaurar os dois arquivos.

### Task T12: `Incidents.jsx` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T8
- **INPUT**: cards de incidente com `t.colors.danger` para aberto e `t.colors.success` para resolvido, e `t.colors.surface`.
- **OUTPUT**: borda do card aberto `border-pg-down`, fechado `border-pg-borda`; fundo `bg-pg-branco`; textos com `text-pg-neblina` ou `text-pg-ardosia`. Badges OPEN/RESOLVED já com cor de status (texto + ícone da Fase 2).
- **VERIFY**: incidente aberto com borda vermelha clara e ícone; fechado neutro.
- **Rollback**: restaurar `Incidents.jsx`.

### Task T13: `Settings.jsx` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T8
- **INPUT**: cards e inputs com `style`; botão de download com `t.colors.surfaceAlt`.
- **OUTPUT**: cards `bg-pg-branco border-pg-borda rounded-pg-card`; botão "Download dump" vira `secondary` do `Button` (§7), trocando o `button` cru por `<Button variant="secondary">`.
- **VERIFY**: "Save Settings" é o único primário da tela; "Download dump" é secundário.
- **Rollback**: restaurar `Settings.jsx`.

### Task T14: `Login.jsx`, `Spinner`, `ConnectionOverlay`, `TimeRangeSlider` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T8
- **INPUT**: Login com `t.colors.bg`, `surface`, `primary` (roxo no título), `danger`, `textMuted`. `Spinner` e `ConnectionOverlay` com fundo escuro (o que a T7 deixou para trás). `TimeRangeSlider` com cores inline.
- **OUTPUT**:
  1. Login: fundo `bg-pg-nuvem`; cartão `bg-pg-branco rounded-pg-card`; título em `text-pg-ciano-fundo` com `font-heading` (Rubik 500, ver §5).
  2. `Spinner`: `border-pg-ciano` (tom de marca, sem contraste de texto exigido).
  3. `ConnectionOverlay`: fundo `bg-pg-nuvem/90` derivado de nuvem; mensagem em `text-pg-ardosia`.
  4. `TimeRangeSlider`: trilho `bg-pg-borda`, preenchimento `bg-pg-ciano`.
- **VERIFY**:
  - Desligar o backend: overlay de conexão aparece claro e legível
  - Login: título em ciano-fundo; campo de chave com anel de foco
- **Rollback**: restaurar os arquivos.

---

### Fase 4C: Limpeza e verificação

### Task T15: Remover `tokens.js` e legados sem uso ✅

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T9, T10, T11, T12, T13, T14
- **PRE-CONDITIONS**:
  - `grep -rn "theme/tokens\|t.colors\|var(--text\|var(--bg\|var(--border\|var(--accent\|var(--social\|var(--code" web/src` retorna vazio
- **OUTPUT**:
  1. Apagar `web/src/theme/tokens.js`.
  2. Remover de `index.css`: legados `--text`, `--text-h`, `--bg`, `--border`, `--code-bg`, `--accent*`, `--social-bg`, `--shadow`; substituir `h1,h2,code` por `var(--pg-ardosia)` / `var(--pg-ciano-fundo)` (mono em ciano-fundo, §5).
  3. `globalStyles.js`: remover import de `tokens`; usar as classes ou `var(--pg-*)` direto. Bloco `.counter`, `.hero` e `App.css` (template Vite) apagados se não tiverem uso (checar com grep).
  4. Remover `darkMode: 'class'` do `tailwind.config.js`.
- **WHY**: sem isso ficam duas fontes de verdade (`tokens.js` e `--pg-*`), que é o risco nº 1 da migração gradual.
- **VERIFY**:
  - `npm run build` passa
  - Telas idênticas às da T7/T18 (sem regressão visual)
- **Rollback**: restaurar `tokens.js` e `index.css`.

### Task T16: Resolver os 2 erros de lint pré-existentes ✅

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: None
- **INPUT**:
  - `TimeRangeSlider.jsx:17` — `unit` declarado e não usado.
  - `useConnection.jsx:43` — exporta função junto com componente (quebra fast refresh).
- **OUTPUT**: remover a variável não usada; mover a função exportada para `hooks/` ou tornar não exportada.
- **VERIFY**: `npm run lint` sem erros.
- **Rollback**: restaurar os dois arquivos.

### Task T17: Auditoria de contraste e a11y no tema claro ✅

- **Agent**: ux-accessibility-auditor
- **Priority**: P3
- **Complexity**: S
- **Dependencies**: T15
- **OUTPUT**: relatório com achados e trechos de correção, nas 4 lentes da skill.
- **VERIFY**:
  - Todo texto ≥ 4,5:1 (`neblina` 5,5:1, `ardosia` 14,5:1, `up` 4,7:1; `down` 5,0:1; ciano puro (`#1E9BB3`) só em ≥ 24 px, conforme §3.4)
  - Foco visível em todos os interativos
- **Rollback**: n/a (revisão).

### Task T18: Verificação final (checklist DESIGN.md §11) ✅

- **Agent**: qa-review
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T17
- **OUTPUT**: screenshots finais das 4 rotas, Login e modal, em 1280 px e 390 px, com o checklist da §11 marcado.
- **VERIFY**:
  - Os 8 itens da §11 marcados
  - Success Criteria do plano todos marcados
- **Rollback**: n/a (verificação).

## Risk Areas

| Risk                                                                                                      | Impact                                                           | Mitigation                                                                                                                                                                                                                               |
| --------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Cores fora do DESIGN.md**: `primaryHov` (hover do primário) e `warning` (amarelo) não têm token oficial | Viola §3.4 se inventar cor. Sem hover, o botão não dá feedback   | Hover por `brightness-90` em `--pg-ciano-fundo`, derivado da própria cor. **Aprovação pendente.** `warning` não tem uso de status previsto no §4: confirmar se algum componente o usa (grep antes da T2). Se usar, decidir com aprovação |
| **Violeta sai do app**: `#6c63ff` é a cor de marca atual                                                  | Mudança de identidade visível na primeira troca                  | É o que o DESIGN.md pede (primário = ciano-fundo). Confirmar antes da T1 que a troca de cor de marca é desejada                                                                                                                          |
| Tela mista durante a 4B (componentes migrados e não migrados)                                             | Aparência inconsistente em PR intermediário                      | 4A já deixa tudo claro. A 4B não muda aparência. Cada task tem VERIFY visual                                                                                                                                                             |
| `var()` não funciona em `rgba()` nem em comparações JS com `t.colors`                                     | Cor quebrada silenciosamente                                     | T2 faz `grep` por uso de `t.colors` fora de `style`. Usos com `rgba` tratados em T6                                                                                                                                                      |
| Remoção de `outline: none` sem substituto (T4) ao mexer no foco                                           | Usuário de teclado perde a posição                               | Anel de foco é requisito da T4, verificado na T7 e na T17                                                                                                                                                                                |
| `Button` descartar `type` (bug existente)                                                                 | Submit de Login/Settings pode não submeter em alguns navegadores | Forwardar `type` na T8                                                                                                                                                                                                                   |
| Contraste do badge sobre `*-bg` não testado                                                               | Texto de status ilegível                                         | T5 verifica contraste sobre o fundo do badge, não só sobre branco                                                                                                                                                                        |

## Rollback Strategy

1. Cada task é um commit. Commit separado para T1–T6 (troca de tema) e para cada componente da 4B.
2. Para desfazer a troca visual sem perder a refatoração: `git revert` apenas dos commits T1–T6. Os componentes da 4B continuam usando `pg-*` e funcionam sobre o tema escuro legado.
3. Se a 4B quebrar uma tela: `git revert` do commit do componente afetado. Os demais não dependem dele, exceto T9–T14 dependerem de T8 (`Button`).
4. Se o tema claro for rejeitado depois da T7: reverter a 4A inteira. O plano de 4B pode ser refeito sobre o tema escuro.

## Phase X: Verification

```bash
cd web

# Build (gera ../internal/handler/dist)
npm run build

# Lint
npm run lint

# Nenhuma cor de tema antigo
grep -rn "theme/tokens\|t\.colors\|#0f1117\|#1a1d27\|#6c63ff\|bg-black\|text-gray" src && echo "FALHOU" || echo "OK"

# Nenhum legado de variável
grep -rn "var(--text)\|var(--bg)\|var(--border)\|var(--accent" src && echo "FALHOU" || echo "OK"

# Nenhum style inline com cor
grep -rn "style={{" src | grep -i "color\|background\|border" && echo "REVISAR" || echo "OK"
```

## Cronograma Sugerido

| Fase                      | Tarefas | Complexity | Esforço | Target                               | Status      |
| ------------------------- | ------- | ---------- | ------- | ------------------------------------ | ----------- |
| 4A: Troca de tema         | T0–T7   | M          | ~2 h    | App claro, sem mudança de componente | ✅ Concluído |
| 4B: Componentes           | T8–T14  | M          | ~3 h    | Nenhum `t.colors` nos componentes    | ✅ Concluído |
| 4C: Limpeza e verificação | T15–T18 | M          | ~2 h    | Lint limpo, checklist §11            | ✅ Concluído |

## Status de Execução

### Concluído

- T0 a T17 (4A, 4B, 4C exceto itens listados abaixo)
- Correções pós-T18 (fora do plano original): menu mobile deslizante, foco do modal, botão de fechar maior, labels associados, textos em pt-BR

### Em andamento

- Nada.

### Pendente

- Redesenho da tela de login (adiado a pedido do usuário)

## Notes

- Esta fase não toca em ícones (Fase 2) nem em fontes (Fase 3). Ambas já estão no tema claro, porque a Fase 2 usa `currentColor` e a Fase 3 não depende de cor.
- O `Shell` mantém o símbolo `pingou-simbolo.svg`. O DESIGN.md §2 diz que esse símbolo é o padrão e funciona sobre fundo claro. Não há versão invertida necessária, porque o `pingou-simbolo-invertido.svg` é só para fundo escuro.
- Consultar o DESIGN.md §11 antes de cada commit de componente.

### Alternativas consideradas e descartadas

- **Migrar todos os componentes de uma vez, sem a camada 4A** → descartada: contraria o pedido de migrar aos poucos, e um PR grande demais é difícil de revisar.
- **Manter o tema escuro e propor tokens escuros (§10)** → descartada pelo pedido explícito de tema claro. A §10 continua válida se o tema escuro voltar a ser pauta.
- **Trocar `t.colors.x` por `var(--pg-...)` direto em cada componente, sem `tokens.js` como ponte** → descartada: 104 edições de uma vez antes da troca visual. A ponte deixa a troca em uma única alteração (T2).
- **Suporte a dois temas com `prefers-color-scheme`** → descartada: o DESIGN.md não define tema escuro oficial (§10).

## Resultado da execução

**Status:** 4A, 4B e 4C executadas. T18 concluída; restam as decisões de design abaixo. Nenhum commit foi feito nesta fase; a troca de tema está só no working tree.

### Por task

| Task | Resultado | Desvio do plano |
|---|---|---|
| T0 | Dispensada: Fases 1–3 já estavam commitadas pelo usuário (`21c1750`, `2e01959`, `0521ab5`) | — |
| T1 | `index.css` com tema claro; legados (`--text`, `--bg`, `--accent`…) apontam para `--pg-*`; bloco dark removido; `color-scheme: light` | `--sans`/`--heading` passaram a Inter/Rubik (titulos eram `system-ui`, violando §5) |
| T2 + T3 | `tokens.js` com `var(--pg-*)`; radius 8/12/12 | `warning` removido (sem uso). `primaryHov` mantido como `--pg-ciano-fundo` com `brightness-90` (ver decisões) |
| T4 | Input branco, borda `--pg-borda`; foco visível com `focus-visible` (outline 2px ciano) | Anel via `outline`, não `box-shadow`. `outline: none` removido |
| T5 | Badge em classes `pg-*`, pílula `rounded-full`, `font-medium` | Texto mudou para ardósia (ver decisões) |
| T6 | Overlay `rgba(30,43,51,0.5)` | — |
| T7 | Verificação visual em 1280 px: login, dashboard, incidentes, configurações | 390 px reprovado: layout quebrado (tratado depois, ver abaixo) |
| T8 | `Button` com variantes `primary`, `ghost`, `danger` (sólido) e `destructive` (contorno); `type` encaminhado, padrão `"button"` | Nomes `ghost`/`danger` mantidos, para não mexer em todos os call sites |
| T9 | `Shell` em classes `pg-*`; `NavLink` com estado ativo por classe | Marca trocada para SVG horizontal oficial (§2) |
| T10 | Dashboard e `MonitorCard` em classes; números com `tabular-nums` | Glow do ponto de status removido (§1 e §6) |
| T11 | `MonitorForm` e `Modal` em classes; switch com `role="switch"` | — |
| T12 | `Incidents` em classes | — |
| T13 | `Settings` em classes; "Baixar dump" virou `Button` ghost | Item selecionado de retenção em tom ciano-névoa, sem preenchimento (evita dois primários) |
| T14 | `Login`, `Spinner`, `ConnectionOverlay`, `TimeRangeSlider` em classes | Redesenho do login adiado a pedido do usuário |
| T15 | `tokens.js` e `App.css` removidos; `globalStyles.js` usa `var(--pg-*)`; legados fora do `index.css`; `h1`/`h2` usam `text-pg-h1`/`text-pg-h2`; `darkMode: 'class'` removido | — |
| T16 | Lint sem erros. `unit` não usado removido; hook `useConnection` separado em `hooks/connectionContext.js` | — |
| T17 | Auditoria de contraste e a11y feita (relatório na conversa). Achados tratados depois: mobile, modal, labels, badges, idioma | — |
| T18 | Checklist §11 verificado (ver abaixo). Login com `px-4` para respiro em 390 px | Sem screenshots novos; verificação por grep/código |

### Verificação

- `npm run lint` e `npm run build` passam.
- Greps do Phase X passam: nenhuma referência a `theme/tokens`, `t.colors`, `var(--text|--bg|--border|--accent)`, `bg-black`, `text-gray`, `font-semibold`, `font-bold` ou `style={{` com cor.
- Contraste medido com a fórmula WCAG (`node`), sobre os pares usados no app:

| Par | Razão | Situação |
|---|---|---|
| Ardósia sobre nuvem | 13,4:1 | ✅ |
| Neblina sobre nuvem | 5,1:1 | ✅ |
| Neblina sobre branco | 5,5:1 | ✅ |
| Branco sobre ciano-fundo (botão) | 5,9:1 | ✅ |
| Ciano-fundo sobre ciano-névoa (item ativo) | 5,1:1 | ✅ |
| Erro sobre branco | 5,0:1 | ✅ |
| Ardósia sobre `*-bg` (badges) | ≥ 12:1 | ✅ (após a troca de cor do texto) |
| Ícone UP sobre up-bg | 4,1:1 | ✅ (ícone é não-texto, mínimo 3:1) |
| Ciano puro (`#1E9BB3`) sobre branco | 3,3:1 | Não usado em texto pequeno |

- Screenshots headless (Chrome) em 1280 px: login, dashboard, incidentes, configurações e modal de criação. Em 390 px: dashboard (antes e depois da correção), menu aberto e modal de criação. Os screenshots ficam em `/tmp`, não versionados.
- Menu mobile e modal foram abertos por um mock que dispara o clique; não há teste automatizado de teclado.

### Checklist §11 (T18)

- Só tokens: nenhum hex/`rgba` fora de `--pg-*`; as 3 exceções são o overlay `rgba(30,43,51,0.5)` (ardósia a 50%, derivado). ✅
- Sem cor quente fora de status/erro. ✅
- Status com texto + ícone (Badge). ✅
- Botão primário `--pg-ciano-fundo`. ✅
- Só Rubik/Inter 400/500; sem `font-semibold`/`bold`. ✅
- Logo via SVG oficial. ✅
- Grafismos: nenhum. ✅
- Contraste ≥ 4,5:1 (tabela acima). ✅
- Limite: sem screenshots novos nesta rodada; `lint` e `build` passam.

### Correções pós-T18 (fora do plano original)

- **Layout mobile:** sidebar vira barra superior com botão de menu abaixo de `md`. O menu desliza por cima do conteúdo, fecha com Esc, overlay ou item, e devolve o foco ao botão. Cards de monitor empilham sem scroll horizontal.
- **Modal:** `role="dialog"` e `aria-modal`; foco inicial no primeiro campo; Tab fica preso no modal; foco volta ao disparador ao fechar; botão de fechar com área de ~36 px.
- **Labels:** todos os campos de formulário, sliders e webhook têm `label` associado; a retenção é um `fieldset` com `legend`. Erros têm `role="alert"`; "Salvo" tem `role="status"`.
- **Idioma:** `lang="pt-BR"`, título "Pingou" e toda a interface em português.

## Decisões tomadas

1. **Texto dos badges:** ardósia, ícone colorido (contraste AA). DESIGN.md §4 atualizado.
2. **Estados:** mantidos em inglês (`UP`/`DOWN`/`UNKNOWN`). Nota adicionada na §9.
3. **Hover do primário:** `brightness-90`, sem token novo. Registrado na §7.
4. **Login:** respiro em 390 px corrigido. Redesenho continua adiado.
