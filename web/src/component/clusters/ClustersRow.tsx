import {Box, Chip, FormControl, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap, Style} from "../../app/types";
import {
    CancelIconButton,
    DeleteIconButton,
    EditIconButton,
    RefreshIconButton,
    SaveIconButton,
} from "../view/IconButtons";
import {DynamicInputs} from "../view/DynamicInputs";
import {useStore} from "../../provider/StoreProvider";
import {activeNode, createColorsMap} from "../../app/utils";
import {useSmartClusterQuery} from "../../app/hooks";

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
    const {store, setStore, isClusterActive} = useStore()

    const [stateName, setStateName] = useState(name);
    const [stateNodes, setStateNodes] = useState(nodes);
    const isActive = isClusterActive(stateName)

    const {instance, instanceResult: {data, isFetching}, update, refetch} = useSmartClusterQuery(stateName, stateNodes)

    const queryClient = useQueryClient();
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
                <RefreshIconButton loading={isFetching} onClick={refetch}/>
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
                <EditIconButton loading={updateCluster.isLoading} disabled={isDisabled} onClick={toggleEdit}/>
                <DeleteIconButton loading={deleteCluster.isLoading} disabled={isDisabled} onClick={handleDelete}/>
            </>
        )
    }

    function renderActionButtons(isDisabled: boolean) {
        return (
            <>
                <CancelIconButton loading={updateCluster.isLoading} onClick={handleCancel}/>
                <SaveIconButton loading={updateCluster.isLoading} disabled={isDisabled} onClick={handleUpdate}/>
            </>
        )
    }

    function handleUpdate() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        updateCluster.mutate({name: stateName, nodes: stateNodes})
        update(stateNodes)
    }

    function handleCancel() {
        if (isNewElement) closeNewElement(); else toggleEdit()
        setStateNodes([...nodes])
    }

    function handleDelete() {
        deleteCluster.mutate(stateName)
    }

    function handleChipClick() {
        const cluster = !isActive ? {
            name: stateName, instance: instance, leader: activeNode(data)
        } : {
            name: "", instance: "", leader: undefined
        }
        setStore({ activeCluster: {...cluster, tab: store.activeCluster.tab}, activeNode: '' })
    }
}
