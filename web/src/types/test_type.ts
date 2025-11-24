import { WorkOrder } from "./work_order";
import { Device } from "./device";

export interface TestType {
  id: number;
  name: string;
  code: string;
  alias_code: string;
  loinc_code: string;
  alternative_codes?: string[]; // Alternative codes from different devices
  unit: string;
  low_ref_range: number;
  high_ref_range: number;
  normal_ref_string: string; // Default/global reference value string
  category: string;
  sub_category: string;
  description: string;
  is_calculated_test: boolean;
  device_id?: number; // Deprecated: use device_ids instead
  device?: Device; // Deprecated: use devices instead
  device_ids?: number[]; // Array of device IDs for many-to-many relationship
  devices?: Device[]; // Array of devices for many-to-many relationship
  types: TestTypeSpecimenType[];
  decimal: number;
  work_order: WorkOrder;
  specific_ref_ranges?: SpecificReferenceRange[]; // Array of specific reference ranges
}

export interface TestTypeSpecimenType {
  type: string;
}

export interface SpecificReferenceRange {
  gender?: string | null; // "M", "F", or null for all genders
  age_min?: number | null; // minimum age in years
  age_max?: number | null; // maximum age in years
  low_ref_range?: number | null; // numeric low value
  high_ref_range?: number | null; // numeric high value
  normal_ref_string?: string | null; // or string value like "Negative", "Positive"
}
