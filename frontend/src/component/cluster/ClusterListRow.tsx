import {Box, Chip, CircularProgress, FormControl, IconButton, OutlinedInput, TableCell, TableRow} from "@mui/material";
import {Add, Check, Delete, Edit, Remove, Visibility} from "@mui/icons-material";
import {cloneElement, ReactElement, useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";

const SX = {
    clusterCellIcon: { fontSize: 10 },
    actionCellIcon: { fontSize: 18 },
    actionCellButton: { height: '32px', width: '32px' },
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
                <ClusterNameCell />
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                <ClusterNodesCell />
            </TableCell>
            <TableCell sx={{ verticalAlign: "top", width: '1%', whiteSpace: 'nowrap' }}>
                <ClusterActionsCell />
            </TableCell>
        </TableRow>
    )

    function ClusterNameCell() {
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

    function ClusterNodesCell() {
        return (
            <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.nodesCellBox.minWidth}, 1fr))`} gap={1}>
                {nodes.map((node, index) => (<ClusterNodesCellElement node={node} index={index} />))}
            </Box>
        )
    }

    function ClusterNodesCellElement(props: { node: string, index: number }) {
        const { node, index } = props

        const isNodeSelected = !!node && store.activeNode === node

        if (isReadOnly) {
            return <Chip
                key={index}
                sx={{ width: '100%' }}
                label={node ? node : `Node ${index}`}
                disabled={!node}
                variant="outlined"
                color={isNodeSelected ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeSelected)}
            />
        }

        return (
            <FormControl key={index}>
                <OutlinedInput
                    sx={SX.nodesCellInput}
                    type="string"
                    endAdornment={<NodeSelectButton node={node} isNodeSelected={isNodeSelected} />}
                    placeholder={`Node ${index}`}
                    size="small"
                    value={node}
                    onChange={(event) => handleChange(index, event.target.value)}
                />
            </FormControl>
        )
    }

    function ClusterActionsCell() {
        return (
            <>
                {isReadOnly ?
                    <CellButton icon={<Edit />} isLoading={updateCluster.isLoading} onClick={toggleEdit} /> :
                    <CellButton icon={<Check />} isLoading={updateCluster.isLoading} onClick={handleUpdate} />
                }
                <CellButton icon={<Delete />} isLoading={deleteCluster.isLoading} onClick={handleDelete} />
            </>
        )
    }

    function CellButton(props: { isLoading: boolean, onClick: () => void, icon: ReactElement }) {
        const { isLoading, icon, onClick } = props;
        return (
            <IconButton sx={SX.actionCellButton} disabled={isLoading || !name} onClick={onClick}>
                {isLoading ?
                    <CircularProgress size={SX.actionCellIcon.fontSize} /> :
                    cloneElement(icon, { sx: SX.actionCellIcon })
                }
            </IconButton>
        )
    }

    function NodeSelectButton(props: { node: string, isNodeSelected: boolean }) {
        const { node, isNodeSelected } = props
        if (!node) return null

        return (
            <IconButton
                sx={SX.nodesCellButton}
                color={isNodeSelected ? "primary" : "default"}
                onClick={handleSelectedNodeChange(node, isNodeSelected)}>
                <Visibility sx={SX.nodesCellIcon} />
            </IconButton>
        )
    }

    function handleSelectedNodeChange(selectedNode: string, isNodeSelected: boolean) {
        return () => setStore({ activeNode: isNodeSelected ? '' : selectedNode })
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
