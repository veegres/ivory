import {Box, CircularProgress, SxProps, TableCell, Theme} from "@mui/material";
import {cloneElement, ReactElement} from "react";
import {SxPropsMap} from "../../app/types";


const SX: SxPropsMap = {
    button: {height: '30px', width: '30px'}
}

type Props = {
    isFetching: boolean,
    children?: ReactElement,
    sx?: SxProps<Theme>
}

export function TableCellLoader({isFetching, children, sx}: Props) {
    return (
        <TableCell sx={sx}>
            <Box display={"flex"} justifyContent={"right"} alignItems={"center"} gap={1}>
                {isFetching ? <CircularProgress size={20}/> : null}
                {children ? cloneElement(children, {sx: SX.button}) : null}
            </Box>
        </TableCell>
    )
}
