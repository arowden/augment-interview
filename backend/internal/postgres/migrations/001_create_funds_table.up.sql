-- 001_create_funds_table.sql
-- Creates the funds table for storing investment fund metadata

CREATE TABLE funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    total_units INTEGER NOT NULL CHECK (total_units > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE funds IS 'Investment funds with fixed total units';
COMMENT ON COLUMN funds.id IS 'Unique identifier for the fund';
COMMENT ON COLUMN funds.name IS 'Display name of the fund';
COMMENT ON COLUMN funds.total_units IS 'Total units issued, immutable after creation';
COMMENT ON COLUMN funds.created_at IS 'Timestamp when the fund was created';
