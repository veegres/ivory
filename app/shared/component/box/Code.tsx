import {Box, Theme} from "@mui/material"
import {SystemStyleObject} from "@mui/system"
import {forwardRef, ReactNode} from "react"

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

export const Code = forwardRef<HTMLDivElement, Props>(function Code({children, sx, ...rest}, ref) {
    return <Box ref={ref} sx={[SX.code, sx ?? {}]} {...rest}>{children}</Box>
})
