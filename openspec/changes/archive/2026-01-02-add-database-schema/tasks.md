## 1. Migration Files
- [x] 1.1 Create `/migrations` directory (moved to `internal/postgres/migrations/` for embed.FS)
- [x] 1.2 Create `001_create_funds_table.sql` with comments
- [x] 1.3 Create `002_create_cap_table_entries.sql` with soft delete column
- [x] 1.4 Create `003_create_transfers_table.sql` with idempotency_key and FK constraints
- [x] 1.5 Create `004_add_indexes.sql` including from_owner, to_owner, idempotency indexes
- [x] 1.6 Create `005_add_timestamp_trigger.sql` for automatic updated_at
- [x] ~~1.7 Create `006_create_schema_migrations.sql`~~ (Removed: table is created programmatically by migrator)

## 2. Database Infrastructure
- [x] 2.1 Create `internal/postgres/pool.go` with connection pool setup
- [x] 2.2 Implement pool configuration from environment (DB_MAX_CONNECTIONS, DB_IDLE_TIMEOUT)
- [x] 2.3 Add OpenTelemetry tracing to pool via otelpgx
- [x] 2.4 Create `internal/postgres/migrate.go` for running migrations
- [x] 2.5 Implement embedded migrations using embed.FS
- [x] 2.6 Implement migration version tracking with schema_migrations table
- [x] 2.7 Implement dirty flag handling for failed migrations
- [x] 2.8 Add advisory lock to prevent concurrent migration execution
- [x] 2.9 Use dedicated connection for entire migration to ensure lock consistency

## 3. Pool Metrics
- [x] 3.1 Create `internal/postgres/metrics.go` for pool metrics
- [x] 3.2 Implement db_pool_size gauge
- [x] 3.3 Implement db_pool_active_connections gauge
- [x] 3.4 Implement db_pool_idle_connections gauge
- [x] 3.5 Add periodic metrics collection goroutine with safe start/stop
- [x] 3.6 Support MetricsCollector restart (recreate done channel on Start)
- [x] ~~3.7 Implement db_pool_wait_count counter~~ (Removed: pgxpool doesn't expose via Stat())
- [x] ~~3.8 Implement db_pool_wait_duration_seconds histogram~~ (Removed: pgxpool doesn't expose via Stat())

## 4. Configuration
- [x] 4.1 Add database configuration to `internal/config/config.go`
- [x] 4.2 Support DATABASE_URL environment variable
- [x] 4.3 Support individual connection parameters (DB_HOST, DB_PORT, etc.)
- [x] 4.4 Support pool size configuration (DB_MAX_CONNECTIONS, DB_MIN_CONNECTIONS)
- [x] 4.5 Support idle timeout configuration (DB_IDLE_TIMEOUT, DB_MAX_CONN_LIFETIME)
- [x] 4.6 Add SafeDSN() method for logging without exposing password
- [x] 4.7 Fix validation to check env vars directly (not defaults)

## 5. Testing Infrastructure
- [x] 5.1 Create `internal/postgres/testcontainer.go` helper
- [x] 5.2 Implement NewTestContainer function
- [x] 5.3 Implement automatic migration for tests
- [x] 5.4 Implement container cleanup
- [x] 5.5 Add test for FK constraint enforcement
- [x] 5.6 Add test for automatic updated_at trigger
- [x] 5.7 Add test for idempotency key uniqueness
- [x] 5.8 Add tests for SafeDSN and password redaction
- [x] 5.9 Add tests for config validation
- [x] 5.10 Add tests for MetricsCollector restart functionality
