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
- **Atomicidade nas operações de state machine** (transações SQL)

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
- Tailwind CSS

### Infra

- Docker
- Docker Compose
- `embed.FS` para embutir os assets do frontend no binário

## Arquitetura de Transações

O Pingou implementa **atomicidade** no fluxo de processamento da state machine usando o padrão **Unit of Work**.

### O problema

Antes desta implementação, `StateMachine.Process` executava operações sequenciais via repositórios:

1. Inserir `checks`
2. Atualizar `monitors` (current_state, last_checked_at)
3. Abrir incidente (se transita para DOWN)
4. Fechar incidente (se transita para UP)

Se uma etapa falhasse depois de outra já ter sido gravada, o banco ficaria em estado parcialmente persistido.

### A solução

A implementação usa **transações SQL** para garantir que todas as operações sejam aplicadas atomicamente:

```
┌─────────────────────────────────────────────────────────────┐
│                     StateMachine.Process                     │
├─────────────────────────────────────────────────────────────┤
│  1. uow.Begin() → inicia transação                          │
│  2. checkRepo.CreateWithTx() → insere check                  │
│  3. monitorRepo.UpdateWithTx() → atualiza monitor            │
│  4. incidentRepo.CreateWithTx() → abre incidente            │
│  5. incidentRepo.Close() → fecha incidente                   │
│  6. uow.Commit() → confirma transação                        │
│  7. notifier.NotifyDown/NotifyRecovery() → webhook (APÓS!)  │
└─────────────────────────────────────────────────────────────┘
```

### Garantias

- **Atomicidade**: Check + mudança de estado + incidente são persistidos juntos
- **Rollback automático**: Em caso de erro, todas as operações são desfeitas
- **Webhook seguro**: Notificações são enviadas APÓS o commit, garantindo que o estado do banco está correto

### Componentes

| Componente         | Arquivo                                   | Descrição                                  |
| ------------------ | ----------------------------------------- | ------------------------------------------ |
| `UnitOfWork`       | `internal/service/unit_of_work.go`        | Interface para gerenciamento de transações |
| `sqliteUnitOfWork` | `internal/service/unit_of_work_impl.go`   | Implementação com `sql.Tx`                 |
| `CheckRepoTx`      | `internal/repository/check_repo_tx.go`    | Repositório de checks com suporte a Tx     |
| `MonitorRepoTx`    | `internal/repository/monitor_repo_tx.go`  | Repositório de monitors com suporte a Tx   |
| `IncidentRepoTx`   | `internal/repository/incident_repo_tx.go` | Repositório de incidents com suporte a Tx  |

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
│   │   ├── main.jsx
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
| `PINGOU_MAX_REDIRECTS`        |         não | `5`         | Número máximo de redirects que o checker HTTP segue                                                                |
| `PINGOU_GLOBAL_TIMEOUT`       |         não | `60`        | Timeout global em segundos para requisições HTTP                                                                   |

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

Baixa um dump do banco SQLite gerado com `VACUUM INTO` em arquivo temporário.

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

Os alvos do `Makefile` e seu propósito (uso típico):

- `start`: Compila e executa o binário usando as variáveis do `.env`.
  - Uso: `make start`

- `dev`: Roda a aplicação em modo desenvolvimento usando as variáveis de `.env`.
  - Uso: `make dev`

- `fmt`: Formata o código Go com `gofumpt`.
  - Uso: `make fmt`

- `build`: Compila o frontend e o binário Go (usa `CGO_ENABLED=1`). Gera `bin/pingou`.
  - Uso: `make build` (invoca `build-web` internamente)

- `test`: Executa os testes Go em todo o módulo.
  - Uso: `make test`

- `clean`: Remove artefatos locais como `bin/`.
  - Uso: `make clean`

- `build-web`: Instala dependências do frontend e gera o bundle via Vite.
  - Uso: `make build-web`

- `docker-build`: Constrói a imagem Docker multi-stage e a marca com `$(VERSION)` e `latest`.
  - Uso: `make docker-build`

- `docker-up`: Sobe os serviços via `docker compose up --build` (desenvolvimento/integração local).
  - Uso: `make docker-up`

- `docker-down`: Derruba os serviços levantados via `docker compose`.
  - Uso: `make docker-down`

- `docker-size`: Constrói uma imagem temporária `pingou:size-check`, imprime informação de tamanho e remove a imagem.
  - Objetivo: validar automaticamente o requisito de tamanho da imagem (PRD).
  - Uso: `make docker-size`

- `docker-startup-test`: Mede o tempo de startup do `docker compose` (aguarda healthcheck), considera < 30s como sucesso.
  - Objetivo: validar o critério de sucesso `docker compose up < 30s` automaticamente.
  - Uso: `make docker-startup-test`

- `release`: Gera o build final (compila e imprime o local do binário em `bin/pingou`).
  - Uso: `make release`

Cada target tem comentários e comportamentos encadeados no `Makefile` (por exemplo, `build` roda `build-web` antes de compilar o binário). Use `make <target>` para executar o fluxo desejado.

## CLI

O Pingou inclui um binário CLI para operações de banco de dados e debug sem passar pela API.

### Instalação

```bash
make build
# ou compilar diretamente
go build -o pingou ./cmd/pingou
```

### Comandos

| Comando                 | Descrição                                                |
| ----------------------- | -------------------------------------------------------- |
| `pingou add`            | Cria um novo monitor                                     |
| `pingou rm`             | Remove um monitor pelo ID                                |
| `pingou serve`          | Inicia o servidor API + dashboard (comportamento padrão) |
| `pingou migrate up`     | Executa migrations pendentes                             |
| `pingou migrate down`   | Reverte a última migration                               |
| `pingou migrate status` | Mostra o status das migrations                           |
| `pingou export-db`      | Exporta o banco SQLite                                   |
| `pingou version`        | Exibe versão do binário                                  |
| `pingou config`         | Exibe configuração atual                                 |

### Detalhes

#### serve

Inicia o servidor HTTP com API REST e dashboard SPA.

```bash
# Porta padrão
./pingou serve

# Porta customizada
PINGOU_PORT=9999 ./pingou serve
```

**Proteção anti-multinstância**: O comando usa um lock file em `~/.pingou/pingou.lock` para evitar múltiplas instâncias. Se já houver um servidor rodando, retorna erro.

#### add

Cria um novo monitor de health check.

```bash
# Exemplo básico
./pingou add -n "Meu Servico" -u "https://api.exemplo.com/health" --key <API_KEY>

# Com flags opcionais
./pingou add -n "API" -u "https://api.exemplo.com/health" -i 30 -t 10 -K 5 --key <API_KEY>
```

**Flags**:

| Flag          | Abrev | Padrão        | Descrição                        |
| ------------- | ----- | ------------- | -------------------------------- |
| `--name`      | `-n`  | (obrigatório) | Nome do monitor                  |
| `--url`       | `-u`  | (obrigatório) | URL para verificar health        |
| `--interval`  | `-i`  | `60`          | Intervalo em segundos            |
| `--timeout`   | `-t`  | `5`           | Timeout em segundos              |
| `--threshold` | `-K`  | `3`           | Falhas antes de marcar como down |

O monitor é sempre criado com `enabled: true`.

#### rm

Remove um monitor pelo ID.

```bash
./pingou rm --id 019def55-0c92-76a5-a7c4-b573ab447ac2 --key <API_KEY>
```

**Flags**:

| Flag   | Abrev | Padrão        | Descrição     |
| ------ | ----- | ------------- | ------------- |
| `--id` | `-i`  | (obrigatório) | ID do monitor |

#### migrate

Gerencia migrations do banco de dados.

```bash
# Aplicar migrations
./pingou migrate up --key <API_KEY>

# Reverter última migration
./pingou migrate down --key <API_KEY>

# Ver status
./pingou migrate status --key <API_KEY>
```

> **Nota**: Requer `--key` igual ao `PINGOU_API_KEY` do `.env`.

#### export-db

Exporta o banco de dados atual para um arquivo SQLite.

```bash
# Sem argumento: cria exported_<banco>.db no PWD
./pingou export-db --key <API_KEY>

#指定 output path
./pingou export-db -o /tmp/backup.db --key <API_KEY>
```

#### version

Exibe informações de versão.

```bash
./pingou version
# Output:
# pingou dev
# commit: none
# build date: unknown
```

#### config

Exibe a configuração atual (sem secrets).

```bash
./pingou config --key <API_KEY>
# Output:
# {
#   "DatabaseURL": "./pingou.db",
#   "Port": "8080",
#   "LogLevel": "debug",
#   "CORSAllowedOrigins": null,
#   "MaxRedirects": 5,
#   "GlobalTimeout": 60
# }
```

### Flags globais

| Flag            | Descrição                                                           |
| --------------- | ------------------------------------------------------------------- |
| `--key`         | API key para comandos protegidos (`migrate`, `export-db`, `config`) |
| `-v, --verbose` | Modo verboso                                                        |

### Exit codes

- `0`: Sucesso
- `1`: Erro

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
