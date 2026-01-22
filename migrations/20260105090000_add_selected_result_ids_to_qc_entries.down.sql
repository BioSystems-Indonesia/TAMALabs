-- Remove selected result ID columns from qc_entries
ALTER TABLE qc_entries DROP COLUMN level1_selected_result_id;

ALTER TABLE qc_entries DROP COLUMN level2_selected_result_id;

ALTER TABLE qc_entries DROP COLUMN level3_selected_result_id;