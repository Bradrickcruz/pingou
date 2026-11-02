# Histórico de implementações (planos aplicados)

Consolida o que os planos `plan-*.md` (já removidos; ficam no `git log`) entregaram. Datas vêm do `git log` (commits e merges de PR). O comportamento atual de API e CLI está no [README.md](../README.md); aqui fica o **o quê, quando e por quê**.

| Data (commits)     | Plano                                                                                           | Branch / PR                     |
| ------------------ | ----------------------------------------------------------------------------------------------- | ------------------------------- |
| 2026-10-27 → 10-28 | [Tipos de check + redesign do painel](#1-tipos-de-check-http-tcp-heartbeat--redesign-do-painel) | `feat/check-types` (#2)         |
| 2026-10-28 → 10-29 | [Página de Incidentes](#2-página-de-incidentes)                                                 | `feat/incidents-page` (#3)      |
| 2026-10-29         | [Página de Configurações](#3-página-de-configurações)                                           | `feat/settings-page` (#4)       |
| 2026-10-29 → 10-30 | [Detalhe do monitor](#4-detalhe-do-monitor)                                                     | `feat/monitor-details` (#5)     |
| 2026-10-30         | [Remodelagem mobile](#5-remodelagem-mobile)                                                     | `feat/mobile-redesign` (#6)     |
| 2026-11-01         | [E-mail (SMTP) e envio fora do lock](#6-e-mail-smtp-e-envio-fora-do-lock)                       | `feat/email-notifications` (#9) |
| 2026-11-01         | [Largura no desktop + busca parcial](#7-largura-no-desktop--busca-parcial)                      | `feature/largura-busca` (#10)   |
| 2026-11-01         | [Notificação via ntfy](#8-notificação-via-ntfy)                                                 | `feat/ntfy-notifications` (#11) |
| 2026-11-01         | [Eixo Y dinâmico + gráfico de pings](#9-eixo-y-dinâmico--gráfico-de-pings)                      | `develop` (sem PR)              |
| 2026-11-01 → 11-02 | [Pausar monitor fecha o incidente](#10-pausar-monitor-fecha-o-incidente)                        | `develop` (sem PR)              |

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

## 6. E-mail (SMTP) e envio fora do lock

**2026-11-01** (`eadc1b6`, docs `65ca5d6`), merge do PR #9 em 2026-11-01.

### Entregue

- **Fase 0 (correção):** o envio da notificação rodava dentro do lock global do banco, então um webhook lento travava todos os monitores. Agora `MultiNotifier` ([multi_notifier.go](../internal/service/multi_notifier.go)) envia cada canal em goroutine própria (falha ou panic de um não afeta os outros) e o shutdown espera até 10 s (`Wait`). `StateMachine` e `Scheduler` não mudaram.
- Canal de e-mail com a stdlib (`net/smtp`, `crypto/tls`, `net/mail`), texto simples, modos `none`, `starttls` (obrigatório, sem rebaixar) e TLS direto. Sem dependência nova.
- Servidor SMTP **todo no ambiente** (`PINGOU_SMTP_*`, inclusive a senha); `config.LoadSMTP()` roda só no `serve`, para um SMTP inválido não travar o `migrate`. A tela guarda só os destinatários (`settings.email_recipients`, até 10, separados por vírgula). Sem migration.
- `POST /api/settings/email/test` (atrás do login, 5 por minuto por IP). Sem SMTP no ambiente, o card aparece bloqueado listando as variáveis que faltam.
- `EmailCard`; botão e resultado de teste extraídos para `TestControls.jsx`, compartilhados com o card do webhook.
- Segurança conferida: senha só no `PlainAuth` (fora de log, `pingou config`, `GET /settings` e export); cabeçalhos sem quebra de linha (`oneLine`, teste com `\r\nBcc:`); `tls.Config` sem `InsecureSkipVerify`, TLS 1.2+; log só com `recipients_count`; e-mail nunca leva o `ping_token`.

### Decisões

- Um canal de cada tipo (1 webhook + 1 e-mail com vários destinatários), em vez de tabela de canais.
- Subir com SMTP incompleto **falha** (mesmo padrão do `PINGOU_API_KEY`), em vez de desabilitar o card em silêncio.
- A goroutine dentro do `MultiNotifier` resolve o lock sem mudar a assinatura do `Process`.

### Limitações conhecidas

- Sem retry nem outbox: aviso perdido em crash não é reenviado. Próximo passo natural: tabela `notification_outbox` gravada na transação do `Process` + worker.
- `webhook_url` também é segredo (Slack/Discord) e continua indo no export.

### Descartado

Tabela `notification_channels` com CRUD; `NotificationEvent` + `Channel` + `Dispatcher`; `go-mail`; host/porta/usuário editáveis na tela; `email_available` em `/api/info`; destinatários como array JSON; tirar o `notify()` de dentro do `withMonitor`.

---

## 7. Largura no desktop + busca parcial

**2026-11-01** (`50c642a`, ajuste `7ab5185`), merge do PR #10 em 2026-11-01. Só frontend.

### Entregue

- Limite de largura de Painel, Incidentes e Detalhe subiu de 1180 px para **1680 px** (`max-w-[1680px]`); Configurações segue em 760 px. Mobile inalterado.
- `matchesQuery` em [monitors.js](../web/src/lib/monitors.js) (com testes): busca parcial sem acento nem caixa ("servico" acha "Serviço") no nome e no alvo (URL ou `host:porta`); heartbeat só pelo nome.
- A busca fica na URL (`?q=`), com mensagem de vazio "Nenhum monitor encontrado para “…”" e botão "Limpar busca".
- Placeholder passou a "Buscar por nome ou endereço". O campo ficou sempre visível (`7ab5185` removeu o limite de 10 monitores, `SEARCH_FROM`).

### Decisões

- Evolui o campo que já existia; filtro combinado (E) com os de estado.
- Aumentar o limite, não remover (linhas ilegíveis em ultrawide).

### Descartado

Busca no backend (a lista já traz até 100 monitores); debounce (filtro em memória).

---

## 8. Notificação via ntfy

**2026-11-01** (`bf317e5` e demais, docs `094b94f`), merge do PR #11 em 2026-11-01.

### Entregue

- Terceiro canal, no molde do e-mail. Servidor no ambiente: `PINGOU_NTFY_URL` (liga o canal) e `PINGOU_NTFY_TOKEN` (opcional); a tela guarda só o tópico (`settings.ntfy_topic`; vazio = desligado). Sem migration, sem dependência.
- Envio: POST de texto em `<URL>/<tópico>` com `Authorization: Bearer`, `Title` (RFC 2047), `Priority`, `Tags` (emoji) e `Click` (link do painel). Queda: prioridade 4; volta: prioridade 3. Texto reaproveita as funções do e-mail.
- `POST /api/settings/ntfy/test` (5 por minuto por IP); `GET/PATCH /api/settings` ganham `ntfy_topic` e `ntfy_available`. Sem `PINGOU_NTFY_URL`, `PATCH` do tópico e o teste dão 400 e o card fica bloqueado.
- `LoadNtfy` recusa URL inválida, com query, fragmento ou usuário (o token vai só em `PINGOU_NTFY_TOKEN`). 401/403 no teste ganham a dica "token inválido ou sem permissão no tópico".
- `slog`, API, export e `pingou config` nunca expõem token, URL completa nem tópico (só `NtfyHost` e `NtfyTokenSet`); o `ping_token` não vai na mensagem.
- README ganhou a seção geral "Entrega das notificações", válida para os três canais.
- Validado com ntfy 2.36.0 em Docker (`deny-all` + token): queda, volta, heartbeat expirado, token errado, tópico inválido, 429.

### Decisões

- Token e servidor no ambiente (não no banco): o export já vaza o segredo do webhook, não repetir.
- Auth por header, não por `?auth=` (não aparece em log de proxy).

### Descartado

Só documentar a URL com template no README; templates genéricos de webhook; token/servidor editáveis na tela; JSON na API do ntfy.

---

## 9. Eixo Y dinâmico + gráfico de pings

**2026-11-01** (`3e8b970`, `d636b4a`), direto em `develop`. Só frontend.

### Entregue

- `niceScale(peak, floor)` em [chartScale.js](../web/src/lib/chartScale.js) (com testes): passo em {1,2,5}×10ⁿ com no máximo 5 intervalos. 70 → 100, 350 → 400, 2300 → 2500; piso de 100 ms em HTTP/TCP. Substitui o `scaleOf` (3 passos de 100 ms, que dava 402 → 600).
- Heartbeat ganhou o card **"Pings recebidos"**: barras por intervalo (`checks - failures`, porque prazo perdido é check com falha), eixo inteiro com piso 1, sem faixa cinza (zero ping é normal); vermelho só nos períodos "fora do ar". Tooltip, legenda e `aria-label` falam de pings.
- O eixo é derivado de `stats`, então trocar 24h/7d/14d/30d recalcula.

### Decisões

- O plano dizia 101 → 120 (passo 20), mas pela própria regra o passo é 50 e dá 150; o teste usa 150 (e 402 → 500).
- Sem componente novo: o modo heartbeat é uma prop do `ResponseChart`.

### Descartado

Lib de gráfico; p95 como teto; manter 3 intervalos fixos.

---

## 10. Pausar monitor fecha o incidente

**2026-11-01** (`19e0f24`), testes em **2026-11-02** (`33e492a`), direto em `develop`.

### Entregue

- **Bug corrigido:** pausar um monitor `DOWN` deixava o incidente aberto para sempre e o monitor em `DOWN`, porque só uma transição `DOWN→UP` fecha incidente e o scheduler não checa monitor pausado.
- `MonitorService.Update` com `enabled:false` (único caminho de pausa) fecha o incidente aberto **sem notificar** (e-mail, webhook e ntfy ficam de fora: o service não tem acesso a notifier) e devolve o monitor a `UNKNOWN`.
- **Migration 0005**: `incidents.closed_manually INTEGER NOT NULL DEFAULT 0`. Incidentes existentes ficam `0`.
- `IncidentRepository.CloseManually(ctx, monitorID, endedAt)`: um `UPDATE` só, sem leitura prévia, que preenche `ended_at`, `closed_manually = 1` e `duration_seconds` (início até a pausa, mínimo 0). Não faz nada se não há incidente aberto.
- `GET /api/incidents` expõe `closed_manually`. `status=resolved` continua incluindo os fechados manualmente.
- `ResolvedTable.jsx`: selo neutro "Fechado manualmente" (tabela e cartão mobile).
- Testes em `monitor_pause_test.go` e `incidents_test.go`: fecha com duração e `UNKNOWN`, pausa repetida, pausa sem incidente, nenhuma notificação na pausa e retomada abrindo incidente novo.

### Decisões

- Booleano, não enum: cobre o único caso novo. Vira enum se surgir um terceiro motivo de fechamento.
- `ended_at` preenchido, para "aberto = `ended_at IS NULL`" continuar valendo em todas as queries e as stats de uptime não mudarem.
- Fechamento **idempotente**: roda a cada `Update` com `enabled:false`. Se a segunda escrita falhar, o erro volta ao usuário e repetir a pausa conclui.
- O estado precisa ir a `UNKNOWN`: `handleFailure` só abre incidente se `prevState != DOWN`, então retomar um monitor ainda `DOWN` nunca abriria outro.
- Check em andamento que termina depois da pausa não reabre incidente (`scheduler.withMonitor` descarta resultado de monitor desabilitado).

### Limitações conhecidas

- Após retomar, a **primeira falha já reabre o incidente**: `countConsecutiveFails` lê o histórico de checks, que a pausa não zera. Aceito, pois o monitor segue fora do ar. Exigir o threshold inteiro pede uma coluna com o instante da pausa.
- Pausa e fechamento são duas escritas sem transação (ver idempotência acima). Atomicidade de verdade pede mover o fechamento para dentro do `UPDATE` do monitor.
- Incidentes que ficaram órfãos em pausas anteriores à 0005 continuam abertos.

### Descartado

Reutilizar a `StateMachine` para fechar (acoplada a notifier e a resultado de check); coluna enum `resolution`; filtro `status=closed_manually`; deixar `ended_at` nulo e marcar só por flag.

---

## Pendências abertas

- **Migrations 0003 e 0004** foram testadas em bancos de teste e de QA, mas **não** rodadas no `pingou.db` real. Fazer backup antes e conferir `failure_count` / `last_error`.
- **Configurações (T12)**: webhook.site real, destino que não responde (5 s), navegação só por teclado, nome do arquivo baixado no navegador.
- **Mobile**: teste em dispositivo real (teclado virtual com a barra "Alterações não salvas", safe-area no iPhone, leitor de tela).
- **Observação de infraestrutura**: o incidente do monitor `c10` vinha de certificado de outro domínio servido em `api.c10.site` (provável proxy roteando errado). Não é bug do Pingou.
- **E-mail (T14)**: QA manual com Mailpit e provedor real (587/465); o ambiente de desenvolvimento não tinha Docker nem SMTP real.
- **Eixo Y**: verificação manual (T5) dos 4 ranges nos 3 tipos, tema claro/escuro e mobile; conferir em um `/stats` real que `checks`/`failures` do heartbeat batem com pings e prazos perdidos.
- **Pausar fecha o incidente**: QA manual (T7) no navegador (HTTP/TCP e heartbeat: incidente sai de "Abertos", aparece em "Resolvidos" com selo, sem e-mail/ntfy; retomar reabre) e `npm run lint`/`build` do frontend. Decidir se vale limpar incidentes órfãos antigos de monitores já pausados.
