import { Card, CardContent, Chip, Typography } from "@mui/material";
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
    TextField,
    WithRecord
} from "react-admin";
import PrintMCUButton from "../../component/PrintReport";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";


const ResultFilterSidebar = () => {
    return (
        <Card sx={{
            order: -1, mr: 2, mt: 2, width: 200, minWidth: 200,
        }}>
            <CardContent>
                <FilterLiveForm>
                    <ReferenceInput source={"work_order_ids"} reference="work-order" label={"Work Order"}>
                        <AutocompleteArrayInput />
                    </ReferenceInput>
                    <ReferenceInput source={"patient_ids"} reference="patient" label={"Patient"}>
                        <AutocompleteArrayInput />
                    </ReferenceInput>
                    <BooleanInput source={"has_result"} label={"Show Only With Result"} />
                </FilterLiveForm>
            </CardContent>
        </Card>
    )
}

export const ResultDataGrid = (props: any) => {
    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"} onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )} />
            <WithRecord label="Work Order" render={(record: any) => (
                <Link to={`/work-order/${record.work_order.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.work_order.id} - ${record.work_order.status}`} color={WorkOrderChipColorMap(record.work_order.status)} />
                </Link>
            )} />
            <TextField source="barcode" />
            <WithRecord label="Request" render={(record: any) => (
                <Typography variant="body2" >
                    {record.observation_requests.length}
                </Typography>
            )} />
            <WithRecord label="Result" render={(record: any) => (
                <Typography variant="body2" >
                    {record.test_result?.length}
                </Typography>
            )} />
            <DateField source="created_at" showDate showTime />
            <WithRecord label="Print Result" render={(record: any) => (
                <PrintMCUButton results={record.test_result} patient={record.patient} workOrder={record.work_order} />
            )} />
        </Datagrid>
    )
}

export const ResultList = () => (
    <List resource="result" sort={{
        field: "id",
        order: "DESC"
    }} aside={<ResultFilterSidebar />} exporter={false} >
        <ResultDataGrid />
    </List>
);



