import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, MenuItem, CircularProgress, Alert } from '@mui/material';
import { useState } from 'react';
import { useNotify } from 'react-admin';

interface ManualQCInputDialogProps {
    open: boolean;
    onClose: () => void;
    deviceId: string;
    testTypeId: string;
    activeEntries: Array<{
        id: number;
        qc_level: 1 | 2 | 3;
        lot_number: string;
        target_mean: number;
        ref_min: number;
        ref_max: number;
    }>;
    onSuccess: () => void;
}

export const ManualQCInputDialog = ({
    open,
    onClose,
    deviceId,
    testTypeId,
    activeEntries,
    onSuccess,
}: ManualQCInputDialogProps) => {
    const notify = useNotify();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        qc_level: 1,
        measured_value: '',
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.measured_value) {
            notify('Please fill all required fields', { type: 'warning' });
            return;
        }

        setLoading(true);

        try {
            const response = await fetch('/api/v1/quality-control/results/manual', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    device_id: parseInt(deviceId),
                    test_type_id: parseInt(testTypeId),
                    qc_level: formData.qc_level,
                    measured_value: parseFloat(formData.measured_value),
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Failed to create QC result');
            }

            await response.json();
            notify('Manual QC result created successfully', { type: 'success' });

            // Reset form
            setFormData({
                qc_level: 1,
                measured_value: '',
            });

            onSuccess();
            onClose();
        } catch (error: any) {
            notify(error.message || 'Error creating manual QC result', { type: 'error' });
        } finally {
            setLoading(false);
        }
    };

    const selectedEntry = activeEntries.find(e => e.qc_level === formData.qc_level);

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <form onSubmit={handleSubmit}>
                <DialogTitle>Manual QC Input</DialogTitle>
                <DialogContent>
                    {activeEntries.length === 0 ? (
                        <Alert severity="warning" sx={{ mt: 2 }}>
                            No active QC entry found. Please create a QC entry first.
                        </Alert>
                    ) : (
                        <>
                            <TextField
                                select
                                fullWidth
                                label="QC Level"
                                value={formData.qc_level}
                                onChange={(e) => setFormData({ ...formData, qc_level: parseInt(e.target.value) as 1 | 2 | 3 })}
                                margin="normal"
                                required
                            >
                                {activeEntries.map(entry => (
                                    <MenuItem key={entry.qc_level} value={entry.qc_level}>
                                        Level {entry.qc_level} - Lot: {entry.lot_number}
                                    </MenuItem>
                                ))}
                            </TextField>

                            {selectedEntry && (
                                <Alert severity="info" sx={{ mt: 2, mb: 2 }}>
                                    <strong>Reference Range:</strong> {selectedEntry.ref_min.toFixed(2)} - {selectedEntry.ref_max.toFixed(2)}
                                    <br />
                                    <strong>Target Mean:</strong> {selectedEntry.target_mean.toFixed(2)}
                                </Alert>
                            )}

                            <TextField
                                fullWidth
                                label="Measured Value"
                                type="number"
                                value={formData.measured_value}
                                onChange={(e) => setFormData({ ...formData, measured_value: e.target.value })}
                                margin="normal"
                                required
                                inputProps={{
                                    step: '0.01',
                                }}
                            />
                        </>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={onClose} disabled={loading}>
                        Cancel
                    </Button>
                    <Button
                        type="submit"
                        variant="contained"
                        disabled={loading || activeEntries.length === 0}
                    >
                        {loading ? <CircularProgress size={24} /> : 'Submit'}
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};
