import { CircularProgress, Stack, Card, CardContent, Typography, Chip, Box as MuiBox, useTheme } from "@mui/material";
import Box from "@mui/material/Box";
import LanIcon from '@mui/icons-material/Lan';

import { useEffect, useState } from 'react';
import {
    AutocompleteInput,
    Create,
    Edit,
    // FilterLiveSearch,
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
// import SideFilter from "../../component/SideFilter.tsx";
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
    const theme = useTheme();
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
        <Box sx={{ p: { xs: 2, sm: 3 } }}>
            <SimpleForm
                disabled={props.readonly}
                toolbar={props.readonly === true ? false : undefined}
                warnWhenUnsavedChanges
                sx={{
                    '& .RaSimpleForm-form': {
                        backgroundColor: 'transparent',
                        boxShadow: 'none',
                        padding: 0
                    }
                }}
            >
                <Stack spacing={3} sx={{ width: '100%' }}>
                    {props.mode !== Action.CREATE && (
                        <Card
                            elevation={0}
                            sx={{
                                border: `1px solid ${theme.palette.divider}`,
                                borderRadius: 2
                            }}
                        >
                            <CardContent sx={{ p: 3 }}>
                                <Typography
                                    variant="subtitle1"
                                    sx={{
                                        fontWeight: 600,
                                        color: theme.palette.text.primary,
                                        mb: 3
                                    }}
                                >
                                    ℹ️ System Information
                                </Typography>
                                <TextInput
                                    source={"id"}
                                    readOnly={true}
                                    sx={{
                                        '& .MuiOutlinedInput-root': {
                                            borderRadius: 2,
                                            transition: 'all 0.2s ease'
                                        }
                                    }}
                                />
                            </CardContent>
                        </Card>
                    )}

                    <Card
                        elevation={0}
                        sx={{
                            border: `1px solid ${theme.palette.divider}`,
                            borderRadius: 2
                        }}
                    >
                        <CardContent sx={{ p: 3 }}>
                            <Box sx={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: 1.5,
                                mb: 3
                            }}>
                                <Typography
                                    variant="subtitle1"
                                    sx={{
                                        fontWeight: 600,
                                        color: theme.palette.text.primary
                                    }}
                                >
                                    ❗Basic Information
                                </Typography>
                                <Chip
                                    label="Required"
                                    size="small"
                                    color="error"
                                    variant="outlined"
                                    sx={{ ml: 'auto', fontSize: '0.75rem' }}
                                />
                            </Box>

                            <Stack>
                                <TextInput
                                    source="name"
                                    validate={[required()]}
                                    readOnly={props.readonly}
                                    sx={{
                                        '& .MuiOutlinedInput-root': {
                                            borderRadius: 2,
                                            transition: 'all 0.2s ease',
                                            ...(!props.readonly && {
                                                '&:hover': {
                                                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                }
                                            })
                                        }
                                    }}
                                />
                                <FeatureList source={"type"} types={"device-type"}>
                                    <AutocompleteInput
                                        source={"type"}
                                        readOnly={props.readonly}
                                        validate={[required()]}
                                        sx={{
                                            '& .MuiOutlinedInput-root': {
                                                borderRadius: 2,
                                                transition: 'all 0.2s ease',
                                                ...(!props.readonly && {
                                                    '&:hover': {
                                                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                    }
                                                })
                                            }
                                        }}
                                    />
                                </FeatureList>
                            </Stack>
                        </CardContent>
                    </Card>

                    <FormDataConsumer<{ type: DeviceTypeValue }>>
                        {({ formData }) => {
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

                            if (deviceTypeFeature.additional_info.can_send || deviceTypeFeature.additional_info.can_receive) {
                                dynamicForm.push(<NetworkConfig key="network" {...props}
                                    canSend={deviceTypeFeature.additional_info.can_send}
                                    canReceive={deviceTypeFeature.additional_info.can_receive}
                                    useSerial={deviceTypeFeature.additional_info.use_serial}
                                    isLoadingSerialPortList={isLoadingSerialPortList}
                                    serialPortList={serialPortList} />)
                            }

                            if (deviceTypeFeature.additional_info.have_authentication) {
                                dynamicForm.push(<AuthenticationConfig key="auth" {...props} />)
                            }

                            if (deviceTypeFeature.additional_info.have_path) {
                                dynamicForm.push(<PathConfig key="path" {...props} />)
                            }

                            return (
                                <>
                                    {dynamicForm}
                                </>
                            )
                        }}
                    </FormDataConsumer>
                </Stack>
            </SimpleForm>
        </Box>
    )
}

type NetworkConfigProps = DeviceFormProps & {
    canSend: boolean
    canReceive: boolean
    useSerial: boolean
    isLoadingSerialPortList: boolean
    serialPortList: string[] | undefined
}

function NetworkConfig(props: NetworkConfigProps) {
    const theme = useTheme();

    if (props.canReceive && props.useSerial) {
        if (props.isLoadingSerialPortList || !props.serialPortList) {
            return (
                <Card
                    elevation={0}
                    sx={{
                        border: `1px solid ${theme.palette.divider}`,
                        borderRadius: 2
                    }}
                >
                    <CardContent sx={{ p: 3, display: 'flex', justifyContent: 'center' }}>
                        <CircularProgress />
                    </CardContent>
                </Card>
            )
        }
    }

    return (
        <Card
            elevation={0}
            sx={{
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2,
            }}
        >
            <CardContent sx={{ p: 3 }}>
                <Box sx={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    mb: 3
                }}>
                    <Typography
                        variant="subtitle1"
                        sx={{
                            fontWeight: 600,
                            color: theme.palette.text.primary
                        }}
                    >
                        🌐 Network Configuration
                    </Typography>

                    <Chip
                        label="Required"
                        size="small"
                        color="error"
                        variant="outlined"
                        sx={{ ml: 'auto', fontSize: '0.75rem' }}
                    />
                </Box>

                <Stack>
                    {props.canSend && (
                        <>
                            <TextInput
                                source="ip_address"
                                validate={[required()]}
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                            <TextInput
                                source="send_port"
                                validate={[required()]}
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                        </>
                    )}

                    {props.canReceive && (
                        <>
                            {props.useSerial ? (
                                <>
                                    <AutocompleteInput
                                        source="receive_port"
                                        choices={props.serialPortList}
                                        validate={[required()]}
                                        freeSolo
                                        readOnly={props.readonly}
                                        disabled={props.readonly}
                                        sx={{
                                            '& .MuiOutlinedInput-root': {
                                                borderRadius: 2,
                                                transition: 'all 0.2s ease',
                                                ...(!props.readonly && {
                                                    '&:hover': {
                                                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                    }
                                                })
                                            }
                                        }}
                                    />
                                    <FeatureList source={"baud_rate"} types={"baud-rate"}>
                                        <AutocompleteInput
                                            source="baud_rate"
                                            validate={[required()]}
                                            defaultValue={9600}
                                            readOnly={props.readonly}
                                            parse={value => value === '' ? undefined : Number(value)}
                                            sx={{
                                                '& .MuiOutlinedInput-root': {
                                                    borderRadius: 2,
                                                    transition: 'all 0.2s ease',
                                                    ...(!props.readonly && {
                                                        '&:hover': {
                                                            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                        }
                                                    })
                                                }
                                            }}
                                        />
                                    </FeatureList>
                                </>
                            ) : (
                                props.mode !== "CREATE" && (
                                    <TextInput
                                        source="receive_port"
                                        validate={[required(), minValue(0), maxValue(65535)]}
                                        readOnly={props.readonly}
                                        sx={{
                                            '& .MuiOutlinedInput-root': {
                                                borderRadius: 2,
                                                transition: 'all 0.2s ease',
                                                ...(!props.readonly && {
                                                    '&:hover': {
                                                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                    }
                                                })
                                            }
                                        }}
                                    />
                                )
                            )}
                        </>
                    )}
                </Stack>
            </CardContent>
        </Card>
    )
}

function AuthenticationConfig(props: DeviceFormProps) {
    const theme = useTheme();

    return (
        <Card
            elevation={0}
            sx={{
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2
            }}
        >
            <CardContent sx={{ p: 3 }}>
                <Box sx={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    mb: 3
                }}>
                    <Typography
                        variant="subtitle1"
                        sx={{
                            fontWeight: 600,
                            color: theme.palette.text.primary
                        }}
                    >
                        🔐 Authentication
                    </Typography>
                    <Chip
                        label="Optional"
                        size="small"
                        color="default"
                        variant="outlined"
                        sx={{ ml: 'auto', fontSize: '0.75rem' }}
                    />
                </Box>

                <Stack>
                    <TextInput
                        source="username"
                        readOnly={props.readonly}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                borderRadius: 2,
                                transition: 'all 0.2s ease',
                                ...(!props.readonly && {
                                    '&:hover': {
                                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                    }
                                })
                            }
                        }}
                    />
                    <PasswordInput
                        source="password"
                        readOnly={props.readonly}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                borderRadius: 2,
                                transition: 'all 0.2s ease',
                                ...(!props.readonly && {
                                    '&:hover': {
                                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                    }
                                })
                            }
                        }}
                    />
                </Stack>
            </CardContent>
        </Card>
    )
}

function PathConfig(props: DeviceFormProps) {
    const theme = useTheme();

    return (
        <Card
            elevation={0}
            sx={{
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2
            }}
        >
            <CardContent sx={{ p: 3 }}>
                <Box sx={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    mb: 3
                }}>
                    <Typography
                        variant="subtitle1"
                        sx={{
                            fontWeight: 600,
                            color: theme.palette.text.primary
                        }}
                    >
                        📁 Path Configuration
                    </Typography>
                    <Chip
                        label="Optional"
                        size="small"
                        color="default"
                        variant="outlined"
                        sx={{ ml: 'auto', fontSize: '0.75rem' }}
                    />
                </Box>
                <TextInput
                    source="path"
                    readOnly={props.readonly}
                    sx={{
                        '& .MuiOutlinedInput-root': {
                            borderRadius: 2,
                            transition: 'all 0.2s ease',
                            ...(!props.readonly && {
                                '&:hover': {
                                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                }
                            })
                        }
                    }}
                />
            </CardContent>
        </Card>
    )
}

export function DeviceCreate() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Create redirect={useRefererRedirect("list")} resource="device">
                <DeviceForm mode={"CREATE"} />
            </Create>
        </Box>
    )
}

export function DeviceShow() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Show resource="device">
                <DeviceForm readonly mode={"SHOW"} />
                <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
                    <ReferenceSection />
                </Box>
            </Show>
        </Box>
    )
}

export function DeviceEdit() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Edit mutationMode={"pessimistic"} resource="device">
                <DeviceForm mode={"EDIT"} />
            </Edit>
        </Box>
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
                <CardContent sx={{ position: "relative", height: 220 }}>
                    <MuiBox sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <LanIcon color="primary" />
                            <Typography variant="h6" component="div">
                                {record.name}
                            </Typography>
                        </Box>
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

                    <MuiBox sx={{ display: 'flex', gap: 2, mt: 2, justifyContent: "space-between", position: "absolute", bottom: 15, width: "92%" }}>
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

// const DeviceFilterSidebar = () => {
//     const theme = useTheme();
//     const isDarkMode = theme.palette.mode === 'dark';
    
//     return (
//         <SideFilter sx={{
//             backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',          
//         }}>
//             <FilterLiveForm debounce={1500}>
//                 <Stack spacing={0}>
//                     <Box>
//                         <Typography variant="h6" sx={{ 
//                             color: theme.palette.text.primary, 
//                             marginBottom: 2, 
//                             fontWeight: 600,
//                             fontSize: '1.1rem',
//                             textAlign: 'center'
//                         }}>
//                             🖥️ Filter Devices
//                         </Typography>
//                     </Box>
//                     <SearchInput 
//                         source="q" 
//                         alwaysOn 
//                         sx={{
//                             '& .MuiOutlinedInput-root': {
//                                 backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
//                                 borderRadius: '12px',
//                                 transition: 'all 0.3s ease',
//                                 border: isDarkMode ? `1px solid ${theme.palette.divider}` : '1px solid #e5e7eb',
//                                 '&:hover': {
//                                     backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
//                                 },
//                                 '&.Mui-focused': {
//                                     backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
//                                 }
//                             },
//                             '& .MuiInputLabel-root': {
//                                 color: theme.palette.text.secondary,
//                                 fontWeight: 500,
//                             }
//                         }} 
//                     />
//                 </Stack>
//             </FilterLiveForm>
//         </SideFilter>
//     )
// };

// const DeviceFilterSidebar = () => (
//     <SideFilter>
//         <FilterLiveSearch />
//     </SideFilter>
// );

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
