import {Datagrid, List,  TextField, WithRecord, WrapperField} from "react-admin";
import Stack from "@mui/material/Stack";
import CustomShowDialog from "../../component/CustomDialog.tsx";
import PrintMCU from "../../component/PrintReport.tsx";


export const ResultList = () => (
    <List >
        <Datagrid>
            <TextField source="patient_id"/>
            <TextField source="patient_name"/>
            <TextField source="date" />
            <WithRecord label="Show Result" render={(record: any) => (
                <CustomShowDialog resource="result" recordId={record.barcode}/>
            )} />
            <WrapperField label="Print Result">
                <PrintMCU/>
            </WrapperField>
        </Datagrid>
    </List>
);
