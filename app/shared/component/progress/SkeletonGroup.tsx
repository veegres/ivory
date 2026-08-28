import {Skeleton} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    grow: {flexGrow: 1},
}

type Props = {
    count: number,
    width?: number | string,
    height?: number | string,
    grow?: boolean,
}

export function SkeletonGroup(props: Props) {
    const {count, width = 200, height = 150, grow = true} = props

    return (
        <>
            {[...Array(count).keys()].map((key) => (
                <Skeleton variant={"rounded"} sx={[grow && SX.grow]} key={key} width={width} height={height}/>
            ))}
        </>
    )
}
