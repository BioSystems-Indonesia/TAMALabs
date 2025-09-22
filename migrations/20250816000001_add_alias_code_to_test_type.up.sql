-- Add alias_code column to test_types
PRAGMA foreign_keys = off;

BEGIN TRANSACTION;

ALTER TABLE `test_types` RENAME TO `old_test_types`;

CREATE TABLE `test_types` (
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

INSERT INTO
    `test_types` (
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
    '' AS `alias_code`,
    `unit`,
    `low_ref_range`,
    `high_ref_range`,
    `decimal`,
    `category`,
    `sub_category`,
    `description`,
    `type`
FROM `old_test_types`;

DROP TABLE `old_test_types`;

CREATE UNIQUE INDEX IF NOT EXISTS `test_types_code` ON `test_types` (`code`);

COMMIT;

PRAGMA foreign_keys = on;