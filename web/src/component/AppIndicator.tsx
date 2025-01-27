import WifiIcon from '@mui/icons-material/Wifi';
import WifiOffIcon from '@mui/icons-material/WifiOff';
import { IconButton, Tooltip } from "@mui/material";
import { useEffect, useRef, useState } from "react";

const AppIndicator = () => {
    var [state, setState] = useState("loading");
    var [detailState, setDetailState] = useState({
        rest: "",
        hl7tcp: "",
    });
    var timer = useRef<ReturnType<typeof setInterval> | null>(null)
    const fetchData = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/server/status`);
            const data = await response.json();
            setDetailState(data);
            setState("online");
        } catch {
            setState("offline")
            setDetailState({
                rest: "offline",
                hl7tcp: "offline",
            })
        }
    };

    useEffect(() => {
        fetchData();
        timer.current = setInterval(fetchData, 5000);
        return () => { if (timer.current != null) clearInterval(timer.current) }
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

export default AppIndicator;