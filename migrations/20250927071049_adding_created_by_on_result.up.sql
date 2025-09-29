-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- add column "created_by" to table: "observation_results"
ALTER TABLE `observation_results` ADD COLUMN `created_by` integer NOT NULL DEFAULT -1;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
