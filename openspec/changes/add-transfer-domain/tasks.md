## 1. Domain Layer
- [ ] 1.1 Create `internal/transfer/entity.go` with Transfer struct
- [ ] 1.2 Add IdempotencyKey field (*uuid.UUID) to Transfer
- [ ] 1.3 Create `internal/transfer/request.go` with TransferRequest value object
- [ ] 1.4 Add IdempotencyKey field to TransferRequest
- [ ] 1.5 Create `internal/transfer/validator.go` with combined Validate method
- [ ] 1.6 Validate fromOwner/toOwner not empty
- [ ] 1.7 Validate units > 0
- [ ] 1.8 Validate fromOwner != toOwner (case-insensitive)
- [ ] 1.9 Validate fromEntry exists and has sufficient units
- [ ] 1.10 Create `internal/transfer/errors.go` with domain errors
- [ ] 1.11 Define ErrInsufficientUnits, ErrSelfTransfer, ErrInvalidUnits, ErrInvalidOwner, ErrOwnerNotFound
- [ ] 1.12 Create `internal/transfer/repository.go` with Repository interface

## 2. Repository Interface
- [ ] 2.1 Define Create(ctx, transfer) method
- [ ] 2.2 Define CreateTx(ctx, tx, transfer) method
- [ ] 2.3 Define FindByFundID(ctx, fundID, limit, offset) returning TransferList
- [ ] 2.4 Define FindByIdempotencyKey(ctx, tx, key) method
- [ ] 2.5 Create TransferList struct with Transfers, TotalCount, Limit, Offset

## 3. Infrastructure Layer
- [ ] 3.1 Create `internal/transfer/postgres.go` implementing Repository
- [ ] 3.2 Implement Create method for recording transfers
- [ ] 3.3 Implement CreateTx method using provided pgx.Tx
- [ ] 3.4 Implement FindByFundID with pagination and count query
- [ ] 3.5 Order transfers by transferred_at ascending
- [ ] 3.6 Implement FindByIdempotencyKey for deduplication lookup

## 4. Application Layer
- [ ] 4.1 Create `internal/transfer/service.go` with Service struct
- [ ] 4.2 Implement ServiceOption type and functional options (WithRepository, WithOwnershipRepository, WithPool)
- [ ] 4.3 Implement NewService with functional options pattern
- [ ] 4.4 Implement ExecuteTransfer with explicit SELECT FOR UPDATE SQL
- [ ] 4.5 Check idempotency key first, return existing transfer if found
- [ ] 4.6 Lock from_owner row with SELECT ... FOR UPDATE
- [ ] 4.7 Validate sufficient units after locking
- [ ] 4.8 Decrement from_owner units
- [ ] 4.9 Upsert to_owner (INSERT ON CONFLICT)
- [ ] 4.10 Record transfer with idempotency key
- [ ] 4.11 Implement ListTransfers with pagination

## 5. HTTP Integration
- [ ] 5.1 Add transfer service to HTTP handler dependencies
- [ ] 5.2 Implement CreateTransfer handler with idempotencyKey from request
- [ ] 5.3 Implement ListTransfers handler with limit/offset query params

## 6. Testing
- [ ] 6.1 Create `internal/transfer/validator_test.go` for unit tests
- [ ] 6.2 Test valid transfer validation
- [ ] 6.3 Test self-transfer rejection
- [ ] 6.4 Test zero/negative units rejection
- [ ] 6.5 Test empty owner name rejection
- [ ] 6.6 Test insufficient units rejection
- [ ] 6.7 Test missing owner rejection
- [ ] 6.8 Create `internal/transfer/service_test.go` with Testcontainers
- [ ] 6.9 Test successful transfer flow
- [ ] 6.10 Test new recipient creation
- [ ] 6.11 Test idempotency key returns existing transfer
- [ ] 6.12 Test concurrent transfers with same from_owner (race condition)
- [ ] 6.13 Test rollback on failure preserves original state
- [ ] 6.14 Test pagination of transfer history
