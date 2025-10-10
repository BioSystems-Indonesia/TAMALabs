import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Barcode from "react-barcode";
// import { trimName } from "../helper/format";

export type LIMSBarcodeProps = {
    barcode: string
    name: string
    width?: number
    height?: number
    birthDt?: string
    sex?: string
}

// Helper function to format date
const formatDate = (dateString?: string): string => {
    if (!dateString) return '';

    try {
        const date = new Date(dateString);
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();
        return `${day}/${month}/${year}`;
    } catch {
        return dateString;
    }
};

// Helper function to format sex
const formatSex = (sex?: string): string => {
    if (!sex) return '';

    switch (sex.toUpperCase()) {
        case 'M':
        case 'MALE':
            return 'L';
        case 'F':
        case 'FEMALE':
            return 'P';
        default:
            return sex;
    }
};

export default function LIMSBarcode(props: LIMSBarcodeProps) {
    return (
        <Stack gap={0} justifyContent={"center"} alignItems={"center"}
            className={"barcode-container"} sx={{}}>
            <Typography
                fontSize={9}
                sx={{
                    margin: 0,
                    padding: 0,
                    textAlign: 'center',
                    fontWeight: 600,
                    lineHeight: 1.1,
                }}>
                {props.name}
            </Typography>
            <Typography
                fontSize={8}
                sx={{
                    margin: 0,
                    padding: 0,
                    textAlign: 'center',
                    fontWeight: 400,
                    lineHeight: 1.1,
                    marginBottom: 1,
                }}>
                {formatDate(props.birthDt)} | {formatSex(props.sex)}
            </Typography>
            <Barcode value={props.barcode} displayValue={false} height={props.height} margin={0} width={props.width} />
            <Typography
                className={"barcode-text"}
                fontSize={11}
                sx={{
                    margin: 0,
                    padding: 0,
                    textAlign: 'center',
                }}>{props.barcode}</Typography>
        </Stack>
    )


}