import {Box, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {grey} from "@mui/material/colors";

const SX: SxPropsMap = {
    label: {borderRadius: 1, padding: "0 10px", textAlign: "center", cursor: "pointer", lineHeight: "1.55"},
}

type Props = {
    label: string,
    title?: string,
    bgColor?: string,
    color?: string,
    opacity?: number,
}

export function InfoColorBox(props: Props) {
    const {title, label, bgColor, color, opacity} = props
    return (
        <Tooltip title={title} placement={"top"} disableInteractive={!!title}>
            <Box sx={{...SX.label, opacity}} bgcolor={bgColor ?? "grey"} color={color ?? grey[50]}>
                {label}
            </Box>
        </Tooltip>
    )
}
