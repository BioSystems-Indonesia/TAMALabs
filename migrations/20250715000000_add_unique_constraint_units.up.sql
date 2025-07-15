-- create unique index "units_base_value_uniq" to table: "units"
CREATE UNIQUE INDEX `units_base_value_uniq` ON `units` (`base`, `value`);
