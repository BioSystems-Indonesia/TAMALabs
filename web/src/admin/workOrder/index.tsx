import PlayCircleFilledIcon from "@mui/icons-material/PlayCircleFilled";
import { CircularProgress, Dialog, DialogContent, DialogTitle } from "@mui/material";
import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useMutation } from "@tanstack/react-query";
import dayjs from "dayjs";
import { useState } from "react";
import {
    Button,
    Create,
    Datagrid,
    DateField,
    DeleteButton,
    Edit,
    List,
    ReferenceField,
    SearchInput,
    ShowButton,
    TextField,
    TopToolbar,
    WithRecord,
    WrapperField,
    useListContext,
    useNotify,
    useRefresh
} from "react-admin";
import { useParams, useSearchParams } from "react-router-dom";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import useAxios from "../../hooks/useAxios.ts";
import type { WorkOrder } from "../../types/work_order.ts";
import { WorkOrderChipColorMap } from "./ChipFieldStatus.tsx";
import WorkOrderForm from "./Form.tsx";
import { RunWorkOrderForm } from "./Show.tsx";



const WorkOrderAction = () => {
    return (
        <TopToolbar>
            <ShowButton />
        </TopToolbar>
    )
}

export function WorkOrderCreate() {
    return (
        <Create redirect={"show"} actions={<WorkOrderAction />} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }}>
            <WorkOrderForm mode={"CREATE"} />
        </Create>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit redirect={"show"} actions={<WorkOrderAction />} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }} mutationMode="pessimistic">
            <WorkOrderForm mode={"EDIT"} />
        </Edit>
    )
}

export function WorkOrderAddTest() {
    const { id } = useParams();
    const [searchParams] = useSearchParams();

    const getTitle = () => {
        const patientIDs = searchParams.getAll("patient_id")!.map(id => parseInt(id));
        if (patientIDs.length === 0) {
            return `Add Test to Work Order #${id}`
        }

        return `Edit Test Work Order #${id} for Patient ID: ${searchParams.getAll("patient_id").join(", ")}`
    }

    return (
        <Create
            title={getTitle()}
            redirect={() => {
                return `work-order/${id}/show`
            }} actions={<WorkOrderAction />} sx={{
                "& .RaCreate-card": {
                    overflow: "visible",
                }
            }} resource={`work-order/${id}/show/add-test`}>
            <WorkOrderForm mode={"ADD_TEST"} />
        </Create>
    )
}


const WorkOrderFilters = [
    <SearchInput source="q" alwaysOn />,
    <CustomDateInput label={"Created At Start"} source="created_at_start" alwaysOn />,
    <CustomDateInput label={"Created At End"} source="created_at_end" alwaysOn />,
];

function getRequestLength(data: WorkOrder): number {
    return data.specimen_list?.reduce((acc, specimen) => acc + specimen.observation_requests.length, 0) || 0
}

function RunWorkOrderButton(props: RunWorkOrderProps) {
    const notify = useNotify();
    const refresh = useRefresh();
    const axios = useAxios();
    const { isPending } = useMutation({
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

    return (
        <Button label="Run Work Order" onClick={() => {
            props.setOpen(true)
        }}>
            {isPending ? <CircularProgress size={12} variant='indeterminate' color='primary' /> : <PlayCircleFilledIcon />}
        </Button>
    )
}

type RunWorkOrderProps = {
    open: boolean
    setOpen: React.Dispatch<React.SetStateAction<boolean>>
    onClose: () => void
}

function RunWorkOrderDialog(props: RunWorkOrderProps) {
    const { selectedIds } = useListContext();

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
                Run Work Order
            </DialogTitle>
            <DialogContent sx={{}}>
                <RunWorkOrderForm workOrderIDs={selectedIds} />
            </DialogContent>
        </Dialog >
    )
}


const WorkOrderListBulkActionButtons = (props: RunWorkOrderProps) => (
    <>
        <RunWorkOrderButton {...props} />
    </>
)

export const WorkOrderList = () => {
    const [open, setOpen] = useState(false)

    return (
        <List sort={{
            field: "id",
            order: "DESC"
        }} filters={WorkOrderFilters} title="Lab Request" filterDefaultValues={{
            created_at_start: dayjs().subtract(7, "day").toISOString(),
            created_at_end: dayjs().toISOString(),
        }} storeKey={false} exporter={false}>
            <Datagrid bulkActionButtons={<WorkOrderListBulkActionButtons
                open={open}
                setOpen={setOpen}
                onClose={() => setOpen(false)}
            />}>
                <TextField source="id" />
                <WithRecord label="Status" render={(record: any) => (
                    <Chip label={`${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                )} />
                <ReferenceField source="patient_id" reference="patient">
                </ReferenceField>
                <WithRecord label="Request" render={(record: any) => (
                    <Typography variant="body2" >
                        {getRequestLength(record)}
                    </Typography>
                )} />
                <DateField source="created_at" />
                <DateField source="updated_at" />
                <WrapperField label="Actions" sortable={false} >
                    <Stack direction={"row"} spacing={2}>
                        <ShowButton variant="contained" />
                        <DeleteButton variant="contained" />
                    </Stack>
                </WrapperField>
            </Datagrid>
            <RunWorkOrderDialog open={open} onClose={() => setOpen(false)} setOpen={setOpen} />
        </List>
    )
};
