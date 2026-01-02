# Project Context

## Purpose
Augment Fund Cap Table System - A full-stack application to manage investment fund cap tables, tracking unit ownership and transfers between parties.

## Tech Stack
- Go 1.22+ (Backend REST API)
- PostgreSQL 16 (Database)
- React 18 + TypeScript (Frontend SPA)
- Terraform (AWS Infrastructure)
- OpenTelemetry (Observability)

## Project Conventions

### Code Style
- Go: Standard library imports first, internal packages second, third-party last (separated by blank lines)
- Go: No comments in production code unless absolutely necessary
- Go: Small, focused packages following DDD bounded contexts
- TypeScript: Strict mode enabled, no any types
- All code generated from OpenAPI spec lives in dedicated /generated directories

### Architecture Patterns
- Contract-first API design (OpenAPI 3.0 is source of truth)
- Domain-Driven Design with bounded contexts: fund, ownership, transfer
- Repository pattern for data access
- Clean architecture: domain entities independent of infrastructure
- Aggregate roots enforce invariants (Fund owns its cap table integrity)

### Testing Strategy
- Go: Testcontainers for PostgreSQL integration tests
- Go: Table-driven tests for unit tests
- Frontend: React Testing Library + MSW for API mocking
- All business logic must have test coverage

### Git Workflow
- Main branch only (no feature branches for this project)
- Build must pass before commit
- Commits pushed directly to main

## Domain Context

### Bounded Contexts
1. **Fund** - The investment fund entity, created with a name and total units
2. **Ownership** - Cap table entries tracking who owns how many units
3. **Transfer** - Transactions moving units between owners

### Key Invariants
- Total units across all cap table entries must equal fund's total units
- Transfers must not result in negative ownership
- Owner cannot transfer to themselves
- Transfer units must be positive

### Glossary
- **Fund**: An investment vehicle with a fixed number of ownership units
- **Cap Table**: The authoritative record of who owns what percentage of a fund
- **Unit**: A single ownership share in a fund
- **Transfer**: Movement of units from one owner to another
- **Initial Owner**: The party who receives all units when a fund is created

## Important Constraints
- No authentication required (demo/interview project)
- HTTP only (no HTTPS/TLS)
- Static IP for API (no domain name)
- Single region deployment (us-east-1)

## External Dependencies
- AWS RDS PostgreSQL
- AWS ECS Fargate
- AWS S3 (static website hosting)
- AWS ECR (container registry)
