-- Revert added last_sync columns

ALTER TABLE `admins` DROP COLUMN `last_sync`;

ALTER TABLE `observation_results` DROP COLUMN `last_sync`;

ALTER TABLE `work_orders` DROP COLUMN `last_sync`;

ALTER TABLE `observation_requests` DROP COLUMN `last_sync`;

ALTER TABLE `test_types` DROP COLUMN `last_sync`;

ALTER TABLE `test_types` DROP COLUMN `updated_at`;

ALTER TABLE `patients` DROP COLUMN `last_sync`;