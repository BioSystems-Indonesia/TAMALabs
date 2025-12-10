import MedicalServicesIcon from '@mui/icons-material/MedicalServices';
import PrintIcon from '@mui/icons-material/Print';
import ScienceIcon from '@mui/icons-material/Science';
import { Avatar, Box, Card, CardContent, Grid, List, ListItem, ListItemAvatar, ListItemText } from "@mui/material";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import React from "react";
import {
    ArrayField,
    Button,
    ChipField,
    Datagrid,
    DateField,
    EditButton,
    Link,
    RecordContextProvider,
    ReferenceArrayField,
    ReferenceField,
    Show,
    SingleFieldList,
    TabbedShowLayout,
    TextField,
    TopToolbar,
    WithRecord,
    WrapperField,
    useGetRecordId
} from "react-admin";
import { useReactToPrint } from "react-to-print";
import LIMSBarcode from '../../component/Barcode';
import useSettings from '../../hooks/useSettings';
import type { BarcodeStyle } from '../../types/general';
import { User } from '../../types/user';
import { workOrderStatusDontShowRun, workOrderStatusShowCancel, type WorkOrder } from '../../types/work_order';
import { PatientForm } from '../patient';
import { WorkOrderStatusChipField } from "./ChipFieldStatus";
import RunWorkOrderForm from './RunWorkOrderForm';

const barcodePageStyle = (style: BarcodeStyle) => `
@media all {
  .page-break {
    display: none;
  }
}

@media print {
  * {
    margin: 0 !important;
    padding: 0 !important;
    box-sizing: border-box !important;
  }
  
  html, body {
    -webkit-print-color-adjust: exact;
    margin: 0 !important;
    padding: 0 !important;
    width: 100% !important;
    height: 100% !important;
  }
  
  body {
    margin: 0 !important;
    padding: 0 !important;
  }

  @page {
    size: ${style.width} ${style.height};
    margin: 0 !important;
    padding: 0 !important;
  }
  
  .barcode-container {
    display: flex !important;
    page-break-before: always;
    flex-direction: column !important;
    justify-content: center !important;
    align-items: center !important;
    transform: rotate(${style.rotate});
    margin: 0 !important;
    padding: 0 !important;
    width: 100vw !important;
    height: 100vh !important;
    box-sizing: border-box !important;
  }
  
  .barcode-container:first-child {
    page-break-before: auto;
  }
  
  .barcode-container p:first-child {
    font-size: 14px !important;
    line-height: 1.1 !important;
  }
  
  .barcode-text {
    margin: 0 !important;
    padding: 0 !important;
    font-size: 12px !important;
    margin-top: 2px !important;
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



function WorkOrderShowActions({ barcodeRef, workOrderID }: { barcodeRef: React.RefObject<any>, workOrderID: number }) {
    return (
        <TopToolbar>
            <PrintBarcodeButton barcodeRef={barcodeRef} />
            <EditButton />
        </TopToolbar>
    )
}


export function WorkOrderShow() {
    const barcodeRef = React.useRef<any>(null);
    const workOrderID = useGetRecordId();
    const [settings] = useSettings();
    const [isProcessing, setIsProcessing] = React.useState(false);

    return (
        <Show actions={<WorkOrderShowActions barcodeRef={barcodeRef} workOrderID={Number(workOrderID)} />}>
            <TabbedShowLayout>
                <TabbedShowLayout.Tab label="Test">
                    <Card elevation={0} sx={{ borderBottom: '2px solid rgba(0, 0, 0, 0.1)' }}>
                        <WithRecord render={(data: WorkOrder) => {
                            return (
                                <CardContent sx={{
                                    display: "flex",
                                    justifyContent: "center",
                                    flexDirection: "column",
                                    gap: 2,
                                }}>
                                    <WorkOrderStatusChipField source='status' />
                                    <RunWorkOrderForm isProcessing={isProcessing} setIsProcessing={setIsProcessing} showCancelButton={workOrderStatusShowCancel.includes(data?.status)}
                                        workOrderIDs={[data?.id]}
                                        showRunButton={!workOrderStatusDontShowRun.includes(data?.status)}
                                        defaultDeviceID={data?.devices && data?.devices.length > 0 ? data?.devices[0]?.id : undefined}
                                    />
                                </CardContent>
                            )
                        }} />
                    </Card>
                    <Card elevation={0} sx={{ borderBottom: '2px solid rgba(0, 0, 0, 0.1)' }}>
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <Typography variant='subtitle1'>Test Info</Typography>
                            <WithRecord render={(workOrder: WorkOrder) => {
                                return (
                                    <ArrayField source={"patient.specimen_list"} label={"Specimen"} textAlign="center" >
                                        <Datagrid bulkActionButtons={false} rowClick={false} hover={false}>
                                            <ChipField source={"type"} textAlign={"center"} />
                                            <WrapperField source={"barcode"} label={"Barcode"} textAlign={"center"}>
                                                <Stack>
                                                    <WithRecord render={(specimen: any) => {
                                                        return (
                                                            <LIMSBarcode
                                                                barcode={specimen.barcode}
                                                                name={workOrder.patient.first_name + " " + workOrder.patient.last_name}
                                                                birthDt={workOrder.patient.birthdate}
                                                                sex={workOrder.patient.sex}
                                                            />
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
                    <Card elevation={0} sx={{ borderBottom: '2px solid rgba(0, 0, 0, 0.1)' }}>
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <Typography variant='subtitle1'>Patient Info</Typography>
                            <WithRecord render={(record: WorkOrder) => {
                                return (
                                    <Box>
                                        <RecordContextProvider value={record.patient}>
                                            <PatientForm readonly mode={"SHOW"} />
                                        </RecordContextProvider>
                                    </Box>
                                )
                            }} />
                        </CardContent>
                    </Card>
                    <Card elevation={0} >
                        <CardContent sx={{
                            display: "flex",
                            justifyContent: "center",
                            flexDirection: "column",
                            gap: 2,
                        }}>
                            <WithRecord render={(record: WorkOrder) => {
                                return (
                                    <Grid container spacing={2}>
                                        <Grid size={6} >
                                            <Stack gap={1}>
                                                <Typography variant='subtitle1' sx={{
                                                    textAlign: "center",
                                                }}>Doctor</Typography>
                                                <RecordContextProvider value={record.doctors}>
                                                    <AdminShow icon={<MedicalServicesIcon />} />
                                                </RecordContextProvider>
                                            </Stack>
                                        </Grid>
                                        <Grid size={6}>
                                            <Stack gap={1}>
                                                <Typography variant='subtitle1' sx={{
                                                    textAlign: "center",
                                                }}>Analyts</Typography>
                                                <RecordContextProvider value={record.analyzers}>
                                                    <AdminShow icon={<ScienceIcon />} />
                                                </RecordContextProvider>
                                            </Stack>
                                        </Grid>
                                    </Grid>
                                )

                            }} />
                        </CardContent>
                    </Card>
                </TabbedShowLayout.Tab>
                <TabbedShowLayout.Tab label="Detail">
                    <TextField source="id" />
                    <TextField source="barcode" label="Barcode" />
                    <TextField source="barcode_simrs" label="Barcode SIMRS" />
                    <TextField source="medical_record_number" label="No. Rekam Medis" />
                    <DateField source="created_at" showTime />
                    <DateField source="updated_at" showTime />
                    <ReferenceField source="created_by" reference='user' />
                    <ReferenceField source="last_updated_by" reference='user' />
                    <ReferenceArrayField source="test_template_ids" reference='test-template' />
                </TabbedShowLayout.Tab>
            </TabbedShowLayout>
            {/*Below is a barcode component for printing only, it will be hidden on the screen*/}
            <WithRecord render={(record: any) => {
                return (
                    <Stack ref={barcodeRef} sx={{
                        display: "none"
                    }}>
                        {
                            record.patient?.specimen_list?.map((specimen: any) => {
                                return (
                                    <LIMSBarcode
                                        barcode={specimen.barcode}
                                        name={record.patient.first_name}
                                        height={settings.barcode_height}
                                        width={settings.barcode_width}
                                        birthDt={record.patient.birthdate}
                                        sex={record.patient.sex}
                                    />
                                    // <Stack gap={0} justifyContent={"center"} alignItems={"center"}
                                    //     className={"barcode-container"}>
                                    //     <Typography
                                    //         className={"barcode-text"}
                                    //         fontSize={8}
                                    //         sx={{
                                    //             margin: 0,
                                    //         }}>{} | {specimen.barcode}</Typography>
                                    //     <Barcode value={specimen.barcode} displayValue={false} height={settings.barcode_height} margin={0} width={settings.barcode_width} />
                                    // </Stack>
                                )
                            })
                        }
                    </Stack>
                )
            }} />
        </Show >
    )
}

type AdminShowProps = {
    icon?: React.ReactNode;
}

function AdminShow(props: AdminShowProps) {
    return (
        <WithRecord render={(record: User[]) => {
            return (
                <List sx={{
                    width: '100%',
                    bgcolor: 'background.paper',
                    borderRadius: 1,
                    boxShadow: 1,
                }}>
                    {record.map((user, index) => (
                        <Link
                            to={`/user/${user.id}`}
                            key={user.id}>
                            <ListItem
                                key={user.id}
                                divider={index !== record.length - 1}
                                sx={{
                                    py: 1.5,
                                    transition: 'all 0.2s ease-in-out',
                                    '&:hover': {
                                        backgroundColor: 'action.hover',
                                        boxShadow: 1,
                                    },
                                    '&:active': {
                                        transform: 'translateX(8px) scale(0.98)',
                                    }
                                }}
                            >
                                <ListItemAvatar>
                                    <Avatar
                                        sx={{
                                            bgcolor: 'primary.main',
                                            transition: 'transform 0.2s ease-in-out',
                                            '&:hover': {
                                                transform: 'scale(1.1)',
                                            }
                                        }}
                                    >
                                        {props.icon}
                                    </Avatar>
                                </ListItemAvatar>
                                <ListItemText
                                    primary={user.fullname}
                                    secondary={user.email}
                                    slotProps={{
                                        primary: {
                                            fontWeight: 'medium',
                                            color: 'text.primary',
                                            sx: {
                                                textDecoration: 'none',
                                            }
                                        },
                                        secondary: {
                                            color: 'text.secondary'
                                        }
                                    }}
                                />
                            </ListItem>
                        </Link>
                    ))}
                </List>
            )
        }} />
    )
}

export type RunWorkOrderFormProps = {
    workOrderIDs?: number[];
}

