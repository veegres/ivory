import {Box, Divider} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../app/types";

const SX: SxPropsMap = {
    body: {padding: "8px 15px", fontSize: "13px"},
}

type Props = {
    children: ReactNode,
    show: boolean,
}

export function QueryItemBody(props: Props) {
    const {show, children} = props
    if (!show) return null

    return (
        <Box>
            <Divider/>
            <Box sx={SX.body}>
                {children}
            </Box>
        </Box>
    )
}
