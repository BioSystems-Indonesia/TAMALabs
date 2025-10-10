import CategoryIcon from '@mui/icons-material/Category';
import PagesIcon from '@mui/icons-material/Pages';
import SegmentIcon from '@mui/icons-material/Segment';
import DeviceHubIcon from '@mui/icons-material/DeviceHub';
import Divider from "@mui/material/Divider";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { FilterList, FilterListItem, FilterLiveSearch, useGetList, useListContext } from "react-admin";
import { stopEnterPropagation } from '../../helper/component';
import { FieldValues, UseFormGetValues, UseFormSetValue } from 'react-hook-form'
import type { ObservationRequestCreateRequest } from '../../types/observation_requests';
import type { Device } from '../../types/device';
import SideFilter from '../../component/SideFilter';
import useAxios from '../../hooks/useAxios';
import { Typography, useTheme } from '@mui/material';

type TestFilterSidebarProps = {
    setSelectedData: React.Dispatch<React.SetStateAction<Record<number, ObservationRequestCreateRequest>>>
    selectedData: Record<number, ObservationRequestCreateRequest>
    setValue?: UseFormSetValue<FieldValues>
    getValues?: UseFormGetValues<FieldValues>
}

export const TestFilterSidebar = ({
    setSelectedData,
    selectedData,
    setValue,
    getValues,
}: TestFilterSidebarProps) => {
    const theme = useTheme();
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

    const { data } = useGetList(
        'test-template',
        {
            pagination: { page: 1, perPage: 1000 },
            sort: { field: 'id', order: 'DESC' }
        }
    );

    const isTemplateSelected = (value: any, filters: any) => {
        if (getValues) {
            const testTemplateIDs = getValues('test_template_ids')
            if (!testTemplateIDs?.includes(value.template.id)) {
                return false
            }

            return true
        }

        const templates = filters.templates || [];
        return templates.includes(value.template.id);
    };

    const toggleTemplateFilter = (value: any, filters: any) => {
        let templates = filters.templates || [];
        if (getValues) {
            const testTemplateIDs = getValues('test_template_ids')
            if (testTemplateIDs) {
                templates = testTemplateIDs
            }
        }

        const removeTemplate = (value: any, templates: any[]) => {
            setSelectedData(v => {
                const newSelectedData = { ...v }
                const testTypes = value.template.test_types as Record<number, ObservationRequestCreateRequest>
                Object.entries(testTypes).forEach(([key, value]) => {
                    delete newSelectedData[value.test_type_id]
                })
                return newSelectedData
            })
            if (setValue && getValues) {
                const testTemplateIDs = getValues('test_template_ids')
                if (testTemplateIDs?.includes(value.template.id)) {
                    setValue('test_template_ids', testTemplateIDs.filter((v: number) => v !== value.template.id))
                }
            }

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
            const doctorIDs = value.template.doctor_ids
            const analyzersIDs = value.template.analyzer_ids
            if (setValue) {
                setValue('doctor_ids', doctorIDs)
                setValue('analyzer_ids', analyzersIDs)
            }

            if (setValue && getValues) {
                const testTemplateIDs = getValues('test_template_ids')
                if (!testTemplateIDs) {
                    setValue('test_template_ids', [value.template.id])
                } else if (!testTemplateIDs.includes(value.template.id)) {
                    setValue('test_template_ids', [...testTemplateIDs, value.template.id])
                }
            }

            return [...templates, value.template.id]
        }

        return {
            ...filters,
            templates: templates.includes(value.template.id)
                ? removeTemplate(value, templates)
                : addTemplate(value, templates),
        };
    };

    const axios = useAxios()
    const { data: filter } = useQuery({
        queryKey: ['filterTestType'],
        queryFn: () => axios.get('/test-type/filter').then(res => res.data),
    });

    // Fetch devices for device filter
    const { data: devices } = useQuery({
        queryKey: ['devices'],
        queryFn: () => axios.get('/device').then(res => res.data),
    });

    // Device filter functions
    const isDeviceSelected = (value: any, filters: any) => {
        const deviceIds = filters.device_id || [];
        // Handle null device (General tests)
        if (value.id === null) {
            return deviceIds.includes(null);
        }
        return deviceIds.includes(value.id);
    };

    const toggleDeviceFilter = (value: any, filters: any) => {
        const deviceIds = filters.device_id || [];
        const targetId = value.id;

        return {
            ...filters,
            device_id: deviceIds.includes(targetId)
                ? deviceIds.filter((v: any) => v !== targetId)
                : [...deviceIds, targetId],
        };
    };

    return (
        <SideFilter
            sx={{
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
            <Typography variant="h6" sx={{
                color: theme.palette.text.primary,
                fontWeight: 600,
                fontSize: '1rem',
                textAlign: 'center'
            }}>
                🧪 Filter Type Test
            </Typography>
            <FilterLiveSearch
                placeholder="Search..."
                source="q"
                sx={{
                    marginTop: 2,
                    '& .MuiInputBase-root': {
                        backgroundColor: 'background.paper',
                    },
                }}
                onKeyDown={stopEnterPropagation}
            />
            <Divider />
            <FilterList
                label="Template"
                icon={<PagesIcon color="primary" />}
                sx={{
                    '& .MuiListItemIcon-root': {
                        minWidth: 36,
                    },
                }}
            >
                {data?.map((val: any, i: number) => (
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
                icon={<CategoryIcon color="primary" />}
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
                icon={<SegmentIcon color="primary" />}
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
            <Divider />
            <FilterList
                label="Device"
                icon={<DeviceHubIcon color="primary" />}
                sx={{
                    '& .MuiListItemIcon-root': {
                        minWidth: 36,
                    },
                }}
            >
                {/* Add "All Devices" option */}
                {devices?.map((device: Device, i: number) => (
                    <FilterListItem
                        key={i}
                        label={`${device.name}`}
                        value={device}
                        isSelected={isDeviceSelected}
                        toggleFilter={toggleDeviceFilter}
                    />
                ))}
            </FilterList>
        </SideFilter>
    )
};
// Removed custom useTheme function, using MUI's useTheme instead.

