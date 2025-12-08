-- Add test_type_id column to observation_requests table
-- This column will store the specific test_type ID that was selected by the user
-- Keeping test_code for backward compatibility and device communication
ALTER TABLE `observation_requests`
ADD COLUMN `test_type_id` INTEGER DEFAULT NULL;

-- Add test_type_id column to observation_results table
ALTER TABLE `observation_results`
ADD COLUMN `test_type_id` INTEGER DEFAULT NULL;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS `idx_observation_requests_test_type_id` ON `observation_requests` (`test_type_id`);

CREATE INDEX IF NOT EXISTS `idx_observation_results_test_type_id` ON `observation_results` (`test_type_id`);

-- Add foreign key constraint (optional, but recommended for data integrity)
-- Note: We use ON DELETE SET NULL because test_code is still the primary reference
-- This allows the observation to exist even if the test_type is deleted (will fall back to code lookup)