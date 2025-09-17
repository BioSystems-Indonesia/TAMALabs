import React, { useState } from 'react';
import {
    Avatar,
    Box,
    Button,
    CardContent,
    Container,
    TextField,
    Typography,
    Alert,
    CircularProgress,
    InputAdornment,
    IconButton,
    Paper,
    Stack,
    ThemeProvider,
} from '@mui/material';
import {
    Lock as LockIcon,
    Person as PersonIcon,
    Visibility,
    VisibilityOff,
    Login as LoginIcon
} from '@mui/icons-material';
import { useLogin, useNotify } from 'react-admin';
import { useForm, Controller } from 'react-hook-form';
import logo from '../../assets/alinda-husada-logo.png'
import { radiantLightTheme } from '../theme.tsx';

interface LoginFormData {
    username: string;
    password: string;
}

const CustomLoginPage: React.FC = () => {
    const [showPassword, setShowPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const login = useLogin();
    const notify = useNotify();

    const { control, handleSubmit, formState: { errors }, setValue } = useForm<LoginFormData>({
        defaultValues: {
            username: '',
            password: ''
        }
    });

    const handleLogin = async (data: LoginFormData) => {
        setLoading(true);
        setError(null);

        try {
            await login(data);
            notify('Login success', { type: 'success' });
        } catch (error) {
            if (error instanceof Error) {
                setError('Username or password wrong');
                setValue('password', ''); 
            }
            notify('Login failed', { type: 'error' });
        } finally {
            setLoading(false);
        }
    };

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    return (
        <ThemeProvider theme={radiantLightTheme}>
            <Box
                sx={{
                    minHeight: '100vh',
                    background: 'white',
                    display: 'flex',
                    width: '100%',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: 2
                }}
            >
                <Container maxWidth="lg">
                    <Paper
                        elevation={5}
                        sx={{
                            display: 'flex',
                            height: '50vh',
                            borderRadius: 4,
                            overflow: 'hidden',
                            backdropFilter: 'blur(10px)',

                        }}
                    >
                        <Box
                            sx={{
                                background: 'linear-gradient(135deg, #4abaab 0%, #2fa091ff 100%)',
                                padding: 4,
                                display: 'flex',
                                flexDirection: 'column',
                                justifyContent: 'center',
                                alignItems: 'center',
                                textAlign: 'center',
                                color: 'white',
                                width: '50%'
                            }}
                        >
                            <Avatar
                                sx={{
                                    width: 80,
                                    height: 80,
                                    margin: '0 auto 16px',
                                    background: 'rgba(255, 255, 255)',
                                    backdropFilter: 'blur(10px)',
                                }}
                            >
                                <img src={logo} width={75} />
                            </Avatar>
                            <Typography variant="h5" fontWeight="bold" gutterBottom>
                                TAMALabs
                            </Typography>
                            <Typography variant="body1" sx={{ opacity: 0.9 }}>
                                Laboratory Information Management System TAMALabs
                            </Typography>
                        </Box>

                        <CardContent
                            sx={{
                                paddingX: 6,
                                width: '50%',
                                margin: 'auto',
                                flexDirection: 'column'
                            }}>
                            <Typography
                                variant="h4"
                                gutterBottom
                                textAlign="center"
                                color="text.primary"
                                fontWeight="600"
                            >
                                Sign In
                            </Typography>

                            <Typography
                                variant="body2"
                                textAlign="center"
                                color="text.secondary"
                                sx={{ mb: 3, fontSize: 13 }}
                            >
                                Please enter your username and password
                            </Typography>

                            {error && (
                                <Alert severity="error" sx={{ mb: 2 }}>
                                    {error}
                                </Alert>
                            )}

                            <form onSubmit={handleSubmit(handleLogin)}>
                                <Stack spacing={3}>
                                    <Controller
                                        name="username"
                                        control={control}
                                        rules={{
                                            required: 'Username required',
                                            minLength: {
                                                value: 3,
                                                message: 'Username must be at least 3 characters long'
                                            }
                                        }}
                                        render={({ field }) => (
                                            <TextField
                                                {...field}
                                                label="Username"
                                                fullWidth
                                                error={!!errors.username}
                                                helperText={errors.username?.message}
                                                InputProps={{
                                                    startAdornment: (
                                                        <InputAdornment position="start">
                                                            <PersonIcon color="action" />
                                                        </InputAdornment>
                                                    ),
                                                }}
                                                InputLabelProps={{
                                                    sx: {
                                                        fontSize: '1.2rem',
                                                        color: 'text.primary',
                                                        '&.Mui-focused': {
                                                            fontSize: '1.2rem',
                                                            color: 'primary.main',
                                                        }
                                                    }
                                                }}
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        borderRadius: 2,
                                                        minWidth: '100%',
                                                        height: 60,
                                                        fontSize: '1.2rem'
                                                    }
                                                }}
                                            />
                                        )}
                                    />

                                    <Controller
                                        name="password"
                                        control={control}
                                        rules={{
                                            required: 'Password required',
                                            minLength: {
                                                value: 6,
                                                message: 'Password must be at least 6 characters long'
                                            }
                                        }}
                                        render={({ field }) => (
                                            <TextField
                                                {...field}
                                                label="Password"
                                                type={showPassword ? 'text' : 'password'}
                                                fullWidth
                                                error={!!errors.password}
                                                helperText={errors.password?.message}
                                                InputProps={{
                                                    startAdornment: (
                                                        <InputAdornment position="start">
                                                            <LockIcon color="action" />
                                                        </InputAdornment>
                                                    ),
                                                    endAdornment: (
                                                        <InputAdornment position="end">
                                                            <IconButton
                                                                onClick={handleClickShowPassword}
                                                                edge="end"
                                                            >
                                                                {showPassword ? <VisibilityOff /> : <Visibility />}
                                                            </IconButton>
                                                        </InputAdornment>
                                                    ),
                                                }}
                                                InputLabelProps={{
                                                    sx: {
                                                        fontSize: '1.2rem',
                                                        color: 'text.primary',
                                                        '&.Mui-focused': {
                                                            fontSize: '1.2rem',
                                                            color: 'primary.main',
                                                        }
                                                    }
                                                }}
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        borderRadius: 2,
                                                        height: 60,
                                                        fontSize: '1.2rem'
                                                    }
                                                }}
                                            />
                                        )}
                                    />

                                    <Button
                                        type="submit"
                                        variant="contained"
                                        fullWidth
                                        size="large"
                                        disabled={loading}
                                        startIcon={loading ? <CircularProgress size={20} /> : <LoginIcon />}
                                        sx={{
                                            borderRadius: 2,
                                            padding: '12px 24px',
                                            fontSize: '1.1rem',
                                            fontWeight: 'bold',
                                            background: 'linear-gradient(135deg, #4abaab 0%, #2fa091ff 100%)',
                                            color: 'white',
                                            '&:hover': {
                                                background: 'linear-gradient(135deg, #42b0a1ff 0%, #2e9688ff 100%)',
                                            }
                                        }}
                                    >
                                        {loading ? 'Logging in...' : 'Login'}
                                    </Button>
                                </Stack>
                            </form>
                        </CardContent>

                        {/* <Box
                        sx={{
                            textAlign: 'center',
                            padding: 2,
                            borderTop: '1px solid',
                            borderColor: 'divider',
                            background: 'rgba(0, 0, 0, 0.02)'
                        }}
                    >
                        <Typography variant="body2" color="text.secondary">
                            © 2025 RS Alinda Husada. All rights reserved.
                        </Typography>
                    </Box> */}
                    </Paper>
                </Container>
            </Box>
        </ThemeProvider>
    );
};

export default CustomLoginPage;
