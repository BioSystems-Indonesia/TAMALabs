import { alpha, createTheme, PaletteOptions, Theme } from '@mui/material';
import { ThemeOptions as MuiThemeOptions } from '@mui/material';

export type ComponentsTheme = {
    [key: string]: any;
};

export interface RaThemeOptions extends MuiThemeOptions {
    sidebar?: {
        width?: number;
        closedWidth?: number;
    };
    components?: ComponentsTheme;
}

export type ThemeType = 'light' | 'dark';


/**
 * Radiant: A theme emphasizing clarity and ease of use.
 *
 * Uses generous margins, outlined inputs and buttons, no uppercase, and an acid color palette.
 */

const componentsOverrides = (theme: Theme) => {
    const shadows = [
        alpha(theme.palette.primary.main, 0.2),
        alpha(theme.palette.primary.main, 0.1),
        alpha(theme.palette.primary.main, 0.05),
    ];
    
    const isDarkMode = theme.palette.mode === 'dark';
    const primaryColor = isDarkMode ? '#2661BF' : '#4abaab';
    
    return {
        MuiAppBar: {
            styleOverrides: {
                colorSecondary: {
                    backgroundColor: theme.palette.background.default,
                    color: theme.palette.text.primary,
                },
            },
        },
        MuiAutocomplete: {
            defaultProps: {
                fullWidth: true,
            },
        },
        MuiButton: {
            defaultProps: {
                variant: 'outlined' as const,
            },
            styleOverrides: {
                sizeSmall: {
                    padding: `${theme.spacing(0.5)} ${theme.spacing(1.5)}`,
                },
                contained: {
                    '&.MuiButton-containedPrimary': {
                        backgroundColor: primaryColor, 
                        color: '#ffffff',
                        '&:hover': {
                            backgroundColor: isDarkMode ? '#1a4d9f' : '#43b1a3', 
                            color: '#ffffff',
                        },
                        '&:disabled': {
                            backgroundColor: theme.palette.action.disabledBackground,
                            color: theme.palette.action.disabled,
                        },
                    },
                },
                root: {
                    
                    '&.MuiButton-containedPrimary': {
                        backgroundColor: primaryColor, 
                        color: '#ffffff',
                        '&:hover': {
                            backgroundColor: isDarkMode ? '#1a4d9f' : '#43b1a3',
                            color: '#ffffff',
                        },
                    },
                    '&[style*="background-color: rgb(74, 186, 171)"], &[style*="background-color: #4abaab"]': {
                        color: '#ffffff !important',
                    },
                    '&[style*="background-color: rgb(38, 97, 191)"], &[style*="background-color: #2661BF"]': {
                        color: '#ffffff !important',
                    },
                },
            },
        },
        MuiFormControl: {
            defaultProps: {
                variant: 'outlined' as const,
                margin: 'dense' as const,
                size: 'small' as const,
                fullWidth: true,
            },
        },
        MuiPaper: {
            styleOverrides: {
                elevation1: {
                    boxShadow: `${shadows[0]} -2px 2px, ${shadows[1]} -4px 4px,${shadows[2]} -6px 6px`,
                },
                root: {
                    backgroundClip: 'padding-box',
                },
            },
        },
        MuiTableCell: {
            styleOverrides: {
                root: {
                    padding: theme.spacing(1.5),
                    '&.MuiTableCell-sizeSmall': {
                        padding: theme.spacing(1),
                    },
                    '&.MuiTableCell-paddingNone': {
                        padding: 0,
                    },
                },
            },
        },
        MuiTableRow: {
            styleOverrides: {
                root: {
                    '&:last-child td': { border: 0 },
                },
            },
        },
        MuiTextField: {
            defaultProps: {
                variant: 'outlined' as const,
                margin: 'dense' as const,
                size: 'small' as const,
                fullWidth: true,
            },
        },
        RaDatagrid: {
            styleOverrides: {
                root: {
                    '& .RaDatagrid-headerCell': {
                        color: primaryColor,
                    },
                },
            },
        },
        RaFilterForm: {
            styleOverrides: {
                root: {
                    [theme.breakpoints.up('sm')]: {
                        minHeight: theme.spacing(6),
                    },
                },
            },
        },
        RaAppBar: {
            styleOverrides: {
                root: {
                    '& .RaAppBar-title': {
                        fontSize: '1.5rem',
                        fontWeight: 'bold',
                        color: primaryColor,
                    },
                },
            },
        },
        RaLayout: {
            styleOverrides: {
                root: {
                    '& .RaLayout-appFrame': { 
                        marginTop: theme.spacing(5),
                    },
                    '& .RaLayout-sidebar': {
                        '& .MuiToolbar-root': {
                            backgroundColor: primaryColor,
                            color: theme.palette.primary.contrastText,
                            boxShadow: theme.shadows[1],
                        },
                    },
                },
            },
        },
        RaMenuItemLink: {
            styleOverrides: {
                root: {
                    borderLeft: `3px solid transparent`,       
                    margin: theme.spacing(0.5, 1.5, 0.5, -0),
                    borderRadius: theme.spacing(1),          
                    padding: theme.spacing(1, 2, 1, 2),      
                    transition: 'all 0.3s ease',           
                    '&:hover': {
                        // borderRadius: '0px 100px 100px 0px', 
                        backgroundColor: theme.palette.action.hover,
                        transform: 'translateX(5px)',          
                    },
                    '&.RaMenuItemLink-active': {
                        borderLeft: `3px solid ${primaryColor}`, 
                        // borderRadius: '0px 50px 50px 0px',
                        backgroundColor: primaryColor,
      
                        color: '#ffffff',
                        transform: 'translateX(8px)',        
                        '& .MuiSvgIcon-root': {
                            fill: '#ffffff',
                        },
                        '& .MuiListItemText-primary': {
                            fontWeight: 'bold',         
                            color: '#ffffff',
                        },
                    },
                    '.RaSidebar-closed &': {
                        margin: theme.spacing(0.5, 0, 0.5, -0.6), 
                        padding: theme.spacing(1, 0, 1, 2),     
                        '&.RaMenuItemLink-active': {
                            borderRadius: '0px',             
                            transform: 'translateX(0px)', 
                        },
                        '&:hover': {
                            borderRadius: '0px',              
                            transform: 'translateX(0px)',   
                        },
                    },
                },
            },
        },
        RaSimpleFormIterator: {
            defaultProps: {
                fullWidth: true,
            },
        },
        RaTranslatableInputs: {
            defaultProps: {
                fullWidth: true,
            },
        },
        RaSidebar: {
            styleOverrides: {
                root: {
                    '& .RaSidebar-fixed': {
                        height: '110%',
                        backgroundColor: theme.palette.background.paper,
                        borderRight: `1px solid ${theme.palette.divider}`,
                        boxShadow: theme.shadows[2],
                        // Responsive sidebar
                        [theme.breakpoints.down('md')]: {
                            width: '240px',
                        },
                        [theme.breakpoints.down('sm')]: {
                            width: '200px',
                        },
                    },
                    '& .MuiDrawer-paper': {
                        marginTop: '30px', 
                        height: 'calc(100vh - 80px)',
                        backgroundColor: theme.palette.background.paper,
                        backgroundImage: 'linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.02) 100%)',
                        position: 'relative',
                        [theme.breakpoints.down('md')]: {
                            width: '240px',
                        },
                        [theme.breakpoints.down('sm')]: {
                            width: '200px',
                        },
                    },
                }
            },
        },
        RaMenu: {
            styleOverrides: {
                root: {
                    '& .MuiList-root': {
                        paddingTop: theme.spacing(1),
                        paddingBottom: theme.spacing(1),
                    },
                },
            },
        },
        RaList: {
            styleOverrides: {
                root: {
                    marginTop: theme.spacing(-5),
                },
            },
        },
        RaButton: {
            styleOverrides: {
                root: {
                    '&.MuiButton-containedPrimary': {
                        backgroundColor: primaryColor,
                        color: '#ffffff',
                        '&:hover': {
                            backgroundColor: isDarkMode ? '#1a4d9f' : '#43b1a3',
                            color: '#ffffff',
                        },
                    },
                },
            },
        },
    };
};

const alert = {
    error: { main: '#DB488B' },
    warning: { main: '#F2E963' },
    info: { main: '#3ED0EB' },
    success: { main: '#0FBF9F' },
};

const darkPalette: PaletteOptions = {
    primary: { main: '#2661BF', dark: '#1a4d9f' },
    secondary: { main: '#4ABAAB' },
    background: { default: '#110e1c', paper: '#151221' },
    ...alert,
    mode: 'dark' as 'dark',
};

const lightPalette: PaletteOptions = {
    primary: { main: '#4abaab', dark: '#43b1a3ff' }, 
    secondary: { main: '#2661BF' },
    background: { default: '#f9f9f9ff' },
    text: {
        primary: '#544f5a',
        secondary: '#89868D',
    },
    ...alert,
    mode: 'light' as 'light',
};

const createRadiantTheme = (palette: RaThemeOptions['palette']) => {
    const themeOptions = {
        palette,
        shape: { borderRadius: 6 },
        sidebar: { 
            width: 250,             
            closedWidth: 55,  
      
        },
        spacing: 10,
        typography: {
            fontFamily: 'Gabarito, tahoma, sans-serif',
            h1: {
                fontWeight: 500,
                fontSize: '6rem',
            },
            h2: { fontWeight: 600 },
            h3: { fontWeight: 700 },
            h4: { fontWeight: 800 },
            h5: { fontWeight: 900 },
            button: { textTransform: undefined, fontWeight: 700 },
        },
    };
    const theme = createTheme(themeOptions);
    theme.components = componentsOverrides(theme);
    return theme;
};

export const radiantLightTheme = createRadiantTheme(lightPalette);
export const radiantDarkTheme = createRadiantTheme(darkPalette);
