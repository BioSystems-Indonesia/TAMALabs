import { GridLegacy as Grid, LinearProgress, SxProps, Alert } from '@mui/material';
import CircularProgress from '@mui/material/CircularProgress';
import Stack from '@mui/material/Stack';

import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import WarningIcon from '@mui/icons-material/Warning';
import Typography from '@mui/material/Typography';
import { useCallback, useEffect, useRef, useState } from 'react';
import { AutocompleteInput, BooleanInput, Button, Form, InputHelperText, Link, RecordContextProvider, ReferenceInput, required, useNotify, useRefresh } from 'react-admin';
import { SubmitHandler, useFormContext } from 'react-hook-form';
import { getRefererParam } from '../../hooks/useReferer';

import { DeviceForm } from '../device';

type WorkOrderStatus = 'IDLE' | 'PENDING' | 'IN_PROGRESS' | 'DONE' | 'INCOMPLETE' | 'ERROR';
type StreamStatus = 'DONE' | 'IN_PROGRESS' | 'INCOMPLETE';
type ProgressCallback = (percentage: number, status: StreamStatus) => void;
const WorkOrderAction = {
    run: "run",
    cancel: "cancel",
} as const
type WorkOrderActionValue = typeof WorkOrderAction[keyof typeof WorkOrderAction];



type RunWorkOrder = {
    work_order_ids: number[];
    device_id: number;
    urgent: boolean;
    action: WorkOrderActionValue;
}

interface StreamResult {
    percentage: number;
    status: StreamStatus;
    errorCause?: string;
}


// --- Helper Function to Process Server-Sent Events Stream (TypeScript) ---
/**
 * Processes a ReadableStream from a fetch response expecting SSE.
 * @param {Response} response - The fetch Response object.
 * @param {ProgressCallback} onProgress - Callback function invoked with (percentage, status).
 * @param {AbortSignal | undefined} signal - AbortSignal to allow cancellation.
 * @returns {Promise<StreamResult>} - Promise resolving with the final percentage and status.
 */
async function processSSEStream(
    response: Response,
    onProgress: ProgressCallback,
    signal: AbortSignal | undefined
): Promise<StreamResult> {
    if (!response.body) {
        throw new Error('Response body is null');
    }
    const contentType: string | null = response.headers.get('Content-Type');
    if (!contentType || !contentType.includes('text/event-stream')) {
        // Try to read the response body as text for more detailed error info
        let errorText = `Expected text/event-stream, but received ${contentType || 'N/A'}.`;
        try {
            const text: string = await response.text();
            errorText += ` Server response: ${text}`;
        } catch (e: unknown) {
            // Ignore error reading body if it fails
            console.warn("Could not read error response body:", e);
        }
        throw new Error(errorText);
    }

    const reader: ReadableStreamDefaultReader<Uint8Array> = response.body.getReader();
    const decoder: TextDecoder = new TextDecoder('utf-8');
    let buffer: string = '';
    let lastKnownPercentage: number = 0;
    let lastKnownStatus: StreamStatus | '' = '';

    try {
        // Loop indefinitely to read from the stream
        while (true) {
            if (signal?.aborted) {
                console.log('Aborting stream read due to signal.');
                try { await reader.cancel('Request aborted by client'); } catch (cancelError: unknown) { console.warn("Error cancelling reader:", cancelError); }
                const abortError = new Error('Request aborted');
                abortError.name = 'AbortError';
                throw abortError;
            }

            const { value, done }: ReadableStreamReadResult<Uint8Array> = await reader.read();

            if (done) {
                console.log('Stream finished.');
                if (buffer.trim()) {
                    console.warn("Stream ended with unprocessed buffer:", buffer);
                }
                break;
            }

            buffer += decoder.decode(value, { stream: true });

            let boundaryIndex: number;
            while ((boundaryIndex = buffer.indexOf('\n\n')) >= 0) {
                const message: string = buffer.substring(0, boundaryIndex);
                buffer = buffer.substring(boundaryIndex + 2);
                console.log('Processing message:', message); // Log the message for debugging purposes

                const lines: string[] = message.split('\n');
                for (const line of lines) {
                    if (line.startsWith('event:')) {
                        const eventName: string = line.substring(6).trim();
                        if (eventName === 'error') {
                            const errorCause = lines.length > 1 ? lines[1].trim() : 'Unknown error';

                            const abortError = new Error(errorCause);
                            abortError.name = 'AbortError';
                            throw abortError;
                        }
                    }

                    if (line.startsWith('data:')) {
                        const dataContent: string = line.substring(5).trim();
                        if (dataContent) {
                            try {
                                const params: URLSearchParams = new URLSearchParams(dataContent);
                                const percentageStr: string | null = params.get('percentage');
                                const statusStr: string | null = params.get('status');
                                const percentage: number = percentageStr ? parseInt(percentageStr, 10) : NaN;

                                const isValidStatus = (s: string | null): s is StreamStatus =>
                                    s === 'DONE' || s === 'IN_PROGRESS' || s === 'INCOMPLETE';

                                if (!isNaN(percentage) && isValidStatus(statusStr)) {
                                    lastKnownPercentage = percentage;
                                    lastKnownStatus = statusStr;
                                    onProgress(percentage, statusStr);
                                } else {
                                    console.warn('Received invalid or incomplete data:', { dataContent });
                                }
                            } catch (parseError: unknown) {
                                console.error('Error parsing SSE data content:', dataContent, parseError);
                            }
                        }
                    }
                }
            }
        }
    } catch (error: unknown) { // Catch block parameter must be any or unknown in TS
        console.error('Error reading or processing stream:', error);
        if (error instanceof Error && error.name === 'AbortError') {
            console.log('Stream processing aborted.');
            return { percentage: lastKnownPercentage, status: 'INCOMPLETE', errorCause: error.message };
        }

        if (error instanceof Error) {
            throw error;
        } else {
            throw new Error(`An unknown error occurred during stream processing: ${String(error)}`);
        }
    } finally {
        if (!signal?.aborted) {
            try {
                if (reader && typeof reader.closed === 'boolean' && !reader.closed) {
                    await reader.cancel();
                }
            } catch (e: unknown) {
                console.warn("Error during final reader cancellation (might be expected if stream ended normally):", e);
            }
        }
    }

    const finalStatus = lastKnownStatus || 'INCOMPLETE';
    console.log("Stream processing complete. Final state:", { percentage: lastKnownPercentage, status: finalStatus });
    return { percentage: lastKnownPercentage, status: finalStatus };
}

type RunWorkOrderFormProps = {
    workOrderIDs: number[];
    showCancelButton?: boolean;
    showRunButton?: boolean;
    isProcessing: boolean;
    defaultDeviceID?: number;
    setIsProcessing: (isProcessing: boolean) => void;
}

export default function RunWorkOrderForm(props: RunWorkOrderFormProps) {
    const [percentage, setPercentage] = useState<number>(0);
    const [status, setStatus] = useState<WorkOrderStatus>('IDLE');
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const abortControllerRef = useRef<AbortController | null>(null);
    const refresh = useRefresh();

    // Callback passed to the stream processor to update component state during streaming
    const handleProgressUpdate: ProgressCallback = useCallback((newPercentage, newStatus) => {
        setPercentage(prev => newPercentage ?? prev);
        setStatus(newStatus);
        setError(null);
    }, []);

    // Function to initiate the work order request and stream processing
    const startWorkOrder = useCallback(async (data: RunWorkOrder) => {
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            console.log('Previous request aborted.');
        }

        const controller = new AbortController();
        abortControllerRef.current = controller;
        const signal = controller.signal;

        setIsLoading(true);
        setError(null);
        setPercentage(0);
        setStatus('PENDING');

        try {
            let url
            switch (data.action) {
                case WorkOrderAction.run:
                    url = `${import.meta.env.VITE_BACKEND_BASE_URL}/work-order/run`;
                    break;
                case WorkOrderAction.cancel:
                    url = `${import.meta.env.VITE_BACKEND_BASE_URL}/work-order/cancel`;
                    break;
            }

            const response: Response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Accept': 'text/event-stream',
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Use cookies instead of Authorization header
                body: JSON.stringify(data),
                signal,
            });

            if (!response.ok) {
                let errorBody = `Server responded with status ${response.status}`;
                try {
                    const text: string = await response.text();
                    errorBody += `: ${text || '(no details provided)'}`;
                } catch (_) { }
                throw new Error(errorBody);
            }

            const finalResult: StreamResult = await processSSEStream(response, handleProgressUpdate, signal);

            setStatus(finalResult.status);
            setPercentage(finalResult.percentage);
            if (finalResult.status === 'INCOMPLETE') {
                setError(`Failed to send request, message: '${finalResult.errorCause?.split("data: ")[1]}'`);
                refresh();
            } else {
                notify("Work order sent successfully", {
                    type: 'success',
                    autoHideDuration: 3000,
                })
                refresh();
            }

        } catch (error: unknown) {
            if (error instanceof Error && error.name === 'AbortError') {
                console.log('Fetch operation was aborted.');
                setStatus('INCOMPLETE');
            } else if (error instanceof Error) {
                console.error("Error during work order processing:", error);
                setError(error.message);
                setStatus('ERROR');
                setPercentage(0);
            } else {
                console.error("An unknown error occurred:", error);
                setError('An unexpected error occurred.');
                setStatus('ERROR');
                setPercentage(0);
            }
        } finally {
            setIsLoading(false);
            if (abortControllerRef.current === controller && !signal.aborted) {
                abortControllerRef.current = null;
            }
        }
    }, [handleProgressUpdate]);

    useEffect(() => {
        return () => {
            const controller = abortControllerRef.current;
            if (controller) {
                console.log('Component unmounting, aborting fetch...');
                controller.abort();
                abortControllerRef.current = null;
            }
        };
    }, []);

    useEffect(() => {
        if (props.setIsProcessing) {
            props.setIsProcessing(isLoading || status === 'PENDING' || status === 'IN_PROGRESS');
        }
    }, [isLoading, status, props.setIsProcessing]);
    const notify = useNotify();

    const onSubmit: SubmitHandler<any> = (data, event) => {
        console.log("Run work order", data)
        if (!data.device_id) {
            notify('Please select device to run', {
                type: 'error',
            });
            return;
        }

        const submitter = (event?.nativeEvent as SubmitEvent | undefined)?.submitter as HTMLButtonElement | undefined;
        const action = submitter?.value as WorkOrderActionValue | undefined;

        startWorkOrder({
            work_order_ids: props.workOrderIDs ?? [data.id],
            device_id: data.device_id,
            urgent: data.urgent,
            action: action ?? WorkOrderAction.run,
        });
    }

    return <RecordContextProvider value={{
        device_id: props.defaultDeviceID ?? undefined,
    }}>
        <Form disabled={props.isProcessing} onSubmit={onSubmit}>
            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>)}
            <Grid direction={"row"} sx={{
                width: "100%",
            }} container>
                <Grid item xs={12} md={9} sx={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                }}>
                    <ReferenceInput source={"device_id"} reference={"device"} disabled={props.isProcessing} >
                        <AutocompleteInput source={"device_id"} validate={[required()]} create={<DeviceForm />} sx={{
                            margin: 0,
                        }} disabled={props.isProcessing} helperText={
                            <Link to={"/device/create?" + getRefererParam()}>
                                <InputHelperText helperText="Create new device"></InputHelperText>
                            </Link>
                        }
                        />
                    </ReferenceInput>
                </Grid>
                <Grid item xs={12} md={3} sx={{
                    display: "flex",
                    paddingLeft: "24px",
                    justifyContent: "start",
                    alignItems: "center",
                }}>
                    <BooleanInput source={"urgent"} disabled={props.isProcessing} label="Urgent" />
                </Grid>
            </Grid>
            <RunWorkOrderSubmit isPending={props.isProcessing} sx={{
                marginTop: "12px",
            }} percentage={percentage} status={status} showCancelButton={props.showCancelButton} showRunButton={props.showRunButton ?? true} />
        </Form>
    </RecordContextProvider>
}


type RunWorkOrderSubmitProps = {
    isPending: boolean;
    percentage: number;
    status: WorkOrderStatus;
    sx?: SxProps;
    showCancelButton?: boolean;
    showRunButton?: boolean;
}

function RunWorkOrderSubmit({ isPending, percentage, status, sx, showCancelButton, showRunButton }: RunWorkOrderSubmitProps) {
    const { watch } = useFormContext()

    return (
        <Stack sx={sx}>
            <Stack spacing={1} width="100%">
                <Stack direction={"row"} width={"100%"} spacing={2}>
                    {
                        showCancelButton &&
                        <Button
                            label="Cancel Work Order"
                            disabled={isPending || !watch("device_id")}
                            variant="contained"
                            type='submit'
                            name="formAction"
                            value={WorkOrderAction.cancel}
                            color='error'
                            sx={{
                                width: showRunButton ? "20%" : "100%",
                                cursor: "pointer"
                            }}
                        >
                            {isPending ? <CircularProgress size={12} variant='indeterminate' color='primary' /> : <PlayCircleFilledIcon />}
                        </Button>
                    }
                    {
                        showRunButton &&
                        <Button
                            label="Run Work Order"
                            disabled={isPending || !watch("device_id")}
                            variant="contained"
                            type='submit'
                            name="formAction"
                            value={WorkOrderAction.run}
                            sx={{
                                width: showCancelButton ? "80%" : "100%",
                                cursor: "pointer"
                            }}
                        >
                            {isPending ? <CircularProgress size={12} variant='indeterminate' color='primary' /> : <PlayCircleFilledIcon />}
                        </Button>
                    }
                </Stack>
                {isPending && (
                    <Stack spacing={1} width="100%">
                        <Stack
                            direction="row"
                            spacing={1}
                            alignItems="center"
                            justifyContent="center"
                            sx={{
                                backgroundColor: 'warning.light',
                                padding: '12px 20px',
                                borderRadius: '8px',
                                border: '2px solid',
                                borderColor: 'warning.main',
                                boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
                                transform: 'scale(1.02)',
                                transition: 'all 0.2s ease-in-out',
                                '&:hover': {
                                    transform: 'scale(1.03)',
                                    boxShadow: '0 6px 12px rgba(0,0,0,0.15)',
                                }
                            }}
                        >
                            <WarningIcon
                                sx={{
                                    color: 'dark',
                                    fontSize: '24px',
                                    animation: 'pulse 2s infinite'
                                }}
                            />
                            <Typography
                                variant="body2"
                                sx={{
                                    color: 'dark',
                                    fontWeight: 'bold',
                                    letterSpacing: '0.5px',
                                    '@keyframes pulse': {
                                        '0%': { opacity: 1 },
                                        '50%': { opacity: 0.7 },
                                        '100%': { opacity: 1 },
                                    }
                                }}
                            >
                                Please wait and don't close this tab browser
                            </Typography>
                        </Stack>
                        <LinearProgress
                            variant="determinate"
                            value={percentage}
                            sx={{
                                height: 8,
                                borderRadius: 4
                            }}
                        />
                        <Typography variant="caption" textAlign="center">
                            {`${Math.round(percentage)}%`}
                        </Typography>
                    </Stack>
                )}
            </Stack>
            {!isPending && !watch("device_id") && <Typography color='error' fontSize={12}>Please pick device to run</Typography>}
        </Stack>
    );
}
