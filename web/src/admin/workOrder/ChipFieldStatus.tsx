import React from "react";
import { useEffect } from "react"
import { ChipField, useRecordContext, type ChipFieldProps } from "react-admin"

export const WorkOrderStatusChipField = (props: Partial<ChipFieldProps>) => {
    const data = useRecordContext()
    const [color, setColor] = React.useState<'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' | undefined>(undefined);

    useEffect(() => {
        switch (data?.status) {
            case 'NEW':
                setColor('default');
                break;
            case 'PENDING':
                setColor('secondary');
                break;
            case 'COMPLETED':
                setColor('success');
                break;
            case 'CANCELLED':
                setColor('error');
                break
            default:
                setColor('default');
                break;
        }
    }, [data]);


    return (
        <ChipField {...props} sx={{

        }} textAlign="center" color={color} source="status" />
    )
}