# Architecture Design v2

## Overview

This document outlines the refined architecture for a fully-functional web API server boilerplate using oapi-codegen and go-chi, following Hexagonal Architecture principles and addressing all requirements from PROMPT.md.

## Architecture Layers

### 1. HTTP Layer (Ports & Adapters)
- **Router**: go-chi router with oapi-codegen generated handlers
- **Middleware Stack** (in order):
  1. Recovery middleware (panic recovery)
  2. Request ID middleware (assigns unique UUID v4 to each request)
  3. Logging middleware (slog-based with contextual logging)
  4. CORS middleware (configurable)
  5. Request validation middleware (oapi-codegen validation)
  6. Error handling middleware (centralized error translation)
- **Generated Handlers**: oapi-codegen strict-server mode handlers
- **Response Writers**: JSON response handling with proper status codes and headers

### 2. Application Layer (Use Cases)
- **Post Service**: Business logic for blog post operations (CRUD)
- **Debug Service**: Business logic for debug operations (metrics, logs, pprof)
- **Error Handling**: Domain-specific error types and centralized error translation
- **Request/Response DTOs**: Data transfer objects for API communication
- **Pagination**: Cursor-based pagination implementation

### 3. Domain Layer (Core Business Logic)
- **Post Entity**: Core blog post domain model with validation
- **Post Repository Interface**: Storage abstraction
- **Domain Errors**: Business-specific error types
- **Value Objects**: ID, Title, Content, Timestamps
- **Business Rules**: Post validation, ID generation

### 4. Infrastructure Layer (Adapters)
- **In-Memory Store**: Thread-safe map-based implementation of Store interface
- **Configuration**: YAML-based with environment variable overrides using `//go:embed`
- **Logging**: slog implementation with contextual logging and environment-specific formats
- **Metrics**: Prometheus metrics collection
- **Graceful Shutdown**: Signal handling and resource cleanup with 5-second timeout

## API Design

### OpenAPI Specification Structure
```yaml
openapi: "3.0.0"
info:
  title: Blog API Server
  version: 1.0.0
  description: A fully-functional blog API server with debug capabilities
paths:
  /debug/metrics:
    get:
      summary: Get Prometheus metrics
      description: Returns metrics in Prometheus format for scraping
      responses:
        '200':
          description: Metrics in Prometheus format
          content:
            text/plain:
              schema:
                type: string
  /debug/logs:
    post:
      summary: Adjust log level at runtime
      description: Changes the application log level without restart
      parameters:
        - name: level
          in: query
          required: true
          schema:
            type: string
            enum: [debug, info, warn, error]
      responses:
        '200':
          description: Log level changed successfully
        '400':
          description: Invalid log level
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
  /debug/pprof/{profile}:
    get:
      summary: Get pprof profile
      description: Returns Go pprof profiling data
      parameters:
        - name: profile
          in: path
          required: true
          schema:
            type: string
            enum: [allocs, block, cmdline, goroutine, heap, mutex, profile, threadcreate, trace]
      responses:
        '200':
          description: Pprof profile data
          content:
            application/octet-stream:
              schema:
                type: string
                format: binary
  /posts:
    get:
      summary: List posts with pagination
      description: Returns a paginated list of blog posts
      parameters:
        - name: cursor
          in: query
          description: Cursor for pagination
          schema:
            type: string
        - name: limit
          in: query
          description: Maximum number of posts to return
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: List of posts
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PostList'
    post:
      summary: Create a new blog post
      description: Creates a new blog post with the provided data
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePostRequest'
      responses:
        '201':
          description: Post created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Post'
        '400':
          description: Invalid post data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
  /posts/{id}:
    get:
      summary: Get a specific post
      description: Returns the blog post with the specified ID
      parameters:
        - name: id
          in: path
          required: true
          description: Post ID
          schema:
            type: string
            pattern: '^[a-zA-Z0-9-]+$'
      responses:
        '200':
          description: Post found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Post'
        '404':
          description: Post not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    put:
      summary: Update a blog post
      description: Updates the blog post with the specified ID
      parameters:
        - name: id
          in: path
          required: true
          description: Post ID
          schema:
            type: string
            pattern: '^[a-zA-Z0-9-]+$'
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePostRequest'
      responses:
        '200':
          description: Post updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Post'
        '400':
          description: Invalid post data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: Post not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    delete:
      summary: Delete a blog post
      description: Deletes the blog post with the specified ID
      parameters:
        - name: id
          in: path
          required: true
          description: Post ID
          schema:
            type: string
            pattern: '^[a-zA-Z0-9-]+$'
      responses:
        '204':
          description: Post deleted successfully
        '404':
          description: Post not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

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
          description: Unique post identifier
          example: "post-123"
        title:
          type: string
          description: Post title
          minLength: 1
          maxLength: 200
          example: "My First Blog Post"
        content:
          type: string
          description: Post content
          minLength: 1
          maxLength: 10000
          example: "This is the content of my first blog post..."
        createdAt:
          type: string
          format: date-time
          description: Creation timestamp
          example: "2024-01-01T12:00:00Z"
        updatedAt:
          type: string
          format: date-time
          description: Last update timestamp
          example: "2024-01-01T12:00:00Z"
    CreatePostRequest:
      type: object
      required:
        - title
        - content
      properties:
        title:
          type: string
          description: Post title
          minLength: 1
          maxLength: 200
          example: "My First Blog Post"
        content:
          type: string
          description: Post content
          minLength: 1
          maxLength: 10000
          example: "This is the content of my first blog post..."
    UpdatePostRequest:
      type: object
      required:
        - title
        - content
      properties:
        title:
          type: string
          description: Post title
          minLength: 1
          maxLength: 200
          example: "Updated Blog Post Title"
        content:
          type: string
          description: Post content
          minLength: 1
          maxLength: 10000
          example: "This is the updated content of my blog post..."
    PostList:
      type: object
      required:
        - posts
      properties:
        posts:
          type: array
          description: List of posts
          items:
            $ref: '#/components/schemas/Post'
        nextCursor:
          type: string
          description: Cursor for next page
          example: "post-456"
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: string
          description: Error code
          example: "POST_NOT_FOUND"
        message:
          type: string
          description: Error message
          example: "Post with ID 'post-123' not found"
        requestId:
          type: string
          description: Request ID for tracing
          example: "550e8400-e29b-41d4-a716-446655440000"
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
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  readTimeout: "30s"
  writeTimeout: "30s"
  idleTimeout: "60s"

logging:
  level: "info"
  format: "json"  # "text" for development, "json" for production
  output: "stdout"

storage:
  type: "memory"
  # Future: database connection details

debug:
  metrics:
    enabled: true
    path: "/debug/metrics"
  pprof:
    enabled: true
    path: "/debug/pprof"

cors:
  allowedOrigins: ["*"]
  allowedMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowedHeaders: ["Content-Type", "Authorization"]
  maxAge: 86400
```

### Loading Strategy
1. Load embedded `defaults.yaml` using `//go:embed`
2. Override with environment variables (YAML_PATH format, e.g., `SERVER_PORT=9090`)
3. Validate configuration on startup
4. Provide configuration to all layers
5. Support hot-reloading for log level changes

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
- Thread-safe map with RWMutex for concurrent access
- Key-based storage with prefix support for efficient listing
- JSON serialization for complex objects
- Proper cleanup on Close()
- Memory usage monitoring
- Configurable key expiration (future enhancement)

## Error Handling Strategy

### Domain Errors
```go
type PostNotFoundError struct{ ID string }
type InvalidPostDataError struct{ Field string, Value string }
type StorageError struct{ Err error }
type ValidationError struct{ Field string, Message string }
type PaginationError struct{ Cursor string }
```

### HTTP Error Mapping
- `PostNotFoundError` → 404 Not Found
- `InvalidPostDataError` → 400 Bad Request
- `ValidationError` → 400 Bad Request
- `PaginationError` → 400 Bad Request
- `StorageError` → 500 Internal Server Error
- Generic errors → 500 Internal Server Error

### Error Response Format
```json
{
  "code": "POST_NOT_FOUND",
  "message": "Post with ID 'post-123' not found",
  "requestId": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Logging Strategy

### Contextual Logging
- Request ID in all log entries using slog
- Structured logging with consistent field names
- Environment-specific formats (text for dev, JSON for prod)
- Log levels: debug, info, warn, error
- Runtime log level adjustment via API

### Log Fields
- `request_id`: Unique request identifier
- `method`: HTTP method
- `path`: Request path
- `status_code`: HTTP status code
- `duration_ms`: Request duration in milliseconds
- `user_agent`: Client user agent
- `remote_addr`: Client IP address
- `error`: Error details (when applicable)
- `level`: Log level

### Log Format Examples
**Development (text):**
```
2024-01-01T12:00:00.000Z INFO  request completed request_id=550e8400-e29b-41d4-a716-446655440000 method=GET path=/posts status_code=200 duration_ms=15
```

**Production (JSON):**
```json
{
  "time": "2024-01-01T12:00:00.000Z",
  "level": "INFO",
  "msg": "request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/posts",
  "status_code": 200,
  "duration_ms": 15
}
```

## Metrics Collection

### Prometheus Metrics
- `http_requests_total`: Counter by method, path, status code
- `http_request_duration_seconds`: Histogram by method, path
- `http_requests_in_flight`: Gauge of active requests
- `storage_operations_total`: Counter by operation type
- `storage_operation_duration_seconds`: Histogram by operation type
- `posts_total`: Gauge of total posts
- `log_level`: Gauge of current log level

### Metrics Endpoint
- `/debug/metrics` returns Prometheus format
- No authentication required
- Real-time metrics collection
- Configurable collection interval

## Request ID Strategy

### Generation
- UUID v4 for uniqueness and security
- Generated in middleware layer before any processing
- Added to request context for propagation

### Propagation
- All business logic receives context with request ID
- Logging includes request ID in all entries
- Response headers include `X-Request-ID`
- Error responses include request ID in JSON
- Metrics include request ID for correlation

## Graceful Shutdown

### Signal Handling
- Listen for SIGTERM and SIGINT
- Wait 5 seconds for in-flight requests to complete
- Cancel root context to stop new requests
- Close all resources (storage, HTTP server)
- Log shutdown progress with timestamps

### Resource Cleanup
- Close HTTP server gracefully (stop accepting new connections)
- Close storage connections
- Flush logs
- Exit with appropriate code (0 for graceful, 1 for error)

## Testing Strategy

### Test Coverage Requirements
- ≥80% test coverage for all functions
- Exclude trivially simple functions (getters, setters)
- Focus on business logic, middleware, and error handling
- Integration tests for HTTP handlers
- Unit tests for domain logic

### Test Types
- **Unit tests**: Business logic, configuration, storage
- **Integration tests**: HTTP handlers, middleware chain
- **Middleware tests**: Request ID, logging, error handling
- **Configuration tests**: Loading, validation, environment overrides
- **Storage tests**: CRUD operations, concurrency, error handling

### Test Structure
```
tests/
├── unit/           # Business logic tests
├── integration/    # HTTP handler tests
├── middleware/     # Middleware tests
├── fixtures/       # Test data
└── helpers/        # Test utilities
```

## Security Considerations

### Input Validation
- Request body validation using oapi-codegen validation
- Path parameter validation (ID format, length)
- Query parameter validation (limits, formats)
- Content-Type validation
- Request size limits

### Output Sanitization
- Response header sanitization
- Error message sanitization (no internal details)
- Log data sanitization (no sensitive data)
- CORS header validation

### CORS Configuration
- Configurable allowed origins
- Configurable allowed methods
- Configurable allowed headers
- Preflight request handling

## Performance Considerations

### Caching
- No caching for MVP (can be added later)
- Consider Redis for production post caching
- Consider in-memory caching for frequently accessed data

### Database
- In-memory storage for MVP
- Consider PostgreSQL for production with proper indexing
- Connection pooling for database connections

### Monitoring
- Prometheus metrics for performance monitoring
- Structured logging for debugging
- Request tracing with request ID
- Memory usage monitoring

## Deployment Considerations

### Containerization
- Multi-stage Docker build for minimal image size
- Non-root user for security
- Health check endpoint
- Graceful shutdown handling

### Configuration
- Environment variable overrides
- Config file mounting
- Secrets management (future enhancement)
- Configuration validation on startup

### Monitoring
- Prometheus metrics endpoint
- Health check endpoint (`/health`)
- Structured logging
- Request tracing

## File Structure

```
.
├── api/
│   ├── openapi.yaml      # OpenAPI specification
│   ├── gen.go           # Generated code (go:generate)
│   └── impl.go          # Handler implementations
├── cmd/
│   └── server/
│       └── main.go      # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/          # Domain models and interfaces
│   │   ├── post.go
│   │   ├── errors.go
│   │   └── store.go
│   ├── application/     # Use cases and services
│   │   ├── post_service.go
│   │   ├── debug_service.go
│   │   └── pagination.go
│   ├── infrastructure/  # Adapters and implementations
│   │   ├── memory_store.go
│   │   ├── logger.go
│   │   └── metrics.go
│   └── middleware/      # HTTP middleware
│       ├── request_id.go
│       ├── logging.go
│       ├── error_handler.go
│       └── cors.go
├── config/
│   └── defaults.yaml    # Default configuration (go:embed)
├── docs/                # Documentation
├── tests/               # Test files
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

## Dependencies

### Core Dependencies
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/oapi-codegen/oapi-codegen/v2` - OpenAPI code generation
- `github.com/oapi-codegen/nethttp-middleware` - Request validation
- `gopkg.in/yaml.v3` - YAML configuration
- `log/slog` - Structured logging (Go 1.21+)
- `github.com/google/uuid` - UUID generation

### Development Dependencies
- `github.com/stretchr/testify` - Testing utilities
- `github.com/prometheus/client_golang` - Metrics
- `golang.org/x/net/http2` - HTTP/2 support

## Implementation Phases

### Phase 1: Foundation (Week 1)
1. Project setup and dependencies
2. Configuration management with `//go:embed`
3. Basic logging setup with slog
4. Storage interface and in-memory implementation
5. Basic error handling

### Phase 2: HTTP Layer (Week 2)
1. OpenAPI specification
2. Code generation setup with oapi-codegen
3. Basic middleware stack
4. Error handling middleware
5. CORS middleware

### Phase 3: Business Logic (Week 3)
1. Domain models and interfaces
2. Application services
3. Handler implementations
4. Request/response validation
5. Pagination implementation

### Phase 4: Debug Features (Week 4)
1. Metrics collection with Prometheus
2. Pprof endpoints
3. Log level adjustment API
4. Request ID tracking
5. Health check endpoint

### Phase 5: Testing & Polish (Week 5)
1. Unit tests (≥80% coverage)
2. Integration tests
3. Graceful shutdown implementation
4. Documentation and README
5. Docker containerization

## Success Criteria

1. **Functional Requirements**: All APIs work as specified in PROMPT.md
2. **Non-Functional Requirements**: Performance, security, maintainability
3. **Test Coverage**: ≥80% coverage achieved for all non-trivial functions
4. **Documentation**: Complete API documentation and comprehensive README
5. **Deployment**: Ready for containerized deployment with health checks
6. **Architecture**: Clean Hexagonal Architecture with proper separation of concerns
7. **Error Handling**: Centralized error handling with proper HTTP status codes
8. **Logging**: Contextual logging with request ID propagation
9. **Configuration**: Flexible configuration with environment variable overrides
10. **Graceful Shutdown**: Proper resource cleanup with 5-second timeout