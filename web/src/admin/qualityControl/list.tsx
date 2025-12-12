import { Box as MuiBox, Card, CardContent, Typography, Chip, CircularProgress, useTheme } from '@mui/material';
import { List, useListContext, Link } from 'react-admin';
import { useEffect, useState } from 'react';
import LanIcon from '@mui/icons-material/Lan';
import { Device } from '../../types/device';
import { DeviceConnectionManager, ConnectionResponse } from '../device/DeviceConnectionManager';

interface DeviceCardProps {
    record: Device;
    connectionStatuses: Record<number, ConnectionResponse>;
}

const DeviceCard = ({ record, connectionStatuses }: DeviceCardProps) => {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';

    return (
        <Link to={`/quality-control/${record.id}`} style={{ textDecoration: 'none' }}>
            <Card
                sx={{
                    position: 'relative',
                    transition: 'all 0.3s ease',
                    cursor: 'pointer',
                    minHeight: 230,
                    backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
                    border: `1px solid ${isDarkMode ? theme.palette.divider : '#e5e7eb'}`,
                    '&:hover': {
                        transform: 'translateY(-4px)',
                        boxShadow: isDarkMode
                            ? '0 10px 25px rgba(0,0,0,0.5)'
                            : '0 10px 25px rgba(0,0,0,0.1)',
                        borderColor: theme.palette.primary.main,
                    }
                }}
            >
                <CardContent>
                    <MuiBox sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 1 }}>
                        <LanIcon sx={{ color: theme.palette.primary.main }} />
                        <MuiBox sx={{ flex: 1 }}>
                            <Typography variant="h6" component="div" sx={{ fontWeight: 400 }}>
                                {record.name}
                            </Typography>
                        </MuiBox>
                        <Chip
                            label={`ID: ${record.id}`}
                            size="small"
                            color="primary"
                            variant="outlined"
                        />
                    </MuiBox>

                    <MuiBox sx={{ display: 'flex', flexDirection: 'column', gap: 1.5, mb: 2 }}>
                        <MuiBox sx={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            p: 1.5,
                            borderRadius: 1,
                            backgroundColor: isDarkMode ? 'rgba(76, 175, 80, 0.1)' : 'rgba(76, 175, 80, 0.05)',
                            border: `1px solid ${isDarkMode ? 'rgba(76, 175, 80, 0.3)' : 'rgba(76, 175, 80, 0.2)'}`
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                <strong>QC Today:</strong>
                            </Typography>
                            <Chip
                                label="Done"
                                size="small"
                                sx={{
                                    backgroundColor: '#4caf50',
                                    color: 'white',
                                    fontWeight: 600,
                                    fontSize: '0.75rem'
                                }}
                            />
                        </MuiBox>

                        <MuiBox sx={{ display: 'flex', gap: 2 }}>
                            <MuiBox sx={{ flex: 1 }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Total QC:</strong> 45 times
                                </Typography>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>QC This Month:</strong> 12 times
                                </Typography>
                            </MuiBox>
                            <MuiBox sx={{ flex: 1, textAlign: 'end' }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Last QC:</strong> Today
                                </Typography>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Status:</strong> <span style={{ color: '#4caf50', fontWeight: 600 }}>Normal</span>
                                </Typography>
                            </MuiBox>
                        </MuiBox>
                    </MuiBox>

                    <MuiBox sx={{
                        display: 'flex',
                        gap: 1,
                        mt: 2,
                        position: "absolute",
                        bottom: 15,
                        width: "92%",
                        flexWrap: 'wrap'
                    }}>
                        <Chip
                            label={`Level 1: ✓`}
                            size="small"
                            sx={{
                                backgroundColor: isDarkMode ? 'rgba(33, 150, 243, 0.2)' : 'rgba(33, 150, 243, 0.1)',
                                color: '#2196f3',
                                fontWeight: 500
                            }}
                        />
                        <Chip
                            label={`Level 2: ✓`}
                            size="small"
                            sx={{
                                backgroundColor: isDarkMode ? 'rgba(76, 175, 80, 0.2)' : 'rgba(76, 175, 80, 0.1)',
                                color: '#4caf50',
                                fontWeight: 500
                            }}
                        />
                        <Chip
                            label={`Level 3: ✓`}
                            size="small"
                            sx={{
                                backgroundColor: isDarkMode ? 'rgba(255, 152, 0, 0.2)' : 'rgba(255, 152, 0, 0.1)',
                                color: '#ff9800',
                                fontWeight: 500
                            }}
                        />
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

export const QualityControlList = () => {
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