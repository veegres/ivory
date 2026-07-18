import {Box, CircularProgress} from "@mui/material"
import {cloneElement, ReactElement, ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center", flex: "1 1 auto", minWidth: 0},
    actions: {display: "flex", justifyContent: "space-between", alignItems: "center", gap: "5px"},
    progress: {display: "flex", justifyContent: "center", alignItems: "center", padding: "0 5px"},
}

type Children = ReactElement<{size?: number}> | ReactElement<{size?: number}>[]
type Props = {
    loading: boolean,
    label?: ReactNode,
    children?: Children,
    size?: number,
}

export function ActionsLoader(props: Props) {
    const {loading, children, size = 32, label} = props
    return (
        <Box sx={SX.box}>
            <Box>{label}</Box>
            <Box sx={SX.actions}>
                {loading && <Box sx={SX.progress}><CircularProgress size={size - 15}/></Box>}
                {children && renderChildren(children)}
            </Box>
        </Box>
    )

    function renderChildren(children: Children) {
        if (Array.isArray(children)) return children.map((el, key) => cloneElement(el, {key, size}))
        else return cloneElement(children, {size})
    }
}
