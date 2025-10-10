import { WorkOrder } from "./work_order";
import { Device } from "./device";

export interface TestType {
  id: number;
  name: string;
  code: string;
  alias_code: string;
  unit: string;
  low_ref_range: number;
  high_ref_range: number;
  normal_ref_string: string; // New field for string reference values
  category: string;
  sub_category: string;
  description: string;
  is_calculated_test: boolean;
  device_id?: number;
  device?: Device;
  types: TestTypeSpecimenType[];
  decimal: number;
  work_order: WorkOrder;
}

export interface TestTypeSpecimenType {
  type: string;
}
