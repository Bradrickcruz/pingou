# Architecture Overview

## High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │   HTTP Server   │    │   Background    │
│   (React SPA)   │◄──►│   (Go Handler)  │◄──►│   Scheduler     │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Service Layer │    │   HTTP Checker  │
                       │                 │    │                 │
                       └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Repository    │    │   State Machine │
                       │     Layer       │    │                 │
                       └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────────────────────────────┐
                       │              SQLite Database             │
                       │                                         │
                       └─────────────────────────────────────────┘
```

## Backend Architecture Layers

### 1. Handler Layer (`internal/handler/`)

- **HTTP Server**: Main entry point for all HTTP requests
- **Authentication**: API key validation for protected routes
- **Response formatting**: JSON responses and error handling
- **Static file serving**: Embedded frontend assets

### 2. Service Layer (`internal/service/`)

- **MonitorService**: Monitor CRUD operations and lifecycle management
- **IncidentService**: Incident management and resolution
- **SettingsService**: Global configuration management
- **StateMachine**: State transition logic for monitor status changes
- **WebhookNotifier**: External notification delivery

### 3. Repository Layer (`internal/repository/`)

- **MonitorRepo**: Monitor data access operations
- **CheckRepo**: Check result persistence
- **IncidentRepo**: Incident data access
- **SettingsRepo**: Configuration data access

### 4. Domain Layer (`internal/domain/`)

- **Monitor**: Core monitor entity with configuration
- **Check**: Check result entity with timing and status
- **Incident**: Incident entity with start/end times
- **Settings**: Configuration entity

## Background Processing Architecture

### Scheduler (`internal/scheduler/`)

- **MonitorScheduler**: Orchestrates periodic health checks
- **RetentionWorker**: Manages data retention policies
- **Context-based lifecycle**: Graceful shutdown handling

### Health Checking (`internal/checker/`)

- **HTTPChecker**: Performs HTTP requests to monitored endpoints
- **Timeout handling**: Configurable request timeouts
- **Response validation**: Status code and response time analysis

## Frontend Architecture

### Component Structure

```
src/
├── components/     # Reusable UI components
├── pages/         # Route-level components
├── hooks/         # Custom React hooks
├── api/           # API client and types
└── theme/         # Styling and UI theming
```

### Data Flow

1. **API Client** (Axios) communicates with backend
2. **React State** manages local component state
3. **Router** handles client-side navigation
4. **Components** render based on API responses

## Data Flow Patterns

### Monitor Check Flow

```
Scheduler → HTTPChecker → StateMachine → Repository → WebhookNotifier
    ↓              ↓              ↓           ↓              ↓
  Trigger      Perform       Analyze     Store Results   Send Notification
```

### API Request Flow

```
Client → Handler → Service → Repository → Database
   ↓        ↓        ↓         ↓          ↓
Request  Auth   Business   Data      Persist
        Logic   Logic     Access
```

## Key Architectural Decisions

### 1. Single Binary Architecture

- **Rationale**: Simplifies deployment and distribution
- **Implementation**: `embed.FS` for frontend assets
- **Trade-offs**: Tighter coupling vs. deployment simplicity

### 2. SQLite as Primary Database

- **Rationale**: Self-hosted, lightweight, no external dependencies
- **Implementation**: Single file database with migrations
- **Trade-offs**: Limited concurrency vs. simplicity

### 3. Domain-Driven Design

- **Rationale**: Clear separation of concerns and business logic
- **Implementation**: Layered architecture with clear boundaries
- **Trade-offs**: More boilerplate vs. maintainability

### 4. Background Scheduler

- **Rationale**: Independent health checking from HTTP serving
- **Implementation**: Context-based goroutine management
- **Trade-offs**: Complexity vs. reliability

## State Management

### Monitor State Machine

```
UNKNOWN ──► UP ──► DOWN ──► UP
    │        │        │
    └────────┴────────┘
```

### State Transitions

- **UNKNOWN → UP**: First successful check
- **UNKNOWN → DOWN**: First failed check
- **UP → DOWN**: Service degradation detected
- **DOWN → UP**: Service recovery detected

## Concurrency Model

### Goroutine Usage

- **HTTP Server**: Main goroutine for request handling
- **Scheduler**: Separate goroutine for periodic checks
- **Retention Worker**: Background goroutine for cleanup
- **Signal Handling**: Context-based cancellation

### Synchronization

- **Database**: SQLite handles concurrent access
- **Channels**: Used for graceful shutdown coordination
- **Context**: Propagates cancellation signals

## Security Architecture

### Authentication

- **API Key**: Single key for all API endpoints
- **Header-based**: `X-API-Key` header validation
- **Frontend Storage**: LocalStorage for API key persistence

### Data Protection

- **No sensitive data logging**: Avoids exposing credentials
- **Input validation**: Basic request validation
- **CORS handling**: Basic cross-origin request handling

## Deployment Architecture

### Container Model

- **Single Container**: All components in one container
- **Volume Mount**: SQLite database persistence
- **Port Exposure**: HTTP port configurable

### Binary Distribution

- **Static Compilation**: Includes embedded frontend
- **Single File**: No external dependencies
- **Cross-platform**: Supports multiple architectures
