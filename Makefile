.PHONY: all build build-instrumented test test-coverage lint generate generate-api clean help \
        docker-build docker-push verify-size scan-vulnerabilities sbom verify-binary run

# Build configuration
BINARY_NAME := server
BUILD_DIR := bin
GO := go

# Tool versions
OAPI_CODEGEN_VERSION := v2.4.1
ORCHESTRION_VERSION := v0.18.0

# Version from git
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Docker configuration
DOCKER_REGISTRY ?=
IMAGE_NAME ?= augment-fund-api
IMAGE_TAG ?= $(VERSION)
FULL_IMAGE_NAME := $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/$(IMAGE_NAME),$(IMAGE_NAME))

# Default target
all: generate build test

# Build the server binary (local, no instrumentation)
build:
	$(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Build with Orchestrion instrumentation
build-instrumented:
	orchestrion go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Build Docker image
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(FULL_IMAGE_NAME):$(IMAGE_TAG) -t $(FULL_IMAGE_NAME):latest .
	@$(MAKE) verify-size

# Push Docker image to registry
docker-push:
	docker push $(FULL_IMAGE_NAME):$(IMAGE_TAG)
	docker push $(FULL_IMAGE_NAME):latest

# Verify image size (<50MB target)
verify-size:
	@if [ -n "$$(docker images -q $(FULL_IMAGE_NAME):latest 2>/dev/null)" ]; then \
		SIZE=$$(docker images $(FULL_IMAGE_NAME):latest --format "{{.Size}}"); \
		echo "Image size: $$SIZE"; \
	fi

# Scan for vulnerabilities
scan-vulnerabilities:
	@command -v trivy >/dev/null 2>&1 && \
		trivy image $(FULL_IMAGE_NAME):latest --severity HIGH,CRITICAL || \
		echo "Install trivy for vulnerability scanning: brew install trivy"

# Generate software bill of materials
sbom:
	@command -v syft >/dev/null 2>&1 && \
		syft $(FULL_IMAGE_NAME):latest -o json > sbom.json || \
		echo "Install syft for SBOM generation: brew install syft"

# Verify binary is statically linked
verify-binary:
	@file $(BUILD_DIR)/$(BINARY_NAME)
	@ldd $(BUILD_DIR)/$(BINARY_NAME) 2>&1 || echo "Statically linked (good)"

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

# Install orchestrion
install-orchestrion:
	$(GO) install github.com/DataDog/orchestrion@$(ORCHESTRION_VERSION)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html sbom.json

# Run the server locally
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Display help
help:
	@echo "Available targets:"
	@echo "  all                  - Generate code, build, and test"
	@echo "  build                - Build the server binary (local)"
	@echo "  build-instrumented   - Build with Orchestrion instrumentation"
	@echo "  docker-build         - Build Docker image"
	@echo "  docker-push          - Push Docker image to registry"
	@echo "  verify-size          - Show Docker image size"
	@echo "  scan-vulnerabilities - Scan image for vulnerabilities"
	@echo "  sbom                 - Generate software bill of materials"
	@echo "  verify-binary        - Check if binary is statically linked"
	@echo "  test                 - Run all tests"
	@echo "  test-coverage        - Run tests with coverage report"
	@echo "  lint                 - Run golangci-lint"
	@echo "  generate             - Generate all code"
	@echo "  generate-api         - Generate Go server code from OpenAPI spec"
	@echo "  validate-api         - Validate OpenAPI spec"
	@echo "  tools                - Install required tools"
	@echo "  install-orchestrion  - Install Orchestrion"
	@echo "  clean                - Remove build artifacts"
	@echo "  run                  - Build and run the server"
	@echo "  help                 - Display this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  DOCKER_REGISTRY=$(DOCKER_REGISTRY)"
	@echo "  IMAGE_NAME=$(IMAGE_NAME)"
