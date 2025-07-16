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
    Stack
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
import  logo  from '../../assets/elgatama-logo.png'

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
            notify('Login berhasil', { type: 'success' });
        } catch (error) {
            if (error instanceof Error) {
                setError('Username atau password salah');
                setValue('password', ''); 
            }
            notify('Login gagal', { type: 'error' });
        } finally {
            setLoading(false);
        }
    };

    const handleClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    return (
        <Box
            sx={{
                minHeight: '100vh',
                background: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                padding: 2
            }}
        >
            <Container maxWidth="sm">
                <Paper
                    elevation={5}
                    sx={{
                        borderRadius: 4,
                        overflow: 'hidden',
                        background: 'rgba(255, 255, 255, 0.95)',
                        backdropFilter: 'blur(10px)',

                    }}
                >
                    <Box
                        sx={{
                            background: 'linear-gradient(135deg, #4abaab 0%, #2fa091ff 100%)',
                            padding: 4,
                            textAlign: 'center',
                            color: 'white',
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
                        <Typography variant="h4" fontWeight="bold" gutterBottom>
                           Laboratory Information System
                        </Typography>
                        <Typography variant="body1" sx={{ opacity: 0.9, padding: '0 4rem' }}>
                            Laboratory Information Management System Elga Tama
                        </Typography>
                    </Box>

                    <CardContent sx={{ padding: 4 }}>
                        <Typography 
                            variant="h5" 
                            gutterBottom 
                            textAlign="center"
                            color="text.primary"
                            fontWeight="600"
                        >
                            Masuk ke Sistem
                        </Typography>
                        
                        <Typography 
                            variant="body2" 
                            textAlign="center" 
                            color="text.secondary"
                            sx={{ mb: 3 }}
                        >
                            Silakan masukkan username dan password Anda
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
                                        required: 'Username wajib diisi',
                                        minLength: {
                                            value: 3,
                                            message: 'Username minimal 3 karakter'
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
                                            sx={{
                                                '& .MuiOutlinedInput-root': {
                                                    borderRadius: 2,
                                                }
                                            }}
                                        />
                                    )}
                                />

                                <Controller
                                    name="password"
                                    control={control}
                                    rules={{ 
                                        required: 'Password wajib diisi',
                                        minLength: {
                                            value: 6,
                                            message: 'Password minimal 6 karakter'
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
                                            sx={{
                                                '& .MuiOutlinedInput-root': {
                                                    borderRadius: 2,
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
                                    {loading ? 'Sedang masuk...' : 'Masuk'}
                                </Button>
                            </Stack>
                        </form>
                    </CardContent>

                    <Box
                        sx={{
                            textAlign: 'center',
                            padding: 2,
                            borderTop: '1px solid',
                            borderColor: 'divider',
                            background: 'rgba(0, 0, 0, 0.02)'
                        }}
                    >
                        <Typography variant="body2" color="text.secondary">
                            © 2025 PT Elga Tama. All rights reserved.
                        </Typography>
                    </Box>
                </Paper>
            </Container>
        </Box>
    );
};

export default CustomLoginPage;