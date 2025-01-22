import { CircularProgress, Stack, Typography } from "@mui/material";
import { useQuery } from '@tanstack/react-query';
import { Datagrid, Edit, List, SimpleForm, TextField, TextInput, useGetOne } from "react-admin";
import MUITextField from "@mui/material/TextField";

export const ConfigList = () => {
    const { data, isPending, error } = useQuery({
        queryKey: ['server-info'],
        queryFn: async ({ signal }) => {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/ping`)

            return await response.json()
        }
    });


    return (
        <Stack gap={5}>
            <Stack gap={2}>
                <Typography>
                    Server Info
                </Typography>
                {isPending ? <CircularProgress /> :
                    <Stack>
                        <MUITextField label="Server IP" value={data.serverIP} aria-readonly InputProps={{
                            readOnly: true,
                        }} />
                        <MUITextField label="Server Port" value={data.port} aria-readonly InputProps={{
                            readOnly: true,
                        }} />
                        <MUITextField label="Revision" value={data.revision} aria-readonly InputProps={{
                            readOnly: true,
                        }} />
                        <MUITextField label="Version" value={data.version} aria-readonly InputProps={{
                            readOnly: true,
                        }} />
                    </Stack>
                }
            </Stack>

            <Stack gap={2}>
                <Typography>
                    Config
                </Typography>

                <List resource="config" actions={false} pagination={false}>
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="id" />
                        <TextField source="value" />
                    </Datagrid>
                </List>
            </Stack>
        </Stack>
    )
};


export const ConfigEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="id" readOnly />
            <TextInput source="value" />
        </SimpleForm>
    </Edit>
);
