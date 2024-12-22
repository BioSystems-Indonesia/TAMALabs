import {Datagrid, DateField, List, ReferenceField, Show, SimpleShowLayout, TextField} from "react-admin";

export function ObservationRequestShow() {
    return (
        <Show>
            <SimpleShowLayout sx={{}}>
                <TextField source="id"/>
                <TextField source="test_code"/>
                <TextField source="test_description"/>
                <ReferenceField reference={"work-order"} source={"order_id"}/>
                <ReferenceField reference={"specimen"} source={"specimen_id"}/>
                <DateField source="created_at" showTime/>
                <DateField source="updated_at" showTime/>
            </SimpleShowLayout>
        </Show>
    )
}


export const ObservationRequestList = () => (
    <List>
        <Datagrid>
            <TextField source="id"/>
            <TextField source="test_code"/>
            <TextField source="test_description"/>
            <ReferenceField reference={"work-order"} source={"order_id"}/>
            <ReferenceField reference={"specimen"} source={"specimen_id"}/>
            <DateField source="created_at" showTime/>
            <DateField source="updated_at" showTime/>
        </Datagrid>
    </List>
);
