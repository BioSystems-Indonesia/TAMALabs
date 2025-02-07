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
  result: number;
  reference: string;
  unit : string;
  abnormality: ReportDataAbnormality;
}
