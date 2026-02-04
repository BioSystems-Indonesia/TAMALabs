-- Drop index
DROP INDEX IF EXISTS idx_work_order_simrs_sent_status;

-- Drop simrs_sent_at column
ALTER TABLE work_orders DROP COLUMN simrs_sent_at;

-- Drop simrs_sent_status column
ALTER TABLE work_orders DROP COLUMN simrs_sent_status;