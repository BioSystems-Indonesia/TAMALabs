-- reverse: create "admin_roles" table
DROP TABLE `admin_roles`;
-- reverse: create index "idx_role_name_uniq" to table: "roles"
DROP INDEX `idx_role_name_uniq`;
-- reverse: create "roles" table
DROP TABLE `roles`;
-- reverse: create index "idx_admin_fullname" to table: "admins"
DROP INDEX `idx_admin_fullname`;
-- reverse: create index "idx_admin_username" to table: "admins"
DROP INDEX `idx_admin_username`;
-- reverse: create index "idx_admins_email" to table: "admins"
DROP INDEX `idx_admins_email`;
-- reverse: create "admins" table
DROP TABLE `admins`;