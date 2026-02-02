-- Remove source, verified_at, verified_by, and completed_at fields from work_orders table
ALTER TABLE work_orders DROP COLUMN source;
ALTER TABLE work_orders DROP COLUMN verified_at;
ALTER TABLE work_orders DROP COLUMN verified_by;
ALTER TABLE work_orders DROP COLUMN completed_at;
