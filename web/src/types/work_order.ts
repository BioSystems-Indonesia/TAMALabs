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
}
