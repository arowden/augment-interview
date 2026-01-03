-- 006_add_funds_indexes.down.sql
-- Removes fund indexes

DROP INDEX IF EXISTS idx_funds_created_at_id;
DROP INDEX IF EXISTS idx_funds_name;
