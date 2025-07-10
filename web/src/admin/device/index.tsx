import { CircularProgress, Stack } from "@mui/material";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import { useEffect, useState } from 'react';
import {
    AutocompleteInput,
    Create,
    Datagrid,
    Edit,
    FilterLiveSearch,
    FormDataConsumer,
    List,
    maxValue,
    minValue,
    NumberInput,
    PasswordInput,
    required,

    Show,
    SimpleForm,
    TextField,
    TextInput,
    useGetList,
    WithRecord
} from "react-admin";
import FeatureList from "../../component/FeatureList.tsx";
import SideFilter from "../../component/SideFilter.tsx";
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import { Device, DeviceTypeFeatureList, DeviceTypeValue } from "../../types/device.ts";
import { Action, ActionKeys } from "../../types/props.ts";
import { ConnectionStatus } from './ConnectionStatus';
import { ConnectionResponse, DeviceConnectionManager } from './DeviceConnectionManager';

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


export function DeviceForm(props: DeviceFormProps) {
    const { data: deviceTypeFeatureList, isLoading: isLoadingDeviceTypeFeatureList } = useGetList<DeviceTypeFeatureList>("feature-list-device-type", {
        pagination: {
            page: 1,
            perPage: 1000
        }
    });

    const { data: serialPortList, isLoading: isLoadingSerialPortList } = useGetList("server/serial-port-list", {
        pagination: {
            page: 1,
            perPage: 1000
        }
    });

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


            <FormDataConsumer<{ type: DeviceTypeValue }>>
                {({ formData, ...rest }) => {
                    const dynamicForm = []
                    if (isLoadingDeviceTypeFeatureList || !deviceTypeFeatureList) {
                        return <Stack>
                            <CircularProgress />
                        </Stack>
                    }

                    const deviceTypeFeature = deviceTypeFeatureList?.find(item => item.id === formData.type);
                    if (!deviceTypeFeature) {
                        console.error(`Device type feature not found for ${formData.type}`)
                        return null
                    }

                    if (deviceTypeFeature.additional_info.can_send) {
                        dynamicForm.push(<SendConfig {...props} />)
                    }

                    if (deviceTypeFeature.additional_info.can_receive) {
                        dynamicForm.push(<ReceiveConfig {...props}
                            useSerial={deviceTypeFeature.additional_info.use_serial}
                            isLoadingSerialPortList={isLoadingSerialPortList}
                            serialPortList={serialPortList} />)
                    }

                    if (deviceTypeFeature.additional_info.have_authentication) {
                        dynamicForm.push(<AuthenticationConfig {...props} />)
                    }

                    if (deviceTypeFeature.additional_info.have_path) {
                        dynamicForm.push(<PathConfig {...props} />)
                    }

                    return (
                        <>
                            {dynamicForm.map((item, index) => (
                                <div key={index}>
                                    {item}
                                </div>
                            ))}
                        </>
                    )
                }}
            </FormDataConsumer>
        </SimpleForm>
    )
}

function AuthenticationConfig(props: DeviceFormProps) {
    return (
        <>
            <TextInput source="username" readOnly={props.readonly} />
            <PasswordInput source="password" readOnly={props.readonly} />
        </>
    )
}

function SendConfig(props: DeviceFormProps) {
    return (
        <>
            <TextInput source="ip_address" validate={[required()]} readOnly={props.readonly} />
            <TextInput source="send_port" validate={[required()]} readOnly={props.readonly} />
        </>
    )
}

type ReceiveConfigProps = DeviceFormProps & {
    useSerial: boolean
    isLoadingSerialPortList: boolean
    serialPortList: string[] | undefined
}

function ReceiveConfig(props: ReceiveConfigProps) {
    if (props.useSerial) {
        if (props.isLoadingSerialPortList || !props.serialPortList) {
            return <Stack>
                <CircularProgress />
            </Stack>
        }

        return (
            <>
                <AutocompleteInput
                    source="receive_port"
                    choices={props.serialPortList}
                    validate={[required()]}
                    freeSolo
                />
            </>
        )
    }


    return (
        <>
            {props.mode !== "CREATE" && <TextInput source="receive_port" validate={[required(), minValue(0), maxValue(65535)]} readOnly={props.readonly} />}
        </>
    )
}

function PathConfig(props: DeviceFormProps) {
    return (
        <>
            <TextInput source="path" readOnly={props.readonly} />
        </>
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
    const [connectionStatuses, setConnectionStatuses] = useState<Record<number, ConnectionResponse>>({});

    const handleStatusUpdate = (deviceId: number, status: ConnectionResponse) => {
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
                    <TextField source="send_port" />
                    <TextField source="receive_port" />
                    <WithRecord label="Connection Status Sender" render={(record: Device) => {
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
                                status={{
                                    device_id: record.id,
                                    message: connectionStatuses[record.id]?.sender_message,
                                    status: connectionStatuses[record.id]?.sender_status
                                }}
                            />
                        );
                    }} />
                    <WithRecord label="Connection Status Receiver" render={(record: Device) => {
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
                                status={{
                                    device_id: record.id,
                                    message: connectionStatuses[record.id]?.receiver_message,
                                    status: connectionStatuses[record.id]?.receiver_status
                                }}
                            />
                        );
                    }} />
                </Datagrid>
            </List>
        </>
    );
};
