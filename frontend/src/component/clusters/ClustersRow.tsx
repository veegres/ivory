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

export function ClustersRow(globalProps: { name: string, nodes: string[], edit: { isReadOnly: boolean, toggleEdit: () => void } }) {
    const { store, setStore } = useStore()
    const { isReadOnly, toggleEdit } = globalProps.edit
    const [name, setName] = useState(globalProps.name);
    const [nodes, setNodes] = useState(globalProps.nodes.length ? globalProps.nodes : ['']);

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
        if (globalProps.name) {
            return <Chip sx={{ width: '100%' }} label={name} />
        }

        return (
            <FormControl fullWidth>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    placeholder="Name"
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                />
            </FormControl>
        )
    }

    function renderClusterNodesCell() {
        return (
            <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.nodesCellBox.minWidth}, 1fr))`} gap={1}>
                {nodes.map((node, index) => renderClusterNodesCellElement(node, index))}
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
        const isDisabled = !name
        return isReadOnly ? (
            <>
                <ClustersActionButton icon={<Edit />} tooltip={'Edit'} loading={updateCluster.isLoading} disabled={isDisabled} onClick={toggleEdit} />
                <ClustersActionButton icon={<Delete />} tooltip={'Delete'} loading={deleteCluster.isLoading} disabled={isDisabled} onClick={handleDelete} />
            </>
        ) : (
            <>
                <ClustersActionButton icon={<Cancel />} tooltip={'Cancel'} loading={updateCluster.isLoading} disabled={!globalProps.name} onClick={toggleEdit} />
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
        setNodes((previous) => { previous[index] = value; return [...previous] })
    }

    function handleAutoIncreaseElement() {
        if (nodes.length > 0) {
            const lastIndex = nodes.length - 1
            const preLastIndex = nodes.length - 2
            if (nodes[lastIndex] !== '') {
                setNodes([...nodes, ''])
            }

            if (nodes.length > 1) {
                if (nodes[lastIndex] === '' && nodes[preLastIndex] === '') {
                    setNodes(nodes.slice(0, nodes.length - 1))
                }
            }
        }
    }

    function handleUpdate() {
        const newNodes = nodes[nodes.length - 1] === '' ? nodes.slice(0, nodes.length - 1) : nodes
        toggleEdit()
        updateCluster.mutate({ name, nodes: newNodes })
    }

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
