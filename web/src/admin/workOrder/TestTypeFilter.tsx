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

    const { data } = useGetList(
        'test-template',
        {
            pagination: { page: 1, perPage: 1000 },
            sort: { field: 'id', order: 'DESC' }
        }
    );

    const isTemplateSelected = (value: any, filters: any) => {
        const templates = filters.templates || [];
        return templates.includes(value.template.id);
    };

    const toggleTemplateFilter = (value: any, filters: any) => {

        const templates = filters.templates || [];
        const removeTemplate = (value: any, templates: any[]) => {
            console.log("removeTemplate", value, templates);
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
            console.log("addTemplate", value, templates);
            setSelectedData(v => {
                const newSelectedData = { ...v }
                const testTypes = value.template.test_types as Record<number, ObservationRequestCreateRequest>
                Object.entries(testTypes).forEach(([key, value]) => {
                    newSelectedData[value.test_type_id] = value
                })
                console.log(newSelectedData);

                return newSelectedData
            })

            return [...templates, value.template.id]
        }

        return {
            ...filters,
            templates: templates.includes(value.template.id)
                // Remove the category if it was already present
                ? removeTemplate(value, templates)
                // Add the category if it wasn't already present
                : addTemplate(value, templates),
        };
    };

    const { data: filter, isLoading: isFilterLoading } = useQuery({
        queryKey: ['filterTestType'],
        queryFn: () => fetch(import.meta.env.VITE_BACKEND_BASE_URL + '/test-type/filter').then(res => res.json()),
    });

    return (
        <Card sx={{
            order: -1,  width: 200, minWidth: 200,
            overflow: "visible",
        }}>
            <CardContent sx={{
                position: "sticky",
                top: 0,
            }}>
                <SavedQueriesList />
                <FilterLiveSearch onSubmit={(event) => event.preventDefault()} onKeyDown={(e) => { e.key === 'Enter' && e.preventDefault() }}  />
                <FilterList label="Template" icon={<PagesIcon />} >
                    {data?.map((val: any, i) => {
                        return (
                            <FilterListItem key={i} label={val.name} value={{ template: val }}
                                toggleFilter={toggleTemplateFilter} isSelected={isTemplateSelected} />
                        )
                    })}
                </FilterList>
                <Divider sx={{ marginY: 1 }} />
                {isFilterLoading ? <LinearProgress /> : (
                    <>
                        <FilterList label="Category" icon={<CategoryIcon />}>
                            {filter?.categories.map((val: string, i: number) => {
                                return (
                                    <FilterListItem key={i} label={val} value={{ category: val}}
                                        toggleFilter={toggleCategoryFilter} isSelected={isCategorySelected} />
                                )
                            })}
                        </FilterList>
                        <FilterList label="Sub Category" icon={<SegmentIcon />}>
                            {filter?.sub_categories.map((val: string, i: number) => {
                                return (
                                    <FilterListItem key={i} label={val} value={{ sub_category: val}}
                                        toggleFilter={toggleSubCategoryFilter} isSelected={isSubCategorySelected} />
                                )
                            })}
                        </FilterList>
                    </>

                )}
            </CardContent>
        </Card>
    )
};