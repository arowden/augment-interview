## ADDED Requirements

### Requirement: Transfer Entity with Idempotency Key
The system SHALL define an immutable Transfer entity representing a completed transfer with id, fundId, fromOwner, toOwner, units, idempotencyKey, and transferredAt fields.

#### Scenario: Transfer structure
- **WHEN** a Transfer is created
- **THEN** it contains UUID id, UUID fundId, string fromOwner, string toOwner, int units, optional UUID idempotencyKey, and timestamp transferredAt

#### Scenario: Transfer immutability
- **WHEN** a Transfer is persisted
- **THEN** it cannot be modified or deleted (audit record)

### Requirement: Transfer Request Value Object with Idempotency Key
The system SHALL define a TransferRequest value object with fundId, fromOwner, toOwner, units, and optional idempotencyKey fields.

#### Scenario: TransferRequest creation
- **WHEN** a TransferRequest is created with fundId, fromOwner "Alice", toOwner "Bob", units 100
- **THEN** all fields are accessible

#### Scenario: TransferRequest with idempotency key
- **WHEN** a TransferRequest is created with idempotencyKey UUID
- **THEN** the key is stored for deduplication

### Requirement: Combined Transfer Validation
The system SHALL validate transfer requests with a single Validate method combining all business rules.

#### Scenario: Valid transfer
- **WHEN** Validate is called with valid request and fromEntry with sufficient units
- **THEN** no error is returned

#### Scenario: Self-transfer rejection
- **WHEN** Validate is called with fromOwner == toOwner
- **THEN** ErrSelfTransfer is returned

#### Scenario: Zero units rejection
- **WHEN** Validate is called with units 0
- **THEN** ErrInvalidUnits is returned

#### Scenario: Negative units rejection
- **WHEN** Validate is called with units -50
- **THEN** ErrInvalidUnits is returned

#### Scenario: Empty owner name rejection
- **WHEN** Validate is called with empty fromOwner or toOwner
- **THEN** ErrInvalidOwner is returned

#### Scenario: Insufficient units rejection
- **WHEN** Validate is called with fromEntry having fewer units than requested
- **THEN** ErrInsufficientUnits is returned

#### Scenario: Missing owner rejection
- **WHEN** Validate is called with nil fromEntry
- **THEN** ErrOwnerNotFound is returned

### Requirement: Transfer Repository Interface with Pagination
The transfer package SHALL define a Repository interface with Create, CreateTx, FindByFundID (paginated), and FindByIdempotencyKey methods.

#### Scenario: Repository interface defined
- **WHEN** the transfer package is imported
- **THEN** a Repository interface is available with all required methods

#### Scenario: Idempotency lookup
- **WHEN** FindByIdempotencyKey is called with existing key
- **THEN** the existing Transfer is returned

### Requirement: Transfer Recording
The system SHALL record completed transfers for audit purposes.

#### Scenario: Create transfer record
- **WHEN** Repository.Create is called with a Transfer
- **THEN** the transfer is persisted with generated UUID and current timestamp

#### Scenario: Transfer history retrieval with pagination
- **WHEN** FindByFundID is called with limit=10, offset=0 for a fund with 25 transfers
- **THEN** TransferList with 10 transfers, TotalCount=25 is returned

#### Scenario: Chronological ordering
- **WHEN** FindByFundID returns transfers
- **THEN** they are ordered by transferredAt ascending (oldest first)

### Requirement: Transactional Transfer Recording
The system SHALL support recording transfers within an existing database transaction.

#### Scenario: CreateTx within transaction
- **WHEN** CreateTx is called with an active transaction
- **THEN** the transfer is recorded within that transaction
- **AND** no commit or rollback is performed by the repository

### Requirement: Transfer Service with Functional Options DI
The transfer package SHALL provide a Service using functional options for dependency injection.

#### Scenario: Service creation with options
- **WHEN** NewService is called with WithRepository, WithOwnershipRepository, WithPool
- **THEN** a Service instance is returned with all dependencies configured

#### Scenario: ExecuteTransfer success
- **WHEN** Service.ExecuteTransfer is called with valid request and sufficient units
- **THEN** the transfer is executed and the Transfer record is returned

#### Scenario: ExecuteTransfer with new recipient
- **WHEN** Service.ExecuteTransfer is called where toOwner does not exist
- **THEN** a new ownership entry is created for toOwner with transferred units

#### Scenario: ExecuteTransfer with idempotency key (new)
- **WHEN** Service.ExecuteTransfer is called with idempotencyKey first time
- **THEN** transfer executes normally and stores the key

#### Scenario: ExecuteTransfer with idempotency key (duplicate)
- **WHEN** Service.ExecuteTransfer is called with same idempotencyKey second time
- **THEN** the original Transfer is returned without re-executing

#### Scenario: ListTransfers via service with pagination
- **WHEN** Service.ListTransfers is called with fundId, limit, offset
- **THEN** paginated TransferList is returned

### Requirement: Atomic Transfer Execution with Explicit Locking
The system SHALL execute transfers atomically using SELECT FOR UPDATE for pessimistic locking.

#### Scenario: Full transfer atomicity
- **WHEN** a transfer is executed
- **THEN** all operations (lock, decrement, upsert, record) succeed or fail together

#### Scenario: Rollback on failure
- **WHEN** any step of transfer execution fails
- **THEN** all changes are rolled back and original ownership is preserved

#### Scenario: Concurrent transfer protection
- **WHEN** two transfers attempt to use the same from_owner's units simultaneously
- **THEN** one succeeds and one fails with ErrInsufficientUnits (no double-spending)

### Requirement: Ownership Locking with SELECT FOR UPDATE
The system SHALL lock the from_owner's entry during transfer using explicit SQL.

#### Scenario: Row-level lock acquisition
- **WHEN** a transfer begins execution
- **THEN** SELECT ... FOR UPDATE is executed on the from_owner's cap_table_entries row

#### Scenario: Lock prevents concurrent modification
- **WHEN** another transaction attempts to modify the locked entry
- **THEN** it blocks until the first transaction commits or rolls back

#### Scenario: Lock scope documented
- **WHEN** the service implementation is examined
- **THEN** the exact SELECT FOR UPDATE SQL is documented in design.md

### Requirement: Cross-Aggregate Transaction Documentation
The system SHALL document the intentional DDD boundary crossing where Transfer modifies Ownership.

#### Scenario: Trade-off documented
- **WHEN** design.md is examined
- **THEN** it explicitly documents why cross-aggregate modification is accepted

#### Scenario: Alternatives documented
- **WHEN** design.md is examined
- **THEN** it explains why event-driven/saga pattern was rejected

### Requirement: PostgreSQL Repository Implementation
The transfer package SHALL provide a PostgresRepository that implements the Repository interface.

#### Scenario: PostgresRepository creation
- **WHEN** NewPostgresRepository is called with a pgxpool.Pool
- **THEN** a PostgresRepository instance is returned

#### Scenario: Count query for pagination
- **WHEN** FindByFundID executes
- **THEN** it performs a count query to populate TotalCount
