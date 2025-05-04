import AddIcon from '@mui/icons-material/Add';
import { Box, Button, Checkbox, Dialog, DialogContent, DialogTitle, MenuItem, Select, type ButtonProps, GridLegacy as Grid, ListItem, ListItemText, Paper, Tooltip } from "@mui/material";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import MUIList from "@mui/material/List";
import { DataGrid as MuiDatagrid, useGridApiRef, type GridRenderCellParams } from "@mui/x-data-grid";
import TouchAppIcon from '@mui/icons-material/TouchApp';
import React, { useEffect, useState } from "react";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import {
    AutocompleteInput,
    DateInput,
    Form,
    List,
    ReferenceInput,
    SaveButton,
    SimpleForm,
    TextInput,
    Toolbar,
    useCreate,
    useListContext,
    useNotify,
    useRecordContext,
    useSaveContext
} from "react-admin";
import { useFormContext } from "react-hook-form";
import useAxios from "../../hooks/useAxios.ts";
import type { ObservationRequest, ObservationRequestCreateRequest } from "../../types/observation_requests.ts";
import { ActionKeys } from "../../types/props.ts";
import type { Specimen } from '../../types/specimen.ts';
import type { TestType, TestTypeSpecimenType } from "../../types/test_type.ts";
import { WorkOrder } from '../../types/work_order.ts';
import { PatientFormField } from "../patient/index.tsx";
import FormStepper from "./Stepper.tsx";
import { TestFilterSidebar } from "./TestTypeFilter.tsx";
import ListAltIcon from '@mui/icons-material/ListAlt';


type WorkOrderActionKeys = ActionKeys | "ADD_TEST";

type WorkOrderFormProps = {
    readonly?: boolean
    mode: WorkOrderActionKeys
}

function TestSelectorPrompt() {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                padding: 1,
                margin: 1,
                border: '2px',
                borderRadius: 2,
                textAlign: 'center',
                maxWidth: '200px',
                backgroundColor: 'background.paper',
                boxShadow: 1,
            }}
        >
            <TouchAppIcon color="primary" sx={{ fontSize: 32, mb: 1 }} />
            <Typography variant="h6" component="p" gutterBottom>
                Please select tests
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Check the tests you want to run
            </Typography>
            <Divider sx={{ width: '80%', mb: 2 }} />

            <Box sx={{ textAlign: 'left', width: '100%', px: 2 }}>
                <Typography variant="subtitle1" gutterBottom sx={{ display: 'flex', alignItems: 'center', color: 'secondary.main' }}>
                    <InfoOutlinedIcon sx={{ mr: 1, fontSize: '1.2rem' }} /> Quick Tips:
                </Typography>
                <MUIList dense>
                    <ListItem disablePadding>
                        <ListItemText
                            primary="You can filter with the sidebar"
                            primaryTypographyProps={{ variant: 'body2', color: 'text.secondary' }}
                        />
                    </ListItem>
                    <ListItem disablePadding>
                        <ListItemText
                            primary={<p>For recurring tests across multiple patients, consider using <u>Test Templates</u></p>}
                            primaryTypographyProps={{ variant: 'body2', color: 'text.secondary' }}
                        />
                    </ListItem>
                </MUIList >
            </Box>
        </Box>
    );
}


export const testTypesField = "test_types";


const PickedTest = ({ selectedData }: { selectedData: Record<number, ObservationRequestCreateRequest> }) => {
    if (Object.keys(selectedData).length === 0) {
        return (
            <TestSelectorPrompt />
        )
    }

    return (
        <Paper elevation={2} sx={{ p: 2, backgroundColor: 'background.paper' /* Or 'grey.50' for slight contrast */ }}>
            <Stack spacing={2}>
                {/* Enhanced Title Section */}
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <ListAltIcon color="primary" sx={{ mr: 1 }} /> {/* Icon before title */}
                    <Typography variant="h6" fontSize={18} fontWeight="medium"> {/* More prominent title */}
                        Selected Tests
                    </Typography>
                </Box>
                <Divider/>
                {/* Grid container for the chips */}
                <Grid container spacing={1}> {/* Adjust spacing as needed */}
                    {Object.entries(selectedData).map(([key, value]) => {
                        // Convert the key back to a number if you need it for handlers like onRemoveTest
                        const testId = Number(key);

                        // Determine the text to show in the tooltip (use full name if available, otherwise code)
                        const tooltipTitle = `${value.test_type_code} (${value.specimen_type})`;

                        return (
                            <Grid item key={testId}>
                                <Tooltip title={tooltipTitle} arrow placement="top">
                                    <Chip
                                        label={value.test_type_code}
                                        color="primary"
                                        variant="filled"
                                        size="medium"
                                    />
                                </Tooltip>
                            </Grid>
                        );
                    })}
                </Grid>
            </Stack>
        </Paper>
    )
}

const patientIDField = "patient_id";

type TableTestType = TestType & {
    checked: boolean
    picked_type: string
}

type TestTableProps = {
    setSelectedData: React.Dispatch<React.SetStateAction<Record<number, ObservationRequestCreateRequest>>>
    selectedData: Record<number, ObservationRequestCreateRequest>
} & TestInputProps

function TestTable({
    setSelectedData,
    selectedData,
    ...props
}: TestTableProps) {
    const data = useRecordContext<WorkOrder>();
    const { setValue } = useFormContext();
    const { data: testType, isPending } = useListContext();
    const [rows, setRows] = useState<TableTestType[]>([]);
    const [allTestType, setAllTestType] = useState<TestType[] | undefined>(undefined);
    const [alreadySetTestType, setAlreadySetTestType] = useState<boolean>(false);

    useEffect(() => {
        if (testType && !alreadySetTestType) {
            setAllTestType(testType);
            setAlreadySetTestType(true);
        }

        if (testType) {
            setRows(testType)
        }
    }, [testType]);

    useEffect(() => {
        if (props.initSelectedType && allTestType) {
            const initSelectedType = props.initSelectedType;
            setSelectedData(initSelectedType);
        }
    }, [props.initSelectedType, allTestType])

    useEffect(() => {
        setValue(testTypesField, selectedData);

        if (testType && selectedData) {
            setRows(rows => {
                const newRows = rows.map(row => {
                    if (selectedData[row.id]) {
                        return {
                            ...row,
                            checked: true,
                            picked_type: selectedData[row.id].specimen_type,
                        }
                    }

                    return {
                        ...row,
                        checked: false,
                        picked_type: row.types[0].type
                    }
                })

                return newRows
            })
        }
    }, [selectedData, testType])

    useEffect(() => {
        if (data && data.patient && allTestType) {
            const observationRequestCodeList = data.patient?.specimen_list?.map((specimen: Specimen) => {
                return specimen.observation_requests.map((observationRequest: ObservationRequest) => {
                    return {
                        test_type_id: observationRequest.test_type.id,
                        specimen_type: specimen.type,
                        test_type_code: observationRequest.test_type.code,
                    } as ObservationRequestCreateRequest;
                })
            }).flat()
            if (observationRequestCodeList?.length === 0) {
                return;
            }

            const observationRequestMap: Record<string, ObservationRequestCreateRequest> = {};
            observationRequestCodeList?.forEach((v: ObservationRequestCreateRequest) => {
                observationRequestMap[v.test_type_id] = v;
            })

            setSelectedData(observationRequestMap);
        }
    }, [data, allTestType]);

    function updateSelectedData(testType: TableTestType) {
        setSelectedData(val => {
            const newSelectedData = { ...val }

            if (testType.checked) {
                const newData: ObservationRequestCreateRequest = {
                    test_type_id: testType.id,
                    test_type_code: testType.code,
                    specimen_type: testType.picked_type ?? testType.types[0].type,
                }

                newSelectedData[testType.id] = newData
            } else {
                delete newSelectedData[testType.id]
            }

            return newSelectedData
        });
    }

    const apiRef = useGridApiRef();

    return <Grid container spacing={2}>
        <Grid item xs={12} md={9}>
            <MuiDatagrid
                pageSizeOptions={[-1]}
                hideFooter
                editMode="row"
                apiRef={apiRef}
                loading={isPending}
                columns={[
                    {
                        field: 'checked',
                        headerName: 'Checked',
                        flex: 1,
                        filterable: false,
                        renderCell: (params: GridRenderCellParams) => {
                            return <Checkbox checked={params.row.checked} onChange={(e: any) => {
                                const newRow = { ...params.row, checked: e.target.checked }
                                updateSelectedData(newRow)
                            }} />
                        },
                    },
                    {
                        field: 'name',
                        headerName: 'Name',
                        flex: 1,
                        filterable: false,
                    },
                    {
                        field: 'category',
                        headerName: 'Category',
                        flex: 1,
                        filterable: false,
                    },
                    {
                        field: 'sub_category',
                        headerName: 'Sub Category',
                        flex: 1,
                    },
                    {
                        field: 'picked_type',
                        headerName: 'Specimen Type',
                        flex: 1,
                        renderCell: (params: GridRenderCellParams) => {
                            const testType = params.row as TestType

                            if (testType.types.length == 1) {
                                return testType.types[0].type
                            }

                            return <Select defaultValue={testType.types[0]?.type} value={params.row.picked_type} onChange={(e) => {
                                const newRow = { ...params.row, picked_type: e.target.value }
                                updateSelectedData(newRow)
                            }} sx={{
                                padding: '0px',
                                maxHeight: "30px",
                            }}>
                                {testType.types.map((type: TestTypeSpecimenType) => {
                                    return <MenuItem value={type.type}>{type.type}</MenuItem>
                                })}
                            </Select>
                        },
                    }
                ]}
                disableColumnSelector
                rows={rows}
            />
        </Grid>
        <Grid item xs={12} md={3}>
            <PickedTest selectedData={selectedData} />
        </Grid>
    </Grid>;
}

type TestInputProps = {
    initSelectedType?: Record<number, ObservationRequestCreateRequest>
}

export function TestInput(props: TestInputProps) {
    const [selectedData, setSelectedData] = useState<Record<number, ObservationRequestCreateRequest>>({});
    return (
        <Box sx={{
            maxHeight: "calc(70vh - 48px)",
            overflow: "scroll",
            width: "100%",
        }}>
            <List resource={"test-type"} exporter={false} aside={<TestFilterSidebar setSelectedData={setSelectedData} selectedData={selectedData} />}
                perPage={999999}
                storeKey={false}
                actions={false}
                title={false}
                pagination={false}
                disableSyncWithLocation
                sx={{
                    width: "100%",
                }}
            >
                <TestTable selectedData={selectedData} setSelectedData={setSelectedData} {...props} />
            </List>
        </Box>
    );
}


type PatientFormProps = {
    open: boolean
    onClose: () => void
    setPatientID: React.Dispatch<React.SetStateAction<number | undefined>>
}

function PatientFormModal(props: PatientFormProps) {
    const [create, { isPending }] = useCreate("patient");

    const PatientToolbar = () => {

        return (
            <Toolbar>
                <SaveButton
                    label="Save Patient"
                    resource="patient"
                    disabled={isPending}
                />

            </Toolbar>
        );
    };


    const notify = useNotify();
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
                Create Patient
            </DialogTitle>
            <DialogContent sx={{}}>
                <SimpleForm resource="patient" sx={{
                    width: "100%",
                }} toolbar={<PatientToolbar />} onSubmit={async (data: any) => {
                    console.log(data);
                    create("patient", {
                        data: data,
                    }, {
                        onSuccess: (data) => {
                            notify("Success create patient", {
                                type: 'success',
                            });
                            props.setPatientID(data.id);
                            props.onClose();
                        },
                        onError: () => {
                            notify("Error create patient", {
                                type: 'error',
                            });
                        }
                    }
                    );
                }} >
                    <PatientFormField mode="CREATE" />
                </SimpleForm>
            </DialogContent>
        </Dialog >
    )
}


function NoPatient(props: CreatePatientButtonProps) {
    return (
        <Stack sx={{ width: "100%" }} spacing={2}>
            <Typography fontSize={16}>No Patient found </Typography>
            <CreatePatientButton setOpen={props.setOpen} />
        </Stack>
    )
}


function PatientInput(props: WorkOrderFormProps) {
    const [open, setOpen] = useState(false);
    const [patientID, setPatientID] = useState<number | undefined>(undefined);
    const { setValue, watch, getValues } = useFormContext()
    const notify = useNotify();

    useEffect(() => {
        if (patientID) {
            setValue("patient_id", patientID);
        }
    }, [patientID])
    const axios = useAxios();

    useEffect(() => {
        const patientID = getValues("patient_id");
        if (!patientID) {
            return;
        }

        axios.get(`patient/${patientID}`).then((res) => {
            console.log(res)
            for (const [key, value] of Object.entries(res.data)) {
                setValue(key, value);
            }
        }).catch((err) => {
            notify("Error get patient info", {
                type: 'error',
            });
        })
    }, [watch("patient_id")])

    return (
        <>
            <Stack sx={{
                marginBottom: "2rem",
            }}>
                <ReferenceInput source="patient_id" reference="patient" target="patient_id" label="Patient Name">
                    <AutocompleteInput
                        shouldRenderSuggestions={(val: string) => { return val.trim().length > 2 }}
                        suggestionLimit={10}
                        noOptionsText={<NoPatient setOpen={setOpen} />}
                    />
                </ReferenceInput>
                {
                    watch("patient_id") && <>
                        <Divider sx={{ my: "1rem" }} />
                        <Stack>
                            <Typography variant="subtitle1">Patient Info</Typography>
                            <Stack direction={"row"} gap={5} width={"100%"}>
                                <TextInput source="patient_id" label="Patient ID" readOnly />
                                <TextInput source="first_name" readOnly />
                                <TextInput source="last_name" readOnly />
                            </Stack>
                            <Stack direction={"row"} gap={3} width={"100%"}>
                                <DateInput source="birthdate" label="Birth Date" readOnly />
                                <TextInput source="sex" readOnly />
                            </Stack>
                        </Stack>
                    </>
                }
            </Stack>
            <PatientFormModal open={open} onClose={() => setOpen(false)} setPatientID={setPatientID} />
        </>
    )
}

type CreatePatientButtonProps = {
    setOpen: React.Dispatch<React.SetStateAction<boolean>>
} & Partial<ButtonProps>

function CreatePatientButton(props: CreatePatientButtonProps) {
    return <Button variant="contained" color="secondary" sx={{
        maxWidth: "200px",
    }} endIcon={<AddIcon />} onClick={() => props.setOpen(true)} {...props}>
        Create Patient
    </Button>;
}

export const WorkOrderSaveButton = ({ disabled }: { disabled?: boolean }) => {
    const { getValues } = useFormContext();
    const { save } = useSaveContext();
    const notify = useNotify();
    const handleClick = (e: any) => {
        e.preventDefault(); // necessary to prevent default SaveButton submit logic
        const { ...data } = getValues();

        if (data == undefined) {
            notify("Please fill in all required fields", {
                type: "error",
            });
            return;
        }

        if (!data[patientIDField] || data[patientIDField].length === 0) {
            notify("Please select patient", {
                type: "error",
            });
            return;
        }

        if (!data[testTypesField] || data[testTypesField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[testTypesField] as TestType[]
            save({
                ...data,
                observation_requests: observationRequest.map((test: TestType) => {
                    return test.code
                })
            });
        }
    };


    return <SaveButton type="button" onClick={handleClick} alwaysEnable size="small" />
}


const steps = ['Patient', 'Test'];


export default function WorkOrderForm(props: WorkOrderFormProps) {
    const [activeStep, setActiveStep] = React.useState(0);
    const { save } = useSaveContext();
    const notify = useNotify();
    const onFinish = (data: any) => {
        if (data == undefined) {
            notify("Please fill in all required fields", {
                type: "error",
            });
            return;
        }

        if (!data[patientIDField]) {
            notify("Please select patient", {
                type: "error",
            });
            return;
        }

        if (!data[testTypesField] || data[testTypesField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[testTypesField] as Record<number, ObservationRequestCreateRequest>
            save({
                patient_id: data[patientIDField],
                test_types: Object.entries(observationRequest).map(([_, value]) => {
                    return value
                })
            });
        }
    };

    return (
        <Form>
            <Box sx={{
                margin: '24px',
            }}>
                <FormStepper activeStep={activeStep} setActiveStep={setActiveStep} steps={steps} onFinish={onFinish} >
                    {
                        activeStep === 0 && <PatientInput {...props} />
                    }
                    {
                        activeStep === 1 && <TestInput />
                    }
                </FormStepper>
            </Box>
        </Form>
    )
}
