import {Skeleton} from "@mui/material"

import {SxPropsMap} from "../../helper/type"

const SX: SxPropsMap = {
    skeleton: {transform: "unset", flexGrow: 1},
}

type Props = {
    count: number
}

export function SkeletonGroup(props: Props) {
    return (
        <>
            {[...Array(props.count).keys()].map((key) => (
                <Skeleton sx={SX.skeleton} key={key} width={200} height={150}/>
            ))}
        </>
    )
}
