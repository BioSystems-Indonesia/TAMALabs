import {
    Button,
    ChipField,
    Create,
    Datagrid,
    DateField,
    DeleteButton,
    Edit,
    List,
    ShowButton,
    TextField,
    TopToolbar,
    useListContext,
    useNotify,
    useRecordContext,
    useRefresh,
    useUnselectAll
} from "react-admin";
import {useSearchParams} from "react-router-dom";
import WorkOrderForm, { WorkOrderSaveButton } from "./Form.tsx";
import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import { useMutation } from "@tanstack/react-query";
import { Stack } from "@mui/material";
import { useParams } from "react-router-dom";

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
            }} resource={`work-order/${id}/add-test`}>
            <WorkOrderForm mode={"ADD_TEST"} />
        </Create>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit mutationMode={"pessimistic"} actions={<WorkOrderAction />} sx={{
            "& .RaEdit-card": {
                overflow: "visible",
            }
        }}>
            <WorkOrderForm mode={"EDIT"} />
        </Edit>
    )
}

export const WorkOrderList = () => (
    <List>
        <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <ChipField source="status" />
            <DateField source="created_at" />
            <DateField source="updated_at" />
            <DeleteButton />
        </Datagrid>
    </List>
);
