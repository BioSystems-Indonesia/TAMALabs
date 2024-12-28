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

const RunWorkOrderButton = () => {
    const { selectedIds } = useListContext();
    const refresh = useRefresh();
    const notify = useNotify();
    const unselectAll = useUnselectAll('posts');

    const { mutate, isPending } = useMutation({
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
            unselectAll();
        },
        onError: (error) => {
            notify('Error:' + error.message, {
                type: 'error',
            });
            refresh();
        },
    })

    const handleClick = () => {
        mutate({
            work_order_ids: selectedIds
        });
    }

    return (
        <Button label="Run Work Order" onClick={handleClick} disabled={isPending}>
            <PlayCircleFilledIcon />
        </Button>
    );
}
    ;


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
