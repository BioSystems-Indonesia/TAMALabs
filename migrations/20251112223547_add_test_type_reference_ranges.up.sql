-- Add column for specific reference ranges to test_types table
-- This will store JSON array of reference range criteria
ALTER TABLE `test_types` ADD COLUMN `specific_ref_ranges` text NULL;