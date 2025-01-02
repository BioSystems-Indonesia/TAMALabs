import {Datagrid, List,  TextField, WithRecord, WrapperField} from "react-admin";
import Stack from "@mui/material/Stack";
import Barcode from "react-barcode";
import Typography from "@mui/material/Typography";
import CustomShowDialog from "../../component/CustomDialog.tsx";



export const ResultList = () => (
    <List >
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
            <TextField source="date" />
            <WithRecord render={(record: any) => (
                <CustomShowDialog resource="result" recordId={record.barcode} />
            )} />
        </Datagrid>
    </List>
);