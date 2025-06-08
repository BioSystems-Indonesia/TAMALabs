-- create "work_order_test_templates" table
CREATE TABLE `work_order_test_templates` (
  `work_order_id` integer NULL,
  `test_template_id` integer NULL,
  PRIMARY KEY (`work_order_id`, `test_template_id`),
  CONSTRAINT `fk_work_order_test_templates_test_template` FOREIGN KEY (`test_template_id`) REFERENCES `test_templates` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_work_order_test_templates_work_order` FOREIGN KEY (`work_order_id`) REFERENCES `work_orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
