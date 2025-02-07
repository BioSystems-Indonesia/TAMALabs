import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useEffect, useState } from "react";
import {
    CreateButton,
    Datagrid,
    DateField,
    DateTimeInput,
    DeleteButton,
    FilterListSection,
    FilterLiveForm,
    FilterLiveSearch,
    List,
    RadioButtonGroupInput,
    SaveButton,
    TabbedForm,
    TextField,
    TextInput,
    Toolbar,
    TopToolbar,
    useGetMany,
    useGetOne,
    useListContext,
    useNotify,
    useSaveContext
} from "react-admin";
import { FieldValues, UseFormWatch, useFormContext } from "react-hook-form";
import { useParams, useSearchParams } from "react-router-dom";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import FeatureList from "../../component/FeatureList.tsx";
import { getRefererParam } from "../../hooks/useReferer.ts";
import type { ObservationRequest } from "../../types/observation_requests.ts";
import { ActionKeys } from "../../types/props.ts";
import type { TestType } from "../../types/test_type.ts";
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

const patientIDsField = "patient_ids";


function TestTable(props: TestInputProps) {
    const { selectedIds, onSelect, data: testList } = useListContext();
    const { setValue } = useFormContext();

    const { id } = useParams()
    const { data, isLoading } = useGetOne('work-order', { id: id });
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
        if (data && searchParams.getAll("patient_id").length > 0) {
            if (!allTestType) {
                console.error("testList is undefined");
                return
            }

            const patientIDs = searchParams.getAll("patient_id")!.map(id => parseInt(id));
            const patients = data.patient_list.filter((patient: any) => {
                return patientIDs.includes(patient.id);
            })
            const observationRequestCodeList = patients.map((patient: any) => {
                return patient.specimen_list.map((specimen: any) => {
                    return specimen.observation_requests.map((observationRequest: ObservationRequest) => {
                        return observationRequest.test_code;
                    })
                }).flat()
            }) as string[][]
            if (observationRequestCodeList.length === 0) {
                return;
            }

            const observationRequestCodesLongest = observationRequestCodeList.reduce((acc, cur) => {
                return acc.length > cur.length ? acc : cur;
            }, observationRequestCodeList[0]);
            console.debug("observationRequestCodesLongest", observationRequestCodesLongest);
            setValue(observationRequestField, observationRequestCodesLongest);


            const observationRequestIDs = allTestType.filter((test: TestType) => {
                return observationRequestCodesLongest.includes(test.code);
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
                isLoading={isLoading}
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
    return (<List resource={"test-type"} exporter={false} aside={<TestFilterSidebar />}
        perPage={999999}
        storeKey={false}
        actions={false}
        title={false}
        pagination={false}
        disableSyncWithLocation
        sx={{
            marginTop: "48px",
            width: "100%"
        }}
    >
        <TestTable {...props} />
    </List>);
}

const PatientFilterSidebar = () => (
    <Card sx={{
        order: -1, mr: 2, mt: 2, width: 200, minWidth: 200,

    }}>
        <CardContent>
            <FilterLiveSearch onKeyDown={(e) => { e.key === 'Enter' && e.preventDefault() }} />
            <FilterListSection label="Birth Date" icon={<CalendarMonthIcon />}>
                <FilterLiveForm debounce={1500}>
                    <CustomDateInput source={"birthdate"} label={"Birth Date"} clearable />
                </FilterLiveForm>
            </FilterListSection>
        </CardContent>
    </Card>
);

function PickedPatient() {
    const { selectedIds, onToggleItem } = useListContext();

    if (selectedIds.length === 0) {
        return (
            <Typography fontSize={16}>Please select patient</Typography>
        )
    }

    const { data, isPending, error } = useGetMany("patient", {
        ids: selectedIds,
    });

    if (isPending) {
        return <Typography fontSize={16}>Loading...</Typography>
    }
    if (error) {
        return <Typography fontSize={16} color="error">{error.message}</Typography>
    }

    return (
        <Stack spacing={2}>
            <Typography fontSize={16}>Selected patient</Typography>
            <Grid container spacing={1}>
                {
                    data?.map((v: any) => {
                        return (
                            <Grid item key={v.id}>
                                <Chip label={`${v.id} - ${v.first_name} ${v.last_name}`}
                                    onDelete={() => {
                                        const currentId = v.id;
                                        console.log(currentId)
                                        onToggleItem(currentId);
                                    }}
                                />
                            </Grid>
                        )
                    })
                }
            </Grid>
        </Stack>
    )
}


function PatientTable(props: WorkOrderFormProps) {
    const BulkActionButtons = () => {
        return (
            <></>
        );
    };
    const { selectedIds, onSelect } = useListContext();
    const { setValue } = useFormContext();
    const [searchParams] = useSearchParams();
    const [havePatientIDsInQueryParam, setHavePatientIDsInQueryParam] = useState(false);

    useEffect(() => {
        if (searchParams.get("patient_id")) {
            const patientIDs = searchParams.getAll("patient_id")!.map(id => parseInt(id));
            console.debug("patientIDs", patientIDs);
            onSelect(patientIDs);
            setValue(patientIDsField, patientIDs);
            setHavePatientIDsInQueryParam(true);
        }
    }, [searchParams]);


    useEffect(() => {
        console.debug("test selected ids", selectedIds);
        setValue(patientIDsField, selectedIds);
    }, [selectedIds]);


    return <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
            <Datagrid rowClick={"toggleSelection"} bulkActionButtons={<BulkActionButtons />}

                // Disable selection when have patient IDs in query param
                isRowSelectable={(record: any) => {
                    if (havePatientIDsInQueryParam) {
                        const patientIDs = searchParams.getAll("patient_id")!.map(id => parseInt(id));
                        return patientIDs.includes(record.id);
                    }

                    return true
                }}
                onToggleItem={!havePatientIDsInQueryParam ? undefined : () => { }}
                onSelect={!havePatientIDsInQueryParam ? undefined : () => { }}
            >
                <TextField source="id" />
                <TextField source="first_name" />
                <TextField source="last_name" />
                <DateField source="birthdate" locales={["id-ID"]} />
                <TextField source="sex" />
                <DateField source="created_at" showTime />
            </Datagrid>
        </Grid>
        <Grid item xs={12} md={4}>
            <PickedPatient />
        </Grid>
    </Grid>;
}

const PatientListActions = () => {
    return (
        <TopToolbar>
            <CreateButton to={`/patient/create?${getRefererParam()}`} label="Create Patient" />
        </TopToolbar>
    )
};

function PatientInput(props: WorkOrderFormProps) {

    return (
        <List aside={<PatientFilterSidebar />} resource={"patient"}
            actions={<PatientListActions />}
            exporter={false}
            title={false}
            perPage={25}
            sx={{
                width: "100%"
            }}
            disableSyncWithLocation
            storeKey={false}
            empty={false}
        >
            <PatientTable {...props} />
        </List>
    )
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

        if (!data[patientIDsField] || data[patientIDsField].length === 0) {
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

const validForm = (watch: UseFormWatch<FieldValues>): boolean => {
    return (!!watch(observationRequestField) && watch(observationRequestField).length !== 0)
        && (!!watch(patientIDsField) && watch(patientIDsField).length !== 0)
}

const WorkOrderToolbar = () => {
    const { watch } = useFormContext();
    return (
        <Stack width={"100%"}
            sx={{
                position: "sticky",
                top: 48,
                borderBottom: "1px solid #ccc",
                zIndex: 2147483647,
                marginBottom: 1,
            }}
        >
            <Toolbar sx={{
                gap: 2,
                width: "100%",
                display: "flex",
                justifyContent: "flex-end",
            }}>
                <DeleteButton variant="contained" size="small" />
                {validForm(watch) && <WorkOrderSaveButton disabled={!validForm(watch)} />}
            </Toolbar>
        </Stack>
    )
};

const showDetailOnMode: Array<WorkOrderActionKeys> = ["SHOW", "EDIT"];
export default function WorkOrderForm(props: WorkOrderFormProps) {
    return (
        <TabbedForm toolbar={false} >
            <TabbedForm.Tab label="Patient">
                <WorkOrderToolbar />
                <PatientInput {...props} />
            </TabbedForm.Tab>
            <TabbedForm.Tab label="Test" sx={{
                position: "relative",
                overflow: "visible",
            }}>
                <WorkOrderToolbar />
                <TestInput />
            </TabbedForm.Tab>
            {showDetailOnMode.includes(props.mode) && (
                <TabbedForm.Tab label="Detail">
                    <WorkOrderToolbar />
                    <div>
                        <TextInput source={"id"} readOnly={true} size={"small"} />
                        <DateTimeInput source={"created_at"} readOnly={true} size={"small"} />
                        <DateTimeInput source={"updated_at"} readOnly={true} size={"small"} />
                        <FeatureList types={"work-order-status"} source={"status"}>
                            <RadioButtonGroupInput source="status" readOnly={true} size={"small"} />
                        </FeatureList>
                        <Divider />
                    </div>
                </TabbedForm.Tab>
            )}
        </TabbedForm>
    )
}
