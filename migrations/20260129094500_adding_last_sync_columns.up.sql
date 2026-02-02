-- Add last_sync column to multiple tables

ALTER TABLE `admins`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;

ALTER TABLE `observation_results`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;

ALTER TABLE `work_orders`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;

ALTER TABLE `observation_requests`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;

ALTER TABLE `test_types`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;

ALTER TABLE `test_types`
ADD COLUMN `updated_at` datetime NULL DEFAULT NULL;

ALTER TABLE `patients`
ADD COLUMN `last_sync` datetime NULL DEFAULT NULL;