import { Card, CardContent, Typography, Box, Button } from '@mui/material';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import HourglassEmptyIcon from '@mui/icons-material/HourglassEmpty';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import ScienceIcon from '@mui/icons-material/Science';
import DevicesIcon from '@mui/icons-material/Lan';
import PeopleIcon from '@mui/icons-material/People';
import ListAltIcon from '@mui/icons-material/ListAlt';
import DashboardOutlinedIcon from "@mui/icons-material/Summarize";
import ShowChartIcon from "@mui/icons-material/ShowChart";
import { useEffect, useRef, useState } from 'react';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';

// Components
import { WorkOrderTrend } from './components/workOrderTrend';
import { TestTypeDistribution } from './components/testTypeDistribution';
import { TopTestsOrdered } from './components/top10Ordered';
import { AgeGroupDistribution } from './components/patientDemographic';
import { AbnormalCriticalChart } from './components/abnormalCriticalChart';
import { GenderDistribution } from './components/genderDistribution';

interface DashboardPageProps {
    isWindow: boolean;
}

export const DashboardPage = ({ isWindow }: DashboardPageProps) => {
    const [summary, setSummary] = useState({
        total: 0,
        completed: 0,
        pending: 0,
        incomplete: 0,
        tests: 0,
        devices: 0,
        patients: 0,
        parameters: 0,
    });

    const [analyticData, setAnalyticData] = useState({
        work_order_trend: [],
        abnormal_summary: [],
        gender_summary: [],
        age_group: [],
        top_test_ordered: [],
        test_type_distribution: [],
    });

    useEffect(() => {

        // Open websocket for live updates
        let mounted = true;
        const backendBase = import.meta.env.VITE_BACKEND_BASE_URL || window.location.origin;
        const wsBase = backendBase.replace(/^http/, 'ws').replace(/\/$/, '');
        // backendBase may already include the API prefix (/api/v1). Avoid duplicating it.
        const apiPrefix = '/api/v1';
        const wsUrl = wsBase + (wsBase.includes(apiPrefix) ? '/summary/ws' : `${apiPrefix}/summary/ws`);

        const wsRef: { current: WebSocket | null } = { current: null };
        let reconnectHandle: number | null = null;

        const connect = () => {
            try {
                const socket = new WebSocket(wsUrl);
                wsRef.current = socket;

                socket.onopen = () => {
                    console.debug('summary WS connected', wsUrl);
                    if (reconnectHandle) {
                        clearTimeout(reconnectHandle);
                        reconnectHandle = null;
                    }
                };

                socket.onmessage = (ev: MessageEvent) => {
                    try {
                        const payload = JSON.parse(ev.data || '{}');
                        const s = payload.Summary || payload.summary;
                        const a = payload.Analytics || payload.analytics;
                        if (s) {
                            setSummary({
                                total: s.total_work_orders || 0,
                                completed: s.completed_work_orders || 0,
                                incomplete: s.incomplate_work_orders || 0,
                                pending: s.pending_work_orders || 0,
                                tests: s.total_test || 0,
                                devices: s.devices_connected || 0,
                                parameters: s.total_test_parameters || 0,
                                patients: s.total_patients || 0
                            });
                        }
                        if (a) {
                            setAnalyticData({
                                work_order_trend: a.work_order_trend || [],
                                abnormal_summary: a.abnormal_summary || [],
                                gender_summary: a.gender_summary || [],
                                age_group: a.age_group || [],
                                top_test_ordered: a.top_test_ordered || [],
                                test_type_distribution: a.test_type_distribution || [],
                            });
                        }
                    } catch (err) {
                        console.error('Failed to parse summary WS message', err);
                    }
                };

                socket.onclose = () => {
                    if (!mounted) return;
                    console.warn('summary WS closed, reconnecting in 2s');
                    reconnectHandle = window.setTimeout(() => connect(), 2000) as unknown as number;
                };

                socket.onerror = (err) => {
                    console.error('summary WS error', err);
                };
            } catch (err) {
                console.error('Failed to connect summary WS', err);
                reconnectHandle = window.setTimeout(() => connect(), 2000) as unknown as number;
            }
        };

        connect();

        return () => {
            mounted = false;
            if (wsRef.current) {
                try { wsRef.current.close(); } catch (e) { /* ignore */ }
            }
            if (reconnectHandle) clearTimeout(reconnectHandle);
        };
    }, []);

    const rootRef = useRef<HTMLDivElement | null>(null);
    const [isFullscreen, setIsFullscreen] = useState<boolean>(false);

    const toggleFullscreen = async () => {
        try {
            if (!document.fullscreenElement) {
                if (rootRef.current) {
                    await rootRef.current.requestFullscreen();
                } else {
                    await document.documentElement.requestFullscreen();
                }
            } else {
                await document.exitFullscreen();
            }
        } catch (err) {
            console.error('Failed to toggle fullscreen:', err);
        }
    };

    useEffect(() => {
        const onFsChange = () => setIsFullscreen(document.fullscreenElement === rootRef.current);
        document.addEventListener('fullscreenchange', onFsChange);
        return () => document.removeEventListener('fullscreenchange', onFsChange);
    }, []);

    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Box
            sx={{
                display: 'grid',
                backgroundColor: isFullscreen ? "#f9f9f9ff" : "",
                padding: isFullscreen ? 5 : 0,
                ...(isFullscreen
                    ? {
                        height: '100vh',
                        width: '100vw',
                        overflow: 'auto',
                        boxSizing: 'border-box',
                        position: 'relative',
                    }
                    : {}),
            }}
            ref={rootRef}
        >
            {/* SUMMARY SECTION */}
            <Box sx={{ mb: 2, display: 'flex', alignItems: "center", gap: 1, justifyContent: 'space-between' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <DashboardOutlinedIcon />
                    <Typography variant='h5'>SUMMARY</Typography>
                </Box>
                {!isWindow ? (
                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <Button
                            variant="contained"
                            color={isFullscreen ? 'secondary' : 'primary'}
                            startIcon={isFullscreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
                            onClick={toggleFullscreen}
                        >
                            {isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}
                        </Button>

                    </Box>
                ) : ""}
            </Box>

            {/* Cards */}
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(16, 1fr)', gap: 1 }}>
                {[
                    { icon: <AssignmentTurnedInIcon sx={{ mr: 1, color: 'primary.main' }} />, title: 'Total Work Orders Today', value: summary.total },
                    { icon: <CheckCircleIcon sx={{ mr: 1, color: 'success.main' }} />, title: 'Completed Work Orders Today', value: summary.completed },
                    { icon: <HourglassEmptyIcon sx={{ mr: 1, color: 'warning.main' }} />, title: 'Pending Work Orders Today', value: summary.pending },
                    { icon: <ErrorOutlineIcon sx={{ mr: 1, color: 'error.main' }} />, title: 'Incomplete Work Orders Today', value: summary.incomplete },
                    { icon: <ScienceIcon sx={{ mr: 1, color: 'primary.main' }} />, title: 'Total Tests (OBR/OBX) Today', value: summary.tests },
                    { icon: <DevicesIcon sx={{ mr: 1, color: 'secondary.main' }} />, title: 'Devices Connected', value: summary.devices },
                    { icon: <PeopleIcon sx={{ mr: 1, color: 'primary.main' }} />, title: 'Total Patients', value: summary.patients },
                    { icon: <ListAltIcon sx={{ mr: 1, color: 'info.main' }} />, title: 'Total Test Parameters', value: summary.parameters },
                ].map((item, idx) => (
                    <Box key={idx} sx={{ gridColumn: { xs: '1 / -1', md: 'span 4' } }}>
                        <Card>
                            <CardContent>
                                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                    {item.icon}
                                    <Typography variant="h6" color='gray'>{item.title}</Typography>
                                </Box>
                                <Typography variant="h4">{formatNumber(item.value)}</Typography>
                            </CardContent>
                        </Card>
                    </Box>
                ))}
            </Box>

            {/* ANALYTICS SECTION */}
            <Box sx={{ mt: 4, mb: 2, display: 'flex', alignItems: "center", gap: 1 }}>
                <ShowChartIcon />
                <Typography variant='h5'>DATA ANALYTICS</Typography>
            </Box>

            <Box>
                <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 1 }}>
                    <WorkOrderTrend data={analyticData.work_order_trend} />
                    <AbnormalCriticalChart data={analyticData.abnormal_summary} />
                </Box>

                <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 1, mt: 1 }}>
                    <TopTestsOrdered data={analyticData.top_test_ordered} />
                    <TestTypeDistribution data={analyticData.test_type_distribution} />
                </Box>

                <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 1, mt: 1 }}>
                    <AgeGroupDistribution data={analyticData.age_group} />
                    <GenderDistribution data={analyticData.gender_summary} />
                </Box>
            </Box>
        </Box>
    );
};
