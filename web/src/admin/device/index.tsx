import { CircularProgress, Stack, Card, CardContent, Typography, Chip, Box as MuiBox } from "@mui/material";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import { useEffect, useState } from 'react';
import {
    AutocompleteInput,
    Create,
    Edit,
    FilterLiveSearch,
    FormDataConsumer,
    List,
    maxValue,
    minValue,
    PasswordInput,
    required,
    Show,
    SimpleForm,
    TextInput,
    useGetList,
    useListContext,
    Link
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

        // Common baud rates for serial communication
        const baudRates = [
            { id: 9600, name: "9600" },
            { id: 19200, name: "19200" },
            { id: 38400, name: "38400" },
            { id: 57600, name: "57600" },
            { id: 115200, name: "115200" },
            { id: 230400, name: "230400" },
            { id: 460800, name: "460800" },
            { id: 921600, name: "921600" }
        ];
      
        return (
            <>
                <AutocompleteInput
                    source="receive_port"
                    choices={props.serialPortList}
                    validate={[required()]}
                    freeSolo
                />
                <AutocompleteInput
                    source="baud_rate"
                    choices={baudRates}
                    validate={[required()]}
                    defaultValue={9600}
                    readOnly={props.readonly}
                />
                    readOnly={props.readonly}
                    disabled={props.readonly}
                <FeatureList source={"baud_rate"} types={"baud-rate"}>
                    <AutocompleteInput
                        source="baud_rate"
                        validate={[required()]}
                        defaultValue={9600}
                        readOnly={props.readonly}
                        parse={value => value === '' ? undefined : Number(value)}
                    />
                </FeatureList>
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

const DeviceCard = ({ record, connectionStatuses }: { record: Device, connectionStatuses: Record<number, ConnectionResponse> }) => {
    return (
        <Link to={`/device/${record.id}`} style={{ textDecoration: 'none' }}>
            <Card 
                elevation={0}
                sx={{ 
                    boxShadow: 'rgba(0, 0, 0, 0.16) 0px 1px 4px',
                    cursor: 'pointer',
                    '&:hover': { 
                        boxShadow: 4,
                        transform: 'translateY(-2px)',
                        transition: 'all 0.2s ease-in-out'
                    } 
                }}
            >
                <CardContent sx={{position: "relative", height: 220}}>
                    <MuiBox sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Typography variant="h5" component="div">
                            {record.name}
                        </Typography>
                        <Chip 
                            label={`ID: ${record.id}`} 
                            size="small" 
                            color="primary" 
                            variant="outlined"
                        />
                    </MuiBox>

                <MuiBox sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2, mb: 2 }}>
                    <MuiBox sx={{ flex: 1 }}>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                            <strong>Type:</strong> {record.type}
                        </Typography>
                        {record.ip_address && (
                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                <strong>IP Address:</strong> {record.ip_address}
                            </Typography>
                        )}
                        
                    </MuiBox>
                    <MuiBox sx={{ flex: 1, textAlign: 'end' }}>
                        {record.send_port !== undefined && Number(record.send_port) > 0 && (
                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                <strong>Send Port:</strong> {record.send_port}
                            </Typography>
                        )}
                        {record.receive_port && (
                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                <strong>Receive Port:</strong> {record.receive_port}
                            </Typography>
                        )}
                        {record.baud_rate !== undefined && Number(record.baud_rate) > 0 && (
                            <Typography variant="body2" color="text.secondary" gutterBottom>
                                <strong>Baud Rate:</strong> {record.baud_rate}
                            </Typography>
                        )}
                    </MuiBox>
                </MuiBox>

                <MuiBox sx={{ display: 'flex', gap: 2, mt: 2 , justifyContent: "space-between", position:"absolute", bottom: 15, width: "92%"}}>
                    <MuiBox>
                        <Typography variant="caption" display="block" gutterBottom>
                            Sender Status:
                        </Typography>
                        <ConnectionStatus
                            deviceId={record.id}
                            status={{
                                device_id: record.id,
                                message: connectionStatuses[record.id]?.sender_message,
                                status: connectionStatuses[record.id]?.sender_status
                            }}
                        />
                    </MuiBox>
                    <MuiBox>
                        <Typography variant="caption" display="block" gutterBottom>
                            Receiver Status:
                        </Typography>
                        <ConnectionStatus
                            deviceId={record.id}
                            status={{
                                device_id: record.id,
                                message: connectionStatuses[record.id]?.receiver_message,
                                status: connectionStatuses[record.id]?.receiver_status
                            }}
                        />
                    </MuiBox>
                </MuiBox>
            </CardContent>
        </Card>
        </Link>
    );
};

const DeviceCardList = ({ connectionStatuses, setDeviceIds }: { 
    connectionStatuses: Record<number, ConnectionResponse>,
    setDeviceIds: React.Dispatch<React.SetStateAction<number[]>>
}) => {
    const { data, isLoading } = useListContext<Device>();

    useEffect(() => {
        if (data) {
            const ids = data.map(device => device.id);
            setDeviceIds(ids);
        }
    }, [data, setDeviceIds]);

    if (isLoading) {
        return (
            <MuiBox sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
            </MuiBox>
        );
    }

    if (!data || data.length === 0) {
        return (
            <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                <Typography variant="body1" color="text.secondary">
                    No devices found
                </Typography>
            </MuiBox>
        );
    }

    return (
        <MuiBox sx={{ p: 2 }}>
            <MuiBox 
                sx={{ 
                    display: 'grid',
                    gridTemplateColumns: {
                        xs: '1fr',
                        md: 'repeat(2, 1fr)',
                        lg: 'repeat(3, 1fr)'
                    },
                    gap: 2
                }}
            >
                {data.map((device) => (
                    <DeviceCard 
                        key={device.id}
                        record={device} 
                        connectionStatuses={connectionStatuses}
                    />
                ))}
            </MuiBox>
        </MuiBox>
    );
};

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
            <List 
                // aside={<DeviceFilterSidebar />} 
                resource="device"
                storeKey={false} 
                exporter={false}
                sort={{
                    field: "id",
                    order: "DESC"
                }}
            >
                <DeviceCardList 
                    connectionStatuses={connectionStatuses}
                    setDeviceIds={setDeviceIds}
                />
            </List>
        </>
    );
};
