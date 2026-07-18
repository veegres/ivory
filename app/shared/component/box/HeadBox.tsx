import {Box} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    head: {
        display: "flex", justifyContent: "space-between", alignItems: "center", gap: 1,
        minHeight: "40px", padding: "4px 0px", borderBottom: 1, borderTop: 1, borderColor: "divider",
    },
    title: {
        display: "flex", alignItems: "center", alignSelf: "stretch",
        fontFamily: "monospace", fontSize: "13px", border: 1, borderRadius: 1,
        borderColor: "divider", padding: "0px 10px",
    },
}

type Props = {
    title?: string,
}

export function HeadBox(props: PropsWithChildren<Props>) {
    const {title, children} = props
    return (
        <Box sx={SX.head}>
            {title && <Box sx={SX.title}>{title}</Box>}
            {children}
        </Box>
    )
}
