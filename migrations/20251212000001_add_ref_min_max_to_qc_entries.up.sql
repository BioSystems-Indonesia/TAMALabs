-- add ref_min and ref_max columns to qc_entries table
ALTER TABLE `qc_entries`
ADD COLUMN `ref_min` real NOT NULL DEFAULT 0;

ALTER TABLE `qc_entries`
ADD COLUMN `ref_max` real NOT NULL DEFAULT 0;