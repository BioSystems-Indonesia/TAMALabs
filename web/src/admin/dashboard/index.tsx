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
import useAxios from "../../hooks/useAxios";
import { useEffect, useState } from 'react';

// Components
import { WorkOrderTrend } from './components/workOrderTrend';
import { TestTypeDistribution } from './components/testTypeDistribution';
import { TopTestsOrdered } from './components/top10Ordered';
import { AgeGroupDistribution } from './components/patientDemographic';
import { AbnormalCriticalChart } from './components/abnormalCriticalChart';
import { GenderDistribution } from './components/genderDistribution';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';

interface DashboardPageProps {
    isWindow: boolean;
}

export const DashboardPage = ({ isWindow }: DashboardPageProps) => {
    const axios = useAxios()
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
        axios.get('/summary/analytics')
            .then(resp => {
                const data = resp.data || {};
                console.log(data)
                setAnalyticData({
                    work_order_trend: data.work_order_trend || [],
                    abnormal_summary: data.abnormal_summary || [],
                    gender_summary: data.gender_summary || [],
                    age_group: data.age_group || [],
                    top_test_ordered: data.top_test_ordered || [],
                    test_type_distribution: data.test_type_distribution || [],
                });
            })
            .catch(err => {
                console.error("Failed to load summary data:", err);
            });
        axios.get('/summary/')
            .then(resp => {
                const data = resp.data || {};
                console.log(data)
                setSummary({
                    total: data.total_work_orders,
                    completed: data.completed_work_orders,
                    incomplete: data.incomplate_work_orders,
                    pending: data.pending_work_orders,
                    tests: data.total_test,
                    devices: data.devices_connected,
                    parameters: data.total_test_parameters,
                    patients: data.total_patients

                })
            })
    }, []);

    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Box sx={{ display: 'grid' }}>
            {/* SUMMARY SECTION */}
            <Box sx={{ mb: 2, display: 'flex', alignItems: "center", gap: 1, justifyContent: 'space-between' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <DashboardOutlinedIcon />
                    <Typography variant='h5'>SUMMARY</Typography>
                </Box>
                {!isWindow ? <Box>
                    <Button
                        variant="contained"
                        color="primary"
                        startIcon={<OpenInNewIcon />}
                        onClick={() => {
                            try {
                                const url = `${window.location.origin}/dashboard-window`;
                                window.open(url, '_blank', 'noopener,noreferrer');
                            } catch (err) {
                                window.open('/dashboard-window', '_blank', 'noopener,noreferrer');
                            }
                        }}
                    >
                        Open Window
                    </Button>
                </Box> : ""}
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
