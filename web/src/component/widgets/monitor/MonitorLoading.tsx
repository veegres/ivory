import {Skeleton} from "@mui/material"

import {SxPropsMap} from "../../../app/type"
import {MonitorRow} from "./MonitorRow"

const SX: SxPropsMap = {
    skeleton: {transform: "unset", flexGrow: 1},
}

type Props = {
    count: number
}

export function MonitorLoading(props: Props) {
    return (
        <MonitorRow>
            {[...Array(props.count).keys()].map((key) => (
                <Skeleton sx={SX.skeleton} key={key} width={200} height={150}/>
            ))}
        </MonitorRow>
    )
}
