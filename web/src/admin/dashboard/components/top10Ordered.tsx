import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    Cell,
} from "recharts";
import { Card, CardContent, Typography, Box } from "@mui/material";
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


type TestData = {
    name: string;
    total: number;
};

export const TopTestsOrdered = ({ data = [] }: { data?: TestData[] }) => {

    const sortedTopTestsData = [...data].sort((a, b) => b.total - a.total);

    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom textAlign={"center"} color='gray'>
                    Top 10 Tests Ordered (Last 7 Days)
                </Typography>

                <Box sx={{ width: "100%", height: 350 }}>
                    <ResponsiveContainer>
                        <BarChart
                            data={sortedTopTestsData}
                            layout="vertical"
                            margin={{ top: 10, right: 30, left: 0, bottom: 10 }}
                        >
                            <XAxis type="number" tickFormatter={(value) => formatNumber(value)} />
                            <YAxis
                                dataKey="name"
                                type="category"
                                width={120}
                                tick={({ x, y, payload }) => {
                                    const text = payload.value.length > 12 ? payload.value.slice(0, 12) + "…" : payload.value;
                                    return (
                                        <text
                                            x={x - 10}
                                            y={y + 4}
                                            textAnchor="end"
                                            fontSize={12}
                                            fill="#555"
                                        >
                                            {text}
                                        </text>
                                    );
                                }}
                            />
                            <Tooltip content={CustomTooltip} />
                            <Bar dataKey="total" radius={[5, 5, 5, 5]}>
                                {sortedTopTestsData.map((_, index) => (
                                    <Cell key={index} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </Box>

                <Box sx={{ textAlign: "center", mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Showing top 10 most frequently ordered lab tests.
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
};
