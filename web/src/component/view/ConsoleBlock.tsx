import React, {ReactNode} from "react";
import {Box} from "@mui/material";
import {SxProps} from "@mui/system";

const SX = {
    box: { fontSize: "13px", width: "100%", background: "#000", padding: "10px 20px", borderRadius: "5px", color: "#e0e0e0" }
}

type Props = {
    children: ReactNode,
    sx?: SxProps,
}
export function ConsoleBlock(props: Props) {
    return (
        <Box sx={{ ...props.sx, ...SX.box }}>
            {props.children}
        </Box>
    )
}
