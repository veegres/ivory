import {Box, Chip, FormControl, IconButton, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {Cancel, CheckCircle, Delete, Edit, Visibility} from "@mui/icons-material";
import {useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";
import {ClustersActionButton} from "./ClustersActionButton";

const SX = {
    nodesCellIcon: { fontSize: 18 },
    nodesCellButton: { height: '22px', width: '22px' },
    nodesCellBox: { minWidth: '150px' },
    nodesCellInput: { height: '32px' }
}

type Props = {
    name?: string,
    nodes?: string[],
    edit?: { isReadOnly: boolean, toggleEdit: () => void }
}

export function ClustersRow({ name = '', nodes = [''], edit = { isReadOnly: false, toggleEdit: () => {} } }: Props) {
    const { store, setStore } = useStore()
    const { isReadOnly, toggleEdit } = edit
    const [stateName, setStateName] = useState(name);
    const [stateNodes, setStateNodes] = useState(nodes.length ? nodes : ['']);

    if (!isReadOnly) handleAutoIncreaseElement()

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
            <TableCell sx={{ verticalAlign: "top", width: '15%' }}>
                {renderClusterNameCell()}
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                {renderClusterNodesCell()}
            </TableCell>
            <TableCell sx={{ verticalAlign: "top", width: '1%', whiteSpace: 'nowrap' }}>
                {renderClusterActionsCell()}
            </TableCell>
        </TableRow>
    )

    function renderClusterNameCell() {
        if (name) {
            return <Chip sx={{ width: '100%' }} label={stateName} />
        }

        return (
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

    function renderClusterNodesCell() {
        return (
            <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.nodesCellBox.minWidth}, 1fr))`} gap={1}>
                {stateNodes.map((node, index) => renderClusterNodesCellElement(node, index))}
            </Box>
        )
    }

    function renderClusterNodesCellElement(node: string, index: number) {
        const isNodeActive = !!node && store.activeNode === node

        if (isReadOnly) {
            return (
                <Chip
                    key={index}
                    sx={{ width: '100%' }}
                    label={node ? node : `Node ${index}`}
                    disabled={!node}
                    variant="outlined"
                    color={isNodeActive ? "primary" : "default"}
                    onClick={handleSelectedNodeChange(node, isNodeActive)}
                />
            )
        }

        return (
            <FormControl key={index}>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    type="string"
                    endAdornment={renderNodeSelectIcon(node, isNodeActive)}
                    placeholder={`Node ${index}`}
                    size="small"
                    value={node}
                    onChange={(event) => handleChange(index, event.target.value)}
                />
            </FormControl>
        )
    }

    function renderClusterActionsCell() {
        const isDisabled = !stateName
        return isReadOnly ? (
            <>
                <ClustersActionButton icon={<Edit />} tooltip={'Edit'} loading={updateCluster.isLoading} disabled={isDisabled} onClick={toggleEdit} />
                <ClustersActionButton icon={<Delete />} tooltip={'Delete'} loading={deleteCluster.isLoading} disabled={isDisabled} onClick={handleDelete} />
            </>
        ) : (
            <>
                <ClustersActionButton icon={<Cancel />} tooltip={'Cancel'} loading={updateCluster.isLoading} disabled={!name} onClick={toggleEdit} />
                <ClustersActionButton icon={<CheckCircle />} tooltip={'Save'} loading={updateCluster.isLoading} disabled={isDisabled} onClick={handleUpdate} />
            </>
        )
    }

    function renderNodeSelectIcon(node: string, isNodeActive: boolean) {
        if (!node) return null

        return (
            <IconButton
                sx={SX.nodesCellButton}
                color={isNodeActive ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeActive)}>
                <Visibility sx={SX.nodesCellIcon} />
            </IconButton>
        )
    }

    function handleSelectedNodeChange(selectedNode: string, isNodeActive: boolean) {
        return () => setStore({ activeNode: isNodeActive ? '' : selectedNode })
    }

    function handleChange(index: number, value: string) {
        setStateNodes((previous) => { previous[index] = value; return [...previous] })
    }

    function handleAutoIncreaseElement() {
        if (stateNodes.length > 0) {
            const lastIndex = stateNodes.length - 1
            const preLastIndex = stateNodes.length - 2
            if (stateNodes[lastIndex] !== '') {
                setStateNodes([...stateNodes, ''])
            }

            if (stateNodes.length > 1) {
                if (stateNodes[lastIndex] === '' && stateNodes[preLastIndex] === '') {
                    setStateNodes(stateNodes.slice(0, stateNodes.length - 1))
                }
            }
        }
    }

    function handleUpdate() {
        const newNodes = stateNodes[stateNodes.length - 1] === '' ? stateNodes.slice(0, stateNodes.length - 1) : stateNodes
        toggleEdit()
        updateCluster.mutate({ name: stateName, nodes: newNodes })
    }

    function handleDelete() {
        deleteCluster.mutate(stateName)
    }
}
