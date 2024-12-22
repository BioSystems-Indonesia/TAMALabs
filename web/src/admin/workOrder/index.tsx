import {
    Button,
    ChipField,
    Create,
    Datagrid,
    DateField,
    DeleteButton,
    Edit,
    List,
    Show,
    TextField,
    useListContext,
    useNotify,
    useRefresh,
    useUnselectAll
} from "react-admin";
import WorkOrderForm from "./Form.tsx";
import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import {useMutation} from "@tanstack/react-query";


export function WorkOrderCreate() {
    return (
        <Create redirect={"list"}>
            <WorkOrderForm mode={"CREATE"}/>
        </Create>
    )
}

export function WorkOrderShow() {
    return (
        <Show>
            <WorkOrderForm readonly mode={"SHOW"}/>
        </Show>
    )
}

export function WorkOrderEdit() {
    return (
        <Edit mutationMode={"pessimistic"}>
            <WorkOrderForm mode={"EDIT"}/>
        </Edit>
    )
}

const RunWorkOrderButton = () => {
        const {selectedIds} = useListContext();
        const refresh = useRefresh();
        const notify = useNotify();
        const unselectAll = useUnselectAll('posts');
        const {mutate, isPending} = useMutation({
            mutationFn: (data: any) => {
                return fetch('/api/work-order/run', {
                    method: 'POST',
                    body: JSON.stringify(data),
                })
            },
            onSuccess: () => {
                notify('Success run');
                unselectAll();
            },
            onError: () => {
                notify('Error: run', {
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
                <PlayCircleFilledIcon/>
            </Button>
        );
    }
;

const WorkOrderBulkActionButtons = () => (
    <>
        <RunWorkOrderButton/>
    </>
);

export const WorkOrderList = () => (
    <List>
        <Datagrid rowClick={"edit"} bulkActionButtons={<WorkOrderBulkActionButtons/>}>
            <TextField source="id"/>
            <ChipField source="status"/>
            <DateField source="created_at"/>
            <DateField source="updated_at"/>
            <DeleteButton/>
        </Datagrid>
    </List>
);
