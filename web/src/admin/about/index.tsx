import { 
    Box, 
    Card, 
    CardContent, 
    Typography, 
    GridLegacy as Grid, 
    Avatar,
    Stack,
    useTheme
} from '@mui/material';
import BusinessIcon from '@mui/icons-material/Business';
import EmailIcon from '@mui/icons-material/Email';
import PhoneIcon from '@mui/icons-material/Phone';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import ScienceIcon from '@mui/icons-material/Science';
import VerifiedIcon from '@mui/icons-material/Verified';
import SecurityIcon from '@mui/icons-material/Security';
import SpeedIcon from '@mui/icons-material/Speed';
import logo from '../../assets/elgatama-logo.png';

export const AboutPage = () => {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';


    return (
        <Box sx={{ 
            margin: '0 auto'
        }}>
            {/* Header Section */}
            <Card sx={{ 
                marginBottom: 4,
                background: isDarkMode 
                    ? `linear-gradient(135deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 100%)`
                    : `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.primary.dark} 100%)`,
                color: 'white',
            }}>
                <CardContent sx={{ padding: 4 }}>
                    <Box sx={{ 
                        display: 'flex', 
                        alignItems: 'center', 
                        gap: 3,
                        marginBottom: 3
                    }}>
                        <Avatar sx={{ 
                            width: 80, 
                            height: 80,
                            backgroundColor: 'white',
                            padding: 0.5
                        }}>
                            <img 
                                src={logo} 
                                alt="Elga Tama Logo" 
                                style={{ 
                                    width: '100%',
                                    height: '100%',
                                    objectFit: 'contain',
                                }} 
                            />
                        </Avatar>
                        <Box>
                            <Typography variant="h3" sx={{ 
                                fontWeight: 700,
                                marginBottom: 1
                            }}>
                                PT ELGA TAMA
                            </Typography>
                            <Typography variant="h6" sx={{ 
                                opacity: 0.9,
                                fontWeight: 400
                            }}>
                                Laboratory Information Management System
                            </Typography>
                        </Box>
                    </Box>
                    
                    <Typography variant="body1" sx={{ 
                        fontSize: '1.1rem',
                        lineHeight: 1.6,
                        opacity: 0.95
                    }}>
                        <strong>PT Elga Tama</strong> provides a smart and reliable <strong>Laboratory Information System (LIS)</strong> designed to simplify and automate laboratory processes. From sample tracking to result reporting, our system ensures speed, accuracy, and seamless integration with hospital systems and lab instruments.
                    </Typography>
                    <Typography variant="body1" sx={{ 
                        fontSize: '1.1rem',
                        lineHeight: 1.6,
                        opacity: 0.95
                    }}>
                        We aim to support healthcare providers with efficient digital tools that improve diagnostic services and patient care.
                    </Typography>
                </CardContent>
            </Card>

            <Grid container spacing={4}>
                {/* Company Information */}
                <Grid item xs={12} md={6}>
                    <Card sx={{ 
                        height: '100%',
                        borderRadius: 2
                    }}>
                        <CardContent sx={{ padding: 3 }}>
                            <Box sx={{ 
                                display: 'flex', 
                                alignItems: 'center', 
                                gap: 2,
                                marginBottom: 3
                            }}>
                                <BusinessIcon sx={{ 
                                    fontSize: 32, 
                                    color: theme.palette.primary.main 
                                }} />
                                <Typography variant="h5" sx={{ fontWeight: 600 }}>
                                    Company Information
                                </Typography>
                            </Box>
                            
                            <Stack spacing={2.5}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <EmailIcon sx={{ color: theme.palette.text.secondary }} />
                                    <Box>
                                        <Typography variant="body2" color="text.secondary">
                                            Email
                                        </Typography>
                                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                            help@elgatama.com
                                        </Typography>
                                    </Box>
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <PhoneIcon sx={{ color: theme.palette.text.secondary }} />
                                    <Box>
                                        <Typography variant="body2" color="text.secondary">
                                            Phone
                                        </Typography>
                                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                            (+62) 8193 8123 234
                                        </Typography>
                                    </Box>
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <LocationOnIcon sx={{ color: theme.palette.text.secondary }} />
                                    <Box>
                                        <Typography variant="body2" color="text.secondary">
                                            Address
                                        </Typography>
                                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                            JI. Kyai Caringin No.18A-20, RT.11/RW.4, Cideng, Kecamatan Gambir Daerah Khusus Ibukota Jakarta 10150
                                        </Typography>
                                    </Box>
                                </Box>
                            </Stack>
                        </CardContent>
                    </Card>
                </Grid>

                {/* System Features */}
                <Grid item xs={12} md={6}>
                    <Card sx={{ 
                        height: '100%',
                        borderRadius: 2
                    }}>
                        <CardContent sx={{ padding: 3 }}>
                            <Box sx={{ 
                                display: 'flex', 
                                alignItems: 'center', 
                                gap: 2,
                                marginBottom: 3
                            }}>
                                <ScienceIcon sx={{ 
                                    fontSize: 32, 
                                    color: theme.palette.success.main 
                                }} />
                                <Typography variant="h5" sx={{ fontWeight: 600 }}>
                                    System Features
                                </Typography>
                            </Box>
                            
                            <Stack spacing={2}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <VerifiedIcon sx={{ color: theme.palette.success.main, fontSize: 20 }} />
                                    <Typography variant="body1">
                                        Sample Tracking & Management
                                    </Typography>
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <SecurityIcon sx={{ color: theme.palette.info.main, fontSize: 20 }} />
                                    <Typography variant="body1">
                                        Secure Data Management
                                    </Typography>
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <SpeedIcon sx={{ color: theme.palette.warning.main, fontSize: 20 }} />
                                    <Typography variant="body1">
                                        Real-time Result Processing
                                    </Typography>
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                                    <VerifiedIcon sx={{ color: theme.palette.primary.main, fontSize: 20 }} />
                                    <Typography variant="body1">
                                        Quality Control & Compliance
                                    </Typography>
                                </Box>
                            </Stack>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};

