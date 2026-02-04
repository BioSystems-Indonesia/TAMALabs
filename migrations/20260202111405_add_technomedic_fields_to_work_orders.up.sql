-- Add source, verified_at, verified_by, and completed_at fields to work_orders table
ALTER TABLE work_orders ADD COLUMN source VARCHAR(100) DEFAULT '';
ALTER TABLE work_orders ADD COLUMN verified_at DATETIME;
ALTER TABLE work_orders ADD COLUMN verified_by VARCHAR(255);
ALTER TABLE work_orders ADD COLUMN completed_at DATETIME;

-- Create index on barcode_simrs for faster lookups (already has index but ensure it exists)
CREATE INDEX IF NOT EXISTS idx_work_order_barcode_simrs ON work_orders(barcode_simrs);
