## Context
Manual OpenTelemetry instrumentation requires wrapping every HTTP handler, database call, and adding spans throughout the codebase. Orchestrion provides compile-time instrumentation that automatically injects tracing, metrics, and structured logging without code modifications.

## Goals / Non-Goals
- Goals: Zero-code instrumentation, consistent telemetry, reproducible builds, small secure container images
- Non-Goals: Runtime instrumentation (eBPF), custom span attributes, vendor-specific exporters

## Decisions
- Decision: Use Datadog Orchestrion for compile-time instrumentation with pinned version
- Alternatives considered: Manual instrumentation (verbose), eBPF auto-instrumentation (complex, kernel dependency), go-instrument (less mature)

- Decision: Multi-stage Dockerfile with distroless base and SHA256 digest
- Alternatives considered: Alpine (larger attack surface), scratch (no CA certs), debian-slim (larger)

- Decision: Embed migrations in binary (not copy to image)
- Alternatives considered: Separate migration container (more complexity), external migration tool (another dependency)

- Decision: Pin all versions (Go, Orchestrion, base images) for reproducibility
- Alternatives considered: Using :latest tags (non-reproducible builds)

## Orchestrion Limitations
**IMPORTANT: Orchestrion has limitations that require manual instrumentation in some cases.**

Orchestrion automatically instruments:
- net/http handlers and clients
- github.com/jackc/pgx/v5 queries
- log/slog calls

Orchestrion does NOT instrument:
- Custom libraries or internal packages
- Prepared statements (some patterns)
- Batch database operations
- gRPC calls (use separate instrumentation)
- Custom business logic spans

**When to add manual instrumentation:**
1. Business operations (fund creation, transfers) - add explicit spans for domain visibility
2. Error context - add error details to spans manually
3. Custom attributes - add fund_id, owner_name to spans
4. Cross-service calls - ensure trace propagation headers

```go
func (s *Service) CreateFund(ctx context.Context, req CreateFundRequest) (*Fund, error) {
    ctx, span := s.tracer.Start(ctx, "CreateFund")
    defer span.End()

    span.SetAttributes(
        attribute.String("fund.name", req.Name),
        attribute.Int("fund.total_units", req.TotalUnits),
    )

    // ... implementation
}
```

## Orchestrion Configuration
```yaml
# .orchestrion.yml
service-name: augment-fund-api

integrations:
  - name: net/http
    enabled: true
  - name: database/sql
    enabled: true
  - name: github.com/jackc/pgx/v5
    enabled: true
  - name: log/slog
    enabled: true

telemetry:
  tracer: otel
  metrics: otel
  logs: otel
```

## Dockerfile Structure
```dockerfile
# Stage 1: Build with instrumentation
FROM golang:1.22-bookworm AS builder

# Pin orchestrion version for reproducibility
RUN go install github.com/DataDog/orchestrion@v0.18.0

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build with instrumentation and stripped debug info
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod \
    orchestrion go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o /app/server ./cmd/server

# Stage 2: Minimal secure runtime
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary with correct ownership
COPY --from=builder --chown=nonroot:nonroot /app/server /server

# Note: migrations are embedded in binary, not copied
USER nonroot:nonroot
WORKDIR /home/nonroot

EXPOSE 8080

# Healthcheck compatible with distroless (no shell)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD ["/server", "--health-check"]

ENTRYPOINT ["/server"]
```

## .dockerignore
```
.git/
.gitignore
*_test.go
testdata/
test_data/
.env*
docker-compose.yml
Dockerfile
.dockerignore
*.md
!README.md
bin/
coverage.out
```

## Makefile Targets
```makefile
.PHONY: build build-instrumented docker-build generate test lint scan-vulnerabilities sbom verify-size

VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -ldflags="-s -w" -o bin/server ./cmd/server

build-instrumented:
	orchestrion go build -ldflags="-s -w -X main.Version=$(VERSION)" -o bin/server ./cmd/server

docker-build: verify-size
	docker build --build-arg VERSION=$(VERSION) -t augment-fund-api:$(VERSION) -t augment-fund-api:latest .

verify-size:
	@if [ -n "$$(docker images -q augment-fund-api:latest)" ]; then \
		SIZE=$$(docker images augment-fund-api:latest --format "{{.Size}}"); \
		echo "Image size: $$SIZE"; \
	fi

generate:
	go generate ./...
	oapi-codegen -generate types,server -package http api/openapi.yaml > internal/http/openapi.gen.go

test:
	go test -race -cover -coverprofile=coverage.out ./...

lint:
	golangci-lint run

scan-vulnerabilities:
	trivy image augment-fund-api:latest --severity HIGH,CRITICAL

sbom:
	syft augment-fund-api:latest -o json > sbom.json

verify-binary:
	@file bin/server
	@ldd bin/server 2>&1 || echo "Statically linked (good)"
```

## Build Artifacts
```
bin/
  server              # Instrumented binary (~15MB, stripped)
```

## Risks / Trade-offs
- Orchestrion adds build-time dependency → Pinned version v0.18.0
- Instrumented binary slightly larger → Acceptable (~2MB overhead)
- Compile time increases ~20% → Acceptable for production builds
- Auto-instrumentation misses business context → Document when to add manual spans

## Open Questions
- None
