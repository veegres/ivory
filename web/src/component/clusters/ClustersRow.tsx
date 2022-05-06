import {Box, Chip, FormControl, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {Cached, Cancel, CheckCircle, Delete, Edit} from "@mui/icons-material";
import {useState} from "react";
import {useMutation, useQuery, useQueryClient} from "react-query";
import {clusterApi, nodeApi} from "../../app/api";
import {ClusterMap, Style} from "../../app/types";
import {ClustersRowButton} from "./ClustersRowButton";
import {DynamicInputs} from "../view/DynamicInputs";
import {useStore} from "../../provider/StoreProvider";
import {activeNode, createColorsMap} from "../../app/utils";

const SX = {
    nodesCellInput: {height: '32px'},
    chipSize: {width: '100%'},
    cell: {verticalAlign: "top"},
}

const style: Style = {
    thirdCell: {whiteSpace: 'nowrap'}
}

type Props = {
    name: string,
    nodes: string[],
    edit?: { isReadOnly?: boolean, toggleEdit?: () => void, closeNewElement?: () => void }
}

export function ClustersRow({name, nodes, edit = {}}: Props) {
    const {isReadOnly = false, toggleEdit = () => {}, closeNewElement = () => {}} = edit
    const isNewElement = !name
    const {setStore, isClusterActive} = useStore()

    const [stateName, setStateName] = useState(name);
    const [stateNodes, setStateNodes] = useState(nodes);
    const isActive = isClusterActive(stateName)

    const queryClient = useQueryClient();
    const { data, isFetching, refetch } = useQuery(
        ['node/cluster', name],
        () => nodeApi.cluster(stateNodes[0]),
        {retry: 0}
    )
    const updateCluster = useMutation(clusterApi.update, {
        onSuccess: (data) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            map[data.name] = data.nodes
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })
    const deleteCluster = useMutation(clusterApi.delete, {
        onSuccess: (_, newName) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            delete map[newName]
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })

    return (
        <TableRow>
            <TableCell sx={SX.cell}>
                {renderClusterNameCell()}
            </TableCell>
            <TableCell sx={SX.cell}>
                <DynamicInputs
                    inputs={stateNodes}
                    colors={createColorsMap(data)}
                    editable={!isReadOnly}
                    placeholder={`Node`}
                    onChange={n => setStateNodes(n)}
                />
            </TableCell>
            <TableCell sx={SX.cell} style={style.thirdCell}>
                {isReadOnly ? renderReadButtons(!stateName) : renderActionButtons(!stateName)}
            </TableCell>
        </TableRow>
    )

    function renderClusterNameCell() {
        return !isNewElement ? (
            <Box display={"flex"} flexDirection={"row"}>
                <Chip
                    sx={SX.chipSize}
                    color={isActive ? "primary" : "default"}
                    label={stateName}
                    onClick={handleChipClick}
                />
                <ClustersRowButton icon={<Cached />} loading={isFetching} onClick={refetch} tooltip={"Reload"} />
            </Box>
        ) : (
            <FormControl fullWidth>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    placeholder="Name"
                    value={stateName}
                    onChange={(event) => setStateName(event.target.value)}
                />
            </FormControl>
        )
    }

    function renderReadButtons(isDisabled: boolean) {
        return (
            <>
                <ClustersRowButton
                    icon={<Edit/>}
                    tooltip={'Edit'}
                    loading={updateCluster.isLoading}
                    disabled={isDisabled}
                    onClick={toggleEdit}
                />
                <ClustersRowButton
                    icon={<Delete />}
                    tooltip={'Delete'}
                    loading={deleteCluster.isLoading}
                    disabled={isDisabled}
                    onClick={handleDelete}
                />
            </>
        )
    }

    function renderActionButtons(isDisabled: boolean) {
        return (
            <>
                <ClustersRowButton
                    icon={<Cancel/>}
                    tooltip={'Cancel'}
                    loading={updateCluster.isLoading}
                    onClick={handleCancel}
                />
                <ClustersRowButton
                    icon={<CheckCircle/>}
                    tooltip={'Save'}
                    loading={updateCluster.isLoading}
                    disabled={isDisabled}
                    onClick={handleUpdate}
                />
            </>
        )
    }

    function handleUpdate() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        updateCluster.mutate({name: stateName, nodes: stateNodes})
    }

    function handleCancel() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        setStateNodes([...nodes])
    }

    function handleDelete() {
        deleteCluster.mutate(stateName)
    }

    function handleChipClick() {
        const cluster = !isActive ? {name: stateName, node: activeNode(data)} : {name: '', node: ''}
        setStore({ activeCluster: {...cluster, tab: 0}, activeNode: '' })
    }
}
