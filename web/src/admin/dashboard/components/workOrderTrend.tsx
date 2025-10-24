import { CardContent, Typography, Box, Card } from "@mui/material"
import { ResponsiveContainer, AreaChart, XAxis, YAxis, Tooltip, Area } from "recharts"
import { CustomTooltip } from "./tooltip";

type WorkOrderTrendData = { date: string; total: number };

export const WorkOrderTrend = ({ data = [] }: { data?: any[] }) => {
    const locale = (typeof navigator !== 'undefined' && navigator.language) ? navigator.language : 'en-US';

    const toDayLabel = (dateStr?: string) => {
        if (!dateStr) return '';
        const parts = String(dateStr).split('-');
        if (parts.length !== 3) return dateStr;
        const y = Number(parts[0]);
        const m = Number(parts[1]) - 1;
        const d = Number(parts[2]);
        const dt = new Date(y, m, d);
        try {
            return dt.toLocaleDateString(locale, { weekday: 'short' });
        } catch {
            return dt.toDateString().slice(0, 3);
        }
    };

    const normalized: (WorkOrderTrendData & { label: string; })[] = (data || []).map((d: any) => {
        const dateStr = d?.date ?? d?.name ?? String(d?.label ?? "");
        const total = typeof d?.total === 'number' ? d.total : Number(d?.total) || 0;
        return { date: dateStr, total, label: toDayLabel(dateStr) };
    });
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
                    <AreaChart data={normalized}>
                        <defs>
                            <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#26a69a" stopOpacity={0.4} />
                                <stop offset="95%" stopColor="#26a69a" stopOpacity={0.05} />
                            </linearGradient>
                        </defs>

                        <XAxis dataKey="label" />
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
                        Total Tests: {formatNumber(normalized.reduce((sum, d) => sum + d.total, 0))}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Avg/Day: {formatNumber(normalized.length ? (normalized.reduce((sum, d) => sum + d.total, 0) / normalized.length) : 0)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {(() => {
                            if (normalized.length === 0) return <>Highest: 0</>;
                            const totals = normalized.map(d => d.total);
                            const highest = Math.max(...totals);
                            const highestDate = normalized.find(d => d.total === highest)?.date ?? '-';
                            return <>Highest: {formatNumber(highest)} ({highestDate})</>;
                        })()}
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    )
}