import { useEffect, useState } from 'react';
import { Box as MuiBox, Card, CardContent, Typography, CircularProgress, Chip, Dialog, DialogTitle, DialogContent, DialogActions, Button, Radio, RadioGroup, FormControlLabel, Tabs, Tab, Badge, IconButton } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { useGetOne, useGetList, useNotify, Datagrid, WithRecord } from 'react-admin';
import { Device } from '../../types/device';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import HistoryIcon from '@mui/icons-material/History';

interface TestType {
    id: number;
    code: string;
    name: string;
    unit: string;
    reference_range_min?: number;
    reference_range_max?: number;
    device_id?: number;
}

interface QCParameter {
    id: number;
    device_id: number;
    test_type_id: number;
    test_type?: TestType;
    level1_today: boolean;
    level2_today: boolean;
    level3_today: boolean;
    level1_value?: number;
    level2_value?: number;
    level3_value?: number;
    level1_results?: QCResult[];
    level2_results?: QCResult[];
    level3_results?: QCResult[];
    level1_selected_id?: number;
    level2_selected_id?: number;
    level3_selected_id?: number;
    level1_entry_id?: number;
    level2_entry_id?: number;
    level3_entry_id?: number;
}

interface QCResult {
    id: number;
    qc_entry_id: number;
    measured_value: number;
    result: string;
    method: string;
    calculated_mean?: number;
    calculated_sd?: number;
    calculated_cv?: number;
    error_sd?: number;
    created_at: string;
    created_by?: string;
    operator: string;
    qc_entry?: {
        id: number;
        qc_level: number;
        level1_selected_result_id?: number;
        level2_selected_result_id?: number;
        level3_selected_result_id?: number;
    };
}

export const QualityControlDetail = () => {
    const { id } = useParams<{ id: string }>();
    const deviceId = parseInt(id || '0');
    const navigate = useNavigate();
    const [qcParameters, setQcParameters] = useState<QCParameter[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectionDialog, setSelectionDialog] = useState<{
        open: boolean;
        testTypeId?: number;
        level1Results?: QCResult[];
        level2Results?: QCResult[];
        level3Results?: QCResult[];
        level1EntryId?: number;
        level2EntryId?: number;
        level3EntryId?: number;
    }>({ open: false });
    const [activeTab, setActiveTab] = useState(0);
    const notify = useNotify();

    // Fetch device data
    const { data: device, isLoading: deviceLoading } = useGetOne<Device>('device', { id: deviceId });

    const { data: testTypes, isLoading: testTypesLoading } = useGetList<TestType>('test-type', {
        filter: { device_id: deviceId },
        pagination: { page: 1, perPage: 1000 },
        sort: { field: 'code', order: 'ASC' }
    });

    useEffect(() => {
        const fetchQCStatus = async () => {
            if (!testTypes || testTypes.length === 0) {
                setLoading(false);
                return;
            }

            try {
                setLoading(true);

                // Get today's date at midnight
                const today = new Date();
                today.setHours(0, 0, 0, 0);

                // Fetch QC results for today for all test types
                const qcDataPromises = testTypes.map(async (testType) => {
                    const response = await fetch(
                        `/api/v1/quality-control/results?device_id=${deviceId}&test_type_id=${testType.id}`,
                        { credentials: 'include' }
                    );

                    if (!response.ok) {
                        return {
                            id: testType.id,
                            device_id: deviceId,
                            test_type_id: testType.id,
                            test_type: testType,
                            level1_today: false,
                            level2_today: false,
                            level3_today: false,
                        };
                    }

                    const data = await response.json();
                    const results: QCResult[] = data.data || [];

                    // Check which levels have results today
                    const todayResults = results.filter((result) => {
                        const resultDate = new Date(result.created_at);
                        resultDate.setHours(0, 0, 0, 0);
                        return resultDate.getTime() === today.getTime();
                    });

                    const level1Results = todayResults.filter(r => r.qc_entry?.qc_level === 1);
                    const level2Results = todayResults.filter(r => r.qc_entry?.qc_level === 2);
                    const level3Results = todayResults.filter(r => r.qc_entry?.qc_level === 3);

                    // Get QC entries for each level
                    const level1Entry = level1Results.length > 0 ? level1Results[0].qc_entry : undefined;
                    const level2Entry = level2Results.length > 0 ? level2Results[0].qc_entry : undefined;
                    const level3Entry = level3Results.length > 0 ? level3Results[0].qc_entry : undefined;

                    // Get selected result IDs from QC entries
                    const level1SelectedId = level1Entry?.level1_selected_result_id;
                    const level2SelectedId = level2Entry?.level2_selected_result_id;
                    const level3SelectedId = level3Entry?.level3_selected_result_id;

                    // Get selected or latest result for each level
                    const level1Result = level1SelectedId
                        ? level1Results.find(r => r.id === level1SelectedId) || level1Results[0]
                        : level1Results[0];
                    const level2Result = level2SelectedId
                        ? level2Results.find(r => r.id === level2SelectedId) || level2Results[0]
                        : level2Results[0];
                    const level3Result = level3SelectedId
                        ? level3Results.find(r => r.id === level3SelectedId) || level3Results[0]
                        : level3Results[0];

                    return {
                        id: testType.id,
                        device_id: deviceId,
                        test_type_id: testType.id,
                        test_type: testType,
                        level1_today: !!level1Result,
                        level2_today: !!level2Result,
                        level3_today: !!level3Result,
                        level1_value: level1Result?.measured_value,
                        level2_value: level2Result?.measured_value,
                        level3_value: level3Result?.measured_value,
                        level1_results: level1Results,
                        level2_results: level2Results,
                        level3_results: level3Results,
                        level1_selected_id: level1SelectedId || level1Result?.id,
                        level2_selected_id: level2SelectedId || level2Result?.id,
                        level3_selected_id: level3SelectedId || level3Result?.id,
                        level1_entry_id: level1Entry?.id,
                        level2_entry_id: level2Entry?.id,
                        level3_entry_id: level3Entry?.id,
                    };
                });

                const qcData = await Promise.all(qcDataPromises);
                setQcParameters(qcData);
            } catch (error) {
                console.error('Error fetching QC status:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchQCStatus();
    }, [testTypes, deviceId]);

    const handleOpenSelection = (parameter: QCParameter) => {
        // Determine which tab to open first (first level with multiple results)
        let initialTab = 0;
        if (parameter.level1_results && parameter.level1_results.length > 1) {
            initialTab = 0;
        } else if (parameter.level2_results && parameter.level2_results.length > 1) {
            initialTab = 1;
        } else if (parameter.level3_results && parameter.level3_results.length > 1) {
            initialTab = 2;
        }

        setActiveTab(initialTab);
        setSelectionDialog({
            open: true,
            testTypeId: parameter.test_type_id,
            level1Results: parameter.level1_results,
            level2Results: parameter.level2_results,
            level3Results: parameter.level3_results,
            level1EntryId: parameter.level1_entry_id,
            level2EntryId: parameter.level2_entry_id,
            level3EntryId: parameter.level3_entry_id,
        });
    };

    const handleSelectResult = (resultId: number, level: number) => {
        if (!selectionDialog.testTypeId) return;

        // Get the QC entry ID for this level
        const qcEntryId = level === 1 ? selectionDialog.level1EntryId :
            level === 2 ? selectionDialog.level2EntryId :
                selectionDialog.level3EntryId;

        if (!qcEntryId) {
            notify('QC entry not found', { type: 'error' });
            return;
        }

        fetch(`/api/v1/quality-control/entries/${qcEntryId}/select-result`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({
                qc_level: level,
                result_id: resultId,
            }),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to update selected result');
                }
                return response.json();
            })
            .then(() => {
                // Update local state
                const updatedParams = qcParameters.map(p => {
                    if (p.test_type_id === selectionDialog.testTypeId) {
                        const results = level === 1 ? selectionDialog.level1Results :
                            level === 2 ? selectionDialog.level2Results :
                                selectionDialog.level3Results;
                        const selectedResult = results?.find(r => r.id === resultId);

                        if (level === 1) {
                            return {
                                ...p,
                                level1_selected_id: resultId,
                                level1_value: selectedResult?.measured_value,
                            };
                        } else if (level === 2) {
                            return {
                                ...p,
                                level2_selected_id: resultId,
                                level2_value: selectedResult?.measured_value,
                            };
                        } else if (level === 3) {
                            return {
                                ...p,
                                level3_selected_id: resultId,
                                level3_value: selectedResult?.measured_value,
                            };
                        }
                    }
                    return p;
                });

                setQcParameters(updatedParams);
                setSelectionDialog({ open: false });
                notify('QC result selected successfully', { type: 'success' });
            })
            .catch(error => {
                notify(`Error: ${error.message}`, { type: 'error' });
            });
    };

    const isLoading = deviceLoading || testTypesLoading || loading;

    if (isLoading) {
        return (
            <MuiBox sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
            </MuiBox>
        );
    }

    if (!device) {
        return (
            <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                <Typography variant="h6" color="text.secondary">
                    Device not found
                </Typography>
            </MuiBox>
        );
    }



    return (
        <MuiBox sx={{ mt: 2 }}>
            <MuiBox sx={{ mb: 3 }}>
                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
                            {device.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" >
                            Device ID: {device.id} | Type: {device.type}
                        </Typography>
                    </CardContent>
                </Card>
            </MuiBox>

            <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                QC Parameters ({qcParameters.length})
            </Typography>

            <Card
                sx={{
                    overflow: 'hidden'
                }}
            >
                <Datagrid
                    data={qcParameters}
                    bulkActionButtons={false}
                    rowClick={(id, resource, record) => {
                        navigate(`/quality-control/${deviceId}/parameter/${record.test_type_id}`);
                        return false;
                    }}
                >
                    <WithRecord label="Test Code" render={(parameter: QCParameter) => (
                        <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                            {parameter.test_type?.code}
                        </Typography>
                    )} />

                    <WithRecord label="Test Name" render={(parameter: QCParameter) => (
                        <Typography variant="body2">
                            {parameter.test_type?.name}
                        </Typography>
                    )} />

                    <WithRecord label="QC Level 1" render={(parameter: QCParameter) => (
                        <MuiBox sx={{ display: 'flex', alignItems: 'center', justifyContent: 'start', gap: 1 }}>
                            {parameter.level1_today ? (
                                <>
                                    <CheckCircleIcon sx={{ color: '#4caf50', fontSize: 24 }} />
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {parameter.level1_value?.toFixed(2)}
                                    </Typography>
                                </>
                            ) : (
                                <CancelIcon sx={{ color: '#f44336', fontSize: 28 }} />
                            )}
                        </MuiBox>
                    )} />

                    <WithRecord label="QC Level 2" render={(parameter: QCParameter) => (
                        <MuiBox sx={{ display: 'flex', alignItems: 'center', justifyContent: 'start', gap: 1 }}>
                            {parameter.level2_today ? (
                                <>
                                    <CheckCircleIcon sx={{ color: '#4caf50', fontSize: 24 }} />
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {parameter.level2_value?.toFixed(2)}
                                    </Typography>
                                </>
                            ) : (
                                <CancelIcon sx={{ color: '#f44336', fontSize: 28 }} />
                            )}
                        </MuiBox>
                    )} />

                    <WithRecord label="QC Level 3" render={(parameter: QCParameter) => (
                        <MuiBox sx={{ display: 'flex', alignItems: 'center', justifyContent: 'start', gap: 1 }}>
                            {parameter.level3_today ? (
                                <>
                                    <CheckCircleIcon sx={{ color: '#4caf50', fontSize: 24 }} />
                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                        {parameter.level3_value?.toFixed(2)}
                                    </Typography>
                                </>
                            ) : (
                                <CancelIcon sx={{ color: '#f44336', fontSize: 28 }} />
                            )}
                        </MuiBox>
                    )} />

                    <WithRecord label="Action" render={(parameter: QCParameter) => {
                        const hasMultiple = (parameter.level1_results && parameter.level1_results.length > 1) ||
                            (parameter.level2_results && parameter.level2_results.length > 1) ||
                            (parameter.level3_results && parameter.level3_results.length > 1);
                        const badgeCount = (parameter.level1_results && parameter.level1_results.length > 1 ? 1 : 0) +
                            (parameter.level2_results && parameter.level2_results.length > 1 ? 1 : 0) +
                            (parameter.level3_results && parameter.level3_results.length > 1 ? 1 : 0);

                        return (
                            <MuiBox onClick={(e) => e.stopPropagation()} sx={{ display: 'flex', justifyContent: 'start' }}>
                                <IconButton
                                    onClick={() => hasMultiple && handleOpenSelection(parameter)}
                                    disabled={!hasMultiple}
                                    size="small"
                                    sx={{
                                        color: hasMultiple ? '#2196f3' : 'rgba(0, 0, 0, 0.26)',
                                        '&:hover': {
                                            backgroundColor: hasMultiple ? 'rgba(33, 150, 243, 0.1)' : 'transparent',
                                        }
                                    }}
                                >
                                    <Badge
                                        badgeContent={badgeCount}
                                        color="error"
                                    >
                                        <HistoryIcon />
                                    </Badge>
                                </IconButton>
                            </MuiBox>
                        );
                    }} />
                </Datagrid>
            </Card>

            {qcParameters.length === 0 && !isLoading && (
                <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                    <Typography variant="body1" color="text.secondary">
                        No QC parameters found for this device
                    </Typography>
                </MuiBox>
            )}

            {/* QC Result Selection Dialog */}
            <Dialog
                open={selectionDialog.open}
                onClose={() => setSelectionDialog({ open: false })}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>
                    Select QC Results
                </DialogTitle>
                <DialogContent>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        Multiple QC results found for today. Please select which result to use for each level:
                    </Typography>

                    <Tabs
                        value={activeTab}
                        onChange={(_, newValue) => setActiveTab(newValue)}
                        sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}
                    >
                        {selectionDialog.level1Results && selectionDialog.level1Results.length > 1 && (
                            <Tab
                                label={
                                    <Badge badgeContent={selectionDialog.level1Results.length} color="primary">
                                        <span style={{ marginRight: 8 }}>Level 1</span>
                                    </Badge>
                                }
                            />
                        )}
                        {selectionDialog.level2Results && selectionDialog.level2Results.length > 1 && (
                            <Tab
                                label={
                                    <Badge badgeContent={selectionDialog.level2Results.length} color="secondary">
                                        <span style={{ marginRight: 8 }}>Level 2</span>
                                    </Badge>
                                }
                            />
                        )}
                        {selectionDialog.level3Results && selectionDialog.level3Results.length > 1 && (
                            <Tab
                                label={
                                    <Badge badgeContent={selectionDialog.level3Results.length} color="warning">
                                        <span style={{ marginRight: 8 }}>Level 3</span>
                                    </Badge>
                                }
                            />
                        )}
                    </Tabs>

                    {/* Level 1 Tab */}
                    {selectionDialog.level1Results && selectionDialog.level1Results.length > 1 &&
                        activeTab === (selectionDialog.level1Results.length > 1 ? 0 : -1) && (
                            <RadioGroup>
                                {selectionDialog.level1Results.map((result, index) => {
                                    const testParam = qcParameters.find(p => p.test_type_id === selectionDialog.testTypeId);
                                    const isSelected = result.id === testParam?.level1_selected_id;
                                    const resultTime = new Date(result.created_at).toLocaleTimeString('en-US', {
                                        hour: '2-digit',
                                        minute: '2-digit',
                                        second: '2-digit'
                                    });
                                    const isRaw = result.method === 'raw';

                                    return (
                                        <Card
                                            key={result.id}
                                            variant="outlined"
                                            sx={{
                                                mb: 2,
                                                border: isSelected ? '2px solid #2196f3' : '1px solid rgba(0, 0, 0, 0.12)',
                                                backgroundColor: isSelected ? 'rgba(33, 150, 243, 0.05)' : 'transparent',
                                            }}
                                        >
                                            <CardContent>
                                                <FormControlLabel
                                                    value={result.id.toString()}
                                                    control={
                                                        <Radio
                                                            checked={isSelected}
                                                            onChange={() => handleSelectResult(result.id, 1)}
                                                        />
                                                    }
                                                    label={
                                                        <MuiBox sx={{ ml: 1 }}>
                                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                                <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                                                                    Result #{index + 1} - {resultTime}
                                                                </Typography>
                                                                <Chip
                                                                    label={isRaw ? 'RAW' : result.method?.toUpperCase()}
                                                                    size="small"
                                                                    sx={{
                                                                        backgroundColor: isRaw ? '#9e9e9e' :
                                                                            result.method === 'manual' ? '#2196f3' : '#4caf50',
                                                                        color: 'white',
                                                                        fontWeight: 600,
                                                                        height: 18,
                                                                        fontSize: '0.7rem',
                                                                    }}
                                                                />
                                                            </MuiBox>
                                                            <MuiBox sx={{ display: 'flex', gap: 3, mt: 1, flexWrap: 'wrap' }}>
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Measured Value
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                        {result.measured_value.toFixed(2)}
                                                                    </Typography>
                                                                </MuiBox>
                                                                {!isRaw && (
                                                                    <>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Result
                                                                            </Typography>
                                                                            <Typography variant="body2">
                                                                                <Chip
                                                                                    label={result.result}
                                                                                    size="small"
                                                                                    sx={{
                                                                                        backgroundColor:
                                                                                            result.result === 'In Control' ? '#4caf50' :
                                                                                                result.result === 'Warning' ? '#ff9800' : '#f44336',
                                                                                        color: 'white',
                                                                                        fontWeight: 600,
                                                                                        height: 20,
                                                                                    }}
                                                                                />
                                                                            </Typography>
                                                                        </MuiBox>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Error SD
                                                                            </Typography>
                                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                                {result.error_sd?.toFixed(2)}
                                                                            </Typography>
                                                                        </MuiBox>
                                                                    </>
                                                                )}
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Created By
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                                                        {result.created_by || result.operator}
                                                                    </Typography>
                                                                </MuiBox>
                                                            </MuiBox>
                                                        </MuiBox>
                                                    }
                                                />
                                            </CardContent>
                                        </Card>
                                    );
                                })}
                            </RadioGroup>
                        )}

                    {/* Level 2 Tab */}
                    {selectionDialog.level2Results && selectionDialog.level2Results.length > 1 &&
                        activeTab === ((selectionDialog.level1Results && selectionDialog.level1Results.length > 1 ? 1 : 0)) && (
                            <RadioGroup>
                                {selectionDialog.level2Results.map((result, index) => {
                                    const testParam = qcParameters.find(p => p.test_type_id === selectionDialog.testTypeId);
                                    const isSelected = result.id === testParam?.level2_selected_id;
                                    const resultTime = new Date(result.created_at).toLocaleTimeString('en-US', {
                                        hour: '2-digit',
                                        minute: '2-digit',
                                        second: '2-digit'
                                    });
                                    const isRaw = result.method === 'raw';

                                    return (
                                        <Card
                                            key={result.id}
                                            variant="outlined"
                                            sx={{
                                                mb: 2,
                                                border: isSelected ? '2px solid #9c27b0' : '1px solid rgba(0, 0, 0, 0.12)',
                                                backgroundColor: isSelected ? 'rgba(156, 39, 176, 0.05)' : 'transparent',
                                            }}
                                        >
                                            <CardContent>
                                                <FormControlLabel
                                                    value={result.id.toString()}
                                                    control={
                                                        <Radio
                                                            checked={isSelected}
                                                            onChange={() => handleSelectResult(result.id, 2)}
                                                        />
                                                    }
                                                    label={
                                                        <MuiBox sx={{ ml: 1 }}>
                                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                                <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                                                                    Result #{index + 1} - {resultTime}
                                                                </Typography>
                                                                <Chip
                                                                    label={isRaw ? 'RAW' : result.method?.toUpperCase()}
                                                                    size="small"
                                                                    sx={{
                                                                        backgroundColor: isRaw ? '#9e9e9e' :
                                                                            result.method === 'manual' ? '#2196f3' : '#4caf50',
                                                                        color: 'white',
                                                                        fontWeight: 600,
                                                                        height: 18,
                                                                        fontSize: '0.7rem',
                                                                    }}
                                                                />
                                                            </MuiBox>
                                                            <MuiBox sx={{ display: 'flex', gap: 3, mt: 1, flexWrap: 'wrap' }}>
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Measured Value
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                        {result.measured_value.toFixed(2)}
                                                                    </Typography>
                                                                </MuiBox>
                                                                {!isRaw && (
                                                                    <>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Result
                                                                            </Typography>
                                                                            <Typography variant="body2">
                                                                                <Chip
                                                                                    label={result.result}
                                                                                    size="small"
                                                                                    sx={{
                                                                                        backgroundColor:
                                                                                            result.result === 'In Control' ? '#4caf50' :
                                                                                                result.result === 'Warning' ? '#ff9800' : '#f44336',
                                                                                        color: 'white',
                                                                                        fontWeight: 600,
                                                                                        height: 20,
                                                                                    }}
                                                                                />
                                                                            </Typography>
                                                                        </MuiBox>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Error SD
                                                                            </Typography>
                                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                                {result.error_sd?.toFixed(2)}
                                                                            </Typography>
                                                                        </MuiBox>
                                                                    </>
                                                                )}
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Created By
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                                                        {result.created_by || result.operator}
                                                                    </Typography>
                                                                </MuiBox>
                                                            </MuiBox>
                                                        </MuiBox>
                                                    }
                                                />
                                            </CardContent>
                                        </Card>
                                    );
                                })}
                            </RadioGroup>
                        )}

                    {/* Level 3 Tab */}
                    {selectionDialog.level3Results && selectionDialog.level3Results.length > 1 &&
                        activeTab === (
                            (selectionDialog.level1Results && selectionDialog.level1Results.length > 1 ? 1 : 0) +
                            (selectionDialog.level2Results && selectionDialog.level2Results.length > 1 ? 1 : 0)
                        ) && (
                            <RadioGroup>
                                {selectionDialog.level3Results.map((result, index) => {
                                    const testParam = qcParameters.find(p => p.test_type_id === selectionDialog.testTypeId);
                                    const isSelected = result.id === testParam?.level3_selected_id;
                                    const resultTime = new Date(result.created_at).toLocaleTimeString('en-US', {
                                        hour: '2-digit',
                                        minute: '2-digit',
                                        second: '2-digit'
                                    });
                                    const isRaw = result.method === 'raw';

                                    return (
                                        <Card
                                            key={result.id}
                                            variant="outlined"
                                            sx={{
                                                mb: 2,
                                                border: isSelected ? '2px solid #ff9800' : '1px solid rgba(0, 0, 0, 0.12)',
                                                backgroundColor: isSelected ? 'rgba(255, 152, 0, 0.05)' : 'transparent',
                                            }}
                                        >
                                            <CardContent>
                                                <FormControlLabel
                                                    value={result.id.toString()}
                                                    control={
                                                        <Radio
                                                            checked={isSelected}
                                                            onChange={() => handleSelectResult(result.id, 3)}
                                                        />
                                                    }
                                                    label={
                                                        <MuiBox sx={{ ml: 1 }}>
                                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                                <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                                                                    Result #{index + 1} - {resultTime}
                                                                </Typography>
                                                                <Chip
                                                                    label={isRaw ? 'RAW' : result.method?.toUpperCase()}
                                                                    size="small"
                                                                    sx={{
                                                                        backgroundColor: isRaw ? '#9e9e9e' :
                                                                            result.method === 'manual' ? '#2196f3' : '#4caf50',
                                                                        color: 'white',
                                                                        fontWeight: 600,
                                                                        height: 18,
                                                                        fontSize: '0.7rem',
                                                                    }}
                                                                />
                                                            </MuiBox>
                                                            <MuiBox sx={{ display: 'flex', gap: 3, mt: 1, flexWrap: 'wrap' }}>
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Measured Value
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                        {result.measured_value.toFixed(2)}
                                                                    </Typography>
                                                                </MuiBox>
                                                                {!isRaw && (
                                                                    <>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Result
                                                                            </Typography>
                                                                            <Typography variant="body2">
                                                                                <Chip
                                                                                    label={result.result}
                                                                                    size="small"
                                                                                    sx={{
                                                                                        backgroundColor:
                                                                                            result.result === 'In Control' ? '#4caf50' :
                                                                                                result.result === 'Warning' ? '#ff9800' : '#f44336',
                                                                                        color: 'white',
                                                                                        fontWeight: 600,
                                                                                        height: 20,
                                                                                    }}
                                                                                />
                                                                            </Typography>
                                                                        </MuiBox>
                                                                        <MuiBox>
                                                                            <Typography variant="caption" color="text.secondary">
                                                                                Error SD
                                                                            </Typography>
                                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                                {result.error_sd?.toFixed(2)}
                                                                            </Typography>
                                                                        </MuiBox>
                                                                    </>
                                                                )}
                                                                <MuiBox>
                                                                    <Typography variant="caption" color="text.secondary">
                                                                        Created By
                                                                    </Typography>
                                                                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                                                        {result.created_by || result.operator}
                                                                    </Typography>
                                                                </MuiBox>
                                                            </MuiBox>
                                                        </MuiBox>
                                                    }
                                                />
                                            </CardContent>
                                        </Card>
                                    );
                                })}
                            </RadioGroup>
                        )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setSelectionDialog({ open: false })}>
                        Close
                    </Button>
                </DialogActions>
            </Dialog>
        </MuiBox>
    );
};
