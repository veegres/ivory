import {Box, Theme} from "@mui/material"
import {SystemStyleObject} from "@mui/system"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    code: {
        display: "inline-block", border: 1, borderColor: "divider",
        color: "text.secondary", padding: "0px 5px", borderRadius: 1,
    },
}

type Props = {
    children: ReactNode,
    sx?: SystemStyleObject<Theme>,
}

export function Code(props: Props) {
    return <Box sx={[SX.code, props.sx ?? {}]}>{props.children}</Box>
}
