import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    head: {display: "flex", justifyContent: "space-between", alignItems: "center", mb: 0.5},
    label: {fontSize: "13px", fontFamily: "monospace", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis"},
    value: {fontSize: "13px", color: "text.secondary", whiteSpace: "nowrap",},
}

type Props = {
    label?: string,
    value?: number,
    unit?: string,
    fixed?: boolean,
}

export function UnitBox(props: Props) {
    const {label, value, unit, fixed = true} = props

    if (!label) return null

    return (
        <Box sx={SX.head}>
            <Box sx={SX.label}>{label}</Box>
            {value !== undefined && <Box sx={SX.value}>{fixed ? value.toFixed(2) : value} {unit ?? ""}</Box>}
        </Box>
    )
}
