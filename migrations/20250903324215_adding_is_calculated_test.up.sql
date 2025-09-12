-- Add column is calculated test in test type table
ALTER TABLE `test_types`
ADD COLUMN `is_calculated_test` BOOLEAN DEFAULT false;