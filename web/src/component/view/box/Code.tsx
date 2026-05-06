import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    code: {
        display: "inline-block", border: 1, borderColor: "divider",
        color: "text.secondary", padding: "0px 5px", borderRadius: 1,
    },
}

type Props = {
    children: ReactNode,
}

export function Code(props: Props) {
    return <Box sx={SX.code}>{props.children}</Box>
}
