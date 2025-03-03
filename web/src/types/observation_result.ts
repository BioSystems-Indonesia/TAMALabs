import type { Patient } from "./patient";

export interface ObservationResult {
  id: number;
  specimen_id: number;
  code: string;
  description: string;
  values: string[];
  type: string;
  unit: string;
  reference_range: string;
  date: string;
  abnormal_flag: string[];
  comments: string;
  created_at: string;
  updated_at: string;
  test_type: TestType;
}

export interface TestResult {
  id: number;
  test_type_id: number;
  specimen_id: number;
  test: string;
  result: number;
  unit: string;
  category: string;
  abnormal: number;
  reference_range: string;
  created_at: string;
  history: TestResult[] | null;
}

export interface TestType {
  id: number;
  name: string;
  code: string;
  unit: string;
  low_ref_range: number;
  high_ref_range: number;
  category: string;
  sub_category: string;
  description: string;
}

export type ReportDataAbnormality = "High" | "Low" | "Normal" | "No Data";
export interface ReportData {
  category: string;
  subCategory: string;
  parameter: string;
  result: string;
  reference: string;
  unit: string;
  abnormality: ReportDataAbnormality;
}

export interface Result {
  id:          number;
  status:      string;
  patient_id:  number;
  device_id:   number;
  created_at:  string;
  updated_at:  string;
  patient:     Patient;
  devices:     null;
  test_result: TestResult;
}