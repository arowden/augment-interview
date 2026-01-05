.PHONY: all build test test-unit test-integration test-all test-coverage lint generate generate-api clean help \
        docker-build docker-push verify-size scan-vulnerabilities sbom verify-binary run \
        deploy deploy-api deploy-frontend ecr-login

# Build configuration
BINARY_NAME := server
BUILD_DIR := bin
GO := go

# Tool versions
OAPI_CODEGEN_VERSION := v2.4.1

# Version from git
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# AWS configuration (override with environment variables)
AWS_REGION ?= us-west-2
ENVIRONMENT ?= dev
ECR_REPO ?= augment-fund-$(ENVIRONMENT)-api
ECS_CLUSTER ?= augment-fund-$(ENVIRONMENT)
ECS_SERVICE ?= augment-fund-$(ENVIRONMENT)-api
S3_BUCKET ?= $(shell cd deploy/terraform && terraform output -raw frontend_bucket_name 2>/dev/null || echo "augment-fund-$(ENVIRONMENT)-frontend")

# Docker configuration
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text 2>/dev/null)
ECR_REGISTRY ?= $(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,)
IMAGE_NAME ?= $(if $(ECR_REGISTRY),$(ECR_REGISTRY)/$(ECR_REPO),augment-fund-api)
IMAGE_TAG ?= $(VERSION)

# Default target
all: generate build test

# Build the server binary (local, no instrumentation)
build:
	$(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Build Docker image (linux/amd64 for ECS Fargate)
docker-build:
	docker build --platform linux/amd64 --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(IMAGE_TAG) -t $(IMAGE_NAME):latest .
	@$(MAKE) verify-size

# Login to ECR
ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)

# Push Docker image to ECR
docker-push: ecr-login
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

# Deploy API to ECS (build, push, update service)
deploy-api: docker-build docker-push
	@echo "Updating ECS service..."
	aws ecs update-service \
		--cluster $(ECS_CLUSTER) \
		--service $(ECS_SERVICE) \
		--force-new-deployment \
		--region $(AWS_REGION)
	@echo "Waiting for deployment to stabilize..."
	aws ecs wait services-stable \
		--cluster $(ECS_CLUSTER) \
		--services $(ECS_SERVICE) \
		--region $(AWS_REGION)
	@echo "API deployed successfully!"

# Deploy frontend to S3
deploy-frontend:
	@echo "Building frontend..."
	cd frontend && npm ci && npm run build
	@echo "Deploying to S3..."
	aws s3 sync frontend/dist/ s3://$(S3_BUCKET) \
		--delete \
		--cache-control "public,max-age=31536000,immutable" \
		--exclude "index.html" \
		--region $(AWS_REGION)
	aws s3 cp frontend/dist/index.html s3://$(S3_BUCKET)/index.html \
		--cache-control "no-cache" \
		--region $(AWS_REGION)
	@echo "Frontend deployed successfully!"

# Deploy everything
deploy: deploy-api deploy-frontend

# Verify image size (<50MB target)
verify-size:
	@if [ -n "$$(docker images -q $(IMAGE_NAME):latest 2>/dev/null)" ]; then \
		SIZE=$$(docker images $(IMAGE_NAME):latest --format "{{.Size}}"); \
		echo "Image size: $$SIZE"; \
	fi

# Scan for vulnerabilities
scan-vulnerabilities:
	@command -v trivy >/dev/null 2>&1 && \
		trivy image $(IMAGE_NAME):latest --severity HIGH,CRITICAL || \
		echo "Install trivy for vulnerability scanning: brew install trivy"

# Generate software bill of materials
sbom:
	@command -v syft >/dev/null 2>&1 && \
		syft $(IMAGE_NAME):latest -o json > sbom.json || \
		echo "Install syft for SBOM generation: brew install syft"

# Verify binary is statically linked
verify-binary:
	@file $(BUILD_DIR)/$(BINARY_NAME)
	@ldd $(BUILD_DIR)/$(BINARY_NAME) 2>&1 || echo "Statically linked (good)"

# Run unit tests only (fast, no Docker required)
test:
	$(GO) test -v -race ./...

# Alias for test
test-unit: test

# Run integration tests (requires Docker)
test-integration:
	$(GO) test -v -race -tags=integration ./...

# Run all tests (unit + integration)
test-all:
	$(GO) test -v -race -tags=integration ./...

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
	rm -f coverage.out coverage.html sbom.json

# Run the server locally
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Display help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build:"
	@echo "  build                - Build the server binary (local)"
	@echo "  docker-build         - Build Docker image"
	@echo ""
	@echo "Deploy:"
	@echo "  deploy               - Deploy API and frontend"
	@echo "  deploy-api           - Build, push to ECR, update ECS"
	@echo "  deploy-frontend      - Build and sync to S3"
	@echo "  ecr-login            - Login to ECR"
	@echo "  docker-push          - Push image to ECR"
	@echo ""
	@echo "Test & Quality:"
	@echo "  test                 - Run unit tests (fast, no Docker)"
	@echo "  test-integration     - Run integration tests (requires Docker)"
	@echo "  test-all             - Run all tests (unit + integration)"
	@echo "  test-coverage        - Run tests with coverage report"
	@echo "  lint                 - Run golangci-lint"
	@echo "  scan-vulnerabilities - Scan image for vulnerabilities"
	@echo "  sbom                 - Generate software bill of materials"
	@echo ""
	@echo "Other:"
	@echo "  generate             - Generate all code"
	@echo "  generate-api         - Generate Go server code from OpenAPI spec"
	@echo "  tools                - Install required tools"
	@echo "  clean                - Remove build artifacts"
	@echo "  run                  - Build and run the server"
	@echo ""
	@echo "Configuration:"
	@echo "  AWS_REGION=$(AWS_REGION)"
	@echo "  ENVIRONMENT=$(ENVIRONMENT)"
	@echo "  VERSION=$(VERSION)"
	@echo "  IMAGE_NAME=$(IMAGE_NAME)"
