# Architecture

This project is a Go web API server generated from an OpenAPI specification using `oapi-codegen` and built on top of the `chi` router. It follows Hexagonal Architecture principles.

## Packages

- `cmd/server`: entrypoint that wires configuration, logging, HTTP server, and graceful shutdown.
- `pkg/config`: loads configuration from an embedded `defaults.yaml` and overrides values from environment variables.
- `pkg/log`: initializes slog logger with contextual request ID support.
- `pkg/middleware`: request ID generation, logging, centralized error handling, and request validation.
- `pkg/store`: in-memory implementation of `Store` interface for persistence.
- `pkg/handler`: business logic for blog posts and debug endpoints generated via `oapi-codegen`'s strict server mode.

## Data Flow

1. `cmd/server` loads configuration and initializes dependencies.
2. Router attaches middleware for request IDs, logging, validation, and error handling.
3. Handlers interact with the `Store` interface for persistence.
4. Logging uses slog with request ID attached to the context.
5. Graceful shutdown waits for 5 seconds on SIGTERM before closing resources.

## Testing

- Each package maintains ≥ 80% test coverage except for trivial functions.
- `go test -timeout` is used to avoid hanging tests.
