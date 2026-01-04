-- 008_convert_varchar_to_text.sql
-- Converts VARCHAR(255) columns to TEXT with CHECK constraints per postgres-pro conventions.
-- TEXT is preferred over VARCHAR in PostgreSQL; length limits are enforced via CHECK constraints.

-- funds.name: VARCHAR(255) -> TEXT with CHECK
ALTER TABLE funds ALTER COLUMN name TYPE TEXT;
ALTER TABLE funds ADD CONSTRAINT chk_funds_name_length CHECK (LENGTH(name) <= 255 AND LENGTH(name) >= 1);

-- cap_table_entries.owner_name: VARCHAR(255) -> TEXT with CHECK
ALTER TABLE cap_table_entries ALTER COLUMN owner_name TYPE TEXT;
ALTER TABLE cap_table_entries ADD CONSTRAINT chk_cap_table_owner_name_length CHECK (LENGTH(owner_name) <= 255 AND LENGTH(owner_name) >= 1);

-- transfers.from_owner: VARCHAR(255) -> TEXT with CHECK
ALTER TABLE transfers ALTER COLUMN from_owner TYPE TEXT;
ALTER TABLE transfers ADD CONSTRAINT chk_transfers_from_owner_length CHECK (LENGTH(from_owner) <= 255 AND LENGTH(from_owner) >= 1);

-- transfers.to_owner: VARCHAR(255) -> TEXT with CHECK
ALTER TABLE transfers ALTER COLUMN to_owner TYPE TEXT;
ALTER TABLE transfers ADD CONSTRAINT chk_transfers_to_owner_length CHECK (LENGTH(to_owner) <= 255 AND LENGTH(to_owner) >= 1);
