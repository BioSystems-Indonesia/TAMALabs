import { CardContent, Typography, Box, Card } from "@mui/material"
import { ResponsiveContainer, AreaChart, XAxis, YAxis, Tooltip, Area } from "recharts"
import { CustomTooltip } from "./tooltip";

type WorkOrderTrendData = { date: string; total: number };

export const WorkOrderTrend = ({ data = [] }: { data?: WorkOrderTrendData[] }) => {
    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom textAlign={"center"} color='gray'>
                    Work Order Trend (Last 7 Days)
                </Typography>

                <ResponsiveContainer width="100%" height={350}>
                    <AreaChart data={data}>
                        <defs>
                            <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#26a69a" stopOpacity={0.4} />
                                <stop offset="95%" stopColor="#26a69a" stopOpacity={0.05} />
                            </linearGradient>
                        </defs>

                        <XAxis dataKey="date" />
                        <YAxis
                            tickFormatter={(value) =>
                                formatNumber(value)
                            }
                        />
                        <Tooltip content={CustomTooltip} />

                        <Area
                            type="monotone"
                            dataKey="total"
                            stroke="#26a69a"
                            fill="url(#colorUv)"
                            dot={{ r: 4, fill: "#4ABAAB", strokeWidth: 1, stroke: "#fff" }}
                            strokeWidth={2}
                        />
                    </AreaChart>
                </ResponsiveContainer>

                {/* Chart summary info */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Total Tests: {formatNumber(data.reduce((sum, d) => sum + d.total, 0))}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Avg/Day: {formatNumber((data.reduce((sum, d) => sum + d.total, 0) / data.length))}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Highest: {formatNumber(Math.max(...data.map(d => d.total)))} ({data.find(d => d.total === Math.max(...data.map(x => x.total)))?.date})
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    )
}