# Technology Stack Analysis

## Backend Stack

### Core Language & Runtime

- **Go 1.25.1** - Primary backend language
- **Standard library** - `net/http`, `log/slog`, `embed`, `context`

### Database & Persistence

- **SQLite** - Primary database for lightweight, self-hosted deployment
- **github.com/mattn/go-sqlite3** - SQLite driver
- **github.com/pressly/goose/v3** - Database migration tool

### HTTP & Networking

- **net/http** - Standard library HTTP server
- **Custom HTTP checker** - For monitoring endpoints

### Configuration & Environment

- **github.com/joho/godotenv** - Environment variable loading
- **github.com/spf13/viper** - Advanced configuration management
- **Environment variables** - Configuration management

### CLI Framework

- **github.com/spf13/cobra** - Command-line interface framework
- **github.com/spf13/pflag** - Command-line flag parsing

### Utilities

- **github.com/google/uuid** - UUID generation
- **log/slog** - Structured logging

## Frontend Stack

### Core Framework

- **React 19.2.5** - UI framework
- **React DOM 19.2.5** - DOM rendering
- **React Router DOM 6.30.3** - Client-side routing

### Build Tools & Development

- **Vite 6.4.2** - Build tool and dev server
- **@vitejs/plugin-react 4.7.0** - React plugin for Vite
- **ESLint 10.2.1** - Code linting

### Styling & CSS Framework

- **Tailwind CSS 3.4.17** - Utility-first CSS framework
- **PostCSS 8.5.13** - CSS transformation tool
- **Autoprefixer 10.5.0** - CSS vendor prefixing

### HTTP Client

- **Axios 1.15.2** - HTTP client for API calls

### Development Tools

- **TypeScript types** - `@types/react`, `@types/react-dom`
- **ESLint plugins** - React hooks and refresh

## Infrastructure & Deployment

### Containerization

- **Docker** - Container runtime
- **Docker Compose** - Multi-container orchestration

### Build & Distribution
- **Makefile** - Enhanced build automation with Docker testing
- **embed.FS** - Frontend embedding in Go binary
- **Single binary distribution** - Self-contained deployment
- **Docker compose** - Multi-container orchestration

### Development Workflow
- **gofumpt** - Go code formatting
- **CGO** - Required for SQLite driver
- **Docker startup testing** - Automated container startup validation
- **Docker size optimization** - Image size monitoring

## Architecture Patterns

### Backend Architecture

- **Domain-Driven Design** - Clear separation of concerns
- **Repository Pattern** - Data access abstraction
- **Service Layer** - Business logic encapsulation
- **Handler Layer** - HTTP request/response handling
- **Scheduler Pattern** - Background job execution

### Frontend Architecture

- **Component-based** - React components
- **API-first** - RESTful API integration
- **Single Page Application** - SPA architecture

## Key Dependencies Analysis

### Critical Dependencies

- **go-sqlite3** - Core data persistence
- **goose** - Database schema management
- **React** - Frontend framework
- **Vite** - Build system
- **Cobra** - CLI framework for command structure
- **Viper** - Configuration management
- **Tailwind CSS** - UI styling framework

### Development Dependencies

- **ESLint** - Code quality
- **TypeScript types** - Development experience

### Optional Dependencies

- **Webhook functionality** - External notifications
- **Docker** - Container deployment

## Stack Strengths

1. **Simplicity** - Minimal external dependencies
2. **Self-contained** - Single binary deployment
3. **Lightweight** - SQLite for data persistence
4. **Modern** - Recent versions of Go and React
5. **Standard library focused** - Leverages Go's robust stdlib

## Stack Considerations

1. **CGO requirement** - SQLite driver needs CGO
2. **Single database** - SQLite limits for high concurrency
3. **Embedded frontend** - Tight coupling between backend and frontend
4. **Manual migrations** - Requires careful database versioning
