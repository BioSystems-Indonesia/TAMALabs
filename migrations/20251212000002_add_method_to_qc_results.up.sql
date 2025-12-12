-- Add method column to qc_results table
-- QC method: statistic (from device) or manual (manual entry)
ALTER TABLE qc_results
ADD COLUMN method VARCHAR(20) NOT NULL DEFAULT 'statistic';