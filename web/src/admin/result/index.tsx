import RefreshIcon from '@mui/icons-material/Refresh';
import WarningAmberIcon from '@mui/icons-material/WarningAmber';
import { Box, Chip, Button as MUIButton, Stack, Typography, useTheme } from "@mui/material";
import dayjs from "dayjs";
import {
    AutocompleteArrayInput,
    BooleanInput,
    Button,
    Datagrid,
    DateField,
    FilterLiveForm,
    Link,
    List,
    NumberField,
    ReferenceInput,
    TopToolbar,
    useNotify,
    WithRecord
} from "react-admin";
import CustomDateInput from "../../component/CustomDateInput";
import PrintReportButton from "../../component/PrintReport";
import SideFilter from "../../component/SideFilter";
import useAxios from "../../hooks/useAxios";
import type { WorkOrder } from "../../types/work_order";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { FilledPercentChip, VerifiedChip } from "./component";


export const ResultList = () => (
    <List
        resource="result"
        sort={{ field: "id", order: "DESC" }}
        aside={<ResultSideFilter />}
        filters={ResultMoreFilter}
        filterDefaultValues={{
            created_at_start: dayjs().subtract(7, "day").toISOString(),
            created_at_end: dayjs().toISOString(),
        }}
        actions={<ResultActions />}
        exporter={false}
        storeKey={false}
        sx={{
            '& .RaList-main': {
                marginTop: '-14px'
            },
            '& .RaList-content': {
                backgroundColor: 'background.paper',
                padding: 2,
                borderRadius: 1,
            },
        }}
    >
        <ResultDataGrid />
    </List>
);

function ResultActions() {
    const axios = useAxios()
    const notify = useNotify()
    return (
        <TopToolbar>
            <Button label={"Refresh"} onClick={() => {
                axios.post("/result/refresh").then(() => {
                    notify("Refresh Result Success", {
                        type: "success"
                    })
                }).catch(() => {
                    notify("Refresh Result Failed", {
                        type: "error"
                    })
                })
            }}>
                <RefreshIcon />
            </Button>
        </TopToolbar>
    )
}

export const ResultDataGrid = (props: any) => {
    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Barcode" render={(record: WorkOrder) => (
                <Typography variant="body2">{record.barcode}</Typography>
            )} />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"}
                    onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )} />
            <WithRecord label="Request" render={(record: any) => (
                <Link to={`/work-order/${record.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.id} - ${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                </Link>
            )} />
            <WithRecord label="Request" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_request}
                </Typography>
            )} />
            <WithRecord label="Result" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_result_filled}
                </Typography>
            )} />
            <WithRecord label="Filled" render={(record: WorkOrder) => (
                <FilledPercentChip percent={record.percent_complete} />
            )} />
            <WithRecord label="Verified" render={(record: WorkOrder) => (
                <VerifiedChip verified={record.verified_status !== '' ? record.verified_status : "VERIFIED"} />
            )} />
            <DateField source="created_at" showDate showTime />
            <WithRecord label="Print Result" render={(record: WorkOrder) => {
                if (record.verified_status !== "" && record.verified_status !== "VERIFIED") {
                    return (
                        <MUIButton
                            variant="contained"
                            color="warning"
                            startIcon={<WarningAmberIcon />}
                            size="small"
                            sx={{
                                textTransform: 'none',
                                fontSize: '12px',
                                whiteSpace: 'nowrap',
                                '&:hover': {
                                    backgroundColor: 'warning.main',
                                    cursor: 'default'
                                }
                            }}
                        >
                            Not verified
                        </MUIButton>
                    )
                }

                return (
                    <PrintReportButton results={record.test_result} patient={record.patient} workOrder={record} />
                )
            }} />
        </Datagrid>
    )
}

function ResultSideFilter() {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';
    
    return (
        <SideFilter sx={{
            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',          
        }}>
            <FilterLiveForm debounce={1500}>
                <Stack spacing={0}>
                    <Box>
                        <Typography variant="h6" sx={{ 
                            color: theme.palette.text.primary, 
                            marginBottom: 2, 
                            fontWeight: 600,
                            fontSize: '1.1rem',
                            textAlign: 'center'
                        }}>
                            🔬 Filter Lab Results
                        </Typography>
                    </Box>
                    <ReferenceInput 
                        source={"patient_ids"} 
                        reference="patient" 
                        label={"Patient"} 
                        alwaysOn
                        sx={{
                            '& .MuiInputLabel-root': {
                                color: theme.palette.text.primary,
                                fontWeight: 500,
                                fontSize: '0.9rem',
                            }
                        }}
                    >
                        <AutocompleteArrayInput 
                            size="small" 
                            sx={{
                                '& .MuiOutlinedInput-root': {
                                    backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                    borderRadius: '12px',
                                    transition: 'all 0.3s ease',
                                    '&:hover': {
                                        backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                       
                                    },
                                  
                                }
                            }}
                        />
                    </ReferenceInput>
                    <Box>
                        <Typography variant="body2" sx={{ 
                            color: theme.palette.text.secondary, 
                            marginBottom: 1.5,
                            fontSize: '0.85rem',
                            fontWeight: 500
                        }}>
                            📅 Date Range
                        </Typography>
                        <Stack>
                            <CustomDateInput 
                                label={"Start Date"} 
                                source="created_at_start" 
                                disableFuture 
                                alwaysOn 
                                size="small" 
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                        borderRadius: '12px',
                                        border: isDarkMode ? `1px solid ${theme.palette.divider}` : '1px solid #e5e7eb',
                                        transition: 'all 0.3s ease',
                                        '&:hover': {
                                            backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                            borderColor: isDarkMode ? theme.palette.primary.main : '#9ca3af',
                                            boxShadow: isDarkMode 
                                                ? '0 4px 12px rgba(255, 255, 255, 0.1)' 
                                                : '0 4px 12px rgba(0, 0, 0, 0.1)',
                                        },
                                        '&.Mui-focused': {
                                            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
                                            borderColor: theme.palette.primary.main,
                                            boxShadow: `0 0 0 3px ${theme.palette.primary.main}30`,
                                        }
                                    },
                                    '& .MuiInputLabel-root': {
                                        color: theme.palette.text.primary,
                                        fontWeight: 500,
                                        fontSize: '0.85rem',
                                    }
                                }} 
                            />
                            <CustomDateInput 
                                label={"End Date"} 
                                source="created_at_end" 
                                disableFuture 
                                alwaysOn 
                                size="small" 
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                        borderRadius: '12px',
                                        border: isDarkMode ? `1px solid ${theme.palette.divider}` : '1px solid #e5e7eb',
                                        transition: 'all 0.3s ease',
                                        '&:hover': {
                                            backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                            borderColor: isDarkMode ? theme.palette.primary.main : '#9ca3af',
                                            boxShadow: isDarkMode 
                                                ? '0 4px 12px rgba(255, 255, 255, 0.1)' 
                                                : '0 4px 12px rgba(0, 0, 0, 0.1)',
                                        },
                                        '&.Mui-focused': {
                                            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
                                            borderColor: theme.palette.primary.main,
                                            boxShadow: `0 0 0 3px ${theme.palette.primary.main}30`,
                                        }
                                    },
                                    '& .MuiInputLabel-root': {
                                        color: theme.palette.text.primary,
                                        fontWeight: 500,
                                        fontSize: '0.85rem',
                                    }
                                }} 
                            />
                        </Stack>
                    </Box>
                </Stack>
            </FilterLiveForm>
        </SideFilter>
    )
}


const ResultMoreFilter = [
    <BooleanInput source={"has_result"} label={"Show Only With Result"} />,
]
