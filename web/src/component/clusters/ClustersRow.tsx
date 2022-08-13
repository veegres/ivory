import {Box, Chip, TableRow} from "@mui/material";
import {useState} from "react";
import {Cluster} from "../../app/types";
import {RefreshIconButton,} from "../view/IconButtons";
import {DynamicInputs} from "../view/DynamicInputs";
import {initialActiveCluster, useStore} from "../../provider/StoreProvider";
import {createColorsMap} from "../../app/utils";
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
    const {store, setStore, isClusterActive} = useStore()

    const [stateNodes, setStateNodes] = useState(cluster.nodes);
    const isActive = isClusterActive(name)

    const {instance, instanceResult, update, refetch} = useSmartClusterQuery(name, stateNodes)
    const {data, isFetching} = instanceResult
    const instanceMap = data ?? {}

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
                    colors={createColorsMap(instanceMap)}
                    editable={editable}
                    placeholder={`Node`}
                    onChange={n => setStateNodes(n)}
               />
            </ClustersCell>
            <ClustersCell>
                {!editable ? (
                    <ClustersRowRead name={name} edit={toggle}/>
                ) : (
                    <ClustersRowUpdate
                        name={name}
                        nodes={stateNodes}
                        toggle={toggle}
                        onUpdate={() => update(stateNodes)}
                    />
                )}
            </ClustersCell>
        </TableRow>
    )

    function handleChipClick() {
        // either find leader or set instance that we were sent request to
        const activeInstance = Object.values(instanceMap).find(node => node.leader) ?? instanceMap[instance]
        const activeCluster = isActive ? initialActiveCluster : {cluster: cluster, instance: activeInstance}

        setStore({
            activeCluster: {...activeCluster, tab: store.activeCluster.tab},
            activeNode: ''
        })
    }
}
