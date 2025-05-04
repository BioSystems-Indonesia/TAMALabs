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

    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DatePicker label={`${label} ${requiredAstrix(required)}`} format={"DD-MM-YYYY"}
                onChange={(date) => {
                    setValue(date);
                    field.onChange(date?.toISOString());
                }}
                slotProps={
                    {
                        field: {
                            clearable: clearable, onBlur: field.onBlur
                        }, 
                        textField: {
                            size: size
                        }
                    }
                }
                sx={{
                    //@ts-ignore this is perfectly fine
                    maxWidth: "280px",
                    ...sx
                }}
                ref={field.ref}
                name={field.name}
                value={value}
                disabled={field.disabled}
                size={'small'}

                readOnly={readonly}
                disableFuture={disableFuture}
                disablePast={disablePast}
            />
            <InputHelperText error={fieldState.error?.message} />
        </LocalizationProvider>
    );
}
