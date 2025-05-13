import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import {
    AutocompleteInput,
    Create,
    Datagrid,
    Edit,
    FilterLiveSearch,
    FormDataConsumer,
    List,
    NumberInput,
    PasswordInput,
    required,

    Show,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import FeatureList from "../../component/FeatureList.tsx";
import { Action, ActionKeys } from "../../types/props.ts";
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import SideFilter from "../../component/SideFilter.tsx";
import { DeviceType, DeviceTypeValue } from "../../types/device.ts";
import { Typography } from "@mui/material";

type DeviceFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function ReferenceSection() {
    return (
        <Box sx={{ width: "100%" }}>
        </Box>
    )
}

const showFileConfig = [DeviceType.A15] as DeviceTypeValue[];

export function DeviceForm(props: DeviceFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
            toolbar={props.readonly === true ? false : undefined}
            warnWhenUnsavedChanges
        >
            {props.mode !== Action.CREATE && (
                <div>
                    <TextInput source={"id"} readOnly={true} />
                    <Divider />
                </div>
            )}

            <TextInput source="name" validate={[required()]} readOnly={props.readonly} />
            <FeatureList source={"type"} types={"device-type"}>
                <AutocompleteInput source={"type"} readOnly={props.readonly} validate={[required()]} />
            </FeatureList>
            <TextInput source="ip_address" validate={[required()]} readOnly={props.readonly} />
            <NumberInput source="port" validate={[required()]} readOnly={props.readonly} />

            <FormDataConsumer<{ type: DeviceTypeValue }>>
                {({ formData, ...rest }) => showFileConfig.includes(formData.type) &&
                    <>
                        <Typography component="p" gutterBottom>File Sender Config</Typography>
                        <Divider />
                        <TextInput source="username" readOnly={props.readonly} />
                        <PasswordInput source="password" readOnly={props.readonly} />
                        <TextInput source="path" readOnly={props.readonly} />
                    </>
                }
            </FormDataConsumer>
        </SimpleForm>
    )
}

export function DeviceCreate() {
    return (
        <Create redirect={useRefererRedirect("list")} resource="device">
            <DeviceForm mode={"CREATE"} />
        </Create>
    )
}

export function DeviceShow() {
    return (
        <Show resource="device">
            <DeviceForm readonly mode={"SHOW"} />
            <ReferenceSection />
        </Show>
    )
}

export function DeviceEdit() {
    return (
        <Edit mutationMode={"pessimistic"} resource="device">
            <DeviceForm mode={"EDIT"} />
        </Edit>
    )
}

const DeviceFilterSidebar = () => (
    <SideFilter>
        <FilterLiveSearch />
    </SideFilter>
);


export const DeviceList = () => (
    <List aside={<DeviceFilterSidebar />} resource="device"
        storeKey={false} exporter={false}
        sort={{
            field: "id",
            order: "DESC"
        }}
    >
        <Datagrid>
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="type" />
            <TextField source="ip_address" />
            <TextField source="port" />
        </Datagrid>
    </List>
);
