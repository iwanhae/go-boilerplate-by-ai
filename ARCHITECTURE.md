# Architecture Design

## Overview

This document outlines the architecture for a fully-functional web API server boilerplate using oapi-codegen and go-chi, following Hexagonal Architecture principles.

## Architecture Layers

### 1. HTTP Layer (Ports & Adapters)
- **Router**: go-chi router with oapi-codegen generated handlers
- **Middleware Stack**:
  - Request ID middleware (assigns unique ID to each request)
  - Logging middleware (slog-based with contextual logging)
  - Validation middleware (request validation before business logic)
  - Error handling middleware (centralized error translation)
  - CORS middleware
  - Recovery middleware
- **Generated Handlers**: oapi-codegen strict-server mode handlers
- **Response Writers**: JSON response handling with proper status codes

### 2. Application Layer (Use Cases)
- **Post Service**: Business logic for blog post operations
- **Debug Service**: Business logic for debug operations (metrics, logs, pprof)
- **Error Handling**: Domain-specific error types and centralized error translation
- **Request/Response DTOs**: Data transfer objects for API communication

### 3. Domain Layer (Core Business Logic)
- **Post Entity**: Core blog post domain model
- **Post Repository Interface**: Storage abstraction
- **Domain Errors**: Business-specific error types
- **Value Objects**: ID, Title, Content, etc.

### 4. Infrastructure Layer (Adapters)
- **In-Memory Store**: Map-based implementation of Store interface
- **Configuration**: YAML-based with environment variable overrides
- **Logging**: slog implementation with contextual logging
- **Metrics**: Prometheus metrics collection
- **Graceful Shutdown**: Signal handling and resource cleanup

## API Design

### OpenAPI Specification Structure
```yaml
openapi: "3.0.0"
info:
  title: Blog API Server
  version: 1.0.0
paths:
  /debug/metrics:
    get:
      summary: Get Prometheus metrics
  /debug/logs:
    post:
      summary: Adjust log level
      parameters:
        - name: level
          in: query
          required: true
          schema:
            type: string
            enum: [debug, info, warn, error]
  /debug/pprof/{profile}:
    get:
      summary: Get pprof profile
      parameters:
        - name: profile
          in: path
          required: true
          schema:
            type: string
  /posts:
    get:
      summary: List posts
      parameters:
        - name: cursor
          in: query
          schema:
            type: string
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
    post:
      summary: Create post
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePostRequest'
  /posts/{id}:
    get:
      summary: Get post
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
    put:
      summary: Update post
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePostRequest'
    delete:
      summary: Delete post
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string

components:
  schemas:
    Post:
      type: object
      required:
        - id
        - title
        - content
        - createdAt
        - updatedAt
      properties:
        id:
          type: string
        title:
          type: string
        content:
          type: string
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
    CreatePostRequest:
      type: object
      required:
        - title
        - content
      properties:
        title:
          type: string
        content:
          type: string
    UpdatePostRequest:
      type: object
      required:
        - title
        - content
      properties:
        title:
          type: string
        content:
          type: string
    PostList:
      type: object
      required:
        - posts
        - nextCursor
      properties:
        posts:
          type: array
          items:
            $ref: '#/components/schemas/Post'
        nextCursor:
          type: string
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: string
        message:
          type: string
```

## Configuration Management

### Structure
```
config/
├── defaults.yaml          # Embedded default configuration
├── config.go             # Configuration struct and loading logic
└── config_test.go        # Configuration tests
```

### Configuration Fields
- **Server**: Port, host, read/write timeouts
- **Logging**: Level, format (text/json), output
- **Storage**: Type (memory), connection details
- **Debug**: Metrics enabled, pprof enabled
- **CORS**: Allowed origins, methods, headers

### Loading Strategy
1. Load embedded `defaults.yaml` using `//go:embed`
2. Override with environment variables (YAML_PATH format)
3. Validate configuration on startup
4. Provide configuration to all layers

## Storage Interface

```go
type Store interface {
    Set(key string, value any) error
    Get(key string) (value any, err error)
    List(keyPrefix string) (values []any, err error)
    Delete(key string) error
    Close() error
}
```

### In-Memory Implementation
- Thread-safe map with RWMutex
- Key-based storage with prefix support
- JSON serialization for complex objects
- Proper cleanup on Close()

## Error Handling Strategy

### Domain Errors
```go
type PostNotFoundError struct{ ID string }
type InvalidPostDataError struct{ Field string }
type StorageError struct{ Err error }
```

### HTTP Error Mapping
- `PostNotFoundError` → 404 Not Found
- `InvalidPostDataError` → 400 Bad Request
- `StorageError` → 500 Internal Server Error
- Generic errors → 500 Internal Server Error

### Error Response Format
```json
{
  "code": "POST_NOT_FOUND",
  "message": "Post with ID '123' not found"
}
```

## Logging Strategy

### Contextual Logging
- Request ID in all log entries
- Structured logging with slog
- Environment-specific formats (text for dev, JSON for prod)
- Log levels: debug, info, warn, error

### Log Fields
- Request ID
- HTTP method and path
- Status code
- Response time
- User agent
- Error details (when applicable)

## Metrics Collection

### Prometheus Metrics
- HTTP request duration
- HTTP request count by status code
- Active requests
- Storage operations count
- Custom business metrics

### Metrics Endpoint
- `/debug/metrics` returns Prometheus format
- No authentication required
- Real-time metrics collection

## Request ID Strategy

### Generation
- UUID v4 for uniqueness
- Generated in middleware layer
- Added to request context

### Propagation
- All business logic receives context with request ID
- Logging includes request ID
- Response headers include request ID
- Error responses include request ID

## Graceful Shutdown

### Signal Handling
- Listen for SIGTERM
- Wait 5 seconds for in-flight requests
- Cancel root context
- Close all resources (storage, HTTP server)
- Log shutdown progress

### Resource Cleanup
- Close HTTP server gracefully
- Close storage connections
- Flush logs
- Exit with appropriate code

## Testing Strategy

### Test Coverage Requirements
- ≥80% test coverage for all functions
- Exclude trivially simple functions
- Focus on business logic and middleware

### Test Types
- Unit tests for business logic
- Integration tests for HTTP handlers
- Middleware tests
- Configuration tests
- Storage tests

### Test Structure
```
tests/
├── unit/           # Business logic tests
├── integration/    # HTTP handler tests
├── middleware/     # Middleware tests
└── fixtures/       # Test data
```

## Security Considerations

### Input Validation
- Request body validation
- Path parameter validation
- Query parameter validation
- Content-Type validation

### Output Sanitization
- Response header sanitization
- Error message sanitization
- Log data sanitization

### CORS Configuration
- Configurable allowed origins
- Configurable allowed methods
- Configurable allowed headers

## Performance Considerations

### Caching
- No caching for MVP (can be added later)
- Consider Redis for production

### Database
- In-memory storage for MVP
- Consider PostgreSQL for production

### Monitoring
- Prometheus metrics
- Structured logging
- Request tracing with request ID

## Deployment Considerations

### Containerization
- Multi-stage Docker build
- Minimal runtime image
- Health check endpoint

### Configuration
- Environment variable overrides
- Config file mounting
- Secrets management

### Monitoring
- Prometheus metrics endpoint
- Health check endpoint
- Structured logging

## File Structure

```
.
├── api/
│   ├── openapi.yaml      # OpenAPI specification
│   ├── gen.go           # Generated code
│   └── impl.go          # Handler implementations
├── cmd/
│   └── server/
│       └── main.go      # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and interfaces
│   ├── application/     # Use cases and services
│   ├── infrastructure/  # Adapters and implementations
│   └── middleware/      # HTTP middleware
├── config/
│   └── defaults.yaml    # Default configuration
├── docs/                # Documentation
├── tests/               # Test files
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Dependencies

### Core Dependencies
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/oapi-codegen/oapi-codegen/v2` - OpenAPI code generation
- `github.com/oapi-codegen/nethttp-middleware` - Request validation
- `gopkg.in/yaml.v3` - YAML configuration
- `go.uber.org/zap` - Structured logging

### Development Dependencies
- `github.com/stretchr/testify` - Testing utilities
- `github.com/prometheus/client_golang` - Metrics
- `golang.org/x/net/http2` - HTTP/2 support

## Implementation Phases

### Phase 1: Foundation
1. Project setup and dependencies
2. Configuration management
3. Basic logging setup
4. Storage interface and in-memory implementation

### Phase 2: HTTP Layer
1. OpenAPI specification
2. Code generation setup
3. Basic middleware stack
4. Error handling middleware

### Phase 3: Business Logic
1. Domain models
2. Application services
3. Handler implementations
4. Request/response validation

### Phase 4: Debug Features
1. Metrics collection
2. Pprof endpoints
3. Log level adjustment
4. Request ID tracking

### Phase 5: Testing & Polish
1. Unit tests
2. Integration tests
3. Graceful shutdown
4. Documentation

## Success Criteria

1. **Functional Requirements**: All APIs work as specified
2. **Non-Functional Requirements**: Performance, security, maintainability
3. **Test Coverage**: ≥80% coverage achieved
4. **Documentation**: Complete API documentation and README
5. **Deployment**: Ready for containerized deployment