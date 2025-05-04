import Box from "@mui/material/Box";
import MUIBreadcrumbs from "@mui/material/Breadcrumbs";
import Typography from "@mui/material/Typography";
import { Link } from "react-admin";

export type BreadcrumbsLink = {
    label: string
    href: string
    icon?: React.ReactNode
    active?: boolean
}

export type BreadcrumbsProps = {
    links : Array<BreadcrumbsLink>
}

export default function Breadcrumbs(props: BreadcrumbsProps) {
    return (
        <Box sx={{ padding: 2, display: 'flex', justifyContent: 'center' }}>
            <MUIBreadcrumbs separator="›" aria-label="breadcrumb">
                {
                    props.links.length > 1 && props.links.map((link, index) => {
                        return (
                            <Link underline="hover" color={
                                index == props.links.length - 1 ? "text.primary" : "inherit"}
                                to={link.href} key={index} >
                                {link.icon}
                                <Typography variant="body2" color={link.active ? "text.primary" : "text.secondary"}
                                fontWeight={link.active ? "bold" : "normal"}
                                
                                >
                                    {link.label}
                                </Typography>
                            </Link>
                        )
                    })
                }
            </MUIBreadcrumbs>
        </Box>
    );

}