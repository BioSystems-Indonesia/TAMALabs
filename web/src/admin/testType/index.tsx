import { Stack, useTheme, Card, CardContent, Chip } from '@mui/material';
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { ArrayInput, AutocompleteInput, Create, Datagrid, Edit, FunctionField, List, NumberInput, SimpleForm, SimpleFormIterator, TextField, TextInput, required } from "react-admin";
import { useFormContext } from "react-hook-form";
import { useSearchParams } from "react-router-dom";
import FeatureList from "../../component/FeatureList";
import type { ActionKeys } from "../../types/props";
import type { Unit } from "../../types/unit";
import { TestFilterSidebar } from "../workOrder/TestTypeFilter";
import useAxios from "../../hooks/useAxios";
import { TestType } from "../../types/test_type";

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
            <FunctionField 
                label="Types" 
                render={(record: TestType) => 
                    record.types && record.types.length > 0 
                        ? record.types.map(t => t.type).join(", ")
                        : "-"
                } 
            />
            <TextField source="decimal" />
        </Datagrid>
    )
}

export const TestTypeList = () => {
    const [selectedData, setSelectedData] = useState<any>([]);

    return (
        <List
            aside={<TestFilterSidebar selectedData={selectedData} setSelectedData={setSelectedData} />}
            title="Test Type"
            sort={{
                field: "id",
                order: "DESC",
            }}
            sx={{
                '& .RaList-main': {},
                '& .RaList-content': {
                    backgroundColor: 'background.paper',
                    padding: 2,
                    borderRadius: 1,
                },
            }}
            storeKey={false} exporter={false}
        >
            <TestTypeDatagrid />
        </List>
    );
};

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
    const theme = useTheme();
    const axios = useAxios()
    const { data: filter, isLoading: isFilterLoading } = useQuery({
        queryKey: ['filterTestType'],
        queryFn: () => axios.get('/test-type/filter').then(res => res.data),
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
        queryFn: () => axios.get('/unit').then(res => res.data),
    });

    const [unit, setUnit] = useState<string[]>([]);
    useEffect(() => {
        if (units && Array.isArray(units)) {
            const unitValues = units.map(unit => unit.value);
            const uniqueUnits = [...new Set(unitValues)]; 
            setUnit(uniqueUnits);
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
        <Stack spacing={3} sx={{ width: '100%' }}>
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
                            ❗Basic Information
                        </Typography>
                        <Chip 
                            label="Required" 
                            size="small" 
                            color="error" 
                            variant="outlined"
                            sx={{ ml: 'auto', fontSize: '0.75rem' }}
                        />
                    </Box>
                    
                    <Stack spacing={3}>
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <TextInput
                                source="name"
                                readOnly={props.readonly}
                                validate={[required()]}
                                fullWidth
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
                                source="code"
                                readOnly={props.readonly}
                                validate={[required()]}
                                fullWidth
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
                        
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <AutocompleteInput
                                source="category"
                                readOnly={props.readonly}
                                filterSelectedOptions={false}
                                loading={isFilterLoading}
                                choices={categories.map(val => ({ id: val, name: val }))}
                                onCreate={val => {
                                    if (!val || categories.includes(val)) return;
                                    const newCategories = [...categories, val];
                                    setCategories(newCategories);
                                    return { id: val, name: val }
                                }}
                                fullWidth
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
                            <AutocompleteInput
                                source="sub_category"
                                readOnly={props.readonly}
                                loading={isFilterLoading}
                                choices={subCategories.map(val => ({ id: val, name: val }))}
                                onCreate={val => {
                                    if (!val || subCategories.includes(val)) return;
                                    const newSubCategories = [...subCategories, val];
                                    setSubCategories(newSubCategories);
                                    return { id: val, name: val }
                                }}
                                fullWidth
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
                            📊 Range & Units
                        </Typography>
                        <Chip 
                            label="Required" 
                            size="small" 
                            color="error" 
                            variant="outlined"
                            sx={{ ml: 'auto', fontSize: '0.75rem' }}
                        />
                    </Box>
                    
                    <Stack spacing={3}>
                        <Stack direction={"row"} gap={3} width={"100%"}>
                            <NumberInput
                                source="low_ref_range"
                                label="Low Range"
                                readOnly={props.readonly}
                                validate={[required()]}
                                fullWidth
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
                            <NumberInput
                                source="high_ref_range"
                                label="High Range"
                                readOnly={props.readonly}
                                validate={[required()]}
                                fullWidth
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
                            <AutocompleteInput
                                source="unit"
                                readOnly={props.readonly}
                                loading={isUnitLoading}
                                choices={[...new Set(unit)].map(val => ({ id: val, name: val }))}
                                onCreate={val => {
                                    if (!val || unit.includes(val)) return;
                                    const newUnit = [...new Set([...unit, val])];
                                    setUnit(newUnit);
                                    return { id: val, name: val }
                                }}
                                fullWidth
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
                            📋 Additional Settings
                        </Typography>
                        <Chip 
                            label="Required" 
                            size="small" 
                            color="error" 
                            variant="outlined"
                            sx={{ ml: 'auto', fontSize: '0.75rem' }}
                        />
                    </Box>
                    
                    <Stack spacing={3}>
                        <NumberInput
                            source="decimal"
                            readOnly={props.readonly}
                            validate={[required()]}
                            fullWidth
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
                        
                        <Box>
                            <Typography 
                                variant="body1" 
                                sx={{ 
                                    fontWeight: 500,
                                    color: theme.palette.text.secondary,
                                    mb: 2
                                }}
                            >
                                Specimen Types
                            </Typography>
                            <ArrayInput source="types">
                                <SimpleFormIterator inline>
                                    <FeatureList source="type" readOnly={props.readonly} types="specimen-type">
                                        <AutocompleteInput 
                                            source="type" 
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
                                </SimpleFormIterator>
                            </ArrayInput>
                        </Box>
                    </Stack>
                </CardContent>
            </Card>
        </Stack>
    )
}

function TestTypeForm(props: TestTypeFormProps) {
    return (
        <Box sx={{ p: { xs: 2, sm: 3 } }}>
            <SimpleForm
                sx={{
                    '& .RaSimpleForm-form': {
                        backgroundColor: 'transparent',
                        boxShadow: 'none',
                        padding: 0
                    }
                }}
            >
                <TestTypeInput {...props} />
            </SimpleForm>
        </Box>
    )
}

export function TestTypeEdit() {
    const theme = useTheme();
    
    return (
        <Box sx={{ 
            minHeight: '100vh', 
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Edit mutationMode="pessimistic" title="Edit Test Type" redirect={"list"}>
                <TestTypeForm readonly={false} mode={"EDIT"} />
                <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
                    <ReferenceSection />
                </Box>
            </Edit>
        </Box>
    )
}

export function TestTypeCreate() {
    const theme = useTheme();
    
    return (
        <Box sx={{ 
            minHeight: '100vh', 
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Create title="Create Test Type" redirect={"list"}>
                <TestTypeForm readonly={false} mode={"CREATE"} />
                <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
                    <ReferenceSection />
                </Box>
            </Create>
        </Box>
    )
}
