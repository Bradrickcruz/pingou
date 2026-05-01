# Testing Strategy & Framework

## Current Testing Status

### Test Coverage Analysis

- **Backend Tests**: No test files currently present in the codebase
- **Frontend Tests**: No test files currently present in the codebase
- **Integration Tests**: No automated integration testing
- **End-to-End Tests**: No E2E testing framework in place

## Recommended Testing Framework

### Backend Testing Stack

- **Go Testing**: Standard Go testing package
- **Testify**: Assertion library and mock generation
- **SQL Mock**: Database mocking for repository tests
- **HTTP Test**: `net/http/httptest` for handler testing
- **Test Containers**: Integration testing with real database

### Frontend Testing Stack

- **Vitest**: Modern testing framework (compatible with Vite)
- **React Testing Library**: Component testing utilities
- **MSW**: API mocking for integration tests
- **Jest DOM**: DOM assertion extensions

## Testing Strategy

### 1. Unit Testing

#### Backend Unit Tests

- **Domain Layer**: Test business logic and entity behavior
- **Service Layer**: Test business rules and orchestration
- **Repository Layer**: Test data access logic with mocks
- **Handler Layer**: Test HTTP request/response handling
- **Checker Layer**: Test health checking logic

#### Frontend Unit Tests

- **Components**: Test component rendering and behavior
- **Hooks**: Test custom hook logic
- **API Client**: Test API interaction functions
- **Utilities**: Test helper functions

### 2. Integration Testing

#### Backend Integration Tests

- **Database Integration**: Test repository with real SQLite
- **API Integration**: Test handler-to-database flow
- **Scheduler Integration**: Test background job execution
- **Webhook Integration**: Test notification delivery

#### Frontend Integration Tests

- **API Integration**: Test component-to-backend communication
- **User Flows**: Test complete user workflows
- **State Management**: Test state persistence and updates

### 3. End-to-End Testing

#### E2E Test Scenarios

- **Monitor Lifecycle**: Create → Check → Incident → Resolve
- **User Authentication**: Login → Access → Logout
- **Settings Management**: Configure → Save → Verify
- **Data Export**: Export → Download → Verify

## Test Organization Structure

### Backend Test Structure

```
internal/
├── checker/
│   └── http_checker_test.go
├── config/
│   └── config_test.go
├── database/
│   └── database_test.go
├── domain/
│   ├── check_test.go
│   ├── incident_test.go
│   ├── monitor_test.go
│   └── settings_test.go
├── handler/
│   ├── auth_test.go
│   ├── export_test.go
│   ├── incidents_test.go
│   ├── monitors_test.go
│   ├── settings_test.go
│   └── static_test.go
├── repository/
│   ├── check_repo_test.go
│   ├── incident_repo_test.go
│   ├── monitor_repo_test.go
│   └── settings_repo_test.go
├── scheduler/
│   ├── retention_worker_test.go
│   └── scheduler_test.go
└── service/
    ├── incident_service_test.go
    ├── monitor_service_test.go
    ├── settings_service_test.go
    ├── state_machine_test.go
    └── webhook_notifier_test.go
```

### Frontend Test Structure

```
src/
├── api/
│   ├── api.test.ts
│   ├── incidents.test.ts
│   ├── monitors.test.ts
│   └── settings.test.ts
├── components/
│   ├── Layout.test.tsx
│   ├── Loading.test.tsx
│   └── Login.test.tsx
├── hooks/
│   ├── useApi.test.ts
│   └── useAuth.test.ts
├── pages/
│   ├── IncidentsPage.test.tsx
│   ├── MonitorsPage.test.tsx
│   └── SettingsPage.test.tsx
└── setupTests.ts
```

## Test Data Management

### Test Fixtures

- **Monitors**: Sample monitor configurations
- **Incidents**: Sample incident data
- **Settings**: Sample configuration data
- **API Responses**: Mock API response data

### Test Database

- **In-memory SQLite**: For fast unit tests
- **Testcontainers**: For integration tests
- **Migration Testing**: Test database migrations
- **Seed Data**: Consistent test data setup

## Testing Best Practices

### Backend Testing Guidelines

- **Table-driven tests**: Use Go's table-driven test pattern
- **Subtests**: Use `t.Run()` for related test cases
- **Mock interfaces**: Mock external dependencies
- **Error testing**: Test both success and error paths
- **Race condition testing**: Use `go test -race` for concurrent code

### Frontend Testing Guidelines

- **Behavior testing**: Test component behavior, not implementation
- **User interactions**: Test user interaction flows
- **Accessibility**: Test accessibility features
- **Responsive design**: Test different screen sizes
- **Error states**: Test error handling and display

## Test Automation

### Continuous Integration

- **Pre-commit hooks**: Run tests before commits
- **Pull request tests**: Run full test suite on PRs
- **Coverage reporting**: Generate and track test coverage
- **Performance tests**: Monitor test execution time

### Test Execution Commands

```bash
# Backend tests
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -race ./...              # Race condition detection
go test -cover ./...             # Coverage report
go test -bench ./...             # Benchmark tests

# Frontend tests
npm test                          # Run all tests
npm run test:coverage            # Coverage report
npm run test:watch               # Watch mode
npm run test:e2e                 # E2E tests
```

## Mocking Strategy

### Backend Mocking

- **Database**: Use `sqlmock` for repository tests
- **HTTP client**: Mock external HTTP calls
- **Time**: Mock time for deterministic testing
- **UUID**: Mock UUID generation for predictable tests

### Frontend Mocking

- **API calls**: Use MSW to mock API responses
- **Browser APIs**: Mock browser APIs when needed
- **Time**: Mock time for consistent testing
- **LocalStorage**: Mock storage for auth tests

## Test Environment Setup

### Backend Test Environment

- **Test database**: In-memory SQLite for speed
- **Test configuration**: Separate test config
- **Logging**: Test-specific logging configuration
- **Environment variables**: Test-specific environment setup

### Frontend Test Environment

- **Test setup**: Configure test environment in `setupTests.ts`
- **Mock service workers**: Setup MSW for API mocking
- **Test utilities**: Custom test utilities and helpers
- **Component testing**: Configure React Testing Library

## Performance Testing

### Load Testing Scenarios

- **Concurrent monitors**: Test with many active monitors
- **API throughput**: Test API endpoint performance
- **Database performance**: Test query performance
- **Memory usage**: Monitor memory consumption

### Performance Test Tools

- **Go benchmarks**: Built-in Go benchmarking
- **Apache Bench (ab)**: HTTP load testing
- **Database profiling**: SQLite performance analysis
- **Memory profiling**: Go memory profiling tools

## Security Testing

### Security Test Cases

- **Authentication**: Test API key validation
- **Input validation**: Test input sanitization
- **SQL injection**: Test parameterized queries
- **XSS prevention**: Test output encoding

### Security Testing Tools

- **Go security scanners**: Static analysis for security issues
- **Dependency scanning**: Check for vulnerable dependencies
- **OWASP ZAP**: Web application security scanning
- **Manual security review**: Code review for security issues

## Quality Gates

### Code Coverage Requirements

- **Backend**: Minimum 80% line coverage
- **Frontend**: Minimum 75% line coverage
- **Critical paths**: 100% coverage for critical business logic
- **Integration tests**: Cover all major user workflows

### Test Quality Standards

- **Test naming**: Clear, descriptive test names
- **Test isolation**: Tests should not depend on each other
- **Test speed**: Unit tests should run quickly
- **Test reliability**: Tests should be flake-free
