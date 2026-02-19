-- Remove simrs_index column from observation_requests (rollback)
ALTER TABLE observation_requests
DROP COLUMN IF EXISTS simrs_index;
