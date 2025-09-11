import type { Patient } from "./patient";
import { User } from "./user";
import { VerifiedStatus } from "./work_order";

export type ObservationResult = {
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

export type TestResult = {
  id: number;
  test_type_id: number;
  specimen_id: number;
  test: string;
  result: number;
  formatted_result: number;
  unit: string;
  category: string;
  specimen_type: string;
  abnormal: number;
  reference_range: string;
  created_at: string;
  history: TestResult[] | null;
}

export type TestType ={
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
export type ReportData = {
  category: string;
  subCategory: string;
  parameter: string;
  result: number;
  reference: string;
  unit: string;
  abnormality: ReportDataAbnormality;
}

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
  next_id: number
} ;
