import {ReactNode} from "react";
import {Box} from "@mui/material";
import {SxPropsMap} from "../../../app/types";

const SX: SxPropsMap = {
    box: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"}
}


type Props = {
    children: ReactNode,
}

export function QueryInfo(props: Props) {
    return (
        <Box sx={SX.box}>{props.children}</Box>
    )
}
