import {Box, CircularProgress, FormControl, IconButton, Input, TableCell, TableRow} from "@mui/material";
import {Add, CheckCircle, Delete, Remove, Visibility} from "@mui/icons-material";
import {ReactNode, useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";
import {ClusterMap} from "../../app/types";
import {useStore} from "../../provider/StoreProvider";

const SX = {
    clusterNameIcon: { fontSize: 10 },
    tableIcon: { fontSize: 14 },
    input: { minWidth: '150px' }
}

export function ClusterListRow(props: { nodes?: string[], name?: string }) {
    const { setStore } = useStore()
    const [name, setName] = useState(props.name ?? '');
    const [nodes, setNodes] = useState(props.nodes && props.nodes.length ? props.nodes : ['']);

    const queryClient = useQueryClient();
    const updateCluster = useMutation(clusterApi.update, {
        onSuccess: (data) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            map[data.name] = data.nodes
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })
    const deleteCluster = useMutation(clusterApi.delete, {
        onSuccess: (_, name) => {
            const map = queryClient.getQueryData<ClusterMap>('cluster/list') ?? {} as ClusterMap
            delete map[name]
            queryClient.setQueryData<ClusterMap>('cluster/list', map)
        }
    })

    return (
        <TableRow>
            <TableCell sx={{ verticalAlign: "top", width: '15%' }}>
                <Box display="flex" sx={SX.input} flexDirection="row" alignItems="center" justifyContent="space-between">
                    <FormControl fullWidth focused>
                        <Input size="small" placeholder="Name" value={name} onChange={(event) => setName(event.target.value)}/>
                    </FormControl>
                    <Box display="flex" flexDirection="column" alignItems="center" justifyContent="space-between">
                        <IconButton sx={{ padding: '2px' }} onClick={handleAdd}><Add sx={SX.clusterNameIcon} /></IconButton>
                        <IconButton sx={{ padding: '2px' }} onClick={handleRemove}><Remove sx={SX.clusterNameIcon} /></IconButton>
                    </Box>
                </Box>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top" }}>
                <Box display="grid" gridTemplateColumns={`repeat(auto-fill, minmax(${SX.input.minWidth}, 1fr))`} gap={1}>
                    {nodes.map((node, index) => (
                        <FormControl key={index} focused>
                            <Input
                                color="secondary"
                                type="string"
                                endAdornment={<ViewEndAdornment node={node} />}
                                placeholder={`Node ${index}`}
                                size="small"
                                value={node}
                                onChange={(event) => handleChange(index, event.target.value)}
                            />
                        </FormControl>
                    ))}
                </Box>
            </TableCell>
            <TableCell sx={{ verticalAlign: "top", width: '1%', whiteSpace: 'nowrap' }}>
                <CellButton isLoading={updateCluster.isLoading} onClick={handleUpdate}>
                    <CheckCircle sx={SX.tableIcon} />
                </CellButton>
                <CellButton isLoading={deleteCluster.isLoading} onClick={handleDelete}>
                    <Delete sx={SX.tableIcon} />
                </CellButton>
            </TableCell>
        </TableRow>
    )

    function CellButton(props: { isLoading: boolean, onClick: () => void, children: ReactNode }) {
        const { isLoading, children, onClick } = props;
        return (
            <IconButton disabled={isLoading} onClick={onClick}>
                {isLoading ? <CircularProgress size={14} /> : children}
            </IconButton>
        )
    }

    function ViewEndAdornment({ node }: { node: string }) {
        if (!node) return null

        return (
            <IconButton onClick={() => setStore({ activeNode: node })}>
                <Visibility sx={SX.tableIcon} />
            </IconButton>
        )
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
        updateCluster.mutate({ name, nodes })
    }

    function handleDelete() {
        deleteCluster.mutate(name)
    }
}
