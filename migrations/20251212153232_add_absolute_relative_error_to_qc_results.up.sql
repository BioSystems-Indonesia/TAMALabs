-- Add absolute_error and relative_error columns to qc_results table
ALTER TABLE qc_results
ADD COLUMN absolute_error REAL NOT NULL DEFAULT 0;

ALTER TABLE qc_results
ADD COLUMN relative_error REAL NOT NULL DEFAULT 0;