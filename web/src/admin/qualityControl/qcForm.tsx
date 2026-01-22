import { Box as MuiBox, Card, CardContent, Typography, Button, Divider, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Chip, CircularProgress } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { useGetOne } from 'react-admin';
import { useState, useEffect, useRef } from 'react';
import html2canvas from 'html2canvas';
import ScienceIcon from '@mui/icons-material/Science';
import AddIcon from '@mui/icons-material/Add';
import EditNoteIcon from '@mui/icons-material/EditNote';
import DownloadIcon from '@mui/icons-material/Download';
import { Line, XAxis, YAxis, CartesianGrid, ResponsiveContainer, ReferenceLine, Scatter, ComposedChart } from 'recharts';
import { ManualQCInputDialog } from './ManualQCInputDialog';
import { QCReportTemplate } from './QCReportTemplate';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { pdf } from '@react-pdf/renderer';

interface TestType {
    id: number;
    code: string;
    name: string;
    unit: string;
    reference_range_min?: number;
    reference_range_max?: number;
}

interface Device {
    id: number;
    name: string;
    serial_number: string;
}

interface QCEntry {
    id: number;
    device_id: number;
    test_type_id: number;
    qc_level: 1 | 2 | 3;
    lot_number: string;
    target_mean: number;
    target_sd?: number;
    ref_min: number;
    ref_max: number;
    method: string; // 'statistic' or 'manual'
    is_active: boolean;
    created_by: string;
    created_at: string;
    updated_at: string;
    level1_selected_result_id?: number;
    level2_selected_result_id?: number;
    level3_selected_result_id?: number;
}

interface QCResult {
    id: number;
    qc_entry_id: number;
    measured_value: number;
    calculated_mean: number;
    calculated_sd: number;
    calculated_cv: number;
    error_sd: number; // Error in SD units for Levey-Jennings chart (calculated by backend)
    absolute_error: number; // Hasil QC - Target Mean
    relative_error: number; // (ABS Error / Target Mean) × 100%
    sd_1: number;
    sd_2: number;
    sd_3: number;
    result: string;
    method: string; // 'raw', 'statistic' or 'manual'
    operator: string;
    created_by?: string; // New field for audit trail
    message_control_id?: string;
    created_at: string;
    qc_entry?: QCEntry;
    result_count?: number; // To determine if we show calculated SD
}

export const QCForm = () => {
    const { deviceId, testTypeId } = useParams<{ deviceId: string; testTypeId: string }>();
    const navigate = useNavigate();
    const [selectedQcLevelFilter, setSelectedQcLevelFilter] = useState<'all' | 1 | 2 | 3>('all');
    const [visibleLevels, setVisibleLevels] = useState<{ 1: boolean; 2: boolean; 3: boolean }>({ 1: true, 2: true, 3: true });
    const [selectedMethod, setSelectedMethod] = useState<'statistic' | 'manual'>('statistic');
    const [hoveredPoint, setHoveredPoint] = useState<{ level: number; index: number; data: any } | null>(null);
    const [tooltipPos, setTooltipPos] = useState<{ x: number; y: number } | null>(null);
    const [manualInputDialogOpen, setManualInputDialogOpen] = useState(false);
    const [startDate, setStartDate] = useState<Date | null>(null);
    const [endDate, setEndDate] = useState<Date | null>(null);
    const chartRef = useRef<HTMLDivElement>(null);

    const { data: testType, isLoading: testTypeLoading } = useGetOne<TestType>('test-type', { id: parseInt(testTypeId || '0') });
    const { data: device, isLoading: deviceLoading } = useGetOne<Device>('device', { id: parseInt(deviceId || '0') });

    const [qcResults, setQcResults] = useState<QCResult[]>([]);
    const [qcEntries, setQcEntries] = useState<QCEntry[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchQCData();
    }, [deviceId, testTypeId, selectedMethod, startDate, endDate]);

    const fetchQCData = async () => {
        try {
            setLoading(true);

            // Fetch QC entries
            const entriesResponse = await fetch(
                `/api/v1/quality-control/entries?device_id=${deviceId}&test_type_id=${testTypeId}`,
                {
                    credentials: 'include',
                }
            );

            let entries: QCEntry[] = [];
            if (!entriesResponse.ok) {
                //
            } else {
                const entriesData = await entriesResponse.json();
                entries = entriesData.data || [];
                setQcEntries(entries);
            }

            // Fetch QC results (include method if selected)
            const resultsUrl = new URL(`/api/v1/quality-control/results`, window.location.origin);
            resultsUrl.searchParams.set('device_id', String(deviceId));
            resultsUrl.searchParams.set('test_type_id', String(testTypeId));
            resultsUrl.searchParams.set('method', selectedMethod);

            // Add date range filters if set
            if (startDate) {
                resultsUrl.searchParams.set('start_date', startDate.toISOString().split('T')[0]);
            }
            if (endDate) {
                resultsUrl.searchParams.set('end_date', endDate.toISOString().split('T')[0]);
            }

            const resultsResponse = await fetch(resultsUrl.toString(), {
                credentials: 'include',
            });

            if (!resultsResponse.ok) {
                //
            } else {
                const resultsData = await resultsResponse.json();
                const allResults: QCResult[] = resultsData.data || [];

                // Filter to show only selected results for each day
                const filteredResults = filterSelectedResults(allResults, entries);
                setQcResults(filteredResults);
            }
        } catch (error) {
            console.error('Error fetching QC data:', error);
        } finally {
            setLoading(false);
        }
    };

    const filterSelectedResults = (results: QCResult[], entries: QCEntry[]): QCResult[] => {
        // Group results by date and level
        const resultsByDateAndLevel = results.reduce((acc, result) => {
            if (!result.qc_entry) return acc;

            const date = new Date(result.created_at);
            date.setHours(0, 0, 0, 0);
            const dateKey = date.toISOString();
            const level = result.qc_entry.qc_level;
            const key = `${dateKey}-${level}`;

            if (!acc[key]) {
                acc[key] = [];
            }
            acc[key].push(result);
            return acc;
        }, {} as Record<string, QCResult[]>);

        // For each date+level group, select only the selected result
        const selectedResults: QCResult[] = [];

        Object.entries(resultsByDateAndLevel).forEach(([key, dayResults]) => {
            if (dayResults.length === 1) {
                // Only one result, include it
                selectedResults.push(dayResults[0]);
            } else {
                // Multiple results, find the selected one
                const level = dayResults[0].qc_entry?.qc_level;
                const entry = entries.find(e => e.qc_level === level && e.is_active);

                let selectedId: number | undefined;
                if (level === 1) selectedId = entry?.level1_selected_result_id;
                else if (level === 2) selectedId = entry?.level2_selected_result_id;
                else if (level === 3) selectedId = entry?.level3_selected_result_id;

                // If a result is selected, use it; otherwise use the first one
                const selectedResult = selectedId
                    ? dayResults.find(r => r.id === selectedId) || dayResults[0]
                    : dayResults[0];

                selectedResults.push(selectedResult);
            }
        });

        return selectedResults;
    };

    const formatDate = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            hour12: false
        });
    };

    const prepareChartData = () => {
        const allChartData: any[] = [];

        [1, 2, 3].forEach(level => {
            const levelResults = qcResults
                .filter(r => r.method === selectedMethod && r.qc_entry?.qc_level === level)
                .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()); // Descending: newest first

            if (levelResults.length === 0) return;

            const levelEntry = activeEntries.find(e => e.qc_level === level);

            levelResults.forEach((result, index) => {
                // Use calculated mean/sd for statistic method, target mean/sd for manual method
                const useMean = selectedMethod === 'statistic' && result.calculated_mean
                    ? result.calculated_mean
                    : (levelEntry?.target_mean || 0);
                const useSD = selectedMethod === 'statistic' && result.calculated_sd
                    ? result.calculated_sd
                    : (levelEntry?.target_sd || 0);

                // Determine status based on error_sd
                const absErrorSD = Math.abs(result.error_sd);
                let status = 'In Control';
                if (absErrorSD >= 3) status = 'Reject';
                else if (absErrorSD >= 2) status = 'Warning';

                // Determine source
                const source = result.message_control_id ? 'Analyzer' : 'Manual';

                allChartData.push({
                    level: level,
                    levelName: `Level ${level}`,
                    index: index + 1,
                    date: formatDate(result.created_at),
                    errorSD: result.error_sd,
                    result: result.measured_value,
                    mean: useMean,
                    sd: useSD,
                    lotNumber: levelEntry?.lot_number || '-',
                    absoluteError: result.absolute_error,
                    relativeError: result.relative_error,
                    status: status,
                    source: source,
                    // For debugging: include both calculated and target values
                    _debug: {
                        calculated_mean: result.calculated_mean,
                        calculated_sd: result.calculated_sd,
                        target_mean: levelEntry?.target_mean,
                        target_sd: levelEntry?.target_sd,
                        method: selectedMethod,
                    }
                });
            });
        });

        console.log('Chart data for PDF:', allChartData);
        return allChartData;
    };

    const handleDownloadPDF = async () => {
        if (!testType) return;

        try {
            const chartData = prepareChartData();

            // Capture chart as image
            let chartImageUrl = '';
            if (chartRef.current && chartData.length > 0) {
                const canvas = await html2canvas(chartRef.current, {
                    backgroundColor: '#ffffff',
                    scale: 2, // Higher quality
                    logging: false,
                });
                chartImageUrl = canvas.toDataURL('image/png');
            }

            const blob = await pdf(
                <QCReportTemplate
                    testType={testType}
                    device={device}
                    entries={activeEntries}
                    qcResults={qcResults.filter(r => r.method === selectedMethod)}
                    selectedMethod={selectedMethod}
                    startDate={startDate}
                    endDate={endDate}
                    chartData={chartData}
                    chartImageUrl={chartImageUrl}
                />
            ).toBlob();

            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = `QC_Report_${testType.code}_${new Date().toISOString().split('T')[0]}.pdf`;
            link.click();
            URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error generating PDF:', error);
        }
    };

    const isLoading = testTypeLoading || deviceLoading || loading;

    if (isLoading) {
        return (
            <MuiBox sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <CircularProgress />
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

    const activeEntries = qcEntries.filter(e => e.is_active);
    const hasActiveEntry = activeEntries.length > 0;

    return (
        <MuiBox sx={{ mt: 2 }}>
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <MuiBox sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                        <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <ScienceIcon sx={{ fontSize: 40, color: '#2196f3' }} />
                            <MuiBox>
                                <Typography variant="h5" sx={{ fontWeight: 600 }}>
                                    {testType.code} - {testType.name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    Unit: {testType.unit}
                                </Typography>
                            </MuiBox>
                        </MuiBox>
                        <MuiBox sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                            <Button
                                variant="contained"
                                startIcon={<AddIcon />}
                                onClick={() => navigate(`/quality-control/${deviceId}/parameter/${testTypeId}/entry/new`)}
                            >
                                New QC Entry
                            </Button>
                            <Button
                                variant="outlined"
                                startIcon={<EditNoteIcon />}
                                onClick={() => setManualInputDialogOpen(true)}
                                disabled={!hasActiveEntry}
                            >
                                Manual Input
                            </Button>
                            <Button
                                variant="outlined"
                                startIcon={<DownloadIcon />}
                                onClick={handleDownloadPDF}
                                disabled={!hasActiveEntry || qcResults.length === 0}
                                color="primary"
                            >
                                Download PDF
                            </Button>
                        </MuiBox>
                    </MuiBox>
                </CardContent>
            </Card>

            {hasActiveEntry && (
                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                            Active QC Entries
                        </Typography>
                        <MuiBox sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                            {activeEntries.map(entry => {
                                const entryResults = qcResults.filter(r => r.qc_entry_id === entry.id);
                                const resultCount = entryResults.length;
                                const hasCalculated = resultCount >= 5;

                                let calculatedMean = 0;
                                let calculatedSD = 0;
                                let calculatedCV = 0;

                                if (hasCalculated && entryResults.length > 0) {
                                    const latestResult = entryResults[0];
                                    calculatedMean = latestResult.calculated_mean;
                                    calculatedSD = latestResult.calculated_sd;
                                    calculatedCV = latestResult.calculated_cv;
                                }

                                const quotedCV = entry.target_sd && entry.target_mean
                                    ? (entry.target_sd / entry.target_mean) * 100
                                    : 0;

                                return (
                                    <Card key={entry.id} variant="outlined" sx={{ flex: '1 1 400px' }}>
                                        <CardContent>
                                            <MuiBox sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', mb: 2 }}>
                                                <Chip
                                                    label={`Level ${entry.qc_level}`}
                                                    size="small"
                                                    sx={{
                                                        backgroundColor:
                                                            entry.qc_level === 1 ? '#2196f3' :
                                                                entry.qc_level === 2 ? '#9c27b0' : '#ff9800',
                                                        color: 'white',
                                                        fontWeight: 600
                                                    }}
                                                />
                                                <Chip
                                                    label="Active"
                                                    size="small"
                                                    sx={{
                                                        backgroundColor: 'rgba(76, 175, 80, 0.1)',
                                                        color: '#4caf50',
                                                        fontWeight: 600
                                                    }}
                                                />
                                            </MuiBox>

                                            <Typography variant="body2" sx={{ mb: 2 }}>
                                                <strong>Lot:</strong> {entry.lot_number} | <strong>Count:</strong> {resultCount} {"|"} <strong>Ref:</strong> {entry.ref_min?.toFixed ? entry.ref_min.toFixed(2) : entry.ref_min} - {entry.ref_max?.toFixed ? entry.ref_max.toFixed(2) : entry.ref_max}
                                            </Typography>

                                            <MuiBox sx={{ mb: 2, p: 1.5, backgroundColor: 'rgba(33, 150, 243, 0.05)', borderRadius: 1 }}>
                                                <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1, color: '#2196f3' }}>
                                                    Quoted (Datasheet)
                                                </Typography>
                                                <MuiBox sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 1 }}>
                                                    <MuiBox>
                                                        <Typography variant="caption" color="text.secondary">Mean</Typography>
                                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                            {entry.target_mean.toFixed(2)}
                                                        </Typography>
                                                    </MuiBox>
                                                    <MuiBox>
                                                        <Typography variant="caption" color="text.secondary">SD</Typography>
                                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                            {entry.target_sd?.toFixed(2) || '-'}
                                                        </Typography>
                                                    </MuiBox>
                                                    <MuiBox>
                                                        <Typography variant="caption" color="text.secondary">CV (%)</Typography>
                                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                            {quotedCV > 0 ? quotedCV.toFixed(2) : '-'}
                                                        </Typography>
                                                    </MuiBox>
                                                </MuiBox>
                                            </MuiBox>

                                            <MuiBox sx={{ p: 1.5, backgroundColor: 'rgba(76, 175, 80, 0.05)', borderRadius: 1 }}>
                                                <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1, color: '#4caf50' }}>
                                                    Calculated (System)
                                                </Typography>
                                                {hasCalculated ? (
                                                    <MuiBox sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 1 }}>
                                                        <MuiBox>
                                                            <Typography variant="caption" color="text.secondary">Mean</Typography>
                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                {calculatedMean.toFixed(2)}
                                                            </Typography>
                                                        </MuiBox>
                                                        <MuiBox>
                                                            <Typography variant="caption" color="text.secondary">SD</Typography>
                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                {calculatedSD.toFixed(2)}
                                                            </Typography>
                                                        </MuiBox>
                                                        <MuiBox>
                                                            <Typography variant="caption" color="text.secondary">CV (%)</Typography>
                                                            <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                                                {calculatedCV.toFixed(2)}
                                                            </Typography>
                                                        </MuiBox>
                                                    </MuiBox>
                                                ) : (
                                                    <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                                                        Waiting for {5 - resultCount} more measurements...
                                                    </Typography>
                                                )}
                                            </MuiBox>
                                        </CardContent>
                                    </Card>
                                );
                            })}
                        </MuiBox>
                    </CardContent>
                </Card>
            )}
            {/* Filters: Level + Method */}
            <MuiBox sx={{ display: 'flex', gap: 2, mb: 2, alignItems: 'center', flexWrap: 'wrap' }}>
                <MuiBox sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                    <Typography variant="body2" sx={{ mr: 1, fontWeight: 500 }}>Method:</Typography>
                    <Chip
                        label="Statistic"
                        onClick={() => setSelectedMethod('statistic')}
                        size="small"
                        sx={{
                            backgroundColor: selectedMethod === 'statistic' ? '#2196f3' : 'rgba(33, 150, 243, 0.1)',
                            color: selectedMethod === 'statistic' ? 'white' : '#2196f3',
                            fontWeight: 500,
                            cursor: 'pointer'
                        }}
                    />
                    <Chip
                        label="Manual"
                        onClick={() => setSelectedMethod('manual')}
                        size="small"
                        sx={{
                            backgroundColor: selectedMethod === 'manual' ? '#9c27b0' : 'rgba(156, 39, 176, 0.1)',
                            color: selectedMethod === 'manual' ? 'white' : '#9c27b0',
                            fontWeight: 500,
                            cursor: 'pointer'
                        }}
                    />
                </MuiBox>

                <Divider orientation="vertical" flexItem />

                <MuiBox sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>Date Range:</Typography>
                    <LocalizationProvider dateAdapter={AdapterDateFns}>
                        <DatePicker
                            label="Start Date"
                            value={startDate}
                            onChange={(newValue) => setStartDate(newValue)}
                            slotProps={{
                                textField: {
                                    size: 'small',
                                    sx: { width: 160 }
                                }
                            }}
                        />
                        <DatePicker
                            label="End Date"
                            value={endDate}
                            onChange={(newValue) => setEndDate(newValue)}
                            minDate={startDate || undefined}
                            slotProps={{
                                textField: {
                                    size: 'small',
                                    sx: { width: 160 }
                                }
                            }}
                        />
                    </LocalizationProvider>
                    {(startDate || endDate) && (
                        <Button
                            size="small"
                            onClick={() => {
                                setStartDate(null);
                                setEndDate(null);
                            }}
                            sx={{ textTransform: 'none' }}
                        >
                            Clear
                        </Button>
                    )}
                </MuiBox>
            </MuiBox>

            {!hasActiveEntry && (
                <Card sx={{ mb: 3, backgroundColor: 'rgba(255, 152, 0, 0.05)' }}>
                    <CardContent>
                        <Typography variant="body1" sx={{ mb: 2 }}>
                            No active QC entry found. Please create a QC entry first before running QC measurements.
                        </Typography>
                        <Button
                            variant="contained"
                            startIcon={<AddIcon />}
                            onClick={() => navigate(`/quality-control/${deviceId}/parameter/${testTypeId}/entry/new`)}
                        >
                            Create QC Entry
                        </Button>
                    </CardContent>
                </Card>
            )}

            {/* Levey-Jennings Charts */}
            {qcResults.length > 0 && (
                <Card sx={{ mb: 3 }}>
                    <CardContent>
                        <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                            Levey-Jennings Chart
                        </Typography>

                        {(() => {
                            const allChartData: any[] = [];

                            [1, 2, 3].forEach(level => {
                                const levelResults = qcResults
                                    .filter(r => r.method === selectedMethod && r.qc_entry?.qc_level === level)
                                    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());

                                if (levelResults.length === 0) return;

                                const levelEntry = activeEntries.find(e => e.qc_level === level);

                                levelResults.forEach((result, index) => {
                                    // Use calculated mean/sd for statistic method, target mean/sd for manual method
                                    const useMean = selectedMethod === 'statistic' && result.calculated_mean
                                        ? result.calculated_mean
                                        : (levelEntry?.target_mean || 0);
                                    const useSD = selectedMethod === 'statistic' && result.calculated_sd
                                        ? result.calculated_sd
                                        : (levelEntry?.target_sd || 0);

                                    allChartData.push({
                                        level: level,
                                        levelName: `Level ${level}`,
                                        index: index + 1,
                                        date: new Date(result.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
                                        value: result.measured_value,
                                        error: result.error_sd, // Use error_sd calculated by backend
                                        result: result.result,
                                        mean: useMean,
                                        sd: useSD,
                                        lotNumber: levelEntry?.lot_number || '-',
                                    });
                                });
                            });

                            if (allChartData.length === 0) return null;

                            const level1Data = allChartData.filter(d => d.level === 1);
                            const level2Data = allChartData.filter(d => d.level === 2);
                            const level3Data = allChartData.filter(d => d.level === 3);

                            return (
                                <MuiBox sx={{ mb: 4 }}>
                                    <MuiBox sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 2, flexWrap: 'wrap' }}>
                                        {level1Data.length > 0 && (
                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                <Chip
                                                    label="Level 1"
                                                    size="small"
                                                    onClick={() => setVisibleLevels(prev => ({ ...prev, 1: !prev[1] }))}
                                                    sx={{
                                                        backgroundColor: visibleLevels[1] ? '#2196f3' : 'rgba(33, 150, 243, 0.3)',
                                                        color: 'white',
                                                        fontWeight: 600,
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            backgroundColor: visibleLevels[1] ? '#1976d2' : 'rgba(33, 150, 243, 0.5)',
                                                        }
                                                    }}
                                                />
                                                <Typography variant="caption" color="text.secondary">
                                                    Lot: {level1Data[0].lotNumber} | {level1Data[0].mean.toFixed(2)} ± {level1Data[0].sd.toFixed(2)}
                                                </Typography>
                                            </MuiBox>
                                        )}
                                        {level2Data.length > 0 && (
                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                <Chip
                                                    label="Level 2"
                                                    size="small"
                                                    onClick={() => setVisibleLevels(prev => ({ ...prev, 2: !prev[2] }))}
                                                    sx={{
                                                        backgroundColor: visibleLevels[2] ? '#9c27b0' : 'rgba(156, 39, 176, 0.3)',
                                                        color: 'white',
                                                        fontWeight: 600,
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            backgroundColor: visibleLevels[2] ? '#7b1fa2' : 'rgba(156, 39, 176, 0.5)',
                                                        }
                                                    }}
                                                />
                                                <Typography variant="caption" color="text.secondary">
                                                    Lot: {level2Data[0].lotNumber} | {level2Data[0].mean.toFixed(2)} ± {level2Data[0].sd.toFixed(2)}
                                                </Typography>
                                            </MuiBox>
                                        )}
                                        {level3Data.length > 0 && (
                                            <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                <Chip
                                                    label="Level 3"
                                                    size="small"
                                                    onClick={() => setVisibleLevels(prev => ({ ...prev, 3: !prev[3] }))}
                                                    sx={{
                                                        backgroundColor: visibleLevels[3] ? '#ff9800' : 'rgba(255, 152, 0, 0.3)',
                                                        color: 'white',
                                                        fontWeight: 600,
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            backgroundColor: visibleLevels[3] ? '#f57c00' : 'rgba(255, 152, 0, 0.5)',
                                                        }
                                                    }}
                                                />
                                                <Typography variant="caption" color="text.secondary">
                                                    Lot: {level3Data[0].lotNumber} | {level3Data[0].mean.toFixed(2)} ± {level3Data[0].sd.toFixed(2)}
                                                </Typography>
                                            </MuiBox>
                                        )}
                                    </MuiBox>

                                    <MuiBox ref={chartRef}>
                                        <ResponsiveContainer width="100%" height={400}>
                                            <ComposedChart
                                                margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                                                onMouseLeave={() => {
                                                    setHoveredPoint(null);
                                                    setTooltipPos(null);
                                                }}
                                            >
                                                <CartesianGrid strokeDasharray="3 3" />
                                                <XAxis
                                                    dataKey="index"
                                                    label={{ value: 'N (Measurement)', position: 'insideBottom', offset: -5 }}
                                                    tick={{ fontSize: 12 }}
                                                    type="number"
                                                    domain={[1, 'dataMax']}
                                                />
                                                <YAxis
                                                    label={{ value: 'Error (SD)', angle: -90, position: 'insideLeft' }}
                                                    domain={[-4, 4]}
                                                    ticks={[-3, -2, -1, 0, 1, 2, 3]}
                                                />

                                                <ReferenceLine y={0} stroke="#000" strokeWidth={2} />
                                                <ReferenceLine y={1} stroke="#FFA500" strokeDasharray="5 5" />
                                                <ReferenceLine y={-1} stroke="#FFA500" strokeDasharray="5 5" />
                                                <ReferenceLine y={2} stroke="#FF6347" strokeDasharray="3 3" />
                                                <ReferenceLine y={-2} stroke="#FF6347" strokeDasharray="3 3" />
                                                <ReferenceLine y={3} stroke="#DC143C" strokeWidth={2} />
                                                <ReferenceLine y={-3} stroke="#DC143C" strokeWidth={2} />

                                                {level1Data.length > 0 && visibleLevels[1] && (
                                                    <>
                                                        <Line
                                                            type="monotone"
                                                            data={level1Data}
                                                            dataKey="error"
                                                            stroke="#2196f3"
                                                            strokeWidth={2}
                                                            dot={false}
                                                            activeDot={false}
                                                            connectNulls
                                                        />
                                                        <Scatter
                                                            id="level1-scatter"
                                                            data={level1Data}
                                                            dataKey="error"
                                                            fill="#2196f3"
                                                            isAnimationActive={false}
                                                            shape={(props: any) => {
                                                                const { cx, cy, payload, index } = props;
                                                                const errorValue = payload.error;
                                                                const isHovered = hoveredPoint?.level === 1 && hoveredPoint?.index === index;
                                                                let color;
                                                                if (Math.abs(errorValue) <= 2) {
                                                                    color = '#4caf50';
                                                                } else if (Math.abs(errorValue) <= 3) {
                                                                    color = '#ff9800';
                                                                } else {
                                                                    color = '#f44336';
                                                                }
                                                                return (
                                                                    <circle
                                                                        cx={cx}
                                                                        cy={cy}
                                                                        r={isHovered ? 7 : 5}
                                                                        fill={color}
                                                                        stroke={color}
                                                                        strokeWidth={isHovered ? 2 : 1}
                                                                        style={{ cursor: 'pointer', }}
                                                                        onMouseOver={(e: any) => {
                                                                            setHoveredPoint({ level: 1, index, data: payload });
                                                                            setTooltipPos({ x: e.clientX, y: e.clientY });
                                                                        }}
                                                                        onMouseOut={() => {
                                                                            setHoveredPoint(null);
                                                                            setTooltipPos(null);
                                                                        }}
                                                                    />
                                                                );
                                                            }}
                                                        />

                                                    </>
                                                )}

                                                {level2Data.length > 0 && visibleLevels[2] && (
                                                    <>
                                                        <Line
                                                            type="monotone"
                                                            data={level2Data}
                                                            dataKey="error"
                                                            stroke="#9c27b0"
                                                            strokeWidth={2}
                                                            dot={false}
                                                            activeDot={false}
                                                            connectNulls
                                                        />
                                                        <Scatter
                                                            id="level2-scatter"
                                                            data={level2Data}
                                                            dataKey="error"
                                                            fill="#9c27b0"
                                                            isAnimationActive={false}
                                                            shape={(props: any) => {
                                                                const { cx, cy, payload, index } = props;
                                                                const errorValue = payload.error;
                                                                const isHovered = hoveredPoint?.level === 2 && hoveredPoint?.index === index;
                                                                let color;
                                                                if (Math.abs(errorValue) <= 2) {
                                                                    color = '#4caf50';
                                                                } else if (Math.abs(errorValue) <= 3) {
                                                                    color = '#ff9800';
                                                                } else {
                                                                    color = '#f44336';
                                                                }
                                                                return (
                                                                    <circle
                                                                        cx={cx}
                                                                        cy={cy}
                                                                        r={isHovered ? 7 : 5}
                                                                        fill={color}
                                                                        stroke={color}
                                                                        strokeWidth={isHovered ? 2 : 1}
                                                                        style={{ cursor: 'pointer' }}
                                                                        onMouseOver={(e: any) => {
                                                                            setHoveredPoint({ level: 2, index, data: payload });
                                                                            setTooltipPos({ x: e.clientX, y: e.clientY });
                                                                        }}
                                                                        onMouseOut={() => {
                                                                            setHoveredPoint(null);
                                                                            setTooltipPos(null);
                                                                        }}
                                                                    />
                                                                );
                                                            }}
                                                        />

                                                    </>
                                                )}

                                                {level3Data.length > 0 && visibleLevels[3] && (
                                                    <>
                                                        <Scatter
                                                            id="level3-scatter"
                                                            data={level3Data}
                                                            dataKey="error"
                                                            fill="#ff9800"
                                                            isAnimationActive={false}
                                                            shape={(props: any) => {
                                                                const { cx, cy, payload, index } = props;
                                                                const errorValue = payload.error;
                                                                const isHovered = hoveredPoint?.level === 3 && hoveredPoint?.index === index;
                                                                let color;
                                                                if (Math.abs(errorValue) <= 2) {
                                                                    color = '#4caf50';
                                                                } else if (Math.abs(errorValue) <= 3) {
                                                                    color = '#ff9800';
                                                                } else {
                                                                    color = '#f44336';
                                                                }
                                                                return (
                                                                    <circle
                                                                        cx={cx}
                                                                        cy={cy}
                                                                        r={isHovered ? 7 : 5}
                                                                        fill={color}
                                                                        stroke={color}
                                                                        strokeWidth={isHovered ? 2 : 1}
                                                                        style={{ cursor: 'pointer' }}
                                                                        onMouseOver={(e: any) => {
                                                                            setHoveredPoint({ level: 3, index, data: payload });
                                                                            setTooltipPos({ x: e.clientX, y: e.clientY });
                                                                        }}
                                                                        onMouseOut={() => {
                                                                            setHoveredPoint(null);
                                                                            setTooltipPos(null);
                                                                        }}
                                                                    />
                                                                );
                                                            }}
                                                        />
                                                        <Line
                                                            type="monotone"
                                                            data={level3Data}
                                                            dataKey="error"
                                                            stroke="#ff9800"
                                                            strokeWidth={2}
                                                            dot={false}
                                                            activeDot={false}
                                                            connectNulls
                                                        />
                                                    </>
                                                )}
                                            </ComposedChart>
                                        </ResponsiveContainer>
                                    </MuiBox>

                                    {hoveredPoint && tooltipPos && (
                                        <Paper
                                            sx={{
                                                position: 'fixed',
                                                left: tooltipPos.x + 10,
                                                top: tooltipPos.y + 10,
                                                p: 1.5,
                                                border: '1px solid #ccc',
                                                zIndex: 9999,
                                                pointerEvents: 'none'
                                            }}
                                        >
                                            <Typography variant="caption" display="block">
                                                <strong>{hoveredPoint.data.levelName}</strong> - {hoveredPoint.data.lotNumber}
                                            </Typography>
                                            <Typography variant="caption" display="block">
                                                <strong>Date:</strong> {hoveredPoint.data.date}
                                            </Typography>
                                            <Typography variant="caption" display="block">
                                                <strong>Measurement #{hoveredPoint.data.index}:</strong> {hoveredPoint.data.value.toFixed(2)} {testType.unit}
                                            </Typography>
                                            <Typography variant="caption" display="block">
                                                <strong>Error:</strong> {hoveredPoint.data.error.toFixed(2)} SD
                                            </Typography>
                                            <Typography variant="caption" display="block" sx={{
                                                color:
                                                    hoveredPoint.data.result === 'In Control' ? '#4caf50' :
                                                        hoveredPoint.data.result === 'Warning' ? '#ff9800' : '#f44336',
                                                fontWeight: 600
                                            }}>
                                                {hoveredPoint.data.result}
                                            </Typography>
                                        </Paper>
                                    )}

                                    <MuiBox sx={{ display: 'flex', gap: 5, justifyContent: 'center', mt: 5, mb: 2 }}>
                                        <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                            <MuiBox sx={{
                                                width: 18,
                                                height: 18,
                                                borderRadius: '50%',
                                                backgroundColor: '#4caf50',
                                                border: '2px solid #fff',
                                                boxShadow: '0 0 0 1px #ccc'
                                            }} />
                                            <Typography variant="body2">
                                                <strong>In Control</strong> (±2 SD)
                                            </Typography>
                                        </MuiBox>
                                        <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                            <MuiBox sx={{
                                                width: 18,
                                                height: 18,
                                                borderRadius: '50%',
                                                backgroundColor: '#ff9800',
                                                border: '2px solid #fff',
                                                boxShadow: '0 0 0 1px #ccc'
                                            }} />
                                            <Typography variant="body2">
                                                <strong>Warning</strong> (2-3 SD)
                                            </Typography>
                                        </MuiBox>
                                        <MuiBox sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                                            <MuiBox sx={{
                                                width: 18,
                                                height: 18,
                                                borderRadius: '50%',
                                                backgroundColor: '#f44336',
                                                border: '2px solid #fff',
                                                boxShadow: '0 0 0 1px #ccc'
                                            }} />
                                            <Typography variant="body2">
                                                <strong>Reject</strong> (&gt;3 SD)
                                            </Typography>
                                        </MuiBox>
                                    </MuiBox>
                                </MuiBox>
                            );
                        })()}
                    </CardContent>
                </Card>
            )}

            {/* QC History */}
            <Card>
                <CardContent>
                    <MuiBox sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>
                            QC History
                        </Typography>

                        <MuiBox sx={{ display: 'flex', gap: 2 }}>
                            <MuiBox sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                                <Typography variant="body2" sx={{ mr: 1, fontWeight: 500 }}>Level:</Typography>
                                <Chip
                                    label="All"
                                    onClick={() => setSelectedQcLevelFilter('all')}
                                    size="small"
                                    sx={{
                                        backgroundColor: selectedQcLevelFilter === 'all' ? '#555' : 'rgba(0, 0, 0, 0.08)',
                                        color: selectedQcLevelFilter === 'all' ? 'white' : 'text.primary',
                                        fontWeight: 500,
                                        cursor: 'pointer',
                                        '&:hover': {
                                            backgroundColor: selectedQcLevelFilter === 'all' ? '#333' : 'rgba(0, 0, 0, 0.12)',
                                        }
                                    }}
                                />
                                <Chip
                                    label="Level 1"
                                    onClick={() => setSelectedQcLevelFilter(1)}
                                    size="small"
                                    sx={{
                                        backgroundColor: selectedQcLevelFilter === 1 ? '#2196f3' : 'rgba(33, 150, 243, 0.1)',
                                        color: selectedQcLevelFilter === 1 ? 'white' : '#2196f3',
                                        fontWeight: 500,
                                        cursor: 'pointer',
                                        '&:hover': {
                                            backgroundColor: selectedQcLevelFilter === 1 ? '#1976d2' : 'rgba(33, 150, 243, 0.2)',
                                        }
                                    }}
                                />
                                <Chip
                                    label="Level 2"
                                    onClick={() => setSelectedQcLevelFilter(2)}
                                    size="small"
                                    sx={{
                                        backgroundColor: selectedQcLevelFilter === 2 ? '#9c27b0' : 'rgba(156, 39, 176, 0.1)',
                                        color: selectedQcLevelFilter === 2 ? 'white' : '#9c27b0',
                                        fontWeight: 500,
                                        cursor: 'pointer',
                                        '&:hover': {
                                            backgroundColor: selectedQcLevelFilter === 2 ? '#7b1fa2' : 'rgba(156, 39, 176, 0.2)',
                                        }
                                    }}
                                />
                                <Chip
                                    label="Level 3"
                                    onClick={() => setSelectedQcLevelFilter(3)}
                                    size="small"
                                    sx={{
                                        backgroundColor: selectedQcLevelFilter === 3 ? '#ff9800' : 'rgba(255, 152, 0, 0.1)',
                                        color: selectedQcLevelFilter === 3 ? 'white' : '#ff9800',
                                        fontWeight: 500,
                                        cursor: 'pointer',
                                        '&:hover': {
                                            backgroundColor: selectedQcLevelFilter === 3 ? '#f57c00' : 'rgba(255, 152, 0, 0.2)',
                                        }
                                    }}
                                />
                            </MuiBox>
                        </MuiBox>
                    </MuiBox>

                    <Divider sx={{ mb: 2 }} />

                    <TableContainer component={Paper} variant="outlined">
                        <Table>
                            <TableHead>
                                <TableRow sx={{ backgroundColor: 'rgba(0, 0, 0, 0.08)' }}>
                                    <TableCell><strong>Date</strong></TableCell>
                                    <TableCell><strong>Level</strong></TableCell>
                                    <TableCell><strong>Lot Number</strong></TableCell>
                                    <TableCell align="center"><strong>Measurement</strong></TableCell>
                                    <TableCell align="center"><strong>Absolute Error</strong></TableCell>
                                    <TableCell align="center"><strong>Relative Error (%)</strong></TableCell>
                                    <TableCell align="center"><strong>Result</strong></TableCell>
                                    <TableCell><strong>Source</strong></TableCell>
                                    <TableCell><strong>Created By</strong></TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {qcResults
                                    .filter((result) => {
                                        if (result.method !== selectedMethod) return false;
                                        if (selectedQcLevelFilter !== 'all' && result.qc_entry?.qc_level !== selectedQcLevelFilter) {
                                            return false;
                                        }
                                        return true;
                                    })
                                    .map((result, index) => (
                                        <TableRow
                                            key={result.id}
                                            sx={{
                                                '&:hover': { backgroundColor: 'rgba(0, 0, 0, 0.04)' },
                                                backgroundColor:
                                                    result.result === 'Error' ? 'rgba(244, 67, 54, 0.08)' :
                                                        result.result === 'Warning' ? 'rgba(255, 152, 0, 0.08)' :
                                                            index % 2 === 0 ? 'rgba(0, 0, 0, 0.02)' : 'white'
                                            }}
                                        >
                                            <TableCell>{formatDate(result.created_at)}</TableCell>
                                            <TableCell>
                                                <Chip
                                                    label={`Level ${result.qc_entry?.qc_level || '-'}`}
                                                    size="small"
                                                    sx={{
                                                        backgroundColor:
                                                            result.qc_entry?.qc_level === 1 ? 'rgba(33, 150, 243, 0.2)' :
                                                                result.qc_entry?.qc_level === 2 ? 'rgba(156, 39, 176, 0.2)' :
                                                                    'rgba(255, 152, 0, 0.2)',
                                                        color:
                                                            result.qc_entry?.qc_level === 1 ? '#2196f3' :
                                                                result.qc_entry?.qc_level === 2 ? '#9c27b0' :
                                                                    '#ff9800',
                                                        fontWeight: 600
                                                    }}
                                                />
                                            </TableCell>
                                            <TableCell>{result.qc_entry?.lot_number || '-'}</TableCell>
                                            <TableCell align="center">
                                                {result.measured_value.toFixed(2)}
                                            </TableCell>
                                            <TableCell align="center">
                                                {result.absolute_error !== undefined && result.absolute_error !== null ? (
                                                    <Typography variant="body2" sx={{
                                                        color: Math.abs(result.absolute_error) > 2 ? '#f44336' : 'text.primary',
                                                        fontWeight: Math.abs(result.absolute_error) > 2 ? 600 : 400
                                                    }}>
                                                        {result.absolute_error > 0 ? '+' : ''}{result.absolute_error.toFixed(2)}
                                                    </Typography>
                                                ) : (
                                                    <Typography variant="body2" color="text.secondary">-</Typography>
                                                )}
                                            </TableCell>
                                            <TableCell align="center">
                                                {result.relative_error !== undefined && result.relative_error !== null ? (
                                                    <Typography variant="body2" sx={{
                                                        color: result.relative_error > 10 ? '#f44336' : 'text.primary',
                                                        fontWeight: result.relative_error > 10 ? 600 : 400
                                                    }}>
                                                        {result.relative_error.toFixed(2)}%
                                                    </Typography>
                                                ) : (
                                                    <Typography variant="body2" color="text.secondary">-</Typography>
                                                )}
                                            </TableCell>
                                            <TableCell align="center">
                                                <Chip
                                                    label={result.result}
                                                    size="small"
                                                    sx={{
                                                        backgroundColor:
                                                            result.result === 'In Control' ? '#4caf50' :
                                                                result.result === 'Warning' ? '#ff9800' : '#f44336',
                                                        color: 'white',
                                                        fontWeight: 600
                                                    }}
                                                />
                                            </TableCell>
                                            <TableCell>
                                                <Chip
                                                    label={result.message_control_id ? 'Analyzer' : 'Manual'}
                                                    size="small"
                                                    sx={{
                                                        backgroundColor: result.message_control_id ? 'rgba(33, 150, 243, 0.1)' : 'rgba(156, 39, 176, 0.1)',
                                                        color: result.message_control_id ? '#2196f3' : '#9c27b0',
                                                        fontWeight: 500
                                                    }}
                                                />
                                            </TableCell>
                                            <TableCell>{result.created_by || result.operator}</TableCell>
                                        </TableRow>
                                    ))}
                            </TableBody>
                        </Table>
                    </TableContainer>

                    {qcResults.length === 0 && (
                        <MuiBox sx={{ textAlign: 'center', p: 3 }}>
                            <Typography variant="body1" color="text.secondary">
                                No QC results found. QC results will appear here automatically when measurements are received from the analyzer.
                            </Typography>
                        </MuiBox>
                    )}
                </CardContent>
            </Card>

            <ManualQCInputDialog
                open={manualInputDialogOpen}
                onClose={() => setManualInputDialogOpen(false)}
                deviceId={deviceId || ''}
                testTypeId={testTypeId || ''}
                activeEntries={activeEntries}
                onSuccess={fetchQCData}
            />
        </MuiBox >
    );
};
