## 1. Domain Layer
- [x] 1.1 Create `internal/ownership/entry.go` with Entry struct
- [x] 1.2 Implement NewCapTableEntry constructor with validation returning (*Entry, error)
- [x] 1.3 Add ownerName trimming and whitespace validation
- [x] 1.4 Add units >= 0 validation (zero is valid)
- [x] 1.5 Add DeletedAt field for soft delete support
- [x] 1.6 Create `internal/ownership/captableview.go` with CapTableView read model
- [x] 1.7 Add TotalCount, Limit, Offset fields for pagination
- [x] 1.8 Implement TotalUnits() and FindOwner() methods
- [x] 1.9 Create `internal/ownership/errors.go` with ErrOwnerNotFound, ErrInvalidOwner, ErrInvalidUnits
- [x] 1.10 Create `internal/ownership/repository.go` with Repository interface

## 2. Repository Interface
- [x] 2.1 Define Create(ctx, entry) method
- [x] 2.2 Define CreateTx(ctx, tx, entry) method for transaction support
- [x] 2.3 Define FindByFundID(ctx, fundID, limit, offset) method with pagination
- [x] 2.4 Define FindByFundAndOwner(ctx, fundID, owner) method
- [x] 2.5 Define Upsert(ctx, entry) method
- [x] 2.6 Define UpsertTx(ctx, tx, entry) method

## 3. Infrastructure Layer
- [x] 3.1 Create `internal/ownership/postgres.go` implementing Repository
- [x] 3.2 Implement Create method for initial ownership
- [x] 3.3 Implement CreateTx method using provided pgx.Tx
- [x] 3.4 Implement FindByFundID with pagination and count query
- [x] 3.5 Return empty slice (not nil) for non-existent fund
- [x] 3.6 Order entries by units descending
- [x] 3.7 Implement FindByFundAndOwner returning ErrOwnerNotFound when missing
- [x] 3.8 Implement Upsert preserving acquiredAt, updating updatedAt
- [x] 3.9 Implement UpsertTx using provided pgx.Tx

## 4. Application Layer
- [x] 4.1 Create `internal/ownership/service.go` with Service struct
- [x] 4.2 Implement ServiceOption type and WithRepository function
- [x] 4.3 Implement NewService with functional options pattern
- [x] 4.4 Implement GetCapTable with pagination (default limit 100, max 1000)
- [x] 4.5 Implement GetOwnership for single owner lookup

## 5. HTTP Integration
- [x] 5.1 Add ownership repository to HTTP handler
- [x] 5.2 Implement GetCapTable handler with limit/offset query params
- [x] 5.3 Return paginated response with total count

## 6. Testing
- [x] 6.1 Create `internal/ownership/entry_test.go` for constructor validation
- [x] 6.2 Test NewCapTableEntry with valid inputs
- [x] 6.3 Test NewCapTableEntry with empty owner name
- [x] 6.4 Test NewCapTableEntry with negative units
- [x] 6.5 Test NewCapTableEntry with zero units (valid)
- [x] 6.6 Create `internal/ownership/postgres_test.go` with Testcontainers
- [x] 6.7 Test FindByFundID pagination (limit, offset, TotalCount)
- [x] 6.8 Test FindByFundID returns empty slice for missing fund
- [x] 6.9 Test FindByFundAndOwner returns ErrOwnerNotFound
- [x] 6.10 Test Upsert preserves acquiredAt on update
- [x] 6.11 Test UpsertTx with transaction rollback
- [x] 6.12 Add concurrency test for simultaneous upserts
