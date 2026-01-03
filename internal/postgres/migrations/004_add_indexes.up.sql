-- 004_add_indexes.up.sql
-- Creates indexes for query performance

-- Cap table indexes
-- Note: idx_cap_table_fund_owner covers fund_id lookups via leftmost prefix
CREATE INDEX idx_cap_table_fund_owner ON cap_table_entries(fund_id, owner_name);
CREATE INDEX idx_cap_table_active ON cap_table_entries(fund_id) WHERE deleted_at IS NULL;

-- Transfer indexes
CREATE INDEX idx_transfers_fund ON transfers(fund_id);
CREATE INDEX idx_transfers_fund_date ON transfers(fund_id, transferred_at DESC);
CREATE INDEX idx_transfers_from_owner ON transfers(fund_id, from_owner, transferred_at DESC);
CREATE INDEX idx_transfers_to_owner ON transfers(fund_id, to_owner, transferred_at DESC);
CREATE UNIQUE INDEX idx_transfers_idempotency ON transfers(idempotency_key) WHERE idempotency_key IS NOT NULL;
