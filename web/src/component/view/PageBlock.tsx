import {Grid, Paper} from "@mui/material";
import {ReactNode} from "react";

const SX = {
    paper: {width: '100%', margin: '10px 5%', minWidth: "750px"}
}

type Props = {
    children: ReactNode
    withPadding?: boolean
    visible?: boolean
    elevation?: number
}

export function PageBlock(props: Props) {
    const elevation = props.elevation ?? 4
    const visible = props.visible ?? true
    const padding = props.withPadding ? "10px 20px" : "0";

    if (!visible) return null

    return (
        <Grid container>
            <Paper elevation={elevation} sx={{ ...SX.paper, padding}}>
                {props.children}
            </Paper>
        </Grid>
    )
}
