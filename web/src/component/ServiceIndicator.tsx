import RouterIcon from '@mui/icons-material/Router';
import { IconButton, Tooltip } from "@mui/material";
import { useServiceStatus } from '../context/ServiceContext';

const ServiceIndicator = () => {
    const { status } = useServiceStatus();

    const getIcon = () => {
        switch (status) {
            case "connected":
                return <RouterIcon style={{ color: "#4CAF50", width: 30, height: 'auto' }} />;
            case "disconnected":
                return <RouterIcon style={{ color: "#af4c4cff", width: 30, height: 'auto' }} />;
            case "loading":
                return <RouterIcon style={{ color: "#FFA726", width: 30, height: 'auto' }} />;
            default:
                return <RouterIcon />;
        }
    };

    const getTooltipTitle = () => {
        switch (status) {
            case "connected":
                return "Integration Service Connected";
            case "disconnected":
                return "Integration Service Down";
            case "loading":
                return "Checking...";
            default:
                return "Service (Port 8214)";
        }
    };

    return (
        <Tooltip title={getTooltipTitle()}>
            <IconButton color="inherit">
                {getIcon()}
            </IconButton>
        </Tooltip>
    );
};

export default ServiceIndicator;
