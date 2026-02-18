-- Revert indexes added in 20260218095600_adding_indexes_observation_results.up.sql

DROP INDEX IF EXISTS `idx_obs_specimen`;

DROP INDEX IF EXISTS `idx_obs_specimen_picked`;

DROP INDEX IF EXISTS `idx_obs_sync`;