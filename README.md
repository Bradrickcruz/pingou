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
├── go.mod
└── README.md
```
