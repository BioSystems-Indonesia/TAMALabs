-- reverse: add column "path" to table: "devices"
ALTER TABLE `devices` DROP COLUMN `path`;
-- reverse: add column "password" to table: "devices"
ALTER TABLE `devices` DROP COLUMN `password`;
-- reverse: add column "username" to table: "devices"
ALTER TABLE `devices` DROP COLUMN `username`;