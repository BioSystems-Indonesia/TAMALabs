import {
    ArrayField,
    AutocompleteInput,
    ChipField,
    Create,
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
    SingleFieldList,
    TextField,
    TextInput,
    TopToolbar
} from "react-admin";
import Divider from "@mui/material/Divider";
import {Action, ActionKeys} from "../../types/props.ts";
import FeatureList from "../../component/FeatureList.tsx";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

type SpecimenFormProps = {
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

function SpecimenForm(props: SpecimenFormProps) {
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

            <ReferenceInput source={"patient_id"} reference={"patient"}>
                <AutocompleteInput source={"patient_id"} readOnly={props.readonly}/>
            </ReferenceInput>
            <FeatureList source={"type"} types={"Specimen-type"}>
                <AutocompleteInput source={"type"} readOnly={props.readonly}/>
            </FeatureList>
            <FeatureList source={"test"} types={"Specimen-test"}>
                <AutocompleteInput source={"test"} readOnly={props.readonly}/>
            </FeatureList>
        </SimpleForm>
    )
}

export function SpecimenCreate() {
    return (
        <Create redirect={"list"}>
            <SpecimenForm mode={"CREATE"}/>
        </Create>
    )
}

export function SpecimenShow() {
    return (
        <Show>
            <SpecimenForm readonly mode={"SHOW"}/>
            <ReferenceSection/>
        </Show>
    )
}

export function SpecimenEdit() {
    return (
        <Edit>
            <SpecimenForm mode={"EDIT"}/>
        </Edit>
    )
}

const SpecimenListActions = () => (
    <TopToolbar>
        <SelectColumnsButton/>
        <FilterButton/>
        <ExportButton/>
    </TopToolbar>
);

const SpecimenFilters = [
    <SearchInput source="q" alwaysOn/>
];


export const SpecimenList = () => (
    <List actions={<SpecimenListActions/>} filters={SpecimenFilters}>
        <DatagridConfigurable bulkActionButtons={false}>
            <TextField source="id"/>
            <TextField source="barcode"/>
            <TextField source="type"/>
            <ReferenceField reference={"work-order"} source={"order_id"}/>
            <ReferenceField reference={"patient"} source={"patient_id"}/>
            <ArrayField source={"observation_requests"} label={"Observation Requests"}>
                <SingleFieldList>
                    <ChipField source={"test_code"}/>
                </SingleFieldList>
            </ArrayField>
            <DateField source="created_at" showTime/>
            <DateField source="updated_at" showTime/>
        </DatagridConfigurable>
    </List>
);
