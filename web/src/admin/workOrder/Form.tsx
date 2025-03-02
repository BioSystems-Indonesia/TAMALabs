import AddIcon from '@mui/icons-material/Add';
import { Box, Button, Dialog, DialogContent, DialogTitle, type ButtonProps } from "@mui/material";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import React, { useEffect, useState } from "react";
import {
    AutocompleteInput,
    Datagrid,
    DateInput,
    Form,
    List,
    ReferenceInput,
    SaveButton,
    SimpleForm,
    TextField,
    TextInput,
    Toolbar,
    useCreate,
    useListContext,
    useNotify,
    useRecordContext,
    useSaveContext
} from "react-admin";
import { useFormContext } from "react-hook-form";
import { useSearchParams } from "react-router-dom";
import useAxios from "../../hooks/useAxios.ts";
import type { ObservationRequest } from "../../types/observation_requests.ts";
import { ActionKeys } from "../../types/props.ts";
import type { TestType } from "../../types/test_type.ts";
import { PatientFormField } from "../patient/index.tsx";
import FormStepper from "./Stepper.tsx";
import { TestFilterSidebar } from "./TestTypeFilter.tsx";


type WorkOrderActionKeys = ActionKeys | "ADD_TEST";

type WorkOrderFormProps = {
    readonly?: boolean
    mode: WorkOrderActionKeys
}


const observationRequestField = "observation_requests";


const PickedTest = ({ allTestType }: { allTestType: TestType[] }) => {
    const { watch } = useFormContext()
    const [selectedData, setSelectedData] = useState<any[]>([]);

    useEffect(() => {
        if (!allTestType) {
            return;
        }

        const testTypes = watch(observationRequestField) as TestType[] | undefined
        if (!testTypes) {
            return;
        }

        setSelectedData(testTypes);
    }, [watch(observationRequestField), allTestType]);

    if (watch(observationRequestField)?.length === 0) {
        return (
            <Typography fontSize={16}>Please select test to run</Typography>
        )
    }

    return (
        <Stack spacing={2}>
            <Typography fontSize={16}>Selected Test</Typography>
            <Grid container spacing={1}>
                {
                    selectedData.map((v: any) => {
                        return (
                            <Grid item key={v.id}>
                                <Chip label={v.code} />
                            </Grid>
                        )
                    })
                }
            </Grid>
        </Stack>
    )
}

const patientIDField = "patient_id";


function TestTable(props: TestInputProps) {
    const { selectedIds, onSelect, data: testList } = useListContext();
    const { setValue } = useFormContext();

    const data = useRecordContext();
    const [searchParams] = useSearchParams();
    const [allTestType, setAllTestType] = useState<TestType[]>([]);
    const [allTestTypeSet, setAllTestTypeSet] = useState<boolean>(false);

    useEffect(() => {
        if (props.initSelectedIds) {
            onSelect(props.initSelectedIds);
        }
    }, [props.initSelectedIds])

    // This will store all test type (without filter)
    useEffect(() => {
        if ((testList && testList.length > 0) && !allTestTypeSet) {
            setAllTestTypeSet(true);
            setAllTestType(testList);
        }
    }, [testList])

    useEffect(() => {
        if (data) {
            if (!allTestType) {
                console.error("testList is undefined");
                return
            }

            const observationRequestCodeList = data.patient.specimen_list.map((specimen: any) => {
                return specimen.observation_requests.map((observationRequest: ObservationRequest) => {
                    return observationRequest.test_code;
                })
            }).flat()
            if (observationRequestCodeList.length === 0) {
                return;
            }

            console.debug("observationRequestCodesLongest", observationRequestCodeList);
            setValue(observationRequestField, observationRequestCodeList);


            const observationRequestIDs = allTestType.filter((test: TestType) => {
                return observationRequestCodeList.includes(test.code);
            }).map((test: any) => {
                return test?.id
            })
            console.debug("observationRequestIDs", observationRequestIDs);
            onSelect(observationRequestIDs);
        }
    }, [data, searchParams, allTestType]);

    useEffect(() => {
        console.debug("selected ids", selectedIds);
        const pickedTestType = allTestType?.filter((test: any) => {
            return selectedIds.includes(test?.id);
        })
        console.debug("observationRequest", pickedTestType);
        setValue(observationRequestField, pickedTestType);
    }, [selectedIds, allTestType]);

    const BulkActionButtons = () => {
        return (
            <></>
        );
    };

    return <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
            <Datagrid width={"100%"}
                bulkActionButtons={<BulkActionButtons />}
                rowClick={"toggleSelection"}
            >
                <TextField label={"Name"} source={"name"} />
                <TextField label={"Code"} source={"code"} />
                <TextField label={"Category"} source={"category"} />
            </Datagrid>
        </Grid>
        <Grid item xs={12} md={4}>
            <PickedTest allTestType={allTestType} />
        </Grid>
    </Grid>;
}

type TestInputProps = {
    initSelectedIds?: number[]
}

export function TestInput(props: TestInputProps) {
    return (
        <Box sx={{
            maxHeight: "calc(70vh - 48px)",
            overflow: "scroll",
        }}>
            <List resource={"test-type"} exporter={false} aside={<TestFilterSidebar />}
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
                <TestTable {...props} />
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

        if (!data[observationRequestField] || data[observationRequestField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[observationRequestField] as TestType[]
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

        if (!data[observationRequestField] || data[observationRequestField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[observationRequestField] as TestType[]
            save({
                patient_id: data[patientIDField],
                test_ids: observationRequest.map((test: TestType) => {
                    return test.id
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
