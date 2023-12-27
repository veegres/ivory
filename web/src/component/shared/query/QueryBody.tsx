import {Box, Divider} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    body: {padding: "8px 15px", fontSize: "13px"},
}

type Props = {
    children: ReactNode,
    show: boolean,
    render?: boolean,
}

export function QueryBody(props: Props) {
    const {show, children, render = false} = props
    if (!show && !render) return null

    return (
        <Box>
            <Divider/>
            <Box sx={SX.body} display={!show && render ? "none" : "default"}>
                {children}
            </Box>
        </Box>
    )
}
