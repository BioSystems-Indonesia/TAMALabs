import PlayCircleFilledIcon from "@mui/icons-material/PlayCircleFilled";
import { Box, CircularProgress, Dialog, DialogContent, DialogTitle, Divider, useTheme } from "@mui/material";
import SyncIcon from "@mui/icons-material/Sync";
import ScienceIcon from '@mui/icons-material/Science';
import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useMutation } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
    AutocompleteArrayInput,
    Button,
    Create,
    CreateButton,
    Datagrid,
    DateField,
    DeleteButton,
    Edit,
    FilterLiveForm,
    List,
    ReferenceArrayField,
    ReferenceField,
    ReferenceInput,
    ShowButton,
    TextField,
    TopToolbar,
    WithRecord,
    WrapperField,
    useListContext,
    useNotify,
    useRefresh,

} from "react-admin";
import { useParams, useSearchParams } from "react-router-dom";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import SideFilter from "../../component/SideFilter.tsx";
import useAxios from "../../hooks/useAxios.ts";
import { workOrderStatusDontShowRun, workOrderStatusShowCancel, type WorkOrder } from "../../types/work_order.ts";
import { WorkOrderChipColorMap } from "./ChipFieldStatus.tsx";
import WorkOrderForm from "./Form.tsx";
import RunWorkOrderForm from "./RunWorkOrderForm.tsx";

const WorkOrderAction = () => {
    return (
        <TopToolbar sx={{
            '& .RaTopToolbar-root': {
                padding: '16px 24px',
                minHeight: '64px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'flex-end',
                gap: '8px',
            }
        }}>
            <ShowButton />
        </TopToolbar>
    )
}

export function WorkOrderCreate() {
    return (
        <Create redirect={"show"} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }}>
            <WorkOrderForm mode={"CREATE"} />
        </Create>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit redirect={"show"} actions={<WorkOrderAction />} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }} mutationMode="pessimistic">
            <WorkOrderForm mode={"EDIT"} />
        </Edit>
    )
}

export function WorkOrderAddTest() {
    const { id } = useParams();
    const [searchParams] = useSearchParams();

    const getTitle = () => {
        const patientIDs = searchParams.getAll("patient_id")!.map(id => parseInt(id));
        if (patientIDs.length === 0) {
            return `Add Test to Work Order #${id}`
        }

        return `Edit Test Work Order #${id} for Patient ID: ${searchParams.getAll("patient_id").join(", ")}`
    }

    return (
        <Create
            title={getTitle()}
            redirect={() => {
                return `work-order/${id}/show`
            }} actions={<WorkOrderAction />} sx={{
                "& .RaCreate-card": {
                    overflow: "visible",
                }
            }} resource={`work-order/${id}/show/add-test`}>
            <WorkOrderForm mode={"ADD_TEST"} />
        </Create>
    )
}

function WorkOrderSideFilters() {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';

    return (
        <SideFilter sx={{
            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
        }}>
            <FilterLiveForm debounce={1500}>
                <Stack spacing={0}>
                    {/* Judul filter */}
                    <Box>
                        <Typography variant="h6" sx={{
                            color: theme.palette.text.primary,
                            marginBottom: 2,
                            fontWeight: 600,
                            fontSize: '1.1rem',
                            textAlign: 'center'
                        }}>
                            🔍 Filter Lab Requests
                        </Typography>
                    </Box>

                    {/* Filter Patient */}
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

                    {/* Filter Barcode */}
                    <ReferenceInput
                        source={"barcode_ids"} reference={`work-order/barcode`} alwaysOn
                        sx={{
                            '& .MuiInputLabel-root': {
                                color: theme.palette.text.primary,
                                fontWeight: 500,
                                fontSize: '0.9rem',
                            }
                        }}>
                        <AutocompleteArrayInput size="small"
                            sx={{
                                '& .MuiOutlinedInput-root': {
                                    backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                    borderRadius: '12px',
                                    transition: 'all 0.3s ease',
                                    '&:hover': {
                                        backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                    },
                                }
                            }} />
                    </ReferenceInput>

                    <Divider sx={{ marginBottom: 2 }} />

                    {/* Filter Date Range */}
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
                                clearable
                                sx={{
                                    marginBottom: '4px',
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
                                clearable
                                sx={{
                                    marginBottom: '4px',
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


function getRequestLength(data: WorkOrder): number {
    return data.specimen_list?.reduce((acc, specimen) => acc + specimen.observation_requests.length, 0) || 0
}

function RunWorkOrderButton(props: RunWorkOrderProps) {
    const notify = useNotify();
    const refresh = useRefresh();
    const axios = useAxios();
    const { isPending } = useMutation({
        mutationFn: async (data: any) => {
            const response = await axios.post(`/work-order/run`, data, {
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status != 200) {
                throw new Error(response.data?.error);
            }

            return response.data;
        },
        onSuccess: () => {
            notify('Success run', {
                type: 'success',
            });
            refresh()
        },
        onError: (error) => {
            notify('Error:' + error.message, {
                type: 'error',
            });
        },
    })

    return (
        <Button label="Run Work Order" onClick={() => {
            props.setOpen(true)
        }}>
            {isPending ? <CircularProgress size={12} variant='indeterminate' color='primary' /> : <PlayCircleFilledIcon />}
        </Button>
    )
}

type RunWorkOrderProps = {
    open: boolean
    setOpen: React.Dispatch<React.SetStateAction<boolean>>
    onClose: () => void
}

function RunWorkOrderDialog(props: RunWorkOrderProps) {
    const { selectedIds, data } = useListContext<WorkOrder>();
    const [processing, setProcessing] = useState(false)
    const notify = useNotify();
    const [dataMap, setDataMap] = useState<Record<number, WorkOrder>>({})
    useEffect(() => {
        if (data) {
            const map: Record<number, WorkOrder> = {}
            data.forEach((workOrder) => {
                map[workOrder.id] = workOrder
            })
            setDataMap(map)
        }
    }, [data])

    function determineShowCancelButton(selectedIds: number[], dataMap: Record<number, WorkOrder>): boolean | undefined {
        if (selectedIds.length === 0) {
            return undefined
        }

        for (const id of selectedIds) {
            if (workOrderStatusShowCancel.includes(dataMap[id].status)) {
                return true
            }
        }

        return false
    }

    function determineShowRunButton(selectedIds: number[], dataMap: Record<number, WorkOrder>): boolean | undefined {
        if (selectedIds.length === 0) {
            return undefined
        }

        for (const id of selectedIds) {
            if (!workOrderStatusDontShowRun.includes(dataMap[id].status)) {
                return true
            }
        }

        return false
    }

    function determineDefaultDeviceID(selectedIds: number[], dataMap: Record<number, WorkOrder>): number | undefined {
        if (selectedIds.length === 0) {
            return undefined
        }

        for (const id of selectedIds) {
            const workOrder = dataMap[id]
            if (workOrder.devices && workOrder?.devices?.length > 0) {
                return workOrder?.devices[0].id
            }
        }

        return undefined
    }

    return (
        <Dialog
            open={props.open}
            onClose={() => {
                if (processing) {
                    notify('Cannot close dialog while processing', {
                        type: 'error',
                    });
                    return;
                }

                props.onClose()
            }}
            fullWidth
            sx={{
                width: "100%",
                margin: 0,
            }}
            maxWidth="lg"
        >
            <DialogTitle id="alert-dialog-title">
                Run Work Order
            </DialogTitle>
            <DialogContent sx={{}}>
                <RunWorkOrderForm workOrderIDs={selectedIds} setIsProcessing={setProcessing} isProcessing={processing}
                    showCancelButton={determineShowCancelButton(selectedIds, dataMap)}
                    showRunButton={determineShowRunButton(selectedIds, dataMap)}
                    defaultDeviceID={determineDefaultDeviceID(selectedIds, dataMap)}
                />
            </DialogContent>
        </Dialog >
    )
}


const WorkOrderListBulkActionButtons = (props: RunWorkOrderProps) => (
    <>
        <RunWorkOrderButton {...props} />
    </>
)


function WorkOrderListActions() {
    const axios = useAxios()
    const notify = useNotify()
    const refresh = useRefresh()
    return (
        <TopToolbar>
            <Button label={"Sync request from SIMRS"} onClick={async () => {
                const response = await axios.post("/external/sync-all-requests", {})

                refresh()
                notify("Sync Success " + response.statusText, {
                    type: "success"
                })
            }}>
                <SyncIcon />
            </Button>
            <CreateButton/>
        </TopToolbar>
    )
}

export const WorkOrderList = () => {
    const [open, setOpen] = useState(false)

    return (
        <List sort={{
            field: "id",
            order: "DESC"
        }} aside={<WorkOrderSideFilters />} 
        actions={<WorkOrderListActions/>}
        title="Lab Request" exporter={false}
            storeKey={false}
            sx={{
                '& .RaList-content': {

    const WorkOrderDataGrid = () => {
    const { isLoading, isFetching, data } = useListContext();
    const [open, setOpen] = useState(false);
    const [initialLoading, setInitialLoading] = useState(true);

    useEffect(() => {
        if (data && data.length > 0) {
            const timer = setTimeout(() => {
                setInitialLoading(false);
            }, 500);
            return () => clearTimeout(timer);
        }
    }, [data]);

    const shouldShowLoading = isLoading || isFetching || initialLoading || !data;

    if (shouldShowLoading) {
        return (
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
                                },
                            }}
                        />
                    </Box>
                </Box>

                <Typography
                    variant="h6"
                    color="text.primary"
                    sx={{ fontWeight: 500 }}
                >
                    Loading Lab Requests...
                </Typography>

                <Typography
                    variant="body2"
                    color="text.secondary"
                    textAlign="center"
                >
                    Please wait while we fetch your data
                </Typography>
            </Box>
        );
    }

    return (
        <>
            <Datagrid
                rowClick={(id, resource, record) => {
                    return false
                }}
                bulkActionButtons={<WorkOrderListBulkActionButtons
                    open={open}
                    setOpen={setOpen}
                    onClose={() => setOpen(false)}
                />}>
                <TextField source="id" />
                <WithRecord label="Status" render={(record: any) => (
                    <Chip label={`${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                )} />
                <ReferenceField source="patient_id" reference="patient">
                </ReferenceField>
                <TextField source="barcode" />
                <WithRecord label="Barcode SIMRS" render={(record: any) => (
                    <Typography variant="body2">{record.barcode_simrs || '-'}</Typography>
                )} />
                <WithRecord label="Request" render={(record: any) => (
                    <Typography variant="body2" >
                        {getRequestLength(record)}
                    </Typography>
                )} />
                <ReferenceArrayField source="doctor_ids" reference="user" />
                <ReferenceArrayField source="analyzer_ids" reference="user" />
                <DateField source="created_at" />
                <WrapperField label="Actions" sortable={false} >
                    <Stack direction={"row"} spacing={2}>
                        <ShowButton variant="contained" />
                        <DeleteButton variant="contained" />
                    </Stack>
                </WrapperField>
            </Datagrid>
            <RunWorkOrderDialog open={open} onClose={() => setOpen(false)} setOpen={setOpen} />
        </>
    );
};

export const WorkOrderList = () => {
    return (
        <List sort={{
            field: "id",
            order: "DESC"
        }} aside={<WorkOrderSideFilters />} title="Lab Request" exporter={false}
            storeKey={false}
            sx={{
                '& .RaList-content': {
                    backgroundColor: 'background.paper',
                    padding: 2,
                    borderRadius: 1,
                },
            }}
        >
            <WorkOrderDataGrid />
        </List>
    )
};
