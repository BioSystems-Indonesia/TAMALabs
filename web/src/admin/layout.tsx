import SettingsIcon from '@mui/icons-material/Settings';
import { Breadcrumbs, Stack } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useEffect, useState, type ReactNode } from 'react';
import { AppBar, Button, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, TitlePortal, ToggleThemeButton } from 'react-admin';
import { useLocation, useNavigate } from "react-router-dom";
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { toTitleCase } from '../helper/format';
import AppIndicator from '../component/AppIndicator';


const SettingsButton = () => (
    <Link to="/settings" color={"inherit"}>
        <IconButton color="inherit">
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

        if (location.pathname.includes('/test-template/')) {
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
                    <AppIndicator />
                    <ToggleThemeButton />
                    <LoadingIndicator />
                </>
            }
        >
            <TitlePortal />
        </AppBar>
    )
};


type PathConfiguration = {
    path: string
    type: string
}

const DynamicBreadcrumbs = () => {
    const location = useLocation();
    const [paths, setPaths] = useState<Array<PathConfiguration>>([])

    useEffect(() => {
        // Remove first slash then split
        const pathSplit = location.pathname.substring(1).split("/")

        const pathsConfig: Array<PathConfiguration> = []
        pathSplit.forEach((val: string) => {
            pathsConfig.push({
                path: val,
                type: "general",
            })
        })

        setPaths(pathsConfig)
    }, [location.pathname])

    return (
        <Breadcrumbs aria-label="breadcrumb">
            {paths.map((val, i) => {
                const generateHref = (): string => {
                    const pathUntil = paths.slice(0, i + 1)
                    return "/" + pathUntil.map(v => v.path).filter(v => v).join("/") + location.search
                }

                return (
                    <Link underline="hover" color={
                        i == paths.length - 1 ? "text.primary" : "inherit"}
                        to={generateHref()} key={i} >
                        {toTitleCase(val.path)}
                    </Link>
                )
            })}
        </Breadcrumbs>
    )
}

export const DefaultLayout = ({ children }: { children: ReactNode }) => {
    const navigate = useNavigate();
    const location = useLocation();

    return (
        <Layout sx={{}} appBar={MyAppBar}>
            <Stack direction={"row"} gap={2}>
                <Button label='Back' variant='contained' onClick={() => navigate(-1)} sx={{
                    display: location.pathname.split("/").length > 2 ? 'flex' : 'none'
                }}>
                    <ArrowBackIcon />
                </Button>

                <DynamicBreadcrumbs />
            </Stack>
            {children}
            <CheckForApplicationUpdate />
        </Layout>
    )
};
