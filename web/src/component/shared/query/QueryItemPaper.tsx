import {ReactNode} from "react";
import {SxPropsMap} from "../../../type/common";
import {Paper} from "@mui/material";

const SX: SxPropsMap = {
    item: {fontSize: "15px"},
}

type Props = {
    children: ReactNode,
}

export function QueryItemPaper(props: Props) {
    const {children} = props
    return (
        <Paper sx={SX.item} variant={"outlined"}>
            {children}
        </Paper>
    )
}
