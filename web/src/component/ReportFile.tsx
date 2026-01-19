import {
    Document,
    Font,
    Image,
    Page,
    StyleSheet,
    Text,
    View,
    // Svg,
    // Text as SvgText,
} from '@react-pdf/renderer';
// import useSettings from '../hooks/useSettings';
import logo from '../assets/trigas-logo.png'
// import yt from '../assets/youtube.png'
// import fb from '../assets/facebook.png'
// import ig from '../assets/instagram.png'
import trigasKop from '../assets/terigas-kop.png'
import aptos from '../assets/font/Aptos.ttf'
import aptosBold from '../assets/font/Aptos-Bold.ttf'
import aptosSemiBold from '../assets/font/Aptos-SemiBold.ttf'
import type { ReportData } from '../types/observation_result';
import { Patient } from "../types/patient.ts";
import { WorkOrder } from '../types/work_order.ts';

Font.register({
    family: 'Aptos',
    fonts: [
        { src: aptos, fontWeight: 400 },
        { src: aptosSemiBold, fontWeight: 500 },
        { src: aptosBold, fontWeight: 700 },
    ],
});


const styles = StyleSheet.create({
    page: {
        fontSize: 10,
        fontFamily: 'Aptos',
        paddingTop: 10,
        paddingLeft: 40,
        paddingRight: 40,
        paddingBottom: 49,
    },
    header: {
        marginBottom: 20,
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
    },
    headerTitleContainer: {
        position: 'relative',
        textAlign: 'center',
        height: 36,
    },
    titleStroke: {
        position: 'absolute',
        fontSize: 28,
        fontWeight: 500,
        color: '#4C94D8',
        opacity: 1,
        fontFamily: 'Aptos',
        textAlign: 'center',


    },
    title: {
        textAlign: 'center',
        fontSize: 28,
        fontWeight: 500,
        color: '#a5c9eb',
        fontFamily: 'Aptos',
    },
    companyInfo: {
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: 12,
        textAlign: 'center',
        color: '#4C94D8'
    },
    logo: {
        width: 120,
        height: 120,
        marginBottom: 10,
    },
    footer: {
        position: 'absolute',
        bottom: 30,
        padding: '0 40px',
        left: 0,
        right: 0,
        textAlign: 'center',
        fontSize: 9,
        color: '#666666',
    },
    category: {
        fontSize: 14,
        fontWeight: 'bold',
        marginTop: 15,
        marginBottom: 8,
        color: '#2d3748',
    },
    subCategory: {
        fontSize: 12,
        fontWeight: 'semibold',
        marginTop: 5,
        marginBottom: 5,
        color: '#4a5568',
    },
    tableHeader: {
        color: 'white',
        flexDirection: 'row',
        backgroundColor: '#4abaab',
        fontWeight: 'bold',
        paddingVertical: 8,
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
        borderRadius: 2,
    },
    tableRow: {
        flexDirection: 'row',
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
        paddingVertical: 8,
    },
    columnHeader: {
        width: '30%',
        paddingHorizontal: 6,
    },
    columnResult: {
        width: '20%',
        paddingHorizontal: 6,
    },
    columnUnit: {
        width: '20%',
        paddingHorizontal: 6,
    },
    columnReference: {
        width: '20%',
        paddingHorizontal: 6,
    },
    cell: {
        paddingHorizontal: 6,
    },
    rectangleContainer: {
        width: '100%',
        borderWidth: 1,
        borderColor: '#cad5e2',
        borderRadius: 2,
        padding: 10,
    },
    gridContainer: {
        flexDirection: 'column',
    },
    row: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        padding: 2,
    },
    leftColumn: {
        flex: 1,
        paddingRight: 10,
    },
    rightColumn: {
        flex: 1,
    },
    labelValue: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    label: {
        fontWeight: 'bold',
        width: 80,
        display: 'flex',
    },
    value: {
        flex: 1,
    },
});

const formatBirthdate = (birthdate: string) => {
    const date = new Date(birthdate);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
};

const calculateAge = (birthdate: string) => {
    const now = new Date();
    const birthDate = new Date(birthdate);

    let years = now.getFullYear() - birthDate.getFullYear();
    let months = now.getMonth() - birthDate.getMonth();
    let days = now.getDate() - birthDate.getDate();

    if (months < 0 || (months === 0 && days < 0)) {
        years--;
        months += 12;
    }
    if (days < 0) {
        const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 0);
        days += lastMonth.getDate();
        months--;
    }

    return `${years} year(s)`;
};

const formatGender = (gender: string) => {
    if (gender === 'F') return 'Female';
    if (gender === 'M') return 'Male';
    return '-';
};

const Header = () => {
    // const [settings] = useSettings();

    return (
        <View style={styles.header} fixed>
            <Image
                style={styles.logo}
                src={logo}
            />
            <View style={{ width: '85%' }}>
                <Image src={trigasKop} style={{ width: '350px', marginLeft: 26, marginBottom: 6 }} />

                <View style={styles.companyInfo}>

                    <Text
                        wrap={true}
                        style={{ width: '100%', fontSize: 12 }}
                    >Jln. Trans Kalimantan Kec. Bulik, Kab. Lamandau Provinsi Kalimantan Tengah</Text>
                    <Text
                        wrap={true}
                        style={{
                            width: '100%',
                            fontSize: 12
                        }}
                    >Hp. 0895-3210-65003</Text>
                </View>
            </View>
        </View>
    )
};

const PatientInfo = ({ patient, workOrder }: { patient: Patient, workOrder: WorkOrder }) => (
    <View style={styles.rectangleContainer}>
        <View style={styles.gridContainer}>
            {/* Row 1 */}
            <View style={styles.row}>
                <View style={styles.leftColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Barcode No</Text>
                        <Text style={styles.value}>: {workOrder.barcode}</Text>
                    </Text>
                </View>
                <View style={styles.rightColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Date of Birth</Text>
                        <Text style={styles.value}>: {formatBirthdate(patient.birthdate)}</Text>
                    </Text>
                </View>
            </View>

            {/* Row 2 */}
            <View style={styles.row}>
                <View style={styles.leftColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Patient Name</Text>
                        <Text style={styles.value}>: {patient.first_name} {patient.last_name}</Text>
                    </Text>
                </View>
                <View style={styles.rightColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Age</Text>
                        <Text style={styles.value}>: {calculateAge(patient.birthdate)}</Text>
                    </Text>
                </View>
            </View>

            {/* Row 3 */}
            <View style={styles.row}>
                <View style={styles.leftColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Address</Text>
                        <Text style={styles.value}>: {patient.address}</Text>
                    </Text>
                </View>
                <View style={styles.rightColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Gender</Text>
                        <Text style={styles.value}>: {formatGender(patient.sex)}</Text>
                    </Text>
                </View>
            </View>

            {/* Row 4 */}
            <View style={styles.row}>
                <View style={styles.leftColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Doctor</Text>
                        <Text style={styles.value}>: {workOrder.doctors?.length > 0 ? workOrder.doctors[0].fullname : ""}</Text>
                    </Text>
                </View>
                <View style={styles.rightColumn}>
                    <Text style={styles.labelValue}>
                        <Text style={styles.label}>Analyst</Text>
                        <Text style={styles.value}>: {workOrder.analyzers?.length > 0 ? workOrder.analyzers[0].fullname : ""}</Text>
                    </Text>
                </View>
            </View>
        </View>
    </View>
);

// const Footer = () => (
//     <View style={styles.footer} fixed>
//         <View style={{
//             height: '0.2rem',
//             backgroundColor: 'rgb(74, 186, 171)'
//         }}>
//         </View>
//         <View style={{
//             marginTop: 4,
//             display: 'flex',
//             flexDirection: 'row',
//             alignItems: 'center',
//             justifyContent: 'space-between',
//         }}>
//             <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
//                 <Image src={yt} style={{ width: 15, height: 15 }} />
//                 <Text>BioSystems Indonesia</Text>
//             </View>
//             <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
//                 <Image src={ig} style={{ width: 15, height: 15 }} />
//                 <Text>@biosystems.ind</Text>
//             </View>
//             <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
//                 <Image src={fb} style={{ width: 15, height: 15 }} />
//                 <Text>BioSystems Indonesia</Text>
//             </View>
//         </View>
//     </View>
// );

const groupData = (data: ReportData[]) => {
    const grouped: Record<string, ReportData[]> = {};
    data?.forEach((item) => {
        if (!grouped[item.category]) {
            grouped[item.category] = [];
        }
        grouped[item.category].push(item);
    });
    return grouped;
};

export const ReportDocument = ({ data, patientData, workOrderData }: {
    data: ReportData[], patientData: Patient, workOrderData: WorkOrder, groupedData?: { [category: string]: ReportData[] },
}) => {
    const groupedData = groupData(data);

    return (
        <Document>
            <Page size={"A4"} style={styles.page} wrap>
                <Header />
                <View style={{
                    marginBottom: 15,
                }}>
                    <Text style={{
                        textAlign: 'center',
                        fontSize: 12,
                        fontWeight: 'bold',
                    }}>TEST RESULT</Text>
                </View>
                <PatientInfo patient={patientData} workOrder={workOrderData} />
                {Object.entries(groupedData).map(([category, items]) => (
                    <View key={category} wrap>
                        <Text style={styles.category}>{category}</Text>

                        {/* Table Header */}
                        <View style={styles.tableHeader}>
                            <Text style={[styles.columnHeader, styles.cell]}>Parameter</Text>
                            <Text style={[styles.columnResult, styles.cell]}>Result</Text>
                            <Text style={[styles.columnUnit, styles.cell]}>Unit</Text>
                            <Text style={[styles.columnReference, styles.cell]}>Reference</Text>
                            <Text style={[styles.columnReference, styles.cell]}>Status</Text>
                        </View>

                        {/* Table Rows */}
                        {[...items].sort((a, b) => {
                            const aKey = ((a.alias_code || a.parameter) || '').toString().toUpperCase();
                            const bKey = ((b.alias_code || b.parameter) || '').toString().toUpperCase();
                            if (aKey === 'WBC' && bKey !== 'WBC') return -1;
                            if (bKey === 'WBC' && aKey !== 'WBC') return 1;
                            return 0;
                        }).map((item, index) => {
                            const abnormalColor = {
                                color:
                                    item.abnormality === 'High' ? '#e53e3e' : // Red for High
                                        item.abnormality === 'Low' ? '#3182ce' :  // Blue for Low
                                            '#222222', // Default color for No Data or other cases
                            };
                            return (
                                <View key={index} style={styles.tableRow}>
                                    <Text style={[styles.columnHeader, styles.cell, abnormalColor]}>{item.parameter}</Text>
                                    <Text style={[styles.columnResult, styles.cell, abnormalColor]}>{item.result}</Text>
                                    <Text style={[styles.columnResult, styles.cell, abnormalColor]}>{item.unit}</Text>
                                    <Text style={[styles.columnReference, styles.cell, abnormalColor]}>
                                        {item.reference}
                                    </Text>
                                    <Text style={[styles.columnReference, styles.cell, abnormalColor]}>
                                        {item.abnormality}
                                    </Text>
                                </View>
                            );
                        })}
                    </View>
                ))}
                {/* <Footer /> */}
            </Page>
        </Document>
    );
};
