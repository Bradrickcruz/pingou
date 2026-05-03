# Project Structure Analysis

## Root Directory Structure

```
pingou-health-checker/
├── .git/                              # Git repository metadata
├── .agents/                           # Agent-related configurations
├── .cache/                            # Build and development cache
├── .codex/                            # CodeX IDE configurations
├── .specs/                            # Project specifications and documentation
├── .windsurf/                         # Windsurf IDE configurations
├── bin/                               # Compiled binaries
├── cmd/                               # Application entry points
├── internal/                          # Internal application code
├── tasks/                             # Development task specifications
├── web/                               # Frontend React application
├── .dockerignore                      # Docker ignore patterns
├── .editorconfig                      # Editor configuration
├── .env                               # Local environment variables
├── .env.example                       # Environment variable template
├── .gitignore                         # Git ignore patterns
├── DIVERGENCIAS-PRD-IMPLEMENTACAO.md  # Detailed implementation divergences
├── docker-compose.yml                 # Multi-container orchestration
├── Dockerfile                         # Container build configuration
├── go.mod                             # Go module definition
├── go.sum                             # Go module checksums
├── LICENSE                            # Apache 2.0 license
├── Makefile                           # Enhanced build automation
├── pingou.db                          # SQLite database (runtime)
├── PRD-divergencias-implementacao.md  # PRD implementation divergences
├── PRD.md                             # Product Requirements Document
├── README.md                          # Project documentation
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
│   ├── settings_repo.go      # Settings data access
│   ├── health_check_repo.go  # Health check data access
│   ├── notification_repo.go  # Notification data access
│   └── user_repo.go          # User data access (if multi-user support)
├── scheduler/                 # Background job scheduling
│   ├── retention_worker.go   # Data retention management
│   └── scheduler.go          # Monitor check scheduling
└── service/                   # Business logic layer
    ├── incident_service.go    # Incident management logic
    ├── monitor_service.go    # Monitor management logic
    ├── settings_service.go   # Settings management logic
    ├── state_machine.go      # State transition logic
    ├── webhook_notifier.go   # Webhook notification logic
    ├── notification_service.go # Notification management logic
    ├── health_check_service.go # Health check orchestration
    └── export_service.go     # Data export service
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
│   ├── styles/              # Styling with Tailwind CSS
│   │   ├── globals.css      # Global Tailwind styles
│   │   └── components.css   # Component-specific styles
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
├── postcss.config.js        # PostCSS configuration
├── tailwind.config.js       # Tailwind CSS configuration
└── vite.config.js           # Vite build configuration
```

## Application Entry Point (`cmd/`)

```
cmd/
└── pingou/
    ├── main.go              # Application main entry point
    └── commands/            # CLI command implementations
        ├── root.go          # Root command and global flags
        ├── serve.go         # HTTP server command
        ├── migrate.go       # Database migration command
        ├── export.go        # Database export command
        └── info.go          # Application information command
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
- **Makefile**: Enhanced build automation with Docker testing

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
