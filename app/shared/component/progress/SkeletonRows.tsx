import {Box, Skeleton} from "@mui/material"
import {ReactNode} from "react"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    row: {padding: "5px"},
}

type Props = {
    isLoading: boolean,
    rowCount?: number,
    children?: ReactNode,
    height?: number,
}

export function SkeletonRows(props: Props) {
    const {isLoading, rowCount = 3, children, height} = props

    return isLoading ? renderLoading() : <>{children}</>

    function renderLoading() {
        return (
            <>
                {Array.from({length: rowCount}).map((_, row) => (
                    <Box key={row} sx={SX.row}>
                        <Skeleton variant={"rounded"} height={height} width={"100%"}/>
                    </Box>
                ))}
            </>
        )
    }
}
