-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_devices" table
CREATE TABLE `new_devices` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL,
  `type` text NULL,
  `ip_address` text NULL,
  `send_port` integer NULL,
  `receive_port` integer NULL,
  `username` text NULL,
  `password` text NULL,
  `path` text NULL
);
-- copy rows from old table "devices" to new temporary table "new_devices"
INSERT INTO `new_devices` (`id`, `name`, `type`, `ip_address`,`send_port`, `receive_port`, `username`, `password`, `path`) 
SELECT `id`, `name`, `type`, `ip_address`, `port`, `id`+10123, `username`, `password`, `path` FROM `devices`;
-- drop "devices" table after copying rows
DROP TABLE `devices`;
-- rename temporary table "new_devices" to "devices"
ALTER TABLE `new_devices` RENAME TO `devices`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
