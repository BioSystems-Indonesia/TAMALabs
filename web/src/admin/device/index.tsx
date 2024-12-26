import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Divider from "@mui/material/Divider";
import {
    Create,
    Datagrid,
    Edit,
    FilterLiveSearch,
    List,
    required,
    Show,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import { Action, ActionKeys } from "../../types/props.ts";

type DeviceFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
        </Box>
    )
}

function DeviceForm(props: DeviceFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
                    toolbar={props.readonly === true ? false : undefined}
                    warnWhenUnsavedChanges
        >
            {props.mode !== Action.CREATE && (
                <div>
                    <TextInput source={"id"} readOnly={true}/>
                    <Divider/>
                </div>
            )}

            <TextInput source="name" validate={[required()]} readOnly={props.readonly}/>
            <TextInput source="ip_address" validate={[required()]} readOnly={props.readonly}/>
        </SimpleForm>
    )
}

export function DeviceCreate() {
    return (
        <Create redirect={"list"}>
            <DeviceForm mode={"CREATE"}/>
        </Create>
    )
}

export function DeviceShow() {
    return (
        <Show>
            <DeviceForm readonly mode={"SHOW"}/>
            <ReferenceSection/>
        </Show>
    )
}

export function DeviceEdit() {
    return (
        <Edit>
            <DeviceForm mode={"EDIT"}/>
        </Edit>
    )
}

const DeviceFilterSidebar = () => (
    <Card sx={{order: -1, mr: 2, mt: 2, width: 300}}>
        <CardContent>
            <FilterLiveSearch/>
        </CardContent>
    </Card>
);


export const DeviceList = () => (
    <List aside={<DeviceFilterSidebar/>}>
        <Datagrid>
            <TextField source="id"/>
            <TextField source="name"/>
            <TextField source="ip_address"/>
        </Datagrid>
    </List>
);