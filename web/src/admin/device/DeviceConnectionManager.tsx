import React, { useEffect, useState } from 'react';
import { Box, Typography } from '@mui/material';
import { LOCAL_STORAGE_ACCESS_TOKEN } from '../../types/constant';

export interface DeviceConnectionManagerProps {
    deviceIds: number[];
    onStatusUpdate?: (deviceId: number, status: ConnectionResponse) => void;
}

export interface ConnectionResponse {
    device_id: number;
    sender_message: string;
    sender_status: 'connected' | 'not_supported' | 'disconnected';
    receiver_message: string;
    receiver_status: 'connected' | 'not_supported' | 'disconnected';
}

export const DeviceConnectionManager: React.FC<DeviceConnectionManagerProps> = ({ deviceIds, onStatusUpdate }) => {
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (deviceIds.length === 0) return;

        const checkConnections = async () => {
            try {
                const params = new URLSearchParams();
                deviceIds.forEach(id => params.append('device_ids', id.toString()));
                const response = await fetch(`${import.meta.env.VITE_BACKEND_BASE_URL}/device/connection?${params.toString()}`, {
                    headers: {
                        'Accept': 'text/event-stream',
                        'Authorization': `Bearer ${localStorage.getItem(LOCAL_STORAGE_ACCESS_TOKEN)}`
                    },
                });

                const reader = response.body?.getReader();
                if (!reader) return;

                while (true) {
                    const { done, value } = await reader.read();
                    if (done) break;

                    const text = new TextDecoder().decode(value);
                    const events = text.split('\n\n').filter(Boolean);

                    for (const event of events) {
                        if (event.startsWith('data: ')) {
                            const data = event.slice(6);
                            console.log("received data", data)
                            const params = new URLSearchParams(data);
                            const deviceId = parseInt(params.get('device_id') || '0');
                            const senderMessage = params.get('sender_message') || '';
                            const senderStatus = params.get('sender_status') as ConnectionResponse['sender_status'];
                            const receiverMessage = params.get('receiver_message') || '';
                            const receiverStatus = params.get('receiver_status') as ConnectionResponse['receiver_status'];

                            onStatusUpdate?.(deviceId, {
                                device_id: deviceId,
                                sender_message: senderMessage,
                                sender_status: senderStatus,
                                receiver_message: receiverMessage,
                                receiver_status: receiverStatus
                            });
                        }
                    }
                }
            } catch (error) {
                console.error('Failed to check connections:', error);
                setError('Failed to check device connections');
            }
        };

        checkConnections();
        const interval = setInterval(checkConnections, 10000); // Check every 10 seconds
        return () => clearInterval(interval);
    }, [deviceIds.join(',')]);

    if (error) {
        return (
            <Box sx={{ p: 1, bgcolor: 'error.light', borderRadius: 1 }}>
                <Typography variant="body2" color="error">
                    {error}
                </Typography>
            </Box>
        );
    }

    return null;
}; 