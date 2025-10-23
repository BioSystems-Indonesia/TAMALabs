import { Card, CardContent, Typography, Box, useTheme } from "@mui/material";
import {
    ResponsiveContainer,
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Cell,
} from "recharts";
import { CustomTooltip } from "./tooltip";



export function AbnormalCriticalChart({ data = [] }) {
    console.log(data)
    const theme = useTheme();
    const isDark = theme.palette.mode === "dark";

    const COLORS = ["#E8A87C", "#4ABAAB"]
    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom textAlign={"center"} color='gray'>
                    Abnormal / Critical Results Summary (Last 7 Days)
                </Typography>

                <Box sx={{ width: "100%", height: 350 }}>
                    <ResponsiveContainer>
                        <BarChart
                            data={data}
                            margin={{ top: 20, right: 30, left: 10, bottom: 10 }}
                        >
                            <CartesianGrid strokeDasharray="3 3" stroke={isDark ? "#444" : "#e0e0e0"} />
                            <XAxis dataKey="name" tick={{ fill: isDark ? "#ccc" : "#555" }} />
                            <YAxis
                                tickFormatter={formatNumber}
                                tick={{ fill: isDark ? "#ccc" : "#555" }}
                            />
                            <Tooltip content={CustomTooltip} />
                            <Bar dataKey="total" radius={[6, 6, 0, 0]}>
                                {data.map((_, index) => (
                                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </Box>

                <Box sx={{ textAlign: "center", mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Distribution of test results based on normality.
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
}
