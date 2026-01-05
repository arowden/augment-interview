-- 008_convert_varchar_to_text.down.sql
-- Reverts TEXT columns back to VARCHAR(255) and removes CHECK constraints.

-- Remove CHECK constraints first.
ALTER TABLE funds DROP CONSTRAINT IF EXISTS chk_funds_name_length;
ALTER TABLE cap_table_entries DROP CONSTRAINT IF EXISTS chk_cap_table_owner_name_length;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS chk_transfers_from_owner_length;
ALTER TABLE transfers DROP CONSTRAINT IF EXISTS chk_transfers_to_owner_length;

-- Revert to VARCHAR(255).
ALTER TABLE funds ALTER COLUMN name TYPE VARCHAR(255);
ALTER TABLE cap_table_entries ALTER COLUMN owner_name TYPE VARCHAR(255);
ALTER TABLE transfers ALTER COLUMN from_owner TYPE VARCHAR(255);
ALTER TABLE transfers ALTER COLUMN to_owner TYPE VARCHAR(255);
