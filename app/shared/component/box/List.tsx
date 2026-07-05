import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    wrapper: {
        borderRadius: 2, border: 1, borderColor: "divider",
        "li:not(:last-child)": {borderBottom: 1, borderColor: "divider"}
    },
    name: {padding: "0px 12px", fontFamily: "monospace", fontWeight: "500"},
    container: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    name?: string,
    children: ReactNode,
}

export function List(props: Props) {
    const {name, children} = props

    return (
        <Box sx={SX.container}>
            {name && <Box sx={SX.name}>{name}</Box>}
            <Box sx={SX.wrapper}>{children}</Box>
        </Box>
    )
}
