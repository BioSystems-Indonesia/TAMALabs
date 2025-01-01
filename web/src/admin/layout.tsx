import SettingsIcon from '@mui/icons-material/Settings';
import IconButton from '@mui/material/IconButton';
import { type ReactNode } from 'react';
import { AppBar, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, LocalesMenuButton, TitlePortal, ToggleThemeButton } from 'react-admin';
import { useLocation } from "react-router-dom";

const SettingsButton = () => (
    <Link to="/settings" color={"inherit"} LinkComponent={Link}>
        <IconButton color="inherit" LinkComponent={Link}>
            <SettingsIcon />
        </IconButton>
    </Link>
);

const MyAppBar = () => {
    const location = useLocation()
    const appBarAlwaysOn = () => {
        if (location.pathname.includes('/work-order/')) {
            return true;
        }

        return false
    }

    return (
        <AppBar color="primary" sx={{
            position: 'fixed',
        }} alwaysOn={appBarAlwaysOn()}
            toolbar={
                <>
                    <SettingsButton />
                    <ToggleThemeButton />
                    <LoadingIndicator />
                </>
            }
        >
            <TitlePortal />
        </AppBar>
    )
};


export const DefaultLayout = ({ children }: { children: ReactNode }) => (
    <Layout sx={{}} appBar={MyAppBar}>
        {children}
        <CheckForApplicationUpdate />
    </Layout>
);