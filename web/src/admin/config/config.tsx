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

    const [simrsBridgingActive, setSimrsBridgingActive] = useState<boolean>(false);
    const simrsList = [
        { id: "khanza", label: "Khanza" },
        { id: "simrs", label: "SIMRS (Database Sharing)" },
        { id: "softmedix", label: "Softmedix" },
        { id: "simrs-local", label: "Local SIMRS" },
    ];

    const [selectedSimrs, setSelectedSimrs] = useState<string>("");
    const [isDirty, setIsDirty] = useState<boolean>(false);
    const [isSaving, setIsSaving] = useState<boolean>(false);
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

    const markDirty = () => setIsDirty(true);

    const closeSnackbar = () => setSnackbarOpen(false);

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

        // Override with SelectedSimrs if available
        if (selectedEntry && (selectedEntry as any).value !== undefined && (selectedEntry as any).value !== "" && (selectedEntry as any).value !== "none") {
            selectedSimrsType = (selectedEntry as any).value;
            bridgingEnabled = true;
        } else if (selectedEntry && (selectedEntry as any).value === "none") {
            // If SelectedSimrs is "none", bridging is disabled
            bridgingEnabled = false;
            selectedSimrsType = "";
        }

        // If no integration is enabled, ensure bridging is disabled
        if (!configEntry || (configEntry as any).value !== "true") {
            if (!simrsEnabledEntry || (simrsEnabledEntry as any).value !== "true") {
                bridgingEnabled = false;
                selectedSimrsType = "";
            }
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
        });
    }, [configEntry, simrsEnabledEntry, selectedEntry, khanzaBridgeEntry, khanzaMainEntry, khanzaMethodEntry, simrsDsnEntry, bridgeUser, bridgePassword, bridgeHost, bridgePort, bridgeDb, bridgeParams, mainUser, mainPassword, mainHost, mainPort, mainDb, mainParams]);

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

                    {/* SIMRS Configuration */}
                    {simrsBridgingActive && selectedSimrs === "simrs" && (
                        <Stack direction="column" gap={2} alignItems="flex-start" style={{ width: "100%" }}>
                            <Typography variant="body1" color="primary">SIMRS Database Configuration</Typography>
                            <MUITextField
                                label="Database DSN"
                                value={simrsDsn}
                                size="small"
                                fullWidth
                                onChange={(e) => { setSimrsDsn(e.target.value); markDirty(); }}
                                placeholder="username:password@tcp(hostname:port)/database?params"
                                helperText="Format: username:password@tcp(hostname:port)/database"
                                required
                            />

                            <Card sx={{ width: "100%", border: "1px solid #e6e6e6ff" }} elevation={0}>
                                <Accordion>
                                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                        <Typography variant="subtitle1">SIMRS Integration Documentation</Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>SIMRS Database Integration</Typography>
                                            <Divider sx={{ mb: 1 }} />
                                            <Typography variant="body2" paragraph>
                                                This integration uses database sharing to sync lab requests and results with SIMRS.
                                            </Typography>

                                            <Typography variant="subtitle2" gutterBottom>Required Tables:</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap', mt: 1 }}>
                                                {`• patients - Patient master data
• lab_requests - Lab test requests from SIMRS
• lab_results - Lab test results to SIMRS

Database DSN Examples:
• MySQL: root:password@tcp(localhost:3306)/simrs_db
• With params: user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local`}
                                            </Typography>

                                            <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>Sync Process:</Typography>
                                            <Typography variant="body2">
                                                • Automatic sync runs every few minutes<br />
                                                • Lab requests are pulled from SIMRS and converted to work orders<br />
                                                • Lab results are pushed back to SIMRS when tests are completed<br />
                                                • Manual sync is available for immediate synchronization
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
                                    value: (simrsBridgingActive && selectedSimrs === "simrs") ? "true" : "false",
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

                                    // Save SIMRS configuration
                                    if (selectedSimrs === "simrs") {
                                        await axios.put(`/config/SimrsDatabaseDSN`, {
                                            id: "SimrsDatabaseDSN",
                                            value: simrsDsn,
                                        });
                                    }
                                } else {
                                    // When bridging is disabled, set SelectedSimrs to "none" instead of empty string
                                    await axios.put(`/config/SelectedSimrs`, {
                                        id: "SelectedSimrs",
                                        value: "none",
                                    });
                                }

                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaIntegrationEnabled"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SelectedSimrs"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaConnectionMethod"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaMainDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "KhanzaBridgeDatabaseDSN"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimrsIntegrationEnabled"] });
                                queryClient.invalidateQueries({ queryKey: ["config", "SimrsDatabaseDSN"] });

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
                                };
                                setInitialSnapshot(newSnapshot);
                                setSnackbarSeverity("success");
                                setSnackbarMessage("Configuration saved");
                                setSnackbarOpen(true);
                                setIsDirty(false);
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
