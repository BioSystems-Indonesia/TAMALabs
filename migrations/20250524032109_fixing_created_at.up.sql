UPDATE test_templates
SET created_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE created_at='autoCreateTime' OR updated_at='autoUpdateTime';

-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_test_templates" table
CREATE TABLE `new_test_templates` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `test_types` text NOT NULL DEFAULT '{}',
  `doctor_ids` text NOT NULL DEFAULT '{}',
  `analyzer_ids` text NOT NULL DEFAULT '{}',
  `created_by` integer NOT NULL DEFAULT 0,
  `last_updated_by` integer NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `fk_test_templates_last_updated_by_user` FOREIGN KEY (`last_updated_by`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_test_templates_created_by_user` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "test_templates" to new temporary table "new_test_templates"
INSERT INTO `new_test_templates` (`id`, `name`, `description`, `test_types`, `doctor_ids`, `analyzer_ids`, `created_by`, `last_updated_by`, `created_at`, `updated_at`) SELECT `id`, `name`, `description`, `test_types`, `doctor_ids`, `analyzer_ids`, `created_by`, `last_updated_by`, IFNULL(`created_at`, CURRENT_TIMESTAMP) AS `created_at`, IFNULL(`updated_at`, CURRENT_TIMESTAMP) AS `updated_at` FROM `test_templates`;
-- drop "test_templates" table after copying rows
DROP TABLE `test_templates`;
-- rename temporary table "new_test_templates" to "test_templates"
ALTER TABLE `new_test_templates` RENAME TO `test_templates`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
