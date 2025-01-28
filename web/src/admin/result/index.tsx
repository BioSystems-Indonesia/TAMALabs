import {
    AutocompleteArrayInput,
    AutocompleteInput,
    BooleanInput,
    Button,
    CheckboxGroupInput,
    Create,
    Datagrid,
    DateField,
    DeleteButton,
    Edit,
    FilterList,
    FilterListItem,
    FilterListSection,
    FilterLiveForm,
    FilterLiveSearch,
    Labeled,
    Link,
    List,
    NumberField,
    NumberInput,
    ReferenceInput,
    SelectArrayInput,
    SimpleForm,
    TextField,
    WithRecord,
    WrapperField,
    useNotify,
    useRefresh,
} from "react-admin";
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import { Card, CardContent, Chip, Grid, Stack, Typography } from "@mui/material";
import { DataGrid as MuiDatagrid, type GridRenderCellParams } from '@mui/x-data-grid';
import type { ResultColumn } from "../../types/general";
import AddIcon from '@mui/icons-material/Add';
import { useSearchParams } from "react-router-dom";
import { useFormContext } from "react-hook-form";
import { useEffect, useState } from "react";
import { getRefererParam, useRefererRedirect } from "../../hooks/useReferer";
import PrintMCUButton from "../../component/PrintReport";
import BeenhereIcon from '@mui/icons-material/Beenhere';
import ChecklistRtlIcon from '@mui/icons-material/ChecklistRtl';
import FeatureList from "../../component/FeatureList";


const ResultFilterSidebar = () => {
    return (
        <Card sx={{
            order: -1, mr: 2, mt: 2, width: 200, minWidth: 200,
        }}>
            <CardContent>
                <FilterLiveForm>
                    <FeatureList source={"work_order_status"} types={"work-order-status"}>
                        <CheckboxGroupInput />
                    </FeatureList>
                    <ReferenceInput source={"work_order_ids"} reference="work-order" label={"Work Order"}>
                        <AutocompleteArrayInput />
                    </ReferenceInput>
                    <ReferenceInput source={"patient_ids"} reference="patient" label={"Patient"}>
                        <AutocompleteArrayInput />
                    </ReferenceInput>
                    <BooleanInput source={"has_result"} label={"Show Only With Result"} />
                </FilterLiveForm>
            </CardContent>
        </Card>
    )
}

export const ResultList = () => (
    <List resource="result" sort={{
        field: "id",
        order: "DESC"
    }} aside={<ResultFilterSidebar />} exporter={false} >
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.order_id}/show`} resource="patient" label={"Patient"} onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )} />
            <WithRecord label="Work Order" render={(record: any) => (
                <Link to={`/work-order/${record.order_id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.order_id} - ${record.work_order.status}`} color={WorkOrderChipColorMap(record.work_order.status)} />
                </Link>
            )} />
            <TextField source="barcode" />
            <WithRecord label="Request" render={(record: any) => (
                <Typography variant="body2" >
                    {record.observation_requests.length}
                </Typography>
            )} />
            <WithRecord label="Result" render={(record: any) => (
                <Typography variant="body2" >
                    {record.observation_result.length}
                </Typography>
            )} />
            <DateField source="created_at" showDate showTime />
            <WithRecord label="Print Result" render={(record: any) => (
                <PrintMCUButton results={record.observation_result} patient={record.patient} workOrder={record.work_order} />
            )} />
        </Datagrid>
    </List>
);




export const ResultEdit = () => {
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

            return respJSON[0]
        }

    }

    const refresh = useRefresh();

    return (
        <Edit title="Edit Result">
            <SimpleForm toolbar={false}>
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
                                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                                </Link>
                            )} />
                        </Labeled>
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <Labeled>
                            <WithRecord label="Work Order" render={(record: any) => (
                                <Link to={`/work-order/${record.order_id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                                    <Chip label={`#${record.order_id} - ${record.work_order.status}`} color={WorkOrderChipColorMap(record.work_order.status)} />
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
                        return <Typography variant="subtitle1" sx={{
                            width: "100%",
                            textAlign: "center",
                            marginY: 2,
                        }}>No Test Result found</Typography>
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
                                            columns={[
                                                {
                                                    field: 'test',
                                                    headerName: 'Test',
                                                    flex: 1,
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
                                                {
                                                    field: 'action',
                                                    headerName: 'Action',
                                                    flex: 1,
                                                    renderCell: (params: GridRenderCellParams) => {
                                                        return <DeleteButton mutationMode="pessimistic" resource="result" record={{
                                                            id: params.row.id,
                                                            code: params.row.test,
                                                        }} confirmTitle={`Delete test ${params.row.test}?`}
                                                            confirmColor="warning"
                                                            confirmContent="This will delete the test result. This action cannot be undone."
                                                            redirect={false}
                                                            mutationOptions={{
                                                                onSettled: () => {
                                                                    notify(`Success delete test ${params.row.test}`, {
                                                                        type: 'success',
                                                                    });
                                                                    refresh();
                                                                },
                                                            }}
                                                        />
                                                    },
                                                }
                                            ]} />
                                    </Stack>
                                ))
                            }
                        </>
                    )
                }} />
                {/* <WithRecord label="Test Result" render={(record: any) => (
                <>
                    {
                        Object.entries(record.test_result).map(([category, result]: any) => (
                            <Stack sx={{
                                marginY: 2,
                            }}>
                                <Typography variant="subtitle1" gutterBottom>
                                    {category}
                                </Typography>
                                <Table>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>Name</TableCell>
                                            <TableCell>Result</TableCell>
                                            <TableCell>Unit</TableCell>
                                            <TableCell>Reference Range</TableCell>
                                            <TableCell>Status</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {result.map((test: any, index: number) => (
                                            <TableRow key={index}>
                                                <TableCell>{test.test}</TableCell>
                                                <TableCell>{test.result}</TableCell>
                                                <TableCell>{test.unit}</TableCell>
                                                <TableCell>{test.reference_range}</TableCell>
                                                <TableCell>
                                                    {test.abnormal === 1 ? (
                                                        <Typography color="error">Abnormal</Typography>
                                                    ) : (
                                                        <Typography color="primary">Normal</Typography>
                                                    )}
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </Stack>
                        ))
                    }
                </>
            )} /> */}
            </SimpleForm>
        </Edit>
    )
}

export const ObservationResultAdd = () => {

    const FormField = () => {
        const [params] = useSearchParams();
        const [specimenReadonly, setSpecimenReadonly] = useState(false);
        const { setValue } = useFormContext()

        useEffect(() => {
            if (params.get("specimen_id")) {
                setSpecimenReadonly(true);
                setValue("specimen_id", Number(params.get("specimen_id")));
            }
        }, [params])

        return (
            <>
                <NumberInput source="specimen_id" label="Specimen ID" readOnly={specimenReadonly} />
                <ReferenceInput source="code" reference="test-type" >
                    <AutocompleteInput optionValue="code" optionText={record => `${record.code} - ${record.category}`} />
                </ReferenceInput >
                <NumberInput source="value" label="Result" />
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
