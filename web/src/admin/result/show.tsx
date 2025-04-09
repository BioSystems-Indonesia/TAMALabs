import HistoryIcon from '@mui/icons-material/History';
import { Badge, Box, ButtonGroup, Checkbox, Chip, Dialog, DialogContent, DialogTitle, GridLegacy as Grid, IconButton, Stack, Tooltip, Typography } from "@mui/material";
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
    useRefresh
} from "react-admin";
import useAxios from '../../hooks/useAxios';
import type { ResultColumn } from "../../types/general";
import { Result, TestResult } from '../../types/observation_result';
import type { Specimen } from '../../types/specimen';
import type { WorkOrder } from '../../types/work_order';
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { FilledPercentChip } from './component';

export const ResultShow = (props: any) => {
    const [openHistory, setOpenHistory] = useState(false);
    const [history, setHistory] = useState<HistoryChangeProps>({
        rows: [],
        title: '',
    });

    return (
        <Show title="Edit Result">
            <SimpleShowLayout >
                <HeaderInfo />
                <WithRecord label="Test Result" render={(record: Result) => (
                    <>
                        {
                            Object.entries(record?.test_result).map(([category, rows]) => (
                                <TestResultTableGroup
                                    key={category}
                                    category={category}
                                    rows={rows}
                                    setHistory={setHistory}
                                    setOpenHistory={setOpenHistory}
                                />
                            ))
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

const HeaderInfo = (props: any) => (
    <Grid sx={{
        display: "flex",
        border: "1px solid #ccc",
        padding: "12px",
    }} container rowGap={1}>
        <Grid item xs={12} md={12} >
            <Labeled>
                <WithRecord label="Barcodes" render={(record: WorkOrder) => {
                    return (
                        <Stack direction={"row"} gap={1}>
                            {record.specimen_list.map((specimen: Specimen) => {
                                return (
                                    <Chip label={specimen.barcode} />
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
                    <Link to={`/patient/${record.order_id}/show`} resource="patient" label={"Patient"} onClick={e => e.stopPropagation()}>
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

    function onUpdateError(error: any): void {
        notify(`Error update ${error}`, { type: 'error' });
    }

    async function putResult(newRow: TestResult, _oldRow: TestResult) {
        if (newRow.result === _oldRow.result && newRow.unit === _oldRow.unit) {
            return _oldRow
        }

        const url = `${import.meta.env.VITE_BACKEND_BASE_URL}/result/${newRow.specimen_id}/test`

        const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newRow)
        });

        if (!response.ok) {
            const responseJson = await response.json();
            throw new Error(responseJson?.error);
        }

        const respJSON = await response.json();
        notify(`Success update ${newRow.test}`, { type: 'success' });

        return respJSON
    }

    let negID = -1

    // support id == 0 when the TestResult is not set yet
    // TODO find better hack than this
    const [rows, setRows] = useState<any>([])
    useEffect(() => {
        if (!props?.rows) {
            return
        }

        // Check if rows is array
        if (!Array.isArray(props.rows)) {
            return
        }

        setRows(props.rows.map((r: any) => ({ ...r, id: r.id || negID-- })))
    }, [props?.rows])

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
                    field: 'test',
                    headerName: 'Test',
                    flex: 2,
                    renderCell: (params: GridRenderCellParams) => {
                        if (!params.row.test_type_id) {
                            return <Link to={`/test-type/create?code=${params.row.test}`} onClick={e => e.stopPropagation()} sx={{
                                textAlign: "center",
                                width: "100%",
                                height: "100%",
                                display: "flex",
                                justifyContent: "start",
                                textDecoration: "unset",
                                alignItems: "center",
                            }}>
                                <Tooltip title="Test Type is Unset, please create it first">
                                    <Chip label={params.row.test} color="error" />
                                </Tooltip>
                            </Link >
                        }

                        return <>
                            <Link to={`/test-type/${params.row.test_type_id}/show`} onClick={e => e.stopPropagation()} sx={{
                                width: "100%",
                                height: "100%",
                                textAlign: "center",
                                display: "flex",
                                justifyContent: "start",
                                alignItems: "center",
                            }}>
                                <p>{params.row.test}</p>
                            </Link>
                        </>
                    },
                },
                {
                    field: 'result',
                    headerName: 'Result',
                    type: 'number',
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
                    field: 'abnormal',
                    headerName: 'Status',
                    flex: 1,
                    renderCell: (params: GridRenderCellParams) => {
                        switch (params.value) {
                            case 1: return <Chip color="error" label="High" />
                            case 2: return <Chip color="primary" label="Low" />
                            default: return <Chip color="success" label="Normal" />
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
                        let resultDifference = !params.row.history
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

    // Some Hack because Dialog is on top, and it will refresh
    // some updated value, then we need to make the values
    // not rollback into previous value
    const onClose = () => {
        props.onClose()
    }

    const axios = useAxios();
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
                            field: 'result',
                            headerName: 'Result',
                            type: 'number',
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
                            field: '',
                            headerName: 'Action',
                            flex: 1,
                            renderCell: (params: GridRenderCellParams) => {
                                return <ButtonGroup sx={{
                                    gap: 2,
                                }}>
                                    <DeleteButton
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
