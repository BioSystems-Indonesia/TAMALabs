import { SxProps } from '@mui/material';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs, { Dayjs } from "dayjs";
import React from "react";
import { DateInputProps, InputHelperText, useInput } from "react-admin";
import { requiredAstrix } from "../helper/format.ts";

type CustomDateInputProps = {
    source: string
    label: string
    required?: boolean
    readonly?: boolean
    clearable?: boolean
    sx?: SxProps
    disableFuture?: boolean
    disablePast?: boolean
    size?: string
} & DateInputProps

export default function CustomDateInput({
    source, label, required, readonly, clearable, sx, disableFuture,
    disablePast, size
}: CustomDateInputProps) {
    const { field, fieldState } = useInput({ source });
    const [value, setValue] = React.useState<Dayjs | null>(field.value ? dayjs(field.value) : null);

    const labelText = label + ' ' + requiredAstrix(required);
    const slotProps ={
        field: {
            clearable: clearable,
            onBlur: field.onBlur
        },
        textField: { size: size }
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DatePicker
                label={labelText}
                format={"DD-MM-YYYY"}
                onChange={(date) => {
                    setValue(date);
                    field.onChange(date?.toISOString());
                }}
                slotProps={slotProps}
                // sx={{
                //     //@ts-ignore this is perfectly fine
                //     maxWidth: "280px",
                //     ...sx
                // }}
                // size={'small'}
                ref={field.ref}
                name={field.name}
                value={value}
                disabled={field.disabled}
                disableFuture={disableFuture}
                disablePast={disablePast}
            />
            <InputHelperText error={fieldState.error?.message} />
        </LocalizationProvider>
    );
}
