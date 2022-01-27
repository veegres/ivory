import {Paper} from "@mui/material";

export function Item(props: any) {
    return (
        <Paper elevation={4} sx={{ width: '100%', margin: '10px 5%' }}>
            {props.children}
        </Paper>
    )
}
