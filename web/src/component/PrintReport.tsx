import FileDownloadIcon from '@mui/icons-material/FileDownload';
import PrintIcon from '@mui/icons-material/Print';
import { IconButton, Stack, Tooltip } from '@mui/material';
import {
    BlobProvider,
    Font
} from '@react-pdf/renderer';
import dayjs from 'dayjs';
import { useEffect, useState } from 'react';
import type { ReportData, ReportDataAbnormality, TestResult } from '../types/observation_result';
import type { Patient } from '../types/patient';
import type { WorkOrder } from '../types/work_order';
import { ReportDocument } from './ReportFile';

// Optional: Register custom fonts if required
Font.register({
    family: 'Roboto',
    src: 'https://fonts.gstatic.com/s/roboto/v27/KFOmCnqEu92Fr1Me5Q.ttf',
});

type PrintReportButtonProps = {
    patient: Patient
    workOrder: WorkOrder
    results: TestResult[]
}

const PrintReportButton = (prop: PrintReportButtonProps) => {
    const [data, setData] = useState<ReportData[]>([])
    // const [groupedData, setGroupedData] = useState<{ [category: string]: ReportData[] }>({})

    useEffect(() => {
        const reportData = prop.results?.map(v => {
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

            const aliasCode = v.test_type?.alias_code;
            let displayResult: string = v.result || ""; // Use string result instead of formatted_result
            let displayUnit = v.unit;
            let displayReference = v.computed_reference_range || v.reference_range;

            // Convert specific test results to /µL
            if (aliasCode === "Jumlah Trombosit" || aliasCode === "Jumlah Leukosit") {
                // Parse the string result and multiply by 1000
                const numericResult = parseFloat(v.result || "0");
                displayResult = (numericResult * 1000).toLocaleString();
                displayUnit = '/µL';

                // Convert reference range by multiplying by 1000
                if (displayReference) {
                    const rangeMatch = displayReference.match(/(\d+\.?\d*)\s*-\s*(\d+\.?\d*)/);
                    if (rangeMatch) {
                        const lowRef = parseFloat(rangeMatch[1]) * 1000;
                        const highRef = parseFloat(rangeMatch[2]) * 1000;
                        displayReference = `${lowRef.toLocaleString()} - ${highRef.toLocaleString()}`;
                    }
                }
            }

            const reportData: ReportData = {
                category: v.category,
                parameter: v.test_type?.name || v.test,
                alias_code: aliasCode,
                reference: displayReference,
                unit: displayUnit,
                result: displayResult,
                abnormality: abnormality,
                subCategory: v.category,
            }

            return reportData
        }) || []

        setData(reportData)

        // Group data by category
        const grouped = reportData.reduce((acc, item) => {
            const category = item.category || 'Other';
            if (!acc[category]) {
                acc[category] = [];
            }
            acc[category].push(item);
            return acc;
        }, {} as { [category: string]: ReportData[] });

        console.log('PrintReport - reportData:', reportData);
        console.log('PrintReport - grouped data:', grouped);
        // setGroupedData(grouped);

    }, [prop.results]);

    // const [patientData, setPatientData] = useState<Patient | null>(null)
    // useEffect(() => {
    //     setPatientData(prop.patient)
    // }, [prop.patient]);

    // groupedData={groupedData}
    return (
        <BlobProvider document={<ReportDocument data={data} patientData={prop.patient} workOrderData={prop.workOrder} />}>
            {({ url, loading, error }) => {
                if (error) {
                    return <span color='red'>Error generating PDF: {error.message}</span>;
                }

                return (
                    <Stack gap={1} direction={"row"}>
                        {/* Download PDF Button */}
                        <Tooltip title={
                            loading ? "Loading..." :
                                // !prop.workOrder.have_complete_data ? `Hasil belum lengkap. Unduh akan tersedia ketika semua tes selesai.` :
                                "Download PDF"
                        }>
                            <span>
                                <IconButton
                                    onClick={e => e.stopPropagation()}
                                    color='primary'
                                    download={`LAB_Test_Result_${dayjs(prop.workOrder.created_at).format("YYYYMMDD")}_${prop.patient.id}_${prop.patient.first_name}_${prop.patient.last_name}.pdf`}
                                    href={url || ''}
                                    disabled={loading
                                        // || !prop.workOrder.have_complete_data
                                    }
                                >
                                    <FileDownloadIcon />
                                </IconButton>
                            </span>
                        </Tooltip>

                        {/* Print PDF Button */}
                        <Tooltip title={
                            loading ? "Loading..." :
                                // !prop.workOrder.have_complete_data ? `Hasil belum lengkap. Unduh akan tersedia ketika semua tes selesai.` :
                                "Print PDF"
                        }>
                            <span>
                                <IconButton
                                    color='secondary'
                                    onClick={(e) => {
                                        e.stopPropagation()
                                        if (url
                                            // && prop.workOrder.have_complete_data
                                        ) {
                                            window.open(url, '_blank')?.focus();
                                        }
                                    }}
                                    disabled={loading
                                        // || !prop.workOrder.have_complete_data
                                    }
                                >
                                    <PrintIcon />
                                </IconButton>
                            </span>
                        </Tooltip>
                    </Stack>
                );
            }}
        </BlobProvider>
    );
};

export default PrintReportButton;
