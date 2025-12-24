# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies first (for caching)
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o fizzbuzz-service cmd/server/main.go

# Final stage - minimal image
FROM scratch

# Copy CA certificates for HTTPS calls (if needed in future)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/fizzbuzz-service /fizzbuzz-service

# Expose port
EXPOSE 8080

# Run as non-root user for security
USER 65534:65534

ENTRYPOINT ["/fizzbuzz-service"]
