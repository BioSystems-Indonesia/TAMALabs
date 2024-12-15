import {
    AutocompleteInput,
    BulkDeleteButton,
    Create,
    CreateButton,
    Datagrid,
    DatagridConfigurable,
    DateField,
    DateTimeInput,
    Edit,
    ExportButton,
    FilterButton,
    List,
    ReferenceField,
    ReferenceInput,
    ReferenceManyField,
    SearchInput,
    SelectColumnsButton,
    Show,
    SimpleForm,
    TextField,
    TextInput,
    TopToolbar
} from "react-admin";
import Divider from "@mui/material/Divider";
import {Action, ActionKeys} from "../../types/props.ts";
import FeatureList from "../../component/FeatureList.tsx";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import SendSpecimentToWorkOrder from "./SendSpecimentToWorkOrder.tsx";

type SpecimentFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}


function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
            <Divider sx={{my: "1rem"}}/>
            <Typography variant={"h6"}>Patient</Typography>
            <ReferenceManyField reference="patient" target="id" source={"patient_id"}>
                <Datagrid>
                    <TextField source="id"/>
                    <TextField source="first_name"/>
                    <TextField source="last_name"/>
                    <DateField source="birthdate"/>
                    <TextField source="sex"/>
                    <TextField source="location"/>
                    <DateField source="created_at" showTime/>
                    <DateField source="updated_at" showTime/>
                </Datagrid>
            </ReferenceManyField>
        </Box>
    )
}

function SpecimentForm(props: SpecimentFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
                    toolbar={props.readonly === false ? false : undefined}
                    warnWhenUnsavedChanges>
            {props.mode !== Action.CREATE && (
                <div>
                    <TextInput source={"id"} readOnly={true}/>
                    <DateTimeInput source={"created_at"} readOnly={true}/>
                    <DateTimeInput source={"updated_at"} readOnly={true}/>
                    <TextInput source="barcode" readOnly={true}/>
                    <Divider/>
                </div>
            )}

            <TextInput source={"description"}/>
            <ReferenceInput source={"patient_id"} reference={"patient"}>
                <AutocompleteInput source={"patient_id"} readOnly={props.readonly}/>
            </ReferenceInput>
            <FeatureList source={"type"} types={"speciment-type"}>
                <AutocompleteInput source={"type"} readOnly={props.readonly}/>
            </FeatureList>
            <FeatureList source={"test"} types={"speciment-test"}>
                <AutocompleteInput source={"test"} readOnly={props.readonly}/>
            </FeatureList>
        </SimpleForm>
    )
}

export function SpecimentCreate() {
    return (
        <Create redirect={"list"}>
            <SpecimentForm mode={"CREATE"}/>
        </Create>
    )
}

export function SpecimentShow() {
    return (
        <Show>
            <SpecimentForm readonly mode={"SHOW"}/>
            <ReferenceSection/>
        </Show>
    )
}

export function SpecimentEdit() {
    return (
        <Edit>
            <SpecimentForm mode={"EDIT"}/>
        </Edit>
    )
}

const SpecimentListActions = () => (
    <TopToolbar>
        <SelectColumnsButton/>
        <FilterButton/>
        <CreateButton/>
        <ExportButton/>
    </TopToolbar>
);

const SpecimentFilters = [
    <SearchInput source="q" alwaysOn/>
];

const SpecimentBulkAction = () => (
    <>
        <SendSpecimentToWorkOrder/>
        <BulkDeleteButton/>
    </>
);


export const SpecimentList = () => (
    <List actions={<SpecimentListActions/>} filters={SpecimentFilters}>
        <DatagridConfigurable bulkActionButtons={<SpecimentBulkAction/>}>
            <TextField source="id"/>
            <TextField source="description"/>
            <TextField source="barcode"/>
            <TextField source="type"/>
            <TextField source="test"/>
            <ReferenceField reference={"patient"} source={"patient_id"}/>
            <DateField source="created_at" showTime/>
            <DateField source="updated_at" showTime/>
        </DatagridConfigurable>
    </List>
);
