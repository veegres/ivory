import {Box, Tooltip} from "@mui/material"
import {grey} from "@mui/material/colors"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    label: {display: "flex", alignItems: "center", borderRadius: 2, padding: "2px 10px 0", cursor: "pointer", minHeight: "15px", textWrap: "nowrap"},
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
