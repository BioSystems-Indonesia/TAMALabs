-- Add created_by column to qc_results table for audit trail
ALTER TABLE qc_results
ADD COLUMN created_by VARCHAR(255) NOT NULL DEFAULT '';

-- Set value from operator for existing records
UPDATE qc_results SET created_by = operator WHERE created_by = '';