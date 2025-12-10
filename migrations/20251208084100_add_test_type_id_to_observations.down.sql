-- Rollback: Remove test_type_id columns and indexes

-- Drop indexes
DROP INDEX IF EXISTS `idx_observation_requests_test_type_id`;

DROP INDEX IF EXISTS `idx_observation_results_test_type_id`;

-- Drop columns
ALTER TABLE `observation_requests` DROP COLUMN `test_type_id`;

ALTER TABLE `observation_results` DROP COLUMN `test_type_id`;