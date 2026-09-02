import {Box, SxProps, Theme} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {SxPropsFormatter} from "../../helper/HelperUtils"

const SX: SxPropsMap = {
    label: {
        typography: "button", lineHeight: 1, color: "action.active",
        display: "flex", alignItems: "center", justifyContent: "center",
        border: 1, borderColor: "divider", borderRadius: 1,
        paddingX: "10px", height: "32px", whiteSpace: "nowrap",
    },
}

type Props = {
    sx?: SxProps<Theme>,
}

// FieldLabel names a whole row of controls, the way a field's own label names
// one field. It is a control's height, so the row it heads lines up with it.
export function FieldLabel(props: PropsWithChildren<Props>) {
    return <Box sx={SxPropsFormatter.merge(SX.label, props.sx)}>{props.children}</Box>
}
