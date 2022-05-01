import {Paper} from "@mui/material";
import {ReactNode} from "react";

const SX = {
    paper: {width: '100%', margin: '10px 5%'}
}

type Props = {
    children: ReactNode
    withPadding?: boolean
}

export function Block(props: Props) {
    const padding = props.withPadding ? "10px 20px" : "0";
    return (
        <Paper elevation={4} sx={{ ...SX.paper, padding}}>
            {props.children}
        </Paper>
    )
}
