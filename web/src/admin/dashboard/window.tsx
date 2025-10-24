import { Card } from "@mui/material"
import { DashboardPage } from "."

export const DashboardWindow = () => {
    return (
        <Card sx={{ padding: 5, margin: 2, backgroundColor: "#f9f9f9ff" }}>
            <DashboardPage isWindow={true} />
        </Card>
    )
}