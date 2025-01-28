import {
    ChipField,
    Create,
    Datagrid,
    DateField,
    DeleteButton,
    List,
    ShowButton,
    TextField,
    TopToolbar,
    WithRecord
} from "react-admin";
import { useParams, useSearchParams } from "react-router-dom";
import WorkOrderForm from "./Form.tsx";
import { WorkOrderChipColorMap } from "./ChipFieldStatus.tsx";
import Chip from "@mui/material/Chip";

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
            }} resource={`work-order/${id}/show/add-test`}>
            <WorkOrderForm mode={"ADD_TEST"} />
        </Create>
    )
}


export const WorkOrderList = () => (
    <List sort={{
        field: "id",
        order: "DESC"
    }}>
        <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <WithRecord render={(record: any) => (
                <Chip label={`${record.status}`} color={WorkOrderChipColorMap(record.status)} />
            )} />
            <DateField source="created_at" />
            <DateField source="updated_at" />
            <DeleteButton />
        </Datagrid>
    </List>
);
