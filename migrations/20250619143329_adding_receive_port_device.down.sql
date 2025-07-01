-- reverse: create "new_devices" table
ALTER TABLE `devices` DROP COLUMN `receive_port`;
ALTER TABLE `devices` RENAME COLUMN `send_port` TO `port`;

