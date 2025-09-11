import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import CloseIcon from '@mui/icons-material/Close';
import MedicalServicesIcon from '@mui/icons-material/MedicalServices';
import ScienceIcon from '@mui/icons-material/Science';
import { Button, Grid, IconButton, ListItem, ListItemIcon, ListItemText, Modal, List as MuiList, Paper, Stack, Typography, useTheme, Card, CardContent, Chip } from '@mui/material';
import Box from "@mui/material/Box";
import { useQuery } from '@tanstack/react-query';
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
    FilterLiveForm,
    List,
    PasswordInput,
    ReferenceArrayInput,
    required,
    SearchInput,
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
import useAxios from '../../hooks/useAxios.ts';
import { useRefererRedirect } from "../../hooks/useReferer.ts";
import { Action, ActionKeys } from "../../types/props.ts";
import { Role, RoleName, RoleNameValue } from '../../types/role.ts';

export type UserFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

export function UserFormField(props: UserFormProps) {
    const theme = useTheme();
    const axios = useAxios();
    const { data: roleData } = useQuery<Role[]>({
        queryKey: ['roles'],
        queryFn: async () => {
            const response = await axios.get('/role');
            if (response.status !== 200) {
                throw new Error('Failed to fetch roles');
            }
            return response.data;
        }
    });

    const [openModal, setOpenModal] = useState(false);
    const handleOpenModal = () => setOpenModal(true);
    const handleCloseModal = () => setOpenModal(false);

    return (
        <Box sx={{ p: { xs: 2, sm: 3, width: '100%' } }}>
            <Stack spacing={3} sx={{ width: '100%' }}>
                {props.mode !== Action.CREATE && (
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
                                    ℹ️ System Information
                                </Typography>
                                <Chip 
                                    label="Read Only" 
                                    size="small" 
                                    color="default" 
                                    variant="outlined"
                                    sx={{ ml: 'auto', fontSize: '0.75rem' }}
                                />
                            </Box>
                            
                            <Stack>
                                <TextInput 
                                    source={"id"} 
                                    readOnly={true}
                                    sx={{
                                        '& .MuiOutlinedInput-root': {
                                            borderRadius: 2,
                                            transition: 'all 0.2s ease'
                                        }
                                    }}
                                />
                                <Stack direction={"row"} gap={2} width={"100%"}>
                                    <DateTimeInput 
                                        source={"created_at"} 
                                        readOnly={true}
                                        sx={{
                                            flex: 1,
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
                                            flex: 1,
                                            '& .MuiOutlinedInput-root': {
                                                borderRadius: 2,
                                                transition: 'all 0.2s ease'
                                            }
                                        }}
                                    />
                                </Stack>
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
                                👤 Basic Information
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
                                source="username" 
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
                                source="fullname" 
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
                                source="email" 
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
                                🔐 Security & Access
                            </Typography>
                            <Chip 
                                label={props.mode === Action.EDIT ? "Optional" : "Required"} 
                                size="small" 
                                color={props.mode === Action.EDIT ? "default" : "error"} 
                                variant="outlined"
                                sx={{ ml: 'auto', fontSize: '0.75rem' }}
                            />
                        </Box>
                        
                        <Stack>
                            {props.mode !== Action.SHOW && (
                                <PasswordInput 
                                    source="password" 
                                    validate={props.mode === Action.EDIT ? [] : [required()]} 
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
                            )}
                            <BooleanInput 
                                source="is_active" 
                                defaultValue={true} 
                                readOnly={props.readonly}
                                sx={{
                                    '& .MuiFormControlLabel-root': {
                                        ml: 0,
                                        '& .MuiSwitch-root': {
                                            mr: 1
                                        }
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
                                🎭 Role Management
                            </Typography>
                            <Chip 
                                label="Required" 
                                size="small" 
                                color="error" 
                                variant="outlined"
                                sx={{ ml: 'auto', fontSize: '0.75rem' }}
                            />
                        </Box>
                        
                        <Grid sx={{ width: '100%' }} >
                            <Grid size={8}>
                                <ReferenceArrayInput source="roles_id" reference="role">
                                    <SelectArrayInput 
                                        optionText="name" 
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
                                </ReferenceArrayInput>
                            </Grid>
                            <Grid >
                                <Button 
                                    variant="outlined" 
                                    onClick={handleOpenModal} 
                                    sx={{ 
                                        borderRadius: 2,
                                        textTransform: 'none',
                                        fontWeight: 500
                                    }} 
                                    size="small"
                                    
                                >
                                    View Roles
                                </Button>
                            </Grid>
                        </Grid>
                    </CardContent>
                </Card>
            </Stack>

            <Modal
                open={openModal}
                onClose={handleCloseModal}
                aria-labelledby="role-list-modal-title"
                aria-describedby="role-list-modal-description"
            >
                <Paper sx={{
                    position: 'absolute' as const,
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)',
                    width: { xs: '90%', sm: '75%', md: '60%' },
                    maxWidth: 700,
                    bgcolor: 'background.paper',
                    borderRadius: 2,
                    boxShadow: 24,
                    p: { xs: 2, sm: 3, md: 4 },
                    outline: 'none',
                    maxHeight: '90vh',
                    overflowY: 'auto',
                }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Typography id="role-list-modal-title" variant="h5" component="h2" sx={{ fontWeight: 'bold' }}>
                            Available Roles
                        </Typography>
                        <IconButton onClick={handleCloseModal} aria-label="close modal">
                            <CloseIcon />
                        </IconButton>
                    </Box>
                    <RolePresentationList roles={roleData ?? []} />
                    <Typography id="role-list-modal-description" sx={{ display: 'none' }}>
                        A list of roles with their descriptions.
                    </Typography>
                </Paper>
            </Modal>
        </Box>
    )
}


const getRoleIcon = (roleName: RoleName) => {
    switch (roleName) {
        case RoleNameValue.ADMIN:
            return <AdminPanelSettingsIcon fontSize="medium" />;
        case RoleNameValue.ANALYZER:
            return <ScienceIcon fontSize="medium" />;
        case RoleNameValue.DOCTOR:
            return <MedicalServicesIcon fontSize="medium" />;
        default:
            return <AssignmentIndIcon fontSize="medium" />;
    }
};


interface RolePresentationListProps {
    roles: Role[];
}

const RolePresentationList: React.FC<RolePresentationListProps> = ({ roles }) => {
    if (!roles || roles.length === 0) {
        return <Typography sx={{ textAlign: 'center', p: 2 }}>No roles to display.</Typography>;
    }

    return (

        // Removed margin: '20px auto' from MuiList as it will be controlled by Modal
        <MuiList sx={{
            display: 'flex',
            flexDirection: 'column',
            gap: 2,
            p: { xs: 1, sm: 2 },
            listStyle: 'none',
            width: '100%',
            // maxWidth is good to keep for content within modal
            maxWidth: { xs: '100%', sm: 500, md: 600 },
            backgroundColor: 'background.paper',
            borderRadius: 2,
            // boxShadow: '0 2px 10px rgba(0,0,0,0.05)', // Shadow is now on the Modal's Paper
        }}>
            {roles.map((role) => {
                const isPrivilegedRole = ['Admin', 'Doctor'].includes(role.name);
                const backgroundColor = {
                    Admin: 'primary.main',
                    Doctor: 'warning.main',
                    Analyzer: 'info.main',
                }[role.name] || 'grey.100';
                const textColor = isPrivilegedRole ? 'primary.contrastText' : 'text.primary';
                const iconColor = isPrivilegedRole ? 'primary.contrastText' : 'action.active';

                return (
                    <ListItem
                        key={role.id}
                        sx={{
                            bgcolor: backgroundColor,
                            color: textColor,
                            borderRadius: 2,
                            boxShadow: 2,
                            transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
                            '&:hover': {
                                transform: 'translateY(-3px) scale(1.01)',
                                boxShadow: 5,
                            },
                            p: 0,
                        }}
                    >
                        <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', p: 2, gap: 2 }}>
                            <ListItemIcon sx={{
                                minWidth: 'auto',
                                color: iconColor,
                            }}>
                                {getRoleIcon(role.name)}
                            </ListItemIcon>
                            <ListItemText
                                primary={role.name}
                                secondary={role.description}
                                primaryTypographyProps={{
                                    fontWeight: 'bold',
                                    color: 'inherit',
                                    component: 'h3',
                                    variant: 'subtitle1'
                                }}
                                secondaryTypographyProps={{
                                    color: 'inherit',
                                    // @ts-ignore
                                    opacity: isPrivilegedRole ? 0.85 : 0.75,
                                    variant: 'body2',
                                    lineHeight: 1.4
                                }}
                            />
                        </Box>
                    </ListItem>
                );
            })}
        </MuiList>
    );
};

export function UserForm(props: UserFormProps) {
    return (
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
            <UserFormField {...props} />
        </SimpleForm>
    )
}

export function UserCreate() {
    const theme = useTheme();
    const redirect = useRefererRedirect("show");

    return (
        <Box sx={{ 
            minHeight: '100vh', 
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Create redirect={redirect} resource="user" mutationMode='pessimistic'>
                <UserForm mode={"CREATE"} />
            </Create>
        </Box>
    )
}

export function UserShow() {
    const theme = useTheme();
    
    return (
        <Box sx={{ 
            minHeight: '100vh', 
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Show resource="user">
                <UserForm readonly mode={"SHOW"} />
            </Show>
        </Box>
    )
}

export function UserEdit() {
    const theme = useTheme();
    
    return (
        <Box sx={{ 
            minHeight: '100vh', 
            bgcolor: theme.palette.background.default,
            pb: 4
        }}>
            <Edit resource="user" mutationMode='pessimistic'>
                <UserForm mode={"EDIT"} />
            </Edit>
        </Box>
    )
}

const UserFilterSidebar = () => {
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
                            👥 Filter Users
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
            <ArrayField source="roles" sortable={false}>
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
        case 'Analyzer':
            return 'error';
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
