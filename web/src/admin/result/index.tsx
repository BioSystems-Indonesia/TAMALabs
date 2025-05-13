import {Chip, Divider, Stack, Typography} from "@mui/material";
import dayjs from "dayjs";
import {
    AutocompleteArrayInput,
    BooleanInput,
    Button,
    Datagrid,
    DateField,
    FilterLiveForm,
    Link,
    List,
    NumberField,
    ReferenceInput,
    TopToolbar,
    useNotify,
    WithRecord
} from "react-admin";
import RefreshIcon from '@mui/icons-material/Refresh';
import CustomDateInput from "../../component/CustomDateInput";
import PrintReportButton from "../../component/PrintReport";
import SideFilter from "../../component/SideFilter";
import type {WorkOrder} from "../../types/work_order";
import {WorkOrderChipColorMap} from "../workOrder/ChipFieldStatus";
import {FilledPercentChip} from "./component";
import useAxios from "../../hooks/useAxios";


export const ResultList = () => (
    <List
        resource="result"
        sort={{field: "id", order: "DESC"}}
        aside={<ResultSideFilter/>}
        filters={ResultMoreFilter}
        filterDefaultValues={{
            created_at_start: dayjs().subtract(7, "day").toISOString(),
            created_at_end: dayjs().toISOString(),
        }}
        actions={<ResultActions/>}
        storeKey={false}
        exporter={false}
        disableSyncWithLocation
        sx={{
            '& .RaList-main': {
                marginTop: '-14px'
            },
            '& .RaList-content': {
                backgroundColor: 'background.paper',
                padding: 2,
                borderRadius: 1,
            },
        }}
    >
        <ResultDataGrid/>
    </List>
);

function ResultActions() {
    const axios = useAxios()
    const notify = useNotify()
    return (
        <TopToolbar>
            <Button label={"Refresh"} onClick={() => {
                axios.post("/result/refresh").then(() => {
                    notify("Refresh Result Success", {
                        type: "success"
                    })
                }).catch(() => {
                    notify("Refresh Result Failed", {
                        type: "error"
                    })
                })
            }}>
                <RefreshIcon/>
            </Button>
        </TopToolbar>
    )
}

export const ResultDataGrid = (props: any) => {
    return (
        <Datagrid bulkActionButtons={false} >
            <NumberField source="id"/>
            <WithRecord label="Patient" render={(record: any) => (
                <Link to={`/patient/${record.patient.id}/show`} resource="patient" label={"Patient"}
                      onClick={e => e.stopPropagation()}>
                    #{record.patient.id}-{record.patient.first_name} {record.patient.last_name}
                </Link>
            )}/>
            <WithRecord label="Request" render={(record: any) => (
                <Link to={`/work-order/${record.id}/show`} label={"Work Order"} onClick={e => e.stopPropagation()}>
                    <Chip label={`#${record.id} - ${record.status}`} color={WorkOrderChipColorMap(record.status)}/>
                </Link>
            )}/>
            <WithRecord label="Request" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_request}
                </Typography>
            )}/>
            <WithRecord label="Result" render={(record: WorkOrder) => (
                <Typography variant="body2">
                    {record.total_result_filled}
                </Typography>
            )}/>
            <WithRecord label="Filled" render={(record: WorkOrder) => (
                <FilledPercentChip percent={record.percent_complete}/>
            )}/>
            <DateField source="created_at" showDate showTime/>
            <WithRecord label="Print Result" render={(record: any) => (
                <PrintReportButton results={record.test_result} patient={record.patient} workOrder={record}/>
            )}/>
        </Datagrid>
    )
}

function ResultSideFilter() {
    return (
        <SideFilter>
            <FilterLiveForm debounce={1500}>
                <Stack>
                    <ReferenceInput source={"patient_ids"} reference="patient" label={"Patient"} alwaysOn>
                        <AutocompleteArrayInput size="small"/>
                    </ReferenceInput>
                    <Divider sx={{
                        marginBottom: 2,
                    }}/>
                    <CustomDateInput label={"Created At Start"} source="created_at_start" alwaysOn size="small"
                                     disableFuture/>
                    <CustomDateInput label={"Created At End"} source="created_at_end" alwaysOn size="small"
                                     disableFuture/>
                </Stack>
            </FilterLiveForm>
        </SideFilter>
    )

}


const ResultMoreFilter = [
    <BooleanInput source={"has_result"} label={"Show Only With Result"}/>,
]
