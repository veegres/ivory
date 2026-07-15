import {Box} from "@mui/material"

import {useRouterClusterListKey} from "../../../../features/cluster/api/ClusterHook"
import {ClusterApi} from "../../../../features/cluster/api/ClusterRouter"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {Refresher} from "../../../widgets/browser/Refresher"

const SX: SxPropsMap = {
    box: {padding: "0px 5px"},
}

type Props = {
    clusters: string[],
}

export function ListTableRefresher(props: Props) {
    const {clusters} = props
    const clusterListKeys = useRouterClusterListKey()
    const clusterListOverviewKeys = clusters.map(c => ClusterApi.overview.key(c))
    return (
        <Box sx={SX.box}>
            <Refresher queryKeys={[clusterListKeys, ...clusterListOverviewKeys]}/>
        </Box>
    )
}