import Box from "@mui/material/Box";
import { Create, Datagrid, Edit, List, NumberInput, SimpleForm, TextField, TextInput } from "react-admin";
import type { ActionKeys } from "../../types/props";
import { TestFilterSidebar } from "../workOrder/TestTypeFilter";

export const TestTypeList = () => (
    <List aside={<TestFilterSidebar />} title="Test Type" sort={{
        field: "id",
        order: "DESC",
    }}>
        <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="code" />
            <TextField source="category" />
            <TextField source="sub_category" />
            <TextField source="low_ref_range" label="low" />
            <TextField source="high_ref_range" label="high" />
            <TextField source="unit" />
            <TextField source="description" />
        </Datagrid>
    </List>
);


function ReferenceSection() {
    return (
        <Box sx={{ width: "100%" }}>
        </Box>
    )
}

type TestTypeFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function TestTypeForm(props: TestTypeFormProps) {
    return (
        <SimpleForm>
            <TextInput source="name" readOnly={props.readonly} />
            <TextInput source="code" readOnly={props.readonly} />
            <TextInput source="category" readOnly={props.readonly} />
            <TextInput source="sub_category" readOnly={props.readonly} />
            <NumberInput source="low_ref_range" label="low" readOnly={props.readonly} />
            <NumberInput  source="high_ref_range" label="high" readOnly={props.readonly} />
            <TextInput source="unit" readOnly={props.readonly} />
            <TextInput source="description" readOnly={props.readonly} />
        </SimpleForm>
    )
}


export function TestTypeEdit() {
    return (
        <Edit mutationMode="pessimistic" title="Edit Test Type">
            <TestTypeForm readonly={false} mode={"EDIT"} />
            <ReferenceSection />
        </Edit>
    )
}

export function TestTypeCreate() {
    return (
        <Create title="Create Test Type">
            <TestTypeForm readonly={false} mode={"CREATE"} />
            <ReferenceSection />
        </Create>
    )
}
