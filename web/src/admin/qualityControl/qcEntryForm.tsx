import {
    Box as MuiBox,
    Card,
    CardContent,
    Typography,
    Button,
    TextField,
    Stack,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
} from '@mui/material';
import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGetOne, useNotify } from 'react-admin';
import SaveIcon from '@mui/icons-material/Save';
import ScienceIcon from '@mui/icons-material/Science';

interface TestType {
    id: number;
    code: string;
    name: string;
    unit: string;
    reference_range_min?: number;
    reference_range_max?: number;
}

export const QCEntryForm = () => {
    const { deviceId, testTypeId } = useParams<{ deviceId: string; testTypeId: string }>();
    const navigate = useNavigate();
    const notify = useNotify();

    const { data: testType, isLoading } = useGetOne<TestType>('test-type', { id: parseInt(testTypeId || '0') });

    const [qcLevel, setQcLevel] = useState<1 | 2 | 3>(1);
    const [formData, setFormData] = useState({
        sd_target: '',
        lot_number: '',
        range_min: '',
        range_max: '',
    });
    const [saving, setSaving] = useState(false);

    const handleInputChange = (field: string, value: string) => {
        setFormData(prev => ({
            ...prev,
            [field]: value
        }));
    };

    const handleSubmit = async () => {
        try {
            setSaving(true);

            // Validate inputs
            if (!formData.lot_number.trim()) {
                notify('Lot number is required', { type: 'error' });
                return;
            }

            if (!formData.sd_target.trim()) {
                notify('SD Target is required', { type: 'error' });
                return;
            }

            const sdTarget = parseFloat(formData.sd_target);
            if (isNaN(sdTarget) || sdTarget <= 0) {
                notify('Please enter a valid SD Target greater than 0', { type: 'error' });
                return;
            }

            const rangeMin = parseFloat(formData.range_min);
            const rangeMax = parseFloat(formData.range_max);

            if (isNaN(rangeMin) || isNaN(rangeMax)) {
                notify('Please enter valid numbers for range', { type: 'error' });
                return;
            }

            if (rangeMin >= rangeMax) {
                notify('Range min must be less than range max', { type: 'error' });
                return;
            }

            // Calculate target_mean from range
            const targetMean = (rangeMin + rangeMax) / 2;

            if (targetMean <= 0) {
                notify('Target mean must be greater than 0', { type: 'error' });
                return;
            }

            const payload = {
                device_id: parseInt(deviceId || '0'),
                test_type_id: parseInt(testTypeId || '0'),
                qc_level: qcLevel,
                lot_number: formData.lot_number.trim(),
                target_mean: targetMean,
                target_sd: sdTarget,
                ref_min: rangeMin,
                ref_max: rangeMax,
                created_by: 'Admin', // TODO: Get from auth context
            };

            console.log('Sending QC Entry:', payload);

            const response = await fetch('/api/v1/quality-control/entries', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload),
            });

            let responseData;
            const responseText = await response.text();
            console.log('Raw response:', response.status, responseText);

            try {
                responseData = JSON.parse(responseText);
            } catch (e) {
                console.error('Failed to parse JSON:', responseText);
                responseData = { error: responseText || 'Failed to parse server response' };
            }

            console.log('Parsed response:', responseData);

            if (!response.ok) {
                const errorMsg = responseData.message || responseData.error || responseText || `Server error: ${response.status}`;
                throw new Error(errorMsg);
            }

            notify('QC Entry created successfully', { type: 'success' });
            navigate(`/quality-control/${deviceId}/parameter/${testTypeId}`);
        } catch (error: any) {
            console.error('Error creating QC entry:', error);
            notify(error.message || 'Failed to create QC entry', { type: 'error' });
        } finally {
            setSaving(false);
        }
    };

    if (isLoading) {
        return (
            <MuiBox sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <Typography>Loading...</Typography>
            </MuiBox>
        );
    }

    if (!testType) {
        return (
            <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                <Typography variant="h6" color="text.secondary">
                    Test type not found
                </Typography>
            </MuiBox>
        );
    }

    return (
        <MuiBox sx={{ mt: 2 }}>

            {/* Test Type Info */}
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <ScienceIcon sx={{ fontSize: 40, color: '#2196f3' }} />
                        <MuiBox>
                            <Typography variant="h5" sx={{ fontWeight: 600 }}>
                                {testType.code} - {testType.name}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                Unit: {testType.unit} | Reference Range: {testType.reference_range_min} - {testType.reference_range_max}
                            </Typography>
                        </MuiBox>
                    </MuiBox>
                </CardContent>
            </Card>

            {/* QC Entry Form */}
            <Card>
                <CardContent>
                    <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                        Create New QC Entry
                    </Typography>

                    <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                        Configure the QC entry for this test parameter. Enter the acceptable range from the QC material certificate.
                    </Typography>

                    <Stack spacing={2}>
                        <FormControl fullWidth required>
                            <InputLabel>QC Level</InputLabel>
                            <Select
                                value={qcLevel}
                                label="QC Level"
                                onChange={(e) => setQcLevel(e.target.value as 1 | 2 | 3)}
                            >
                                <MenuItem value={1}>Level 1 - Low concentration</MenuItem>
                                <MenuItem value={2}>Level 2 - Normal concentration</MenuItem>
                                <MenuItem value={3}>Level 3 - High concentration</MenuItem>
                            </Select>
                        </FormControl>

                        <TextField
                            fullWidth
                            label="Lot Number"
                            value={formData.lot_number}
                            onChange={(e) => handleInputChange('lot_number', e.target.value)}
                            placeholder="e.g., LOT-2024-1234"
                            required
                            helperText="Enter the lot number from the QC material package"
                        />


                        <MuiBox sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                            <TextField
                                fullWidth
                                label="Reference Range Min"
                                type="number"
                                value={formData.range_min}
                                onChange={(e) => handleInputChange('range_min', e.target.value)}
                                placeholder="e.g., 13.5"
                                required
                                InputProps={{
                                    endAdornment: <Typography variant="body2" color="text.secondary">{testType.unit}</Typography>
                                }}
                                helperText="Minimum acceptable value"
                            />

                            <TextField
                                fullWidth
                                label="Reference Range Max"
                                type="number"
                                value={formData.range_max}
                                onChange={(e) => handleInputChange('range_max', e.target.value)}
                                placeholder="e.g., 15.5"
                                required
                                InputProps={{
                                    endAdornment: <Typography variant="body2" color="text.secondary">{testType.unit}</Typography>
                                }}
                                helperText="Maximum acceptable value"
                            />
                        </MuiBox>

                        <TextField
                            fullWidth
                            label="SD Target"
                            value={formData.sd_target}
                            onChange={(e) => handleInputChange('sd_target', e.target.value)}
                            placeholder="e.g., 3.23"
                            required
                            helperText="Enter the SD target from the Product Documentation"
                        />



                        {formData.range_min && formData.range_max && (
                            <MuiBox sx={{
                                p: 2,
                                borderRadius: 1,
                                backgroundColor: 'rgba(33, 150, 243, 0.05)',
                                border: '1px solid rgba(33, 150, 243, 0.2)'
                            }}>
                                <Typography variant="body2" color="text.secondary">
                                    <strong>Target Mean (calculated):</strong> {((parseFloat(formData.range_min) + parseFloat(formData.range_max)) / 2).toFixed(2)} {testType.unit}
                                </Typography>
                                <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 1 }}>
                                    SD will be calculated automatically from incoming measurements
                                </Typography>
                            </MuiBox>
                        )}

                        <MuiBox sx={{
                            p: 2,
                            borderRadius: 1,
                            backgroundColor: 'rgba(33, 150, 243, 0.1)',
                            border: '1px solid rgba(33, 150, 243, 0.3)'
                        }}>
                            <Typography variant="body2" color="text.secondary">
                                <strong>Note:</strong> Creating this entry will deactivate any previous QC entry for the same device, test type, and level.
                                All new QC measurements will use this configuration until a new entry is created.
                            </Typography>
                        </MuiBox>

                        <MuiBox sx={{ display: 'flex', gap: 2 }}>
                            <Button
                                variant="outlined"
                                onClick={() => navigate(`/quality-control/${deviceId}/parameter/${testTypeId}`)}
                                size="large"
                                fullWidth
                            >
                                Cancel
                            </Button>
                            <Button
                                variant="contained"
                                startIcon={<SaveIcon />}
                                onClick={handleSubmit}
                                disabled={!formData.lot_number || !formData.range_min || !formData.range_max || !formData.sd_target || saving}
                                size="large"
                                fullWidth
                            >
                                {saving ? 'Creating...' : 'Create QC Entry'}
                            </Button>
                        </MuiBox>
                    </Stack>
                </CardContent>
            </Card>
        </MuiBox>
    );
};
