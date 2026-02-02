-- Drop index on sub_category_id
DROP INDEX IF EXISTS idx_test_types_sub_category_id;

-- Remove sub_category_id column from test_types
ALTER TABLE test_types DROP COLUMN sub_category_id;

-- Drop indexes from sub_categories
DROP INDEX IF EXISTS idx_sub_categories_name;
DROP INDEX IF EXISTS idx_sub_categories_code;

-- Drop sub_categories table
DROP TABLE IF EXISTS sub_categories;
