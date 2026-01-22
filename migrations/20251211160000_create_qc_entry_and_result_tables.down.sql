-- drop indexes
DROP INDEX IF EXISTS `idx_qc_results_created_at`;

DROP INDEX IF EXISTS `idx_qc_results_qc_entry_id`;

DROP INDEX IF EXISTS `idx_qc_entries_device_test_level`;

DROP INDEX IF EXISTS `idx_qc_entries_is_active`;

DROP INDEX IF EXISTS `idx_qc_entries_test_type_id`;

DROP INDEX IF EXISTS `idx_qc_entries_device_id`;

-- drop tables
DROP TABLE IF EXISTS `qc_results`;

DROP TABLE IF EXISTS `qc_entries`;