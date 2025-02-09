import FileDownloadIcon from '@mui/icons-material/FileDownload';
import PrintIcon from '@mui/icons-material/Print';
import { IconButton, Stack, Tooltip } from '@mui/material';
import {
    BlobProvider,
    Font
} from '@react-pdf/renderer';
import dayjs from 'dayjs';
import { useEffect, useState } from 'react';
import type { ObservationResult, ReportData, ReportDataAbnormality, TestResult } from '../types/observation_result';
import type { Patient } from '../types/patient';
import type { WorkOrder } from '../types/work_order';
import { ReportDocument } from './ReportFile';

// Optional: Register custom fonts if required
Font.register({
    family: 'Roboto',
    src: 'https://fonts.gstatic.com/s/roboto/v27/KFOmCnqEu92Fr1Me5Q.ttf',
});

type PrintMCUProps = {
    patient: Patient
    workOrder: WorkOrder
    results: TestResult[]
}

const PrintMCUButton = (prop: PrintMCUProps) => {
    const [data, setData] = useState<ReportData[]>([])
    useEffect(() => {

        setData(prop.results.map(v => {
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
                default:
                    abnormality = "Normal"
                    break
            }

            const reportData: ReportData =  {
                category: v.category,
                parameter: v.test,
                reference: v.reference_range,
                unit: v.unit,
                result: v.result,
                abnormality: abnormality,
                subCategory: v.category,
            } 

            return reportData
        }))

    }, [prop.results]);

    // const [patientData, setPatientData] = useState<Patient | null>(null)
    // useEffect(() => {
    //     setPatientData(prop.patient)
    // }, [prop.patient]);

    return (
        <BlobProvider document={<ReportDocument data={data} patientData={prop.patient} />}>
            {({ url, loading, error }) => {
                if (error) {
                    return <span color='red'>Error generating PDF: {error.message}</span>;
                }

                return (
                    <Stack gap={1} direction={"row"}>
                        {/* Download PDF Button */}
                        <Tooltip title={loading ? "Loading..." : "Download PDF"}>
                            <IconButton
                                onClick={e => e.stopPropagation()}
                                color='primary'
                                download={`MCU_Result_${dayjs(prop.workOrder.created_at).format("YYYYMMDD")}_${prop.patient.id}_${prop.patient.first_name}_${prop.patient.last_name}.pdf`}
                                href={url || ''}
                                disabled={loading}
                            >
                                <FileDownloadIcon />
                            </IconButton>
                        </Tooltip>

                        {/* Print PDF Button */}
                        <Tooltip title={loading ? "Loading..." : "Print PDF"}>
                            <IconButton
                                color='secondary'
                                onClick={(e) => {
                                    e.stopPropagation()
                                    if (url) {
                                        window.open(url, '_blank')?.focus();
                                    }
                                }}
                                disabled={loading}
                            >
                                <PrintIcon />
                            </IconButton>
                        </Tooltip>
                    </Stack>
                );
            }}
        </BlobProvider>
    );
};

export default PrintMCUButton;
