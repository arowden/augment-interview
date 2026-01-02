# Augment Fund Cap Table System - Project Index

## Overview

A full-stack application for managing investment fund cap tables, tracking unit ownership and transfers between parties. Built as a demo/interview project with a focus on clean architecture and domain-driven design.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+ (REST API) |
| Database | PostgreSQL 16 |
| Frontend | React 18 + TypeScript (planned) |
| Infrastructure | Terraform on AWS |
| Observability | OpenTelemetry |

## Project Structure

```
augment-interview/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   └── postgres/
│       ├── config.go            # Database config with DSN builder
│       ├── pool.go              # Connection pool with OTel tracing
│       ├── migrate.go           # Embedded SQL migrations
│       ├── metrics.go           # Pool metrics for OTel
│       ├── testcontainer.go     # Integration test helpers
│       └── migrations/
│           ├── 001_create_funds_table.up.sql
│           ├── 002_create_cap_table_entries.up.sql
│           ├── 003_create_transfers_table.up.sql
│           ├── 004_add_indexes.up.sql
│           └── 005_add_timestamp_trigger.up.sql
├── deploy/
│   └── terraform/
│       ├── providers.tf
│       ├── variables.tf
│       ├── secrets.tf
│       └── outputs.tf
├── openspec/                    # Spec-driven development framework
│   ├── project.md               # Project conventions and domain context
│   ├── AGENTS.md                # AI assistant instructions
│   └── changes/                 # Change proposals (9 active)
│       ├── add-api-contract/
│       ├── add-fund-domain/
│       ├── add-ownership-domain/
│       ├── add-transfer-domain/
│       ├── add-database-schema/
│       ├── add-observability/
│       ├── add-frontend-app/
│       ├── add-aws-infrastructure/
│       └── add-build-config/
└── go.mod                       # Module: github.com/arowden/augment-fund
```

## Domain Model

### Bounded Contexts

| Context | Purpose | Key Entity |
|---------|---------|------------|
| **Fund** | Investment fund management | Fund (aggregate root) |
| **Ownership** | Cap table records | CapTableEntry |
| **Transfer** | Unit movement between owners | Transfer |

### Database Schema

```
┌─────────────────┐       ┌───────────────────────┐       ┌──────────────────┐
│     funds       │       │   cap_table_entries   │       │    transfers     │
├─────────────────┤       ├───────────────────────┤       ├──────────────────┤
│ id (PK)         │◄──────│ fund_id (FK)          │◄──────│ fund_id (FK)     │
│ name            │       │ id (PK)               │       │ id (PK)          │
│ total_units     │       │ owner_name            │◄──────│ from_owner (FK)  │
│ created_at      │       │ units                 │◄──────│ to_owner (FK)    │
└─────────────────┘       │ acquired_at           │       │ units            │
                          │ updated_at            │       │ idempotency_key  │
                          │ deleted_at (soft del) │       │ transferred_at   │
                          └───────────────────────┘       └──────────────────┘
```

### Key Invariants

- Total units across all cap table entries must equal fund's total units
- Transfers must not result in negative ownership
- Owner cannot transfer to themselves
- Transfer units must be positive

## Core Packages

### `cmd/server/main.go`

Application entry point with graceful shutdown handling.

```go
// Key flow:
1. Load config from environment
2. Connect to PostgreSQL with connection pooling
3. Run database migrations
4. Register OTel metrics
5. Wait for shutdown signal
```

### `internal/config`

Environment-based configuration using `envconfig`.

| Variable | Default | Required |
|----------|---------|----------|
| `DB_HOST` | - | Yes |
| `DB_PORT` | 5432 | No |
| `DB_USER` | - | Yes |
| `DB_PASSWORD` | - | Yes |
| `DB_NAME` | - | Yes |
| `DB_SSLMODE` | require | No |
| `SERVER_HOST` | 0.0.0.0 | No |
| `SERVER_PORT` | 8080 | No |

### `internal/postgres`

PostgreSQL utilities with comprehensive observability:

- **pool.go** - Connection pool with `otelpgx` instrumentation
- **config.go** - DSN builder with proper URL encoding
- **migrate.go** - Embedded migrations using `golang-migrate`
- **metrics.go** - Pool metrics (size, active, idle connections)

## OpenSpec Change Proposals

Active changes pending implementation:

| Change | Status | Description |
|--------|--------|-------------|
| `add-api-contract` | Proposed | OpenAPI 3.0 spec with code generation |
| `add-fund-domain` | Proposed | Fund entity, repository, service |
| `add-ownership-domain` | Proposed | Cap table entries and queries |
| `add-transfer-domain` | Proposed | Transfer operations with validation |
| `add-database-schema` | Applied | PostgreSQL migrations (implemented) |
| `add-observability` | Proposed | OTel setup with OTLP exporters |
| `add-frontend-app` | Proposed | React + TypeScript SPA |
| `add-aws-infrastructure` | Proposed | Terraform for VPC, RDS, ECS, S3 |
| `add-build-config` | Proposed | Dockerfile with Orchestrion |

## API Endpoints (Planned)

Based on `add-api-contract` proposal:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/funds` | Create fund with initial owner |
| GET | `/api/funds` | List all funds |
| GET | `/api/funds/{fundId}` | Get fund by ID |
| GET | `/api/funds/{fundId}/cap-table` | Get cap table (paginated) |
| POST | `/api/funds/{fundId}/transfers` | Execute transfer |
| GET | `/api/funds/{fundId}/transfers` | List transfers |

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `jackc/pgx/v5` | PostgreSQL driver |
| `exaring/otelpgx` | OTel tracing for pgx |
| `golang-migrate/migrate/v4` | Database migrations |
| `kelseyhightower/envconfig` | Environment config |
| `testcontainers/testcontainers-go` | Integration testing |
| `go.opentelemetry.io/otel` | Observability |

## Development Workflow

### Git Workflow

- Main branch only (no feature branches)
- Build must pass before commit
- Commits pushed directly to main

### Code Style

```go
// Import order: stdlib, internal, third-party
import (
    "context"
    "fmt"

    "github.com/arowden/augment-fund/internal/config"

    "github.com/jackc/pgx/v5/pgxpool"
)
```

### Testing Strategy

- **Unit tests**: Table-driven tests
- **Integration tests**: Testcontainers for PostgreSQL
- All business logic requires test coverage

## Quick Reference

### Run Migrations

Migrations run automatically on server start via embedded SQL files.

### Environment Variables

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=secret
export DB_NAME=augment_fund
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
```

### Build

```bash
go build -o bin/server ./cmd/server
```

## Related Documentation

- [Project Context](openspec/project.md) - Domain glossary and conventions
- [OpenSpec Guide](openspec/AGENTS.md) - Spec-driven development workflow
- [API Contract Spec](openspec/changes/add-api-contract/specs/api-contract/spec.md) - Full API requirements
