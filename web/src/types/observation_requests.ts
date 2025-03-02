import type { TestType } from "./test_type";

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
};
