-- +migrate Up
-- create "test_type_devices" junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS `test_type_devices` (
    `test_type_id` integer NOT NULL,
    `device_id` integer NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`test_type_id`, `device_id`),
    CONSTRAINT `fk_test_type_devices_test_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `fk_test_type_devices_device` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);

-- create index for faster lookups
CREATE INDEX IF NOT EXISTS `idx_test_type_devices_test_type_id` ON `test_type_devices` (`test_type_id`);

CREATE INDEX IF NOT EXISTS `idx_test_type_devices_device_id` ON `test_type_devices` (`device_id`);

-- +migrate Down
DROP TABLE IF EXISTS `test_type_devices`;