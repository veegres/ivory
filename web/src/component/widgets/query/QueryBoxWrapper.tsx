import {Box, SxProps, Theme} from "@mui/material"
import {ReactNode, useState} from "react"

import {SxPropsMap} from "../../../app/type"
import {SxPropsFormatter} from "../../../app/utils"

const SX: SxPropsMap = {
    box: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px", minHeight: "40px"},
    edit: {outline: 1, outlineColor: "divider"},
    hover: {":hover": {outlineColor: "text.secondary"}},
    focus: {outlineColor: "primary.main"},
}


type Props = {
    editable?: boolean,
    children: ReactNode,
    sx?: SxProps<Theme>,
}

export function QueryBoxWrapper(props: Props) {
    const {children, editable, sx} = props
    const [focus, setFocus] = useState(false)

    const editSx = !editable ? [] : !focus ? [SX.edit, SX.hover] : [SX.edit, SX.focus]
    const commonSx = SxPropsFormatter.merge(SX.box, editSx)
    return (
        <Box sx={SxPropsFormatter.merge(commonSx, sx)} onFocus={() => setFocus(true)} onBlur={() => setFocus(false)}>
            {children}
        </Box>
    )
}
