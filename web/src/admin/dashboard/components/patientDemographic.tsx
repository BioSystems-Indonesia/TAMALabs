import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    CartesianGrid,
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





export const AgeGroupDistribution = ({ data = [] }) => {


    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom textAlign={"center"} color='gray'>
                    Age Group Distribution
                </Typography>

                <Box sx={{ width: "100%", height: 350 }}>
                    <ResponsiveContainer>
                        <BarChart
                            data={data}
                            margin={{ top: 10, right: 30, left: 0, bottom: 10 }}
                        >
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="name" />
                            <YAxis type="number" tickFormatter={(total) => formatNumber(total)} />
                            <Tooltip content={CustomTooltip} />

                            <Bar dataKey="total" radius={[5, 5, 0, 0]}>
                                {data.map((_, index) => (
                                    <Cell key={index} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </Box>

                <Box sx={{ textAlign: "center", mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Patient total by age group.
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
};

