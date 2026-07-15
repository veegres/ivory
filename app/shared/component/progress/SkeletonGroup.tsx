import {Skeleton} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    skeleton: {flexGrow: 1},
}

type Props = {
    count: number
}

export function SkeletonGroup(props: Props) {
    return (
        <>
            {[...Array(props.count).keys()].map((key) => (
                <Skeleton variant={"rounded"} sx={SX.skeleton} key={key} width={200} height={150}/>
            ))}
        </>
    )
}
