import {Box, Paper} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex"},
    paper: {width: "100%", margin: "0 5%", minWidth: "750px"},
}

type Props = {
    children: ReactNode,
    withPadding?: boolean,
    withMarginTop?: string,
    withMarginBottom?: string,
    visible?: boolean,
    elevation?: number,
}

export function PageBox(props: Props) {
    const {withMarginTop, withMarginBottom} = props
    const elevation = props.elevation ?? 4
    const visible = props.visible ?? true
    const padding = props.withPadding ? "10px 20px" : "0";
    const margin = `${withMarginTop ?? 0} 5% ${withMarginBottom ?? 0}`

    if (!visible) return null

    return (
        <Box sx={SX.box}>
            <Paper elevation={elevation} sx={{ ...SX.paper, padding, margin}}>
                {props.children}
            </Paper>
        </Box>
    )
}
