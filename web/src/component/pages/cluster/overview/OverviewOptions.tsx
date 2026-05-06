import {Stack} from "@mui/material"

import {useRouterClusterUpdate} from "../../../../api/cluster/hook"
import {Cluster, Options as ClusterOptions, Overview as ClusterOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Options} from "../../../widgets/options/Options"
import {OverviewOptionsNode} from "./OverviewOptionsNode"

const SX: SxPropsMap = {
    settings: {width: "250px", gap: 1, padding: "8px 0"},
}

type Props = {
    cluster: Cluster,
    overview?: ClusterOverview,
    mainKeeper?: string,
    manualKeeper?: string,
}

export function OverviewOptions(props: Props) {
    const {cluster, overview, mainKeeper, manualKeeper} = props

    const updateCluster = useRouterClusterUpdate()

    return (
        <Stack sx={SX.settings}>
            <OverviewOptionsNode
                nodes={overview?.nodes ?? cluster.nodesOverview ?? {}}
                mainKeeper={mainKeeper}
                manualKeeper={manualKeeper}
            />
            <LinearProgressStateful loading={updateCluster.isPending} line={true} color={"inherit"}/>
            <Options options={cluster} onUpdate={handleClusterUpdate}/>
        </Stack>
    )

    function handleClusterUpdate(opt: ClusterOptions) {
        updateCluster.mutate({...cluster, ...opt})
    }
}
