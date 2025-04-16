import CategoryIcon from '@mui/icons-material/Category';
import PagesIcon from '@mui/icons-material/Pages';
import SegmentIcon from '@mui/icons-material/Segment';
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Divider from "@mui/material/Divider";
import LinearProgress from "@mui/material/LinearProgress";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { FilterList, FilterListItem, FilterLiveSearch, SavedQueriesList, useGetList, useListContext } from "react-admin";
import type { ObservationRequestCreateRequest } from '../../types/observation_requests';
import { Typography } from "@mui/material";

type TestFilterSidebarProps = {
    setSelectedData: React.Dispatch<React.SetStateAction<Record<number, ObservationRequestCreateRequest>>>
    selectedData: Record<number, ObservationRequestCreateRequest>
}

export const TestFilterSidebar = ({
    setSelectedData,
    selectedData,
}: TestFilterSidebarProps) => {
    const list = useListContext();
    const [_dataUniqueCategory, setDataUniqueCategory] = useState<Array<any>>([])
    const [_dataUniqueSubCategory, setDataUniqueSubCategory] = useState<Array<any>>([])
    const hasRunEffect = useRef(false);

    useEffect(() => {
        if (!list.data || hasRunEffect.current) {
            return;
        }

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
                ? categories.filter((v: any) => v !== value.category)
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
                ? subCategories.filter((v: any) => v !== value.sub_category)
                : [...subCategories, value.sub_category],
        };
    };

    const isTemplateSelected = (value: any, filters: any) => {
        const templates = filters.templates || [];
        return templates.includes(value.template.id);
    };

    const toggleTemplateFilter = (value: any, filters: any) => {
        const templates = filters.templates || [];
        const removeTemplate = (value: any, templates: any[]) => {
            setSelectedData(v => {
                const newSelectedData = { ...v }
                const testTypes = value.template.test_types as Record<number, ObservationRequestCreateRequest>
                Object.entries(testTypes).forEach(([key, value]) => {
                    delete newSelectedData[value.test_type_id]
                })
                return newSelectedData
            })
            return templates.filter((v: any) => v !== value.template.id)
        }

        const addTemplate = (value: any, templates: any[]) => {
            setSelectedData(v => {
                const newSelectedData = { ...v }
                const testTypes = value.template.test_types as Record<number, ObservationRequestCreateRequest>
                Object.entries(testTypes).forEach(([key, value]) => {
                    newSelectedData[value.test_type_id] = value
                })
                return newSelectedData
            })
            return [...templates, value.template.id]
        }

        return {
            ...filters,
            templates: templates.includes(value.template.id)
                ? removeTemplate(value, templates)
                : addTemplate(value, templates),
        };
    };

    const { data: filter, isLoading: isFilterLoading } = useQuery({
        queryKey: ['filterTestType'],
        queryFn: () => fetch(import.meta.env.VITE_BACKEND_BASE_URL + '/test-type/filter').then(res => res.json()),
    });

    return (
        <Card 
            sx={{ 
                order: -1,
                width: 280,
                position: 'sticky',
                top: 0,
                marginRight: 2,
                '& .MuiCardContent-root': {
                    padding: 2,
                },
                '& .RaFilterList-root': {
                    marginTop: 2,
                },
                '& .MuiDivider-root': {
                    margin: '16px 0',
                },
                '& .RaFilterListItem-root': {
                    marginBottom: 1,
                    transition: 'all 0.2s',
                    borderRadius: 1,
                    '&:hover': {
                        backgroundColor: 'action.hover',
                    },
                },
            }}
        >
            <CardContent>
                <FilterLiveSearch 
                    placeholder="Search..."
                    source="q"
                    sx={{
                        marginTop: 2,
                        '& .MuiInputBase-root': {
                            backgroundColor: 'background.paper',
                        },
                    }}
                />
                <Divider />
                <FilterList 
                    label="Template" 
                    icon={<PagesIcon />}
                    sx={{
                        '& .MuiListItemIcon-root': {
                            minWidth: 36,
                        },
                    }}
                >
                    {filter?.templates?.map((val: any, i: number) => (
                        <FilterListItem 
                            key={i} 
                            label={val.name} 
                            value={{ template: val }}
                            isSelected={isTemplateSelected}
                            toggleFilter={toggleTemplateFilter}
                        />
                    ))}
                </FilterList>
                <Divider />
                <FilterList 
                    label="Category" 
                    icon={<CategoryIcon />}
                    sx={{
                        '& .MuiListItemIcon-root': {
                            minWidth: 36,
                        },
                    }}
                >
                    {filter?.categories.map((val: string, i: number) => (
                        <FilterListItem 
                            key={i} 
                            label={val} 
                            value={{ category: val }}
                            isSelected={isCategorySelected}
                            toggleFilter={toggleCategoryFilter}
                        />
                    ))}
                </FilterList>
                <Divider />
                <FilterList 
                    label="Sub Category" 
                    icon={<SegmentIcon />}
                    sx={{
                        '& .MuiListItemIcon-root': {
                            minWidth: 36,
                        },
                    }}
                >
                    {filter?.sub_categories.map((val: string, i: number) => (
                        <FilterListItem 
                            key={i} 
                            label={val} 
                            value={{ sub_category: val }}
                            isSelected={isSubCategorySelected}
                            toggleFilter={toggleSubCategoryFilter}
                        />
                    ))}
                </FilterList>
            </CardContent>
        </Card>
    )
};