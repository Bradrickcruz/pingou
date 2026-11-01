# 💡 Ideias de evolução

---

# 🧩 Task: Notificações via ntfy

## 💡 Idea

Receber os alertas de queda e recuperação num servidor ntfy, como mais um canal além do webhook e do e-mail.

## 🎯 Goal

Alertas em tempo real no celular/desktop via ntfy, com mensagem legível (não o JSON cru).

## 📦 Initial Scope

**Etapa 1 — testar sem código (fazer primeiro):**

- O ntfy suporta templates para mensagens recebidas via webhook (procurar "message templating" na documentação do ntfy).
- Montar uma URL de ntfy com template que leia os campos do payload do Pingou (`event`, `monitor.name`, `monitor.target`, `last_error`, `downtime_duration_seconds`, `dashboard_url`).
- Configurar essa URL como `webhook_url` e validar com **Enviar teste** (`POST /api/settings/webhook/test`).
- Testar também com servidor ntfy protegido por token.

**Se funcionar:** a task vira documentação. Adicionar uma seção "ntfy" no README com a URL de exemplo pronta.

**Se não funcionar (etapa 2):** canal próprio, seguindo o padrão do e-mail:

- Servidor, tópico e token **só no ambiente** (`PINGOU_NTFY_*`). Nada no banco nem no export.
- Qualquer servidor ntfy (próprio ou ntfy.sh).
- Autenticação por token.
- Configuração global (como webhook e e-mail).
- Envia na queda e na recuperação.
- Mesma entrega dos outros canais: assíncrona, uma vez por transição, sem retry.
- Botão **Enviar teste** em Configurações, com rate limit igual aos outros (5 por minuto por IP).
- Token de ping do Heartbeat **nunca** vai na mensagem.

## 🚧 Out of Scope

- Outros canais (Telegram, Slack, Discord)
- Regras avançadas de prioridade ou ações do ntfy
- Configuração por monitor

## 🔍 Context

Sistema de notificações existente (webhook + e-mail). Roadmap "Depois" já prevê templates para Discord e Slack.

## ⚠️ Notes

- Se a etapa 1 falhar por limitação do template, avaliar se vale resolver de forma genérica (templates de webhook) em vez de um canal só para o ntfy. Isso adiantaria o item de Discord/Slack do roadmap.

## 🏷️ Type

feature

## ⏱️ Priority

medium (etapa 1 é rápida, fazer antes das outras)

---

# 🧩 Task: Eixo Y dinâmico no gráfico de tempo de resposta

## 💡 Idea

Na página de detalhe do monitor, o eixo Y do gráfico de tempo de resposta se ajusta aos dados do período exibido, em vez de usar uma escala fixa.

## 🎯 Goal

Gráfico legível e proporcional aos dados reais de cada monitor.

## 📦 Initial Scope

**HTTP e TCP:**

- Máximo do eixo Y = maior `avg_latency_ms` da `series` retornada por `/api/monitors/:id/stats` no `range` selecionado (24h / 7d / 14d / 30d).
- Ignorar baldes com `avg_latency_ms: null`.
- Arredondar para cima com "números bonitos" (passos 1, 2, 5 × 10ⁿ). Ex.: 70 → 100, 350 → 400, 2300 → 2500.
- Máximo do eixo nunca abaixo de 100ms.
- Recalcular ao trocar o `range`.

**Heartbeat:**

- Não mostra tempo de resposta. O gráfico mostra **só os pings recebidos** por balde (campo `checks` da `series`).
- Balde com zero pings é **normal** (ex.: tarefa diária). Não pintar como erro.
- Vermelho / indicação de problema vem **só** dos períodos em `down`.

## 🚧 Out of Scope

- Zoom ou seleção de intervalo no gráfico
- Alterar outros gráficos
- Mudanças no endpoint `/stats`

## 🔍 Context

Página de detalhe do monitor (`/monitors/:id`), gráfico de tempo de resposta.

## ⚠️ Notes

- Outliers já chegam suavizados: cada ponto é média de balde e timeouts não entram em `avg_latency_ms`. Por isso o máximo da série é suficiente (não precisa de p95).
- Verificar se a lib de gráfico já tem domínio automático e se o eixo fixo atual foi definido manualmente.
- Se todos os baldes forem `null` (sem checks com sucesso), usar o eixo mínimo de 100ms.

## 🏷️ Type

improvement

## ⏱️ Priority

medium

---

# 🧩 Task: Aproveitar melhor a largura no desktop

## 💡 Idea

No desktop, as páginas têm um limite de largura baixo e sobra espaço vazio à direita.

## 🎯 Goal

Páginas usando bem o espaço horizontal, sem aparência desbalanceada e sem perder legibilidade.

## 📦 Initial Scope

- **Aumentar** o limite de largura das páginas (faixa entre ~1600px e ~1920px, definir validando visualmente). Não remover o limite.
- Ajustar proporções dos elementos de cada página quando necessário.
- Validar principalmente o dashboard (lista de monitores): linhas muito largas dificultam a leitura.

## 🚧 Out of Scope

- Redesign das páginas
- Mudanças no layout mobile

## 🔍 Context

Todas as páginas do frontend, em telas desktop. Seguir `DESIGN.md`.

## ⚠️ Notes

- Testar em 1366px, 1920px e ultrawide (2560px+).
- Registrar o novo limite no `DESIGN.md`.

## 🏷️ Type

improvement

## ⏱️ Priority

medium

---

# 🧩 Task: Busca parcial por nome no dashboard

## 💡 Idea

Campo de busca no dashboard para filtrar monitores por parte do nome. Substitui a ideia de tags.

## 🎯 Goal

Encontrar monitores rápido conforme a lista cresce, com o mínimo de complexidade.

## 📦 Initial Scope

- Campo de busca no topo do dashboard.
- Filtro **só no frontend**: `GET /api/monitors` já traz todos (limite de 100 monitores).
- Busca parcial (contém), ignorando maiúsculas/minúsculas e acentos ("servico" encontra "Serviço").
- Buscar também no alvo: URL (http) e `host:porta` (tcp).
- Guardar o termo na URL (`?q=...`) para sobreviver ao recarregar e ao voltar da página de detalhe.
- Estado vazio claro quando nada for encontrado, com opção de limpar a busca.

## 🚧 Out of Scope

- Tags
- Busca no backend
- Filtros por estado, tipo ou outros campos

## 🔍 Context

Dashboard (lista de monitores).

## ⚠️ Notes

- Tags foram descartadas por enquanto: com no máximo 100 monitores, a busca resolve o problema com muito menos esforço. Reavaliar só se surgir necessidade real de agrupar.

## 🏷️ Type

improvement

## ⏱️ Priority

medium
