import {Box} from "@mui/material"
import {PropsWithChildren} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    caption: {
        fontSize: "13px", fontWeight: 600, textTransform: "uppercase", lineHeight: 1,
        fontFamily: "monospace", color: "text.secondary",
    },
}

// Caption labels a section nested inside a box, where TitledBox or
// SubContentBox would add a second frame around content that already has one.
// It deliberately matches their uppercase monospace label a size down, so a
// nested section reads as part of the same family rather than a new style.
export function Caption(props: PropsWithChildren) {
    return <Box sx={SX.caption}>{props.children}</Box>
}
