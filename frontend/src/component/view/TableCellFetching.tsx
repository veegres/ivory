import {Box, CircularProgress, TableCell} from "@mui/material";
import React, {cloneElement, ReactElement} from "react";

const SX = {
    icon: { fontSize: 20 },
    button: { height: '30px', width: '30px' }
}

export function TableCellFetching({ isFetching, children }: { isFetching: boolean, children?: ReactElement }) {
    return (
        <TableCell>
            <Box display={"flex"} justifyContent={"right"} alignItems={"center"} gap={1}>
                {isFetching ? <CircularProgress size={SX.icon.fontSize}/> : null}
                {children ? cloneElement(children, { sx: SX.button }) : null}
            </Box>
        </TableCell>
    )
}
