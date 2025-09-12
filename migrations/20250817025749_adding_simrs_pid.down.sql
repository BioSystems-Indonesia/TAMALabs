-- reverse: create index "idx_patient_simrs_pid" to table: "patients"
DROP INDEX `idx_patient_simrs_pid`;
-- reverse: add column "simrsp_id" to table: "patients"
ALTER TABLE `patients` DROP COLUMN `simrsp_id`;
-- reverse: create index "test_types_code" to table: "test_types"
DROP INDEX `test_types_code`;
-- reverse: create "new_test_types" table
DROP TABLE `new_test_types`;
