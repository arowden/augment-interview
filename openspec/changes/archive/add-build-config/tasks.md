## 1. Orchestrion Setup
- [x] 1.1 Add orchestrion@v0.18.0 to go.mod as tool dependency
- [x] 1.2 Create `.orchestrion.yml` configuration file
- [x] 1.3 Configure automatic HTTP instrumentation
- [x] 1.4 Configure automatic database instrumentation
- [x] 1.5 Configure automatic logging injection
- [x] 1.6 Document Orchestrion limitations in design.md
- [x] 1.7 Add examples of when manual instrumentation is needed

## 2. Dockerfile
- [x] 2.1 Create multi-stage Dockerfile
- [x] 2.2 Stage 1: Builder with golang:1.22-bookworm (pinned version)
- [x] 2.3 Install orchestrion@v0.18.0 (pinned, not @latest)
- [x] 2.4 Add --mount=type=cache for /go/pkg/mod
- [x] 2.5 Build with -ldflags="-s -w" to strip debug info
- [x] 2.6 Stage 2: distroless/static-debian12:nonroot with SHA256 digest
- [x] 2.7 Add USER nonroot:nonroot directive
- [x] 2.8 Configure healthcheck with --start-period for distroless
- [x] 2.9 Use `["/server", "--health-check"]` syntax (no shell)

## 3. Dockerignore
- [x] 3.1 Create .dockerignore file
- [x] 3.2 Exclude .git/ directory
- [x] 3.3 Exclude *_test.go files
- [x] 3.4 Exclude testdata/ directories
- [x] 3.5 Exclude .env* files
- [x] 3.6 Exclude docker-compose.yml
- [x] 3.7 Exclude bin/ directory

## 4. Makefile Targets
- [x] 4.1 Add `make build` for local non-instrumented build with -ldflags
- [x] 4.2 Add `make build-instrumented` using orchestrion
- [x] 4.3 Add `make docker-build` with VERSION arg and size verification
- [x] 4.4 Add `make generate` for all code generation
- [x] 4.5 Add `make test` with race detection and coverage
- [x] 4.6 Add `make lint` for golangci-lint
- [x] 4.7 Add `make scan-vulnerabilities` using trivy
- [x] 4.8 Add `make sbom` for software bill of materials
- [x] 4.9 Add `make verify-binary` to check static linking

## 5. Server Health Check
- [x] 5.1 Implement --health-check flag in cmd/server/main.go
- [x] 5.2 Return exit code 0 for healthy, 1 for unhealthy
- [x] 5.3 Check database connectivity in health check
- [x] 5.4 Document health check behavior

## 6. CI Integration (GitHub Actions)
- [x] 6.1 Create CI workflow for PRs (test, lint, build)
- [x] 6.2 Create deploy workflow for main branch
- [x] 6.3 Configure ECR push and ECS deployment
- [x] 6.4 Configure frontend S3 deployment
- [x] 6.5 Add image tagging strategy (VERSION from git)
