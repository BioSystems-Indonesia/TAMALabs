import React, { useState, useEffect } from 'react';
import {
    Card,
    CardContent,
    Typography,
    Button,
    Box,
    Alert,
    Divider,
    Paper
} from '@mui/material';
import {
    Refresh as RefreshIcon,
    Info as InfoIcon,
    Computer as ComputerIcon,
    Schedule as ScheduleIcon
} from '@mui/icons-material';

interface LicenseInfo {
    valid: boolean;
    machine_id?: string;
    expires_at?: string;
    license_type?: string;
    company?: string;
    issued_at?: string;
    error?: string;
}

const LicenseStatusPage: React.FC = () => {
    const [licenseInfo, setLicenseInfo] = useState<LicenseInfo | null>(null);
    const [loading, setLoading] = useState(false);
    const [lastChecked, setLastChecked] = useState<Date | null>(null);

    const checkLicenseStatus = async () => {
        setLoading(true);
        try {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/license/check`, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
                credentials: 'include'
            });

            if (!response.ok) {
                throw new Error('Failed to check license');
            }

            const data = await response.json();
            console.log(data)
            setLicenseInfo(data);
            setLastChecked(new Date());
        } catch (error) {
            console.error('Error checking license:', error);
            setLicenseInfo({
                valid: false,
                error: 'Failed to check license status'
            });
            setLastChecked(new Date());
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        checkLicenseStatus();
    }, []);

    const formatDate = (dateString?: string) => {
        if (!dateString) return 'N/A';
        try {
            return new Date(dateString).toLocaleString();
        } catch {
            return dateString;
        }
    };


    return (
        <Box sx={{ margin: '0 auto' }}>
            {/* Status Overview Card */}
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Typography variant="h6">Current License Status</Typography>
                        <Button
                            variant="outlined"
                            startIcon={<RefreshIcon />}
                            onClick={checkLicenseStatus}
                            disabled={loading}
                            size="small"
                        >
                            Refresh
                        </Button>
                    </Box>
                    {licenseInfo?.valid && (
                        <>
                            <Divider sx={{ mb: 1 }} />
                            <Alert severity="success">
                                <Typography variant="body2">
                                    Your license is valid and active.
                                </Typography>
                            </Alert>
                        </>
                    )}

                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                        {lastChecked && (
                            <Typography variant="body2" color="text.secondary">
                                Last checked: {lastChecked.toLocaleString()}
                            </Typography>
                        )}
                    </Box>

                    {licenseInfo?.error && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {licenseInfo.error}
                        </Alert>
                    )}

                    {!licenseInfo?.valid && !loading && (
                        <Alert severity="warning" sx={{ mb: 2 }}>
                            Your license is not valid. Please contact your administrator or activate a new license.
                        </Alert>
                    )}
                </CardContent>
            </Card>

            {/* License Details Card */}
            {licenseInfo && (
                <Card>
                    <CardContent>
                        <Typography variant="h6" gutterBottom>
                            License Details
                        </Typography>

                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
                                <Paper sx={{ p: 2, bgcolor: 'grey.50', flex: '1 1 300px' }}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                        <ComputerIcon sx={{ mr: 1, color: 'primary.main' }} />
                                        <Typography variant="subtitle2">Machine Information</Typography>
                                    </Box>
                                    <Typography variant="body2" color="text.secondary">
                                        Machine ID: {licenseInfo.machine_id || 'N/A'}
                                    </Typography>
                                </Paper>

                                {licenseInfo.valid && (
                                    <Paper sx={{ p: 2, bgcolor: 'grey.50', flex: '1 1 300px' }}>
                                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                            <InfoIcon sx={{ mr: 1, color: 'primary.main' }} />
                                            <Typography variant="subtitle2">License Type</Typography>
                                        </Box>
                                        <Typography variant="body2" color="text.secondary">
                                            {licenseInfo.license_type || 'Standard'}
                                        </Typography>
                                    </Paper>
                                )}
                            </Box>

                            {licenseInfo.valid && (
                                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
                                    {licenseInfo.company && (
                                        <Paper sx={{ p: 2, bgcolor: 'grey.50', flex: '1 1 300px' }}>
                                            <Typography variant="subtitle2" gutterBottom>
                                                Licensed To
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {licenseInfo.company}
                                            </Typography>
                                        </Paper>
                                    )}

                                    <Paper sx={{ p: 2, bgcolor: 'grey.50', flex: '1 1 300px' }}>
                                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                            <ScheduleIcon sx={{ mr: 1, color: 'primary.main' }} />
                                            <Typography variant="subtitle2">Issue Date</Typography>
                                        </Box>
                                        <Typography variant="body2" color="text.secondary">
                                            {formatDate(licenseInfo.issued_at)}
                                        </Typography>
                                    </Paper>

                                    <Paper sx={{ p: 2, bgcolor: 'grey.50', flex: '1 1 300px' }}>
                                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                            <ScheduleIcon sx={{ mr: 1, color: 'warning.main' }} />
                                            <Typography variant="subtitle2">Expiration Date</Typography>
                                        </Box>
                                        <Typography variant="body2" color="text.secondary">
                                            {formatDate(licenseInfo.expires_at)}
                                        </Typography>
                                    </Paper>
                                </Box>
                            )}
                        </Box>


                    </CardContent>
                </Card>
            )}

            {/* Actions Card for Invalid License */}
            {licenseInfo && !licenseInfo.valid && (
                <Card sx={{ mt: 3 }}>
                    <CardContent>
                        <Typography variant="h6" gutterBottom>
                            Need Help?
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                            If you need to activate a new license or resolve license issues, please contact your system administrator.
                        </Typography>
                        <Box sx={{ display: 'flex', gap: 2 }}>
                            <Button
                                variant="contained"
                                onClick={() => window.location.href = '/#/license'}
                                size="small"
                            >
                                Go to License Activation
                            </Button>
                            <Button
                                variant="outlined"
                                onClick={checkLicenseStatus}
                                disabled={loading}
                                size="small"
                            >
                                Check Again
                            </Button>
                        </Box>
                    </CardContent>
                </Card>
            )}
        </Box>
    );
};

export default LicenseStatusPage;