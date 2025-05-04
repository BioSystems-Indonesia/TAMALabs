import { SxProps } from "@mui/material"
import Card from "@mui/material/Card"
import CardContent from "@mui/material/CardContent"
import { ReactNode } from "react"


type SideFilterProps = {
    children: ReactNode
    sx?: SxProps
}

export default function SideFilter(props: SideFilterProps) {
    return (
        <Card sx={{ order: -1, mr: 2, mt: 4.6, maxWidth: 300 , ...props.sx}}>
            <CardContent>
                {props.children}
            </CardContent>
        </Card>
    )
}