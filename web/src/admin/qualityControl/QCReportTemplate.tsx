import {
    Document,
    Page,
    Text,
    View,
    StyleSheet,
    Font,
    Image,
} from '@react-pdf/renderer';

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
        paddingTop: 40,
        paddingLeft: 40,
        paddingRight: 40,
        paddingBottom: 60,
        backgroundColor: 'white',
    },
    header: {
        marginBottom: 20,
        textAlign: 'center',
        borderBottom: '2px solid #4abaab',
        paddingBottom: 10,
    },
    title: {
        fontSize: 20,
        fontWeight: 'bold',
        color: '#4abaab',
        marginBottom: 3,
    },
    subtitle: {
        fontSize: 14,
        fontWeight: 'bold',
        marginBottom: 3,
    },

    generatedText: {
        fontSize: 8,
        color: '#666666',
    },
    section: {
        marginBottom: 15,
    },
    sectionTitle: {
        fontSize: 12,
        fontWeight: 'bold',
        marginBottom: 8,
        color: '#2d3748',
    },
    infoBox: {
        backgroundColor: '#f5f5f5',
        padding: 10,
        borderRadius: 2,
        marginBottom: 15,
    },
    infoGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
    },
    infoItem: {
        width: '50%',
        marginBottom: 5,
    },
    infoLabel: {
        fontSize: 8,
        color: '#666666',
        marginBottom: 2,
    },
    infoValue: {
        fontSize: 9,
        fontWeight: 'bold',
    },
    statBox: {
        border: '1px solid #e0e0e0',
        borderRadius: 2,
        padding: 10,
        marginBottom: 10,
    },
    statHeader: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 8,
    },
    levelChip: {
        padding: '3 8',
        borderRadius: 2,
        marginRight: 8,
    },
    levelChipText: {
        fontSize: 8,
        fontWeight: 'bold',
        color: 'white',
    },
    statInfo: {
        fontSize: 8,
    },
    statGrid: {
        flexDirection: 'row',
        gap: 10,
    },
    statColumn: {
        flex: 1,
        padding: 8,
        borderRadius: 2,
    },
    statColumnTitle: {
        fontSize: 9,
        fontWeight: 'bold',
        marginBottom: 5,
    },
    statRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        marginBottom: 3,
    },
    statLabel: {
        fontSize: 7,
        color: '#666666',
    },
    statValue: {
        fontSize: 8,
        fontWeight: 'bold',
    },
    chartImage: {
        width: '100%',
        marginTop: 10,
        marginBottom: 10,
    },
    tableHeader: {
        flexDirection: 'row',
        backgroundColor: '#f5f5f5',
        padding: 4,
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
        fontWeight: 'bold',
    },
    tableRow: {
        flexDirection: 'row',
        borderBottomWidth: 1,
        borderBottomColor: '#e2e8f0',
        paddingTop: 8,
        paddingBottom: 8,
        paddingLeft: 4,
        paddingRight: 4,
    },
    tableCell: {
        fontSize: 8,
        paddingLeft: 4,
        paddingRight: 4,
    },
    footer: {
        position: 'absolute',
        bottom: 10,
        left: 40,
        right: 40,
        textAlign: 'center',
        fontSize: 8,
        color: '#666666',
        borderTop: '1px solid #e0e0e0',
        paddingTop: 10,
    },
});

interface ChartDataPoint {
    index: number;
    date: string;
    result: number;
    level: number;
    errorSD: number;
    mean: number;
    sd: number;
    lotNumber: string;
    absoluteError: number;
    relativeError: number;
    status: string;
    source: string;
}

interface QCReportProps {
    testType: {
        code: string;
        name: string;
        unit: string;
    };
    device?: {
        id: number;
        name: string;
        serial_number: string;
    };
    entries: Array<{
        id: number;
        qc_level: number;
        lot_number: string;
        target_mean: number;
        target_sd?: number;
        ref_min: number;
        ref_max: number;
    }>;
    qcResults: Array<any>;
    selectedMethod: string;
    startDate: Date | null;
    endDate: Date | null;
    chartData: ChartDataPoint[];
    chartImageUrl?: string;
}

// Helper function to format dates
const formatDate = (date: Date | null) => {
    if (!date) return 'N/A';
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
    });
};

export const QCReportTemplate = (props: QCReportProps) => {
    const { testType, device, entries, qcResults, selectedMethod, startDate, endDate, chartData, chartImageUrl } = props;

    return (
        <Document>
            <Page size="A4" style={styles.page}>
                {/* Header */}
                <View style={styles.header} fixed>
                    <Text style={styles.title}>Quality Control Report</Text>
                    <Text style={styles.subtitle}>{testType.code} - {testType.name}</Text>
                    {device && <Text style={styles.generatedText}>Device: {device.name}</Text>}
                    <Text style={styles.generatedText}>Generated: {new Date().toLocaleString()}</Text>
                </View>

                {/* Report Information */}
                <View style={styles.infoBox}>
                    <Text style={styles.sectionTitle}>Report Parameters</Text>
                    <View style={styles.infoGrid}>
                        <View style={styles.infoItem}>
                            <Text style={styles.infoLabel}>Test Unit:</Text>
                            <Text style={styles.infoValue}>{testType.unit}</Text>
                        </View>
                        <View style={styles.infoItem}>
                            <Text style={styles.infoLabel}>Calculation Method:</Text>
                            <Text style={styles.infoValue}>{selectedMethod}</Text>
                        </View>
                        <View style={styles.infoItem}>
                            <Text style={styles.infoLabel}>Date Range:</Text>
                            <Text style={styles.infoValue}>{formatDate(startDate)} - {formatDate(endDate)}</Text>
                        </View>
                        <View style={styles.infoItem}>
                            <Text style={styles.infoLabel}>Total Data Points:</Text>
                            <Text style={styles.infoValue}>{chartData.length}</Text>
                        </View>
                    </View>
                </View>

                {/* QC Entry Statistics */}
                <View style={styles.section}>
                    <Text style={styles.sectionTitle}>QC Entry Statistics</Text>
                    {entries.filter(e => {
                        const entryResults = qcResults.filter(r => r.qc_entry_id === e.id);
                        return entryResults.length > 0;
                    }).map(entry => {
                        const entryResults = qcResults.filter(r => r.qc_entry_id === entry.id);
                        const resultCount = entryResults.length;
                        const hasCalculated = resultCount >= 5;

                        let calculatedMean = 0;
                        let calculatedSD = 0;
                        let calculatedCV = 0;

                        if (hasCalculated && entryResults.length > 0) {
                            const latestResult = entryResults[0];
                            calculatedMean = latestResult.calculated_mean;
                            calculatedSD = latestResult.calculated_sd;
                            calculatedCV = latestResult.calculated_cv;
                        }

                        const quotedCV = entry.target_sd && entry.target_mean
                            ? (entry.target_sd / entry.target_mean) * 100
                            : 0;

                        const levelColor = entry.qc_level === 1 ? '#2196f3' : entry.qc_level === 2 ? '#9c27b0' : '#ff9800';

                        return (
                            <View key={entry.id} style={styles.statBox}>
                                <View style={styles.statHeader}>
                                    <View style={[styles.levelChip, { backgroundColor: levelColor }]}>
                                        <Text style={styles.levelChipText}>Level {entry.qc_level}</Text>
                                    </View>
                                    <Text style={styles.statInfo}>
                                        Lot: {entry.lot_number} | Count: {resultCount} | Ref: {entry.ref_min?.toFixed ? entry.ref_min.toFixed(2) : entry.ref_min} - {entry.ref_max?.toFixed ? entry.ref_max.toFixed(2) : entry.ref_max}
                                    </Text>
                                </View>

                                <View style={styles.statGrid}>
                                    <View style={[styles.statColumn, { backgroundColor: 'rgba(33, 150, 243, 0.05)' }]}>
                                        <Text style={[styles.statColumnTitle, { color: '#2196f3' }]}>Quoted (Datasheet)</Text>
                                        <View style={styles.statRow}>
                                            <Text style={styles.statLabel}>Mean:</Text>
                                            <Text style={styles.statValue}>{entry.target_mean.toFixed(2)}</Text>
                                        </View>
                                        <View style={styles.statRow}>
                                            <Text style={styles.statLabel}>SD:</Text>
                                            <Text style={styles.statValue}>{entry.target_sd?.toFixed(2) || '-'}</Text>
                                        </View>
                                        <View style={styles.statRow}>
                                            <Text style={styles.statLabel}>CV (%):</Text>
                                            <Text style={styles.statValue}>{quotedCV > 0 ? quotedCV.toFixed(2) : '-'}</Text>
                                        </View>
                                    </View>

                                    <View style={[styles.statColumn, { backgroundColor: 'rgba(76, 175, 80, 0.05)' }]}>
                                        <Text style={[styles.statColumnTitle, { color: '#4caf50' }]}>Calculated (System)</Text>
                                        {hasCalculated ? (
                                            <>
                                                <View style={styles.statRow}>
                                                    <Text style={styles.statLabel}>Mean:</Text>
                                                    <Text style={styles.statValue}>{calculatedMean.toFixed(2)}</Text>
                                                </View>
                                                <View style={styles.statRow}>
                                                    <Text style={styles.statLabel}>SD:</Text>
                                                    <Text style={styles.statValue}>{calculatedSD.toFixed(2)}</Text>
                                                </View>
                                                <View style={styles.statRow}>
                                                    <Text style={styles.statLabel}>CV (%):</Text>
                                                    <Text style={styles.statValue}>{calculatedCV.toFixed(2)}</Text>
                                                </View>
                                            </>
                                        ) : (
                                            <Text style={styles.statLabel}>Minimum 5 data points required</Text>
                                        )}
                                    </View>
                                </View>
                            </View>
                        );
                    })}
                </View>

                {/* Levey-Jennings Chart */}
                {chartImageUrl && (
                    <View style={styles.section}>
                        <Text style={styles.sectionTitle}>Levey-Jennings Chart</Text>
                        <Image src={chartImageUrl} style={styles.chartImage} />
                    </View>
                )}

                {/* Footer */}
                <View style={styles.footer} fixed>
                    <Text>This is an automated report generated by TAMALabs LIMS</Text>
                </View>
            </Page>

            {/* QC History Table - New Page */}
            <Page size="A4" style={styles.page}>
                <View style={styles.header} fixed>
                    <Text style={styles.title}>Quality Control Report</Text>
                    <Text style={styles.subtitle}>{testType.code} - {testType.name}</Text>
                    {device && <Text style={styles.generatedText}>Device: {device.name}</Text>}
                </View>

                <View style={styles.section}>
                    <Text style={styles.sectionTitle}>QC Results History</Text>

                    {/* Table Header */}
                    <View style={styles.tableHeader}>
                        <Text style={[styles.tableCell, { width: '22%' }]}>Date</Text>
                        <Text style={[styles.tableCell, { width: '8%' }]}>Level</Text>
                        <Text style={[styles.tableCell, { width: '10%' }]}>Lot</Text>
                        <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>Value</Text>
                        <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>Abs Err</Text>
                        <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>Rel Err (%)</Text>
                        <Text style={[styles.tableCell, { width: '24%', paddingLeft: 15, }]}>Status</Text>
                        <Text style={[styles.tableCell, { width: '13%' }]}>Source</Text>
                    </View>

                    {/* Table Rows */}
                    {chartData.map((data, idx) => {
                        const levelColor = data.level === 1 ? '#2196f3' : data.level === 2 ? '#9c27b0' : '#ff9800';

                        let statusColor = '#4caf50';
                        if (data.status === 'Reject') statusColor = '#f44336';
                        else if (data.status === 'Warning') statusColor = '#ff9800';

                        return (
                            <View key={idx} style={styles.tableRow}>
                                <Text style={[styles.tableCell, { width: '22%' }]}>{data.date}</Text>
                                <Text style={[styles.tableCell, { width: '8%', color: levelColor, fontWeight: 'bold' }]}>L{data.level}</Text>
                                <Text style={[styles.tableCell, { width: '10%' }]}>{data.lotNumber}</Text>
                                <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>{data.result.toFixed(2)}</Text>
                                <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>{data.absoluteError.toFixed(2)}</Text>
                                <Text style={[styles.tableCell, { width: '13%', textAlign: 'right' }]}>{data.relativeError.toFixed(2)}</Text>
                                <Text style={[styles.tableCell, { width: '24%', paddingLeft: 15, color: statusColor, fontWeight: 'bold' }]}>{data.status}</Text>
                                <Text style={[styles.tableCell, { width: '13%' }]}>{data.source}</Text>
                            </View>
                        );
                    })}
                </View>

                {/* Footer */}
                <View style={styles.footer} fixed>
                    <Text>This is an automated report generated by TAMALabs</Text>
                    <Text style={{ marginTop: 2 }}>© {new Date().getFullYear()} PT ELGA TAMA. All rights reserved.</Text>
                </View>
            </Page>
        </Document>
    );
};
