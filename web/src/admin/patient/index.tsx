import { Stack } from '@mui/material';
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import {
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    FilterLiveForm, FilterLiveSearch,
    List,
    ReferenceManyField,
    required,
    SelectInput,
    Show,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import FeatureList from "../../component/FeatureList.tsx";
import SideFilter from '../../component/SideFilter.tsx';
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import { Action, ActionKeys } from "../../types/props.ts";
import { ResultDataGrid } from "../result/index.tsx";

export type PatientFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function ReferenceSection() {
    return (
        <Box sx={{ width: "100%", padding: 2 }}>
            <Divider sx={{ my: "1rem" }} />
            <Typography variant={"subtitle1"}>Result</Typography>
            <ReferenceManyField label={"Result"} reference="result" target="patient_ids">
                <ResultDataGrid />
            </ReferenceManyField>
        </Box>
    )
}

export function PatientFormField(props: PatientFormProps) {
    return (
        <>
            {props.mode !== Action.CREATE && (
                <Box sx={{
                    width: "100%",
                }}>
                    <TextInput source={"id"} readOnly={true} />
                    <Stack direction={"row"} gap={5} width={"100%"}>
                        <DateTimeInput source={"created_at"} readOnly={true} />
                        <DateTimeInput source={"updated_at"} readOnly={true} />
                    </Stack>
                    <Divider />
                </Box>
            )}

            <Typography variant={"subtitle1"}>Required</Typography>
            <Divider sx={{ my: "1rem" }} />
            <Stack direction={"row"} gap={5} width={"100%"}>
                <TextInput source="first_name" validate={[required()]} readOnly={props.readonly} />
                <TextInput source="last_name" readOnly={props.readonly} />
            </Stack>
            <Stack direction={"row"} gap={3} width={"100%"}>
                <CustomDateInput source={"birthdate"} label={"Birth Date"} required sx={{
                    maxWidth: null
                }} />
                <FeatureList source={"sex"} types={"sex"}>
                    <SelectInput source="sex" validate={[required()]} readOnly={props.readonly} />
                </FeatureList>
            </Stack>
            <Typography variant={"subtitle1"}>Optional</Typography>
            <Divider sx={{ my: "1rem" }} />
            <Stack direction={"row"} gap={5} width={"100%"}>
                <TextInput source="phone_number" readOnly={props.readonly} />
                <TextInput source="location" readOnly={props.readonly} />
            </Stack>
            <TextInput source="address" readOnly={props.readonly} />
        </>
    )

}

export function PatientForm(props: PatientFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
            toolbar={props.readonly === true ? false : undefined}
            warnWhenUnsavedChanges
        >
            <PatientFormField {...props} />
        </SimpleForm>
    )
}

export function PatientCreate() {
    const redirect = useRefererRedirect("show");

    return (
        <Create redirect={redirect} resource="patient">
            <PatientForm mode={"CREATE"} />
        </Create>
    )
}

export function PatientShow() {
    return (
        <Show resource="patient">
            <PatientForm readonly mode={"SHOW"} />
            <ReferenceSection />
        </Show>
    )
}

export function PatientEdit() {
    return (
        <Edit resource="patient">
            <PatientForm mode={"EDIT"} />
        </Edit>
    )
}

const PatientFilterSidebar = () => (
    <SideFilter>
        <FilterLiveSearch />
        <FilterLiveForm debounce={1500}>
            <Stack sx={{
                marginTop: "1rem"
            }}>
                <CustomDateInput source={"birthdate"} label={"Birth Date"} clearable size="small" />
            </Stack>
        </FilterLiveForm>
    </SideFilter>
);


export const PatientList = () => (
    <List aside={<PatientFilterSidebar />} sort={{
        field: "id",
        order: "DESC"
    }}
        storeKey={false} exporter={false}
        sx={{
            '& .RaList-content': {
                backgroundColor: 'background.paper',
                padding: 2,
                borderRadius: 1,
            },
        }}>
        <Datagrid>
            <TextField source="id" />
            <TextField source="first_name" />
            <TextField source="last_name" />
            <DateField source="birthdate" locales={["id-ID"]} />
            <TextField source="sex" />
            <TextField source="location" />
            <DateField source="created_at" showTime />
            <DateField source="updated_at" showTime />
        </Datagrid>
    </List>
);
