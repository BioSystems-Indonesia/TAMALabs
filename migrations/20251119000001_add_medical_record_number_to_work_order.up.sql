-- Add medical_record_number column to work_orders table
ALTER TABLE work_orders
ADD COLUMN medical_record_number VARCHAR(255) DEFAULT '';