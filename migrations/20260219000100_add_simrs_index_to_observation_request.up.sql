-- Add simrs_index to observation_requests to store index provided by external SIMRS (Nuha)
ALTER TABLE observation_requests
ADD COLUMN simrs_index INTEGER DEFAULT NULL;

-- no-op index required; column is small and only used for mapping when sending results back to SIMRS
