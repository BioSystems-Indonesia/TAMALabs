import WifiIcon from '@mui/icons-material/Wifi';
import WifiOffIcon from '@mui/icons-material/WifiOff';
import { IconButton, Tooltip } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import useAxios from '../hooks/useAxios';
// import { useNavigate } from 'react-router-dom';

const AppIndicator = () => {
    var [state, setState] = useState("loading");
    var [detailState, setDetailState] = useState({
        rest: "",
        hl7tcp: "",
    });
    var timer = useRef<ReturnType<typeof setInterval> | null>(null)
    const axios = useAxios()
    const fetchData = async () => {
        try {
            const response = await axios.get(`/server/status`);
            setDetailState(response.data);
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
            case "online": return <WifiIcon style={{color: "#4CAF50", width: 30, height:'auto'}}/>
            case "offline": return <WifiOffIcon style={{color: "#af4c4cff", width: 30, height:'auto'}}/>
            case "loading": return <WifiOffIcon style={{width: 30, height:'auto'}}/>
            default: <WifiOffIcon />
        };
    }

    const tooltipTitle = () => {
        return "server: " + detailState.rest ;
    }

    // const navigate = useNavigate()

    return (
        <Tooltip title={tooltipTitle()}>
            <IconButton color="inherit" onClick={() => {
                // navigate("/device")
            }}>
                {icon()}
            </IconButton>
        </Tooltip>
    )
}

export default AppIndicator;