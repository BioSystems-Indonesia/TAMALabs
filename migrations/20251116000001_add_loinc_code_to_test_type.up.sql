-- Add loinc_code column to test_types table
ALTER TABLE test_types ADD COLUMN loinc_code VARCHAR(255) DEFAULT '';

-- Create index on loinc_code for faster searches
CREATE INDEX test_type_loinc_code ON test_types (loinc_code);