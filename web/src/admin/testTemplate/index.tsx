import { Divider } from "@mui/material";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import { useEffect, useState } from "react";
import { Create, Datagrid, DeleteButton, Edit, FilterLiveSearch, List, NumberField, SaveButton, SimpleForm, TextField, TextInput, Toolbar, required, useEditContext, useNotify, useSaveContext } from "react-admin";
import { useFormContext } from "react-hook-form";
import type { ObservationRequestCreateRequest } from "../../types/observation_requests";
import type { ActionKeys } from "../../types/props";
import { TestInput, testTypesField } from '../workOrder/Form';

export const TestTemplateList = () => (
    <List aside={<TestTemplateFilterSidebar />} title="Test Template">
        <Datagrid bulkActionButtons={false}>
            <NumberField source="id" />
            <TextField source="name" />
            <TextField source="description" />
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
    const [isLoading, setIsLoading] = useState(false);
    const [selectedType, setSelectedType] = useState<Record<number, ObservationRequestCreateRequest>>({});

    if (props.mode === "EDIT") {
        const { record, isPending } = useEditContext();

        useEffect(() => {
            if (isPending) {
                setIsLoading(true);
            } else {
                setIsLoading(false);
            }
        }, [isPending])

        useEffect(() => {
            if (record) {
                setIsLoading(true);
                const newSelectedType: Record<number, ObservationRequestCreateRequest> = {};
                record.test_types.forEach((v: ObservationRequestCreateRequest) => {
                    newSelectedType[v.test_type_id] = v
                })

                setSelectedType(newSelectedType)
                setIsLoading(false);
            }
        }, [record])
    }


    if (isLoading) {
        return <></>
    }

    return (
        <SimpleForm disabled={props.readonly} toolbar={false}>
            <TestTypeToolbar />
            <Divider sx={{
                marginBottom: "36px",
            }} />
            <TextInput source="name" readOnly={props.readonly} validate={[required()]} />
            <TextInput source="description" readOnly={props.readonly} multiline />
            <TestInput initSelectedType={selectedType} />
        </SimpleForm>
    )
}

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

        if (!data[testTypesField] || data[testTypesField].length === 0) {
            notify("Please select test", {
                type: "error",
            });
            return;
        }

        if (save) {
            const observationRequest = data[testTypesField] as Record<number, ObservationRequestCreateRequest>
            save({
                ...data,
                test_types: Object.entries(observationRequest).map(([_, value]) => {
                    return value
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
        }} emptyWhileLoading>
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
