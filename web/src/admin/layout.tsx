import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import SettingsIcon from '@mui/icons-material/Settings';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import LogoutIcon from '@mui/icons-material/Logout';
import AssessmentIcon from '@mui/icons-material/Assessment';
import BiotechIcon from '@mui/icons-material/Biotech';
import ApprovalIcon from '@mui/icons-material/Approval';
import BuildIcon from '@mui/icons-material/Build';
import LanIcon from '@mui/icons-material/Lan';
import PersonIcon from '@mui/icons-material/Person';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import TableViewIcon from '@mui/icons-material/TableView';
import DashboardIcon from '@mui/icons-material/Dashboard';
import InfoIcon from '@mui/icons-material/Info';
import { Box, Stack, Tooltip, Typography, Avatar, Menu, MenuItem, ListItemIcon, ListItemText } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useEffect, useState, type ReactNode } from 'react';
import React from 'react';
import { AppBar, Button, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, useTheme, useLogout } from 'react-admin';
import { useLocation, useNavigate } from "react-router-dom";
import AppIndicator from '../component/AppIndicator';
import Breadcrumbs, { type BreadcrumbsLink } from '../component/Breadcrumbs';
import { toTitleCase } from '../helper/format';
import FileOpenIcon from '@mui/icons-material/FileOpen';
import logo from '../assets/alinda-husada-logo.png';
import { useCurrentUser, useCurrentUserRole } from '../hooks/currentUser';
import logo from '../assets/elgatama-logo.png';
import { useCurrentUser } from '../hooks/currentUser';


const SettingsButton = () => (
    <Link to="/settings" color={"inherit"}>
        <IconButton style={{ color: '#555555' }}>
            <SettingsIcon sx={{ width: 30, height: 'auto' }} />
        </IconButton>
    </Link>
);

const LogButton = () => (
    <Tooltip title="Logs">
        <Link to="/logs" color={"inherit"}>
            <IconButton style={{ color: '#1E88E5' }}>
                <FileOpenIcon sx={{ width: 30, height: 'auto' }} />
            </IconButton>
        </Link>
    </Tooltip>
);

const CustomToggleThemeButton = () => {
    const [theme, setTheme] = useTheme();

    // Debug: Log tema saat ini setiap kali berubah
    useEffect(() => {
        console.log(`Current theme mode: ${theme}`);
        console.log(`Is dark mode: ${theme === 'dark'}`);
        console.log(`Is light mode: ${theme === 'light'}`);
    }, [theme]);

    return (
        <Tooltip title={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}>
            <IconButton
                onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
                style={{
                    color: theme === 'light' ? '#FFA726' : '#FFD54F'
                }}
            >
                {theme === 'light' ? <LightModeIcon sx={{ width: 30, height: 'auto' }} /> : <DarkModeIcon sx={{ width: 30, height: 'auto' }} />}
            </IconButton>
        </Tooltip>
    );
};

const UserProfile = () => {
    const currentUser = useCurrentUser();
    const currentUserRole = useCurrentUserRole();
    const logout = useLogout();
    const [theme] = useTheme();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    if (!currentUser) {
        return null;
    }

    const isDarkMode = theme === 'dark';

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = () => {
        logout();
        handleClose();
    };

    let avatarColor;

    switch (currentUserRole) {
        case "Analyzer":
            avatarColor = "#d9db3aff";
            break;
        case "Doctor":
            avatarColor = "#2196F3";
            break;
        default:
            avatarColor = "#4abaab";
    }

    return (
        <>
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1,
                    ml: 1,
                    cursor: 'pointer',
                    padding: '4px 8px',
                    borderRadius: 1,
                    '&:hover': {
                        backgroundColor: 'rgba(74, 186, 171, 0.08)',
                    }
                }}
                onClick={handleClick}
            >
                <Avatar
                    sx={{
                        width: 36,
                        height: 36,
                        bgcolor: avatarColor,
                        fontSize: '14px'
                    }}
                >
                    {currentUser.fullname?.charAt(0)?.toUpperCase() || <AccountCircleIcon />}
                </Avatar>
                <Typography
                    sx={{
                        color: isDarkMode ? '#ffffff' : '#1d293d',
                        fontWeight: 500,
                        display: { xs: 'none', md: 'block' }
                    }}
                >
                    {currentUser.fullname}
                </Typography>
            </Box>

            <Menu
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                onClick={handleClose}
                PaperProps={{
                    elevation: 3,
                    sx: {
                        mt: 1.5,
                        minWidth: 200,
                        '& .MuiAvatar-root': {
                            width: 24,
                            height: 24,
                            ml: -0.5,
                            mr: 1,
                        },
                    },
                }}
                transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
                <MenuItem onClick={handleLogout}>
                    <ListItemIcon>
                        <LogoutIcon fontSize="small" sx={{ color: '#f44336' }} />
                    </ListItemIcon>
                    <ListItemText>Logout</ListItemText>
                </MenuItem>
            </Menu>
        </>
    );
};

const CompanyLogo = () => (
    <Box sx={{
        display: 'flex',
        alignItems: 'center',
        gap: 1,
        mr: 2,
        '&:hover': {
            opacity: 0.8,
            cursor: 'pointer'
        }
    }}>
        <img
            src={logo}
            alt="Elga Tama Logo"
            style={{
                height: '40px',
                width: 'auto'
            }}
        />
        <Typography variant="h6">PT ELGA TAMA</Typography>
    </Box>
);


const MyAppBar = () => {
    const currentUserRole = useCurrentUserRole();
    const location = useLocation();
    const [theme] = useTheme();

    const appBarAlwaysOn = () => {
        if (location.pathname.includes('/work-order/')) return true;
        if (location.pathname.includes('/test-template/')) return true;
        return false;
    };

    const isDarkMode = theme === 'dark';

    return (
        <AppBar
            userMenu={false}
            color="primary"
            sx={{
                position: 'fixed',
                color: isDarkMode ? '#ffffff' : '#1d293d',
                backgroundColor: isDarkMode ? '#151221' : 'white',
                height: 80,
                justifyContent: 'center',
                boxShadow: 'rgba(0, 0, 0, 0.15) 1.95px 1.95px 2.6px',
                '& .RaAppBar-title': {
                    display: 'none',
                },
                '& .MuiIconButton-root[aria-label="Open drawer"]': {
                    display: 'none !important',
                },
            }}
            alwaysOn={appBarAlwaysOn()}
            toolbar={
                <>
                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 4
                    }}>
                        <Box sx={{
                            display: 'flex',
                            alignItems: 'center'
                        }}>
                            {currentUserRole === "Admin" &&
                                <>
                                    <SettingsButton />
                                    <LogButton />
                                </>
                            }
                            <AppIndicator />
                            <CustomToggleThemeButton />
                            <LoadingIndicator sx={{ scale: 1.2 }} />
                        </Box>
                        <UserProfile />
                    </Box>
                </>
            }
        ><CompanyLogo /></AppBar>
    );
};

const PageTitle = () => {
    const location = useLocation();

    const getPageData = (pathname: string) => {
        const pathParts = pathname.substring(1).split('/');
        const mainPath = pathParts[0];

        const resourceData: { [key: string]: { title: string; icon: React.ReactElement } } = {
            'work-order': { title: 'Lab Request', icon: <BiotechIcon /> },
            'result': { title: 'Result', icon: <AssessmentIcon /> },
            'approval': { title: 'Approval', icon: <ApprovalIcon /> },
            'patient': { title: 'Patients', icon: <PersonIcon /> },
            'test-type': { title: 'Test Type', icon: <BiotechIcon /> },
            'test-template': { title: 'Test Template', icon: <TableViewIcon /> },
            'device': { title: 'Devices', icon: <LanIcon /> },
            'user': { title: 'Users', icon: <AdminPanelSettingsIcon /> },
            'config': { title: 'Configuration', icon: <BuildIcon /> },
            'settings': { title: 'Settings', icon: <SettingsIcon /> },
            'logs': { title: 'Logs', icon: <FileOpenIcon /> },
            'about': { title: 'About Us', icon: <InfoIcon /> }
        };

        const defaultData = { title: 'Dashboard', icon: <DashboardIcon /> };
        const resourceInfo = resourceData[mainPath] || { title: toTitleCase(mainPath), icon: <DashboardIcon /> };

        if (pathParts.length > 1) {
            const action = pathParts[1];

            if (action === 'create') return { title: `Create ${resourceInfo.title}`, icon: resourceInfo.icon };
            if (action === 'edit') return { title: `Edit ${resourceInfo.title}`, icon: resourceInfo.icon };
            if (action === 'show') return { title: `View ${resourceInfo.title}`, icon: resourceInfo.icon };
            if (pathParts[2] === 'show') return { title: `View ${resourceInfo.title} #${pathParts[1]}`, icon: resourceInfo.icon };
            if (pathParts[2] === 'edit') return { title: `Edit ${resourceInfo.title} #${pathParts[1]}`, icon: resourceInfo.icon };
        }

        return mainPath ? resourceInfo : defaultData;
    };

    const pageData = getPageData(location.pathname);

    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            mb: 1
        }}>
            <Box sx={{
                color: 'primary.main',
                display: 'flex',
                alignItems: 'center'
            }}>
                {React.cloneElement(pageData.icon, { sx: { fontSize: 32 } })}
            </Box>
            <Typography variant='h5' sx={{
                fontWeight: 600,
                color: 'text.primary'
            }}>
                {pageData.title.toUpperCase()}
            </Typography>
        </Box>
    );
};

type PathConfiguration = {
    path: string
    type: string
}

const DynamicBreadcrumbs = () => {
    const location = useLocation();
    const [breadCrumbsLinks, setBreadCrumbsLinks] = useState<Array<BreadcrumbsLink>>([])

    useEffect(() => {
        const pathSplit = location.pathname.substring(1).split("/")

        const pathsConfig: Array<PathConfiguration> = []
        pathSplit.forEach((val: string) => {
            pathsConfig.push({
                path: val,
                type: "general",
            })
        })


        const currentBreadCrumb: Array<BreadcrumbsLink> = []
        pathsConfig.forEach((val, i) => {
            const generateHref = (): string => {
                const pathUntil = pathsConfig.slice(0, i + 1)
                return "/" + pathUntil.map(v => v.path).filter(v => v).join("/") + location.search
            }

            currentBreadCrumb.push({
                label: toTitleCase(val.path),
                href: generateHref(),
                // icon: val.type == "general" ? <Icon /> : <Icon />,
                active: i == pathsConfig.length - 1,
            })
        })

        console.log(currentBreadCrumb)

        setBreadCrumbsLinks(currentBreadCrumb)
    }, [location.pathname])

    const navigate = useNavigate();

    return (
        <Stack direction={"row"} sx={{ marginLeft: 3 }}>
            <Button label='Back' variant='contained' onClick={() => navigate(-1)} sx={{
                display: location.pathname.split("/").length > 2 ? 'flex' : 'none',
                my: 1.25,
            }}>
                <ArrowBackIcon />
            </Button>
            <Breadcrumbs links={breadCrumbsLinks}>
            </Breadcrumbs>
        </Stack>
    )
}

const Footer = () => {
    const [theme] = useTheme();
    const isDarkMode = theme === 'dark';

    return (
        <Box
            component="footer"
            sx={{
                position: 'absolute',
                bottom: 0,
                left: 0,
                right: 0,
                height: 50,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                zIndex: 1000,
            }}
        >
            <Typography
                variant="body2"
                sx={{
                    color: isDarkMode ? '#ffffff' : 'text.secondary'
                }}
            >
                © {new Date().getFullYear()} PT ELGA TAMA. All rights reserved.
            </Typography>
        </Box>
    );
};

export const DefaultLayout = ({ children }: { children: ReactNode }) => {
    return (
        <Layout sx={{}} appBar={MyAppBar}>
            <Box sx={{
                marginTop: 7,
                marginLeft: 3,
                marginBottom: 1,
            }}>
                <PageTitle />
            </Box>
            <Stack direction={"row"} gap={2}>
                <DynamicBreadcrumbs />
            </Stack>
            <Box sx={{ paddingLeft: 3, paddingRight: 3, paddingBottom: 8 }}>
                {children}
            </Box>
            <CheckForApplicationUpdate />
            <Footer />
        </Layout>
    )
};
