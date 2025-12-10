-- Add alternative_codes column to test_types table
-- This allows storing multiple codes from different devices for the same test type
ALTER TABLE test_types
ADD COLUMN alternative_codes TEXT DEFAULT NULL;