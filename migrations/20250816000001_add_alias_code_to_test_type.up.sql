-- Add alias_code column to test_types table for SIMRS integration
ALTER TABLE test_types ADD COLUMN alias_code VARCHAR(255) DEFAULT '';
