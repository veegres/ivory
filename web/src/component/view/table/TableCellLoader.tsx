import {Box, CircularProgress, SxProps, TableCell, Theme} from "@mui/material";
import {cloneElement, ReactElement} from "react";
import {SxPropsMap} from "../../../api/management/type";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "right", alignItems: "center"},
    progress: {display: "flex", justifyContent: "center", alignItems: "center", padding: "0 5px"},
}

type Children = ReactElement<{size?: number}> | ReactElement<{size?: number}>[]
type Props = {
    isFetching: boolean,
    children?: Children,
    sx?: SxProps<Theme>,
    size?: number,
}

export function TableCellLoader(props: Props) {
    const {isFetching, children, sx, size = 32} = props
    return (
        <TableCell sx={sx}>
            <Box sx={SX.box}>
                {isFetching && <Box sx={SX.progress}><CircularProgress size={size - 15}/></Box>}
                {children && renderChildren(children)}
            </Box>
        </TableCell>
    )

    function renderChildren(children: Children) {
        if (Array.isArray(children)) return children.map((el, key) => cloneElement(el, {key, size}))
        else return cloneElement(children, {size})
    }
}
