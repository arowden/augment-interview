# Augment Fund Cap Table - Project Index

> **Quick Reference Guide** | See [README.md](README.md) for comprehensive documentation

## Project Overview

Full-stack investment fund cap table management system with:
- **Backend**: Go 1.24 REST API with Chi router, PostgreSQL, OpenTelemetry
- **Frontend**: React 18 + TypeScript + Vite + Tailwind CSS
- **Infrastructure**: Terraform for AWS (ECS Fargate, RDS, S3, ALB)

## Quick Start

```bash
# Start everything with Docker Compose
docker-compose up -d

# API: http://localhost:8080
# Frontend: npm run dev (in frontend/)
```

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Frontend                                   │
│  React SPA → TanStack Query → Generated API Client → REST API       │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Backend (Go)                                  │
│  HTTP Handler → Services → Repositories → PostgreSQL                │
│       │                                                              │
│       └── OpenTelemetry → Traces & Metrics                          │
└─────────────────────────────────────────────────────────────────────┘
```

## Domain Model

| Entity | Location | Description |
|--------|----------|-------------|
| **Fund** | `backend/internal/fund/` | Investment fund with fixed units |
| **CapTableEntry** | `backend/internal/ownership/` | Owner's stake in a fund |
| **Transfer** | `backend/internal/transfer/` | Unit movement between owners |

### Key Invariants
- Sum of cap table entries = fund's total units
- No negative ownership (transfers validated)
- No self-transfers
- Idempotency keys prevent duplicate transfers

## Project Structure

### Backend (`backend/`)

| Path | Purpose |
|------|---------|
| `api/openapi.yaml` | OpenAPI 3.0 spec (source of truth) |
| `cmd/server/main.go` | Application entrypoint |
| `internal/config/` | Environment configuration |
| `internal/fund/` | Fund domain (entity, service, store) |
| `internal/ownership/` | Cap table domain |
| `internal/transfer/` | Transfer domain with validation |
| `internal/http/` | HTTP handlers + generated code |
| `internal/postgres/` | Database utilities + migrations |
| `internal/otel/` | OpenTelemetry setup |

### Frontend (`frontend/`)

| Path | Purpose |
|------|---------|
| `src/api/generated/` | Auto-generated API client |
| `src/components/` | React components (12 components) |
| `src/pages/` | Page components (Dashboard, FundPage, OwnersPage) |
| `src/hooks/` | TanStack Query hooks (useFunds, useTransfers, useCapTable) |
| `src/schemas/` | Zod validation schemas |
| `e2e/` | Playwright E2E tests |

### Infrastructure (`deploy/terraform/`)

| Module | Resources |
|--------|-----------|
| `vpc/` | VPC, subnets, NAT gateway |
| `rds/` | PostgreSQL RDS instance |
| `ecr/` | Container registry |
| `ecs/` | Fargate cluster + service |
| `frontend/` | S3 + CloudFront |

## API Quick Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/funds` | List funds (paginated) |
| `POST` | `/api/funds` | Create fund with initial owner |
| `GET` | `/api/funds/{id}` | Get fund details |
| `GET` | `/api/funds/{id}/cap-table` | Get ownership table |
| `POST` | `/api/funds/{id}/transfers` | Execute transfer |
| `GET` | `/api/funds/{id}/transfers` | List transfers |
| `POST` | `/api/reset` | Reset database (dev only) |
| `GET` | `/healthz` | Health check |

## Key Files Reference

### Backend Entry Points

| File | Purpose |
|------|---------|
| `backend/cmd/server/main.go` | Server startup, DI, graceful shutdown |
| `backend/internal/http/handler.go` | HTTP request handlers |
| `backend/internal/http/openapi.gen.go` | Generated types & interfaces |

### Service Layer

| File | Key Functions |
|------|---------------|
| `fund/service.go` | `CreateFundWithInitialOwner`, `GetFund`, `ListFunds` |
| `ownership/service.go` | `GetCapTable`, `GetAllOwnersWithHoldings` |
| `transfer/service.go` | `ExecuteTransfer`, `ListTransfers` |

### Repository Layer

| File | Database Operations |
|------|---------------------|
| `fund/store.go` | CRUD for funds table |
| `ownership/store.go` | Cap table entries, upserts, soft delete |
| `transfer/store.go` | Transfer records, idempotency lookup |

### Frontend Components

| Component | Purpose |
|-----------|---------|
| `Layout.tsx` | App shell with navigation |
| `Dashboard.tsx` | Fund list, stats, create fund |
| `FundPage.tsx` | Fund detail with cap table + transfers |
| `CapTable.tsx` | Ownership percentage visualization |
| `TransferForm.tsx` | Transfer creation with validation |
| `TransferHistory.tsx` | Transfer timeline |

## Development Commands

```bash
# Build & Run
make build              # Build server binary
make run                # Build and run server
make docker-build       # Build Docker image

# Testing
make test               # Unit tests (fast)
make test-integration   # Integration tests (Docker required)
make test-all           # All tests
make test-coverage      # Coverage report

# Code Quality
make lint               # golangci-lint
make generate-api       # Regenerate from OpenAPI

# Frontend
cd frontend
npm run dev             # Development server
npm run build           # Production build
npm run test            # Vitest tests
npm run e2e             # Playwright E2E
```

## Environment Variables

### Required
| Variable | Description |
|----------|-------------|
| `DB_HOST` | PostgreSQL host |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_NAME` | Database name |

### Optional
| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PORT` | `5432` | Database port |
| `SERVER_PORT` | `8080` | API server port |
| `OTEL_ENABLED` | `false` | Enable telemetry |

See [README.md#environment-variables](README.md#environment-variables) for full list.

## Code Patterns

### Functional Options (Service Construction)
```go
svc, err := transfer.NewService(
    transfer.WithRepository(transferStore),
    transfer.WithOwnershipRepository(ownershipStore),
    transfer.WithPool(pool),
)
```

### Transaction Management
```go
tx, err := s.pool.Begin(ctx)
defer tx.Rollback(ctx)
// ... operations ...
tx.Commit(ctx)
```

### Domain Errors
```go
var ErrNotFound = errors.New("fund not found")

if errors.Is(err, fund.ErrNotFound) {
    // Handle not found
}
```

## Tech Stack Summary

| Layer | Technology |
|-------|------------|
| Language | Go 1.24 |
| Router | Chi v5 |
| Database | PostgreSQL 16 |
| DB Driver | pgx v5 |
| Migrations | golang-migrate |
| API Spec | OpenAPI 3.0 |
| Code Gen | oapi-codegen |
| Testing | testify, testcontainers-go |
| Linting | golangci-lint |
| Telemetry | OpenTelemetry |
| Frontend | React 18, Vite, Tailwind CSS |
| State | TanStack Query |
| Validation | Zod + react-hook-form |
| E2E Tests | Playwright |
| Infrastructure | Terraform, AWS ECS Fargate |

## Related Documentation

| Document | Contents |
|----------|----------|
| [README.md](README.md) | Comprehensive docs with Mermaid diagrams |
| [openspec/project.md](openspec/project.md) | Domain glossary, conventions |
| [openspec/AGENTS.md](openspec/AGENTS.md) | AI assistant instructions |
| [backend/api/openapi.yaml](backend/api/openapi.yaml) | Full API specification |
| [.golangci.yml](backend/.golangci.yml) | Linter configuration |

## OpenSpec Changes (Archived)

All planned features have been implemented:

| Change | Status |
|--------|--------|
| `add-api-contract` | Implemented |
| `add-fund-domain` | Implemented |
| `add-ownership-domain` | Implemented |
| `add-transfer-domain` | Implemented |
| `add-database-schema` | Implemented |
| `add-observability` | Implemented |
| `add-frontend-app` | Implemented |
| `add-aws-infrastructure` | Implemented |
| `add-build-config` | Implemented |
