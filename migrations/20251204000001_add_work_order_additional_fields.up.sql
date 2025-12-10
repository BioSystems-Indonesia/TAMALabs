-- Add visit_number, specimen_collection_date, result_release_date, and diagnosis columns to work_orders table
ALTER TABLE work_orders
ADD COLUMN visit_number VARCHAR(255) DEFAULT '';

ALTER TABLE work_orders ADD COLUMN specimen_collection_date DATETIME;

ALTER TABLE work_orders ADD COLUMN result_release_date DATETIME;

ALTER TABLE work_orders ADD COLUMN diagnosis TEXT DEFAULT '';