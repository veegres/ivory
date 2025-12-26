import {TableCell} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    cell: {verticalAlign: "top"},
}

type Props = {
    children: ReactNode,
    width?: number | string,
}

export function ListCell(props: Props) {
    const {children, width} = props
    return (
        <TableCell sx={SX.cell} width={width}>{children}</TableCell>
    )
}
