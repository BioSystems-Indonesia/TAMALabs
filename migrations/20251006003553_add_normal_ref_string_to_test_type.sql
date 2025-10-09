-- Add normal_ref_string column to test_type table
ALTER TABLE test_type
ADD COLUMN normal_ref_string TEXT DEFAULT '' COMMENT 'String reference values for qualitative tests like negative, positive, 1+, etc.';