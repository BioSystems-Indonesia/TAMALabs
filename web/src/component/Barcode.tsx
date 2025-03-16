import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Barcode from "react-barcode";
import { trimName } from "../helper/format";

export type LIMSBarcodeProps = {
    barcode: string
    name: string
    width?: number
    height?: number
}

export default function LIMSBarcode(props: LIMSBarcodeProps) {
    return (
        <Stack gap={0} justifyContent={"center"} alignItems={"center"}
            className={"barcode-container"} sx={{}}>
            <Typography
                className={"barcode-text"}
                fontSize={11}
                sx={{
                    margin: 0,
                }}>{trimName(props.name, 14)} | {props.barcode}</Typography>
            <Barcode value={props.barcode} displayValue={false} height={props.height} margin={0} width={props.width} />
        </Stack>
    )


}