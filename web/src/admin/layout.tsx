import Box from '@mui/material/Box';
import { type ReactNode } from 'react';
import { AppBar, CheckForApplicationUpdate, Layout, TitlePortal } from 'react-admin';
import { useLocation } from "react-router-dom";

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
        }} alwaysOn={appBarAlwaysOn()}>
            <TitlePortal />
            <Box flex="1" />
        </AppBar>
    )
};


export const DefaultLayout = ({ children }: { children: ReactNode }) => (
    <Layout sx={{}} appBar={MyAppBar}>
        {children}
        <CheckForApplicationUpdate />
    </Layout>
);