import { Device } from "./device";
import { TestResult } from "./observation_result";
import type { Patient } from "./patient";
import type { Specimen } from "./specimen";
import { User } from "./user";

export type WorkOrderStatus =
  | "PENDING"
  | "INCOMPLETE_SEND"
  | "COMPLETE"
  | "CANCELLED";
export type VerifiedStatus = "PENDING" | "VERIFIED" | "REJECTED" | "";

export type WorkOrder = {
  id: number;
  status: WorkOrderStatus;
  device_id: number;
  verified_status: VerifiedStatus;
  barcode: string;
  barcode_simrs: string;
  medical_record_number: string;
  visit_number: string;
  specimen_collection_date: string;
  result_release_date: string;
  diagnosis: string;
  simrs_sent_status: string;
  simrs_sent_at: string | null;
  devices: Device[] | null;
  created_at: string;
  updated_at: string;
  created_by: number;
  last_updated_by: number;

  patient: Patient;
  doctors: User[];
  analyst: User[];
  specimen_list: Specimen[];

  test_result: TestResult[];
  total_request: number;
  total_result_filled: number;
  percent_complete: number;
  have_complete_data: boolean;
  analyzers: any[];
};

export const workOrderStatusShowCancel: WorkOrderStatus[] = [
  "PENDING",
  "INCOMPLETE_SEND",
] as const;
export const workOrderStatusDontShowRun: WorkOrderStatus[] = [
  "INCOMPLETE_SEND",
] as const;
