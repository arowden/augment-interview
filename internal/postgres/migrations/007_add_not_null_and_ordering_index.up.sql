-- 007_add_not_null_and_ordering_index.sql
-- Adds NOT NULL constraint to funds.created_at and index for cap_table ordering

-- Add NOT NULL constraint to funds.created_at
-- Safe: DEFAULT NOW() ensures no NULL values exist
ALTER TABLE funds ALTER COLUMN created_at SET NOT NULL;

-- Add index for cap_table ordering by units DESC, owner_name ASC
-- Improves performance for FindByFundID queries with ORDER BY
CREATE INDEX idx_cap_table_fund_units ON cap_table_entries(fund_id, units DESC, owner_name ASC) WHERE deleted_at IS NULL;
