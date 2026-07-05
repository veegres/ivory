import {Box} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"
import {UnitBox} from "./UnitBox"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", borderRadius: 2, padding: "8px 10px 0px 10px",
        border: 1, borderColor: "divider", flexGrow: 1, gap: 1,
    },
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
        <Box sx={[SX.box, {width}]}>
            <UnitBox label={label} value={value} unit={unit} fixed={fixed}/>
            {children}
        </Box>
    )
}
