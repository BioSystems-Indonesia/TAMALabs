import { useState } from "react";
import {
    Box,
    Button,
    Card,
    CardContent,
    TextField,
    Typography,
    Alert,
} from "@mui/material";

export default function LicensePage() {
    const [licenseCode, setLicenseCode] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const handleActivate = async () => {
        setError("");
        setSuccess("");
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
                    >
                        Enter your license code to activate the application.
                    </Typography>

                    <TextField
                        label="License Code"
                        variant="outlined"
                        fullWidth
                        value={licenseCode}
                        onChange={(e) => setLicenseCode(e.target.value)}
                        sx={{
                            mb: 2, mt: 2, "& label.Mui-focused": {
                                color: "#4abaab",
                            },
                            "& .MuiOutlinedInput-root": {
                                "&.Mui-focused fieldset": {
                                    borderColor: "#4abaab",
                                },
                            },
                        }}
                    />

                    <Button
                        fullWidth
                        variant="contained"
                        onClick={handleActivate}
                        sx={{
                            borderRadius: 2,
                            padding: '12px 24px',
                            fontSize: '1.1rem',
                            fontWeight: 'bold',
                            background: '#4abaab',
                            color: 'white',
                            '&:hover': {
                                background: 'rgba(76, 192, 176, 1)'
                            }
                        }}
                    >
                        Activate
                    </Button>
                </CardContent>
            </Card>
        </Box>
    );
}
