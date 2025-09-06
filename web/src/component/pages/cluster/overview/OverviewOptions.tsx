import {Stack} from "@mui/material";
import {OverviewOptionsInstance} from "./OverviewOptionsInstance";
import {getDomain} from "../../../../app/utils";
import {Options} from "../../../widgets/options/Options";
import {ActiveCluster, ClusterOptions} from "../../../../api/cluster/type";
import {useRouterClusterUpdate} from "../../../../api/cluster/hook";
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful";
import {SxPropsMap} from "../../../../app/type";

const SX: SxPropsMap = {
    settings: {width: "250px", gap: "12px", padding: "8px 0"},
}

type Props = {
    info: ActiveCluster
}

export function OverviewOptions(props: Props) {
    const {info} = props
    const {combinedInstanceMap, defaultInstance, cluster, detection} = info

    const updateCluster = useRouterClusterUpdate()

    return (
        <Stack sx={SX.settings}>
            <OverviewOptionsInstance instance={getDomain(defaultInstance.sidecar)} instances={combinedInstanceMap} detection={detection}/>
            <LinearProgressStateful loading={updateCluster.isPending} line={true} color={"inherit"}/>
            <Options cluster={cluster} onUpdate={handleClusterUpdate}/>
        </Stack>
    )

    function handleClusterUpdate(opt: ClusterOptions) {
        updateCluster.mutate({...cluster, ...opt})
    }
}
