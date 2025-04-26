-- create "test_types" table
CREATE TABLE `test_types` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL,
  `code` text NULL,
  `unit` text NULL,
  `low_ref_range` real NULL,
  `high_ref_range` real NULL,
  `decimal` integer NULL,
  `category` text NULL,
  `sub_category` text NULL,
  `description` text NULL,
  `type` text NULL
);
-- create index "test_types_code" to table: "test_types"
CREATE UNIQUE INDEX `test_types_code` ON `test_types` (`code`);
-- create "observation_requests" table
CREATE TABLE `observation_requests` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `test_code` text NOT NULL,
  `test_description` text NULL,
  `requested_date` datetime NULL,
  `result_status` text NULL,
  `specimen_id` integer NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  CONSTRAINT `fk_specimens_observation_request` FOREIGN KEY (`specimen_id`) REFERENCES `specimens` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_observation_requests_test_type` FOREIGN KEY (`test_code`) REFERENCES `test_types` (`code`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "observation_request_uniq" to table: "observation_requests"
CREATE UNIQUE INDEX `observation_request_uniq` ON `observation_requests` (`specimen_id`, `test_code`);
-- create "observation_results" table
CREATE TABLE `observation_results` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `specimen_id` integer NULL,
  `code` integer NULL,
  `description` text NULL,
  `values` json NULL,
  `type` text NULL,
  `unit` text NULL,
  `reference_range` text NULL,
  `date` datetime NULL,
  `abnormal_flag` json NULL,
  `comments` text NULL,
  `picked` numeric NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  CONSTRAINT `fk_specimens_observation_result` FOREIGN KEY (`specimen_id`) REFERENCES `specimens` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_observation_results_test_type` FOREIGN KEY (`code`) REFERENCES `test_types` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "patients" table
CREATE TABLE `patients` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `first_name` text NOT NULL,
  `last_name` text NOT NULL,
  `birthdate` datetime NOT NULL,
  `sex` text NOT NULL,
  `phone_number` text NOT NULL,
  `location` text NOT NULL,
  `address` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);
-- create "work_orders" table
CREATE TABLE `work_orders` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `status` text NOT NULL,
  `patient_id` integer NULL DEFAULT 0,
  `device_id` blob NOT NULL DEFAULT '0',
  `created_at` datetime NULL,
  `barcode` varchar NULL DEFAULT '',
  `updated_at` datetime NULL,
  CONSTRAINT `fk_work_orders_patient` FOREIGN KEY (`patient_id`) REFERENCES `patients` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "work_order_barcode" to table: "work_orders"
CREATE UNIQUE INDEX `work_order_barcode` ON `work_orders` (`barcode`);
-- create index "work_order_created_at" to table: "work_orders"
CREATE INDEX `work_order_created_at` ON `work_orders` (`created_at`);
-- create "specimens" table
CREATE TABLE `specimens` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `patient_id` integer NOT NULL,
  `order_id` integer NOT NULL,
  `type` text NOT NULL,
  `collection_date` text NOT NULL,
  `received_date` datetime NOT NULL,
  `source` text NOT NULL,
  `condition` text NOT NULL,
  `method` text NOT NULL,
  `comments` text NOT NULL,
  `barcode` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  CONSTRAINT `fk_patients_specimen` FOREIGN KEY (`patient_id`) REFERENCES `patients` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_work_orders_specimen` FOREIGN KEY (`order_id`) REFERENCES `work_orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "specimen_barcode_uniq" to table: "specimens"
CREATE UNIQUE INDEX `specimen_barcode_uniq` ON `specimens` (`barcode`);
-- create index "specimen_uniq" to table: "specimens"
CREATE UNIQUE INDEX `specimen_uniq` ON `specimens` (`order_id`, `patient_id`, `type`);
-- create "work_order_devices" table
CREATE TABLE `work_order_devices` (
  `work_order_id` integer NOT NULL,
  `device_id` integer NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);
-- create index "work_order_device_uniq" to table: "work_order_devices"
CREATE UNIQUE INDEX `work_order_device_uniq` ON `work_order_devices` (`work_order_id`, `device_id`);
-- create "devices" table
CREATE TABLE `devices` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL,
  `type` text NULL,
  `ip_address` text NULL,
  `port` integer NULL
);
-- create "units" table
CREATE TABLE `units` (
  `base` text NULL,
  `value` text NULL
);
-- create "configs" table
CREATE TABLE `configs` (
  `id` text NULL,
  `value` text NOT NULL,
  PRIMARY KEY (`id`)
);
-- create "test_templates" table
CREATE TABLE `test_templates` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `test_types` text NOT NULL DEFAULT '{}'
);
-- create "test_template_test_types" table
CREATE TABLE `test_template_test_types` (
  `test_template_id` integer NULL,
  `test_type_id` integer NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`test_template_id`, `test_type_id`)
);

-- create "sequence_daily" table
CREATE TABLE `sequence_daily` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `sequence_type` text NOT NULL,
  `current_value` integer NULL,
  `last_updated` datetime NULL
);
