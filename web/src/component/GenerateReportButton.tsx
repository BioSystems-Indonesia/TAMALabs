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
    const [isGenerating, setIsGenerating] = useState(false)

    const buttonId = prop.uniqueId || `${prop.workOrder.id}-${prop.patient.id}`

    const isGenerated = prop.currentGeneratedId === buttonId

    const reportDocument = useMemo(() =>
        <ReportDocument data={data} patientData={prop.patient} workOrderData={prop.workOrder} />,
        [data, prop.patient, prop.workOrder]
    )

    useEffect(() => {
        setData(prop.results?.map(v => {
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
                default:
                    abnormality = "Normal"
                    break
            }

            const reportData: ReportData = {
                category: v.category,
                parameter: v.test,
                reference: v.reference_range,
                unit: v.unit,
                result: v.formatted_result,
                abnormality: abnormality,
                subCategory: v.category,
            }

            return reportData
        }))
    }, [prop.results]);

    useEffect(() => {
        if (prop.currentGeneratedId && prop.currentGeneratedId !== buttonId) {
            setIsGenerating(false);
        }
    }, [prop.currentGeneratedId, buttonId]);

    const handleGenerateReport = () => {
        setIsGenerating(true)
        setTimeout(() => {
            setIsGenerating(false)
            if (prop.onGenerate) {
                prop.onGenerate(buttonId)
            }
        }, 800) // 1.5 seconds simulation
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
                                        onClick={e => e.stopPropagation()}
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
                                            e.stopPropagation()
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
        </div>
    );
};

export default GenerateReportButton;
