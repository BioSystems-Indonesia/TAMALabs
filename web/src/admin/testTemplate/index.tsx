import { Divider } from "@mui/material";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import { Create, Datagrid, DeleteButton, Edit, FilterLiveSearch, List, NumberField, ReferenceArrayField, SaveButton, SimpleForm, TextField, TextInput, Toolbar, required, useNotify, useRecordContext, useSaveContext } from "react-admin";
import { useFormContext } from "react-hook-form";
import type { ActionKeys } from "../../types/props";
import type { TestType } from '../../types/test_type';
import { TestInput } from '../workOrder/Form';

export const TestTemplateList = () => (
    <List aside={<TestTemplateFilterSidebar />} title="Test Template">
        <Datagrid bulkActionButtons={false}>
            <NumberField source="id" />
            <TextField source="name" />
            <TextField source="description" />
            <ReferenceArrayField reference="test-type" source="test_type_id" />
        </Datagrid>
    </List>
);

const TestTemplateFilterSidebar = () => {
    return (
        <Card sx={{ order: -1, mr: 2, mt: 2, width: 300 }}>
            <CardContent>
                <FilterLiveSearch />
            </CardContent>
        </Card>
    )
};


type TestTemplateFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function TestTemplateForm(props: TestTemplateFormProps) {
    const record = useRecordContext();

    return (
        <SimpleForm disabled={props.readonly} toolbar={false}>
            <TestTypeToolbar />
            <Divider sx={{
                marginBottom: "36px",
            }} />
            <TextInput source="name" readOnly={props.readonly} validate={[required()]} />
            <TextInput source="description" readOnly={props.readonly} multiline />
            <TestInput initSelectedIds={record?.test_type_id} />
        </SimpleForm>
    )
}

const observationRequestField = "observation_requests";
const TestTemplateSaveButton = ({ disabled }: { disabled?: boolean }) => {
    const { getValues } = useFormContext();
    const { save } = useSaveContext();
    const notify = useNotify();
    const handleClick = (e: any) => {
        e.preventDefault(); // necessary to prevent default SaveButton submit logic
        const { ...data } = getValues();

        if (data == undefined) {
            notify("Please fill in all required fields", {
                type: "error",
            });
            return;
        }

        if (!data[observationRequestField] || data[observationRequestField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[observationRequestField] as TestType[]
            save({
                ...data,
                test_type_id: observationRequest.map((test: TestType) => {
                    return test.id
                })
            });
        }
    };


    return <SaveButton type="button" onClick={handleClick} alwaysEnable size="small" />
}

const TestTypeToolbar = () => {
    return (
        <Stack width={"100%"}
            sx={{
                position: "sticky",
                top: 48,
                borderBottom: "1px solid #ccc",
                zIndex: 2147483647,
                marginBottom: 1,
            }}
        >
            <Toolbar sx={{
                gap: 2,
                width: "100%",
                display: "flex",
                justifyContent: "flex-end",
            }}>
                <DeleteButton variant="contained" size="small" />
                <TestTemplateSaveButton />
            </Toolbar>
        </Stack>
    )
};


export function TestTemplateEdit() {
    return (
        <Edit mutationMode="pessimistic" title="Edit Test Template" sx={{
            "& .RaEdit-card": {
                overflow: "visible",
            }
        }}>
            <TestTemplateForm readonly={false} mode={"EDIT"} />
        </Edit>
    )
}

export function TestTemplateCreate() {
    return (
        <Create title="Create Test Template" redirect={"show"} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }}>
            <TestTemplateForm readonly={false} mode={"CREATE"} />
        </Create>
    )
}
