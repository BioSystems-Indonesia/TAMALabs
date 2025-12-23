import React from 'react';
import { Box as MuiBox, Card, CardContent, Typography, CircularProgress, List, ListItem, ListItemAvatar, Avatar, ListItemText, Divider } from '@mui/material';
import { useParams } from 'react-router-dom';
import { useGetOne, useGetList, Link } from 'react-admin';
import { Device } from '../../types/device';
import ScienceIcon from '@mui/icons-material/Science';

interface TestType {
    id: number;
    code: string;
    name: string;
    unit: string;
    reference_range_min?: number;
    reference_range_max?: number;
    device_id?: number;
}

interface QCParameter {
    id: number;
    device_id: number;
    test_type_id: number;
    test_type?: TestType;
    qc_status: 'Pass' | 'Fail' | 'Pending';
    last_qc_date: string;
    total_runs: number;
    pass_rate: number;
}

const QCParameterCard = ({ parameter, deviceId }: { parameter: QCParameter; deviceId: number }) => {
    return (
        <Card>
            <MuiBox>
                <Link to={`/quality-control/${deviceId}/parameter/${parameter.test_type_id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                    <ListItem
                        sx={{
                            px: 2,
                            py: 1.25,
                            borderRadius: 1,
                            transition: 'all 0.15s ease',
                            '&:hover': { backgroundColor: 'action.hover', transform: 'translateY(-2px)' }
                        }}
                        disableGutters
                    >
                        <ListItemAvatar sx={{ minWidth: 44 }}>
                            <Avatar sx={{ bgcolor: 'transparent' }}>
                                <ScienceIcon sx={{ color: '#2196f3' }} />
                            </Avatar>
                        </ListItemAvatar>

                        <ListItemText
                            primary={<Typography variant="subtitle1" sx={{ fontWeight: 600 }}>{parameter.test_type?.code}</Typography>}
                            secondary={<Typography variant="body2" color="text.secondary">{parameter.test_type?.name}</Typography>}
                        />

                    </ListItem>
                </Link>
            </MuiBox>

        </Card>
    );
};

export const QualityControlDetail = () => {
    const { id } = useParams<{ id: string }>();
    const deviceId = parseInt(id || '0');

    // Fetch device data
    const { data: device, isLoading: deviceLoading } = useGetOne<Device>('device', { id: deviceId });

    const { data: testTypes, isLoading: testTypesLoading } = useGetList<TestType>('test-type', {
        filter: { device_id: deviceId },
        pagination: { page: 1, perPage: 1000 },
        sort: { field: 'code', order: 'ASC' }
    });

    const qcParameters: QCParameter[] = (testTypes || []).map((testType, index) => {
        const statuses: Array<'Pass' | 'Fail' | 'Pending'> = ['Pass', 'Pass', 'Pass', 'Fail', 'Pending'];
        return {
            id: testType.id,
            device_id: deviceId,
            test_type_id: testType.id,
            test_type: testType,
            qc_status: statuses[index % statuses.length],
            last_qc_date: index % 3 === 0 ? 'Today' : index % 2 === 0 ? 'Yesterday' : '2 days ago',
            total_runs: Math.floor(Math.random() * 100) + 50,
            pass_rate: Math.floor(Math.random() * 20) + 80,
        };
    });

    const isLoading = deviceLoading || testTypesLoading;

    if (isLoading) {
        return (
            <MuiBox sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
            </MuiBox>
        );
    }

    if (!device) {
        return (
            <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                <Typography variant="h6" color="text.secondary">
                    Device not found
                </Typography>
            </MuiBox>
        );
    }



    return (
        <MuiBox sx={{ mt: 2 }}>
            <MuiBox sx={{ mb: 3 }}>
                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
                            {device.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" >
                            Device ID: {device.id} | Type: {device.type}
                        </Typography>
                    </CardContent>
                </Card>
            </MuiBox>

            <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                QC Parameters ({qcParameters.length})
            </Typography>

            <List sx={{ width: '100%', bgcolor: 'transparent' }}>
                {qcParameters.map((parameter: QCParameter, idx: number) => (
                    <React.Fragment key={parameter.id}>
                        <QCParameterCard parameter={parameter} deviceId={deviceId} />
                        {idx < qcParameters.length - 1 && <Divider component="li" />}
                    </React.Fragment>
                ))}
            </List>

            {qcParameters.length === 0 && (
                <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                    <Typography variant="body1" color="text.secondary">
                        No QC parameters found for this device
                    </Typography>
                </MuiBox>
            )}
        </MuiBox>
    );
};
