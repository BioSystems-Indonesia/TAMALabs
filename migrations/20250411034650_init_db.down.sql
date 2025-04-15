-- reverse: create "test_template_test_types" table
DROP TABLE `test_template_test_types`;
-- reverse: create "test_templates" table
DROP TABLE `test_templates`;
-- reverse: create "configs" table
DROP TABLE `configs`;
-- reverse: create "units" table
DROP TABLE `units`;
-- reverse: create "devices" table
DROP TABLE `devices`;
-- reverse: create index "work_order_device_uniq" to table: "work_order_devices"
DROP INDEX `work_order_device_uniq`;
-- reverse: create "work_order_devices" table
DROP TABLE `work_order_devices`;
-- reverse: create index "specimen_uniq" to table: "specimens"
DROP INDEX `specimen_uniq`;
-- reverse: create index "specimen_barcode_uniq" to table: "specimens"
DROP INDEX `specimen_barcode_uniq`;
-- reverse: create "specimens" table
DROP TABLE `specimens`;
-- reverse: create index "work_order_created_at" to table: "work_orders"
DROP INDEX `work_order_created_at`;
-- reverse: create index "work_order_barcode" to table: "work_orders"
DROP INDEX `work_order_barcode`;
-- reverse: create "work_orders" table
DROP TABLE `work_orders`;
-- reverse: create "patients" table
DROP TABLE `patients`;
-- reverse: create "observation_results" table
DROP TABLE `observation_results`;
-- reverse: create index "observation_request_uniq" to table: "observation_requests"
DROP INDEX `observation_request_uniq`;
-- reverse: create "observation_requests" table
DROP TABLE `observation_requests`;
-- reverse: create index "test_types_code" to table: "test_types"
DROP INDEX `test_types_code`;
-- reverse: create "test_types" table
DROP TABLE `test_types`;
