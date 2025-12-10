-- Add alternative_codes column to test_types table as JSON
ALTER TABLE test_types ADD COLUMN alternative_codes TEXT DEFAULT '';