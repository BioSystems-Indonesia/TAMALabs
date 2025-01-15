import React from 'react';
import {
    Page,
    Text,
    View,
    Document,
    StyleSheet,
    BlobProvider,
    Font,
} from '@react-pdf/renderer';

// Optional: Register custom fonts if required
Font.register({
    family: 'Roboto',
    src: 'https://fonts.gstatic.com/s/roboto/v27/KFOmCnqEu92Fr1Me5Q.ttf',
});

// Define data types for clarity
interface ReportData {
    parameter: string;
    result: string;
    reference: string;
}

interface MCUReportProps {
    data: ReportData[];
}

// Styles for PDF layout
const styles = StyleSheet.create({
    page: {
        padding: 30,
        fontFamily: 'Roboto',
    },
    header: {
        fontSize: 20,
        textAlign: 'center',
        marginBottom: 20,
    },
    section: {
        margin: 10,
        padding: 10,
        borderBottom: '1px solid #ccc',
    },
    label: {
        fontSize: 12,
        fontWeight: 'bold',
    },
    value: {
        fontSize: 12,
    },
});

// PDF Document Component
const MCUReport: React.FC<MCUReportProps> = ({ data }) => (
    <Document>
        <Page size="A4" style={styles.page}>
            <Text style={styles.header}>MCU Result</Text>
            {data.map((item, index) => (
                <View key={index} style={styles.section}>
                    <Text style={styles.label}>Parameter: {item.parameter}</Text>
                    <Text style={styles.value}>Result: {item.result}</Text>
                    <Text style={styles.value}>Reference: {item.reference}</Text>
                </View>
            ))}
        </Page>
    </Document>
);

const PrintMCU: React.FC = () => {
    const mockData: ReportData[] = [
        { parameter: 'Blood Pressure', result: '120/80', reference: 'Normal' },
        { parameter: 'Heart Rate', result: '72 bpm', reference: 'Normal' },
        { parameter: 'Cholesterol', result: '190 mg/dL', reference: 'Desirable' },
    ];

    return (
        <div style={{ textAlign: 'center', marginTop: '50px' }}>
            <h1>Generate Your MCU Report</h1>
            <BlobProvider document={<MCUReport data={mockData} />}>
                {({ url, loading, error }) => {
                    if (loading) {
                        return (
                            <button
                                disabled
                                style={{
                                    padding: '10px',
                                    fontSize: '14px',
                                    cursor: 'not-allowed',
                                    backgroundColor: '#ccc',
                                    border: 'none',
                                    borderRadius: '5px',
                                }}
                            >
                                Generating PDF...
                            </button>
                        );
                    }

                    if (error) {
                        return <span>Error generating PDF: {error.message}</span>;
                    }

                    return (
                        <div>
                            {/* Download PDF Button */}
                            <a
                                href={url || ''}
                                download="MCU_Result.pdf"
                                style={{
                                    display: 'inline-block',
                                    padding: '10px',
                                    fontSize: '14px',
                                    backgroundColor: '#007bff',
                                    color: '#fff',
                                    textDecoration: 'none',
                                    borderRadius: '5px',
                                    marginRight: '10px',
                                    cursor: 'pointer',
                                }}
                            >
                                Download PDF
                            </a>

                            {/* Print PDF Button */}
                            <button
                                onClick={() => {
                                    if (url) {
                                        window.open(url, '_blank')?.focus();
                                    }
                                }}
                                style={{
                                    padding: '10px',
                                    fontSize: '14px',
                                    backgroundColor: '#28a745',
                                    color: '#fff',
                                    border: 'none',
                                    borderRadius: '5px',
                                    cursor: 'pointer',
                                }}
                            >
                                Print PDF
                            </button>
                        </div>
                    );
                }}
            </BlobProvider>
        </div>
    );
};

export default PrintMCU;