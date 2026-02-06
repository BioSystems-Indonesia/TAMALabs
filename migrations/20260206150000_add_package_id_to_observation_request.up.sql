-- Add package_id column to observation_requests table
ALTER TABLE observation_requests
ADD COLUMN package_id INT DEFAULT NULL;