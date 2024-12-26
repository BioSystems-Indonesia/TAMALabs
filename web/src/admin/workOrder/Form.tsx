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
    SavedQueriesList,
    TabbedForm,
    TextField,
    TextInput,
    Toolbar,
    TopToolbar,
    useGetMany,
    useListContext,
    useNotify,
    useRecordContext,
    useSaveContext
} from "react-admin";
import { useFormContext } from "react-hook-form";
import { Action, ActionKeys } from "../../types/props.ts";
import FeatureList from "../../component/FeatureList.tsx";
import Divider from "@mui/material/Divider";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
import Grid from "@mui/material/Grid";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import { useEffect } from "react";


type WorkOrderFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}


const observationRequestField = "observation_requests";

const TestFilterSidebar = () => (
    <Card sx={{
        order: -1, mr: 1, mt: 2, width: 200, minWidth: 200,
        overflow: "visible",
    }}>
        <CardContent sx={{
            position: "sticky",
            top: 96,
        }}>
            <SavedQueriesList />
            <FilterLiveSearch onSubmit={(event) => event.preventDefault()} />
        </CardContent>
    </Card>
);

const PickedTest = () => {
    const { selectedIds } = useListContext();

    if (selectedIds.length === 0) {
        return (
            <Typography fontSize={16}>Please select observation request</Typography>
        )
    }

    return (
        <Stack spacing={2}>
            <Typography fontSize={16}>Selected observation requests</Typography>
            <Grid container spacing={1}>
                {
                    selectedIds.map((v: any) => {
                        return (
                            <Grid item key={v}>
                                <Chip label={v} />
                            </Grid>
                        )
                    })
                }
            </Grid>
        </Stack>
    )
}

const patientIDsField = "patient_ids";

// eslint-disable-next-line no-unused-vars
function TestTable(props: WorkOrderFormProps) {
    const { selectedIds, onSelect } = useListContext();
    const { setValue } = useFormContext();

    const data = useRecordContext()
    useEffect(() => {
        if (data) {
            console.debug("setDataToValue", data[observationRequestField]);
            onSelect(data[observationRequestField])
            setValue(observationRequestField, data[observationRequestField]);
        }
    }, [data]);

    useEffect(() => {
        console.debug("test selected ids", selectedIds);
        setValue(observationRequestField, selectedIds);
    }, [selectedIds]);

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
                <TextField label={"ID"} source={"id"} />
                <TextField label={"Name"} source={"name"} />
                <TextField label={"Type"} source={"additional_info.type"} />
            </Datagrid>
        </Grid>
        <Grid item xs={12} md={4}>
            <PickedTest />
        </Grid>
    </Grid>;
}

function TestInput(props: WorkOrderFormProps) {

    return (<List resource={"feature-list-observation-type"} exporter={false} aside={<TestFilterSidebar />}
        perPage={999999}
        storeKey={false}
        disableSyncWithLocation
        sx={{
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
            <FilterLiveSearch />
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

// eslint-disable-next-line no-unused-vars
function PatientTable(props: WorkOrderFormProps) {
    const BulkActionButtons = () => {
        return (
            <></>
        );
    };
    const { selectedIds, onSelect } = useListContext();
    const { setValue } = useFormContext();
    const data = useRecordContext()

    useEffect(() => {
        if (data) {
            console.debug("setDataToValue", data[patientIDsField]);
            onSelect(data[patientIDsField])
            setValue(patientIDsField, data[patientIDsField]);
        }
    }, [data]);

    useEffect(() => {
        console.debug("test selected ids", selectedIds);
        setValue(patientIDsField, selectedIds);
    }, [selectedIds]);


    return <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
            <Datagrid rowClick={"toggleSelection"} bulkActionButtons={<BulkActionButtons />}
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

const PatientListActions = () => (
    <TopToolbar>
        <CreateButton target={"_blank"} rel={"noopener"} />
    </TopToolbar>
);

function PatientInput(props: WorkOrderFormProps) {

    return (
        <List aside={<PatientFilterSidebar />} resource={"patient"}
            actions={<PatientListActions />}
            exporter={false}
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

export const WorkOrderSaveButton = () => {
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
            save(data);
        }
    };


    return <SaveButton type="button" onClick={handleClick} alwaysEnable size="small" />
}

const WorkOrderToolbar = () => {
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
                <WorkOrderSaveButton />
            </Toolbar>
        </Stack>
    )
};

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
                <TestInput {...props} />
            </TabbedForm.Tab>
            {props.mode !== Action.CREATE && (
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
