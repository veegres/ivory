import {Grid, Paper} from "@mui/material";
import {ReactNode} from "react";

const SX = {
    paper: {width: '100%', margin: '10px 5%'}
}

type Props = {
    children: ReactNode
    withPadding?: boolean
    visible?: boolean
}

export function Block(props: Props) {
    const visible = props.visible ?? true
    const padding = props.withPadding ? "10px 20px" : "0";

    if (!visible) return null

    return (
        <Grid container>
            <Paper elevation={4} sx={{ ...SX.paper, padding}}>
                {props.children}
            </Paper>
        </Grid>
    )
}
