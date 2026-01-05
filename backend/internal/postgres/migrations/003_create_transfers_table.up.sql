-- 003_create_transfers_table.sql
-- Creates the transfers table for immutable transfer history with FK constraints

CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fund_id UUID NOT NULL REFERENCES funds(id) ON DELETE CASCADE,
    from_owner VARCHAR(255) NOT NULL,
    to_owner VARCHAR(255) NOT NULL,
    units INTEGER NOT NULL CHECK (units > 0),
    idempotency_key UUID,
    transferred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_transfer_from_owner
        FOREIGN KEY (fund_id, from_owner)
        REFERENCES cap_table_entries(fund_id, owner_name),
    CONSTRAINT fk_transfer_to_owner
        FOREIGN KEY (fund_id, to_owner)
        REFERENCES cap_table_entries(fund_id, owner_name),
    CONSTRAINT chk_different_owners
        CHECK (from_owner <> to_owner)
);

COMMENT ON TABLE transfers IS 'Immutable transfer history for audit';
COMMENT ON COLUMN transfers.id IS 'Unique identifier for the transfer';
COMMENT ON COLUMN transfers.fund_id IS 'Reference to the fund';
COMMENT ON COLUMN transfers.from_owner IS 'Owner transferring units';
COMMENT ON COLUMN transfers.to_owner IS 'Owner receiving units';
COMMENT ON COLUMN transfers.units IS 'Number of units transferred';
COMMENT ON COLUMN transfers.idempotency_key IS 'Client-generated UUID for deduplication';
COMMENT ON COLUMN transfers.transferred_at IS 'Timestamp when transfer was executed';
