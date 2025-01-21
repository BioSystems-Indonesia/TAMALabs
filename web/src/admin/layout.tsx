import SettingsIcon from '@mui/icons-material/Settings';
import { Breadcrumbs } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useEffect, type ReactNode, useState } from 'react';
import { AppBar, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, TitlePortal, ToggleThemeButton } from 'react-admin';
import { useLocation } from "react-router-dom";
import { toTitleCase } from '../helper/format';


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

export const DefaultLayout = ({ children }: { children: ReactNode }) => (
    <Layout sx={{}} appBar={MyAppBar}>
        <DynamicBreadcrumbs />
        {children}
        <CheckForApplicationUpdate />
    </Layout>
);
