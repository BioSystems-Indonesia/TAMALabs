-- Remove absolute_error and relative_error columns from qc_results table
ALTER TABLE qc_results DROP COLUMN absolute_error;

ALTER TABLE qc_results DROP COLUMN relative_error;