-- 002_create_cap_table_entries.sql
-- Creates the cap_table_entries table for ownership records with soft delete support

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
COMMENT ON COLUMN cap_table_entries.id IS 'Unique identifier for the entry';
COMMENT ON COLUMN cap_table_entries.fund_id IS 'Reference to the parent fund';
COMMENT ON COLUMN cap_table_entries.owner_name IS 'Name of the unit owner';
COMMENT ON COLUMN cap_table_entries.units IS 'Current units owned, 0 = sold all';
COMMENT ON COLUMN cap_table_entries.acquired_at IS 'Timestamp when ownership was first acquired';
COMMENT ON COLUMN cap_table_entries.updated_at IS 'Timestamp of last modification';
COMMENT ON COLUMN cap_table_entries.deleted_at IS 'Soft delete timestamp for audit trail';
