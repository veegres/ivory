import {Skeleton, TableBody as MuiTableBody, TableCell, TableRow} from "@mui/material"
import {ReactNode} from "react"

type Props = {
    isLoading: boolean,
    cellCount?: number,
    rowCount?: number,
    children?: ReactNode,
    height?: number,
}

export function TableBody(props: Props) {
    const {isLoading, cellCount, rowCount, children, height} = props
    const cellCountLocal = cellCount ?? 1, rowCountLocal = rowCount ?? 3

    return (
        <MuiTableBody>
            {isLoading ? renderLoading() : children}
        </MuiTableBody>
    )

    function renderLoading() {
        return Array.from({length: rowCountLocal}).map((_, row) => (
            <TableRow key={row}>
                {Array.from({length: cellCountLocal}).map((_, cell) => (
                    <TableCell key={cell}><Skeleton variant={"rounded"} height={height} width={"100%"}/></TableCell>
                ))}
            </TableRow>
        ))
    }
}
