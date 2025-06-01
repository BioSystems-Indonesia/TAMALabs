import Chip from "@mui/material/Chip";
import { VerifiedStatus } from "../../types/work_order";

export function FilledPercentChip(props: { percent: number }) {
    // const [color, setColor] = useState<'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' | undefined>(undefined);
    // useEffect(() => {
    let color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' | undefined = 'default';
    if (props.percent < 0.5) {
        color = 'error';
    } else if (props.percent < 1) {
        color = 'warning';
    } else {
        color = 'success';
    }

    // setColor(color);
    // }, [props.percent]);


    return (
        <Chip label={`${(props.percent * 100).toFixed(2)}%`} color={color} />
    )
}

export function VerifiedChip(props: { verified: VerifiedStatus }) {
    let color: 'default' | 'primary' |'secondary' | 'error' | 'info' |'success' | 'warning' | undefined = 'default';
    if (props.verified === "PENDING") {
        color = 'default';
    } else if (props.verified === "VERIFIED") {
        color = 'success';
    } else if (props.verified === "REJECTED") {
        color = 'error';
    }

    return (
        <Chip label={props.verified} color={color} />
    )
}