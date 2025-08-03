# Architecture v2

This version refines the design to explicitly cover all requirements from `docs/PROMPT.md`.

## Components

### cmd/server
- Initializes configuration (`pkg/config`), logger (`pkg/log`), in-memory store (`pkg/store`), and HTTP router.
- Registers generated handlers from `oapi-codegen` in strict-server mode.
- Exposes debug endpoints:
  - `/debug/metrics` from `promhttp`.
  - `/debug/logs` for runtime log level adjustment.
  - `/debug/pprof/*` from `net/http/pprof`.
- Handles graceful shutdown on SIGTERM with a 5 second timeout.

### pkg/config
- Embeds `defaults.yaml` using `//go:embed`.
- Loads defaults and overrides them with environment variables.
- Exposes strongly typed configuration struct.

### pkg/log
- Provides `slog.Logger` configured for text (development) or JSON (production).
- Attaches request ID from context to every log entry.

### pkg/middleware
- `RequestID`: generates UUID and stores in context & `X-Request-ID` response header.
- `Logger`: logs incoming requests and responses with duration and status code.
- `ErrorHandler`: converts business errors into HTTP responses.
- `Validator`: runs oapi-codegen request validation before handler execution.

### pkg/store
- Defines `Store` interface and an in-memory map-based implementation.

### pkg/handler
- Generated from an OpenAPI spec describing blog post CRUD operations.
- Uses strict-server interface to ensure compile time handler completeness.

## Testing Strategy
- Each package includes unit tests achieving ≥80% coverage.
- Integration tests cover HTTP endpoints using `httptest`.
- `go test -timeout 30s ./...` and `go build ./...` run in CI and locally.

