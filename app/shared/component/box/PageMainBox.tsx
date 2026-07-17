import {Box, Paper} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex"},
    paper: {width: "100%", margin: "0 5%", minWidth: "750px"},
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
            <Paper elevation={elevation} sx={{...SX.paper, padding}}>
                {props.children}
            </Paper>
        </Box>
    )
}
