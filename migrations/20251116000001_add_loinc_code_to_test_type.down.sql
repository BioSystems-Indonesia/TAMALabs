-- Drop index first
DROP INDEX IF EXISTS test_type_loinc_code ON test_types;

-- Drop loinc_code column from test_types table
ALTER TABLE test_types DROP COLUMN loinc_code;