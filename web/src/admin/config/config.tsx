import { Datagrid, Edit, List, SimpleForm, TextField, TextInput } from "react-admin";

export const ConfigList = () => (
    <List>
        <Datagrid bulkActionButtons={false}>
            <TextField source="id"/>
            <TextField source="value"/>
        </Datagrid>
    </List>
);


export const ConfigEdit = () => (
    <Edit>
    <SimpleForm>
        <TextInput source="id" />
        <TextInput source="value" />
    </SimpleForm>
    </Edit>
);
