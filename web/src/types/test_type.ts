import { WorkOrder } from "./work_order";

export interface TestType {
    id:             number;
    name:           string;
    code:           string;
    unit:           string;
    low_ref_range:  number;
    high_ref_range: number;
    category:       string;
    sub_category:   string;
    description:    string;
    types: TestTypeSpecimenType[];
    decimal: number;
    work_order: WorkOrder;
}

export interface TestTypeSpecimenType {
    type: string;
}
