import {Box, Chip, TableRow, Tooltip} from "@mui/material";
import {useEffect, useState} from "react";
import {Cluster, DetectionType} from "../../app/types";
import {RefreshIconButton,} from "../view/IconButtons";
import {DynamicInputs} from "../view/DynamicInputs";
import {useStore} from "../../provider/StoreProvider";
import {useAutoInstanceDetection, useManualInstanceDetection} from "../../app/hooks";
import {ClustersRowRead} from "./ClustersRowRead";
import {ClustersRowUpdate} from "./ClustersRowUpdate";
import {ClustersCell} from "./ClustersCell";
import {initialInstance, InstanceColor} from "../../app/utils";
import {InfoTitle} from "../view/InfoTitle";
import {orange, purple} from "@mui/material/colors";

const SX = {
    chip: {width: "100%"},
    clusterName: {display: "flex", justifyContent: "center", alignItems: "center", gap: "3px"}
}

type Props = {
    name: string,
    cluster: Cluster,
    editable: boolean,
    toggle: () => void,
}

export function ClustersRow({name, cluster, editable, toggle}: Props) {
    const {setCluster, isClusterActive, store} = useStore()
    const [detection, setDetection] = useState<DetectionType>("auto")
    const [manualInstance, setManualInstance] = useState(initialInstance(""))
    const isActive = isClusterActive(name)
    const isDetectionManual = detection === "manual"

    const [stateNodes, setStateNodes] = useState(cluster.nodes);
    const auto = useAutoInstanceDetection(!isDetectionManual, cluster)
    const manual = useManualInstanceDetection(isDetectionManual, cluster, manualInstance)

    const {active: {warning, instance, instances}, colors, fetching, refetch} = isDetectionManual ? manual : auto

    // we need disable cause it thinks that function can be updated, and it causes endless recursion
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectStoreUpdate, [isActive, warning, cluster, instance, instances, detection])
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectRefetch, [isActive, detection])

    useEffect(handleEffectDetection, [isActive, store.activeCluster, detection])
    useEffect(handleEffectManualInstance, [isActive, store.activeCluster, manualInstance])

    return (
        <TableRow>
            <ClustersCell>
                <Box sx={SX.clusterName}>
                    <Tooltip title={renderChipTooltip()} placement={"right-start"} arrow disableInteractive>
                        <Chip
                            sx={SX.chip}
                            color={isActive ? "primary" : "default"}
                            label={name}
                            onClick={() => setCluster(isActive ? undefined : {warning, cluster, instance, instances, detection})}
                        />
                    </Tooltip>
                    <RefreshIconButton loading={fetching} onClick={refetch}/>
                </Box>
            </ClustersCell>
            <ClustersCell>
                <DynamicInputs
                    inputs={stateNodes}
                    colors={colors}
                    editable={editable}
                    placeholder={`Instance`}
                    onChange={n => setStateNodes(n)}
                />
            </ClustersCell>
            <ClustersCell>
                {!editable ? (
                    <ClustersRowRead name={name} toggle={toggle}/>
                ) : (
                    <ClustersRowUpdate
                        name={name}
                        nodes={stateNodes}
                        toggle={toggle}
                        onUpdate={refetch}
                        onClose={() => setStateNodes(cluster.nodes)}
                    />
                )}
            </ClustersCell>
        </TableRow>
    )

    function renderChipTooltip() {
        const items = [
            { name: "Detection", value: detection, bgColor: purple[400] },
            { name: "Instance", value: instance.api_domain, bgColor: InstanceColor[instance.role] },
            { name: "Warning", value: warning ? "Yes" : "No", bgColor: orange[500] }
        ]

        return <InfoTitle items={items} />
    }

    function handleEffectRefetch() {
        if (isActive) refetch()
    }

    function handleEffectStoreUpdate() {
        if (isActive){
            setCluster({warning, cluster, instance, instances, detection})
        }
    }

    function handleEffectDetection() {
        if (isActive && store.activeCluster) {
            const storeDetection = store.activeCluster.detection
            if (detection !== storeDetection) {
                setDetection(storeDetection)
            }
        }
    }

    function handleEffectManualInstance() {
        if (isActive && store.activeCluster) {
            const storeManualInstance = store.activeCluster.instance
            if (manualInstance !== storeManualInstance) {
                setManualInstance(storeManualInstance)
            }
        }
    }
}
