-- Add method column to qc_entries table
-- QC method: statistic (statistical calculation) or manual (manual entry)
ALTER TABLE qc_entries ADD COLUMN method VARCHAR(20) NOT NULL DEFAULT 'statistic';