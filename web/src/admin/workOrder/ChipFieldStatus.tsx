import React from "react";
import { useEffect } from "react"
import { ChipField, useRecordContext, type ChipFieldProps } from "react-admin"
import { getNestedValue } from "../../helper/accessor";

export type WorkOrderStatusChipFieldProps = Partial<ChipFieldProps> & {
    record?: any
    source: string
}

export function WorkOrderChipColorMap(value: string) {
    switch (value) {
        case 'NEW':
            return 'default';
        case 'PENDING':
            return 'secondary';
        case 'SUCCESS':
            return 'success';
        case 'CANCELLED':
            return 'error';
        default:
            return 'default';
    }
}

export const WorkOrderStatusChipField = (props: WorkOrderStatusChipFieldProps) => {
    const [color, setColor] = React.useState<'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' | undefined>(undefined);
    const record = props.record ?? useRecordContext();

    useEffect(() => {
        if (record === undefined) {
            return;
        }

        const value = getNestedValue(record, props.source);
        const color = WorkOrderChipColorMap(value);
        setColor(color);
    }, [record, props.source]);


    return (
        <ChipField {...props} sx={{}} textAlign="center" color={color} source={props.source} />
    )
}