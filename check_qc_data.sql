-- Check QC Results for lot 0003980
SELECT
    id,
    qc_entry_id,
    measured_value,
    calculated_mean,
    calculated_sd,
    calculated_cv,
    result,
    created_at
FROM qc_results
WHERE
    qc_entry_id = 11
ORDER BY id ASC;

-- Check how many results per entry
SELECT
    qc_entry_id,
    COUNT(*) as total_count,
    MIN(measured_value) as min_value,
    MAX(measured_value) as max_value,
    AVG(measured_value) as avg_value
FROM qc_results
WHERE
    qc_entry_id IN (10, 11)
GROUP BY
    qc_entry_id;

-- Check QC Entry details
SELECT
    id,
    lot_number,
    qc_level,
    target_mean,
    target_sd
FROM qc_entries
WHERE
    id IN (10, 11);