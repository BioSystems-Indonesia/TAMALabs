import {
    ChipField,
    Create,
    Datagrid,
    DateField,
    Edit,
    List,
    ReferenceField,
    ReferenceManyField,
    Show,
    TextField
} from "react-admin";
import Divider from "@mui/material/Divider";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import WorkOrderForm from "./Form.tsx";


function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
            <Divider sx={{my: "1rem"}}/>
            <Typography variant={"h6"}>Specimens</Typography>
            <ReferenceManyField reference={"Specimen"} target={"id"} source={"Specimen_ids"}>
                <Datagrid>
                    <TextField source="id"/>
                    <TextField source="barcode"/>
                    <TextField source="type"/>
                    <ReferenceField reference={"patient"} source={"patient_id"}/>
                    <DateField source="created_at" showTime/>
                    <DateField source="updated_at" showTime/>
                </Datagrid>
            </ReferenceManyField>
        </Box>
    )
}


export function WorkOrderCreate() {
    return (
        <Create redirect={"list"}>
            <WorkOrderForm mode={"CREATE"}/>
        </Create>
    )
}

export function WorkOrderShow() {
    return (
        <Show>
            <WorkOrderForm readonly mode={"SHOW"}/>
            <ReferenceSection/>
        </Show>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit>
            <WorkOrderForm mode={"EDIT"}/>
        </Edit>
    )
}

export const WorkOrderList = () => (
    <List>
        <Datagrid>
            <TextField source="id"/>
            <ChipField source="status"/>
            <DateField source="created_at"/>
            <DateField source="updated_at"/>
        </Datagrid>
    </List>
);
