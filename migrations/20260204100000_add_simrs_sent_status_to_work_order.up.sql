-- Add simrs_sent_status column to work_orders table
ALTER TABLE work_orders
ADD COLUMN simrs_sent_status VARCHAR(50) DEFAULT '';

-- Add simrs_sent_at timestamp column to track when result was sent
ALTER TABLE work_orders ADD COLUMN simrs_sent_at DATETIME;

-- Add index for better query performance
CREATE INDEX idx_work_order_simrs_sent_status ON work_orders (simrs_sent_status);