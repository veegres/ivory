import {Box} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    // NOTE: the size lives here rather than on InfoColorBox, which is used at
    // the ambient size elsewhere - this is the row that makes a tag a tag
    row: {display: "flex", alignItems: "center", gap: 0.5, flexShrink: 0, paddingX: 0.5, fontSize: "11px"},
}

// InfoColorBoxRow lays out one or more InfoColorBox tags as a row of metadata,
// at the one size they are read at wherever a template shows them.
export function InfoColorBoxRow(props: PropsWithChildren) {
    return <Box sx={SX.row}>{props.children}</Box>
}
