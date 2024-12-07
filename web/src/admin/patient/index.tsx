import {
    Create,
    Datagrid,
    DateField,
    DateInput,
    Edit,
    List,
    RadioButtonGroupInput,
    required,
    SimpleForm,
    TextField,
    TextInput
} from "react-admin";
import {dateFormatter, dateParser} from "../../helper/format.ts";


function PatientForm() {
    return (
        <SimpleForm>
            <TextInput source="first_name" validate={[required()]} />
            <TextInput source="last_name" validate={[required()]}/>
            <DateInput source="birthdate" defaultValue={new Date()} validate={[required()]}
                       format={dateFormatter}
                       parse={dateParser}
            />
            <RadioButtonGroupInput source="sex" validate={[required()]} choices={[
                {
                    "id": "M",
                    "name": "Male"
                },
                {
                    "id": "F",
                    "name": "Female"
                },
                {
                    "id": "U",
                    "name": "Unknown"
                }
            ]}/>
            <TextInput source="location"/>
        </SimpleForm>
    )
}

export function PatientCreate() {
    return (
        <Create>
            <PatientForm/>
        </Create>
    )
}

export function PatientEdit() {
    return (
        <Edit>
            <PatientForm/>
        </Edit>
    )
}

export const PatientList = () => (
    <List>
        <Datagrid>
            <TextField source="id"/>
            <TextField source="first_name"/>
            <TextField source="last_name"/>
            <DateField source="birthdate"/>
            <TextField source="sex"/>
            <TextField source="location"/>
            <DateField source="created_at"/>
            <DateField source="updated_at"/>
        </Datagrid>
    </List>
);