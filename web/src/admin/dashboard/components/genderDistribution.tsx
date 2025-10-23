import { Card, CardContent, Typography, Box } from "@mui/material";
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts";
import { CustomTooltip } from "./tooltip";

const COLORS = [
    "#4ABAAB",
    "#A7E6D7",
    "#80CBC4",
    "#B2DFDB",
    "#26A69A",
    "#64B5F6",
    "#81C784",
    "#FFD54F",
    "#FFB74D",
    "#E57373",
];

type GenderDataItem = {
    name: string;
    total: number;
};

export function GenderDistribution({ data = [] }: { data?: GenderDataItem[] }) {
    const totalPatients = data.reduce((sum, d) => sum + d.total, 0);

    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    // Label persentase di tengah potongan
    const renderLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent }: any) => {
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
                fontSize={13}
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
                    Gender Distribution
                </Typography>

                <Box sx={{ width: "100%", height: 350 }}>
                    <ResponsiveContainer>
                        <PieChart>
                            <Legend verticalAlign="top" height={36} />

                            <Pie
                                data={data}
                                dataKey="total"
                                nameKey="name"
                                cx="50%"
                                cy="50%"
                                outerRadius={120}
                                labelLine={false}
                                label={renderLabel}
                                stroke="none"
                            >
                                {data.map((_, index) => (
                                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Pie>
                            <Tooltip content={CustomTooltip} />
                        </PieChart>
                    </ResponsiveContainer>
                </Box>

                <Box sx={{ textAlign: "center", mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Total Patients: {formatNumber(totalPatients)}
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
}
