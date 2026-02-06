import RefreshIcon from '@mui/icons-material/Refresh';
import SyncIcon from '@mui/icons-material/Sync';
import SendIcon from '@mui/icons-material/Send';
import ScienceIcon from '@mui/icons-material/Science';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PrintIcon from '@mui/icons-material/Print';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import PendingIcon from '@mui/icons-material/Pending';
import CloudDoneIcon from '@mui/icons-material/CloudDone';
import CloudOffIcon from '@mui/icons-material/CloudOff';
import CloudQueueIcon from '@mui/icons-material/CloudQueue';
import { Box, Chip, Button as MUIButton, Stack, Typography, useTheme, CircularProgress, Menu, MenuItem, IconButton, Tooltip } from "@mui/material";
import { useState } from "react";
import {
    AutocompleteArrayInput,
    BooleanInput,
    Button,
    Datagrid,
    FilterLiveForm,
    Link,
    List,
    NumberField,
    ReferenceInput,
    TopToolbar,
    useNotify,
    WithRecord,
    useListContext,
    useRefresh
} from "react-admin";
import { useMutation } from "@tanstack/react-query";
import CustomDateInput from "../../component/CustomDateInput";
import GenerateReportButton from "../../component/GenerateReportButton";
import SideFilter from "../../component/SideFilter";
import useAxios from "../../hooks/useAxios";
import type { WorkOrder } from "../../types/work_order";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { FilledPercentChip } from "./component";


export const ResultList = () => (
    <List
        resource="result"
        sort={{ field: "id", order: "DESC" }}
        aside={<ResultSideFilter />}
        filters={ResultMoreFilter}
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

function SyncAllResultButton() {
    const notify = useNotify();
    const refresh = useRefresh();
    const axios = useAxios();

    const { mutate: syncAllResult, isPending } = useMutation({
        mutationFn: async () => {
            try {
                const response = await axios.post('/external/sync-all-results');
                if (!response || response.status !== 200) {
                    throw new Error(response?.data?.error || 'Failed to sync results');
                }
                return response.data;
            } catch (error: any) {
                // Handle axios errors or network errors
                if (error.response) {
                    // Server responded with error status
                    throw new Error(error.response.data?.error || `Server error: ${error.response.status}`);
                } else if (error.request) {
                    // Network error
                    throw new Error('Network error: Unable to connect to server');
                } else {
                    // Other error
                    throw new Error(error.message || 'Unknown error occurred');
                }
            }
        },
        onSuccess: () => {
            notify('Successfully synced all results to external systems', {
                type: 'success',
            });
            refresh();
        },
        onError: (error) => {
            notify(`Sync failed: ${error.message}`, {
                type: 'error',
            });
        },
    });

    return (
        <Button
            label="Sync All Result"
            onClick={() => syncAllResult()}
            disabled={isPending}
            sx={{
                backgroundColor: 'primary.main',
                color: 'white',
                '&:hover': {
                    backgroundColor: 'secondary.dark',
                },
                '&:disabled': {
                    backgroundColor: 'action.disabled',
                },
            }}
        >
            {isPending ? (
                <CircularProgress size={16} sx={{ color: 'white' }} />
            ) : (
                <SyncIcon />
            )}
        </Button>
    );
}

// ActionMenuButton component - dropdown menu for Generate Report and Send to SIMRS
function ActionMenuButton({ record, currentGeneratedId, onGenerate }: { record: WorkOrder; currentGeneratedId: string | null; onGenerate: (id: string) => void }) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const axios = useAxios();
    const notify = useNotify();
    const refresh = useRefresh();

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        event.stopPropagation();
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const mutation = useMutation({
        mutationFn: async () => {
            const response = await axios.post(`/nuha-simrs/send-result/${record.id}`);
            return response.data;
        },
        onSuccess: (data) => {
            notify(data?.message || 'Successfully sent results to Nuha SIMRS', { type: 'success' });
            refresh();
            handleClose();
        },
        onError: (error: any) => {
            console.error('Send to Nuha error:', error);

            if (error.response) {
                const status = error.response.status;
                const errorMsg = error.response.data?.error || error.response.data?.message || 'Unknown error';

                if (status === 403) {
                    notify('⚠️ Nuha SIMRS integration is not enabled!\n\nPlease enable it in: Settings → Config → SIMRS Bridging → Select "Nuha SIMRS"', {
                        type: 'warning',
                        autoHideDuration: 8000,
                        multiLine: true,
                    });
                } else {
                    notify(`❌ Failed to send results: ${errorMsg}`, {
                        type: 'error',
                        autoHideDuration: 6000,
                    });
                }
            } else if (error.request) {
                notify('❌ Network error: Unable to connect to server. Please check if the server is running.', {
                    type: 'error',
                    autoHideDuration: 8000,
                });
            } else {
                notify(`❌ Failed to send results: ${error.message || 'Unknown error'}`, {
                    type: 'error',
                    autoHideDuration: 6000,
                });
            }
            handleClose();
        },
    });

    const handleSendToSimrs = () => {
        if (record.total_result_filled === 0) {
            notify('⚠️ No results to send. Please fill in test results first.', {
                type: 'warning',
                autoHideDuration: 5000,
            });
            handleClose();
            return;
        }
        mutation.mutate();
    };

    const hasResults = record.total_result_filled > 0;
    const isFromNuha = record.barcode_simrs && record.barcode_simrs.startsWith('NUHA-');

    return (
        <>
            <Tooltip title="Actions Menu" arrow>
                <IconButton
                    size="small"
                    onClick={handleClick}
                    sx={{
                        backgroundColor: 'primary.main',
                        color: 'white',
                        borderRadius: '8px',
                        width: 32,
                        height: 32,
                        '&:hover': {
                            backgroundColor: 'primary.dark',
                            transform: 'scale(1.1)',
                            transition: 'all 0.2s',
                        },
                    }}
                >
                    <MoreVertIcon fontSize="small" />
                </IconButton>
            </Tooltip>
            <Menu
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                onClick={(e) => e.stopPropagation()}
                PaperProps={{
                    sx: {
                        minWidth: 200,
                    }
                }}
            >
                <MenuItem
                    onClick={(e) => {
                        e.stopPropagation();
                        handleClose();
                    }}
                    disabled={record.verified_status !== "" && record.verified_status !== "VERIFIED"}
                >
                    <PrintIcon fontSize="small" sx={{ mr: 1 }} />
                    <GenerateReportButton
                        results={record.test_result}
                        patient={record.patient}
                        workOrder={record}
                        currentGeneratedId={currentGeneratedId}
                        onGenerate={onGenerate}
                    />
                </MenuItem>

                {isFromNuha && (
                    <MenuItem
                        onClick={(e) => {
                            e.stopPropagation();
                            handleSendToSimrs();
                        }}
                        disabled={!hasResults || mutation.isPending}
                    >
                        {mutation.isPending ? (
                            <CircularProgress size={16} sx={{ mr: 1 }} />
                        ) : (
                            <SendIcon fontSize="small" sx={{ mr: 1 }} />
                        )}
                        <Typography variant="body2">
                            {mutation.isPending ? 'Sending to SIMRS...' : 'Send to SIMRS'}
                        </Typography>
                    </MenuItem>
                )}
            </Menu>
        </>
    );
}

function ResultActions() {
    const axios = useAxios()
    const notify = useNotify()
    return (
        <TopToolbar>
            <SyncAllResultButton />
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
    const { isLoading, isFetching, data } = useListContext();
    const [currentGeneratedId, setCurrentGeneratedId] = useState<string | null>(null);

    const handleGenerate = (buttonId: string) => {
        setCurrentGeneratedId(buttonId);
    };


    // Show loading only when actually loading and data hasn't been fetched yet
    const shouldShowLoading = (isLoading || isFetching) && data === undefined;

    if (shouldShowLoading) {
        return (
            <Box>
                <Box
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    minHeight="200px"
                    flexDirection="column"
                    gap={2}
                    sx={{
                        backgroundColor: 'background.paper',
                        borderRadius: 2,
                        mb: 2,
                        p: 3
                    }}
                >
                    <Box position="relative">
                        <CircularProgress
                            size={60}
                            thickness={4}
                            sx={{
                                color: 'primary.main',
                                animationDuration: '1.5s',
                            }}
                        />
                        <Box
                            position="absolute"
                            top="50%"
                            left="50%"
                            sx={{
                                transform: 'translate(-50%, -50%)',
                            }}
                        >
                            <ScienceIcon
                                sx={{
                                    fontSize: 24,
                                    color: 'primary.main',
                                    animation: 'pulse 2s infinite',
                                    '@keyframes pulse': {
                                        '0%': { opacity: 1 },
                                        '50%': { opacity: 0.5 },
                                        '100%': { opacity: 1 },
                                    }
                                }}
                            />
                        </Box>
                    </Box>

                    <Typography
                        variant="h6"
                        color="text.primary"
                        sx={{ fontWeight: 500 }}
                    >
                        Loading Lab Results...
                    </Typography>

                    <Typography
                        variant="body2"
                        color="text.secondary"
                        textAlign="center"
                    >
                        Please wait while we prepare your data and PDF reports...
                    </Typography>
                </Box>


            </Box>
        );
    }

    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Barcode" render={(record: WorkOrder) => (
                <Typography variant="body2">{record.barcode}</Typography>
            )} />
            <WithRecord label="SIMRS No Order" render={(record: WorkOrder) => (
                <Typography variant="body2">{record.barcode_simrs || '-'}</Typography>
            )} />
            <WithRecord label="No. RM" render={(record: WorkOrder) => (
                <Typography variant="body2">{record.medical_record_number || '-'}</Typography>
            )} />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"}
                    onClick={e => e.stopPropagation()}>
                    {record.patient.first_name} {record.patient.last_name}
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
            <WithRecord label="Status" render={(record: WorkOrder) => {
                const status = record.verified_status;
                const isVerified = status === "" || status === "VERIFIED";
                const isPending = status === "PENDING";
                const isRejected = status === "REJECTED";

                let color: "success" | "warning" | "error" = "success";
                let icon = <CheckCircleIcon />;
                let label = "Verified";

                if (isPending) {
                    color = "warning";
                    icon = <PendingIcon />;
                    label = "Not Verified";
                } else if (isRejected) {
                    color = "error";
                    icon = <CancelIcon />;
                    label = "Rejected";
                } else if (!isVerified) {
                    color = "error";
                    icon = <CancelIcon />;
                    label = "Not Verified";
                }

                return (
                    <MUIButton
                        variant="contained"
                        color={color}
                        startIcon={icon}
                        size="small"
                        sx={{
                            textTransform: 'none',
                            fontSize: '12px',
                            fontWeight: 600,
                            whiteSpace: 'nowrap',
                            borderRadius: '8px',
                            px: 2,
                            boxShadow: 2,
                            color: "white",
                            '&:hover': {
                                cursor: 'default',
                                boxShadow: 3,
                            }
                        }}
                    >
                        {label}
                    </MUIButton>
                );
            }} />
            <WithRecord label="SIMRS Status" render={(record: WorkOrder) => {
                const simrsStatus = record.simrs_sent_status || "";
                const isFromNuha = record.barcode_simrs && record.barcode_simrs.startsWith('NUHA-');

                // Don't show SIMRS status if not from Nuha
                if (!isFromNuha) {
                    return <Typography variant="body2" color="text.secondary">-</Typography>;
                }

                let color: "success" | "warning" | "error" | "info" = "info";
                let icon = <CloudQueueIcon />;
                let label = "Not Sent";

                if (simrsStatus === "SENT") {
                    color = "success";
                    icon = <CloudDoneIcon />;
                    label = "Sent";
                } else if (simrsStatus === "PARTIAL") {
                    color = "warning";
                    icon = <CloudQueueIcon />;
                    label = "Partial";
                } else if (simrsStatus === "FAILED") {
                    color = "error";
                    icon = <CloudOffIcon />;
                    label = "Failed";
                }

                return (
                    <MUIButton
                        variant="contained"
                        color={color}
                        startIcon={icon}
                        size="small"
                        sx={{
                            textTransform: 'none',
                            fontSize: '12px',
                            fontWeight: 600,
                            whiteSpace: 'nowrap',
                            borderRadius: '8px',
                            px: 2,
                            boxShadow: 2,
                            color: "white",
                            '&:hover': {
                                cursor: 'default',
                                boxShadow: 3,
                            }
                        }}
                    >
                        {label}
                    </MUIButton>
                );
            }} />
            <WithRecord label="Actions" render={(record: WorkOrder) => {
                // Show warning if not verified
                if (record.verified_status !== "" && record.verified_status !== "VERIFIED") {
                    return (
                        <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>

                            <ActionMenuButton
                                record={record}
                                currentGeneratedId={currentGeneratedId}
                                onGenerate={handleGenerate}
                            />
                        </Box>
                    )
                }

                // Show action menu for verified results
                return (
                    <ActionMenuButton
                        record={record}
                        currentGeneratedId={currentGeneratedId}
                        onGenerate={handleGenerate}
                    />
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
