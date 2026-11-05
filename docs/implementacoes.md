# Histórico de implementações (planos aplicados)

Consolida o que os planos `plan-*.md` (já removidos; ficam no `git log`) entregaram. Datas vêm do `git log` (commits e merges de PR). O comportamento atual de API e CLI está no [README.md](../README.md); aqui fica o **o quê, quando e por quê**.

| Data (commits)     | Plano                                                                                               | Branch / PR                     |
| ------------------ | --------------------------------------------------------------------------------------------------- | ------------------------------- |
| 2026-10-27 → 10-28 | [Tipos de check + redesign do painel](#1-tipos-de-check-http-tcp-heartbeat--redesign-do-painel)     | `feat/check-types` (#2)         |
| 2026-10-28 → 10-29 | [Página de Incidentes](#2-página-de-incidentes)                                                     | `feat/incidents-page` (#3)      |
| 2026-10-29         | [Página de Configurações](#3-página-de-configurações)                                               | `feat/settings-page` (#4)       |
| 2026-10-29 → 10-30 | [Detalhe do monitor](#4-detalhe-do-monitor)                                                         | `feat/monitor-details` (#5)     |
| 2026-10-30         | [Remodelagem mobile](#5-remodelagem-mobile)                                                         | `feat/mobile-redesign` (#6)     |
| 2026-11-01         | [E-mail (SMTP) e envio fora do lock](#6-e-mail-smtp-e-envio-fora-do-lock)                           | `feat/email-notifications` (#9) |
| 2026-11-01         | [Largura no desktop + busca parcial](#7-largura-no-desktop--busca-parcial)                          | `feature/largura-busca` (#10)   |
| 2026-11-01         | [Notificação via ntfy](#8-notificação-via-ntfy)                                                     | `feat/ntfy-notifications` (#11) |
| 2026-11-01         | [Eixo Y dinâmico + gráfico de pings](#9-eixo-y-dinâmico--gráfico-de-pings)                          | `develop` (sem PR)              |
| 2026-11-01 → 11-02 | [Pausar monitor fecha o incidente](#10-pausar-monitor-fecha-o-incidente)                            | `develop` (sem PR)              |
| 2026-11-02         | [Fuso horário + desligamento ordenado](#11-fuso-horário--desligamento-ordenado)                     | `develop` (sem PR)              |
| 2026-11-02         | [Índices e desempenho do banco](#12-índices-e-desempenho-do-banco)                                  | `ccr-dd100985-raj41c` (#15)     |
| 2026-11-02         | [Aviso de expiração do certificado SSL](#13-aviso-de-expiração-do-certificado-ssl)                  | `feature/ssl-expiry` (#16)      |
| 2026-11-04         | [Manutenção: pausa com motivo e janela agendada](#14-manutenção-pausa-com-motivo-e-janela-agendada) | `feature/maintenance-windows`   |
| 2026-11-05         | [Modo escuro (Claro, Escuro, Sistema)](#15-modo-escuro-claro-escuro-sistema)                        | `ccr-f77fa1bb-9acipc`           |
| 2026-11-05         | [Retry de notificações (outbox)](#16-retry-de-notificações-outbox)                                  | `feature/notification-retry`    |

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

## 11. Fuso horário + desligamento ordenado

**2026-11-02** (desligamento `65ecef3`, fuso `39e5f55`, painel `88f38da`, docs `403c89c`; detecção do navegador `bad6890`/`bf92800`; combobox `d4a4ee9`), direto em `develop`.

### Entregue

**Fuso horário**

- `settings.timezone` (nome IANA) em `GET/PATCH /api/settings`, validado com `time.LoadLocation`: vazio vira `UTC`, máximo de 64 caracteres, `Local` e nomes desconhecidos dão 400. Sem migration (as settings são chave-valor; chave ausente = `UTC`).
- `import _ "time/tzdata"` no `main`: os fusos vão no binário (~450 KB) e funcionam em imagem mínima.
- E-mail e ntfy leem o fuso **a cada envio** (`service.LocationLoader`), então mudar vale na próxima notificação, sem reiniciar. Layout `2006-01-02 15:04:05 MST` (`UTC`, `-03`, `+09`) em `Caiu em` e `Voltou em`. Fuso salvo que não carrega cai em UTC com aviso no log, e a notificação sai mesmo assim.
- Painel: `lib/timezone.js` (formatação com `Intl.DateTimeFormat`), `format.js` e eixo do gráfico no fuso salvo, `TimeZoneProvider` carrega o fuso antes de mostrar o painel; trocar vale sem recarregar.
- Card **Fuso horário** em Configurações: combobox digitável (`ui/Combobox.jsx`, ARIA 1.2, ↑ ↓ Enter Esc) que filtra a lista ao digitar. O filtro ignora acento, caixa, `/` e `_` e aceita palavras em qualquer ordem ("sao paulo" acha `America/Sao_Paulo`). Mostra a hora atual no fuso escolhido antes de salvar. Registrado no `DESIGN.md`.
- **Fuso do navegador na primeira abertura:** `GET /api/settings` ganhou `timezone_set` (`false` = nunca salvo). Nesse caso o painel salva o fuso do navegador, uma vez; sem fuso válido, fica em `UTC`. Quem escolheu `UTC` de propósito já tem `timezone_set: true` e não é sobrescrito.
- Continua sempre em UTC: API (`*_at` em RFC3339), banco, JSON do webhook (`timestamp`) e baldes do `/stats`.

**Desligamento ordenado**

- `runShutdown` com orçamento de **15 s** para tudo, na ordem: scheduler (nenhum check novo) → servidor HTTP (inclui `/ping/{token}` em andamento) → checks e retenção → notificações → banco → lock. Uma etapa que falha ou estoura não pula as seguintes: banco e lock sempre fecham. Logs `shutdown started`, uma linha por etapa e `shutdown complete`.
- Código de saída `0` se limpo, `1` se alguma etapa falhou. Segundo sinal durante o desligamento sai na hora com `1`. Falha ao subir o servidor também sai com `1` (antes saía com 0), e o `ErrServerClosed` normal deixou de virar log de erro.
- `Scheduler.Stop()` (cancela e bloqueia jobs novos) + `Scheduler.Wait(ctx)`; `RetentionWorker.Wait(ctx)`. **Check cancelado pelo desligamento é descartado**, não vira falha nem incidente. A gravação em andamento usa `context.WithoutCancel` e termina inteira. O avaliador de heartbeat não gera `DOWN` durante o desligamento.
- `stop_grace_period: 20s` no `docker-compose.yml` (o padrão do Docker é 10 s e mataria o processo antes). O `ENTRYPOINT` já estava em formato exec, então o binário é o PID 1.
- **Correção extra:** lock com o PID do próprio processo era tratado como "outro processo rodando". Num container reiniciado (PID 1 de novo) depois de um `kill -9`, o Pingou não subia.
- README: seções "Fuso horário" e "Desligamento".

### Decisões

- Fuso é uma configuração **global** (uma por instalação), como o webhook: com uma única API Key, fuso por pessoa não faz sentido.
- O padrão deixou de ser "UTC fixo" durante a execução: salvar o fuso do navegador evita que quem já usava o painel passe a ver UTC, e mantém painel e e-mail no mesmo fuso (por isso o fuso é **salvo**, não só exibido).
- Resultado de check cortado por edição ou pausa também é descartado (antes era gravado se o monitor continuasse ativo); o job novo faz um check na hora, então nada se perde.
- Select nativo trocado por combobox: num `<select>` digitar só pula por prefixo, e quase todos os fusos começam com `America/` ou `Europe/`.

### Limitações conhecidas

- Quem abrir o painel primeiro define o fuso de todos; até esse primeiro acesso, e-mail e ntfy ficam em UTC.
- Baldes do `/stats` seguem alinhados em UTC: em fusos com meia hora (ex.: `Asia/Kolkata`) os rótulos de 1 h não caem em hora cheia. Revisar junto com a futura janela de manutenção.
- Mudar o fuso com notificações em andamento: a mensagem sai com o fuso antigo.
- O teste de integração com sinal real (`kill -TERM` num binário com check lento) foi feito à mão, não automatizado; os testes cobrem `runShutdown` e o descarte do check cancelado.
- `PATCH /api/settings` grava os campos em sequência: com fuso inválido, os campos anteriores já foram gravados e a resposta é 400 (padrão que já existia nos outros campos).

### Descartado

Migration de coluna `timezone` (T2: as settings são chave-valor); `PINGOU_TIMEZONE` por variável de ambiente (exigiria reiniciar); fuso do navegador só na exibição (painel e e-mail mostrariam horários diferentes); `os.Exit` logo após o timeout, sem log; esperar mais de 15 s no desligamento (o Docker mataria o processo antes); fuso por usuário; mudar os baldes do `/stats`.

### Validação

Testes Go (incluindo `-race` em scheduler e comandos) e `npm test` (22) verdes; `npm run lint` e `build` limpos. Binário real com check de 5 s e `kill -TERM` no meio: etapas na ordem esperada, `exit=0`, lock removido, 0 checks e 0 incidentes gravados. Docker testado pelo dono do projeto: `docker compose stop` sem `Killed`, fuso na imagem, **Enviar teste** e queda real no e-mail e no ntfy.

---

## 12. Índices e desempenho do banco

**2026-11-02** (ferramentas `f8d3597`, índice e lista `bdd5b0f`, retenção `59d1c51`, Docker e RAM `d060795`, correção `d3a00b8`, docs `55bbf75`), branch `ccr-dd100985-raj41c`.

### Entregue

- **Medição antes de mudar.** `scripts/seed` gera, de forma determinística, os perfis típico (10 monitores × 60 s, 397 mil checks, 95 MB) e limite (100 × 10 s, 24,7 milhões de checks, 6 GB), pelas migrations reais e com os monitores intercalados em ordem cronológica. `scripts/dbbench.sh` roda `TestDBBench` (p50/p95 por consulta, gravação, concorrência leitura + gravação, retenção) e `TestPrintQueryPlans` numa **cópia** do banco; `scripts/dbmem.sh` mede o pico de RAM do app real. Sem `PINGOU_BENCH_DB` os testes pulam e o `make test` segue rápido. O SQL das consultas quentes foi para `internal/repository/queries.go` (mesmo texto), para o teste de plano usar as mesmas strings.
- **Migration 0006**: cria `idx_checks_monitor_time (monitor_id, checked_at, id, success, latency_ms)` e apaga `idx_checks_monitor_checked`; o down faz o contrário. O índice novo lê de trás para frente na ordem de `ORDER BY checked_at DESC, id DESC` (antes o SQLite ordenava à parte) e cobre `success` e `latency_ms`, então `/stats` e uptime 24 h leem só o índice. Substitui o antigo em vez de somar um segundo. `idx_checks_checked_at` (retenção) ficou como estava.
- **Lista do painel (`Summaries`)**: o histórico de 30 checks deixou de usar `ROW_NUMBER()` sobre todos os checks (lia ~25 milhões de linhas) e passou a `LIMIT` por monitor com `UNION ALL`.
- **State machine**: `CheckRepositoryTx.Latest` busca os últimos checks sem o `COUNT(*)` que era descartado a cada falha.
- **Retenção em lotes** de 5 mil linhas por transação, com `ctx` conferido entre eles (para rápido num desligamento). Apaga o mesmo conjunto de antes.
- **Log de migrations pendentes** (versões e tempo), para quem atualiza um banco grande não achar que o app travou.
- **Dockerfile**: `HEALTHCHECK --start-period=5m`, porque a migration no perfil limite passou de 30 s.
- **Correção**: com o banco já em dia, `goose.CollectMigrations` devolvia `ErrNoMigrationFiles` e o app não subia na segunda vez (`TestOpenTwice`).
- Testes: `TestHotQueriesUseIndexes` (confere o **nome** do índice no `EXPLAIN QUERY PLAN`; falha sem a 0006), `TestSummariesHistoryMatchesLegacy`, `TestResultsSameWithOldIndex` (respostas de `/stats`, lista, checks e incidentes iguais antes e depois do `goose down`), `TestDeleteOlderThanBatches`, `TestOpenTwice`.
- README: "Desempenho e tamanho do banco" e "Atualizando um banco grande". `scripts/README.md`: seção do dbbench.

### Resultados (perfil limite, mesma máquina: Xeon 2,8 GHz, 4 núcleos, 16 GB, SQLite 3.51.3)

| Consulta (p95)                     |                          Antes |                         Depois |
| ---------------------------------- | -----------------------------: | -----------------------------: |
| `/stats` 30d http                  |                   815 a 855 ms |                   208 a 291 ms |
| `GET /api/monitors`                |   91 a 107 s (falhava, > 30 s) |                   114 a 203 ms |
| checks, offset 0 / 1.000 / 10.000  | 36 a 41 / 28 a 34 / 72 a 77 ms | 13 a 22 / 13 a 21 / 13 a 19 ms |
| gravação com falha                 |                          40 ms |                   0,4 a 0,5 ms |
| gravação com `/stats` 30d ao lado  |                   740 a 759 ms |                   247 a 258 ms |
| retenção: maior espera de gravação |                       2.358 ms |                    81 a 104 ms |

- Índices de `checks`: 4.091 MB → 2.719 MB (o novo tem 1.878 MB contra 3.251 MB do antigo; parte do ganho é de índice recém-criado, que tende a crescer com inserções intercaladas). O arquivo foi de 5.972 MB para 7.850 MB, porque o novo é escrito antes de o antigo ser apagado; as páginas livres voltam com `VACUUM`.
- Migration 0006: 1 min 18 s a 1 min 27 s no perfil limite; 0,7 s no típico.
- Custo de gravação com sucesso: +9 % na média (alvo: no máximo +20 %).
- Nenhuma consulta quente faz `SCAN checks`.

### Decisões

- Não criar índices "por precaução": cada candidato precisou de número antes e depois. `checks(checked_at)`, incidentes e `ping_token` já atendiam e não foram tocados.
- Sem `ANALYZE`/`PRAGMA optimize` na migration: o planejador já escolhe os índices certos e o `ANALYZE` varreria 25 milhões de linhas a mais na atualização. O teste de plano roda `ANALYZE` no banco pequeno, para ficar estável.
- Índice cobrindo **substitui** o antigo (mesmo prefixo), então o total de índices não aumenta.

### Limitações conhecidas

- **RAM acima da meta (pendente):** nos testes do perfil limite a RAM ficou **entre 50 e 60 MB**, o que **quebra a meta de < 50 MB** do README. Esse valor será corrigido futuramente; até lá, a meta não vale para esse cenário (25 milhões de checks, painel consultado sem parar). Medições do `scripts/dbmem.sh` que sustentam isso:
  - em regime (banco já migrado, `/stats` 30d de 3 monitores + lista sem parar): pico de 39 a 40 MB;
  - na **primeira subida logo depois da migration**: 47 a 51 MB (RssAnon 31 MB e RssFile 15 MB, contra 27 e 11 em regime), que some ao reiniciar.
  - O teste do README (100 monitores por 1 hora, sem banco grande) segue em 29 MB.
- A lista do painel deu 203 ms de p95 em 1 de 3 rodadas (alvo de 200 ms); p50 de 82 a 135 ms.
- A retenção em lotes leva mais tempo total (6,0 s contra 2,5 a 2,7 s para apagar 1 dia), em troca de não segurar a única conexão.
- `HTTPS_N=50 scripts/memtest.sh 3600` foi interrompido em ~6 min (pico de 27 MB, RSS estável em 23 a 24 MB); a hora completa não foi rodada.
- Paginação com `offset` grande segue sem cursor (fora do escopo).

### Descartado

Tabela de agregados por balde (rollup), revisar só se os alvos falharem; pool de leitura separado; índice cobrindo desde o início com dois índices; `ANALYZE` na migration; gerar os dados com o `sqlite3` CLI; copiar o SQL para o arquivo de benchmark (sai de sincronia com o código).

---

## 13. Aviso de expiração do certificado SSL

**2026-11-02** (migration `b9087c2`, checker `8777d30`, regra `c60d702`, API `7f14563`, notificações `a276a06`, state machine `b63ae81`, painel `56ab8b5`, `certserver` `1e929e8`, docs `f5ab7dd`), merge do PR #16 em 2026-11-02.

### Entregue

- **Aviso 30 e 7 dias antes** de o certificado expirar, **uma vez por limite**, em todos os canais configurados (webhook, e-mail e ntfy). Opcional por monitor, só HTTP com URL `https://`, desligado por padrão.
- **Sem conexão nem check a mais:** o checker HTTP lê `resp.TLS.PeerCertificates[0].NotAfter` da resposta que já usa (`CheckResult.CertNotAfter`, em UTC). Preenche também quando o status HTTP é de erro (ex.: 500), para a data continuar atualizada durante um incidente. Falha de conexão ou de TLS não traz resposta e deixa o campo vazio. TCP e Heartbeat não mudam.
- **Migration 0007**: 4 colunas em `monitors`: `cert_alert_enabled` (padrão 0), `cert_expires_at`, `cert_alerted_30_for` e `cert_alerted_7_for` (aceitam NULL). Monitores antigos ficam com a opção desligada. Down remove as colunas (`DROP COLUMN`).
- **Regra pura** `EvaluateCertAlert(now, expiresAt, alerted30For, alerted7For)` em `service/cert_alert.go`, sem banco, relógio global nem envio:
  - `restante <= 7 d` (inclui vencido) e 7 d ainda não avisado → aviso de 7 dias, marcando **os dois** limites;
  - senão, `restante <= 30 d` e 30 d ainda não avisado → aviso de 30 dias;
  - senão, nada. O aviso leva dias restantes (arredondados para baixo) e o limite.
- **Evento `cert_expiring`** nos três canais (`NotifyCertExpiring`). Webhook: mesmo formato dos outros eventos, com `cert_expires_at`, `days_left` e `threshold_days`; `down` e `up` não mudam. E-mail e ntfy seguem o estilo atual; ntfy usa prioridade normal no aviso de 30 dias e alta no de 7. Datas no fuso configurado.
- **Integração na state machine** (`Process`), dentro da mesma transação do check: atualiza `cert_expires_at` quando muda, avalia a regra, grava as marcas e só **depois do commit** envia o aviso. Rollback não envia nada nem grava marca. Log `info` ao enviar (monitor, limite, data).
- **Escrita guardada pela URL:** os campos de certificado têm `UpdateCertWithTx` (`UPDATE ... WHERE id = ? AND url = ?`), separado do `UpdateWithTx`, que grava só estado. Se o usuário trocar a URL durante um check, a escrita é ignorada. `StateMachine` ganhou `WithClock` para testes.
- **API:** `POST`/`PATCH /api/monitors` aceitam `cert_alert_enabled`; `GET /api/monitors` e `GET /api/monitors/:id` devolvem `cert_alert_enabled` e `cert_expires_at` (RFC 3339, `null` se ainda não visto). As marcas de aviso não saem da API. Ligar em monitor não HTTP ou com URL não `https://` dá 400 `VALIDATION_ERROR`. Trocar a URL zera a data e as marcas dentro do próprio `UPDATE`.
- **Painel:** checkbox "Avisar quando o certificado SSL estiver perto de expirar" no formulário HTTP (desabilitado e desmarcado sozinho se a URL não é https) e bloco **Certificado SSL** no detalhe, com data e dias restantes, "Aguardando o primeiro check" quando vazio, "Aviso ativo (30 e 7 dias antes)" e "último valor conhecido" enquanto o monitor está DOWN. Sem rolagem horizontal a 390 px. `DESIGN.md` §4.1.
- **`scripts/certserver`**: servidor HTTPS local com CA de teste e validade escolhida (`go run ./scripts/certserver --days 25 --port 9443`), para QA com `SSL_CERT_FILE=<ca.pem>` (Linux).
- README (recurso, campos da API, evento `cert_expiring` com aviso para quem consome o webhook) e `scripts/README.md`.

### Decisões

- **Guardar a data do certificado avisado, não "sim/não":** `cert_alerted_30_for = <data alertada>` rearma sozinho quando o certificado é renovado (a data nova não bate com a guardada). Não precisa de lógica de reset.
- **Ativação tardia:** ligar a opção com 20 dias restantes manda o aviso de 30 dias na hora; com 5 dias manda só o de 7 (não dois seguidos). Desligar e ligar de novo não repete o mesmo aviso.
- **A data aparece em todo monitor HTTP https**, mesmo com a opção desligada (ler é grátis); o checkbox controla só o aviso.
- **Certificado lido é o da resposta final**, depois dos redirects: é o que o visitante vê.
- **Sem escrita extra por check:** a data só é gravada quando muda e as marcas só quando um aviso sai. Importa pela conexão única do SQLite. Medido com trigger: 1 escrita de certificado em 23 checks.
- Sem retry, como os outros canais: a redundância dos dois limites cobre um aviso perdido.

### Desvios do plano

1. `cert_alert_enabled` e `cert_expires_at` também saem na **lista** `GET /api/monitors` (o plano previa só o detalhe). O formulário de edição abre a partir da lista; sem o campo, salvar pelo painel desligaria a opção sem querer. Custo: ~60 bytes por monitor.
2. Trocar o tipo do monitor não se aplica: o `PATCH` já recusa troca de tipo.
3. `DESIGN.md` não tem token de "atenção". O nível de 30 dias usa o neutro (`unknown`) com o texto "ATENÇÃO"; vermelho fica para 7 dias ou vencido. Uma cor âmbar precisaria de aprovação (§3.4).

### Limitações conhecidas

- **Renovação automática (Let's Encrypt/Certbot)** renova quando faltam ~30 dias, o limite do primeiro aviso. Quem usa isso pode receber o aviso de 30 dias pouco antes da renovação; o de 7 dias é o que indica falha. Saída, se incomodar: carência de ~24 h antes do aviso de 30 dias (exige guardar quando o limite foi cruzado).
- Só o certificado do próprio site (não os intermediários), só HTTP com `https://` (não TCP com TLS). Limites fixos em 30 e 7; sem flag no CLI nem badge na lista.
- Balanceador ou CDN com certificados diferentes por servidor pode fazer a data oscilar e repetir o aviso, só quando volta a um certificado já avisado depois de outro.
- Certificado já vencido: o Go recusa o handshake, então o monitor cai como DOWN pelo fluxo normal e a data não é relida (fica o último valor conhecido).
- Monitor pausado não faz check, então não gera aviso.

### Descartado

Conexão TLS separada só para ler o certificado; gravar "enviado sim/não" por limite; `cert_checked_at` gravado a cada check (uma escrita por check); dois avisos juntos na ativação tardia; limites configuráveis; mostrar a data só com a opção ligada; capturar o certificado do primeiro host com `httptrace`; flag `--cert-alert` no CLI; badge na lista.

### Validação

`go vet`, `go test ./...` e `npm run lint`/`test`/`build` verdes. Simulação de 90 dias (1 check por hora): exatamente 1 aviso de 30 dias e 1 de 7, no horário certo. QA manual com `certserver`, Pingou real e receptores falsos de webhook, ntfy e SMTP: `--days 25` (aviso de 30 dias nos três canais, sem repetir), `--days 5` (aviso de 7 dias uma vez, ntfy com prioridade alta), `--days 90` (renovação: data atualizada, nenhum aviso novo), troca de URL (data zerada e relida), opção desligada com `--days 3` (data atualizada, nenhum aviso). Detalhe conferido em desktop e celular nos estados válido, atenção, alerta, vencido, aguardando e fora do ar.

---

## 14. Manutenção: pausa com motivo e janela agendada

Plano `plan-maintenance-windows.md` (removido; fica no `git log`), em três etapas: pausa com motivo, janela agendada e página de manutenções com edição.

### Entregue

**Etapa 1: pausa com motivo**

- Migração 0008: `monitors.pause_reason`, `paused_at`, `resumed_at` e `incidents.close_note`. `closed_manually` passa a guardar 0/1/2 (nenhum/pausa/manutenção), sem mudar o esquema.
- `domain.Incident.ClosedManually bool` virou `CloseReason` (enum) + `CloseNote`. A API ganhou `close_reason`/`close_note` e manteve `closed_manually`.
- Modal de pausa com motivo opcional (contador, Enter confirma). Motivo na linha, no card e no detalhe; incidente fechado mostra o selo certo e o motivo.
- **Bug corrigido:** a state machine contava falhas de antes da pausa junto com as de depois. Agora só conta checks a partir de `resumed_at` (consulta nova `sqlChecksLatestSince`, coberta pelo teste de plano).

**Etapa 2: janela agendada**

- Migração 0009: tabela `maintenance_windows` com índice parcial único (uma janela viva por monitor) e índice `(monitor_id, starts_at)`. `ON DELETE CASCADE` com o monitor.
- `MaintenanceService`: criação em lote tudo ou nada, validação (início futuro, fim depois, até 7 dias, motivo até 200), conflito `409` com os monitores, cancelamento.
- Gate único em `Scheduler.withMonitor`: cobre check HTTP/TCP, heartbeat vencido e ping, e por consequência incidente e todos os avisos (inclusive o de certificado). `runCheck` e o avaliador de heartbeat pulam antes, para não fazer requisição à toa.
- Tick de 10 s (`runMaintenance`, no molde do `runHeartbeats`): aplica início (fecha incidente como manutenção, estado `UNKNOWN`) e fim (`resumed_at`, `UNKNOWN`, prazo do heartbeat, `Reload` com check imediato). Roda uma vez no `Start`.
- Web: status `MAINTENANCE` (selo próprio, filtro, contagem na faixa), "Em manutenção até HH:MM" no fuso das Configurações, modal de agendar com multi-seleção, card no detalhe com cancelar/encerrar, `zonedToUtc` para ler o `datetime-local` no fuso das Configurações.

**Etapa 3: página de manutenções e edição**

- `GET /api/maintenance-windows` (janelas vivas com o nome do monitor e `state` = `scheduled` | `active`) e `PATCH /api/maintenance-windows/{id}`.
- Agendada: edita início, fim e motivo. Iniciada (pelo relógio do servidor): só o motivo; horário vira `409 MAINTENANCE_LOCKED`. A guarda `starts_at > now` no próprio `UPDATE` fecha a corrida com o início.
- Página `/maintenance` (menu lateral e aba do celular): lista, editar, cancelar/desligar, vira "Em curso" sozinha e recarrega nas viradas de início e fim.
- `MaintenanceModal` em modo editar: monitor fixo, só envia o que mudou, horário desabilitado em janela iniciada; o 409 mantém o modal aberto com o horário travado.
- "Editar" no card de manutenção do detalhe, com o mesmo modal e o mesmo diálogo de cancelar/desligar da página.
- D10/D11 (no modal de criar: nome do fuso no rótulo; pausados e com janela viva desabilitados) foram commitados direto pelo autor do plano.

### Decisões

D1–D9 do plano: teto de 7 dias; uma janela viva por monitor; pausa manual vence a janela; motivo opcional de até 200 caracteres; lote tudo ou nada; "em manutenção" derivado do relógio; uma instância só; uptime de 24 h sem mudança; `closed_manually` 0/1/2.

D10–D17 (ajustes e Etapa 3):

- **D10:** no modal de agendar, o rótulo de Início e Fim mostra o nome do fuso das Configurações (`America/Sao_Paulo`), não o offset (`-03`).
- **D11:** no modal de criar, monitores pausados ou com janela viva aparecem desabilitados ("pausado" / "já agendada"), pois não entram numa janela.
- **D12:** janela iniciada tem o horário travado e só o motivo muda. "Iniciada" vale pelo relógio do servidor (`starts_at <= agora`), não pelo tick. Para parar antes, desligar.
- **D13:** editar não troca o monitor; muda só `starts_at`, `ends_at` e `reason`, revalidando D1 e D4. Para outro monitor: desligar e agendar de novo.
- **D14:** a página lista só janelas vivas (agendadas e em curso), por início. Sem histórico: cancelar antes do início apaga a linha.
- **D15:** o modal serve a criar e a editar (monitor fixo; horário desabilitado em janela iniciada). Se a janela começar com o modal aberto e o horário alterado, o servidor responde `409 MAINTENANCE_LOCKED`.
- **D16:** o plano previa a Etapa 3 em branch e PR empilhado; foi aplicada no mesmo PR (desvio 8).
- **D17:** "Editar" também no card de manutenção do detalhe do monitor, com o mesmo modal.

### Desvios do plano

1. **O bloqueio vai até o fim ser aplicado, não até `ends_at`.** O plano previa bloqueio exato em `ends_at`. Mas o fim (estado `UNKNOWN`, `resumed_at`) só roda no tick, até 10 s depois. Nesse intervalo um check com falha abriria um incidente, e o fim, ao zerar o estado para `UNKNOWN`, deixaria esse incidente órfão. Custo: até 1 tick a mais de manutenção.
2. **`MonitorService.Ping` não consulta a janela.** O ping passa por `Record` → `withMonitor`, que já tem o gate e devolve `nil` (200 sem registrar). Uma segunda checagem no service seria a mesma regra duas vezes.
3. **`Start`/`Finalize` são métodos transacionais do repositório**, no lugar de `MarkStarted`/`MarkFinalized` soltos. Assim fechar o incidente, zerar o estado e marcar a janela ficam na mesma transação. `Finalize` aplica o início antes se ele não rodou (janela que começou e terminou com o Pingou fora).
4. **Contagem desde a retomada usa `>=`**, não `>`: a precisão é de segundos e o primeiro check após retomar roda no mesmo segundo.
5. **Cancelar no segundo exato do início** grava `ends_at = starts_at + 1 s` (o `CHECK (ends_at > starts_at)` recusaria igual).
6. `MaintenanceModal` ficou em arquivo próprio (o plano o colocava em `MonitorModals.jsx`), e a validação do formulário em `lib/maintenance.js`, para ter teste.
7. No desktop, o selo e o motivo do incidente fechado foram para a coluna "Motivo da queda": "Fechado por manutenção" não cabia na coluna de duração (100 px).
8. **Etapa 3 no mesmo PR das etapas 1 e 2**, e não num PR empilhado (D16), por decisão do autor na hora de aplicar.
9. O card do detalhe trocou "Encerrar agora" por **"Desligar"**, o mesmo nome da página, e as duas telas usam o mesmo diálogo (`CancelWindowDialog`).
10. `PATCH` com `monitor_id` no corpo responde `400` (em vez de ignorar o campo), para o erro ficar visível.
11. A listagem da página varre `monitors` (no máximo 100 linhas); o teste de plano de consulta registra essa exceção.

### Limitações conhecidas

- Editar o motivo de uma janela em curso não muda o `close_note` dos incidentes que ela já fechou (cópia de propósito).
- Duas abas editando a mesma janela: a última gravação vence, sem aviso.
- D11 bloqueia monitores pausados só na tela: a API ainda aceita agendar janela para monitor pausado (a pausa vence a janela, D3).
- Efeitos de início e fim atrasam até 10 s (tick). O bloqueio de checks começa na hora.
- Logo depois do `ends_at`, até o tick rodar, criar outra janela para o mesmo monitor dá `409` (a antiga ainda está viva).
- No celular, agendar fica no menu ⋮ do detalhe; o botão "Agendar manutenção" do topo é só do desktop.
- Baldes do `/stats` em fusos de meia hora seguem pendentes (fora do escopo).

### Descartado

Flag `is_in_maintenance` em `monitors`; um timer por janela (`time.AfterFunc`); coluna `canceled_at`; biblioteca de date picker; lote parcial; coluna nova `close_reason TEXT`; janela recorrente; aviso de "começou/terminou"; histórico de janelas encerradas (exigiria migração); estender ou encurtar o horário de janela em curso (reaplicaria os efeitos do tick); "editar" como cancelar e criar outra (perde o ID e pode falhar no meio); montar a página com `GET /monitors` e filtrar `maintenance` (payload grande, sem ordem por início).

### Validação

`go vet`, `go test -race ./...` (scheduler 5× seguidas, sem falha intermitente) e `npm run lint`/`format:check`/`test`/`build` verdes. QA manual com o Pingou real, fuso `America/Sao_Paulo` e navegador em UTC: pausa com motivo; janela de 3 min em 2 monitores (um DOWN com incidente aberto) gravada no UTC certo; conflito `409` com o nome do monitor; incidente fechado por manutenção com o motivo; nenhum check em 25 s de janela; ping `200` sem registro; fim natural com check imediato e prazo do heartbeat recomeçado; reinício no meio da janela mantendo a manutenção; "Encerrar agora" pelo celular com check imediato. Etapa 3: página vazia; duas janelas ordenadas; edição de horário e motivo com o UTC certo; janela de 2 min virando "Em curso" sozinha; editar o horário com o início passando com o modal aberto → mensagem do 409 e horário travado, e salvar só o motivo em seguida; nenhum incidente novo ao editar o motivo em curso; editar pelo card do detalhe; desligar pela página com check imediato; `curl` do `PATCH` (409 com horário, 200 só com motivo); página no celular com a aba nova.

## 15. Modo escuro (Claro, Escuro, Sistema)

Plano `tasks/plano-dark-mode.md`. Tema escuro no painel, escolhido em Configurações > Aparência. Só frontend (`web/`); nenhuma mudança no Go nem na API.

### Entregue

- **Opções:** Claro, Escuro e **Sistema** (padrão; segue `prefers-color-scheme` e troca ao vivo quando o sistema muda). A escolha fica no `localStorage` (`pingou_theme`), por navegador, e **não** passa pelo "Salvar" das Configurações. Valor inválido ou `localStorage` bloqueado vira "Sistema".
- **Sem atalho no cabeçalho do painel.** O seletor fica só em Configurações (decisão de 2026-11-04).
- **Mecanismo:** atributo `data-theme="light|dark"` no `<html>`. Um script inline no `<head>` do `index.html` aplica o tema **antes da primeira pintura** (sem piscar). `lib/theme.js` tem a mesma regra e um teste confere que os dois concordam. `ThemeProvider` (na raiz, inclusive no Login) acompanha o sistema e trocas feitas em outra aba.
- **Tokens de papel:** os componentes usam `--pg-cor-*` (`bg-pg-superficie`, `text-pg-texto`, `border-pg-linha`, `bg-pg-status-down-bg`...), definidos em `index.css` para os dois temas e mapeados em `tailwind.config.js`. As cores cruas (`--pg-ardosia`, `--pg-up`...) continuam constantes. Motivo: a mesma cor crua serve a papéis diferentes (`--pg-ardosia` é texto e também fundo de tooltip).
- **Trava:** `scripts/check-tokens.sh` falha em cor crua, hexadecimal ou `rgba()` em `components/` e `pages/` (`npm run check:tokens`, também no `make build-web`).
- **Logos:** `ThemedImage` mostra marca negativa / símbolo invertido só por CSS (`dark:hidden`). **Toasts** (react-toastify) têm variáveis próprias no escuro.
- **Paleta:** `DESIGN.md` §3.6 e §3.6.6 (papéis aprovados pelo dono do projeto em 2026-11-05). Cor nova: `--pg-down-borda-escuro` (`#684953`); `--pg-down-borda` (`#EDC4C8`) virou token no claro.

### Decisões

1. **Tokens de papel, e não `dark:` em cada classe nem redefinir as cores cruas.** `dark:` dobraria ~900 classes; redefinir `--pg-ardosia` quebraria os fundos escuros fixos.
2. **Trava de segurança (`THEME_ENABLED` e `var ON`)** usada durante a migração, para o padrão "Sistema" não ligar o escuro com telas pela metade. Foi ligada no mesmo commit do seletor. Pode ser removida; ela segue coberta por teste.
3. **Sem `<meta name="color-scheme">`:** o script já define `color-scheme` no `<html>`.
4. **Página pública de status (ideia 3 do `IDEAS.md`)** reaproveita tudo isso e terá botão de tema no cabeçalho.

### Limitações conhecidas

- Com CSP no futuro, o script inline precisa de hash ou nonce.
- No escuro o modal se destaca pouco do véu (1,1:1); a borda `--pg-ardosia-linha` faz a separação.
- As barras de histórico (UP/DOWN) e o gráfico dependem de cor para o dado; o estado em si sempre tem texto e ícone (selo). Em acromatopsia, verde e vermelho claros ficam parecidos.
- No tema **claro** (já existia): o texto "N fora do ar" do banner do Painel usa `--pg-down` sobre `--pg-down-bg` (4,18:1, abaixo de AA). O `DESIGN.md` §4 já registra esse par.
- O estilo global de `input` (`globalStyles.js`) usa `--pg-borda`, e o `DESIGN.md` manda `--pg-neblina` na borda de input. Mantido como estava.

### Descartado

Atalho de tema no cabeçalho do painel; guardar o tema no backend (`settings`); aplicar o tema só no React (piscaria); temas personalizados; tema por monitor ou por página.

### Validação

- **Tema claro idêntico:** 29 telas (desktop e celular, modais, menus, 4 níveis de certificado), antes e depois, com respostas da API gravadas e hora congelada. Únicas diferenças: porta na URL de ping (origem da captura) e a altura das Configurações (bloco novo).
- **axe-core** (WCAG 2 A/AA): **0 violações** no escuro nas 29 telas. No claro, 1 (a do banner, acima, anterior a esta entrega).
- **Seletor no navegador (16 verificações):** primeira visita em "Sistema"; troca ao vivo com o sistema; Escuro/Claro aplicam na hora e persistem; recarregar mantém; navegar mantém; valor inválido volta a "Sistema"; teclado; outra aba acompanha; nenhum atalho no cabeçalho (desktop e celular).
- **Sem piscar:** primeira pintura medida com o JavaScript do app atrasado em 1,5 s; já sai no tema certo nos 3 casos testados.
- **Daltonismo** (deuteranopia, protanopia, acromatopsia no Chromium): estados legíveis por selo (texto e ícone) e texto do banner.
- `npm run lint`, `format:check`, `test` (49), `check:tokens` e `build` verdes.

---

## 16. Retry de notificações (outbox)

Plano `plan-notification-retry.md` (ideia 5 do `IDEAS.md`). Antes, o aviso saía uma vez, depois do commit: canal fora do ar ou crash entre o commit e o envio perdia o aviso, e o webhook registrava "sent" até com 500.

### Entregue

- **Outbox:** tabela `notification_outbox` (migration `0010`), gravada na **mesma transação** do `StateMachine.Process`, uma linha por canal. Rollback não deixa aviso; commit sempre deixa.
- **Worker** (`OutboxWorker`): lê os vencidos em lotes de 20 (conexão solta antes de enviar), envia, marca `sent` ou agenda a próxima. Acorda depois do commit (`Wake`) e a cada 30 s. Ao subir, retoma os pendentes.
- **Esperas:** 30 s, 2 min, 10 min, 30 min. **Prazos:** queda e recuperação 60 min do evento; SSL 24 h. Se a próxima tentativa passaria do prazo, expira na hora. "Expirado" só no log.
- **Canais:** `Channel.Send(ctx, OutboxEvent) error` substitui `Notifier`. Webhook com status fora de 2xx agora é falha. Texto de e-mail e ntfy montado no envio a partir do snapshot (compatível com i18n).
- **Webhook:** campo novo `id` (UUID, igual em todas as tentativas). `timestamp` passa a ser a hora do evento; no `cert_expiring` antes era a hora do envio.
- **Retenção:** linhas `sent`/`expired` com mais de 7 dias saem na rotina horária.
- **Desligamento:** etapa "notifications" espera o worker; pendentes ficam no banco.
- `MultiNotifier` removido (o fan-out e o `WaitGroup` viraram o worker).

### Decisões

1. **Sem estado `sending`.** Um único processo por banco (lock de arquivo) e um único worker: crash nunca deixa linha presa. Custo: envio sequencial (um webhook lento atrasa o e-mail da fila).
2. **Snapshot do evento, não só `monitor_id`:** o monitor pode ser editado ou apagado antes do envio. O **destino** é lido na hora.
3. **Sem segredo na fila:** nada de URL do webhook, destinatários, tópico, token do ntfy ou `ping_token`. `target` do monitor pode ter token na query; já é assim no banco e no payload do webhook.
4. **Linha com payload ilegível** vira `expired` (senão ocuparia o lote para sempre).
5. **`event_type`:** `down`, `up`, `cert_expiring` (os nomes que o webhook usa hoje). Se a ideia 7 (nomes de evento) entrar depois, a troca é nas constantes e as linhas antigas vivem no máximo 24 h.

### Limitações conhecidas

- Entrega "pelo menos uma vez": duplicata é possível (por isso o `id`).
- Webhook que responde 4xx por URL errada gasta as tentativas até o prazo de 60 min (o log mostra o status).
- `last_error` do e-mail pode trazer texto do servidor SMTP (até 200 caracteres).
- Ideia 1 (backup): `notification_outbox` deve ficar fora do export.

### Descartado

Reenvio com `time.AfterFunc` em memória; gravar o texto pronto do e-mail/ntfy; uma linha por evento com o estado de cada canal em JSON; estado `sending` com reclamação de linha; fila externa; retry infinito.

### Validação

`go vet` e `go test -race ./...` verdes. Testes novos: migration (sobe/desce, índices parciais), plano de consulta usa `idx_outbox_due` e `idx_outbox_done`, repositório, worker com relógio falso (esperas, prazos de 60 min e 24 h, só o canal com falha repete, canal desconhecido, pânico, aviso vencido, reinício), webhook (2xx/não-2xx, erro sem URL, `id`, `timestamp`), state machine (uma linha por canal com o mesmo `event_id`, rollback, nada sem mudança de estado) e concorrência entre `Process` e worker. Mutação conferida: status não-2xx como sucesso e esperas diferentes derrubam os testes. **Não rodados:** carga com 1 000 linhas pendentes e QA manual com o binário real (`kill -9`, receptor instável); estão em "Pendências abertas".

---

## Pendências abertas

- **Retry de notificações (item 16):** medir a gravação de checks com 1 000 linhas `pending` (receptor fora do ar) e rodar o QA manual com o binário real: receptor com 500 por ~3 min, `kill -9` logo após a queda, troca da URL do webhook durante as tentativas.

- **Migrations 0003 e 0004** foram testadas em bancos de teste e de QA, mas **não** rodadas no `pingou.db` real. Fazer backup antes e conferir `failure_count` / `last_error`.
- **Configurações (T12)**: webhook.site real, destino que não responde (5 s), navegação só por teclado, nome do arquivo baixado no navegador.
- **Mobile**: teste em dispositivo real (teclado virtual com a barra "Alterações não salvas", safe-area no iPhone, leitor de tela).
- **Observação de infraestrutura**: o incidente do monitor `c10` vinha de certificado de outro domínio servido em `api.c10.site` (provável proxy roteando errado). Não é bug do Pingou.
- **E-mail (T14)**: QA manual com Mailpit e provedor real (587/465); o ambiente de desenvolvimento não tinha Docker nem SMTP real.
- **Eixo Y**: verificação manual (T5) dos 4 ranges nos 3 tipos, tema claro/escuro e mobile; conferir em um `/stats` real que `checks`/`failures` do heartbeat batem com pings e prazos perdidos.
- **Pausar fecha o incidente**: QA manual (T7) no navegador (HTTP/TCP e heartbeat: incidente sai de "Abertos", aparece em "Resolvidos" com selo, sem e-mail/ntfy; retomar reabre) e `npm run lint`/`build` do frontend. Decidir se vale limpar incidentes órfãos antigos de monitores já pausados.
- **Modo escuro**: conferir em dispositivo real (iPhone/Android com tema do sistema trocando de dia para noite) e com leitor de tela no seletor "Aparência". Decidir se o texto "N fora do ar" do banner (tema claro) deve usar `--pg-ardosia` para passar em AA.
- **Manutenção**: QA em dispositivo real e com leitor de tela no modal de agendar (multi-seleção, erros com foco); janela atravessando a troca de horário de verão num fuso que ainda a tem.
