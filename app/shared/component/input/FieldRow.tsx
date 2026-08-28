import {Box} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    // NOTE: auto columns rather than a fixed count, so a row of two and a row
    // of three line up on the same grid and a conditionally hidden field just
    // leaves one column fewer instead of breaking the alignment
    row: {display: "grid", gridAutoFlow: "column", gridAutoColumns: "1fr", gap: 1},
}

// FieldRow lays form fields out side by side in equal columns. It is the only
// way the deploy dialog puts two inputs on one line, so every such row lines up
// with every other one.
export function FieldRow(props: PropsWithChildren) {
    return <Box sx={SX.row}>{props.children}</Box>
}
