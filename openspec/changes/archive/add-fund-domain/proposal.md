# Change: Add Fund Domain

## Why
Funds are the core aggregate root of the cap table system. A well-defined Fund domain enables clean separation of concerns and enforces business invariants at the domain level.

## What Changes
- Add `internal/fund` package with domain entity, repository interface, and service
- Implement PostgreSQL repository for fund persistence
- Wire fund service into HTTP handlers
- Add integration tests using Testcontainers

## Impact
- Affected specs: fund-domain (new)
- Affected code: `internal/fund/`, `internal/http/handler.go`
