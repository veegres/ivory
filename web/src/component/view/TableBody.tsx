import {Skeleton, TableBody as TableBodyMui, TableCell, TableRow} from "@mui/material";
import {ReactNode} from "react";

type Props = {
    isLoading: boolean,
    cellCount?: number,
    rowCount?: number,
    children?: ReactNode
}

export function TableBody(props: Props) {
    const {isLoading, cellCount, rowCount, children} = props
    const cellCountLocal = cellCount ?? 1, rowCountLocal = rowCount ?? 3

    return (
        <TableBodyMui>
            {isLoading ? <Loading/> : children}
        </TableBodyMui>
    )

    function Loading() {
        const cells = []
        for (let i = 0; i < cellCountLocal; i++) {
            cells.push(<TableCell key={i}><Skeleton width={'100%'}/></TableCell>)
        }
        const rows = []
        for (let i = 0; i < rowCountLocal; i++) {
            rows.push(<TableRow key={i}>{cells}</TableRow>)
        }
        return <>{rows}</>
    }
}
