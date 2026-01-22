-- remove ref_min and ref_max columns from qc_entries table
ALTER TABLE `qc_entries` DROP COLUMN `ref_min`;

ALTER TABLE `qc_entries` DROP COLUMN `ref_max`;