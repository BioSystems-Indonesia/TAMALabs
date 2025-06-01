-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_work_orders" table
CREATE TABLE `new_work_orders` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `status` text NOT NULL,
  `patient_id` integer NULL DEFAULT 0,
  `device_id` blob NOT NULL DEFAULT '0',
  `barcode` varchar NULL DEFAULT '',
  `verified_status` varchar NULL DEFAULT '',
  `created_by` integer NULL DEFAULT 0,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  CONSTRAINT `fk_work_orders_patient` FOREIGN KEY (`patient_id`) REFERENCES `patients` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_work_orders_created_by_user` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- copy rows from old table "work_orders" to new temporary table "new_work_orders"
INSERT INTO `new_work_orders` (`id`, `status`, `patient_id`, `device_id`, `barcode`, `created_at`, `updated_at`) SELECT `id`, `status`, `patient_id`, `device_id`, `barcode`, `created_at`, `updated_at` FROM `work_orders`;
-- drop "work_orders" table after copying rows
DROP TABLE `work_orders`;
-- rename temporary table "new_work_orders" to "work_orders"
ALTER TABLE `new_work_orders` RENAME TO `work_orders`;
-- create index "work_order_barcode" to table: "work_orders"
CREATE UNIQUE INDEX `work_order_barcode` ON `work_orders` (`barcode`);
-- create index "work_order_created_at" to table: "work_orders"
CREATE INDEX `work_order_created_at` ON `work_orders` (`created_at`);
-- create "work_order_analyzers" table
CREATE TABLE `work_order_analyzers` (
  `work_order_id` integer NULL,
  `admin_id` integer NULL,
  PRIMARY KEY (`work_order_id`, `admin_id`),
  CONSTRAINT `fk_work_order_analyzers_admin` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_work_order_analyzers_work_order` FOREIGN KEY (`work_order_id`) REFERENCES `work_orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "work_order_doctors" table
CREATE TABLE `work_order_doctors` (
  `work_order_id` integer NULL,
  `admin_id` integer NULL,
  PRIMARY KEY (`work_order_id`, `admin_id`),
  CONSTRAINT `fk_work_order_doctors_admin` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_work_order_doctors_work_order` FOREIGN KEY (`work_order_id`) REFERENCES `work_orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
