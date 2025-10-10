import type { TestResult } from "./observation_result";
import type { Patient } from "./patient";

// Raw test result data from API
export interface RawTestResult {
  abnormal: number;
  category: string;
  created_at: string;
  formatted_result: number;
  history: TestResult[] | null;
  id: number;
  picked: boolean;
  reference_range: string;
  computed_reference_range?: string; // New field for computed reference range from TestType
  result: string; // Changed from number to string to support qualitative values like "negative", "1+", etc.
  specimen_id: number;
  test: string;
  test_type_id: number;
  unit: string;

  egfr: EGFRCalculation;
}

export type EGFRCalculation = {
  value: number;
  formula: string;
  unit: string;
  category: string;
};

// API response structure for patient result history
export interface PatientResultHistoryResponse {
  patient: Patient;
  test_result: RawTestResult[];
}

// Processed data structure for DataTable display
export interface ProcessedTestResult {
  test: string;
  reference_range: string;
  computed_reference_range?: string; // New field for computed reference range from TestType
  unit: string;
  category: string;
  isCategory?: boolean; // Flag to identify category header rows
  [key: string]: string | number | boolean | undefined | EGFRCalculation; // Dynamic date columns: "2025-01-01_result", "2025-01-01_color", etc.
}

// Color types for abnormal flags
export type AbnormalColor = "default" | "error" | "secondary" | "success";

// Abnormal flag enum
export enum AbnormalFlag {
  Normal = 0,
  High = 1,
  Low = 2,
  Critical = 3,
}

// Utility types for grouping functions
export type GroupedByTest = Record<string, RawTestResult[]>;
export type GroupedByDate = Record<string, RawTestResult[]>;

// Hook return type for patient result history
export interface UsePatientResultHistoryReturn {
  data: ProcessedTestResult[];
  allDates: string[];
  isLoading: boolean;
  error: string | null;
}

// Component props
export interface PatientResultHistoryProps {
  patientId: number;
  startDate?: string;
  endDate?: string;
}
