-- Add indexes to observation_results for faster lookups
-- Adds:
--  - idx_obs_specimen on (specimen_id)
--  - idx_obs_specimen_picked on (specimen_id, picked)
--  - idx_obs_sync on (specimen_id, last_sync, updated_at)

CREATE INDEX IF NOT EXISTS `idx_obs_specimen` ON `observation_results` (`specimen_id`);

CREATE INDEX IF NOT EXISTS `idx_obs_specimen_picked` ON `observation_results` (`specimen_id`, `picked`);

CREATE INDEX IF NOT EXISTS `idx_obs_sync` ON `observation_results` (
    `specimen_id`,
    `last_sync`,
    `updated_at`
);