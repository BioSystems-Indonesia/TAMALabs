import { useEffect, useState } from 'react';
import {
    useList,
    useNotify,
    ListContextProvider,
    DataTable,
    WithRecord,
    AutocompleteInput,
    ReferenceInput,
    Form,

} from 'react-admin';
import { Card, Stack, Typography, Box, Tooltip, Chip } from '@mui/material';
import { PersonSearch } from '@mui/icons-material';
import useAxios from "../../hooks/useAxios.ts";
import PatientInfoCard from './PatientInfoCard.tsx';
import type {
    RawTestResult,
    PatientResultHistoryResponse,
    ProcessedTestResult,
    AbnormalColor,
    GroupedByDate,
    EGFRCalculation
} from '../../types/patient_result_history';
import { AbnormalFlag } from '../../types/patient_result_history';

const obj: RawTestResult = {
    abnormal: 2,
    category: "Hematology",
    created_at: "2025-08-22T09:27:10+07:00",
    formatted_result: 2,
    history: null,
    id: 99,
    picked: true,
    reference_range: "14 - 18",
    result: 2,
    specimen_id: 11,
    test: "HGB",
    test_type_id: 2,
    unit: "g/dL",
    egfr: {
        value: 85,
        category: "eGFR",
        formula: "CKD-EPI",
        unit: "mL/min/1.73m²"
    },
};

const initialData: RawTestResult[] = [obj];

const ResultHistory = () => {
    const listContext = useList({ data: initialData });
    const axios = useAxios()
    const notify = useNotify();
    const [data, setData] = useState<ProcessedTestResult[]>([]);
    const [allDates, setAllDates] = useState<string[]>([]);

    const [patientID, setPatientID] = useState<number>(0);

    useEffect(() => {
        const patientID = localStorage.getItem("patient_history_patient_id");
        if (patientID) {
            setPatientID(parseInt(patientID));
        }
    }, []);

    useEffect(() => {
        if (!patientID) {
            return;
        }
        localStorage.setItem("patient_history_patient_id", patientID.toString());

        axios.get<PatientResultHistoryResponse>(`patient/${patientID}/result/history?start_date=2025-01-01T00:00:00Z`).then((res) => {
            const uniqueDates = getUniqueDates(res.data.test_result.map(extractDateFromTestResult));
            setAllDates(uniqueDates);

            // Group by category first, then by test within each category
            const testResultsGroupedByCategory = groupTestResultsByTestName(res.data.test_result, (i: RawTestResult) => i.category);
            const processedTestResults = transformTestResultsGroupedByCategory(testResultsGroupedByCategory);
            setData(processedTestResults)
        }).catch((err) => {
            console.error(err)
            notify("Error processing data", {
                type: 'error',
            });
        })
    }, [axios, notify, patientID]);

    const resultCol = (d: string, params: ProcessedTestResult) => {
        const displayValue = params[d + "_result"];
        const egfr = params[d + "_egfr"] as EGFRCalculation;
        return (
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', gap: 0.5 }}>
                <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
                    {displayValue?.toLocaleString() || '-'}
                </Typography>

                {egfr && (
                    <Tooltip
                        title={`${egfr.category} (${egfr.formula})`}
                        arrow
                        placement="top"
                    >
                        <Chip
                            label={`eGFR: ${egfr.value.toFixed(1)}`}
                            size="small"
                            variant="filled"
                            color={
                                egfr.value >= 90 ? 'success' :
                                    egfr.value >= 60 ? 'info' :
                                        egfr.value >= 45 ? 'warning' : 'error'
                            }
                            sx={{
                                fontSize: '0.7rem',
                                height: '18px',
                                cursor: 'help',
                                '& .MuiChip-label': {
                                    px: 0.75
                                }
                            }}
                        />
                    </Tooltip>
                )}
            </Box>
        );
    }

    return (
        <Stack gap={2}>
            <Form>
                <ReferenceInput
                    source={"patient_ids"}
                    reference="patient"
                    label={"Patient"}
                    alwaysOn
                    sx={{
                        '& .MuiInputLabel-root': {
                            // color: theme.palette.text.primary,
                            fontWeight: 500,
                            fontSize: '0.9rem',
                        }
                    }}
                >
                    <AutocompleteInput
                        onChange={(value) => {
                            setPatientID(value)
                        }}
                        size="small"
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                // backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                borderRadius: '12px',
                                transition: 'all 0.3s ease',
                                '&:hover': {
                                    // backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                },
                            }
                        }}
                    />
                </ReferenceInput>
            </Form>

            {!patientID ? (
                <Card sx={{ p: 4, textAlign: 'center', mt: 2 }}>
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
                        <PersonSearch sx={{ fontSize: 48, color: 'text.secondary' }} />
                        <Typography variant="h6" color="textSecondary">
                            Please select a patient to view their test results
                        </Typography>
                    </Box>
                </Card>
            ) : (
                <>
                    <PatientInfoCard patientId={patientID} />
                    <ListContextProvider value={{ ...listContext, data, total: data.length, isPending: false }}>
                        <Card>
                            <DataTable resource="actors" bulkActionButtons={false}>
                                <DataTable.Col source="test">
                                    <WithRecord render={(record: ProcessedTestResult) => (
                                        <Typography
                                            sx={{
                                                fontWeight: record.isCategory ? 700 : 400,
                                                fontSize: record.isCategory ? '1.1rem' : '0.875rem',
                                                color: record.isCategory ? 'primary.main' : 'text.primary',
                                                backgroundColor: record.isCategory ? 'action.hover' : 'transparent',
                                                padding: record.isCategory ? '8px 16px' : '0',
                                                borderRadius: record.isCategory ? 1 : 0,
                                                marginY: record.isCategory ? 1 : 0
                                            }}
                                        >
                                            {record.isCategory ? `📊 ${record.test}` : record.test}
                                        </Typography>
                                    )} />
                                </DataTable.Col>
                                <DataTable.Col source="reference_range">
                                    <WithRecord render={(record: ProcessedTestResult) => (
                                        <Typography sx={{ opacity: record.isCategory ? 0 : 1 }}>
                                            {record.isCategory ? '' : record.reference_range}
                                        </Typography>
                                    )} />
                                </DataTable.Col>
                                <DataTable.Col source="unit">
                                    <WithRecord render={(record: ProcessedTestResult) => (
                                        <Typography sx={{ opacity: record.isCategory ? 0 : 1 }}>
                                            {record.isCategory ? '' : record.unit}
                                        </Typography>
                                    )} />
                                </DataTable.Col>
                                {allDates.map(d =>
                                    <DataTable.Col key={d} label={d}>
                                        <WithRecord label={d} render={(record: ProcessedTestResult) => (
                                            <Typography
                                                color={record.isCategory ? 'transparent' : (record[d + "_color"] as AbnormalColor)}
                                                sx={{ opacity: record.isCategory ? 0 : 1 }}
                                            >
                                                {record.isCategory ? '' :
                                                    record[`${d}_result`] != null ? resultCol(d, record) : ''}
                                            </Typography>
                                        )} />
                                    </DataTable.Col>
                                )}
                            </DataTable>
                        </Card>
                    </ListContextProvider>
                </>
            )}
        </Stack>
    );
};

export default ResultHistory;

function groupTestResultsByTestName<T>(a: T[], fn: (item: T) => string): Record<string, T[]> {
    return a.reduce((acc: Record<string, T[]>, obj: T) => {
        const keyValue = fn(obj);
        if (!acc[keyValue]) {
            acc[keyValue] = [];
        }
        acc[keyValue].push(obj);
        return acc;
    }, {});
}

function convertGroupedObjectToArrays<T>(object: Record<string, T[]>): T[][] {
    return Object.entries(object).map(([key, value]) => value);
}

function transformTestResultsGroupedByCategory(categoryGroups: Record<string, RawTestResult[]>): ProcessedTestResult[] {
    const result: ProcessedTestResult[] = [];

    // Sort categories alphabetically
    const sortedCategories = Object.keys(categoryGroups).sort();

    for (const category of sortedCategories) {
        // Add category header row
        const categoryHeader: ProcessedTestResult = {
            test: category,
            reference_range: '',
            unit: '',
            category: category,
            isCategory: true,
        };
        result.push(categoryHeader);

        // Group tests within this category by test name
        const testsInCategory = categoryGroups[category];
        const testResultsGroupedByTest = groupTestResultsByTestName(testsInCategory, (i: RawTestResult) => i.test);
        const testResultArrays = convertGroupedObjectToArrays(testResultsGroupedByTest);
        const processedTests = transformTestResultsForDataTable(testResultArrays);

        // Add all tests in this category
        result.push(...processedTests);
    }

    return result;
}

function transformTestResultsForDataTable(ar: RawTestResult[][]): ProcessedTestResult[] {
    return ar.map((li: RawTestResult[]) => {
        const dates: GroupedByDate = groupTestResultsByTestName(li, extractDateFromTestResult);

        const result: ProcessedTestResult = {
            test: li[0].test,
            reference_range: li[0].reference_range,
            unit: li[0].unit,
            category: li[0].category,
        }

        for (const [date, items] of Object.entries(dates)) {
            // TODO make it to the last? maybe?
            const i = items as RawTestResult[];

            let color: AbnormalColor = "success";
            switch (i[0].abnormal) {
                case AbnormalFlag.Normal: color = "default"; break;
                case AbnormalFlag.High: color = "error"; break;
                case AbnormalFlag.Low: color = "secondary"; break;
                case AbnormalFlag.Critical: color = "default"; break;
                default: color = "success";
            }

            result[date + "_result"] = i[0].result;
            result[date + "_color"] = color;
            result[date + "_egfr"] = i[0].egfr;
        }
        return result;
    });

}

function extractDateFromTestResult(d: RawTestResult): string {
    return new Date(d.created_at).toISOString().substring(0, 10)
}

function getUniqueDates(a: string[]) {
    return Array.from(new Set(a));
}
