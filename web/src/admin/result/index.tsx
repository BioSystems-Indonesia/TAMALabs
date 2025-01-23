import {Datagrid, DateField, List,  TextField, WithRecord, WrapperField} from "react-admin";
import Stack from "@mui/material/Stack";
import CustomShowDialog from "../../component/CustomDialog.tsx";
import PrintMCU from "../../component/PrintReport.tsx";


export const ResultList = () => (
    <List resource="result">
        <Datagrid>
            <TextField source="patient_id"/>
            <TextField source="patient_name"/>
            <DateField source="date" showDate showTime/>
            <WithRecord label="Show Result" render={(record: any) => (
                <CustomShowDialog resource="result" recordId={record.barcode}/>
            )} />
            <WrapperField label="Print Result">
                <PrintMCU/>
            </WrapperField>
        </Datagrid>
    </List>
);
