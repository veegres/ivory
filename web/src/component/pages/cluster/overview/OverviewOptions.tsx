import {Divider, Stack} from "@mui/material";
import {OverviewOptionsInstance} from "./OverviewOptionsInstance";
import {getDomain} from "../../../../app/utils";
import {SxPropsMap} from "../../../../type/general";
import {Options} from "../../../shared/options/Options";
import {useMutationOptions} from "../../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {ClusterApi} from "../../../../app/api";
import {ActiveCluster, ClusterOptions} from "../../../../type/cluster";

const SX: SxPropsMap = {
    settings: {width: "250px", gap: "12px", padding: "8px 0"}
}

type Props = {
    info: ActiveCluster
}

export function OverviewOptions(props: Props) {
    const {info} = props
    const {combinedInstanceMap, defaultInstance, cluster, detection} = info

    const updateMutationOptions = useMutationOptions([["cluster", "list"], ["tag", "list"]])
    const updateCluster = useMutation({mutationFn: ClusterApi.update, ...updateMutationOptions})

    return (
        <Stack sx={SX.settings}>
            <OverviewOptionsInstance instance={getDomain(defaultInstance.sidecar)} instances={combinedInstanceMap} detection={detection}/>
            <Divider variant={"middle"}/>
            <Options cluster={cluster} onUpdate={handleClusterUpdate} />
        </Stack>
    )

    function handleClusterUpdate(opt: ClusterOptions) {
        updateCluster.mutate({...cluster, ...opt})
    }
}
