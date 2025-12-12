import { Box as MuiBox, Card, CardContent, Typography, Chip, CircularProgress } from '@mui/material';
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
    const getStatusColor = () => {
        switch (parameter.qc_status) {
            case 'Pass': return '#4caf50';
            case 'Fail': return '#f44336';
            case 'Pending': return '#ff9800';
            default: return '#9e9e9e';
        }
    };

    const getStatusBgColor = () => {
        switch (parameter.qc_status) {
            case 'Pass': return 'rgba(76, 175, 80, 0.1)';
            case 'Fail': return 'rgba(244, 67, 54, 0.1)';
            case 'Pending': return 'rgba(255, 152, 0, 0.1)';
            default: return 'rgba(158, 158, 158, 0.1)';
        }
    };

    return (
        <Link to={`/quality-control/${deviceId}/parameter/${parameter.test_type_id}`} style={{ textDecoration: 'none' }}>
            <Card
                sx={{
                    transition: 'all 0.2s ease',
                    cursor: 'pointer',
                    '&:hover': {
                        transform: 'translateY(-2px)',
                        boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                    }
                }}
            >
                <CardContent>
                    <MuiBox sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                        <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <ScienceIcon sx={{ color: '#2196f3' }} />
                            <MuiBox>
                                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                    {parameter.test_type?.code}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {parameter.test_type?.name}
                                </Typography>
                            </MuiBox>
                        </MuiBox>
                        <Chip
                            label={parameter.qc_status}
                            size="small"
                            sx={{
                                backgroundColor: getStatusBgColor(),
                                color: getStatusColor(),
                                fontWeight: 600,
                                border: `1px solid ${getStatusColor()}`
                            }}
                        />
                    </MuiBox>

                    <MuiBox sx={{ display: 'flex', gap: 3, mb: 1 }}>
                        <MuiBox>
                            <Typography variant="caption" color="text.secondary" display="block">
                                Unit
                            </Typography>
                            <Typography variant="body2" fontWeight={500}>
                                {parameter.test_type?.unit}
                            </Typography>
                        </MuiBox>
                        <MuiBox>
                            <Typography variant="caption" color="text.secondary" display="block">
                                Reference Range
                            </Typography>
                            <Typography variant="body2" fontWeight={500}>
                                {parameter.test_type?.reference_range_min} - {parameter.test_type?.reference_range_max}
                            </Typography>
                        </MuiBox>
                    </MuiBox>

                    <MuiBox sx={{ display: 'flex', gap: 3, mt: 2 }}>
                        <MuiBox>
                            <Typography variant="caption" color="text.secondary" display="block">
                                Last QC
                            </Typography>
                            <Typography variant="body2" fontWeight={500}>
                                {parameter.last_qc_date}
                            </Typography>
                        </MuiBox>
                        <MuiBox>
                            <Typography variant="caption" color="text.secondary" display="block">
                                Total Runs
                            </Typography>
                            <Typography variant="body2" fontWeight={500}>
                                {parameter.total_runs}
                            </Typography>
                        </MuiBox>
                        <MuiBox>
                            <Typography variant="caption" color="text.secondary" display="block">
                                Pass Rate
                            </Typography>
                            <Typography
                                variant="body2"
                                fontWeight={500}
                                sx={{
                                    color: parameter.pass_rate >= 95 ? '#4caf50' : parameter.pass_rate >= 85 ? '#ff9800' : '#f44336'
                                }}
                            >
                                {parameter.pass_rate}%
                            </Typography>
                        </MuiBox>
                    </MuiBox>
                </CardContent>
            </Card>
        </Link>
    );
};

export const QualityControlDetail = () => {
    const { id } = useParams<{ id: string }>();
    const deviceId = parseInt(id || '0');

    // Fetch device data
    const { data: device, isLoading: deviceLoading } = useGetOne<Device>('device', { id: deviceId });

    // Fetch test types yang berelasi dengan device ini
    const { data: testTypes, isLoading: testTypesLoading } = useGetList<TestType>('test-type', {
        filter: { device_id: deviceId },
        pagination: { page: 1, perPage: 1000 },
        sort: { field: 'code', order: 'ASC' }
    });

    // Generate QC parameters from test types
    // TODO: Replace with actual QC data from backend when API is ready
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

    const passCount = qcParameters.filter((p: QCParameter) => p.qc_status === 'Pass').length;
    const failCount = qcParameters.filter((p: QCParameter) => p.qc_status === 'Fail').length;
    const pendingCount = qcParameters.filter((p: QCParameter) => p.qc_status === 'Pending').length;

    return (
        <MuiBox sx={{ mt: 2 }}>
            <MuiBox sx={{ mb: 3 }}>


                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
                            {device.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                            Device ID: {device.id} | Type: {device.type}
                        </Typography>

                        <MuiBox sx={{ display: 'flex', gap: 2, mt: 3 }}>
                            <Chip
                                label={`${passCount} Pass`}
                                sx={{
                                    backgroundColor: 'rgba(76, 175, 80, 0.1)',
                                    color: '#4caf50',
                                    fontWeight: 600,
                                }}
                            />
                            <Chip
                                label={`${failCount} Fail`}
                                sx={{
                                    backgroundColor: 'rgba(244, 67, 54, 0.1)',
                                    color: '#f44336',
                                    fontWeight: 600,
                                }}
                            />
                            <Chip
                                label={`${pendingCount} Pending`}
                                sx={{
                                    backgroundColor: 'rgba(255, 152, 0, 0.1)',
                                    color: '#ff9800',
                                    fontWeight: 600,
                                }}
                            />
                        </MuiBox>
                    </CardContent>
                </Card>
            </MuiBox>

            <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                QC Parameters ({qcParameters.length})
            </Typography>

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
                {qcParameters.map((parameter: QCParameter) => (
                    <QCParameterCard key={parameter.id} parameter={parameter} deviceId={deviceId} />
                ))}
            </MuiBox>

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
