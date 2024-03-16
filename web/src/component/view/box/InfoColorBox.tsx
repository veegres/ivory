import {Box, Tooltip} from "@mui/material";
import {SxPropsMap} from "../../../type/common";
import {grey} from "@mui/material/colors";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    label: {borderRadius: 1, padding: "0 10px", textAlign: "center", cursor: "pointer", lineHeight: "1.55", textWrap: "nowrap"},
}

type Props = {
    label: string,
    title?: ReactNode,
    bgColor?: string,
    color?: string,
    opacity?: number,
}

export function InfoColorBox(props: Props) {
    const {title, label, bgColor, color, opacity} = props
    return (
        <Tooltip title={title} placement={"top"} disableInteractive={!!title}>
            <Box sx={{...SX.label, opacity}} bgcolor={bgColor ?? grey[800]} color={color ?? grey[50]}>
                {label}
            </Box>
        </Tooltip>
    )
}
