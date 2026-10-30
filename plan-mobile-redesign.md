# Plano de Desenvolvimento: Remodelagem mobile (Painel, Incidentes, Configurações e Detalhe do monitor)

## Overview

Remodelar a versão **mobile** (< 768 px, breakpoint `md` do Tailwind) das quatro páginas para chegar nos mockups:

- [mobile-Detalhe do monitor · aberto pelo alerta@1x.png](concept-design/mobile-Detalhe%20do%20monitor%20%C2%B7%20aberto%20pelo%20alerta@1x.png)
- [mobile-Painel · 2 fora do ar (filtros funcionam)@1x.png](<concept-design/mobile-Painel%20%C2%B7%202%20fora%20do%20ar%20(filtros%20funcionam)@1x.png>)
- [mobile-Incidentes@1x.png](concept-design/mobile-Incidentes@1x.png)
- [mobile-Configurações · retenção alterada@1x.png](concept-design/mobile-Configura%C3%A7%C3%B5es%20%C2%B7%20reten%C3%A7%C3%A3o%20alterada@1x.png)

O desktop (≥ 768 px) **não muda**. Hoje o mobile é o desktop espremido: header com logo + hambúrguer, drawer lateral, tabela do painel com colunas escondidas, filtros com contagem em texto.

### O que os mockups mudam (estrutura compartilhada)

| Elemento  | Hoje (mobile)                                                                                     | Mockup                                                                                                                             |
| --------- | ------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Navegação | Header com logo + hambúrguer, drawer de 260 px ([Shell.jsx](web/src/components/layout/Shell.jsx)) | **Barra inferior fixa** com 3 abas (Painel, Incidentes, **Ajustes**) e badge vermelho de incidentes abertos                        |
| Cabeçalho | Logo, sem título da página                                                                        | **Barra superior com símbolo + título da página** (+ ações da página no Painel). O `<h1>` desktop some no mobile para não duplicar |
| Sair      | Botão no drawer                                                                                   | Botão **"Sair"** ao fim de Configurações + rodapé `pingou v1.0 · GitHub · Relatar problema`                                        |
| Fundo     | Branco                                                                                            | Fundo `nuvem` (cinza claro) com cards brancos                                                                                      |

### Por página

**Painel**

- Header: símbolo, "Monitores", "Atualizado há 12 s", botão atualizar (ícone) e **botão "+" (ícone, primário)** no lugar de "Novo monitor".
- Banner compacto: título + **uma linha de contagens** (`3 no ar · 2 fora · 1 desconhecido · 1 pausado`), sem a frase longa.
- Filtros: **linha de chips rolável na horizontal** (sem quebrar), em vez de `flex-wrap`.
- Lista: **cards** no lugar da tabela. Card = badge + nome + chevron → detalhe; alvo em mono; caixa com o erro (só DOWN); `HistoryBar` em largura total; rodapé de uma linha por estado:
  - DOWN: `Fora há 1 dia e 4 h · 90,0% em 24 h`
  - UP: `99,2% em 24 h · 214 ms · há 41 s`
  - UNKNOWN: `Aguardando o primeiro check` (barras cinza)
  - Heartbeat: `Heartbeat · espera um ping por dia`
- Card DOWN com borda avermelhada.
- Sem menu `⋯` no card: editar/pausar/excluir ficam no detalhe, que já tem os botões ([MonitorDetail.jsx](web/src/pages/MonitorDetail.jsx) linhas 338–345).

**Incidentes**

- Header com título; sem o parágrafo "Cada vez que um monitor caiu…" no mobile.
- Abas **segmentadas** (controle único de 3 partes, "Abertos 2") no lugar dos chips; o `<select>` de monitor vira **campo com rótulo acima**, largura total.
- "Abertos agora" (badge 2): cards **empilhados**, com botões lado a lado "Ver monitor" + "Abrir URL" no rodapé (hoje só existe "Abrir URL"; "Ver monitor" é novo no mobile).
- "Resolvidos": **lista em um card** agrupada por dia (`Hoje · 30 out`), cada item = ícone ✓, nome, duração à direita, `10:12 → 10:13` e o erro em mono; botão "Carregar mais" centralizado. Substitui a tabela.

**Configurações**

- Header com título; sem o `<h1>` e o subtítulo.
- Cards com padding menor (`p-4`). "Enviar teste" em **largura total** abaixo do input; resultado em caixa verde (`✓ Teste enviado. Resposta 200 OK em 142 ms.`).
- Chips de retenção **numéricos em linha única** (`7 14 30 60 90`) com o rótulo "dias" embaixo, não "7 dias" em cada chip.
- **Aviso de impacto** ao mudar a retenção: "Ao salvar, checks com mais de **14 dias** serão apagados. Isso não pode ser desfeito." (só quando o valor novo é **menor** que o salvo).
- Backup: botão largura total `Baixar backup · 12,4 MB`.
- **Barra "Alterações não salvas" + Descartar / Salvar** fixa acima da barra de abas, só quando `dirty`. Substitui o botão "Salvar alterações" no mobile.

**Detalhe do monitor** ([MonitorDetail.jsx](web/src/pages/MonitorDetail.jsx))

- **Barra superior própria** (não a do `MobileHeader` padrão): `‹ Monitores` à esquerda (volta ao Painel) e menu `⋮` à direita (`RowMenu`, com editar/pausar/excluir). Substitui o link "Monitores" solto e o botão de menu do cabeçalho. A aba "Painel" continua ativa na barra inferior.
- Título grande + `Badge`; URL com ícone de globo (link em nova aba); linha `HTTP · a cada 30 s · último check há 11 s` (o tipo entra na linha, por isso some da "Configuração").
- **Cartão de alerta** quando DOWN: "Fora do ar há 1 dia e 4 h", `Desde ontem, 11:43 · 3.365 falhas seguidas` e o erro numa caixa branca. Vem do incidente aberto (já carregado na lista de incidentes do detalhe); sem incidente aberto, não aparece.
- **Pausar** e **Editar** lado a lado, largura total, abaixo do alerta.
- "Desempenho": abas de período **segmentadas** em largura total (24 h, 7 d, 14 d, 30 d; o 30 d fica riscado e não selecionável por estar acima da retenção de 14 dias).
- 4 cards de métrica em grade 2×2, rótulos curtos (`Uptime 24 h`, `Resposta média`, `Incidentes`, `Tempo fora`), **sem a nota** de rodapé do desktop; valor em vermelho para uptime baixo e para incidentes > 0.
- Gráfico "Tempo de resposta" com a legenda numa linha (`média`, `fora do ar`, `sem dados`) e 3 rótulos de tempo no eixo.
- "Últimos checks": lista (hora, `× Falha`/`✓ OK`, latência à direita, erro em mono truncado embaixo) + "Ver mais checks" centralizado.
- "Incidentes": linhas simples, só `quando` à esquerda e `duração`/`em aberto` à direita, com "Ver todos" no título.
- "Configuração": lista chave/valor **sem a linha "Tipo"** e sem o link "Editar" (o botão Editar já está no topo).

### Divergências entre mockup e código (decisões já tomadas)

| Mockup                                            | Código                                                                                                            | Decisão                                                                                                                                         |
| ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Aba "Ajustes"                                     | Item chama "Configurações"                                                                                        | Rótulo curto **só na barra inferior** ("Ajustes"); a página segue "Configurações"                                                               |
| Aviso "checks com mais de 14 dias serão apagados" | Retenção apaga checks **e incidentes resolvidos** ([retention.go](internal/scheduler/retention.go))               | Cópia do aviso cita os dois: "checks e incidentes resolvidos com mais de **14 dias**…". Corrige o mockup, que diz só "checks"                   |
| Texto "Incidentes ficam guardados" na retenção    | Falso: resolvidos também são apagados (já corrigido no desktop em [plan-settings-page.md](plan-settings-page.md)) | Manter a cópia atual do desktop também no mobile                                                                                                |
| Card DOWN mostra o erro                           | `last_error` vem no incidente, **não está confirmado** no payload de `GET /api/monitors`                          | T1 verifica. Se faltar, adiciona o campo (só adicionar, não quebrar)                                                                            |
| "Fora há 1 dia e 4 h" no card                     | O payload tem estado/`last_checked_at`; início da queda vem do incidente aberto                                   | T1 resolve junto com `last_error` (mesmo ponto de origem)                                                                                       |
| Heartbeat "espera um ping por dia"                | `interval_seconds` existe                                                                                         | Derivar o texto de `interval_seconds` com `formatDuration`; `86400` → "por dia"                                                                 |
| Mockup de Configurações sem sidebar/versão        | Rodapé com versão, commit e licença (desktop)                                                                     | Mobile mostra a linha curta do mockup; desktop mantém o rodapé completo                                                                         |
| Detalhe: "Uptime 24 h" e "Tempo fora"             | Rótulos atuais: "Uptime" + nota `24 h`, "Tempo fora do ar"                                                        | Rótulo composto e curto só no mobile (classes `md:hidden`/`hidden md:block` no `StatCard`)                                                      |
| Detalhe: incidente aberto mostra "em aberto"      | `RecentIncidents` usa "em andamento"                                                                              | **Única exceção ao "desktop idêntico"**: a cópia passa a ser "em aberto" nos dois (mesmo texto da página de Incidentes). Mudança de uma palavra |
| Detalhe: alerta com "3.365 falhas seguidas"       | `failure_count` existe no incidente, mas `RecentIncidents` não o mostra                                           | Usa o incidente aberto da lista já carregada; sem chamada nova                                                                                  |
| Detalhe: menu `⋮` no topo                         | `RowMenu` já existe e é reaproveitado                                                                             | Sem componente de menu novo                                                                                                                     |
| Gráfico do mockup com eixo Y cortado ("00 ms")    | O gráfico real usa `w-12` para o eixo                                                                             | É corte do mockup. Conferir que o rótulo cabe (`1.000 ms`) na T12                                                                               |
| Heartbeat no detalhe                              | Sem gráfico e com 3 cards (sem "Resposta média")                                                                  | Mantido; mockup só cobre HTTP. No mobile heartbeat vira grade 2×2 com 3 cards (o último ocupa a linha inteira)                                  |

**Limites (fora do escopo):** mudanças no desktop (salvo a palavra "em aberto" acima); tema escuro; pull-to-refresh; swipe actions nos cards; busca (continua só com ≥ 10 monitores, como hoje); login mobile (tem plano próprio); novas métricas ou novo gráfico; qualquer mudança de comportamento da retenção.

**Restrições:**

- Sem dependência nova; Tailwind + `lucide-react` já instalados.
- Reaproveitar `Badge`, `HistoryBar`, `Button`, `ChipGroup`, `lib/format.js`, `lib/monitors.js`. Dois layouts (tabela/cards) renderizados da **mesma** lista; alternados por classes `md:` (CSS), sem `matchMedia` nem estado de viewport.
- Acessibilidade mantida: alvos ≥ 44 px, `aria-current` na aba ativa, `aria-pressed` nos filtros, foco visível, `sr-only` nos ícones.
- Respeitar `safe-area-inset-bottom` na barra inferior e na barra "Alterações não salvas".
- Branch `feat/mobile-redesign` a partir de `develop`; um commit por fase.

## Project Type

**FRONTEND** (React 19 + Vite + Tailwind em `web/`), com uma possível checagem/ajuste mínimo de API em Go (T1).

## Definition of Done (DoD)

Este plano está concluído quando:

- [ ] Todos os Success Criteria estão marcados
- [ ] Build, lint e type-check passam
- [ ] Nenhum item pendente no Status de Execução
- [ ] Evidências de verificação registradas

> Revisar o plano se: escopo mudar, nova dependência for descoberta, ou mais de 2 tasks forem bloqueadas.

## Success Criteria

- [ ] Em 390 px de largura, as 4 páginas ficam visualmente equivalentes aos mockups (comparação lado a lado por screenshot)
- [ ] Detalhe: barra `‹ Monitores` + `⋮`, cartão de alerta (só quando DOWN), Pausar/Editar lado a lado, abas de período segmentadas, métricas 2×2, checks em lista, incidentes e configuração enxutos
- [ ] Detalhe: monitor UP, PAUSED, UNKNOWN e heartbeat renderizam sem quebrar (sem alerta, sem gráfico no heartbeat)
- [ ] Barra inferior com 3 abas, aba ativa destacada e badge de incidentes abertos; drawer e hambúrguer não existem no mobile
- [ ] Painel: banner compacto, chips roláveis, cards por monitor nos 4 estados (DOWN, UP, UNKNOWN, heartbeat), filtros continuam funcionando
- [ ] Incidentes: abas segmentadas, filtro de monitor com rótulo, abertos empilhados com "Ver monitor" + "Abrir URL", resolvidos agrupados por dia com "Carregar mais"
- [ ] Configurações: teste do webhook com resultado em caixa, chips numéricos, aviso de impacto ao reduzir retenção, barra "Alterações não salvas" com Descartar/Salvar, botão Sair
- [ ] Desktop (≥ 768 px) idêntico ao de antes (comparar com `Painel.png`, `Incidentes.png`, `Configuracoes.png`)
- [ ] Sem scroll horizontal da página em 360, 390 e 430 px
- [ ] Nenhuma ação se perdeu no mobile: criar, editar, pausar/retomar, excluir (via detalhe), sair

## Tech Stack

| Technology                  | Choice Rationale                                                                                                                 |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| Tailwind (`md:` breakpoint) | Já é o padrão do projeto; troca de layout por CSS evita estado de viewport e flash na hidratação                                 |
| `lucide-react`              | Ícones dos mockups (`LayoutDashboard`, `TriangleAlert`, `Settings`, `ChevronRight`, `Plus`, `RefreshCw`, `LogOut`) já são usados |
| `react-router` `NavLink`    | A barra inferior reaproveita a lógica de aba ativa do `Shell`                                                                    |
| Componentes existentes      | `Badge`, `HistoryBar`, `Button`, `ChipGroup`, formatadores. Só ganham props de largura                                           |

## File Structure

```
web/src/
├── components/
│   ├── layout/
│   │   ├── Shell.jsx              # edit: remove drawer; header com título; renderiza BottomNav
│   │   ├── BottomNav.jsx          # NOVO: 3 abas + badge + safe-area
│   │   └── MobileHeader.jsx       # NOVO: símbolo + título + slot de ações (md:hidden)
│   ├── monitors/
│   │   ├── MonitorCard.jsx        # NOVO: card mobile (md:hidden)
│   │   ├── MonitorRow.jsx         # edit: tabela vira hidden md:table-row
│   │   ├── HistoryBar.jsx         # edit: prop `fluid` (largura total)
│   │   └── StatusBanner.jsx       # edit: variante compacta no mobile
│   ├── incidents/
│   │   ├── OpenIncidentCard.jsx   # edit: botão "Ver monitor", ações lado a lado
│   │   ├── ResolvedTable.jsx      # edit: lista agrupada no mobile
│   │   └── SegmentedTabs.jsx      # NOVO (ou variante em ui/): abas segmentadas
│   ├── monitor-detail/
│   │   ├── DetailTopBar.jsx       # NOVO: ‹ Monitores + RowMenu (md:hidden)
│   │   ├── DownAlert.jsx          # NOVO: cartão "Fora do ar há…" (incidente aberto)
│   │   ├── RangeTabs.jsx          # edit: largura total e segmentado no mobile
│   │   ├── ChecksTable.jsx        # edit: lista no mobile
│   │   ├── RecentIncidents.jsx    # edit: linha simples no mobile, "em aberto"
│   │   └── ConfigCard.jsx         # edit: sem "Tipo" no mobile
│   └── settings/
│       ├── UnsavedBar.jsx         # NOVO: barra fixa Descartar/Salvar
│       ├── RetentionCard.jsx      # edit: chips numéricos + aviso de impacto
│       └── WebhookCard.jsx        # edit: botão largura total + caixa de resultado
├── pages/
│   ├── Dashboard.jsx              # edit: cabeçalho mobile, lista de cards
│   ├── Incidents.jsx              # edit: abas segmentadas, select com rótulo
│   ├── Settings.jsx               # edit: Sair, rodapé curto, UnsavedBar
│   └── MonitorDetail.jsx          # edit: cabeçalho, alerta, ações, métricas, ordem das seções
└── index.css                      # edit: fundo nuvem no mobile, utilitário safe-area
```

## Task Breakdown

| Task ID | Agent                   | Name                                                             | Priority | Complexity | Dependencies         |
| ------- | ----------------------- | ---------------------------------------------------------------- | -------- | ---------- | -------------------- |
| T1      | `go-backend-specialist` | Confirmar `last_error` e início da queda no payload de monitores | P0       | S          | -                    |
| T2      | `frontend-specialist`   | Shell mobile: `MobileHeader` + `BottomNav`, remover drawer       | P0       | M          | -                    |
| T3      | `frontend-specialist`   | Painel: cabeçalho, banner compacto e chips roláveis              | P1       | M          | T2                   |
| T4      | `frontend-specialist`   | Painel: `MonitorCard` nos 4 estados                              | P1       | L          | T1, T3               |
| T5      | `frontend-specialist`   | Incidentes: abas segmentadas, filtro e abertos                   | P1       | M          | T2                   |
| T6      | `frontend-specialist`   | Incidentes: resolvidos em lista agrupada                         | P1       | M          | T5                   |
| T7      | `frontend-specialist`   | Configurações: cards, webhook, chips e aviso de retenção         | P1       | M          | T2                   |
| T8      | `frontend-specialist`   | Configurações: `UnsavedBar`, Sair e rodapé                       | P1       | M          | T7                   |
| T11     | `frontend-specialist`   | Detalhe: barra superior, título, alerta e ações                  | P1       | M          | T2                   |
| T12     | `frontend-specialist`   | Detalhe: período, métricas e gráfico                             | P1       | M          | T11                  |
| T13     | `frontend-specialist`   | Detalhe: checks, incidentes e configuração                       | P1       | M          | T11                  |
| T9      | — (`design-ux`)         | Revisão de a11y e estados (loading/vazio/erro)                   | P2       | S          | T4, T6, T8, T12, T13 |
| T10     | — (`qa-review`)         | QA visual em 360/390/430 px + regressão desktop                  | P3       | M          | T9                   |

### Task T1: Confirmar `last_error` e início da queda no payload de monitores

- **Agent**: `go-backend-specialist`
- **Priority**: P0
- **Complexity**: S
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Branch `feat/mobile-redesign` criada a partir de `develop`
  - Backend sobe (`make dev`)
- **INPUT**: `GET /api/monitors` devolve `current_state`, `history`, `uptime_24h`, `last_latency_ms`, `last_checked_at`. Não está confirmado se devolve o erro nem o início da queda
- **OUTPUT**:
  1. Verificação documentada: o payload tem `last_error` e `down_since` (ou equivalente)?
  2. Se faltar: adicionar os dois campos **só quando `current_state == DOWN`**, vindos do incidente aberto (uma consulta agregada para a lista, nunca N+1), `null` nos outros estados. Formato atual preservado. Teste no handler/repositório e linha nova na tabela de campos do [README.md](README.md)
  3. Se já existir: nenhuma mudança, registrar o nome real dos campos nas Notes
- **VERIFY**:
  - `curl -s localhost:4000/api/monitors | jq '.data[] | {name, current_state, last_error, down_since}'` mostra os campos para um monitor DOWN
  - `go test ./...` passa

### Task T2: Shell mobile: `MobileHeader` + `BottomNav`, remover drawer

- **Agent**: `frontend-specialist`
- **Priority**: P0
- **Complexity**: M
- **Dependencies**: None
- **PRE-CONDITIONS**:
  - Mockups lidos; breakpoint mobile = `< md`
- **INPUT**: [Shell.jsx](web/src/components/layout/Shell.jsx) tem header com hambúrguer, `<aside>` com `translate-x` e handlers de Escape/foco só para o drawer; `useOpenIncidents` já calcula o badge
- **OUTPUT**:
  1. `BottomNav.jsx`: `<nav aria-label="Principal" class="md:hidden fixed bottom-0 …">`, 3 `NavLink` (ícone + rótulo: Painel, Incidentes, **Ajustes**), `aria-current="page"` (vem do `NavLink`), aba "Painel" ativa também em `/monitors/*` (regra que o `Shell` já tem), badge vermelho no ícone de Incidentes com `sr-only`, `padding-bottom: env(safe-area-inset-bottom)`
  2. `MobileHeader.jsx`: `md:hidden sticky top-0`, símbolo da marca + `<h1>` com o título + slot `actions`; subtítulo opcional ("Atualizado há 12 s")
  3. `Shell.jsx`: remove `open`, overlay, botão de menu, Escape e o foco do drawer; `aside` fica `hidden md:flex` (desktop idêntico, **inclui o botão Sair do desktop**); `main` ganha `pb-24 md:pb-8` para o conteúdo não ficar sob a barra; `useOpenIncidents` passa a alimentar sidebar e barra pela mesma chamada (um polling só)
  4. Cada página fornece título/ações ao header (props de `Shell` ou `MobileHeader` renderizado dentro da página, o que for menor; preferir o segundo, sem contexto novo)
  5. `index.css`: fundo `nuvem` abaixo de `md`
- **VERIFY**:
  - 390 px: sem hambúrguer, barra inferior visível, aba ativa muda ao navegar, badge mostra o número
  - ≥ 768 px: sidebar idêntica à anterior
  - Tab percorre as 3 abas; foco visível
  - `npm run lint && npm run build`

### Task T3: Painel: cabeçalho, banner compacto e chips roláveis

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**: T2 mergeada na branch
- **INPUT**: [Dashboard.jsx](web/src/pages/Dashboard.jsx) tem `<h1>`, "Atualizado", botão atualizar e "Novo monitor" em linha; [StatusBanner.jsx](web/src/components/monitors/StatusBanner.jsx) mostra título + frase + lista de contagens; filtros com `flex-wrap`
- **OUTPUT**:
  1. `MobileHeader` do Painel com "Monitores", `Atualizado {timeAgo}`, botão atualizar (44 px) e botão "+" primário (`aria-label="Novo monitor"`) que abre o mesmo modal. O cabeçalho desktop vira `hidden md:flex`. Os dois usam os mesmos handlers
  2. `StatusBanner`: no mobile, esconde a frase de detalhe e mostra as contagens numa linha (`3 no ar · 2 fora · …`, "fora" em vermelho); desktop igual. Texto vem do `stats` que já existe
  3. Filtros: `flex-nowrap overflow-x-auto` no mobile (`-mx-4 px-4`, sem barra de rolagem visível), `flex-shrink-0` nos chips; `md:flex-wrap` mantém o desktop
  4. Os chips continuam com `aria-pressed` e a contagem; contagem some do texto visível só se o mockup não a mostrar (o mockup corta o texto dos chips: manter a contagem, é informação útil)
- **VERIFY**:
  - 390 px: título, "Atualizado há N s" atualiza a cada segundo, "+" abre o modal, atualizar gira enquanto busca
  - Chips rolam na horizontal sem causar scroll da página
  - Filtro "Fora do ar" mostra só DOWN

### Task T4: Painel: `MonitorCard` nos 4 estados

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: L
- **Dependencies**: T1, T3
- **PRE-CONDITIONS**: Nomes reais dos campos de erro e início da queda definidos na T1
- **INPUT**: A lista é só tabela ([MonitorRow.jsx](web/src/components/monitors/MonitorRow.jsx)); colunas somem por breakpoint
- **OUTPUT**:
  1. `MonitorCard.jsx` (`md:hidden`): `<li>` com `Link` cobrindo o card para `/monitors/:id` (chevron à direita; `aria-label` com nome e estado); `Badge` + nome (truncado); alvo em `font-pg-mono` truncado, ou "Heartbeat · espera um ping por {intervalo}"; caixa `bg-pg-nuvem` com o erro em mono quando DOWN; `HistoryBar` com `fluid`; rodapé de estado (textos do Overview). Borda `#EDC4C8` quando DOWN, opacidade reduzida quando pausado
  2. Helper `cardFooter(monitor, now)` em [lib/monitors.js](web/src/lib/monitors.js) devolvendo a string por estado, reusando `formatUptime`, `formatDuration`, `timeAgo`. Fica testável sem React
  3. `HistoryBar`: prop `fluid` troca `w-[198px]` por `w-full`; barras continuam 30
  4. Dashboard renderiza `<ul class="md:hidden">` de cards **e** a tabela `hidden md:block` a partir do mesmo `visible`; estados vazios ("Nenhum monitor neste filtro") nos dois
  5. Ações editar/pausar/excluir: só pelo detalhe no mobile (`RowMenu` continua na tabela)
- **VERIFY**:
  - Com monitores DOWN, UP, UNKNOWN, PAUSED e heartbeat, cada card mostra o texto do mockup
  - Toque no card abre o detalhe; Enter no teclado também
  - Erro longo quebra linha sem estourar a largura
  - Teste de `cardFooter` para os 5 casos (UP com/sem latência, DOWN, UNKNOWN, heartbeat, PAUSED)
  - Desktop: tabela igual à anterior

### Task T5: Incidentes: abas segmentadas, filtro e abertos

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**: T2 mergeada
- **INPUT**: [Incidents.jsx](web/src/pages/Incidents.jsx) usa chips `aria-pressed` ligados ao parâmetro `?tab=`; `<select>` com rótulo ao lado; [OpenIncidentCard.jsx](web/src/components/incidents/OpenIncidentCard.jsx) só tem "Abrir URL"
- **OUTPUT**:
  1. `SegmentedTabs`: controle de 3 partes (`Todos`, `Abertos 2`, `Resolvidos`), mesma lógica `?tab=` e mesmos `aria-pressed`; no mobile ocupa a largura toda, no desktop mantém os chips atuais (a variante é só classe, não componente duplicado, se couber)
  2. Filtro de monitor: rótulo acima do campo, largura total no mobile (`flex-col`), lado a lado no desktop
  3. `<h1>` e parágrafo de retenção `hidden md:block`; título no `MobileHeader`
  4. `OpenIncidentCard`: no mobile, botões "Ver monitor" (`Link` para `/monitors/:id`, só se `monitor` existe) e "Abrir URL" lado a lado (`grid grid-cols-2`), abaixo da caixa de erro; `dl` de "Começou/Falhas seguidas" some do mobile (mockup mostra "Desde ontem, 11:43" sob o título, que é o `startedLabel`)
- **VERIFY**:
  - Trocar de aba atualiza a URL (`?tab=`) e recarrega o estado ao voltar
  - Monitor removido: só "Abrir URL"/nada, sem link quebrado
  - "Abrir URL" só aparece para http e abre em nova aba
  - Desktop igual

### Task T6: Incidentes: resolvidos em lista agrupada

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T5
- **PRE-CONDITIONS**: T5 concluída
- **INPUT**: [ResolvedTable.jsx](web/src/components/incidents/ResolvedTable.jsx) já faz `groupByDay` e tem "Carregar mais"; a tabela tem 4 colunas
- **OUTPUT**:
  1. No mobile, cada grupo vira cabeçalho de dia (`Hoje · 30 out`) seguido de itens: ícone ✓ verde, nome (peso médio), duração à direita, linha `início → fim`, erro em mono truncado em 1–2 linhas. Reusa `groupByDay`, `dayLabel` e os formatadores atuais
  2. Tabela `hidden md:table`, lista `md:hidden`, os dois do mesmo `groups`
  3. "Carregar mais" centralizado, largura total, mesmo `onLoadMore`/`loadingMore`
- **VERIFY**:
  - Incidente que cruza a meia-noite (`23:58 → 00:04`) mostra a duração certa e fica no grupo do `ended_at`
  - "Carregar mais" acrescenta sem duplicar (a deduplicação por `id` já existe)
  - Erro de carga mantém a lista e mostra o alerta

### Task T7: Configurações: cards, webhook, chips e aviso de retenção

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**: T2 mergeada
- **INPUT**: [Settings.jsx](web/src/pages/Settings.jsx), [WebhookCard.jsx](web/src/components/settings/WebhookCard.jsx), [RetentionCard.jsx](web/src/components/settings/RetentionCard.jsx), [BackupCard.jsx](web/src/components/settings/BackupCard.jsx)
- **OUTPUT**:
  1. `Card`: `p-4 md:p-6`
  2. `WebhookCard`: input e "Enviar teste" empilhados, botão largura total com ícone; resultado em caixa (`bg-pg-up-bg` / vermelha no erro) com `role="status"`, no formato do mockup: `Teste enviado. Resposta 200 OK em 142 ms.` (hoje: "Entregue · 200 · 142 ms"; ajustar `resultText` nos dois layouts para uma cópia só)
  3. `RetentionCard`: no mobile, chips mostram só o número e o rótulo "dias" aparece uma vez embaixo; no desktop continua "7 dias". Opção: sempre número no chip + "dias" na legenda nos dois, se o desktop aceitar (decidir olhando `Configuracoes.png`: se mostra "7 dias", manter)
  4. Aviso de impacto abaixo dos chips quando `value < settings.retention_days`: "Ao salvar, checks e incidentes resolvidos com mais de **N dias** serão apagados. Isso não pode ser desfeito." (`role="status"`). A página passa `savedValue` ao card
  5. `BackupCard`: botão largura total `Baixar backup · ≈ 12,4 MB`
- **VERIFY**:
  - Escolher 14 com 30 salvo mostra o aviso; escolher 60 não mostra; voltar ao salvo esconde
  - Teste com URL inválida e com 429 mostra a mensagem certa em caixa vermelha
  - Teclado: setas trocam a retenção (rádios `sr-only` continuam)

### Task T8: Configurações: `UnsavedBar`, Sair e rodapé

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T7
- **PRE-CONDITIONS**: T7 concluída
- **INPUT**: Salvar é um botão no fim do formulário; `form`/`dirty`/`saving`/`error` já vivem em `Settings`; `onLogout` chega no `Shell`, não na página
- **OUTPUT**:
  1. `UnsavedBar.jsx`: `fixed inset-x-0 bottom-[calc(56px+env(safe-area-inset-bottom))] md:hidden`, fundo ardósia, texto "Alterações não salvas", **Descartar** (limpa `form`) e **Salvar** (submit do `form` por `form="settings-form"`), `role="status"`; some quando `!dirty`; estado "Salvando…" desabilita os botões; mostra o erro de salvar na própria barra
  2. Botão "Salvar alterações" atual fica `hidden md:inline-flex`; o `Salvo ✓` continua nos dois
  3. Botão **Sair** no fim da página (largura total, `md:hidden`), chamando `onLogout`: passar via prop de `Shell` para a rota de Settings ou via um contexto já existente (checar [App/main.jsx](web/src/main.jsx); sem contexto novo se a prop servir)
  4. Rodapé curto no mobile: `pingou v{versão} · GitHub · Relatar problema`; o completo (commit, Apache 2.0) fica `hidden md:flex`
  5. `main` ganha `pb` extra para a barra de alterações não cobrir o botão Sair
- **VERIFY**:
  - Mudar a retenção ou a URL faz a barra aparecer; Descartar volta ao valor salvo e esconde a barra; Salvar persiste e mostra "Salvo"
  - A barra não cobre o conteúdo nem a barra de abas; respeita a safe-area
  - Sair encerra a sessão e leva ao login

### Task T11: Detalhe: barra superior, título, alerta e ações

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T2
- **PRE-CONDITIONS**: T2 mergeada (`BottomNav` e fundo `nuvem` prontos)
- **INPUT**: [MonitorDetail.jsx](web/src/pages/MonitorDetail.jsx): link `back` solto, `<header>` com `h1`, badge, alvo, "a cada … · último check" e botões Pausar/Editar `hidden sm:inline-flex` + `RowMenu`. A lista de incidentes recentes (3) já é carregada na página
- **OUTPUT**:
  1. `DetailTopBar.jsx` (`md:hidden sticky top-0`): `Link` para `/` com `ChevronLeft` + "Monitores" (alvo ≥ 44 px) e o `RowMenu` à direita. O `back` e o `RowMenu` do cabeçalho viram `hidden md:…`; o `back` continua nos estados de loading/erro, também com a barra no mobile
  2. Cabeçalho: no mobile, título em `text-pg-h1` + `Badge` na mesma linha (quebra se longo), URL com globo (link em nova aba, já existe), e a linha `{TIPO} · a cada {intervalo} · {último check}` (`HTTP`/`TCP`/`Heartbeat` com o `typeLabel` do `ConfigCard`, exportado para não duplicar). Desktop igual
  3. `DownAlert.jsx`: renderiza quando `status === "DOWN"` **e** existe incidente aberto; título `Fora do ar há {formatDuration}` (reusa a conta do `OpenIncidentCard`), linha `{startedLabel} · {failure_count} falhas seguidas` (`plural` + `toLocaleString("pt-BR")`), erro em caixa branca mono. `md:hidden`: o desktop não ganha o cartão
  4. Botões **Pausar/Retomar** e **Editar** em `grid grid-cols-2` abaixo do alerta (`md:hidden`); mesmos handlers (`handleToggle`, `openModal("edit")`)
  5. Se o incidente aberto não vier nos 3 mais recentes (caso raro), o cartão simplesmente não aparece
- **VERIFY**:
  - Monitor DOWN: alerta com tempo, início, falhas e erro; Pausar vira "Retomar" depois de pausar e o alerta some (PAUSED)
  - Monitor UP/UNKNOWN/PAUSED: sem alerta, botões presentes
  - `‹ Monitores` volta ao Painel; `⋮` abre editar/pausar/excluir; excluir leva ao Painel
  - Erro longo e nome longo não geram scroll horizontal
  - Desktop igual

### Task T12: Detalhe: período, métricas e gráfico

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T11
- **PRE-CONDITIONS**: T11 concluída
- **INPUT**: `Performance` em `MonitorDetail.jsx`, [RangeTabs.jsx](web/src/components/monitor-detail/RangeTabs.jsx) (`inline-flex`, tooltip em CSS), `StatCard`, [ResponseChart.jsx](web/src/components/monitor-detail/ResponseChart.jsx), `ChartLegend`
- **OUTPUT**:
  1. `RangeTabs`: no mobile, `flex w-full` e botões `flex-1 text-center`; o título "Desempenho" fica numa linha e as abas na de baixo (`flex-col` no mobile). Opção acima da retenção continua `aria-disabled`, riscada (`line-through`) e com tooltip; **como não há hover no toque**, o motivo vai também no `aria-describedby` e aparece ao focar/tocar (o `group-has-[:focus-visible]` atual já cobre o foco; testar o toque)
  2. `StatCard`: rótulos compostos no mobile (`Uptime {período}`, `Resposta média`, `Incidentes`, `Tempo fora`); a `note` fica `hidden md:block`. Cor vermelha (`text-pg-down`) para uptime < 100 % e para incidentes > 0. Grade `grid-cols-2` (heartbeat: 3 cards, o último com `col-span-2`)
  3. Cartão do gráfico: título + legenda empilhados no mobile; legenda com os 3 itens do mockup. Altura 200 px mantida; padding `p-4`
  4. Eixo X: no mobile, 3 rótulos (início, meio, `agora`), pela regra que já oculta os de índice ímpar (`hidden sm:inline`); conferir que não se sobrepõem em 360 px
  5. Eixo Y: garantir que o maior rótulo (`1.000 ms`) cabe em `w-12`; se não, alargar só no mobile
- **VERIFY**:
  - 24 h/7 d/14 d trocam o período e atualizam a URL (`?range=`); 30 d com retenção 14 não seleciona e explica o motivo
  - Uptime 0 % e incidentes 1 aparecem em vermelho no monitor do mockup
  - Gráfico com falhas mostra a faixa vermelha e "sem dados" cinza como no mockup
  - Heartbeat: sem gráfico, 3 cards sem buraco na grade
  - Desktop igual

### Task T13: Detalhe: checks, incidentes e configuração

- **Agent**: `frontend-specialist`
- **Priority**: P1
- **Complexity**: M
- **Dependencies**: T11
- **PRE-CONDITIONS**: T11 concluída
- **INPUT**: [ChecksTable.jsx](web/src/components/monitor-detail/ChecksTable.jsx) (tabela com a coluna Detalhe escondida em `< sm` e o erro repetido embaixo), [RecentIncidents.jsx](web/src/components/monitor-detail/RecentIncidents.jsx), [ConfigCard.jsx](web/src/components/monitor-detail/ConfigCard.jsx), grade `lg:grid-cols-3` da página
- **OUTPUT**:
  1. `ChecksTable`: no mobile, cada item é `hora` + `✓ OK`/`× Falha` (cor e ícone, não só cor) na primeira linha, latência (ou `—`) à direita, e o erro em mono truncado na segunda. A estrutura de tabela atual já mostra o detalhe embaixo em `< sm`; ajustar espaçamento e a fonte para bater com o mockup, sem trocar o markup se der. "Ver mais checks" centralizado, `onMore`/`loadingMore` atuais. Primeiros 5 checks no mobile (hoje `FIRST_CHECKS`; se for maior, usar constante por breakpoint só para a exibição, sem mudar o fetch)
  2. `RecentIncidents`: no mobile, sem o círculo de ícone nem a linha de erro (o erro aparece no alerta e na página de Incidentes); `quando` à esquerda, `em aberto` (vermelho) ou duração à direita. Copy "em andamento" → **"em aberto"** nos dois layouts; o `sr-only` continua dizendo "Em andamento"/"Resolvido"
  3. `ConfigCard`: no mobile, esconde a linha "Tipo"; remove o link "Editar" do título da seção (`hidden md:inline`)
  4. Ordem das seções no mobile (coluna única): Desempenho, Últimos checks, Incidentes, Configuração. É a ordem do DOM atual, a grade só passa a 1 coluna abaixo de `lg`; confirmar sem reordenar
  5. Espaçamento entre seções `gap-8` → ajustado para bater com o mockup (≈ 24 px)
- **VERIFY**:
  - Monitor com 3.365 falhas: lista de 5 + "Ver mais checks" carrega mais sem duplicar
  - Check bem-sucedido mostra `✓ OK` e a latência; falha mostra `× Falha` e `—`
  - Incidente aberto mostra "em aberto"; "Ver todos" leva a `/incidents?monitor=<id>`
  - Heartbeat: o "Tipo" some e a tolerância aparece no lugar do tempo limite (lógica atual do `ConfigCard`)
  - Desktop igual, salvo a palavra "em aberto"

### Task T9: Revisão de a11y e estados

- **Agent**: `design-ux`
- **Priority**: P2
- **Complexity**: S
- **Dependencies**: T4, T6, T8, T12, T13
- **PRE-CONDITIONS**: Fases 1–4 prontas
- **INPUT**: As quatro páginas no mobile
- **OUTPUT**:
  1. Checklist: alvos ≥ 44 px, contraste dos textos `neblina` sobre `nuvem`, ordem de foco (header → conteúdo → barra de alterações → abas), `aria-current`, rótulos dos botões só-ícone
  2. Estados conferidos: carregando, lista vazia, erro de carga, sem conexão (`ConnectionOverlay` não pode ficar sob a barra inferior; subir o z-index se preciso)
  3. Correções aplicadas
- **VERIFY**:
  - Navegação completa só por teclado nas três páginas
  - Leitor de tela (VoiceOver/TalkBack ou o inspetor de acessibilidade do Chrome) anuncia aba ativa, badge e estado dos cards

### Task T10: QA visual e regressão desktop

- **Agent**: `qa-review`
- **Priority**: P3
- **Complexity**: M
- **Dependencies**: T9
- **PRE-CONDITIONS**: T9 concluída
- **INPUT**: Branch completa
- **OUTPUT**:
  1. Screenshots em 360, 390 e 430 px das 4 páginas, nos estados dos mockups, comparados lado a lado com os PNGs de `concept-design/`. No detalhe: monitor DOWN (mockup), UP, PAUSED, UNKNOWN e heartbeat
  2. Screenshots do desktop (1280 px) comparados com `Painel.png`, `Incidentes.png`, `Configuracoes.png` e `Detalhes.png`
  3. Teste de orientação paisagem e de teclado virtual aberto no campo do webhook (a barra "Alterações não salvas" e a barra de abas não podem tapar o input)
  4. Lista de diferenças restantes ou "nenhuma"
- **VERIFY**:
  - Relatório anexado em Notes; zero scroll horizontal; zero regressão no desktop

## Risk Areas

| Risk                                                                 | Impact                                                                       | Mitigation                                                                                                                               |
| -------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| Payload de monitores sem `last_error`/início da queda                | Card DOWN do mockup (erro + "Fora há 1 dia e 4 h") não dá para fazer         | T1 vem antes e decide; backend só adiciona campos                                                                                        |
| Duas renderizações (tabela + cards) da mesma lista                   | Duplica nós no DOM, leitor de tela pode ler duas vezes                       | `hidden`/`md:hidden` usam `display:none`, que sai da árvore de acessibilidade. Conferir na T9                                            |
| Barra inferior e barra "Alterações não salvas" fixas                 | Cobrem conteúdo/campo com teclado virtual aberto                             | Padding inferior no `main`, safe-area, teste na T10                                                                                      |
| Remover o drawer tira o acesso ao "Sair"                             | Usuário mobile fica sem logout                                               | T8 coloca o "Sair" em Configurações; T2 só remove o drawer depois que o botão existe (se T8 atrasar, manter o Sair provisório no header) |
| Menu `⋯` some no card                                                | Editar/pausar/excluir exigem um toque a mais                                 | O detalhe já tem os botões; se for incômodo, adicionar o menu depois (fora do escopo)                                                    |
| Aviso de retenção desatualizado com o mockup                         | Mockup diz só "checks"                                                       | Cópia corrigida (checks **e** incidentes resolvidos), decidido no Overview                                                               |
| Tooltip do período desabilitado depende de hover                     | No toque o usuário não vê por que o "30 d" não abre                          | Motivo ligado ao foco/toque (T12) e testado em dispositivo real                                                                          |
| Cartão de alerta depende do incidente estar entre os 3 mais recentes | Em monitor com muitos incidentes recentes o alerta pode faltar               | O aberto é sempre o mais novo; se ainda assim faltar, o cartão não aparece e o resto da página segue completo                            |
| Mudar "em andamento" para "em aberto"                                | Única diferença visível no desktop                                           | Uma palavra, alinhada à página de Incidentes; registrar na T10                                                                           |
| Dois breakpoints diferentes de `md` no CSS atual                     | Mistura de `sm:`/`lg:` no Painel causa estados intermediários feios (tablet) | Mobile = `< md`; testar 600–800 px na T10                                                                                                |

## Rollback Strategy

1. Um commit por fase na branch `feat/mobile-redesign`: `git revert <sha>` da fase problemática
2. Sem migração e sem mudança de contrato (T1 só adiciona campos opcionais): reverter o front não exige reverter o back
3. Se a barra inferior der problema em produção, reverter só T2: o `Shell` antigo (drawer) volta sozinho, pois as demais tasks só adicionam variantes `md:hidden`
4. `git checkout develop -- web/` restaura o front inteiro

## Phase X: Verification

```bash
# Front
cd web && npm run lint && npm run build

# Back (se T1 mexer em Go)
cd .. && go vet ./... && go test ./...

# Build completo com o front embutido
make build

# Manual
make dev   # abrir em 390 px (DevTools) e comparar com concept-design/mobile-*.png
```

## Cronograma Sugerido

| Fase                  | Tarefas       | Complexity | Esforço | Target                        | Status      |
| --------------------- | ------------- | ---------- | ------- | ----------------------------- | ----------- |
| Fase 0: Dados         | T1            | S          | ~15 min | Payload confirmado            | ✅ Concluído |
| Fase 1: Shell         | T2            | M          | ~30 min | Barra inferior + header       | ✅ Concluído |
| Fase 2: Painel        | T3, T4        | M + L      | ~1 h 30 | Painel igual ao mockup        | ✅ Concluído |
| Fase 3: Incidentes    | T5, T6        | M + M      | ~1 h    | Incidentes igual ao mockup    | ✅ Concluído |
| Fase 4: Configurações | T7, T8        | M + M      | ~1 h    | Configurações igual ao mockup | ✅ Concluído |
| Fase 4b: Detalhe      | T11, T12, T13 | M + M + M  | ~1 h 30 | Detalhe igual ao mockup       | ✅ Concluído |
| Fase 5: Polimento     | T9, T10       | S + M      | ~45 min | A11y + QA + regressão         | ✅ Concluído |

> Fases 2, 3, 4 e 4b são independentes entre si (só dependem da Fase 1) e podem rodar em paralelo. Dentro da 4b, T12 e T13 rodam em paralelo depois da T11.

## Status de Execução

### Concluído

- T1: `GET /api/monitors` **não** tinha os campos. Adicionados `last_error` e `down_since` (só DOWN, do incidente aberto, uma query para a lista via `FindOpenByMonitors`), com teste e README
- T2 a T8, T11 a T13: implementados, um commit por fase
- T9/T10: conferidos 360/390/430/600 px (sem scroll horizontal, barra visível) e comparação pixel a pixel do desktop (1280 e 800 px) contra `develop`: diferenças só de dados e "em aberto"

### Em andamento

-

### Pendente

- Teste em dispositivo real (teclado virtual com a barra "Alterações não salvas", safe-area no iPhone, leitor de tela)

## Notes

- O mockup do Detalhe é de um monitor **DOWN aberto pelo alerta** (webhook leva a `/monitors/:id`). A barra `‹ Monitores` sempre volta ao Painel, não ao histórico do navegador, para o caso de o usuário ter chegado por link externo.

- Mobile = `< 768 px` (`md`). Tudo que é exclusivo de um lado usa classes `md:` (`hidden md:flex`, `md:hidden`); nenhum `matchMedia`.
- Os mockups são de 390 px de largura (@1x). Comparar nessa largura.
- O mockup do Painel mostra a barra de abas sobre o último card ("Backup diário"): o `main` precisa de `pb-24` para o fim da lista ficar alcançável.
- Decisões na execução: o aviso de retenção aparece também no desktop (protege ação irreversível); a cópia do teste de webhook mudou nos dois layouts (como previsto na T7); a faixa cinza "sem dados" do gráfico é só mobile; o `RowMenu` mantém 40 px (componente compartilhado com a tabela).
- Itens deixados de fora por não aparecerem nos mockups: busca (≥ 10 monitores), menu `⋯` no card, `Ver formato do JSON` e o "Baixar backup (.db)" com nome de arquivo. Os dois últimos continuam no desktop; no mobile, o link do JSON fica na descrição do card (é um link inline, sem custo de espaço).

### Alternativas consideradas e descartadas

- `Componente separado por página para mobile (rotas /m/...)` → descartada: duplica dados, estado e rotas; as variantes por CSS reaproveitam tudo
- `matchMedia`/hook `useIsMobile` para escolher o layout` → descartada: causa flash no primeiro render e estado extra; CSS resolve
- `Biblioteca de bottom sheet / tab bar` → descartada: uma `<nav>` fixa com 3 links é suficiente, sem dependência nova
- `Manter o drawer além da barra inferior` → descartada: duas navegações no mesmo viewport; o mockup só tem a barra
- `Manter o menu ⋯ no card` → descartada: o mockup não tem; o detalhe já cobre as ações
