# DESIGN.md — Design system do Pingou

> Fonte da verdade visual para agentes de IA e contribuidores.
> Resumo do "Manual de Identidade Visual v1.0" (Bryan H. Cruz, set/2026).
> Em caso de conflito, **este arquivo vence qualquer preferência de estilo do agente**.

## 0. Como usar este arquivo

- Leia tudo antes de criar ou alterar qualquer UI, README, doc ou asset.
- Use **somente** os tokens daqui. Não invente cores, fontes, pesos ou sombras.
- Se precisar de algo que não existe aqui, **pare e pergunte** em vez de improvisar.
- Regras marcadas com **SEMPRE** e **NUNCA** não têm exceção.

## 1. Conceito

Atributos: **sério · acolhedor · leve · confiável.**

- Formas arredondadas, cores frias e suaves, muito espaço em branco.
- Visual calmo: o ciano aparece pouco; o vermelho só quando algo quebra.
- Nada de visual "gamer", neon, degradê, glassmorphism ou sombra pesada.
- Teste rápido: se a tela parecer pesada, fria demais ou infantil, está errada.

## 2. Marca

| Termo | O que é | Pode aparecer sozinho? |
|---|---|---|
| Símbolo | Gota + 2 ondas em quadrado arredondado | Sim (favicon, avatar) |
| Logotipo | "pingou" em Rubik Medium, em curvas | **NUNCA** |
| Marca | Símbolo + logotipo | Sim (assinatura padrão) |

Regras:

- **SEMPRE** escreva "pingou" em minúsculas no logotipo.
- Em texto corrido, escreva "Pingou" (nome próprio). **NUNCA** "PingoU", "PINGOU" ou "PinGoU".
- **NUNCA** recrie o logo com `<span>` e fonte. Use os SVGs oficiais (estão em curvas).
- **NUNCA** distorça, gire, recolora, aplique sombra/brilho/degradê ou coloque sobre fundo poluído.
- Área de proteção: espaço livre de **1/4 da altura do símbolo** em todos os lados.
- Tamanho mínimo: símbolo **16 px**; marca horizontal **24 px** de altura; marca vertical **48 px**.

### Arquivos oficiais

| Arquivo | Uso |
|---|---|
| `pingou-simbolo.svg` | Padrão. Favicon, avatar, header do dashboard |
| `pingou-simbolo-livre.svg` | Símbolo sem caixa, sobre fundo claro |
| `pingou-simbolo-invertido.svg` | Caixa branca + arte ciano. **Só** sobre fundo escuro |
| `pingou-simbolo-mono.svg` | Uma cor só |
| `pingou-marca-horizontal.svg` | Assinatura padrão em fundo claro |
| `pingou-marca-horizontal-negativa.svg` | Fundo escuro (ardósia, tema escuro) |
| `pingou-marca-horizontal-mono.svg` | Uma cor só |
| `pingou-marca-vertical.svg` | Espaços quadrados/estreitos |
| `favicon-16/32/180/512.png` | Favicons e ícones de app |

> Ajuste o caminho conforme onde os assets forem colocados no repositório.

## 3. Cores

### 3.1 Tokens

```css
:root {
  /* marca */
  --pg-ciano: #1E9BB3;        /* primária: símbolo, ícones, textos >= 24px */
  --pg-ciano-fundo: #146E80;  /* botões com texto, links, hover, texto em ciano */
  --pg-ciano-nevoa: #E3F2F5;  /* fundos de destaque, cards, header de tabela */

  /* neutros */
  --pg-ardosia: #1E2B33;      /* texto principal, fundos escuros */
  --pg-neblina: #5B6B75;      /* texto secundário, legendas, UNKNOWN */
  --pg-nuvem: #F3F7F8;        /* fundo geral */
  --pg-branco: #FFFFFF;       /* superfícies (cards, inputs) */
  --pg-borda: #D9E3E6;        /* bordas e divisórias */

  /* status — uso exclusivo para estado */
  --pg-up: #1F8456;           /* verde sinal */
  --pg-up-bg: #E4F3EC;
  --pg-down: #C93A46;         /* alerta */
  --pg-down-bg: #F8E6E8;
  --pg-unknown: #5B6B75;
  --pg-unknown-bg: #EDF1F3;
}
```

### 3.2 Tailwind

Tailwind v3 (`tailwind.config.js`):

```js
theme: {
  extend: {
    colors: {
      pg: {
        ciano: '#1E9BB3', 'ciano-fundo': '#146E80', 'ciano-nevoa': '#E3F2F5',
        ardosia: '#1E2B33', neblina: '#5B6B75', nuvem: '#F3F7F8', borda: '#D9E3E6',
        up: '#1F8456', 'up-bg': '#E4F3EC',
        down: '#C93A46', 'down-bg': '#F8E6E8',
        unknown: '#5B6B75', 'unknown-bg': '#EDF1F3',
      },
    },
    fontFamily: {
      heading: ['Rubik', 'system-ui', 'sans-serif'],
      sans: ['Inter', 'system-ui', 'sans-serif'],
      mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
    },
  },
}
```

Tailwind v4: declare os mesmos valores em `@theme` (`--color-pg-ciano: #1E9BB3;` etc.).

### 3.3 Proporção

~60% nuvem/branco · ~25% ardósia · ~10% ciano · ~5% ciano-fundo/névoa. Status só onde houver estado.

### 3.4 Regras de cor

- **NUNCA** use cores quentes (vermelho, laranja, amarelo, rosa) fora de status/erro.
- **NUNCA** use verde ou vermelho em logo, grafismos, ilustrações ou peças de divulgação.
- **NUNCA** use preto puro (`#000`). O mais escuro é `--pg-ardosia`.
- **NUNCA** use `--pg-ciano` para texto menor que 24 px (contraste 3,3:1, reprova AA).
- **SEMPRE** use `--pg-ciano-fundo` como fundo de botão com texto branco (5,9:1).
- Cores fora desta lista só com aprovação. Tints novos derivam da própria cor, não de outra paleta.

### 3.5 Contraste sobre branco

| Cor | Contraste | Texto pequeno? |
|---|---|---|
| Ardósia | 14,5:1 | Sim |
| Ciano fundo | 5,9:1 | Sim |
| Neblina | 5,5:1 | Sim |
| Alerta | 5,0:1 | Sim |
| Verde sinal | 4,7:1 | Sim |
| Ciano pingo | 3,3:1 | **Não** (só ≥ 24 px ou ícones) |

## 4. Status

| Estado | Cor | Fundo | Ícone | Texto |
|---|---|---|---|---|
| UP | `--pg-up` | `--pg-up-bg` | ✓ (check) | `UP` |
| DOWN | `--pg-down` | `--pg-down-bg` | × (x) | `DOWN` |
| UNKNOWN | `--pg-unknown` | `--pg-unknown-bg` | ? | `UNKNOWN` |

- **SEMPRE** mostre status com **texto + ícone**, além da cor. **NUNCA** só a cor (daltonismo).
- Badge de status: pílula (`border-radius: 9999px`), fundo `*-bg`, texto e ícone na cor forte, Inter Medium 12–13 px.
- Bolinha de status (dot) pode acompanhar, mas nunca substitui o badge com texto.

## 5. Tipografia

| Papel | Fonte | Tamanho | Linha | Peso |
|---|---|---|---|---|
| Título 1 | Rubik | 32 px | 1.2 | 500 |
| Título 2 | Rubik | 24 px | 1.25 | 500 |
| Título 3 | Rubik | 20 px | 1.3 | 500 |
| Texto | Inter | 16 px | 1.6 | 400 |
| Interface | Inter | 14 px | 1.5 | 400 |
| Legenda | Inter | 12 px | 1.4 | 400 |
| Código/URL/ID | mono do sistema | 14 px | 1.5 | 400 |

- Só **dois pesos**: 400 e 500. **NUNCA** 300, 600, 700 ou mais.
- Títulos **SEMPRE** em Rubik 500. Textos e UI **SEMPRE** em Inter.
- **NUNCA** carregue fonte mono externa. Use a pilha `ui-monospace, ...`.
- URLs, IDs e logs: fonte mono na cor `--pg-ciano-fundo`.
- Números no dashboard (tempo de resposta, %, contagens): `font-variant-numeric: tabular-nums`.
- Nada menor que 12 px.
- Fontes: Google Fonts, licença SIL OFL 1.1. Carregar só `Rubik:wght@400;500` e `Inter:wght@400;500`.

## 6. Forma e espaço

- Cantos arredondados, ecoando o símbolo: controles `8px`, cards `12px`, badges `9999px`.
- Bordas finas: `1px solid var(--pg-borda)`.
- Sombras: evite. No máximo uma sombra muito sutil em elementos flutuantes (menus, modais).
- Espaçamento em escala de 4 px (4, 8, 12, 16, 24, 32, 48).
- Prefira espaço em branco a divisórias e caixas extras.

## 7. Componentes base

- **Botão primário**: fundo `--pg-ciano-fundo`, texto branco, Inter 500 14 px, raio 8 px. No máximo **um** por tela.
- **Botão secundário**: fundo transparente, borda `--pg-borda`, texto `--pg-ardosia`.
- **Botão destrutivo** (ex.: remover monitor): texto/borda `--pg-down`; fundo sólido só na confirmação.
- **Link**: `--pg-ciano-fundo`, sublinhado no hover.
- **Input**: fundo branco, borda `--pg-borda`, raio 8 px; foco com anel `--pg-ciano`.
- **Card**: fundo branco sobre `--pg-nuvem`, borda `--pg-borda`, raio 12 px.
- **Tabela**: header com fundo `--pg-ciano-nevoa` e texto `--pg-ciano-fundo`; linhas separadas por `--pg-borda`.
- **Header do dashboard**: símbolo (24–32 px) + logotipo, via SVG oficial.

## 8. Grafismos de apoio

| Grafismo | Função | Onde |
|---|---|---|
| Ondas | Destaque/fundo (ping se espalhando) | Banners, capa do README, slides, posts |
| Arco divisor | Sublinhar um título (uma onda só) | Títulos de seção, login |
| Barras de histórico | Separar blocos | Divisórias, rodapés |
| Gota marcador | Marcador de lista | Listas de features |

- **NUNCA** mais de 2 grafismos por peça.
- **NUNCA** ondas e arco divisor na mesma peça.
- Ondas nascem de um ponto (canto, borda ou símbolo), máx. 5 anéis, transparência crescente.
- Grafismos **só em ciano** (ou branco sobre ciano). Nunca verde/vermelho.
- Barras decorativas têm **altura igual**; variar altura parece dado real.
- **NUNCA** grafismo atrás de texto pequeno.
- Telas de trabalho (listas, formulários) ficam limpas: no máximo o arco divisor.

## 9. Tom de voz (UI)

- Português, frases curtas, voz ativa, verbo primeiro: "Criar monitor", "Salvar".
- Sem "por favor", sem "com sucesso", sem exclamação em mensagens do sistema.
- Erros: o que aconteceu + o que fazer. Ex.: "Não foi possível conectar. Verifique a URL."
- Slogan oficial: **"Rodou, pingou."**

## 10. Pendências (não definidas na v1.0)

- **Tema escuro**: não há tokens oficiais. Se for implementar, use ardósia como base, `--pg-nuvem` para texto, marca negativa e símbolo invertido, e **proponha os valores para aprovação** antes de mergear.
- Estados de gráfico/métricas e ilustrações ainda não definidos.

## 11. Checklist antes de entregar

- [ ] Só usei tokens deste arquivo?
- [ ] Nenhuma cor quente fora de status/erro?
- [ ] Status com texto + ícone, não só cor?
- [ ] Botões com texto usam `--pg-ciano-fundo`?
- [ ] Só Rubik/Inter, só pesos 400/500?
- [ ] Logo via SVG oficial, em minúsculas, com área de proteção?
- [ ] No máximo 2 grafismos, e sem ondas + arco juntos?
- [ ] Contraste de texto ≥ 4,5:1?
