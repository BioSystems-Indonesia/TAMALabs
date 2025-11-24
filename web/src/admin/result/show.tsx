import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import CheckIcon from '@mui/icons-material/CheckCircleOutline'; // Using outline for a slightly different style
import CloseIcon from '@mui/icons-material/HighlightOff'; // Using a different close icon for variety

import HistoryIcon from '@mui/icons-material/History';
import {
    Badge,
    Box,
    Button,
    ButtonGroup,
    Card,
    CardActions,
    Checkbox,
    Chip,
    Dialog,
    DialogContent,
    DialogTitle,
    GridLegacy as Grid,
    IconButton,
    Stack,
    Tooltip,
    Typography
} from "@mui/material";
import { DataGrid as MuiDatagrid, type DataGridProps, type GridRenderCellParams } from '@mui/x-data-grid';
import dayjs from 'dayjs';
import { memo, useEffect, useState } from 'react';
import {
    DateField,
    DeleteButton,
    Labeled,
    Link,
    Show,
    SimpleShowLayout,
    TextField,
    WithRecord,
    useNotify,
    useRecordContext,
    useRedirect,
    useRefresh
} from "react-admin";
import { useCurrentUser } from '../../hooks/currentUser';
import useAxios from '../../hooks/useAxios';
import type { ResultColumn } from "../../types/general";
import { Result, TestResult } from '../../types/observation_result';
import type { Specimen } from '../../types/specimen';
import { User } from '../../types/user';
import type { WorkOrder } from '../../types/work_order';
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { FilledPercentChip, VerifiedChip } from './component';
import { CreatedBy } from '../../types/constant';

export const ResultShow = (props: any) => {
    const [openHistory, setOpenHistory] = useState(false);
    const [history, setHistory] = useState<HistoryChangeProps>({
        rows: [],
        title: '',
    });

    return (
        <Show title="Edit Result" sx={{
            overflow: 'visible',
            "& .RaShow-main": {
                overflow: 'visible',
            },
            "& .RaShow-card": {
                overflow: 'visible',
            }
        }}>
            <SimpleShowLayout sx={{
                overflow: 'visible',
                position: 'relative',
                "& .RaSimpleShowLayout-stack": {
                    display: 'block',
                    overflow: 'visible',
                },
                "& .RaSimpleShowLayout-row": {
                    overflow: 'visible',
                },
            }}>
                <ActionButton />
                <HeaderInfo />
                <WithRecord label="Test Result" render={(record: Result) => (
                    <>
                        {
                            record?.test_result ? Object.entries(record.test_result).map(([category, rows]) => (
                                <TestResultTableGroup
                                    key={category}
                                    category={category}
                                    rows={rows}
                                    setHistory={setHistory}
                                    setOpenHistory={setOpenHistory}
                                />
                            )) : null
                        }

                        <HistoryDialog
                            workOrderID={record.id}
                            title={history.title}
                            open={openHistory}
                            onClose={() => setOpenHistory(false)}
                            rows={history.rows}
                            setHistory={setHistory}
                        />
                    </>)
                } />

            </SimpleShowLayout>
        </Show>
    )
}

const ActionButton = () => {
    const record = useRecordContext<Result>()
    const redirect = useRedirect()
    const notify = useNotify();
    const axios = useAxios();
    const refresh = useRefresh();
    const currentUser = useCurrentUser();

    return (
        <Card
            sx={{
                position: 'sticky',
                overflow: 'visible',
                top: 0,
                zIndex: 100,
                borderRadius: '12px', // Rounded corners for the card
                p: 1, // Padding inside the card
                boxShadow: 'rgba(0, 0, 0, 0.15) 1.95px 1.95px 2.6px', // A subtle shadow
                width: '100%', // Make card take full width of its container
            }}
        >
            <CardActions
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    width: '100%',
                    p: 0, // No padding for CardActions to use the card's padding
                    flexWrap: 'wrap', // Allow buttons to wrap on very small screens
                    gap: 1, // Gap between button groups if they wrap
                }}
            >
                {/* Left side: Previous and Next buttons */}
                <Box sx={{ display: 'flex', gap: 1, flexGrow: 1, justifyContent: { xs: 'center', sm: 'flex-start' } }}>
                    <Button
                        variant="outlined"
                        startIcon={<ArrowBackIcon />}
                        onClick={() => redirect(`/result/${record?.prev_id}/show`)}
                        disabled={record?.prev_id === 0}
                        sx={{ borderRadius: '8px', textTransform: 'none' }}
                        aria-label="Previous item"
                    >
                        Previous
                    </Button>
                    <Button
                        variant="outlined"
                        endIcon={<ArrowForwardIcon />}
                        onClick={() => redirect(`/result/${record?.next_id}/show`)}
                        disabled={record?.next_id === 0}
                        sx={{ borderRadius: '8px', textTransform: 'none' }}
                        aria-label="Next item"
                    >
                        Next
                    </Button>
                </Box>

                {
                    record?.doctors?.map(v => v.id).includes(currentUser?.id ?? 0) && (
                        <Box sx={{ display: 'flex', gap: 1, flexGrow: 1, justifyContent: { xs: 'center', sm: 'flex-end' } }}>
                            <Button
                                variant="contained"
                                color="success"

                                startIcon={<CheckIcon />}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    axios.post(`/result/${record?.id}/approve`, {})
                                        .then((response) => {
                                            notify(`Success approve result`, { type: 'success' });
                                            refresh()
                                        })
                                        .catch((error) => {
                                            notify(`Error approve result ${error}`, { type: 'error' });
                                        });
                                }}
                                disabled={record?.verified_status === "VERIFIED"}
                                sx={{
                                    color: 'white !important',
                                    borderRadius: '8px',
                                    textTransform: 'none',
                                }}
                                aria-label="Approve item"
                            >
                                Approve
                            </Button>
                            <Button
                                variant="contained"
                                color="error"
                                startIcon={<CloseIcon />}
                                disabled={record?.verified_status === "REJECTED"}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    axios.post(`/result/${record?.id}/reject`, {})
                                        .then((response) => {
                                            notify(`Success reject result`, { type: 'success' });
                                            refresh()
                                        })
                                        .catch((error) => {
                                            notify(`Error reject result ${error}`, { type: 'error' });
                                        });
                                }}
                                sx={{
                                    borderRadius: '8px',
                                    textTransform: 'none',
                                }}
                            >
                                Reject
                            </Button>
                        </Box>
                    )
                }
            </CardActions>
        </Card>
    );
};


const HeaderInfo = (props: any) => (
    <Grid sx={{
        display: "flex",
        border: "1px solid #ccc",
        padding: "26px",
        borderRadius: '0.5rem',
        margin: '2rem 0'
    }} container rowGap={1}>
        <Grid item xs={12} md={12} >
            <Labeled>
                <WithRecord label="Barcodes" render={(record: WorkOrder) => {
                    return (
                        <Stack direction={"row"} gap={1}>
                            {(Array.isArray(record?.specimen_list) ? record.specimen_list : []).map((specimen: Specimen, idx: number) => {
                                return (
                                    <Chip key={specimen.barcode || idx} label={specimen.barcode} />
                                )
                            })}
                        </Stack>
                    )
                }} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4}>
            <Labeled>
                <WithRecord label="Patient" render={(record: any) => (
                    <Link to={`/patient/${record.patient?.id}/show`} resource="patient" label={"Patient"} onClick={e => e.stopPropagation()}>
                        #{record.patient?.id}-{record.patient?.first_name} {record.patient?.last_name}
                    </Link>
                )} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4}>
            <Labeled>
                <WithRecord label="Work Order" render={(record: any) => (
                    <Link to={`/work-order/${record.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                        <Chip label={`#${record.id} - ${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                    </Link>
                )} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4}>
            <Labeled>
                <DateField source="created_at" showTime />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <TextField source="total_request" label="Total Request" />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <TextField source="total_result_filled" label="Total Result Filled" />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <WithRecord label="Filled" render={(record: WorkOrder) => (
                    <FilledPercentChip percent={record.percent_complete} />
                )} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <WithRecord label="Doctors" render={(record: WorkOrder) => {
                    return (
                        <Stack direction={"row"} gap={1}>
                            {record?.doctors?.map((user: User) => {
                                return (
                                    <Chip label={`${user.id} - ${user.fullname}`} />
                                )
                            })}
                        </Stack>
                    )
                }} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <WithRecord label="Analysts" render={(record: WorkOrder) => {
                    return (
                        <Stack direction={"row"} gap={1}>
                            {record?.analyzers?.map((user: User) => {
                                return (
                                    <Chip label={`${user.id} - ${user.fullname}`} />
                                )
                            })}
                        </Stack>
                    )
                }} />
            </Labeled>
        </Grid>
        <Grid item xs={12} md={4} >
            <Labeled>
                <WithRecord label="Verified" render={(record: WorkOrder) => {
                    return (
                        <VerifiedChip verified={record.verified_status !== '' ? record.verified_status : "VERIFIED"} />
                    )
                }} />
            </Labeled>
        </Grid>
    </Grid>
)

// using memo because we don't want this table rerendered when dialog appear or disappear
const TestResultTableGroup = memo((props: TestResultTableProps) => {
    return <Stack sx={{
        marginY: 2,
        width: "100%",
    }} key={props.category}>
        <Typography variant="subtitle1" gutterBottom>
            {props.category}
        </Typography>
        <Typography variant="caption" gutterBottom>
            Click result column to edit
        </Typography>
        <TestResultTable {...props} />
    </Stack>
})

type TestResultTableProps = {
    category: string
    rows: TestResult[]
    setHistory: (props: HistoryChangeProps) => void
    setOpenHistory: (open: boolean) => void
}

const TestResultTable = (props: TestResultTableProps) => {
    const notify = useNotify();
    const axios = useAxios();

    function onUpdateError(error: any): void {
        notify(`Error update ${error}`, { type: 'error' });
    }

    async function putResult(newRow: TestResult, _oldRow: TestResult) {
        if (newRow.result === _oldRow.result && newRow.unit === _oldRow.unit) {
            return _oldRow
        }
        const url = `/result/${newRow.specimen_id}/test`

        try {
            const response = await axios.put(url, newRow);
            notify(`Success update ${newRow.test}`, { type: 'success' });
            return response.data;
        } catch (error: any) {
            notify(`Error update ${error}`, { type: 'error' });
        }
    }

    let negID = -1

    // support id == 0 when the TestResult is not set yet
    // TODO find better hack than this
    const [rows, setRows] = useState<any>([])
    useEffect(() => {
        if (!props?.rows) return;
        if (!Array.isArray(props.rows)) return;

        console.log(props.rows)

        setRows(props.rows.map((r: any) => ({
            ...r,
            id: r.id || negID--,
            name: r?.test_type?.name || r?.history?.[0]?.test_type?.name || r.test,
            specimen_type: r?.test_type?.types?.[0]?.type || '',
            alias: r?.test_type?.alias_code || r?.history?.[0]?.test_type?.alias_code || r.alias || r.test,
        })));
    }, [props?.rows]);

    return (
        <MuiDatagrid rows={rows}
            pageSizeOptions={[-1]}
            hideFooter
            editMode="row"
            processRowUpdate={putResult}
            onProcessRowUpdateError={onUpdateError}
            rowHeight={60}
            columns={[
                {
                    field: 'name',
                    headerName: 'Parameter Name',
                    flex: 1,
                },
                {
                    field: 'specimen_type',
                    headerName: 'Specimen Type',
                    flex: 1,
                },
                {
                    field: 'alias',
                    headerName: 'Alias/Code',
                    flex: 1,
                },

                {
                    field: 'result',
                    headerName: 'Result',
                    type: 'string',
                    editable: true,
                    flex: 1,

                },
                {
                    field: 'unit',
                    headerName: 'Unit',
                    flex: 1,

                },
                {
                    field: 'reference_range',
                    headerName: 'Reference Range',
                    flex: 2,
                },
                {
                    // Error
                    field: 'abnormal',
                    headerName: 'Status',
                    flex: 1,
                    renderCell: (params: GridRenderCellParams) => {
                        switch (params.value) {
                            case 0: return <Chip color="success" label="Normal" />
                            case 1: return <Chip color="error" label="High" />
                            case 2: return <Chip color="secondary" label="Low" />
                            case 3: return <Chip color="default" label="No Data" />
                            case 4: return <Chip color="warning" label="Positive" />
                            case 5: return <Chip color="info" label="Negative" />
                            default: return <Chip color="success" label="Normal" />
                        }
                    },
                },
                {
                    field: 'created_by',
                    headerName: 'Input By',
                    flex: 2,
                    renderCell: (params: GridRenderCellParams) => {
                        switch (params.value?.id) {
                            case 0: return ""
                            case CreatedBy.Unknown: return <Chip label="Unknown" />
                            case CreatedBy.System: return <Chip color='primary' label="System" />
                            default: return <Chip color='info' label={`${params.value?.fullname}`} />
                        }
                    },
                },
                {
                    field: 'created_at',
                    headerName: 'Date',
                    flex: 2,
                    renderCell: (params: GridRenderCellParams) => {
                        return dayjs(params.value).format("YYYY-MM-DD HH:mm:ss")
                    },
                },
                {
                    field: '',
                    headerName: 'Action',
                    flex: 1,
                    renderCell: (params: GridRenderCellParams) => {
                        const resultDifference = !params.row.history
                            .map((h: TestResult) => "" + h.result + h.unit)
                            .every((v: string, _: number, a: string[]) => v === a[0])

                        return <Box>
                            <Tooltip title={resultDifference ? "History has different result" : "Show History"}>
                                <Badge badgeContent={params.row.history.length} color={resultDifference ? "warning" : "primary"}>
                                    <IconButton color={resultDifference ? "warning" : "primary"}
                                        onClick={() => {
                                            props.setHistory({
                                                rows: params.row.history,
                                                title: `History of ${params.row.test}`
                                            })
                                            props.setOpenHistory(true)
                                        }} >
                                        <HistoryIcon />
                                    </IconButton>
                                </Badge>
                            </Tooltip>
                        </Box>
                    },
                },
            ]} />
    );
}

type HistoryChangeProps = {
    title: string
    rows: ResultColumn[]
}

type HistoryDialogProps = {
    workOrderID: number
    title: string
    rows: ResultColumn[]
    open: boolean
    onClose: () => void
    setHistory: (props: HistoryChangeProps) => void
} & Partial<DataGridProps<ResultColumn>>

const HistoryDialog = (props: HistoryDialogProps) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const record = useRecordContext<Result>();
    const currentUser = useCurrentUser();

    // Some Hack because Dialog is on top, and it will refresh
    // some updated value, then we need to make the values
    // not rollback into previous value
    const onClose = () => {
        props.onClose()
    }

    const axios = useAxios();

    // Check if current user is a doctor
    const isDoctor = record?.doctors?.map(v => v.id).includes(currentUser?.id ?? 0);
    const pickTestResult = async (testResultID: number) => {
        try {
            const url = `/result/${props.workOrderID}/test/${testResultID}/pick`
            const response = await axios.put(url);
            if (response.status !== 200) {
                throw new Error("Error pick test result");
            }

            notify("Success pick test result", {
                type: 'success',
            });

            refresh()

            props.setHistory({
                rows: props.rows.map(v => v.id === testResultID ? { ...v, picked: true } : { ...v, picked: false }),
                title: props.title,
            })
        } catch (err) {
            notify("Error pick test result", {
                type: 'error',
            });
        }
    }

    return (
        <Dialog
            open={props.open}
            onClose={onClose}
            fullWidth
            sx={{
                width: "100%",
                margin: 0,
            }}
            maxWidth="lg"
        >
            <DialogTitle id="alert-dialog-title">
                {props.title}
            </DialogTitle>
            <DialogContent>
                <MuiDatagrid
                    pageSizeOptions={[-1]}
                    hideFooter
                    rowHeight={60}
                    rows={props.rows}
                    columns={[
                        {
                            field: 'test',
                            headerName: 'Test',
                            flex: 1,
                        },
                        {
                            field: 'specimen_type',
                            headerName: 'Specimen Type',
                            flex: 1,
                        },
                        {
                            field: 'result',
                            headerName: 'Result',
                            type: 'string',
                            flex: 1,
                        },
                        {
                            field: 'unit',
                            headerName: 'Unit',
                            flex: 1,
                        },
                        {
                            field: 'reference_range',
                            headerName: 'Reference Range',
                            flex: 1,
                        },
                        {
                            field: 'created_at',
                            headerName: 'Created At',
                            flex: 1,
                            renderCell: (params: GridRenderCellParams) => {
                                return dayjs(params.value).format("YYYY-MM-DD HH:mm:ss")
                            },
                        },
                        {
                            field: 'picked',
                            headerName: 'Picked',
                            flex: 1,
                            renderCell: (params: GridRenderCellParams) => {
                                return <Checkbox checked={params.value} readOnly onClick={() => {
                                    pickTestResult(params.row.id)
                                }} />
                            },
                        },
                        {
                            field: 'created_by',
                            headerName: 'Input By',
                            flex: 1,
                            renderCell: (params: GridRenderCellParams) => {
                                switch (params.value?.id) {
                                    case 0: return ""
                                    case CreatedBy.Unknown: return <Chip label="Unknown" />
                                    case CreatedBy.System: return <Chip color='primary' label="System" />
                                    default: return <Chip color='info' label={`${params.value?.fullname}`} />
                                }
                            },
                        },
                        {
                            field: '',
                            headerName: 'Action',
                            flex: 1,
                            renderCell: (params: GridRenderCellParams) => {

                                return <ButtonGroup sx={{
                                    gap: 2,
                                }}>
                                    <DeleteButton
                                        disabled={!isDoctor}
                                        sx={{ marginLeft: 2 }}
                                        label={''}
                                        mutationMode="pessimistic"
                                        size='medium'
                                        resource={`result/${props.workOrderID}/test`}
                                        variant='text'
                                        record={{ id: params.row.id }}
                                        confirmTitle={`Delete test ${params.row.test}?`}
                                        confirmColor="warning"
                                        confirmContent="This will delete the test result. This action cannot be undone."
                                        redirect={false}
                                        mutationOptions={{
                                            onError: () => {
                                                notify(`Error delete test ${params.row.test}`, {
                                                    type: 'error',
                                                });
                                            },
                                            onSuccess: () => {
                                                notify(`Success delete test ${params.row.test}`, {
                                                    type: 'success',
                                                });

                                                refresh()
                                                const newHistoryRows = props.rows.filter(v => v.id !== params.row.id)
                                                props.setHistory({
                                                    rows: newHistoryRows,
                                                    title: props.title,
                                                })
                                            },
                                        }}
                                    />
                                </ButtonGroup>
                            }
                        }
                    ]} />
            </DialogContent>
        </Dialog>
    )
}