import {Box, CircularProgress, SxProps, TableCell, Theme} from "@mui/material";
import {cloneElement, ReactElement} from "react";
import {SxPropsMap} from "../../type/common";

const SX: SxPropsMap = {
    box: {display: "flex", justifyContent: "right", alignItems: "center"},
    progress: {display: "flex", justifyContent: "center", alignItems: "center", padding: "0 5px"},
}

type Props = {
    isFetching: boolean,
    children?: ReactElement[],
    sx?: SxProps<Theme>,
    size?: number,
}

export function TableCellLoader(props: Props) {
    const {isFetching, children, sx, size = 32} = props
    return (
        <TableCell sx={sx}>
            <Box sx={SX.box}>
                {isFetching && <Box sx={SX.progress}><CircularProgress size={size - 15}/></Box>}
                {children && children.map((el, key) => cloneElement(el, {key, size}))}
            </Box>
        </TableCell>
    )
}
