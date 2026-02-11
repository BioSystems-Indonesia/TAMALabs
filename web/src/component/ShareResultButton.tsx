import ShareIcon from '@mui/icons-material/Share';
import { Button as MUIButton, CircularProgress, Tooltip } from "@mui/material";
import { useState } from "react";
import { useNotify } from "react-admin";
import { useServiceStatus } from '../context/ServiceContext';
import type { WorkOrder } from '../types/work_order';

interface ShareResultButtonProps {
    workOrder: WorkOrder;
}

const ShareResultButton = ({ workOrder }: ShareResultButtonProps) => {
    const { status } = useServiceStatus();
    const notify = useNotify();
    const [isLoading, setIsLoading] = useState(false);

    const isServiceConnected = status === "connected";

    const handleShare = async (e: React.MouseEvent) => {
        e.stopPropagation();
        setIsLoading(true);
        try {
            const response = await fetch('http://localhost:8214/generate_result_public', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "barcode": workOrder.barcode,
                }),
            });

            if (!response.ok) {
                throw new Error('Failed to generate public link');
            }

            const data = await response.text();
            await navigator.clipboard.writeText(data);

            notify('Public link copied to clipboard!', {
                type: 'success',
            });
        } catch (error) {
            notify('Failed to generate public link', {
                type: 'error',
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Tooltip
            title={
                !isServiceConnected
                    ? "Integration service is not available"
                    : "Generate and copy public link"
            }
        >
            <span>
                <MUIButton
                    variant="contained"
                    color="primary"
                    startIcon={isLoading ? <CircularProgress size={16} color="inherit" /> : <ShareIcon />}
                    size="small"
                    onClick={handleShare}
                    disabled={!isServiceConnected || isLoading}
                    sx={{
                        textTransform: 'none',
                        fontSize: '12px',
                        whiteSpace: 'nowrap',
                        '&:disabled': {
                            backgroundColor: 'action.disabled',
                        }
                    }}
                >
                    Share
                </MUIButton>
            </span>
        </Tooltip>
    );
};

export default ShareResultButton;
