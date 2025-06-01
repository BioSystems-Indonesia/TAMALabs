-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_observation_results" table
CREATE TABLE `new_observation_results` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `specimen_id` integer NULL,
  `code` text NULL,
  `description` text NULL,
  `values` json NULL,
  `type` text NULL,
  `unit` text NULL,
  `reference_range` text NULL,
  `date` datetime NULL,
  `abnormal_flag` json NULL,
  `comments` text NULL,
  `picked` numeric NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  CONSTRAINT `fk_observation_results_test_type` FOREIGN KEY (`code`) REFERENCES `test_types` (`code`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_specimens_observation_result` FOREIGN KEY (`specimen_id`) REFERENCES `specimens` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "observation_results" to new temporary table "new_observation_results"
INSERT INTO `new_observation_results` (`id`, `specimen_id`, `code`, `description`, `values`, `type`, `unit`, `reference_range`, `date`, `abnormal_flag`, `comments`, `picked`, `created_at`, `updated_at`) SELECT `id`, `specimen_id`, `code`, `description`, `values`, `type`, `unit`, `reference_range`, `date`, `abnormal_flag`, `comments`, `picked`, `created_at`, `updated_at` FROM `observation_results`;
-- drop "observation_results" table after copying rows
DROP TABLE `observation_results`;
-- rename temporary table "new_observation_results" to "observation_results"
ALTER TABLE `new_observation_results` RENAME TO `observation_results`;
-- create "new_devices" table
CREATE TABLE `new_devices` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL,
  `type` text NULL,
  `ip_address` text NULL,
  `port` integer NULL,
  `username` text NULL,
  `password` text NULL,
  `path` text NULL
);
-- copy rows from old table "devices" to new temporary table "new_devices"
INSERT INTO `new_devices` (`id`, `name`, `type`, `ip_address`, `port`, `username`, `password`, `path`) SELECT `id`, `name`, `type`, `ip_address`, `port`, `username`, `password`, `path` FROM `devices`;
-- drop "devices" table after copying rows
DROP TABLE `devices`;
-- rename temporary table "new_devices" to "devices"
ALTER TABLE `new_devices` RENAME TO `devices`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
