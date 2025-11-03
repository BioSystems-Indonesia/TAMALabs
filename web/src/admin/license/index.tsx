import { useState, useEffect } from "react";
import {
    Box,
    Button,
    Card,
    CardContent,
    TextField,
    Typography,
    Alert,
    CircularProgress,
} from "@mui/material";

export default function LicensePage() {
    const [licenseCode, setLicenseCode] = useState(["", "", "", ""]);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [licenseStatus, setLicenseStatus] = useState<{
        checking: boolean;
        valid: boolean;
        message: string;
    }>({
        checking: true,
        valid: false,
        message: ""
    });

    useEffect(() => {
        checkLicenseStatus();
    }, []);

    const checkLicenseStatus = async () => {
        try {
            setLicenseStatus(prev => ({ ...prev, checking: true }));

            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/license/check`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                throw new Error('Failed to check license status');
            }

            const data = await response.json();

            setLicenseStatus({
                checking: false,
                valid: data.valid || false,
                message: data.message || "Unknown status"
            });

            if (data.valid) {
                setSuccess("License is already activated and valid!");
                setTimeout(() => {
                    window.location.reload();
                }, 2000);
            }

        } catch (err) {
            console.error('License check error:', err);
            setLicenseStatus({
                checking: false,
                valid: false,
                message: err instanceof Error ? err.message : 'Failed to check license status'
            });
        }
    };

    const handleLicenseChange = (index: number, value: string) => {
        const cleanValue = value.toUpperCase().replace(/[^A-Z0-9]/g, '').slice(0, 4);
        const newLicenseCode = [...licenseCode];
        newLicenseCode[index] = cleanValue;
        setLicenseCode(newLicenseCode);

        if (cleanValue.length === 4 && index < 3) {
            const nextField = document.getElementById(`license-field-${index + 1}`);
            nextField?.focus();
        }
    };

    const handleActivate = async () => {
        setError("");
        setSuccess("");
        setIsLoading(true);

        try {
            const fullLicenseKey = licenseCode.join('-');
            if (licenseCode.some(code => code.length !== 4)) {
                setError("Please enter all 4 parts of the license code");
                setIsLoading(false);
                return;
            }

            const requestBody = {
                license_code: fullLicenseKey
            };

            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/license/activate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody)
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'License activation failed');
            }

            const data = await response.json();

            if (data.payload && data.signature) {
                setSuccess("License activated successfully!");
                setLicenseCode(["", "", "", ""]);

                setTimeout(async () => {
                    await checkLicenseStatus();
                }, 1000);

            } else {
                throw new Error('Invalid response format from server');
            }

        } catch (err) {
            console.error('License activation error:', err);
            setError(err instanceof Error ? err.message : 'License activation failed. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box
            sx={{
                minHeight: "100vh",
                bgcolor: "grey.100",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                p: 2,
            }}
        >
            <Card sx={{ maxWidth: 400, width: "100%", boxShadow: 3 }}>
                <CardContent>
                    {/* License Status Check */}
                    {licenseStatus.checking && (
                        <Alert severity="info" sx={{ mb: 2, display: 'flex', alignItems: 'center' }}>
                            <CircularProgress size={20} sx={{ mr: 1 }} />
                            Checking license status...
                        </Alert>
                    )}

                    {!licenseStatus.checking && licenseStatus.valid && (
                        <Alert severity="success" sx={{ mb: 2 }}>
                            {licenseStatus.message}
                        </Alert>
                    )}

                    {!licenseStatus.checking && !licenseStatus.valid && (
                        <Alert severity="warning" sx={{ mb: 2 }}>
                            {licenseStatus.message}
                        </Alert>
                    )}

                    {error && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {error}
                        </Alert>
                    )}
                    {success && (
                        <Alert severity="success" sx={{ mb: 2 }}>
                            {success}
                        </Alert>
                    )}

                    <Typography variant="h5" align="center" gutterBottom>
                        License Activation
                    </Typography>
                    <Typography
                        variant="body2"
                        align="center"
                        color="text.secondary"
                        gutterBottom
                        sx={{ mb: 3 }}
                    >
                        Enter your license code to activate the application.
                    </Typography>

                    {/* License Code Input - 4 separate fields */}
                    <Box
                        sx={{
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                            mb: 3,
                            mt: 2
                        }}
                    >
                        {licenseCode.map((code, index) => (
                            <Box key={index} sx={{ display: 'flex', alignItems: 'center' }}>
                                <TextField
                                    id={`license-field-${index}`}
                                    variant="outlined"
                                    value={code}
                                    onChange={(e) => handleLicenseChange(index, e.target.value)}
                                    placeholder="XXXX"
                                    inputProps={{
                                        maxLength: 4,
                                        style: {
                                            textAlign: 'center',
                                            fontSize: '16px',
                                        }
                                    }}
                                    sx={{
                                        width: '70px',
                                        "& label.Mui-focused": {
                                            color: "#4abaab",
                                        },
                                        "& .MuiOutlinedInput-root": {
                                            "&.Mui-focused fieldset": {
                                                borderColor: "#4abaab",
                                            },
                                        },
                                    }}
                                />
                                {index < 3 && (
                                    <Typography
                                        variant="h6"
                                        sx={{
                                            mx: 1,
                                            color: 'text.secondary',
                                            fontWeight: 'bold'
                                        }}
                                    >
                                        -
                                    </Typography>
                                )}
                            </Box>
                        ))}
                    </Box>

                    <Typography
                        variant="caption"
                        align="center"
                        color="text.secondary"
                        display="block"
                        sx={{ mb: 2 }}
                    >
                    </Typography>

                    <Button
                        fullWidth
                        variant="contained"
                        onClick={handleActivate}
                        disabled={isLoading || licenseStatus.valid}
                        sx={{
                            borderRadius: 2,
                            padding: '12px 24px',
                            fontSize: '1.1rem',
                            fontWeight: 'bold',
                            background: (isLoading || licenseStatus.valid) ? '#ccc' : '#4abaab',
                            color: 'white',
                            '&:hover': {
                                background: (isLoading || licenseStatus.valid) ? '#ccc' : 'rgba(76, 192, 176, 1)'
                            },
                            mb: 2
                        }}
                    >
                        {isLoading ? 'Activating...' :
                            licenseStatus.valid ? 'License Already Active' : 'Activate'}
                    </Button>

                    {/* Re-check License Button */}
                    <Button
                        fullWidth
                        variant="outlined"
                        onClick={checkLicenseStatus}
                        disabled={licenseStatus.checking}
                        sx={{
                            borderRadius: 2,
                            padding: '8px 16px',
                            fontSize: '0.9rem',
                            borderColor: '#4abaab',
                            color: '#4abaab',
                            '&:hover': {
                                borderColor: '#4abaab',
                                backgroundColor: 'rgba(74, 186, 171, 0.1)'
                            }
                        }}
                    >
                        {licenseStatus.checking ? (
                            <>
                                <CircularProgress size={16} sx={{ mr: 1 }} />
                                Checking...
                            </>
                        ) : (
                            'Re-check License Status'
                        )}
                    </Button>
                </CardContent>
            </Card>
        </Box>
    );
}