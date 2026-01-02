## ADDED Requirements

### Requirement: Funds Table
The database SHALL contain a funds table with id (UUID), name (VARCHAR), total_units (INTEGER), and created_at (TIMESTAMPTZ) columns.

#### Scenario: Funds table structure
- **WHEN** the funds table is described
- **THEN** it has columns: id (UUID, PK, default gen_random_uuid()), name (VARCHAR(255), NOT NULL), total_units (INTEGER, NOT NULL, CHECK > 0), created_at (TIMESTAMPTZ, default NOW())

#### Scenario: Total units constraint
- **WHEN** an INSERT with total_units = 0 is attempted
- **THEN** a constraint violation error occurs

#### Scenario: Negative units constraint
- **WHEN** an INSERT with total_units = -100 is attempted
- **THEN** a constraint violation error occurs

### Requirement: Cap Table Entries Table
The database SHALL contain a cap_table_entries table with id, fund_id, owner_name, units, acquired_at, updated_at, and deleted_at columns.

#### Scenario: Cap table entries structure
- **WHEN** the cap_table_entries table is described
- **THEN** it has columns: id (UUID, PK), fund_id (UUID, FK to funds), owner_name (VARCHAR(255), NOT NULL), units (INTEGER, NOT NULL, CHECK >= 0), acquired_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ), deleted_at (TIMESTAMPTZ, nullable)

#### Scenario: Foreign key constraint
- **WHEN** an INSERT with non-existent fund_id is attempted
- **THEN** a foreign key violation error occurs

#### Scenario: Unique owner per fund
- **WHEN** an INSERT with duplicate (fund_id, owner_name) is attempted
- **THEN** a unique constraint violation error occurs

#### Scenario: Zero units allowed
- **WHEN** an UPDATE sets units = 0
- **THEN** the operation succeeds (owner sold all units)

#### Scenario: Negative units not allowed
- **WHEN** an UPDATE sets units = -10
- **THEN** a constraint violation error occurs

#### Scenario: Soft delete support
- **WHEN** deleted_at is set to a timestamp
- **THEN** the entry is considered soft-deleted and excluded from active queries

### Requirement: Automatic Updated Timestamp
The database SHALL automatically update the updated_at column on cap_table_entries modifications.

#### Scenario: Trigger exists
- **WHEN** the database triggers are listed
- **THEN** update_cap_table_entries_timestamp trigger exists

#### Scenario: Automatic timestamp update
- **WHEN** a cap_table_entries row is updated
- **THEN** updated_at is automatically set to NOW()

#### Scenario: Trigger function exists
- **WHEN** the database functions are listed
- **THEN** update_timestamp() function exists

### Requirement: Transfers Table with Referential Integrity
The database SHALL contain a transfers table with foreign key constraints to cap_table_entries for data integrity.

#### Scenario: Transfers table structure
- **WHEN** the transfers table is described
- **THEN** it has columns: id (UUID, PK), fund_id (UUID, FK to funds), from_owner (VARCHAR(255), NOT NULL), to_owner (VARCHAR(255), NOT NULL), units (INTEGER, NOT NULL, CHECK > 0), idempotency_key (UUID, nullable, unique), transferred_at (TIMESTAMPTZ)

#### Scenario: Transfer units constraint
- **WHEN** an INSERT with units = 0 is attempted
- **THEN** a constraint violation error occurs

#### Scenario: Idempotency key uniqueness
- **WHEN** an INSERT with duplicate idempotency_key is attempted
- **THEN** a unique constraint violation error occurs

#### Scenario: Null idempotency key allowed
- **WHEN** multiple INSERTs with null idempotency_key are performed
- **THEN** all insertions succeed (null is not considered duplicate)

#### Scenario: From owner FK constraint
- **WHEN** an INSERT references a from_owner that doesn't exist in cap_table_entries for the fund
- **THEN** a foreign key violation error occurs

#### Scenario: To owner FK constraint
- **WHEN** an INSERT references a to_owner that doesn't exist in cap_table_entries for the fund
- **THEN** a foreign key violation error occurs

### Requirement: Database Indexes
The database SHALL have indexes for common query patterns including transfer lookups by owner.

#### Scenario: Cap table fund index
- **WHEN** querying cap_table_entries by fund_id
- **THEN** an index on (fund_id) is used

#### Scenario: Cap table fund and owner index
- **WHEN** querying cap_table_entries by fund_id and owner_name
- **THEN** an index on (fund_id, owner_name) is used

#### Scenario: Transfers fund index
- **WHEN** querying transfers by fund_id
- **THEN** an index on (fund_id) is used

#### Scenario: Transfers date index
- **WHEN** querying transfers by fund_id ordered by date
- **THEN** an index on (fund_id, transferred_at) is used

#### Scenario: Transfers from owner index
- **WHEN** querying transfers by fund_id and from_owner
- **THEN** an index on (fund_id, from_owner) is used

#### Scenario: Transfers to owner index
- **WHEN** querying transfers by fund_id and to_owner
- **THEN** an index on (fund_id, to_owner) is used

#### Scenario: Idempotency key index
- **WHEN** querying transfers by idempotency_key
- **THEN** a unique index on (idempotency_key) is used

### Requirement: Cascade Delete
The database SHALL cascade delete related records when a fund is deleted.

#### Scenario: Fund deletion cascades to cap table
- **WHEN** a fund is deleted
- **THEN** all cap_table_entries for that fund are deleted

#### Scenario: Fund deletion cascades to transfers
- **WHEN** a fund is deleted
- **THEN** all transfers for that fund are deleted

### Requirement: Connection Pool with Metrics
The system SHALL use pgxpool for database connection pooling with configurable pool size and metrics.

#### Scenario: Pool creation
- **WHEN** NewPool is called with database configuration
- **THEN** a pgxpool.Pool is returned

#### Scenario: Pool with OTel tracing
- **WHEN** the pool is created
- **THEN** it has OpenTelemetry tracing configured via otelpgx

#### Scenario: Pool configuration
- **WHEN** pool is created with max connections = 10
- **THEN** the pool limits concurrent connections to 10

#### Scenario: Pool size metric
- **WHEN** metrics are collected
- **THEN** db_pool_size gauge shows configured max connections

#### Scenario: Active connections metric
- **WHEN** metrics are collected during queries
- **THEN** db_pool_active_connections gauge shows current active count

#### Scenario: Idle connections metric
- **WHEN** metrics are collected
- **THEN** db_pool_idle_connections gauge shows available connections

#### Scenario: Wait count metric
- **WHEN** a client waits for a connection
- **THEN** db_pool_wait_count counter is incremented

### Requirement: Database Migration with Versioning
The system SHALL support running SQL migrations on startup with version tracking.

#### Scenario: Run migrations
- **WHEN** RunMigrations is called with a pool
- **THEN** all migration files in /migrations are executed in order

#### Scenario: Migration idempotency
- **WHEN** RunMigrations is called twice
- **THEN** migrations are not re-applied

#### Scenario: Migration ordering
- **WHEN** migrations exist with prefixes 001_, 002_, 003_
- **THEN** they are executed in numerical order

#### Scenario: Migration version tracking
- **WHEN** migrations are run
- **THEN** a schema_migrations table tracks which versions have been applied

#### Scenario: Migration failure handling
- **WHEN** a migration fails midway
- **THEN** the transaction is rolled back and dirty flag is set

### Requirement: Embedded Migrations
The system SHALL embed migration files in the binary using embed.FS.

#### Scenario: Embedded migration files
- **WHEN** the binary is built
- **THEN** migration SQL files are embedded and accessible at runtime

### Requirement: Test Container Support
The system SHALL provide a test helper for creating PostgreSQL containers.

#### Scenario: NewTestContainer creation
- **WHEN** NewTestContainer is called in a test
- **THEN** a PostgreSQL testcontainer is started with migrations applied

#### Scenario: Container cleanup
- **WHEN** the test completes
- **THEN** the container is stopped and removed

#### Scenario: Isolated test databases
- **WHEN** multiple tests run in parallel
- **THEN** each test has its own isolated database container

### Requirement: Database Configuration
The system SHALL support database configuration via environment variables.

#### Scenario: DATABASE_URL configuration
- **WHEN** DATABASE_URL environment variable is set
- **THEN** the connection string is parsed and used

#### Scenario: Individual parameter configuration
- **WHEN** DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD are set
- **THEN** these values are used to construct the connection string

#### Scenario: Default values
- **WHEN** DB_PORT is not set
- **THEN** it defaults to 5432

#### Scenario: Pool size configuration
- **WHEN** DB_MAX_CONNECTIONS is set
- **THEN** the pool uses that maximum size

#### Scenario: Idle timeout configuration
- **WHEN** DB_IDLE_TIMEOUT is set
- **THEN** idle connections are closed after that duration
