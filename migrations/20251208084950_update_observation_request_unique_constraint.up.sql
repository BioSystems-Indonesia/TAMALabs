-- Drop old unique constraint that only uses (specimen_id, test_code)
-- This prevented multiple observation requests with same code but different test_type_id
DROP INDEX IF EXISTS `observation_request_uniq`;

-- Create new unique constraint that includes test_type_id
-- This allows multiple tests with same code (e.g., GLUCOSE for GDS, GDP, G2JPP)
-- but ensures uniqueness for the specific test_type selected
CREATE UNIQUE INDEX `observation_request_uniq` ON `observation_requests` (`specimen_id`, `test_type_id`);

-- Note: For backward compatibility with records that don't have test_type_id,
-- we only enforce uniqueness when test_type_id is present.
-- Records with NULL test_type_id will fall back to the old behavior.