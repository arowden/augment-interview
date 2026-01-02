## Context
PostgreSQL serves as the single source of truth for fund, ownership, and transfer data. The schema must enforce business constraints and support the transactional requirements of transfers.

## Goals / Non-Goals
- Goals: Data integrity, referential integrity, query performance, migration support, audit trail
- Non-Goals: Multi-database support, read replicas, sharding

## Decisions
- Decision: Use gen_random_uuid() for UUID generation in database
- Alternatives considered: Application-generated UUIDs (less atomic), serial IDs (less distributed-friendly)

- Decision: Use TIMESTAMP WITH TIME ZONE for all timestamps
- Alternatives considered: TIMESTAMP without timezone (timezone ambiguity)

- Decision: Unique constraint on (fund_id, owner_name) in cap_table_entries
- Alternatives considered: Allowing duplicate entries (violates domain model)

- Decision: pgxpool for connection pooling with OTel instrumentation
- Alternatives considered: database/sql (less pgx features), pgbouncer (external dependency)

- Decision: Foreign key constraints on transfers referencing cap_table_entries
- Alternatives considered: String-only owner references (no referential integrity)

- Decision: Soft delete via deleted_at column for audit trail
- Alternatives considered: Hard delete only (loses history), separate audit table (more complexity)

- Decision: Idempotency key column with unique constraint for transfer deduplication
- Alternatives considered: Application-level deduplication only (race conditions possible)

## Schema Design

### funds table
```sql
CREATE TABLE funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    total_units INTEGER NOT NULL CHECK (total_units > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE funds IS 'Investment funds with fixed total units';
COMMENT ON COLUMN funds.total_units IS 'Total units issued, immutable after creation';
```

### cap_table_entries table
```sql
CREATE TABLE cap_table_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fund_id UUID NOT NULL REFERENCES funds(id) ON DELETE CASCADE,
    owner_name VARCHAR(255) NOT NULL,
    units INTEGER NOT NULL CHECK (units >= 0),
    acquired_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(fund_id, owner_name)
);

COMMENT ON TABLE cap_table_entries IS 'Cap table ownership records';
COMMENT ON COLUMN cap_table_entries.units IS 'Current units owned, 0 = sold all';
COMMENT ON COLUMN cap_table_entries.deleted_at IS 'Soft delete timestamp for audit';
```

### Automatic updated_at trigger
```sql
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cap_table_entries_timestamp
    BEFORE UPDATE ON cap_table_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
```

### transfers table
```sql
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fund_id UUID NOT NULL REFERENCES funds(id) ON DELETE CASCADE,
    from_owner VARCHAR(255) NOT NULL,
    to_owner VARCHAR(255) NOT NULL,
    units INTEGER NOT NULL CHECK (units > 0),
    idempotency_key UUID UNIQUE,
    transferred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_transfer_from_owner
        FOREIGN KEY (fund_id, from_owner)
        REFERENCES cap_table_entries(fund_id, owner_name),
    CONSTRAINT fk_transfer_to_owner
        FOREIGN KEY (fund_id, to_owner)
        REFERENCES cap_table_entries(fund_id, owner_name)
);

COMMENT ON TABLE transfers IS 'Immutable transfer history for audit';
COMMENT ON COLUMN transfers.idempotency_key IS 'Client-generated UUID for deduplication';
```

### Indexes
```sql
CREATE INDEX idx_cap_table_fund ON cap_table_entries(fund_id);
CREATE INDEX idx_cap_table_fund_owner ON cap_table_entries(fund_id, owner_name);
CREATE INDEX idx_cap_table_active ON cap_table_entries(fund_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_transfers_fund ON transfers(fund_id);
CREATE INDEX idx_transfers_fund_date ON transfers(fund_id, transferred_at DESC);
CREATE INDEX idx_transfers_from_owner ON transfers(fund_id, from_owner, transferred_at DESC);
CREATE INDEX idx_transfers_to_owner ON transfers(fund_id, to_owner, transferred_at DESC);
CREATE UNIQUE INDEX idx_transfers_idempotency ON transfers(idempotency_key) WHERE idempotency_key IS NOT NULL;
```

### Schema migrations table
```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT false,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE schema_migrations IS 'Tracks applied database migrations';
COMMENT ON COLUMN schema_migrations.dirty IS 'True if migration failed midway';
```

## Package Structure
```
/migrations/
  001_create_funds_table.sql
  002_create_cap_table_entries.sql
  003_create_transfers_table.sql
  004_add_indexes.sql
  005_add_timestamp_trigger.sql
  006_add_transfer_fk_constraints.sql

internal/postgres/
  pool.go          - Connection pool setup with OTel and metrics
  migrate.go       - Migration runner with version tracking
  metrics.go       - Pool metrics collection
  testcontainer.go - Test helper
```

## Connection Pool Configuration
```go
type PoolConfig struct {
    DSN             string
    MaxConns        int32         // DB_MAX_CONNECTIONS, default 25
    MinConns        int32         // DB_MIN_CONNECTIONS, default 5
    MaxConnLifetime time.Duration // DB_MAX_CONN_LIFETIME, default 1h
    MaxConnIdleTime time.Duration // DB_IDLE_TIMEOUT, default 10m
}
```

## Pool Metrics
```go
var (
    poolSize        metric.Int64Gauge      // db_pool_size
    activeConns     metric.Int64Gauge      // db_pool_active_connections
    idleConns       metric.Int64Gauge      // db_pool_idle_connections
    waitCount       metric.Int64Counter    // db_pool_wait_count
    waitDuration    metric.Float64Histogram // db_pool_wait_duration_seconds
)

func collectPoolMetrics(pool *pgxpool.Pool) {
    stat := pool.Stat()
    poolSize.Record(ctx, int64(stat.MaxConns()))
    activeConns.Record(ctx, int64(stat.AcquiredConns()))
    idleConns.Record(ctx, int64(stat.IdleConns()))
}
```

## Risks / Trade-offs
- CASCADE delete propagates from funds → Acceptable, fund deletion removes all related data
- FK constraints on transfers require cap_table_entries to exist first → Handled by transaction ordering
- Soft delete adds query complexity → Mitigated by partial index on active entries

## Open Questions
- None
