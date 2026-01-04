## ADDED Requirements

### Requirement: Fund Entity with Constructor Validation
The system SHALL define a Fund entity with id (UUID), name (string), totalUnits (positive integer), and createdAt (timestamp). The constructor SHALL validate inputs and return errors for invalid data.

#### Scenario: Valid fund creation
- **WHEN** NewFund is called with name "Growth Fund" and totalUnits 1000
- **THEN** a Fund entity is returned with a generated UUID and current timestamp
- **AND** the name is trimmed of whitespace

#### Scenario: Invalid total units - zero
- **WHEN** NewFund is called with totalUnits 0
- **THEN** ErrInvalidFund error is returned

#### Scenario: Invalid total units - negative
- **WHEN** NewFund is called with totalUnits -100
- **THEN** ErrInvalidFund error is returned

#### Scenario: Empty name
- **WHEN** NewFund is called with an empty name ""
- **THEN** ErrInvalidFund error is returned

#### Scenario: Whitespace-only name
- **WHEN** NewFund is called with name "   "
- **THEN** ErrInvalidFund error is returned

### Requirement: Fund Repository Interface with Transaction Support
The fund package SHALL define a Repository interface with Create, CreateTx, FindByID, and FindAll methods. Fund repository does NOT handle initial ownership - that is a separate aggregate.

#### Scenario: Repository interface defined
- **WHEN** the fund package is imported
- **THEN** a Repository interface is available with Create, CreateTx, FindByID, and FindAll methods

#### Scenario: CreateTx for transaction coordination
- **WHEN** CreateTx is called with a pgx.Tx
- **THEN** the fund is created within the provided transaction
- **AND** no commit or rollback is performed by the repository

### Requirement: Fund Creation (Standalone Aggregate)
The fund repository SHALL create only the fund record. Initial ownership creation is handled separately by the ownership aggregate, coordinated at the handler level.

#### Scenario: Successful fund creation
- **WHEN** Repository.Create is called with a valid Fund
- **THEN** only the fund is persisted to the funds table
- **AND** no cap_table_entries are created by this operation

#### Scenario: Fund creation with transaction
- **WHEN** Repository.CreateTx is called with a valid Fund and pgx.Tx
- **THEN** the fund is persisted within the provided transaction
- **AND** the caller controls commit/rollback

### Requirement: Fund Retrieval by ID
The system SHALL retrieve a fund by its UUID.

#### Scenario: Fund exists
- **WHEN** FindByID is called with an existing fund's UUID
- **THEN** the Fund entity is returned

#### Scenario: Fund does not exist
- **WHEN** FindByID is called with a non-existent UUID
- **THEN** ErrFundNotFound is returned

### Requirement: List All Funds
The system SHALL retrieve all funds ordered by creation time descending.

#### Scenario: Multiple funds exist
- **WHEN** FindAll is called with 3 funds in the database
- **THEN** all 3 funds are returned ordered by createdAt descending

#### Scenario: No funds exist
- **WHEN** FindAll is called with no funds in the database
- **THEN** an empty slice is returned (not nil)

### Requirement: Fund Service with Functional Options DI
The fund package SHALL provide a Service that orchestrates fund operations using functional options for dependency injection.

#### Scenario: Service creation with options
- **WHEN** NewService is called with WithRepository(repo)
- **THEN** a Service instance is returned with the repository configured

#### Scenario: Service creation with multiple options
- **WHEN** NewService is called with multiple ServiceOptions
- **THEN** all options are applied in order

#### Scenario: CreateFund via service
- **WHEN** Service.CreateFund is called with valid name and totalUnits
- **THEN** NewFund constructor is called for validation
- **AND** the fund is created via the repository and returned

#### Scenario: CreateFund via service with invalid input
- **WHEN** Service.CreateFund is called with invalid parameters
- **THEN** ErrInvalidFund is returned without calling repository

#### Scenario: GetFund via service
- **WHEN** Service.GetFund is called with a valid UUID
- **THEN** the fund is retrieved via the repository and returned

#### Scenario: ListFunds via service
- **WHEN** Service.ListFunds is called
- **THEN** all funds are retrieved via the repository and returned

### Requirement: PostgreSQL Repository Implementation
The fund package SHALL provide a PostgresRepository that implements the Repository interface using pgxpool.

#### Scenario: PostgresRepository creation
- **WHEN** NewPostgresRepository is called with a pgxpool.Pool
- **THEN** a PostgresRepository instance is returned

#### Scenario: Create executes insert
- **WHEN** PostgresRepository.Create is called
- **THEN** it executes INSERT against only the funds table

#### Scenario: CreateTx uses provided transaction
- **WHEN** PostgresRepository.CreateTx is called with a pgx.Tx
- **THEN** it executes INSERT using the provided transaction
- **AND** it does not begin, commit, or rollback the transaction

### Requirement: Handler-Level Cross-Aggregate Orchestration
Fund and ownership creation SHALL be coordinated at the HTTP handler level, not within individual repositories.

#### Scenario: Handler creates fund with initial owner
- **WHEN** HTTP handler receives CreateFund request with initialOwner
- **THEN** handler begins a transaction
- **AND** calls fundRepo.CreateTx to create fund
- **AND** calls ownershipRepo.CreateTx to create initial cap table entry
- **AND** commits the transaction

#### Scenario: Handler rollback on ownership failure
- **WHEN** fund creation succeeds but ownership creation fails
- **THEN** the handler rolls back the entire transaction
- **AND** no fund or ownership records exist
