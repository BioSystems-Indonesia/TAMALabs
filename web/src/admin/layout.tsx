import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import SettingsIcon from '@mui/icons-material/Settings';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import LogoutIcon from '@mui/icons-material/Logout';
import { Box, Stack, Tooltip, Typography, Avatar, Menu, MenuItem, ListItemIcon, ListItemText } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useEffect, useState, type ReactNode } from 'react';
import { AppBar, Button, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, useTheme, useLogout } from 'react-admin';
import { useLocation, useNavigate } from "react-router-dom";
import AppIndicator from '../component/AppIndicator';
import Breadcrumbs, { type BreadcrumbsLink } from '../component/Breadcrumbs';
import { toTitleCase } from '../helper/format';
import FileOpenIcon from '@mui/icons-material/FileOpen';
import logo from '../assets/elgatama-logo.png';
import { useCurrentUser } from '../hooks/currentUser';


const SettingsButton = () => (
    <Link to="/settings" color={"inherit"}>
        <IconButton style={{color: '#555555'}}>
            <SettingsIcon sx={{width: 30, height:'auto'}}/>
        </IconButton>
    </Link>
);

const LogButton = () => (
    <Tooltip title="Logs">
    <Link to="/logs" color={"inherit"}>
        <IconButton style={{color:'#1E88E5'}}>
            <FileOpenIcon  sx={{width: 30, height:'auto'}} />
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
                {theme === 'light' ? <LightModeIcon sx={{width: 30, height:'auto'}} /> : <DarkModeIcon sx={{width: 30, height:'auto'}} />}
            </IconButton>
        </Tooltip>
    );
};

const UserProfile = () => {
    const currentUser = useCurrentUser();
    const logout = useLogout();
    const [theme] = useTheme(); // Deteksi tema saat ini
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    
    if (!currentUser) {
        return null;
    }

    // Contoh penggunaan deteksi tema
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
                        bgcolor: '#4abaab',
                        fontSize: '14px'
                    }}
                >
                    {currentUser.fullname?.charAt(0)?.toUpperCase() || <AccountCircleIcon />}
                </Avatar>
                <Typography 
                    sx={{ 
                        color: isDarkMode ? '#ffffff' : '#1d293d', // Warna berubah berdasarkan tema
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
        marginLeft: 21,
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
    const location = useLocation();
    const [theme] = useTheme(); // Deteksi tema di AppBar
    
    const appBarAlwaysOn = () => {
        if (location.pathname.includes('/work-order/')) return true;
        if (location.pathname.includes('/test-template/')) return true;
        return false;
    };

    // Variabel untuk deteksi mode
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
                    <CompanyLogo />
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
                            <SettingsButton />
                            <LogButton />
                            <AppIndicator />
                            <CustomToggleThemeButton />
                            <LoadingIndicator sx={{scale: 1.2}} />
                        </Box>
                        <UserProfile />
                    </Box>
                </>
            }
        />
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
        <Stack direction={"row"} sx={{marginTop: 5}}>
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
    const [theme] = useTheme(); // Deteksi tema di Footer
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
                    color: isDarkMode ? '#ffffff' : 'text.secondary' // Warna text berubah berdasarkan tema
                }}
            >
                © 2025 PT ELGA TAMA. All rights reserved.
            </Typography>
        </Box>
    );
};

export const DefaultLayout = ({ children }: { children: ReactNode }) => {
    return (
        <Layout sx={{}} appBar={MyAppBar}>
            <Stack direction={"row"} gap={2}>
                <DynamicBreadcrumbs />
            </Stack>
            <Box sx={{ paddingBottom: '60px' }}>
                {children}
            </Box>
            <Footer />
            <CheckForApplicationUpdate />
        </Layout>
    )
};
