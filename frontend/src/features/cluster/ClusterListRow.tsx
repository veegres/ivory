import {Box, FormControl, IconButton, Input, TableCell, TableRow} from "@mui/material";
import {Add, Delete, Remove, Save, Visibility} from "@mui/icons-material";
import {Dispatch, useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {clusterApi} from "../../app/api";

const SX = {
    clusterNameIcon: { fontSize: 10 },
    tableIcon: { fontSize: 14 },
    input: { minWidth: '150px' }
}

export function ClusterListRow(props: { nodes?: string[], name?: string, setNode?: Dispatch<string> }) {
    const [name, setName] = useState(props.name ?? '');
    const [nodes, setNodes] = useState(props.nodes && props.nodes.length ? props.nodes : ['']);

    const queryClient = useQueryClient();
    const { mutate } = useMutation(clusterApi.create, {
        onSuccess: () => queryClient.invalidateQueries('cluster/list')
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
                <IconButton onClick={() => {}}><Delete sx={SX.tableIcon} /></IconButton>
                <IconButton onClick={handleUpdate}><Save sx={SX.tableIcon} /></IconButton>
            </TableCell>
        </TableRow>
    )

    function ViewEndAdornment({ node }: { node: string }) {
        const update = props.setNode
        if (!update) return null

        return (
            <IconButton onClick={() => update(node)}>
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
        mutate({ name, nodes })
    }
}
