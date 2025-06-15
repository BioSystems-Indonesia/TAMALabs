import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import {
    AutocompleteInput,
    Button,
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
    TextInput,
    WithRecord
} from "react-admin";
import FeatureList from "../../component/FeatureList.tsx";
import { Action, ActionKeys } from "../../types/props.ts";
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import SideFilter from "../../component/SideFilter.tsx";
import { Device, DeviceType, DeviceTypeValue } from "../../types/device.ts";
import { Typography } from "@mui/material";
import useAxios from "../../hooks/useAxios.ts";
import { ConnectionStatus } from './ConnectionStatus';
import { DeviceConnectionManager } from './DeviceConnectionManager';
import { useState, useEffect } from 'react';

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

export const DeviceList = () => {
    const [deviceIds, setDeviceIds] = useState<number[]>([]);
    const [connectionStatuses, setConnectionStatuses] = useState<Record<number, any>>({});

    const handleStatusUpdate = (deviceId: number, status: any) => {
        setConnectionStatuses(prev => ({
            ...prev,
            [deviceId]: status
        }));
    };

    return (
        <>
            <DeviceConnectionManager 
                deviceIds={deviceIds}
                onStatusUpdate={handleStatusUpdate}
            />
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
                    <WithRecord label="Connection Status" render={(record: Device) => {
                        useEffect(() => {
                            setDeviceIds(prev => {
                                if (!prev.includes(record.id)) {
                                    return [...prev, record.id];
                                }
                                return prev;
                            });
                        }, [record.id]);

                        return (
                            <ConnectionStatus 
                                deviceId={record.id}
                                status={connectionStatuses[record.id]}
                            />
                        );
                    }}/>
                </Datagrid>
            </List>
        </>
    );
};
