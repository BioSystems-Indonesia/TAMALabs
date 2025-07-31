import CancelOutlinedIcon from '@mui/icons-material/CancelOutlined';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import { Button, Chip, Divider, Stack, Typography } from "@mui/material";
import { useMutation } from "@tanstack/react-query";
import dayjs from "dayjs";
import {
    AutocompleteArrayInput,
    Datagrid,
    DateField,
    FilterLiveForm,
    Link,
    List,
    NumberField,
    ReferenceInput,
    useNotify,
    useRefresh,
    WithRecord
} from "react-admin";
import CustomDateInput from "../../component/CustomDateInput";
import SideFilter from "../../component/SideFilter";
import { useCurrentUser } from "../../hooks/currentUser";
import type { WorkOrder } from "../../types/work_order";
import { FilledPercentChip, VerifiedChip } from '../result/component';
import { WorkOrderChipColorMap } from "../workOrder/ChipFieldStatus";
import useAxios from '../../hooks/useAxios';


export const ApprovalList = () => {
    const currentUser = useCurrentUser();

    return (
    <List
        resource="result"
        sort={{ field: "id", order: "DESC" }}
        aside={<ApprovalSideFilter />}
        filterDefaultValues={{
            created_at_start: dayjs().subtract(7, "day").toISOString(),
            created_at_end: dayjs().toISOString(),
        }}
        filter={{
            doctor_ids: [currentUser?.id],
        }}
        exporter={false}
        sx={{
            '& .RaList-main': {
                // marginTop: '-14px'
            },
            '& .RaList-content': {
                backgroundColor: 'background.paper',
                padding: 2,
                borderRadius: 1,
            },
        }}
    >
        <ApprovalDataGrid />
    </List>
)};

export const ApprovalDataGrid = (props: any) => {
    const axios = useAxios();
    const refresh = useRefresh();
    const notify = useNotify();
    const { mutate: verifyResult, isPending: verifyIsPending } = useMutation({
        mutationFn: (id: number) => axios.post(`/result/${id}/approve`),
        onSuccess: () => {
            refresh()
            notify('Result verified successfully', {
                type: 'success',
            })
        },
        onError: () => {
            notify('Result verified failed', {
                type:'error',
            })
        }
    });

    const { mutate: rejectResult, isPending: rejectIsPending } = useMutation({
        mutationFn: (id: number) => axios.post(`/result/${id}/reject`),
        onSuccess: () => {
            refresh()
            notify('Result rejected successfully', {
                type:'success',
            })
        },
        onError: () => {
            notify('Result rejected failed', {
                type:'error',
            })
        }
    });
    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id" />
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"}
                    onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )} />
            <WithRecord label="Request" render={(record: any) => (
                <Link to={`/work-order/${record.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.id} - ${record.status}`} color={WorkOrderChipColorMap(record.status)} />
                </Link>
            )} />
            <WithRecord label="Request" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_request}
                </Typography>
            )} />
            <WithRecord label="Result" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_result_filled}
                </Typography>
            )} />
            <WithRecord label="Filled" render={(record: WorkOrder) => (
                <FilledPercentChip percent={record.percent_complete} />
            )} />
            <WithRecord label="Verified" render={(record: WorkOrder) => (
                <VerifiedChip verified={record.verified_status !== '' ? record.verified_status : "VERIFIED"} />
            )} />
            <DateField source="created_at" showDate showTime />
            <WithRecord label="Action" render={(record: WorkOrder) => (
                <Stack direction="row" spacing={1}>
                    <Button
                        variant="contained"
                        loading={verifyIsPending}
                        color="success"
                        sx={{color: 'white !important'}}
                        startIcon={<CheckCircleOutlineIcon />}
                        onClick={(e) => {
                            e.stopPropagation();
                            verifyResult(record.id);
                        }}
                        size="small"
                        disabled={record.verified_status === "VERIFIED"} 
                    >
                        Verify
                    </Button>
                    <Button
                        variant="contained"
                        color="error"
                        loading={rejectIsPending}
                        startIcon={<CancelOutlinedIcon />}
                        onClick={(e) => {
                            e.stopPropagation();
                            rejectResult(record.id);
                        }}
                        size="small"
                        disabled={record.verified_status === "REJECTED"} 
                    >
                        Reject
                    </Button>
                </Stack>
            )} />
        </Datagrid>
    )
}

function ApprovalSideFilter() {
    return (
        <SideFilter>
            <FilterLiveForm debounce={1500}>
                <Stack>
                    <ReferenceInput source={"patient_ids"} reference="patient" label={"Patient"} alwaysOn>
                        <AutocompleteArrayInput size="small" />
                    </ReferenceInput>
                    <Divider sx={{
                        marginBottom: 2,
                    }} />
                    <CustomDateInput label={"Created At Start"} source="created_at_start" alwaysOn size="small"
                        disableFuture />
                    <CustomDateInput label={"Created At End"} source="created_at_end" alwaysOn size="small"
                        disableFuture />
                </Stack>
            </FilterLiveForm>
        </SideFilter>
    )

}

