# Code Conventions & Standards

## Go Code Conventions

### Formatting & Style

- **gofumpt**: Used for code formatting (stricter than gofmt)
- **Standard Go conventions**: Following Go's official style guide
- **Package naming**: Lowercase, short, descriptive names
- **Exported names**: PascalCase for public APIs
- **Private names**: camelCase for internal use

### File Organization

```
internal/
├── checker/        # Health checking logic
├── config/         # Configuration management
├── database/       # Database connection and setup
├── domain/         # Domain entities and models
├── handler/        # HTTP handlers and routing
├── repository/     # Data access layer
├── scheduler/      # Background job scheduling
└── service/        # Business logic layer
```

### Naming Patterns

- **Interfaces**: Usually end with `er` suffix (e.g., `Repository`, `Checker`)
- **Structs**: Descriptive names, often domain-specific (e.g., `Monitor`, `Incident`)
- **Functions**: Verb-noun pattern for actions (e.g., `NewMonitorService`)
- **Variables**: Short but meaningful names
- **Constants**: UPPER_SNAKE_CASE

### Error Handling

- **Explicit error handling**: Always check and handle errors
- **Error wrapping**: Use `fmt.Errorf` with `%w` for error wrapping
- **Structured logging**: Use `log/slog` with context and key-value pairs
- **Error messages**: Descriptive, actionable messages

### Concurrency Patterns

- **Context propagation**: Pass context through all layers
- **Goroutine lifecycle**: Proper cleanup and shutdown handling
- **Channel usage**: Buffered channels where appropriate
- **Mutex usage**: Protect shared state when necessary

## Frontend Code Conventions

### React Patterns

- **Functional components**: Use function components with hooks
- **Custom hooks**: Extract reusable logic into custom hooks
- **Props interface**: TypeScript interfaces for props
- **State management**: Local state with useState/useEffect

### File Structure

```
src/
├── api/           # API client functions and types
├── components/    # Reusable UI components
├── hooks/         # Custom React hooks
├── pages/         # Route-level components
└── theme/         # Styling and theme configuration
```

### Naming Conventions

- **Components**: PascalCase (e.g., `MonitorList`, `IncidentCard`)
- **Files**: camelCase for components, kebab-case for utilities
- **Functions**: camelCase for functions and methods
- **Constants**: UPPER_SNAKE_CASE
- **CSS classes**: kebab-case following BEM methodology

### API Integration

- **Axios client**: Centralized API client configuration
- **Error handling**: Consistent error handling across components
- **Loading states**: Proper loading and error state management
- **TypeScript**: Type definitions for API responses

## Database Conventions

### Schema Design

- **Table naming**: plural_noun format (e.g., `monitors`, `incidents`)
- **Column naming**: snake_case format
- **Primary keys**: `id` column with auto-increment
- **Foreign keys**: `table_id` format (e.g., `monitor_id`)
- **Timestamps**: `created_at`, `updated_at` columns

### Migration Conventions

- **Goose migrations**: Version-controlled database migrations
- **Up/Down methods**: Both upgrade and downgrade paths
- **Descriptive names**: Clear migration file names
- **Idempotent operations**: Migrations should be re-runnable

### Query Patterns

- **Prepared statements**: Use parameterized queries
- **Transaction handling**: Proper transaction management
- **Error handling**: Check database errors appropriately
- **Connection management**: Proper connection lifecycle

## API Conventions

### REST API Design

- **Resource naming**: Plural nouns for collections (e.g., `/api/monitors`)
- **HTTP methods**: Standard RESTful method usage
  - `GET`: Retrieve resources
  - `POST`: Create resources
  - `PATCH`: Partial updates
  - `DELETE`: Remove resources
- **Status codes**: Appropriate HTTP status codes
- **Response format**: Consistent JSON responses

### Authentication

- **API key header**: `X-API-Key` header for authentication
- **Global key**: Single API key for all endpoints
- **Error responses**: 401 for missing/invalid keys
- **Public endpoints**: Only health check endpoint is public

### Response Format

```json
{
  "data": { ... },
  "error": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Configuration Conventions

### Environment Variables

- **Naming**: UPPER_SNAKE_CASE
- **Required vs optional**: Clear documentation of required variables
- **Default values**: Sensible defaults where appropriate
- **Validation**: Input validation for configuration values

### Configuration Structure

- **Centralized config**: Single config struct
- **Environment loading**: Use godotenv for development
- **Type safety**: Strong typing for configuration
- **Validation**: Runtime validation of configuration

## Testing Conventions

### Go Testing

- **Table-driven tests**: Use table-driven test patterns
- **Test naming**: `TestFunctionName_Condition_ExpectedResult`
- **Subtests**: Use `t.Run()` for related test cases
- **Mock interfaces**: Mock external dependencies

### Frontend Testing

- **Component testing**: Test component behavior
- **Integration testing**: Test API integration
- **User interactions**: Test user interaction flows
- **Error scenarios**: Test error handling paths

## Documentation Conventions

### Code Comments

- **Package documentation**: Package-level documentation
- **Exported functions**: Document public APIs
- **Complex logic**: Comment non-obvious code
- **TODO comments**: Mark future work clearly

### API Documentation

- **Endpoint documentation**: Clear endpoint descriptions
- **Request/response examples**: Provide example payloads
- **Error documentation**: Document error responses
- **Authentication notes**: Document auth requirements

## Git Conventions

### Commit Messages

- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- **Description**: Imperative mood, lowercase
- **Body**: Detailed explanation when necessary

### Branch Strategy

- **Main branch**: `main` for production code
- **Feature branches**: `feature/description` for new features
- **Bug fixes**: `fix/description` for bug fixes
- **Release branches**: `release/version` for releases

## Build and Deployment Conventions

### Makefile Targets

- **Standard targets**: `build`, `run`, `test`, `clean`
- **Development targets**: `fmt`, `lint`, `dev`
- **Docker targets**: `docker-build`, `docker-up`, `docker-down`
- **Release targets**: `release` for production builds

### Version Management

- **Semantic versioning**: Follow SemVer principles
- **Build info**: Include version, commit, and build date
- **Release process**: Automated release process
- **Tagging**: Git tags for releases

## Security Conventions

### Input Validation

- **API validation**: Validate all input data
- **SQL injection prevention**: Use parameterized queries
- **XSS prevention**: Proper output encoding
- **CSRF protection**: Implement CSRF protection

### Data Protection

- **Sensitive data**: Avoid logging sensitive information
- **API key security**: Secure API key storage and transmission
- **Database security**: Proper database file permissions
- **Network security**: HTTPS in production
