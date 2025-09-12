import {
    Document,
    Font,
    Image,
    Page,
    StyleSheet,
    Text,
    View,
} from '@react-pdf/renderer';
import useSettings from '../hooks/useSettings';
import logo from '../assets/elgatama-logo.png'
import yt from '../assets/youtube.png'
import fb from '../assets/facebook.png'
import ig from '../assets/instagram.png'
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
        fontSize: 10,
        fontFamily: 'Helvetica',
        padding: 40,
    },
    header: {
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
    },
    companyInfo: {
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-end',
        justifyContent: 'center',
        fontSize: 7.5,
        textAlign: 'right',
    },
    logo: {
        width: 64,
        height: 64,
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

// Helper function to format birthdate
const formatBirthdate = (birthdate: string) => {
    const date = new Date(birthdate);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
};

// Helper function to calculate age
const calculateAge = (birthdate: string) => {
    const now = new Date();
    const birthDate = new Date(birthdate);

    let years = now.getFullYear() - birthDate.getFullYear();
    let months = now.getMonth() - birthDate.getMonth();
    let days = now.getDate() - birthDate.getDate();

    // Adjust for negative months or days
    if (months < 0 || (months === 0 && days < 0)) {
        years--;
        months += 12;
    }
    if (days < 0) {
        const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 0);
        days += lastMonth.getDate();
        months--;
    }

    // return `${years} year(s), ${months} month(s), ${days} day(s)`;
    return `${years} year(s)`;
};

// Helper function to format gender
const formatGender = (gender: string) => {
    if (gender === 'F') return 'Female';
    if (gender === 'M') return 'Male';
    return '-';
};

const Header = () => {
    const [settings] = useSettings();

    return (
        <View style={styles.header} fixed>
            <Image
                style={styles.logo}
                src={logo}
            />
            <View style={{ width: '85%' }}>
                <Text style={{
                    fontSize: 24
                }}>{settings.company_name}</Text>
                <View style={{
                    width: '100%',
                    height: '0.2rem',
                    backgroundColor: 'rgb(74, 186, 171)'
                }}>
                </View>
                <View style={styles.companyInfo}>
                    <Text
                        wrap={true}
                        style={{ width: '45%' }}
                    >{settings.company_address}</Text>
                    <Text
                        wrap={true}
                        style={{ width: '45%' }}
                    >{settings.company_contact_phone}</Text>
                    <Text
                        wrap={true}
                        style={{ width: '45%' }}
                    >{settings.company_contact_email}</Text>
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
                        <Text style={styles.value}>: {workOrder.analyst?.length > 0 ? workOrder.analyst[0].fullname : ""}</Text>
                    </Text>
                </View>
            </View>
        </View>
    </View>
);

const Footer = () => (
    <View style={styles.footer}>
        <View style={{
            height: '0.2rem',
            backgroundColor: 'rgb(74, 186, 171)'
        }}>
        </View>
        <View style={{
            marginTop: 4,
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: 'space-between',
        }}>
            <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
                <Image src={yt} style={{ width: 15, height: 15 }} />
                <Text>BioSystems Indonesia</Text>
            </View>
            <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
                <Image src={ig} style={{ width: 15, height: 15 }} />
                <Text>@biosystems.ind</Text>
            </View>
            <View style={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 6 }}>
                <Image src={fb} style={{ width: 15, height: 15 }} />
                <Text>BioSystems Indonesia</Text>
            </View>
        </View>
    </View>
);

// Helper function to group data by category (if needed)
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

export const ReportDocument = ({ data, patientData, workOrderData }: { data: ReportData[], patientData: Patient, workOrderData: WorkOrder }) => {
    const groupedData = groupData(data);

    return (
        <Document>
            <Page size={"A4"} style={styles.page} wrap>
                <Header />
                <View style={{
                    marginTop: 20,
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
                        {items.map((item, index) => {
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
                <Footer />
            </Page>
        </Document>
    );
};

// Usage example remains the same