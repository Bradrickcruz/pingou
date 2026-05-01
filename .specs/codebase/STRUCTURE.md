# Project Structure Analysis

## Root Directory Structure

```
pingou-health-checker/
├── .agents/                    # Agent-related configurations
├── .claude/                    # Claude/AI assistant files
├── .cursor/                    # Cursor IDE configurations
├── .dockerignore              # Docker ignore patterns
├── .env.example               # Environment variable template
├── .git/                      # Git repository metadata
├── .gitignore                 # Git ignore patterns
├── .windsurf/                 # Windsurf IDE configurations
├── Dockerfile                 # Container build configuration
├── LICENSE                    # Apache 2.0 license
├── PRD.md                     # Product Requirements Document
├── README.md                  # Project documentation
├── cmd/                       # Application entry points
├── docker-compose.yml         # Multi-container orchestration
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── internal/                  # Internal application code
├── makefile                   # Build automation
├── migrations/                # Database migration files
├── web/                       # Frontend React application
├── pingou.db                  # SQLite database (runtime)
├── pingou.db-shm              # SQLite shared memory (runtime)
├── pingou.db-wal              # SQLite write-ahead log (runtime)
└── bin/                       # Compiled binaries
```

## Backend Structure (`internal/`)

```
internal/
├── checker/                   # Health checking components
│   └── http_checker.go       # HTTP health checker implementation
├── config/                    # Configuration management
│   └── config.go             # Configuration loading and validation
├── database/                  # Database setup and utilities
│   ├── database.go           # Database connection and initialization
│   ├── migrations.go         # Migration management
│   └── schemas.go            # Database schema definitions
├── domain/                    # Domain entities and models
│   ├── check.go              # Check result entity
│   ├── incident.go           # Incident entity
│   ├── monitor.go            # Monitor entity
│   └── settings.go           # Settings entity
├── handler/                   # HTTP request handlers
│   ├── auth.go               # Authentication middleware
│   ├── export.go             # Database export handler
│   ├── incidents.go          # Incident API handlers
│   ├── monitors.go           # Monitor API handlers
│   ├── server.go             # HTTP server setup
│   ├── settings.go           # Settings API handlers
│   └── static.go             # Static file serving
├── repository/                # Data access layer
│   ├── check_repo.go         # Check data access
│   ├── incident_repo.go      # Incident data access
│   ├── monitor_repo.go       # Monitor data access
│   └── settings_repo.go      # Settings data access
├── scheduler/                 # Background job scheduling
│   ├── retention_worker.go   # Data retention management
│   └── scheduler.go          # Monitor check scheduling
└── service/                   # Business logic layer
    ├── incident_service.go    # Incident management logic
    ├── monitor_service.go    # Monitor management logic
    ├── settings_service.go   # Settings management logic
    ├── state_machine.go      # State transition logic
    └── webhook_notifier.go   # Webhook notification logic
```

## Frontend Structure (`web/`)

```
web/
├── dist/                      # Build output directory
├── node_modules/              # NPM dependencies
├── public/                    # Static assets
│   ├── favicon.ico           # Site favicon
│   └── vite.svg              # Vite logo
├── src/                       # Source code
│   ├── api/                  # API client layer
│   │   ├── api.ts           # Axios client configuration
│   │   ├── index.ts         # API exports
│   │   ├── incidents.ts     # Incident API functions
│   │   ├── monitors.ts      # Monitor API functions
│   │   └── settings.ts      # Settings API functions
│   ├── components/           # Reusable UI components
│   │   ├── Layout.tsx       # Main layout component
│   │   ├── Loading.tsx      # Loading indicator
│   │   └── Login.tsx        # Login form component
│   ├── hooks/               # Custom React hooks
│   │   ├── useApi.ts        # API interaction hook
│   │   └── useAuth.ts       # Authentication hook
│   ├── pages/               # Route-level components
│   │   ├── IncidentsPage.tsx # Incident management page
│   │   ├── MonitorsPage.tsx  # Monitor management page
│   │   └── SettingsPage.tsx  # Settings configuration page
│   ├── theme/               # Styling and theming
│   │   ├── index.css        # Global styles
│   │   └── variables.css    # CSS variables
│   ├── App.tsx              # Main application component
│   ├── index.css            # Application styles
│   ├── main.tsx             # Application entry point
│   └── vite-env.d.ts        # Vite type definitions
├── .gitignore               # Frontend git ignore
├── README.md                # Frontend documentation
├── eslint.config.js         # ESLint configuration
├── index.html               # HTML template
├── package.json             # NPM package configuration
├── package-lock.json        # NPM dependency lock
└── vite.config.js           # Vite build configuration
```

## Application Entry Point (`cmd/`)

```
cmd/
└── pingou/
    └── main.go              # Application main entry point
```

## Database Structure (`migrations/`)

```
migrations/
├── 00001_create_monitors.sql    # Monitors table creation
├── 00002_create_checks.sql      # Checks table creation
├── 00003_create_incidents.sql   # Incidents table creation
└── 00004_create_settings.sql    # Settings table creation
```

## Configuration Files

### Build and Deployment

- **Dockerfile**: Multi-stage build for containerized deployment
- **docker-compose.yml**: Local development with volume persistence
- **makefile**: Build automation with development and production targets

### Development Configuration

- **.env.example**: Environment variable template with documentation
- **.gitignore**: Git ignore patterns for Go, Node.js, and IDE files
- **go.mod/go.sum**: Go module management and dependency locking

### IDE Configuration

- **.windsurf/**: Windsurf IDE specific configurations
- **.cursor/**: Cursor IDE specific configurations
- **.claude/**: Claude AI assistant configurations

## Key Architectural Patterns

### Layered Architecture

1. **Handler Layer**: HTTP request/response handling
2. **Service Layer**: Business logic and orchestration
3. **Repository Layer**: Data access abstraction
4. **Domain Layer**: Core business entities

### Separation of Concerns

- **Backend**: Pure Go with clear module boundaries
- **Frontend**: Separate React application with its own build process
- **Database**: SQLite with migration-based schema management
- **Infrastructure**: Docker-based deployment with volume persistence

### Module Organization

- **Domain-driven**: Business logic organized by domain concepts
- **Feature-based**: Related functionality grouped together
- **Infrastructure concerns**: Separated from business logic

## Data Flow Architecture

### Request Processing

```
HTTP Request → Handler → Service → Repository → Database
```

### Background Processing

```
Scheduler → Checker → State Machine → Repository → Webhook
```

### Frontend Data Flow

```
Component → API Client → HTTP Request → Backend Response → State Update
```

## Build and Distribution Structure

### Development Build

- **Backend**: `go build` creates binary in `bin/`
- **Frontend**: `vite build` creates assets in `web/dist/`
- **Integration**: Frontend embedded in Go binary via `embed.FS`

### Production Distribution

- **Single binary**: All components embedded
- **Docker image**: Containerized with volume mounting
- **Database**: SQLite file persisted externally

## Configuration Management

### Environment-based Configuration

- **Development**: `.env` file with local settings
- **Production**: Environment variables or container environment
- **Runtime**: Settings stored in database and loaded at startup

### Feature Flags

- **Webhook integration**: Configurable via settings
- **Data retention**: Configurable retention policies
- **Check intervals**: Per-monitor configurable intervals

## Security Structure

### Authentication Layer

- **API Key**: Single key for all protected endpoints
- **Header-based**: `X-API-Key` header validation
- **Frontend storage**: LocalStorage for key persistence

### Data Protection

- **Input validation**: Request validation at handler layer
- **SQL safety**: Parameterized queries throughout
- **Error handling**: Secure error responses without information leakage
