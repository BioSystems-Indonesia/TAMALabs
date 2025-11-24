import { Card, CardContent, Chip, Stack, useTheme } from '@mui/material';
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import {
    Create,
    DataTable,
    DateField,
    DateTimeInput,
    Edit,
    FilterLiveForm,
    List,
    ReferenceManyField,
    required,
    SearchInput,
    SelectInput,
    Show,
    SimpleForm,
    TextInput,
} from "react-admin";
import CustomDateInput from "../../component/CustomDateInput.tsx";
import FeatureList from "../../component/FeatureList.tsx";
import SideFilter from '../../component/SideFilter.tsx';
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import { Action, ActionKeys } from "../../types/props.ts";
import { ResultDataGrid } from "../result/index.tsx";

export type PatientFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

const NullableField = ({ value }: { value: any }) => (
    <span style={{
        color: !value || value === '' ? '#888' : 'inherit',
        fontStyle: !value || value === '' ? 'italic' : 'normal',
        opacity: !value || value === '' ? 0.6 : 1,
        fontSize: !value || value === '' ? '0.875rem' : 'inherit'
    }}>
        {value || 'null'}
    </span>
);

function ReferenceSection() {
    const theme = useTheme();

    return (
        <Card
            elevation={0}
            sx={{
                margin: 'auto',
                width: '100%',
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: 2,
                overflow: 'hidden',
                mb: 4
            }}
        >
            <CardContent sx={{ p: 3 }}>
                <Typography
                    variant="h6"
                    sx={{
                        fontWeight: 600,
                        color: theme.palette.text.primary,
                        mb: 3
                    }}
                >
                    📊 Test Results
                </Typography>
                <ReferenceManyField label={""} reference="result" target="patient_ids">
                    <ResultDataGrid />
                </ReferenceManyField>
            </CardContent>
        </Card>
    )
}

export function PatientFormField(props: PatientFormProps) {
    const theme = useTheme();

    return (
        <Stack spacing={3} sx={{ width: '100%' }}>
            {props.mode !== Action.CREATE && (
                <Card
                    elevation={0}
                    sx={{
                        border: `1px solid ${theme.palette.divider}`,
                        borderRadius: 2,
                    }}
                >
                    <CardContent sx={{ p: 3 }}>
                        <Typography
                            variant="subtitle1"
                            sx={{
                                fontWeight: 600,
                                color: theme.palette.text.primary
                            }}
                        >
                            ℹ️ System Information
                        </Typography>

                        <TextInput
                            source={"id"}
                            readOnly={true}
                            sx={{
                                mt: 3,
                                '& .MuiOutlinedInput-root': {
                                    borderRadius: 2,
                                    transition: 'all 0.2s ease',
                                }
                            }}
                        />
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <DateTimeInput
                                source={"created_at"}
                                readOnly={true}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease'
                                    }
                                }}
                            />
                            <DateTimeInput
                                source={"updated_at"}
                                readOnly={true}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease'
                                    }
                                }}
                            />
                        </Stack>
                    </CardContent>
                </Card>
            )}

            <Card
                elevation={0}
                sx={{
                    border: `1px solid ${theme.palette.divider}`,
                    borderRadius: 2
                }}
            >
                <CardContent sx={{ p: 3 }}>
                    <Box sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1.5,
                        mb: 3
                    }}>

                        <Typography
                            variant="subtitle1"
                            sx={{
                                fontWeight: 600,
                                color: theme.palette.text.primary
                            }}
                        >
                            ❗Required Information
                        </Typography>
                        <Chip
                            label="Required"
                            size="small"
                            color="error"
                            variant="outlined"
                            sx={{ ml: 'auto', fontSize: '0.75rem' }}
                        />
                    </Box>
                    <Stack>
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <TextInput
                                source="first_name"
                                validate={[required()]}
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                            <TextInput
                                source="last_name"
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                        </Stack>

                        <Stack direction={"row"} gap={1.9} width={"100%"}>
                            <CustomDateInput
                                source={"birthdate"}
                                label={"Birth Date"}
                                required
                                sx={{
                                    maxWidth: null,
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                            <FeatureList source={"sex"} types={"sex"}>
                                <SelectInput
                                    source="sex"
                                    validate={[required()]}
                                    readOnly={props.readonly}
                                    sx={{
                                        '& .MuiOutlinedInput-root': {
                                            borderRadius: 2,
                                            transition: 'all 0.2s ease',
                                            ...(!props.readonly && {
                                                '&:hover': {
                                                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                                }
                                            })
                                        }
                                    }}
                                />
                            </FeatureList>
                        </Stack>
                    </Stack>
                </CardContent>
            </Card>

            <Card
                elevation={0}
                sx={{
                    border: `1px solid ${theme.palette.divider}`,
                    borderRadius: 2
                }}
            >
                <CardContent sx={{ p: 3 }}>
                    <Box sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1.5,
                        mb: 3
                    }}>
                        <Typography
                            variant="subtitle1"
                            sx={{
                                fontWeight: 600,
                                color: theme.palette.text.primary
                            }}
                        >
                            📋 Additional Information
                        </Typography>
                        <Chip
                            label="Optional"
                            size="small"
                            color="default"
                            variant="outlined"
                            sx={{ ml: 'auto', fontSize: '0.75rem' }}
                        />
                    </Box>

                    <Stack>
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <TextInput
                                source="phone_number"
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                            <TextInput
                                source="location"
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: 2,
                                        transition: 'all 0.2s ease',
                                        ...(!props.readonly && {
                                            '&:hover': {
                                                boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                            }
                                        })
                                    }
                                }}
                            />
                        </Stack>

                        <TextInput
                            source="address"
                            readOnly={props.readonly}
                            multiline
                            rows={3}
                            sx={{
                                '& .MuiOutlinedInput-root': {
                                    borderRadius: 2,
                                    transition: 'all 0.2s ease',
                                    ...(!props.readonly && {
                                        '&:hover': {
                                            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                                        }
                                    })
                                }
                            }}
                        />
                    </Stack>
                </CardContent>
            </Card>
        </Stack>
    )
}

export function PatientForm(props: PatientFormProps) {
    return (
        <Box sx={{ p: { xs: 2, sm: 3 } }}>
            <SimpleForm
                disabled={props.readonly}
                toolbar={props.readonly === true ? false : undefined}
                warnWhenUnsavedChanges
                sx={{
                    '& .RaSimpleForm-form': {
                        backgroundColor: 'transparent',
                        boxShadow: 'none',
                        padding: 0
                    }
                }}
            >
                <PatientFormField {...props} />
            </SimpleForm>
        </Box>
    )
}

export function PatientCreate() {
    const redirect = useRefererRedirect("show");
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Create redirect={redirect} resource="patient">
                <PatientForm mode={"CREATE"} />
            </Create>
        </Box>
    )
}

export function PatientShow() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Show resource="patient">
                <PatientForm readonly mode={"SHOW"} />
                <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
                    <ReferenceSection />
                </Box>
            </Show>
        </Box>
    )
}

export function PatientEdit() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Edit resource="patient">
                <PatientForm mode={"EDIT"} />
            </Edit>
        </Box>
    )
}

const PatientFilterSidebar = () => {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';

    return (
        <SideFilter sx={{
            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
        }}>
            <FilterLiveForm debounce={1500}>
                <Stack spacing={0}>
                    <Box>
                        <Typography variant="h6" sx={{
                            color: theme.palette.text.primary,
                            marginBottom: 2,
                            fontWeight: 600,
                            fontSize: '1.1rem',
                            textAlign: 'center'
                        }}>
                            👤 Filter Patients
                        </Typography>
                    </Box>
                    <SearchInput
                        source="q"
                        alwaysOn
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                borderRadius: '12px',
                                transition: 'all 0.3s ease',
                                border: isDarkMode ? `1px solid ${theme.palette.divider}` : '1px solid #e5e7eb',
                                '&:hover': {
                                    backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                },
                                '&.Mui-focused': {
                                    backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
                                }
                            },
                            '& .MuiInputLabel-root': {
                                color: theme.palette.text.secondary,
                                fontWeight: 500,
                            }
                        }}
                    />
                    <Box>
                        <Typography variant="body2" sx={{
                            color: theme.palette.text.secondary,
                            marginBottom: 1.5,
                            fontSize: '0.85rem',
                            fontWeight: 500
                        }}>
                            📅 Birth Date Filter
                        </Typography>
                        <Stack>
                            <CustomDateInput
                                source={"birthdate"}
                                label={"Birth Date"}
                                clearable
                                size="small"
                                sx={{
                                    '& .MuiOutlinedInput-root': {
                                        backgroundColor: isDarkMode ? theme.palette.action.hover : '#f9fafb',
                                        borderRadius: '12px',
                                        border: isDarkMode ? `1px solid ${theme.palette.divider}` : '1px solid #e5e7eb',
                                        transition: 'all 0.3s ease',
                                        '&:hover': {
                                            backgroundColor: isDarkMode ? theme.palette.action.selected : '#f3f4f6',
                                            borderColor: isDarkMode ? theme.palette.primary.main : '#9ca3af',
                                            boxShadow: isDarkMode
                                                ? '0 4px 12px rgba(255, 255, 255, 0.1)'
                                                : '0 4px 12px rgba(0, 0, 0, 0.1)',
                                        },
                                        '&.Mui-focused': {
                                            backgroundColor: isDarkMode ? theme.palette.background.paper : 'white',
                                            borderColor: theme.palette.primary.main,
                                            boxShadow: `0 0 0 3px ${theme.palette.primary.main}30`,
                                        }
                                    },
                                    '& .MuiInputLabel-root': {
                                        color: theme.palette.text.primary,
                                        fontWeight: 500,
                                        fontSize: '0.85rem',
                                    }
                                }}
                            />
                        </Stack>
                    </Box>
                </Stack>
            </FilterLiveForm>
        </SideFilter>
    )
};


export const PatientList = () => {
    return (
        <List aside={<PatientFilterSidebar />} sort={{
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
            <DataTable>
                <DataTable.Col source="id" />
                <DataTable.Col source="first_name" render={(record: any) => (
                    <NullableField value={record.first_name} />
                )} />
                <DataTable.Col source="last_name" render={(record: any) => (
                    <NullableField value={record.last_name} />
                )} />
                <DataTable.Col source="birthdate" >
                    <DateField source="birthdate" locales={["id-ID"]} />
                </DataTable.Col>
                <DataTable.Col source="sex" render={(record: any) => (
                    <NullableField value={record.sex} />
                )} />
                <DataTable.Col source="location" render={(record: any) => (
                    <NullableField value={record.location} />
                )} />
                <DataTable.Col source="created_at">
                    <DateField source="created_at" showTime />
                </DataTable.Col>
                <DataTable.Col source="updated_at">
                    <DateField source="updated_at" showTime />
                </DataTable.Col>
            </DataTable>
        </List>
    )
};

// Export PatientInfoCard component
export { default as PatientInfoCard } from './PatientInfoCard';

