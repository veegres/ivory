import {Box, Paper} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SxPropsFormatter} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    box: {display: "flex"},
    paper: {width: "100%", minWidth: 0, overflowX: "auto"},
}

type Props = {
    children: ReactNode,
    withPadding?: boolean,
    visible?: boolean,
    elevation?: number,
}

export function PageMainBox(props: Props) {
    const elevation = props.elevation ?? 4
    const visible = props.visible ?? true
    const padding = props.withPadding ? "5px 10px" : "0"

    if (!visible) return null

    return (
        <Box sx={SX.box}>
            <Paper elevation={elevation} sx={[SX.paper, SxPropsFormatter.style.pageMargin, {padding}]}>
                {props.children}
            </Paper>
        </Box>
    )
}
