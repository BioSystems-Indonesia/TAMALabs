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
import type { ReportData, ReportDataAbnormality } from '../types/observation_result';

Font.register({
    family: 'Helvetica',
    fonts: [
        { src: 'https://fonts.gstatic.com/s/roboto/v30/KFOmCnqEu92Fr1Mu4mxP.ttf', fontWeight: 400 },
        { src: 'https://fonts.gstatic.com/s/roboto/v30/KFOlCnqEu92Fr1MmEU9fBBc9.ttf', fontWeight: 700 },
    ],
});


const styles = StyleSheet.create({
    page: {
        padding: 40,
        fontSize: 10,
        fontFamily: 'Helvetica',
    },
    header: {
        display: "flex",
        flexDirection: "row",
        gap: "12px",
        justifyContent: "space-around",
        alignItems: "center",
        marginBottom: 20,
        paddingBottom: 10,
        borderBottomWidth: 1,
        borderBottomColor: '#112131',
    },
    companyInfo: {
        maxWidth: '80%', 
        display: 'flex',
        gap: "2px",
    },
    logo: {
        width: 64,
        height: 64,
    },
    footer: {
        position: 'absolute',
        bottom: 30,
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
        flexDirection: 'row',
        backgroundColor: '#f7fafc',
        fontWeight: 'bold',
        paddingVertical: 8,
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
    },
    tableRow: {
        flexDirection: 'row',
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
        paddingVertical: 8,
    },
    columnHeader: {
        width: '40%',
        paddingHorizontal: 6,
    },
    columnResult: {
        width: '30%',
        paddingHorizontal: 6,
    },
    columnReference: {
        width: '30%',
        paddingHorizontal: 6,
    },
    cell: {
        paddingHorizontal: 6,
    },
});

const Header = () => {
    const [settings] = useSettings();

    return (
        <View style={styles.header} fixed>
            <Image
                style={styles.logo}
                src={logo}
            />
            <View >
                <Text wrap style={{
                    fontSize: "16px",
                    fontWeight: 'black',
                    marginBottom: '4px',
                }}>{settings.company_name}</Text>
                <View style={styles.companyInfo}>
                    <Text wrap>{settings.company_address}</Text>
                    {settings.company_contact_phone ? <Text wrap>Phone: {settings.company_contact_phone}</Text> : undefined}
                    {settings.company_contact_email ? <Text wrap>Email: {settings.company_contact_email}</Text> : undefined}
                </View>
            </View>
        </View>
    )
};

const Footer = () => (
    <Text style={styles.footer} fixed render={({ pageNumber, totalPages }) => (
        <Text>Page {pageNumber} of {totalPages}</Text>
    )} />
);

const groupData = (data: ReportData[]) => {
    return data.reduce((acc, item) => {
        const { category, subCategory } = item;
        if (!acc[category]) acc[category] = {};
        if (!acc[category][subCategory]) acc[category][subCategory] = [];
        acc[category][subCategory].push(item);
        return acc;
    }, {} as Record<string, Record<string, ReportData[]>>);
};

export const ReportDocument = ({ data }: { data: ReportData[] }) => {
    const groupedData = groupData(data);

    return (
        <Document>
            <Page size={"A4"} style={styles.page} wrap >
                <Header />
                {Object.entries(groupedData).map(([category, subCategories]) => (
                    <View key={category} wrap>
                        <Text style={styles.category}>{category}</Text>

                        {Object.entries(subCategories).map(([subCategory, items]) => (
                            <View key={subCategory}>
                                <Text style={styles.subCategory}>{subCategory}</Text>

                                {/* Table Header */}
                                <View style={styles.tableHeader}>
                                    <Text style={[styles.columnHeader, styles.cell]}>Parameter</Text>
                                    <Text style={[styles.columnResult, styles.cell]}>Result</Text>
                                    <Text style={[styles.columnReference, styles.cell]}>Reference</Text>
                                    <Text style={[styles.columnReference, styles.cell]}>Status</Text>
                                </View>

                                {/* Table Rows */}
                                {items.map((item, index) => {
                                    const abnormal = ['High', 'Low', "No Data"] as ReportDataAbnormality[];
                                    const isAbnormal = abnormal.includes(item.abnormality);
                                    const abnormalColor = {
                                        color: isAbnormal ? '#e53e3e' : '#222222',
                                    }
                                    return (
                                        <View key={index} style={styles.tableRow}>
                                            <Text style={[styles.columnHeader, styles.cell, abnormalColor]}>{item.parameter}</Text>
                                            <Text style={[styles.columnResult, styles.cell, abnormalColor]}>{item.result}</Text>
                                            <Text style={[
                                                styles.columnReference,
                                                styles.cell,
                                                abnormalColor
                                            ]}>
                                                {item.reference}
                                            </Text>
                                            <Text style={[
                                                styles.columnReference,
                                                styles.cell,
                                                abnormalColor
                                            ]}>
                                                {item.abnormality}
                                            </Text>
                                        </View>
                                    );
                                })}
                            </View>
                        ))}
                    </View>
                ))}
                <Footer />
            </Page>
        </Document>
    );
};

// Usage example remains the same