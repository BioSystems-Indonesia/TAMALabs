-- Add barcode_simrs column to work_orders table for SIMRS integration
ALTER TABLE work_orders ADD COLUMN barcode_simrs VARCHAR(255) DEFAULT '';
