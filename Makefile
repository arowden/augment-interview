.PHONY: all build test lint generate generate-api clean help

# Build configuration
BINARY_NAME := server
BUILD_DIR := bin
GO := go

# Tool versions
OAPI_CODEGEN_VERSION := v2.4.1

# Default target
all: generate build test

# Build the server binary
build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run all tests
test:
	$(GO) test -v -race ./...

# Run tests with coverage
test-coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Generate all code
generate: generate-api

# Generate Go server code from OpenAPI spec
generate-api: tools
	@echo "Generating Go server code from OpenAPI spec..."
	@mkdir -p internal/http
	$(GO) run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		--config api/oapi-codegen.yaml \
		api/openapi.yaml

# Validate OpenAPI spec
validate-api:
	@echo "Validating OpenAPI spec..."
	@command -v openapi-generator-cli >/dev/null 2>&1 && \
		openapi-generator-cli validate -i api/openapi.yaml || \
		echo "Install openapi-generator-cli for validation: npm install -g @openapitools/openapi-generator-cli"

# Install/verify tools
tools:
	@echo "Ensuring oapi-codegen is available..."
	@$(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the server locally
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Display help
help:
	@echo "Available targets:"
	@echo "  all           - Generate code, build, and test"
	@echo "  build         - Build the server binary"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run golangci-lint"
	@echo "  generate      - Generate all code"
	@echo "  generate-api  - Generate Go server code from OpenAPI spec"
	@echo "  validate-api  - Validate OpenAPI spec"
	@echo "  tools         - Install required tools"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the server"
	@echo "  help          - Display this help message"
