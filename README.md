# Go Boilerplate by AI

This project provides a starter web API server written in Go using [chi](https://github.com/go-chi/chi) and code generated from an OpenAPI specification via [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Features

- Configuration loader with embedded defaults and environment overrides
- Structured logging with `slog` and request IDs
- In-memory store implementing a simple key/value interface
- Generated REST API for blog post CRUD operations
- Debug endpoints for metrics, pprof, and recent logs
- Graceful shutdown on SIGTERM

## Development

Generate server code from the OpenAPI spec:

```sh
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
```

Run tests:

```sh
go test -timeout 30s -cover ./...
```

Start the server:

```sh
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080` by default.
