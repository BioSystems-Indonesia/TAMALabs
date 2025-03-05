import { Card, CardContent, Chip, Stack, Typography } from "@mui/material";
import {
    AutocompleteArrayInput,
    BooleanInput,
    Datagrid,
    DateField,
    FilterLiveForm,
    Link,
    List,
    NumberField,
    ReferenceInput,
    WithRecord
} from "react-admin";
import CustomDateInput from "../../component/CustomDateInput";
import PrintReportButton from "../../component/PrintReport";
import type { WorkOrder } from "../../types/work_order";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { FilledPercentChip } from "./component";



export const ResultList = () => (
    <List 
        resource="result" 
        sort={{ field: "id", order: "DESC" }} 
        filters={ResultFilters}
        storeKey={false} 
        exporter={false} 
        disableSyncWithLocation 
        >
        <ResultDataGrid />
    </List>
);

export const ResultDataGrid = (props: any) => {
    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"} onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )} />
            <WithRecord label="Request" render={(record: any) => (
                <Link to={`/work-order/${record.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.id} - ${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                </Link>
            )} />
            <WithRecord label="Request" render={(record: WorkOrder) => (
                <Typography variant="body2" >
                    {record.total_request}
                </Typography>
            )} />
            <WithRecord label="Result" render={(record: WorkOrder) => (
                <Typography variant="body2" >
                    {record.total_result_filled}
                </Typography>
            )} />
            <WithRecord label="Filled" render={(record: WorkOrder) => (
                <FilledPercentChip percent={record.percent_complete} />
            )} />
            <DateField source="created_at" showDate showTime />
            <WithRecord label="Print Result" render={(record: any) => (
                <PrintReportButton results={record.test_result} patient={record.patient} workOrder={record} />
            )} />
        </Datagrid>
    )
}

const ResultFilters = [
    <ReferenceInput source={"patient_ids"} reference="patient" label={"Patient"} >
        <AutocompleteArrayInput />
    </ReferenceInput>,
    <CustomDateInput label={"Created At Start"} source="created_at_start" />,
    <CustomDateInput label={"Created At End"} source="created_at_end" />,
    <BooleanInput source={"has_result"} label={"Show Only With Result"} />,
];
