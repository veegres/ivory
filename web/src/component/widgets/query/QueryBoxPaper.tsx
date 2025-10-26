import {Paper} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

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
