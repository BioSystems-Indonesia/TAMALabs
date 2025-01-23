import {Datagrid, DateField, List,  TextField, WithRecord, WrapperField} from "react-admin";
import Stack from "@mui/material/Stack";
import Barcode from "react-barcode";
import Typography from "@mui/material/Typography";
import CustomShowDialog from "../../component/CustomDialog.tsx";
import PrintMCU from "../../component/PrintReport.tsx";


export const ResultList = () => (
    <List resource="result">
        <Datagrid>
            <TextField source="patient_id"/>
            <TextField source="patient_name"/>
            <WrapperField source={"barcode"} label={"Barcode"} textAlign={"center"}>
                <Stack>
                    <WithRecord render={(record: any) => {
                        return (
                            <Stack gap={0} justifyContent={"center"} alignItems={"center"}>
                                <Barcode value={record.barcode} displayValue={false}/>
                                <Typography
                                    className={"barcode-text"}
                                    fontSize={12}
                                    sx={{
                                        margin: 0,
                                    }}>{record.barcode}</Typography>
                            </Stack>
                        )
                    }}/>
                </Stack>
            </WrapperField>
            <DateField source="date" showDate showTime/>
            <WithRecord render={(record: any) => (
                <>
                    <CustomShowDialog resource="result" recordId={record.barcode}/>
                    <PrintMCU/>
                </>
            )} />
        </Datagrid>
    </List>
);