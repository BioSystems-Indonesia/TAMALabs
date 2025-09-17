import React from 'react';
import {
    Card,
    CardContent,
    Typography,
    Box,
    Stack,
    Chip,
    Avatar,
    useTheme,
    Skeleton
} from '@mui/material';
import {
    Phone,
    LocationOn,
    Home,
    Cake
} from '@mui/icons-material';
import { useGetOne } from 'react-admin';

interface PatientInfoCardProps {
    patientId: string | number;
    compact?: boolean;
}

interface PatientData {
    id: number;
    first_name: string;
    last_name: string;
    birthdate: string;
    sex: string;
    phone_number?: string;
    location?: string;
    address?: string;
    created_at: string;
    updated_at: string;
}

const PatientInfoCard: React.FC<PatientInfoCardProps> = ({ patientId, compact = false }) => {

    const theme = useTheme();
    const { data: patient, isLoading, error } = useGetOne<PatientData>('patient', { id: Number(patientId) });

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('id-ID', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    const calculateAge = (birthdate: string) => {
        const birth = new Date(birthdate);
        const today = new Date();
        let age = today.getFullYear() - birth.getFullYear();
        const monthDiff = today.getMonth() - birth.getMonth();
        
        if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
            age--;
        }
        
        return age;
    };

    const getInitials = (firstName: string, lastName: string) => {
        return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();
    };

    const getSexIcon = (sex: string) => {
        return sex === 'M' ? '♂️' : sex === 'F' ? '♀️' : '⚧️';
    };

    const getSexColor = (sex: string) => {
        return sex === 'M' ? theme.palette.info.main : 
               sex === 'F' ? theme.palette.secondary.main : 
               theme.palette.grey[500];
    };

    if (isLoading) {
        return (
            <Card
                elevation={0}
                sx={{
                    border: `1px solid ${theme.palette.divider}`,
                    borderRadius: 2,
                    height: compact ? 200 : 'auto'
                }}
            >
                <CardContent sx={{ p: 3 }}>
                    <Stack spacing={2}>
                        <Box display="flex" alignItems="center" gap={2}>
                            <Skeleton variant="circular" width={compact ? 40 : 60} height={compact ? 40 : 60} />
                            <Box flex={1}>
                                <Skeleton variant="text" width="60%" height={compact ? 24 : 32} />
                                <Skeleton variant="text" width="40%" height={20} />
                            </Box>
                        </Box>
                        <Skeleton variant="rectangular" height={compact ? 80 : 120} />
                    </Stack>
                </CardContent>
            </Card>
        );
    }

    if (error || !patient) {
        return (
            <Card
                elevation={0}
                sx={{
                    border: `1px solid ${theme.palette.error.main}`,
                    borderRadius: 2,
                    bgcolor: theme.palette.error.light + '10'
                }}
            >
                <CardContent sx={{ p: 3, textAlign: 'center' }}>
                    <Typography color="error" variant="h6">
                        ⚠️ Error Loading Patient
                    </Typography>
                    <Typography color="error" variant="body2" sx={{ mt: 1 }}>
                        {error ? 'Failed to load patient data' : 'Patient not found'}
                    </Typography>
                </CardContent>
            </Card>
        );
    }

    const InfoItem = ({ icon, label, value }: { icon: React.ReactNode, label: string, value: string | undefined }) => {
        if (!value) return null;
        
        return (
            <Box display="flex" alignItems="center" gap={1.5}>
                <Box
                    sx={{
                        color: theme.palette.text.secondary,
                        display: 'flex',
                        alignItems: 'center',
                        minWidth: compact ? 20 : 24
                    }}
                >
                    {icon}
                </Box>
                <Box flex={1}>
                    <Typography
                        variant={compact ? "caption" : "body2"}
                        color="text.secondary"
                        sx={{ fontWeight: 500, fontSize: compact ? '0.7rem' : '0.875rem' }}
                    >
                        {label}
                    </Typography>
                    <Typography
                        variant={compact ? "body2" : "body1"}
                        sx={{ fontWeight: 500, fontSize: compact ? '0.8rem' : '1rem' }}
                    >
                        {value}
                    </Typography>
                </Box>
            </Box>
        );
    };

    return (
        <Card
            elevation={0}
            sx={{
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2,
                transition: 'all 0.2s ease',
                '&:hover': {
                    boxShadow: theme.shadows[4],
                    borderColor: theme.palette.primary.main + '40'
                }
            }}
        >
            <CardContent sx={{ p: compact ? 2 : 3 }}>
                {/* Horizontal layout: Avatar on left, info on right */}
                <Box display="flex" gap={3} alignItems="flex-start">
                    {/* Avatar section */}
                    <Avatar
                        sx={{
                            width: compact ? 50 : 80,
                            height: compact ? 50 : 80,
                            bgcolor: theme.palette.primary.main,
                            fontSize: compact ? '1.2rem' : '2rem',
                            fontWeight: 600,
                            flexShrink: 0
                        }}
                    >
                        {getInitials(patient.first_name, patient.last_name)}
                    </Avatar>

                    {/* Name and Information section */}
                    <Box flex={1} minWidth={0}>
                        {/* Name and chips */}
                        <Box mb={2}>
                            <Typography
                                variant={compact ? "h6" : "h5"}
                                sx={{
                                    fontWeight: 700,
                                    color: theme.palette.text.primary,
                                    fontSize: compact ? '1.1rem' : '1.5rem',
                                    mb: 1
                                }}
                            >
                                {patient.first_name} {patient.last_name}
                            </Typography>
                            <Box display="flex" alignItems="center" gap={1} flexWrap="wrap">
                                <Chip
                                    label={`ID: ${patient.id}`}
                                    size="small"
                                    variant="outlined"
                                    sx={{ fontSize: compact ? '0.65rem' : '0.75rem' }}
                                />
                                <Chip
                                    icon={<span style={{ fontSize: compact ? '0.8rem' : '1rem' }}>{getSexIcon(patient.sex)}</span>}
                                    label={patient.sex === 'M' ? 'Male' : patient.sex === 'F' ? 'Female' : 'Other'}
                                    size="small"
                                    sx={{
                                        bgcolor: getSexColor(patient.sex) + '20',
                                        color: getSexColor(patient.sex),
                                        fontSize: compact ? '0.65rem' : '0.75rem'
                                    }}
                                />
                            </Box>
                        </Box>

                        {/* Patient Information */}
                        <Stack spacing={compact ? 1.5 : 2}>
                            {/* Birth Date and Phone side by side */}
                            <Box display="flex" gap={3} flexWrap="wrap">
                                <Box flex={1} minWidth="200px">
                                    <InfoItem
                                        icon={<Cake fontSize={compact ? "small" : "medium"} />}
                                        label="Birth Date & Age"
                                        value={`${formatDate(patient.birthdate)} (${calculateAge(patient.birthdate)} years old)`}
                                    />
                                </Box>
                                {patient.phone_number && (
                                    <Box flex={1} minWidth="200px">
                                        <InfoItem
                                            icon={<Phone fontSize={compact ? "small" : "medium"} />}
                                            label="Phone Number"
                                            value={patient.phone_number}
                                        />
                                    </Box>
                                )}
                            </Box>

                            {/* Location and Address side by side */}
                            <Box display="flex" gap={3} flexWrap="wrap">
                                {patient.location && (
                                    <Box flex={1} minWidth="200px">
                                        <InfoItem
                                            icon={<LocationOn fontSize={compact ? "small" : "medium"} />}
                                            label="Location"
                                            value={patient.location}
                                        />
                                    </Box>
                                )}
                                {patient.address && (
                                    <Box flex={1} minWidth="200px">
                                        <InfoItem
                                            icon={<Home fontSize={compact ? "small" : "medium"} />}
                                            label="Address"
                                            value={patient.address}
                                        />
                                    </Box>
                                )}
                            </Box>
                        </Stack>
                    </Box>
                </Box>
           </CardContent>
        </Card>
    );
};

export default PatientInfoCard;
