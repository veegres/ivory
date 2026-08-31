import {Box, CircularProgress} from "@mui/material"

import {useRouterClusterUpdate} from "../../../../features/cluster/api/ClusterHook"
import {Cluster, Options, Overview as ClusterOverview} from "../../../../features/cluster/api/ClusterType"
import {ClusterOptions} from "../../../../features/cluster/component/ClusterOptions"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {OverviewClusterConfigNode} from "./OverviewClusterConfigNode"

const SX: SxPropsMap = {
    settings: {display: "flex", flexDirection: "column", gap: 1, padding: "8px 0"},
    head: {display: "flex", flexWrap: "wrap", justifyContent: "space-between", alignItems: "center", gap: 1},
    saving: {
        display: "flex", alignItems: "center", gap: 1, minHeight: "26px", fontSize: "12px",
        color: "text.secondary", padding: "0px 5px", textTransform: "uppercase",
    },
}

type Props = {
    cluster: Cluster,
    overview?: ClusterOverview,
    mainKeeper?: string,
    manualKeeper?: string,
}

export function OverviewClusterConfig(props: Props) {
    const {cluster, overview, mainKeeper, manualKeeper} = props

    const updateCluster = useRouterClusterUpdate(cluster.name)

    return (
        <Box sx={SX.settings}>
            <Box sx={SX.head}>
                <OverviewClusterConfigNode
                    nodes={overview?.nodes ?? cluster.nodesOverview ?? {}}
                    mainKeeper={mainKeeper}
                    manualKeeper={manualKeeper}
                />
                {renderSaving()}
            </Box>
            <ClusterOptions options={cluster} onUpdate={handleClusterUpdate}/>
        </Box>
    )

    function renderSaving() {
        if (!updateCluster.isPending) return <Box sx={SX.saving}/>
        return (
            <Box sx={SX.saving}>
                <CircularProgress size={12} color={"inherit"}/>
                Saving Changes
            </Box>
        )
    }

    function handleClusterUpdate(opt: Options) {
        updateCluster.mutate({...cluster, ...opt})
    }
}
