import {Box, Chip, TableRow} from "@mui/material";
import {useEffect, useMemo, useState} from "react";
import {Cluster} from "../../app/types";
import {RefreshIconButton,} from "../view/IconButtons";
import {DynamicInputs} from "../view/DynamicInputs";
import {useStore} from "../../provider/StoreProvider";
import {useSmartClusterQuery} from "../../app/hooks";
import {ClustersRowRead} from "./ClustersRowRead";
import {ClustersRowUpdate} from "./ClustersRowUpdate";
import {ClustersCell} from "./ClustersCell";

const SX = {
    chipSize: {width: '100%'},
}

type Props = {
    name: string,
    cluster: Cluster,
    editable: boolean,
    toggle: () => void,
}

export function ClustersRow({name, cluster, editable, toggle}: Props) {
    const {setCluster, isClusterActive} = useStore()
    const active = isClusterActive(name)

    const [stateNodes, setStateNodes] = useState(cluster.nodes);
    const {instance, instances, colors, isFetching, warning, refetch} = useSmartClusterQuery(name, cluster.nodes)
    const activeCluster = useMemo(() => ({cluster, warning, instance, instances}), [cluster, warning, instance, instances])

    // we need disable cause it thinks that function can be updated, and it causes endless recursion
    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(handleEffectStoreUpdate, [active, activeCluster])

    return (
        <TableRow>
            <ClustersCell>
                <Box display={"flex"} flexDirection={"row"} alignItems={"center"}>
                    <Chip
                        sx={SX.chipSize}
                        color={active ? "primary" : "default"}
                        label={name}
                        onClick={() => setCluster(active ? undefined : activeCluster)}
                    />
                    <RefreshIconButton loading={isFetching} onClick={refetch}/>
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

    function handleEffectStoreUpdate() {
        if (active) setCluster(activeCluster)
    }
}
