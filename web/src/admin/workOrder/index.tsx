import {
    AutocompleteArrayInput,
    ChipField,
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    List,
    RadioButtonGroupInput,
    ReferenceArrayInput,
    ReferenceField,
    ReferenceManyField,
    Show,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import Divider from "@mui/material/Divider";
import {Action, ActionKeys} from "../../types/props.ts";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import FeatureList from "../../component/FeatureList.tsx";

type WorkOrderFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
            <Divider sx={{my: "1rem"}}/>
            <Typography variant={"h6"}>Speciments</Typography>
            <ReferenceManyField reference={"speciment"} target={"id"} source={"speciment_ids"}>
                <Datagrid>
                    <TextField source="id"/>
                    <TextField source="description"/>
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

function WorkOrderForm(props: WorkOrderFormProps) {
    return (
        <SimpleForm>
            {props.mode !== Action.CREATE && (
                <div>
                    <TextInput source={"id"} readOnly={true}/>
                    <DateTimeInput source={"created_at"} readOnly={true}/>
                    <DateTimeInput source={"updated_at"} readOnly={true}/>
                    <FeatureList types={"work-order-status"} source={"status"}>
                        <RadioButtonGroupInput source="status" readOnly={props.readonly}/>
                    </FeatureList>
                    <Divider/>
                </div>
            )}

            <TextInput source={"description"} readOnly={props.readonly}/>
            <ReferenceArrayInput source="speciment_ids" reference={"speciment"} readOnly={props.readonly}>
                <AutocompleteArrayInput source={"speciment_ids"} readOnly={props.readonly}/>
            </ReferenceArrayInput>
        </SimpleForm>
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
            <TextField source="description"/>
            <ChipField source="status"/>
            <DateField source="created_at"/>
            <DateField source="updated_at"/>
        </Datagrid>
    </List>
);
