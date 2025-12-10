-- add column "medical_record_number" to table: "patients"
ALTER TABLE `patients` ADD COLUMN `medical_record_number` varchar(255) NULL DEFAULT '';
-- create index "idx_patient_medical_record_number" to table: "patients"
CREATE INDEX `idx_patient_medical_record_number` ON `patients` (`medical_record_number`);