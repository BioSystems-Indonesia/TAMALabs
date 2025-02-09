import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import HistoryIcon from '@mui/icons-material/History';
import { Badge, Box, ButtonGroup, Chip, Dialog, DialogContent, DialogTitle, Grid, IconButton, Stack, Tooltip, Typography } from "@mui/material";
import { DataGrid as MuiDatagrid, type DataGridProps, type GridRenderCellParams } from '@mui/x-data-grid';
import dayjs from 'dayjs';
import { useState } from 'react';
import {
    ArrayInput,
    AutocompleteInput,
    Button,
    Confirm,
    Create,
    DateField,
    DeleteButton,
    Labeled,
    Link,
    NumberInput,
    ReferenceInput,
    Show,
    SimpleForm,
    SimpleFormIterator,
    SimpleShowLayout,
    TextField,
    WithRecord,
    useDeleteMany,
    useNotify,
    useRefresh,
    type ButtonProps
} from "react-admin";
import { useSearchParams } from "react-router-dom";
import { getRefererParam, useRefererRedirect } from "../../hooks/useReferer";
import type { ResultColumn } from "../../types/general";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";


export const ResultShow = () => {
    const notify = useNotify();

    function onUpdateError(error: any): void {
        notify(`Error update ${error}`, {
            type: 'error',
        });
    }

    async function updateResult(newRow: ResultColumn, oldRow: ResultColumn) {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/result`, {
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                data: [newRow]
            }),
            method: 'PUT',
        });

        if (!response.ok) {
            const responseJson = await response.json();
            throw new Error(responseJson?.error);
        } else {
            const respJSON = await response.json();
            if (Array.isArray(respJSON) === false) {
                console.error("respJSON is not array");
                return newRow
            }

            notify(`Success update row ${newRow.test}`, {
                type: 'success',
            });
            refresh();

            return respJSON[0]
        }

    }

    const refresh = useRefresh();
    const [openHistory, setOpenHistory] = useState(false);
    const [history, setHistory] = useState<HistoryChangeProps>({
        rows: [],
        title: '',
    });
    const baseColumns = [
        {
            field: 'test',
            headerName: 'Test',
            flex: 1,
            renderCell: (params: GridRenderCellParams) => {
                if (!params.row.test_type_id) {
                    return <Link to={`/test-type/create?code=${params.row.test}`} onClick={e => e.stopPropagation()} sx={{
                        textAlign: "center",
                        width: "100%",
                        height: "100%",
                        display: "flex",
                        justifyContent: "center",
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
                        justifyContent: "center",
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
            flex: 1,
            editable: true,
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
            field: 'abnormal',
            headerName: 'Status',
            flex: 1,
            renderCell: (params: GridRenderCellParams) => {
                return params.value === 1 ? (
                    <Chip color="error" label="High" />
                ) : params.value === 2 ? (
                    <Chip color="warning" label="Low" />
                ) :
                    (
                        <Chip color="primary" label="Normal" />
                    )
            },
        },
    ]

    return (
        <Show title="Edit Result">
            <SimpleShowLayout >
                <Grid sx={{
                    display: "flex",
                    border: "1px solid #ccc",
                    padding: "12px",
                }} container>
                    <Grid item xs={12} md={12} sx={{
                    }}>
                        <Labeled>
                            <TextField source="barcode" label="Barcode" />
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
                                <Link to={`/work-order/${record.order_id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                                    <Chip label={`#${record?.order_id} - ${record.work_order?.status}`} color={WorkOrderChipColorMap(record.work_order?.status)} />
                                </Link>
                            )} />
                        </Labeled>
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <Labeled>
                            <DateField source="created_at" showTime />
                        </Labeled>
                    </Grid>
                </Grid>
                <WithRecord label="Test Result" render={(record: any) => {
                    if (!record?.test_result || Object.keys(record?.test_result).length === 0) {
                        return <Stack sx={{
                            width: "100%",
                            marginY: 2,
                        }} gap={0.5}>
                            <Typography variant="h5" color={"text.primary"} sx={{
                                width: "100%",
                                textAlign: "center",
                            }}>No Test Result found</Typography>
                            <Link to={`add-result?specimen_id=${record?.id}&${getRefererParam()}`}>
                                <Button label="Add Result" variant="contained" sx={{
                                    width: "default",
                                }}>
                                    <AddIcon />
                                </Button>
                            </Link>
                        </Stack>
                    }
                    const [expandedRows, setExpandedRows] = useState<string[]>([]);

                    const checkHistoryHaveDifferentResult = (row: ResultColumn, history: ResultColumn[]) => {
                        let resultDifference = false
                        for (const v of history) {
                            if (v.result !== row.result) {
                                resultDifference = true
                                break
                            }
                        }
                        return resultDifference
                    }

                    return (
                        <>
                            <Stack sx={{
                                display: "flex",
                                alignItems: "flex-end",
                                width: "100%",
                                marginTop: 2,
                                paddingRight: 5,
                            }}>
                                <Link to={`add-result?specimen_id=${record.id}&${getRefererParam()}`}>
                                    <Button label="Add Result" variant="contained" sx={{
                                        width: "default",
                                    }}>
                                        <AddIcon />
                                    </Button>
                                </Link>
                            </Stack>
                            {
                                Object.entries(record?.test_result).map(([category, result]: any) => (
                                    <Stack sx={{
                                        marginY: 2,
                                        width: "100%",
                                    }} key={category}>
                                        <Typography variant="subtitle1" gutterBottom>
                                            {category}
                                        </Typography>
                                        <Typography variant="caption" gutterBottom>
                                            Click result column to edit
                                        </Typography>
                                        <MuiDatagrid rows={result}
                                            pageSizeOptions={[-1]}
                                            hideFooter
                                            editMode="row"
                                            processRowUpdate={updateResult}
                                            onProcessRowUpdateError={onUpdateError}
                                            rowHeight={60}
                                            columns={[
                                                ...baseColumns,
                                                {
                                                    field: '',
                                                    headerName: 'Action',
                                                    flex: 1,
                                                    renderCell: (params: GridRenderCellParams) => {
                                                        let resultDifference = checkHistoryHaveDifferentResult(params.row, params.row.history)

                                                        return <Box sx={{
                                                        }}>
                                                            <Tooltip title={resultDifference ? "History has different result" : "Show History"}>
                                                                <Badge badgeContent={params.row.history.length} color={resultDifference ? "warning" : "primary"}>
                                                                    <IconButton color={resultDifference ? "warning" : "primary"}
                                                                        onClick={() => {
                                                                            setHistory({
                                                                                rows: params.row.history,
                                                                                title: `History of ${params.row.test}`
                                                                            })
                                                                            setOpenHistory(true)
                                                                        }} >
                                                                        <HistoryIcon />
                                                                    </IconButton>
                                                                </Badge>
                                                            </Tooltip>


                                                            <DeleteTestButton label={''} variant='text' size='large' resource="result"
                                                                ids={params.row.history.map((v: ResultColumn) => v.id)}
                                                                code={params.row.test}
                                                                onError={() => {
                                                                    notify(`Error delete test ${params.row.test}`, {
                                                                        type: 'error',
                                                                    });
                                                                }}
                                                                onSuccess={() => {
                                                                    notify(`Success delete test ${params.row.test}`, {
                                                                        type: 'success',
                                                                    });
                                                                    refresh();

                                                                    const newHistoryRows = history.rows.filter(v => v.id !== params.row.id)
                                                                    setHistory({
                                                                        rows: newHistoryRows,
                                                                        title: history.title,
                                                                    })
                                                                }}
                                                            />
                                                        </Box>
                                                    },
                                                }
                                            ]} />
                                    </Stack>
                                ))
                            }
                            <HistoryDialog
                                pageSizeOptions={[-1]}
                                hideFooter
                                editMode="row"
                                processRowUpdate={updateResult}
                                onProcessRowUpdateError={onUpdateError}
                                rowHeight={60}
                                title={history.title} open={openHistory} onClose={() => setOpenHistory(false)} rows={history.rows} columns={[
                                    ...baseColumns,
                                    {
                                        field: '',
                                        headerName: 'Action',
                                        flex: 1,
                                        renderCell: (params: GridRenderCellParams) => {
                                            return <ButtonGroup sx={{
                                                gap: 2,
                                            }}>
                                                <DeleteButton sx={{
                                                    marginLeft: 2,
                                                }} label={''} mutationMode="pessimistic" size='medium' resource="result" variant='text' record={{
                                                    id: params.row.id,
                                                }} confirmTitle={`Delete test ${params.row.test}?`}
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
                                                            refresh();

                                                            const newHistoryRows = history.rows.filter(v => v.id !== params.row.id)
                                                            setHistory({
                                                                rows: newHistoryRows,
                                                                title: history.title,
                                                            })
                                                        },
                                                    }}
                                                />
                                            </ButtonGroup>
                                        }
                                    }
                                ]} />
                        </>
                    )
                }} />
            </SimpleShowLayout>
        </Show>
    )
}

type DeleteTestButtonProps = {
    ids: number[]
    code: string
    onSuccess?: () => void
    onError?: () => void
} & Partial<ButtonProps>

const DeleteTestButton = (props: DeleteTestButtonProps) => {
    const [open, setOpen] = useState(false);

    const [remove, { isPending }] = useDeleteMany(
        'result',
        { ids: props.ids },
        {
            onError: props.onError,
            onSuccess: props.onSuccess,
        }
    );

    const handleClick = () => setOpen(true);
    const handleDialogClose = () => setOpen(false);
    const handleConfirm = () => {
        remove();
        setOpen(false);
    };

    return (
        <>
            <Tooltip title="Delete Test (All History)">
                <IconButton  {...props} onClick={handleClick} color='error' resource="result" sx={{
                    marginLeft: 2,
                }}>
                    <DeleteIcon />
                </IconButton>
            </Tooltip>
            <Confirm
                isOpen={open}
                loading={isPending}
                title={`Delete test ${props.code}?`}
                content="This will delete the test result. This action cannot be undone."
                confirmColor="warning"
                onConfirm={handleConfirm}
                onClose={handleDialogClose}
            />
        </>
    );
};


type HistoryChangeProps = {
    rows: ResultColumn[]
    title: string
}

type HistoryDialogProps = {
    title: string
    rows: ResultColumn[]
    columns: any[]
    open: boolean
    onClose: () => void
} & Partial<DataGridProps<ResultColumn>>

const HistoryDialog = (props: HistoryDialogProps) => {
    return (
        <Dialog
            open={props.open}
            onClose={props.onClose}
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
            <DialogContent sx={{
            }}>
                <MuiDatagrid  {...props} rows={props.rows} columns={props.columns} />
            </DialogContent>
        </Dialog>
    )
}

export const ObservationResultAdd = () => {
    const FormField = () => {
        const [params] = useSearchParams();
        const specimenID = Number(params.get("specimen_id"));

        return (
            <>
                <NumberInput source="specimen_id" label="Specimen ID" defaultValue={specimenID} readOnly={true} />
                <ArrayInput source="tests" >
                    <SimpleFormIterator inline disableReordering disableClear>
                        <ReferenceInput source="test_type_id" reference="test-type" >
                            <AutocompleteInput optionValue="id" optionText={record => `${record.code} - ${record.category}`}
                                TextFieldProps={{
                                    autoFocus: true,
                                }}
                            />
                        </ReferenceInput >
                        <NumberInput source="value" label="Result" autoFocus />
                    </SimpleFormIterator>
                </ArrayInput>
            </>

        )
    }

    const useReferer = useRefererRedirect("/result");
    return (
        <Create title="Create Observation Result" resource="result" redirect={useReferer}>
            <SimpleForm>
                <FormField />
            </SimpleForm>
        </Create>
    )
}