import {
    Card,
    CardContent,
    CircularProgress,
    Divider,
    Stack,
    Typography,
    // Box,
} from "@mui/material";
import { useQuery } from "@tanstack/react-query";
import { Datagrid, Edit, List, SimpleForm, TextField, TextInput } from "react-admin";
import MUITextField from "@mui/material/TextField";
import useAxios from "../../hooks/useAxios";
import StorageIcon from "@mui/icons-material/Storage";
import SettingsIcon from "@mui/icons-material/Settings";

export const ConfigList = () => {
    const axios = useAxios();
    const { data, isPending } = useQuery({
        queryKey: ["server-info"],
        queryFn: async ({ signal }) => {
            const { data } = await axios.get("/ping");
            return data;
        },
    });

    return (
        <Stack gap={4}>
            {/* Server Info */}
            <Card sx={{ p: 1 }}>
                <CardContent>
                    <Stack gap={2}>
                        <Stack direction="row" alignItems="center" gap={1}>
                            <StorageIcon color="primary" />
                            <Typography variant="h6">Server Info</Typography>
                        </Stack>
                        <Divider />

                        {isPending ? (
                            <Stack direction="row" alignItems="center" gap={2}>
                                <CircularProgress size={24} />
                                <Typography variant="body2" color="text.secondary">
                                    Loading server info...
                                </Typography>
                            </Stack>
                        ) : (
                            <Stack gap={2}>
                                <MUITextField
                                    label="Server IP"
                                    value={data.serverIP}
                                    variant="outlined"
                                    size="small"
                                    InputProps={{ readOnly: true }}
                                />
                                <MUITextField
                                    label="Server Port"
                                    value={data.port}
                                    variant="outlined"
                                    size="small"
                                    InputProps={{ readOnly: true }}
                                />
                                <MUITextField
                                    label="Revision"
                                    value={data.revision}
                                    variant="outlined"
                                    size="small"
                                    InputProps={{ readOnly: true }}
                                />
                                <MUITextField
                                    label="Version"
                                    value={data.version}
                                    variant="outlined"
                                    size="small"
                                    InputProps={{ readOnly: true }}
                                />
                            </Stack>
                        )}
                    </Stack>
                </CardContent>
            </Card>

            {/* Config List */}
            <Card sx={{ p: 1 }}>
                <CardContent>
                    <Stack gap={2}>
                        <Stack direction="row" alignItems="center" gap={1}>
                            <SettingsIcon color="primary" />
                            <Typography variant="h6">Config</Typography>
                        </Stack>
                        <Divider />

                        <List resource="config" actions={false} pagination={false} sx={{ mt: 1 }} >
                            <Datagrid
                                bulkActionButtons={false}
                                rowClick="edit"
                                sx={{
                                    "& .RaDatagrid-row:hover": {
                                        backgroundColor: "action.hover",
                                        cursor: "pointer",
                                    },
                                }}
                            >
                                <TextField source="id" />
                                <TextField source="value" />
                            </Datagrid>
                        </List>
                    </Stack>
                </CardContent>
            </Card>
        </Stack>
    );
};

export const ConfigEdit = () => (
    <Edit>
        <SimpleForm>
            <TextInput source="id" readOnly />
            <TextInput source="value" />
        </SimpleForm>
    </Edit>
);
