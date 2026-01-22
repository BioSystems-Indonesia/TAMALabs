-- drop indexes
DROP INDEX IF EXISTS `idx_quality_controls_device_test`;

DROP INDEX IF EXISTS `idx_quality_controls_created_at`;

DROP INDEX IF EXISTS `idx_quality_controls_qc_level`;

DROP INDEX IF EXISTS `idx_quality_controls_test_type_id`;

DROP INDEX IF EXISTS `idx_quality_controls_device_id`;

-- drop table
DROP TABLE IF EXISTS `quality_controls`;