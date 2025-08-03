# Blog API Server Boilerplate

A fully-functional web API server boilerplate built with Go, using oapi-codegen and go-chi, following Hexagonal Architecture principles.

## Features

### Core API Endpoints
- **Blog Posts CRUD**: Full CRUD operations for blog posts with cursor-based pagination
- **Debug Endpoints**: Metrics, profiling, and runtime configuration
- **Health Check**: Application health monitoring

### Technical Features
- **OpenAPI 3.0 Specification**: Auto-generated server code using oapi-codegen
- **Structured Logging**: slog-based logging with contextual request tracking
- **Metrics Collection**: Prometheus metrics for monitoring
- **Request ID Tracking**: Unique request IDs for tracing
- **Graceful Shutdown**: Proper resource cleanup with 5-second timeout
- **Configuration Management**: YAML-based with environment variable overrides
- **In-Memory Storage**: Thread-safe storage implementation
- **CORS Support**: Configurable CORS middleware
- **Error Handling**: Centralized error handling with proper HTTP status codes

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters) principles:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Ports)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Router    │  │ Middleware  │  │   Generated Code    │  │
│  │  (go-chi)   │  │   Stack     │  │  (oapi-codegen)     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                 Application Layer (Use Cases)               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │Post Service │  │Debug Service│  │   Pagination        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer (Core)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Post     │  │   Errors    │  │   Store Interface   │  │
│  │   Entity    │  │   Types     │  │                     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│              Infrastructure Layer (Adapters)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Memory    │  │   Logger    │  │     Metrics         │  │
│  │   Store     │  │  (slog)     │  │  (Prometheus)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## API Endpoints

### Blog Posts
- `GET /posts` - List posts with pagination
- `POST /posts` - Create a new post
- `GET /posts/{id}` - Get a specific post
- `PUT /posts/{id}` - Update a post
- `DELETE /posts/{id}` - Delete a post

### Debug Endpoints
- `GET /debug/metrics` - Prometheus metrics
- `POST /debug/logs?level=debug` - Adjust log level at runtime
- `GET /debug/pprof/*` - Go pprof profiling endpoints
- `GET /health` - Health check

## Quick Start

### Prerequisites
- Go 1.21 or later
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd boilerplate
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Build the application**
   ```bash
   go build ./cmd/server
   ```

4. **Run the server**
   ```bash
   ./server
   ```

The server will start on `http://localhost:8080` by default.

### Using Docker

```bash
# Build the image
docker build -t blog-api-server .

# Run the container
docker run -p 8080:8080 blog-api-server
```

## Configuration

The application uses a YAML-based configuration system with environment variable overrides.

### Default Configuration (`config/defaults.yaml`)

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

### Environment Variables

You can override any configuration value using environment variables in the format `SECTION_SUBSECTION_KEY`:

```bash
export SERVER_PORT=9090
export LOGGING_LEVEL=debug
export LOGGING_FORMAT=text
```

## API Usage Examples

### Create a Post

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post..."
  }'
```

### List Posts

```bash
curl "http://localhost:8080/posts?limit=10"
```

### Get a Specific Post

```bash
curl http://localhost:8080/posts/post-123
```

### Update a Post

```bash
curl -X PUT http://localhost:8080/posts/post-123 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Blog Post Title",
    "content": "This is the updated content..."
  }'
```

### Delete a Post

```bash
curl -X DELETE http://localhost:8080/posts/post-123
```

### Adjust Log Level

```bash
curl -X POST "http://localhost:8080/debug/logs?level=debug"
```

### Get Metrics

```bash
curl http://localhost:8080/debug/metrics
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Development

### Project Structure

```
.
├── api/                    # OpenAPI specification and generated code
│   ├── openapi.yaml       # API specification
│   ├── handlers.go        # HTTP handler implementations
│   └── oapi-codegen.yaml  # Code generation configuration
├── cmd/
│   └── server/
│       └── main.go        # Application entry point
├── internal/
│   ├── application/       # Use cases and services
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models and interfaces
│   ├── infrastructure/   # Adapters and implementations
│   └── middleware/       # HTTP middleware
├── config/
│   └── defaults.yaml     # Default configuration
└── docs/                 # Documentation
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/config

# Run tests with verbose output
go test -v ./...
```

### Code Generation

The API code is generated from the OpenAPI specification:

```bash
# Generate code from OpenAPI spec
go generate ./api
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Check for race conditions
go test -race ./...
```

## Monitoring and Observability

### Metrics

The application exposes Prometheus metrics at `/debug/metrics`:

- **HTTP Metrics**: Request counts, durations, and in-flight requests
- **Storage Metrics**: Operation counts and durations
- **Business Metrics**: Post counts and application state
- **Application Metrics**: Log level and system health

### Logging

The application uses structured logging with the following features:

- **Request ID Tracking**: Every request gets a unique ID
- **Contextual Logging**: Request context is included in all log entries
- **Environment-Specific Formats**: JSON for production, text for development
- **Runtime Level Adjustment**: Change log level without restart

### Health Checks

The `/health` endpoint provides:

- Overall application status
- Storage health check
- Logger availability check
- Dependency status

## Error Handling

The application uses a centralized error handling system:

### Error Types

- `PostNotFoundError` → 404 Not Found
- `InvalidPostDataError` → 400 Bad Request
- `ValidationError` → 400 Bad Request
- `PaginationError` → 400 Bad Request
- `StorageError` → 500 Internal Server Error

### Error Response Format

```json
{
  "code": "POST_NOT_FOUND",
  "message": "Post with ID 'post-123' not found",
  "requestId": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Performance Considerations

### Storage

- **In-Memory Storage**: Fast access for development and testing
- **Thread-Safe Operations**: Concurrent access support
- **JSON Serialization**: Efficient data storage and retrieval

### HTTP Server

- **Graceful Shutdown**: 5-second timeout for in-flight requests
- **Connection Timeouts**: Configurable read/write/idle timeouts
- **Middleware Optimization**: Efficient middleware chain

### Monitoring

- **Request Metrics**: Track performance and identify bottlenecks
- **Storage Metrics**: Monitor storage operation performance
- **Business Metrics**: Track application usage patterns

## Security

### Input Validation

- **Request Body Validation**: JSON schema validation
- **Path Parameter Validation**: ID format and length checks
- **Query Parameter Validation**: Limits and format validation

### Output Sanitization

- **Response Headers**: Sanitized HTTP headers
- **Error Messages**: No internal details exposed
- **Log Data**: Sensitive data filtering

### CORS Configuration

- **Configurable Origins**: Control allowed origins
- **Method Restrictions**: Limit HTTP methods
- **Header Validation**: Validate allowed headers

## Deployment

### Production Considerations

1. **Configuration**: Use environment variables for production settings
2. **Logging**: Set log format to JSON and appropriate level
3. **Metrics**: Configure Prometheus scraping
4. **Health Checks**: Set up monitoring for the health endpoint
5. **Storage**: Consider replacing in-memory storage with a database
6. **Security**: Review and configure CORS settings

### Docker Deployment

```dockerfile
# Multi-stage build for minimal image size
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) for OpenAPI code generation
- [go-chi](https://github.com/go-chi/chi) for HTTP routing
- [Prometheus](https://prometheus.io/) for metrics collection
- [slog](https://pkg.go.dev/log/slog) for structured logging