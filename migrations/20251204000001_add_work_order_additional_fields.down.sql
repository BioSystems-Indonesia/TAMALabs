-- Remove visit_number, specimen_collection_date, result_release_date, and diagnosis columns from work_orders table
ALTER TABLE work_orders
DROP COLUMN IF EXISTS visit_number,
DROP COLUMN IF EXISTS specimen_collection_date,
DROP COLUMN IF EXISTS result_release_date,
DROP COLUMN IF EXISTS diagnosis;