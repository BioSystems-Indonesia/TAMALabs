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
                        color: theme.palette.primary.main,
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
                        color: theme.palette.primary.main,
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
                            minHeight: '64px',
                            backgroundColor: theme.palette.primary.main,
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
                    borderLeft: `3px solid transparent`,        // Border default
                    margin: theme.spacing(0.5, 2, 0.5, -0.6),            // Margin hanya atas-bawah
                    borderRadius: theme.spacing(1),           // Border radius
                    padding: theme.spacing(1, 2, 1, 2),       // Padding: top, right, bottom, left (0 untuk left)
                    transition: 'all 0.3s ease',             // Transisi smooth
                    '&:hover': {
                        borderRadius: '0px 100px 100px 0px',  // Rounded saat hover
                        backgroundColor: theme.palette.action.hover,
                        transform: 'translateX(5px)',          // Efek slide
                    },
                    '&.RaMenuItemLink-active': {
                        borderLeft: `3px solid ${theme.palette.primary.main}`,
                        borderRadius: '0px 50px 50px 0px',
                        backgroundImage: `linear-gradient(98deg, ${theme.palette.primary.light}, ${theme.palette.primary.dark} 94%)`,
                        boxShadow: theme.shadows[3],           // Shadow lebih tebal
                        color: theme.palette.primary.contrastText,
                        transform: 'translateX(8px)',          // Efek slide lebih jauh
                        '& .MuiSvgIcon-root': {
                            fill: theme.palette.primary.contrastText,
                        },
                        '& .MuiListItemText-primary': {
                            fontWeight: 'bold',                // Teks bold saat aktif
                        },
                    },
                    // Style khusus untuk sidebar yang di-collapse
                    '.RaSidebar-closed &': {
                        margin: theme.spacing(0.5, 0, 0.5, -0.6),  // Hapus margin kanan
                        padding: theme.spacing(1, 0, 1, 2),        // Hapus padding kanan
                        '&.RaMenuItemLink-active': {
                            borderRadius: '0px',              // Tidak ada rounded saat sidebar collapse
                            transform: 'translateX(0px)',      // Tidak ada efek slide
                        },
                        '&:hover': {
                            borderRadius: '0px',               // Tidak ada rounded saat hover di sidebar collapse
                            transform: 'translateX(0px)',      // Tidak ada efek slide
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
                        backgroundColor: theme.palette.background.paper,
                        backgroundImage: 'linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.02) 100%)',
                        // Responsive drawer
                        [theme.breakpoints.down('md')]: {
                            width: '240px',
                        },
                        [theme.breakpoints.down('sm')]: {
                            width: '200px',
                        },
                    },
                },
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
    };
};

const alert = {
    error: { main: '#DB488B' },
    warning: { main: '#F2E963' },
    info: { main: '#3ED0EB' },
    success: { main: '#0FBF9F' },
};

const darkPalette: PaletteOptions = {
    primary: { main: '#2661BF' },
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
