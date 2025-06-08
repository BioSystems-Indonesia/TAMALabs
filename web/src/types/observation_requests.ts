import type { TestType } from "./test_type";
import { WorkOrder } from "./work_order";

export type ObservationRequest = {
  id: number;
  test_code: string;
  test_description: string;
  requested_date: string;
  result_status: string;
  specimen_id: number;
  created_at: string;
  updated_at: string;

  test_type: TestType;
  work_order?: WorkOrder;
};

export type ObservationRequestCreateRequest = {
  test_type_id: number;
  test_type_code: string;
  specimen_type: string;
};