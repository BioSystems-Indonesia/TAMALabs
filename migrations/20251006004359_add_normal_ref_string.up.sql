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
  `created_by` integer NOT NULL DEFAULT -1,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  CONSTRAINT `fk_specimens_observation_result` FOREIGN KEY (`specimen_id`) REFERENCES `specimens` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_observation_results_created_by_admin` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_observation_results_test_type` FOREIGN KEY (`code`) REFERENCES `test_types` (`code`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "observation_results" to new temporary table "new_observation_results"
INSERT INTO `new_observation_results` (`id`, `specimen_id`, `code`, `description`, `values`, `type`, `unit`, `reference_range`, `date`, `abnormal_flag`, `comments`, `picked`, `created_by`, `created_at`, `updated_at`) SELECT `id`, `specimen_id`, `code`, `description`, `values`, `type`, `unit`, `reference_range`, `date`, `abnormal_flag`, `comments`, `picked`, `created_by`, `created_at`, `updated_at` FROM `observation_results`;
-- drop "observation_results" table after copying rows
DROP TABLE `observation_results`;
-- rename temporary table "new_observation_results" to "observation_results"
ALTER TABLE `new_observation_results` RENAME TO `observation_results`;
-- create "new_test_types" table
CREATE TABLE `new_test_types` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL,
  `code` text NULL,
  `alias_code` text NULL,
  `unit` text NULL,
  `low_ref_range` real NULL,
  `high_ref_range` real NULL,
  `normal_ref_string` text NULL,
  `decimal` integer NULL,
  `category` text NULL,
  `sub_category` text NULL,
  `description` text NULL,
  `is_calculated_test` numeric NULL DEFAULT false,
  `type` text NULL
);
-- copy rows from old table "test_types" to new temporary table "new_test_types"
INSERT INTO `new_test_types` (`id`, `name`, `code`, `alias_code`, `unit`, `low_ref_range`, `high_ref_range`, `decimal`, `category`, `sub_category`, `description`, `is_calculated_test`, `type`) SELECT `id`, `name`, `code`, `alias_code`, `unit`, `low_ref_range`, `high_ref_range`, `decimal`, `category`, `sub_category`, `description`, `is_calculated_test`, `type` FROM `test_types`;
-- drop "test_types" table after copying rows
DROP TABLE `test_types`;
-- rename temporary table "new_test_types" to "test_types"
ALTER TABLE `new_test_types` RENAME TO `test_types`;
-- create index "test_types_code" to table: "test_types"
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);
-- create index "test_type_alias_code" to table: "test_types"
CREATE INDEX `test_type_alias_code` ON `test_types` (`alias_code`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
