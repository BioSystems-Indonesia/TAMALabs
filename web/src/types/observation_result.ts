import type { Patient } from "./patient";
import { User } from "./user";
import { VerifiedStatus } from "./work_order";

export type EGFRCalculation = {
  value: number;
  formula: string;
  unit: string;
  category: string;
};

export type ObservationResult = {
  id: number;
  specimen_id: number;
  code: string;
  name: string;
  description: string;
  values: string[];
  type: string;
  unit: string;
  reference_range: string;
  computed_reference_range: string; // New field that uses TestType.GetReferenceRange()
  date: string;
  abnormal_flag: string[];
  comments: string;
  created_at: string;
  updated_at: string;
  test_type: TestType;
  egfr?: EGFRCalculation;
};

export type TestResult = {
  id: number;
  test_type_id: number;
  alias_code: string;
  specimen_id: number;
  test: string;
  result: string;
  formatted_result: number;
  unit: string;
  category: string;
  specimen_type: string;
  abnormal: number;
  reference_range: string;
  computed_reference_range?: string; // New field for computed reference range from TestType
  created_at: string;
  history: TestResult[] | null;
  test_type: TestType;
  egfr?: EGFRCalculation;
};

export type TestType = {
  id: number;
  name: string;
  code: string;
  alias_code: string;
  unit: string;
  low_ref_range: number;
  high_ref_range: number;
  category: string;
  sub_category: string;
  description: string;
};

export type ReportDataAbnormality =
  | "High"
  | "Low"
  | "Normal"
  | "Positive"
  | "Negative"
  | "No Data";
export type ReportData = {
  category: string;
  subCategory: string;
  parameter: string;
  alias_code?: string;
  result: string; // Changed from number to string to support non-numeric values like "3+"
  reference: string;
  unit: string;
  abnormality: ReportDataAbnormality;
};

export type Result = {
  id: number;
  status: string;
  patient_id: number;
  device_id: number;
  verified_status: VerifiedStatus;
  created_at: string;
  updated_at: string;
  patient: Patient;
  devices: null;
  created_by: number;
  last_updated_by: number;

  doctors: User[];
  analyzers: User[];

  test_result: Record<string, TestResult[]>;
  prev_id: number;
  next_id: number;
};
