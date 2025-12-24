# Swagger/OpenAPI Documentation

This document explains how to generate and use the Swagger/OpenAPI documentation for the FizzBuzz REST API.

## Table of Contents

- [Installation](#installation)
- [Generating Documentation](#generating-documentation)
- [Viewing Documentation](#viewing-documentation)
- [API Specification](#api-specification)
- [Annotation Reference](#annotation-reference)

## Installation

### Install go-swagger CLI

The project uses [go-swagger](https://github.com/go-swagger/go-swagger) for generating OpenAPI 2.0 specifications from code annotations.

```bash
# Install the swagger CLI tool
go install github.com/go-swagger/go-swagger/cmd/swagger@latest

# Verify installation
swagger version
```

### Alternative: Using Docker

If you prefer not to install the CLI globally:

```bash
docker pull quay.io/goswagger/swagger

# Use it with an alias
alias swagger="docker run --rm -it -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
```

## Generating Documentation

### Generate Swagger Spec

Generate both JSON and YAML formats:

```bash
make swagger
```

This will create:
- `docs/swagger.json` - OpenAPI 2.0 spec in JSON format
- `docs/swagger.yaml` - OpenAPI 2.0 spec in YAML format

### Manual Generation

```bash
# Generate JSON
swagger generate spec -o ./docs/swagger.json --scan-models

# Generate YAML
swagger generate spec -o ./docs/swagger.yaml --scan-models
```

### Validate Spec

Ensure your generated spec is valid:

```bash
make swagger-validate
```

## Viewing Documentation

### Option 1: Swagger UI (Recommended)

Serve an interactive Swagger UI:

```bash
make swagger-serve
```

Then open http://localhost:8081/docs in your browser.

Features:
- **Interactive API explorer** - Try out API endpoints directly from the UI
- **Request/Response examples** - See example payloads
- **Schema documentation** - Detailed model definitions

### Option 2: Online Swagger Editor

1. Generate the spec: `make swagger`
2. Go to [editor.swagger.io](https://editor.swagger.io)
3. Upload `docs/swagger.yaml` or copy/paste its contents



## API Specification

### Overview

The FizzBuzz API provides three main endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/fizzbuzz` | POST | Generate a FizzBuzz sequence |
| `/statistics` | GET | Get the most frequent request |
| `/health` | GET | Health check |

### Example Usage

#### Generate FizzBuzz

```bash
curl -X POST http://localhost:8080/fizzbuzz \
  -H "Content-Type: application/json" \
  -d '{
    "int1": 3,
    "int2": 5,
    "limit": 15,
    "str1": "fizz",
    "str2": "buzz"
  }'
```

Response (200 OK):
```json
{
  "result": [
    "1", "2", "fizz", "4", "buzz",
    "fizz", "7", "8", "fizz", "buzz",
    "11", "fizz", "13", "14", "fizzbuzz"
  ]
}
```

Validation Error (400 Bad Request):
```json
{
  "error": "invalid parameters",
  "details": [
    "int1 must be greater than 0",
    "str2 cannot be empty"
  ]
}
```

#### Get Statistics

```bash
curl http://localhost:8080/statistics
```

Response (200 OK):
```json
{
  "most_frequent_request": {
    "int1": 3,
    "int2": 5,
    "limit": 15,
    "str1": "fizz",
    "str2": "buzz"
  },
  "hits": 42
}
```

#### Health Check

```bash
curl http://localhost:8080/health
```

Response (200 OK):
```json
{
  "status": "healthy"
}
```

## Annotation Reference

### Supported Annotations

The project uses go-swagger's declarative comments for API documentation.

#### Package-level (docs/swagger.go)

```go
// Package classification FizzBuzz REST API
//
// Documentation for FizzBuzz API
//
//     Schemes: http, https
//     Host: localhost:8080
//     BasePath: /
//     Version: 1.0.0
//
// swagger:meta
package docs
```

#### Route Documentation

```go
// swagger:route POST /fizzbuzz fizzbuzz generateFizzBuzz
//
// Generate FizzBuzz Sequence
//
// Long description of the endpoint...
//
// Responses:
//   200: generateResponse
//   400: errorResponse
//   500: errorResponse
func (h *FizzBuzzHandler) Generate(w http.ResponseWriter, r *http.Request) {
```

#### Model Documentation

```go
// GenerateRequest represents input for FizzBuzz generation
// swagger:model
type generateRequest struct {
    // First divisor (must be > 0)
    // required: true
    // example: 3
    Int1 int `json:"int1"`
}
```

#### Parameter Documentation

```go
// swagger:parameters generateFizzBuzz
type generateFizzBuzzParams struct {
    // FizzBuzz generation parameters
    // in: body
    // required: true
    Body generateRequest
}
```

#### Response Documentation

```go
// swagger:response generateResponse
type generateResponseWrapper struct {
    // in: body
    Body generateResponse
}
```

### Common Annotation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `swagger:meta` | Package-level metadata | Applied to docs package |
| `swagger:route` | Endpoint definition | Applied to handler methods |
| `swagger:model` | Data model | Applied to structs |
| `swagger:parameters` | Request parameters | Wrapper type for params |
| `swagger:response` | Response definition | Wrapper type for responses |
| `required` | Mark field as required | `// required: true` |
| `example` | Example value | `// example: fizz` |
| `minimum` | Minimum value | `// minimum: 1` |
| `maximum` | Maximum value | `// maximum: 10000` |

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Generate Swagger Docs

on:
  push:
    branches: [main]

jobs:
  swagger:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install swagger
        run: go install github.com/go-swagger/go-swagger/cmd/swagger@latest
      
      - name: Generate Swagger
        run: make swagger
      
      - name: Validate Swagger
        run: make swagger-validate
      
      - name: Upload Swagger Spec
        uses: actions/upload-artifact@v3
        with:
          name: swagger-spec
          path: docs/swagger.yaml
```

## Best Practices

### 1. Keep Documentation Up-to-Date

Always regenerate Swagger docs after API changes:

```bash
# After modifying handlers
make swagger swagger-validate
```

### 2. Validate Before Committing

Add a pre-commit hook:

```bash
#!/bin/bash
# .git/hooks/pre-commit
make swagger-validate || exit 1
```

### 3. Version Your API

Update version in `docs/swagger.go`:

```go
//     Version: 2.0.0
```

### 4. Provide Examples

Always include examples for request/response models:

```go
// example: ["1", "2", "fizz"]
Result []string `json:"result"`
```

### 5. Document Validation Rules

Use annotations to document constraints:

```go
// First divisor (must be > 0)
// required: true
// minimum: 1
// example: 3
Int1 int `json:"int1"`
```

## Troubleshooting

### Issue: "swagger: command not found"

**Solution**: Install the swagger CLI:
```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

### Issue: "WARNING: No operations defined"

**Solution**: Ensure you have:
1. Added `swagger:route` annotations to handlers
2. Defined response types with `swagger:response`
3. Run `make swagger` to regenerate

### Issue: Models not showing in spec

**Solution**: Use `swagger:model` annotation:
```go
// swagger:model
type MyModel struct { ... }
```

### Issue: Invalid spec validation errors

**Solution**: 
1. Check for circular references in models
2. Ensure all referenced types are defined
3. Use `make swagger-validate` to see specific errors

## Additional Resources

- [go-swagger Documentation](https://goswagger.io/)
- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
- [Swagger Editor](https://editor.swagger.io/)
- [Example Projects](https://github.com/go-swagger/go-swagger/tree/master/examples)

## Migration to OpenAPI 3.0

While go-swagger generates OpenAPI 2.0, you can convert to 3.0:

```bash
# Using swagger-cli
npm install -g @apidevtools/swagger-cli
swagger-cli bundle -o docs/openapi-v3.yaml -t yaml docs/swagger.yaml
```

Or use online converters:
- [Swagger Converter](https://converter.swagger.io/)
- [Mermade OpenAPI Converter](https://mermade.org.uk/openapi-converter)
