## 1. Domain Layer
- [x] 1.1 Create `internal/fund/entity.go` with Fund struct
- [x] 1.2 Implement NewFund constructor with validation returning (*Fund, error)
- [x] 1.3 Add name trimming and whitespace validation
- [x] 1.4 Add totalUnits positive integer validation
- [x] 1.5 Create `internal/fund/errors.go` with ErrNotFound, ErrInvalidFund
- [x] 1.6 Create `internal/fund/repository.go` with Repository interface

## 2. Repository Interface
- [x] 2.1 Define Create(ctx, fund) method
- [x] 2.2 Define CreateTx(ctx, tx, fund) method for transaction support
- [x] 2.3 Define FindByID(ctx, id) method
- [x] 2.4 Define FindAll(ctx) method

## 3. Infrastructure Layer
- [x] 3.1 Create `internal/fund/postgres.go` implementing Repository
- [x] 3.2 Implement Create method (fund table only, no cap_table_entries)
- [x] 3.3 Implement CreateTx method using provided pgx.Tx
- [x] 3.4 Implement FindByID method
- [x] 3.5 Implement FindAll method with ORDER BY created_at DESC

## 4. Application Layer
- [x] 4.1 Create `internal/fund/service.go` with Service struct
- [x] 4.2 Implement ServiceOption type and WithRepository function
- [x] 4.3 Implement NewService with functional options pattern
- [x] 4.4 Implement CreateFund calling NewFund for validation
- [x] 4.5 Implement GetFund
- [x] 4.6 Implement ListFunds

## 5. HTTP Handler Integration
- [ ] 5.1 Add fund repository and ownership repository to HTTP handler
- [ ] 5.2 Add pgxpool.Pool to handler for transaction management
- [ ] 5.3 Implement CreateFund handler with cross-aggregate transaction
- [ ] 5.4 Begin transaction, create fund, create ownership, commit
- [ ] 5.5 Implement rollback on any failure
- [ ] 5.6 Implement GetFund handler method
- [ ] 5.7 Implement ListFunds handler method

NOTE: Section 5 is blocked on add-ownership-domain. The CreateFund handler requires cross-aggregate transaction coordination with the ownership repository.

## 6. Testing
- [x] 6.1 Create `internal/fund/entity_test.go` for constructor validation
- [x] 6.2 Test NewFund with valid inputs
- [x] 6.3 Test NewFund with empty name
- [x] 6.4 Test NewFund with whitespace-only name
- [x] 6.5 Test NewFund with zero totalUnits
- [x] 6.6 Test NewFund with negative totalUnits
- [x] 6.7 Create `internal/fund/postgres_test.go` with Testcontainers
- [x] 6.8 Test Create persists only to funds table
- [x] 6.9 Test CreateTx with transaction
- [x] 6.10 Test FindByID returns ErrNotFound for missing
- [x] 6.11 Test FindAll returns empty slice (not nil) when no funds
