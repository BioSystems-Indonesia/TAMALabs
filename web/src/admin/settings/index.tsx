import { zodResolver } from '@hookform/resolvers/zod';
import CircularProgress from '@mui/material/CircularProgress';
import Stack from "@mui/material/Stack";
import { useEffect, useState } from "react";
import { Labeled, NumberInput, RadioButtonGroupInput, SaveButton, TabbedForm, Toolbar, required, useNotify, useStoreContext, type SaveHandler } from "react-admin";
import { Settings, defaultSettings, orientationChoices, settingSchema, settingsStoreKey } from "../../types/setting";
import useSettings from '../../hooks/useSettings';


function SettingsToolbar() {
    return (
        <Toolbar>
            <SaveButton />
        </Toolbar>
    )
}



export default function SettingsPage() {
    const store = useStoreContext();
    const [loading, setLoading] = useState(true);
    const [settings, setSettings] = useSettings();

    useEffect(() => {
        setLoading(true);

        const localSetting = store.getItem<Settings>(settingsStoreKey);
        setSettings(settingSchema.parse({
            ...defaultSettings,
            ...localSetting,
        }));

        setLoading(false);
    }, [store])

    const notify = useNotify();
    const onSubmit: SaveHandler<Settings> = async (data: Partial<Settings>): Promise<void> => {
        try {
            store.setItem(settingsStoreKey, data);
            notify('Success update settings', {
                type: 'success',
            });
        } catch (error) {
            console.error(error);
            notify('Error update settings', {
                type: 'error',
            });
        }
    }

    if (loading) {
        return (
            <Stack sx={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                height: "100vh",
            }}>
                <CircularProgress />
            </Stack>
        )
    }

    return (
        <TabbedForm onSubmit={onSubmit} toolbar={<SettingsToolbar />} resolver={zodResolver(settingSchema)} record={settings}
        >
            <TabbedForm.Tab label="Personal">
                <Labeled label="Barcode">
                    <Stack>
                        <Stack direction={"row"} gap={1} sx={{
                            xs: {
                                width: "100%",
                            },
                            md: {
                                width: "50%",
                            },
                        }} >
                            <NumberInput source="barcode_page_width" label="Page Width (mm)" validate={[required()]} />
                            <NumberInput source="barcode_page_height" label="Page Height (mm)" validate={[required()]} />
                        </Stack>
                        <Stack direction={"row"} gap={1} sx={{
                            xs: {
                                width: "100%",
                            },
                            md: {
                                width: "50%",
                            },
                        }} >
                            <NumberInput source="barcode_width" label="Width" validate={[required()]} />
                            <NumberInput source="barcode_height" label="Height" validate={[required()]} />
                        </Stack>
                        <RadioButtonGroupInput choices={[...orientationChoices]} source="barcode_orientation" label="Orientation" validate={[required()]} />
                    </Stack>
                </Labeled>
            </TabbedForm.Tab>
        </TabbedForm>
    )
}