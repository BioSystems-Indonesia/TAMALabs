-- add column "username" to table: "devices"
ALTER TABLE `devices` ADD COLUMN `username` text NOT NULL DEFAULT '';
-- add column "password" to table: "devices"
ALTER TABLE `devices` ADD COLUMN `password` text NOT NULL DEFAULT '';
-- add column "path" to table: "devices"
ALTER TABLE `devices` ADD COLUMN `path` text NOT NULL DEFAULT '';
