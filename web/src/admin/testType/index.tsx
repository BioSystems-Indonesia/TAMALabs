import {Datagrid, FilterLiveSearch, List, Show, TextField} from "react-admin";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Box from "@mui/material/Box";

export const TestTypeList = () => (
    <List aside={<TestTypeFilterSidebar/>}>
        <Datagrid>
            <TextField source="name"/>
            <TextField source="code"/>
            <TextField source="category" />
            <TextField source="sub_category" />
            <TextField source="low_ref_range" label="low"/>
            <TextField source="high_ref_range" label="high"/>
            <TextField source="unit" />
            <TextField source="description" />
        </Datagrid>
    </List>
);

const TestTypeFilterSidebar = () => (
    <Card sx={{order: -1, mr: 2, mt: 2, width: 300}}>
        <CardContent>
            <FilterLiveSearch/>
        </CardContent>
    </Card>
);

function ReferenceSection() {
    return (
        <Box sx={{width: "100%"}}>
        </Box>
    )
}

export function TestTypeShow() {
    return (
        <Show>
            <ReferenceSection/>
        </Show>
    )
}