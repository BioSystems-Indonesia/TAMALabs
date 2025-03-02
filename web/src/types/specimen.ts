import type { ObservationRequest } from "./observation_requests";

export interface Specimen {
    id:                   number;
    specimen_hl7_id:      string;
    patient_id:           number;
    order_id:             number;
    type:                 string;
    collection_date:      string;
    received_date:        string;
    source:               string;
    condition:            string;
    method:               string;
    comments:             string;
    barcode:              string;
    created_at:           string;
    updated_at:           string;
    observation_result:   any[];
    observation_requests: ObservationRequest[];
    test_result:          null;
}