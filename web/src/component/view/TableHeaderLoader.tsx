import {Box, CircularProgress, TableCell} from "@mui/material";
import React, {cloneElement, ReactElement} from "react";


const SX = {
    button: { height: '30px', width: '30px' }
}

type Props = { isFetching: boolean, children?: ReactElement }

export function TableHeaderLoader({ isFetching, children }: Props) {
    return (
        <TableCell>
            <Box display={"flex"} justifyContent={"right"} alignItems={"center"} gap={1}>
                {isFetching ? <CircularProgress size={20}/> : null}
                {children ? cloneElement(children, { sx: SX.button }) : null}
            </Box>
        </TableCell>
    )
}
