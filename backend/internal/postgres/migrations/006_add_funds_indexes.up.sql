-- 006_add_funds_indexes.sql
-- Adds unique constraint on fund name and index for created_at ordering

-- Unique constraint to prevent duplicate fund names
CREATE UNIQUE INDEX idx_funds_name ON funds(name);

-- Composite index for efficient ordering in List queries (created_at DESC, id DESC)
CREATE INDEX idx_funds_created_at_id ON funds(created_at DESC, id DESC);
