-- Remove alternative_codes column from test_types table
ALTER TABLE test_types DROP COLUMN IF EXISTS alternative_codes;