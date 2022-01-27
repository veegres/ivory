import {Box, Chip, CircularProgress, FormControl, IconButton, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {Add, Check, Delete, Edit, Remove, Visibility} from "@mui/icons-material";
import {cloneElement, ReactElement, useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";
import {ClusterListActionButton} from "./ClusterListActionButton";

const SX = {
    clusterCellIcon: { fontSize: 10 },
    nodesCellIcon: { fontSize: 18 },
    nodesCellButton: { height: '22px', width: '22px' },
    nodesCellBox: { minWidth: '150px' },
    nodesCellInput: { height: '32px' }
}

export function ClusterListRow(globalProps: { name: string, nodes: string[], edit: { isReadOnly: boolean, toggleEdit: () => void } }) {
    const { store, setStore } = useStore()
    const { isReadOnly, toggleEdit } = globalProps.edit
    const [name, setName] = useState(globalProps.name);
    const [nodes, setNodes] = useState(globalProps.nodes.length ? globalProps.nodes : ['']);

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
            <Box display="flex" sx={SX.nodesCellBox} flexDirection="row" alignItems="center" justifyContent="space-between">
                <FormControl fullWidth>
                    <OutlinedInput
                        sx={SX.nodesCellInput}
                        placeholder="Name"
                        value={name}
                        onChange={(event) => setName(event.target.value)}
                    />
                </FormControl>
                <Box display="flex" flexDirection="column" alignItems="center" justifyContent="space-between">
                    <IconButton sx={{ padding: '2px' }} onClick={handleAdd}><Add sx={SX.clusterCellIcon} /></IconButton>
                    <IconButton sx={{ padding: '2px' }} onClick={handleRemove}><Remove sx={SX.clusterCellIcon} /></IconButton>
                </Box>
            </Box>
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
            return <Chip
                key={index}
                sx={{ width: '100%' }}
                label={node ? node : `Node ${index}`}
                disabled={!node}
                variant="outlined"
                color={isNodeActive ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeActive)}
            />
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
        return (
            <>
                {isReadOnly ?
                    <ClusterListActionButton icon={<Edit />} loading={updateCluster.isLoading} disabled={isDisabled} onClick={toggleEdit} /> :
                    <ClusterListActionButton icon={<Check />} loading={updateCluster.isLoading} disabled={isDisabled} onClick={handleUpdate} />
                }
                <ClusterListActionButton icon={<Delete />} loading={deleteCluster.isLoading} disabled={isDisabled} onClick={handleDelete} />
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

    function handleAdd() {
        nodes.push("")
        setNodes([...nodes])
    }

    function handleRemove() {
        if (nodes.length <= 1) return
        nodes.pop()
        setNodes([...nodes])
    }

    function handleUpdate() {
        toggleEdit()
        updateCluster.mutate({ name, nodes })
    }

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
