import {Skeleton, TableBody, TableCell, TableRow} from "@mui/material";
import {ReactNode} from "react";

export function TableBodySkeleton({ isLoading, cellCount, children }: { isLoading: boolean, cellCount?: number, children?: ReactNode }) {
    return (
        <TableBody>
            {isLoading ? <Loading /> : children}
        </TableBody>
    )

    function Loading() {
        return (
            <>
                <TableRow><TableCell colSpan={cellCount}><Skeleton width={'100%'} /></TableCell></TableRow>
                <TableRow><TableCell colSpan={cellCount}><Skeleton width={'100%'} /></TableCell></TableRow>
                <TableRow><TableCell colSpan={cellCount}><Skeleton width={'100%'} /></TableCell></TableRow>
            </>
        )
    }
}
