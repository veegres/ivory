import {Box} from "@mui/material"

import {useRouterClusterListKey} from "../../../../features/cluster/api/ClusterHook"
import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Refresher} from "../../../widgets/browser/Refresher"

const SX: SxPropsMap = {
    box: {padding: "0px 5px"},
}

// TODO it updates only current overview, probably it is not correct and we should update all of them...
export function ListTableRefresher() {
    const clusterListKeys = useRouterClusterListKey()
    return (
        <Box sx={SX.box}>
            <Refresher queryKeys={[clusterListKeys, ClusterApi.overview.key()]}/>
        </Box>
    )
}