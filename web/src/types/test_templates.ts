import { ObservationRequest } from "./observation_requests"
import { TestType } from "./test_type"

export interface TestTemplateDiff {
  ToCreate: ObservationRequest[]
  ToDelete: ObservationRequest[]
}

export interface TestTemplate {
  id: number
  name: string
  description: string
  created_by: number
  last_updated_by: number
  created_at: string
  updated_at: string
  test_types: TestType[]
  doctor_ids: number[]
  analyzer_ids: number
}
