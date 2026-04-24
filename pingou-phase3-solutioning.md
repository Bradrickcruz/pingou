# Pingou — Fase 3: Solutioning (v1.0)

> **Status:** Aprovado  
> **Escopo:** Design detalhado antes da implementação — schema, contratos de API, estrutura de UI.

---

## 3.1 — Modelo de Dados (Schema SQL)

### Filosofia de design

- **IDs:** UUID v7 (ordenável por tempo) para entidades de produto; INTEGER autoincrement para tabela de alto volume (`checks`).
- **Timestamps:** `TEXT` ISO 8601 UTC (`2026-04-23T14:30:00Z`). Ordenação lexicográfica = ordenação temporal.
- **Booleans:** `INTEGER` com `CHECK (x IN (0,1))`. SQLite não tem BOOLEAN nativo.
- **Soft delete:** não usado no MVP. Delete é delete.
- **Foreign keys:** sempre `ON DELETE CASCADE` para evitar órfãos.

### Tabela `monitors`

| Campo               | Tipo    | Constraints                                                   | Notas                           |
| ------------------- | ------- | ------------------------------------------------------------- | ------------------------------- |
| `id`                | TEXT    | PK                                                            | UUID v7                         |
| `name`              | TEXT    | NOT NULL, 1–100 chars                                         | Label legível                   |
| `url`               | TEXT    | NOT NULL, 1–2048 chars                                        | URL a ser checada               |
| `interval_seconds`  | INTEGER | NOT NULL, CHECK 10–86400                                      | Frequência de check             |
| `timeout_seconds`   | INTEGER | NOT NULL, DEFAULT 10, CHECK 1–60                              | Timeout HTTP                    |
| `failure_threshold` | INTEGER | NOT NULL, DEFAULT 3, CHECK 1–10                               | N falhas consecutivas para DOWN |
| `enabled`           | INTEGER | NOT NULL, DEFAULT 1, CHECK (0,1)                              | Liga/desliga scheduler          |
| `current_state`     | TEXT    | NOT NULL, DEFAULT 'UNKNOWN', CHECK IN ('UNKNOWN','UP','DOWN') | Cache do state machine          |
| `last_checked_at`   | TEXT    | NULL                                                          | ISO 8601 UTC                    |
| `created_at`        | TEXT    | NOT NULL                                                      | ISO 8601 UTC                    |
| `updated_at`        | TEXT    | NOT NULL                                                      | ISO 8601 UTC                    |

**Índices:**

- `idx_monitors_enabled` em `(enabled)` — scheduler filtra `enabled = 1`
- `idx_monitors_current_state` em `(current_state)` — dashboard UI

### Tabela `checks`

| Campo           | Tipo    | Constraints                         | Notas                                 |
| --------------- | ------- | ----------------------------------- | ------------------------------------- |
| `id`            | INTEGER | PK AUTOINCREMENT                    | Alto volume: economiza espaço vs UUID |
| `monitor_id`    | TEXT    | NOT NULL, FK → monitors(id) CASCADE |                                       |
| `success`       | INTEGER | NOT NULL, CHECK (0,1)               | Boolean                               |
| `status_code`   | INTEGER | NULL                                | Null se request não chegou            |
| `latency_ms`    | INTEGER | NOT NULL                            | Tempo até resposta ou timeout         |
| `error_message` | TEXT    | NULL                                | Preenchido em falhas                  |
| `checked_at`    | TEXT    | NOT NULL                            | ISO 8601 UTC                          |

**Índices:**

- `idx_checks_monitor_checked` em `(monitor_id, checked_at DESC)` — histórico por monitor
- `idx_checks_checked_at` em `(checked_at)` — retention cleanup

### Tabela `incidents`

| Campo               | Tipo    | Constraints                         | Notas                      |
| ------------------- | ------- | ----------------------------------- | -------------------------- |
| `id`                | TEXT    | PK                                  | UUID v7                    |
| `monitor_id`        | TEXT    | NOT NULL, FK → monitors(id) CASCADE |                            |
| `started_at`        | TEXT    | NOT NULL                            | Quando virou DOWN          |
| `ended_at`          | TEXT    | NULL                                | Null = incident aberto     |
| `last_error`        | TEXT    | NULL                                | Mensagem de erro no início |
| `notification_sent` | INTEGER | NOT NULL, DEFAULT 0, CHECK (0,1)    | Webhook enviado?           |

**Índices:**

- `idx_incidents_monitor_ended` em `(monitor_id, ended_at)`
- Regra de negócio: máximo 1 incident aberto (`ended_at IS NULL`) por monitor — enforced no código.

### Tabela `settings` (key-value)

| Campo        | Tipo | Constraints |
| ------------ | ---- | ----------- |
| `key`        | TEXT | PK          |
| `value`      | TEXT | NOT NULL    |
| `updated_at` | TEXT | NOT NULL    |

**Seeds:**

- `webhook_url` = `""` (vazio = desabilitado)
- `retention_days` = `"30"`

### Relacionamentos

```
monitors 1 ─── N checks      (CASCADE DELETE)
monitors 1 ─── N incidents   (CASCADE DELETE)
settings  (standalone)
```

### Migration `0001_init.sql`

```sql
-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE monitors (
  id                TEXT PRIMARY KEY,
  name              TEXT NOT NULL CHECK(length(name) BETWEEN 1 AND 100),
  url               TEXT NOT NULL CHECK(length(url) BETWEEN 1 AND 2048),
  interval_seconds  INTEGER NOT NULL CHECK(interval_seconds BETWEEN 10 AND 86400),
  timeout_seconds   INTEGER NOT NULL DEFAULT 10 CHECK(timeout_seconds BETWEEN 1 AND 60),
  failure_threshold INTEGER NOT NULL DEFAULT 3 CHECK(failure_threshold BETWEEN 1 AND 10),
  enabled           INTEGER NOT NULL DEFAULT 1 CHECK(enabled IN (0,1)),
  current_state     TEXT NOT NULL DEFAULT 'UNKNOWN' CHECK(current_state IN ('UNKNOWN','UP','DOWN')),
  last_checked_at   TEXT,
  created_at        TEXT NOT NULL,
  updated_at        TEXT NOT NULL
);
CREATE INDEX idx_monitors_enabled ON monitors(enabled);
CREATE INDEX idx_monitors_current_state ON monitors(current_state);

CREATE TABLE checks (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  monitor_id    TEXT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  success       INTEGER NOT NULL CHECK(success IN (0,1)),
  status_code   INTEGER,
  latency_ms    INTEGER NOT NULL,
  error_message TEXT,
  checked_at    TEXT NOT NULL
);
CREATE INDEX idx_checks_monitor_checked ON checks(monitor_id, checked_at DESC);
CREATE INDEX idx_checks_checked_at ON checks(checked_at);

CREATE TABLE incidents (
  id                 TEXT PRIMARY KEY,
  monitor_id         TEXT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  started_at         TEXT NOT NULL,
  ended_at           TEXT,
  last_error         TEXT,
  notification_sent  INTEGER NOT NULL DEFAULT 0 CHECK(notification_sent IN (0,1))
);
CREATE INDEX idx_incidents_monitor_ended ON incidents(monitor_id, ended_at);

CREATE TABLE settings (
  key        TEXT PRIMARY KEY,
  value      TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

INSERT INTO settings (key, value, updated_at) VALUES
  ('webhook_url',    '',   strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  ('retention_days', '30', strftime('%Y-%m-%dT%H:%M:%fZ','now'));

-- +goose Down
DROP TABLE settings;
DROP TABLE incidents;
DROP TABLE checks;
DROP TABLE monitors;
```

---

## 3.2 — Contratos da API REST

### Convenções globais

- **Base URL:** `/api`
- **Auth:** header `X-API-Key: <key>` em todas as rotas exceto `/healthz`
- **Content-Type:** `application/json`
- **Timestamps:** ISO 8601 UTC
- **Error format:** `{ "error": "mensagem", "code": "SLUG_CODE" }`

### HTTP status codes

| Código           | Uso                                           |
| ---------------- | --------------------------------------------- |
| 200 OK           | Sucesso com body                              |
| 201 Created      | POST de sucesso                               |
| 204 No Content   | DELETE de sucesso                             |
| 400 Bad Request  | Validação                                     |
| 401 Unauthorized | API key inválida/ausente                      |
| 404 Not Found    | Recurso inexistente                           |
| 409 Conflict     | Regra de negócio (ex: limite de 100 monitors) |
| 500 Internal     | Erro inesperado                               |

### Endpoints

#### Health

| Método | Path       | Auth | Output                              |
| ------ | ---------- | ---- | ----------------------------------- |
| GET    | `/healthz` | Não  | `{"status":"ok","version":"0.1.0"}` |

#### Monitors

| Método | Path                           | Input                                        | Output                         | Status |
| ------ | ------------------------------ | -------------------------------------------- | ------------------------------ | ------ |
| POST   | `/api/monitors`                | `CreateMonitorDTO`                           | `Monitor`                      | 201    |
| GET    | `/api/monitors`                | Query: `enabled`, `state`, `limit`, `offset` | `{data, total, limit, offset}` | 200    |
| GET    | `/api/monitors/{id}`           | —                                            | `Monitor`                      | 200    |
| PATCH  | `/api/monitors/{id}`           | `UpdateMonitorDTO`                           | `Monitor`                      | 200    |
| DELETE | `/api/monitors/{id}`           | —                                            | —                              | 204    |
| GET    | `/api/monitors/{id}/checks`    | Query: `limit`, `offset`                     | Paginado                       | 200    |
| GET    | `/api/monitors/{id}/incidents` | Query: `open`, `limit`, `offset`             | Paginado                       | 200    |

#### Incidents

| Método | Path             | Input                            |
| ------ | ---------------- | -------------------------------- |
| GET    | `/api/incidents` | Query: `open`, `limit`, `offset` |

#### Settings

| Método | Path            | Input               | Output                          |
| ------ | --------------- | ------------------- | ------------------------------- |
| GET    | `/api/settings` | —                   | `{webhook_url, retention_days}` |
| PUT    | `/api/settings` | `UpdateSettingsDTO` | Settings atualizado             |

#### Export

| Método | Path          | Output                                  |
| ------ | ------------- | --------------------------------------- |
| GET    | `/api/export` | `application/octet-stream` — dump `.db` |

### DTOs

**CreateMonitorDTO**

```json
{
  "name": "API de Produção",
  "url": "https://api.exemplo.com/health",
  "interval_seconds": 60,
  "timeout_seconds": 10,
  "failure_threshold": 3,
  "enabled": true
}
```

**UpdateMonitorDTO** (todos os campos opcionais)

```json
{
  "name": "Novo nome",
  "interval_seconds": 120,
  "enabled": false
}
```

**Monitor (response)**

```json
{
  "id": "01HXY...",
  "name": "API de Produção",
  "url": "https://api.exemplo.com/health",
  "interval_seconds": 60,
  "timeout_seconds": 10,
  "failure_threshold": 3,
  "enabled": true,
  "current_state": "UP",
  "last_checked_at": "2026-04-23T14:30:00Z",
  "created_at": "2026-04-23T10:00:00Z",
  "updated_at": "2026-04-23T14:30:00Z"
}
```

**Check (response)**

```json
{
  "id": 12345,
  "monitor_id": "01HXY...",
  "success": true,
  "status_code": 200,
  "latency_ms": 147,
  "error_message": null,
  "checked_at": "2026-04-23T14:30:00Z"
}
```

**Incident (response)**

```json
{
  "id": "01HXZ...",
  "monitor_id": "01HXY...",
  "started_at": "2026-04-23T14:25:00Z",
  "ended_at": null,
  "last_error": "dial tcp: timeout",
  "notification_sent": true,
  "duration_seconds": 300
}
```

**UpdateSettingsDTO**

```json
{
  "webhook_url": "https://discord.com/api/webhooks/...",
  "retention_days": 30
}
```

### Webhook Payload

**Evento `down`**

```json
{
  "event": "down",
  "monitor": {
    "id": "01HXY...",
    "name": "API de Produção",
    "url": "https://api.exemplo.com/health"
  },
  "timestamp": "2026-04-23T14:25:00Z",
  "last_error": "dial tcp: timeout",
  "downtime_duration_seconds": null
}
```

**Evento `up`**

```json
{
  "event": "up",
  "monitor": { "id": "...", "name": "...", "url": "..." },
  "timestamp": "2026-04-23T14:30:00Z",
  "last_error": null,
  "downtime_duration_seconds": 300
}
```

### Exemplos de erro

```json
{ "error": "interval_seconds must be between 10 and 86400", "code": "VALIDATION_ERROR" }
{ "error": "monitor limit reached (max 100)", "code": "LIMIT_REACHED" }
{ "error": "invalid API key", "code": "UNAUTHORIZED" }
{ "error": "monitor not found", "code": "NOT_FOUND" }
```

---

## 3.3 — Estrutura de Componentes React

### Decisões de UX

- **Dashboard-centric:** rota `/` concentra list + create + edit + detail via **modais**. Deep-link via query string (`/?monitor=<id>&action=edit`).
- **Dark mode default** com toggle persistido em `localStorage` (chave `pingou.theme`).
- **Tema customizado exclusivo** centralizado em `src/theme/` — design tokens + Tailwind config estendido.

### Rotas (React Router)

| Path         | Página          | Descrição                              |
| ------------ | --------------- | -------------------------------------- |
| `/login`     | `LoginPage`     | Input da API key (localStorage)        |
| `/`          | `DashboardPage` | Lista + modais CRUD/detail de monitors |
| `/incidents` | `IncidentsPage` | Histórico global de incidentes         |
| `/settings`  | `SettingsPage`  | Webhook URL + retention + toggle tema  |
| `/export`    | `ExportPage`    | Download do dump `.db`                 |

### Árvore de pastas

```
src/
├── api/
│   ├── client.ts              # fetch wrapper + X-API-Key + ApiError
│   ├── monitors.ts
│   ├── checks.ts
│   ├── incidents.ts
│   ├── settings.ts
│   └── types.ts
│
├── theme/
│   ├── tokens.ts              # cores, espaçamentos, radius, shadows
│   ├── ThemeProvider.tsx      # contexto light/dark + toggle
│   ├── useTheme.ts            # hook consumidor
│   └── index.ts               # barrel export
│
├── components/
│   ├── ui/                    # primitivos (usam tokens do tema)
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Select.tsx
│   │   ├── Badge.tsx
│   │   ├── Card.tsx
│   │   ├── Modal.tsx          # base reutilizável
│   │   ├── Pagination.tsx
│   │   ├── Toast.tsx
│   │   └── ThemeToggle.tsx
│   │
│   ├── layout/
│   │   ├── AppShell.tsx
│   │   ├── Header.tsx         # logo + tema toggle + logout
│   │   └── Sidebar.tsx
│   │
│   ├── monitors/
│   │   ├── MonitorCard.tsx          # card no dashboard
│   │   ├── MonitorStateBadge.tsx    # UP/DOWN/UNKNOWN
│   │   ├── MonitorForm.tsx          # form reutilizável (create/edit)
│   │   ├── MonitorFormModal.tsx     # wrapper modal p/ create/edit
│   │   ├── MonitorDetailModal.tsx   # detail + histórico + incidents
│   │   ├── MonitorDeleteDialog.tsx  # confirm dialog
│   │   └── ChecksTable.tsx
│   │
│   └── incidents/
│       ├── IncidentsList.tsx
│       └── IncidentRow.tsx
│
├── pages/
│   ├── LoginPage.tsx
│   ├── DashboardPage.tsx      # orquestra modais via query params
│   ├── IncidentsPage.tsx
│   ├── SettingsPage.tsx
│   └── ExportPage.tsx
│
├── hooks/
│   ├── useAuth.ts
│   ├── useMonitors.ts
│   ├── useMonitor.ts
│   ├── useChecks.ts
│   ├── useIncidents.ts
│   ├── useSettings.ts
│   └── useModalRoute.ts       # sincroniza modal ↔ query string
│
├── lib/
│   ├── format.ts
│   └── constants.ts
│
├── App.tsx
└── main.tsx
```

### Sistema de tema

**`src/theme/tokens.ts`**

```ts
export const tokens = {
  colors: {
    dark: {
      bg: "#0b0f14",
      surface: "#121821",
      surfaceAlt: "#1a2230",
      border: "#263041",
      text: "#e6edf3",
      textMuted: "#8b95a5",
      primary: "#4f8cff",
      primaryHover: "#6aa0ff",
      success: "#2ecc71",
      danger: "#ff5c5c",
      warning: "#f5a623",
    },
    light: {
      bg: "#f7f9fc",
      surface: "#ffffff",
      surfaceAlt: "#f1f4f9",
      border: "#dde3ec",
      text: "#0b0f14",
      textMuted: "#5b6676",
      primary: "#2563eb",
      primaryHover: "#1d4ed8",
      success: "#18a957",
      danger: "#e03131",
      warning: "#d97706",
    },
  },
  radius: { sm: "4px", md: "8px", lg: "12px", full: "9999px" },
  space: { 1: "4px", 2: "8px", 3: "12px", 4: "16px", 6: "24px", 8: "32px" },
  shadow: {
    sm: "0 1px 2px rgba(0,0,0,0.1)",
    md: "0 4px 12px rgba(0,0,0,0.15)",
  },
  font: {
    family: "'Inter', ui-sans-serif, system-ui, sans-serif",
  },
} as const;

export type ThemeMode = "dark" | "light";
```

**`src/theme/ThemeProvider.tsx`**

- Contexto React que expõe `{ mode, toggle, tokens }`.
- Aplica `<html data-theme="dark">` para Tailwind estratégia `class` ou CSS vars.
- Default: `dark`. Persistência em `localStorage["pingou.theme"]`.

**Integração Tailwind:** `tailwind.config.js` estende `theme.extend.colors` consumindo CSS variables (`var(--color-primary)`) definidas em `:root[data-theme="dark"]` e `:root[data-theme="light"]`. Isso permite trocar tema sem re-render completo.

### Padrão de modais no Dashboard

**Query string como fonte de verdade:**

- `/` — lista
- `/?action=create` — modal de criação
- `/?monitor=<id>` — modal de detalhe
- `/?monitor=<id>&action=edit` — modal de edição
- `/?monitor=<id>&action=delete` — dialog de confirmação

**Hook `useModalRoute`:**

```ts
const { action, monitorId, open, close } = useModalRoute();
// open({ action: 'edit', monitorId: '01HXY...' })
// close()
```

**Vantagens:** deep-link compartilhável, back-button do browser funciona, refresh preserva estado.

### Componentes-chave

**`MonitorCard`** — `{ monitor, onClick }` → click abre `MonitorDetailModal`.

**`MonitorStateBadge`** — `{ state: 'UP'|'DOWN'|'UNKNOWN' }` → cor via tokens do tema.

**`MonitorForm`** — `{ initialValue?, onSubmit, submitLabel }` — validações espelham backend (interval 10–86400, threshold 1–10, URL válida).

**`MonitorFormModal`** — wrapper que consome `useModalRoute` e renderiza `MonitorForm` em `create` ou `edit`.

**`MonitorDetailModal`** — abas: **Overview** (campos + estado atual), **Histórico** (`ChecksTable`), **Incidentes** (`IncidentsList`).

**`useAuth`** — `{ apiKey, setApiKey, logout, isAuthenticated }`. Guarda em `localStorage["pingou.apiKey"]`, redireciona pra `/login` quando ausente.

**`client.ts`** — base URL `/api`, injeta `X-API-Key`, parse de erro → `throw new ApiError(message, code, status)`.

### Estado da aplicação

Sem Redux/Zustand/TanStack Query no MVP. `useState` + refetch manual após mutations. Revisitar se surgir polling em tempo real ou dados compartilhados entre 3+ lugares.

---

## 📐 Decisões Arquiteturais (consolidadas nesta fase)

| ID      | Decisão                                                                          | Rationale                                                              |
| ------- | -------------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| **D9**  | IDs mistos: UUID v7 para monitors/incidents, INTEGER autoincrement para `checks` | Otimiza espaço/performance em tabela de alto volume (~13M rows em 90d) |
| **D10** | `monitors.current_state` materializado (cache)                                   | O(1) vs query derivada; atualizado em transação                        |
| **D11** | `settings` como key-value                                                        | Evita migration por config; tipagem forte no service layer             |
| **D12** | Sem state manager no MVP                                                         | Complexidade desnecessária p/ 7 telas simples                          |
| **D13** | Dashboard-centric UX com modais + query string                                   | Elimina navegação para tarefas rotineiras; preserva deep-link          |
| **D14** | Tema customizado isolado em `src/theme/` com CSS vars                            | Troca de paleta sem tocar em componentes; dark-default                 |
| **D15** | `monitors.enabled` BOOL em vez de `status` enum                                  | Só há 2 estados; boolean é mais honesto que enum com 2 valores         |

---

## 🔴 Exit Gate da Fase 3

- [x] Schema SQL finalizado
- [x] Contratos da API definidos
- [x] Estrutura React definida
- [x] Tema e padrões de UX documentados
- [ ] **Aprovação final do Bryan para avançar à Fase 4 (Implementação)**

**Próximo passo (após aprovação):** Fase 4 — Fase 1 do plano de implementação (`pingou-mvp.md`): setup do projeto Go, módulos, estrutura de pastas, Makefile, slog.
