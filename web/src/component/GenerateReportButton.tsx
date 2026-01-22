import ArticleIcon from '@mui/icons-material/Article';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import PrintIcon from '@mui/icons-material/Print';
import { Button, CircularProgress, IconButton, Stack, Tooltip } from '@mui/material';
import {
    BlobProvider,
    Font
} from '@react-pdf/renderer';
import dayjs from 'dayjs';
import { useEffect, useState, useMemo } from 'react';
import { useNotify, useRefresh } from 'react-admin';
import useAxios from '../hooks/useAxios';
import type { ReportData, ReportDataAbnormality, TestResult } from '../types/observation_result';
import type { Patient } from '../types/patient';
import type { WorkOrder } from '../types/work_order';
import { ReportDocument } from './ReportFile';

Font.register({
    family: 'Roboto',
    src: 'https://fonts.gstatic.com/s/roboto/v27/KFOmCnqEu92Fr1Me5Q.ttf',
});

type GenerateReportButtonProps = {
    patient: Patient
    workOrder: WorkOrder
    results: TestResult[]
    uniqueId?: string
    currentGeneratedId?: string | null
    onGenerate?: (id: string) => void
}

const GenerateReportButton = (prop: GenerateReportButtonProps) => {
    const [data, setData] = useState<ReportData[]>([])
    // const [groupedData, setGroupedData] = useState<{ [category: string]: ReportData[] }>({})
    const [isGenerating, setIsGenerating] = useState(false)
    const axios = useAxios();
    const notify = useNotify();
    const refresh = useRefresh();

    const buttonId = prop.uniqueId || `${prop.workOrder.id}-${prop.patient.id}`

    const isGenerated = prop.currentGeneratedId === buttonId

    const updateReleaseDate = async () => {
        // Only update if result_release_date is not set yet
        if (!prop.workOrder.result_release_date || prop.workOrder.result_release_date === '') {
            const currentDate = new Date().toISOString();
            try {
                await axios.patch(`work-order/${prop.workOrder.id}/release-date`, {
                    result_release_date: currentDate
                });
                console.log('Release date updated successfully');

                // Fetch updated work order
                const response = await axios.get(`work-order/${prop.workOrder.id}`);
                if (response.data) {
                    // Update the prop by triggering a refresh
                    prop.workOrder.result_release_date = response.data.result_release_date;
                    refresh();
                }
            } catch (error) {
                console.error('Failed to update release date:', error);
                notify('Failed to update release date', { type: 'error' });
            }
        }
    };

    const reportDocument = useMemo(() => {
        // Format release date for QR code
        // const releaseDate = prop.workOrder.result_release_date
        //     ? dayjs(prop.workOrder.result_release_date).format('DD-MM-YYYY HH:mm')
        //     : dayjs().format('DD-MM-YYYY HH:mm');

        // const doctorQRText = `Dikeluarkan di UPT.RSUD KH. HAYYUNG, Kabupaten Kepulauan Selayar. Ditandatangani secara elektronik oleh dr.Hj. Misnah, M.Kes, Sp.PK(K), Pada tanggal: ${releaseDate}`;

        return (
            <ReportDocument
                data={data}
                patientData={prop.patient}
                workOrderData={prop.workOrder}
            // customDoctorQRText={doctorQRText}
            />
        );
    }, [data, prop.patient, prop.workOrder])

    useEffect(() => {
        // Add null check for prop.results
        if (!prop.results || !Array.isArray(prop.results)) {
            setData([]);
            return;
        }

        setData(prop.results.map(v => {
            // Add null check for each result item
            if (!v) return null;

            let abnormality = "Normal" as ReportDataAbnormality
            switch (v.abnormal) {
                case 0:
                    abnormality = "Normal"
                    break
                case 1:
                    abnormality = "High"
                    break
                case 2:
                    abnormality = "Low"
                    break
                case 3:
                    abnormality = "No Data"
                    break
                case 4:
                    abnormality = "Positive"
                    break
                case 5:
                    abnormality = "Negative"
                    break
                default:
                    abnormality = "Normal"
                    break
            }

            const reportData: ReportData = {
                category: v.category || '',
                parameter: v.test_type?.name || v.test || '',
                alias_code: v.test_type?.alias_code,
                reference: v.computed_reference_range || v.reference_range || '',
                unit: v.unit || '',
                result: String(v.formatted_result || v.result || ''), // Use formatted_result with proper decimal places
                abnormality: abnormality,
                subCategory: v.category || '',
            }

            return reportData
        }).filter(Boolean) as ReportData[]); // Filter out null values

        // Group data by category - same logic as PrintReport.tsx
        const reportData = prop.results.map(v => {
            if (!v) return null;

            let abnormality = "Normal" as ReportDataAbnormality
            switch (v.abnormal) {
                case 0:
                    abnormality = "Normal"
                    break
                case 1:
                    abnormality = "High"
                    break
                case 2:
                    abnormality = "Low"
                    break
                case 3:
                    abnormality = "No Data"
                    break
                case 4:
                    abnormality = "Positive"
                    break
                case 5:
                    abnormality = "Negative"
                    break
                default:
                    abnormality = "Normal"
                    break
            }

            const reportData: ReportData = {
                category: v.category || '',
                parameter: v.test_type?.name || v.test || '',
                alias_code: v.test_type?.alias_code,
                reference: v.computed_reference_range || v.reference_range || '',
                unit: v.unit || '',
                result: String(v.formatted_result || v.result || ''), // Use formatted_result with proper decimal places
                abnormality: abnormality,
                subCategory: v.category || '',
            }

            return reportData
        }).filter(Boolean) as ReportData[];

        setData(reportData);

        // Group data by category
        const grouped = reportData.reduce((acc, item) => {
            const category = item.category || 'Other';
            if (!acc[category]) {
                acc[category] = [];
            }
            acc[category].push(item);
            return acc;
        }, {} as { [category: string]: ReportData[] });

        console.log('GenerateReportButton - reportData:', reportData);
        console.log('GenerateReportButton - grouped data:', grouped);
        // setGroupedData(grouped);
    }, [prop.results]);

    useEffect(() => {
        if (prop.currentGeneratedId && prop.currentGeneratedId !== buttonId) {
            setIsGenerating(false);
        }
    }, [prop.currentGeneratedId, buttonId]);

    const handleGenerateReport = async () => {
        setIsGenerating(true)

        // Update release date when generating report
        await updateReleaseDate();

        setTimeout(() => {
            setIsGenerating(false)
            if (prop.onGenerate) {
                prop.onGenerate(buttonId)
            }
        }, 800) // 0.8 seconds simulation
    }
    if (!isGenerated) {
        return (
            <Button
                variant="contained"
                startIcon={isGenerating ? <CircularProgress size={16} color="inherit" /> : <ArticleIcon />}
                onClick={(e) => {
                    e.stopPropagation()
                    handleGenerateReport()
                }}
                disabled={isGenerating}
                size="small"
                sx={{
                    textTransform: 'none',
                    fontSize: '12px',
                    whiteSpace: 'nowrap',
                    minWidth: 'auto',
                    px: 2,
                    py: 0.5,
                    backgroundColor: isGenerating ? 'action.disabled' : 'primary.main',
                    '&:hover': {
                        backgroundColor: isGenerating ? 'action.disabled' : 'primary.dark',
                    },
                    '&:disabled': {
                        backgroundColor: 'action.disabled',
                        color: 'text.disabled'
                    },
                    transition: 'all 0.3s ease',
                }}
            >
                {isGenerating ? 'Generating...' : 'Generate Report'}
            </Button>
        )
    }

    // Show Print and Download buttons after generation
    return (
        <div style={{
            minWidth: '80px',
            display: 'flex',
            transition: 'all 0.3s ease'
        }}>
            {/* Only render BlobProvider when data is ready */}
            {data && Array.isArray(data) && data.length >= 0 && prop.patient && prop.workOrder ? (
                <BlobProvider document={reportDocument}>
                    {({ url, loading, error }) => {
                        if (error) {
                            return <span style={{ color: 'red', fontSize: '12px' }}>Error: {error.message}</span>;
                        }

                        return (
                            <Stack gap={0.5} direction={"row"} sx={{
                                opacity: loading ? 0.8 : 1,
                                transition: 'opacity 0.3s ease'
                            }}>
                                {/* Download PDF Button */}
                                <Tooltip title={
                                    loading ? "Preparing PDF..." : "Download PDF Report"
                                }>
                                    <span>
                                        <IconButton
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                updateReleaseDate();
                                            }}
                                            download={`MCU_Result_${dayjs(prop.workOrder.created_at).format("YYYYMMDD")}_${prop.patient.id}_${prop.patient.first_name}_${prop.patient.last_name}.pdf`}
                                            href={url || ''}
                                            disabled={loading}
                                            size="small"
                                            sx={{
                                                backgroundColor: loading ? 'action.disabled' : '#e3f2fd',
                                                color: loading ? 'text.disabled' : '#1976d2',
                                                '&:hover': {
                                                    backgroundColor: loading ? 'action.disabled' : '#bbdefb',
                                                    color: loading ? 'text.disabled' : '#0d47a1',
                                                },
                                                '&:disabled': {
                                                    backgroundColor: 'action.disabled',
                                                    color: 'text.disabled',
                                                },
                                                transition: 'all 0.3s ease'
                                            }}
                                        >
                                            {loading ? <CircularProgress size={16} color="inherit" /> : <FileDownloadIcon fontSize="small" />}
                                        </IconButton>
                                    </span>
                                </Tooltip>

                                {/* Print PDF Button */}
                                <Tooltip title={
                                    loading ? "Preparing PDF..." : "Print PDF Report"
                                }>
                                    <span>
                                        <IconButton
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                updateReleaseDate();
                                                if (url) {
                                                    window.open(url, '_blank')?.focus();
                                                }
                                            }}
                                            disabled={loading}
                                            size="small"
                                            sx={{
                                                backgroundColor: loading ? 'action.disabled' : '#f3e5f5',
                                                color: loading ? 'text.disabled' : '#7b1fa2',
                                                '&:hover': {
                                                    backgroundColor: loading ? 'action.disabled' : '#e1bee7',
                                                    color: loading ? 'text.disabled' : '#4a148c',
                                                },
                                                '&:disabled': {
                                                    backgroundColor: 'action.disabled',
                                                    color: 'text.disabled',
                                                },
                                                transition: 'all 0.3s ease'
                                            }}
                                        >
                                            {loading ? <CircularProgress size={16} color="inherit" /> : <PrintIcon fontSize="small" />}
                                        </IconButton>
                                    </span>
                                </Tooltip>
                            </Stack>
                        );
                    }}
                </BlobProvider>
            ) : (
                <Stack gap={0.5} direction={"row"}>
                    <Tooltip title="Loading report data...">
                        <span>
                            <IconButton disabled size="small">
                                <CircularProgress size={16} />
                            </IconButton>
                        </span>
                    </Tooltip>
                </Stack>
            )}
        </div>
    );
};

export default GenerateReportButton;
