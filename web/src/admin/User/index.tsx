import { Stack } from '@mui/material';
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import { useEffect, useState } from 'react';
import {
    ArrayField,
    BooleanField,
    BooleanInput,
    ChipField,
    ChipFieldProps,
    Create,
    Datagrid,
    DateField,
    DateTimeInput,
    Edit,
    FilterLiveSearch,
    List,
    PasswordInput,
    ReferenceArrayInput,
    required,
    SelectArrayInput,
    Show,
    SimpleForm,
    SingleFieldList,
    TextField,
    TextInput,
    useRecordContext
} from "react-admin";
import SideFilter from '../../component/SideFilter.tsx';
import { getNestedValue } from '../../helper/accessor.ts';
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import { Action, ActionKeys } from "../../types/props.ts";

export type UserFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

export function UserFormField(props: UserFormProps) {
    return (
        <>
            {props.mode !== Action.CREATE && (
                <Box sx={{
                    width: "100%",
                }}>
                    <TextInput source={"id"} readOnly={true} />
                    <Stack direction={"row"} gap={5} width={"100%"}>
                        <DateTimeInput source={"created_at"} readOnly={true} />
                        <DateTimeInput source={"updated_at"} readOnly={true} />
                    </Stack>
                    <Divider />
                </Box>
            )}

            <TextInput source="username" validate={[required()]} readOnly={props.readonly} />
            <TextInput source="fullname" validate={[required()]} readOnly={props.readonly} />
            <TextInput source="email" readOnly={props.readonly} />
            {props.mode !== Action.SHOW && <PasswordInput source="password" validate={
                props.mode === Action.EDIT ? [] : [required()]} readOnly={props.readonly} />}
            <BooleanInput source="is_active" defaultValue={true} readOnly={props.readonly} />
            <ReferenceArrayInput source="roles_id" reference="role">
                <SelectArrayInput optionText="name" readOnly={props.readonly} validate={[required()]} />
            </ReferenceArrayInput>
        </>
    )

}

export function UserForm(props: UserFormProps) {
    return (
        <SimpleForm disabled={props.readonly}
            toolbar={props.readonly === true ? false : undefined}
            warnWhenUnsavedChanges
        >
            <UserFormField {...props} />
        </SimpleForm>
    )
}

export function UserCreate() {
    const redirect = useRefererRedirect("show");

    return (
        <Create redirect={redirect} resource="user" mutationMode='pessimistic'>
            <UserForm mode={"CREATE"} />
        </Create>
    )
}

export function UserShow() {
    return (
        <Show resource="user">
            <UserForm readonly mode={"SHOW"} />
        </Show>
    )
}

export function UserEdit() {
    return (
        <Edit resource="user" mutationMode='pessimistic'>
            <UserForm mode={"EDIT"} />
        </Edit>
    )
}

const UserFilterSidebar = () => (
    <SideFilter>
        <FilterLiveSearch />
    </SideFilter>
);


export const UserList = () => (
    <List aside={<UserFilterSidebar />} sort={{
        field: "id",
        order: "DESC"
    }}
        storeKey={false} exporter={false}
        sx={{
            '& .RaList-content': {
                backgroundColor: 'background.paper',
                padding: 2,
                borderRadius: 1,
            },
        }}>
        <Datagrid>
            <TextField source="id" />
            <TextField source="username" />
            <TextField source="fullname" />
            <TextField source="email" />
            <BooleanField source="is_active" />
            <ArrayField source="roles">
                <SingleFieldList linkType={false}>
                    <UserRolesChipField source="name" size="small" />
                </SingleFieldList>
            </ArrayField>
            <DateField source="created_at" showTime />
            <DateField source="updated_at" showTime />
        </Datagrid>
    </List>
);

export type UserRolesChipFieldProps = Partial<ChipFieldProps> & {
    record?: any
    source: string
}

export function UserRolesChipColorMap(value: string) {
    switch (value) {
        case 'Admin':
            return 'primary';
        case 'Doctor':
            return 'secondary';
        default:
            return 'default';
    }
}

export const UserRolesChipField = (props: UserRolesChipFieldProps) => {
    const [color, setColor] = useState<'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' | undefined>(undefined);
    const record = props.record ?? useRecordContext();

    useEffect(() => {
        if (record === undefined) {
            return;
        }

        const value = getNestedValue(record, props.source);
        const color = UserRolesChipColorMap(value);
        setColor(color);
    }, [record, props.source]);


    return (
        <ChipField {...props} sx={{}} textAlign="center" color={color} source={props.source} />
    )
}
