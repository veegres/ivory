import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/type"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", borderRadius: 2, padding: "10px 10px 0px 10px",
        border: 1, borderColor: "divider", flexGrow: 1,
    },
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", mb: 0.5},
    label: {fontSize: "12px", fontFamily: "monospace"},
    value: {fontSize: "12px", color: "text.secondary"},
}

type Props = {
    label: string,
    value?: number,
    unit?: string,
    children: ReactNode,
    width?: string,
    fixed?: boolean,
}

export function ChartBox(props: Props) {
    const {children, label, value, unit, width = "200px", fixed = true} = props

    return (
        <Box sx={SX.box} width={width}>
            {renderHead()}
            {children}
        </Box>
    )

    function renderHead() {
        if (!label) return
        return (
            <Box sx={SX.head}>
                <Box sx={SX.label}>{label}</Box>
                {value !== undefined && <Box sx={SX.value}>{fixed ? value.toFixed(2) : value} {unit ?? ""}</Box>}
            </Box>
        )
    }
}
