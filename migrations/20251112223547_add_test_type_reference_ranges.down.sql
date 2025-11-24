-- Remove specific_ref_ranges column from test_types table
-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
-- For safety, we'll just set it to NULL in down migration
UPDATE `test_types` SET `specific_ref_ranges` = NULL;