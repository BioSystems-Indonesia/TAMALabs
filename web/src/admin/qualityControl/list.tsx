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

interface DeviceStats {
    total_qc?: number;
    qc_this_month?: number;
    last_qc?: string;
    last_qc_status?: string;
    qc_today_status?: string;
    level_1_complete?: boolean;
    level_2_complete?: boolean;
    level_3_complete?: boolean;
    level_1_today?: boolean;
    level_2_today?: boolean;
    level_3_today?: boolean;
}

const DeviceCard = ({ record, connectionStatuses }: DeviceCardProps) => {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';
    const [stats, setStats] = useState<DeviceStats | null>(null);
    const [loadingStats, setLoadingStats] = useState(false);

    const formatDate = (iso?: string) => {
        if (!iso) return '-';
        const d = new Date(iso);
        if (isNaN(d.getTime())) return iso;
        return d.toLocaleString();
    };

    const mapStatusColor = (s?: string) => {
        if (!s) return '#9e9e9e';
        const v = s.toLowerCase();
        if (v.includes('in control') || v.includes('normal') || v.includes('good')) return '#4caf50';
        if (v.includes('warn') || v.includes('partial')) return '#ff9800';
        return '#f44336';
    };

    const l1Today = stats?.level_1_today;
    const l2Today = stats?.level_2_today;
    const l3Today = stats?.level_3_today;

    const l1Done = l1Today === true;
    const l2Done = l2Today === true;
    const l3Done = l3Today === true;

    // Fetch device QC statistics
    useEffect(() => {
        let mounted = true;
        const load = async () => {
            setLoadingStats(true);
            try {
                const res = await fetch(`/api/v1/quality-control/statistics?device_id=${record.id}`, { credentials: 'include' });
                if (!res.ok) {
                    if (mounted) setStats(null);
                    return;
                }
                const data = await res.json();
                if (mounted) setStats(data);
            } catch (e) {
                if (mounted) setStats(null);
            } finally {
                if (mounted) setLoadingStats(false);
            }
        };

        load();
        return () => { mounted = false; };
    }, [record.id]);

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
                            backgroundColor: isDarkMode ? 'rgba(76, 175, 80, 0.06)' : 'rgba(76, 175, 80, 0.03)',
                            border: `1px solid ${isDarkMode ? 'rgba(76, 175, 80, 0.18)' : 'rgba(76, 175, 80, 0.12)'}`
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                <strong>QC Today:</strong>
                            </Typography>
                            {loadingStats ? (
                                <CircularProgress size={18} />
                            ) : (
                                (() => {
                                    const status = stats?.qc_today_status ?? 'Not Done';
                                    const lower = status.toLowerCase();
                                    const bg = lower === 'done' ? '#4caf50' : lower.includes('partial') ? '#ff9800' : '#f44336';
                                    const label = status;
                                    return (
                                        <Chip
                                            label={label}
                                            size="small"
                                            sx={{
                                                backgroundColor: bg,
                                                color: 'white',
                                                fontWeight: 600,
                                                fontSize: '0.75rem'
                                            }}
                                        />
                                    );
                                })()
                            )}
                        </MuiBox>

                        <MuiBox sx={{ display: 'flex', gap: 2 }}>
                            <MuiBox sx={{ flex: 1 }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Total QC:</strong>
                                    {' '}
                                    {loadingStats ? <CircularProgress size={14} /> : `${stats?.total_qc ?? 0} times`}
                                </Typography>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>QC This Month:</strong>
                                    {' '}
                                    {loadingStats ? <CircularProgress size={14} /> : stats?.qc_this_month ?? 0}
                                </Typography>
                            </MuiBox>
                            <MuiBox sx={{ flex: 1, textAlign: 'end' }}>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Last QC:</strong>
                                    {' '}
                                    {loadingStats ? <CircularProgress size={14} /> : formatDate(stats?.last_qc)}
                                </Typography>
                                <Typography variant="body2" color="text.secondary" gutterBottom>
                                    <strong>Status:</strong>
                                    {' '}
                                    {loadingStats ? (
                                        <CircularProgress size={14} />
                                    ) : (
                                        <span style={{ color: mapStatusColor(stats?.last_qc_status), fontWeight: 600 }}>{stats?.last_qc_status ? stats.last_qc_status : 'N/A'}</span>
                                    )}
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
                            label={`Level 1: ${l1Done ? '✓' : '✗'}`}
                            size="small"
                            sx={{
                                backgroundColor: l1Done ? (isDarkMode ? 'rgba(76,175,80,0.18)' : 'rgba(76,175,80,0.12)') : (isDarkMode ? 'rgba(244,67,54,0.12)' : 'rgba(244,67,54,0.06)'),
                                color: l1Done ? '#4caf50' : '#f44336',
                                fontWeight: 500
                            }}
                        />
                        <Chip
                            label={`Level 2: ${l2Done ? '✓' : '✗'}`}
                            size="small"
                            sx={{
                                backgroundColor: l2Done ? (isDarkMode ? 'rgba(76,175,80,0.18)' : 'rgba(76,175,80,0.12)') : (isDarkMode ? 'rgba(244,67,54,0.12)' : 'rgba(244,67,54,0.06)'),
                                color: l2Done ? '#4caf50' : '#f44336',
                                fontWeight: 500
                            }}
                        />
                        <Chip
                            label={`Level 3: ${l3Done ? '✓' : '✗'}`}
                            size="small"
                            sx={{
                                backgroundColor: l3Done ? (isDarkMode ? 'rgba(255,152,0,0.18)' : 'rgba(255,152,0,0.12)') : (isDarkMode ? 'rgba(244,67,54,0.12)' : 'rgba(244,67,54,0.06)'),
                                color: l3Done ? '#ff9800' : '#f44336',
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