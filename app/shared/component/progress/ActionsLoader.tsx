import {Box, CircularProgress} from "@mui/material"
import {PropsWithChildren, ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center", flex: "1 1 auto", minWidth: 0},
    actions: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: "5px"},
    progress: {display: "flex", justifyContent: "center", alignItems: "center", padding: "0 5px"},
}

type Props = {
    loading: boolean,
    label?: ReactNode,
    size?: number,
}

export function ActionsLoader(props: PropsWithChildren<Props>) {
    const {loading, children, size = 32, label} = props
    return (
        <Box sx={SX.box}>
            <Box>{label}</Box>
            <Box sx={SX.actions}>
                {loading && <Box sx={SX.progress}><CircularProgress size={size - 15}/></Box>}
                {children}
            </Box>
        </Box>
    )
}
