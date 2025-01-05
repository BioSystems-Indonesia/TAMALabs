import AddIcon from '@mui/icons-material/Add';
import CancelIcon from '@mui/icons-material/Cancel';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import PlayCircleFilledIcon from "@mui/icons-material/PlayCircleFilled";
import PrintIcon from '@mui/icons-material/Print';
import { Card, CardContent, Grid } from "@mui/material";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useMutation } from "@tanstack/react-query";
import React from "react";
import {
    ArrayField,
    AutocompleteInput,
    Button,
    ChipField,
    Datagrid,
    DateField,
    Form,
    InputHelperText,
    Link,
    ReferenceField,
    ReferenceInput,
    Show,
    SingleFieldList,
    TabbedShowLayout,
    TextField,
    TopToolbar,
    WithRecord,
    WrapperField,
    required,
    useDeleteMany,
    useGetRecordId,
    useListContext,
    useNotify,
    useRecordContext,
    useRefresh,
    type ButtonProps} from "react-admin";
import Barcode from "react-barcode";
import { useReactToPrint } from "react-to-print";
import { DeviceForm } from "../device";
import { WorkOrderStatusChipField } from "./ChipFieldStatus";
import type { BarcodeStyle } from '../../types/general';
import useSettings from '../../hooks/useSettings';

const barcodePageStyle = (style: BarcodeStyle) => `
@media all {
  .page-break {
    display: none;
  }
}

@media print {
  html, body {
    -webkit-print-color-adjust: exact;
  }
}

@media print {
    @page {
        size: ${style.width} ${style.height};
        margin: 0;
        text-align: center;
    }
    
    .barcode-container {
        display: flex;
        padding-top: 5px;
        page-break-before: always;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        transform: rotate(${style.rotate});
    }
    
    .barcode-text {
        font-size: 12px;
        margin: 0;
    }
}`

const PrintBarcodeButton = ({ barcodeRef }: { barcodeRef: React.RefObject<any> }) => {
    const [settings] = useSettings();
    const reactToPrint = useReactToPrint({
        contentRef: barcodeRef,
        pageStyle: barcodePageStyle({
            width: `${settings.barcode_size_width}mm`,
            height: `${settings.barcode_size_height}mm`,
            rotate: `${settings.barcode_orientation === "portrait" ? "-90" : "0"}deg`,
        }),
        documentTitle: "Barcode",
        ignoreGlobalStyles: true,
    });

    const handleClick = () => {
        reactToPrint();
    }

    return (
        <Button label="Print Barcode" onClick={handleClick}>
            <PrintIcon />
        </Button>
    );
}

const AddTestButton = (props: ButtonProps) => {
    const record = useRecordContext();

    return (
        <Link to={`/work-order/${record?.id}/add-test`}>
            <Button label="Add Test" {...props}>
                <AddIcon />
            </Button>
        </Link>
    );
}

const BulkEditButton = ({ patientIDs, workOrderID }: { patientIDs?: number[], workOrderID?: number }) => {
    const recordId = workOrderID ?? useGetRecordId();

    const generateUrl = (patientIDs?: number[]) => {
        const urlParams = new URLSearchParams();
        if (patientIDs) {
            patientIDs.forEach(id => urlParams.append('patient_id', id.toString()));
        }

        return `/work-order/${recordId}/add-test/1?${urlParams.toString()}`;
    }

    return (
        <Link to={generateUrl(patientIDs)}>
            <Button label="Edit" color="secondary" variant="contained">
                <EditIcon />
            </Button>
        </Link>
    );
}

const BulkDeleteButton = ({ patientIDs, workOrderID }: { patientIDs?: number[], workOrderID?: number }) => {
    const recordId = workOrderID ?? useGetRecordId();
    const notify = useNotify();
    const refresh = useRefresh();

    const [deleteMany, { isPending }] = useDeleteMany(`work-order/${recordId}/test`, { ids: patientIDs }, {
        onError: (error: Error) => {
            notify('Error:' + error.message, {
                type: 'error',
            });
        },
        onSuccess: () => {
            refresh();
            notify('Success delete', {
                type: 'success',
            });
        },
    });

    return (
        <Button label="Delete" onClick={() => {
            deleteMany();
        }} disabled={isPending} color="error" variant="contained">
            <DeleteIcon />
        </Button>
    );
}

const CancelButton = ({ workOrderID }: { workOrderID: number }) => {
    const notify = useNotify();
    const refresh = useRefresh();
    const { mutate, isPending } = useMutation({
        mutationFn: async (data: any) => {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/work-order/cancel`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });
            if (!response.ok) {
                const responseJson = await response.json();
                throw new Error(responseJson.error);
            }
            return response.json();
        },
        onSuccess: () => {
            notify('Success cancel', {
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
        <Button label="Cancel" onClick={() => {
            mutate({
                work_order_id: workOrderID
            });
        }} disabled={isPending} color="error" variant="contained">
            <CancelIcon />
        </Button>
    );
}


function WorkOrderShowActions({ barcodeRef, workOrderID }: { barcodeRef: React.RefObject<any>, workOrderID: number }) {
    const data = useRecordContext()
    return (
        <TopToolbar>
            {data?.status === "PENDING" &&
                <CancelButton workOrderID={workOrderID} />
            }
            <PrintBarcodeButton barcodeRef={barcodeRef} />
            <AddTestButton />
        </TopToolbar>
    )
}

function PatientListBulkAction() {
    const { selectedIds } = useListContext();

    return (
        <>
            <BulkEditButton patientIDs={selectedIds} />
            <BulkDeleteButton patientIDs={selectedIds} />
        </>
    )
}

const PatientTestEmpty = () => {
    return (
        <Card>
            <CardContent sx={{
                display: "flex",
                justifyContent: "center",
                flexDirection: "column",
                gap: 2,
            }}>
                <Typography fontSize={24} >No patient test found</Typography>
                <AddTestButton variant="contained" color="primary" />
            </CardContent>
        </Card>
    )
}


export function WorkOrderShow() {
    const barcodeRef = React.useRef<any>(null);
    const workOrderID = useGetRecordId();
    const [settings] = useSettings();

    return (
        <Show actions={<WorkOrderShowActions barcodeRef={barcodeRef} workOrderID={Number(workOrderID)} />}>
            <TabbedShowLayout>
                <TabbedShowLayout.Tab label="Test">
                    <Card>
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <WorkOrderStatusChipField />
                            <RunWorkOrderForm />
                        </CardContent>
                    </Card>
                    <ArrayField source={"patient_list"} label={"Patient Test"} >
                        <Datagrid rowClick={false} bulkActionButtons={<PatientListBulkAction />} empty={<PatientTestEmpty />}>
                            <ReferenceField reference={"patient"} source={"id"} label={"Patient"} textAlign="center" />
                            <ArrayField source={"specimen_list"} label={"Specimen"} textAlign="center" >
                                <Datagrid bulkActionButtons={false} rowClick={false} hover={false}>
                                    <ChipField source={"type"} textAlign={"center"} />
                                    <WrapperField source={"barcode"} label={"Barcode"} textAlign={"center"}>
                                        <Stack>
                                            <WithRecord render={(record: any) => {
                                                return (
                                                    <Stack gap={0} justifyContent={"center"} alignItems={"center"}>
                                                        <Barcode value={record.barcode} displayValue={false} />
                                                        <Typography
                                                            className={"barcode-text"}
                                                            fontSize={12}
                                                            sx={{
                                                                margin: 0,
                                                            }}>{record.barcode}</Typography>
                                                    </Stack>
                                                )
                                            }} />
                                        </Stack>
                                    </WrapperField>
                                    <ArrayField source={"observation_requests"} label={`Observation Requests`} textAlign="center">
                                        <SingleFieldList linkType="false" sx={{
                                            maxHeight: "200px",
                                            overflow: "scroll",
                                        }}>
                                            <ChipField source={"test_code"} textAlign={"center"} clickable={false} />
                                        </SingleFieldList>
                                    </ArrayField>
                                </Datagrid>
                            </ArrayField>
                            <DateField source="updated_at" showTime textAlign="center" />
                            <WrapperField label="Actions">
                                <WithRecord render={(record: any) => {
                                    return (
                                        <Stack gap={1}>
                                            <BulkEditButton patientIDs={[record.id]} workOrderID={Number(workOrderID)} />
                                            <BulkDeleteButton patientIDs={[record.id]} workOrderID={Number(workOrderID)} />
                                        </Stack>
                                    )
                                }} />
                            </WrapperField>
                        </Datagrid>
                    </ArrayField>
                </TabbedShowLayout.Tab>
                <TabbedShowLayout.Tab label="Detail">
                    <TextField source="id" />
                    <DateField source="created_at" showTime />
                    <DateField source="updated_at" showTime />
                </TabbedShowLayout.Tab>
            </TabbedShowLayout>
            {/*Below is a barcode component for printing only, it will be hidden on the screen*/}
            <WithRecord render={(record: any) => {
                return (
                    <Stack ref={barcodeRef}>
                        {record?.patient_list?.map((patient: any) => {
                            return patient?.specimen_list?.map((specimen: any) => {
                                // The width of the barcode paper is approximately 3.91 cm.
                                // The height of the barcode paper is approximately 2.93 cm.
                                return (
                                    <Stack gap={0} justifyContent={"center"} alignItems={"center"}
                                        className={"barcode-container"} sx={{
                                            display: "none",
                                        }}>
                                        <Typography
                                            className={"barcode-text"}
                                            fontSize={12}
                                            sx={{
                                                margin: 0,
                                            }}>{patient.first_name} {patient.last_name}</Typography>
                                        <Barcode value={specimen.barcode} displayValue={false} />
                                        <Typography
                                            className={"barcode-text"}
                                            fontSize={12}
                                            sx={{
                                                margin: 0,
                                            }}>{specimen.barcode}</Typography>
                                    </Stack>
                                )
                            })
                        })}
                    </Stack>
                )
            }} />
        </Show>
    )
}

function RunWorkOrderForm() {
    const notify = useNotify();
    const refresh = useRefresh();
    const { mutate, isPending } = useMutation({
        mutationFn: async (data: any) => {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/work-order/run`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });

            if (!response.ok) {
                const responseJson = await response.json();
                throw new Error(responseJson.error);
            }

            return response.json();
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

    const onSubmit = (data: any) => {
        if (!data.device_id) {
            notify('Please select device to run', {
                type: 'error',
            });
            return;
        }

        mutate({
            work_order_id: data.id,
            device_id: data.device_id
        });
    }

    return <Form disabled={isPending} onSubmit={onSubmit}>
        <Grid direction={"row"} sx={{
            width: "100%",
        }} container>
            <Grid item xs={12} md={9} sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
            }}>
                <Stack width={"100%"}>
                    <ReferenceInput source={"device_id"} reference={"device"}>
                        <AutocompleteInput source={"device_id"} validate={[required()]} create={<DeviceForm />} sx={{
                            margin: 0,
                        }} helperText={
                            <Link to={"/device/create"} target="_blank" rel="noopener noreferrer">
                                <InputHelperText helperText="Create new device"></InputHelperText>
                            </Link>
                        }
                        />
                    </ReferenceInput>
                </Stack>
            </Grid>
            <Grid item xs={12} md={3} sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
            }}>
                <Button label="Run Work Order" disabled={isPending} variant="contained" type="submit">
                    <PlayCircleFilledIcon />
                </Button>
            </Grid>
        </Grid>
    </Form>;
}

