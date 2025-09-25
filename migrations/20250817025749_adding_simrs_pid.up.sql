-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;

-- create "new_test_types" table (adds alias_code column)
CREATE TABLE `new_test_types` (
    `id` integer NULL PRIMARY KEY AUTOINCREMENT,
    `name` text NULL,
    `code` text NULL,
    `alias_code` text NULL,
    `unit` text NULL,
    `low_ref_range` real NULL,
    `high_ref_range` real NULL,
    `decimal` integer NULL,
    `category` text NULL,
    `sub_category` text NULL,
    `description` text NULL,
    `type` text NULL
);

-- copy rows from old table "test_types" to new temporary table "new_test_types"
-- old `test_types` may not have `alias_code`, so insert NULL for that column
INSERT INTO
    `new_test_types` (
        `id`,
        `name`,
        `code`,
        `alias_code`,
        `unit`,
        `low_ref_range`,
        `high_ref_range`,
        `decimal`,
        `category`,
        `sub_category`,
        `description`,
        `type`
    )
SELECT
    `id`,
    `name`,
    `code`,
    NULL AS `alias_code`,
    `unit`,
    `low_ref_range`,
    `high_ref_range`,
    `decimal`,
    `category`,
    `sub_category`,
    `description`,
    `type`
FROM `test_types`;

-- drop "test_types" table after copying rows
DROP TABLE `test_types`;

-- rename temporary table "new_test_types" to "test_types"
ALTER TABLE `new_test_types` RENAME TO `test_types`;

-- create index "test_types_code" to table: "test_types"
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);

-- add column "simrs_pid" to table: "patients"
ALTER TABLE `patients` ADD COLUMN `simrs_pid` text NULL;

-- create index "idx_patient_simrs_pid" to table: "patients"
CREATE UNIQUE INDEX `idx_patient_simrs_pid` ON `patients` (`simrs_pid`);

-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;