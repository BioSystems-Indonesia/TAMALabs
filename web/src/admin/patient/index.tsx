import {
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    FilterListSection,
    FilterLiveForm,
    FilterLiveSearch,
    List,
    RadioButtonGroupInput,
    ReferenceManyField,
    required,
    Show,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import Divider from "@mui/material/Divider";
import {Action, ActionKeys} from "../../types/props.ts";
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import FeatureList from "../../component/FeatureList.tsx";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";

type PatientFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
            <Divider sx={{my: "1rem"}}/>
            <Typography variant={"h6"}>Specimens</Typography>
            <ReferenceManyField label={"Specimens"} reference="Specimen" target="patient_id">
                <Datagrid>
                    <TextField source="id"/>
                    <TextField source="barcode"/>
                    <TextField source="type"/>
                    <DateField source="created_at" showTime/>
                    <DateField source="updated_at" showTime/>
                </Datagrid>
            </ReferenceManyField>
        </Box>
    )
}

function PatientForm(props: PatientFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
                    toolbar={props.readonly === true ? false : undefined}
                    warnWhenUnsavedChanges
        >
            {props.mode !== Action.CREATE && (
                <div>
                    <TextInput source={"id"} readOnly={true}/>
                    <DateTimeInput source={"created_at"} readOnly={true}/>
                    <DateTimeInput source={"updated_at"} readOnly={true}/>
                    <Divider/>
                </div>
            )}

            <TextInput source="first_name" validate={[required()]} readOnly={props.readonly}/>
            <TextInput source="last_name" validate={[required()]} readOnly={props.readonly}/>
            <CustomDateInput source={"birthdate"} label={"Birth Date"} required/>
            <FeatureList source={"sex"} types={"sex"}>
                <RadioButtonGroupInput source="sex" validate={[required()]} readOnly={props.readonly}/>
            </FeatureList>
            <TextInput source="phone_number" readOnly={props.readonly}/>
            <TextInput source="location" readOnly={props.readonly}/>
            <TextInput source="address" readOnly={props.readonly}/>
        </SimpleForm>
    )
}

export function PatientCreate() {
    return (
        <Create redirect={"list"}>
            <PatientForm mode={"CREATE"}/>
        </Create>
    )
}

export function PatientShow() {
    return (
        <Show>
            <PatientForm readonly mode={"SHOW"}/>
            <ReferenceSection/>
        </Show>
    )
}

export function PatientEdit() {
    return (
        <Edit>
            <PatientForm mode={"EDIT"}/>
        </Edit>
    )
}

const PatientFilterSidebar = () => (
    <Card sx={{order: -1, mr: 2, mt: 2, width: 300}}>
        <CardContent>
            <FilterLiveSearch/>
            <FilterListSection label="Birth Date" icon={<CalendarMonthIcon/>}>
                <FilterLiveForm debounce={1500}>
                    <CustomDateInput source={"birthdate"} label={"Birth Date"} clearable/>
                </FilterLiveForm>
            </FilterListSection>
        </CardContent>
    </Card>
);


export const PatientList = () => (
    <List aside={<PatientFilterSidebar/>}>
        <Datagrid>
            <TextField source="id"/>
            <TextField source="first_name"/>
            <TextField source="last_name"/>
            <DateField source="birthdate" locales={["id-ID"]}/>
            <TextField source="sex"/>
            <TextField source="location"/>
            <DateField source="created_at" showTime/>
            <DateField source="updated_at" showTime/>
        </Datagrid>
    </List>
);