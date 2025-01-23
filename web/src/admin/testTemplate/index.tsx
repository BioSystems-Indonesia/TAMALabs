import CategoryIcon from '@mui/icons-material/Category';
import { Divider } from "@mui/material";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Chip from "@mui/material/Chip";
import Grid from "@mui/material/Grid";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useEffect, useRef, useState } from "react";
import { Create, Datagrid, DeleteButton, Edit, FilterList, FilterListItem, FilterLiveSearch, List, NumberField, ReferenceArrayField, SaveButton, SavedQueriesList, SimpleForm, TextField, TextInput, Toolbar, required, useListContext, useRecordContext } from "react-admin";
import SegmentIcon from '@mui/icons-material/Segment';
import { useFormContext } from "react-hook-form";
import type { ActionKeys } from "../../types/props";

export const TestTemplateList = () => (
    <List aside={<TestTemplateFilterSidebar />} title="Test Template">
        <Datagrid bulkActionButtons={false}>
            <NumberField source="id" />
            <TextField source="name" />
            <TextField source="description" />
            <ReferenceArrayField reference="test-type" source="test_type_id" />
        </Datagrid>
    </List>
);

const TestTemplateFilterSidebar = () => {
    return (
        <Card sx={{ order: -1, mr: 2, mt: 2, width: 300 }}>
            <CardContent>
                <FilterLiveSearch />
            </CardContent>
        </Card>
    )
};


type TestTemplateFormProps = {
    readonly?: boolean
    mode?: ActionKeys
}

function TestTemplateForm(props: TestTemplateFormProps) {
    return (
        <SimpleForm disabled={props.readonly} toolbar={false}>
            <TestTypeToolbar />
            <Divider sx={{
                marginBottom: "36px",
            }} />
            <TextInput source="name" readOnly={props.readonly} validate={[required()]} />
            <TextInput source="description" readOnly={props.readonly} multiline />
            <TestInput {...props} />
        </SimpleForm>
    )
}

const TestFilterSidebar = () => {
    const list = useListContext();
    const [dataUniqueCategory, setDataUniqueCategory] = useState<Array<any>>([])
    const [dataUniqueSubCategory, setDataUniqueSubCategory] = useState<Array<any>>([])
    const hasRunEffect = useRef(false); // Ref to track if the effect has run

    useEffect(() => {
        if (!list.data || hasRunEffect.current) {
            return; 
        }


        // Use Map to ensure uniqueness by 'id'
        const uniqueCategoryMap = new Map<string, any>();
        list.data.forEach((item: any) => {
            uniqueCategoryMap.set(item.category, item);
        });
        const uniqueCategoryArray = Array.from(uniqueCategoryMap.values());
        setDataUniqueCategory(uniqueCategoryArray)

        const uniqueSubCategoryMap = new Map<string, any>();
        list.data.forEach((item: any) => {
            uniqueSubCategoryMap.set(item.sub_category, item);
        });
        const uniqueSubCategoryArray = Array.from(uniqueSubCategoryMap.values());
        setDataUniqueSubCategory(uniqueSubCategoryArray)

        hasRunEffect.current = true;
    }, [list.data])

    const isCategorySelected = (value: any, filters: any) => {
        const categories = filters.categories || [];
        return categories.includes(value.category);
    };

    const toggleCategoryFilter = (value: any, filters: any) => {
        const categories = filters.categories || [];
        return {
            ...filters,
            categories: categories.includes(value.category)
                // Remove the category if it was already present
                ? categories.filter((v: any) => v !== value.category)
                // Add the category if it wasn't already present
                : [...categories, value.category],
        };
    };


    const isSubCategorySelected = (value: any, filters: any) => {
        const subCategories = filters.subCategories || [];
        return subCategories.includes(value.sub_category);
    };

    const toggleSubCategoryFilter = (value: any, filters: any) => {
        const subCategories = filters.subCategories || [];
        return {
            ...filters,
            subCategories: subCategories.includes(value.sub_category)
                // Remove the category if it was already present
                ? subCategories.filter((v: any) => v !== value.sub_category)
                // Add the category if it wasn't already present
                : [...subCategories, value.sub_category],
        };
    };


    return (
        <Card sx={{
            order: -1, mr: 1, mt: 2, width: 200, minWidth: 200,
            overflow: "visible",
        }}>
            <CardContent sx={{
                position: "sticky",
                top: 96,
            }}>
                <SavedQueriesList />
                <FilterLiveSearch onSubmit={(event) => event.preventDefault()} />
                <FilterList label="Category" icon={<CategoryIcon />}>
                    {dataUniqueCategory.map((val: any, i) => {
                        return (
                            <FilterListItem key={i} label={val.category} value={{ category: val.category }}
                                toggleFilter={toggleCategoryFilter} isSelected={isCategorySelected} />
                        )
                    })}
                </FilterList>
                <FilterList label="Sub Category" icon={<SegmentIcon />}>
                    {dataUniqueSubCategory.map((val: any, i) => {
                        return (
                            <FilterListItem key={i} label={val.sub_category} value={{ sub_category: val.sub_category }}
                                toggleFilter={toggleSubCategoryFilter} isSelected={isSubCategorySelected} />
                        )
                    })}
                </FilterList>
            </CardContent>
        </Card>
    )
};

const testTypeField = "test_type_id"

function TestTable(props: TestTemplateFormProps) {
    const { selectedIds, onSelect, data: testList } = useListContext();
    const { setValue } = useFormContext();

    const data = useRecordContext()
    useEffect(() => {
        if (data === undefined) {
            return;
        }
        console.log("data", data);

        onSelect(data.test_type_id);
        setValue(testTypeField, data.test_type_id);
    }, [data]);

    useEffect(() => {
        console.debug("selected ids", selectedIds);
        setValue(testTypeField, selectedIds);
    }, [selectedIds, testList]);

    const BulkActionButtons = () => {
        return (
            <></>
        );
    };

    return <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
            <Datagrid width={"100%"}
                bulkActionButtons={<BulkActionButtons />}
                rowClick={"toggleSelection"}
            >
                <TextField label={"Name"} source={"name"} />
                <TextField label={"Code"} source={"code"} />
                <TextField label={"Category"} source={"category"} />
                <TextField label={"Sub Category"} source={"sub_category"} />
                <TextField label={"Description"} source={"description"} />
            </Datagrid>
        </Grid>
        <Grid item xs={12} md={4}>
            <PickedTest />
        </Grid>
    </Grid>;
}

function TestInput(props: TestTemplateFormProps) {
    return (<List resource={"test-type"} exporter={false} aside={<TestFilterSidebar />}
        perPage={999999}
        storeKey={false}
        actions={false}
        title={false}
        pagination={false}
        disableSyncWithLocation
        sx={{
            marginTop: "48px",
            width: "100%"
        }}
    >
        <TestTable {...props} />
    </List>);
}

const PickedTest = () => {
    const { selectedIds, data } = useListContext();
    const [selectedData, setSelectedData] = useState<any[]>([]);

    useEffect(() => {
        if (!data) {
            return;
        }

        const selectedData = data.filter((v: any) => {
            return selectedIds.includes(v.id);
        });

        setSelectedData(selectedData);
    }, [selectedIds, data]);

    if (selectedIds.length === 0) {
        return (
            <Typography fontSize={16}>Please select test</Typography>
        )
    }

    return (
        <Stack spacing={2}>
            <Typography fontSize={16}>Selected test</Typography>
            <Grid container spacing={1}>
                {
                    selectedData.map((v: any) => {
                        return (
                            <Grid item key={v.id}>
                                <Chip label={v.code} />
                            </Grid>
                        )
                    })
                }
            </Grid>
        </Stack>
    )
}

const TestTypeToolbar = () => {
    return (
        <Stack width={"100%"}
            sx={{
                position: "sticky",
                top: 48,
                borderBottom: "1px solid #ccc",
                zIndex: 2147483647,
                marginBottom: 1,
            }}
        >
            <Toolbar sx={{
                gap: 2,
                width: "100%",
                display: "flex",
                justifyContent: "flex-end",
            }}>
                <DeleteButton variant="contained" size="small" />
                <SaveButton variant="contained" size="small" alwaysEnable />
            </Toolbar>
        </Stack>
    )
};


export function TestTemplateEdit() {
    return (
        <Edit mutationMode="pessimistic" title="Edit Test Template" sx={{
            "& .RaEdit-card": {
                overflow: "visible",
            }
        }}>
            <TestTemplateForm readonly={false} mode={"EDIT"} />
        </Edit>
    )
}

export function TestTemplateCreate() {
    return (
        <Create title="Create Test Template" redirect={"show"} sx={{
            "& .RaCreate-card": {
                overflow: "visible",
            }
        }}>
            <TestTemplateForm readonly={false} mode={"CREATE"} />
        </Create>
    )
}
