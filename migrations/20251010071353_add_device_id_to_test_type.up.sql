-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
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
  `device_id` integer NULL,
  `type` text NULL,
  CONSTRAINT `fk_test_types_device` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "test_types" to new temporary table "new_test_types"
INSERT INTO `new_test_types` (`id`, `name`, `code`, `alias_code`, `unit`, `low_ref_range`, `high_ref_range`, `normal_ref_string`, `decimal`, `category`, `sub_category`, `description`, `is_calculated_test`, `type`) SELECT `id`, `name`, `code`, `alias_code`, `unit`, `low_ref_range`, `high_ref_range`, `normal_ref_string`, `decimal`, `category`, `sub_category`, `description`, `is_calculated_test`, `type` FROM `test_types`;
-- drop "test_types" table after copying rows
DROP TABLE `test_types`;
-- rename temporary table "new_test_types" to "test_types"
ALTER TABLE `new_test_types` RENAME TO `test_types`;
-- create index "test_types_code" to table: "test_types"
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);
-- create index "test_type_device_id" to table: "test_types"
CREATE INDEX `test_type_device_id` ON `test_types` (`device_id`);
-- create index "test_type_alias_code" to table: "test_types"
CREATE INDEX `test_type_alias_code` ON `test_types` (`alias_code`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
