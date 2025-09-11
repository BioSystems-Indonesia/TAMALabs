-- Add back unique constraint to code column in test_types table
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);
