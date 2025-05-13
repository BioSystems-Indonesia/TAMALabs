import { Device } from "./device";
import type { Patient } from "./patient";
import type { Specimen } from "./specimen";

type WorkOrderStatus = "PENDING" | "INCOMPLETE_SEND" | "COMPLETE" | "CANCELLED"

export interface WorkOrder {
  id: number;
  status: WorkOrderStatus;
  device_id: number;
  devices: Device[] | null;
  created_at: string;
  updated_at: string;

  patient: Patient;
  specimen_list: Specimen[];

  total_request: number;
  total_result_filled: number;
  percent_complete: number;
  have_complete_data: boolean;
}

export const workOrderStatusShowCancel: WorkOrderStatus[] = ["PENDING", "INCOMPLETE_SEND"] as const
export const workOrderStatusDontShowRun: WorkOrderStatus[] = ["INCOMPLETE_SEND"] as const