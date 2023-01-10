import {Box, Chip, TableRow, Tooltip} from "@mui/material";
import {useEffect, useRef, useState} from "react";
import {Cluster, DetectionType, InstanceLocal} from "../../../app/types";
import {RefreshIconButton,} from "../../view/IconButtons";
import {DynamicInputs} from "../../view/DynamicInputs";
import {useStore} from "../../../provider/StoreProvider";
import {ListRowRead} from "./ListRowRead";
import {ListRowUpdate} from "./ListRowUpdate";
import {ListCell} from "./ListCell";
import {getDomain, initialInstance, InstanceColor} from "../../../app/utils";
import {InfoTitle} from "../../view/InfoTitle";
import {grey, orange, purple} from "@mui/material/colors";
import {useAutoInstanceDetection, useManualInstanceDetection} from "../../../hook/InstanceDetection";

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

export function ListRow({name, cluster, editable, toggle}: Props) {
    const {setCluster, isClusterActive, store} = useStore()
    const [detection, setDetection] = useState<DetectionType>("auto")
    const detectionRef = useRef(detection)
    const [manualState, setManualState] = useState({ instances: {}, instance: initialInstance() })
    const isActive = isClusterActive(name)
    const isDetectionManual = detection === "manual"

    const [stateNodes, setStateNodes] = useState(cluster.nodes);
    const auto = useAutoInstanceDetection(!isDetectionManual, cluster)
    const manual = useManualInstanceDetection(isDetectionManual, cluster, manualState)

    const {active: {warning, instance, instances}, colors, fetching, refetch} = isDetectionManual ? manual : auto

    // we need disable cause it thinks that function can be updated, and it causes endless recursion
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectStoreUpdate, [isActive, warning, cluster, instance, instances, detection])
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectRefetch, [isActive, detection, detectionRef])

    useEffect(handleEffectDetection, [isActive, store.activeCluster, detection])
    useEffect(handleEffectManualInstance, [isActive, store.activeCluster, manualState])

    return (
        <TableRow>
            <ListCell>
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
            </ListCell>
            <ListCell>
                <DynamicInputs
                    inputs={stateNodes}
                    colors={colors}
                    editable={editable}
                    placeholder={`Instance`}
                    onChange={n => setStateNodes(n)}
                />
            </ListCell>
            <ListCell>
                {!editable ? (
                    <ListRowRead name={name} toggle={toggle}/>
                ) : (
                    <ListRowUpdate
                        name={name}
                        nodes={stateNodes}
                        toggle={toggle}
                        onUpdate={refetch}
                        onClose={() => setStateNodes(cluster.nodes)}
                    />
                )}
            </ListCell>
        </TableRow>
    )

    function renderChipTooltip() {
        const items = [
            { name: "Detection", value: detection, bgColor: purple[400] },
            { name: "Instance", value: getDomain(instance.sidecar), bgColor: InstanceColor[instance.role] },
            { name: "Warning", value: warning ? "Yes" : "No", bgColor: warning ? orange[500] : grey[500] }
        ]

        return <InfoTitle items={items} />
    }

    function handleEffectRefetch() {
        if (isActive && detection !== detectionRef.current) {
            refetch()
            detectionRef.current = detection
        }
    }

    function handleEffectStoreUpdate() {
        if (isActive) {
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
            const storeManualInstance = {
                instance: store.activeCluster.instance,
                instances: store.activeCluster.instances,
            }
            if (!isInstancesEqual(manualState.instance, storeManualInstance.instance)) {
                setManualState(storeManualInstance)
            }
        }
    }

    function isInstancesEqual(i1: InstanceLocal, i2: InstanceLocal) {
        return i1.sidecar.host === i2.sidecar.host && i1.sidecar.port === i2.sidecar.port
    }
}
