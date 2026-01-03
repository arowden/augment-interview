-- 006_add_funds_indexes.sql
-- Adds unique constraint on fund name and index for created_at ordering

-- Unique constraint to prevent duplicate fund names
CREATE UNIQUE INDEX idx_funds_name ON funds(name);

-- Index for efficient ordering by created_at (used in List queries)
CREATE INDEX idx_funds_created_at ON funds(created_at DESC);
