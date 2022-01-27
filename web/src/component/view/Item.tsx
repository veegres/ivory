import {Paper} from "@mui/material";

const SX = {
    paper: { width: '100%', margin: '10px 5%' }
}

export function Item(props: any) {
    return (
        <Paper elevation={4} sx={SX.paper}>
            {props.children}
        </Paper>
    )
}
