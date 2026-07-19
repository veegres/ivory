import {Box, SxProps, Theme} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SxPropsFormatter} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", alignItems: "flex-start", whiteSpace: "nowrap", lineHeight: 1.2, minWidth: 0, padding: "0 6px"},
    label: {fontSize: "8px", textTransform: "uppercase", color: "text.secondary"},
}

type Props = {
    label: string,
    sx?: SxProps<Theme>,
    children?: ReactNode,
}

export function InfoLabelBox(props: Props) {
    const {label, sx, children} = props
    return (
        <Box sx={SxPropsFormatter.merge(SX.box, sx)}>
            <Box sx={SX.label}>{label}</Box>
            {children ?? "-"}
        </Box>
    )
}
