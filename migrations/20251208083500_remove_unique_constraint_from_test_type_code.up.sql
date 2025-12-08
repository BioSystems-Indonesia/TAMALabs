-- Drop the unique index on code column to allow multiple test types with the same code
DROP INDEX IF EXISTS `test_types_code`;

-- Drop the unique composite index on code and name if it exists
DROP INDEX IF EXISTS `idx_code_name`;

-- Create a regular index on code for better query performance
CREATE INDEX IF NOT EXISTS `idx_test_type_code` ON `test_types` (`code`);

-- Note: Removing unique constraint on code column allows multiple test types to share the same code.
-- This is useful for cases like GLUCOSE which can be used for:
-- 1. Glukosa Sewaktu
-- 2. Glukosa Darah Puasa
-- 3. Glukosa 2 Jam PP
-- Users can now select the appropriate test without code conflicts.
--
-- WARNING: Foreign keys that reference test_types(code) from observation_requests and observation_results
-- will still work in SQLite, but queries must be updated to handle multiple matches by using additional
-- criteria such as test name or specimen type to differentiate between tests with the same code.