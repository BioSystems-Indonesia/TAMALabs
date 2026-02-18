import {
    Document,
    Font,
    // Image,
    Page,
    StyleSheet,
    Text,
    View,
} from '@react-pdf/renderer';
// import useSettings from '../hooks/useSettings';
// import logo from '../assets/elgatama-logo.png'
// import yt from '../assets/youtube.png'
// import fb from '../assets/facebook.png'
// import ig from '../assets/instagram.png'
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
    pageWithTopSpacing: {
        fontSize: 10,
        fontFamily: 'Helvetica',
        paddingTop: 120,
        paddingRight: 40,
        paddingBottom: 40,
        paddingLeft: 40,
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
    categoryRow: {
        flexDirection: 'row',
        paddingVertical: 12,
        paddingHorizontal: 8,
        marginTop: 15,
        marginBottom: 0,
    },
    categoryText: {
        fontSize: 12,
        fontWeight: 'bold',
        color: '#2d3748',
        textTransform: 'uppercase',
        width: '100%',
        textAlign: 'left',
        letterSpacing: 0.5,
    },
    subCategory: {
        fontSize: 12,
        fontWeight: 'semibold',
        marginTop: 5,
        marginBottom: 5,
        color: '#4a5568',
    },
    tableHeader: {
        color: 'black',
        flexDirection: 'row',
        fontWeight: 'bold',
        paddingVertical: 6,
        borderRadius: 2,
        backgroundColor: '#f8f9fa',
        borderBottom: '1px solid #e2e8f0',
        marginBottom: 2,
    },
    tableRow: {
        flexDirection: 'row',
        paddingVertical: 1,
        minHeight: 10,
    },
    columnHeader: {
        width: '25%',
        paddingHorizontal: 4,
    },
    columnResult: {
        width: '10%',
        paddingHorizontal: 4,
        textAlign: 'center',
    },
    columnUnit: {
        width: '15%',
        paddingHorizontal: 4,
        textAlign: 'center',
    },
    columnStatus: {
        width: '10%',
        textAlign: 'left',
        alignItems: 'center',
    },
    columnReference: {
        width: '20%',
        paddingHorizontal: 4,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    cell: {
        paddingHorizontal: 2,
        fontSize: 9,
        textAlign: 'left',
        fontWeight: 200
    },
    gridContainer: {
        flexDirection: 'column',
    },
    row: {
        flexDirection: 'row',
        justifyContent: 'space-between',
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

// // Helper function to format birthdate
// Helper function to format birthdate (show only date part, not time)
const formatBirthdate = (birthdate?: string | null) => {
    if (!birthdate) return '-';
    const d = new Date(birthdate);
    if (isNaN(d.getTime())) return birthdate;
    // Format as DD/MM/YYYY
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    return `${day}/${month}/${year}`;
};

// Helper to extract time (HH:MM:SS) from an ISO datetime string
const formatTime = (datetime?: string | null) => {
    if (!datetime) return '-';
    const d = new Date(datetime);
    if (isNaN(d.getTime())) return datetime;
    const hh = String(d.getHours()).padStart(2, '0');
    const mm = String(d.getMinutes()).padStart(2, '0');
    const ss = String(d.getSeconds()).padStart(2, '0');
    return `${hh}:${mm}:${ss}`;
};

// // Helper function to calculate age
// const calculateAge = (birthdate: string) => {
//     const now = new Date();
//     const birthDate = new Date(birthdate);

//     let years = now.getFullYear() - birthDate.getFullYear();
//     let months = now.getMonth() - birthDate.getMonth();
//     let days = now.getDate() - birthDate.getDate();

//     // Adjust for negative months or days
//     if (months < 0 || (months === 0 && days < 0)) {
//         years--;
//         months += 12;
//     }
//     if (days < 0) {
//         const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 0);
//         days += lastMonth.getDate();
//         months--;
//     }

//     // return `${years} year(s), ${months} month(s), ${days} day(s)`;
//     return `${years} year(s)`;
// };

// Helper function to format gender
// const formatGender = (gender: string) => {
//     if (gender === 'F') return 'Female';
//     if (gender === 'M') return 'Male';
//     return '-';
// };

// Helper function to calculate age from birthdate until today
const calculateAge = (birthdate?: string | null) => {
    if (!birthdate) return '-';
    const now = new Date();
    const birth = new Date(birthdate);
    if (isNaN(birth.getTime())) return '-';

    let years = now.getFullYear() - birth.getFullYear();
    let months = now.getMonth() - birth.getMonth();
    let days = now.getDate() - birth.getDate();

    if (days < 0) {
        // borrow days from previous month
        const prevMonth = new Date(now.getFullYear(), now.getMonth(), 0);
        days += prevMonth.getDate();
        months -= 1;
    }
    if (months < 0) {
        months += 12;
        years -= 1;
    }

    if (years <= 0) {
        if (months <= 0) {
            return `${days} day(s)`;
        }
        return `${months} month(s)` + (days > 0 ? ` ${days} day(s)` : '');
    }

    return `${years} year(s)` + (months > 0 ? ` ${months} month(s)` : '');
};

// const Header = () => {
//     const [settings] = useSettings();

//     return (
//         <View style={styles.header} fixed>
//             <Image
//                 style={styles.logo}
//                 src={logo}
//             />
//             <View style={{ width: '85%' }}>
//                 <Text style={{
//                     fontSize: 24
//                 }}>{settings.company_name}</Text>
//                 <View style={{
//                     width: '100%',
//                     height: '0.2rem',
//                     backgroundColor: 'rgb(74, 186, 171)'
//                 }}>
//                 </View>
//                 <View style={styles.companyInfo}>
//                     <Text
//                         wrap={true}
//                         style={{ width: '45%' }}
//                     >{settings.company_address}</Text>
//                     <Text
//                         wrap={true}
//                         style={{ width: '45%' }}
//                     >{settings.company_contact_phone}</Text>
//                     <Text
//                         wrap={true}
//                         style={{ width: '45%' }}
//                     >{settings.company_contact_email}</Text>
//                 </View>
//             </View>
//         </View>
//     )
// };

const PatientInfo = ({ patient, workOrder }: { patient: Patient, workOrder: WorkOrder }) => (
    <View>
        <View style={styles.gridContainer}>
            {/* Single column layout - vertical */}
            <View style={{ flexDirection: 'row', marginBottom: 4, marginTop: 120 }}>
                <Text style={[styles.label, { width: 60 }]}>ID 1</Text>
                <Text style={styles.value}>: {workOrder.barcode}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>ID 2</Text>
                <Text style={styles.value}>: {patient.first_name} {patient.last_name}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Birth Date</Text>
                <Text style={styles.value}>: {formatBirthdate(patient.birthdate)}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Age</Text>
                <Text style={styles.value}>: {calculateAge(patient.birthdate)}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Seq.</Text>
                <Text style={styles.value}>: {workOrder.id || ''}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Date</Text>
                <Text style={styles.value}>: {formatBirthdate(workOrder.updated_at)}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Time</Text>
                <Text style={styles.value}>: {formatTime(workOrder.updated_at)}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Prof.</Text>
                <Text style={styles.value}>: Blood</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Asp.</Text>
                <Text style={styles.value}>: Open Tube</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Oper.</Text>
                <Text style={styles.value}>: {workOrder.analyzers?.length > 0 ? workOrder.analyzers[0].fullname : ""}</Text>
            </View>

            <View style={{ flexDirection: 'row', marginBottom: 4 }}>
                <Text style={[styles.label, { width: 60 }]}>Notes</Text>
                <Text style={styles.value}>: </Text>
            </View>
        </View>
    </View>
);

// const Footer = () => (
//     <View style={styles.footer}>
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

// Helper function to format reference range for better readability
const formatReferenceRange = (reference: string) => {
    if (!reference || reference === '-') return { value: reference, hasRange: false };

    // Check if it contains a dash (range)
    if (reference.includes(' - ')) {
        const [low, high] = reference.split(' - ');
        return { low: low.trim(), high: high.trim(), hasRange: true };
    }

    return { value: reference, hasRange: false };
};

// Helper function to sort data according to primaryOrder
const sortLabData = (data: ReportData[]) => {
    console.log('Input data:', data)
    // Define the exact order for all parameters - updated to match your requirements
    const primaryOrder = [
        'WBC', 'LYM', 'LYM%', 'MID', 'MID%', 'GRA', 'GRA%',
        'HGB', 'MCH', 'MCHC', 'RBC', 'MCV', 'HCT',
        'RDW%', 'RDWa', 'PLT', 'MPV', 'PDW%', 'PDWa',
        'PCT', 'P-LCR', 'P-LCC'
    ];

    // Create array to hold ordered parameters
    const orderedParams: (ReportData | null)[] = new Array(primaryOrder.length).fill(null);
    const remainingParams: ReportData[] = [];

    // First, place all items according to primaryOrder
    data?.forEach((item) => {
        const paramCode = (item.alias_code || item.parameter).toUpperCase().trim();
        const orderIndex = primaryOrder.findIndex(p => p.toUpperCase() === paramCode);

        console.log(`Checking parameter: "${paramCode}", found at index: ${orderIndex}`);

        if (orderIndex !== -1) {
            // If parameter exists in primaryOrder, place it in correct position
            orderedParams[orderIndex] = item;
        } else {
            // If parameter doesn't exist in primaryOrder, add to remaining
            remainingParams.push(item);
            console.log(`Parameter "${paramCode}" not found in primaryOrder, added to remaining`);
        }
    });

    // Filter out null slots and combine with remaining parameters
    const sortedPrimary = orderedParams.filter((item): item is ReportData => item !== null);
    console.log('Sorted according to primaryOrder:', sortedPrimary.map(item => item.parameter || item.alias_code));

    // Return all parameters in order: primaryOrder first, then remaining alphabetically
    const sortedRemaining = remainingParams.sort((a, b) => {
        const aName = (a.alias_code || a.parameter).toUpperCase();
        const bName = (b.alias_code || b.parameter).toUpperCase();
        return aName.localeCompare(bName);
    });

    console.log('Final order:', [...sortedPrimary, ...sortedRemaining].map(item => item.parameter || item.alias_code));

    return [...sortedPrimary, ...sortedRemaining];
};

// New function to create grouped data with category headers
const createGroupedDisplayData = (groupedData: { [category: string]: ReportData[] }) => {
    const displayData: (ReportData | { isCategory: true, category: string })[] = [];

    // Filter available categories
    const availableCategories = Object.keys(groupedData)
        .filter(category => groupedData[category] && groupedData[category].length > 0);

    // Check if Hematology exists in the data
    const hasHematology = availableCategories.some(category =>
        category.toLowerCase() === 'hematology'
    );

    // Sort categories with different logic based on Hematology presence
    const sortedCategories = availableCategories.sort((a, b) => {
        if (hasHematology) {
            // If Hematology exists, put it first
            if (a.toLowerCase() === 'hematology') return -1;
            if (b.toLowerCase() === 'hematology') return 1;
        }
        // All other categories sorted alphabetically
        return a.localeCompare(b);
    });

    console.log('Available categories from data:', sortedCategories);
    console.log('Has Hematology:', hasHematology);

    // Add each category as it appears in the data
    sortedCategories.forEach(category => {
        // Add category header
        displayData.push({ isCategory: true, category });
        // Add sorted data for this category
        displayData.push(...sortLabData(groupedData[category]));
    });

    return displayData;
}; export const ReportDocument = ({
    data,
    groupedData,
    patientData,
    workOrderData
}: {
    data: ReportData[],
    groupedData?: { [category: string]: ReportData[] },
    patientData: Patient,
    workOrderData: WorkOrder
}) => {
    console.log(data)
    // Use grouped data if available, otherwise fall back to sorted data
    console.log('ReportDocument - Available categories:', groupedData ? Object.keys(groupedData) : 'No groupedData');
    console.log('ReportDocument - groupedData:', groupedData);

    let displayData: (ReportData | { isCategory: true, category: string })[];

    if (groupedData && Object.keys(groupedData).length > 0) {
        // Use the grouped data as-is, with categories from the database
        console.log('Using grouped data with categories (from prop)');
        displayData = createGroupedDisplayData(groupedData);
    } else {
        // If groupedData not provided, try to derive grouping from `data` (if category fields exist)
        const groupedFromData = (data || []).reduce((acc: { [k: string]: ReportData[] }, item) => {
            const category = (item.category || 'Other').toString();
            if (!acc[category]) acc[category] = [];
            acc[category].push(item);
            return acc;
        }, {} as { [category: string]: ReportData[] });

        if (Object.keys(groupedFromData).length > 0) {
            // If categories are meaningful (more than 1 or first key not 'Other'), use grouped view so headers appear
            const meaningful = Object.keys(groupedFromData).some(k => k && k.toLowerCase().trim() !== 'other');
            if (meaningful) {
                console.log('Deriving grouped data from `data` (categories detected)');
                displayData = createGroupedDisplayData(groupedFromData);
            } else {
                console.log('No meaningful categories in data, falling back to sorted list');
                displayData = sortLabData(data);
            }
        } else {
            console.log('No grouped data and data is empty, using sorted (empty)');
            displayData = sortLabData(data);
        }
    }

    // Detect if displayData contains category headers (grouped) or plain list (ungrouped)
    const isGrouped = displayData.some(item => 'isCategory' in item && item.isCategory);

    // Determine the first category robustly (prefer 'hematology' if present)
    let firstCategory: string | undefined = undefined;
    if (isGrouped) {
        // look for exact 'hematology' category first
        const hemat = displayData.find(item => 'isCategory' in item && item.isCategory && (item.category || '').toLowerCase().trim() === 'hematology') as { isCategory: true, category: string } | undefined;
        if (hemat) {
            firstCategory = 'hematology';
        } else {
            const firstCatItem = displayData.find(item => 'isCategory' in item && item.isCategory) as { isCategory: true, category: string } | undefined;
            firstCategory = firstCatItem?.category?.toLowerCase().trim();
        }
    }

    console.log('Is grouped data:', isGrouped);
    console.log('Has Hematology:', firstCategory === 'hematology');
    console.log('First Category for Page 1:', firstCategory);

    // Check if there's content for page 2 (only possible when grouped)
    const hasSecondPageContent = isGrouped && displayData.some(item => {
        if ('isCategory' in item && item.isCategory) {
            return (item.category || '').toLowerCase().trim() !== firstCategory;
        } else {
            const dataItem = item as ReportData;
            return (dataItem.category || '').toLowerCase().trim() !== firstCategory;
        }
    });

    console.log('Has second page content:', hasSecondPageContent);

    return (
        <Document>
            {/* First Page - Patient Info + First Category (Hematology or first available) */}
            <Page size={"A4"} style={styles.page} wrap>
                <PatientInfo patient={patientData} workOrder={workOrderData} />

                <View style={{ marginTop: 15, marginBottom: 10 }}>
                    {(() => {
                        // Ensure we show the table header at least once on this page.
                        let tableHeaderShownPage1 = false;

                        return displayData.map((item, index) => {
                            // Check if this is a category header
                            if ('isCategory' in item && item.isCategory) {
                                // Only show the first category on first page
                                const catName = (item.category || '').toLowerCase().trim();
                                if (catName !== firstCategory) {
                                    return null;
                                }

                                // Render category header and table header
                                tableHeaderShownPage1 = true;
                                return (
                                    <View key={`category-${index}`}>
                                        {/* Category header */}
                                        <View style={styles.categoryRow}>
                                            <Text style={styles.categoryText}>{item.category}</Text>
                                        </View>

                                        {/* Table header for this category */}
                                        <View style={styles.tableHeader}>
                                            <Text style={[styles.columnHeader, styles.cell]}>Parameter</Text>
                                            <Text style={[styles.columnResult, styles.cell]}>Result</Text>
                                            <Text style={[styles.columnStatus, styles.cell]}>Status</Text>
                                            <Text style={[styles.columnUnit, styles.cell]}>Unit</Text>
                                            <Text style={[styles.columnReference, styles.cell, { textAlign: 'center' }]}>Ranges</Text>
                                        </View>
                                    </View>
                                );
                            }

                            // Regular data row - check if it belongs to first category
                            const dataItem = item as ReportData;

                            // If grouped, only show items belonging to the firstCategory on page 1.
                            // If not grouped, show everything.
                            if (isGrouped && (dataItem.category || '').toLowerCase().trim() !== firstCategory) {
                                return null;
                            }

                            // If table header hasn't been shown yet on this page (e.g. ungrouped data), render it once before first data row
                            const nodes: any[] = [];
                            if (!tableHeaderShownPage1) {
                                tableHeaderShownPage1 = true;
                                nodes.push(
                                    <View key={`table-header-${index}`} style={styles.tableHeader}>
                                        <Text style={[styles.columnHeader, styles.cell]}>Parameter</Text>
                                        <Text style={[styles.columnResult, styles.cell]}>Result</Text>
                                        <Text style={[styles.columnStatus, styles.cell]}>Status</Text>
                                        <Text style={[styles.columnUnit, styles.cell]}>Unit</Text>
                                        <Text style={[styles.columnReference, styles.cell, { textAlign: 'center' }]}>Ranges</Text>
                                    </View>
                                );
                            }

                            const isHigh = dataItem.abnormality === 'High';
                            const isLow = dataItem.abnormality === 'Low';
                            const isPositive = dataItem.abnormality === 'Positive';
                            const isNegative = dataItem.abnormality === 'Negative';
                            const abnormalStyle = {};

                            // Format the result to show proper decimals like in the image
                            const formattedResult = typeof dataItem.result === 'string'
                                ? dataItem.result
                                : String(dataItem.result || '');

                            nodes.push(
                                <View key={index} style={styles.tableRow}>
                                    <Text style={[styles.columnHeader, styles.cell, abnormalStyle]}>
                                        {dataItem.parameter}
                                    </Text>
                                    <Text style={[styles.columnResult, styles.cell, abnormalStyle]}>
                                        {formattedResult}
                                    </Text>
                                    <Text style={[styles.columnStatus, styles.cell, abnormalStyle]}>
                                        {isHigh ? 'H' : isLow ? 'L' : isPositive ? '+' : isNegative ? '-' : ''}
                                    </Text>
                                    <Text style={[styles.columnUnit, styles.cell, abnormalStyle]}>
                                        {dataItem.unit}
                                    </Text>
                                    <View style={[styles.columnReference, styles.cell]}>
                                        {(() => {
                                            const rangeData = formatReferenceRange(dataItem.reference);
                                            if ('hasRange' in rangeData && rangeData.hasRange && 'low' in rangeData && 'high' in rangeData) {
                                                return (
                                                    <>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'left', flex: 1 }]}>
                                                            {rangeData.low}
                                                        </Text>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'center', flex: 0 }]}>
                                                            -
                                                        </Text>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'right', flex: 1 }]}>
                                                            {rangeData.high}
                                                        </Text>
                                                    </>
                                                );
                                            } else {
                                                return (
                                                    <Text style={[styles.cell, abnormalStyle]}>
                                                        {'value' in rangeData ? rangeData.value : dataItem.reference}
                                                    </Text>
                                                );
                                            }
                                        })()}
                                    </View>
                                </View>
                            );

                            return nodes;
                        });
                    })()}
                </View>
            </Page>

            {/* Second Page - Other Categories (only if there's content) */}
            {hasSecondPageContent && (
                <Page size={"A4"} style={styles.pageWithTopSpacing} wrap>
                    <View style={{ marginBottom: 10 }}>
                        {displayData.map((item, index) => {
                            // Check if this is a category header
                            if ('isCategory' in item && item.isCategory) {
                                // Skip the first category on second page (already shown on page 1)
                                const catName = (item.category || '').toLowerCase().trim();
                                if (catName === firstCategory) {
                                    return null;
                                }

                                return (
                                    <View key={`category-page2-${index}`}>
                                        {/* Category header */}
                                        <View style={styles.categoryRow}>
                                            <Text style={styles.categoryText}>{item.category}</Text>
                                        </View>

                                        {/* Table header for this category */}
                                        <View style={styles.tableHeader}>
                                            <Text style={[styles.columnHeader, styles.cell]}>Parameter</Text>
                                            <Text style={[styles.columnResult, styles.cell]}>Result</Text>
                                            <Text style={[styles.columnStatus, styles.cell]}>Status</Text>
                                            <Text style={[styles.columnUnit, styles.cell]}>Unit</Text>
                                            <Text style={[styles.columnReference, styles.cell, { textAlign: 'center' }]}>Ranges</Text>
                                        </View>
                                    </View>
                                );
                            }

                            // Regular data row - check if it belongs to categories other than first category
                            const dataItem = item as ReportData;

                            // If grouped, skip first category data on second page (already shown on page 1).
                            // If not grouped, there should be no second page.
                            if (isGrouped && (dataItem.category || '').toLowerCase().trim() === firstCategory) {
                                return null;
                            }

                            const isHigh = dataItem.abnormality === 'High';
                            const isLow = dataItem.abnormality === 'Low';
                            const isPositive = dataItem.abnormality === 'Positive';
                            const isNegative = dataItem.abnormality === 'Negative';
                            const abnormalStyle = {};

                            // Format the result to show proper decimals like in the image
                            const formattedResult = typeof dataItem.result === 'string'
                                ? dataItem.result
                                : String(dataItem.result || '');

                            return (
                                <View key={`page2-${index}`} style={styles.tableRow}>
                                    <Text style={[styles.columnHeader, styles.cell, abnormalStyle]}>
                                        {dataItem.parameter}
                                    </Text>
                                    <Text style={[styles.columnResult, styles.cell, abnormalStyle]}>
                                        {formattedResult}
                                    </Text>
                                    <Text style={[styles.columnStatus, styles.cell, abnormalStyle]}>
                                        {isHigh ? 'H' : isLow ? 'L' : isPositive ? '+' : isNegative ? '-' : ''}
                                    </Text>
                                    <Text style={[styles.columnUnit, styles.cell, abnormalStyle]}>
                                        {dataItem.unit}
                                    </Text>
                                    <View style={[styles.columnReference, styles.cell]}>
                                        {(() => {
                                            const rangeData = formatReferenceRange(dataItem.reference);
                                            if ('hasRange' in rangeData && rangeData.hasRange && 'low' in rangeData && 'high' in rangeData) {
                                                return (
                                                    <>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'left', flex: 1 }]}>
                                                            {rangeData.low}
                                                        </Text>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'center', flex: 0 }]}>
                                                            -
                                                        </Text>
                                                        <Text style={[styles.cell, abnormalStyle, { textAlign: 'right', flex: 1 }]}>
                                                            {rangeData.high}
                                                        </Text>
                                                    </>
                                                );
                                            } else {
                                                return (
                                                    <Text style={[styles.cell, abnormalStyle]}>
                                                        {'value' in rangeData ? rangeData.value : dataItem.reference}
                                                    </Text>
                                                );
                                            }
                                        })()}
                                    </View>
                                </View>
                            );
                        })}
                    </View>
                </Page>
            )}
        </Document >
    );
};

// Usage example remains the same