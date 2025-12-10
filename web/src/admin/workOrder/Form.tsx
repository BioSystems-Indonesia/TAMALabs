import AddIcon from '@mui/icons-material/Add';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import ListAltIcon from '@mui/icons-material/ListAlt';
import TouchAppIcon from '@mui/icons-material/TouchApp';
import { Box, Button, Checkbox, Dialog, DialogContent, DialogTitle, GridLegacy as Grid, ListItem, ListItemText, MenuItem, Paper, Select, Tooltip, type ButtonProps } from "@mui/material";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import MUIList from "@mui/material/List";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { GridRowId, DataGrid as MuiDatagrid, useGridApiRef, type GridRenderCellParams } from "@mui/x-data-grid";
import React, { useEffect, useState } from "react";
import {
    AutocompleteArrayInput,
    AutocompleteInput,
    DateInput,
    DateTimeInput,
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
import { useCurrentUser } from '../../hooks/currentUser.ts';
import useAxios from "../../hooks/useAxios.ts";
import type { ObservationRequest, ObservationRequestCreateRequest } from "../../types/observation_requests.ts";
import { ActionKeys } from "../../types/props.ts";
import { RoleNameValue } from '../../types/role.ts';
import type { Specimen } from '../../types/specimen.ts';
import type { TestType, TestTypeSpecimenType } from "../../types/test_type.ts";
import { WorkOrder } from '../../types/work_order.ts';
import { PatientFormField } from "../patient/index.tsx";
import FormStepper from "./Stepper.tsx";
import { TestFilterSidebar } from "./TestTypeFilter.tsx";


type WorkOrderActionKeys = ActionKeys | "ADD_TEST";

type WorkOrderFormProps = {
    readonly?: boolean
    mode: WorkOrderActionKeys
}

type InputProps = {
    setDisableNext: React.Dispatch<React.SetStateAction<boolean>>
} & WorkOrderFormProps


function TestSelectorPrompt() {
    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                border: '2px',
                borderRadius: 2,
                textAlign: 'center',
                maxWidth: '250px',
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
        <Paper elevation={2} sx={{ p: 2, backgroundColor: 'background.paper', maxWidth: "250px" }}>
            <Stack spacing={2}>
                {/* Enhanced Title Section */}
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <ListAltIcon color="primary" sx={{ mr: 1 }} /> {/* Icon before title */}
                    <Typography variant="h6" fontSize={18} fontWeight="medium"> {/* More prominent title */}
                        Selected Tests
                    </Typography>
                </Box>
                <Divider />
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
const analystIDField = "analyzer_ids";
const doctorIDField = "doctor_ids";
const testTemplateIDField = "test_template_ids";

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
        if (Object.keys(selectedData).length > 0) {
            setValue(testTypesField, selectedData);
        }

        if (testType && Object.keys(selectedData).length > 0) {
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
    const lastClickedRowIdRef = React.useRef<GridRowId | null>(null); // To track for shift-click
    // Main handler for checkbox changes (click, shift+click, keyboard)
    const handleCheckboxInteraction = (rowId: GridRowId, newCheckedState: boolean, isShiftKey: boolean) => {
        const newRows = [...rows];

        if (isShiftKey && lastClickedRowIdRef.current !== null && lastClickedRowIdRef.current !== rowId) {
            const lastIdx = newRows.findIndex(r => r.id === lastClickedRowIdRef.current);
            const currentIdx = newRows.findIndex(r => r.id === rowId);

            if (lastIdx !== -1 && currentIdx !== -1) {
                const start = Math.min(lastIdx, currentIdx);
                const end = Math.max(lastIdx, currentIdx);

                // Get the checked state of the row that initiated the shift-click range.
                // The desired behavior might be to set all in range to `newCheckedState` (state of current click)
                // or to the state of `newRows[lastIdx].checked`. For simplicity, using `newCheckedState`.
                const targetCheckedStateForRowRange = newRows[lastIdx].checked;


                for (let i = start; i <= end; i++) {
                    newRows[i] = { ...newRows[i], checked: targetCheckedStateForRowRange };
                    updateSelectedData(newRows[i])
                }
                // Also ensure the currently clicked row (the end of the range) gets the intended newCheckedState
                // const clickedRowInRange = newRows.find(r => r.id === rowId);
                // if (clickedRowInRange) {
                //     newRows = newRows.map(r => r.id === rowId ? { ...r, checked: newCheckedState } : r);
                //     updateSelectedData(clickedRowInRange)
                // }

            } else { // Fallback to single toggle if indices are invalid
                const targetRowIndex = newRows.findIndex(r => r.id === rowId);
                if (targetRowIndex !== -1) {
                    newRows[targetRowIndex] = { ...newRows[targetRowIndex], checked: newCheckedState };
                    updateSelectedData(newRows[targetRowIndex])
                }
            }
        } else {
            // Single toggle
            const targetRowIndex = newRows.findIndex(r => r.id === rowId);
            if (targetRowIndex !== -1) {
                newRows[targetRowIndex] = { ...newRows[targetRowIndex], checked: newCheckedState };
                updateSelectedData(newRows[targetRowIndex])
            }
        }

        setRows(newRows);

        // Update lastClickedRowIdRef only on non-shift clicks or if it's the first click in a potential shift-sequence
        if (!isShiftKey) {
            lastClickedRowIdRef.current = rowId;
        }
    };

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
                        field: 'customChecked',
                        headerName: 'Checked',
                        flex: 0.5,
                        filterable: false,
                        sortable: false,

                        renderHeader: () => {
                            const allRowsChecked = rows.length > 0 && rows.every(row => row.checked);
                            const someRowsChecked = rows.some(row => row.checked);
                            const isIndeterminate = someRowsChecked && !allRowsChecked;

                            return (
                                <Checkbox
                                    checked={allRowsChecked}
                                    indeterminate={isIndeterminate}
                                    onChange={(event, checked) => {
                                        setRows(prevRows => prevRows.map(row => {
                                            const newRow = { ...row, checked }
                                            updateSelectedData(newRow)
                                            return newRow
                                        }));
                                        lastClickedRowIdRef.current = null; // Reset shift-click anchor
                                    }}
                                />
                            );
                        },
                        renderCell: (params: GridRenderCellParams<any, TestType>) => {
                            const handleKeyDown = (event: React.KeyboardEvent) => {
                                if (event.key === ' ' || event.key === 'Enter') {
                                    event.preventDefault();  // Prevent default browser action (e.g., scrolling on space)
                                    event.stopPropagation(); // Prevent event from bubbling to DataGrid, which might move focus
                                    handleCheckboxInteraction(params.row.id, !params.row.checked, false);
                                }
                            };

                            return (
                                // Wrapper to make the cell focusable and handle keyboard events
                                <Box
                                    tabIndex={0} // Make it focusable
                                    onKeyDown={handleKeyDown}
                                    sx={{
                                        width: '100%',
                                        height: '100%',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        outline: 'none', // Remove default focus outline if desired, or style it
                                        '&:focus-visible': { // Modern way to style keyboard focus
                                            boxShadow: `0 0 0 2px rgba(0,123,255,.5)`, // Example focus ring
                                        }
                                    }}
                                    onClick={(e) => {
                                        // This outer click can be used if you want the whole cell to be clickable
                                        // For now, relying on Checkbox's own click
                                    }}
                                >
                                    <Checkbox
                                        checked={!!params.row.checked} // Ensure it's a boolean
                                        onChange={(event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => {
                                            const isShift = (event.nativeEvent instanceof MouseEvent && (event.nativeEvent as MouseEvent).shiftKey);
                                            handleCheckboxInteraction(params.row.id, checked, isShift);
                                        }}
                                        // Prevent clicks on the checkbox from propagating to the row's onRowClick (if any)
                                        onClick={(e) => e.stopPropagation()}
                                        inputProps={{ 'aria-label': `Select row ${params.row.name || params.row.id}` }}

                                    />
                                </Box>
                            );
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
                        field: 'is_calculated_test',
                        headerName: 'Calc Test',
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
    const { getValues, setValue } = useFormContext();
    useEffect(() => {
        const testFields = getValues(testTypesField);
        if (testFields && Object.keys(testFields).length > 0) {
            setSelectedData(testFields);
        }
    }, [])
    return (
        <Box sx={{
            maxHeight: "calc(60vh - 48px)",
            overflow: "scroll",
            width: "100%",
            paddingTop: 10
        }}>
            <List resource={"test-type"} exporter={false}
                aside={
                    <TestFilterSidebar setSelectedData={setSelectedData}
                        selectedData={selectedData} setValue={setValue} getValues={getValues} />
                }
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


// function NoDoctor(props: CreatePatientButtonProps) {
//     return (
//         <Stack sx={{ width: "100%" }} spacing={2}>
//             <Typography fontSize={16}>No Doctor Found</Typography>
//             <CreatePatientButton setOpen={props.setOpen} />
//         </Stack>
//     )
// }

function PatientInput(props: InputProps) {
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
            props.setDisableNext(true);
            return;
        }

        props.setDisableNext(false);
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
                {/* <Divider sx={{
                    my: "0.5rem",
                }}/>
                <Typography variant="subtitle1" sx={{
                    mb: "0.5rem",
                }}>Required</Typography> */}
                <ReferenceInput source="patient_id" reference="patient" target="patient_id" label="Patient Name">
                    <AutocompleteInput
                        suggestionLimit={10}
                        noOptionsText={<NoPatient setOpen={setOpen} />}
                    />
                </ReferenceInput>
                {
                    watch("patient_id") && <>
                        <Divider sx={{ my: "1rem" }} />
                        <Stack>
                            <Typography variant="subtitle1">Patient Info Preview</Typography>
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
                <Divider sx={{
                    my: "0.5rem",
                }} />
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

function SIMRSInformationInput(props: InputProps) {
    return (
        <>
            <Stack sx={{
                marginBottom: "2rem",
            }}>
                <Typography variant="subtitle1" sx={{
                    mb: "0.5rem",
                }}>SIMRS Informations</Typography>
                <Stack gap={1}>
                    <TextInput source="barcode_simrs" label="No Order Lab" helperText="No Order Lab for SIMRS integration" fullWidth />
                    <TextInput source="medical_record_number" label="Medical Record Number" helperText="Medical record number from SIMRS" fullWidth />
                    <TextInput source="visit_number" label="Visit Number" helperText="Patient visit number" fullWidth />
                    <DateTimeInput
                        source="specimen_collection_date"
                        label="Collection Date & Time"
                        helperText="Specimen collection date and time"
                        fullWidth
                        inputProps={{
                            max: "2100-12-31T23:59",
                            min: "1900-01-01T00:00"
                        }}
                    />
                    <TextInput source="diagnosis" label="Diagnosis / Clinical Notes" helperText="Clinical diagnosis or notes" fullWidth multiline rows={2} />
                </Stack>
            </Stack>
        </>
    )
}

function AdditionalInput(props: InputProps) {
    const currentUser = useCurrentUser();

    return (
        <>
            <Stack sx={{
                marginBottom: "2rem",
            }}>
                <Typography variant="subtitle1" sx={{
                    mb: "0.5rem",
                }}>Optional</Typography>
                <Stack gap={1}>
                    <ReferenceInput source={doctorIDField} reference="user" resource='user' target="id" label="Doctor" filter={{
                        role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
                    }}>
                        <AutocompleteArrayInput
                            suggestionLimit={10}
                            filterToQuery={(searchText) => ({
                                q: searchText,
                                role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
                            })}
                            helperText="Leave blank if do not need verification"
                        />
                    </ReferenceInput>
                    <ReferenceInput source={analystIDField} reference="user" target="id" label="Analyst">
                        <AutocompleteArrayInput
                            label="Analyst"
                            suggestionLimit={10}
                            defaultValue={[currentUser?.id]}
                            helperText="Default to current user"
                        />
                    </ReferenceInput>
                </Stack>
            </Stack>
        </>
    )
}


const steps = ['Patient', 'Parameter Tests', 'SIMRS Information', 'Additional'];


export default function WorkOrderForm(props: WorkOrderFormProps) {
    const [activeStep, setActiveStep] = React.useState(0);
    const { save } = useSaveContext();
    const notify = useNotify();
    const [disableNext, setDisableNext] = React.useState(false);
    const currentUser = useCurrentUser();
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

        if (!data[testTypesField] || Object.entries(data[testTypesField]).length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[testTypesField] as Record<number, ObservationRequestCreateRequest>
            const testTypes = Object.entries(observationRequest).map(([_, value]) => {
                return value
            })

            const payload = {
                patient_id: data[patientIDField],
                test_types: testTypes,
                created_by: currentUser?.id,
                analyzer_ids: data[analystIDField],
                doctor_ids: data[doctorIDField],
                test_template_ids: data[testTemplateIDField],
                barcode: data.barcode,
                barcode_simrs: data.barcode_simrs,
                medical_record_number: data.medical_record_number,
                visit_number: data.visit_number,
                specimen_collection_date: data.specimen_collection_date ? new Date(data.specimen_collection_date).toISOString() : '',
                diagnosis: data.diagnosis,
            };

            console.log('Work Order Payload:', payload);
            save(payload);
        }
    };

    return (
        <Form>
            <Box sx={{
                margin: '24px',
            }}>
                <FormStepper activeStep={activeStep} setActiveStep={setActiveStep} steps={steps} onFinish={onFinish} disableNext={disableNext}>
                    {
                        activeStep === 0 && <PatientInput {...props} setDisableNext={setDisableNext} />
                    }
                    {
                        activeStep === 1 && <TestInput />
                    }
                    {
                        activeStep === 2 && <SIMRSInformationInput {...props} setDisableNext={setDisableNext} />
                    }
                    {
                        activeStep === 3 && <AdditionalInput {...props} setDisableNext={setDisableNext} />
                    }
                </FormStepper>
            </Box>
        </Form>
    )
}
