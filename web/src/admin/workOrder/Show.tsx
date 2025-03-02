import CancelIcon from '@mui/icons-material/Cancel';
import PlayCircleFilledIcon from "@mui/icons-material/PlayCircleFilled";
import PrintIcon from '@mui/icons-material/Print';
import { Card, CardContent, CircularProgress, Grid, type SxProps } from "@mui/material";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useMutation } from "@tanstack/react-query";
import React from "react";
import {
    ArrayField,
    AutocompleteInput,
    BooleanInput,
    Button,
    ChipField,
    Datagrid,
    DateField,
    EditButton,
    Form,
    InputHelperText,
    Link,
    RecordContextProvider,
    ReferenceInput,
    Show,
    SingleFieldList,
    TabbedShowLayout,
    TextField,
    TopToolbar,
    WithRecord,
    WrapperField,
    required,
    useGetRecordId,
    useNotify,
    useRecordContext,
    useRefresh
} from "react-admin";
import Barcode from "react-barcode";
import { useFormContext } from 'react-hook-form';
import { useReactToPrint } from "react-to-print";
import { trimName } from '../../helper/format';
import useAxios from '../../hooks/useAxios';
import { getRefererParam } from '../../hooks/useReferer';
import useSettings from '../../hooks/useSettings';
import type { BarcodeStyle } from '../../types/general';
import type { WorkOrder } from '../../types/work_order';
import { DeviceForm } from "../device";
import { PatientForm } from '../patient';
import { WorkOrderStatusChipField } from "./ChipFieldStatus";

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
    body {
        margin: 0;
    }

    @page {
        size: ${style.width} ${style.height};
        margin: 0;
        text-align: center;
    }
    
    .barcode-container {
        display: flex;
        page-break-before: always;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        transform: rotate(${style.rotate});
    }
    
    .barcode-text {
        font-size: 8px;
        margin: 0;
    }
}`

const PrintBarcodeButton = ({ barcodeRef }: { barcodeRef: React.RefObject<any> }) => {
    const [settings] = useSettings();
    const reactToPrint = useReactToPrint({
        contentRef: barcodeRef,
        pageStyle: barcodePageStyle({
            width: `${settings.barcode_page_width}mm`,
            height: `${settings.barcode_page_height}mm`,
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
            <EditButton />
        </TopToolbar>
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
                            <WorkOrderStatusChipField source='status' />
                            <RunWorkOrderForm />
                        </CardContent>
                    </Card>
                    <Card>
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <Typography variant='subtitle1'>Patient Info</Typography>
                            <WithRecord render={(record: WorkOrder) => {
                                return (
                                    <RecordContextProvider value={record.patient}>
                                        <PatientForm readonly mode={"SHOW"} />
                                    </RecordContextProvider>
                                )
                            }} />
                        </CardContent>
                    </Card>
                    <Card>
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <Typography variant='subtitle1'>Test Info</Typography>
                            <WithRecord render={(record: WorkOrder) => {

                                return (
                                    <ArrayField source={"patient.specimen_list"} label={"Specimen"} textAlign="center" >
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
                                                <SingleFieldList linkType={false} sx={{
                                                    maxHeight: "200px",
                                                    overflow: "scroll",
                                                }}>
                                                    <ChipField source={"test_code"} textAlign={"center"} />
                                                </SingleFieldList>
                                            </ArrayField>
                                        </Datagrid>
                                    </ArrayField>
                                )
                            }} />
                        </CardContent>
                    </Card>
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
                        {
                            record.patient?.specimen_list?.map((specimen: any) => {
                                return (
                                    <Stack gap={0} justifyContent={"center"} alignItems={"center"}
                                        className={"barcode-container"} sx={{
                                            display: "none",
                                        }}>
                                        <Typography
                                            className={"barcode-text"}
                                            fontSize={8}
                                            sx={{
                                                margin: 0,
                                            }}>{trimName(`${record.patient.first_name} ${record.patient.last_name}`, 14)} | {specimen.barcode}</Typography>
                                        <Barcode value={specimen.barcode} displayValue={false} height={settings.barcode_height} margin={0} width={settings.barcode_width} />
                                    </Stack>
                                )
                            })
                        }
                    </Stack>
                )
            }} />
        </Show >
    )
}

export type RunWorkOrderFormProps = {
    workOrderIDs?: number[];
}

export function RunWorkOrderForm(props: RunWorkOrderFormProps) {
    const notify = useNotify();
    const refresh = useRefresh();
    const axios = useAxios();
    const { mutate, isPending } = useMutation({
        mutationFn: async (data: any) => {
            const response = await axios.post(`${import.meta.env.VITE_BACKEND_BASE_URL}/work-order/run`, data, {
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

    const onSubmit = (data: any) => {
        if (!data.device_id) {
            notify('Please select device to run', {
                type: 'error',
            });
            return;
        }

        mutate({
            work_order_ids: props.workOrderIDs ?? [data.id],
            device_id: data.device_id,
            urgent: data.urgent,
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
                <ReferenceInput source={"device_id"} reference={"device"} disabled={isPending} >
                    <AutocompleteInput source={"device_id"} validate={[required()]} create={<DeviceForm />} sx={{
                        margin: 0,
                    }} disabled={isPending} helperText={
                        <Link to={"/device/create?" + getRefererParam()}>
                            <InputHelperText helperText="Create new device"></InputHelperText>
                        </Link>
                    }
                    />
                </ReferenceInput>
            </Grid>
            <Grid item xs={12} md={3} sx={{
                display: "flex",
                paddingLeft: "24px",
                justifyContent: "start",
                alignItems: "center",
            }}>
                <BooleanInput source={"urgent"} disabled={isPending} label="Urgent" />
            </Grid>
        </Grid>
        <RunWorkOrderSubmit isPending={isPending} sx={{
            marginTop: "12px",
        }} />
    </Form>;
}

function RunWorkOrderSubmit({ isPending, sx }: { isPending: boolean, sx: SxProps }) {
    const { watch } = useFormContext()


    return (
        <Stack sx={sx}>
            <Button label="Run Work Order" disabled={isPending || !watch("device_id")} variant="contained" type='submit' sx={{
                cursor: "pointer"
            }} >
                {isPending ? <CircularProgress size={12} variant='indeterminate' color='primary' /> : <PlayCircleFilledIcon />}
            </Button>
            {!isPending && !watch("device_id") && <Typography color='error' fontSize={12}>Please pick device to run</Typography>}
        </Stack>
    );
}

