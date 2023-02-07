import {ReactNode} from "react";
import {Box, SxProps, Theme} from "@mui/material";
import {SxPropsMap} from "../../app/types";
import {mergeSxProps} from "../../app/utils";

const SX: SxPropsMap = {
    box: { fontSize: "13px", width: "100%", background: "#000", padding: "10px 20px", borderRadius: "5px", color: "#e0e0e0" }
}

type Props = {
    children: ReactNode,
    sx?: SxProps<Theme>,
}

export function ConsoleBlock(props: Props) {
    return (
        <Box sx={mergeSxProps(props.sx, SX.box)}>
            {props.children}
        </Box>
    )
}
