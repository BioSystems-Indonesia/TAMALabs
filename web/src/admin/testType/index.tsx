import Box from "@mui/material/Box";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { ArrayInput, AutocompleteInput, Create, Datagrid, Edit, List, NumberInput, SimpleForm, SimpleFormIterator, TextField, TextInput, required } from "react-admin";
import { useFormContext } from "react-hook-form";
import { useSearchParams } from "react-router-dom";
import FeatureList from "../../component/FeatureList";
import type { ActionKeys } from "../../types/props";
import type { Unit } from "../../types/unit";
import { TestFilterSidebar } from "../workOrder/TestTypeFilter";

export const TestTypeDatagrid = (props: any) => {
    return (
        <Datagrid bulkActionButtons={false}>
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="code" />
            <TextField source="category" />
            <TextField source="sub_category" />
            <TextField source="low_ref_range" label="low" />
            <TextField source="high_ref_range" label="high" />
            <TextField source="unit" />
            <TextField source="type" />
            <TextField source="decimal" />
            <TextField source="description" />
        </Datagrid>
    )
}

export const TestTypeList = () => {
    const [selectedData, setSelectedData] = useState<any>([]);
    
    return (
    <List aside={<TestFilterSidebar selectedData={selectedData} setSelectedData={setSelectedData} />} title="Test Type" sort={{
        field: "id",
        order: "DESC",
    }}>
        <TestTypeDatagrid />
    </List>
)};


function ReferenceSection() {
    return (
        <Box sx={{ width: "100%" }}>
        </Box>
    )
}

type TestTypeFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function TestTypeInput(props: TestTypeFormProps) {
    const { data: filter, isLoading: isFilterLoading } = useQuery({
        queryKey: ['filterTestType'],
        queryFn: () => fetch(import.meta.env.VITE_BACKEND_BASE_URL + '/test-type/filter').then(res => res.json()),
    });

    const [categories, setCategories] = useState<string[]>([]);
    const [subCategories, setSubCategories] = useState<string[]>([]);

    useEffect(() => {
        if (filter) {
            setCategories(filter.categories);
            setSubCategories(filter.sub_categories);
        }
    }, [filter, isFilterLoading]);


    const { data: units, isLoading: isUnitLoading } = useQuery<Unit[]>({
        queryKey: ['unit'],
        queryFn: () => fetch(import.meta.env.VITE_BACKEND_BASE_URL + '/unit').then(res => res.json()),
    });

    const [unit, setUnit] = useState<string[]>([]);
    useEffect(() => {
        if (units && Array.isArray(units)) {
            // Extract the `value` property from each object in the `units` array
            const unitValues = units.map(unit => unit.value);
            setUnit(unitValues);
        }
    }, [units, isUnitLoading]);

    const { setValue } = useFormContext();
    const [params] = useSearchParams()
    useEffect(() => {
        if (params.has("code")) {
            const code = params.get("code")
            setValue("code", code)
            setValue("name", code)
        }
    }, [params.has("code")])

    return (
        <>
            <TextInput source="name" readOnly={props.readonly} validate={[required()]} />
            <TextInput source="code" readOnly={props.readonly} validate={[required()]} />
            <AutocompleteInput source="category" readOnly={props.readonly} filterSelectedOptions={false}
                loading={isFilterLoading}
                choices={categories.map(val => {
                    return { id: val, name: val }
                })}
                onCreate={val => {
                    if (!val || categories.includes(val)) {
                        return;
                    }

                    const newCategories = [...categories, val];
                    setCategories(newCategories);
                    return { id: val, name: val }
                }} />
            <AutocompleteInput source="sub_category"
                readOnly={props.readonly}
                loading={isFilterLoading}
                choices={subCategories.map(val => {
                    return { id: val, name: val }
                })}
                onCreate={val => {
                    console.log("onCreate", val, subCategories);
                    if (!val || subCategories.includes(val)) {
                        return;
                    }

                    const newSubCategories = [...subCategories, val];
                    setSubCategories(newSubCategories);
                    return { id: val, name: val }
                }} />
            <NumberInput source="low_ref_range" label="Low Range" readOnly={props.readonly} validate={[required()]} />
            <NumberInput source="high_ref_range" label="High Range" readOnly={props.readonly} validate={[required()]} />
            <AutocompleteInput source="unit"
                readOnly={props.readonly}
                loading={isUnitLoading}
                choices={unit.map(val => {
                    return { id: val, name: val }
                })}
                onCreate={val => {
                    console.log("onCreate", val, unit);
                    if (!val || unit.includes(val)) {
                        return;
                    }

                    const newUnit = [...unit, val];
                    setUnit(newUnit);
                    return { id: val, name: val }
                }} />
            <NumberInput source="decimal" readOnly={props.readonly} validate={[required()]} />
            <ArrayInput source="types">
                <SimpleFormIterator inline>
                    <FeatureList source="type" readOnly={props.readonly} types="specimen-type" >
                        <AutocompleteInput source="type" readOnly={props.readonly} />
                    </FeatureList>
                </SimpleFormIterator>
            </ArrayInput>
            <TextInput source="description" readOnly={props.readonly} />
        </>

    )
}

function TestTypeForm(props: TestTypeFormProps) {
    return (
        <SimpleForm>
            <TestTypeInput {...props} />
        </SimpleForm>
    )
}


export function TestTypeEdit() {
    return (
        <Edit mutationMode="pessimistic" title="Edit Test Type" redirect={"list"}>
            <TestTypeForm readonly={false} mode={"EDIT"} />
            <ReferenceSection />
        </Edit>
    )
}

export function TestTypeCreate() {
    return (
        <Create title="Create Test Type" redirect={"list"}>
            <TestTypeForm readonly={false} mode={"CREATE"} />
            <ReferenceSection />
        </Create>
    )
}
