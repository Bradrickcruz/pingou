# 🏓 Pingou

> **Rodou, Pingou** — health checker self-hosted, leve e open-source em Go.

Pingou é um monitor de uptime self-hosted que executa checks HTTP periódicos em URLs configuradas, persiste resultados em SQLite, detecta incidentes com base em transições de estado e oferece um dashboard web embutido no próprio binário Go.

O objetivo do projeto é ser:

- **simples de rodar**
- **simples de manter**
- **leve**
- **sem dependências externas obrigatórias**, além do próprio binário ou container

## Visão geral

O Pingou monitora endpoints HTTP em intervalos configuráveis e registra:

- status atual do monitor
- histórico de checks
- incidentes abertos e resolvidos
- configurações globais, como webhook e retenção

A aplicação pode ser distribuída de duas formas:

- **binário único** com frontend embutido via `embed.FS`
- **container Docker** com persistência em volume

## Funcionalidades

- Monitoramento de URLs HTTP
- Intervalo, timeout e threshold configuráveis por monitor
- Detecção de transições de estado:
  - `UNKNOWN → UP`
  - `UNKNOWN → DOWN`
  - `UP → DOWN`
  - `DOWN → UP`
- Registro de incidentes
- Dashboard web embutido no backend
- Tela de login com API Key
- Configuração global de webhook
- Configuração de retenção de checks
- Export do banco SQLite
- Execução self-hosted com SQLite

## Stack atual

### Backend

- Go
- `net/http`
- SQLite
- `github.com/mattn/go-sqlite3`
- `goose` para migrations
- `log/slog`

### Frontend

- React
- Vite
- React Router
- Axios

### Infra

- Docker
- Docker Compose
- `embed.FS` para embutir os assets do frontend no binário

## Estrutura do projeto

```text
pingou/
├── cmd/
│   └── pingou/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── domain/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   └── checker/
├── migrations/
├── web/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── pages/
│   │   └── theme/
│   ├── package.json
│   └── vite.config.js
├── bin/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
├── .editorconfig
├── go.mod
└── README.md
```

## Como rodar em desenvolvimento

### Requisitos

- Go
- Node.js + npm
- GCC ou toolchain compatível com CGO
- SQLite
- `gofumpt` instalado, se você for usar o target `fmt`

## Convenções de formatação

O projeto inclui `.editorconfig` na raiz para padronizar configuracoes basicas entre editores:

- UTF-8, line endings LF e newline final.
- Remocao de espacos finais em editores compativeis.
- Tabs para arquivos Go e Makefile.
- 2 espacos para JS, JSX, JSON, CSS, Markdown, YAML, SQL, Dockerfile e `.env*`.

Essa configuracao complementa `gofumpt` no backend e ESLint/Vite no frontend.

### 1. Configurar ambiente

Crie um arquivo `.env` na raiz do projeto com base no `.env.example`.

Exemplo:

```env
PINGOU_PORT=8080
PINGOU_DATABASE_URL=./pingou.db
PINGOU_API_KEY=dev-api-key
PINGOU_LOG_LEVEL=info
PINGOU_CORS_ALLOWED_ORIGINS=
```

### 2. Rodar a aplicação

```bash
make dev
```

Esse comando:

- formata o código Go
- sobe o backend
- usa as variáveis do `.env`

### 3. Build do frontend

```bash
make build-web
```

### 4. Build completo

```bash
make build
```

O binário será gerado em:

```bash
bin/pingou
```

Ainda não há workflow de GitHub Actions neste momento. Também não há suíte de testes automatizados no primeiro momento; a validação do projeto é local e manual, via `make build` e `make docker-build` quando necessário.

## Como rodar com Docker

O Dockerfile inclui um `HEALTHCHECK` que verifica o status da aplicação a cada 30 segundos via endpoint público `/healthz`.

### Build da imagem

```bash
make docker-build
```

### Subir com Docker Compose

```bash
make docker-up
```

Verifique o status de health do container com:

```bash
docker ps
```

Você deverá ver o status `healthy` na coluna `STATUS` se tudo estiver funcionando.

### Derrubar containers

```bash
make docker-down
```

## Login no dashboard

O dashboard exige autenticação por **API Key**.

Ao abrir a aplicação no navegador, a tela de login solicitará a chave.
Essa chave é validada contra a API e armazenada no `localStorage`.

Use o mesmo valor definido em:

```env
PINGOU_API_KEY=...
```

## Variáveis de ambiente

| Variável                      | Obrigatória | Default     | Descrição                                                                                                          |
| ----------------------------- | ----------: | ----------- | ------------------------------------------------------------------------------------------------------------------ |
| `PINGOU_PORT`                 |         não | `8080`      | Porta HTTP da aplicação                                                                                            |
| `PINGOU_DATABASE_URL`         |         não | `pingou.db` | Caminho do arquivo SQLite                                                                                          |
| `PINGOU_API_KEY`              |         sim | -           | Chave usada para proteger o dashboard e as rotas `/api/*`                                                          |
| `PINGOU_LOG_LEVEL`            |         não | `info`      | Nível de log (`DEBUG` habilita logs debug)                                                                         |
| `PINGOU_CORS_ALLOWED_ORIGINS` |         não | (vazio)     | Lista separada por vírgula de origins permitidas para CORS; vazio = CORS desabilitado. Ex: `http://localhost:5173` |

## SQLite

O backend abre o SQLite com:

- `_foreign_keys=on`
- `_journal_mode=WAL`
- `_busy_timeout=5000`
- `SetMaxOpenConns(1)`

Se `PINGOU_DATABASE_URL` ja incluir query params, eles sao preservados e os parametros operacionais acima sao aplicados pelo app.

## Autenticação da API

Todas as rotas protegidas exigem o header:

```http
X-API-Key: sua-chave
```

### Exemplo com `curl`

```bash
curl -H "X-API-Key: dev-api-key" http://localhost:8080/api/monitors
```

## Middlewares HTTP

### Recover

O servidor inclui um middleware de `recover` que captura `panic` em handlers, loga internamente com `request_id` e stack trace, e responde ao cliente com JSON `500` (mensagem genérica).

Resposta de erro:

```json
{
  "error": "internal server error",
  "code": "INTERNAL_ERROR"
}
```

### CORS

CORS é controlado via variável de ambiente `PINGOU_CORS_ALLOWED_ORIGINS` (lista separada por vírgula). Se vazia (padrão), nenhum cabeçalho CORS é adicionado.

Exemplo para desenvolvimento com Vite:

```env
PINGOU_CORS_ALLOWED_ORIGINS=http://localhost:5173
```

O middleware:

- Responde a preflight `OPTIONS` com cabeçalhos apropriados.
- Permite header `X-API-Key` em requisições cross-origin.
- Valida a origem contra a lista de permitidas.

---

## Endpoints principais

### Público

#### `GET /healthz`

Retorna o status básico da aplicação.

### Protegidos

#### `GET /api/monitors`

Lista os monitores.

#### `POST /api/monitors`

Cria um novo monitor.

#### `GET /api/monitors/:id`

Busca um monitor por ID.

#### `PATCH /api/monitors/:id`

Atualiza um monitor.

#### `DELETE /api/monitors/:id`

Remove um monitor.

#### `GET /api/incidents`

Lista incidentes.

#### `GET /api/settings`

Busca as configurações globais.

#### `PATCH /api/settings`

Atualiza configurações globais.

#### `GET /api/export`

Baixa um dump do banco SQLite.

## Webhook

Quando `webhook_url` estiver configurado em settings, o Pingou envia um `POST` JSON nas transições de estado.

Evento de queda:

```json
{
  "event": "down",
  "monitor": {
    "id": "018f2f7a-...",
    "name": "API",
    "url": "https://api.example.com/health"
  },
  "timestamp": "2026-05-01T12:00:00Z",
  "last_error": "unexpected status code: 500",
  "downtime_duration_seconds": null
}
```

Evento de recuperação:

```json
{
  "event": "up",
  "monitor": {
    "id": "018f2f7a-...",
    "name": "API",
    "url": "https://api.example.com/health"
  },
  "timestamp": "2026-05-01T12:05:00Z",
  "last_error": null,
  "downtime_duration_seconds": 300
}
```

## Comandos disponíveis no Makefile

```bash
make dev
make fmt
make build
make test
make clean
make build-web
make docker-build
make docker-up
make docker-down
make release
```

## Release

Para gerar uma release local:

```bash
make release
```

Esse processo gera o build da aplicação e deixa o binário pronto em:

```bash
bin/pingou
```

## Objetivos do projeto

O Pingou nasceu com alguns objetivos bem claros:

- ser um health checker simples
- rodar com poucos recursos
- usar SQLite para reduzir complexidade operacional
- embutir o frontend no backend
- funcionar bem como projeto real e também como estudo de Go

## Limites e foco do MVP

O projeto foi pensado para manter escopo controlado.
A ideia não é competir com ferramentas enterprise, e sim entregar um monitor funcional, leve e compreensível.

Foco do MVP:

- checks HTTP
- incidentes básicos
- webhook global
- dashboard embutido
- export do banco

## Pós-MVP

Ideias futuras:

- TCP checks
- status page pública
- métricas e gráficos
- integração com Prometheus
- webhook por monitor
- templates para Discord e Slack
- multi-tenancy

## Licença

Apache 2.0
