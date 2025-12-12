-- create "quality_controls" table
CREATE TABLE `quality_controls` (
    `id` integer NULL PRIMARY KEY AUTOINCREMENT,
    `device_id` integer NOT NULL,
    `test_type_id` integer NOT NULL,
    `qc_level` integer NOT NULL,
    `lot_number` text NOT NULL,
    `ref_value` real NOT NULL,
    `sd_value` real NOT NULL,
    `measured_value` real NOT NULL,
    `cv_value` real NOT NULL,
    `result` text NOT NULL,
    `operator` text NOT NULL,
    `device_identifier` text NULL,
    `message_control_id` text NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_quality_controls_device` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT `fk_quality_controls_test_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create index on device_id for faster queries
CREATE INDEX `idx_quality_controls_device_id` ON `quality_controls` (`device_id`);

-- create index on test_type_id for faster queries
CREATE INDEX `idx_quality_controls_test_type_id` ON `quality_controls` (`test_type_id`);

-- create index on qc_level for faster filtering
CREATE INDEX `idx_quality_controls_qc_level` ON `quality_controls` (`qc_level`);

-- create index on created_at for date-based queries
CREATE INDEX `idx_quality_controls_created_at` ON `quality_controls` (`created_at`);

-- create composite index for common queries (device + test_type)
CREATE INDEX `idx_quality_controls_device_test` ON `quality_controls` (`device_id`, `test_type_id`);