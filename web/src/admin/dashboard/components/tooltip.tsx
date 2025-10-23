import { Box, Typography, useTheme } from "@mui/material";
import { TooltipProps } from "recharts";

export const CustomTooltip = (props: TooltipProps<any, any>) => {
    const { active, payload, label } = props as any;
    const theme = useTheme();

    const formatNumber = (num: number): string => {
        if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M";
        if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, "") + "k";
        return num.toString();
    };

    if (!active || !payload?.length) return null;

    const isDark = theme.palette.mode === "dark";
    const item = payload[0];
    const name = item.payload?.name || label || "";

    return (
        <Box
            sx={{
                backgroundColor: isDark
                    ? theme.palette.background.paper
                    : theme.palette.background.default,
                border: `1px solid ${isDark ? theme.palette.divider : theme.palette.grey[300]
                    }`,
                borderRadius: 1,
                p: "6px 10px",
                boxShadow: isDark
                    ? "0 2px 8px rgba(0,0,0,0.6)"
                    : "0 2px 6px rgba(0,0,0,0.1)",
                color: theme.palette.text.primary,
            }}
        >
            <Typography
                variant="body2"
                sx={{
                    fontWeight: 600,
                    color: theme.palette.primary.main,
                }}
            >
                {name}
            </Typography>

            <Typography variant="body2" color="text.secondary">
                Total:{" "}
                <strong style={{ color: theme.palette.primary.main }}>
                    {formatNumber(Number(payload?.[0]?.value ?? 0))}
                </strong>
            </Typography>
        </Box>
    );
};
