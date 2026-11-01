# Histórico de implementações (planos aplicados)

Consolida o que os planos `plan-*.md` entregaram. Datas vêm do `git log` (commits e merges de PR). O comportamento atual de API e CLI está no [README.md](../README.md); aqui fica o **o quê, quando e por quê**.

| Data (commits)     | Plano                                                                                           | Branch / PR                 |
| ------------------ | ----------------------------------------------------------------------------------------------- | --------------------------- |
| 2026-10-27 → 10-28 | [Tipos de check + redesign do painel](#1-tipos-de-check-http-tcp-heartbeat--redesign-do-painel) | `feat/check-types` (#2)     |
| 2026-10-28 → 10-29 | [Página de Incidentes](#2-página-de-incidentes)                                                 | `feat/incidents-page` (#3)  |
| 2026-10-29         | [Página de Configurações](#3-página-de-configurações)                                           | `feat/settings-page` (#4)   |
| 2026-10-29 → 10-30 | [Detalhe do monitor](#4-detalhe-do-monitor)                                                     | `feat/monitor-details` (#5) |
| 2026-10-30         | [Remodelagem mobile](#5-remodelagem-mobile)                                                     | `feat/mobile-redesign` (#6) |

Pendências de validação manual ainda abertas estão no [fim do documento](#pendências-abertas).

---

## 1. Tipos de check (HTTP, TCP, Heartbeat) + redesign do painel

**2026-10-27** (backend `d593fe3`, frontend `996c542`), merge do PR #2 em 2026-10-28.

### Entregue

- Monitores de três tipos: `http`, `tcp`, `heartbeat`. Sem `type` na criação, cria `http`.
- **Migration 0003**: recria `monitors` (URL deixa de ser obrigatória; intervalo vai até 1 semana = 604800 s; `config` JSON por tipo; `last_ping_at`). Checks e incidentes são preservados.
- `checker.Registry` despacha por tipo, sem mexer no `Scheduler`. `TCPChecker` mede o tempo de conexão e grava erro legível em `last_error`.
- **Heartbeat** (direção invertida): rota pública `GET|POST /ping/{token}` (200 e check `UP`; token inválido 404; pausado 200 sem registrar; excesso por IP 429), avaliador de prazos em um único ticker, token regenerável (o antigo morre na hora).
  - Sem ping além de `interval + grace` → `DOWN` **uma vez** (um incidente, sem checks repetidos); volta a `UP` no próximo ping.
  - Antes do primeiro ping fica `UNKNOWN`.
- Webhook ganha `type` e `target`; `url` só para HTTP; o token **nunca** vai no payload nem no log (`/ping/[redacted]`).
- `GET /api/monitors` devolve `history` (últimos 30), `uptime_24h` e `last_latency_ms` sem N+1.
- CLI `add` aceita `--type/--host/--port/--grace`.
- Painel, Novo monitor e Editar monitor refeitos para os 3 tipos; contador de incidentes abertos na sidebar.

### Decisões

- **Down da 0003** falha se existir monitor `tcp`/`heartbeat` (apagar em silêncio seria pior; o backup é o rollback real).
- **Heartbeat retomado** da pausa: o prazo recomeça de "agora".
- `UpdateWithTx` (state machine) grava só estado; `Update` (API) grava só campos editáveis, para um não desfazer o outro.
- Ordenação de checks por `checked_at DESC, id DESC` (resolução de segundos causava contagem errada de falhas seguidas).
- `enabled` ausente na criação = `true`.
- URL de ping montada no front com `window.location.origin`; o CLI imprime só `/ping/<token>`.
- Busca no painel só aparece com 10+ monitores.

### Limitações conhecidas

- Uptime 24 h do heartbeat é sobre checks gravados (1 por ping + 1 por queda), não sobre tempo.
- Mensagens de validação da API continuam em inglês.

### Descartado

URL fictícia `heartbeat://id`; uma tabela por tipo; um job por heartbeat; heartbeat respeitando `failure_threshold` (o `grace_seconds` já faz esse papel); lib de rate limit; websocket para o painel.

---

## 2. Página de Incidentes

**2026-10-28** (plano) → **2026-10-29** (`a1d0e70`, `8aac267`, `0da71e8`), merge do PR #3 em 2026-10-29.

### Entregue

- **Migration 0004**: coluna `failure_count` em `incidents` + índice; backfill e limpeza do prefixo `Get "<url>": ` em erros já gravados.
- `GET /api/incidents` com `?status=open|resolved` e `?monitor_id=` combináveis (`?open=true` segue valendo; `status` inválido → 400), `total` real via `COUNT(*)`, e `monitor {id,name,type,target}` + `failure_count` numa única consulta (sem N+1).
- Erro HTTP limpo (`requestErrorMessage`); timeout aparece como `context deadline exceeded (timeout 10s)`.
- `Monitor.Target()` extraído do notifier (reuso no webhook e na API).
- State machine incrementa `failure_count` enquanto o incidente está aberto.
- Retention apaga incidentes **resolvidos** antigos; abertos nunca.
- Badge da sidebar passou a refletir o total real (antes mostrava 1 com 2 abertos).
- `Incidents.jsx` refeito: cards de abertos ("Fora do ar há 1 dia e 4 h", falhas seguidas, erro em caixa, "Abrir URL"), resolvidos em tabela agrupada por dia com duração, abas Todos/Abertos/Resolvidos, filtro por monitor, "Carregar mais"; estados vazio/carregando/erro. Filtro e aba ficam na URL.
- "Ver histórico" do menu ⋮ abre `/incidents?monitor=<id>`.

### Decisões

- `last_error` guarda o erro da **abertura** do incidente.
- Contagem de falhas gravada na coluna, não por `COUNT(*)` em `checks` (o retention apagaria checks e distorceria incidentes longos).
- Datas formatadas à mão (o `Intl` em pt-BR gera "28 de out."); quedas de mais de 1 dia mostram a data de início.

### Descartado

Endpoint `/api/incidents/summary`; paginação por cursor; recriar a tabela `incidents`; date-fns/dayjs.

---

## 3. Página de Configurações

**2026-10-29** (`e8caa18` backend, `4654a77` frontend, `13e97fb` docs), merge do PR #4 em 2026-10-29.

### Entregue

- `GET /api/info` → `version`, `commit`, `build_date`, `database_size_bytes` (banco + `-wal` + `-shm`).
- `GET /api/export` com `Content-Disposition: attachment; filename="pingou-YYYY-MM-DD.db"` (UTC); o front usa esse nome.
- `POST /api/settings/webhook/test`: testa a URL do corpo (ou a salva), responde sempre `200 {ok, status_code, duration_ms, error}`; `400` sem URL ou esquema não http(s). Payload com `"event": "test"`, nada é gravado.
  - Timeout 5 s, sem redirects, lê no máximo 4 KB, não devolve o corpo do destino, loga só o host, 5 testes/min por IP (429 com `Retry-After`).
- `PATCH /api/settings` agora recusa `webhook_url` que não seja http/https (aceitava `ftp://`).
- Página nova com três cards (Notificações, Retenção com chips de 7 a 90 dias, Backup), "Salvar alterações" com dirty state, rodapé com versão/commit/licença.
- Cópia da retenção corrigida: apaga checks e incidentes resolvidos; abertos nunca (roda na subida e a cada 1 h).

### Decisões

- SSRF para rede interna no teste do webhook: **aceito** (bloquear quebraria n8n em rede local; o envio real já chama qualquer URL salva).
- Tamanho do backup mostrado como "≈" usando `database_size_bytes` (o `VACUUM INTO` só para medir seria desperdício).
- `/api/session` e login ficaram fora do escopo.

### Descartado

`x/time/rate`; bloquear IPs privados; reaproveitar `send()` do notifier (timeout/redirect do envio real mudariam); `BuildInfo` como parâmetro de `NewServer`.

---

## 4. Detalhe do monitor

**2026-10-29** (`e5b4ff0` API, `9428a74` frontend, `87d0115` docs), merge do PR #5 em 2026-10-30, com ajustes de 10-30 (rótulos do eixo X, datas com hora, link do menu ⋮).

### Entregue

- `GET /api/monitors/{id}/stats?range=24h|7d|14d|30d`: `uptime`, `avg_latency_ms`, `checks`, `incidents`, `downtime_seconds`, `series` e `down`.
  - Baldes de 10 min / 1 h / 3 h / 6 h → 144 / 168 / 112 / 120 pontos, vazios incluídos (`avg_latency_ms: null` = buraco na linha).
  - Média só de checks com sucesso; incidente que começou antes da janela ou segue aberto entra recortado à janela.
  - `range` inválido → 400; monitor inexistente → 404. ~35 ms com 43 200 checks (índice coberto).
- `CheckRepo.Buckets`, `IncidentRepo.FindInWindow`, `MonitorService.Stats`.
- Rota `/monitors/:id`: cabeçalho com estado, ações Pausar/Editar/Excluir, 4 cards, gráfico SVG de resposta com faixas "fora do ar" e tooltip (hover, toque, setas do teclado), Últimos checks (8 + "ver mais"), Incidentes recentes (3 + "Ver todos"), Configuração.
- `RangeTabs`: períodos acima de `retention_days` ficam `aria-disabled` com tooltip (hover e foco) explicando o motivo e link para Configurações; `?range=` acima da retenção cai no maior permitido. Período fica na URL.
- Heartbeat sem card de resposta média e sem gráfico. Estados: carregando, não encontrado, sem checks, pausado, erro com "Tentar de novo".
- `MonitorModals` extraído do Dashboard; nome do monitor no painel e nos incidentes leva à página.

### Desvios na implementação

- Incidente aberto termina em `now`, não no fim do balde (senão contaria minutos futuros).
- Escala Y em 3 passos de 100 ms arredondados.
- "ver mais" aumenta o `limit` (8 → 28 → … até 200) em vez de paginar por `offset`.
- `RANGES`/`clampRange` em `lib/ranges.js` (regra de lint do fast refresh).

### Descartado

Estatísticas dentro de `/monitors/{id}`; `IncidentFilter.Since`; percentis (p95); `disabled` nativo nos períodos; esconder períodos acima da retenção; API rejeitar `range` acima da retenção; `title` nativo como tooltip.

---

## 5. Remodelagem mobile

**2026-10-30** (`b909230` → `4ebf4d8`), merge do PR #6 em 2026-10-30.

### Entregue

- Mobile = `< 768 px` (`md`), só com classes CSS (`hidden md:flex`, `md:hidden`); sem `matchMedia`. **Desktop inalterado** (comparado pixel a pixel a 1280 e 800 px).
- API: `GET /api/monitors` ganhou `last_error` e `down_since` (só DOWN, uma query para a lista via `FindOpenByMonitors`).
- Shell: `MobileHeader` + barra inferior com 3 abas e badge de incidentes; drawer e hambúrguer removidos.
- Painel: banner compacto, chips roláveis, cards por monitor nos 4 estados (DOWN, UP, UNKNOWN, heartbeat).
- Incidentes: abas segmentadas, filtro com rótulo, abertos empilhados ("Ver monitor" + "Abrir URL"), resolvidos agrupados por dia.
- Configurações: teste do webhook com resultado em caixa, chips numéricos, aviso de impacto ao reduzir retenção, barra "Alterações não salvas" (Descartar/Salvar), botão Sair.
- Detalhe: barra `‹ Monitores` + `⋮`, cartão de alerta só quando DOWN, Pausar/Editar lado a lado, métricas 2×2, checks em lista.
- Conferido em 360 / 390 / 430 / 600 px sem scroll horizontal.

### Decisões

- Aviso de retenção também no desktop (protege ação irreversível).
- `main` com `pb-24` para o fim da lista não ficar sob a barra inferior.
- Fora do mobile por não estarem nos mockups: busca, menu `⋯` no card, "Baixar backup (.db)" com nome de arquivo.

### Descartado

Rotas `/m/...`; hook `useIsMobile` (flash no primeiro render); lib de bottom sheet; manter o drawer junto da barra; menu `⋯` no card.

---

## Pendências abertas

- **Migrations 0003 e 0004** foram testadas em bancos de teste e de QA, mas **não** rodadas no `pingou.db` real. Fazer backup antes e conferir `failure_count` / `last_error`.
- **Configurações (T12)**: webhook.site real, destino que não responde (5 s), navegação só por teclado, nome do arquivo baixado no navegador.
- **Mobile**: teste em dispositivo real (teclado virtual com a barra "Alterações não salvas", safe-area no iPhone, leitor de tela).
- **Observação de infraestrutura**: o incidente do monitor `c10` vinha de certificado de outro domínio servido em `api.c10.site` (provável proxy roteando errado). Não é bug do Pingou.
