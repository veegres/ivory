import {ReactNode} from "react";
import {TableCell} from "@mui/material";

const SX = {
    cell: {verticalAlign: "top"},
}

type Props = {
    children: ReactNode
}

export function ListCell(props: Props) {
    const { children } = props
    return (
        <TableCell sx={SX.cell}>{children}</TableCell>
    )
}
