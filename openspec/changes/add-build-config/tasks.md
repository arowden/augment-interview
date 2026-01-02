## 1. Orchestrion Setup
- [ ] 1.1 Add orchestrion@v0.18.0 to go.mod as tool dependency
- [ ] 1.2 Create `.orchestrion.yml` configuration file
- [ ] 1.3 Configure automatic HTTP instrumentation
- [ ] 1.4 Configure automatic database instrumentation
- [ ] 1.5 Configure automatic logging injection
- [ ] 1.6 Document Orchestrion limitations in design.md
- [ ] 1.7 Add examples of when manual instrumentation is needed

## 2. Dockerfile
- [ ] 2.1 Create multi-stage Dockerfile
- [ ] 2.2 Stage 1: Builder with golang:1.22-bookworm (pinned version)
- [ ] 2.3 Install orchestrion@v0.18.0 (pinned, not @latest)
- [ ] 2.4 Add --mount=type=cache for /go/pkg/mod
- [ ] 2.5 Build with -ldflags="-s -w" to strip debug info
- [ ] 2.6 Stage 2: distroless/static-debian12:nonroot with SHA256 digest
- [ ] 2.7 Add USER nonroot:nonroot directive
- [ ] 2.8 Configure healthcheck with --start-period for distroless
- [ ] 2.9 Use `["/server", "--health-check"]` syntax (no shell)

## 3. Dockerignore
- [ ] 3.1 Create .dockerignore file
- [ ] 3.2 Exclude .git/ directory
- [ ] 3.3 Exclude *_test.go files
- [ ] 3.4 Exclude testdata/ directories
- [ ] 3.5 Exclude .env* files
- [ ] 3.6 Exclude docker-compose.yml
- [ ] 3.7 Exclude bin/ directory

## 4. Makefile Targets
- [ ] 4.1 Add `make build` for local non-instrumented build with -ldflags
- [ ] 4.2 Add `make build-instrumented` using orchestrion
- [ ] 4.3 Add `make docker-build` with VERSION arg and size verification
- [ ] 4.4 Add `make generate` for all code generation
- [ ] 4.5 Add `make test` with race detection and coverage
- [ ] 4.6 Add `make lint` for golangci-lint
- [ ] 4.7 Add `make scan-vulnerabilities` using trivy
- [ ] 4.8 Add `make sbom` for software bill of materials
- [ ] 4.9 Add `make verify-binary` to check static linking

## 5. Server Health Check
- [ ] 5.1 Implement --health-check flag in cmd/server/main.go
- [ ] 5.2 Return exit code 0 for healthy, 1 for unhealthy
- [ ] 5.3 Check database connectivity in health check
- [ ] 5.4 Document health check behavior

## 6. CI Integration
- [ ] 6.1 Document orchestrion requirements
- [ ] 6.2 Add build verification step
- [ ] 6.3 Configure image tagging strategy (VERSION from git)
- [ ] 6.4 Add image size validation (<50MB)
- [ ] 6.5 Add vulnerability scan step
