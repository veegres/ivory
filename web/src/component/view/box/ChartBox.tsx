import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {
        flex: 1, display: "flex", flexDirection: "column", borderRadius: 1, padding: "8px 10px 0px 10px",
        border: "1px solid", borderColor: "divider", minWidth: 0,
    },
}

type Props = {
    children: ReactNode,
}

export function ChartBox(props: Props) {
    const {children} = props

    return (
        <Box sx={SX.box}>
            {children}
        </Box>
    )
}
