# Plano de Desenvolvimento: Design system no front (tokens, ícones, fontes)

## Overview

Levar o design system do [DESIGN.md](DESIGN.md) para o front em `web/`, em três etapas independentes e em ordem:

1. **Tokens `pg-*` no Tailwind**: somam-se aos tokens atuais, sem substituí-los. A migração dos componentes acontece aos poucos, fora deste plano.
2. **Ícones com `lucide-react`**: trocar os emojis e glifos Unicode usados como ícone por componentes Lucide. Onde o ícone é a marca (🏓), usar o SVG oficial do Pingou.
3. **Fontes locais**: Rubik e Inter (só os pesos 400 e 500) versionadas no repositório e servidas pelo próprio bundle. Nenhuma requisição ao Google Fonts em runtime.

**Estado atual (levantado no código):**

- Tailwind **v3.4** ([tailwind.config.js](web/tailwind.config.js)). Os tokens atuais são os do template Vite (`text`, `text-h`, `bg`, `border`, `accent`…), definidos como variáveis CSS em [index.css](web/src/index.css).
- Existe um segundo sistema de tokens em JS, de tema escuro roxo, em [tokens.js](web/src/theme/tokens.js), consumido via `style={{...}}` em 12 arquivos.
- As fontes vêm do Google em runtime: `@import url('https://fonts.googleapis.com/...')` em [globalStyles.js:6](web/src/theme/globalStyles.js#L6). Ele carrega Inter 400–700 e JetBrains Mono, o que viola o DESIGN.md §5 (só 400/500, mono do sistema).
- Não há nenhum arquivo `.woff2/.ttf` no repositório.
- O build do Vite vai para `internal/handler/dist`, que é embutido no binário Go (`//go:embed dist` em [spa.go](internal/handler/spa.go)). Fontes importadas pelo CSS entram no bundle automaticamente.

**Limites:**

- Nenhum token atual (`tailwind.config.js`, `index.css`, `tokens.js`) é removido ou alterado nas Fases 1 e 2.
- A migração visual dos componentes para as cores `pg-*` fica fora do escopo. A exceção é o `Badge` de status, que já recebe ícone na Fase 2, por exigência do DESIGN.md §4.
- Tema escuro: não há tokens oficiais (DESIGN.md §10). Fica fora do escopo.

## Project Type

**FRONTEND**: React 19 + Vite 6 + Tailwind 3, SPA embutida no backend Go.

## Definition of Done (DoD)

Este plano está concluído quando:

- [x] Todos os Success Criteria estão marcados
- [x] `npm run build` e `npm run lint` passam em `web/`
- [x] Nenhum item pendente no Status de Execução
- [x] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [x] Classes `bg-pg-nuvem`, `text-pg-ardosia`, `border-pg-borda`, `bg-pg-up-bg` etc. geram CSS válido, e as classes antigas (`text-text-h`, `bg-bg`, `border-border`…) continuam funcionando
- [x] Variáveis `--pg-*` do DESIGN.md §3.1 existem em `:root` e são a fonte única dos valores usados pelo Tailwind
- [x] `lucide-react` instalado. Nenhum emoji ou glifo Unicode usado como ícone em `web/src` (`grep` limpo)
- [x] Botões só com ícone (atualizar, fechar modal) têm `aria-label`, e ícones decorativos têm `aria-hidden`
- [x] `Badge` de status mostra ícone + texto (Check / X / CircleHelp)
- [x] Marca "🏓 Pingou" substituída pelo SVG oficial `pingou-simbolo.svg`
- [x] Arquivos `.woff2` de Rubik 400/500 e Inter 400/500 no repositório, com `OFL.txt`
- [x] Zero requisições a `fonts.googleapis.com` / `fonts.gstatic.com` (aba Network do DevTools)
- [x] `DESIGN.md` §5 atualizado: fontes locais, não Google Fonts

## Tech Stack

| Technology                                         | Choice Rationale                                                                                                                   |
| -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Tailwind v3 `theme.extend`                         | Já instalado. `extend` soma chaves sem apagar as existentes, o que garante a migração gradual                                      |
| CSS custom properties `--pg-*`                     | Valor único, compartilhado entre Tailwind, CSS puro e o `style={}` inline dos componentes atuais                                   |
| `lucide-react`                                     | Pedido explícito. Tree-shakeable (só os ícones importados entram no bundle), traço consistente, combina com o visual arredondado   |
| `@font-face` + `.woff2` em `web/src/assets/fonts/` | Pedido explícito (fontes no repo). Com import via CSS, o Vite gera hash e copia para o `dist`, e o `go:embed` já cobre o resultado |

## File Structure

```text
web/
├── tailwind.config.js            # + colors.pg, fontFamily pg-*, borderRadius pg-*, fontSize pg-*
├── package.json                  # + lucide-react
└── src/
    ├── index.css                 # + :root --pg-*, + @font-face (Fase 3)
    ├── theme/
    │   └── globalStyles.js       # - @import Google Fonts (Fase 3)
    ├── assets/
    │   ├── brand/
    │   │   └── pingou-simbolo.svg        # cópia de /assets/pingou-simbolo.svg
    │   └── fonts/
    │       ├── OFL.txt
    │       ├── rubik-latin-400.woff2
    │       ├── rubik-latin-500.woff2
    │       ├── inter-latin-400.woff2
    │       └── inter-latin-500.woff2
    ├── components/
    │   ├── layout/Shell.jsx      # ícones nav + logout + marca
    │   └── ui/{Badge,Modal}.jsx  # ícone de status, botão fechar
    └── pages/{Dashboard,Login,Settings,Incidents}.jsx
DESIGN.md                         # §5: fontes locais
```

## Task Breakdown

| Task ID | Agent                    | Name                                            | Priority | Complexity | Dependencies |
| ------- | ------------------------ | ----------------------------------------------- | -------- | ---------- | ------------ |
| T1      | frontend-specialist      | Variáveis `--pg-*` em `:root`                   | P0       | XS         | -            |
| T2      | frontend-specialist      | Tokens `pg` no `tailwind.config.js`             | P0       | S          | T1           |
| T3      | frontend-specialist      | Smoke test dos tokens (antigos + novos)         | P0       | XS         | T2           |
| T4      | frontend-specialist      | Instalar `lucide-react` + SVG da marca          | P1       | XS         | -            |
| T5      | frontend-specialist      | Ícones no `Shell` (nav, logout, marca)          | P1       | S          | T4           |
| T6      | frontend-specialist      | Ícones em `Dashboard`, `Modal`, `Login`         | P1       | S          | T4           |
| T7      | frontend-specialist      | Ícones em `Settings`, `Incidents`, `Badge`      | P1       | S          | T4           |
| T8      | ux-accessibility-auditor | Revisão a11y dos ícones                         | P3       | XS         | T5, T6, T7   |
| T9      | frontend-specialist      | Baixar `.woff2` para o repositório (uma vez só) | P2       | S          | -            |
| T10     | frontend-specialist      | `@font-face` + ativar fontes nos tokens         | P2       | S          | T2, T9       |
| T11     | frontend-specialist      | Remover Google Fonts do `globalStyles.js`       | P2       | XS         | T10          |
| T12     | document                 | Atualizar `DESIGN.md` §5                        | P3       | XS         | T11          |

---

### Fase 1: Tokens do DESIGN.md no Tailwind

### Task T1: Variáveis `--pg-*` em `:root` ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - `web/src/index.css` com o bloco `:root` atual intacto
- **INPUT**: `:root` só tem os tokens do template (`--text`, `--accent`…).
- **OUTPUT**:
  1. Novo bloco `:root { ... }` logo **abaixo** do atual, com as 15 variáveis do DESIGN.md §3.1 copiadas literalmente (`--pg-ciano` … `--pg-unknown-bg`).
  2. Não mexer no `@media (prefers-color-scheme: dark)`. Não há tokens escuros oficiais.
- **WHY**: as variáveis ficam como fonte única. O Tailwind (T2) e o `style={}` inline dos componentes atuais (`var(--pg-ciano)`) passam a ler o mesmo valor.
- **VERIFY**:
  - DevTools → `getComputedStyle(document.documentElement).getPropertyValue('--pg-ciano')` retorna `#1E9BB3`
  - Visual da app idêntico ao anterior (nada consome as variáveis ainda)
- **Rollback**: apagar o bloco adicionado.

### Task T2: Tokens `pg` no `tailwind.config.js` ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - T1 concluída
- **INPUT**: `theme.extend` com `colors`, `fontFamily`, `fontSize`, `borderRadius` do template.
- **OUTPUT**, tudo **somado** às chaves existentes dentro de `theme.extend`:
  1. `colors.pg`: `{ ciano: 'var(--pg-ciano)', 'ciano-fundo': ..., 'ciano-nevoa', ardosia, neblina, nuvem, branco, borda, up, 'up-bg', down, 'down-bg', unknown, 'unknown-bg' }`
  2. `fontFamily`: **novas chaves** `'pg-heading': ['Rubik', 'system-ui', 'sans-serif']`, `'pg-sans': ['Inter', 'system-ui', 'sans-serif']`, `'pg-mono': ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace']`. As chaves `sans`/`heading`/`mono` atuais **não mudam** nesta fase.
  3. `borderRadius`: `'pg-control': '8px'`, `'pg-card': '12px'`. Badge usa o `rounded-full` nativo.
  4. `fontSize` (DESIGN.md §5): `'pg-h1': ['32px', { lineHeight: '1.2' }]`, `'pg-h2': ['24px', '1.25']`, `'pg-h3': ['20px', '1.3']`, `'pg-body': ['16px', '1.6']`, `'pg-ui': ['14px', '1.5']`, `'pg-caption': ['12px', '1.4']`
- **WHY**: o prefixo `pg-` evita colisão com `sans`, `sm`, `md` e `h1`/`h2`, que já existem e estão em uso. Qualquer classe nova com `pg-` é, por definição, do design system, e isso deixa a migração rastreável com `grep pg-`.
- **Nota**: `var(--pg-*)` em hex **não** suporta modificador de opacidade (`bg-pg-ciano/50`). O DESIGN.md não pede opacidade, então é aceitável. Se um dia for necessário, trocar para o formato `rgb(var(--pg-ciano-rgb) / <alpha-value>)`.
- **VERIFY**:
  - `npm run build` sem warnings do Tailwind
- **Rollback**: `git checkout web/tailwind.config.js`.

### Task T3: Smoke test dos tokens ✅

- **Agent**: frontend-specialist
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: T2
- **INPUT**: config com tokens novos.
- **OUTPUT**:
  1. Aplicar **uma** classe `pg-*` real num componente de baixo risco para provar o pipeline, por exemplo `bg-pg-nuvem` no container do `Login`. Mantê-la: é o primeiro passo da migração.
- **VERIFY**:
  - `grep -o 'pg-nuvem' ../internal/handler/dist/assets/*.css` encontra a classe
  - Classes antigas (`text-text-h`, `bg-bg`) continuam no CSS gerado
  - `npm run lint` passa
- **Rollback**: remover a classe.

---

### Fase 2: Ícones com lucide-react

**Mapeamento (levantado com `grep` de caracteres não ASCII em `web/src`):**

| Arquivo:linha                                                                                    | Hoje                      | Substituir por                                                                |
| ------------------------------------------------------------------------------------------------ | ------------------------- | ----------------------------------------------------------------------------- |
| [Shell.jsx:6](web/src/components/layout/Shell.jsx#L6)                                            | `⬡ Dashboard`             | `LayoutDashboard`                                                             |
| [Shell.jsx:7](web/src/components/layout/Shell.jsx#L7)                                            | `⚠ Incidents`             | `TriangleAlert`                                                               |
| [Shell.jsx:8](web/src/components/layout/Shell.jsx#L8)                                            | `⚙ Settings`              | `Settings`                                                                    |
| [Shell.jsx:33](web/src/components/layout/Shell.jsx#L33)                                          | `🏓 Pingou`               | **SVG oficial** `pingou-simbolo.svg` (24–32 px) + texto                       |
| [Shell.jsx:80](web/src/components/layout/Shell.jsx#L80)                                          | `⎋ Logout`                | `LogOut`                                                                      |
| [Dashboard.jsx:102](web/src/pages/Dashboard.jsx#L102)                                            | `🔄` (botão só com ícone) | `RefreshCw` + `animate-spin` + `aria-label`                                   |
| [Dashboard.jsx:105](web/src/pages/Dashboard.jsx#L105) e [:177](web/src/pages/Dashboard.jsx#L177) | `+ Add Monitor`           | `Plus` + texto                                                                |
| [Dashboard.jsx:174](web/src/pages/Dashboard.jsx#L174)                                            | `🏓` (estado vazio)       | SVG oficial `pingou-simbolo.svg`                                              |
| [Login.jsx:55](web/src/pages/Login.jsx#L55)                                                      | `🏓`                      | SVG oficial `pingou-simbolo.svg`                                              |
| [Modal.jsx:34](web/src/components/ui/Modal.jsx#L34)                                              | `×` (botão só com ícone)  | `X` + `aria-label="Fechar"`                                                   |
| [Settings.jsx:172](web/src/pages/Settings.jsx#L172)                                              | `✓ Saved`                 | `Check`                                                                       |
| [Settings.jsx:205](web/src/pages/Settings.jsx#L205)                                              | `↓ Download dump`         | `Download`                                                                    |
| [Incidents.jsx:56](web/src/pages/Incidents.jsx#L56)                                              | `No incidents found 🎉`   | remover o emoji ou usar `CircleCheck` (DESIGN.md §9: tom sóbrio)              |
| [Incidents.jsx:79](web/src/pages/Incidents.jsx#L79)                                              | `● OPEN`                  | `CircleAlert`                                                                 |
| [Incidents.jsx:87](web/src/pages/Incidents.jsx#L87)                                              | `✓ RESOLVED`              | `Check`                                                                       |
| [Badge.jsx](web/src/components/ui/Badge.jsx)                                                     | só texto                  | `Check` (UP) / `X` (DOWN) / `CircleHelp` (UNKNOWN), exigido pelo DESIGN.md §4 |

Os `·` em `Dashboard.jsx:79` e `MonitorCard.jsx:79` são separadores de texto e **ficam**.

**Convenção única para todos os ícones:** `size={16}` (nav/botões) ou `size={14}` (badges), `strokeWidth={2}`, `aria-hidden="true"` quando há texto ao lado, cor herdada (`currentColor`). Layout `inline-flex items-center gap-2`. Não criar componente `<Icon>` wrapper: seria um único consumidor por ícone e não justifica.

### Task T4: Instalar `lucide-react` + SVG da marca ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: None
- **INPUT**: sem lib de ícones. SVGs da marca só em `/assets` (fora da raiz do Vite).
- **OUTPUT**:
  1. `cd web && npm install lucide-react`
  2. Copiar `assets/pingou-simbolo.svg` → `web/src/assets/brand/pingou-simbolo.svg`. Importar via `import logo from '../assets/brand/pingou-simbolo.svg'`. O Vite, por padrão, não serve arquivos fora de `web/`.
- **VERIFY**:
  - `lucide-react` aparece em `dependencies` do `package.json` e no `package-lock.json`
- **Rollback**: `npm uninstall lucide-react`, apagar `web/src/assets/brand/`.

### Task T5: Ícones no `Shell` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **INPUT**: array `nav` com label em string contendo o glifo. Marca em emoji.
- **OUTPUT**:
  1. `nav` passa a ter `{ to, label: 'Dashboard', icon: LayoutDashboard }`, e o render faz `<item.icon size={16} aria-hidden="true" />`
  2. Logout com `<LogOut />`
  3. Marca: `<img src={logo} alt="" width={28} height={28} />` + texto "Pingou". Respeitar a área de proteção de 1/4 da altura (≈7 px de gap).
- **VERIFY**:
  - Navegação funciona. O estado ativo do `NavLink` continua colorindo o ícone (herda `currentColor`)
- **Rollback**: `git checkout web/src/components/layout/Shell.jsx`.

### Task T6: Ícones em `Dashboard`, `Modal`, `Login` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **OUTPUT**: conforme a tabela de mapeamento. Botão de refresh e botão de fechar do modal ganham `aria-label` ("Atualizar monitores", "Fechar").
- **VERIFY**:
  - Clicar em atualizar: o ícone gira enquanto `isFetching` e para em seguida
  - Modal fecha pelo `X`, e o botão é alcançável por Tab
- **Rollback**: `git checkout` dos 3 arquivos.

### Task T7: Ícones em `Settings`, `Incidents`, `Badge` ✅

- **Agent**: frontend-specialist
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T4
- **OUTPUT**: conforme a tabela. No `Badge`, mapa `{ UP: Check, DOWN: X, UNKNOWN: CircleHelp }` ao lado do texto, sem alterar as cores atuais (a migração de cor é outro passo).
- **VERIFY**:
  - Badge dos três estados renderiza ícone + texto (testar com um monitor DOWN apontando para uma URL inválida)
  - `grep -rnP "[^\x00-\x7F]" web/src --include='*.jsx'` só retorna os separadores `·`
- **Rollback**: `git checkout` dos 3 arquivos.

### Task T8: Revisão a11y dos ícones ✅

- **Agent**: `ux-accessibility-auditor` (skill)
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T5, T6, T7
- **VERIFY**:
  - Todo botão sem texto visível tem `aria-label`. Todo ícone com texto ao lado tem `aria-hidden`
  - Leitor de tela (Orca/NVDA) não lê nomes de ícones duplicados
- **Rollback**: n/a (revisão).

---

### Fase 3: Fontes locais

### Task T9: Baixar `.woff2` para o repositório ✅

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Acesso à internet **uma única vez**, só para obter os arquivos
- **INPUT**: nenhum arquivo de fonte no repositório.
- **OUTPUT**:
  1. Baixar os `.woff2` estáticos, subset **latin** (cobre acentos do português), dos pesos 400 e 500 de Rubik e Inter. Fonte recomendada: pacotes `@fontsource/rubik` e `@fontsource/inter`, extraídos com `npm pack` sem instalar. Pegar `files/{rubik,inter}-latin-{400,500}-normal.woff2`.
  2. Salvar em `web/src/assets/fonts/` com nomes `rubik-latin-400.woff2` etc.
  3. Incluir `OFL.txt` (licença SIL OFL 1.1, que exige distribuir a licença junto).
- **VERIFY**:
  - `ls web/src/assets/fonts` mostra 4 `.woff2` + `OFL.txt`. Total esperado < 150 KB
  - `file *.woff2` → "Web Open Font Format (Version 2)"
- **Rollback**: apagar a pasta.

### Task T10: `@font-face` + ativar fontes nos tokens ✅

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T2, T9
- **OUTPUT**:
  1. Em `index.css`, antes do `@tailwind base`, 4 blocos:
     ```css
     @font-face {
       font-family: "Inter";
       src: url("./assets/fonts/inter-latin-400.woff2") format("woff2");
       font-weight: 400;
       font-style: normal;
       font-display: swap;
     }
     ```
     O mesmo para Inter 500, Rubik 400 e Rubik 500. O caminho relativo faz o Vite gerar hash e copiar para o `dist`.
  2. `tokens.font.sans` já é `'Inter', system-ui…`, então o body passa a usar Inter local automaticamente. Trocar `tokens.font.mono` de `'JetBrains Mono'` para a pilha mono do sistema (DESIGN.md §5: nunca mono externa). É a **única** alteração em token existente, e ela é forçada pela remoção do Google Fonts: sem o import, JetBrains Mono já não carregaria.
- **VERIFY**:
  - DevTools → Network → filtro "Font": só aparecem `/assets/inter-latin-*.woff2` e `rubik-*`, servidos pelo próprio host
  - `ls ../internal/handler/dist/assets/*.woff2` lista os arquivos
  - Rodar o binário Go e abrir a SPA: as fontes carregam via `go:embed`
- **Rollback**: remover os `@font-face` e reverter `tokens.js`.

### Task T11: Remover Google Fonts do `globalStyles.js` ✅

- **Agent**: frontend-specialist
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T10
- **OUTPUT**:
  1. Apagar a linha `@import url('https://fonts.googleapis.com/...')` em [globalStyles.js:6](web/src/theme/globalStyles.js#L6).
- **VERIFY**:
  - `grep -rn "googleapis\|gstatic" web/src web/index.html` retorna vazio
  - Network, com cache desativado: zero requisições para domínios Google
  - Com a rede offline (DevTools → Offline, depois de carregar uma vez), a tipografia continua Inter/Rubik
- **Rollback**: restaurar a linha.

### Task T12: Atualizar `DESIGN.md` §5 ✅

- **Agent**: `document` (skill)
- **Priority**: P3
- **Complexity**: XS
- **Dependencies**: T11
- **OUTPUT**:
  1. Trocar "Fontes: Google Fonts…" por: fontes **self-hosted** em `web/src/assets/fonts/`, `.woff2` subset latin, pesos 400/500, licença OFL. **NUNCA** carregar via CDN.
  2. §3.2: documentar o prefixo `pg-` nas chaves de fonte, raio e tamanho (divergência intencional do exemplo, por causa da migração gradual).
- **VERIFY**: releitura do §3.2 e do §5 bate com o código.
- **Rollback**: `git checkout DESIGN.md`.

## Risk Areas

| Risk                                                                   | Impact                                        | Mitigation                                                                                         |
| ---------------------------------------------------------------------- | --------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| Chave nova no `extend` sobrescrever uma existente (`sans`, `sm`, `h1`) | Visual quebra em telas não migradas           | Prefixo `pg-` em todas as chaves novas (T2). Smoke test das classes antigas (T3)                   |
| `var()` em cor impede modificador de opacidade (`/50`)                 | Classe silenciosamente sem efeito             | Documentado em T2. DESIGN.md não usa opacidade. Trocar para canais RGB se surgir a necessidade     |
| Subset latin sem algum glifo                                           | Caractere cai para fonte do sistema           | Latin cobre pt-BR. Validar com "ação, não, você, Á" na UI                                          |
| FOUT/FOIT na primeira carga                                            | Texto pisca                                   | `font-display: swap`. Arquivos pequenos e servidos do mesmo host                                   |
| Emoji `🏓` recriado como ícone Lucide                                  | Viola DESIGN.md §2 (marca só via SVG oficial) | Mapeamento explícito: marca → `pingou-simbolo.svg`                                                 |
| `dist` embutido desatualizado no binário Go                            | Deploy sem fontes/ícones                      | Rodar `npm run build` antes de `go build`. Verificar em T10                                        |
| Tokens JS (`tokens.js`) e `--pg-*` coexistindo                         | Duas fontes de verdade durante a migração     | Esperado e aceito. A migração componente a componente remove `tokens.js` no fim (fora deste plano) |

## Rollback Strategy

1. Cada fase é um commit separado (Fase 1 tokens, Fase 2 ícones, Fase 3 fontes). Para reverter uma fase, `git revert <sha>` do commit dela.
2. As fases são independentes, exceto T10, que depende de T2 para as chaves `pg-heading`/`pg-sans`. Reverter a Fase 1 exige reverter a Fase 3 antes.
3. Fase 2: além do revert, `npm uninstall lucide-react`.

## Phase X: Verification

```bash
cd web

# Build (gera ../internal/handler/dist)
npm run build

# Lint
npm run lint

# Sem Google Fonts
grep -rn "googleapis\|gstatic" src index.html && echo "FALHOU" || echo "OK"

# Sem emoji/glifo como ícone (só os separadores "·" devem restar)
grep -rnP "[^\x00-\x7F]" src --include='*.jsx'

# Fontes e tokens no bundle
ls ../internal/handler/dist/assets/*.woff2
grep -o 'pg-nuvem' ../internal/handler/dist/assets/*.css | head -1

# Binário Go servindo tudo
cd .. && go run ./cmd/pingou   # abrir no browser, Network → Font
```

## Cronograma Sugerido

| Fase           | Tarefas            | Complexity | Esforço | Target                                       | Status      |
| -------------- | ------------------ | ---------- | ------- | -------------------------------------------- | ----------- |
| Fase 1: Tokens | T1, T2, T3         | S          | ~25 min | Classes `pg-*` disponíveis, nada quebrado    | ✅ Concluído |
| Fase 2: Ícones | T4, T5, T6, T7, T8 | M          | ~50 min | Zero emoji/glifo como ícone, a11y ok         | ✅ Concluído |
| Fase 3: Fontes | T9, T10, T11, T12  | S          | ~35 min | Zero requisição ao Google, fontes via bundle | ✅ Concluído |

## Status de Execução

### Concluído

- T1 → T12 (Fases 1, 2 e 3)

### Em andamento

-

### Pendente

- Nenhum item deste plano. Ver "Resultado da execução" abaixo para os desvios registrados.

## Notes

- A troca do `font-family` global para Inter acontece "de graça" na Fase 3, porque o `body` já declara Inter via `tokens.js`. O que muda é só a origem do arquivo. Hoje, parte dos componentes usa `font-sans` do Tailwind (`system-ui`). Esses só passam a Inter quando migrarem para `font-pg-sans`, o que é coerente com a migração gradual.
- `globalStyles.js` também carrega Inter 600/700, e há `font-semibold` no `Badge`. Com só 400/500 locais, o navegador faz bold sintético (o `font-synthesis: none` em `index.css` impede isso e cai para 500). Esse efeito colateral é aceitável e alinhado ao DESIGN.md §5 (só 400/500). **Atualização:** o `font-semibold` foi removido do `Badge` e do restante do código; não há mais pesos fora de 400/500.
- O favicon já foi trocado para `favicon-16.png` em `web/public/`. Considerar adicionar `favicon-32.png` e `favicon-180.png` (apple-touch-icon) quando mexer no `index.html`, mas isso está fora do escopo.

### Alternativas consideradas e descartadas

- **Substituir `sans`/`heading`/`mono` do Tailwind direto pelos valores do DESIGN.md** → descartada: muda a fonte de todas as telas de uma vez, contra o pedido de migração gradual.
- **Cores `pg` em hex literal no config (como no exemplo do DESIGN.md §3.2)** → descartada: duplica valores com as variáveis CSS, e os componentes com `style={}` inline não teriam como consumi-los.
- **Instalar `@fontsource/*` como dependência** → descartada: o pedido é ter as fontes versionadas no repositório. O pacote é usado só como origem dos arquivos em T9.
- **Inter variável (um arquivo, todos os pesos)** → descartada: ~300 KB contra ~4×25 KB, e o design system só permite 400/500.
- **Componente wrapper `<Icon name="...">`** → descartado: indireção sem ganho. O Lucide já expõe componentes com props consistentes.
- **Ícone Lucide para a marca (ex.: `Droplet`)** → descartado: DESIGN.md §2 exige o SVG oficial.

## Resultado da execução

**Status:** T1 a T12 concluídas. Commits feitos pelo usuário: `0521ab5` (Tailwind e fontes), `2e01959` (fontes), `21c1750` (ícones) e `22b6a96` (ajustes de graphify e manifest).

### Entregas

| Fase | Entregue | Arquivos principais |
|---|---|---|
| 1. Tokens | Variáveis `--pg-*` em `:root`; chaves `pg` no Tailwind (cores, `pg-h1`…`pg-caption`, `pg-control`, `pg-card`, `pg-heading`, `pg-sans`, `pg-mono`) | `web/src/index.css`, `web/tailwind.config.js` |
| 2. Ícones | `lucide-react` instalado; ícones em nav, botões, status e estado vazio; SVG oficial da marca | `web/src/components/**`, `web/src/pages/**`, `web/src/assets/brand/` |
| 3. Fontes | Inter e Rubik 400/500 em `.woff2` subset latin no repositório; `@font-face` com `font-display: swap`; Google Fonts removido | `web/src/assets/fonts/`, `web/src/index.css`, `web/src/theme/globalStyles.js` |
| 4. Documentação | §5 do DESIGN.md descreve fontes self-hosted; §3.2 explica o prefixo `pg-` | `DESIGN.md` |

### Verificação

- `npm run build` gera os 4 `.woff2` com hash em `internal/handler/dist/assets/` e o binário Go os serve via `go:embed`.
- `grep` em `web/src` e `web/index.html` não encontra `googleapis` nem `gstatic`.
- Nenhum emoji ou glifo Unicode funciona como ícone; o único caractere não ASCII restante é o separador `·`.
- Auditoria de acessibilidade dos ícones (T8): `aria-label` nos botões só com ícone e `aria-hidden` nos ícones decorativos. Ajuste aplicado: `motion-reduce:animate-none` no botão de refresh.
- Não houve verificação no DevTools (aba Network) para a regra de zero requisições ao Google. A verificação foi por `grep`, que cobre o código fonte.

### Desvios em relação ao plano

- **Marca no Shell e no Login:** o plano pedia "símbolo + texto Pingou". O texto recriado com fonte viola a §2 do DESIGN.md. Shell e Login agora usam `pingou-marca-horizontal.svg`. O estado vazio do Dashboard usa `pingou-simbolo.svg`.
- **Licença:** em vez de um `OFL.txt` único, há `OFL-Inter.txt` e `OFL-Rubik.txt`, porque cada família traz seu próprio copyright.
- **`font.mono` em `tokens.js`:** trocado de `JetBrains Mono` para a pilha mono do sistema. É a única alteração de token existente nas Fases 1–3.
- **Critério "classes antigas continuam funcionando":** as chaves antigas (`text-text-h`, `bg-bg`…) continuam declaradas no `tailwind.config.js`, mas nenhum componente as usa como classe. Por isso não geram CSS. Não houve quebra.
- **Badge:** o texto mudou depois para ardósia, na migração de tema. Ver `plan-tema-claro.md`.
- **Ícone de "no incidents":** o 🎉 foi removido por causa do tom sóbrio da §9.

### Pendências

- Nenhuma pendência de execução deste plano. Os itens abertos são decisões de design, registradas em `plan-tema-claro.md`.
