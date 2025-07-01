import { Box, Tooltip, Typography } from '@mui/material';
import React from 'react';

export interface ConnectionStatusProps {
    deviceId: number;
    status: ConnectionResponse;
}

export interface ConnectionResponse {
    device_id: number;
    message: string;
    status: 'connected' | 'not_supported' | 'disconnected' | 'standby';
}

export const ConnectionStatus: React.FC<ConnectionStatusProps> = ({ deviceId, status }) => {
    const getStatusColor = () => {
        switch (status?.status) {
            case 'connected':
                return '#4caf50'; // Green
            case 'not_supported':
                return '#9e9e9e'; // Grey
            case 'standby':
                return '#ff9800'; // Yellow
            case 'disconnected':
                return '#f44336'; // Red
            default:
                return '#9e9e9e'; // Grey
        }
    };

    return (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box
                sx={{
                    width: 12,
                    height: 12,
                    borderRadius: '50%',
                    backgroundColor: getStatusColor(),
                    transition: 'background-color 0.3s ease'
                }}
            />
            <Tooltip title={status?.message || 'Checking connection...'}>
                <Typography variant="body2" sx={{ minWidth: 100 }}>
                    {status?.status || 'Checking...'}
                </Typography>
            </Tooltip>
        </Box>
    );
}; 