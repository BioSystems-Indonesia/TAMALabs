
import React, { useState } from 'react';
import {
    useShowController,
    SimpleShowLayout,
    TextField,
    RecordContextProvider
} from 'react-admin';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography, Table, TableBody, TableRow, TableCell } from '@mui/material';

interface CustomShowDialogProps {
    resource: string;
    recordId: string | number;
}

const CustomShowDialog: React.FC<CustomShowDialogProps> = ({ resource, recordId }) => {
    const [open, setOpen] = useState(false);
    const { record, isLoading } = useShowController({ id: recordId, resource });

    const handleOpen = () => setOpen(true);
    const handleClose = () => setOpen(false);

    return (
        <>
            <Button variant="outlined" onClick={handleOpen}>
                Show Test Results
            </Button>
            <Dialog open={open} onClose={handleClose} fullWidth maxWidth="md">
                <DialogTitle>Test Results</DialogTitle>
                <DialogContent>
                    {isLoading ? (
                        <p>Loading...</p>
                    ) : record ? (
                        <RecordContextProvider value={record}>
                            <SimpleShowLayout>
                                <TextField source="barcode" label="Barcode" />
                                <TextField source="patient_name" label="Patient Name" />
                                <TextField source="patient_id" label="Patient ID" />
                                <Typography variant="h6" gutterBottom>
                                    Biochemistry
                                </Typography>
                                {record.detail.biochemistry?.length > 0 ? (
                                    <Table>
                                        <TableBody>
                                            {record.detail.biochemistry.map((test: any, index: number) => (
                                                <TableRow key={index}>
                                                    <TableCell>{test.test}</TableCell>
                                                    <TableCell>{test.result}</TableCell>
                                                    <TableCell>{test.unit}</TableCell>
                                                    <TableCell>{test.computed_reference_range || test.reference_range}</TableCell>
                                                    <TableCell>
                                                        {test.abnormal === 1 ? (
                                                            <Typography color="error">Abnormal</Typography>
                                                        ) : (
                                                            <Typography color="primary">Normal</Typography>
                                                        )}
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                ) : (
                                    <Typography>No Biochemistry Data</Typography>
                                )}
                                <Typography variant="h6" gutterBottom>
                                    Observation
                                </Typography>
                                {record.detail.observation?.length > 0 ? (
                                    <Table>
                                        <TableBody>
                                            {record.detail.observation.map((test: any, index: number) => (
                                                <TableRow key={index}>
                                                    <TableCell>{test.test}</TableCell>
                                                    <TableCell>{test.result}</TableCell>
                                                    <TableCell>{test.unit}</TableCell>
                                                    <TableCell>{test.computed_reference_range || test.reference_range}</TableCell>
                                                    <TableCell>
                                                        {test.abnormal === 1 ? (
                                                            <Typography color="error">Abnormal</Typography>
                                                        ) : (
                                                            <Typography color="primary">Normal</Typography>
                                                        )}
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                ) : (
                                    <Typography>No Observation Data</Typography>
                                )}
                            </SimpleShowLayout>
                        </RecordContextProvider>
                    ) : (
                        <Typography>No record found</Typography>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Close</Button>
                </DialogActions>
            </Dialog>
        </>
    );
};

export default CustomShowDialog;