.PHONY: all build run test test-coverage test-race lint clean docker-build docker-run swagger swagger-serve help

# Default target
all: test build

# Build the binary
build:
	go build -ldflags="-w -s" -o bin/fizzbuzz-service cmd/server/main.go

# Run the server locally
run:
	go run cmd/server/main.go

# Run all tests
test:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race -v ./...

# Run tests with coverage report
test-coverage:
	go test -race -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total

# Run unit tests only
test-unit:
	go test -v ./test/unit/...

# Run integration tests only
test-integration:
	go test -v ./test/integration/...

# Run e2e tests only
test-e2e:
	go test -v ./test/e2e/...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Tidy dependencies
tidy:
	go mod tidy

# Generate Swagger documentation (requires swagger CLI)
swagger:
	@echo "Generating Swagger documentation..."
	@command -v swagger >/dev/null 2>&1 || { echo "swagger CLI not found. Install with: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; exit 1; }
	swagger generate spec -o ./docs/swagger.json --scan-models
	@echo "Swagger spec generated at ./docs/swagger.json"
	@echo "Converting to YAML..."
	swagger generate spec -o ./docs/swagger.yaml --scan-models
	@echo "Swagger spec generated at ./docs/swagger.yaml"

# Serve Swagger UI (requires swagger CLI)
swagger-serve:
	@command -v swagger >/dev/null 2>&1 || { echo "swagger CLI not found. Install with: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; exit 1; }
	@echo "Serving Swagger UI at http://localhost:8081/docs"
	@echo "Press Ctrl+C to stop"
	swagger serve -F=swagger --port=8081 ./docs/swagger.yaml

# Validate Swagger spec
swagger-validate:
	@command -v swagger >/dev/null 2>&1 || { echo "swagger CLI not found. Install with: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; exit 1; }
	swagger validate ./docs/swagger.yaml

# Build Docker image
docker-build:
	docker build -t fizzbuzz-service:latest .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker containers
docker-stop:
	docker-compose down

# Clean build artifacts
clean:
	go clean
	rm -f bin/fizzbuzz-service
	rm -f coverage.out coverage.html
	rm -f docs/swagger.json docs/swagger.yaml
	docker rmi fizzbuzz-service:latest 2>/dev/null || true

# Show help
help:
	@echo "Available targets:"
	@echo "  all            - Run tests and build"
	@echo "  build          - Build the binary"
	@echo "  run            - Run the server locally"
	@echo "  test           - Run all tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-e2e       - Run e2e tests only"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  swagger-serve  - Serve Swagger UI"
	@echo "  swagger-validate - Validate Swagger spec"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  clean          - Clean build artifacts"
