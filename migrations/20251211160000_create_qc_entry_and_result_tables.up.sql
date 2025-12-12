-- create "qc_entries" table
CREATE TABLE `qc_entries` (
    `id` integer NULL PRIMARY KEY AUTOINCREMENT,
    `device_id` integer NOT NULL,
    `test_type_id` integer NOT NULL,
    `qc_level` integer NOT NULL,
    `lot_number` text NOT NULL,
    `target_mean` real NOT NULL,
    `target_sd` real NULL,
    `is_active` numeric NOT NULL DEFAULT true,
    `created_by` text NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_qc_entries_device` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT `fk_qc_entries_test_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create "qc_results" table
CREATE TABLE `qc_results` (
    `id` integer NULL PRIMARY KEY AUTOINCREMENT,
    `qc_entry_id` integer NOT NULL,
    `measured_value` real NOT NULL,
    `calculated_mean` real NOT NULL,
    `calculated_sd` real NOT NULL,
    `calculated_cv` real NOT NULL,
    `sd_1` real NOT NULL,
    `sd_2` real NOT NULL,
    `sd_3` real NOT NULL,
    `result` text NOT NULL,
    `operator` text NOT NULL,
    `message_control_id` text NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_qc_results_entry` FOREIGN KEY (`qc_entry_id`) REFERENCES `qc_entries` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);

-- create indexes for qc_entries
CREATE INDEX `idx_qc_entries_device_id` ON `qc_entries` (`device_id`);

CREATE INDEX `idx_qc_entries_test_type_id` ON `qc_entries` (`test_type_id`);

CREATE INDEX `idx_qc_entries_is_active` ON `qc_entries` (`is_active`);

CREATE INDEX `idx_qc_entries_device_test_level` ON `qc_entries` (
    `device_id`,
    `test_type_id`,
    `qc_level`
);

-- create indexes for qc_results
CREATE INDEX `idx_qc_results_qc_entry_id` ON `qc_results` (`qc_entry_id`);

CREATE INDEX `idx_qc_results_created_at` ON `qc_results` (`created_at`);