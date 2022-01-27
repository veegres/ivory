import {Box, CircularProgress, TableCell} from "@mui/material";
import React from "react";

export function TableCellFetching({ isFetching }: { isFetching: boolean }) {
    if (!isFetching) return <TableCell />
    return (
        <TableCell>
            <Box display={"flex"} justifyContent={"right"} alignItems={"center"}>
                <CircularProgress size={20} />
            </Box>
        </TableCell>
    )
}
