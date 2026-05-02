import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", borderRadius: 1, padding: "8px 10px 0px 10px",
        border: "1px solid", borderColor: "divider", width: "200px", flexGrow: 1,
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
}

export function ChartBox(props: Props) {
    const {children, label, value, unit} = props

    return (
        <Box sx={SX.box}>
            {renderHead()}
            {children}
        </Box>
    )

    function renderHead() {
        if (!label) return
        return (
            <Box sx={SX.head}>
                <Box sx={SX.label}>{label}</Box>
                {value !== undefined && <Box sx={SX.value}>{value.toFixed(2)} {unit ?? ""}</Box>}
            </Box>
        )
    }
}
