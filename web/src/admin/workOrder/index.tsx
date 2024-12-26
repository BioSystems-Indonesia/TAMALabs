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
    useRefresh,
    useUnselectAll
} from "react-admin";
import WorkOrderForm, { WorkOrderSaveButton } from "./Form.tsx";
import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import { useMutation } from "@tanstack/react-query";
//import React from "react";
//import {Simulate} from "react-dom/test-utils";

const WorkOrderAction = () => {
    return (
        <TopToolbar>
            <ShowButton />
        </TopToolbar>
    )
}

export function WorkOrderCreate() {
    return (
        <Create redirect={"show"} actions={<WorkOrderAction />}>
            <WorkOrderForm mode={"CREATE"} />
        </Create>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit mutationMode={"pessimistic"} actions={<WorkOrderAction />}>
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


const WorkOrderBulkActionButtons = () => (
    <>
        <RunWorkOrderButton />
    </>
);

export const WorkOrderList = () => (
    <List>
        <Datagrid bulkActionButtons={<WorkOrderBulkActionButtons />}>
            <TextField source="id" />
            <ChipField source="status" />
            <DateField source="created_at" />
            <DateField source="updated_at" />
            <DeleteButton />
        </Datagrid>
    </List>
);
