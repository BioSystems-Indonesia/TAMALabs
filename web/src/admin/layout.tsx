import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import SettingsIcon from '@mui/icons-material/Settings';
import { Stack } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useEffect, useState, type ReactNode } from 'react';
import { AppBar, Button, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, TitlePortal, ToggleThemeButton } from 'react-admin';
import { useLocation, useNavigate } from "react-router-dom";
import AppIndicator from '../component/AppIndicator';
import Breadcrumbs, { type BreadcrumbsLink } from '../component/Breadcrumbs';
import { toTitleCase } from '../helper/format';


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
        <Stack direction={"row"}>
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

export const DefaultLayout = ({ children }: { children: ReactNode }) => {
    return (
        <Layout sx={{}} appBar={MyAppBar}>
            <Stack direction={"row"} gap={2}>

                <DynamicBreadcrumbs />
            </Stack>
            {children}
            <CheckForApplicationUpdate />
        </Layout>
    )
};
