import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { DateInputProps, InputHelperText, useInput } from "react-admin";
import dayjs, { Dayjs } from "dayjs";
import React from "react";
import { requiredAstrix } from "../helper/format.ts";
import type { SxProps } from '@mui/material';

type CustomDateInputProps = {
    source: string
    label: string
    required?: boolean
    readonly?: boolean
    clearable?: boolean
    sx?: SxProps
} & DateInputProps

export default function CustomDateInput({ source, label, required, readonly, clearable, sx }: CustomDateInputProps) {
    const { field, fieldState } = useInput({ source });
    const [value, setValue] = React.useState<Dayjs | null>(field.value ? dayjs(field.value) : null);

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DatePicker label={`${label} ${requiredAstrix(required)}`} format={"DD-MM-YYYY"}
                onChange={(date) => {
                    setValue(date);
                    field.onChange(date?.toDate())
                }}
                slotProps={{ field: { clearable: clearable, onBlur: field.onBlur } }}
                sx={{
                    //@ts-ignore this is perfectly fine
                    maxWidth: "280px",
                    ...sx
                }}
                ref={field.ref}
                name={field.name}
                value={value}
                disabled={field.disabled}
                readOnly={readonly}
            />
            <InputHelperText error={fieldState.error?.message} />
        </LocalizationProvider>
    );
}
