import { Box as MuiBox, Card, CardContent, Typography, Button, Divider, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Chip, CircularProgress } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { useGetOne } from 'react-admin';
import { useState, useEffect } from 'react';
import ScienceIcon from '@mui/icons-material/Science';
import AddIcon from '@mui/icons-material/Add';
import { Line, XAxis, YAxis, CartesianGrid, ResponsiveContainer, ReferenceLine, Scatter, ComposedChart } from 'recharts';

interface TestType {
    id: number;
    code: string;
    name: string;
    unit: string;
    reference_range_min?: number;
    reference_range_max?: number;
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
    method: string; // 'statistic' or 'manual'
    operator: string;
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

    const { data: testType, isLoading: testTypeLoading } = useGetOne<TestType>('test-type', { id: parseInt(testTypeId || '0') });

    const [qcResults, setQcResults] = useState<QCResult[]>([]);
    const [qcEntries, setQcEntries] = useState<QCEntry[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchQCData();
    }, [deviceId, testTypeId, selectedMethod]);

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

            if (!entriesResponse.ok) {
                //
            } else {
                const entriesData = await entriesResponse.json();
                setQcEntries(entriesData.data || []);
            }

            // Fetch QC results (include method if selected)
            const resultsUrl = new URL(`/api/v1/quality-control/results`, window.location.origin);
            resultsUrl.searchParams.set('device_id', String(deviceId));
            resultsUrl.searchParams.set('test_type_id', String(testTypeId));
            resultsUrl.searchParams.set('method', selectedMethod);

            const resultsResponse = await fetch(resultsUrl.toString(), {
                credentials: 'include',
            });

            if (!resultsResponse.ok) {
                //
            } else {
                const resultsData = await resultsResponse.json();
                setQcResults(resultsData.data || []);
            }
        } catch (error) {
            console.error('Error fetching QC data:', error);
        } finally {
            setLoading(false);
        }
    };

    const formatDate = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: '2-digit' });
    };

    const isLoading = testTypeLoading || loading;

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
                        <Button
                            variant="contained"
                            startIcon={<AddIcon />}
                            onClick={() => navigate(`/quality-control/${deviceId}/parameter/${testTypeId}/entry/new`)}
                        >
                            New QC Entry
                        </Button>
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
            <MuiBox sx={{ display: 'flex', gap: 2, mb: 2, alignItems: 'center' }}>
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
                                const targetMean = levelEntry?.target_mean || 0;
                                const targetSD = levelEntry?.target_sd || 0;

                                levelResults.forEach((result, index) => {
                                    allChartData.push({
                                        level: level,
                                        levelName: `Level ${level}`,
                                        index: index + 1,
                                        date: new Date(result.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
                                        value: result.measured_value,
                                        error: result.error_sd, // Use error_sd calculated by backend
                                        result: result.result,
                                        mean: targetMean,
                                        sd: targetSD,
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
                                    <TableCell><strong>Operator</strong></TableCell>
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
                                            <TableCell>{result.operator}</TableCell>
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
        </MuiBox >
    );
};
