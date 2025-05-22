-- create "admins" table
CREATE TABLE `admins` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `fullname` text NULL,
  `username` text NULL,
  `email` text NULL,
  `password_hash` text NULL,
  `is_active` numeric NULL DEFAULT true,
  `created_at` datetime NULL,
  `updated_at` datetime NULL
);
-- create index "idx_admins_email" to table: "admins"
CREATE UNIQUE INDEX `idx_admins_email` ON `admins` (`email`);
-- create index "idx_admin_username" to table: "admins"
CREATE UNIQUE INDEX `idx_admin_username` ON `admins` (`username`);
-- create index "idx_admin_fullname" to table: "admins"
CREATE INDEX `idx_admin_fullname` ON `admins` (`fullname`);
-- create "roles" table
CREATE TABLE `roles` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `description` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL
);
-- create index "idx_role_name_uniq" to table: "roles"
CREATE UNIQUE INDEX `idx_role_name_uniq` ON `roles` (`name`);
-- create "admin_roles" table
CREATE TABLE `admin_roles` (
  `role_id` integer NULL,
  `admin_id` integer NULL,
  PRIMARY KEY (`role_id`, `admin_id`),
  CONSTRAINT `fk_admin_roles_admin` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_admin_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
