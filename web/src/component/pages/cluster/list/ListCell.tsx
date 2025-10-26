import {TableCell} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../../../app/type"

const SX: SxPropsMap = {
    cell: {verticalAlign: "top"},
}

type Props = {
    children: ReactNode
}

export function ListCell(props: Props) {
    const {children} = props
    return (
        <TableCell sx={SX.cell}>{children}</TableCell>
    )
}
