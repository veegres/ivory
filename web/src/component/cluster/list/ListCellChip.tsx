import {Box, Chip, Tooltip} from "@mui/material";
import {RefreshIconButton} from "../../view/IconButtons";
import {Cluster, InstanceDetection, SxPropsMap} from "../../../app/types";
import {grey, orange, purple} from "@mui/material/colors";
import {getDomain, InstanceColor} from "../../../app/utils";
import {InfoTitle} from "../../view/InfoTitle";
import {useStore} from "../../../provider/StoreProvider";
import {useEffect} from "react";

const SX: SxPropsMap = {
    chip: {width: "100%"},
    clusterName: {display: "flex", justifyContent: "center", alignItems: "center", gap: "3px"}
}

type Props = {
    cluster: Cluster,
    instanceDetection: InstanceDetection,
}

export function ListCellChip(props: Props) {
    const {cluster, instanceDetection} = props
    const {defaultInstance, combinedInstanceMap, warning, fetching, refetch} = instanceDetection

    const {isClusterActive, setCluster, store} = useStore()
    const detection = store.activeCluster?.detection ?? "auto"
    const isActive = isClusterActive(cluster.name)

    // NOTE: we don't want to use `setCluster` as deps cause it causes endless recursion/updates
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectStoreUpdate, [isActive, cluster, defaultInstance, combinedInstanceMap, warning, detection])

    return (
        <Box sx={SX.clusterName}>
            <Tooltip title={renderChipTooltip()} placement={"right-start"} arrow disableInteractive>
                <Chip
                    sx={SX.chip}
                    color={isActive ? "primary" : "default"}
                    label={cluster.name}
                    onClick={() => setCluster(isActive ? undefined : {cluster, defaultInstance, combinedInstanceMap, warning, detection})}
                />
            </Tooltip>
            <RefreshIconButton loading={fetching} onClick={refetch}/>
        </Box>
    )

    function renderChipTooltip() {
        const items = [
            {label: "Detection", value: detection, bgColor: purple[400]},
            {label: "Instance", value: getDomain(defaultInstance.sidecar), bgColor: InstanceColor[defaultInstance.role]},
            {label: "Warning", value: warning ? "Yes" : "No", bgColor: warning ? orange[500] : grey[500]}
        ]

        return <InfoTitle items={items}/>
    }

    function handleEffectStoreUpdate() {
        if (isActive) setCluster({cluster, defaultInstance, combinedInstanceMap, warning, detection})

        return () => {
            if (isActive) setCluster(undefined)
        }
    }

}
