-- Rollback: Restore the unique constraint on code column
-- Note: This may fail if there are duplicate codes in the database
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);

-- Drop the non-unique index created in the up migration
DROP INDEX IF EXISTS `idx_test_type_code`;