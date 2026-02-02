import { useEffect, useState } from "react";
import {
    Card,
    CardContent,
    CircularProgress,
    Divider,
    Stack,
    Typography,
    Button,
    Switch,
    FormControlLabel,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Snackbar,
    Alert,
    Chip,
    Box,
    IconButton,
    Tooltip,
} from "@mui/material";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Edit, SimpleForm, TextInput } from "react-admin";
import MUITextField from "@mui/material/TextField";
import useAxios from "../../hooks/useAxios";
import StorageIcon from "@mui/icons-material/Storage";
import PublicIcon from '@mui/icons-material/Public';
import SettingsEthernetIcon from '@mui/icons-material/SettingsEthernet';
import HistoryIcon from '@mui/icons-material/History';
import InfoIcon from '@mui/icons-material/Info';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { IntegrationInstructionsRounded } from "@mui/icons-material";
import ScheduleIcon from '@mui/icons-material/Schedule';
import RefreshIcon from '@mui/icons-material/Refresh';

export const ConfigList = () => {
    const axios = useAxios();
    const { data, isPending } = useQuery({
        queryKey: ["server-info"],
        queryFn: async ({ signal }) => {
            const { data } = await axios.get("/ping");
            return data;
        },
    });

    const queryClient = useQueryClient();
    const { data: configEntry } = useQuery({
        queryKey: ["config", "KhanzaIntegrationEnabled"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/KhanzaIntegrationEnabled`);
            return data;
        },
        enabled: true,
    });

    const { data: selectedEntry } = useQuery({
        queryKey: ["config", "SelectedSimrs"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/SelectedSimrs`);
            return data;
        },
        enabled: true,
    });

    const { data: khanzaBridgeEntry } = useQuery({
        queryKey: ["config", "KhanzaBridgeDatabaseDSN"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/KhanzaBridgeDatabaseDSN`);
            return data;
        },
        enabled: true,
    });

    const { data: khanzaMainEntry } = useQuery({
        queryKey: ["config", "KhanzaMainDatabaseDSN"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/KhanzaMainDatabaseDSN`);
            return data;
        },
        enabled: true,
    });

    const { data: khanzaMethodEntry } = useQuery({
        queryKey: ["config", "KhanzaConnectionMethod"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/KhanzaConnectionMethod`);
            return data;
        },
        enabled: true,
    });

    // SIMRS Configuration Queries
    const { data: simrsEnabledEntry } = useQuery({
        queryKey: ["config", "SimrsIntegrationEnabled"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/SimrsIntegrationEnabled`);
            return data;
        },
        enabled: true,
    });

    const { data: simrsDsnEntry } = useQuery({
        queryKey: ["config", "SimrsDatabaseDSN"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/SimrsDatabaseDSN`);
            return data;
        },
        enabled: true,
    });

    // SIMGOS Configuration Queries
    const { data: simgosEnabledEntry } = useQuery({
        queryKey: ["config", "SimgosIntegrationEnabled"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/SimgosIntegrationEnabled`);
            return data;
        },
        enabled: true,
    });

    const { data: simgosDsnEntry } = useQuery({
        queryKey: ["config", "SimgosDatabaseDSN"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/SimgosDatabaseDSN`);
            return data;
        },
        enabled: true,
    });

    // TechnoMedic Configuration Queries
    const { data: technomedicEnabledEntry } = useQuery({
        queryKey: ["config", "TechnoMedicIntegrationEnabled"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/TechnoMedicIntegrationEnabled`);
            return data;
        },
        enabled: true,
    });

    // Backup Configuration Queries
    const { data: backupScheduleTypeEntry } = useQuery({
        queryKey: ["config", "BackupScheduleType"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/BackupScheduleType`);
            return data;
        },
        enabled: true,
    });

    const { data: backupIntervalEntry } = useQuery({
        queryKey: ["config", "BackupInterval"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/BackupInterval`);
            return data;
        },
        enabled: true,
    });

    const { data: backupTimeEntry } = useQuery({
        queryKey: ["config", "BackupTime"],
        queryFn: async () => {
            const { data } = await axios.get(`/config/BackupTime`);
            return data;
        },
        enabled: true,
    });

    // Cron jobs query
    const { data: cronJobsData, isLoading: cronJobsLoading, refetch: refetchCronJobs } = useQuery({
        queryKey: ["cron", "jobs"],
        queryFn: async () => {
            const { data } = await axios.get(`/cron/jobs`);
            return data;
        },
        enabled: true,
        refetchInterval: 10000, // Auto-refresh every 10 seconds
    });

    const [simrsBridgingActive, setSimrsBridgingActive] = useState<boolean>(false);
    const simrsList = [
        { id: "khanza", label: "Khanza" },
        { id: "simrs-api", label: "SIMRS (API)" },
        { id: "simgos", label: "SIMRS (Database Sharing)" },
        { id: "technomedic", label: "TechnoMedic (API)" },
        { id: "softmedix", label: "Softmedix" },
        { id: "simrs-local", label: "Local SIMRS" },
    ];

    const [selectedSimrs, setSelectedSimrs] = useState<string>("");
    const [isDirty, setIsDirty] = useState<boolean>(false);
    const [isSaving, setIsSaving] = useState<boolean>(false);
    const [isSavingBackup, setIsSavingBackup] = useState<boolean>(false);
    const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
    const [snackbarMessage, setSnackbarMessage] = useState<string>("");
    const [snackbarSeverity, setSnackbarSeverity] = useState<"success" | "error" | "info">("success");
    const [initialSnapshot, setInitialSnapshot] = useState<null | {
        simrsBridgingActive: boolean;
        selectedSimrs: string;
        khanzaMethod: string;
        bridge: { user: string; pass: string; host: string; port: string; db: string; params: string };
        main: { user: string; pass: string; host: string; port: string; db: string; params: string };
        simrs: { dsn: string };
        simgos: { dsn: string };
        backup: { scheduleType: string; interval: string; time: string };
    }>(null);
    const [khanzaMethod, setKhanzaMethod] = useState<string>("api");

    // Khanza states
    const [bridgeUser, setBridgeUser] = useState<string>("");
    const [bridgePassword, setBridgePassword] = useState<string>("");
    const [bridgeHost, setBridgeHost] = useState<string>("");
    const [bridgePort, setBridgePort] = useState<string>("");
    const [bridgeDb, setBridgeDb] = useState<string>("");
    const [bridgeParams, setBridgeParams] = useState<string>("");
    const [mainUser, setMainUser] = useState<string>("");
    const [mainPassword, setMainPassword] = useState<string>("");
    const [mainHost, setMainHost] = useState<string>("");
    const [mainPort, setMainPort] = useState<string>("");
    const [mainDb, setMainDb] = useState<string>("");
    const [mainParams, setMainParams] = useState<string>("");

    // SIMRS states
    const [simrsDsn, setSimrsDsn] = useState<string>("");

    // SIMGOS states
    const [simgosDsn, setSimgosDsn] = useState<string>("");

    // Backup configuration states
    const [backupScheduleType, setBackupScheduleType] = useState<string>("daily");
    const [backupInterval, setBackupInterval] = useState<string>("6");
    const [backupTime, setBackupTime] = useState<string>("02:00");

    // Cron jobs state
    const [cronJobs, setCronJobs] = useState<any[]>([]);
    const [loadingCrons, setLoadingCrons] = useState<{ [key: string]: boolean }>({});

    // Separate dirty states
    const [isBackupDirty, setIsBackupDirty] = useState(false);

    const markDirty = () => setIsDirty(true);
    const markBackupDirty = () => setIsBackupDirty(true);

    const closeSnackbar = () => setSnackbarOpen(false);

    // Convert cron schedule to human-readable format
    const formatCronSchedule = (schedule: string): string => {
        try {
            const parts = schedule.split(' ');
            if (parts.length !== 6) return schedule;

            const [second, minute, hour, dayOfMonth, month] = parts;

            // Check for interval patterns
            if (hour.startsWith('*/')) {
                const hours = hour.substring(2);
                return `Every ${hours} hour${hours !== '1' ? 's' : ''}`;
            }

            // Check for specific time patterns
            if (second === '0' && minute !== '*' && hour !== '*') {
                const h = parseInt(hour);
                const m = parseInt(minute);
                const time = `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}`;

                if (dayOfMonth === '*' && month === '*') {
                    return `Daily at ${time}`;
                }
            }

            // Check for every N seconds/minutes
            if (second.startsWith('*/')) {
                const seconds = second.substring(2);
                return `Every ${seconds} second${seconds !== '1' ? 's' : ''}`;
            }

            if (minute.startsWith('*/') && hour === '*') {
                const minutes = minute.substring(2);
                return `Every ${minutes} minute${minutes !== '1' ? 's' : ''}`;
            }

            return schedule;
        } catch (e) {
            return schedule;
        }
    };

    const parseKhanzaDSN = (dsn: string) => {
        try {
            const result: any = {
                user: "",
                pass: "",
                host: "",
                port: "",
                db: "",
                params: "",
            };

            if (!dsn) return result;

            const [main, params] = dsn.split("?");
            result.params = params || "";

            const slashIndex = main.lastIndexOf("/");
            if (slashIndex !== -1) {
                result.db = main.substring(slashIndex + 1);
            }

            const beforeDb = slashIndex !== -1 ? main.substring(0, slashIndex) : main;

            const atIndex = beforeDb.indexOf("@");
            if (atIndex !== -1) {
                const up = beforeDb.substring(0, atIndex);
                const [u, p] = up.split(":");
                result.user = u || "";
                result.pass = p || "";
            }

            const lp = beforeDb.indexOf("(");
            const rp = beforeDb.indexOf(")");
            if (lp !== -1 && rp !== -1 && rp > lp) {
                const hostport = beforeDb.substring(lp + 1, rp);
                const [h, pt] = hostport.split(":");
                result.host = h || "";
                result.port = pt || "";
            } else {
                const afterProto = beforeDb.split(")").pop() || beforeDb;
                if (afterProto) {
                    const hp = afterProto.replace(/^[^0-9a-zA-Z]*/, "");
                    const [h, pt] = hp.split(":");
                    result.host = h || "";
                    result.port = pt || "";
                }
            }

            return result;
        } catch (e) {
            return { user: "", pass: "", host: "", port: "", db: "", params: "" };
        }
    };

    const composeKhanzaDSN = (parts: { user: string; pass: string; host: string; port: string; db: string; params: string }) => {
        const userpart = parts.user ? parts.user : "";
        const passpart = parts.pass ? ":" + parts.pass : "";
        const auth = userpart || passpart ? `${userpart}${passpart}@` : "";
        const proto = "tcp";
        const hostport = parts.port ? `${parts.host}:${parts.port}` : parts.host;
        const dbpart = parts.db ? `/${parts.db}` : "";
        const params = parts.params ? `?${parts.params}` : "";
        return `${auth}${proto}(${hostport})${dbpart}${params}`;
    };


    useEffect(() => {
        let bridgingEnabled = false;
        let selectedSimrsType = "";

        // Check if any integration is enabled
        if (configEntry && (configEntry as any).value === "true") {
            bridgingEnabled = true;
            selectedSimrsType = "khanza"; // If KhanzaIntegrationEnabled is true, Khanza is selected
        }

        if (simrsEnabledEntry && (simrsEnabledEntry as any).value === "true") {
            bridgingEnabled = true;
            selectedSimrsType = "simrs"; // If SimrsIntegrationEnabled is true, SIMRS is selected
        }

        if (simgosEnabledEntry && (simgosEnabledEntry as any).value === "true") {
            bridgingEnabled = true;
            selectedSimrsType = "simgos"; // If SimgosIntegrationEnabled is true, SIMGOS is selected
        }

        if (technomedicEnabledEntry && (technomedicEnabledEntry as any).value === "true") {
            bridgingEnabled = true;
            selectedSimrsType = "technomedic"; // If TechnoMedicIntegrationEnabled is true, TechnoMedic is selected
        }

        // Override with SelectedSimrs if available
        if (selectedEntry && (selectedEntry as any).value !== undefined && (selectedEntry as any).value !== "" && (selectedEntry as any).value !== "none") {
            selectedSimrsType = (selectedEntry as any).value;
            bridgingEnabled = true;
        } else if (selectedEntry && (selectedEntry as any).value === "none") {
            // If SelectedSimrs is "none", bridging is disabled
            bridgingEnabled = false;
            selectedSimrsType = "";
        }

        setSimrsBridgingActive(bridgingEnabled);
        setSelectedSimrs(selectedSimrsType);

        // Load Khanza configuration
        if (khanzaBridgeEntry && (khanzaBridgeEntry as any).value !== undefined) {
            const d = (khanzaBridgeEntry as any).value;
            const parts = parseKhanzaDSN(d);
            setBridgeUser(parts.user);
            setBridgePassword(parts.pass);
            setBridgeHost(parts.host);
            setBridgePort(parts.port);
            setBridgeDb(parts.db);
            setBridgeParams(parts.params);
        }
        if (khanzaMainEntry && (khanzaMainEntry as any).value !== undefined) {
            const dm = (khanzaMainEntry as any).value;
            const mparts = parseKhanzaDSN(dm);
            setMainUser(mparts.user);
            setMainPassword(mparts.pass);
            setMainHost(mparts.host);
            setMainPort(mparts.port);
            setMainDb(mparts.db);
            setMainParams(mparts.params);
        }
        if (khanzaMethodEntry && (khanzaMethodEntry as any).value !== undefined) {
            setKhanzaMethod((khanzaMethodEntry as any).value || "api");
        }

        // Load SIMRS config
        const currentSimrsDsn = (simrsDsnEntry as any)?.value || "";
        if (simrsDsnEntry && (simrsDsnEntry as any).value !== undefined) {
            setSimrsDsn((simrsDsnEntry as any).value);
        }

        // Load SIMGOS config
        const currentSimgosDsn = (simgosDsnEntry as any)?.value || "";
        if (simgosDsnEntry && (simgosDsnEntry as any).value !== undefined) {
            setSimgosDsn((simgosDsnEntry as any).value);
        }

        // Load Backup configuration
        if (backupScheduleTypeEntry && (backupScheduleTypeEntry as any).value !== undefined) {
            setBackupScheduleType((backupScheduleTypeEntry as any).value || "daily");
        }
        if (backupIntervalEntry && (backupIntervalEntry as any).value !== undefined) {
            setBackupInterval((backupIntervalEntry as any).value || "6");
        }
        if (backupTimeEntry && (backupTimeEntry as any).value !== undefined) {
            setBackupTime((backupTimeEntry as any).value || "02:00");
        }

        setInitialSnapshot({
            simrsBridgingActive: bridgingEnabled,
            selectedSimrs: selectedSimrsType,
            khanzaMethod: (khanzaMethodEntry as any)?.value || "api",
            bridge: {
                user: bridgeUser,
                pass: bridgePassword,
                host: bridgeHost,
                port: bridgePort,
                db: bridgeDb,
                params: bridgeParams,
            },
            main: {
                user: mainUser,
                pass: mainPassword,
                host: mainHost,
                port: mainPort,
                db: mainDb,
                params: mainParams,
            },
            simrs: {
                dsn: currentSimrsDsn,
            },
            simgos: {
                dsn: currentSimgosDsn,
            },
            backup: {
                scheduleType: (backupScheduleTypeEntry as any)?.value || "daily",
                interval: (backupIntervalEntry as any)?.value || "6",
                time: (backupTimeEntry as any)?.value || "02:00",
            },
        });
    }, [configEntry, simrsEnabledEntry, simgosEnabledEntry, selectedEntry, khanzaBridgeEntry, khanzaMainEntry, khanzaMethodEntry, simrsDsnEntry, simgosDsnEntry, backupScheduleTypeEntry, backupIntervalEntry, backupTimeEntry]);

    // Track SIMRS Bridging configuration changes
    useEffect(() => {
        if (!initialSnapshot) return;
        const same =
            initialSnapshot.simrsBridgingActive === simrsBridgingActive &&
            initialSnapshot.selectedSimrs === selectedSimrs &&
            initialSnapshot.khanzaMethod === khanzaMethod &&
            initialSnapshot.bridge.user === bridgeUser &&
            initialSnapshot.bridge.pass === bridgePassword &&
            initialSnapshot.bridge.host === bridgeHost &&
            initialSnapshot.bridge.port === bridgePort &&
            initialSnapshot.bridge.db === bridgeDb &&
            initialSnapshot.bridge.params === bridgeParams &&
            initialSnapshot.main.user === mainUser &&
            initialSnapshot.main.pass === mainPassword &&
            initialSnapshot.main.host === mainHost &&
            initialSnapshot.main.port === mainPort &&
            initialSnapshot.main.db === mainDb &&
            initialSnapshot.main.params === mainParams &&
            initialSnapshot.simrs.dsn === simrsDsn;

        setIsDirty(!same);
    }, [initialSnapshot, simrsBridgingActive, selectedSimrs, khanzaMethod, bridgeUser, bridgePassword, bridgeHost, bridgePort, bridgeDb, bridgeParams, mainUser, mainPassword, mainHost, mainPort, mainDb, mainParams, simrsDsn]);

    // Track Backup configuration changes separately
    useEffect(() => {
        if (!initialSnapshot) return;
        const same =
            initialSnapshot.backup.scheduleType === backupScheduleType &&
            initialSnapshot.backup.interval === backupInterval &&
            initialSnapshot.backup.time === backupTime;

        setIsBackupDirty(!same);
    }, [initialSnapshot, backupScheduleType, backupInterval, backupTime]);

    // Update cronJobs state when data is loaded
    useEffect(() => {
        if (cronJobsData?.data) {
            setCronJobs(cronJobsData.data);
        }
    }, [cronJobsData]);

    // Handler for toggling cron job
    const handleToggleCronJob = async (jobName: string, currentActive: boolean) => {
        // Optimistic update
        setLoadingCrons(prev => ({ ...prev, [jobName]: true }));

        try {
            const endpoint = currentActive
                ? `/cron/jobs/${jobName}/disable`
                : `/cron/jobs/${jobName}/enable`;

            await axios.post(endpoint);

            // Update local state
            setCronJobs(prevJobs =>
                prevJobs.map(job =>
                    job.name === jobName
                        ? { ...job, active: !currentActive }
                        : job
                )
            );

            setSnackbarSeverity("success");
            setSnackbarMessage(`Job ${currentActive ? 'disabled' : 'enabled'} successfully`);
            setSnackbarOpen(true);
        } catch (err) {
            setSnackbarSeverity("error");
            setSnackbarMessage(`Failed to ${currentActive ? 'disable' : 'enable'} job`);
            setSnackbarOpen(true);
            console.error(err);
        } finally {
            setLoadingCrons(prev => ({ ...prev, [jobName]: false }));
        }
    };

    const hasRequiredFilled = (() => {
        if (!simrsBridgingActive) return true;
        if (!selectedSimrs) return false;
        if (selectedSimrs === "khanza") {
            if (!khanzaMethod) return false;
            if (!mainUser || !mainPassword || !mainHost || !mainPort || !mainDb) return false;
            if (khanzaMethod === "db") {
                if (!bridgeUser || !bridgePassword || !bridgeHost || !bridgePort || !bridgeDb) return false;
            }
        }
        if (selectedSimrs === "simrs") {
            if (!simrsDsn) return false;
        }
        if (selectedSimrs === "simrs-api") {
            // No required fields for SIMRS API - just enabled
            return true;
        }
        if (selectedSimrs === "simgos") {
            if (!simgosDsn) return false;
        }
        return true;
    })();

    return (
        <Stack gap={4}>
            {/* Server Info */}
            <Card sx={{ p: 1 }}>
                <CardContent>
                    <Stack gap={0.5}>
                        <Stack direction="row" alignItems="center" gap={1}>
                            <StorageIcon color="primary" />
                            <Typography variant="h6">Server Info</Typography>
                        </Stack>
                        <Divider sx={{ mt: 2 }} />

                        {isPending ? (
                            <Stack direction="row" alignItems="center" gap={2}>
                                <CircularProgress size={24} />
                                <Typography variant="body2" color="text.secondary">
                                    Loading server info...
                                </Typography>
                            </Stack>
                        ) : (
                            <Stack gap={0.5}>
                                <Stack direction="column" gap={0.5} sx={{ borderBottom: "1px solid #e3e3e3ff", pb: 0.5 }}>
                                    <Stack direction="row" alignItems="center" gap={1}>
                                        <PublicIcon color="action" fontSize="small" />
                                        <Typography variant="subtitle2" color="text.secondary">Server IP</Typography>
                                    </Stack>
                                    <Typography variant="body1" sx={{ ml: 3 }}>{data.serverIP}</Typography>
                                </Stack>

                                <Stack direction="column" gap={0.5} sx={{ borderBottom: "1px solid #e3e3e3ff", pb: 0.5 }}>
                                    <Stack direction="row" alignItems="center" gap={1}>
                                        <SettingsEthernetIcon color="action" fontSize="small" />
                                        <Typography variant="subtitle2" color="text.secondary">Port</Typography>
                                    </Stack>
                                    <Typography variant="body1" sx={{ ml: 3 }}>{data.port}</Typography>
                                </Stack>

                                <Stack direction="column" gap={0.5} sx={{ borderBottom: "1px solid #e3e3e3ff", pb: 0.5 }}>
                                    <Stack direction="row" alignItems="center" gap={1}>
                                        <HistoryIcon color="action" fontSize="small" />
                                        <Typography variant="subtitle2" color="text.secondary">Revision</Typography>
                                    </Stack>
                                    <Typography variant="body1" sx={{ ml: 3 }}>{data.revision}</Typography>
                                </Stack>

                                <Stack direction="column" gap={0.5} sx={{ borderBottom: "1px solid #e3e3e3ff", pb: 0.5 }}>
                                    <Stack direction="row" alignItems="center" gap={1}>
                                        <InfoIcon color="action" fontSize="small" />
                                        <Typography variant="subtitle2" color="text.secondary">Version</Typography>
                                    </Stack>
                                    <Typography variant="body1" sx={{ ml: 3 }}>{data.version}</Typography>
                                </Stack>
                            </Stack>
                        )}
                    </Stack>
                </CardContent>
            </Card>

            {/* Backup Configuration Section */}
            <Card sx={{ p: 2 }}>
                <Stack direction="row" alignItems="center" gap={1} sx={{ mb: 2 }}>
                    <StorageIcon color="primary" />
                    <Typography variant="h6">Database Backup Configuration</Typography>
                </Stack>
                <Divider sx={{ mb: 2 }} />

                <Stack gap={2}>
                    <FormControl size="small" sx={{ minWidth: 240 }}>
                        <InputLabel id="backup-schedule-type-label">Backup Schedule Type</InputLabel>
                        <Select
                            labelId="backup-schedule-type-label"
                            value={backupScheduleType}
                            label="Backup Schedule Type"
                            onChange={(e) => { setBackupScheduleType(e.target.value); markBackupDirty(); }}
                        >
                            <MenuItem value="daily">Daily at specific time</MenuItem>
                            <MenuItem value="interval">Every N hours</MenuItem>
                        </Select>
                    </FormControl>

                    {backupScheduleType === "interval" && (
                        <Stack direction="column" gap={1}>
                            <Typography variant="subtitle2" color="text.secondary">
                                Backup Interval (hours)
                            </Typography>
                            <FormControl size="small" sx={{ minWidth: 200 }}>
                                <InputLabel id="backup-interval-label">Hours</InputLabel>
                                <Select
                                    labelId="backup-interval-label"
                                    value={backupInterval}
                                    label="Hours"
                                    onChange={(e) => { setBackupInterval(e.target.value); markBackupDirty(); }}
                                >
                                    <MenuItem value="1">Every 1 hour</MenuItem>
                                    <MenuItem value="2">Every 2 hours</MenuItem>
                                    <MenuItem value="3">Every 3 hours</MenuItem>
                                    <MenuItem value="4">Every 4 hours</MenuItem>
                                    <MenuItem value="6">Every 6 hours</MenuItem>
                                    <MenuItem value="8">Every 8 hours</MenuItem>
                                    <MenuItem value="12">Every 12 hours</MenuItem>
                                    <MenuItem value="24">Every 24 hours</MenuItem>
                                </Select>
                            </FormControl>
                            <Typography variant="caption" color="text.secondary">
                                Database will be backed up every {backupInterval} hour(s)
                            </Typography>
                        </Stack>
                    )}

                    {backupScheduleType === "daily" && (
                        <Stack direction="column" gap={1}>
                            <Typography variant="subtitle2" color="text.secondary">
                                Daily Backup Time
                            </Typography>
                            <MUITextField
                                type="time"
                                value={backupTime}
                                size="small"
                                onChange={(e) => { setBackupTime(e.target.value); markBackupDirty(); }}
                                sx={{ maxWidth: 200 }}
                                InputLabelProps={{
                                    shrink: true,
                                }}
                            />
                            <Typography variant="caption" color="text.secondary">
                                Database will be backed up daily at {backupTime}
                            </Typography>
                        </Stack>
                    )}

                    <Typography variant="body2" color="warning.main" sx={{ display: 'block' }}>
                        ⚠️ Backup files are stored in %LOCALAPPDATA%\TAMALabs\backup. Ensure sufficient disk space is available.
                    </Typography>

                    <Button
                        variant="contained"
                        color="primary"
                        disabled={!isBackupDirty || isSavingBackup}
                        onClick={async () => {
                            setIsSavingBackup(true);
                            try {
                                // Save Backup configuration
                                await axios.put(`/config/BackupScheduleType`, {
                                    id: "BackupScheduleType",
                                    value: backupScheduleType,
                                });
                                await axios.put(`/config/BackupInterval`, {
                                    id: "BackupInterval",
                                    value: backupInterval,
                                });
                                await axios.put(`/config/BackupTime`, {
                                    id: "BackupTime",
                                    value: backupTime,
                                });

                                queryClient.invalidateQueries({ queryKey: ["config", "BackupScheduleType"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "BackupInterval"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "BackupTime"] });
                                queryClient.invalidateQueries({ queryKey: ["cron", "jobs"] }); // Refresh cron jobs to show updated schedule

                                // Update backup schedule in cron manager
                                try {
                                    await axios.post(`/cron/backup/update-schedule`);
                                } catch (err) {
                                    console.warn("Failed to update backup schedule in cron, will apply on restart", err);
                                }

                                const newSnapshot = {
                                    ...initialSnapshot!,
                                    backup: {
                                        scheduleType: backupScheduleType,
                                        interval: backupInterval,
                                        time: backupTime,
                                    },
                                };
                                setInitialSnapshot(newSnapshot);
                                setSnackbarSeverity("success");
                                setSnackbarMessage("Backup configuration saved successfully. Schedule updated.");
                                setSnackbarOpen(true);
                                setIsBackupDirty(false);
                            } catch (err) {
                                setSnackbarSeverity("error");
                                setSnackbarMessage("Failed to save backup configuration");
                                setSnackbarOpen(true);
                            } finally {
                                setIsSavingBackup(false);
                            }
                        }}
                    >
                        {isSavingBackup ? "Saving..." : "Save Backup Configuration"}
                    </Button>
                </Stack>
            </Card>

            {/* Config List */}
            <Card sx={{ p: 2 }}>
                <Stack direction="row" alignItems="center" gap={1}>
                    <IntegrationInstructionsRounded color="primary" />
                    <Typography variant="h6">SIMRS Bridging</Typography>
                </Stack>
                <Divider sx={{ my: 2 }} />
                {/* SIMRS bridging toggle */}
                <Stack direction="row" alignItems="center" gap={2} sx={{ display: "flex", flexDirection: "column", alignItems: "start" }}>
                    <FormControlLabel
                        control={<Switch checked={simrsBridgingActive} onChange={() => setSimrsBridgingActive(s => !s)} color="primary" />}
                        label="SIMRS Bridging"
                    />

                    {simrsBridgingActive && (
                        <FormControl size="small" sx={{ minWidth: 200 }}>
                            <InputLabel id="simrs-select-label">SIMRS</InputLabel>
                            <Select
                                labelId="simrs-select-label"
                                value={selectedSimrs}
                                label="SIMRS"
                                onChange={(e) => { setSelectedSimrs(e.target.value); markDirty(); }}
                                required
                            >
                                <MenuItem value="" disabled>
                                    Select one
                                </MenuItem>
                                {simrsList.map((s) => (
                                    <MenuItem key={s.id} value={s.id}>
                                        {s.label}
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    )}

                    {simrsBridgingActive && selectedSimrs === "khanza" && (
                        <FormControl size="small" sx={{ minWidth: 240 }}>
                            <InputLabel id="khanza-method-label">Connection Method</InputLabel>
                            <Select
                                labelId="khanza-method-label"
                                value={khanzaMethod}
                                label="Connection Method"
                                onChange={(e) => { setKhanzaMethod(e.target.value); markDirty(); }}
                                required
                            >
                                <MenuItem value="api">API</MenuItem>
                                <MenuItem value="db">DB Sharing</MenuItem>
                            </Select>
                        </FormControl>
                    )}

                    {simrsBridgingActive && selectedSimrs === "khanza" && khanzaMethod === "db" && (

                        <Stack direction="column" gap={1}>
                            <Typography variant="body1" color="primary">Khanza Bridging DB</Typography>

                            <Stack direction="row" gap={1} alignItems="flex-start">
                                <Stack direction="column" gap={1}>
                                    <MUITextField
                                        label="User"
                                        value={bridgeUser}
                                        size="small"
                                        onChange={(e) => { setBridgeUser(e.target.value); markDirty(); }}
                                        required
                                    />
                                    <MUITextField
                                        label="Password"
                                        value={bridgePassword}
                                        size="small"
                                        onChange={(e) => { setBridgePassword(e.target.value); markDirty(); }}
                                        type="password"
                                        required
                                    />
                                </Stack>

                                <Stack direction="column" gap={1}>
                                    <MUITextField
                                        label="Host"
                                        value={bridgeHost}
                                        size="small"
                                        onChange={(e) => { setBridgeHost(e.target.value); markDirty(); }}
                                        required
                                    />
                                    <MUITextField
                                        label="Port"
                                        value={bridgePort}
                                        size="small"
                                        onChange={(e) => { setBridgePort(e.target.value); markDirty(); }}
                                        required
                                    />
                                </Stack>

                                <Stack direction="column" gap={1} sx={{ minWidth: 240 }}>
                                    <MUITextField
                                        label="Database"
                                        value={bridgeDb}
                                        size="small"
                                        onChange={(e) => { setBridgeDb(e.target.value); markDirty(); }}
                                        required
                                    />
                                    <MUITextField
                                        label="Params"
                                        value={bridgeParams}
                                        size="small"
                                        onChange={(e) => { setBridgeParams(e.target.value); markDirty(); }}
                                    />
                                </Stack>
                            </Stack>

                        </Stack>
                    )}

                    {simrsBridgingActive && selectedSimrs === "khanza" && (
                        khanzaMethod === "api" ? (
                            <Stack direction="column" gap={2} alignItems="flex-start" style={{ width: "100%" }}>
                                <Stack direction="column" gap={1} sx={{ flex: 1 }}>
                                    <Typography variant="body1" color="primary">Khanza Main DB</Typography>
                                    <Stack direction="row" gap={1} alignItems="flex-start">
                                        <Stack direction="column" gap={1}>
                                            <MUITextField
                                                label="User"
                                                value={mainUser}
                                                size="small"
                                                onChange={(e) => { setMainUser(e.target.value); markDirty(); }}
                                                required
                                            />
                                            <MUITextField
                                                label="Password"
                                                value={mainPassword}
                                                size="small"
                                                onChange={(e) => { setMainPassword(e.target.value); markDirty(); }}
                                                type="password"
                                                required
                                            />
                                        </Stack>

                                        <Stack direction="column" gap={1}>
                                            <MUITextField
                                                label="Host"
                                                value={mainHost}
                                                size="small"
                                                onChange={(e) => { setMainHost(e.target.value); markDirty(); }}
                                                required
                                            />
                                            <MUITextField
                                                label="Port"
                                                value={mainPort}
                                                size="small"
                                                onChange={(e) => { setMainPort(e.target.value); markDirty(); }}
                                                required
                                            />
                                        </Stack>

                                        <Stack direction="column" gap={1} sx={{ minWidth: 240 }}>
                                            <MUITextField
                                                label="Database"
                                                value={mainDb}
                                                size="small"
                                                onChange={(e) => { setMainDb(e.target.value); markDirty(); }}
                                                required
                                            />
                                            <MUITextField
                                                label="Params"
                                                value={mainParams}
                                                size="small"
                                                onChange={(e) => setMainParams(e.target.value)}
                                            />
                                        </Stack>
                                    </Stack>
                                </Stack>

                                <Card sx={{ width: "100%", border: "1px solid #e6e6e6ff" }} elevation={0}>
                                    <Accordion>
                                        <AccordionSummary>
                                            <Typography variant="subtitle1">API Documentation (click to expand)</Typography>
                                        </AccordionSummary>
                                        <AccordionDetails>
                                            <CardContent>
                                                <Typography variant="subtitle1" gutterBottom>API Documentation</Typography>
                                                <Divider sx={{ mb: 1 }} />
                                                <Typography sx={{ mb: 1 }} variant="subtitle1">Server {`${data.serverIP}:${data.port}`}</Typography>
                                                <Typography variant="subtitle2">POST /api/v1/khanza/order</Typography>
                                                <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                    {`Request example:
{
    "order": {
        "pid": {
            "pname": "Budi",
            "sex": "M",
            "birth_dt": "21.07.1985"
        },
        "obr": {
            "order_lab": "LAB001",
            "order_test": ["HGB","WBC"]
        }
    }
}
`}
                                                </Typography>

                                                <Typography variant="subtitle2" sx={{ mt: 2 }}>GET /api/v1/khanza/result/:user/:key/:id</Typography>
                                                <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                    {`Response example:
{
    "response": {
        "sampel": {
            "result_test": [
                {
                    "test_id": "123",
                    "nama_test": "Hemoglobin",
                    "hasil": "13.5",
                    "nilai_normal": "12-16",
                    "satuan": "g/dL",
                    "flag": ""
                }
                // dst
            ]
        }
    },
    "result": {
        "obx": {
            "order_lab": "ONO123456"
        }
    }
}`}
                                                </Typography>
                                            </CardContent>
                                        </AccordionDetails>
                                    </Accordion>
                                </Card>
                            </Stack>
                        ) : (
                            <Stack direction="column" gap={1} >
                                <Typography variant="body1" color="primary">Khanza Main DB</Typography>
                                <Stack direction="row" gap={1} alignItems="flex-start">
                                    <Stack direction="column" gap={1}>
                                        <MUITextField
                                            label="User"
                                            value={mainUser}
                                            size="small"
                                            onChange={(e) => { setMainUser(e.target.value); markDirty(); }}
                                            required
                                        />
                                        <MUITextField
                                            label="Password"
                                            value={mainPassword}
                                            size="small"
                                            onChange={(e) => { setMainPassword(e.target.value); markDirty(); }}
                                            type="password"
                                            required
                                        />
                                    </Stack>

                                    <Stack direction="column" gap={1}>
                                        <MUITextField
                                            label="Host"
                                            value={mainHost}
                                            size="small"
                                            onChange={(e) => { setMainHost(e.target.value); markDirty(); }}
                                            required
                                        />
                                        <MUITextField
                                            label="Port"
                                            value={mainPort}
                                            size="small"
                                            onChange={(e) => { setMainPort(e.target.value); markDirty(); }}
                                            required
                                        />
                                    </Stack>

                                    <Stack direction="column" gap={1} sx={{ minWidth: 240 }}>
                                        <MUITextField
                                            label="Database"
                                            value={mainDb}
                                            size="small"
                                            onChange={(e) => { setMainDb(e.target.value); markDirty(); }}
                                            required
                                        />
                                        <MUITextField
                                            label="Params"
                                            value={mainParams}
                                            size="small"
                                            onChange={(e) => { setMainParams(e.target.value); markDirty(); }}
                                        />
                                    </Stack>
                                </Stack>

                            </Stack>
                        )
                    )}

                    {/* SIMRS API Configuration */}
                    {simrsBridgingActive && selectedSimrs === "simrs-api" && (
                        <Stack direction="column" gap={2} alignItems="flex-start" style={{ width: "100%" }}>
                            <Card sx={{ width: "100%", border: "1px solid #e6e6e6ff" }} elevation={0}>
                                <Accordion defaultExpanded={false}>
                                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                        <Typography variant="subtitle1">SIMRS API Documentation</Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>SIMRS External API Endpoints</Typography>
                                            <Divider sx={{ mb: 1 }} />
                                            <Typography variant="body2" paragraph>
                                                External API endpoints that SIMRS can call to send lab orders and retrieve results.
                                            </Typography>
                                            <Typography sx={{ mb: 1 }} variant="subtitle1">Server {`${data.serverIP}:${data.port}`}</Typography>

                                            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>POST /api/v1/his/order</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Create a new lab order from SIMRS

Request Body:
{
    "order": {
        "pid": {
            "pname": "Rudi Santoso",
            "sex": "M",
            "birth_dt": "21.07.1985",
            "no_rm": "983929"
        },
        "obr": {
            "order_lab": "LAB20250122001",
            "order_test": [
                206,
                205
            ],
            "doctor": [2],
            "analyst": [3]
        }
    }
}`}
                                            </Typography>

                                            <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>GET /api/v1/his/result/:orderlab</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Retrieve lab results for a specific order

Example: GET /api/v1/his/result/LAB20250122001

Response:
{
    "response": {
        "sampel": {
            "result_test": [
                {
                    "loinc": "302-1",
                    "test_id": 206,
                    "nama_test": "ALBUMIN",
                    "hasil": "10.0",
                    "nilai_normal": "3.5 - 5.0",
                    "satuan": "mg/dl",
                    "flag": "H"
                },
                {
                    "loinc": "2313-1",
                    "test_id": 205,
                    "nama_test": "LDL DIRECT TOOS",
                    "hasil": "99",
                    "nilai_normal": "0 - 100",
                    "satuan": "mg/dL",
                    "flag": "N"
                }
            ]
        }
    },
    "result": {
        "obx": {
            "order_lab": "LAB20250122001"
        }
    }
}`}
                                            </Typography>

                                            <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>DELETE /api/v1/his/order/:orderlab</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Cancel/delete a lab order

Example: DELETE /api/v1/his/order/LAB20250122001

Response:
{
    "status": "success",
    "message": "Order deleted successfully"
}`}
                                            </Typography>

                                            <Typography variant="caption" color="warning.main" sx={{ mt: 2, display: 'block' }}>
                                                ⚠️ Note: These endpoints are only accessible when SIMRS API integration is enabled in the configuration.
                                            </Typography>
                                        </CardContent>
                                    </AccordionDetails>
                                </Accordion>
                            </Card>
                        </Stack>
                    )}

                    {/* Database Sharing Configuration */}
                    {simrsBridgingActive && selectedSimrs === "simgos" && (
                        <Stack direction="column" gap={2} alignItems="flex-start" style={{ width: "100%" }}>
                            <Typography variant="body1" color="primary">Database Sharing Configuration</Typography>
                            <MUITextField
                                label="Database DSN"
                                value={simgosDsn}
                                size="small"
                                fullWidth
                                onChange={(e) => { setSimgosDsn(e.target.value); markDirty(); }}
                                placeholder="username:password@tcp(hostname:port)/database?params"
                                helperText="Format: username:password@tcp(hostname:port)/database"
                                required
                            />

                            <Card sx={{ width: "100%", border: "1px solid #e6e6e6ff" }} elevation={0}>
                                <Accordion>
                                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                        <Typography variant="subtitle1">Database Sharing Integration Documentation</Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>Database Sharing Integration</Typography>
                                            <Divider sx={{ mb: 1 }} />
                                            <Typography variant="body2" paragraph>
                                                This integration uses database sharing to sync lab requests and results with SIMGOS SIMRS.
                                            </Typography>

                                            <Typography variant="subtitle2" gutterBottom>Required Tables:</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`• lab_order - Lab order master data
  - Columns: id, no_lab_order, no_rm, patient_name, birth_date, sex, doctor, analyst, status, created_at
  - Status: NEW, PENDING, LIS_SUCCESS, SIMRS_SUCCESS

• order_detail - Lab order details/test parameters
  - Columns: id, no_lab_order, parameter_code, parameter_name, result_value, unit, reference_range, flag, created_at

Database DSN Examples:
• MySQL: client:LabBridgingLIS001@tcp(192.168.0.10:3306)/lis_bridging
• With params: client:password@tcp(host:3306)/lis_bridging?charset=utf8mb4&parseTime=True&loc=Local`}
                                            </Typography>

                                            <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>Status Flow:</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap' }}>
                                                {`NEW → PENDING → LIS_SUCCESS → SIMRS_SUCCESS

1. NEW: Order inserted by SIMRS
2. PENDING: Order fetched by LIS
3. LIS_SUCCESS: LIS completed tests
4. SIMRS_SUCCESS: SIMRS retrieved results`}
                                            </Typography>

                                            <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>Sync Process:</Typography>
                                            <Typography variant="body2">
                                                • Automatic sync runs every few seconds<br />
                                                • LIS fetches NEW orders and updates to PENDING<br />
                                                • LIS updates test results in order_detail table<br />
                                                • When all tests complete, status updates to LIS_SUCCESS<br />
                                                • SIMRS can then fetch results and update to SIMRS_SUCCESS
                                            </Typography>
                                        </CardContent>
                                    </AccordionDetails>
                                </Accordion>
                            </Card>
                        </Stack>
                    )}

                    {/* TechnoMedic API Configuration */}
                    {simrsBridgingActive && selectedSimrs === "technomedic" && (
                        <Stack direction="column" gap={2} alignItems="flex-start" style={{ width: "100%" }}>
                            <Card sx={{ width: "100%", border: "1px solid #e6e6e6ff" }} elevation={0}>
                                <Accordion defaultExpanded={true}>
                                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                        <Typography variant="subtitle1">TechnoMedic API Documentation</Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>TechnoMedic External API Endpoints</Typography>
                                            <Divider sx={{ mb: 1 }} />
                                            <Typography variant="body2" paragraph>
                                                External API endpoints that TechnoMedic can call to manage lab orders and retrieve results.
                                            </Typography>
                                            <Typography sx={{ mb: 1 }} variant="subtitle1">Server {`${data.serverIP}:${data.port}`}</Typography>

                                            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>GET /api/v1/technomedic/test-types</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Get all available test types

Response:
{
    "code": 200,
    "status": "success",
    "data": [
        {
            "id": "1",
            "code": "HB",
            "name": "Hemoglobin",
            "category": "Hematologi",
            "sub_category": "Complete Blood Count",
            "specimen_type": "Whole Blood",
            "unit": "g/dL"
        }
    ]
}`}
                                            </Typography>

                                            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>GET /api/v1/technomedic/sub-categories</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Get all sub-categories

Response:
{
    "code": 200,
    "status": "success",
    "data": [
        {
            "id": "1",
            "code": "CBC",
            "name": "Complete Blood Count",
            "category": "Hematologi"
        }
    ]
}`}
                                            </Typography>

                                            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>POST /api/v1/technomedic/order</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Create a new lab order from TechnoMedic

Request Body:
{
    "no_order": "TM-2024-001",
    "patient": {
        "full_name": "John Doe",
        "sex": "M",
        "birthdate": "1990-01-15",
        "medical_record_number": "MR001"
    },
    "test_type_ids": [1, 2, 3],
    "sub_category_ids": [1]
}

Response:
{
    "code": 201,
    "status": "success",
    "message": "Order created successfully",
    "data": {
        "no_order": "TM-2024-001"
    }
}`}
                                            </Typography>

                                            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>GET /api/v1/technomedic/order/:no_order</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`Get order details with results

Response:
{
    "code": 200,
    "status": "success",
    "data": {
        "no_order": "TM-2024-001",
        "status": "SUCCESS",
        "patient": {...},
        "sub_categories": [...],
        "parameters_result": [...]
    }
}`}
                                            </Typography>

                                            <Typography variant="caption" color="info.main" sx={{ mt: 2, display: 'block' }}>
                                                ℹ️ Note: All endpoints require TechnoMedic integration to be enabled. See full documentation for detailed request/response schemas.
                                            </Typography>
                                        </CardContent>
                                    </AccordionDetails>
                                </Accordion>
                            </Card>
                        </Stack>
                    )}



                    <Button
                        variant="contained"
                        color="primary"
                        disabled={!isDirty || isSaving || !hasRequiredFilled}
                        onClick={async () => {
                            setIsSaving(true);
                            try {
                                // Always save the main bridging status
                                await axios.put(`/config/KhanzaIntegrationEnabled`, {
                                    id: "KhanzaIntegrationEnabled",
                                    value: (simrsBridgingActive && selectedSimrs === "khanza") ? "true" : "false",
                                });

                                // Always save SIMRS integration status
                                await axios.put(`/config/SimrsIntegrationEnabled`, {
                                    id: "SimrsIntegrationEnabled",
                                    value: (simrsBridgingActive && (selectedSimrs === "simrs" || selectedSimrs === "simrs-api")) ? "true" : "false",
                                });

                                // Always save SIMGOS integration status
                                await axios.put(`/config/SimgosIntegrationEnabled`, {
                                    id: "SimgosIntegrationEnabled",
                                    value: (simrsBridgingActive && selectedSimrs === "simgos") ? "true" : "false",
                                });

                                // Always save TechnoMedic integration status
                                await axios.put(`/config/TechnoMedicIntegrationEnabled`, {
                                    id: "TechnoMedicIntegrationEnabled",
                                    value: (simrsBridgingActive && selectedSimrs === "technomedic") ? "true" : "false",
                                });

                                if (simrsBridgingActive) {
                                    await axios.put(`/config/SelectedSimrs`, {
                                        id: "SelectedSimrs",
                                        value: selectedSimrs,
                                    });

                                    if (selectedSimrs === "khanza") {
                                        await axios.put(`/config/KhanzaConnectionMethod`, {
                                            id: "KhanzaConnectionMethod",
                                            value: khanzaMethod,
                                        });

                                        if (khanzaMethod === "db") {
                                            const composed = composeKhanzaDSN({ user: bridgeUser, pass: bridgePassword, host: bridgeHost, port: bridgePort, db: bridgeDb, params: bridgeParams });
                                            await axios.put(`/config/KhanzaBridgeDatabaseDSN`, {
                                                id: "KhanzaBridgeDatabaseDSN",
                                                value: composed,
                                            });
                                        }

                                        const composedMain = composeKhanzaDSN({ user: mainUser, pass: mainPassword, host: mainHost, port: mainPort, db: mainDb, params: mainParams });
                                        await axios.put(`/config/KhanzaMainDatabaseDSN`, {
                                            id: "KhanzaMainDatabaseDSN",
                                            value: composedMain,
                                        });
                                    }

                                    // Save SIMRS configuration (Database Sharing)
                                    if (selectedSimrs === "simrs") {
                                        await axios.put(`/config/SimrsDatabaseDSN`, {
                                            id: "SimrsDatabaseDSN",
                                            value: simrsDsn,
                                        });
                                    }

                                    // SIMRS API doesn't need additional config - just enabled flag

                                    // Save SIMGOS configuration
                                    if (selectedSimrs === "simgos") {
                                        await axios.put(`/config/SimgosDatabaseDSN`, {
                                            id: "SimgosDatabaseDSN",
                                            value: simgosDsn,
                                        });
                                    }
                                } else {
                                    // When bridging is disabled, set SelectedSimrs to "none" instead of empty string
                                    await axios.put(`/config/SelectedSimrs`, {
                                        id: "SelectedSimrs",
                                        value: "none",
                                    });
                                }

                                // Note: Backup configuration is saved separately via "Save Backup Configuration" button

                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaIntegrationEnabled"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SelectedSimrs"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaConnectionMethod"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaMainDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaBridgeDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimrsIntegrationEnabled"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimrsDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimgosIntegrationEnabled"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimgosDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "TechnoMedicIntegrationEnabled"] });

                                // Reload cron jobs to register/unregister based on new config
                                try {
                                    await axios.post(`/cron/reload`);
                                    console.log("Cron jobs reloaded successfully");
                                } catch (err) {
                                    console.warn("Failed to reload cron jobs, will apply on restart", err);
                                }

                                queryClient.invalidateQueries({ queryKey: ["cron", "jobs"] }); // Refresh cron jobs list

                                const newSnapshot = {
                                    simrsBridgingActive,
                                    selectedSimrs: simrsBridgingActive ? selectedSimrs : "",
                                    khanzaMethod,
                                    bridge: {
                                        user: bridgeUser,
                                        pass: bridgePassword,
                                        host: bridgeHost,
                                        port: bridgePort,
                                        db: bridgeDb,
                                        params: bridgeParams,
                                    },
                                    main: {
                                        user: mainUser,
                                        pass: mainPassword,
                                        host: mainHost,
                                        port: mainPort,
                                        db: mainDb,
                                        params: mainParams,
                                    },
                                    simrs: {
                                        dsn: simrsDsn,
                                    },
                                    simgos: {
                                        dsn: simgosDsn,
                                    },
                                    backup: {
                                        scheduleType: backupScheduleType,
                                        interval: backupInterval,
                                        time: backupTime,
                                    },
                                };
                                setInitialSnapshot(newSnapshot);
                                setSnackbarSeverity("success");
                                setSnackbarMessage("SIMRS bridging configuration saved successfully.");
                                setSnackbarOpen(true);
                                setIsDirty(false);
                                // Don't reset isBackupDirty - backup config is saved separately
                            } catch (err) {
                                setSnackbarSeverity("error");
                                setSnackbarMessage("Failed to save configuration");
                                setSnackbarOpen(true);
                            } finally {
                                setIsSaving(false);
                            }
                        }}
                    >
                        {isSaving ? "Saving..." : "Save"}
                    </Button>
                </Stack>
                {/* <CardContent>
                    <Stack gap={2}>
                        <Stack direction="row" alignItems="center" gap={1}>
                            <SettingsIcon color="primary" />
                            <Typography variant="h6">Config</Typography>
                        </Stack>
                        <Divider />

                        <List resource="config" actions={false} pagination={false} sx={{ mt: 1 }} >
                            <Datagrid
                                bulkActionButtons={false}
                                rowClick="edit"
                                sx={{
                                    "& .RaDatagrid-row:hover": {
                                        backgroundColor: "action.hover",
                                        cursor: "pointer",
                                    },
                                }}
                            >
                                <TextField source="id" />
                                <TextField source="value" />
                            </Datagrid>
                        </List>
                    </Stack>
                </CardContent> */}
            </Card>

            {/* Cron Jobs Section */}
            <Card sx={{ p: 2 }}>
                <Stack direction="row" alignItems="center" gap={1} sx={{ mb: 2 }}>
                    <ScheduleIcon color="primary" />
                    <Typography variant="h6">Scheduled Sync Jobs</Typography>
                    <Chip
                        label={cronJobsLoading ? "Loading..." : `${cronJobs.filter((j: any) => j.active).length} Active`}
                        size="small"
                        color="primary"
                        sx={{ ml: 1 }}
                    />
                    <Box sx={{ flexGrow: 1 }} />
                    <Tooltip title="Refresh job list">
                        <IconButton size="small" onClick={() => refetchCronJobs()} disabled={cronJobsLoading}>
                            <RefreshIcon />
                        </IconButton>
                    </Tooltip>
                </Stack>
                <Divider sx={{ mb: 2 }} />

                {cronJobsLoading ? (
                    <Stack direction="row" alignItems="center" gap={2}>
                        <CircularProgress size={24} />
                        <Typography variant="body2" color="text.secondary">
                            Loading scheduled jobs...
                        </Typography>
                    </Stack>
                ) : cronJobs.length === 0 ? (
                    <Typography variant="body2" color="text.secondary">
                        No scheduled jobs found
                    </Typography>
                ) : (
                    <Stack gap={1.5}>
                        {cronJobs.map((job: any) => {
                            const isSimgosJob = job.name === 'sync_all_request_simrs' || job.name === 'sync_all_result_simrs';
                            const isSimgosSelected = selectedSimrs === 'simgos';
                            const isConditionallyInactive = isSimgosJob && (!simrsBridgingActive || !isSimgosSelected);

                            return (
                                <Card key={job.name} variant="outlined" sx={{ p: 2 }}>
                                    <Stack direction="row" alignItems="center" gap={2}>
                                        <Box sx={{ flex: 1 }}>
                                            <Stack direction="row" alignItems="center" gap={1}>
                                                <Typography variant="subtitle1" fontWeight="bold">
                                                    {job.name.replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase())}
                                                </Typography>
                                                <Chip
                                                    label={job.active ? "Active" : "Inactive"}
                                                    size="small"
                                                    color={job.active ? "success" : "default"}
                                                />
                                                {job.name === 'license_heartbeat' && (
                                                    <Chip
                                                        label="Required"
                                                        size="small"
                                                        color="error"
                                                        variant="outlined"
                                                    />
                                                )}
                                                {isSimgosJob && (
                                                    <Chip
                                                        label="Database Sharing Only"
                                                        size="small"
                                                        color="info"
                                                        variant="outlined"
                                                    />
                                                )}
                                            </Stack>
                                            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                                                {job.description}
                                            </Typography>
                                            {isConditionallyInactive && (
                                                <Typography variant="caption" color="warning.main" sx={{ mt: 0.5, display: 'block' }}>
                                                    ⚠️ This job is only active when Database Sharing bridging is enabled
                                                </Typography>
                                            )}
                                            <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
                                                <strong>Schedule:</strong> {formatCronSchedule(job.schedule)}
                                            </Typography>
                                            <Typography variant="caption" color="text.secondary" sx={{ mt: 0.25, display: 'block', fontFamily: 'monospace', fontSize: '0.65rem' }}>
                                                Cron: {job.schedule}
                                            </Typography>
                                        </Box>
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                            {job.name === 'license_heartbeat' ? (
                                                <Tooltip title="License heartbeat is required and cannot be disabled">
                                                    <Box>
                                                        <FormControlLabel
                                                            control={
                                                                <Switch
                                                                    checked={true}
                                                                    disabled={true}
                                                                    color="primary"
                                                                />
                                                            }
                                                            label="Required"
                                                            labelPlacement="start"
                                                        />
                                                    </Box>
                                                </Tooltip>
                                            ) : isSimgosJob ? (
                                                <Tooltip title={isConditionallyInactive ? "Enable SIMGOS bridging to activate this job" : "Toggle job on/off"}>
                                                    <Box>
                                                        <FormControlLabel
                                                            control={
                                                                <Switch
                                                                    checked={job.active}
                                                                    onChange={() => handleToggleCronJob(job.name, job.active)}
                                                                    disabled={loadingCrons[job.name] || isConditionallyInactive}
                                                                    color="primary"
                                                                />
                                                            }
                                                            label={loadingCrons[job.name] ? "..." : job.active ? "On" : "Off"}
                                                            labelPlacement="start"
                                                        />
                                                    </Box>
                                                </Tooltip>
                                            ) : (
                                                <FormControlLabel
                                                    control={
                                                        <Switch
                                                            checked={job.active}
                                                            onChange={() => handleToggleCronJob(job.name, job.active)}
                                                            disabled={loadingCrons[job.name]}
                                                            color="primary"
                                                        />
                                                    }
                                                    label={loadingCrons[job.name] ? "..." : job.active ? "On" : "Off"}
                                                    labelPlacement="start"
                                                />
                                            )}
                                        </Box>
                                    </Stack>
                                </Card>
                            );
                        })}
                    </Stack>
                )}

                <Typography variant="caption" color="text.secondary" sx={{ mt: 2, display: 'block' }}>
                    ℹ️ Sync jobs run automatically in the background. Database Sharing jobs sync data from today only and are active only when Database Sharing bridging is enabled. Jobs marked as "Required" cannot be disabled.
                </Typography>
            </Card>


            <Snackbar open={snackbarOpen} autoHideDuration={4000} onClose={closeSnackbar}>
                <Alert onClose={closeSnackbar} severity={snackbarSeverity} sx={{ width: '100%' }}>
                    {snackbarMessage}
                </Alert>
            </Snackbar>
        </Stack>
    );
};

export const ConfigEdit = () => (
    <Edit >
        <SimpleForm>
            <TextInput source="id" readOnly />
            <TextInput source="value" />
        </SimpleForm>
    </Edit>
);
