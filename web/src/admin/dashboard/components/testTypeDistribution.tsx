import { Card, CardContent, Typography, Box } from "@mui/material";
import {
    PieChart,
    Pie,
    Tooltip,
    Cell,
    Legend,
    ResponsiveContainer,
} from "recharts";
import { CustomTooltip } from "./tooltip";

const COLORS = [
    "#4ABAAB", // main teal
    "#6DC8B8", // light teal
    "#89D9C2", // mint aqua
    "#A7E8BD", // soft mint
    "#BDE6C6", // pale green
    "#F6DFA7", // light sand
    "#F4BFBF", // pastel coral
    "#B4D6FF", // pastel blue
    "#CBB2FE", // lavender purple
    "#E8A87C", // warm beige
];



type TestTypeDataItem = {
    total: number;
    name: string;
};

export function TestTypeDistribution({ data = [] }: { data?: TestTypeDataItem[] }) {
    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    const totalTests = data.reduce((sum, d) => sum + d.total, 0);
    console.log(totalTests)

    // Fungsi untuk menampilkan label persentase
    const renderCustomizedLabel = (props: any) => {
        const { cx, cy, midAngle, innerRadius, outerRadius, percent } = props;
        const RADIAN = Math.PI / 180;
        const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
        const x = cx + radius * Math.cos(-midAngle * RADIAN);
        const y = cy + radius * Math.sin(-midAngle * RADIAN);

        return (
            <text
                x={x}
                y={y}
                fill="#fff"
                textAnchor="middle"
                dominantBaseline="central"
                fontSize={12}
                fontWeight={500}
            >
                {`${(percent * 100).toFixed(0)}%`}
            </text>
        );
    };


    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom textAlign={"center"} color='gray'>
                    Test Type Distribution (Last 7 Days)
                </Typography>

                <Box sx={{ width: "100%", height: 350 }}>
                    <ResponsiveContainer>
                        <PieChart>
                            <Pie
                                data={data}
                                dataKey="total"
                                nameKey="name"
                                innerRadius={70}
                                outerRadius={120}
                                labelLine={false}
                                paddingAngle={1}
                                stroke="none"
                                strokeWidth={1}
                                label={renderCustomizedLabel}
                            >
                                {data.map((_, index) => (
                                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Pie>
                            <Tooltip content={CustomTooltip} />
                            <Legend verticalAlign="top" height={36} />
                        </PieChart>
                    </ResponsiveContainer>
                </Box>

                <Box sx={{ textAlign: "center", mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Total Tests: {formatNumber(totalTests)}
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
}
