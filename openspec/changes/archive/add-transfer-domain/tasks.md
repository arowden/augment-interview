## 1. Domain Layer
- [x] 1.1 Create `internal/transfer/entity.go` with Transfer struct
- [x] 1.2 Add IdempotencyKey field (*uuid.UUID) to Transfer
- [x] 1.3 Create `internal/transfer/request.go` with TransferRequest value object
- [x] 1.4 Add IdempotencyKey field to TransferRequest
- [x] 1.5 Create `internal/transfer/validator.go` with combined Validate method
- [x] 1.6 Validate fromOwner/toOwner not empty
- [x] 1.7 Validate units > 0
- [x] 1.8 Validate fromOwner != toOwner (case-sensitive, matches DB)
- [x] 1.9 Validate fromEntry exists and has sufficient units
- [x] 1.10 Create `internal/transfer/errors.go` with domain errors
- [x] 1.11 Define ErrInsufficientUnits, ErrSelfTransfer, ErrInvalidUnits, ErrInvalidOwner, ErrOwnerNotFound
- [x] 1.12 Create `internal/transfer/repository.go` with Repository interface

## 2. Repository Interface
- [x] 2.1 Define Create(ctx, transfer) method
- [x] 2.2 Define CreateTx(ctx, tx, transfer) method
- [x] 2.3 Define FindByFundID(ctx, fundID, limit, offset) returning TransferList
- [x] 2.4 Define FindByIdempotencyKey(ctx, tx, key) method
- [x] 2.5 Create TransferList struct with Transfers, TotalCount, Limit, Offset

## 3. Infrastructure Layer
- [x] 3.1 Create `internal/transfer/postgres.go` implementing Repository
- [x] 3.2 Implement Create method for recording transfers
- [x] 3.3 Implement CreateTx method using provided pgx.Tx
- [x] 3.4 Implement FindByFundID with pagination and count query
- [x] 3.5 Order transfers by transferred_at ascending
- [x] 3.6 Implement FindByIdempotencyKey for deduplication lookup

## 4. Application Layer
- [x] 4.1 Create `internal/transfer/service.go` with Service struct
- [x] 4.2 Implement ServiceOption type and functional options (WithRepository, WithPool)
- [x] 4.3 Implement NewService with functional options pattern
- [x] 4.4 Implement ExecuteTransfer with explicit SELECT FOR UPDATE SQL
- [x] 4.5 Check idempotency key first, return existing transfer if found
- [x] 4.6 Lock from_owner row with SELECT ... FOR UPDATE
- [x] 4.7 Validate sufficient units after locking
- [x] 4.8 Decrement from_owner units
- [x] 4.9 Upsert to_owner (INSERT ON CONFLICT)
- [x] 4.10 Record transfer with idempotency key
- [x] 4.11 Implement ListTransfers with pagination

## 5. HTTP Integration
- [x] 5.1 Add transfer service to HTTP handler dependencies
- [x] 5.2 Implement CreateTransfer handler with idempotencyKey from request
- [x] 5.3 Implement ListTransfers handler with limit/offset query params

## 6. Testing
- [x] 6.1 Create `internal/transfer/validator_test.go` for unit tests
- [x] 6.2 Test valid transfer validation
- [x] 6.3 Test self-transfer rejection
- [x] 6.4 Test zero/negative units rejection
- [x] 6.5 Test empty owner name rejection
- [x] 6.6 Test insufficient units rejection
- [x] 6.7 Test missing owner rejection
- [x] 6.8 Create `internal/transfer/service_test.go` with Testcontainers
- [x] 6.9 Test successful transfer flow
- [x] 6.10 Test new recipient creation
- [x] 6.11 Test idempotency key returns existing transfer
- [x] 6.12 Test concurrent transfers with same from_owner (race condition)
- [x] 6.13 Test rollback on failure preserves original state
- [x] 6.14 Test pagination of transfer history
