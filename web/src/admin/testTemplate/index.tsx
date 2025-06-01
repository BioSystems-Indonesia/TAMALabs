import { Divider } from "@mui/material";
import Stack from "@mui/material/Stack";
import { useEffect, useState } from "react";
import { AutocompleteArrayInput, Create, Datagrid, DateField, DeleteButton, Edit, FilterLiveSearch, List, NumberField, ReferenceField, ReferenceInput, SaveButton, SimpleForm, TextField, TextInput, Toolbar, required, useEditContext, useNotify, useSaveContext } from "react-admin";
import { useFormContext } from "react-hook-form";
import type { ObservationRequestCreateRequest } from "../../types/observation_requests";
import type { ActionKeys } from "../../types/props";
import { TestInput, testTypesField } from '../workOrder/Form';
import SideFilter from "../../component/SideFilter";
import { RoleNameValue } from "../../types/role";
import { useCurrentUser } from "../../hooks/currentUser";

export const TestTemplateList = () => (
    <List aside={<TestTemplateFilterSidebar />} title="Test Template" sort={{
        field: "id",
        order: "DESC"
    }}
        storeKey={false} exporter={false}
    >
        <Datagrid bulkActionButtons={false}>
            <NumberField source="id" />
            <TextField source="name" />
            <TextField source="description" />
            <DateField source="created_at" showTime/>
            <DateField source="updated_at" showTime/>
            <ReferenceField source="created_by" reference="user"/>
            <ReferenceField source="last_updated_by" reference="user"/>
        </Datagrid>
    </List>
);

const TestTemplateFilterSidebar = () => {
    return (
        <SideFilter>
            <FilterLiveSearch />
        </SideFilter>
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

    const currentUser = useCurrentUser()
    return (
        <SimpleForm disabled={props.readonly} toolbar={false}>
            <TestTypeToolbar />
            <Divider sx={{
                marginBottom: "36px",
            }} />
            <TextInput source="name" readOnly={props.readonly} validate={[required()]} />
            <TextInput source="description" readOnly={props.readonly} multiline />
            <ReferenceInput source={"doctor_ids"} reference="user" resource='user' target="id" label="Doctor" filter={{
                role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
            }}>
                <AutocompleteArrayInput
                    suggestionLimit={10}
                    filterToQuery={(searchText) => ({
                        q: searchText,
                        role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
                    })}
                />
            </ReferenceInput>
            <ReferenceInput source={"analyzers_ids"} reference="user" resource='user' target="id" label="Analyzer" filter={{
            }}>
                <AutocompleteArrayInput
                    suggestionLimit={10}
                    filterToQuery={(searchText) => ({
                        q: searchText,
                    })}
                    defaultValue={[currentUser?.id]}
                />
            </ReferenceInput>
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
