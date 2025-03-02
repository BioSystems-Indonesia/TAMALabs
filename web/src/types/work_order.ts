import type { Patient } from "./patient";
import type { Specimen } from "./specimen";

export interface WorkOrder {
  id: number;
  status: string;
  device_id: number;
  created_at: string;
  updated_at: string;

  patient: Patient;
  specimen_list: Specimen[];

  total_request: number;
  total_result_filled: number;
  percent_complete: number;
  have_complete_data: boolean;
}
