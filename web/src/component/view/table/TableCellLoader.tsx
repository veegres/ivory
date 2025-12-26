import {Box, CircularProgress, SxProps, TableCell, Theme} from "@mui/material"
import {cloneElement, ReactElement} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "space-between", alignItems: "center"},
    progress: {display: "flex", justifyContent: "center", alignItems: "center", padding: "0 5px"},
}

type Children = ReactElement<{size?: number}> | ReactElement<{size?: number}>[]
type Props = {
    loading: boolean,
    label?: string,
    children?: Children,
    width?: number | string,
    sx?: SxProps<Theme>,
    size?: number,
    colSpan?: number,
}

export function TableCellLoader(props: Props) {
    const {loading, children, sx, width, size = 32, colSpan, label} = props
    return (
        <TableCell sx={sx} width={width} colSpan={colSpan}>
            <Box sx={SX.box}>
                <Box>{label}</Box>
                <Box sx={SX.box}>
                    {loading && <Box sx={SX.progress}><CircularProgress size={size - 15}/></Box>}
                    {children && renderChildren(children)}
                </Box>
            </Box>
        </TableCell>
    )

    function renderChildren(children: Children) {
        if (Array.isArray(children)) return children.map((el, key) => cloneElement(el, {key, size}))
        else return cloneElement(children, {size})
    }
}
