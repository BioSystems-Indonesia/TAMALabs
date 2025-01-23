import SettingsIcon from '@mui/icons-material/Settings';
import WifiIcon from '@mui/icons-material/Wifi';
import WifiOffIcon from '@mui/icons-material/WifiOff';
import { Breadcrumbs, Stack } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import { useEffect, useState, type ReactNode, useRef } from 'react';
import { AppBar, Button, CheckForApplicationUpdate, Layout, Link, LoadingIndicator, TitlePortal, ToggleThemeButton } from 'react-admin';
import { useLocation, useNavigate } from "react-router-dom";
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { toTitleCase } from '../helper/format';


const SettingsButton = () => (
    <Link to="/settings" color={"inherit"} LinkComponent={Link}>
        <IconButton color="inherit" LinkComponent={Link}>
            <SettingsIcon />
        </IconButton>
    </Link>
);

const AppIndicator = () => {
    var [state, setState] = useState("loading");
    var [detailState, setDetailState] = useState({});
    var timer = useRef(0);
    const fetchData = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/server/status`);
            const data = await response.json();
            setDetailState(data);
            setState("online");
        } catch {
            setState("offline")
        }
    };

    useEffect(() => {
        fetchData();
        timer.current = setInterval(fetchData, 5000);
        console.log(timer.current)
        return () => { if (timer.current != 0) clearInterval(timer.current) }
    }, [])

    const icon = () => {
        switch(state) {
            case "online": return <WifiIcon />
            case "offline": return <WifiOffIcon />
            case "loading": return <WifiOffIcon />
            default: <WifiOffIcon />
        };
    }

    const tooltipTitle = () => {
        return "server: " + detailState.rest + " " + "hl7tcp: " + detailState.hl7tcp;
    }

    return (
        <Tooltip title={tooltipTitle()}>
            <IconButton color="inherit">
                {icon()}
            </IconButton>
        </Tooltip>
    )
}

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
