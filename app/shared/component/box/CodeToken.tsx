import {Tooltip} from "@mui/material"
import {ReactNode} from "react"

import {Code} from "./Code"

// NOTE: not typed as SxPropsMap - Code takes a plain SystemStyleObject, and
// the annotation is what makes the two disagree
const SX = {
    token: {fontFamily: "monospace", fontSize: "13px", padding: "2px 6px"},
}

type Props = {
    children: ReactNode,
    tooltip?: ReactNode,
}

// CodeToken is a single identifier shown inline as code - a deployment
// variable, a mask standing in for one. It is a Code at one fixed size so a
// token reads the same in a reference list as it does inside a sentence.
export function CodeToken(props: Props) {
    const {children, tooltip} = props
    if (!tooltip) return <Code sx={SX.token}>{children}</Code>
    return (
        <Tooltip title={tooltip} placement={"top-start"}>
            <Code sx={SX.token}>{children}</Code>
        </Tooltip>
    )
}
