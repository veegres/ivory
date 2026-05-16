import {Stack} from "@mui/material"

import {Feature} from "../../../features/feature"
import {SxPropsMap} from "../../../shared/helper/type"
import {Access} from "../../widgets/access/Access"
import {List as ClusterList} from "./list/List"
import {Node as ClusterNode} from "./node/Node"
import {Overview as ClusterOverview} from "./overview/Overview"

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4},
}

export function ClusterBody() {
    return (
        <Stack sx={SX.stack}>
            <ClusterList/>
            <Access feature={Feature.ViewNodeDbOverview}>
                <ClusterOverview/>
                <ClusterNode/>
            </Access>
        </Stack>
    )
}
