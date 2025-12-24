.PHONY: all build run test test-coverage test-race lint clean docker-build docker-run help

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
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  clean          - Clean build artifacts"
