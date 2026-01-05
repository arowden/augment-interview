## ADDED Requirements

### Requirement: Orchestrion Toolchain Integration
The build system SHALL use Datadog Orchestrion for compile-time automatic instrumentation of OpenTelemetry tracing, metrics, and logging with a pinned version.

#### Scenario: Orchestrion installation with pinned version
- **WHEN** the builder Docker stage runs
- **THEN** orchestrion is installed via `go install github.com/DataDog/orchestrion@v0.18.0`

#### Scenario: Instrumented compilation
- **WHEN** `orchestrion go build` is executed
- **THEN** the resulting binary has automatic tracing instrumentation injected

#### Scenario: HTTP instrumentation
- **WHEN** the instrumented binary handles HTTP requests
- **THEN** spans are automatically created with http.method, http.route, http.status_code

#### Scenario: Database instrumentation
- **WHEN** the instrumented binary executes pgx queries
- **THEN** spans are automatically created with db.statement, db.operation

#### Scenario: Log correlation
- **WHEN** the instrumented binary logs via slog
- **THEN** trace_id and span_id are automatically injected into log records

#### Scenario: Instrumentation verification
- **WHEN** the build completes
- **THEN** the binary contains orchestrion instrumentation markers

### Requirement: Orchestrion Limitations Documentation
The build system SHALL document known limitations of automatic instrumentation.

#### Scenario: Limitations documented
- **WHEN** design.md is examined
- **THEN** it lists: unsupported libraries, missing custom span attributes, prepared statement limitations

#### Scenario: Manual instrumentation guidance
- **WHEN** design.md is examined
- **THEN** it specifies when manual instrumentation is needed (business operations, error context)

### Requirement: Orchestrion Configuration File
The project SHALL include a `.orchestrion.yml` configuration file specifying instrumentation settings.

#### Scenario: Configuration file exists
- **WHEN** the project root is examined
- **THEN** `.orchestrion.yml` exists with service-name and integrations

#### Scenario: Service name configuration
- **WHEN** `.orchestrion.yml` is parsed
- **THEN** service-name is set to "augment-fund-api"

#### Scenario: Integration enablement
- **WHEN** `.orchestrion.yml` is parsed
- **THEN** net/http, pgx, and slog integrations are enabled

#### Scenario: Telemetry backend configuration
- **WHEN** `.orchestrion.yml` is parsed
- **THEN** tracer, metrics, and logs are configured to use otel

### Requirement: Multi-Stage Dockerfile with Security
The project SHALL include a multi-stage Dockerfile that builds an instrumented binary and creates a minimal, secure runtime image.

#### Scenario: Dockerfile exists
- **WHEN** the project root is examined
- **THEN** a Dockerfile exists

#### Scenario: Builder stage with pinned Go version
- **WHEN** the Dockerfile is parsed
- **THEN** a builder stage based on golang:1.22-bookworm (not :latest) exists

#### Scenario: Runtime stage with image digest
- **WHEN** the Dockerfile is parsed
- **THEN** a runtime stage uses gcr.io/distroless/static-debian12:nonroot with SHA256 digest

#### Scenario: Non-root execution enforced
- **WHEN** the container runs
- **THEN** USER nonroot:nonroot is set and process runs as non-root

#### Scenario: Stripped binary
- **WHEN** the build completes
- **THEN** the binary is built with -ldflags="-s -w" to strip debug info

### Requirement: Builder Stage Configuration
The builder stage SHALL install orchestrion with pinned version and compile the application with instrumentation.

#### Scenario: Orchestrion installation with version pin
- **WHEN** builder stage executes
- **THEN** orchestrion@v0.18.0 is installed (not @latest)

#### Scenario: Build cache mounts
- **WHEN** builder stage executes
- **THEN** --mount=type=cache is used for /go/pkg/mod

#### Scenario: Dependency caching
- **WHEN** builder stage executes
- **THEN** go.mod and go.sum are copied first for layer caching

#### Scenario: Module download
- **WHEN** builder stage executes
- **THEN** `go mod download` runs before copying source

#### Scenario: Instrumented build with ldflags
- **WHEN** builder stage compiles
- **THEN** `orchestrion go build -ldflags="-s -w -X main.Version=${VERSION}" -o /app/server ./cmd/server` is executed

### Requirement: Runtime Stage Configuration
The runtime stage SHALL contain only the binary and CA certificates with proper ownership.

#### Scenario: Binary copy with ownership
- **WHEN** runtime stage is built
- **THEN** /server binary is copied from builder with nonroot ownership

#### Scenario: Migrations embedded
- **WHEN** runtime stage is built
- **THEN** migrations are embedded in binary (not copied as files)

#### Scenario: Minimal image size with validation
- **WHEN** the final image is built
- **THEN** it is less than 50MB (validated by make target)

#### Scenario: Port exposure
- **WHEN** Dockerfile is examined
- **THEN** port 8080 is exposed

### Requirement: Container Healthcheck for Distroless
The Dockerfile SHALL configure a healthcheck compatible with distroless (no shell).

#### Scenario: Healthcheck defined
- **WHEN** Dockerfile is examined
- **THEN** HEALTHCHECK instruction exists with --start-period

#### Scenario: Healthcheck interval
- **WHEN** healthcheck runs
- **THEN** it executes every 30 seconds with 5s start period

#### Scenario: Healthcheck command for distroless
- **WHEN** healthcheck executes
- **THEN** it calls `["/server", "--health-check"]` (not shell command)

### Requirement: Dockerignore File
The project SHALL include a .dockerignore file to exclude unnecessary files from build context.

#### Scenario: Dockerignore exists
- **WHEN** the project root is examined
- **THEN** .dockerignore file exists

#### Scenario: Git excluded
- **WHEN** .dockerignore is parsed
- **THEN** .git/ directory is excluded

#### Scenario: Test files excluded
- **WHEN** .dockerignore is parsed
- **THEN** *_test.go and test data files are excluded

#### Scenario: Environment files excluded
- **WHEN** .dockerignore is parsed
- **THEN** .env* files are excluded

#### Scenario: Docker files excluded
- **WHEN** .dockerignore is parsed
- **THEN** docker-compose.yml and Dockerfile are excluded

### Requirement: Makefile Build Targets with Security
The project SHALL include a Makefile with standard build targets and security scanning.

#### Scenario: Makefile exists
- **WHEN** the project root is examined
- **THEN** a Makefile exists

#### Scenario: Build target
- **WHEN** `make build` is executed
- **THEN** a non-instrumented binary is built to bin/server

#### Scenario: Build instrumented target
- **WHEN** `make build-instrumented` is executed
- **THEN** an instrumented binary is built using orchestrion

#### Scenario: Docker build target with size check
- **WHEN** `make docker-build` is executed
- **THEN** a Docker image is built and size is validated <50MB

#### Scenario: Generate target
- **WHEN** `make generate` is executed
- **THEN** OpenAPI code generation runs

#### Scenario: Test target
- **WHEN** `make test` is executed
- **THEN** all tests run with race detection and coverage

#### Scenario: Lint target
- **WHEN** `make lint` is executed
- **THEN** golangci-lint runs against the codebase

#### Scenario: Security scan target
- **WHEN** `make scan-vulnerabilities` is executed
- **THEN** trivy scans the image for HIGH and CRITICAL vulnerabilities

#### Scenario: SBOM generation target
- **WHEN** `make sbom` is executed
- **THEN** syft generates software bill of materials

### Requirement: Embedded Migrations
The binary SHALL embed migration SQL files for runtime execution.

#### Scenario: Migrations embedded
- **WHEN** the binary is built
- **THEN** migrations/*.sql files are embedded via embed.FS

#### Scenario: Migration access at runtime
- **WHEN** the server starts
- **THEN** embedded migrations are accessible without external files

#### Scenario: No migration files in image
- **WHEN** the final image is inspected
- **THEN** no /migrations directory exists (embedded in binary)

### Requirement: Local Development Build
The build system SHALL support non-instrumented builds for faster local development.

#### Scenario: Fast local build
- **WHEN** `make build` is executed locally
- **THEN** build completes without orchestrion overhead

#### Scenario: Local binary execution
- **WHEN** bin/server is executed locally
- **THEN** it runs without requiring OTel collector

### Requirement: Build Reproducibility
The build SHALL be reproducible given the same source and dependencies.

#### Scenario: Pinned orchestrion version
- **WHEN** Dockerfile is examined
- **THEN** orchestrion version is pinned to specific version (v0.18.0)

#### Scenario: Go version specification
- **WHEN** Dockerfile is examined
- **THEN** Go version is explicitly specified (golang:1.22-bookworm)

#### Scenario: Deterministic output
- **WHEN** the same source is built twice
- **THEN** the resulting binaries are identical

#### Scenario: Base image pinned
- **WHEN** Dockerfile is examined
- **THEN** distroless image uses SHA256 digest for reproducibility
