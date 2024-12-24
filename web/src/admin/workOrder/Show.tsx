import {
    ArrayField,
    Button,
    Datagrid,
    DateField,
    EditButton,
    ReferenceField,
    ReferenceManyField,
    Show,
    SimpleShowLayout,
    TextField,
    TopToolbar,
    useNotify,
    useRecordContext,
    WithRecord,
    WrapperField
} from "react-admin";
import PrintIcon from '@mui/icons-material/Print';
import Stack from "@mui/material/Stack";
import Barcode from "react-barcode";
import React from "react";
import {useReactToPrint} from "react-to-print";
import Typography from "@mui/material/Typography";
import {useMutation} from "@tanstack/react-query";
import PlayCircleFilledIcon from "@mui/icons-material/PlayCircleFilled";

const barcodePageStyle = `
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
        size: 60mm 45mm;
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
    }
    
    .barcode-text {
        font-size: 12px;
        margin: 0;
    }
}
`


const PrintBarcodeButton = ({barcodeRef}: { barcodeRef: React.RefObject<any> }) => {
    const reactToPrint = useReactToPrint({
        contentRef: barcodeRef,
        pageStyle: barcodePageStyle,
        documentTitle: "Barcode",
        ignoreGlobalStyles: true,
    });

    const handleClick = () => {
        reactToPrint();
    }

    return (
        <Button label="Print Barcode" onClick={handleClick}>
            <PrintIcon/>
        </Button>
    );
}

const RunSingleWorkOrderButton = () => {
        const notify = useNotify();
        const data = useRecordContext()
        const {mutate, isPending} = useMutation({
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
                notify('Success run');
            },
            onError: (error) => {
                notify('Error:' + error.message, {
                    type: 'error',
                });
            },
        })

        const handleClick = () => {
            mutate({
                work_order_ids: [data?.id]
            });
        }

        return (
            <Button label="Run Work Order" onClick={handleClick} disabled={isPending}>
                <PlayCircleFilledIcon/>
            </Button>
        );
    }
;

function WorkOrderShowActions({barcodeRef}: { barcodeRef: React.RefObject<any> }) {
    return (
        <TopToolbar>
            <RunSingleWorkOrderButton/>
            <PrintBarcodeButton barcodeRef={barcodeRef}/>
            <EditButton/>
        </TopToolbar>
    )
}

export function WorkOrderShow() {
    const barcodeRef = React.useRef<any>(null);

    return (
        <Show actions={<WorkOrderShowActions barcodeRef={barcodeRef}/>}>
            <SimpleShowLayout>
                <TextField source="id"/>
                <TextField source="status"/>
                <TextField source="created_at"/>
                <TextField source="updated_at"/>
                <ReferenceManyField reference={"feature-list-observation-type"} target={"id"}
                                    source={"observation_requests"}>
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="id"/>
                        <TextField source="name"/>
                        <TextField source="additional_info.type"/>
                    </Datagrid>
                </ReferenceManyField>
                <ArrayField source={"specimen_list"} label={"Specimens"}>
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="id"/>
                        <ReferenceField reference={"patient"} source={"patient_id"}/>
                        <TextField source="type"/>
                        <WrapperField source={"barcode"} label={"Barcode"} textAlign={"center"}>
                            <Stack>
                                <WithRecord render={(record: any) => {
                                    return (
                                        <Stack gap={0} justifyContent={"center"} alignItems={"center"}>
                                            <Barcode value={record.barcode} displayValue={false}/>
                                            <Typography
                                                className={"barcode-text"}
                                                fontSize={12}
                                                sx={{
                                                    margin: 0,
                                                }}>{record.barcode}</Typography>
                                        </Stack>
                                    )
                                }}/>
                            </Stack>
                        </WrapperField>
                        <DateField source="created_at"/>
                        <DateField source="updated_at"/>
                    </Datagrid>
                </ArrayField>
                {/*Below is a barcode component for printing only, it will be hidden on the screen*/}
                <WithRecord render={(record: any) => {
                    return (
                        <Stack ref={barcodeRef} sx={{
                            display: "none",
                        }}>
                            {record?.specimen_list?.map((specimen: any) => {
                                return (
                                    <Stack gap={0} justifyContent={"center"} alignItems={"center"}
                                           className={"barcode-container"}>
                                        <Typography
                                            className={"barcode-text"}
                                            fontSize={12}
                                            sx={{
                                                margin: 0,
                                            }}>{specimen.patient.first_name} {specimen.patient.last_name}</Typography>
                                        <Barcode value={specimen.barcode} displayValue={false}/>
                                        <Typography
                                            className={"barcode-text"}
                                            fontSize={12}
                                            sx={{
                                                margin: 0,
                                            }}>{specimen.barcode}</Typography>
                                    </Stack>
                                )
                            })}
                        </Stack>
                    )
                }}/>
            </SimpleShowLayout>
        </Show>
    )
}

