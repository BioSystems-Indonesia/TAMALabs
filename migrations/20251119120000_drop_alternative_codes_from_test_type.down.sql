-- Add back alternative_codes column to test_types table
ALTER TABLE test_types ADD COLUMN alternative_codes TEXT DEFAULT '';