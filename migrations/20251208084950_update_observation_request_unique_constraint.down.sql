-- Rollback: Restore old unique constraint on (specimen_id, test_code)
DROP INDEX IF EXISTS `observation_request_uniq`;

CREATE UNIQUE INDEX `observation_request_uniq` ON `observation_requests` (`specimen_id`, `test_code`);