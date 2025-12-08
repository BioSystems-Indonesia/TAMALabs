import React from 'react';
import {
    Document,
    Font,
    Image,
    Page,
    StyleSheet,
    Text,
    View,
} from '@react-pdf/renderer';
import QRCode from 'qrcode';
import useSettings from '../hooks/useSettings';
import logo from '../assets/logo-selayar.png'
import type { ReportData } from '../types/observation_result';
import { Patient } from "../types/patient.ts";
import { WorkOrder } from '../types/work_order.ts';

Font.register({
    family: 'Helvetica',
    fonts: [
        { src: 'https://fonts.gstatic.com/s/roboto/v30/KFOmCnqEu92Fr1Mu4mxP.ttf', fontWeight: 400 },
        { src: 'https://fonts.gstatic.com/s/roboto/v30/KFOlCnqEu92Fr1MmEU9fBBc9.ttf', fontWeight: 700 },
    ],
});


const styles = StyleSheet.create({
    page: {
        fontSize: 9,
        fontFamily: 'Helvetica',
        paddingTop: 30,
        paddingLeft: 30,
        paddingRight: 30,
        paddingBottom: 30,
    },
    header: {
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'flex-start',
        justifyContent: 'flex-start',
        paddingBottom: 10,
        paddingLeft: 10
    },
    companyInfo: {
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: 9,
        textAlign: 'center',
    },
    logo: {
        width: 50,
        height: 50,
        marginRight: 15,
    },
    footer: {
        marginTop: 12,
        paddingVertical: 0,
        fontSize: 8,
    },
    category: {
        fontSize: 10,
        fontWeight: 'bold',
        marginTop: 10,
        marginBottom: 5,
        color: '#000000',
    },
    subCategory: {
        fontSize: 9,
        fontWeight: 'semibold',
        marginTop: 5,
        marginBottom: 5,
        color: '#000000',
    },
    tableHeader: {
        flexDirection: 'row',
        fontWeight: 'bold',
        paddingHorizontal: 3,
        borderWidth: 1,
        borderColor: '#000000',
    },
    tableRow: {
        flexDirection: 'row',
        paddingHorizontal: 3,
        borderLeftWidth: 1,
        borderRightWidth: 1,
        borderBottomWidth: 1,
        borderColor: '#000000',
    },
    columnNo: {
        width: '8%',
        textAlign: 'center',
    },
    columnParameter: {
        width: '25%',
    },
    columnParameterHeader: {
        width: '25%',
        textAlign: 'center'
    },
    columnResult: {
        width: '15%',
        textAlign: 'center',
    },
    columnMethod: {
        width: '22%',
        textAlign: 'center',
    },
    columnReference: {
        width: '30%',
        textAlign: 'center',
    },
    cell: {
        paddingVertical: 5,
        paddingHorizontal: 3,
        borderRight: 1,
        borderBottomColor: '#000',
    },
    cellEnd: {
        paddingVertical: 5,
        paddingHorizontal: 3,
    },
    rectangleContainer: {
        width: '100%',
        padding: 0,
        marginTop: 5,
        marginBottom: 10,
    },
    infoRow: {
        flexDirection: 'row',
        padding: 5,
    },
    infoRowLast: {
        flexDirection: 'row',
        padding: 5,
    },
    infoLabel: {
        width: '25%',
        fontSize: 9,
    },
    infoValue: {
        width: '25%',
        fontSize: 9,
    },
    interpretationBox: {
        marginTop: 15,
        borderWidth: 1,
        borderColor: '#000000',
        padding: 10,
        minHeight: 80,
    },
    signatureSection: {
        marginTop: 15,
        flexDirection: 'row',
        justifyContent: 'space-between',
    },
    signatureBox: {
        width: '45%',
        alignItems: 'center',
    },
    qrCode: {
        width: 60,
        height: 60,
        marginBottom: 5,
    },
});

// Helper function to generate QR code as data URL
const generateQRCode = async (text: string): Promise<string> => {
    try {
        return await QRCode.toDataURL(text, {
            width: 200,
            margin: 1,
            color: {
                dark: '#000000',
                light: '#FFFFFF',
            },
        });
    } catch (err) {
        console.error('Error generating QR code:', err);
        return '';
    }
};

// Helper function to calculate age and format birthdate
const calculateAgeAndBirthdate = (birthdate: string) => {
    const now = new Date();
    const birthDate = new Date(birthdate);

    // Calculate age
    let years = now.getFullYear() - birthDate.getFullYear();
    let months = now.getMonth() - birthDate.getMonth();
    let days = now.getDate() - birthDate.getDate();

    // Adjust for negative months or days
    if (months < 0 || (months === 0 && days < 0)) {
        years--;
        months += 12;
    }
    if (days < 0) {
        // borrow days from previous month
        const prevMonth = new Date(now.getFullYear(), now.getMonth(), 0);
        days += prevMonth.getDate();
        months--;
    }

    // Format birthdate as DD-MM-YYYY
    const day = birthDate.getDate();
    const month = birthDate.getMonth() + 1;
    const year = birthDate.getFullYear();

    if (years < 1) {
        // show months only if less than 1 year
        const displayMonths = months > 0 ? months : 0;
        return `${displayMonths} bulan / ${day}-${month}-${year}`;
    }

    return `${years} tahun / ${day}-${month}-${year}`;
};

// Helper function to format gender
const formatGender = (gender: string) => {
    if (gender === 'F') return 'Perempuan';
    if (gender === 'M') return 'Laki-laki';
    return '-';
};

// Helper function to determine method based on parameter name and who added it
const getMethodForParameter = (parameter?: string, addedBy?: 'System' | 'user') => {
    if (!parameter) return 'Spektrofotometer';
    const p = parameter.toLowerCase();

    // Glucose tests: system uses Spektrofotometer, user uses ICT
    if (p.includes('gula') || p.includes('glucose') || p.includes('glukosa')) {
        return addedBy === 'System' ? 'Spektrofotometer' : 'ICT';
    }

    // All other tests use Spektrofotometer
    return 'Spektrofotometer';
};

// Hardcoded reference values by parameter name
const DEFAULT_REFERENCE = '-';
const getReferenceForParameter = (parameter?: string, referenceFromDB?: string) => {
    // If no parameter, use reference from DB or default
    if (!parameter) return referenceFromDB || DEFAULT_REFERENCE;

    const p = parameter.toLowerCase();

    // Check hardcoded reference values first (highest priority)
    if (p.includes('gula darah puasa') || p.includes('gula darah (puasa)') || p.includes('glukosa puasa')) return '<126 mg/dl';
    if (p.includes('gula darah 2') || p.includes('2 jam')) return '<200 mg/dl';
    if (p.includes('asam urat')) return 'L=3.4-7.0 mg/dl; P=2.4-5.7 mg/dl';
    if (p.includes('sgot') || p.includes('ast')) return 'L<=42 U/L; P<=37 U/L';
    if (p.includes('sgpt') || p.includes('alt')) return 'L<=42 U/L; P<=32 U/L';
    if (p.includes('ureum') || p.includes('urea')) return '10-50 mg/dl';
    if (p.includes('kreatinin') || p.includes('creatinine')) return 'L<=1.1 mg/dl; P<=0.9 mg/dl';
    if (p.includes('trigliserid') || p.includes('trigliserida') || p.includes('triglyceride')) return '<200 mg/dl';
    if (p.includes('cholest') || p.includes('kolesterol') || p.includes('cholesterol')) {
        if (p.includes('hdl')) return 'L>=55 mg/dl; P>=65 mg/dl';
        if (p.includes('ldl')) return '<130 mg/dl';
        return '<200 mg/dl';
    }
    if (p.includes('albumin')) return '3.8-5.1 g/dl';
    if (p.includes('bilirubin total') || p.includes('bilirubin')) return '<1.1 mg/dl';
    if (p.includes('bilirubin direct') || p.includes('direct')) return '<0.25 mg/dl';
    if (p.includes('protein total') || p.includes('total protein')) return '6.6-8.7 g/dl';

    // If not found in hardcoded list, use reference from DB or default
    return referenceFromDB || DEFAULT_REFERENCE;
};

const Header = () => {
    const [settings] = useSettings();

    return (
        <View style={styles.header} fixed>
            <Image
                style={styles.logo}
                src={logo}
            />
            <View style={styles.companyInfo}>
                <Text style={{ fontSize: 12, fontWeight: 'bold', marginBottom: 2 }}>
                    {settings.company_name?.toUpperCase() || 'PEMERINTAH KABUPATEN KEPULAUAN SELAYAR'}
                </Text>
                <Text style={{ fontSize: 12, fontWeight: 'bold', marginBottom: 2 }}>
                    DINAS KESEHATAN
                </Text>
                <Text style={{ fontSize: 12, fontWeight: 'bold', marginBottom: 3 }}>
                    UPT.RSUD KH. HAYYUNG
                </Text>
                <Text style={{ fontSize: 9.5, fontWeight: 'bold' }}>
                    {settings.company_address || 'JL.KH.ABDUL KADIR HASIM TELP (0414)2707366 KEPULAUAN SELAYAR KODE POS : 92812'}
                </Text>
            </View>
        </View>
    )
};

const PatientInfo = ({ patient, workOrder }: { patient: Patient, workOrder: WorkOrder }) => {
    const formatDate = (date: Date | string) => {
        const d = typeof date === 'string' ? new Date(date) : date;
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const year = d.getFullYear();
        return `${day}-${month}-${year}`;
    };

    const formatDateTime = (date: Date | string) => {
        const d = typeof date === 'string' ? new Date(date) : date;
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        return `${hours}:${minutes}`;
    };

    const currentDate = new Date();

    return (
        <View style={styles.rectangleContainer}>
            {/* Row 1 */}
            <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Nama</Text>
                <Text style={styles.infoValue}>: {patient.first_name} {patient.last_name}</Text>
                <Text style={styles.infoLabel}>No.Kunjungan</Text>
                <Text style={styles.infoValue}>: {workOrder.visit_number || '-'}</Text>
            </View>

            {/* Row 2 */}
            <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Umur / Tgl Lahir</Text>
                <Text style={styles.infoValue}>: {calculateAgeAndBirthdate(patient.birthdate)}</Text>
                <Text style={styles.infoLabel}>No.RM</Text>
                <Text style={styles.infoValue}>: {workOrder.medical_record_number || '-'}</Text>
            </View>

            {/* Row 3 */}
            <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Jenis Kelamin</Text>
                <Text style={styles.infoValue}>: {formatGender(patient.sex)}</Text>
                <Text style={styles.infoLabel}>Tanggal</Text>
                <Text style={styles.infoValue}>: {formatDate(currentDate)}</Text>
            </View>

            {/* Row 4 */}
            <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Alamat</Text>
                <Text style={styles.infoValue}>: {patient.address || patient.location || '-'}</Text>
                <Text style={styles.infoLabel}>Jam Pengambilan</Text>
                <Text style={styles.infoValue}>: {workOrder.specimen_collection_date ? formatDateTime(workOrder.specimen_collection_date) : '-'}</Text>
            </View>

            {/* Row 5 */}
            <View style={styles.infoRowLast}>
                <Text style={styles.infoLabel}>Diagnosa / Ket. Klinis</Text>
                <Text style={styles.infoValue}>: {workOrder.diagnosis || '-'}</Text>
                <Text style={styles.infoLabel}>Jam Pengeluaran</Text>
                <Text style={styles.infoValue}>: {workOrder.result_release_date ? formatDateTime(workOrder.result_release_date) : '-'}</Text>
            </View>
        </View>
    );
};

const Footer = ({ workOrder, doctorQRCode, analyzerQRCode }: {
    workOrder: WorkOrder,
    doctorQRCode?: string,
    analyzerQRCode?: string
}) => (
    <View style={styles.footer}>
        {/* Interpretation Box */}
        <View style={styles.interpretationBox}>
            <Text style={{ fontSize: 9, marginBottom: 5 }}>Interpretasi Hasil :</Text>
        </View>

        {/* Interpretation Box */}
        <View style={styles.interpretationBox}>
            <Text style={{ fontSize: 9, marginBottom: 5 }}>Saran :</Text>
        </View>

        {/* Signature Section */}
        <View style={styles.signatureSection}>
            <View style={styles.signatureBox}>
                <Text style={{ fontSize: 9, marginBottom: 5 }}>Penanggung Jawab</Text>
                {doctorQRCode && <Image style={styles.qrCode} src={doctorQRCode} />}
                <Text style={{ fontSize: 9, fontWeight: 'bold', textDecoration: 'underline', marginTop: 2 }}>
                    {workOrder.doctors?.length > 0 ? workOrder.doctors[0].fullname : 'dr.Hj. Misnah, M.Kes, Sp.PK(K)'}
                </Text>
                <Text style={{ fontSize: 8 }}>
                    Nip : 19771204 200502 2 008
                </Text>
            </View>
            <View style={styles.signatureBox}>
                <Text style={{ fontSize: 9, marginBottom: 70 }}>Analis</Text>
                <Text style={{ fontSize: 9, fontWeight: 'bold' }}>
                    {workOrder.analyzers?.length > 0 ? workOrder.analyzers[0].fullname : ''}
                </Text>
            </View>
        </View>
    </View>
);

export const ReportDocument = ({ data, patientData, workOrderData, customDoctorQRText, customAnalyzerQRText }: {
    data: ReportData[],
    patientData: Patient,
    workOrderData: WorkOrder,
    groupedData?: { [category: string]: ReportData[] },
    customDoctorQRText?: string,
    customAnalyzerQRText?: string,
}) => {
    // Flatten data - remove grouping
    let rowNumber = 1;

    // Generate QR codes for doctor and analyzer
    const [doctorQRCode, setDoctorQRCode] = React.useState<string>('');
    const [analyzerQRCode, setAnalyzerQRCode] = React.useState<string>('');

    React.useEffect(() => {
        const generateQRCodes = async () => {
            // Use custom text if provided, otherwise use default logic
            const doctorText = customDoctorQRText ||
                (workOrderData.doctors?.length > 0
                    ? workOrderData.doctors[0].fullname
                    : 'dr.Hj. Misnah, M.Kes, Sp.PK(K)');

            const analyzerText = customAnalyzerQRText ||
                (workOrderData.analyzers?.length > 0
                    ? workOrderData.analyzers[0].fullname
                    : '');

            const doctorQR = await generateQRCode(doctorText);
            const analyzerQR = analyzerText ? await generateQRCode(analyzerText) : '';

            setDoctorQRCode(doctorQR);
            setAnalyzerQRCode(analyzerQR);
        };

        generateQRCodes();
    }, [workOrderData, customDoctorQRText, customAnalyzerQRText]);

    return (
        <Document>
            <Page size={"A4"} style={styles.page} wrap>
                <Header />
                <View style={{ width: '100%', height: 1, backgroundColor: '#000000', marginTop: -4 }} />
                <View style={{ width: '100%', height: 1, backgroundColor: '#000000', marginTop: 1 }} />
                <PatientInfo patient={patientData} workOrder={workOrderData} />
                <View style={{ width: '100%', height: 1, backgroundColor: '#000000', marginTop: -4 }} />
                <View style={{ width: '100%', height: 1, backgroundColor: '#000000', marginTop: 1 }} />
                <Text style={{
                    textAlign: 'center',
                    fontSize: 13,
                    fontWeight: 'bold',
                    marginVertical: 10
                }}>
                    HASIL PEMERIKSAAN KIMIA KLINIK
                </Text>

                {/* Table Header */}
                <View style={styles.tableHeader}>
                    <Text style={[styles.columnNo, styles.cell]}>No</Text>
                    <Text style={[styles.columnParameterHeader, styles.cell]}>Jenis Parameter</Text>
                    <Text style={[styles.columnResult, styles.cell]}>Hasil</Text>
                    <Text style={[styles.columnMethod, styles.cell]}>Metode</Text>
                    <Text style={[styles.columnReference, styles.cellEnd]}>Nilai Rujukan</Text>
                </View>

                {/* Table Rows */}
                {data.map((item, index) => {
                    const refValue = getReferenceForParameter(item.parameter, item.reference);
                    const methodValue = getMethodForParameter(item.parameter, item.added_by);
                    return (
                        <View key={index} style={styles.tableRow}>
                            <Text style={[styles.columnNo, styles.cell]}>{rowNumber++}</Text>
                            <Text style={[styles.columnParameter, styles.cell]}>{item.parameter}</Text>
                            <Text style={[styles.columnResult, styles.cell]}>{item.result || ''}</Text>
                            <Text style={[styles.columnMethod, styles.cell]}>{methodValue}</Text>
                            <Text style={[styles.columnReference, styles.cellEnd]}>{refValue}</Text>
                        </View>
                    );
                })}

                <Footer workOrder={workOrderData} doctorQRCode={doctorQRCode} analyzerQRCode={analyzerQRCode} />
            </Page>
        </Document>
    );
};

// Usage example remains the same