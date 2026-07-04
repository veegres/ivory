import {Box, Tooltip} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/type"

const SX: SxPropsMap = {
    label: {
        display: "inline-flex", alignItems: "center", borderRadius: 1, padding: "3px 8px", width: "fit-content",
        cursor: "pointer", textWrap: "nowrap", lineHeight: 1, gap: 0.75,
    },
    dot: {width: 5, height: 5, borderRadius: "50%"},
}

const COLORS: {[k: string]: {dot: string, bg: string, text: string}} = {
    "success": {dot: "success.main", bg: "rgba(46, 125, 50, 0.1)", text: "success.main"},
    "warning": {dot: "warning.main", bg: "rgba(211, 47, 47, 0.1)", text: "warning.main"},
    "error": {dot: "error.main", bg: "rgba(211, 47, 47, 0.1)", text: "error.main"},
    "info": {dot: "info.main", bg: "rgba(2, 136, 209, 0.1)", text: "info.main"},
    "primary": {dot: "primary.main", bg: "rgb(0 162 232 / 0.1)", text: "primary.main"},
    "secondary": {dot: "secondary.main", bg: "rgba(156, 39, 176, 0.1)", text: "secondary.main"},
    "default": {dot: "text.disabled", bg: "action.hover", text: "text.secondary"},
}

type Props = {
    label: string,
    title?: ReactNode,
    dot?: boolean,
    color?: "success" | "warning" | "error" | "primary" | "secondary" | "default" | "info" | string,
}

export function InfoColorBox(props: Props) {
    const {title, label, dot = false, color = "default"} = props
    const c = COLORS[color] ?? COLORS["default"]
    return (
        <Tooltip title={title} placement={"top"} disableInteractive={!!title}>
            <Box sx={[SX.label, {bgcolor: c.bg, color: c.text}]}>
                {dot && <Box sx={[SX.dot, {bgcolor: c.dot}]}/>}
                {label}
            </Box>
        </Tooltip>
    )
}

