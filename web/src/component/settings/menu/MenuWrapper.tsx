import {Box} from "@mui/material";
import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1, height: "100%", overflow: "hidden"},
}

type Props = {
    children: ReactNode,
}

export function MenuWrapper(props: Props) {
    return (
        <Box sx={SX.box}>
            {props.children}
        </Box>
    )
}
