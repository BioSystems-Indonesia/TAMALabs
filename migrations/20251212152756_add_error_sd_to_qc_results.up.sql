-- Add error_sd column to qc_results table for Levey-Jennings chart
ALTER TABLE qc_results ADD COLUMN error_sd REAL NOT NULL DEFAULT 0;