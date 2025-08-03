.PHONY: build test clean generate run docker-build docker-run help

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run tests with coverage
test:
	go test -v -cover ./...

# Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate code from OpenAPI spec
generate:
	go generate ./...

# Run the server
run:
	go run cmd/server/main.go

# Build Docker image
docker-build:
	docker build -t boilerplate:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 boilerplate:latest

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Run vet
vet:
	go vet ./...

# Show help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests with coverage"
	@echo "  test-coverage  - Run tests and generate coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  generate       - Generate code from OpenAPI spec"
	@echo "  run            - Run the server"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  deps           - Install dependencies"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo "  help           - Show this help"