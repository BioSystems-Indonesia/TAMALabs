-- Create sub_categories table
CREATE TABLE IF NOT EXISTS sub_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL,
    category VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_sub_categories_name ON sub_categories(name);
CREATE INDEX idx_sub_categories_code ON sub_categories(code);

-- Migrate existing sub_category data from test_types to sub_categories table
INSERT INTO sub_categories (name, code, category)
SELECT DISTINCT
    sub_category as name,
    UPPER(SUBSTR(sub_category, 1, 3)) as code,
    category
FROM test_types
WHERE sub_category IS NOT NULL AND sub_category != ''
ORDER BY sub_category;

-- Add sub_category_id column to test_types
ALTER TABLE test_types ADD COLUMN sub_category_id INTEGER;

-- Update test_types to reference sub_categories
UPDATE test_types
SET sub_category_id = (
    SELECT id FROM sub_categories
    WHERE sub_categories.name = test_types.sub_category
)
WHERE sub_category IS NOT NULL AND sub_category != '';

-- Create index on sub_category_id
CREATE INDEX idx_test_types_sub_category_id ON test_types(sub_category_id);
