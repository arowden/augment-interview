## 1. Domain Layer
- [ ] 1.1 Create `internal/ownership/entry.go` with Entry struct
- [ ] 1.2 Implement NewCapTableEntry constructor with validation returning (*Entry, error)
- [ ] 1.3 Add ownerName trimming and whitespace validation
- [ ] 1.4 Add units >= 0 validation (zero is valid)
- [ ] 1.5 Add DeletedAt field for soft delete support
- [ ] 1.6 Create `internal/ownership/captableview.go` with CapTableView read model
- [ ] 1.7 Add TotalCount, Limit, Offset fields for pagination
- [ ] 1.8 Implement TotalUnits() and FindOwner() methods
- [ ] 1.9 Create `internal/ownership/errors.go` with ErrOwnerNotFound, ErrInvalidOwner, ErrInvalidUnits
- [ ] 1.10 Create `internal/ownership/repository.go` with Repository interface

## 2. Repository Interface
- [ ] 2.1 Define Create(ctx, entry) method
- [ ] 2.2 Define CreateTx(ctx, tx, entry) method for transaction support
- [ ] 2.3 Define FindByFundID(ctx, fundID, limit, offset) method with pagination
- [ ] 2.4 Define FindByFundAndOwner(ctx, fundID, owner) method
- [ ] 2.5 Define Upsert(ctx, entry) method
- [ ] 2.6 Define UpsertTx(ctx, tx, entry) method

## 3. Infrastructure Layer
- [ ] 3.1 Create `internal/ownership/postgres.go` implementing Repository
- [ ] 3.2 Implement Create method for initial ownership
- [ ] 3.3 Implement CreateTx method using provided pgx.Tx
- [ ] 3.4 Implement FindByFundID with pagination and count query
- [ ] 3.5 Return empty slice (not nil) for non-existent fund
- [ ] 3.6 Order entries by units descending
- [ ] 3.7 Implement FindByFundAndOwner returning ErrOwnerNotFound when missing
- [ ] 3.8 Implement Upsert preserving acquiredAt, updating updatedAt
- [ ] 3.9 Implement UpsertTx using provided pgx.Tx

## 4. Application Layer
- [ ] 4.1 Create `internal/ownership/service.go` with Service struct
- [ ] 4.2 Implement ServiceOption type and WithRepository function
- [ ] 4.3 Implement NewService with functional options pattern
- [ ] 4.4 Implement GetCapTable with pagination (default limit 100, max 1000)
- [ ] 4.5 Implement GetOwnership for single owner lookup

## 5. HTTP Integration
- [ ] 5.1 Add ownership repository to HTTP handler
- [ ] 5.2 Implement GetCapTable handler with limit/offset query params
- [ ] 5.3 Return paginated response with total count

## 6. Testing
- [ ] 6.1 Create `internal/ownership/entry_test.go` for constructor validation
- [ ] 6.2 Test NewCapTableEntry with valid inputs
- [ ] 6.3 Test NewCapTableEntry with empty owner name
- [ ] 6.4 Test NewCapTableEntry with negative units
- [ ] 6.5 Test NewCapTableEntry with zero units (valid)
- [ ] 6.6 Create `internal/ownership/postgres_test.go` with Testcontainers
- [ ] 6.7 Test FindByFundID pagination (limit, offset, TotalCount)
- [ ] 6.8 Test FindByFundID returns empty slice for missing fund
- [ ] 6.9 Test FindByFundAndOwner returns ErrOwnerNotFound
- [ ] 6.10 Test Upsert preserves acquiredAt on update
- [ ] 6.11 Test UpsertTx with transaction rollback
- [ ] 6.12 Add concurrency test for simultaneous upserts
