-- Remove non-unique rows before creating the unique index
DELETE FROM `units`
WHERE rowid NOT IN (
  SELECT MIN(rowid)
  FROM `units`
  GROUP BY `base`, `value`
);

-- create index "units_base_value_uniq" to table: "units"
CREATE UNIQUE INDEX `units_base_value_uniq` ON `units` (`base`, `value`);
