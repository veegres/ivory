import {Stack} from "@mui/material"

import {useRouterClusterUpdate} from "../../../../api/cluster/hook"
import {ActiveCluster, ClusterOptions, ClusterOverview} from "../../../../api/cluster/type"
import {SxPropsMap} from "../../../../app/type"
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful"
import {Options} from "../../../widgets/options/Options"
import {OverviewOptionsInstance} from "./OverviewOptionsInstance"

const SX: SxPropsMap = {
    settings: {width: "250px", gap: 1, padding: "8px 0"},
}

type Props = {
    info: ActiveCluster,
    overview?: ClusterOverview,
}

export function OverviewOptions(props: Props) {
    const {info, overview} = props
    const {detectBy, cluster} = info

    const updateCluster = useRouterClusterUpdate()

    return (
        <Stack sx={SX.settings}>
            <OverviewOptionsInstance
                instances={overview?.instances ?? cluster.sidecarsOverview}
                mainInstance={overview?.mainInstance}
                detectBy={detectBy}
            />
            <LinearProgressStateful loading={updateCluster.isPending} line={true} color={"inherit"}/>
            <Options cluster={cluster} onUpdate={handleClusterUpdate}/>
        </Stack>
    )

    function handleClusterUpdate(opt: ClusterOptions) {
        updateCluster.mutate({...cluster, ...opt})
    }
}
