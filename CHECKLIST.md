# Implementation Checklist

This checklist contains all tasks needed to implement the web API server boilerplate according to ARCHITECTUREv2.md.

## Phase 1: Foundation

### 1.1 Project Setup
- [x] **Done** Initialize Go module: `go mod init gosuda.org/boilerplate`
- [x] **Done** Create basic directory structure
- [x] **Done** Add core dependencies to go.mod
- [x] **Done** Create Makefile with common commands
- [x] **Done** Create .gitignore file

### 1.2 Configuration Management
- [x] **Done** Create `config/defaults.yaml` with embedded configuration
- [x] **Done** Implement `internal/config/config.go` with struct definitions
- [x] **Done** Add `//go:embed` directive for defaults.yaml
- [x] **Done** Implement environment variable override logic
- [x] **Done** Add configuration validation
- [x] **Done** Write `internal/config/config_test.go` (≥80% coverage)

### 1.3 Basic Logging Setup
- [x] **Done** Implement `internal/infrastructure/logger.go`
- [x] **Done** Add slog-based structured logging
- [x] **Done** Support environment-specific formats (text/json)
- [x] **Done** Add contextual logging with request ID support
- [x] **Done** Write logger tests

### 1.4 Storage Interface and Implementation
- [x] **Done** Define `internal/domain/store.go` interface
- [x] **Done** Implement `internal/infrastructure/memory_store.go`
- [x] **Done** Add thread-safe map with RWMutex
- [x] **Done** Implement JSON serialization for complex objects
- [x] **Done** Add prefix-based listing support
- [x] **Done** Write storage tests (≥80% coverage)

### 1.5 Basic Error Handling
- [x] **Done** Define `internal/domain/errors.go` with domain errors
- [ ] **WIP** Implement error types: PostNotFoundError, InvalidPostDataError, etc.
- [ ] Add error code constants
- [ ] Write error tests

## Phase 2: HTTP Layer

### 2.1 OpenAPI Specification
- [ ] Create `api/openapi.yaml` with complete API specification
- [ ] Define all endpoints: /debug/*, /posts, /posts/{id}
- [ ] Add proper request/response schemas
- [ ] Include validation rules and examples
- [ ] Add error response schemas

### 2.2 Code Generation Setup
- [ ] Add oapi-codegen to go.mod
- [ ] Create `api/gen.go` with go:generate directive
- [ ] Configure oapi-codegen for strict-server mode
- [ ] Generate initial code
- [ ] Verify generated code compiles

### 2.3 Basic Middleware Stack
- [ ] Implement `internal/middleware/recovery.go`
- [ ] Implement `internal/middleware/request_id.go`
- [ ] Implement `internal/middleware/logging.go`
- [ ] Implement `internal/middleware/cors.go`
- [ ] Write middleware tests (≥80% coverage)

### 2.4 Error Handling Middleware
- [ ] Implement `internal/middleware/error_handler.go`
- [ ] Add centralized error translation logic
- [ ] Map domain errors to HTTP status codes
- [ ] Add request ID to error responses
- [ ] Write error handler tests

## Phase 3: Business Logic

### 3.1 Domain Models
- [ ] Create `internal/domain/post.go` with Post entity
- [ ] Add validation logic for Post fields
- [ ] Implement ID generation strategy
- [ ] Add timestamp handling
- [ ] Write domain model tests

### 3.2 Application Services
- [ ] Implement `internal/application/post_service.go`
- [ ] Add CRUD operations for posts
- [ ] Implement business logic validation
- [ ] Add pagination logic
- [ ] Write service tests (≥80% coverage)

### 3.3 Debug Service
- [ ] Implement `internal/application/debug_service.go`
- [ ] Add metrics collection logic
- [ ] Add log level adjustment logic
- [ ] Add pprof endpoint handling
- [ ] Write debug service tests

### 3.4 Pagination Implementation
- [ ] Create `internal/application/pagination.go`
- [ ] Implement cursor-based pagination
- [ ] Add limit validation
- [ ] Add cursor encoding/decoding
- [ ] Write pagination tests

### 3.5 Handler Implementations
- [ ] Implement `api/impl.go` with all handlers
- [ ] Add strict-server interface implementation
- [ ] Connect handlers to application services
- [ ] Add proper error handling
- [ ] Write handler tests

## Phase 4: Debug Features

### 4.1 Metrics Collection
- [ ] Implement `internal/infrastructure/metrics.go`
- [ ] Add Prometheus metrics collection
- [ ] Define HTTP request metrics
- [ ] Add storage operation metrics
- [ ] Add custom business metrics
- [ ] Write metrics tests

### 4.2 Pprof Endpoints
- [ ] Add pprof import to main.go
- [ ] Configure pprof routes
- [ ] Add security considerations
- [ ] Test pprof endpoints

### 4.3 Log Level Adjustment API
- [ ] Implement runtime log level change
- [ ] Add validation for log levels
- [ ] Add atomic log level switching
- [ ] Write log level tests

### 4.4 Request ID Tracking
- [ ] Ensure request ID propagation through all layers
- [ ] Add request ID to response headers
- [ ] Add request ID to error responses
- [ ] Add request ID to metrics
- [ ] Test request ID propagation

### 4.5 Health Check Endpoint
- [ ] Add `/health` endpoint
- [ ] Include basic health status
- [ ] Add dependency health checks
- [ ] Write health check tests

## Phase 5: Testing & Polish

### 5.1 Unit Tests
- [ ] Write unit tests for all business logic
- [ ] Write unit tests for configuration
- [ ] Write unit tests for storage
- [ ] Write unit tests for services
- [ ] Achieve ≥80% test coverage

### 5.2 Integration Tests
- [ ] Write integration tests for HTTP handlers
- [ ] Test complete request/response cycles
- [ ] Test error scenarios
- [ ] Test middleware chain
- [ ] Test graceful shutdown

### 5.3 Middleware Tests
- [ ] Test request ID middleware
- [ ] Test logging middleware
- [ ] Test error handling middleware
- [ ] Test CORS middleware
- [ ] Test recovery middleware

### 5.4 Graceful Shutdown Implementation
- [ ] Implement signal handling in main.go
- [ ] Add 5-second timeout logic
- [ ] Add resource cleanup
- [ ] Add shutdown logging
- [ ] Test graceful shutdown

### 5.5 Documentation
- [ ] Write comprehensive README.md
- [ ] Add API documentation
- [ ] Add deployment instructions
- [ ] Add development setup guide
- [ ] Add testing instructions

### 5.6 Docker Containerization
- [ ] Create Dockerfile with multi-stage build
- [ ] Add non-root user for security
- [ ] Add health check
- [ ] Optimize image size
- [ ] Test containerized deployment

## Final Verification

### 5.7 Build and Test Verification
- [ ] Run `go build` - verify successful compilation
- [ ] Run `go test ./...` - verify all tests pass
- [ ] Run `go test -cover ./...` - verify ≥80% coverage
- [ ] Run `go vet ./...` - verify no issues
- [ ] Run `golangci-lint run` - verify code quality

### 5.8 API Testing
- [ ] Test all debug endpoints (/debug/metrics, /debug/logs, /debug/pprof/*)
- [ ] Test all post endpoints (GET /posts, POST /posts, GET /posts/{id}, PUT /posts/{id}, DELETE /posts/{id})
- [ ] Test error scenarios and proper HTTP status codes
- [ ] Test request ID propagation
- [ ] Test graceful shutdown

### 5.9 Performance Testing
- [ ] Test concurrent requests
- [ ] Test memory usage
- [ ] Test response times
- [ ] Test graceful shutdown timing
- [ ] Verify metrics collection

## Task Status Legend

- [ ] **TODO** - Task not started
- [ ] **WIP** - Task in progress
- [x] **Done** - Task completed successfully

## Notes

- Each task should be marked as **WIP** when started
- After completing a task, run `go test` and `go build` to verify everything works
- Only mark a task as **Done** after all tests pass and coverage meets requirements
- Update this checklist as tasks are completed