import {Stack} from "@mui/material"

import {Feature} from "../../../features/feature"
import {ManageAccess} from "../../../features/management/component/ManageAccess"
import {SxPropsMap} from "../../../shared/helper/type"
import {List as ClusterList} from "./list/List"
import {Node as ClusterNode} from "./node/Node"
import {Overview as ClusterOverview} from "./overview/Overview"

const SX: SxPropsMap = {
    stack: {width: "100%", height: "100%", gap: 4},
}

export function ClusterBody() {
    return (
        <Stack sx={SX.stack}>
            <ManageAccess feature={Feature.ViewClusterList} error={true}>
                <ClusterList/>
            </ManageAccess>
            <ManageAccess feature={Feature.ViewClusterOverview} error={true}>
                <ClusterOverview/>
                <ClusterNode/>
            </ManageAccess>
        </Stack>
    )
}
