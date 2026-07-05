import {Box, Chip, Divider} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    row: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", gap: 2},
    label: {textTransform: "uppercase"},
}

type Props = {
    label?: ReactNode,
    children: ReactNode,
}

export function DividerBox(props: Props) {
    const {label, children} = props

    return (
        <>
            {label && <Divider sx={SX.label} textAlign={"left"}><Chip label={label}/></Divider>}
            <Box sx={SX.row}>
                {children}
            </Box>
        </>
    )
}
