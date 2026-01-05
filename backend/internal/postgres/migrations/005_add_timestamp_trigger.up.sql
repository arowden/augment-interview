-- 005_add_timestamp_trigger.sql
-- Creates trigger for automatic updated_at maintenance

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cap_table_entries_timestamp
    BEFORE UPDATE ON cap_table_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

COMMENT ON FUNCTION update_timestamp() IS 'Automatically sets updated_at to current timestamp';
