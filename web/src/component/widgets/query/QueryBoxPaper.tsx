import {ReactNode} from "react";
import {SxPropsMap} from "../../../api/management/type";
import {Paper} from "@mui/material";

const SX: SxPropsMap = {
    item: {fontSize: "15px"},
}

type Props = {
    children: ReactNode,
}

export function QueryBoxPaper(props: Props) {
    const {children} = props
    return (
        <Paper sx={SX.item} variant={"outlined"}>
            {children}
        </Paper>
    )
}
