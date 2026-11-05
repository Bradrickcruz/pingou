# Plano de Desenvolvimento: Modo escuro (dark mode) no painel

## Overview

**Contexto.** O painel (`web/`) só tem tema claro. O `DESIGN.md` v1.1 já define a paleta escura (§3.6: tokens, mapa de papéis claro → escuro, status, contraste e regras). O `IDEAS.md` (task "Modo escuro") define o produto.

**Objetivo.** A pessoa escolhe o tema do painel em **Configurações**, com três opções:

| Opção   | O que faz                                            |
| ------- | ---------------------------------------------------- |
| Claro   | Sempre claro                                         |
| Escuro  | Sempre escuro                                        |
| Sistema | Segue o sistema (`prefers-color-scheme`). **Padrão** |

**Limites (boundary).**

- A escolha vale **por navegador** (`localStorage`), não por instalação. Não vai para o backend nem para `settings`.
- **Sem atalho no cabeçalho do painel.** O seletor fica só em Configurações (decisão de 2026-11-04 no `IDEAS.md`).
- Vale para todas as telas: Painel, Incidentes, Detalhe, Manutenções, Configurações, Login e modais, desktop e celular.
- O tema claro **não muda** (nem um pixel).

**Restrições.**

- Só tokens do `DESIGN.md` §3.6. Cor nova só com aprovação do dono do projeto (Bryan).
- Aplicar o tema **antes da primeira pintura** (sem piscar claro → escuro).
- Contraste WCAG AA nos dois temas. Status sempre com texto + ícone.
- O sistema de temas deve servir também para a futura página pública de status (ideia 3), que terá o próprio botão de tema.

## Project Type

**FRONTEND** — só `web/` (React 19 + Vite 6 + Tailwind 3.4). Nenhuma mudança no Go nem na API. O binário continua servindo o build via `go:embed`.

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os itens do Success Criteria estão marcados
- [ ] Build, lint e testes passam (`npm run build`, `npm run lint`, `npm test`)
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas (capturas dos 2 temas × telas × desktop/celular)

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Configurações tem um bloco "Aparência" com Claro, Escuro e Sistema; Sistema vem marcado na primeira visita
- [ ] Trocar a opção aplica o tema **na hora**, sem passar pelo botão "Salvar" das Configurações
- [ ] A escolha sobrevive a recarregar a página e vale só naquele navegador
- [ ] Em "Sistema", mudar o tema do sistema operacional troca o painel sem recarregar
- [ ] Nenhum piscar de tema claro ao abrir qualquer tela no escuro (incluindo Login)
- [ ] Nenhum atalho de tema no cabeçalho (desktop `Shell` e `MobileHeader`)
- [ ] Todas as telas e modais usam só tokens de papel; nenhum hexadecimal solto nos componentes
- [ ] Tema claro idêntico ao atual (comparação de capturas antes/depois)
- [ ] Logo negativo / símbolo invertido no escuro
- [ ] Gráfico de resposta, barras de histórico e selos (UP, DOWN, UNKNOWN, MAINTENANCE) com as cores do §3.6.3
- [ ] Contraste AA conferido nos dois temas
- [ ] `DESIGN.md` atualizado com os papéis novos aprovados

## Tech Stack

| Technology                                               | Choice Rationale                                                                                       |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| CSS custom properties                                    | As classes `pg-*` do Tailwind já apontam para `var(--pg-*)`. Basta trocar o valor da variável por tema |
| Atributo `data-theme` no `<html>`                        | Um único ponto de troca. Funciona antes do React carregar (script inline)                              |
| Tailwind `darkMode: ["selector", '[data-theme="dark"]']` | Só para exceções raras (ex.: trocar o logo). Não é o caminho principal                                 |
| `localStorage` (`pingou_theme`)                          | Escolha por navegador, como pede o `IDEAS.md`. Mesmo padrão do `pingou_api_key`                        |
| `matchMedia("(prefers-color-scheme: dark)")`             | Opção "Sistema" e troca ao vivo                                                                        |
| `node --test`                                            | Já usado em `src/lib/*.test.js`. Testa a lógica pura do tema                                           |
| Playwright (pré-instalado)                               | Capturas para comparar tema claro antes/depois e revisar o escuro                                      |

## Decisão de arquitetura (o ponto mais importante)

**Problema.** Hoje os componentes usam tokens de **cor** (`text-pg-ardosia`, `bg-pg-branco`). Mas no escuro a troca depende do **papel**, não da cor:

- `text-pg-ardosia` (119 usos) é texto → vira `--pg-nuvem`.
- `bg-pg-ardosia` (8 usos) é fundo invertido (tooltip, aba ativa) → **não** pode virar `--pg-nuvem` do mesmo jeito.
- `text-pg-branco` (14 usos) às vezes é texto de botão (→ `--pg-ardosia` no escuro), às vezes texto sobre ardósia.

Ou seja: só redefinir `--pg-ardosia` no escuro quebra a tela.

**Solução escolhida: tokens de papel.** Uma camada nova de variáveis, uma por papel do §3.6.2. Cada uma tem valor no claro e no escuro. Os componentes passam a usar o papel.

```css
:root,
[data-theme="light"] {
  --pg-cor-fundo: var(--pg-nuvem);
  --pg-cor-superficie: var(--pg-branco);
  --pg-cor-linha: var(--pg-borda);
  --pg-cor-texto: var(--pg-ardosia);
  --pg-cor-texto-2: var(--pg-neblina);
  --pg-cor-link: var(--pg-ciano-fundo);
  /* ...demais papéis (ver T2) */
}
[data-theme="dark"] {
  --pg-cor-fundo: var(--pg-ardosia);
  --pg-cor-superficie: var(--pg-ardosia-clara);
  --pg-cor-linha: var(--pg-ardosia-linha);
  --pg-cor-texto: var(--pg-nuvem);
  --pg-cor-texto-2: var(--pg-neblina-clara);
  --pg-cor-link: var(--pg-ciano-luz);
}
```

Os tokens crus (`--pg-ardosia`, `--pg-up`...) **continuam constantes**, como manda o §3.6.1.

Nomes (`pg-cor-*` e cada papel) são proposta: confirmar no T1.

## File Structure

```text
web/
├── index.html                       # + script inline de tema, + <meta name="color-scheme">
├── tailwind.config.js               # + cores de papel (pg.fundo, pg.texto...), + darkMode selector
└── src/
    ├── index.css                    # + tokens §3.6.1, + tokens de papel claro/escuro
    ├── theme/globalStyles.js        # passa a usar tokens de papel
    ├── lib/
    │   ├── theme.js                 # NOVO: ler/salvar/resolver/aplicar tema
    │   └── theme.test.js            # NOVO
    ├── hooks/
    │   └── useTheme.js              # NOVO: estado do tema + escuta do sistema
    ├── assets/brand/
    │   └── pingou-simbolo-invertido.svg   # NOVO (copiado de /assets)
    ├── components/settings/
    │   └── ThemeCard.jsx            # NOVO: bloco "Aparência"
    └── components/**, pages/**      # migração para tokens de papel
scripts/
└── check-tokens.sh                  # NOVO: falha se componente usar token cru ou hex solto
```

## Task Breakdown

| Task ID | Agent          | Name                                         | Priority | Complexity | Dependencies |
| ------- | -------------- | -------------------------------------------- | -------- | ---------- | ------------ |
| T0      | docs           | Corrigir `DESIGN.md` duplicado               | P0       | XS         | -            |
| T1      | design (Bryan) | Fechar lacunas de design do escuro           | P0       | S          | T0           |
| T2      | frontend       | Camada de tokens de papel                    | P0       | S          | T1           |
| T3      | frontend       | Núcleo do tema (lib, hook, script inline)    | P0       | S          | T2           |
| T4      | frontend       | Migrar base: `ui/*`, layout, estilos globais | P0       | M          | T2           |
| T5      | frontend       | Migrar telas, modais e cores soltas          | P0       | L          | T4           |
| T6      | frontend       | Gráficos SVG e barras de histórico           | P1       | S          | T2           |
| T7      | frontend       | Logos e toasts no escuro                     | P1       | XS         | T2           |
| T8      | frontend       | Bloco "Aparência" em Configurações           | P0       | S          | T3           |
| T9      | frontend       | Trava contra regressão de tokens             | P1       | XS         | T5           |
| T10     | qa             | Verificação visual, contraste e piscar       | P0       | M          | T5–T9        |
| T11     | docs           | Documentação final                           | P2       | XS         | T10          |

### Task T0: Corrigir `DESIGN.md` duplicado

- **Agent**: docs
- **Priority**: P0
- **Complexity**: XS
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Nenhuma
- **INPUT**: O `DESIGN.md` tem a v1.1 inteira e, colada no fim (linha ~345 em diante), a v1.0 antiga, a partir do meio de uma frase ("tes, pesos ou sombras."). A v1.0 diz que "não há tokens oficiais" de escuro, o que contradiz a v1.1.
- **OUTPUT**:
  1. Remover o trecho da v1.0 duplicado. Fica só a v1.1.
- **VERIFY**:
  - `grep -c "^## 3. Cores" DESIGN.md` retorna `1`
  - Arquivo termina no checklist §11 da v1.1

### Task T1: Fechar lacunas de design do escuro

- **Agent**: design (aprovação do Bryan)
- **Priority**: P0 — **bloqueia o código**, como decidido no `IDEAS.md` ("o código só começa a usar o escuro depois da aprovação")
- **Complexity**: S
- **Dependencies**: T0
- **PRE-CONDITIONS**:
  - §3.6 do `DESIGN.md` lido
- **INPUT**: O §3.6.2 cobre os papéis principais. Mas o código usa coisas que o §3.6 **não** define:

  | Caso no código                                                            | Onde                                                                                                                                                            | Pergunta para o design                              |
  | ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
  | Superfície invertida: `bg-pg-ardosia text-pg-branco`                      | Tooltip (`RangeTabs`), aba ativa (`SegmentedTabs`, `RangeTabs`, chips do `Dashboard`), `UnsavedBar` mobile, `ConnectionOverlay`, bloco `<pre>` do `MonitorForm` | No escuro, ardósia some no fundo ardósia. Qual cor? |
  | Borda de card DOWN `#EDC4C8`                                              | `OpenIncidentCard`, `DownAlert`, `MonitorCard`                                                                                                                  | Já é cor fora do token. Vira token? Qual no escuro? |
  | Texto de erro `#F3B5BB` sobre ardósia                                     | `UnsavedBar`                                                                                                                                                    | Idem                                                |
  | Véu de modal `rgba(30,43,51,0.5)`                                         | `Modal`, `ConnectionOverlay`                                                                                                                                    | Sobre fundo ardósia quase não aparece               |
  | Sombras `rgba(30,43,51,…)`                                                | `RowMenu`, `ResponseChart`, `RangeTabs`                                                                                                                         | §3.6.5 diz "sem sombras". Remover no escuro?        |
  | Hover de linha/botão `hover:bg-pg-nuvem`                                  | vários                                                                                                                                                          | §3.6.5 indica `--pg-ardosia-linha`. Confirmar       |
  | Selo MAINTENANCE (`ciano-nevoa` + `ciano-fundo`) e selo neutro (`branco`) | `Badge`                                                                                                                                                         | §3.6.3 só cobre UP, DOWN, UNKNOWN                   |
  | Faixas do gráfico (`down-bg`, `unknown-bg`) e linha ciano                 | `ResponseChart`                                                                                                                                                 | §10 diz "estados de gráfico ainda não definidos"    |
  | Rótulo da opção                                                           | `ThemeCard`                                                                                                                                                     | `IDEAS.md` diz "Automático"; o pedido diz "Sistema" |

- **OUTPUT**:
  1. Uma resposta (aprovada) para cada linha da tabela.
  2. `DESIGN.md` §3.6.2 com os papéis novos (ex.: "Superfície invertida", "Borda de alerta", "Véu de modal", "Hover").
  3. Lista final de tokens de papel e seus nomes.
- **VERIFY**:
  - Toda linha da tabela tem resposta registrada no `DESIGN.md`
  - Contraste das cores novas anotado no §3.6.4

### Task T2: Camada de tokens de papel

- **Agent**: frontend
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T1
- **PRE-CONDITIONS**:
  - Lista de papéis aprovada no T1
- **INPUT**: `web/src/index.css` tem só os tokens do claro (§3.1). `tailwind.config.js` mapeia `pg.*` para `var(--pg-*)`.
- **OUTPUT**:
  1. `index.css`: adicionar os tokens constantes do §3.6.1 (`--pg-ardosia-clara`, `--pg-ciano-luz`, `--pg-up-claro`...).
  2. `index.css`: tokens de papel em `:root, [data-theme="light"]` e em `[data-theme="dark"]` (ver "Decisão de arquitetura"). Inclui status: `--pg-cor-up`, `--pg-cor-up-bg`, `--pg-cor-down`... e `color-scheme: light` / `dark`.
  3. `tailwind.config.js`: novas cores de papel (`bg-pg-superficie`, `text-pg-texto`...) e `darkMode: ["selector", '[data-theme="dark"]']`.
  4. **Nenhum componente muda ainda.** O escuro ainda não é ligado por ninguém.
- **VERIFY**:
  - `npm run build` passa
  - Em DevTools, `document.documentElement.dataset.theme = "dark"` muda as variáveis de papel
  - Tela clara sem nenhuma diferença (nada usa os papéis ainda)

### Task T3: Núcleo do tema (lib, hook, script inline)

- **Agent**: frontend
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - Tokens de papel existem
- **INPUT**: Não há lógica de tema. `index.html` não tem script antes do bundle. Não há CSP no servidor Go (script inline é permitido; se um CSP entrar no futuro, precisa de hash/nonce para este script).
- **OUTPUT**:
  1. `src/lib/theme.js` (funções puras):
     - `getPreference()` → `"light" | "dark" | "system"`; valor inválido ou ausente = `"system"`.
     - `setPreference(p)` → salva em `localStorage` (`pingou_theme`).
     - `resolve(p, systemIsDark)` → `"light" | "dark"`.
     - `apply(resolved)` → `data-theme` no `<html>`.
     - Todo acesso a `localStorage` em `try/catch` (aba anônima, armazenamento bloqueado → segue "system").
  2. `src/lib/theme.test.js`: valores inválidos, ausência de `localStorage`, resolução dos 3 casos.
  3. `index.html`: script inline **no `<head>`, antes do CSS**, com a mesma regra mínima (lê `pingou_theme`, olha `matchMedia`, põe `data-theme`). Mais `<meta name="color-scheme" content="light dark">`.
  4. `src/hooks/useTheme.js`: devolve `{ preference, setPreference }`; em "system", escuta `matchMedia(...).addEventListener("change")` e reaplica.
  5. Escutar o evento `storage` para outra aba do mesmo navegador acompanhar a troca.
- **VERIFY**:
  - `npm test` passa com os testes novos
  - Com `pingou_theme=dark` e rede lenta (DevTools → Slow 3G), a primeira pintura já é escura
  - Atenção: **não publicar T3 sozinho** (ver Risk Areas, risco 1)

### Task T4: Migrar base (`ui/*`, layout, estilos globais)

- **Agent**: frontend
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - Capturas de referência do tema claro tiradas **antes** da migração (ver T10)
- **INPUT**: `src/theme/globalStyles.js` usa `var(--pg-nuvem)`, `var(--pg-branco)`... direto. Componentes `ui/` (`Button`, `Badge`, `Modal`, `Combobox`, `ChipGroup`, `SegmentedTabs`, `ConnectionOverlay`, `Spinner`) e `layout/` (`Shell`, `MobileHeader`, `BottomNav`) usam tokens de cor.
- **OUTPUT**:
  1. Trocar token de cor por token de papel, **olhando o contexto de cada uso** (não é `sed` cego: `text-pg-branco` em botão ≠ `text-pg-branco` em tooltip).
  2. `Button` primário: no escuro fundo `--pg-ciano-luz` e texto `--pg-ardosia` (§3.6.5).
  3. Input: borda `--pg-neblina` → `--pg-neblina-clara`.
  4. `Badge`: status pelo §3.6.3 (ícone na cor forte, texto em `--pg-nuvem`).
  5. Scrollbar e inputs do `globalStyles.js`.
- **VERIFY**:
  - Captura do tema claro igual à referência
  - Componentes base legíveis com `data-theme="dark"` forçado

### Task T5: Migrar telas, modais e cores soltas

- **Agent**: frontend
- **Priority**: P0
- **Complexity**: L (o maior volume: ~900 classes `pg-*` em ~60 arquivos)
- **Dependencies**: T4
- **PRE-CONDITIONS**:
  - Base migrada
- **INPUT**: Páginas (`Dashboard`, `Incidents`, `MonitorDetail`, `Maintenance`, `Settings`, `Login`) e pastas `incidents/`, `maintenance/`, `monitor-detail/`, `monitors/`, `settings/`. Cores fora do token: `#EDC4C8` (3×), `#F3B5BB` (1×), `rgba(30,43,51,…)` (5×).
- **OUTPUT**:
  1. Migrar por pasta, um commit por pasta (facilita revisão e rollback).
  2. Trocar os 4 hexadecimais e os 5 `rgba` pelos tokens aprovados no T1.
  3. `CertCard`: equivalências do `DESIGN.md` §4.1 (parágrafo do tema escuro).
  4. URLs/IDs em mono: `--pg-ciano-fundo` → `--pg-ciano-luz`.
- **VERIFY**:
  - `grep -rnE "#[0-9a-fA-F]{3,8}\b|rgba?\(" web/src --include=*.jsx` sem resultados
  - Captura do tema claro igual à referência, tela por tela

### Task T6: Gráficos SVG e barras de histórico

- **Agent**: frontend
- **Priority**: P1
- **Complexity**: S
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - Cores de gráfico aprovadas no T1
- **INPUT**: `ResponseChart.jsx` usa `var(--pg-borda)`, `var(--pg-unknown-bg)`, `var(--pg-down-bg)`, `var(--pg-down)`, `var(--pg-ciano)` em `fill`/`stroke`. `HistoryBar.jsx` usa `bg-pg-up`/`bg-pg-down`/`bg-pg-borda`.
- **OUTPUT**:
  1. `fill`/`stroke` passam a usar variáveis de papel (SVG inline herda a variável do `<html>`, troca sem rerender).
  2. Barras de histórico com `--pg-up-claro`/`--pg-down-claro` no escuro.
  3. Tooltip do gráfico como superfície (sem sombra no escuro, se T1 confirmar).
- **VERIFY**:
  - Gráfico com faixa "fora do ar" legível nos dois temas
  - Linha ciano sobre `--pg-ardosia-clara`: 3,8:1, acima do mínimo 3:1 para elemento gráfico

### Task T7: Logos e toasts no escuro

- **Agent**: frontend
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T2
- **PRE-CONDITIONS**:
  - Nenhuma
- **INPUT**: `Shell` usa `pingou-marca-horizontal.svg`; `MobileHeader` e `Dashboard` usam `pingou-simbolo.svg`. `pingou-simbolo-invertido.svg` existe em `/assets`, mas **não** em `web/src/assets/brand/`. `ToastContainer` (react-toastify) em `main.jsx` está no tema claro.
- **OUTPUT**:
  1. Copiar `pingou-simbolo-invertido.svg` para `web/src/assets/brand/`.
  2. Renderizar as duas versões e esconder uma por CSS (`dark:hidden` / `hidden dark:block`). Assim não depende do React saber o tema e não pisca.
  3. Toasts: usar as variáveis CSS do react-toastify (`--toastify-color-light`...) apontando para tokens de papel, em vez de trocar a prop `theme`.
- **VERIFY**:
  - Marca negativa e símbolo invertido aparecem no escuro, desde a primeira pintura
  - Toast de sucesso e de erro legíveis nos dois temas

### Task T8: Bloco "Aparência" em Configurações

- **Agent**: frontend
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: T3
- **PRE-CONDITIONS**:
  - `useTheme` pronto
  - Rótulo final da 3ª opção decidido no T1
- **INPUT**: `Settings.jsx` tem um formulário com estado `dirty` e botão "Salvar" (valores globais, vão para a API). O tema é por navegador e não vai para a API.
- **OUTPUT**:
  1. `components/settings/ThemeCard.jsx` usando o `Card` existente:
     - Título: "Aparência". Descrição: "Vale só neste navegador." (deixa clara a diferença para "Valem para todos os monitores.").
     - Três opções exclusivas: Claro, Escuro, Sistema. Usar grupo de rádio acessível (`role="radiogroup"`, setas do teclado) ou o `SegmentedTabs` existente (`aria-pressed`).
  2. Aplica **na hora** (chama `setPreference`). **Fora** do `dirty`, do "Salvar" e da `UnsavedBar`.
  3. Posição: primeiro ou último card da tela (definir com o design; sugestão: último, por ser pessoal e não global).
  4. **Não** adicionar nada em `Shell`, `MobileHeader` ou `BottomNav`.
- **VERIFY**:
  - Trocar a opção não mostra a barra "alterações não salvas"
  - Recarregar mantém a escolha
  - Navegação por teclado funciona; leitor de tela anuncia a opção marcada

### Task T9: Trava contra regressão de tokens

- **Agent**: frontend
- **Priority**: P1
- **Complexity**: XS
- **Dependencies**: T5
- **PRE-CONDITIONS**:
  - Migração completa
- **INPUT**: Nada impede alguém de voltar a escrever `text-pg-ardosia` ou um hexadecimal num componente novo, e o escuro fica "furado" de novo.
- **OUTPUT**:
  1. `scripts/check-tokens.sh`: falha se `web/src/components` ou `web/src/pages` tiverem hexadecimal, `rgba(` ou classe de token de cor crua (lista definida no T2).
  2. Script npm `check:tokens` e chamada no `Makefile` (alvo `build-web` ou `test`).
- **VERIFY**:
  - Inserir `text-pg-ardosia` num componente faz o check falhar; remover faz passar

### Task T10: Verificação visual, contraste e piscar

- **Agent**: qa
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: T5–T9
- **PRE-CONDITIONS**:
  - Capturas de referência do claro tiradas antes do T4
- **INPUT**: Telas: Login, Painel (com e sem DOWN, sem conexão), Incidentes, Detalhe (com SSL em cada nível), Manutenções, Configurações, todos os modais e menus. Larguras: desktop e celular.
- **OUTPUT**:
  1. Script Playwright que tira capturas com `colorScheme: "light"` e `"dark"` (opção Sistema) e com `pingou_theme` forçado.
  2. Comparação claro antes × depois (diferença zero esperada).
  3. Revisão do escuro com o `DESIGN.md` §3.6 e o checklist §11.
  4. Contraste: rodar axe-core (via Playwright) nas telas escuras; zero erros de contraste.
  5. Teste de piscar: gravar o carregamento com rede lenta e conferir que o primeiro quadro já é escuro.
  6. Teste de daltonismo: emular deuteranopia/protanopia no DevTools e conferir UP × DOWN (texto + ícone).
- **VERIFY**:
  - Capturas anexadas ao PR
  - Relatório do axe sem violações de contraste

### Task T11: Documentação final

- **Agent**: docs
- **Priority**: P2
- **Complexity**: XS
- **Dependencies**: T10
- **PRE-CONDITIONS**:
  - Feature verificada
- **INPUT**: `DESIGN.md` §10 lista o mecanismo do tema como pendente. `IDEAS.md` marca a ideia como "Para detalhar".
- **OUTPUT**:
  1. `DESIGN.md` §10: remover a pendência do tema escuro; §11: incluir "usei token de papel, não de cor?".
  2. `IDEAS.md`: mover a ideia para feita / tirar do topo.
  3. `docs/implementacoes.md`: registrar o mecanismo (`data-theme`, `pingou_theme`, tokens de papel) para a página pública reaproveitar.
  4. Rodar `graphify update .`.
- **VERIFY**:
  - Docs revisados no PR

## Risk Areas

| Risk                                                                                                                                             | Impact                                    | Mitigation                                                                                                                                   |
| ------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| 1. **Padrão é "Sistema"**: se o T3 for publicado antes da migração, todo usuário com sistema escuro vê um painel quebrado sem ter escolhido nada | Alto                                      | Publicar T3 a T8 **juntos**, numa única versão. Enquanto não terminar, o script inline fica sem efeito (ou atrás de uma constante desligada) |
| 2. Migração por `sed` cego troca papel errado (`text-pg-branco` de botão × de tooltip)                                                           | Médio: escuro ilegível em pontos isolados | Migrar por contexto, um commit por pasta, revisão com capturas                                                                               |
| 3. Tema claro muda sem querer                                                                                                                    | Alto (fora do escopo, `IDEAS.md` proíbe)  | Capturas de referência antes do T4 e comparação após cada pasta                                                                              |
| 4. Lacunas de design (T1) sem resposta travam o código                                                                                           | Médio: atraso                             | T0, T1 primeiro; T2 e T3 podem começar com papéis já definidos no §3.6.2 enquanto o resto é aprovado                                         |
| 5. Cores fora do token já existentes (`#EDC4C8`, `#F3B5BB`) ganham "equivalente escuro" inventado                                                | Médio: viola §3.4                         | Só valores aprovados no T1                                                                                                                   |
| 6. Script inline e `lib/theme.js` divergem com o tempo                                                                                           | Baixo                                     | Comentário cruzado nos dois; teste em `theme.test.js` cobre a mesma regra                                                                    |
| 7. `localStorage` indisponível (aba anônima, bloqueio)                                                                                           | Baixo                                     | `try/catch` e cair em "system"                                                                                                               |
| 8. CSP futura bloqueia o script inline                                                                                                           | Baixo (hoje não há CSP)                   | Registrar no `docs/implementacoes.md`; usar hash do script se uma CSP entrar                                                                 |

## Rollback Strategy

### Fase 1 — tokens e migração (T2, T4–T7)

1. No tema claro os papéis resolvem para os mesmos valores de hoje. Reverter um commit de pasta não afeta as outras.
2. Se algo der errado no claro, `git revert` do commit da pasta.

### Fase 2 — ligar o tema (T3, T8)

1. Desligar o escuro sem reverter a migração: remover o script inline do `index.html` e o `ThemeCard`. Sem `data-theme="dark"`, tudo fica claro.
2. O valor salvo em `localStorage` fica parado e não causa efeito.

## Phase X: Verification

### Build Verification

```bash
cd web && npm ci && npm run build
make build
```

### Lint & Type Check

```bash
cd web && npm run lint && npm run format:check && npm test && npm run check:tokens
```

(O projeto é JavaScript, sem type-check.)

## Cronograma Sugerido

| Fase   | Tarefas        | Complexity | Esforço   | Target                                 | Status      |
| ------ | -------------- | ---------- | --------- | -------------------------------------- | ----------- |
| Fase 0 | T0, T1         | S          | 0,5–1 dia | Design fechado e aprovado              | ✅ Feito    |
| Fase 1 | T2, T3         | S          | 1 dia     | Infra de tema pronta (ainda desligada) | ⬜ Pendente |
| Fase 2 | T4, T5, T6, T7 | L          | 3–4 dias  | Todas as telas em tokens de papel      | ⬜ Pendente |
| Fase 3 | T8, T9         | S          | 0,5–1 dia | Seletor em Configurações + trava       | ⬜ Pendente |
| Fase 4 | T10, T11       | M          | 1–1,5 dia | Verificado, documentado, versão única  | ⬜ Pendente |

Total estimado: **6 a 8,5 dias** de trabalho.

## Status de Execução

### Concluído

- T2 (2026-11-05): tokens de papel criados em `web/src/index.css` (`--pg-cor-*`, claro e escuro) e mapeados em `tailwind.config.js` (`bg-pg-fundo`, `text-pg-texto`, `border-pg-linha`, `shadow-pg-flutuante`...). `darkMode` por `[data-theme="dark"]`. Conferidos 31 papéis × 3 estados (sem `data-theme`, claro, escuro) no navegador. Nenhum componente alterado. Build, lint, format e testes passam.
- T1 (2026-11-05): lacunas de design aprovadas pelo Bryan e registradas no `DESIGN.md` §3.6.1, 3.6.4 e 3.6.6. Nome da 3ª opção: **Sistema**. Cor nova aprovada: `--pg-down-borda-escuro` (#684953); `--pg-down-borda` (#EDC4C8) vira token no claro.
- T0 (2026-11-05): removida a v1.0 duplicada do `DESIGN.md`. Nada da v1.0 faltava na v1.1.

### Em andamento

-

### Pendente

- T3 a T11

## Notes

- Levantamento feito em 2026-11-05 no `web/src`: ~900 usos de classes `pg-*`, 4 hexadecimais soltos, 5 `rgba(...)`, 6 `var(--pg-*)` no `ResponseChart` e 9 no `globalStyles.js`. Nenhuma classe de cor nativa do Tailwind (`gray-*`, `red-*`) e nenhum token antigo (`text-text`, `bg-bg`) em uso.
- Os tokens antigos do `tailwind.config.js` (`text`, `bg`, `accent`...) apontam para variáveis que nem existem no `index.css`. Podem sair junto, mas **fora** deste plano (limpeza separada).
- Página pública de status (ideia 3): reaproveita `lib/theme.js`, os tokens de papel e o script inline. Lá terá botão no cabeçalho, com chave de `localStorage` própria se a decisão for não misturar com o painel.
- Não há mockups de tema escuro em `concept-design/`. A revisão visual do T10 usa o `DESIGN.md` §3.6 como referência.

### Alternativas consideradas e descartadas

- `dark:` do Tailwind em cada classe → descartada: dobra ~900 classes, fácil esquecer uma, e espalha a regra de cor por todo o código.
- Redefinir os tokens de cor (`--pg-ardosia`, `--pg-branco`) no escuro → descartada: o mesmo token tem papéis diferentes (texto × fundo invertido), e o §3.6.1 diz que eles são constantes.
- Guardar o tema em `settings` no backend → descartada: o `IDEAS.md` decidiu "por navegador", e cada pessoa escolhe o seu.
- Aplicar o tema só no React (`useEffect`) → descartada: pisca o claro antes de virar escuro.
- Atalho no cabeçalho do painel → descartado por decisão de produto (2026-11-04).
