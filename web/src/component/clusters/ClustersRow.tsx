import {Box, Chip, TableRow} from "@mui/material";
import {useEffect, useState} from "react";
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
    const isActive = isClusterActive(name)

    const [stateNodes, setStateNodes] = useState(cluster.nodes);
    const {instance, instances, colors, isFetching, warning, update, refetch} = useSmartClusterQuery(name, stateNodes)

    useEffect(handleEffectStoreUpdate, [isActive, cluster, warning, instance, instances, setCluster])

    return (
        <TableRow>
            <ClustersCell>
                <Box display={"flex"} flexDirection={"row"} alignItems={"center"}>
                    <Chip
                        sx={SX.chipSize}
                        color={isActive ? "primary" : "default"}
                        label={name}
                        onClick={handleChipClick}
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
                        onUpdate={() => update(stateNodes)}
                        onClose={() => setStateNodes(cluster.nodes)}
                    />
                )}
            </ClustersCell>
        </TableRow>
    )

    function handleChipClick() {
        setCluster(!isActive ? {cluster, warning, instance, instances} : undefined)
    }

    function handleEffectStoreUpdate() {
        if (isActive) setCluster({cluster, warning, instance, instances})
    }
}
