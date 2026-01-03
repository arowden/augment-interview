# ownership-domain Specification

## Purpose
TBD - created by archiving change add-ownership-domain. Update Purpose after archive.
## Requirements
### Requirement: Ownership Entry Entity with Constructor Validation
The system SHALL define an Entry entity representing a single ownership record with id, fundId, ownerName, units, acquiredAt, updatedAt, and deletedAt fields. The constructor SHALL validate inputs.

#### Scenario: Entry creation with valid inputs
- **WHEN** NewCapTableEntry is called with valid fundId, ownerName "Alice", and units 1000
- **THEN** an Entry is returned with generated UUID, current acquiredAt and updatedAt

#### Scenario: Entry creation with empty owner name
- **WHEN** NewCapTableEntry is called with empty ownerName ""
- **THEN** ErrInvalidOwner is returned

#### Scenario: Entry creation with whitespace owner name
- **WHEN** NewCapTableEntry is called with ownerName "   "
- **THEN** ErrInvalidOwner is returned

#### Scenario: Entry creation with negative units
- **WHEN** NewCapTableEntry is called with units -100
- **THEN** ErrInvalidUnits is returned

#### Scenario: Entry with zero units is valid
- **WHEN** NewCapTableEntry is called with units = 0
- **THEN** Entry is created successfully (owner sold all units but record remains)

#### Scenario: Entry soft delete support
- **WHEN** an Entry is created
- **THEN** it has DeletedAt field of type *time.Time (nil for active entries)

### Requirement: Cap Table View (Read Model)
The system SHALL define a CapTableView read model (not aggregate) that holds ownership entries with pagination metadata.

#### Scenario: CapTableView structure
- **WHEN** a CapTableView is returned
- **THEN** it contains FundID, Entries slice, TotalCount, Limit, and Offset fields

#### Scenario: TotalUnits calculation
- **WHEN** CapTableView.TotalUnits() is called on a view with entries of 500, 300, and 200 units
- **THEN** it returns 1000

#### Scenario: FindOwner lookup
- **WHEN** CapTableView.FindOwner("Alice") is called
- **THEN** it returns the Entry for Alice or nil if not found

### Requirement: Ownership Repository Interface with Transaction Support
The ownership package SHALL define a Repository interface with Create, CreateTx, FindByFundID, FindByFundAndOwner, Upsert, and UpsertTx methods.

#### Scenario: Repository interface defined
- **WHEN** the ownership package is imported
- **THEN** a Repository interface is available with all required methods

#### Scenario: CreateTx for initial owner
- **WHEN** CreateTx is called with a pgx.Tx
- **THEN** the entry is created within the provided transaction
- **AND** no commit or rollback is performed by the repository

### Requirement: Cap Table Retrieval with Pagination
The system SHALL retrieve ownership entries for a fund as a CapTableView with pagination support.

#### Scenario: Fund with multiple owners - paginated
- **WHEN** FindByFundID is called with limit=10, offset=0 for a fund with 25 owners
- **THEN** a CapTableView with 10 entries, TotalCount=25, Limit=10, Offset=0 is returned

#### Scenario: Fund with single owner
- **WHEN** FindByFundID is called for a newly created fund
- **THEN** a CapTableView with 1 entry (initial owner) is returned

#### Scenario: Non-existent fund returns empty view
- **WHEN** FindByFundID is called for a non-existent fund
- **THEN** a CapTableView with empty entries slice (not nil) and TotalCount=0 is returned

#### Scenario: Pagination bounds respected
- **WHEN** FindByFundID is called with offset > total entries
- **THEN** a CapTableView with empty entries slice and correct TotalCount is returned

### Requirement: Single Owner Lookup with Error Handling
The system SHALL retrieve a single ownership entry by fund and owner name.

#### Scenario: Owner exists
- **WHEN** FindByFundAndOwner is called for existing owner "Alice"
- **THEN** the Entry for Alice is returned

#### Scenario: Owner does not exist
- **WHEN** FindByFundAndOwner is called for non-existent owner
- **THEN** ErrOwnerNotFound is returned (not nil entry)

### Requirement: Ownership Create for Initial Owner
The system SHALL support creating initial ownership entries via dedicated Create methods.

#### Scenario: Create initial owner
- **WHEN** Create is called for new owner "Alice" with all fund units
- **THEN** a new Entry is created with acquiredAt and updatedAt set to current time

#### Scenario: CreateTx for coordinated fund creation
- **WHEN** CreateTx is called within a transaction (coordinated with fund creation)
- **THEN** the entry is created in the provided transaction

### Requirement: Ownership Upsert
The system SHALL create or update ownership entries using an upsert operation.

#### Scenario: Create new owner via upsert
- **WHEN** Upsert is called for a new owner "Bob"
- **THEN** a new Entry is created with acquiredAt set to current time

#### Scenario: Update existing owner
- **WHEN** Upsert is called for existing owner "Alice" with different units
- **THEN** the Entry is updated with new units and updatedAt set to current time
- **AND** acquiredAt is NOT modified

#### Scenario: Upsert with zero units (soft state)
- **WHEN** Upsert is called with units = 0
- **THEN** the Entry is updated to 0 units (not deleted, for audit trail)

### Requirement: Transactional Upsert
The system SHALL support upserting within an existing database transaction for atomic multi-entry updates.

#### Scenario: UpsertTx within transaction
- **WHEN** UpsertTx is called with an active transaction
- **THEN** the upsert is performed within that transaction
- **AND** no commit or rollback is performed by the repository

#### Scenario: Transaction rollback
- **WHEN** a transaction containing UpsertTx calls is rolled back
- **THEN** no ownership changes are persisted

### Requirement: Ownership Service with Functional Options DI
The ownership package SHALL provide a Service using functional options for dependency injection.

#### Scenario: Service creation with options
- **WHEN** NewService is called with WithRepository(repo)
- **THEN** a Service instance is returned with the repository configured

#### Scenario: GetCapTable via service with pagination
- **WHEN** Service.GetCapTable is called with fundId, limit, offset
- **THEN** the CapTableView is retrieved via the repository with pagination

#### Scenario: GetCapTable default limit
- **WHEN** Service.GetCapTable is called with limit=0
- **THEN** limit defaults to 100

#### Scenario: GetCapTable max limit
- **WHEN** Service.GetCapTable is called with limit > 1000
- **THEN** limit is capped at 1000

#### Scenario: GetOwnership via service
- **WHEN** Service.GetOwnership is called with fundId and ownerName
- **THEN** the Entry is retrieved via the repository

### Requirement: PostgreSQL Repository Implementation
The ownership package SHALL provide a PostgresRepository that implements the Repository interface.

#### Scenario: PostgresRepository creation
- **WHEN** NewPostgresRepository is called with a pgxpool.Pool
- **THEN** a PostgresRepository instance is returned

#### Scenario: Entries ordered by units descending
- **WHEN** FindByFundID returns entries
- **THEN** they are ordered by units descending (largest stakeholders first)

#### Scenario: Count query for pagination
- **WHEN** FindByFundID executes
- **THEN** it performs a count query to populate TotalCount

