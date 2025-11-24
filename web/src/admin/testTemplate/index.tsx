import CloseIcon from '@mui/icons-material/Close';
import WarningIcon from '@mui/icons-material/WarningAmber';
import { Box, Button, Chip, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, useTheme, Card, CardContent } from "@mui/material";
import Stack from "@mui/material/Stack";
import { AxiosError } from "axios";
import { useEffect, useMemo, useState } from "react";
import { AutocompleteArrayInput, Create, Datagrid, DateField, DeleteButton, Edit, FilterLiveForm, FunctionField, List, NumberField, ReferenceInput, SaveButton, SearchInput, SimpleForm, TextInput, Toolbar, required, useEditContext, useNotify, useRedirect, useSaveContext } from "react-admin";
import { useFormContext } from "react-hook-form";
import SideFilter from "../../component/SideFilter";
import { useCurrentUser } from "../../hooks/currentUser";
import useAxios from "../../hooks/useAxios";
import type { ObservationRequest, ObservationRequestCreateRequest } from "../../types/observation_requests";
import type { ActionKeys } from "../../types/props";
import { RoleNameValue } from "../../types/role";
import { TestTemplate, TestTemplateDiff } from "../../types/test_templates";
import { WorkOrder } from "../../types/work_order";
import { TestInput, testTypesField } from '../workOrder/Form';

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

export const TestTemplateList = () => (
    <List aside={<TestTemplateFilterSidebar />} title="Test Template" sort={{
        field: "id",
        order: "DESC"
    }}
        storeKey={false} exporter={false}
    >
        <Datagrid bulkActionButtons={false} rowClick="edit">
            <NumberField source="id" />
            <FunctionField label="Name" render={(record: any) => <NullableField value={record.name} />} />
            <FunctionField label="Description" render={(record: any) => <NullableField value={record.description} />} />
            <DateField source="created_at" showTime />
            <DateField source="updated_at" showTime />
            <FunctionField
                label="Created By"
                render={(record: any) => <NullableField value={record.created_by_user?.fullname} />}
            />
            <FunctionField
                label="Last Updated By"
                render={(record: any) => <NullableField value={record.last_updated_by_user?.fullname} />}
            />
        </Datagrid>
    </List>
);

const TestTemplateFilterSidebar = () => {
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
                            🧪 Filter Test Templates
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
                </Stack>
            </FilterLiveForm>
        </SideFilter>
    )
};


type TestTemplateFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function TestTemplateForm(props: TestTemplateFormProps) {
    const theme = useTheme();
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
        <Box sx={{ ml: 3, mr: 3 }}>
            <SimpleForm
                disabled={props.readonly}
                toolbar={false}
                sx={{
                    '& .RaSimpleForm-form': {
                        backgroundColor: 'transparent',
                        boxShadow: 'none',
                        padding: 0
                    }
                }}
            >
                <TestTypeToolbar />


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

                            <Stack>
                                <TextInput
                                    source="name"
                                    readOnly={props.readonly}
                                    validate={[required()]}
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
                                    source="description"
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
                                    👨‍⚕️ Assignments
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
                                <ReferenceInput
                                    source={"doctor_ids"}
                                    reference="user"
                                    resource='user'
                                    target="id"
                                    label="Doctor"
                                    filter={{
                                        role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
                                    }}
                                >
                                    <AutocompleteArrayInput
                                        suggestionLimit={10}
                                        filterToQuery={(searchText) => ({
                                            q: searchText,
                                            role: [RoleNameValue.DOCTOR, RoleNameValue.ADMIN]
                                        })}
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
                                </ReferenceInput>

                                <ReferenceInput
                                    source={"analyzers_ids"}
                                    reference="user"
                                    resource='user'
                                    target="id"
                                    label="Analyst"
                                    filter={{}}
                                >
                                    <AutocompleteArrayInput
                                        suggestionLimit={10}
                                        filterToQuery={(searchText) => ({
                                            q: searchText,
                                        })}
                                        defaultValue={[currentUser?.id]}
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
                                </ReferenceInput>
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
                                    🧪 Test Configuration
                                </Typography>
                                <Chip
                                    label="Required"
                                    size="small"
                                    color="error"
                                    variant="outlined"
                                    sx={{ ml: 'auto', fontSize: '0.75rem' }}
                                />
                            </Box>

                            <TestInput initSelectedType={selectedType} />
                        </CardContent>
                    </Card>
                </Stack>
            </SimpleForm>
        </Box>
    )
}

const TestTemplateSaveButton = ({ disabled }: { disabled?: boolean }) => {
    const [open, setOpen] = useState(false);
    const [diffData, setDiffData] = useState<TestTemplateDiff | null>(null)

    const { getValues } = useFormContext();
    const { save } = useSaveContext();
    const axios = useAxios();
    const notify = useNotify();
    const handleClick = async (e: any) => {
        e.preventDefault(); // necessary to prevent default SaveButton submit logic
        await handleSave();
    };
    const redirect = useRedirect()
    const buildPayload = () => {
        const data = getValues() as TestTemplate;
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

        // @ts-ignore
        const observationRequest = data[testTypesField] as Record<number, ObservationRequestCreateRequest>
        const payload = {
            ...data,
            test_types: Object.entries(observationRequest).map(([_, value]) => {
                return value
            })
        }

        return payload
    }

    const handleSave = async () => {
        const payload = buildPayload();
        if (!payload) {
            return;
        }

        try {
            const isEditMode = !!payload.id;
            const url = isEditMode
                ? `/test-template/${payload.id}/update-diff`
                : `/test-template`;
            const method = isEditMode ? 'put' : 'post';

            const response = await axios({
                method,
                url,
                data: payload,
                headers: {
                    "Content-Type": "application/json",
                },
            });

            const respDiff = response.data as TestTemplateDiff;

            if (respDiff.ToCreate?.length > 0 || respDiff.ToDelete?.length > 0) {
                setOpen(true);
                setDiffData(respDiff);
                return;
            }

            if (isEditMode) {
                submitTestTemplate();
            } else {
                redirect("/test-template");
            }
        } catch (error) {
            const message = error instanceof AxiosError
                ? error.response?.data?.error || "Something went wrong"
                : "An unexpected error occurred";

            notify(message, { type: "error" });
            return;
        }
    }

    const submitTestTemplate = () => {
        const payload = buildPayload();
        if (!payload) {
            return;
        }

        if (save) {
            save(payload);
        }
    }

    return <>
        <SaveButton type="button" onClick={handleClick} alwaysEnable size="small" />
        <ConfirmTemplateModificationModal
            isOpen={open}
            onClose={() => setOpen(false)}
            onConfirm={submitTestTemplate}
            testTemplateDiff={diffData}
        />
    </>
}

const TestTypeToolbar = () => {
    const theme = useTheme();
    const isDarkMode = theme.palette.mode === 'dark';

    return (
        <Stack width={'100%'}
            sx={{
                position: "sticky",
                top: 80,
                display: 'flex',
                zIndex: 2147483647,
                marginBottom: 1,
                '& .MuiToolbar-root': {
                    backgroundColor: isDarkMode ? '#211e2c' : 'white',
                    color: isDarkMode ? theme.palette.text.primary : theme.palette.text.primary,
                    boxShadow: isDarkMode ? theme.shadows[2] : '0 2px 4px rgba(0,0,0,0.1)',
                },
                overflow: 'hidden'
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
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Edit mutationMode="pessimistic" title="Edit Test Template" sx={{
                "& .RaEdit-card": {
                    overflow: "visible",
                }
            }} emptyWhileLoading>
                <TestTemplateForm readonly={false} mode={"EDIT"} />
            </Edit>
        </Box>
    )
}

export function TestTemplateCreate() {
    const theme = useTheme();

    return (
        <Box sx={{
            minHeight: '100vh',
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Create title="Create Test Template" redirect={"show"} sx={{
                "& .RaCreate-card": {
                    overflow: "visible",
                }
            }}>
                <TestTemplateForm readonly={false} mode={"CREATE"} />
            </Create>
        </Box>
    )
}

const defaultGetWorkOrderName = (workOrder: WorkOrder): string => {
    if (workOrder.patient && (workOrder.patient.first_name || workOrder.patient.last_name)) {
        return `${workOrder.patient.first_name} ${workOrder.patient.last_name}`;
    }
    return `Work Order ID: ${workOrder.id}`;
};

interface ConfirmTemplateModificationModalProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: () => void;
    testTemplateDiff: TestTemplateDiff | null;
    getWorkOrderName?: (workOrder: WorkOrder) => string; // Optional custom name getter
}

interface ProcessedWorkOrderView {
    id: number; // Work Order ID
    name: string; // Work Order Name
    barcode: string; // Work Order Barcode
    changes: AffectedChange[];
}

interface AffectedChange {
    test_code: string;
    test_description: string;
    change_type: 'Added' | 'Removed';
}

const ConfirmTemplateModificationModal: React.FC<ConfirmTemplateModificationModalProps> = ({
    isOpen,
    onClose,
    onConfirm,
    testTemplateDiff,
    getWorkOrderName = defaultGetWorkOrderName,
}) => {
    // Memoize processed data to avoid re-computation on every render unless dependencies change
    const processedData = useMemo(() => {
        if (!testTemplateDiff) {
            return [];
        }

        const affectedWorkOrdersMap: Map<number, {
            id: number;
            barcode: string;
            originalWorkOrder: WorkOrder;
            changes: AffectedChange[];
        }> = new Map();

        // Helper function to process a list of observation requests
        const processObservationList = (
            list: ObservationRequest[],
            change_type: 'Added' | 'Removed'
        ) => {
            list?.forEach(obsReq => {
                if (obsReq.work_order) {
                    const wo = obsReq.work_order;
                    if (!affectedWorkOrdersMap.has(wo.id)) {
                        affectedWorkOrdersMap.set(wo.id, {
                            id: wo.id,
                            barcode: wo.barcode,
                            originalWorkOrder: wo,
                            changes: [],
                        });
                    }
                    affectedWorkOrdersMap.get(wo.id)!.changes.push({
                        test_code: obsReq.test_code,
                        test_description: obsReq.test_description,
                        change_type: change_type,
                    });
                }
            });
        };

        processObservationList(testTemplateDiff.ToDelete || [], 'Removed');
        processObservationList(testTemplateDiff.ToCreate || [], 'Added');

        const result: ProcessedWorkOrderView[] = [];
        affectedWorkOrdersMap?.forEach((value) => {
            result.push({
                id: value.id,
                name: getWorkOrderName(value.originalWorkOrder),
                barcode: value.barcode,
                changes: value.changes.sort((a, b) => a.test_code.localeCompare(b.test_code)),
            });
        });

        return result.sort((a, b) => a.id - b.id);
    }, [testTemplateDiff, getWorkOrderName]);

    if (!isOpen || !testTemplateDiff) {
        return null;
    }

    const hasAffectedWorkOrders = processedData.length > 0;
    const hasAnyChanges = (testTemplateDiff?.ToCreate?.length || 0) > 0 || (testTemplateDiff?.ToDelete?.length || 0) > 0;

    return (
        <Dialog open={isOpen} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle sx={{ display: 'flex', alignItems: 'center', bgcolor: 'warning.light', color: 'warning.contrastText' }}>
                <WarningIcon sx={{ mr: 1 }} />
                Confirm Test Template Modification
                <IconButton
                    aria-label="close"
                    onClick={onClose}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <CloseIcon />
                </IconButton>
            </DialogTitle>
            <DialogContent dividers>
                <DialogContentText component="div" sx={{ mb: 2 }}>
                    <Typography variant="body1" gutterBottom>
                        Modifying this test template will apply the changes to its definition.
                        {hasAffectedWorkOrders && " Additionally, it will directly alter the tests scheduled for the following current work orders:"}
                        {!hasAffectedWorkOrders && hasAnyChanges && " While no current work orders will have tests added or removed, the template changes will apply to new work orders."}
                        {!hasAffectedWorkOrders && !hasAnyChanges && " There are no changes to be made to the template or any linked work orders."}
                    </Typography>
                    {hasAnyChanges && (
                        <Box sx={{ p: 2, bgcolor: 'warning.lighter', borderRadius: 1, display: 'flex', alignItems: 'center', color: 'warning.darker' }}>
                            <WarningIcon fontSize="small" sx={{ mr: 1 }} />
                            <Typography variant="body2">
                                These changes may impact ongoing processes and data. This action may not be easily undone.
                            </Typography>
                        </Box>
                    )}
                </DialogContentText>

                {hasAffectedWorkOrders ? (
                    processedData.map(wo => (
                        <Paper elevation={2} sx={{ mb: 2, p: 2 }} key={wo.id}>
                            <Typography variant="h6" gutterBottom component="div">
                                Lab Request: {wo.name}
                            </Typography>
                            <Typography variant="body2" color="textSecondary" gutterBottom>
                                ID: {wo.id} | Barcode: {wo.barcode}
                            </Typography>

                            {(wo.changes?.length || 0) > 0 ? (
                                <TableContainer component={Paper} variant="outlined">
                                    <Table size="small" aria-label={`Changes for work order ${wo.id}`}>
                                        <TableHead sx={{ bgcolor: 'grey.200' }}>
                                            <TableRow>
                                                <TableCell>Test Code</TableCell>
                                                <TableCell>Description</TableCell>
                                                <TableCell>Change</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {wo.changes.map((change, index) => (
                                                <TableRow key={`${change.test_code}-${index}`}>
                                                    <TableCell component="th" scope="row">
                                                        {change.test_code}
                                                    </TableCell>
                                                    <TableCell>{change.test_description}</TableCell>
                                                    <TableCell>
                                                        <Chip
                                                            label={change.change_type}
                                                            color={change.change_type === 'Added' ? 'success' : 'error'}
                                                            size="small"
                                                            variant="outlined"
                                                        />
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            ) : (
                                <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'text.secondary' }}>
                                    No specific test changes identified for this work order based on the template modification.
                                </Typography>
                            )}
                        </Paper>
                    ))
                ) : (
                    hasAnyChanges && (
                        <Typography variant="body1" sx={{ mt: 2 }}>
                            No existing work orders have been identified that will have tests added or removed as a direct result of these specific template changes.
                            The modification will apply to <strong>new work orders</strong> created using this template.
                        </Typography>
                    )
                )}
                {!hasAnyChanges && (
                    <Box sx={{ p: 2, bgcolor: 'info.lighter', borderRadius: 1, display: 'flex', alignItems: 'center', color: 'info.darker', mt: 2 }}>
                        <Typography variant="body2">
                            No changes (additions or deletions of tests) are specified in the template modification.
                        </Typography>
                    </Box>
                )}
            </DialogContent>
            <DialogActions sx={{ p: '16px 24px' }}>
                <Button onClick={onClose} color="inherit">
                    Cancel
                </Button>
                <Button
                    onClick={onConfirm}
                    variant="contained"
                    color="warning"
                    disabled={!hasAnyChanges} // Disable confirm if there are no changes to apply
                >
                    Confirm Modification
                </Button>
            </DialogActions>
        </Dialog>
    );
};