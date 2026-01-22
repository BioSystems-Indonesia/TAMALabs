-- Add columns to track which QC result is selected for each level when multiple results exist for the same day
ALTER TABLE qc_entries ADD COLUMN level1_selected_result_id INTEGER;

ALTER TABLE qc_entries ADD COLUMN level2_selected_result_id INTEGER;

ALTER TABLE qc_entries ADD COLUMN level3_selected_result_id INTEGER;