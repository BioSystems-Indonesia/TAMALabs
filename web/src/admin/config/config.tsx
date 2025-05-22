import { Card, CardContent, CircularProgress, Stack, Typography } from "@mui/material";
import { useQuery } from '@tanstack/react-query';
import { Datagrid, Edit, List, SimpleForm, TextField, TextInput } from "react-admin";
import MUITextField from "@mui/material/TextField";
import useAxios from "../../hooks/useAxios";

export const ConfigList = () => {
    const axios = useAxios()
    const { data, isPending } = useQuery({
        queryKey: ['server-info'],
        queryFn: async ({ signal }) => {
            const { data } = await axios.get('/ping')
            return data
        }
    });


    return (
        <Stack gap={5}>
            <Card>
                <CardContent>
                    <Stack gap={2}>
                        <Typography>
                            Server Info
                        </Typography>
                        {isPending ? <CircularProgress /> :
                            <Stack gap={2}>
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
                </CardContent>
            </Card>

            <Card>
                <CardContent>
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
                </CardContent>
            </Card>
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
