-- 007_add_not_null_and_ordering_index.down.sql
-- Reverts NOT NULL constraint and ordering index

DROP INDEX IF EXISTS idx_cap_table_fund_units;

ALTER TABLE funds ALTER COLUMN created_at DROP NOT NULL;
